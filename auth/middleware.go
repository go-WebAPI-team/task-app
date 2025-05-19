package auth

import (
	"log"
	"net/http"

	"github.com/go-webapi-team/task-app/sessions"
)

// 認証チェック Middleware: Cookie を検証し、正当なら Context に userID を注入する
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		sessions.SessionMutex.Lock()
		uid, ok := sessions.Sessions[cookie.Value]
		log.Printf("Session check performed for session ID: %s", cookie.Value)
		sessions.SessionMutex.Unlock()
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		log.Printf("Session ID from Cookie: %s", cookie.Value)

		ctx := WithUserID(r.Context(), uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
