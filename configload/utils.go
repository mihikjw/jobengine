package configload

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/models"
)

//version is the default version to use if no config is found
const version float64 = 1.0

//defaultPort is the default port to use if no config is found
const defaultPort int = 6006

//LoadConfig loads the given config and returns it
func LoadConfig(cfgPath, loaderType string) (*models.Config, error) {
	var cfg *models.Config
	var err error

	if loader := NewConfigLoader(cfgPath, loaderType); loader != nil {
		if cfg, err = loader.LoadFromFile(version); err == nil {
			return cfg, nil
		} else if err.Error() == "Not Found" {
			cfg = &models.Config{
				Version: version,
				Port:    defaultPort,
				Queues:  make(map[string]*models.QueuePermissions),
			}

			if err = loader.SaveToFile(cfg); err != nil {
				return nil, fmt.Errorf("Failed To Save Default Cfg: %s", err.Error())
			}

			return cfg, nil
		}
	} else {
		err = fmt.Errorf("Unable To Create ConfigLoader")
	}

	return nil, err
}
