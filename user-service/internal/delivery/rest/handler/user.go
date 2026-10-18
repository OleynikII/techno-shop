package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"user-service/internal/delivery/rest/request"
	"user-service/internal/delivery/rest/response"
	"user-service/internal/domain"
)

type Users interface {
	GetPaginated(ctx context.Context, param domain.PaginationParam) (*domain.PaginationResult[domain.User], error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)

	Create(ctx context.Context, dto domain.CreateUser) (*domain.User, error)
	Update(ctx context.Context, id int64, dto domain.UpdateUser) error

	Delete(ctx context.Context, id int64) error
}

type UserHandler struct {
	*BaseHandler
	users Users
}

func NewUserHandler(base *BaseHandler, users Users) *UserHandler {
	return &UserHandler{
		BaseHandler: base,
		users:       users}
}

// GetAll godoc
// @Summary      Получить при помощи пагинации пользователей
// @Description  Возвращает пагинационный ответ пользователей
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page  query     int  false  "Номер страницы (дефолт: 1)"
// @Param        size  query     int  false  "Размер страницы (дефолт: 10, макс: 100)"
// @Success      200  {object}  response.PaginationResult[response.User]
// @Failure      500  {object}  response.Error
// @Router       /users [get]
func (h *UserHandler) GetPaginated(w http.ResponseWriter, r *http.Request) {
	paginationParam := h.ParsePagination(r)
	paginatedResult, err := h.users.GetPaginated(r.Context(), paginationParam)
	if err != nil {
		h.HandleError(w, err, "failed to get paginated users")
		return
	}

	mappedUsers := make([]response.User, 0, len(paginatedResult.Data))
	for _, user := range paginatedResult.Data {
		mappedUsers = append(mappedUsers, response.User{
			ID:        user.ID,
			UserName:  user.UserName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	res := response.Paginated[response.User]{
		TotalCount: paginatedResult.TotalCount,
		Data:       mappedUsers,
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
// @Failure      400  {object}  response.Error
// @Failure      404  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := h.ParseID(w, r, "id")
	if !ok {
		return
	}

	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		h.HandleError(w, err, "failed to get user by id", "user_id", id)
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
// @Success      201      {object}  response.User   "Созданный пользователь"
// @Failure      400      {object}  response.Error
// @Failure      409      {object}  response.Error
// @Failure      500      {object}  response.Error
// @Router       /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var req request.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		validationErr := domain.NewValidationError(map[string]string{
			"body": "invalid json format or empty body",
		})
		h.HandleError(w, validationErr, "failed to decode create request body")
		return
	}

	vailidationErr := h.Validate.Struct(req)
	if vailidationErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(vailidationErr.Error())

		return
	}

	dto := domain.CreateUser{
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.users.Create(r.Context(), dto)
	if err != nil {
		h.HandleError(w, err, "failed to create user")
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/api/v1/users/%d", user.ID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	res := response.User{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

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
// @Failure      400      {object}  response.Error
// @Failure      409      {object}  response.Error
// @Failure      500      {object}  response.Error
// @Router       /users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	id, ok := h.ParseID(w, r, "id")
	if !ok {
		return
	}

	var req request.UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		validationErr := domain.NewValidationError(map[string]string{
			"body": "invalid json format or empty body",
		})
		h.HandleError(w, validationErr, "failed to decode update request body", "user_id", id)
		return
	}

	dto := domain.UpdateUser{
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.users.Update(r.Context(), id, dto); err != nil {
		h.HandleError(w, err, "failed to update user", "user_id", id)
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
// @Failure      400  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := h.ParseID(w, r, "id")
	if !ok {
		return
	}

	if err := h.users.Delete(r.Context(), id); err != nil {
		h.HandleError(w, err, "failed to delete user", "user_id", id)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
