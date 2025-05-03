package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"sync"
)

var (
	Sessions     = make(map[string]bool)
	SessionMutex = &sync.Mutex{}
)

func CreateUniqueSessionID() string {
	for {
		sessionID := generateSessionId()
		SessionMutex.Lock()
		if !Sessions[sessionID] {
			Sessions[sessionID] = true
			SessionMutex.Unlock()
			return sessionID
		}
		SessionMutex.Unlock()
	}
}

func generateSessionId() string {
	session := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, session); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(session)
}

func SetCookie(w http.ResponseWriter, sessionID string) {
	cookie := http.Cookie{
		Name:   "session_id",
		Value:  sessionID,
		Path:   "/",
		MaxAge: 600,
	}
	http.SetCookie(w, &cookie)
}
