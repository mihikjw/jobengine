package configload

import (
	"github.com/MichaelWittgreffe/jobengine/models"
)

//ConfigLoader represents an interface for a struct that can load and save a config
type ConfigLoader interface {
	LoadFromFile(version float64) (*models.Config, error)
	SaveToFile(cfg *models.Config) error
}

//ConfigParser is an interface for marshaling/unmarshaling the config format
type ConfigParser interface {
	Unmarshal(in []byte) (map[interface{}]interface{}, error)
	Marshal(in interface{}) ([]byte, error)
}

//NewConfigParser constructs a parser for the given config file format
func NewConfigParser(parserType string) ConfigParser {
	switch {
	case parserType == "YAML" || parserType == "yaml":
		return new(YAMLHandler)
	default:
		return nil
	}
}
