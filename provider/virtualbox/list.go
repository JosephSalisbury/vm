package virtualbox

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/JosephSalisbury/vm/provider"
)

func (p *virtualBoxProvider) List() ([]provider.Status, error) {
	listVMsOut, err := p.vboxmanage(vboxManageCommand{
		description: "list VMs",
		args: []string{
			"list", "vms",
		},
	})
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, line := range strings.Split(string(listVMsOut), "\n") {
		if line != "" {
			name := strings.Split(line, "\"")[1]
			names = append(names, name)
		}
	}

	statuses := []provider.Status{}
	for _, name := range names {
		// TODO: Clean this shit up.
		var hostPort int

		showVMInfoOut, err := p.vboxmanage(vboxManageCommand{
			description: "show VM info",
			args: []string{
				"showvminfo", name,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, line := range strings.Split(string(showVMInfoOut), "\n") {
			if strings.Contains(line, "host port") {
				hostPort, err = strconv.Atoi(
					strings.TrimSuffix(
						strings.Split(line, " ")[18],
						",",
					),
				)
				if err != nil {
					return nil, err
				}
				break
			}
		}

		status := provider.Status{
			ID:   name,
			IP:   fmt.Sprintf("127.0.0.1"),
			Port: hostPort,
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}
