package handler

import (
	"autopilot/backends/api/internal/service"
	"autopilot/backends/api/internal/validator"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// ErrorCode represents a unique error code that can be used for localization
type ErrorCode string

// ErrorCode represents a unique error code that can be used for localization
const (
	// Field presence
	ErrRequired ErrorCode = "REQUIRED"

	// Format validations
	ErrInvalidDate           ErrorCode = "INVALID_DATE"
	ErrInvalidDateTime       ErrorCode = "INVALID_DATE_TIME"
	ErrInvalidEmail          ErrorCode = "INVALID_EMAIL"
	ErrInvalidHostname       ErrorCode = "INVALID_HOSTNAME"
	ErrInvalidIPv4           ErrorCode = "INVALID_IPV4"
	ErrInvalidIPv6           ErrorCode = "INVALID_IPV6"
	ErrInvalidTime           ErrorCode = "INVALID_TIME"
	ErrInvalidTurnstileToken ErrorCode = "INVALID_TURNSTILE_TOKEN"
	ErrInvalidUUID           ErrorCode = "INVALID_UUID"

	// String validations
	ErrMissingLowercase ErrorCode = "MISSING_LOWERCASE"
	ErrMissingUppercase ErrorCode = "MISSING_UPPERCASE"
	ErrMissingNumber    ErrorCode = "MISSING_NUMBER"
	ErrMissingSpecial   ErrorCode = "MISSING_SPECIAL"
	ErrTooShort         ErrorCode = "TOO_SHORT"
	ErrTooLong          ErrorCode = "TOO_LONG"

	// Array validation
	ErrDuplicateItems ErrorCode = "DUPLICATE_ITEMS"

	// Numeric validations
	ErrTooSmall ErrorCode = "TOO_SMALL"
	ErrTooLarge ErrorCode = "TOO_LARGE"

	// Generic validation
	ErrInvalidValue ErrorCode = "INVALID_VALUE" // For regex, enum, exact length, etc.

	// Turnstile validation
	ErrUnableToVerifyTurnstileToken ErrorCode = "UNABLE_TO_VERIFY_TURNSTILE_TOKEN"
	ErrFailedToVerifyTurnstileToken ErrorCode = "FAILED_TO_VERIFY_TURNSTILE_TOKEN"
)

// Error messages
var (
	ErrMessageInvalidTurnstileToken string = "invalid Turnstile token"
)

// ValidationMetadata contains additional validation information
type ValidationMetadata struct {
	// Length/Count validations (applies to strings, arrays, maps)
	MinLength *int `json:"min_length,omitempty"` // Minimum length/items
	MaxLength *int `json:"max_length,omitempty"` // Maximum length/items

	// Numeric validations
	MinValue *float64 `json:"min_value,omitempty"` // Minimum allowed value
	MaxValue *float64 `json:"max_value,omitempty"` // Maximum allowed value

	// String regex validation
	Regex *string `json:"regex,omitempty"` // Regex pattern for string validation

	// Enum validation
	AllowedValues []string `json:"allowed_values,omitempty"` // List of allowed values
}

// ErrorDetail is a detail of an error
type ErrorDetail struct {
	// Code is a unique identifier for this error that can be used for localization
	Code string `json:"code"`

	// Location is the location of the error
	Location string `json:"location"`

	// Message is a human-readable message to return to the client
	Message string `json:"message"`

	// Metadata contains additional validation information
	Metadata *ValidationMetadata `json:"metadata,omitempty"`
}

// Error implements the error interface
func (e *ErrorDetail) Error() string {
	return e.Message
}

// Error implements the standard error interface for APIs
type Error struct {
	status int

	// Code is a unique identifier for this error that can be used for localization
	Code string `json:"code,omitempty"`

	// Errors is a map of additional details to return to the client
	Errors []ErrorDetail `json:"errors"`

	// Message is a human-readable message to return to the client
	Message string `json:"message"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// GetStatus implements huma.StatusError interface
func (e *Error) GetStatus() int {
	return e.status
}

// NewCustomStatusError creates a new custom error for huma
func NewCustomStatusError(status int, message string, errs ...error) huma.StatusError {
	errors := make([]ErrorDetail, 0)
	if len(errs) > 0 {
		for _, err := range errs {
			switch err := err.(type) {
			case *huma.ErrorDetail:
				location := err.Location
				code, metadata := parseValidationError(err.Message)

				if code == ErrRequired {
					property := strings.ReplaceAll(err.Message, "expected required property ", "")
					property = strings.ReplaceAll(property, " to be present", "")
					location = location + "." + property
				}

				detail := ErrorDetail{
					Code:     string(code),
					Location: location,
					Message:  err.Message,
					Metadata: metadata,
				}
				errors = append(errors, detail)
			}
		}
	}

	return &Error{
		status:  status,
		Message: message,
		Errors:  errors,
	}
}

// ConvertServiceError converts a service error to a huma error
func ConvertServiceError(status int, err *service.Error) huma.StatusError {
	return &Error{
		status:  status,
		Code:    err.Code,
		Message: err.Error(),
		Errors:  []ErrorDetail{},
	}
}

// parseValidationError parses an error message to extract both the error code and validation metadata
func parseValidationError(message string) (ErrorCode, *ValidationMetadata) {
	var metadata *ValidationMetadata

	if strings.Contains(message, "expected required property") {
		return ErrRequired, metadata
	}

	if strings.Contains(message, validator.ErrUnableToVerifyTurnstileToken) {
		return ErrUnableToVerifyTurnstileToken, metadata
	}

	if strings.Contains(message, validator.ErrFailedToVerifyTurnstileToken) {
		return ErrFailedToVerifyTurnstileToken, metadata
	}

	// Format validations
	if strings.Contains(message, "expected string to be RFC 5322 email") {
		return ErrInvalidEmail, metadata
	}

	if strings.Contains(message, "expected string to be RFC 3339 date-time") {
		return ErrInvalidDateTime, metadata
	}

	if strings.Contains(message, "expected string to be RFC 3339 date") {
		return ErrInvalidDate, metadata
	}

	if strings.Contains(message, "expected string to be RFC 3339 time") {
		return ErrInvalidTime, metadata
	}

	if strings.Contains(message, "expected string to be RFC 4122 uuid") {
		return ErrInvalidUUID, metadata
	}

	if strings.Contains(message, "expected string to be RFC 5890 hostname") {
		return ErrInvalidHostname, metadata
	}

	if strings.Contains(message, "expected string to be RFC 2673 ipv4") {
		return ErrInvalidIPv4, metadata
	}

	if strings.Contains(message, "expected string to be RFC 2373 ipv6") {
		return ErrInvalidIPv6, metadata
	}

	// Password validations
	if strings.Contains(message, "expected at least one uppercase letter") {
		return ErrMissingUppercase, metadata
	}

	if strings.Contains(message, "expected at least one lowercase letter") {
		return ErrMissingLowercase, metadata
	}

	if strings.Contains(message, "expected at least one number") {
		return ErrMissingNumber, metadata
	}

	if strings.Contains(message, "expected at least one special character") {
		return ErrMissingSpecial, metadata
	}

	// Regex validation
	if strings.Contains(message, "expected string to match pattern") {
		if pattern := extractQuoted(message, "pattern "); pattern != nil {
			metadata = &ValidationMetadata{
				Regex: pattern,
			}

			return ErrInvalidValue, metadata
		}
	}

	// Enum validation
	if strings.Contains(message, "expected value to be one of") {
		if values := extractList(message, "one of"); len(values) > 0 {
			metadata = &ValidationMetadata{
				AllowedValues: values,
			}

			return ErrInvalidValue, metadata
		}
	}

	// Length validations with metadata
	if strings.Contains(message, "expected length >=") {
		if minLen := extractNumber(message, "length >="); minLen != nil {
			metadata = &ValidationMetadata{
				MinLength: minLen,
			}

			return ErrTooShort, metadata
		}
	}

	if strings.Contains(message, "expected length <=") {
		if maxLen := extractNumber(message, "length <="); maxLen != nil {
			metadata = &ValidationMetadata{
				MaxLength: maxLen,
			}

			return ErrTooLong, metadata
		}
	}

	// Array validations with metadata
	if strings.Contains(message, "expected array length >=") {
		if minItems := extractNumber(message, "array length >="); minItems != nil {
			metadata = &ValidationMetadata{
				MinLength: minItems,
			}

			return ErrTooShort, metadata
		}
	}

	if strings.Contains(message, "expected array length <=") {
		if maxItems := extractNumber(message, "array length <="); maxItems != nil {
			metadata = &ValidationMetadata{
				MaxLength: maxItems,
			}

			return ErrTooLong, metadata
		}
	}

	if strings.Contains(message, "expected array items to be unique") {
		return ErrDuplicateItems, nil
	}

	// Numeric validations with metadata
	if strings.Contains(message, "expected number >=") {
		if minVal := extractFloat(message, "number >="); minVal != nil {
			metadata = &ValidationMetadata{
				MinValue: minVal,
			}

			return ErrTooSmall, metadata
		}
	}

	if strings.Contains(message, "expected number <=") {
		if maxVal := extractFloat(message, "number <="); maxVal != nil {
			metadata = &ValidationMetadata{
				MaxValue: maxVal,
			}

			return ErrTooLarge, metadata
		}
	}

	return "", metadata
}

// extractFloat extracts a float64 from a message containing a phrase
func extractFloat(message, phrase string) *float64 {
	parts := strings.Split(message, phrase)
	if len(parts) != 2 {
		return nil
	}

	numStr := strings.Fields(strings.TrimSpace(parts[1]))[0]
	if num, err := strconv.ParseFloat(numStr, 64); err == nil {
		return &num
	}

	return nil
}

// extractNumber extracts a number from a message containing a phrase
func extractNumber(message, phrase string) *int {
	parts := strings.Split(message, phrase)
	if len(parts) != 2 {
		return nil
	}

	// Split by space and take first part to handle cases like "3 characters"
	numStr := strings.Fields(strings.TrimSpace(parts[1]))[0]
	if num, err := strconv.Atoi(numStr); err == nil {
		return &num
	}

	return nil
}

// extractQuoted extracts a quoted string from a message containing a phrase
func extractQuoted(message, phrase string) *string {
	parts := strings.Split(message, phrase)
	if len(parts) < 1 {
		return nil
	}

	return &parts[len(parts)-1]
}

// extractList extracts a comma-separated list from a message containing a phrase
func extractList(message, phrase string) []string {
	parts := strings.Split(message, phrase)
	if len(parts) != 2 {
		return nil
	}

	// Split by commas and clean up each value
	values := strings.Split(strings.TrimSpace(parts[1]), ",")
	result := make([]string, 0, len(values))
	for _, v := range values {
		// Remove quotes and spaces
		cleaned := strings.Trim(strings.TrimSpace(v), "\"")
		if cleaned != "" {
			result = append(result, cleaned)
		}
	}

	return result
}
