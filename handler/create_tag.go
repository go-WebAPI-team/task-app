package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
	"github.com/go-playground/validator/v10"
)

// ユースケース単位のインタフェース
type TagCreator interface {
	CreateTag(ctx context.Context, db store.Execer, t *entity.Tag) error
}

type CreateTag struct {
	Repo TagCreator
	DB store.Execer
	Validator *validator.Validate
}

// CreateTag godoc
// @Summary      タグを新規作成
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        tag  body  entity.Tag  true  "タグ情報"
// @Success      200 {object} handler.IDResponse
// @Failure      400 {object} handler.ErrResponse
// @Router       /tags [post]
func (ct *CreateTag) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------------------------
	// 認証ユーザー取得：認証チェック済み ctx から userID 抽出  
	// ------------------------------
    userID, ok := auth.GetUserID(ctx)
    if !ok {
        RespondJSON(ctx, w, &ErrResponse{Message: "unauthorized"}, http.StatusUnauthorized)
        return
    }

	// ------------------------------
	// リクエストパース
	// ------------------------------
	var b struct {
		Name string `json:"name" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if err := ct.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// ------------------------------
	// ビジネスロジック
	// ------------------------------
	now := time.Now()
	t := &entity.Tag{
		Name:      b.Name,
		UserID:    entity.UserID(userID),
		CreatedAt:  now,
		UpdatedAt: now,
	}

	if err := ct.Repo.CreateTag(ctx, ct.DB, t); err !=  nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// t.ID に自動設定される ※CreateTag が LastInsertId をセット
    rsp := struct{ ID int `json:"id"` }{ID: int(t.ID)}
	RespondJSON(ctx, w, rsp, http.StatusOK)

}