package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
	"github.com/go-webapi-team/task-app/testutil"
)

type fakeLister struct{ ret entity.Tags }

func (f *fakeLister) ListTags(_ context.Context, _ store.Queryer) (entity.Tags, error) {
	return f.ret, nil
}

func TestListTag(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		tags entity.Tags
		want want
	}{
		"ok": {
			tags: entity.Tags{
				{
					ID:   1,
					Name: "tag1",
				},
				{
					ID:   2,
					Name: "tag2",
				},
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_tag/ok_rsp.json.golden",
			
			},
		},
		"empty": {
			tags: entity.Tags{},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_tag/empty_rsp.json.golden",
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			// テスト用のHTTPレスポンスレコーダとリクエスト作成
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodGet,
				"/tags",
				nil,
			)

			// テスト対象となるハンドラ(ListTag)の用意
			sut := ListTag{
				Repo: &fakeLister{ret: tt.tags},
				DB:   nil,
			}

			// リクエストをハンドラに渡す
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t,
				resp,
				tt.want.status, 
				testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}

