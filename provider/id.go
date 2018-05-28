package provider

import (
	"math/rand"
	"time"
)

var (
	idLength = 6
	letters  = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ID returns a unique identifier suitable for naming VMs.
func ID() string {
	id := make([]rune, idLength)

	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}

	return string(id)
}
