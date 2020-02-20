package configload

import (
	"github.com/MichaelWittgreffe/jobengine/models"
)

//ConfigLoader represents an interface for a struct that can load and save a config
type ConfigLoader interface {
	LoadFromFile(version float64) (*models.Config, error)
	SaveToFile(cfg *models.Config) error
}
