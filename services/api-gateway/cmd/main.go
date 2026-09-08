package main

import (
	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/config"
)

func main() {
	_, err := config.Load()
	if err != nil {
		panic(err)
	}
}
