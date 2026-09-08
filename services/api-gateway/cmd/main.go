package main

import (
	"encoding/json"
	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/config"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := HealthResponse{
		Status:  "Ok",
		Service: "api-gateway",
	}
	json.NewEncoder(w).Encode(response)

}
func main() {
	// Load configs
	_, err := config.Load()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}

}
