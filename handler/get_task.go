package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
)

type TaskGetter interface {
	GetTask(ctx context.Context, db store.Execer, userID int64, id entity.TaskID) (*entity.Task, error)
}

type GetTask struct {
	Repo TaskGetter
	DB   store.Execer
}

// GetTask godoc
// @Summary      タスク詳細取得
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "タスクID"
// @Success      200 {object} entity.Task
// @Failure      404 {object} handler.ErrResponse
// @Router       /tasks/{id} [get]
func (gt *GetTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.ParseInt(idStr, 10, 64) //10進数として解釈し、int64 に収める(DBのPK がBIGINT(64bit) のため)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "invalid id"}, http.StatusBadRequest)
		return
	}
	// TODO: 認証機能実装後にログインユーザーの ID を ctx から取得する
	t, err := gt.Repo.GetTask(ctx, gt.DB, 1, entity.TaskID(idInt))
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusNotFound)
		return
	}
	RespondJSON(ctx, w, t, http.StatusOK)
}
