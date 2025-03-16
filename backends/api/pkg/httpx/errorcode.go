//go:generate go run autopilot/tools stringer --type=ErrorCode --trimprefix Err
package httpx

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

var errCodeMap = map[ErrorCode]Error{
	ErrUnknown:                 mkErr("An unknown error has occured.", http.StatusInternalServerError),
	ErrEntityNotFound:          mkErr("Entity not found.", http.StatusNotFound),
	ErrUnauthenticated:         mkErr("Unauthenticated.", http.StatusUnauthorized),
	ErrInsufficientPermissions: mkErr("Insufficient permissions.", http.StatusForbidden),

	ErrInvalidBody:                  mkErr("Invalid request body", http.StatusBadRequest),
	ErrRequired:                     mkErr("This field is required.", http.StatusBadRequest),
	ErrInvalidValue:                 mkErr("Invalid value provided.", http.StatusBadRequest),
	ErrInvalidDate:                  mkErr("Invalid date format.", http.StatusBadRequest),
	ErrInvalidDateTime:              mkErr("Invalid datetime format.", http.StatusBadRequest),
	ErrInvalidTime:                  mkErr("Invalid time format.", http.StatusBadRequest),
	ErrInvalidEmail:                 mkErr("Invalid email address.", http.StatusBadRequest),
	ErrInvalidHostname:              mkErr("Invalid hostname.", http.StatusBadRequest),
	ErrInvalidIPv4:                  mkErr("Invalid IPv4 address.", http.StatusBadRequest),
	ErrInvalidIPv6:                  mkErr("Invalid IPv6 address.", http.StatusBadRequest),
	ErrInvalidUUID:                  mkErr("Invalid UUID.", http.StatusBadRequest),
	ErrMissingLowercase:             mkErr("Expected at least one lowercase letter.", http.StatusBadRequest),
	ErrMissingUppercase:             mkErr("Expected at least one uppercase letter.", http.StatusBadRequest),
	ErrMissingNumber:                mkErr("Expected at least one number.", http.StatusBadRequest),
	ErrMissingSpecial:               mkErr("Expected at least one special character.", http.StatusBadRequest),
	ErrTooShort:                     mkErr("Password is too short.", http.StatusBadRequest),
	ErrTooLong:                      mkErr("Password is too long.", http.StatusBadRequest),
	ErrDuplicateItems:               mkErr("Duplicate items are not allowed.", http.StatusBadRequest),
	ErrTooSmall:                     mkErr("Value is too small.", http.StatusBadRequest),
	ErrTooLarge:                     mkErr("Value is too large.", http.StatusBadRequest),
	ErrInvalidImageFormat:           mkErr("Invalid image format.", http.StatusBadRequest),
	ErrInvalidCursor:                mkErr("Invalid cursor format", http.StatusBadRequest),
	ErrInvalidTurnstileToken:        mkErr("Invalid Turnstile token.", http.StatusUnauthorized),
	ErrFailedToVerifyTurnstileToken: mkErr("Failed to verify Turnstile token.", http.StatusUnauthorized),
	ErrInvalidCurrency:              mkErr("Invalid currency code.", http.StatusBadRequest),
	ErrInvalidCountry:               mkErr("Invalid country code.", http.StatusBadRequest),
	ErrInvalidFinancialAmount:       mkErr("Invalid financial amount.", http.StatusBadRequest),

	// Identity Errors
	ErrAccountLocked:                mkErr("The account is temporarily locked.", http.StatusTooManyRequests),
	ErrEmailNotVerified:             mkErr("Email verification is required.", http.StatusUnauthorized),
	ErrInvalidCredentials:           mkErr("Invalid login credentials.", http.StatusUnauthorized),
	ErrInvalidRefreshToken:          mkErr("Invalid refresh token.", http.StatusUnauthorized),
	ErrInvalidName:                  mkErr("Invalid API key name", http.StatusUnprocessableEntity),
	ErrConnectionNotFound:           mkErr("Connection not found.", http.StatusNotFound),
	ErrInvalidConnectionCredentials: mkErr("Invalid connection credentials.", http.StatusUnprocessableEntity),

	ErrEmailExists:           mkErr("Email already exists.", http.StatusUnprocessableEntity),
	ErrInvalidOrExpiredToken: mkErr("The verification token is invalid or expired.", http.StatusUnauthorized),
	ErrUserNotFound:          mkErr("User not found.", http.StatusNotFound),

	ErrInvalidTwoFactorCode:    mkErr("Invalid two-factor code.", http.StatusUnauthorized),
	ErrTwoFactorNotEnabled:     mkErr("Two-factor authentication is not enabled.", http.StatusBadRequest),
	ErrTwoFactorAlreadyEnabled: mkErr("Two-factor authentication is already enabled.", http.StatusBadRequest),
	ErrTwoFactorPending:        mkErr("Two-factor authentication verification pending.", http.StatusBadRequest),
	ErrBackupCodeValidation:    mkErr("Invalid or used backup code.", http.StatusUnauthorized),
	ErrTwoFactorLocked:         mkErr("Two-factor authentication is locked.", http.StatusTooManyRequests),

	ErrPaymentNotFound: mkErr("Payment not found", http.StatusNotFound),

	ErrUnused: mkErr("Internal Server Error", http.StatusInternalServerError),
}

