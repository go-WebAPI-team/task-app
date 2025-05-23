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

type TaskDeleter interface {
	DeleteTask(ctx context.Context, db store.Execer, userID int64, id entity.TaskID) error
}

type DeleteTask struct {
	Repo TaskDeleter
	DB   store.Execer
}

// DeleteTask godoc
// @Summary      タスクを削除
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "タスクID"
// @Success      200 {object} handler.EmptyResponse
// @Failure      404 {object} handler.ErrResponse
// @Router       /tasks/{id} [delete]
func (dt *DeleteTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.ParseInt(idStr, 10, 64) //10進数として解釈し、int64 に収める(DBのPK がBIGINT(64bit) のため)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "Invalid task ID format"}, http.StatusBadRequest)
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
	
	if err := dt.Repo.DeleteTask(ctx, dt.DB, userID, entity.TaskID(idInt)); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
