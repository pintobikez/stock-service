package utils

import (
	errors "github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

var (
	ErrInvalidFile       = "Invalid file absolute path"
	ErrUnableToReadFile  = "Unable to read the file storage"
	ErrUnableToParseFile = "Unable to parse the file storage"
)

type YmlConfig struct {
	Driver struct {
		Host   string
		User   string
		Pw     string
		Port   int
		Schema string
	}
}

func LoadConfigFile(filename string) (*YmlConfig, error) {

	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidFile)
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, ErrUnableToReadFile)
	}

	conf := new(YmlConfig)

	if err = yaml.Unmarshal(file, conf); err != nil {
		return nil, errors.Wrap(err, ErrUnableToParseFile)
	}

	return conf, nil
}
