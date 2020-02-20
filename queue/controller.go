package queue

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/models"
)

//Controller is a handler for interacting with queues
type Controller struct {
	queues map[string]*models.Queue
}

//NewController is a constructor for Controller, creating a blank, new controller
func NewController() *Controller {
	return &Controller{
		queues: make(map[string]*models.Queue),
	}
}

//NewControllerFromConfig is a constructor for Controller, creating a controller from a config file
func NewControllerFromConfig(cfg *models.Config) (*Controller, error) {
	result := NewController()

	for name, permissions := range cfg.Queues {
		if err := result.AddNewQueue(name, permissions); err != nil {
			return nil, fmt.Errorf("Error Creating Queue %s: %s", name, err.Error())
		}
	}

	return result, nil
}

//NewControllerFromDB is a constructor for Controller, creating from an existing store rather than from scratch
func NewControllerFromDB(cfg *models.Config, db *DBFile) (*Controller, error) {
	/*
		this is where any existing queue file would be loaded
		we also need to diff this against the config, remove any queues not in the cfg
		and add any that are missing - this is okay to do as when a queue is added from
		the API the cfg is also updated - therefore only changes would be manual by an admin
	*/
	return nil, fmt.Errorf("Not Implemented")
}

//AddNewQueue adds a new queue to the controller, from the given queue name and queue permissions set
func (c *Controller) AddNewQueue(name string, permissions *models.QueuePermissions) error {
	if len(name) <= 0 || permissions == nil {
		return fmt.Errorf("Invalid Args")
	}

	if _, exists := c.queues[name]; exists {
		return fmt.Errorf("Queue %s Already Exists", name)
	}

	c.queues[name] = &models.Queue{
		Jobs:        make([]*models.Job, 0),
		Permissions: permissions,
		Size:        0,
	}

	return nil
}
