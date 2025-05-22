package entity

import "time"

type TagID int64

type Tag struct {
	ID        TagID     `json:"id" example:"10"`
	Name      string    `json:"name" validate:"required" example:"urgent"`
	UserID    UserID    `json:"user_id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2025-04-29T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-04-30T09:00:00Z"`
}

type Tags []*Tag
