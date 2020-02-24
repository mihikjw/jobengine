package api

import (
	"github.com/MichaelWittgreffe/jobengine/queue"
)

//Server defines an interfaces for the application API
type Server interface {
	ListenAndServe(port int) error
}

//NewAPIServer is a factory for creating APIServer derived structs
func NewAPIServer(version float64, con *queue.Controller, write chan bool) Server {
	switch version {
	case 1.0:
		return NewHTTPServer(con, write)
	default:
		return nil
	}
}
