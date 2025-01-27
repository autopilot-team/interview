package validator

import (
	"fmt"
	"unicode"

	"github.com/danielgtaylor/huma/v2"
)

// Password is a huma validation type that must be at least 8 characters
// long and contain at least one uppercase letter, one lowercase letter, one
// number, and one special character.
type Password string

// Resolve implements the huma.ResolverWithPath interface.
func (p Password) Resolve(ctx huma.Context, prefix *huma.PathBuffer) []error {
	var errors []error

	minLength := 8
	if len(p) < minLength {
		errors = append(errors, &huma.ErrorDetail{
			Location: prefix.String(),
			Message:  fmt.Sprintf("expected length >= %d", minLength),
			Value:    p,
		})
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range p {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errors = append(errors, &huma.ErrorDetail{
			Location: prefix.String(),
			Message:  "expected at least one uppercase letter",
			Value:    p,
		})
	}

	if !hasLower {
		errors = append(errors, &huma.ErrorDetail{
			Location: prefix.String(),
			Message:  "expected at least one lowercase letter",
			Value:    p,
		})
	}

	if !hasNumber {
		errors = append(errors, &huma.ErrorDetail{
			Location: prefix.String(),
			Message:  "expected at least one number",
			Value:    p,
		})
	}

	if !hasSpecial {
		errors = append(errors, &huma.ErrorDetail{
			Location: prefix.String(),
			Message:  "expected at least one special character",
			Value:    p,
		})
	}

	return errors
}

// Ensure our resolver meets the expected interface
var _ huma.ResolverWithPath = (*Password)(nil)
