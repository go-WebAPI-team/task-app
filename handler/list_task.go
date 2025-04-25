package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

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

	// TODO: 認証機能実装後にログインユーザーの ID を ctx から取得する
	ts, err := lt.Repo.ListTasks(ctx, lt.DB, 1, filter)
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
