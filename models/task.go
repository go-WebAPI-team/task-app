package models

type Task struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Deadline    string `json:"deadline"`
	IsDone      bool   `json:"is_done"`
	//User        *UserInfo `json:"user,omitempty"`
}
