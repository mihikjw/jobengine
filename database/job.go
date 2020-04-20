package database

// Job represents a job within the database
type Job struct {
	KeepMinutes    int                    `json:"keep_minutes"`
	TimeoutMinutes int                    `json:"timeout_minutes"`
	UID            string                 `json:"uid"`
	Content        map[string]interface{} `json:"content"`
	State          string                 `json:"state"`
	LastUpdated    float64                `json:"last_updated"`
	Created        float64                `json:"created"`
	TimeoutTime    float64                `json:"timeout_time"`
	Priority       int                    `json:"priority"`
}
