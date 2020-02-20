package models

//Job represents a task to be executed
type Job struct {
	UID         string
	Content     string
	State       string
	LastUpdated int64
	Created     int64
	TimeoutTime int64
}
