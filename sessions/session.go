
package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	Sessions     = make(map[string]int64)
	SessionMutex = &sync.Mutex{}
)

// ------------------------------
// 公開 API
// ------------------------------

// NewSession はユニークなセッション ID を生成して userID と紐付ける   // NEW
func NewSession(userID int64) string {
	for {
		sid := generateSessionID()
		SessionMutex.Lock()
		if _, ok := Sessions[sid]; !ok {
			Sessions[sid] = userID
			SessionMutex.Unlock()
			return sid
		}
		SessionMutex.Unlock()
	}
}

// SetCookie はレスポンスに Cookie を設定する
func SetCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   600,             // 10 分
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true,           // HTTPS 運用時に有効化
	}
	http.SetCookie(w, cookie)
}

// ------------------------------
// 内部ユーティリティ
// ------------------------------

func generateSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// For future: map のクリーンアップ goroutine (期限切れ削除) など
func init() {
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		for range ticker.C {
			// TTL 管理を導入した際に GC する想定
		}
	}()
}