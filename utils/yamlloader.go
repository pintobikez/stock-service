package utils

import (
	errors "github.com/pkg/errors"
)

var (
	ErrInvalidFile       = "Invalid file absolute path"
	ErrUnableToReadFile  = "Unable to read the file storage"
	ErrUnableToParseFile = "Unable to parse the file storage"
)

type YmlConfig struct {
    Driver struct {
        Host string
	    User string
	    Pw string
	    Port int
	    Schema string
    }
}

func loadConfigFile(string filename) *YmlConfig, error {
	filename, err := filepath.Abs(filename)

	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidFile)
	}

	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, errors.Wrap(err, ErrUnableToReadFile)
	}

	conf := new(YmlConfig)

	err = yaml.Unmarshal(file, conf)

	if err != nil {
		return nil, errors.Wrap(err, ErrUnableToParseFile)
	}

	return conf, nil
}