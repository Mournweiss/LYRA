package main

import (
	"log"
	"github.com/lyra/api-gateway/internal"
)

func main() {
	cfg, err := internal.LoadConfig()
	
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	internal.StartServer(cfg)
}
