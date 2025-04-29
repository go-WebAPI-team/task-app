package handler

import (
    "context"
    "errors"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/go-webapi-team/task-app/entity"
    "github.com/go-webapi-team/task-app/store"
)

type TagRemover interface {
    DeleteTagFromTask(ctx context.Context, db store.Execer, userID int64, t *entity.TaskTag) error
}

type DeleteTagFromTask struct {
    Repo TagRemover
    DB   store.Execer
}

// DeleteTagFromTask godoc
// @Summary      タスクからタグを解除
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id      path  int  true  "タスクID"
// @Param        tag_id  path  int  true  "タグID"
// @Success      200 {object} handler.EmptyResponse
// @Failure      404 {object} handler.ErrResponse
// @Router       /tasks/{id}/tags/{tag_id} [delete]
func (dt *DeleteTagFromTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    taskIDStr := chi.URLParam(r, "task_id")
    tagIDStr := chi.URLParam(r, "tag_id")

    taskInt, err := strconv.ParseInt(taskIDStr, 10, 64)
    if err != nil {
        RespondJSON(ctx, w, &ErrResponse{Message: "invalid task ID"}, http.StatusBadRequest)
        return
    }
    tagInt, err := strconv.ParseInt(tagIDStr, 10, 64)
    if err != nil {
        RespondJSON(ctx, w, &ErrResponse{Message: "invalid tag ID"}, http.StatusBadRequest)
        return
    }

	userID := int64(1) // TODO 認証実装後に取得

    t := &entity.TaskTag{
        TaskID: entity.TaskID(taskInt),
        TagID:  entity.TagID(tagInt),
    }

    if err := dt.Repo.DeleteTagFromTask(ctx, dt.DB, userID, t); err != nil {
        if errors.Is(err, store.ErrNotFound) {
            RespondJSON(ctx, w, &ErrResponse{Message: "association not found"}, http.StatusNotFound)
            return
        }
        RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
        return
    }

    RespondJSON(ctx, w, struct{}{}, http.StatusOK)
}
