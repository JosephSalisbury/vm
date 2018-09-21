package google

import (
	"strings"

	"github.com/JosephSalisbury/vm/provider"
)

const (
	idIndex = 0
	ipIndex = 4

	port = 22
)

func (p *googleProvider) List() ([]provider.Status, error) {
	listInstancesOut, err := p.gcloud(gCloudCommand{
		description: "list instances",
		args: []string{
			"compute", "instances", "list",
		},
	})
	if err != nil {
		return nil, err
	}

	lines := strings.Split(listInstancesOut, "\n")[1:]

	statuses := []provider.Status{}
	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)

		id := fields[idIndex]
		ip := fields[ipIndex]

		status := provider.Status{
			ID:   id,
			IP:   ip,
			Port: port,
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}
