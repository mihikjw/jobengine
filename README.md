# JobEngine
JobEngine is a job queuing system allowing multiple applications to dynamically create and read 'job queues', whilst being able to mark jobs as 'queued', 'inprogress', 'failed' or 'complete', as well as give them a timeout window (process within this time, else delete). The queue is persistent across application restarts and provides an access system, with queues only accessible with a key. Jobs are persisted in an AES encrypted database file. Changes to jobs/queues are written to the database when a request is made to do so, this may include sorting the queue and writing the changes when a user makes a request. API docs are avalible [here](./api_schema.yml).

The application is distributed with a dockerfile/docker-compose.yml, this is the primary supported way of running the application. You'll be able to get an instance running by simply execting `docker-compose up` at the CLI from the root of the repository. If you're new to Docker, I've written an [introduction document with an example project](https://github.com/MichaelWittgreffe/DockerDemo).

## Build
1. ```go get https://github.com/MichaelWittgreffe/jobengine```
2. ```cd $GOPATH/github.com/MichaelWittgreffe/jobengine```
3. ```make```

## To Do
- Unit Tests!