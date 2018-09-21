// Package google provides a Provider implementation backed by Google Cloud.
package google

import (
	"log"
	"os/exec"

	"github.com/JosephSalisbury/vm/provider"
)

const (
	zone = "us-central1-a"
)

var (
	// Provider is a Provider that uses Google Cloud.
	Provider = provider.Name("Google")
)

type googleProvider struct {
	logger *log.Logger
}

func New(config provider.Config) (provider.Interface, error) {
	if config.Logger == nil {
		return nil, provider.InvalidConfigError
	}

	if _, err := exec.LookPath("gcloud"); err != nil {
		return nil, GCloudMissingError
	}

	p := &googleProvider{
		logger: config.Logger,
	}

	return p, nil
}

type gCloudCommand struct {
	description string
	args        []string
}

func (p *googleProvider) gcloud(command gCloudCommand) (string, error) {
	p.logger.Printf("executing: %s", command.description)

	out, err := exec.Command("gcloud", command.args...).Output()
	stringOut := string(out)

	if err != nil {
		p.logger.Printf("gcloud output: %s %s", err, stringOut)
	}

	return stringOut, err
}
