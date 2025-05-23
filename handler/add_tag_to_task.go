package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
)

type TagAdder interface {
	AddTagToTask(ctx context.Context, db store.Execer, userID int64, t *entity.TaskTag) error
}

type AddTagToTask struct {
	Repo TagAdder
	DB   store.Execer
	// JSONリクエストボディを受け取らないためValidatorは不要
}

// AddTagToTask godoc
// @Summary      タスクにタグを紐付け
// @Description  タスク ID とタグ ID を指定して関連づけます
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id       path     int  true  "タスクID"
// @Param        tag_id   path     int  true  "タグID"
// @Success      200 {object} entity.TaskTag
// @Failure      400 {object} handler.ErrResponse
// @Failure      404 {object} handler.ErrResponse
// @Router       /tasks/{id}/tags/{tag_id} [post]
func (at *AddTagToTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskIDStr := chi.URLParam(r, "task_id")
	tagIDStr := chi.URLParam(r, "tag_id")

	taskInt, err := strconv.ParseInt(taskIDStr, 10, 64) // 10進数として解釈し、int64 に収める(DBのPK がBIGINT(64bit) のため)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "invalid id"}, http.StatusBadRequest)
		return
	}
	tagInt, err := strconv.ParseInt(tagIDStr, 10, 64) // 10進数として解釈し、int64 に収める(DBのPK がBIGINT(64bit) のため)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "invalid tag ID"}, http.StatusBadRequest)
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

	now := time.Now()
	t := &entity.TaskTag{
		TaskID:    entity.TaskID(taskInt),
		TagID:     entity.TagID(tagInt),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := at.Repo.AddTagToTask(ctx, at.DB, userID, t); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			RespondJSON(ctx, w, &ErrResponse{Message: "task or tag not found or unauthorized"}, http.StatusNotFound)
			return
		}
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	RespondJSON(ctx, w, &entity.TaskTag{
		ID:        t.ID,
		TaskID:    t.TaskID,
		TagID:     t.TagID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}, http.StatusOK)
}
