package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type application struct {
	appName string
	logger  *log.Logger
}

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

func main() {

    logger := log.Default()

    app := &application{
        appName: "TaskManagerAPI",
        logger:  logger,
    }

    logger.Println("Servidor iniciado na porta 8080...")

    srv := &http.Server{
        Addr: ":8080",
        Handler: app.routes(),
    }

    if err := srv.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}