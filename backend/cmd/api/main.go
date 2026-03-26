package main

import (
	"log"
	"net/http"

	"github.com/akshzyx/gorum/internal/config"
	"github.com/akshzyx/gorum/internal/container"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Build container
	c, err := container.NewContainer(cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer c.DB.Close()

	// Start server
	log.Printf("🚀 Server running on :%s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, c.Router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
