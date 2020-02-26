package api

import (
	"fmt"
	"net/http"

	"github.com/MichaelWittgreffe/jobengine/models"
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
	public.GET("/test", result.test)             //test endpoint
	public.PUT("/jobs/create", result.createJob) //create job
	public.GET("/jobs/next", result.getNextjob)  //get next queued job
	public.GET("/jobs", result.getJobsAtStatus)  //get all jobs at status
	public.POST("/jobs/:uid", result.updateJob)  //update job
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
	appName, queueName, err := GetNameAndQueueFromContext(gc)
	if err != nil {
		gc.AbortWithError(http.StatusBadRequest, err)
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

	s.write <- true

	result := map[string]interface{}{
		"uid":          job.UID,
		"content":      job.Content,
		"state":        job.State,
		"last_updated": job.LastUpdated,
		"created":      job.Created,
		"timeout_time": job.TimeoutTime,
		"priority":     job.Priority,
	}

	gc.JSON(http.StatusCreated, result)
}

//getNextJob is a handler for requests to /api/v1/jobs/next, returns the next job in the queue and marks as 'Inprogress'
func (s *HTTPServer) getNextjob(gc *gin.Context) {
	appName, queueName, err := GetNameAndQueueFromContext(gc)
	if err != nil {
		gc.AbortWithError(http.StatusBadRequest, err)
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

	if err := s.controller.UpdateQueue(queueName); err != nil {
		gc.AbortWithError(http.StatusInternalServerError, err)
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

//getJobsAtStatus is a handler for requests to /api/v1/jobs, returns all the jobs in the queue, optionally at the given status
func (s *HTTPServer) getJobsAtStatus(gc *gin.Context) {
	appName, queueName, err := GetNameAndQueueFromContext(gc)
	if err != nil {
		gc.AbortWithError(http.StatusBadRequest, err)
	}

	statusFilter := gc.GetHeader("X-Status-Filter")

	queueFound, readAllowed, _ := s.controller.QueueExists(queueName, appName)
	if !queueFound {
		gc.AbortWithError(http.StatusNotFound, fmt.Errorf("Queue %s Not Found", queueName))
		return
	}
	if !readAllowed {
		gc.AbortWithError(http.StatusForbidden, fmt.Errorf("Permission Denied For Read From Queue %s", queueName))
		return
	}

	if err := s.controller.UpdateQueue(queueName); err != nil {
		gc.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	//write the changes in UpdateQueue
	s.write <- true

	jobs, err := s.controller.ExportQueue(queueName, statusFilter)
	if err != nil {
		gc.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Error Getting Jobs: %s", err.Error()))
		return
	} else if jobs == nil {
		gc.AbortWithError(http.StatusNotFound, fmt.Errorf("Queue %s Not Found", queueName))
		return
	}

	gc.JSON(http.StatusOK, jobs["jobs"])
}

func (s *HTTPServer) updateJob(gc *gin.Context) {
	jobUID := gc.Param("uid")
	if len(jobUID) <= 0 {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("No UID Supplied"))
		return
	}

	appName, queueName, err := GetNameAndQueueFromContext(gc)
	if err != nil {
		gc.AbortWithError(http.StatusBadRequest, err)
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

	requestBody, err := GetJSONBody(gc)
	if err != nil {
		gc.AbortWithError(http.StatusBadRequest, err)
		return
	}

	content, found := requestBody["content"].(map[string]interface{})
	if !found {
		content = nil
	}

	var timeoutTime int64 = 0
	tempTimeoutTime, found := requestBody["timeout_time"].(int)
	if found {
		timeoutTime = int64(tempTimeoutTime)
	}

	var priority uint8 = 0
	tempPriority, found := requestBody["priority"].(int)
	if found {
		priority = uint8(tempPriority)
	}

	state, found := requestBody["status"].(string)
	if found && (state != models.Complete && state != models.Failed && state != models.Inprogress && state != models.Queued) {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid Status Requested: %s", state))
		return
	}

	err = s.controller.UpdateJob(queueName, jobUID, state, content, timeoutTime, priority)
	if err != nil {
		if err.Error() == "Not Found" {
			gc.AbortWithError(http.StatusNotFound, fmt.Errorf("Job %s Not Found In Queue %s", jobUID, queueName))
			return
		}
		gc.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Error Updating Job %s For Queue %s: %s", jobUID, queueName, err.Error()))
		return
	}

	s.write <- true
	gc.Status(http.StatusOK)
}
