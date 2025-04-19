package main

import "net/http"

// NewMuxは、どのようなハンドラーの実装をどんなURLパスで公開するかルーティングする
func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析(実装予定)のエラーを回避するため明示的に戻り値を捨てる
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}