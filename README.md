# JobEngine
JobEngine was born out of a need for a queuing process for 'jobs' within a backend system. Traditional 'message queue' software such as RabbitMQ or Apache Kafka did not have the features I required to appropriately manage the concept of a 'job' - I needed the messages to have a state, with the ability to also easily call 'next' on a queue and have some granularity of control over the individual 'jobs', i.e. if the job setup failed, it wouldn't always be removed from the queue by calling 'next'.

The concept of a 'job' within the JobEngine is simply a JSON object, which would contain parameters/implementation details for another process within a backend system to interpret as a request for work.

Queues are created dynamically through an HTTP/1.1 interface, Jobs are then added through this API with a concept of state (queued, inprogress, failed, complete). This is persisted across application restarts through an AES-encrypted database file. The queue being accessed by the user/process through the API is updated accordingly each time the user makes a request and written to the database file. The API provides the ability to call 'GetNextJob' which returns the next job in the queue at status 'queued', which is then optionally set to 'inprogress' upon successfully returning, or this can be resolved by the user/process with subsequent API request. Full API docs are available [here](./api_schema.yml).

JobEngine is distributed with a dockerfile/docker-compose.yml, this is the primary supported way of running the application. You will be able to get an instance running by simply executing `docker-compose up` at the CLI from the root of the repository. If you're new to Docker, I've written an [introduction document with an example project](https://github.com/MichaelWittgreffe/DockerDemo).

**Note:** This project does not yet have a full suite of tests, so I wouldn't recommend for production use just yet :)

## Build
1. ```go get https://github.com/MichaelWittgreffe/jobengine```
2. ```cd $GOPATH/github.com/MichaelWittgreffe/jobengine```
3. ```make```

## To Do
- Unit Tests!
- Web UI