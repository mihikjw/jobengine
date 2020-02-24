package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MichaelWittgreffe/jobengine/models"
	"github.com/MichaelWittgreffe/jobengine/queue"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	public.GET("/jobs/create", result.createJob)

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
	//validate incoming request
	if gc.ContentType() != "application/json" {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid Content Type"))
		return
	}
	if gc.Request.ContentLength <= 0 {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("No Body Recieved"))
		return
	}

	appName := gc.GetHeader("X-Name")
	if len(appName) <= 0 {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Missing Header Field: X-Name"))
		return
	}

	requestBody := make(map[string]interface{})
	err := json.NewDecoder(gc.Request.Body).Decode(&requestBody)
	if err != nil {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid JSON"))
		return
	}

	var found bool
	queueName, found := requestBody["queue"].(string)
	if !found {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Missing JSON Field: queue"))
		return
	}

	//check the request is allowed
	queueFound, _, writeAllowed := s.controller.QueueExists(queueName, appName)
	if !queueFound {
		gc.AbortWithError(http.StatusNotFound, fmt.Errorf("Queue %s Not Found", queueName))
		return
	}
	if !writeAllowed {
		gc.AbortWithError(http.StatusForbidden, fmt.Errorf("Permission Denied For Write To Queue %s", queueName))
		return
	}

	job := new(models.Job)
	job.Content, found = requestBody["content"].(map[string]interface{})
	if !found {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Missing JSON Field: content"))
		return
	}

	if tmp, found := requestBody["priority"].(float64); found {
		job.Priority = uint8(tmp)
	} else {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Missing Or Invalid JSON Field: priority"))
		return
	}

	validFor, found := requestBody["valid_for"].(float64)
	if !found || validFor <= 0 {
		gc.AbortWithError(http.StatusBadRequest, fmt.Errorf("Missing Or Invalid JSON Field: valid_for"))
		return
	}

	job.Created = time.Now().Unix()
	job.LastUpdated = job.Created
	job.TimeoutTime = job.Created + int64(validFor)
	job.State = "queued"
	job.UID = uuid.New().String()

	err = s.controller.AddNewJob(queueName, job)
	if err != nil {
		gc.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Error Adding Job: %s", err.Error()))
		return
	}

	s.write <- true

	gc.Status(http.StatusCreated)
	return
}
