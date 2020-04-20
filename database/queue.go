package database

// Queue represents a configured queue
type Queue struct {
	Jobs      []*Job `json:"jobs"`
	AccessKey string `json:"access_key"`
	Size      int    `json:"size"`
	Name      string `json:"name"`
}
