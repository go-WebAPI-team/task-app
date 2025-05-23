package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
)

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer, userID int64, f store.ListTaskFilter) (entity.Tasks, error)
}

// ハンドラ本体
type ListTask struct {
	Repo TaskLister
	DB   store.Queryer
}

// 外部公開用 DTO（新規追加）
type TaskDTO struct {
	ID       int        `json:"id"`
	Title    string     `json:"title"`
	IsDone   bool       `json:"is_done"`
	Priority int        `json:"priority"`
	Description string  `json:"description,omitempty"`
	Deadline    *time.Time `json:"deadline,omitempty"`
}

// ListTask godoc
// @Summary      タスク一覧取得
// @Description  ログインユーザのタスクをフィルタ付きで取得します
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        is_done  query   bool   false  "完了済みフラグ"
// @Param        tag_id   query   int    false  "タグ ID で絞り込み"
// @Success      200  {array}   entity.Task
// @Failure      400  {object}  handler.ErrResponse
// @Router       /tasks [get]
func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query()
	var filter store.ListTaskFilter

	if v := q.Get("is_done"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			RespondJSON(ctx, w, &ErrResponse{Message: "invalid is_done"}, http.StatusBadRequest)
			return
		}
		filter.IsDone = &b
	}
	if v := q.Get("tag_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64) //10進数として解釈し、int64 に収める(DBのPK がBIGINT(64bit) のため)
		if err != nil {
			RespondJSON(ctx, w, &ErrResponse{Message: "invalid tag_id"}, http.StatusBadRequest)
			return
		}
		filter.TagID = &id
	}
	if v := q.Get("due"); v != "" {
		switch v {
		case "asc":
			t := true
			filter.DueAsc = &t
		case "desc":
			f := false
			filter.DueAsc = &f
		default:
			RespondJSON(ctx, w, &ErrResponse{Message: "invalid due"}, http.StatusBadRequest)
			return
		}
	}

	// ------------------------------
	// 認証ユーザー取得：認証チェック済み ctx から userID 抽出  
	// ------------------------------
    userID, ok := auth.GetUserID(ctx)
    if !ok {
        RespondJSON(ctx, w, &ErrResponse{Message: "unauthorized"}, http.StatusUnauthorized)
        return
    }

	ts, err := lt.Repo.ListTasks(ctx, lt.DB, userID, filter)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	// entity.Tasks -> []TaskDTO へ変換
	out := make([]TaskDTO, 0, len(ts))
	for _, t := range ts {
		dto := TaskDTO{
			ID:       int(t.ID),
			Title:    t.Title,
			IsDone:   t.IsDone,
			Priority: t.Priority,
		}
		if t.Description != "" {
			dto.Description = t.Description
		}
		if t.Deadline != nil {
			dto.Deadline = t.Deadline
		}
		out = append(out, dto)
	}

	RespondJSON(ctx, w, out, http.StatusOK)
}
