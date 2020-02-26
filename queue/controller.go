package queue

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/MichaelWittgreffe/jobengine/models"
)

//Controller is a handler for interacting with queues
type Controller struct {
	queues            map[string]*models.Queue
	mutex             sync.Mutex
	jobKeepMinutes    int
	jobTimeoutMinutes int
}

//NewController is a constructor for Controller, creating a blank, new controller
func NewController(jobKeepMinutes, jobTimeoutMinutes int) *Controller {
	return &Controller{
		queues:            make(map[string]*models.Queue),
		mutex:             sync.Mutex{},
		jobKeepMinutes:    jobKeepMinutes * 60,
		jobTimeoutMinutes: jobTimeoutMinutes * 60,
	}
}

//NewControllerFromConfig is a constructor for Controller, creating a controller from a config file
func NewControllerFromConfig(cfg *models.Config) (*Controller, error) {
	result := NewController(cfg.JobKeepMinutes, cfg.JobTimeoutMinutes)

	for name, permissions := range cfg.Queues {
		if err := result.AddNewQueue(name, permissions); err != nil {
			return nil, fmt.Errorf("Error Creating Queue %s: %s", name, err.Error())
		}
	}

	return result, nil
}

//NewControllerFromDB is a constructor for Controller, creating from an existing store rather than from scratch
func NewControllerFromDB(cfg *models.Config, db *DBFile) (*Controller, error) {
	result := NewController(cfg.JobKeepMinutes, cfg.JobTimeoutMinutes)
	db.LoadController(result)

	err := db.LoadFromFile()
	if err != nil {
		return result, fmt.Errorf("Error Loading DB: %s", err.Error())
	}

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

//ExportQueues returns all the loaded queues as a map
func (c *Controller) ExportQueues() (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(c.queues))
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for name, data := range c.queues {
		result[name] = c.exportQueue(name, "", data)
	}

	return result, nil
}

