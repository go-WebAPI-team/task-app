package handlers

import (
	//repositoryフォルダにsqlファイルがあるとする
	"TASK-APP/models"
	"TASK-APP/repository"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

// タスクを作成する
func CreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var task models.Task
	if err := json.Unmarshal(body, &task); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// タスク一覧を取得する
func GetTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userID := query.Get("user_id")

	tasks, err := repository.Gettasks(userID)
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// 特定のタスクを取得する
func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):]
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := repository.GetTask(taskID)
	if err != nil {
		if err.Error() == "task not found" {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

}

// タスクを更新する
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):]
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var task models.Task
	if err := json.Unmarshal(body, &task); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	task.ID = taskID

	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):]
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var task models.Task
	if err := json.Unmarshal(body, &task); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteTask(taskID); err != nil {
		if err.Error() == "task not found" {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})

}

func ToggleCompleteTask(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len("/tasks/"):]
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteTask(taskID); err != nil {
		if err.Error() == "task not found" {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})

}
