package response

import (
	"errors"
	"net/http"
	"user-service/internal/domain"
)

type Error struct {
	Message string            `json:"message"`
	Details map[string]string `json:"details"`
}

func NewError(err error) (int, Error) {
	if valErr, ok := errors.AsType[*domain.ValidationError](err); ok {
		return http.StatusBadRequest, Error{
			Message: "Validation failed",
			Details: valErr.Fields,
		}
	}

	if conflictErr, ok := errors.AsType[*domain.ConflictError](err); ok {
		return http.StatusConflict, Error{
			Message: "Resource conflict",
			Details: map[string]string{
				conflictErr.Field: conflictErr.Value,
			},
		}
	}

	if errors.Is(err, domain.ErrNotAuthorized) {
		return http.StatusUnauthorized, Error{
			Message: err.Error(),
		}
	}
	if errors.Is(err, domain.ErrAccessDenied) {
		return http.StatusForbidden, Error{
			Message: err.Error(),
		}
	}

	if errors.Is(err, domain.ErrNotFound) {
		return http.StatusNotFound, Error{
			Message: err.Error(),
		}
	}

	return http.StatusInternalServerError, Error{
		Message: "Internal server error",
	}

}
