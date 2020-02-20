package configload

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/MichaelWittgreffe/jobengine/filesystem"
	"github.com/MichaelWittgreffe/jobengine/models"
	"gopkg.in/yaml.v2"
)

//ConfigLoad is a struct for loading the config from file
type ConfigLoad struct {
	filePath    string
	fileHandler filesystem.FileSystem
}

//NewConfigLoader acts as a constructor for the ConfigLoader
func NewConfigLoader(filepath, fsType string) ConfigLoader {
	if len(filepath) <= 0 || len(fsType) <= 0 {
		return nil
	}

	return &ConfigLoad{
		filePath:    filepath,
		fileHandler: filesystem.NewFileSystem(fsType),
	}
}

//LoadFromFile loads the setup config file into memory
func (l *ConfigLoad) LoadFromFile(version float64) (*models.Config, error) {
	exists, err := l.fileHandler.FileExists(l.filePath)
	if !exists {
		return nil, fmt.Errorf("Not Found")
	}

	rawFileContent, err := ioutil.ReadFile(l.filePath)
	if err != nil {
		return nil, err
	}

	configFile := make(map[interface{}]interface{})

	err = yaml.Unmarshal(rawFileContent, &configFile)
	if err != nil {
		return nil, err
	}

	result := new(models.Config)

	switch tempVersion := configFile["version"].(type) {
	case float64:
		result.Version = tempVersion
	case int:
		result.Version = float64(tempVersion)
	default:
		return nil, errors.New("Invalid Version Value")
	}

	if result.Version != version {
		return nil, fmt.Errorf("Invalid Version, Require %f Got %f", version, result.Version)
	}

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
}

//SaveToFile saves the given config to file
func (l *ConfigLoad) SaveToFile(cfg *models.Config) error {
	if exists, err := l.fileHandler.FileExists(l.filePath); exists {
		if err := l.fileHandler.DeleteFile(l.filePath); err != nil {
			return fmt.Errorf("Unable To Delete Cfg: %s", err.Error())
		}
	} else if err != nil {
		return fmt.Errorf("Error Locating Config File: %s", err.Error())
	}

	outputData, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(l.filePath, outputData, 0644)
}

//interfaceSliceToStringSlice converts the given slice of interface types to string types
func (l *ConfigLoad) interfaceSliceToStringSlice(input []interface{}) []string {
	output := make([]string, len(input))

	for i, value := range input {
		if valueToAdd := value.(string); len(valueToAdd) > 0 {
			output[i] = valueToAdd
		}
	}

	return output
}
