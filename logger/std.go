package logger

import (
	"log"
	"os"
)

// StdLogger logs to StdOut, StdErr
type StdLogger struct {
	stdOut *log.Logger
	stdErr *log.Logger
}

// NewStdLogger creates a new instance of a logger to StdOut, StdErr
func NewStdLogger() *StdLogger {
	return &StdLogger{
		stdOut: log.New(os.Stdout, "INFO: ", log.LstdFlags),
		stdErr: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
	}
}

// Info logs a message to StdOut
func (l *StdLogger) Info(msg string) error {
	l.stdOut.Print(msg)
	return nil
}

// Error logs a message to StdErr
func (l *StdLogger) Error(msg string) error {
	l.stdErr.Print(msg)
	return nil
}

// Fatal logs a message and exits the application, designed to be used in startup
func (l *StdLogger) Fatal(msg string) error {
	l.stdErr.Fatal(msg)
	return nil
}
