package models

//Config represents the application config
type Config struct {
	Version      float64
	Port         int
	Queues       map[string]*QueuePermissions
	CryptoSecret string
}
