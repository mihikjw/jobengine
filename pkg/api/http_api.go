package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MichaelWittgreffe/jobengine/pkg/database"
	"github.com/MichaelWittgreffe/jobengine/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

// HTTPAPI is an object for an HTTP/1.1 API for controlling the application
type HTTPAPI struct {
	router  *chi.Mux
	logger  logger.Logger
	monitor database.DBMonitor
	control database.QueryController
	json    *database.JSONDataHandler
}

// NewHTTPAPI is a constructor for an HttpAPI object
func NewHTTPAPI(logger logger.Logger, monitor database.DBMonitor, controller database.QueryController) *HTTPAPI {
	if logger == nil || monitor == nil || controller == nil {
		return nil
	}

	api := &HTTPAPI{
		logger:  logger,
		monitor: monitor,
		control: controller,
		json:    new(database.JSONDataHandler),
	}

	api.router = chi.NewRouter()
	api.router.Use(middleware.RequestID)
	api.router.Use(middleware.RealIP)
	api.router.Use(middleware.Logger)
	api.router.Use(middleware.Recoverer)
	api.router.Use(middleware.Timeout(10 * time.Second))

	api.router.Get("/test", api.Test)

	api.router.Put("/api/v1/queue", api.CreateQueue)
	api.router.Get("/api/v1/queue", api.GetQueue)
	api.router.Delete("/api/v1/queue", api.DeleteQueue)

	api.router.Put("/api/v1/job", api.AddJob)
	api.router.Get("/api/v1/job", api.GetJob)
	api.router.Get("/api/v1/job/next", api.GetNextJob)
	api.router.Post("/api/v1/job", api.UpdateJobStatus)
	api.router.Delete("/api/v1/job", api.DeleteJob)

	return api
}

// ListenAndServe starts the API listening for requets, blocks current goroutine
func (a *HTTPAPI) ListenAndServe(port string) error {
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), a.router)
}

// Test is a simple 'is-alive' endpoint handler, returning a 200 status code
func (a *HTTPAPI) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// CreateQueue is an endpoint handler for API requests to create a queue
func (a *HTTPAPI) CreateQueue(w http.ResponseWriter, r *http.Request) {
	body := new(CreateQueueRequest)
	if err := getRequestBody(body, r, a.json); err != nil {
		returnStatusCode(http.StatusBadRequest, w)
		return
	}

	if err := a.control.CreateQueue(body.Name, body.AccessKey); err != nil {
		errStr := err.Error()
		switch {
		case errStr == "Invalid Arg":
			returnStatusCode(http.StatusBadRequest, w)
		case errStr == "Queue Exists":
			returnStatusCode(http.StatusConflict, w)
		default:
			returnInternalServerError(err, w, a.json)
		}
		return
	}

	a.monitor.Write()
	returnStatusCode(http.StatusCreated, w)
}

// GetQueue is an endpoint handler for API requests to return a copy of a queue
func (a *HTTPAPI) GetQueue(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Query().Get("name")
	accessKey := r.Header.Get("X-Access-Key")
	if len(queueName) == 0 || len(accessKey) == 0 {
		returnStatusCode(http.StatusBadRequest, w)
		return
	}

	// regardless whether the user has access, we should use this time to update the queue
	if !updateQueue(queueName, a.control, w, a.json, a.monitor) {
		return
	}

	queue, err := a.control.GetQueue(queueName, accessKey)
	if err != nil {
		errStr := err.Error()
		switch {
		case errStr == "Invalid Args":
			returnStatusCode(http.StatusBadRequest, w)
		case errStr == "Unauthorized":
			returnStatusCode(http.StatusUnauthorized, w)
		default:
			returnInternalServerError(err, w, a.json)
		}
		return
	} else if queue == nil && err == nil {
		returnStatusCode(http.StatusNotFound, w)
		return
	}

	response := new(GetQueueResponse)
	response.Jobs = queue.Jobs
	response.Name = queue.Name
	response.Size = queue.Size

	if err = returnResponseBody(http.StatusOK, response, w, a.json); err != nil {
		returnInternalServerError(err, w, a.json)
	}
}

