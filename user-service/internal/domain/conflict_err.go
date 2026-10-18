package domain

import "fmt"

type ConflictError struct {
	Field string
	Value string
}

func NewConflictError(field, value string) error {
	return &ConflictError{Field: field, Value: value}
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Value)
}

func (e *ConflictError) Unwrap() error {
	return ErrConflict
}
