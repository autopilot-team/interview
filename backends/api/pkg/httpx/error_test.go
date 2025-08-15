package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testValidation is a helper function to test various validation scenarios
func testValidation[T any](t *testing.T, _ T, payload map[string]any, err error) {
	t.Helper()

	path := "/test"
	_, api := humatest.New(t)

	huma.Register(api, huma.Operation{
		Method: http.MethodPost,
		Path:   path,
	}, func(ctx context.Context, in *struct {
		Body T
	},
	) (*struct{}, error) {
		return &struct{}{}, nil
	})

	resp := api.Post(path, payload)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	assertErr(t, err, resp.Body)
}

func assertErr(t *testing.T, expected error, actual io.Reader, msgAndArgs ...any) {
	t.Helper()
	d := json.NewDecoder(actual)
	switch expected := expected.(type) {
	case ErrorCode:
		var err Error
		dErr := d.Decode(&err)
		require.NoError(t, dErr)
		if expected != err.Code {
			exp, act := fmt.Sprintf("%d (%s)", expected, expected.String()), fmt.Sprintf("%d (%s)", err.Code, err.Code.String())
			assert.Equal(t, exp, act, msgAndArgs...)
		}
	default:
		assert.Failf(t, "testutil.AssertErr: invalid usage", "asserting error to invalid type %T", expected)

	}
}

func TestDateValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Date string `json:"date" format:"date"`
		}{},
		map[string]any{"date": "invalid-date"},
		ErrInvalidBody,
	)
}

func TestDateTimeValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			DateTime string `json:"dateTime" format:"date-time"`
		}{},
		map[string]any{"dateTime": "invalid-date-time"},
		ErrInvalidBody,
	)
}

func TestEmailValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"email"`
		}{},
		map[string]any{"value": "invalid-email"},
		ErrInvalidBody,
	)
}

func TestHostnameValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"hostname"`
		}{},
		map[string]any{"value": ""},
		ErrInvalidBody,
	)
}

func TestIPv4Validation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"ipv4"`
		}{},
		map[string]any{"value": "invalid-ipv4"},
		ErrInvalidBody,
	)
}

func TestIPv6Validation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"ipv6"`
		}{},
		map[string]any{"value": "invalid-ipv6"},
		ErrInvalidBody,
	)
}

func TestMinLengthValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value string `json:"value" minLength:"3" maxLength:"6" pattern:"^[a-zA-Z0-9]+$"`
		}{},
		map[string]any{"value": "aa"},
		ErrInvalidBody,
	)
}

func TestMinValueValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value int `json:"value" minimum:"10"`
		}{},
		map[string]any{"value": 5},
		ErrInvalidBody,
	)
}

func TestMaxLengthValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value string `json:"value" minLength:"3" maxLength:"6" pattern:"^[a-zA-Z0-9]+$"`
		}{},
		map[string]any{"value": "aaabbbb"},
		ErrInvalidBody,
	)
}

func TestMaxValueValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value int `json:"value" maximum:"100"`
		}{},
		map[string]any{"value": 150},
		ErrInvalidBody,
	)
}

func TestTimeValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"time"`
		}{},
		map[string]any{"value": "invalid-time"},
		ErrInvalidBody,
	)
}

func TestUuidValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"uuid"`
		}{},
		map[string]any{"value": "invalid-uuid"},
		ErrInvalidBody,
	)
}

func TestArrayMinItemsValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Values []string `json:"values" minItems:"2"`
		}{},
		map[string]any{"values": []string{"one"}},
		ErrInvalidBody,
	)
}

func TestArrayMaxItemsValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Values []string `json:"values" maxItems:"2"`
		}{},
		map[string]any{"values": []string{"one", "two", "three"}},
		ErrInvalidBody,
	)
}

func TestArrayUniqueItemsValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Values []string `json:"values" uniqueItems:"true"`
		}{},
		map[string]any{"values": []string{"one", "one"}},
		ErrInvalidBody,
	)
}

func TestPasswordValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value Password `json:"value"`
		}{},
		map[string]any{"value": ""},
		ErrInvalidBody,
	)
}

func TestRegexValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value string `json:"value" pattern:"^[a-z]+$"`
		}{},
		map[string]any{"value": "123"},
		ErrInvalidBody,
	)
}

func TestEnumValidation(t *testing.T) {
	t.Parallel()
	testValidation(
		t,
		struct {
			Value string `json:"value" enum:"one,two,three"`
		}{},
		map[string]any{"value": "four"},
		ErrInvalidBody,
	)
}
