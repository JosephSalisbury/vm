package file

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/JosephSalisbury/vm/ignition"
	"github.com/JosephSalisbury/vm/secrets"
	"github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"
)

var (
	// Ignition is an Ignition that uses local files.
	Ignition = ignition.Name("file")
)

type fileIgnition struct {
	logger  *log.Logger
	path    string
	secrets *secrets.Secrets

	tempPath string
}

func New(config ignition.Config) (ignition.Interface, error) {
	if config.Logger == nil {
		return nil, ignition.InvalidConfigError
	}
	if config.Path == "" {
		return nil, ignition.InvalidConfigError
	}
	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		return nil, ignition.InvalidConfigError
	}
	if config.Secrets == nil {
		return nil, ignition.InvalidConfigError
	}

	i := &fileIgnition{
		logger:  config.Logger,
		path:    config.Path,
		secrets: config.Secrets,

		tempPath: "",
	}

	return i, nil
}

// Create performs any necessary setup for the Ignition Config.
func (i *fileIgnition) Create() error {
	i.logger.Printf("creating ignition config")

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
func (i *fileIgnition) Path() (string, error) {
	return i.tempPath, nil
}
