package handler

import (
	"fmt"
	"net/http"
	"task-app/sessions" // Ensure the sessions package is properly imported
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		sessions.SessionMutex.Lock()
		delete(sessions.Sessions, cookie.Value)
		if len(sessions.Sessions) == 0 {
			fmt.Println("logoutHandler: current session ID is nothing")
		}
		sessions.SessionMutex.Unlock()
		http.SetCookie(w, &http.Cookie{
			Name:   "session_id",
			Value:  "",
			Path:   "/",
			MaxAge: -1, // クッキーを削除する設定
		})
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
