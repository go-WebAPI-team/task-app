package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (dt *DeleteTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.ParseInt(idStr, 10, 64) //10進数として解釈し、int64 に収める(DBのPK がBIGINT(64bit) のため)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "invalid id"}, http.StatusBadRequest)
		return
	}
	// TODO: 認証機能実装後にログインユーザーの ID を ctx から取得する
	if err := dt.Repo.DeleteTask(ctx, dt.DB, 1, entity.TaskID(idInt)); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
