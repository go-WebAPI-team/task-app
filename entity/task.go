package entity

import "time"

type TaskID int64

type Task struct {
	ID          TaskID    `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string     `json:"description,omitempty"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	Priority    int       `json:"priority"` // 1:Low 2:Normal 3:High
	IsDone      bool      `json:"is_done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Tasks []*Task
