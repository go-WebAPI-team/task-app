package handler

import (
	"net/http"

	"github.com/go-webapi-team/task-app/sessions"
)

// LogoutHandler はセッションを破棄し Cookie を削除する
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Cookie があればセッション map から削除
	if cookie, err := r.Cookie("session_id"); err == nil {
		sessions.SessionMutex.Lock()
		delete(sessions.Sessions, cookie.Value)
		sessions.SessionMutex.Unlock()
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
