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
			result := new(models.Config)
			configFile := make(map[interface{}]interface{})

			if err = yaml.Unmarshal(rawFileContent, &configFile); err == nil {
				result.Version = configFile["version"].(float64)
				result.Port = configFile["port"].(int)
				queues := configFile["queues"].(map[interface{}]interface{})

				for name, data := range queues {
					newData := data.(map[interface{}]interface{})
					queue := models.Queue{
						Name:  name.(string),
						Read:  c.interfaceSliceToStringSlice(newData["read"].([]interface{})),
						Write: c.interfaceSliceToStringSlice(newData["write"].([]interface{})),
					}

					result.Queues = append(result.Queues, queue)
				}
			}
		}

		return nil, err
	} else if err == nil {
		return nil, errors.New("Config File Not Found")
	} else {
		return nil, errors.New("Unable To Load Config: " + err.Error())
	}
}

//interfaceSliceToStringSlice converts the given slice of interface types to string types
func (c *ConfigLoader) interfaceSliceToStringSlice(input []interface{}) []string {
	return nil
}
