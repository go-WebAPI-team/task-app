package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/handler"
	"github.com/go-webapi-team/task-app/sessions"
	"github.com/go-webapi-team/task-app/store"
)

type Loginer struct{}

func (m Loginer) Login(ctx context.Context, execer store.Execer, username, password string) (*entity.User, error) {
	return &entity.User{ID: 1, Name: username}, nil
}

// TestAuthMiddlewareOK は login → 保護エンドポイントへの一連の流れを確認する
func TestAuthMiddlewareOK(t *testing.T) {
	v := validator.New()
	loginHandler := &handler.LoginHandler{
		Repo:      Loginer{},
		DB:        nil,
		Validator: v,
	}
	// 1) login サーバ
	loginSrv := httptest.NewServer(http.HandlerFunc(loginHandler.ServeHTTP))
	defer loginSrv.Close()

	// 2) protected サーバ
	protected := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := auth.GetUserID(r.Context()); !ok {
			t.Fatalf("userID missing")
		}
		w.WriteHeader(http.StatusOK)
	}))
	protectedSrv := httptest.NewServer(protected)
	defer protectedSrv.Close()

	// --- login
	resp, err := http.Post(loginSrv.URL, "text/plain", nil)
	if err != nil {
		t.Fatalf("login request error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected login status %d", resp.StatusCode)
	}
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		t.Fatalf("cookie not set")
	}

	// --- access protected
	req, _ := http.NewRequest(http.MethodGet, protectedSrv.URL, nil)
	req.AddCookie(cookies[0])
	got, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("protected request error: %v", err)
	}
	if got.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", got.StatusCode)
	}
}

// TestAuthMiddlewareUnauthorized は Cookie なしの場合 401 になることを保証
func TestAuthMiddlewareUnauthorized(t *testing.T) {
	protected := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // 通ったら失敗
	}))
	srv := httptest.NewServer(protected)
	defer srv.Close()

	got, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if got.StatusCode != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", got.StatusCode)
	}
}

// --- 補足: セッション map を掃除しておく ---
func init() {
	sessions.Sessions = make(map[string]int64)
}
