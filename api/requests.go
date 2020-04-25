package api

import "github.com/MichaelWittgreffe/jobengine/database"

// CreateQueueRequest represents the request body for the create queue endpoint
type CreateQueueRequest struct {
	Name      string `json:"name"`
	AccessKey string `json:"access_key"`
}

// AddJobRequest represents the request body for the add job endpoint
type AddJobRequest struct {
	Job       *database.Job `json:"job"`
	QueueName string        `json:"queue_name"`
}

// UpdateJobStatusRequest represents the request body for the update job endpoint
type UpdateJobStatusRequest struct {
	QueueName string `json:"queue_name"`
	UID       string `json:"uid"`
	NewStatus string `json:"new_status"`
}
