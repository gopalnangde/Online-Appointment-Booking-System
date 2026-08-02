package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validate is the shared validator instance used across the application
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

// FormatValidationErrors converts validator.ValidationErrors into human-readable messages.
// Each field error is mapped to a clear, user-friendly string.
func FormatValidationErrors(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, formatFieldError(e))
		}
	}

	return errors
}

// formatFieldError converts a single FieldError into a readable message
func formatFieldError(e validator.FieldError) string {
	field := e.Field()

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, e.Param())
	case "numeric":
		return fmt.Sprintf("%s must contain only digits", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
