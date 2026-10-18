package domain

type User struct {
	BaseEntity
	UserName     string `db:"user_name"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}

type CreateUser struct {
	UserName string
	Email    string
	Password string
}

type UpdateUser struct {
	UserName string
	Email    string
	Password string
}
