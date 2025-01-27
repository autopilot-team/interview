package handler

import (
	"autopilot/backends/api/internal/validator"
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
)

// testValidation is a helper function to test various validation scenarios
func testValidation[T any](t *testing.T, _ T, payload map[string]any, expectedBody string) {
	path := "/test"
	_, api := humatest.New(t)
	huma.NewError = NewCustomStatusError

	huma.Register(api, huma.Operation{
		Method: http.MethodPost,
		Path:   path,
	}, func(ctx context.Context, in *struct {
		Body T
	}) (*struct{}, error) {
		return &struct{}{}, nil
	})

	resp := api.Post(path, payload)
	body := strings.Trim(resp.Body.String(), "\n")

	assert.Equal(t, resp.Code, http.StatusUnprocessableEntity)
	assert.Equal(t, body, expectedBody)
}

func TestDateValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Date string `json:"date" format:"date"`
		}{},
		map[string]any{"date": "invalid-date"},
		`{"errors":[{"code":"INVALID_DATE","location":"body.date","message":"expected string to be RFC 3339 date"}],"message":"validation failed"}`,
	)
}

func TestDateTimeValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			DateTime string `json:"date_time" format:"date-time"`
		}{},
		map[string]any{"date_time": "invalid-date-time"},
		`{"errors":[{"code":"INVALID_DATE_TIME","location":"body.date_time","message":"expected string to be RFC 3339 date-time"}],"message":"validation failed"}`,
	)
}

func TestEmailValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"email"`
		}{},
		map[string]any{"value": "invalid-email"},
		`{"errors":[{"code":"INVALID_EMAIL","location":"body.value","message":"expected string to be RFC 5322 email: mail: missing '@' or angle-addr"}],"message":"validation failed"}`,
	)
}

func TestHostnameValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"hostname"`
		}{},
		map[string]any{"value": ""},
		`{"errors":[{"code":"INVALID_HOSTNAME","location":"body.value","message":"expected string to be RFC 5890 hostname"}],"message":"validation failed"}`,
	)
}

func TestIPv4Validation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"ipv4"`
		}{},
		map[string]any{"value": "invalid-ipv4"},
		`{"errors":[{"code":"INVALID_IPV4","location":"body.value","message":"expected string to be RFC 2673 ipv4"}],"message":"validation failed"}`,
	)
}

func TestIPv6Validation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"ipv6"`
		}{},
		map[string]any{"value": "invalid-ipv6"},
		`{"errors":[{"code":"INVALID_IPV6","location":"body.value","message":"expected string to be RFC 2373 ipv6"}],"message":"validation failed"}`,
	)
}

func TestMinLengthValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value string `json:"value" minLength:"3" maxLength:"6" pattern:"^[a-zA-Z0-9]+$"`
		}{},
		map[string]any{"value": "aa"},
		`{"errors":[{"code":"TOO_SHORT","location":"body.value","message":"expected length \u003e= 3","metadata":{"min_length":3}}],"message":"validation failed"}`,
	)
}

func TestMinValueValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value int `json:"value" minimum:"10"`
		}{},
		map[string]any{"value": 5},
		`{"errors":[{"code":"TOO_SMALL","location":"body.value","message":"expected number \u003e= 10","metadata":{"min_value":10}}],"message":"validation failed"}`,
	)
}

func TestMaxLengthValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value string `json:"value" minLength:"3" maxLength:"6" pattern:"^[a-zA-Z0-9]+$"`
		}{},
		map[string]any{"value": "aaabbbb"},
		`{"errors":[{"code":"TOO_LONG","location":"body.value","message":"expected length \u003c= 6","metadata":{"max_length":6}}],"message":"validation failed"}`,
	)
}

func TestMaxValueValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value int `json:"value" maximum:"100"`
		}{},
		map[string]any{"value": 150},
		`{"errors":[{"code":"TOO_LARGE","location":"body.value","message":"expected number \u003c= 100","metadata":{"max_value":100}}],"message":"validation failed"}`,
	)
}

func TestTimeValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"time"`
		}{},
		map[string]any{"value": "invalid-time"},
		`{"errors":[{"code":"INVALID_TIME","location":"body.value","message":"expected string to be RFC 3339 time"}],"message":"validation failed"}`,
	)
}

func TestUuidValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value string `json:"value" format:"uuid"`
		}{},
		map[string]any{"value": "invalid-uuid"},
		`{"errors":[{"code":"INVALID_UUID","location":"body.value","message":"expected string to be RFC 4122 uuid: invalid UUID length: 12"}],"message":"validation failed"}`,
	)
}

func TestArrayMinItemsValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Values []string `json:"values" minItems:"2"`
		}{},
		map[string]any{"values": []string{"one"}},
		`{"errors":[{"code":"TOO_SHORT","location":"body.values","message":"expected array length \u003e= 2","metadata":{"min_length":2}}],"message":"validation failed"}`,
	)
}

func TestArrayMaxItemsValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Values []string `json:"values" maxItems:"2"`
		}{},
		map[string]any{"values": []string{"one", "two", "three"}},
		`{"errors":[{"code":"TOO_LONG","location":"body.values","message":"expected array length \u003c= 2","metadata":{"max_length":2}}],"message":"validation failed"}`,
	)
}

func TestArrayUniqueItemsValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Values []string `json:"values" uniqueItems:"true"`
		}{},
		map[string]any{"values": []string{"one", "one"}},
		`{"errors":[{"code":"DUPLICATE_ITEMS","location":"body.values","message":"expected array items to be unique"}],"message":"validation failed"}`,
	)
}

func TestPasswordValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value validator.Password `json:"value"`
		}{},
		map[string]any{"value": ""},
		`{"errors":[{"code":"TOO_SHORT","location":"body.value","message":"expected length \u003e= 8","metadata":{"min_length":8}},{"code":"MISSING_UPPERCASE","location":"body.value","message":"expected at least one uppercase letter"},{"code":"MISSING_LOWERCASE","location":"body.value","message":"expected at least one lowercase letter"},{"code":"MISSING_NUMBER","location":"body.value","message":"expected at least one number"},{"code":"MISSING_SPECIAL","location":"body.value","message":"expected at least one special character"}],"message":"validation failed"}`,
	)
}

func TestRegexValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value string `json:"value" pattern:"^[a-z]+$"`
		}{},
		map[string]any{"value": "123"},
		`{"errors":[{"code":"INVALID_VALUE","location":"body.value","message":"expected string to match pattern ^[a-z]+$","metadata":{"regex":"^[a-z]+$"}}],"message":"validation failed"}`,
	)
}

func TestEnumValidation(t *testing.T) {
	testValidation(
		t,
		struct {
			Value string `json:"value" enum:"one,two,three"`
		}{},
		map[string]any{"value": "four"},
		`{"errors":[{"code":"INVALID_VALUE","location":"body.value","message":"expected value to be one of \"one, two, three\"","metadata":{"allowed_values":["one","two","three"]}}],"message":"validation failed"}`,
	)
}
