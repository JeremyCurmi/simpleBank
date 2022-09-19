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
	Balance   float64   `db:"balance"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Transfer struct {
	ID         uint      `db:"id"`
	SenderID   uint      `json:"sender_id" db:"sender_id"`
	ReceiverID uint      `json:"receiver_id" db:"receiver_id"`
	Amount     float64   `json:"amount" db:"amount"`
	Status     string    `db:"status"`
	Timestamp  time.Time `db:"timestamp"`
}

type DBModel interface{}
