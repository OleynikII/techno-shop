package psql

import "database/sql"

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{
		db: db,
	}
}
