// Package google provides a Provider implementation backed by Google Cloud.
package google

import (
	"io/ioutil"
	"log"

	"github.com/JosephSalisbury/vm/provider"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

var (
	// Provider is a Provider that uses Google Cloud.
	Provider = provider.Name("google")
)

type googleProvider struct {
	logger           *log.Logger
	instancesService *compute.InstancesService

	preemptible bool
	project     string
	zone        string
}

func New(config provider.Config) (provider.Interface, error) {
	if config.Logger == nil {
		return nil, provider.InvalidConfigError
	}

	if config.GoogleCredentialsFilePath == "" {
		return nil, provider.InvalidConfigError
	}
	if config.GoogleProject == "" {
		return nil, provider.InvalidConfigError
	}
	if config.GoogleZone == "" {
		return nil, provider.InvalidConfigError
	}

	jsonData, err := ioutil.ReadFile(config.GoogleCredentialsFilePath)
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(jsonData, "https://www.googleapis.com/auth/compute")
	if err != nil {
		return nil, err
	}

	client := conf.Client(oauth2.NoContext)
	service, err := compute.New(client)
	if err != nil {
		return nil, err
	}
	instancesService := compute.NewInstancesService(service)

	p := &googleProvider{
		logger:           config.Logger,
		instancesService: instancesService,

		preemptible: config.GooglePreemptible,
		project:     config.GoogleProject,
		zone:        config.GoogleZone,
	}

	return p, nil
}
