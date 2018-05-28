package virtualbox

import (
	"math/rand"
)

const (
	// TODO: Use actual port limits (/proc?)
	lowPortNumber  = 3000
	highPortNumber = 65000
)

// getFreePort returns a port number that can be used.
func getFreePort() int {
	// TODO: Actually check port is free.
	return rand.Intn(highPortNumber-lowPortNumber) + lowPortNumber
}
