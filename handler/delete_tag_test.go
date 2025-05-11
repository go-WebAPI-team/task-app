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

// ダミー TagDeleter
type fakeTagDeleter struct{}

func (f *fakeTagDeleter) DeleteTag(_ context.Context, _ store.Execer, _ int64, _ entity.TagID) error {
	return nil
}

func TestDeleteTag_OK(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/tags/5", nil)

	// chi の URL パラメータを埋め込む
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "5")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(auth.WithUserID(r.Context(), 1))

	sut := DeleteTag{Repo: &fakeTagDeleter{}, DB: nil}
	sut.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	}
}
