# JobEngine

JobEngine is a 'job-queue', a job queuing system allowing multiple applications to dynamically create and read 'job queues', whilst crucially being able to mark jobs as 'inprogress', 'failed' or 'complete', as well as give them a timeout window (process within this time, else delete). The queue is also persistent across application restarts and provides a loose permissions system, only allowing certian functionality to certian application names (e.g. only appA can push to queueX & only appB can read).

## Build
1. ```go get https://github.com/MichaelWittgreffe/jobengine```
2. ```cd $GOPATH/github.com/MichaelWittgreffe/jobengine```
3. ```make```

## API
### Test: GET /api/v1/test
Simple 'is-alive' endpoint.

### CreateJob: PUT /api/v1/jobs/create
#### Headers
- X-Content-Type: application/json
- X-Name: service name defined in config
- X-Queue: queue name defined in config

#### Content
``` json
{
    "content": {
        "job": "content"
    },
    "valid_for": 600,
    "priority": 80
}
```
- content: the job content, variable and can be any user-defined JSON object
- valid_for: job timeout in seconds if it is waiting at a 'queued' state
- priority: priority of the job (>= 1; <= 100) relative to other jobs in the queue

#### Result
##### 201 CREATED
``` json
{
    "content": {
        "job": "content"
    },
    "created": 1582796995,
    "last_updated": 1582796995,
    "priority": 80,
    "state": "queued",
    "timeout_time": 1582797595,
    "uid": "2024ab76-a16f-4cd4-a494-5019174aa265"
}
```
- content: the job content sent in the request
- created: unix epoch time integer when job was created
- last_updated: unix epoch time integer when job was last updated; you can assert created == last_updated in this response
- priority: the job priority sent in the request
- state: state the job was created at; you can assert state == "queued" in this response
- timeout_time: unix epoch time integer for when the job will be considered 'timed out' if it remains at a queued status
- uid: unique identification string for the job; you can assert this is unique

### GetNextJob: GET /api/v1/jobs/next
#### Headers
- X-Name: service name defined in config
- X-Queue: queue name defined in config

#### Response
##### 200 OK
``` json
{
    "content": {
        "job": "content"
    },
    "created": 1582796995,
    "last_updated": 1582797343,
    "priority": 80,
    "state": "inprogress",
    "timeout_time": 1582797595,
    "uid": "2024ab76-a16f-4cd4-a494-5019174aa265"
}
```
- content: the job content sent in the original create request
- created: unix epoch time integer when job was created
- last_updated: unix epoch time integer when the job was last updated; this will be changed when calling this endpoint; you can assert last_updated > created in this response
- priority: the job priority relative to other jobs in the queue
- state: current state of the job, any returned job from this request will be changed to 'inprogress'; you can assert state == "inprogress" in this response
- timeout_time: unix epoch time integer for when the job would have been considered 'timed out' - no longer relevant as the job is now "inprogress"
- uid: unique identification string for the job; you can assert this is unique

##### 201 NO CONTENT
There is no jobs at status 'queued' in the queue

### GetAllJobs: GET /api/v1/jobs
#### Headers
- X-Name: service name defined in config
- X-Queue: queue name defined in config
- *OPTIONAL:* X-Status-Filter:
    - queued
    - inprogress
    - complete
    - failed

#### Response
##### 200 OK
``` json
{
    "0": {
        "content": {
            "job": "content"
        },
        "created": 1582796995,
        "last_updated": 1582797343,
        "priority": 80,
        "state": "inprogress",
        "timeout_time": 1582797595,
        "uid": "2024ab76-a16f-4cd4-a494-5019174aa265"
    },
    "1": {
        "content": {
            "job": "content"
        },
        "created": 1582796995,
        "last_updated": 1582797343,
        "priority": 80,
        "state": "inprogress",
        "timeout_time": 1582797595,
        "uid": "2024ab76-a16f-4cd4-a494-5019174aa265"
    }
}
```
- index: integer stored as a string of the position of a job within the queue
    - content: the job content sent in the original create request
    - created: unix epoch time integer when job was created
    - last_updated: unix epoch time integer when the job was last updated
    - priority: the job priority relative to other jobs in the queue
    - state: current state of the job; if the header X-Status-Filter was set, you can assert state == X-Status-Filter
    - timeout_time: unix epoch time integer for when the job is considered 'timed out' if it is at status 'queued'
    - uid: unique identification string for the job; you can assert this is unique

### UpdateJob: POST /api/v1/jobs/:uid
#### URL Parameters
- uid: the unique identifier for the job (uid generated in CreateJob)

#### Headers
- Content-Type: application/json
- X-Name: service name defined in config
- X-Queue: queue name defined in config

#### Request
``` json
{
    "status": "complete",
    "content": {
        "job": "content"
    },
    "timeout_time": 1582797595,
    "priority": 90
}
```
- *OPTIONAL:* status: the new status for the job (queued, inprogress, complete, failed - all others return BAD REQUEST)
- *OPTIONAL:* content: new content for the job, overwrites and previously held content
- *OPTIONAL:* timeout_time: new unix epoch time integer for the job to be considered 'timed out' if the status is 'queued'
- *OPTIONAL:* priority: the new priority for the job relative to other jobs in the queue (must be >=1 && <=100)

**NOTE:** the last_updated time of the job will also be set to the unix epoch time of the request.

## TO DO
### Short Term
- Add Job Delete endpoint
- Add Create, Read, Delete to queues
    - Add goroutine to monitor config updates when a queue is added/deleted/modified

### Longer Term
- Login & authentication
- Basic Web UI
    - View Queues, Jobs & Metrics