package models

type Config struct {
	Version int
	Port    int
	Queues  []Queue
}
