package logger

// Logger defines an object suitable for performing logging
type Logger interface {
	Info(msg string) error
	Error(msg string) error
	Fatal(msg string) error
}
