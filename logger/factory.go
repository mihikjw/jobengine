package logger

//Logger defines a class suitable for performing logging, error is mandatory but only has to be used for remote logging solutions
type Logger interface {
	Info(msg string) error
	Error(msg string) error
}

//NewLogger creates a class that can be used for logging
func NewLogger(logType string) Logger {
	switch {
	case logType == "std":
		return NewStdLogger()
	default:
		return nil
	}
}
