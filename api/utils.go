package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/MichaelWittgreffe/jobengine/database"
)

// getRequestBody marshals the incoming request body into the given object pointer
func getRequestBody(bodyObj interface{}, r *http.Request, json *database.JSONDataHandler) error {
	if mimeType := r.Header.Get("Content-Type"); mimeType != "application/json" {
		return fmt.Errorf("Unsupported Content-Type %s", mimeType)
	}

	buffer, err := ioutil.ReadAll(r.Body)
	if err == nil {
		return json.Decode(buffer, bodyObj)
	}
	return err
}

// sendResponseBody unmarshals the given object into the http response and sets the content type
func sendResponseBody(statusCode int, bodyObj interface{}, w http.ResponseWriter, json *database.JSONDataHandler) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	body, err := json.Encode(bodyObj)
	if err != nil {
		return err
	}

	size, err := w.Write(body)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Length", strconv.Itoa(size))
	return nil
}

// returnUnauthorized sets up the response writer with a resonse code
func returnStatusCode(code int, w http.ResponseWriter) {
	w.WriteHeader(code)
}

// returnInternalServerError sets up the response writer with an Internal Server Error response
func returnInternalServerError(err error, w http.ResponseWriter, json *database.JSONDataHandler) error {
	w.WriteHeader(http.StatusInternalServerError)
	if result, err := json.Encode(ErrorResponse{Err: err.Error()}); err == nil {
		w.Write(result)
	} else {
		return err
	}
	return nil
}

// updateQueue performs an update and file write on the queue, returns true on success, false on failure with error response already setup
func updateQueue(queueName string, c database.QueryController, w http.ResponseWriter, j *database.JSONDataHandler, m database.DBMonitor) bool {
	if err := c.UpdateQueue(queueName); err != nil {
		errStr := err.Error()
		switch {
		case errStr == "Invalid Arg":
			returnStatusCode(http.StatusBadRequest, w)
		case errStr == "Not Found":
			returnStatusCode(http.StatusNotFound, w)
		default:
			returnInternalServerError(err, w, j)
		}
		return false
	}

	m.Write()
	return true
}
