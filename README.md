# JobEngine

JobEngine is a 'job-queue', a job queuing system allowing multiple applications to dynamically create and read 'job queues', whilst crucially being able to mark jobs as 'inprogress', 'failed' or 'complete', as well as give them a timeout window (process within this time, else delete). The queue is also persistent across application restarts and provides a loose permissions system, only allowing certain functionality to certain application names (e.g. only appA can push to queueX & only appB can read).

Jobs are persisted in an encrypted database file, this is encrypted with AES. Changes to jobs/queues are written to the database when a request is made to do so, this may include sorting the queue for requests that change the status of a job (such as GetNextJob). Further details are in the API documentation below.

## Build
1. ```go get https://github.com/MichaelWittgreffe/jobengine```
2. ```cd $GOPATH/github.com/MichaelWittgreffe/jobengine```
3. ```make```

## Configuration
In the current version (0.0.1), all queue configuration is handled in a YAML config file. The default for this (and is generated if no config is supplied) is ```/jobengine/config.yml```, without any queues configured. The config file is used to define variables for the application including version, api port, job_keep_minutes, job_timeout_minutes and the queues. Please note that the config is king - if the queue configuration is changed, or a queue deleted, this is reflected in the database and this data is permanently lost. Example config:

``` yaml
version: 1
port: 6010
job_keep_minutes: 60
job_timeout_minutes: 10
queues:
  test_queue_1:
    read:
    - service1
    - service2
    write:
    - service3
  test_queue_2:
    read:
    - service3
    write:
    - service1
    - service2
```
- version: the API mode of the application
    - 1/1.0: HTTP/1.1
- port: the port used to host the API
- job_keep_minutes: the number of minutes to keep jobs that are complete/failed
- job_timeout_minutes: the number of minutes any job can be at status 'inprogress' before being marked as 'failed' and affected by the job_keep_minutes variable
- queues: define the applications queues
    - queue_name: name of the queue
        - read: array of arbitary service names that have read access to the queue
        - write: array of arbitary service names that have write access to the queue

Some variables are defined based on environment variables:
- SECRET: key used to encrypt the database
- GIN_MODE: 'release' or 'debug', modifies the logging level of the Gin framework, used as part of the HTTP/1.1 API
- CONFIG_PATH: define a custom path for the config YAML file
- DB_PATH: define a custom path for the encrypted database file (*.queuedb)

The SECRET is used to encrypt/decrypt the database and as such if this is changed you will no-longer be able to use any existing database.

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

##### 204 NO CONTENT
No jobs exist at the current status, if this filter is not applied, no jobs exist.

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
- Unit Tests; Improve test coverage
- Change API to marshal/unmarshal directly from structs instead of manual parsing and asserting types
- Add Job Delete endpoint
- Add Create, Read, Delete to queues
    - Add goroutine to monitor config updates when a queue is added/deleted/modified

### Longer Term
- Login & authentication
- Basic Web UI
    - View Queues, Jobs & Metrics