package database

import "sync"

// DBFile represents an entire database file
type DBFile struct {
	Queues map[string]*Queue `json:"queues"`
	mutex  sync.Mutex
}

// NewDBFile is a constructor for DBFile
func NewDBFile() *DBFile {
	return &DBFile{
		Queues: make(map[string]*Queue),
		mutex:  sync.Mutex{},
	}
}
