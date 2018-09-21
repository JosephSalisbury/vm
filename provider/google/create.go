package google

import (
	"fmt"

	"github.com/JosephSalisbury/vm/ignition"
	"github.com/JosephSalisbury/vm/provider"
)

func (p *googleProvider) Create(channel string, ignition *ignition.Ignition, cpu int, ram int) error {
	id := provider.ID()

	machineType := getMachineType(cpu, ram)
	p.logger.Printf("using machine type: %s", machineType)

	if err := ignition.Create(); err != nil {
		return err
	}

	if _, err := p.gcloud(gCloudCommand{
		description: "create instance",
		args: []string{
			"compute",
			"instances",
			"create",
			id,
			"--image-project", "coreos-cloud",
			"--image-family", fmt.Sprintf("coreos-%s", channel),
			"--zone", zone,
			"--machine-type", machineType,
			"--metadata-from-file", fmt.Sprintf("user-data=%s", ignition.Path()),
		},
	}); err != nil {
		return err
	}

	return nil
}
