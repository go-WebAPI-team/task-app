package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
	"github.com/go-webapi-team/task-app/testutil"
)

type fakeTaskLister struct{ ret entity.Tasks }

func (f *fakeTaskLister) ListTasks(_ context.Context, _ store.Queryer, _ int64,
	_ store.ListTaskFilter) (entity.Tasks, error) {
	return f.ret, nil
}

func TestListTask(t *testing.T) {
	okTasks := entity.Tasks{
	    {ID: 1, UserID: 1, Title: "task1"}, // CreatedAt/UpdatedAt = zero値
	    {ID: 2, UserID: 1, Title: "task2"},
	}
	tests := map[string]struct {
		tasks   entity.Tasks
		wantSt  int
		wantRsp string
	}{
		"ok":   {okTasks, http.StatusOK, "testdata/list_task/ok_rsp.json.golden"},
		"empty": {entity.Tasks{}, http.StatusOK, "testdata/list_task/empty_rsp.json.golden"},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			r = r.WithContext(auth.WithUserID(r.Context(), 1))

			sut := ListTask{
				Repo: &fakeTaskLister{ret: tt.tasks},
				DB:   nil,
			}
			sut.ServeHTTP(w, r)
			testutil.AssertResponse(t, w.Result(), tt.wantSt,
				testutil.LoadFile(t, tt.wantRsp))
		})
	}
}
