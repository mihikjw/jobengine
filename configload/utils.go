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
func LoadConfig(loader ConfigLoader) (*models.Config, error) {
	if loader == nil {
		return nil, fmt.Errorf("Invalid Arg")
	}

	cfg, err := loader.LoadFromFile(version)

	if err == nil {
		return cfg, nil
	} else if err.Error() == "Not Found" {
		return createDefaultCfg(cfg, loader)
	}

	return nil, err
}

//createDefaultCfg generates and saves to file a default config for the application
func createDefaultCfg(in *models.Config, loader ConfigLoader) (*models.Config, error) {
	in.Version = version
	in.Port = defaultPort
	in.Queues = make(map[string]*models.QueuePermissions)
	in.JobKeepMinutes = 60
	in.JobTimeoutMinutes = 60

	if err := loader.SaveToFile(in); err != nil {
		return nil, fmt.Errorf("Failed To Save Default Cfg: %s", err.Error())
	}

	return in, nil
}
