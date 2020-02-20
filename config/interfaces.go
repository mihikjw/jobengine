package config

import (
	"github.com/MichaelWittgreffe/jobengine/models"
)

type Loader interface {
	LoadFromFile(version float64) (*models.Config, error)
	SaveToFile(cfg *models.Config) error
}
