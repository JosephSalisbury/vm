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

	// TODO: Use exec.LookPath.
	if err := exec.Command("VBoxManage", "--help").Run(); err != nil {
		return nil, VBoxManageMissingError
	}

	// TODO: Check `vbox-configdrive-gen` is available.

	p := &virtualBoxProvider{
		logger: config.Logger,
	}

	return p, nil
}
