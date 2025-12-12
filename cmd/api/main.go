package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/luis-henrique-aguiar/task-manager/internal/data"
)

type config struct {
    port string
    db   struct {
          dsn string
    }
}

type application struct {
	appName string
    cfg     config
	logger  *log.Logger
    tasks   *data.TaskModel
}

func main() {
    logger := log.Default()

    if err := godotenv.Load(); err != nil {
        logger.Println("Aviso: arquivo .env não encontrado, usando variáveis de ambiente do sistema")
    }

    var cfg config

    cfg.port = os.Getenv("PORT")
    if cfg.port == "" {
        cfg.port = "8080"
    }

    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    cfg.db.dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        dbUser, dbPass, dbHost, dbPort, dbName)

    db, err := sql.Open("postgres", cfg.db.dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err = db.Ping(); err != nil {
        log.Fatal("Erro ao conectar no Postgres:", err)
    }
    logger.Println("Conectado ao Postgres com sucesso!")

    createTable(db, logger)

    app := &application{
        appName: "TaskManagerAPI",
        logger:  logger,
        tasks:   data.NewTaskModel(db),
        cfg:     cfg,
    }

    addr := fmt.Sprintf(":%s", cfg.port)

    logger.Printf("Servidor iniciado na porta %s...\n", cfg.port)
    srv := &http.Server{
        Addr: addr,
        Handler: app.routes(),
    }

    if err := srv.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

func createTable(db *sql.DB, logger *log.Logger) {
    query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        content TEXT,
        done BOOLEAN DEFAULT FALSE,
        version INTEGER DEFAULT 1,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );`

    if _, err := db.Exec(query); err != nil {
        logger.Fatal("Erro ao criar tabela:", err)
    }
}