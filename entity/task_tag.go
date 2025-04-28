package entity

import "time"

type TaskTagID int64

type TaskTag struct {
	ID     TaskTagID   `json:"id"`
	TaskID TaskID `json:"task_id"`
	TagID  TagID  `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}