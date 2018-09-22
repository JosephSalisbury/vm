// Package providerset provides utilities for working with multiple Providers.
package providerset

import (
	"errors"
	"sort"
	"strings"

	"github.com/JosephSalisbury/vm/provider"
	"github.com/JosephSalisbury/vm/provider/google"
	"github.com/JosephSalisbury/vm/provider/virtualbox"
)

var (
	// DefaultProvider is the Provider to use if no Provider is specified by the user.
	DefaultProvider = google.Provider

	// providerConstructors is a mapping between ProviderNames and Providers.
	providerConstructors = map[provider.Name]func(provider.Config) (provider.Interface, error){
		google.Provider:     google.New,
		virtualbox.Provider: virtualbox.New,
	}

	// UnknownProviderError is the error returned if the ProviderName given is unknown.
	UnknownProviderError = errors.New("unknown provider")
)

// New takes a ProviderName and a Config, and returns a new Provider of the correct type.
func New(name provider.Name, config provider.Config) (provider.Interface, error) {
	providerConstructor, ok := providerConstructors[name]
	if !ok {
		return nil, UnknownProviderError
	}

	return providerConstructor(config)
}

// Names returns a comma seperated list of all ProviderNames.
// This is mainly useful for showing available Providers.
func Names() string {
	providerNames := []string{}

	for providerName, _ := range providerConstructors {
		providerNames = append(providerNames, string(providerName))
	}

	sort.Strings(providerNames)

	return strings.Join(providerNames, ", ")
}
