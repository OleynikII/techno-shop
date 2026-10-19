package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"user-service/internal/delivery/rest/response"
	"user-service/internal/domain"

	"github.com/go-playground/validator/v10"
)

type BaseHandler struct {
	validate *validator.Validate
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *BaseHandler) ParsePagination(r *http.Request) domain.PaginationParam {
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil || page == 0 {
		page = 1
	}

	size, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil || size == 0 {
		size = 10
	}

	if size > 100 {
		size = 100
	}

	return domain.PaginationParam{
		Limit:  size,
		Offset: (page - 1) * size,
	}
}

func (b *BaseHandler) ParseID(w http.ResponseWriter, r *http.Request, paramName string) (int64, bool) {
	idStr := r.PathValue(paramName)
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		validationErr := domain.NewValidationError(map[string]string{
			paramName: "invalid identifier format",
		})

		b.HandleError(w, validationErr, "failed to parse path id", "param", paramName, "value", idStr)
		return 0, false
	}

	return id, true
}

func (b *BaseHandler) ParseAndValidate(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return domain.NewValidationError(map[string]string{
			"body": "invalid json body",
		})
	}

	if err := b.validate.Struct(dst); err != nil {
		fields := make(map[string]string)

		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				errorMsg := msgForTag(e.Tag(), e.Param())
				fields[e.Field()] = errorMsg
			}
			return domain.NewValidationError(fields)
		}

		return err
	}

	return nil
}

func msgForTag(tag string, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "alphanum":
		return "Must contain only letters and numbers (no spaces)"
	case "min":
		return "Must be at least " + param + " characters long"
	case "max":
		return "Must be at most " + param + " characters long"
	default:
		// Если тег неизвестен, отдаем его как есть или дефолтный текст
		return "failed validation: " + tag
	}
}

func (b *BaseHandler) HandleError(w http.ResponseWriter, err error, logMsg string, args ...any) {
	statusCode, errResp := response.NewError(err)

	logArgs := append(args, "error", err, "status_code", statusCode)

	if statusCode >= 500 {
		slog.Error("CRITICAL: "+logMsg, logArgs...)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errResp)
}
