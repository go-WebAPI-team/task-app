package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
	"github.com/go-playground/validator/v10"
)

type TaskUpdater interface {
	UpdateTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type UpdateTask struct {
	Repo      TaskUpdater
	DB        store.Execer
	Validator *validator.Validate
}

func (ut *UpdateTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.ParseInt(idStr, 10, 64) //10進数として解釈し、int64 に収める(DBのPK がBIGINT(64bit) のため)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "invalid task ID format"}, http.StatusBadRequest)
		return
	}

	var in struct {
		Title       string     `json:"title"       validate:"required,max=255"`
		Description string     `json:"description"`
		Deadline    *time.Time `json:"deadline"`
		Priority    int        `json:"priority"    validate:"oneof=1 2 3"`
		IsDone      bool       `json:"is_done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := ut.Validator.Struct(in); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	t := &entity.Task{
		ID:          entity.TaskID(idInt),
		// TODO: 認証機能実装後にログインユーザーの ID を ctx から取得する
		UserID:      1,
		Title:       in.Title,
		Description: in.Description,
		Deadline:    in.Deadline,
		Priority:    in.Priority,
		IsDone:      in.IsDone,
	}
	if err := ut.Repo.UpdateTask(ctx, ut.DB, t); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	RespondJSON(ctx, w, struct{}{}, http.StatusOK)
}
