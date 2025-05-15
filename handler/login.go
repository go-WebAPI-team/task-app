package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/sessions"
	"github.com/go-webapi-team/task-app/store"
)

type loginer interface {
	Login(ctx context.Context, db store.Execer, email, password string) (*entity.User, error)
}

type LoginHandler struct {
	Repo      loginer
	DB        store.Execer
	Validator *validator.Validate
}

// LoginHandler は簡易ログイン（ダミー認証）を行い Cookie を発行する
func (lh *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// ToDo:本来はフォーム値を DB 認証する(後続開発でメール・パスワードを検証 → userID を取得する実装に置換)
	//const dummyUserID int64 = 1

	// セッション生成
	//sessionID := sessions.NewSession(dummyUserID)
	//sessions.SetCookie(w, sessionID)

	ctx := r.Context()

	var in struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}
	//hashed, _ := bcrypt.GenerateFromPassword(password, 10)
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	// Removed unnecessary err check block
	if err := lh.Validator.Struct(in); err != nil {
		// バリデーションエラーの詳細をレスポンスに含める
		RespondJSON(ctx, w, &ErrResponse{Message: "Invalid email or password"}, http.StatusBadRequest)
		log.Printf("validation error: %v", err)
		return
	}

	user, err := lh.Repo.Login(ctx, lh.DB, in.Email, in.Password)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "err"}, http.StatusInternalServerError)
		log.Printf("database error: %v", err)
		return
	}
	sessionID := sessions.NewSession(int64(user.ID))
	sessions.SetCookie(w, sessionID)

	RespondJSON(ctx, w, IDResponse{ID: int64(user.ID)}, http.StatusOK)
	log.Printf("User %d logged in successfully", user.ID)

	w.Write([]byte("Logged in successfully"))
}
