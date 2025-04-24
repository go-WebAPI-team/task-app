package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
	"github.com/go-playground/validator/v10"
)

type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type AddTask struct {
	Repo      TaskAdder
	DB        store.Execer
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var in struct {
		Title       string     `json:"title"       validate:"required,max=255"`
		Description string     `json:"description"`
		Deadline    *time.Time `json:"deadline"`
		Priority    int        `json:"priority"    validate:"oneof=1 2 3"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := at.Validator.Struct(in); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	t := &entity.Task{
		// TODO: 認証機能実装後にログインユーザーの ID を ctx から取得する
		UserID:      1, 
		Title:       in.Title,
		Description: in.Description,
		Deadline:    in.Deadline,
		Priority:    in.Priority,
		IsDone:      false,
	}

	if err := at.Repo.AddTask(ctx, at.DB, t); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	RespondJSON(ctx, w, struct{ ID int `json:"id"` }{ID: int(t.ID)}, http.StatusOK)
}
