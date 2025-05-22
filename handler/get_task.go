package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
)

type TaskGetter interface {
	GetTask(ctx context.Context, db store.Queryer, userID int64, id entity.TaskID) (*entity.Task, error)
	ListTagIDsByTaskID(ctx context.Context, db store.Queryer, taskID int64) ([]int64, error)
}

type GetTask struct {
	Repo TaskGetter
	DB   store.Queryer
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
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "invalid id"}, http.StatusBadRequest)
		return
	}

	// 認証済みユーザーのIDをctxから取得
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		RespondJSON(ctx, w, &ErrResponse{Message: "unauthorized"}, http.StatusUnauthorized)
		return
	}

	t, err := gt.Repo.GetTask(ctx, gt.DB, userID, entity.TaskID(idInt))
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "task not found"}, http.StatusNotFound)
		return
	}
	tagIDs, err := gt.Repo.ListTagIDsByTaskID(ctx, gt.DB, int64(t.ID))
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "failed to retrieve tags"}, http.StatusInternalServerError)
		return
	}
	type taskWithTags struct {
		*entity.Task
		TagIDs []int64 `json:"tag_ids"`
	}
	RespondJSON(ctx, w, taskWithTags{Task: t, TagIDs: tagIDs}, http.StatusOK)
}
