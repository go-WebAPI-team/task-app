package handler

import (
	"context"
	"net/http"
	"testing"
    "net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
)
type fakeTagAdder struct{}

func (f *fakeTagAdder) AddTagToTask(_ context.Context, _ store.Execer, _ int64, t *entity.TaskTag) error {
	t.ID = entity.TaskTagID(999)
	return nil
}

func TestAddTagToTask_Ok(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tasks/1/tags/2", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("task_id", "1")
	rctx.URLParams.Add("tag_id", "2")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(auth.WithUserID(r.Context(), 1))

	sut := AddTagToTask{Repo: &fakeTagAdder{}, DB: nil}
	sut.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Result().StatusCode)
	}
}
