package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
)

// 一覧取得ユースケースのインタフェース
type TagLister interface {
	ListTags(ctx context.Context, db store.Queryer) (entity.Tags, error)
}

type ListTag struct {
	Repo TagLister
	DB   store.Queryer
}

// DTO（Data Transfer Object）
// ドメインモデル（entity.Tag）と外部公開の（クライアントが読むレスポンス JSONの契約）を分離する
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// CreatedAt *は要件があるときだけ追加*
    CreatedAt *time.Time `json:"created,omitempty"`
}

// ListTag godoc
// @Summary      タグ一覧取得
// @Tags         tags
// @Accept       json
// @Produce      json
// @Success      200 {array} entity.Tag
// @Router       /tags [get]
func (lt *ListTag) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tags, err := lt.Repo.ListTags(ctx, lt.DB)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()},
			http.StatusInternalServerError)
		return
	}

	rsp := make([]Tag, 0) // 要素 0 だが non‑nil → JSON では [] になる
	for _, t := range tags {
		var created *time.Time
		if !t.CreatedAt.IsZero() {
			c := t.CreatedAt
			created = &c
		}
		rsp = append(rsp, Tag{
			ID:   int(t.ID),
			Name: t.Name,
			CreatedAt: created,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
