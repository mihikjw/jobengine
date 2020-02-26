package api

import (
	"fmt"
	"net/http"

	"github.com/MichaelWittgreffe/jobengine/logger"
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
	log        logger.Logger
}

//NewHTTPServer creates a new instance of HTTPServer
func NewHTTPServer(con *queue.Controller, write chan bool, logger logger.Logger) *HTTPServer {
	result := &HTTPServer{
		controller: con,
		write:      write,
		router:     gin.Default(),
		json:       new(queue.JSONHandle),
		log:        logger,
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
		s.log.Error(err.Error())
		gc.Status(http.StatusBadRequest)
		return
	}

	requestBody, err := GetJSONBody(gc)
	if err != nil {
		s.log.Error(err.Error())
		gc.Status(http.StatusBadRequest)
		return
	}

	queueFound, _, writeAllowed := s.controller.QueueExists(queueName, appName)
	if !queueFound {
		s.log.Error(fmt.Sprintf("Queue %s Not Found", queueName))
		gc.Status(http.StatusNotFound)
		return
	}
	if !writeAllowed {
		s.log.Error(fmt.Sprintf("Permission Denied For Write To Queue %s", queueName))
		gc.Status(http.StatusNotFound)
		return
	}

	job, err := CreateJobFromBody(requestBody)
	if err != nil {
		s.log.Error(err.Error())
		gc.Status(http.StatusBadRequest)
		return
	}

	err = s.controller.AddNewJob(queueName, job)
	if err != nil {
		s.log.Error(fmt.Sprintf("Error Adding Job: %s", err.Error()))
		gc.Status(http.StatusInternalServerError)
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
		s.log.Error(err.Error())
		gc.Status(http.StatusBadRequest)
		return
	}

	queueFound, readAllowed, _ := s.controller.QueueExists(queueName, appName)
	if !queueFound {
		s.log.Error(fmt.Sprintf("Queue %s Not Found", queueName))
		gc.Status(http.StatusNotFound)
		return
	}
	if !readAllowed {
		s.log.Error(fmt.Sprintf("Permission Denied For Read From Queue %s", queueName))
		gc.Status(http.StatusForbidden)
		return
	}

	if err := s.controller.UpdateQueue(queueName); err != nil {
		s.log.Error(err.Error())
		gc.Status(http.StatusInternalServerError)
		return
	}

	job, err := s.controller.GetNextJob(queueName)
	s.write <- true

	if err != nil {
		s.log.Error(fmt.Sprintf("Error Getting Next Job: %s", err.Error()))
		gc.Status(http.StatusInternalServerError)
		return
	} else if job == nil {
		s.log.Error(fmt.Sprintf("Queue %s Not Found", queueName))
		gc.Status(http.StatusNotFound)
		return
	}

	jobMap := JobToMap(job)
	gc.JSON(http.StatusOK, jobMap)
}

//getJobsAtStatus is a handler for requests to /api/v1/jobs, returns all the jobs in the queue, optionally at the given status
func (s *HTTPServer) getJobsAtStatus(gc *gin.Context) {
	appName, queueName, err := GetNameAndQueueFromContext(gc)
	if err != nil {
		s.log.Error(err.Error())
		gc.Status(http.StatusBadRequest)
		return
	}

	statusFilter := gc.GetHeader("X-Status-Filter")

	queueFound, readAllowed, _ := s.controller.QueueExists(queueName, appName)
	if !queueFound {
		s.log.Error(fmt.Sprintf("Queue %s Not Found", queueName))
		gc.Status(http.StatusNotFound)
		return
	}
	if !readAllowed {
		s.log.Error(fmt.Sprintf("Permission Denied For Read From Queue %s", queueName))
		gc.Status(http.StatusForbidden)
		return
	}

	if err := s.controller.UpdateQueue(queueName); err != nil {
		if err != nil {
			s.log.Error(err.Error())
			gc.Status(http.StatusInternalServerError)
			return
		}
	}

	//write the changes in UpdateQueue
	s.write <- true

	jobs, err := s.controller.ExportQueue(queueName, statusFilter)
	if err != nil {
		s.log.Error(fmt.Sprintf("Error Getting Jobs: %s", err.Error()))
		gc.Status(http.StatusInternalServerError)
		return
	} else if jobs == nil {
		s.log.Error(fmt.Sprintf("Queue %s Not Found", queueName))
		gc.Status(http.StatusNotFound)
		return
	}

	gc.JSON(http.StatusOK, jobs["jobs"])
}

func (s *HTTPServer) updateJob(gc *gin.Context) {
	jobUID := gc.Param("uid")
	if len(jobUID) <= 0 {
		s.log.Error("No UID Supplied")
		gc.Status(http.StatusBadRequest)
		return
	}

	appName, queueName, err := GetNameAndQueueFromContext(gc)
	if err != nil {
		s.log.Error(err.Error())
		gc.Status(http.StatusBadRequest)
		return
	}

	queueFound, _, writeAllowed := s.controller.QueueExists(queueName, appName)
	if !queueFound {
		s.log.Error(fmt.Sprintf("Queue %s Not Found", queueName))
		gc.Status(http.StatusNotFound)
		return
	}
	if !writeAllowed {
		s.log.Error(fmt.Sprintf("Permission Denied For Write To Queue %s", queueName))
		gc.Status(http.StatusForbidden)
		return
	}

	requestBody, err := GetJSONBody(gc)
	if err != nil {
		s.log.Error(err.Error())
		gc.Status(http.StatusBadRequest)
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
		s.log.Error(fmt.Sprintf("Invalid Status Requested: %s", state))
		gc.Status(http.StatusBadRequest)
		return
	}

	err = s.controller.UpdateJob(queueName, jobUID, state, content, timeoutTime, priority)
	if err != nil {
		if err.Error() == "Not Found" {
			s.log.Error(fmt.Sprintf("Job %s Not Found In Queue %s", jobUID, queueName))
			gc.Status(http.StatusNotFound)
			return
		}
		s.log.Error(fmt.Sprintf("Error Updating Job %s For Queue %s: %s", jobUID, queueName, err.Error()))
		gc.Status(http.StatusInternalServerError)
		return
	}

	s.write <- true
	gc.Status(http.StatusOK)
}
