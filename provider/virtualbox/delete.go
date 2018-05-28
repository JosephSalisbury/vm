package virtualbox

import (
	"os/exec"
	"time"
)

func (p *virtualBoxProvider) Delete(id string) error {
	p.logger.Printf("powering off VM %s", id)
	if err := exec.Command("VBoxManage", "controlvm", id, "poweroff").Run(); err != nil {
		p.logger.Printf("could not poweroff VM: %s", err)
	}

	time.Sleep(1 * time.Second)

	p.logger.Printf("deleting VM %s", id)
	if err := exec.Command("VBoxManage", "unregistervm", id, "--delete").Run(); err != nil {
		return err
	}

	return nil
}
