package handler

import (
	"fmt"
	"net/http"
	"task-app/sessions"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// クッキーからセッションIDを取得
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sessions.SessionMutex.Lock()
	authenticated, ok := sessions.Sessions[cookie.Value]
	if ok {
		fmt.Println("get session id is", cookie.Value)
	}
	sessions.SessionMutex.Unlock()

	// 認証されていない場合はエラーを返す
	if !ok || !authenticated {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Write([]byte("Welcome to your dashboard!"))
}
