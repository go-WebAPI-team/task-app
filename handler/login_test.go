package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/login", nil)
	w := httptest.NewRecorder()
	LoginHandler(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected session_id to have a value, but it was empty")
	}

	cookies := resp.Cookies()
	found := false

	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			found = true
			if cookie.Value == "" {
				t.Errorf("expected session_id to have")
			}
		}
	}
	if !found {
		t.Errorf("expected session_id cookie, but it was not found")
	}

}
