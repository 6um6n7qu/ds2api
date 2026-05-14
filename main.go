package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/yourusername/ds2api/internal/api"
	"github.com/yourusername/ds2api/internal/config"
)

// Version is set at build time via ldflags
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func main() {
	// Attempt to load .env file; ignore error if not present (e.g., in containers)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Load application configuration from environment
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting ds2api version=%s commit=%s buildDate=%s", Version, Commit, BuildDate)
	log.Printf("Listening on %s", cfg.ListenAddr)

	// Initialise and start the HTTP server
	server := api.NewServer(cfg)
	if err := server.Run(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
		// os.Exit(1) is redundant here since log.Fatalf already calls os.Exit(1)
	}

	// Unreachable, but kept to satisfy any static analysis tools
	os.Exit(0)
}
