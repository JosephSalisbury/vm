package virtualbox

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/JosephSalisbury/vm/coreos"
	"github.com/JosephSalisbury/vm/provider"
)

func (p *virtualBoxProvider) Create(channel string, ignitionPath string, cpu int, ram int) error {
	// TODO: Validate channel, ignitionPath, cpu, ram.

	p.logger.Printf("fetching latest CoreOS Container Linux version for channel %v", channel)
	version, err := coreos.LatestVersion(channel)
	if err != nil {
		return err
	}
	p.logger.Printf("latest version for channel %v is %v", channel, version)

	compressedPath := coreos.CompressedPath(channel, version)
	p.logger.Printf("checking if image is already downloaded to %v", compressedPath)

	if _, err := os.Stat(compressedPath); err == nil {
		p.logger.Printf("image is already downloaded")
	} else {
		p.logger.Printf("downloading image")

		if err := coreos.DownloadImage(channel, compressedPath); err != nil {
			return err
		}

		// TODO: Verify image.
	}

	uncompressedPath := coreos.UncompressedPath(channel, version)
	p.logger.Printf("checking if image is already decompressed to %v", uncompressedPath)

	if _, err := os.Stat(uncompressedPath); err == nil {
		p.logger.Printf("image is already decompressed")
	} else {
		p.logger.Printf("decompressing image")

		if err := coreos.DecompressImage(compressedPath); err != nil {
			return err
		}
	}

	p.logger.Printf("building config drive")
	configDir, err := ioutil.TempDir("", "vmConfig")
	if err != nil {
		return err
	}
	configDriveImagePath := path.Join(configDir, "configdrive.img")
	configDriveVMDKPath := path.Join(configDir, "configdrive.vmdk")

	p.logger.Printf("using ignition config at %s", ignitionPath)
	p.logger.Printf("configdrive image is at %s", configDriveImagePath)
	p.logger.Printf("configdrive vmdk is at %s", configDriveVMDKPath)

	p.logger.Printf("creating configdrive")
	out, err := exec.Command("bash", "-c", fmt.Sprintf("cat %s | vbox-configdrive-gen > %s", ignitionPath, configDriveImagePath)).Output()
	if err != nil {
		p.logger.Printf("%s", string(out))
		return err
	}

	// TODO: Abstract vboxmanage commands (better error handling)
	p.logger.Printf("creating configdrive vmdk")
	out2, err := exec.Command("VBoxManage", "internalcommands", "createrawvmdk", "-filename", configDriveVMDKPath, "-rawdisk", configDriveImagePath).Output()
	if err != nil {
		p.logger.Printf("%s", string(out2))
		return err
	}

	p.logger.Printf("creating vm")
	id := provider.ID()

	if err := exec.Command("VBoxManage", "createvm", "--name", id, "--ostype", "Linux26_64", "--register").Run(); err != nil {
		return err
	}

	if err := exec.Command("VBoxManage", "modifyvm", id, "--cpus", strconv.Itoa(cpu)).Run(); err != nil {
		return err
	}
	if err := exec.Command("VBoxManage", "modifyvm", id, "--memory", strconv.Itoa(ram*1024)).Run(); err != nil {
		return err
	}
	if err := exec.Command("VBoxManage", "modifyvm", id, "--audio", "none").Run(); err != nil {
		return err
	}

	if err := exec.Command("VBoxManage", "storagectl", id, "--name", "SATA", "--add", "sata", "--controller", "IntelAHCI", "--portcount", "2").Run(); err != nil {
		return err
	}
	if err := exec.Command("VBoxManage", "storageattach", id, "--storagectl", "SATA", "--port", "0", "--device", "0", "--type", "hdd", "--medium", uncompressedPath, "--mtype", "immutable").Run(); err != nil {
		return err
	}
	if err := exec.Command("VBoxManage", "storageattach", id, "--storagectl", "SATA", "--port", "1", "--device", "0", "--type", "hdd", "--medium", configDriveVMDKPath, "--mtype", "immutable").Run(); err != nil {
		return err
	}

	// TODO: Mount secret volume.

	freePort := getFreePort()
	if err := exec.Command("VBoxManage", "modifyvm", id, "--natpf1", fmt.Sprintf("ssh,tcp,,%v,,22", freePort)).Run(); err != nil {
		return err
	}

	p.logger.Printf("starting vm")
	if err := exec.Command("VBoxManage", "startvm", id, "--type", "headless").Run(); err != nil {
		return err
	}

	return nil
}
