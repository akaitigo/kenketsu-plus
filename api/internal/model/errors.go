package model

import "fmt"

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func ErrFieldRequired(field string) *ValidationError {
	return &ValidationError{Field: field, Message: "is required"}
}

func ErrFieldInvalid(field, reason string) *ValidationError {
	return &ValidationError{Field: field, Message: reason}
}
