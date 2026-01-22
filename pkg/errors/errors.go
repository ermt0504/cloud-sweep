package errors

import (
	"errors"
	"fmt"
)

// Common error types
var (
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternalError     = errors.New("internal error")
	ErrServiceUnavailable = errors.New("service unavailable")
)

// AppError represents an application error with additional context
type AppError struct {
	Err     error
	Message string
	Code    string
	Details map[string]any
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(err error, message string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
	}
}

// NewWithCode creates a new AppError with an error code
func NewWithCode(err error, code, message string) *AppError {
	return &AppError{
		Err:     err,
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Is checks if an error is of a specific type
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target
func As(err error, target any) bool {
	return errors.As(err, target)
}
