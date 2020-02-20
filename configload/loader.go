package configload

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/MichaelWittgreffe/jobengine/filesystem"
	"github.com/MichaelWittgreffe/jobengine/models"
	"gopkg.in/yaml.v2"
)

//FileLoader is a struct for loading the config from file
type FileLoader struct {
	filePath    string
	fileHandler filesystem.FileSystem
}

//NewConfigLoader acts as a constructor for the ConfigLoader
func NewConfigLoader(filepath, fsType string) ConfigLoader {
	if len(filepath) > 0 && len(fsType) > 0 {
		return &FileLoader{
			filePath:    filepath,
			fileHandler: filesystem.NewFileSystem(fsType),
		}
	}

	return nil
}

//LoadFromFile loads the setup config file into memory
func (l *FileLoader) LoadFromFile(version float64) (*models.Config, error) {
	var err error
	var exists bool

	if exists, err = l.fileHandler.FileExists(l.filePath); exists {
		var rawFileContent []byte

		if rawFileContent, err = ioutil.ReadFile(l.filePath); err == nil {
			result := new(models.Config)
			configFile := make(map[interface{}]interface{})

			if err = yaml.Unmarshal(rawFileContent, &configFile); err == nil {
				//support both int and float64 types for version
				if tempFloatVersion, valid := configFile["version"].(float64); valid {
					result.Version = tempFloatVersion
				} else {
					if tempIntVersion, valid := configFile["version"].(int); valid {
						result.Version = float64(tempIntVersion)
					} else {
						err = errors.New("Invalid Version Value")
					}
				}

				if result.Version == version {
					result.Port = configFile["port"].(int)
					queues := configFile["queues"].(map[interface{}]interface{})
					result.Queues = make(map[string]*models.QueuePermissions, len(queues))

					for name, data := range queues {
						newData := data.(map[interface{}]interface{})
						queue := models.QueuePermissions{
							Read:  l.interfaceSliceToStringSlice(newData["read"].([]interface{})),
							Write: l.interfaceSliceToStringSlice(newData["write"].([]interface{})),
						}

						result.Queues[name.(string)] = &queue
					}

					return result, nil
				} else {
					err = fmt.Errorf("Invalid Version, Require %f Got %f", version, result.Version)
				}
			}
		}

		return nil, err
	} else if err == nil {
		return nil, fmt.Errorf("Not Found")
	} else {
		return nil, fmt.Errorf("Unable To Load Config: %s" + err.Error())
	}
}

//SaveToFile saves the given config to file
func (l *FileLoader) SaveToFile(cfg *models.Config) error {
	var err error
	var exists bool

	if exists, err = l.fileHandler.FileExists(l.filePath); exists {
		err = l.fileHandler.DeleteFile(l.filePath)
	}

	if err == nil {
		var outputData []byte

		if outputData, err = yaml.Marshal(cfg); err == nil {
			err = ioutil.WriteFile(l.filePath, outputData, 0644)
		}
	} else {
		err = fmt.Errorf("Unable To Delete Cfg: %s", err.Error())
	}

	return err
}

//interfaceSliceToStringSlice converts the given slice of interface types to string types
func (l *FileLoader) interfaceSliceToStringSlice(input []interface{}) []string {
	output := make([]string, len(input))

	for i, value := range input {
		if valueToAdd := value.(string); len(valueToAdd) > 0 {
			output[i] = valueToAdd
		}
	}

	return output
}
