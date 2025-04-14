package main

import (
	"fmt"
	"net/http"

	"TASK-APP/handlers"
)

func main() {
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	fmt.Println("Server running on http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}