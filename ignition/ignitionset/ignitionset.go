package ignitionset

import (
	"errors"
	"sort"
	"strings"

	"github.com/JosephSalisbury/vm/ignition"
	"github.com/JosephSalisbury/vm/ignition/file"
)

var (
	// DefaultIgnition is the Ignition to use if no Ignition is specified.
	DefaultIgnition = file.Ignition

	ignitionConstructors = map[ignition.Name]func(ignition.Config) (ignition.Interface, error){
		file.Ignition: file.New,
	}

	// UnknownIgnitionError is the error returned if the IgnitionName given is unknown.
	UnknownIgnitionError = errors.New("unknown ignition")
)

// New takes an IgnitionName and a Config, and returns a new Ignition of the correct type.
func New(name ignition.Name, config ignition.Config) (ignition.Interface, error) {
	ignitionConstructor, ok := ignitionConstructors[name]
	if !ok {
		return nil, UnknownIgnitionError
	}

	return ignitionConstructor(config)
}

// Names returns a comma seperated list of all IgnitionNames.
// This is mainly useful for showing available Ignitions.
func Names() string {
	ignitionNames := []string{}

	for ignitionName, _ := range ignitionConstructors {
		ignitionNames = append(ignitionNames, string(ignitionName))
	}

	sort.Strings(ignitionNames)

	return strings.Join(ignitionNames, ", ")
}
