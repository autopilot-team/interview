package service

import (
	"autopilot/backends/api/internal/app"
	"autopilot/backends/api/internal/store"
)

// Manager is a collection of services used by the handlers/workers.
type Manager struct {
	Payment Paymenter
}

// NewManager creates a new service manager
func NewManager(container *app.Container, store *store.Manager) (*Manager, error) {
	paymentService, err := NewPayment(container)
	if err != nil {
		return nil, err
	}

	return &Manager{
		Payment: paymentService,
	}, nil
}

// Error is a custom error type that includes a code and an error
type Error struct {
	// Code is the error code
	Code string

	// Err is the underlying error that should only be used for logging and
	// debugging, not for user-facing messages
	Err error

	// Message is the error message
	Message string

	// Status is the HTTP status code associated with the error
	Status int
}

// Error returns the error message
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a new error with the given code and error
func NewError(code, message string, err error) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

// NewUnknownError creates a new unknown error with the given error
func NewUnknownError(err error) *Error {
	return NewError("unknown_error", "Unknown error", err)
}
