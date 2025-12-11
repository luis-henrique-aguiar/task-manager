package main

import (
	"encoding/json"
	"net/http"

	"github.com/luis-henrique-aguiar/task-manager/internal/data"
)

func (app *application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "available",
		"env":     "development",
		"version": "1.0.0",
		"app":     app.appName,
	}

	app.writeJSON(w, http.StatusOK, data)
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		app.logger.Println("Erro ao escrever JSON: ", err)
	}
}

func (app *application) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	task := &data.Task{
		Title:   input.Title,
		Content: input.Content,
		Done:    false,
	}

	app.tasks.Insert(task)

	app.writeJSON(w, http.StatusCreated, task)
}

func (app *application) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := app.tasks.GetAll()
	app.writeJSON(w, http.StatusOK, tasks)
}