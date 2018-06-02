// Package virtualbox provides a Provider implementation backed by VirtualBox.
package virtualbox

import (
	"log"
	"os/exec"

	"github.com/JosephSalisbury/vm/provider"
)

var (
	// Provider is a Provider that uses VirtualBox.
	Provider = provider.Name("VirtualBox")

	vdiName = "coreos_production_qemu_image.vdi"
)

type virtualBoxProvider struct {
	logger *log.Logger
}

func New(config provider.Config) (provider.Interface, error) {
	if config.Logger == nil {
		return nil, provider.InvalidConfigError
	}

	if _, err := exec.LookPath("VBoxManage"); err != nil {
		return nil, VBoxManageMissingError
	}
	if _, err := exec.LookPath("vbox-configdrive-gen"); err != nil {
		return nil, VBoxConfigdriveGenMissingError
	}

	p := &virtualBoxProvider{
		logger: config.Logger,
	}

	return p, nil
}

type vboxManageCommand struct {
	description string
	args        []string
}

func (p *virtualBoxProvider) vboxmanage(command vboxManageCommand) (string, error) {
	p.logger.Printf("executing: %s", command.description)

	out, err := exec.Command("VBoxManage", command.args...).Output()
	stringOut := string(out)

	if err != nil {
		p.logger.Printf("VBoxManage output: %s", stringOut)
	}

	return stringOut, err
}
