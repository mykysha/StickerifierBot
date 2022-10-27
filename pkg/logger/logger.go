// Package logger is used to log messages.
package logger

import (
	"io"
	"log"
)

type Logger interface {
	// Info logs an info message.
	Info(msg string)
	// Error logs an error message.
	Error(msg string)
	// Fatal logs a fatal message and exits.
	Fatal(msg string)
}

type DefaultLogger struct {
	// logger is the logger.
	logger *log.Logger
}

// NewDefault creates a new logger.
func NewDefault(out io.Writer, prefix string) *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(out, prefix, log.LstdFlags),
	}
}

// Info logs an info message.
func (l *DefaultLogger) Info(msg string) {
	l.logger.Println("INFO " + msg)
}

// Error logs an error message.
func (l *DefaultLogger) Error(msg string) {
	l.logger.Println("ERROR " + msg)
}

// Fatal logs a fatal message and exits.
func (l *DefaultLogger) Fatal(msg string) {
	l.logger.Fatal("FATAL " + msg)
}
