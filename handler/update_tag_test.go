package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
)

type fakeTaskUpdater struct{}

func (f *fakeTaskUpdater) UpdateTask(_ context.Context, _ store.Execer, _ *entity.Task) error {
	return nil
}

func TestUpdateTask_Ok(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/tasks/1",
		bytes.NewReader([]byte(`{"title":"upd","priority":2,"is_done":false}`)))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	sut := UpdateTask{Repo: &fakeTaskUpdater{}, DB: nil, Validator: validator.New()}
	sut.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Result().StatusCode)
	}
}
