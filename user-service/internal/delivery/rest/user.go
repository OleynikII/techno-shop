package rest

import (
	"GoTest/internal/delivery/contracts/request"
	"GoTest/internal/delivery/contracts/response"
	"GoTest/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type Users interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)

	Create(ctx context.Context, dto domain.CreateUser) (int64, error)
	Update(ctx context.Context, id int64, dto domain.UpdateUser) error

	Delete(ctx context.Context, id int64) error
}

type UserHandler struct {
	users Users
}

func NewUserHandler(users Users) *UserHandler {
	return &UserHandler{users: users}
}

var domainToHTTPStatus = map[error]int{
	domain.ErrNotFound:              http.StatusNotFound,
	domain.ErrValidationIssuers:     http.StatusBadRequest,
	domain.ErrEmailAlreadyExists:    http.StatusConflict,
	domain.ErrUserNameAlreadyExists: http.StatusConflict,
}

func (h *UserHandler) handleError(w http.ResponseWriter, err error, logMsg string, args ...any) {
	for domainErr, statusCode := range domainToHTTPStatus {
		if errors.Is(err, domainErr) {
			http.Error(w, err.Error(), statusCode)
			return
		}
	}

	logArgs := append(args, "error", err)
	slog.Error(logMsg, logArgs...)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (h *UserHandler) parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

// GetAll godoc
// @Summary      Получить список всех пользователей
// @Description  Возвращает массив всех пользователей
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   response.User
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /users [get]
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.GetAll(r.Context())
	if err != nil {
		h.handleError(w, err, "failed to get users")
		return
	}

	res := make([]response.User, 0, len(users))
	for _, user := range users {
		res = append(res, response.User{
			ID:        user.ID,
			UserName:  user.UserName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

// GetByID godoc
// @Summary      Получить пользователя по ID
// @Description  Возвращает данные одного пользователя по его уникальному идентификатору
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.User
// @Failure      400  {string}  string "Invalid user ID"
// @Failure      404  {string}  string "Not Found"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := h.parseID(w, r)
	if !ok {
		return
	}

	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err, "failed to get user by id", "user_id", id)
		return
	}

	res := response.User{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

// Create godoc
// @Summary      Создать нового пользователя
// @Description  Принимает JSON с данными и создает нового пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      request.CreateUser  true  "Данные пользователя"
// @Success      201      {object}  map[string]int64    "id пользователя"
// @Failure      400      {string}  string              "Invalid request body / Validation error"
// @Failure      409      {string}  string              "Email or Username already exists"
// @Failure      500      {string}  string              "Internal Server Error"
// @Router       /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dto := domain.CreateUser{
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	id, err := h.users.Create(r.Context(), dto)
	if err != nil {
		h.handleError(w, err, "failed to create user")
		return
	}

	// ИСПРАВЛЕНО: Добавлен префикс /api/v1 в заголовок Location для соответствия новому роутингу
	w.Header().Set("Location", fmt.Sprintf("/api/v1/users/%d", id))
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	res := map[string]int64{"id": id}
	_ = json.NewEncoder(w).Encode(res)
}

// Update godoc
// @Summary      Обновить данные пользователя
// @Description  Обновляет профиль существующего пользователя по его ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      int                 true  "User ID"
// @Param        request  body      request.UpdateUser  true  "Новые данные"
// @Success      204      "Успешно обновлено, без тела ответа"
// @Failure      400      {string}  string              "Invalid ID / Invalid request body"
// @Failure      409      {string}  string              "Email or Username already exists"
// @Failure      500      {string}  string              "Internal Server Error"
// @Router       /users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := h.parseID(w, r)
	if !ok {
		return
	}

	var req request.UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dto := domain.UpdateUser{
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.users.Update(r.Context(), id, dto); err != nil {
		h.handleError(w, err, "failed to update user", "user_id", id)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete godoc
// @Summary      Удалить пользователя
// @Description  Удаляет запись пользователя по его ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204  "Успешно удалено"
// @Failure      400  {string}  string  "Invalid ID"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := h.parseID(w, r)
	if !ok {
		return
	}

	if err := h.users.Delete(r.Context(), id); err != nil {
		h.handleError(w, err, "failed to delete user", "user_id", id)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
