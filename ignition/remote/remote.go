package remote

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"github.com/coreos/container-linux-config-transpiler/config"
	"github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/JosephSalisbury/vm/ignition"
	"github.com/JosephSalisbury/vm/secrets"
)

var (
	// Ignition is an Ignition that fetches a Container Linux Config
	// from a URL, and compiles it to an Ignition Config.
	Ignition = ignition.Name("remote")
)

type remoteIgnition struct {
	logger  *log.Logger
	url     url.URL
	secrets *secrets.Secrets

	tempPath string
}

func New(config ignition.Config) (ignition.Interface, error) {
	if config.Logger == nil {
		return nil, ignition.InvalidConfigError
	}
	if config.URL.String() == "" {
		return nil, ignition.InvalidConfigError
	}
	if config.Secrets == nil {
		return nil, ignition.InvalidConfigError
	}

	i := &remoteIgnition{
		logger:  config.Logger,
		url:     config.URL,
		secrets: config.Secrets,

		tempPath: "",
	}

	return i, nil
}

// Create downloads a Container Linux Config from the URL,
// compiles it to an Ignition Config, and adds any secrets.
func (i *remoteIgnition) Create() error {
	i.logger.Printf("downloading Container Linux Config")

	resp, err := http.Get(i.url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	containerLinuxConfig, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	i.logger.Printf("compiling Container Linux Config to Ignition Config")

	cfg, ast, _ := config.Parse(containerLinuxConfig)
	ignitionConfig, _ := config.Convert(cfg, "", ast)

	i.logger.Printf("adding secrets to Ignition Config")

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

	ignitionConfig.Storage.Files = append(ignitionConfig.Storage.Files, secretFiles...)

	updatedBytes, err := json.MarshalIndent(ignitionConfig, "", "  ")
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

// Path returns the path where the Ignition Config can be read from.
func (i *remoteIgnition) Path() (string, error) {
	// TODO: Error if the Ignition Config has not been created.
	return i.tempPath, nil
}