// DeleteQueue is an endpoint handler for API requests to delete a queue
func (a *HTTPAPI) DeleteQueue(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Query().Get("name")
	accessKey := r.Header.Get("X-Access-Key")
	if len(queueName) == 0 || len(accessKey) == 0 {
		returnStatusCode(http.StatusBadRequest, w)
		return
	}

	if err := a.control.DeleteQueue(queueName, accessKey); err != nil {
		errStr := err.Error()
		switch {
		case errStr == "Invalid Args":
			returnStatusCode(http.StatusBadRequest, w)
		case errStr == "Unauthorized":
			returnStatusCode(http.StatusUnauthorized, w)
		case errStr == "Not Found":
			returnStatusCode(http.StatusNotFound, w)
		default:
			returnInternalServerError(err, w, a.json)
		}
		return
	}

	a.monitor.Write()
	returnStatusCode(http.StatusNoContent, w)
}

// AddJob is an endpoint handler for adding a new job to a queue
func (a *HTTPAPI) AddJob(w http.ResponseWriter, r *http.Request) {
	accessKey := r.Header.Get("X-Access-Key")
	body := new(AddJobRequest)
	err := getRequestBody(body, r, a.json)
	if len(accessKey) == 0 || err != nil {
		returnStatusCode(http.StatusBadRequest, w)
		return
	}

	job := body.Job
	job.Created = time.Now().Unix()
	job.LastUpdated = job.Created
	job.State = database.Queued
	job.UID = uuid.New().String()

	if err = a.control.AddJob(job, body.QueueName, accessKey, false); err != nil {
		errStr := err.Error()
		switch {
		case errStr == "Invalid Args":
			returnStatusCode(http.StatusBadRequest, w)
		case errStr == "Unauthorized":
			returnStatusCode(http.StatusUnauthorized, w)
		case errStr == "Not Found":
			returnStatusCode(http.StatusNotFound, w)
		default:
			returnInternalServerError(err, w, a.json)
		}
		return
	}

	a.control.UpdateQueue(body.QueueName)
	a.monitor.Write()
	if err = returnResponseBody(http.StatusCreated, job, w, a.json); err != nil {
		returnInternalServerError(err, w, a.json)
	}
}

// GetJob is a handler for querying an entry for a specific job
func (a *HTTPAPI) GetJob(w http.ResponseWriter, r *http.Request) {
	accessKey := r.Header.Get("X-Access-Key")
	queueName := r.URL.Query().Get("queueName")
	uid := r.URL.Query().Get("jobUID")
	if len(uid) == 0 || len(queueName) == 0 || len(accessKey) == 0 {
		returnStatusCode(http.StatusBadRequest, w)
		return
	}

	// regardless whether the user has access, we should use this time to update the queue
	if !updateQueue(queueName, a.control, w, a.json, a.monitor) {
		return
	}

	job, err := a.control.GetJob(uid, queueName, accessKey)
	if err != nil {
		errStr := err.Error()
		switch {
		case errStr == "Invalid Args":
			returnStatusCode(http.StatusBadRequest, w)
		case errStr == "Unauthorized":
			returnStatusCode(http.StatusUnauthorized, w)
		default:
			returnInternalServerError(err, w, a.json)
		}
		return
	} else if job == nil {
		returnStatusCode(http.StatusNotFound, w)
		return
	}

	if err := returnResponseBody(http.StatusOK, job, w, a.json); err != nil {
		returnInternalServerError(err, w, a.json)
		return
	}
}

