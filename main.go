package main

import (
	"errors"
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/config"
	"github.com/MichaelWittgreffe/jobengine/models"
)

const version float64 = 1.0
const defaultPort int = 6006

func main() {
	if _, err := loadConfig("./examples/config.yml"); err == nil {
		fmt.Println("Config Loaded")
	} else {
		fmt.Println(err.Error())
	}
}

//loadConfig loads the given config and returns it
func loadConfig(cfgPath string) (*models.Config, error) {
	var cfg *models.Config
	var err error

	if loader := config.NewConfigLoader(cfgPath, "os"); loader != nil {
		if cfg, err = loader.LoadFromFile(version); err == nil {
			return cfg, nil
		} else if err.Error() == "Not Found" {
			cfg = &models.Config{
				Version: version,
				Port:    defaultPort,
				Queues:  make(map[string]*models.Queue),
			}

			if err = loader.SaveToFile(cfg); err != nil {
				return nil, errors.New(fmt.Sprintf("Failed To Save Default Cfg: %s", err.Error()))
			} else {
				return cfg, nil
			}
		}
	} else {
		err = errors.New("Unable To Create ConfigLoader")
	}

	return nil, err
}
