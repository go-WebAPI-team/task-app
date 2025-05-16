package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
	"github.com/go-webapi-team/task-app/testutil"
)

type fakeAdder struct{ nextID int64 }

func (f *fakeAdder) AddTask(_ context.Context, _ store.Execer, t *entity.Task) error {
	f.nextID++
	t.ID = entity.TaskID(f.nextID)
	t.CreatedAt, t.UpdatedAt = time.Now(), time.Now()
	return nil
}

func TestAddTask(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)))
			r = r.WithContext(auth.WithUserID(r.Context(), 1))

			sut := AddTask{
				Repo:      &fakeAdder{},
				DB:        nil,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.status,
				testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
