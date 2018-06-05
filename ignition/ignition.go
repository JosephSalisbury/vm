package ignition

import (
	"errors"
	"log"
	"os"
)

var (
	// InvalidConfigError is the error returned if the configuration is invalid.
	InvalidConfigError = errors.New("invalid config")
)

type Config struct {
	Logger *log.Logger
	// Path is the location of the Ignition Config on disk.
	Path string
}

// Ignition handles Ignition Configs.
type Ignition struct {
	logger *log.Logger
	path   string
}

// New takes a Config, and returns an Ignition.
func New(config Config) (*Ignition, error) {
	if config.Logger == nil {
		return nil, InvalidConfigError
	}
	if config.Path == "" {
		return nil, InvalidConfigError
	}
	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		return nil, InvalidConfigError
	}

	i := &Ignition{
		path: config.Path,
	}

	return i, nil
}

// Path returns a path where the Ignition Config can be read from.
func (i *Ignition) Path() string {
	return i.path
}
