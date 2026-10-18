package response

import "time"

type User struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
