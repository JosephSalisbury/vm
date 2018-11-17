package google

import (
	"fmt"
	"io/ioutil"

	"github.com/JosephSalisbury/vm/ignition"
	"github.com/JosephSalisbury/vm/provider"
	compute "google.golang.org/api/compute/v1"
)

func (p *googleProvider) Create(channel string, ignition ignition.Interface, cpu int, ram int) error {
	id := provider.ID()

	machineType := getMachineType(cpu, ram)
	p.logger.Printf("using machine type: %s", machineType)

	if err := ignition.Create(); err != nil {
		return err
	}

	ignitionPath, err := ignition.Path()
	if err != nil {
		return err
	}
	ignitionData, err := ioutil.ReadFile(ignitionPath)
	if err != nil {
		return err
	}
	ignitionDataString := string(ignitionData)

	p.logger.Printf("creating instance")

	if _, err := p.instancesService.Insert(
		p.project,
		p.zone,
		&compute.Instance{
			Disks: []*compute.AttachedDisk{
				{
					AutoDelete: true,
					Boot:       true,
					InitializeParams: &compute.AttachedDiskInitializeParams{
						SourceImage: fmt.Sprintf(
							"projects/coreos-cloud/global/images/family/coreos-%v",
							channel,
						),
					},
				},
			},
			Name:        id,
			MachineType: fmt.Sprintf("zones/%v/machineTypes/%v", p.zone, machineType),
			Metadata: &compute.Metadata{
				Items: []*compute.MetadataItems{
					{
						Key:   "user-data",
						Value: &ignitionDataString,
					},
				},
			},
			NetworkInterfaces: []*compute.NetworkInterface{
				{
					AccessConfigs: []*compute.AccessConfig{
						{
							Name: "External NAT",
							Type: "ONE_TO_ONE_NAT",
						},
					},
				},
			},
			Scheduling: &compute.Scheduling{
				Preemptible: p.preemptible,
			},
		},
	).Do(); err != nil {
		return err
	}

	p.logger.Printf("instance created")

	return nil
}
