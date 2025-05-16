package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
)

type SignupHandler struct {
	Repo      *store.Repository
	DB        store.Execer
	Validator *validator.Validate
}

func (sh *SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var in struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	if err := sh.Validator.Struct(in); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: "Invalid input"}, http.StatusBadRequest)
		log.Printf("validation error: %v", err)
		return
	}

	user := &entity.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	}
	if err := sh.Repo.CreateUser(ctx, sh.DB, user); err != nil {
		//RespondJSON(ctx, w, &ErrResponse{Message: "Failed to create user"}, http.StatusInternalServerError)
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		log.Printf("database error: %v", err)
		return
	}

	RespondJSON(ctx, w, IDResponse{ID: int64(user.ID)}, http.StatusCreated)
}
