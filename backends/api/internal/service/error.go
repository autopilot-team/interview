package service

// Error is a custom error type that includes a code and an error
type Error struct {
	// Code is the error code
	Code string

	// Err is the underlying error that should only be used for logging and
	// debugging, not for user-facing messages
	Err error

	// Message is the error message
	Message string
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
