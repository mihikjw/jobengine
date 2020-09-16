package logger

// NewLogger creates a class that can be used for logging
func NewLogger(logType string) Logger {
	switch {
	case logType == "std":
		return NewStdLogger()
	default:
		return nil
	}
}
