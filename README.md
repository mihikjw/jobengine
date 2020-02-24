# JobEngine

JobEngine is a 'job-queue', a queue system allowing multiple applications to dynamically create and read 'job queues', whilst crucially being able to mark jobs as 'in-progress', 'failed' or 'complete', as well as give them a timeout window (process within this time, else delete). The queue is also persistent across application restarts and provides a loose permissions system, only allowing certian functionality to certian application names (e.g. only appA can push to queueX & only appB can read)

## Build
1. ```go get https://github.com/MichaelWittgreffe/jobengine```
2. ```cd $GOPATH/github.com/MichaelWittgreffe/jobengine```
3. ```make```

## Requirements
- Internal queue for jobs
    - Job object
        - internal details
        - user supplied JSON
        - timeout flag
    - API (create/get/mark_status)
        - jobs
            - create
            - read
            - update status
            - delete
        - queues
            - create
            - read
            - delete
- Queue should be written to disk regularly
    - file should be encrypted
- Permissions based on application name
- Configuration
    - From file
    - Encryption Key from env
    - Default (then write to file) 