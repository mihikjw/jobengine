package api

import "github.com/MichaelWittgreffe/jobengine/database"

// ErrorResponse can be called from any handler, to notify of an error in the request/processing
type ErrorResponse struct {
	Err string `json:"error"`
}

// GetQueueResponse is a response object for the Get Queue endpoint
type GetQueueResponse struct {
	Jobs []*database.Job `json:"jobs"`
	Size int             `json:"size"`
	Name string          `json:"name"`
}
