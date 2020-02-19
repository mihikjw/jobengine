package config

import (
	"errors"
	"io/ioutil"

	"github.com/MichaelWittgreffe/jobengine/filesystem"
	"github.com/MichaelWittgreffe/jobengine/models"
	"gopkg.in/yaml.v2"
)

//ConfigLoader is a struct for loading the config from file
type ConfigLoader struct {
	filePath    string
	fileHandler filesystem.FileSystem
}

//NewConfig acts as a constructor for the ConfigLoader
func NewConfigLoader(filepath, fsType string) *ConfigLoader {
	if len(filepath) > 0 && len(fsType) > 0 {
		return &ConfigLoader{
			filePath:    filepath,
			fileHandler: filesystem.NewFileSystem(fsType),
		}
	}

	return nil
}

//LoadFromFile loads the setup config file into memory
func (c *ConfigLoader) LoadFromFile() (*models.Config, error) {
	if exists, err := c.fileHandler.FileExists(c.filePath); exists {
		if rawFileContent, err := ioutil.ReadFile(c.filePath); err == nil {
			configFile := make(map[interface{}]interface{})

			if err = yaml.Unmarshal(rawFileContent, &configFile); err == nil {

			}
		}

		return nil, err
	} else if err == nil {
		return nil, errors.New("Config File Not Found")
	} else {
		return nil, errors.New("Unable To Load Config: " + err.Error())
	}
}
