package domain

type User struct {
	base
	UserName     string
	Email        string
	PasswordHash string
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
