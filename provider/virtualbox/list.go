package virtualbox

import (
	"os/exec"
	"strings"

	"github.com/JosephSalisbury/vm/provider"
)

const (
	nameIndex = 1

	newLine       = "\n"
	quotationMark = "\""
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

	// Ping network to refresh arp.
	if err := exec.Command("bash", "-c", "for i in {1..254}; do ping -c 1 192.168.1.$i & done").Run(); err != nil {
		return nil, err
	}

	// Create a slice of VM statuses.
	statuses := []provider.Status{}
	for _, name := range names {
		// Create the VM status.
		status := provider.Status{
			ID:   name,
			Port: 22,
		}

		// Grab IP from arp.
		output, err := exec.Command("bash", "-c", "arp -a | grep -E '192.168.1.*' | grep -v 'incomplete' | grep -v '254' | grep -v '255' | tail -n 1 | awk -F '(' '{print $2}' | awk -F ')' '{print $1}'").Output()
		if err != nil {
			return nil, err
		}

		status.IP = strings.TrimSuffix(string(output), newLine)

		statuses = append(statuses, status)
	}

	return statuses, nil
}
