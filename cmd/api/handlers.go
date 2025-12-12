package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	}

	err := app.tasks.Insert(task)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusCreated, task)
}

func (app *application) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := app.tasks.GetAll()
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, tasks)
}

func (app *application) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	task, err := app.tasks.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Erro interno", http.StatusInternalServerError)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, task)
}

func (app *application) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = app.tasks.Delete(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Erro interno", http.StatusInternalServerError)
		}
		return
	}

	app.writeJSON(w, http.StatusNoContent, map[string]string{"message": "tarefa deletada com sucesso"})
}

func (app *application) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	task, err := app.tasks.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Erro interno", http.StatusInternalServerError)
		}
		return
	}

	var input struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
		Done 	*bool 	`json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Content != nil {
		task.Content = *input.Content
	}
	if input.Done != nil {
		task.Done = *input.Done
	}

	err = app.tasks.Update(task)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Conflito de edição: Tente novamente", http.StatusConflict)
		} else {
			app.logger.Println(err)
			http.Error(w, "Erro interno", http.StatusInternalServerError)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, task)
}

func (app *application) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `jsosn:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	user := &data.User{
		Name:  input.Name,
		Email: strings.ToLower(input.Email),
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		http.Error(w, "Erro ao processar senha", http.StatusInternalServerError)
		return
	}

	err = app.users.Insert(user)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			http.Error(w, "Este email já está em uso", http.StatusConflict)
		} else {
			app.logger.Println(err)
			http.Error(w, "Erro interno", http.StatusInternalServerError)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, user)
}

func (app *application) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := app.users.GetByEmail(strings.ToLower(input.Email))
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		} else {
			app.logger.Println(err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
		}
		return
	}

	matches, err := user.Password.Matches(input.Password)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	if !matches {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iss": "task-manager-api",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		app.logger.Println("Erro ao assinar token:", err)
		http.Error(w, "Erro interno ao gerar token", http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]string{
		"token": tokenString,
	})
}