package model

import "fmt"

// ValidationError represents a field-level validation failure.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error returns a human-readable representation of the validation error.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ErrFieldRequired creates a ValidationError for a missing required field.
func ErrFieldRequired(field string) *ValidationError {
	return &ValidationError{Field: field, Message: "is required"}
}

// ErrFieldInvalid creates a ValidationError for an invalid field value.
func ErrFieldInvalid(field, reason string) *ValidationError {
	return &ValidationError{Field: field, Message: reason}
}
