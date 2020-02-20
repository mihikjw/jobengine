package models

type Config struct {
	Version float64
	Port    int
	Queues  map[string]*Queue
}
