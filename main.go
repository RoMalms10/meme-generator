package main

import (
	"log"

	"github.com/RoMalms10/meme-generator/config"
	"github.com/RoMalms10/meme-generator/server"
)

func main() {
	// Load configuration from environment variables
	cfg := config.LoadConfig()

	// Create and start the server
	srv := server.NewServer(cfg)

	log.Printf("Starting Meme Generator microservice on port %s", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
