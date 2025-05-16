package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
)

type fakeTaskCompleter struct{}

func (f *fakeTaskCompleter) ToggleTaskDone(_ context.Context, _ store.Execer, _ int64, _ entity.TaskID) error {
	return nil
}

func TestToggleCompleteTask_Ok(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/tasks/1/complete", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(auth.WithUserID(r.Context(), 1))

	sut := ToggleCompleteTask{Repo: &fakeTaskCompleter{}, DB: nil}
	sut.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Result().StatusCode)
	}
}
