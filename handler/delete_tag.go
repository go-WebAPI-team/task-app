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

type TagDeleter interface {
	DeleteTag(ctx context.Context, db store.Execer, userID int64, id entity.TagID) error
}

type DeleteTag struct {
	Repo TagDeleter
	DB   store.Execer
}

// DeleteTag godoc
// @Summary      タグを削除
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "タグID"
// @Success      200 {object} handler.EmptyResponse
// @Failure      404 {object} handler.ErrResponse
// @Router       /tags/{id} [delete]
func (dt *DeleteTag) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.ParseInt(idStr, 10, 64) //10進数として解釈し、int64 に収める(DBのPK がBIGINT(64bit) のため)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "Invalid tag ID format"}, http.StatusBadRequest)
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

	if err := dt.Repo.DeleteTag(ctx, dt.DB, userID, entity.TagID(idInt)); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
