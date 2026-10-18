package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"user-service/internal/delivery/rest/response"
	"user-service/internal/domain"

	"github.com/go-playground/validator"
)

type BaseHandler struct {
	Validate *validator.Validate
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		Validate: validator.New(),
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