func mkErr(message string, status int) Error {
	return NewError(message).WithStatus(status)
}

// ErrorCode represents a unique error code that can be used for localization
type ErrorCode int

// General Errors
const (
	ErrUnknown ErrorCode = iota + 1
	ErrUnauthenticated
	ErrEntityNotFound
	ErrInsufficientPermissions
)

// Validation Errors
const (
	ErrInvalidBody ErrorCode = iota + 1_000
	ErrRequired

	ErrInvalidValue

	// Format validation
	ErrInvalidDate
	ErrInvalidDateTime
	ErrInvalidTime
	ErrInvalidEmail
	ErrInvalidHostname
	ErrInvalidIPv4
	ErrInvalidIPv6
	ErrInvalidUUID

	// Password validation
	ErrMissingLowercase
	ErrMissingUppercase
	ErrMissingNumber
	ErrMissingSpecial
	ErrTooShort
	ErrTooLong

	// Array Validation
	ErrDuplicateItems

	// Numeric Validation
	ErrTooSmall
	ErrTooLarge

	// File validation
	ErrInvalidImageFormat

	// Application Validation
	ErrInvalidCursor
	ErrInvalidTurnstileToken
	ErrFailedToVerifyTurnstileToken
	ErrInvalidCurrency
	ErrInvalidCountry
	ErrInvalidFinancialAmount
)

// Service/Module errors
const (
	ErrAccountLocked ErrorCode = iota + 10_000
	ErrEmailNotVerified
	ErrInvalidCredentials
	ErrInvalidRefreshToken
	ErrInvalidName
	ErrConnectionNotFound
	ErrInvalidConnectionCredentials

	ErrEmailExists
	ErrInvalidOrExpiredToken
	ErrUserNotFound

	ErrInvalidTwoFactorCode
	ErrTwoFactorNotEnabled
	ErrTwoFactorAlreadyEnabled
	ErrTwoFactorPending
	ErrBackupCodeValidation
	ErrTwoFactorLocked

	ErrPaymentNotFound

	ErrUnused
)

var errorEnumKeys = func() []any {
	names := make([]any, 0, len(_ErrorCode_map))
	for _, val := range _ErrorCode_map {
		names = append(names, val)
	}
	slices.SortFunc(names, func(a, b any) int { return strings.Compare(a.(string), b.(string)) })
	return names
}()

func (e ErrorCode) Schema(r huma.Registry) *huma.Schema {
	m := r.Map()
	if _, ok := m["ErrorCode"]; !ok {
		m["ErrorCode"] = &huma.Schema{
			Type:        huma.TypeString,
			Enum:        errorEnumKeys,
			Description: "The standardized error code",
		}
	}
	return &huma.Schema{
		Ref: "#/components/schemas/ErrorCode",
	}
}

func (e ErrorCode) Error() string {
	err, ok := errCodeMap[e]
	if !ok {
		return "<nil>"
	}
	return err.Message
}

func (e ErrorCode) LogValue() slog.Value {
	err, ok := errCodeMap[e]
	if !ok {
		return slog.StringValue(e.String())
	}
	return slog.StringValue(err.Message)
}

func (e *ErrorCode) UnmarshalJSON(data []byte) error {
	str, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	for c, v := range _ErrorCode_map {
		if str == v {
			*e = c
			return nil
		}
	}
	*e = 0
	return nil
}

func (e ErrorCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e ErrorCode) Errorf(args ...any) Error {
	err, ok := errCodeMap[e]
	if !ok {
		return Error{}
	}
	if err.msgf == "" {
		return err
	}
	err.Message = fmt.Sprintf(err.msgf, args...)
	return err.WithCode(e)
}

func (e ErrorCode) WithInternal(err error) Error {
	ee, ok := errCodeMap[e]
	if !ok {
		return Error{}
	}
	return ee.WithCode(e).WithInternal(err)
}

func (e ErrorCode) WithDetails(path string, details ...error) Error {
	err, ok := errCodeMap[e]
	if !ok {
		return Error{}
	}
	var errDetails []ErrorDetail
	for _, detail := range details {
		var d ErrorDetail
		if errors.As(detail, &d) {
			if path != "" {
				d.Location = path + "." + d.Location
			}
			errDetails = append(errDetails, d)
			continue
		}
		errDetails = append(errDetails, ErrorDetail{
			Code:    ErrUnknown,
			Message: detail.Error(),
		})

	}
	return err.WithCode(e).WithDetails(errDetails)
}

func (e ErrorCode) WithLocation(location string) ErrorDetail {
	err, ok := errCodeMap[e]
	if !ok {
		return ErrorDetail{}
	}
	return ErrorDetail{
		Code:     e,
		Location: location,
		Message:  err.Message,
	}
}

func init() {
	if len(errCodeMap) != len(_ErrorCode_map) {
		for code := range _ErrorCode_map {
			if _, ok := errCodeMap[code]; !ok {
				panic("error code missing: " + code.String())
			}
		}
	}
}
