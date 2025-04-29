package handler

import (
	"fmt"
	"net/http"
	"task-app/sessions"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//セッションIDを生成
	sessionID := sessions.CreateUniqueSessionID()

	sessions.SessionMutex.Lock()
	sessions.Sessions[sessionID] = true
	for id := range sessions.Sessions {
		fmt.Println("loginHandler: Current sessionID is", id)
	}
	sessions.SessionMutex.Unlock()

	//クッキーを設定
	sessions.SetCookie(w, sessionID)
	w.Write([]byte("Logged in successfully"))
}
