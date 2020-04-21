package api

// CreateQueueRequest represents the request body for the create queue endpoint
type CreateQueueRequest struct {
	Name      string `json:"name"`
	AccessKey string `json:"access_key"`
}
