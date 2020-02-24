package api

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/queue"
	"github.com/gin-gonic/gin"
)

//HTTPServer represents an HTTP/1.1 server for the API
type HTTPServer struct {
	controller *queue.Controller
	write      chan bool
	router     *gin.Engine
}

//NewHTTPServer creates a new instance of HTTPServer
func NewHTTPServer(con *queue.Controller, write chan bool) *HTTPServer {
	result := &HTTPServer{
		controller: con,
		write:      write,
		router:     gin.Default(),
	}

	public := result.router.Group("/api/v1")
	public.GET("/test", result.test)

	return result
}

//ListenAndServe begins the HTTPServer listening on the given port for requests
func (s *HTTPServer) ListenAndServe(port int) error {
	boundPort := fmt.Sprintf(":%d", port)
	return s.router.Run(boundPort)
}

//test is a handler for requests to /api/v1/test
func (s *HTTPServer) test(gc *gin.Context) {
	gc.Status(200)
}
