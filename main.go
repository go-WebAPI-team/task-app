package main

import (
	"TASK-APP/config"
	"fmt"
	"net/http"
	"task-app/handlers"
)

func main() {

	fmt.Println("Starting the server!")

	config.ConnectDB() // ルートとハンドラ関数を定義
	http.HandleFunc("/tasks", handlers.CreateTask)
	http.HandleFunc("/tasks", handlers.GetTasks)
	http.HandleFunc("/tasks/:id", handlers.GetTask)
	http.HandleFunc("/tasks/:id", handlers.UpdateTask)
	http.HandleFunc("/tasks/:id", handlers.DeleteTask)
	http.HandleFunc("/tasks/:id/complete", handlers.ToggleCompleteTask)

	// 8000番ポートでサーバを開始
	http.ListenAndServe(":8080", nil)
}
