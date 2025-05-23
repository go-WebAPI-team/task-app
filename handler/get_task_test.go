package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
	"github.com/go-webapi-team/task-app/testutil"
)

type fakeTaskGetter struct{ t *entity.Task }

func (f *fakeTaskGetter) GetTask(_ context.Context, _ store.Queryer, _ int64, _ entity.TaskID) (*entity.Task, error) {
	return f.t, nil
}

// Add the missing method to satisfy TaskGetter interface
func (f *fakeTaskGetter) ListTagIDsByTaskID(_ context.Context, _ store.Queryer, _ int64) ([]int64, error) {
	return []int64{}, nil
}

func TestGetTask(t *testing.T) {
	now := time.Now()
	task := &entity.Task{ID: 1, UserID: 1, Title: "task", Priority: 2, CreatedAt: now, UpdatedAt: now}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(auth.WithUserID(r.Context(), 1))

	sut := GetTask{Repo: &fakeTaskGetter{t: task}, DB: nil}
	sut.ServeHTTP(w, r)

	testutil.AssertResponse(t, w.Result(), http.StatusOK,
		testutil.MustJSON(t, struct {
			*entity.Task
			TagIDs []int64 `json:"tag_ids"`
		}{
			Task:   task,
			TagIDs: []int64{},
		}),
	)
}