//ExportQueue returns a queue as a map[string]interface
func (c *Controller) ExportQueue(name, status string) (map[string]interface{}, error) {
	if len(name) <= 0 {
		return nil, fmt.Errorf("Invalid Arg")
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	data, located := c.queues[name]
	if !located {
		return nil, nil
	}

	result := c.exportQueue(name, status, data)
	return result, nil
}

//exportQueue takes the queue name and a Queue model, and transforms them into a map[string]interface{}
func (c *Controller) exportQueue(name, status string, data *models.Queue) map[string]interface{} {
	jobs := make(map[string]interface{}, data.Size)

	for i, job := range data.Jobs {
		if len(status) == 0 || status == job.State {
			jobs[strconv.Itoa(i)] = map[string]interface{}{
				"uid":          job.UID,
				"content":      job.Content,
				"state":        job.State,
				"last_updated": job.LastUpdated,
				"created":      job.Created,
				"timeout_time": job.TimeoutTime,
				"priority":     job.Priority,
			}
		}
	}

	permissions := make(map[string][]string, 2)
	permissions["read"] = c.copyStringSlice(data.Permissions.Read)
	permissions["write"] = c.copyStringSlice(data.Permissions.Write)

	return map[string]interface{}{
		"name":        data.Name,
		"size":        data.Size,
		"permissions": permissions,
		"jobs":        jobs,
	}
}

//LoadQueues loads the results of ExportQueues into memory
func (c *Controller) LoadQueues(in map[string]interface{}) error {
	if in == nil {
		return fmt.Errorf("Invalid Arg")
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queues = make(map[string]*models.Queue, len(in))

	for queueName, queueData := range in {
		queueData := queueData.(map[string]interface{})
		permData := queueData["permissions"].(map[string]interface{})
		jobData := queueData["jobs"].(map[string]interface{})
		queueEntry := &models.Queue{
			Name: queueName,
			Size: uint8(queueData["size"].(float64)),
			Jobs: make([]*models.Job, uint8(queueData["size"].(float64))),
			Permissions: &models.QueuePermissions{
				Read:  c.interfaceSliceToStringSlice(permData["read"].([]interface{})),
				Write: c.interfaceSliceToStringSlice(permData["write"].([]interface{})),
			},
		}

		for i, job := range jobData {
			job := job.(map[string]interface{})
			iConv, err := strconv.Atoi(i)
			if err != nil {
				return fmt.Errorf("Failed To Assert Job Queue Key As Int: %s", i)
			}
			queueEntry.Jobs[iConv] = &models.Job{
				UID:         job["uid"].(string),
				Content:     job["content"].(map[string]interface{}),
				State:       job["state"].(string),
				LastUpdated: int64(job["last_updated"].(float64)),
				Created:     int64(job["created"].(float64)),
				TimeoutTime: int64(job["timeout_time"].(float64)),
				Priority:    uint8(job["priority"].(float64)),
			}
		}

		c.queues[queueName] = queueEntry
	}

	return nil
}

//QueueExists checks if the given queue exists, and whether the user is allowed to see it
func (c *Controller) QueueExists(queueName, appName string) (bool, bool, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	queue, found := c.queues[queueName]

	if found {
		readAllowed := false
		writeAllowed := false

		for _, name := range queue.Permissions.Read {
			if name == appName {
				readAllowed = true
				break
			}
		}

		for _, name := range queue.Permissions.Write {
			if name == appName {
				writeAllowed = true
				break
			}
		}

		if !readAllowed && !writeAllowed {
			found = false //user cannot see this queue
		}

		return found, readAllowed, writeAllowed
	}

	return false, false, false
}

//AddNewJob adds the given job entry to the queue - must check write permission yourself first
func (c *Controller) AddNewJob(queueName string, in *models.Job) error {
	if len(queueName) <= 0 || in == nil {
		return fmt.Errorf("Invalid Args")
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	queue, found := c.queues[queueName]
	if !found {
		return fmt.Errorf("Queue Not Found")
	}

	queue.Jobs = append(queue.Jobs, in)
	sort.Slice(queue.Jobs, func(i, j int) bool {
		return queue.Jobs[i].Priority > queue.Jobs[j].Priority
	})
	queue.Size++
	return nil
}

//interfaceSliceToStringSlice converts a slice of interface types to a slice of string types
func (c *Controller) interfaceSliceToStringSlice(in []interface{}) []string {
	output := make([]string, len(in))

	for i, data := range in {
		output[i] = data.(string)
	}

	return output
}

//copyStringSlice copies a slice of strings into a new slice of strings
func (c *Controller) copyStringSlice(in []string) []string {
	result := make([]string, len(in))
	copy(result, in)
	return result
}

//UpdateQueue iterates over the queue and resolves any timeouts/mark as failed, should be called prior to read requests
func (c *Controller) UpdateQueue(queueName string) error {
	if len(queueName) <= 0 {
		return fmt.Errorf("Invalid Arg")
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	queue, found := c.queues[queueName]
	if !found {
		return fmt.Errorf("Queue Not Found")
	}

	currentTime := time.Now().Unix()
	indexToDelete := make([]int, 0)

	for i, job := range queue.Jobs {
		if (job.State == models.Complete || job.State == models.Failed) && job.LastUpdated < (currentTime-int64(c.jobKeepMinutes)) {
			//remove complete/failed jobs that are outside the keep window
			indexToDelete = append(indexToDelete, i)
		} else if job.State == models.Inprogress && (job.LastUpdated < (currentTime - int64(c.jobTimeoutMinutes))) {
			//mark as failed if no update within the timeout cut-off
			job.State = models.Failed
			job.LastUpdated = currentTime
		} else if (job.State == models.Queued) && (currentTime > job.TimeoutTime) {
			//delete queued jobs that are timed out
			indexToDelete = append(indexToDelete, i)
		}
	}

	for _, indexToDelete := range indexToDelete {
		c.deleteJobAtIndex(queue, indexToDelete)
	}

	c.queues[queueName] = queue
	return nil
}

/*GetNextJob returns the next job in the specified queue, at status 'queued', will also remove jobs that are timed-out
or are complete/failed and over the keep period*/
func (c *Controller) GetNextJob(queueName string) (*models.Job, error) {
	if len(queueName) <= 0 {
		return nil, fmt.Errorf("Invalid Arg")
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	queue, found := c.queues[queueName]
	if !found {
		return nil, nil
	}

	for _, job := range queue.Jobs {
		if job.State == models.Queued {
			job.State = models.Inprogress
			job.LastUpdated = time.Now().Unix()
			return job, nil //success return
		}
	}

	return nil, nil
}

//deleteJobAtIndex removes the given index from the job queue, cleans up memory during delete
func (c *Controller) deleteJobAtIndex(queue *models.Queue, i int) {
	queueLenMinus := len(queue.Jobs) - 1
	if i < (queueLenMinus) {
		copy(queue.Jobs[i:], queue.Jobs[i+1:])
	}
	queue.Jobs[queueLenMinus] = nil
	queue.Jobs = queue.Jobs[:queueLenMinus]
	queue.Size = uint8(len(queue.Jobs))
}
