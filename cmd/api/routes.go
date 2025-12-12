package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", app.HealthCheck)
	r.Post("/tasks", app.CreateTaskHandler)
	r.Get("/tasks", app.ListTasksHandler)
	r.Delete("/tasks/{id}", app.DeleteTaskHandler)
	r.Get("/tasks/{id}", app.GetTaskHandler)

	return r
}