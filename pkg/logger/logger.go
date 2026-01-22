package logger

import (
	"log"
	"os"
)

// Level represents a log level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger represents a structured logger
type Logger struct {
	level  Level
	debug  *log.Logger
	info   *log.Logger
	warn   *log.Logger
	error  *log.Logger
}

// New creates a new Logger
func New(level Level) *Logger {
	return &Logger{
		level: level,
		debug: log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
		info:  log.New(os.Stdout, "[INFO]  ", log.Ldate|log.Ltime|log.Lshortfile),
		warn:  log.New(os.Stdout, "[WARN]  ", log.Ldate|log.Ltime|log.Lshortfile),
		error: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...any) {
	if l.level <= LevelDebug {
		l.debug.Printf(msg, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...any) {
	if l.level <= LevelInfo {
		l.info.Printf(msg, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...any) {
	if l.level <= LevelWarn {
		l.warn.Printf(msg, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...any) {
	if l.level <= LevelError {
		l.error.Printf(msg, args...)
	}
}

// Default logger instance
var defaultLogger = New(LevelInfo)

// Debug logs a debug message using the default logger
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// Info logs an info message using the default logger
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Warn logs a warning message using the default logger
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Error logs an error message using the default logger
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// SetLevel sets the log level for the default logger
func SetLevel(level Level) {
	defaultLogger.level = level
}
