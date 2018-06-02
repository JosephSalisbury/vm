package virtualbox

import (
	"time"
)

func (p *virtualBoxProvider) Delete(id string) error {
	if _, err := p.vboxmanage(vboxManageCommand{
		description: "power off VM",
		args: []string{
			"controlvm", id, "poweroff",
		},
	}); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	if _, err := p.vboxmanage(vboxManageCommand{
		description: "delete VM",
		args: []string{
			"unregistervm", id, "--delete",
		},
	}); err != nil {
		return err
	}

	return nil
}
