package service

import (
	"context"
	"user-service/internal/domain"
)

type Users interface {
	GetPaginated(ctx context.Context, param domain.PaginationParam) (*domain.PaginationResult[domain.User], error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)

	Create(ctx context.Context, u *domain.User) error
	Update(ctx context.Context, u *domain.User) error

	Delete(ctx context.Context, id int64) error
}

type UserService struct {
	users Users
}

func NewUserService(users Users) *UserService {
	return &UserService{
		users: users,
	}
}

func (s *UserService) GetPaginated(ctx context.Context, param domain.PaginationParam) (*domain.PaginationResult[domain.User], error) {
	return s.users.GetPaginated(ctx, param)
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Create(ctx context.Context, dto domain.CreateUser) (*domain.User, error) {
	user := domain.User{
		UserName:     dto.UserName,
		Email:        dto.Email,
		PasswordHash: dto.Password + "hashed",
	}

	if err := s.users.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) Update(ctx context.Context, id int64, dto domain.UpdateUser) error {
	user := domain.User{
		ID:           id,
		UserName:     dto.UserName,
		Email:        dto.Email,
		PasswordHash: dto.Password + "hashed",
	}
	if err := s.users.Update(ctx, &user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	if err := s.users.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
