package domain

import "strings"

type ValidationError struct {
	Fields map[string]string `json:"fields"`
}

func NewValidationError(fields map[string]string) error {
	return &ValidationError{Fields: fields}
}

func (e *ValidationError) Error() string {
	var msgs []string
	for k, v := range e.Fields {
		msgs = append(msgs, k+": "+v)
	}
	return "validation failed: " + strings.Join(msgs, "; ")
}

func (e *ValidationError) Unwrap() error {
	return ErrValidation
}
