package main

import (
	"encoding/json"
	"net/http"

	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/config"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status:  "ok",
		Service: "api-gateway",
	}

	json.NewEncoder(w).Encode(response)

}
func main() {
	// Load configs
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET/health", healthHandler)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		panic(err)
	}

}
