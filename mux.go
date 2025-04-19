package main

import (
		"net/http"

		"github.com/go-webapi-team/task-app/handler"
		"github.com/go-webapi-team/task-app/store"

		"github.com/go-chi/chi/v5"
		"github.com/go-playground/validator/v10"
)

// NewMuxは、どのようなハンドラーの実装をどんなURLパスで公開するかルーティングする
func NewMux() http.Handler {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析(実装予定)のエラーを回避するため明示的に戻り値を捨てる
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	v := validator.New()
	mux.Handle("/tags", &handler.CreateTag{Store: store.Tags, Validator: v})
	ct := &handler.CreateTag{Store: store.Tags, Validator: v}
	mux.Post("/tags", ct.ServeHTTP)
	lt := &handler.ListTag{Store: store.Tags}
	mux.Get("/tags", lt.ServeHTTP)
	return mux
}