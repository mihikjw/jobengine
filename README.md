# JobEngine

JobEngine is a job queuing system allowing multiple applications to dynamically create and read 'job queues', whilst being able to mark jobs as 'inprogress', 'failed' or 'complete', as well as give them a timeout window (process within this time, else delete). The queue is also persistent across application restarts and provides a loose permissions system, only allowing certain functionality to certain application names (e.g. only appA can push to queueX & only appB can read).

Jobs are persisted in an encrypted database file, this is encrypted with AES. Changes to jobs/queues are written to the database when a request is made to do so, this may include sorting the queue for requests that change the status of a job (such as GetNextJob). Further details are in the API documentation.

## Build
1. ```go get https://github.com/MichaelWittgreffe/jobengine```
2. ```cd $GOPATH/github.com/MichaelWittgreffe/jobengine```
3. ```make```

## Requirements
- Env Data
    - apiPort
    - encryptionSecret
    - dbPath
- Job
    - jobKeepMinutes
    - jobTimeoutMinutes
    - uid
    - content
    - state
    - lastUpdated
    - created
    - timeoutTime
    - priority
- Queue
    - jobs (sorted)
    - password-protected
    - size
    - name
- Database File
    - encrypted
        - require a secret
    - stores queues with jobs
    - write-on-change in a seperate goroutine
- API
    - Queues
        - Create
        - Read
        - Delete
    - Test endpoint
    - Jobs
        - Add
        - GetJob
        - GetNextJob
        - GetAllJobs
        - UpdateJobStatus
        - DeleteJob