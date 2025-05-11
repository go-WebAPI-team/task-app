package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
	"github.com/go-playground/validator/v10"
)

type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type AddTask struct {
	Repo      TaskAdder
	DB        store.Execer
	Validator *validator.Validate
}

// AddTask godoc
// @Summary      タスクを新規作成
// @Description  JSON で受け取ったタスク情報を保存します
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body     entity.Task  true  "新規タスク"
// @Success      200   {object}  handler.IDResponse
// @Failure      400   {object}  handler.ErrResponse
// @Failure      500   {object}  handler.ErrResponse
// @Router       /tasks [post]
func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var in struct {
		Title       string     `json:"title"       validate:"required,max=255"`
		Description string     `json:"description"`
		Deadline    *time.Time `json:"deadline"`
		Priority    int        `json:"priority"    validate:"oneof=1 2 3"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := at.Validator.Struct(in); err != nil {
		// バリデーションエラーの詳細をレスポンスに含める
		msg := err.Error()
		RespondJSON(ctx, w, &ErrResponse{Message: msg}, http.StatusBadRequest)
		log.Printf("validation error: %v", err)
		return
	}

	// ------------------------------
	// 認証ユーザー取得：認証チェック済み ctx から userID 抽出  
	// ------------------------------
    userID, ok := auth.GetUserID(ctx)
    if !ok {
        RespondJSON(ctx, w, &ErrResponse{Message: "unauthorized"}, http.StatusUnauthorized)
        return
    }

	t := &entity.Task{
		// TODO: 認証機能実装後にログインユーザーの ID を ctx から取得する
		UserID:      userID, 
		Title:       in.Title,
		Description: in.Description,
		Deadline:    in.Deadline,
		Priority:    in.Priority,
		IsDone:      false,
	}

	if err := at.Repo.AddTask(ctx, at.DB, t); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "タスクの追加中にエラーが発生しました"}, http.StatusInternalServerError)
		log.Printf("database error: %v", err)
		return
	}
	RespondJSON(ctx, w, IDResponse{ID: int64(t.ID)}, http.StatusOK)
}
