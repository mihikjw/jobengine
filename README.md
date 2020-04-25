# JobEngine

JobEngine is a job queuing system allowing multiple applications to dynamically create and read 'job queues', whilst being able to mark jobs as 'queued', 'inprogress', 'failed' or 'complete', as well as give them a timeout window (process within this time, else delete). The queue is persistent across application restarts and provides an access system, with queues only accessible with a key. Jobs are persisted in an AES encrypted database file. Changes to jobs/queues are written to the database when a request is made to do so, this may include sorting the queue  and writing the changes when a user makes a request.

## Build
1. ```go get https://github.com/MichaelWittgreffe/jobengine```
2. ```cd $GOPATH/github.com/MichaelWittgreffe/jobengine```
3. ```make```