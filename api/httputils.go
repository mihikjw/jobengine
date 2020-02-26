package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MichaelWittgreffe/jobengine/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//GetJSONBody ensures there is body content and it is application/json
func GetJSONBody(gc *gin.Context) (map[string]interface{}, error) {
	if gc.ContentType() != "application/json" {
		return nil, fmt.Errorf("Mime Type Is Not application/json")
	}
	if gc.Request.ContentLength <= 0 {
		return nil, fmt.Errorf("Content Length Is 0")
	}

	requestBody := make(map[string]interface{})
	err := json.NewDecoder(gc.Request.Body).Decode(&requestBody)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON Body")
	}

	return requestBody, nil
}

//CreateJobFromBody creates a models.Job object from a JSON request body as map[string]interface{}
func CreateJobFromBody(in map[string]interface{}) (*models.Job, error) {
	var found bool
	job := new(models.Job)

	job.Content, found = in["content"].(map[string]interface{})
	if !found {
		return nil, fmt.Errorf("Missing JSON Field: content")
	}

	if tmp, found := in["priority"].(float64); found {
		job.Priority = uint8(tmp)
	} else {
		return nil, fmt.Errorf("Missing Or Invalid JSON Field: priority")
	}

	validFor, found := in["valid_for"].(float64)
	if !found || validFor <= 0 {
		return nil, fmt.Errorf("Missing Or Invalid JSON Field: valid_for")
	}

	job.Created = time.Now().Unix()
	job.LastUpdated = job.Created
	job.TimeoutTime = job.Created + int64(validFor)
	job.State = "queued"
	job.UID = uuid.New().String()
	return job, nil
}

//JobToMap transposes a models.Job object to a map[string]interface{}
func JobToMap(in *models.Job) map[string]interface{} {
	result := make(map[string]interface{}, 7)
	result["uid"] = in.UID
	result["content"] = in.Content
	result["state"] = in.State
	result["last_updated"] = in.LastUpdated
	result["created"] = in.Created
	result["timeout_time"] = in.TimeoutTime
	result["priority"] = in.Priority
	return result
}

//GetNameAndQueueFromContext gets the appName and queueName from an incoming request
func GetNameAndQueueFromContext(gc *gin.Context) (string, string, error) {
	appName := gc.GetHeader("X-Name")
	if len(appName) <= 0 {
		return "", "", fmt.Errorf("Missing Header Field: X-Name")
	}

	queueName := gc.GetHeader("X-Queue")
	if len(queueName) <= 0 {
		return "", "", fmt.Errorf("Missing Header Field: X-Queue")
	}

	return appName, queueName, nil
}
