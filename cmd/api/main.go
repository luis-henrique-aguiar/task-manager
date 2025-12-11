package main

import (
	"log"
	"net/http"

	"github.com/luis-henrique-aguiar/task-manager/internal/data"
)

type application struct {
	appName string
	logger  *log.Logger
    tasks   *data.TaskModel
}

func main() {

    logger := log.Default()

    app := &application{
        appName: "TaskManagerAPI",
        logger:  logger,
        tasks:   data.NewTaskModel(),
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