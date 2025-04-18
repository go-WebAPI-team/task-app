package handler

import (
	"net/http"
	"time"

	"github.com/go-webapi-team/task-app/store"
)

type ListTag struct {
	Store *store.TagStore
}

// ドメインモデル（entity.Tag）と外部公開のData Transfer Object（クライアントが読むレスポンス JSONの契約）を分離する
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Created *は要件があるときだけ追加*
    Created *time.Time `json:"created,omitempty"`
}

func (lt *ListTag) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tags := lt.Store.GetAll()
	rsp := make([]Tag, 0) // 要素 0 だが non‑nil → JSON では [] になる
	for _, t := range tags {
		var created *time.Time
		if !t.Created.IsZero() { // ← 作成済みの値だけポインタを立てる
			created = &t.Created
		}
		rsp = append(rsp, Tag{
			ID:   int(t.ID),
			Name: t.Name,
			Created: created,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
