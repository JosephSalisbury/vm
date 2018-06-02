package virtualbox

import (
	"errors"
)

var (
	VBoxConfigdriveGenMissingError = errors.New("vbox-configdrive-gen missing")
	VBoxManageMissingError         = errors.New("VBoxManage missing")
)
