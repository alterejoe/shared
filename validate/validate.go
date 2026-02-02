package validate

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
)

var (
	v         = validator.New()
	sanitizer = bluemonday.StrictPolicy()
)

type ValidationErrors map[string]string

func (e ValidationErrors) Error() string {
	return "validation failed"
}

type options struct {
	applyDefaults bool
	fromQuery     bool
}

type Option func(*options)

// WithDefaults applies `default` struct tags to empty fields before validation.
func WithDefaults() Option {
	return func(o *options) {
		o.applyDefaults = true
	}
}

// FromQuery reads values from URL query params instead of form body.
func FromQuery() Option {
	return func(o *options) {
		o.fromQuery = true
	}
}

// ParseAndValidate is a generic function that:
// 1. Parses the form/query data from the request
// 2. Maps values to struct fields
// 3. Sanitizes all string fields (XSS protection)
// 4. Optionally applies defaults
// 5. Validates using struct tags
func ParseAndValidate[T any](r *http.Request, opts ...Option) (T, error) {
	var params T
	var cfg options

	for _, opt := range opts {
		opt(&cfg)
	}

	if err := r.ParseForm(); err != nil {
		return params, err
	}

	params = parseFormToStruct[T](r, cfg.fromQuery)
	sanitizeStruct(&params)

	if cfg.applyDefaults {
		applyDefaults(&params)
	}

	if err := v.Struct(params); err != nil {
		return params, formatErrors(err)
	}

	return params, nil
}

// parseFormToStruct uses reflection to map form/query values to struct fields.
// It checks for a `form` tag first (e.g., `form:"note_text"`),
// falling back to the lowercase field name (e.g., Text -> "text").
func parseFormToStruct[T any](r *http.Request, fromQuery bool) T {
	var result T
	val := reflect.ValueOf(&result).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		formKey := fieldType.Tag.Get("form")
		if formKey == "" {
			formKey = strings.ToLower(fieldType.Name)
		}

		if field.Kind() == reflect.String && field.CanSet() {
			if fromQuery {
				field.SetString(r.URL.Query().Get(formKey))
			} else {
				field.SetString(r.FormValue(formKey))
			}
		}
	}

	return result
}

// sanitizeStruct trims whitespace and strips HTML from all string fields.
func sanitizeStruct(ptr any) {
	val := reflect.ValueOf(ptr).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String && field.CanSet() {
			cleaned := strings.TrimSpace(field.String())
			cleaned = sanitizer.Sanitize(cleaned)
			field.SetString(cleaned)
		}
	}
}

// applyDefaults sets empty string fields to their `default` tag value.
func applyDefaults(ptr any) {
	val := reflect.ValueOf(ptr).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		defaultVal := fieldType.Tag.Get("default")
		if defaultVal == "" {
			continue
		}

		if field.Kind() == reflect.String && field.CanSet() && field.String() == "" {
			field.SetString(defaultVal)
		}
	}
}

// formatErrors converts validator errors into human-readable messages.
func formatErrors(err error) ValidationErrors {
	errors := make(ValidationErrors)

	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())

		switch e.Tag() {
		case "required":
			errors[field] = "This field is required"
		case "min":
			errors[field] = "Must be at least " + e.Param() + " characters"
		case "max":
			errors[field] = "Must be at most " + e.Param() + " characters"
		case "email":
			errors[field] = "Invalid email address"
		case "url":
			errors[field] = "Invalid URL"
		default:
			errors[field] = "Invalid value"
		}
	}

	return errors
}
