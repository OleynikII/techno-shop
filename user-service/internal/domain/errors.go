package domain

import "errors"

var (
	ErrNotFound = errors.New("resource not found")

	ErrConflict = errors.New("already exists")

	ErrValidation = errors.New("Validation issues")

	ErrAccessDenied  = errors.New("Access denied")
	ErrNotAuthorized = errors.New("Not authorized")
)
