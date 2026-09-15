package handler

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// HealthHandler returns a simple JSON response used by humans, local curl
// checks, and load balancer health checks.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status:  "ok",
		Service: "api-gateway",
	}

	json.NewEncoder(w).Encode(response)
}
