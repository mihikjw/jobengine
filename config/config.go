package config

import (
	"errors"

	"github.com/MichaelWittgreffe/jobengine/filesystem"
)

//Config represents the application configuration
type Config struct {
	filePath    string
	fileHandler filesystem.FileSystem
}

//NewConfig acts as a constructor for the Config struct
func NewConfig(filepath, fsType string) *Config {
	if len(filepath) > 0 && len(fsType) > 0 {
		return &Config{
			filePath:    filepath,
			fileHandler: filesystem.NewFileSystem(fsType),
		}
	}

	return nil
}

//LoadFromFile loads the setup config file into memory
func (c *Config) LoadFromFile() error {
	if exists, err := c.fileHandler.FileExists(c.filePath); exists {
		return nil
	} else if err == nil {
		return errors.New("Config File Not Found")
	} else {
		return errors.New("Unable To Load Config: " + err.Error())
	}
}
