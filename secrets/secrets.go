package secrets

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var (
	// InvalidConfigError is the error returned if the configuration is invalid.
	InvalidConfigError = errors.New("invalid config")
)

type Config struct {
	Directory string
	Logger    *log.Logger
}

type Secrets struct {
	directory string
	logger    *log.Logger
}

func New(config Config) (*Secrets, error) {
	if config.Directory == "" {
		return nil, InvalidConfigError
	}
	if config.Logger == nil {
		return nil, InvalidConfigError
	}

	s := &Secrets{
		directory: config.Directory,
		logger:    config.Logger,
	}

	return s, nil
}

func (s *Secrets) Files() ([]*os.File, error) {
	fileInfos, err := ioutil.ReadDir(s.directory)
	if err != nil {
		return nil, err
	}

	s.logger.Printf("listed %v secret files", len(fileInfos))

	files := []*os.File{}
	for _, fileInfo := range fileInfos {
		file, err := os.Open(path.Join(s.directory, fileInfo.Name()))
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}
