# JobEngine

JobEngine is a 'job-queue', a job queuing system allowing multiple applications to dynamically create and read 'job queues', whilst crucially being able to mark jobs as 'inprogress', 'failed' or 'complete', as well as give them a timeout window (process within this time, else delete). The queue is also persistent across application restarts and provides a loose permissions system, only allowing certian functionality to certian application names (e.g. only appA can push to queueX & only appB can read).

## Build
1. ```go get https://github.com/MichaelWittgreffe/jobengine```
2. ```cd $GOPATH/github.com/MichaelWittgreffe/jobengine```
3. ```make```

## TO DO:
- Add Create, Read, Delete to queues
    - Add goroutine to monitor config updates

## Requirements
- Internal queue for jobs
    - Job object
        - internal details
        - user supplied JSON
        - timeout flag
    - API (create/get/mark_status)
        - jobs
            - create: yes
            - read
                - next: yes
                - all (optional: at status): yes
            - update status: yes
            - delete: 
        - queues
            - create: 
            - read: 
            - delete: 
- Queue should be written to disk regularly: yes
    - file should be encrypted: yes
- Permissions based on application name: yes
- Configuration: yes
    - From file: yes
    - Encryption Key from env: yes
        - SHA256 HASH THE KEY SO ITS ALWAYS 32-BYTES: yes
    - Default (then write to file): yes