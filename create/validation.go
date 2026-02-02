// internal/validate/validate.go
package create

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
)

// Initialized once at package load, both are thread-safe and reusable.
// Creating these per-request would be wasteful.
var (
	v         = validator.New()
	sanitizer = bluemonday.StrictPolicy()
)

// ValidationErrors is a map of field name -> error message.
// Implementing the error interface lets us return it as an error
// while still being able to type assert it later for field-specific messages.
type ValidationErrors map[string]string

func (e ValidationErrors) Error() string {
	return "validation failed"
}

// ParseAndValidate is a generic function that:
// 1. Parses the form data from the request
// 2. Maps form values to struct fields
// 3. Sanitizes all string fields (XSS protection)
// 4. Validates using struct tags
//
// The generic [T any] lets us return a typed struct instead of interface{},
// so callers get compile-time type safety: params.Text works, params.Foo errors.
func ParseAndValidate[T any](r *http.Request) (T, error) {
	var params T

	if err := r.ParseForm(); err != nil {
		return params, err
	}

	params = parseFormToStruct[T](r)
	sanitizeStruct(&params)

	if err := v.Struct(params); err != nil {
		return params, formatErrors(err)
	}

	return params, nil
}

// parseFormToStruct uses reflection to map form values to struct fields.
// It checks for a `form` tag first (e.g., `form:"note_text"`),
// falling back to the lowercase field name (e.g., Text -> "text").
//
// This removes the repetitive r.FormValue("x") calls from every handler.
func parseFormToStruct[T any](r *http.Request) T {
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
			field.SetString(r.FormValue(formKey))
		}
	}

	return result
}

// sanitizeStruct iterates over all string fields and:
// 1. Trims whitespace (no leading/trailing spaces)
// 2. Strips HTML tags (XSS protection via bluemonday)
//
// This ensures ALL string input is sanitized automatically.
// Without this, you'd need to remember to sanitize in every handler.
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

// formatErrors converts validator's error type into our simple map.
// This gives us human-readable messages instead of cryptic ones like
// "Key: 'NoteParams.Text' Error:Field validation for 'Text' failed on the 'required' tag"
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
