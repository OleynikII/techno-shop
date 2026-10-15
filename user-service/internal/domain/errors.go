package domain

import "errors"

var (
	ErrNotFound = errors.New("resource not found")

	ErrEmailAlreadyExists    = errors.New("Email already exists")
	ErrUserNameAlreadyExists = errors.New("User name already exists")

	ErrValidationIssuers = errors.New("Validation issues")

	ErrAccessDenied  = errors.New("Access denied")
	ErrNotAuthorized = errors.New("Not authorized")
)
