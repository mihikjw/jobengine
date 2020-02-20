package models

//QueuePermissions represents who has read/write access to a queue
type QueuePermissions struct {
	Read  []string
	Write []string
}
