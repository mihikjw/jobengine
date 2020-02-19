package models

type Config struct {
	Version float64
	Port    int
	Queues  []Queue
}
