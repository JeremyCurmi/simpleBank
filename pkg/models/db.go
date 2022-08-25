package models

import "time"

type User struct {
	ID uint
	UserName string
	Password string
	CreatedAt time.Time `db:"created_at"`
}