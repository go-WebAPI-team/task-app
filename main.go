package main

import (
	"fmt"
	"net/http"
	"task-app/handler"
)

func main() {

	fmt.Println("Starting the server!")
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)
	http.HandleFunc("/dashboard", handler.DashboardHandler)

	// 8000番ポートでサーバを開始
	http.ListenAndServe(":8080", nil)
}
