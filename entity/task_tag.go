package entity

import "time"

type TaskTagID int64

type TaskTag struct {
	ID     TaskTagID   `json:"id" example:"100"`
	TaskID TaskID `json:"task_id" example:"1"`
	TagID  TagID  `json:"tag_id" example:"10"`
	CreatedAt time.Time `json:"created_at" example:"2025-04-29T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-04-30T09:00:00Z"`
}