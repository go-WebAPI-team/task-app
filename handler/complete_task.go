package handler

import (
	"context"
	"net/http"
	"errors"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
)

type TaskCompleter interface {
	ToggleTaskDone(ctx context.Context, db store.Execer, userID int64, id entity.TaskID) error
}

type ToggleCompleteTask struct {
	Repo TaskCompleter
	DB   store.Execer
}

// CompleteTask godoc
// @Summary      タスクを完了状態に更新
// @Description  タスクの is_done を true に変更します
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "タスクID"
// @Success      200 {object} handler.EmptyResponse
// @Failure      400 {object} handler.ErrResponse
// @Failure      404 {object} handler.ErrResponse
// @Router       /tasks/{id}/complete [patch]
func (tc *ToggleCompleteTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.ParseInt(idStr, 10, 64) // 10進数として解釈し、int64 に収める(DBのPK がBIGINT(64bit) のため)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "invalid id"}, http.StatusBadRequest)
		return
	}

	// ------------------------------
	// 認証ユーザー取得：認証チェック済み ctx から userID 抽出 
	// ------------------------------
    userID, ok := auth.GetUserID(ctx)
    if !ok {
        RespondJSON(ctx, w, &ErrResponse{Message: "unauthorized"}, http.StatusUnauthorized)
        return
    }

	if err := tc.Repo.ToggleTaskDone(ctx, tc.DB, userID, entity.TaskID(idInt)); err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			RespondJSON(ctx, w, &ErrResponse{Message: "Task not found"}, http.StatusNotFound)
		} else {
			RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
