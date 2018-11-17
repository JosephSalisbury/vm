package google

import (
	"github.com/JosephSalisbury/vm/provider"
)

func (p *googleProvider) List() ([]provider.Status, error) {
	p.logger.Printf("listing instances")

	instanceList, err := p.instancesService.List(p.project, p.zone).Do()
	if err != nil {
		return nil, err
	}

	statuses := []provider.Status{}
	for _, instance := range instanceList.Items {
		status := provider.Status{
			ID:   instance.Name,
			IP:   instance.NetworkInterfaces[0].AccessConfigs[0].NatIP,
			Port: 22,
		}

		statuses = append(statuses, status)
	}

	p.logger.Printf("instances listed")

	return statuses, nil
}
