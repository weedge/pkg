package configparser

import (
	"errors"
	"fmt"
)

type fileProvider struct{}

// NewFile create file provider
func NewFile() Provider {
	return &fileProvider{}
}

// Get file provider
func (fl *fileProvider) Get() (*Parser, error) {
	fileName := getConfigFlag()
	if fileName == "" {
		return nil, errors.New("config file not specified")
	}

	cp, err := NewParserFromFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error loading config file %q: %v", fileName, err)
	}

	return cp, nil
}
