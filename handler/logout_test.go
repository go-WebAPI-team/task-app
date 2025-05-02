package handler

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"task-app/sessions"
	"testing"
)

func TestLogoutHandler(t *testing.T) {
	sessions.Sessions = map[string]bool{"test-session-id": true}
	sessions.SessionMutex = &sync.Mutex{}

	r := httptest.NewRequest(http.MethodPost, "/logout", nil)
	r.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	w := httptest.NewRecorder()

	LogoutHandler(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("ecepted session_id to have a value, but it was empty")
	}

	found := false
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session_id" && cookie.MaxAge == -1 {
			found = true
		}
	}
	if !found {
		t.Errorf("expected session_id cookie to be marked for delction, but it was not")
	}

	sessions.SessionMutex.Lock()
	defer sessions.SessionMutex.Unlock()

	if _, exists := sessions.Sessions["test-session-id"]; exists {
		t.Errorf("expected session id to be deleted, but it still exists")
	}

}
