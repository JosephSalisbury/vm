package virtualbox

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/JosephSalisbury/vm/provider"
)

const (
	nameIndex = 1
	portIndex = 18

	comma         = ","
	newLine       = "\n"
	portLine      = "host port"
	quotationMark = "\""
	space         = " "

	localhost = "127.0.0.1"
)

func (p *virtualBoxProvider) List() ([]provider.Status, error) {
	// List the VMs.
	listVMsOut, err := p.vboxmanage(vboxManageCommand{
		description: "list VMs",
		args: []string{
			"list", "vms",
		},
	})
	if err != nil {
		return nil, err
	}

	// Create a slice of all VM names.
	names := []string{}
	listVMsOutLines := strings.Split(string(listVMsOut), newLine)
	for _, line := range listVMsOutLines {
		if line != "" {
			name := strings.Split(line, quotationMark)[nameIndex]
			names = append(names, name)
		}
	}

	// Create a slice of VM statuses.
	statuses := []provider.Status{}
	for _, name := range names {
		// Create the VM status.
		status := provider.Status{
			ID: name,
			IP: localhost,
		}

		// Get the info for this VM.
		showVMInfoOut, err := p.vboxmanage(vboxManageCommand{
			description: fmt.Sprintf("show VM info for VM '%s'", status.ID),
			args: []string{
				"showvminfo", status.ID,
			},
		})
		if err != nil {
			return nil, err
		}
		showVMInfoLines := strings.Split(showVMInfoOut, newLine)

		// Get the port for this VM.
		for _, line := range showVMInfoLines {
			if strings.Contains(line, portLine) {
				portStringWithSuffix := strings.Split(line, space)[portIndex]
				portString := strings.TrimSuffix(portStringWithSuffix, comma)

				port, err := strconv.Atoi(portString)
				if err != nil {
					return nil, err
				}

				status.Port = port
			}
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}
