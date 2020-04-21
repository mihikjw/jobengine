package database

// Job represents a job within the database
type Job struct {
	Priority       int                    `json:"priority"`
	KeepMinutes    int64                  `json:"keep_minutes"`
	TimeoutMinutes int64                  `json:"timeout_minutes"`
	LastUpdated    int64                  `json:"last_updated"`
	Created        int64                  `json:"created"`
	TimeoutTime    int64                  `json:"timeout_time"`
	UID            string                 `json:"uid"`
	Content        map[string]interface{} `json:"content"`
	State          string                 `json:"state"`
}
