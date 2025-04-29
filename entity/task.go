package entity

import "time"

type TaskID int64

// type Task struct {
// 	ID          TaskID    `json:"id"`
// 	UserID      int64     `json:"user_id"`
// 	Title       string    `json:"title"`
// 	Description string     `json:"description,omitempty"`
// 	Deadline    *time.Time `json:"deadline,omitempty"`
// 	Priority    int       `json:"priority"` // 1:Low 2:Normal 3:High
// 	IsDone      bool      `json:"is_done"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

type Task struct {
	// ID タスク識別子
	ID          TaskID      `json:"id"          example:"1"`
	// UserID タスク所有ユーザ
	UserID      int64      `json:"user_id"     example:"2"`
	// Title タイトル
	Title       string     `json:"title"       example:"買い物"`
	// Description 詳細
	Description string     `json:"description,omitempty" example:"牛乳を買う"`
	// Deadline 期限
	Deadline    *time.Time `json:"deadline,omitempty" example:"2025-05-01T12:00:00+09:00"`
	// Priority 優先度 (1:高 2:中 3:低)
	Priority    int        `json:"priority"    example:"1"`
	// IsDone 完了フラグ
	IsDone      bool       `json:"is_done"     example:"false"`
	CreatedAt     time.Time  `json:"created"     example:"2025-04-10T09:00:00Z"`
	UpdatedAt     time.Time  `json:"updated"     example:"2025-04-10T09:00:00Z"`
}

type Tasks []*Task
