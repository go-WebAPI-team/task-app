package handler

import (
	"fmt"
	"net/http"
	"task-app/sessions"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	sessions.SessionMutex.Lock()
	defer sessions.SessionMutex.Unlock()
	if err == nil {
		delete(sessions.Sessions, cookie.Value)
		if len(sessions.Sessions) == 0 {
			fmt.Println("logoutHandler: current session ID is nothing")
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1, //cookie削除
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
