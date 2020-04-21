package database

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/crypto"
)

// QueryController defines an object used to make queries to the database
type QueryController interface {
	CreateQueue(name, accessKey string) error
	GetQueue(name, accessKey string) (*Queue, error)
	DeleteQueue(name, accessKey string) error
}

// QueryControl object is used to make queries to the database
type QueryControl struct {
	db   *DBFile
	hash crypto.HashHandler
}

// NewQueryController is a constructor for the QueryController interface
func NewQueryController(db *DBFile, hasher crypto.HashHandler) QueryController {
	if db == nil {
		return nil
	}

	return &QueryControl{
		db:   db,
		hash: hasher,
	}
}

// CreateQueue creates a new queue entry
func (c *QueryControl) CreateQueue(name, accessKey string) error {
	if len(name) == 0 || len(accessKey) == 0 {
		return fmt.Errorf("Invalid Arg")
	}

	hashedKey, err := c.hash.Process(accessKey)
	if err != nil {
		return err
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	if _, found := c.db.Queues[name]; found {
		return fmt.Errorf("Queue Exists")
	}

	c.db.Queues[name] = &Queue{
		Name:      name,
		AccessKey: hashedKey,
		Size:      0,
		Jobs:      make([]*Job, 0),
	}

	return nil
}

// GetQueue returns the given queue object entry or nil if it cannot be found
func (c *QueryControl) GetQueue(name, accessKey string) (*Queue, error) {
	if len(name) == 0 || len(accessKey) == 0 {
		return nil, fmt.Errorf("Invalid Arg")
	}

	hashedKey, err := c.hash.Process(accessKey)
	if err != nil {
		return nil, err
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	if result, found := c.db.Queues[name]; found {
		if hashedKey == result.AccessKey {
			return result, nil
		}
		return nil, fmt.Errorf("Unauthorized")
	}

	return nil, nil
}

// DeleteQueue removes the given queue by name if the access token is correct
func (c *QueryControl) DeleteQueue(name, accessKey string) error {
	if len(name) == 0 || len(accessKey) == 0 {
		return fmt.Errorf("Invalid Arg")
	}

	hashedKey, err := c.hash.Process(accessKey)
	if err != nil {
		return err
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	queue, found := c.db.Queues[name]
	if !found {
		return fmt.Errorf("Not Found")
	} else if queue.AccessKey != hashedKey {
		return fmt.Errorf("Unauthorized")
	}

	delete(c.db.Queues, name)
	return nil
}
