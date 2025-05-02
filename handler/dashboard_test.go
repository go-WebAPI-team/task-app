package handler

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"task-app/sessions"
	"testing"
)

func TestDashBoard(t *testing.T) {
	sessions.Sessions = make(map[string]bool)
	sessions.SessionMutex = &sync.Mutex{}

	testSessionID := "test-session-id"
	sessions.Sessions[testSessionID] = true

	tests := []struct {
		name           string
		cookie         *http.Cookie
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid session ID",
			cookie:         &http.Cookie{Name: "session_id", Value: testSessionID},
			expectedStatus: http.StatusOK,
			expectedBody:   "Welcome to your dashboard!",
		},
		{
			name:           "Missing session ID cookie",
			cookie:         nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "unauthorized\n",
		},
		{
			name:           "Invalid session ID",
			cookie:         &http.Cookie{Name: "session_id", Value: "invalid-session-id"},
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Forbidden\n",
		},
	}

	for _, tt := range tests {
		r := httptest.NewRequest("GET", "/dashboard", nil)
		if tt.cookie != nil {
			r.AddCookie(tt.cookie)
		}
		w := httptest.NewRecorder()

		DashboardHandler(w, r)

		if w.Code != tt.expectedStatus {
			t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
		}

		if w.Body.String() != tt.expectedBody {
			t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
		}

	}
}
