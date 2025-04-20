package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-webapi-team/task-app/entity"
	"github.com/go-webapi-team/task-app/store"
	"github.com/go-playground/validator/v10"
)

type CreateTag struct {
	Store *store.TagStore
	Validator *validator.Validate
}

func (ct *CreateTag) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Name string `json:"name" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	err := validator.New().Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	t := &entity.Tag{
		Name:      b.Name,
		UserID:    1, // TODO: ユーザIDを取得する
		CreatedAt:  time.Now(),
	}
	idPtr,err := ct.Store.Create(t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		ID int `json:"id"`
	}{
		ID: int(*idPtr),
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)

}