package form

import (
	"autopilot/backends/api/pkg/httpx"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/tidwall/gjson"
)

var (
	ErrInvalidJSON        = errors.New("invalid json input")
	ErrInvalidDestination = errors.New("invalid scan destination")
)

type Type string

const (
	TypeText     Type = "text"
	TypePassword Type = "password"
	// Split into int/uint/decimal
	// TypeNumber   = "number"
	TypeSelect Type = "select"
)

// Option is an option for a field.
type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Field struct {
	// Label is the name of the field that. It can be used for localization.
	Label string `json:"label"`

	// Options is the options for the field (enum type).
	Options []Option `json:"options"`

	// Required is true if the field is required.
	Required bool `json:"required"`

	// Type is the data type of the field.
	Type Type `json:"type"`

	// Validate is an optional function to validate the content of the field.
	Validate func(raw string) error `json:"-"`
}

// Section is a collection of fields
type Section struct {
	Label       string  `json:"label"`
	Description string  `json:"description"`
	Fields      []Field `json:"fields"`
}

// Parse reads json input data and validates it to the schema.
// If dst is not nil, it will write the data into dst.
// TODO: fix error handling responses
func (s *Section) Parse(dst any, data json.RawMessage) []error {
	if !gjson.ValidBytes(data) {
		return []error{ErrInvalidJSON}
	}
	raw := gjson.ParseBytes(data)

	val := reflect.ValueOf(dst)
	var base reflect.Value
	if dst != nil {
		if val.Kind() != reflect.Pointer {
			return []error{ErrInvalidDestination}
		}
		base = val.Elem()
		if base.Kind() != reflect.Struct {
			return []error{ErrInvalidDestination}
		}
	}

	var errs []error
	for _, f := range s.Fields {
		data := raw.Get(f.Label)
		if err := check(f, data); err != nil {
			errs = append(errs, err)
			continue
		}

		// Validation only, do not set.
		if dst == nil {
			continue
		}
		field := base.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, f.Label)
		})
		if !field.IsValid() || !field.CanSet() {
			return []error{ErrInvalidDestination}
		}
		// FIXME: support more than just string destinations
		field.SetString(data.Str)
	}

	return errs
}

func (s *Section) ParseMap(m map[string]string, data json.RawMessage) []error {
	if !gjson.ValidBytes(data) {
		return []error{ErrInvalidJSON}
	}
	raw := gjson.ParseBytes(data)
	var errs []error
	for _, f := range s.Fields {
		data := raw.Get(f.Label)
		if err := check(f, data); err != nil {
			errs = append(errs, err)
			continue
		}
		m[f.Label] = data.Str
	}
	return errs
}

func check(field Field, data gjson.Result) error {
	if field.Required && data.Str == "" {
		return httpx.ErrRequired.WithLocation(field.Label)
	}

	switch field.Type {
	case TypeText, TypePassword:
		// do nothing
	case TypeSelect:
		var ok bool
		for _, opt := range field.Options {
			if data.Str == opt.Value {
				ok = true
				break
			}
		}
		if !ok {
			return httpx.ErrInvalidValue.WithLocation(field.Label)
		}
	}
	if field.Validate != nil {
		if err := field.Validate(data.Raw); err != nil {
			if errors.As(err, &httpx.ErrorDetail{}) {
				return err
			}
			return httpx.ErrInvalidValue.WithLocation(field.Label)
		}
	}
	return nil
}