// GetNextJob is a handler for returning the next job from the queue head that is queued
func (a *HTTPAPI) GetNextJob(w http.ResponseWriter, r *http.Request) {
	accessKey := r.Header.Get("X-Access-Key")
	queueName := r.URL.Query().Get("queueName")
	if len(queueName) == 0 || len(accessKey) == 0 {
		returnStatusCode(http.StatusBadRequest, w)
		return
	}

	// regardless whether the user has access, we should use this time to update the queue
	if !updateQueue(queueName, a.control, w, a.json, a.monitor) {
		return
	}

	job, err := a.control.GetNextJob(queueName, accessKey)
	if err != nil {
		errStr := err.Error()
		switch {
		case errStr == "Invalid Args":
			returnStatusCode(http.StatusBadRequest, w)
		case errStr == "Unauthorized":
			returnStatusCode(http.StatusUnauthorized, w)
		default:
			returnInternalServerError(err, w, a.json)
		}
		return
	} else if job == nil {
		returnStatusCode(http.StatusNoContent, w)
		return
	}

	// update the job status if the flag is set, default don't update
	markQueued := r.URL.Query().Get("markQueued")
	if len(markQueued) > 0 && strings.ToLower(markQueued) == "true" {
		if err := a.control.UpdateJobStatus(job.UID, database.Inprogress, queueName, accessKey); err != nil {
			returnInternalServerError(err, w, a.json)
			return
		}
		job.State = database.Inprogress
		a.monitor.Write()
	}

	if err := returnResponseBody(http.StatusOK, job, w, a.json); err != nil {
		returnInternalServerError(err, w, a.json)
		return
	}

}

// UpdateJobStatus is a handler for updating the status of a given job
func (a *HTTPAPI) UpdateJobStatus(w http.ResponseWriter, r *http.Request) {
	accessKey := r.Header.Get("X-Access-Key")
	if len(accessKey) == 0 {
		returnStatusCode(http.StatusBadRequest, w)
		return
	}

	body := new(UpdateJobStatusRequest)
	if err := getRequestBody(body, r, a.json); err != nil {
		returnStatusCode(http.StatusBadRequest, w)
		return
	}

	// regardless whether the user has access, we should use this time to update the queue
	if !updateQueue(body.QueueName, a.control, w, a.json, a.monitor) {
		return
	}

	if err := a.control.UpdateJobStatus(body.UID, strings.ToLower(body.NewStatus), body.QueueName, accessKey); err != nil {
		errStr := err.Error()
		switch {
		case errStr == "Invalid Args":
			returnStatusCode(http.StatusBadRequest, w)
		case errStr == "Unauthorized":
			returnStatusCode(http.StatusUnauthorized, w)
		case errStr == "Not Found":
			returnStatusCode(http.StatusNotFound, w)
		default:
			returnInternalServerError(err, w, a.json)
		}
		return
	}

	a.monitor.Write()
	returnStatusCode(http.StatusOK, w)
}

// DeleteJob is a handler for deleting a job entry
func (a *HTTPAPI) DeleteJob(w http.ResponseWriter, r *http.Request) {
	accessKey := r.Header.Get("X-Access-Key")
	queueName := r.URL.Query().Get("queueName")
	uid := r.URL.Query().Get("jobUID")
	if len(uid) == 0 || len(queueName) == 0 || len(accessKey) == 0 {
		returnStatusCode(http.StatusBadRequest, w)
		return
	}

	// regardless whether the user has access, we should use this time to update the queue
	if !updateQueue(queueName, a.control, w, a.json, a.monitor) {
		return
	}

	if err := a.control.DeleteJob(uid, queueName, accessKey); err != nil {
		errStr := err.Error()
		switch {
		case errStr == "Invalid Args":
			returnStatusCode(http.StatusBadRequest, w)
		case errStr == "Unauthorized":
			returnStatusCode(http.StatusUnauthorized, w)
		case errStr == "Not Found":
			returnStatusCode(http.StatusNotFound, w)
		default:
			returnInternalServerError(err, w, a.json)
		}
		return
	}

	a.monitor.Write()
	returnStatusCode(http.StatusNoContent, w)
}
