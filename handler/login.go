package handler

import (
	"net/http"
	"github.com/go-webapi-team/task-app/sessions"
)

// LoginHandler は簡易ログイン（ダミー認証）を行い Cookie を発行する
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// ToDo:本来はフォーム値を DB 認証する(後続開発でメール・パスワードを検証 → userID を取得する実装に置換)
	const dummyUserID int64 = 2

	// セッション生成
	sessionID := sessions.NewSession(dummyUserID)
	sessions.SetCookie(w, sessionID)

	w.Write([]byte("Logged in successfully"))
}
