package provider

import (
	"errors"
	"log"

	"github.com/JosephSalisbury/vm/ignition"
)

var (
	// InvalidConfigError is the error returned if the configuration is invalid.
	InvalidConfigError = errors.New("invalid config")

	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// Name is the name of a Provider.
type Name string

// Interface is an interface for managing CoreOS Container Linux Virtual Machines.
type Interface interface {
	// Create creates a CoreOS Container Linux Virtual Machine,
	// with the latest version of the given channel,
	// the given Ignition Config,
	// `cpu` cores, and `ram` GB of RAM.
	Create(channel string, ignition ignition.Interface, cpu int, ram int) error

	// Delete deletes the specified VM.
	Delete(id string) error

	// List returns a slice of Status,
	// describing all current CoreOS Container Linux Virtual Machines.
	List() ([]Status, error)
}

// Config represents the configuration for creating a Provider.
type Config struct {
	// Logger is a logger that Providers can use for outputting information.
	Logger *log.Logger
}

type Status struct {
	// ID is the unique identifier of the VM.
	ID string

	// IP is the IP that the VM is listening on.
	IP string

	// TODO: This should be a map of host -> guest ports.
	// Port is the port that the VM is listening on.
	Port int
}
