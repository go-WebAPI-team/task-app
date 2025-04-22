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

	// RDB 永続化版ハンドラ
	ct := &handler.CreateTag{
		Repo:      repo,
		DB:        db,
		Validator: v,
	}
	mux.Post("/tags", ct.ServeHTTP)

	lt := &handler.ListTag{
		Repo: repo,
		DB:   db,
	}
	mux.Get("/tags", lt.ServeHTTP)

	// 常に使うヘルスチェック
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析(実装予定)のエラーを回避するため明示的に戻り値を捨てる
        _, _ = w.Write([]byte(`{"status": "ok"}`))
    })

	return mux
}