package models

//Queue represents a queue of jobs to be executed
type Queue struct {
	Jobs        []*Job
	Permissions *QueuePermissions
	Size        uint8
}
