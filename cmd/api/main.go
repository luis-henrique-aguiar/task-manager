package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

        response := HealthCheckResponse {
            Status:  "ok",
            Message: "Servidor rodando liso no Go!",
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)

        if err := json.NewEncoder(w).Encode(response); err != nil {
            log.Printf("Erro ao fazer encode: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
	})

    log.Println("Servidor rodando na porta 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}