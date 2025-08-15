package httpx

import (
	"errors"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// ErrorMetadata contains additional validation error information
type ErrorMetadata struct {
	// Length/Count validations (applies to strings, arrays, maps)
	MinLength *int `json:"minLength,omitempty"` // Minimum length/items
	MaxLength *int `json:"maxLength,omitempty"` // Maximum length/items

	// Numeric validations
	MinValue *float64 `json:"minValue,omitempty"` // Minimum allowed value
	MaxValue *float64 `json:"maxValue,omitempty"` // Maximum allowed value

	// String regex validation
	Regex string `json:"regex,omitempty"` // Regex pattern for string validation

	// Enum validation
	AllowedValues []string `json:"allowedValues,omitempty"` // List of allowed values
}

// Detail represents additional error information
type ErrorDetail struct {
	// Code is the error code for this detail
	Code ErrorCode `json:"code"`

	// Location is the location of the error
	Location string `json:"location,omitempty"`

	// Message is the human-readable message
	Message string `json:"message"`

	// Metadata contains additional context about the error
	Metadata *ErrorMetadata `json:"metadata,omitempty"`
}

// Error implements the error interface
func (e ErrorDetail) Error() string {
	return e.Message
}

// Error represents a custom error in the application
type Error struct {
	// Internal error (not exposed to clients)
	internal error `json:"-"`

	// HTTP status code to return
	status int `json:"-"`

	// Code is the machine-readable error code for localization and client handling
	Code ErrorCode `json:"code,omitempty"`

	// Errors contains additional error details for complex errors
	Errors []ErrorDetail `json:"errors"`

	// Message is the human-readable message
	Message string `json:"message"`
	// msgf is a format string that can be used to generate a more specific message
	msgf string
}

// NewError creates a new error
func NewError(message string) Error {
	return Error{
		Errors:  []ErrorDetail{},
		Message: message,
	}
}

// Error returns the error message
func (e Error) Error() string {
	if e.internal == nil {
		return e.Message
	}
	return e.Message + "; internal=" + e.internal.Error()
}

func (e Error) Is(target error) bool {
	var ee Error
	if errors.As(target, &ee) {
		return ee.Code == e.Code
	}
	var ec ErrorCode
	if errors.As(target, &ec) {
		return e.Code == ec
	}
	return errors.Is(e.internal, target)
}

// LogValue implements slog.LogValuer for the internal error information.
func (e Error) LogValue() slog.Value {
	if e.internal == nil {
		return slog.StringValue(e.Message)
	}
	return slog.StringValue(e.Message + "; internal=" + e.internal.Error())
}

func (e Error) Clone() Error {
	e.Errors = slices.Clone(e.Errors)
	return e
}

// GetStatus returns the HTTP status code for the error
func (e Error) GetStatus() int {
	return e.status
}

// Unwrap returns the internal error
func (e Error) Unwrap() error {
	return e.internal
}

// WithCode sets the code for the error
func (e Error) WithCode(code ErrorCode) Error {
	e.Code = code
	return e
}

// WithDetails sets multiple details for the error
func (e Error) WithDetails(errors []ErrorDetail) Error {
	e.Errors = errors
	return e
}

// WithInternal adds internal error details
func (e Error) WithInternal(err error) Error {
	e.internal = err
	return e
}

// WithStatus sets the HTTP status code for the error
func (e Error) WithStatus(status int) Error {
	e.status = status
	return e
}

func init() {
	huma.NewError = NewStatusError
}

// NewStatusError creates a new custom error for huma
func NewStatusError(status int, message string, errs ...error) huma.StatusError {
	errors := make([]ErrorDetail, 0)
	code := ErrInvalidBody
	if len(errs) > 0 {
		for i, err := range errs {
			if err == nil {
				continue
			}
			switch err := err.(type) {
			case ErrorCode:
				errData := errCodeMap[err]
				code = err
				message = errData.Message
				status = errData.status
			case Error:
				code = err.Code
			case ErrorDetail:
				if i == 0 && status == 500 {
					errData := errCodeMap[code]
					message = errData.Message
					status = errData.status
				}
				errors = append(errors, err)
			case *huma.ErrorDetail:
				location := err.Location
				code, metadata := parseHumaError(err.Message)

				if code == ErrRequired {
					property := strings.ReplaceAll(err.Message, "expected required property ", "")
					property = strings.ReplaceAll(property, " to be present", "")
					if !strings.HasPrefix(location, "header") {
						location = location + "." + property
					}
				}

				detail := ErrorDetail{
					Code:     code,
					Location: location,
					Message:  err.Message,
					Metadata: metadata,
				}
				errors = append(errors, detail)
			default:
				errors = append(errors, ErrorDetail{
					Code:    ErrUnknown,
					Message: err.Error(),
				})
			}
		}
	}

	return NewError(message).WithStatus(status).WithDetails(errors).WithCode(code)
}

func parseHumaError(msg string) (ErrorCode, *ErrorMetadata) {
	var metadata *ErrorMetadata

	if strings.Contains(msg, "expected required property") || strings.Contains(msg, "parameter is missing") {
		return ErrRequired, metadata
	}

	if strings.Contains(msg, "expected string to be RFC 5322 email") {
		return ErrInvalidEmail, metadata
	}

	if strings.Contains(msg, "expected string to be RFC 3339 date-time") {
		return ErrInvalidDateTime, metadata
	}

	if strings.Contains(msg, "expected string to be RFC 3339 date") {
		return ErrInvalidDate, metadata
	}

	if strings.Contains(msg, "expected string to be RFC 3339 time") {
		return ErrInvalidTime, metadata
	}

	if strings.Contains(msg, "expected string to be RFC 4122 uuid") {
		return ErrInvalidUUID, metadata
	}

	if strings.Contains(msg, "expected string to be RFC 5890 hostname") {
		return ErrInvalidHostname, metadata
	}

	if strings.Contains(msg, "expected string to be RFC 2673 ipv4") {
		return ErrInvalidIPv4, metadata
	}

	if strings.Contains(msg, "expected string to be RFC 2373 ipv6") {
		return ErrInvalidIPv6, metadata
	}

	// Regex validation
	if strings.Contains(msg, "expected string to match pattern") {
		pattern := strings.TrimPrefix(msg, "expected string to match pattern ")
		metadata = &ErrorMetadata{Regex: pattern}
		return ErrInvalidValue, metadata
	}

	// Enum validation
	if strings.Contains(msg, "expected value to be one of") {
		items := strings.TrimPrefix(msg, "expected value to be one of ")
		items, _ = strconv.Unquote(items)
		metadata = &ErrorMetadata{AllowedValues: strings.Split(items, ", ")}
		return ErrInvalidValue, metadata
	}

	// Length validations with metadata
	if strings.Contains(msg, "expected length >=") {
		min := strings.TrimPrefix(msg, "expected length >= ")
		minLen, _ := strconv.Atoi(min)
		metadata = &ErrorMetadata{MinLength: &minLen}
		return ErrTooShort, metadata
	}

	if strings.Contains(msg, "expected length <=") {
		max := strings.TrimPrefix(msg, "expected length <= ")
		maxLen, _ := strconv.Atoi(max)
		metadata = &ErrorMetadata{MaxLength: &maxLen}
		return ErrTooLong, metadata
	}

	// Array validations with metadata
	if strings.Contains(msg, "expected array length >=") {
		min := strings.TrimPrefix(msg, "expected array length >= ")
		minLen, _ := strconv.Atoi(min)
		metadata = &ErrorMetadata{MinLength: &minLen}
		return ErrTooShort, metadata
	}

	if strings.Contains(msg, "expected array length <=") {
		max := strings.TrimPrefix(msg, "expected array length <= ")
		maxLen, _ := strconv.Atoi(max)
		metadata = &ErrorMetadata{MaxLength: &maxLen}
		return ErrTooLong, metadata
	}

	if strings.Contains(msg, "expected array items to be unique") {
		return ErrDuplicateItems, nil
	}

	// Numeric validations with metadata
	if strings.Contains(msg, "expected number <=") {
		max := strings.TrimPrefix(msg, "expected number <= ")
		maxVal, _ := strconv.ParseFloat(max, 64)
		metadata = &ErrorMetadata{MaxValue: &maxVal}
		return ErrTooLarge, metadata
	}

	if strings.Contains(msg, "expected number >=") {
		min := strings.TrimPrefix(msg, "expected number >= ")
		minVal, _ := strconv.ParseFloat(min, 64)
		metadata = &ErrorMetadata{MinValue: &minVal}
		return ErrTooSmall, metadata
	}

	return ErrUnknown, metadata
}
