package queue

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/MichaelWittgreffe/jobengine/models"
)

//Controller is a handler for interacting with queues
type Controller struct {
	queues map[string]*models.Queue
	mutex  sync.Mutex
}

//NewController is a constructor for Controller, creating a blank, new controller
func NewController() *Controller {
	return &Controller{
		queues: make(map[string]*models.Queue),
		mutex:  sync.Mutex{},
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
	result := NewController()
	db.LoadController(result)
	db.LoadFromFile()
	// any new queues in the config or removed queues should be resolved here
	// if there's any differences, call db.SaveToFile() before returning
	return result, nil
}

//AddNewQueue adds a new queue to the controller, from the given queue name and queue permissions set
func (c *Controller) AddNewQueue(name string, permissions *models.QueuePermissions) error {
	if len(name) <= 0 || permissions == nil {
		return fmt.Errorf("Invalid Args")
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.queues[name]; exists {
		return fmt.Errorf("Queue %s Already Exists", name)
	}

	c.queues[name] = &models.Queue{
		Jobs:        make([]*models.Job, 0),
		Permissions: permissions,
		Size:        0,
		Name:        name,
	}

	return nil
}

//ExportQueues returns the set of loaded queues as a map
func (c *Controller) ExportQueues() (map[string]interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	result := make(map[string]interface{}, len(c.queues))

	for name, data := range c.queues {
		jobs := make(map[string]interface{}, data.Size)

		for i, job := range data.Jobs {
			jobs[strconv.Itoa(i)] = map[string]interface{}{
				"uid":          job.UID,
				"content":      job.Content,
				"state":        job.State,
				"last_updated": job.LastUpdated,
				"created":      job.Created,
				"timeout_time": job.TimeoutTime,
			}
		}

		permissions := make(map[string][]string, 2)
		tmp1 := make([]string, len(data.Permissions.Read))
		copy(tmp1, data.Permissions.Read)
		permissions["read"] = tmp1
		tmp2 := make([]string, len(data.Permissions.Write))
		copy(tmp2, data.Permissions.Write)
		permissions["write"] = tmp2

		result[name] = map[string]interface{}{
			"name":        data.Name,
			"size":        data.Size,
			"permissions": permissions,
			"jobs":        jobs,
		}
	}

	return result, nil
}
