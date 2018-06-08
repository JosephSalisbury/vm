package ignition

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/JosephSalisbury/vm/secrets"
)

var (
	// InvalidConfigError is the error returned if the configuration is invalid.
	InvalidConfigError = errors.New("invalid config")
)

type Config struct {
	Logger *log.Logger
	// Path is the location of the Ignition Config on disk.
	Path string
	// Secrets are sensitive material that may need to be referenced by the config.
	Secrets *secrets.Secrets
}

// Ignition handles Ignition Configs.
type Ignition struct {
	logger  *log.Logger
	path    string
	secrets *secrets.Secrets

	tempPath string
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
	if config.Secrets == nil {
		return nil, InvalidConfigError
	}

	i := &Ignition{
		logger:  config.Logger,
		path:    config.Path,
		secrets: config.Secrets,
	}

	return i, nil
}

// Create performs any necessary setup for the Ignition Config.
func (i *Ignition) Create() error {
	bytes, err := ioutil.ReadFile(i.path)
	if err != nil {
		return err
	}

	var ignition types.Config
	if err := json.Unmarshal(bytes, &ignition); err != nil {
		return err
	}

	// TODO: Update hostname with VM ID.

	files, err := i.secrets.Files()
	if err != nil {
		return err
	}

	secretFiles := []types.File{}

	for _, file := range files {
		path := path.Join("/secrets/", filepath.Base(file.Name()))
		user := "joe"
		filesystem := "root"
		mode := 420

		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		source := (&url.URL{
			Scheme: "data",
			Opaque: "," + dataurl.Escape(fileContents),
		}).String()

		secretFile := types.File{
			Node: types.Node{
				Filesystem: filesystem,
				Group: &types.NodeGroup{
					Name: user,
				},
				Path: path,
				User: &types.NodeUser{
					Name: user,
				},
			},
			FileEmbedded1: types.FileEmbedded1{
				Contents: types.FileContents{
					Source:       source,
					Verification: types.Verification{},
				},
				Mode: &mode,
			},
		}
		secretFiles = append(secretFiles, secretFile)
	}

	ignition.Storage.Files = append(ignition.Storage.Files, secretFiles...)

	updatedBytes, err := json.MarshalIndent(ignition, "", "  ")
	if err != nil {
		return err
	}

	tempFile, err := ioutil.TempFile("", "ignition")
	if err != nil {
		return err
	}
	if _, err := tempFile.Write(updatedBytes); err != nil {
		return err
	}
	i.tempPath = tempFile.Name()

	i.logger.Printf("wrote ignition config to temp file: %s", i.tempPath)

	return nil
}

// Path returns a path where the Ignition Config can be read from.
func (i *Ignition) Path() string {
	return i.tempPath
}
