package request

type CreateUser struct {
	UserName string `json:"user_name" validate:"required,min=3,max=32,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUser struct {
	UserName string `json:"user_name" validate:"required,min=3,max=32,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
