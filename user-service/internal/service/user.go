package service

import (
	"GoTest/internal/domain"
	"context"
)

type Users interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)

	Create(ctx context.Context, u *domain.User) error
	Update(ctx context.Context, u *domain.User) error

	Delete(ctx context.Context, id int64) error

	ExistByUserName(ctx context.Context, userName string) (bool, error)
	ExistByEmail(ctx context.Context, email string) (bool, error)
	ExistByUserNameForUpdate(ctx context.Context, id int64, userName string) (bool, error)
	ExistByEmailForUpdate(ctx context.Context, id int64, email string) (bool, error)
}

type UserService struct {
	users Users
}

func NewUserService(users Users) *UserService {
	return &UserService{
		users: users,
	}
}

func (s *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.users.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	panic("TEst")

	return users, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Create(ctx context.Context, dto domain.CreateUser) (int64, error) {

	// TODO: Валидация модели CreateUser
	var validationErr error
	if validationErr != nil {
		return 0, domain.ErrValidationIssuers
	}

	if exists, err := s.users.ExistByEmail(ctx, dto.Email); err != nil {
		return 0, err
	} else if exists {
		return 0, domain.ErrEmailAlreadyExists
	}

	if exists, err := s.users.ExistByUserName(ctx, dto.UserName); err != nil {
		return 0, err
	} else if exists {
		return 0, domain.ErrUserNameAlreadyExists
	}

	user := domain.User{
		UserName:     dto.UserName,
		Email:        dto.Email,
		PasswordHash: dto.Password,
	}

	if err := s.users.Create(ctx, &user); err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (s *UserService) Update(ctx context.Context, id int64, dto domain.UpdateUser) error {
	// TODO: Валидация модели CreateUser
	var validationErr error
	if validationErr != nil {
		return domain.ErrValidationIssuers
	}

	if exists, err := s.users.ExistByEmailForUpdate(ctx, id, dto.Email); err != nil {
		return err
	} else if exists {
		return domain.ErrEmailAlreadyExists
	}

	if exists, err := s.users.ExistByUserNameForUpdate(ctx, id, dto.UserName); err != nil {
		return err
	} else if exists {
		return domain.ErrUserNameAlreadyExists
	}

	user := domain.User{
		ID:           id,
		UserName:     dto.UserName,
		Email:        dto.Email,
		PasswordHash: dto.Password,
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
