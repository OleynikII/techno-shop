package domain

import "time"

type BaseEntity struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}
