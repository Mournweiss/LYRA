package main

import (
	"log"
	"github.com/lyra/api-gateway/internal/config"
	"github.com/lyra/api-gateway/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	server.StartServer(cfg)
}
