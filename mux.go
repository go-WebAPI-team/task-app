package main

import (
	"database/sql"
	"net/http"

	"github.com/go-webapi-team/task-app/auth"
	"github.com/go-webapi-team/task-app/handler"
	"github.com/go-webapi-team/task-app/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	// Swagger-UI
	_ "github.com/go-webapi-team/task-app/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// NewMuxは、どのようなハンドラーの実装をどんなURLパスで公開するかルーティングする
func NewMux(db *sql.DB, repo *store.Repository) http.Handler {
	mux := chi.NewRouter()
	v := validator.New()

	// ------------------------------
	// 非認証エンドポイント
	// ------------------------------
	mux.Post("/login", (&handler.LoginHandler{Repo: repo, DB: db, Validator: v}).ServeHTTP)
	mux.Post("/logout", handler.LogoutHandler)
	mux.Post("/signup", (&handler.SignupHandler{Repo: repo, DB: db, Validator: v}).ServeHTTP)
	// swagger & health は公開
	mux.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json"))) // 生成された spec のパス
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	// ------------------------------
	// 認証必須エンドポイント
	// ------------------------------
	mux.Group(func(r chi.Router) {
		r.Use(auth.Middleware)

		// Tag
		createTag := &handler.CreateTag{Repo: repo, DB: db, Validator: v}
		r.Post("/tags", createTag.ServeHTTP)

		listTag := &handler.ListTag{Repo: repo, DB: db}
		r.Get("/tags", listTag.ServeHTTP)

		deleteTag := &handler.DeleteTag{Repo: repo, DB: db}
		r.Delete("/tags/{id}", deleteTag.ServeHTTP)

		// Task
		addTask := &handler.AddTask{Repo: repo, DB: db, Validator: v}
		r.Post("/tasks", addTask.ServeHTTP)

		listTask := &handler.ListTask{Repo: repo, DB: db}
		r.Get("/tasks", listTask.ServeHTTP)

		getTask := &handler.GetTask{Repo: repo, DB: db}
		r.Get("/tasks/{id}", getTask.ServeHTTP)

		updateTask := &handler.UpdateTask{Repo: repo, DB: db, Validator: v}
		r.Put("/tasks/{id}", updateTask.ServeHTTP)

		deleteTask := &handler.DeleteTask{Repo: repo, DB: db}
		r.Delete("/tasks/{id}", deleteTask.ServeHTTP)

		toggleComplete := &handler.ToggleCompleteTask{Repo: repo, DB: db}
		r.Patch("/tasks/{id}/complete", toggleComplete.ServeHTTP)

		addTagToTask := &handler.AddTagToTask{Repo: repo, DB: db}
		r.Put("/tasks/{task_id}/tags/{tag_id}", addTagToTask.ServeHTTP)

		deleteTagFromTask := &handler.DeleteTagFromTask{Repo: repo, DB: db}
		r.Delete("/tasks/{task_id}/tags/{tag_id}", deleteTagFromTask.ServeHTTP)
	})

	// ------------------------------
	// 静的ファイル (フロント確認用)
	// ------------------------------
	mux.Handle("/*", http.FileServer(http.Dir("./public")))

	return mux
}
