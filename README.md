# JobEngine

JobEngine is a 'job-queue', a queue system allowing multiple applications to dynamically create and read 'job queues', whilst crucially being able to mark jobs as 'in-progress', 'failed' or 'complete'. The queue is also persistent across application restarts.

## Build
1. ```go get https://github.com/MichaelWittgreffe/jobengine```
2. ```cd $GOPATH/github.com/MichaelWittgreffe/jobengine```
3. ```make```