package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MichaelWittgreffe/jobengine/database"
	"github.com/MichaelWittgreffe/jobengine/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	w.WriteHeader(http.StatusCreated)
}

// GetQueue is an endpoint handler for API requests to return a copy of a queue
func (a *HTTPAPI) GetQueue(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Query().Get("name")
	accessKey := r.Header.Get("X-Access-Key")
	if len(queueName) == 0 || len(accessKey) == 0 {
		returnStatusCode(http.StatusBadRequest, w)
		return
	}

	queue, err := a.control.GetQueue(queueName, accessKey)
	if err != nil {
		errStr := err.Error()
		switch {
		case errStr == "Invalid Arg":
			returnStatusCode(http.StatusBadRequest, w)
		case errStr == "Unauthorized":
			returnStatusCode(http.StatusUnauthorized, w)
		default:
			returnInternalServerError(err, w, a.json)
		}
	} else if queue == nil && err == nil {
		returnStatusCode(http.StatusNotFound, w)
	}

	response := new(GetQueueResponse)
	response.Jobs = queue.Jobs
	response.Name = queue.Name
	response.Size = queue.Size

	if err = sendResponseBody(http.StatusOK, response, w, a.json); err != nil {
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
		case errStr == "Invalid Arg":
			returnStatusCode(http.StatusBadRequest, w)
		case errStr == "Unauthorized":
			returnStatusCode(http.StatusUnauthorized, w)
		case errStr == "Not Found":
			returnStatusCode(http.StatusNotFound, w)
		default:
			returnInternalServerError(err, w, a.json)
		}
	}

	returnStatusCode(http.StatusNoContent, w)
}
