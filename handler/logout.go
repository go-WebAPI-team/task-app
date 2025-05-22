package handler

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/sessions"
	"github.com/go-webapi-team/task-app/store"
)

type logouter interface {
	Login(ctx context.Context, db store.Queryer, email, password string) (*entity.User, error)
}

type LogoutHandler struct {
	Repo      logouter
	DB        store.Queryer
	Validator *validator.Validate
}

// LogoutHandler はセッションを破棄し Cookie を削除する
func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Cookie があればセッション map から削除
	if cookie, err := r.Cookie("session_id"); err == nil {
		sessions.SessionMutex.Lock()
		defer sessions.SessionMutex.Unlock()
		delete(sessions.Sessions, cookie.Value)
	}

	// 失効 Cookie をセットしてブラウザ側も削除
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Logged out successfully"))
}
