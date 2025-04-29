package handler

import (
	"fmt"
	"net/http"
	"task-app/sessions" // Removed unused import
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := sessions.CreateUniqueSessionID()

	sessions.SessionMutex.Lock()
	sessions.Sessions[sessionID] = true // Removed unused reference to sessions
	for id := range sessions.Sessions {
		fmt.Println("loginHandler: Current sessionID is", id)
	}
	sessions.SessionMutex.Unlock()
	sessions.SetCookie(w, sessionID)
	w.Write([]byte("Logged in successfully"))
}
