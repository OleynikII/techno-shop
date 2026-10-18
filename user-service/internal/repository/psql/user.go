package psql

import (
	"context"
	"database/sql"
	"errors"
	"user-service/internal/domain"
)

type PsqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *PsqlUserRepository {
	return &PsqlUserRepository{
		db: db,
	}
}

func (r *PsqlUserRepository) GetPaginated(ctx context.Context, param domain.PaginationParam) (*domain.PaginationResult[domain.User], error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var totalCount uint64
	countQuery := `SELECT COUNT(*) FROM users`
	if err := tx.QueryRowContext(ctx, countQuery).Scan(&totalCount); err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return &domain.PaginationResult[domain.User]{
			TotalCount: 0,
			Data:       []domain.User{},
		}, nil
	}

	selectQuery := `
		SELECT id, user_name, email, created_at 
		FROM users 
		ORDER BY id ASC 
		LIMIT $1 OFFSET $2`

	rows, err := tx.QueryContext(ctx, selectQuery, param.Limit, param.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.UserName, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Используем ваше новое имя структуры здесь
	return &domain.PaginationResult[domain.User]{
		TotalCount: totalCount,
		Data:       users,
	}, nil
}

func (r *PsqlUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, user_name, email, created_at FROM users WHERE id = $1`

	var u domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.UserName, &u.Email, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *PsqlUserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users(user_name, email, password_hash)
		values ($1,$2,$3)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query, u.UserName, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			if constraintName(err) == "users_email_key" {
				return domain.NewConflictError("email", u.Email)
			}
			if constraintName(err) == "users_user_name_key" {
				return domain.NewConflictError("user_name", u.UserName)
			}
		}
		return err
	}
	return nil
}

func (r *PsqlUserRepository) Update(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE users
		SET user_name = $1, email = $2, password_hash = $3
		WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query, u.UserName, u.Email, u.PasswordHash, u.ID)
	if err != nil {
		if isUniqueViolation(err) {
			if constraintName(err) == "users_email_key" {
				return domain.NewConflictError("email", u.Email)
			}
			if constraintName(err) == "users_user_name_key" {
				return domain.NewConflictError("user_name", u.UserName)
			}
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *PsqlUserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
