package ignition

import (
	"errors"
	"log"
	"net/url"

	"github.com/JosephSalisbury/vm/secrets"
)

var (
	// InvalidConfigError is the error returned if the configuration is invalid.
	InvalidConfigError = errors.New("invalid config")
)

// Name is the name of an Ignition.
type Name string

// Interface is an interface for managing Ignition Configs.
type Interface interface {
	// Create performs any necessary setup for the Ignition Config.
	Create() error

	// Path returns a path where the Ignition Config can be read from.
	Path() (string, error)
}

type Config struct {
	// Logger is a logger that Ignitions can use for outputting information.
	Logger *log.Logger

	// Path is the location of the Ignition Config on disk.
	Path string

	// URL is the location of the Container Linux Config.
	URL url.URL

	// Secrets are sensitive material that may need to be referenced.
	Secrets *secrets.Secrets
}
