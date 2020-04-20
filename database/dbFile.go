package database

import "sync"

// DBFile represents an entire database file
type DBFile struct {
	lock   *sync.Mutex
	Queues map[string]*Queue `json:"queues"`
}

// NewDBFile is a constructor for DBFile
func NewDBFile() *DBFile {
	return &DBFile{
		Queues: make(map[string]*Queue),
		lock:   new(sync.Mutex),
	}
}
