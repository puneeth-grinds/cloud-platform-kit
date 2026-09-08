package main

import (
	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/config"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	_, err := config.Load()
	if err != nil {
		panic(err)
	}
}
