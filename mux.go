package main

import (
	"database/sql"
	"net/http"

	"github.com/go-webapi-team/task-app/handler"
	"github.com/go-webapi-team/task-app/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// NewMuxは、どのようなハンドラーの実装をどんなURLパスで公開するかルーティングする
func NewMux(db *sql.DB, repo *store.Repository) http.Handler {
	mux := chi.NewRouter()
	v := validator.New()

	// Tagハンドラ
	createTag := &handler.CreateTag{Repo: repo, DB: db, Validator: v}
	mux.Post("/tags", createTag.ServeHTTP)

	listTag := &handler.ListTag{Repo: repo, DB: db}
	mux.Get("/tags", listTag.ServeHTTP)

	deleteTag := &handler.DeleteTag{Repo: repo, DB: db}
	mux.DELETE("/tags", deleteTag.ServeHTTP)

	// Task ハンドラ
	addTask := &handler.AddTask{Repo: repo, DB: db, Validator: v}
	mux.Post("/tasks", addTask.ServeHTTP)

	listTask := &handler.ListTask{Repo: repo, DB: db}
	mux.Get("/tasks", listTag.ServeHTTP)

	getTask := &handler.GetTask{Repo: repo, DB: db}
	mux.Get("/tasks/{id}", getTask.ServeHTTP)

	updateTask := &handler.UpdateTask{Repo: repo, DB: db, Validator: v}
	mux.Put("/tasks/{id}", updateTask.ServeHTTP)

	deleteTask := &handler.DeleteTask{Repo: repo, DB: db}
	mux.Delete("/tasks/{id}", deleteTask.ServeHTTP)

	toggleComplete := &handler.ToggleCompleteTask{Repo: repo, DB: db}
	mux.Patch("/tasks/{id}/complete", toggleComplete.ServeHTTP)

	addTagToTask := &handler.AddTagToTask{Repo: repo, DB: db,}
	mux.Put("/tasks/{task_id}/tags/{tag_id}", addTagToTask.ServeHTTP)

	deleteTagFromTask := &handler.DeleteTagFromTask{Repo: repo, DB: db}
    mux.Delete("/tasks/{task_id}/tags/{tag_id}", deleteTagFromTask.ServeHTTP)

	// 常に使うヘルスチェック
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析(実装予定)のエラーを回避するため明示的に戻り値を捨てる
        _, _ = w.Write([]byte(`{"status": "ok"}`))
    })

	return mux
}