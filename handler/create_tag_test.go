package handler

import (
    "context"
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-playground/validator/v10"

    "github.com/go-webapi-team/task-app/auth"
    "github.com/go-webapi-team/task-app/entity"
    "github.com/go-webapi-team/task-app/store"

    "github.com/budougumi0617/go_todo_app/testutil"
)

// fakeCreator は Repository.CreateTag のモック実装
type fakeCreator struct{ nextID int64 }

func (f *fakeCreator) CreateTag(_ context.Context, _ store.Execer, t *entity.Tag) error {
	f.nextID++
	t.ID = entity.TagID(f.nextID)
	return nil
}

func TestCreateTag(t *testing.T) {
    type want struct {
        status  int
        rspFile string
    }
    tests := map[string]struct {
        reqFile string
        want    want
    }{
        "ok": {
            reqFile: "testdata/create_tag/ok_req.json.golden",
            want: want{
                status:  http.StatusOK,
                rspFile: "testdata/create_tag/ok_rsp.json.golden",
            },
        },
        "badRequest": {
            reqFile: "testdata/create_tag/bad_req.json.golden",
            want: want{
                status:  http.StatusBadRequest,
                rspFile: "testdata/create_tag/bad_rsp.json.golden",
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
                http.MethodPost,
                "/tags", 
                bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
            )

            // ★認証済みユーザーを Context に注入
            r = r.WithContext(auth.WithUserID(r.Context(), 1))

            // テスト対象となるハンドラ(CreateTag)の用意
            sut := CreateTag{
                Repo: &fakeCreator{},     // モックリポジトリ
                DB: nil,               // db 引数はモックで使わないので nil OK
                Validator: validator.New(),
            }

            // 実行
            sut.ServeHTTP(w, r)

            // 結果を検証
            resp := w.Result()
            testutil.AssertResponse(t,
                resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
            )
        })
    }
}
