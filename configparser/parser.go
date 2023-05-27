package configparser

import (
	"fmt"

	"github.com/spf13/viper"
)

// Parser loads configuration.
type Parser struct {
	v *viper.Viper
}

func newViper() *viper.Viper {
	return viper.NewWithOptions()
}

// UnmarshalExact unmarshals the config into a struct, erroring if a field is nonexistent.
func (l *Parser) UnmarshalExact(intoCfg interface{}) error {
	return l.v.UnmarshalExact(intoCfg)
}

// NewParserFromFile creates a new Parser by reading the given file.
func NewParserFromFile(fileName string) (*Parser, error) {
	v := newViper()
	v.SetConfigFile(fileName)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("unable to read the file %v: %w", fileName, err)
	}
	return &Parser{v: v}, nil
}
