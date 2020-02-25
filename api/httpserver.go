package api

import (
	"fmt"
	"net/http"

	"github.com/MichaelWittgreffe/jobengine/queue"
	"github.com/gin-gonic/gin"
)

//HTTPServer represents an HTTP/1.1 server for the API
type HTTPServer struct {
	controller *queue.Controller
	write      chan bool
	router     *gin.Engine
	json       queue.JSONHandler
}

//NewHTTPServer creates a new instance of HTTPServer
func NewHTTPServer(con *queue.Controller, write chan bool) *HTTPServer {
	result := &HTTPServer{
		controller: con,
		write:      write,
		router:     gin.Default(),
		json:       new(queue.JSONHandle),
	}

	public := result.router.Group("/api/v1")
	public.GET("/test", result.test)
	public.PUT("/jobs/create", result.createJob)
	public.GET("jobs/next", result.getNextjob)

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

//createJob is a handler for requests to /api/v1/jobs/create, creates a new job
func (s *HTTPServer) createJob(gc *gin.Context) {
	appName := gc.GetHeader("X-Name")
	if len(appName) <= 0 {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Missing Header Field: X-Name"))
		return
	}

	queueName := gc.GetHeader("X-Queue")
	if len(queueName) <= 0 {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Missing Header Field: X-Queue"))
		return
	}

	requestBody, err := GetJSONBody(gc)
	if err != nil {
		gc.AbortWithError(http.StatusBadRequest, err)
		return
	}

	queueFound, _, writeAllowed := s.controller.QueueExists(queueName, appName)
	if !queueFound {
		gc.AbortWithError(http.StatusNotFound, fmt.Errorf("Queue %s Not Found", queueName))
		return
	}
	if !writeAllowed {
		gc.AbortWithError(http.StatusForbidden, fmt.Errorf("Permission Denied For Write To Queue %s", queueName))
		return
	}

	job, err := CreateJobFromBody(requestBody)
	if err != nil {
		gc.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = s.controller.AddNewJob(queueName, job)
	if err != nil {
		gc.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Error Adding Job: %s", err.Error()))
		return
	}

	//request queue change write to disk
	s.write <- true

	gc.Status(http.StatusCreated)
	return
}

//getNextJob is a handler for requests to /api/v1/jobs/next, returns the next job in the queue and marks as 'Inprogress'
func (s *HTTPServer) getNextjob(gc *gin.Context) {
	appName := gc.GetHeader("X-Name")
	if len(appName) <= 0 {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Missing Header Field: X-Name"))
		return
	}

	queueName := gc.GetHeader("X-Queue")
	if len(queueName) <= 0 {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Missing Header Field: X-Queue"))
		return
	}

	queueFound, readAllowed, _ := s.controller.QueueExists(queueName, appName)
	if !queueFound {
		gc.AbortWithError(http.StatusNotFound, fmt.Errorf("Queue %s Not Found", queueName))
		return
	}
	if !readAllowed {
		gc.AbortWithError(http.StatusForbidden, fmt.Errorf("Permission Denied For Read From Queue %s", queueName))
		return
	}

	job, err := s.controller.GetNextJob(queueName)
	s.write <- true

	if err != nil {
		gc.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Error Getting Next Job: %s", err.Error()))
		return
	} else if job == nil {
		gc.AbortWithError(http.StatusNotFound, fmt.Errorf("Queue %s Not Found", queueName))
		return
	}

	jobMap := JobToMap(job)
	gc.JSON(http.StatusOK, jobMap)
}
