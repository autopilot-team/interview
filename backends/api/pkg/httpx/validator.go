package httpx

import (
	"encoding/json"
	"strconv"
	"unicode"

	"github.com/Rhymond/go-money"
	"github.com/danielgtaylor/huma/v2"
)

// Type Validators

var (
	_ huma.SchemaProvider   = (*Currency)(nil)
	_ huma.ResolverWithPath = (*Currency)(nil)
)

type Currency money.Currency

func (c Currency) Schema(r huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type:        huma.TypeString,
		Format:      "currency",
		Description: "The currency code in ISO 4217 format.",
	}
}

func (c *Currency) Resolve(ctx huma.Context, prefix *huma.PathBuffer) (errors []error) {
	cur := money.GetCurrency(c.Code)
	if cur == nil {
		errors = append(errors, ErrInvalidCurrency.WithLocation(prefix.String()))
	} else {
		*c = Currency(*cur)
	}
	return
}

func (c *Currency) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.Code)
}

func (c Currency) MarshalJSON() ([]byte, error) {
	str := "\"\""
	if c.Code != "" {
		str = "\"" + c.Code + "\""
	}
	return []byte(str), nil
}

var (
	_ huma.SchemaProvider   = (*Money)(nil)
	_ huma.ResolverWithPath = (*Money)(nil)
)

type Money int64

func (m Money) Schema(r huma.Registry) *huma.Schema {
	return &huma.Schema{
		OneOf: []*huma.Schema{
			{Type: huma.TypeInteger},
			{Type: huma.TypeString, Pattern: `\d+`},
		},
		Description: `All monetary amounts should be provided in the minor unit. For example:
  - 1000 to charge 10 USD (or any other two-decimal currency)
  - 10 to charge 10 JPY (or any other zero-decimal currency)`,
	}
}

func (m Money) Resolve(ctx huma.Context, prefix *huma.PathBuffer) (errors []error) {
	if m < 0 {
		errors = append(errors, ErrInvalidFinancialAmount.WithLocation(prefix.String()))
	}
	return
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*m = Money(i)
		return nil
	}

	s := "-1"
	_ = json.Unmarshal(data, &s)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		i = -1
	}
	*m = Money(i)
	return nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	str := strconv.FormatInt(int64(m), 10)
	return []byte(str), nil
}

var _ huma.ResolverWithPath = (*Password)(nil)

const MinPasswordLength = 8

// Password is a huma validation type that must be at least 8 characters
// long and contain at least one uppercase letter, one lowercase letter, one
// number, and one special character.
type Password string

// Resolve implements the huma.ResolverWithPath interface.
func (p Password) Resolve(ctx huma.Context, prefix *huma.PathBuffer) (errors []error) {
	if len(p) < MinPasswordLength {
		errors = append(errors, ErrTooShort.WithLocation(prefix.String()))
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

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
		errors = append(errors, ErrMissingUppercase.WithLocation(prefix.String()))
	}
	if !hasLower {
		errors = append(errors, ErrMissingLowercase.WithLocation(prefix.String()))
	}
	if !hasNumber {
		errors = append(errors, ErrMissingNumber.WithLocation(prefix.String()))
	}
	if !hasSpecial {
		errors = append(errors, ErrMissingSpecial.WithLocation(prefix.String()))
	}
	return
}

var _ huma.ResolverWithPath = (*TurnstileToken)(nil)

// TurnstileToken is a huma validation type that validates a Cloudflare Turnstile token.
type TurnstileToken string

func (t TurnstileToken) Resolve(ctx huma.Context, prefix *huma.PathBuffer) (errors []error) {
	ctxx := ctx.Context()
	turnstile, ok := ctxx.Value(TurnstileKey).(Turnstiler)
	if !ok {
		return
	}

	// Verify Cloudflare Turnstile token
	ok, err := turnstile.Verify(ctxx, string(t), "")
	if err != nil {
		errors = append(errors, ErrFailedToVerifyTurnstileToken.WithLocation(prefix.String()))
	} else {
		if !ok {
			errors = append(errors, ErrInvalidTurnstileToken.WithLocation(prefix.String()))
		}
	}

	return
}
