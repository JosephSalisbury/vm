package virtualbox

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/JosephSalisbury/vm/coreos"
	"github.com/JosephSalisbury/vm/ignition"
	"github.com/JosephSalisbury/vm/provider"
)

func (p *virtualBoxProvider) Create(channel string, ignition *ignition.Ignition, cpu int, ram int) error {
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

	if err := ignition.Create(); err != nil {
		return err
	}

	p.logger.Printf("building config drive")
	configDir, err := ioutil.TempDir("", "vmConfig")
	if err != nil {
		return err
	}
	configDriveImagePath := path.Join(configDir, "configdrive.img")
	configDriveVMDKPath := path.Join(configDir, "configdrive.vmdk")

	p.logger.Printf("using ignition config at %s", ignition.Path())
	p.logger.Printf("configdrive image is at %s", configDriveImagePath)
	p.logger.Printf("configdrive vmdk is at %s", configDriveVMDKPath)

	p.logger.Printf("creating configdrive")
	vboxConfigdriveGenOut, err := exec.Command("bash", "-c", fmt.Sprintf("cat %s | vbox-configdrive-gen > %s", ignition.Path(), configDriveImagePath)).Output()
	if err != nil {
		p.logger.Printf("%s", string(vboxConfigdriveGenOut))
		return err
	}

	if _, err := p.vboxmanage(vboxManageCommand{
		description: "create configdrive vmdk",
		args: []string{
			"internalcommands", "createrawvmdk",
			"-filename", configDriveVMDKPath,
			"-rawdisk", configDriveImagePath,
		},
	}); err != nil {
		return err
	}

	id := provider.ID()

	if _, err := p.vboxmanage(vboxManageCommand{
		description: "create virtual machine",
		args: []string{
			"createvm",
			"--name", id,
			"--ostype", "Linux26_64",
			"--register",
		},
	}); err != nil {
		return err
	}

	if _, err := p.vboxmanage(vboxManageCommand{
		description: "configure virtual machine resources",
		args: []string{
			"modifyvm", id,
			"--cpus", strconv.Itoa(cpu),
			"--memory", strconv.Itoa(ram * 1024),
			"--audio", "none",
		},
	}); err != nil {
		return err
	}

	if _, err := p.vboxmanage(vboxManageCommand{
		description: "configure network connection",
		args: []string{
			"modifyvm", id,
			"--nic1", "bridged",
			"--bridgeadapter1", "en0",
		},
	}); err != nil {
		return err
	}

	if _, err := p.vboxmanage(vboxManageCommand{
		description: "attach SATA controller",
		args: []string{
			"storagectl", id,
			"--name", "SATA",
			"--add", "sata",
			"--controller", "IntelAHCI",
			"--portcount", "2",
		},
	}); err != nil {
		return err
	}
	if _, err := p.vboxmanage(vboxManageCommand{
		description: "attach CoreOS disk",
		args: []string{
			"storageattach", id,
			"--storagectl", "SATA",
			"--port", "0",
			"--device", "0",
			"--type", "hdd",
			"--medium", uncompressedPath,
			"--mtype", "immutable",
		},
	}); err != nil {
		return err
	}
	if _, err := p.vboxmanage(vboxManageCommand{
		description: "attach config drive",
		args: []string{
			"storageattach", id,
			"--storagectl", "SATA",
			"--port", "1",
			"--device", "0",
			"--type", "hdd",
			"--medium", configDriveVMDKPath,
			"--mtype", "immutable",
		},
	}); err != nil {
		return err
	}

	if _, err := p.vboxmanage(vboxManageCommand{
		description: "start vm",
		args: []string{
			"startvm", id,
			"--type", "headless",
		},
	}); err != nil {
		return err
	}

	return nil
}
