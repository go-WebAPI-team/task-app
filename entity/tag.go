package entity

import "time"

type TagID int64
type UserID int64

type Tag struct {
	ID      TagID      `json:"id"`
	Name    string     `json:"name"`
	UserID  UserID 	   `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Tags []*Tag