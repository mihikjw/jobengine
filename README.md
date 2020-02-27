# JobEngine

JobEngine is a 'job-queue', a job queuing system allowing multiple applications to dynamically create and read 'job queues', whilst crucially being able to mark jobs as 'inprogress', 'failed' or 'complete', as well as give them a timeout window (process within this time, else delete). The queue is also persistent across application restarts and provides a loose permissions system, only allowing certian functionality to certian application names (e.g. only appA can push to queueX & only appB can read).

## Build
1. ```go get https://github.com/MichaelWittgreffe/jobengine```
2. ```cd $GOPATH/github.com/MichaelWittgreffe/jobengine```
3. ```make```

## TO DO
### Short Term
- Add Job Delete endpoint
- Add Create, Read, Delete to queues
    - Add goroutine to monitor config updates when a queue is added/deleted/modified

### Longer Term
- Login & authentication
- Basic Web UI
    - View Queues, Jobs & Metrics