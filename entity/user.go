package entity

import "time"

// UserID は型エイリアス
type UserID int64

type User struct {
	ID        UserID    `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
