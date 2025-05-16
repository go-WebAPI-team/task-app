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

type fakeTaskDeleter struct{}

func (f *fakeTaskDeleter) DeleteTask(_ context.Context, _ store.Execer, _ int64, _ entity.TaskID) error {
	return nil
}

func TestDeleteTask_Ok(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(auth.WithUserID(r.Context(), 1))

	sut := DeleteTask{Repo: &fakeTaskDeleter{}, DB: nil}
	sut.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Result().StatusCode)
	}
}
