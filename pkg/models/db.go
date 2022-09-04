package models

import "time"

type User struct {
	ID        uint      `db:"id"`
	UserName  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

type Account struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	UserID    uint      `db:"user_id"`
	Balance   int       `db:"balance"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
