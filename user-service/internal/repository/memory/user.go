package memory

import (
	"GoTest/internal/domain"
	"context"
	"sync"
	"time"
)

type MemoryUserRepository struct {
	mu     sync.RWMutex
	users  []domain.User
	nextID int64
}

func NewUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:  make([]domain.User, 0),
		nextID: 1,
	}
}

func (r *MemoryUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]domain.User, len(r.users))
	copy(users, r.users)

	return users, nil
}

func (r *MemoryUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.users {
		if r.users[i].ID == id {
			return &r.users[i], nil
		}
	}

	return nil, domain.ErrNotFound
}

func (r *MemoryUserRepository) Create(ctx context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u.ID = r.nextID
	u.CreatedAt = time.Now().UTC()
	r.users = append(r.users, *u)
	r.nextID++

	return nil
}

func (r *MemoryUserRepository) Update(ctx context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.users {
		if r.users[i].ID == u.ID {
			r.users[i].UserName = u.UserName
			r.users[i].Email = u.Email
			r.users[i].PasswordHash = u.PasswordHash

			return nil
		}
	}

	return domain.ErrNotFound
}

func (r *MemoryUserRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.users {
		if r.users[i].ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}

	return domain.ErrNotFound
}

func (r *MemoryUserRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.users {
		if r.users[i].Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (r *MemoryUserRepository) ExistByUserName(ctx context.Context, userName string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.users {
		if r.users[i].UserName == userName {
			return true, nil
		}
	}
	return false, nil
}

func (r *MemoryUserRepository) ExistByEmailForUpdate(ctx context.Context, id int64, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.users {
		if r.users[i].Email == email && r.users[i].ID != id {
			return true, nil
		}
	}
	return false, nil
}

func (r *MemoryUserRepository) ExistByUserNameForUpdate(ctx context.Context, id int64, userName string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.users {
		if r.users[i].UserName == userName && r.users[i].ID != id {
			return true, nil
		}
	}
	return false, nil
}
