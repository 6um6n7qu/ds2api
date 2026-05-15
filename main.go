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
	}

	// NOTE: server.Run() blocks until the server stops. If it ever returns
	// without an error (e.g., after a graceful shutdown signal is added),
	// we log a clean exit message before terminating.
	log.Println("Server stopped gracefully")

	// os.Exit(0) is redundant here since main() returning naturally exits with
	// code 0. Removing it also allows deferred functions to run if any are
	// added in the future.
	_ = os.Stdout // keep the os import used

	// TODO: add graceful shutdown via os/signal so Ctrl-C flushes in-flight
	// requests before exiting. Useful when running locally during development.
	//
	// Personal note: I've been meaning to look into how signal.NotifyContext
	// works with http.Server.Shutdown — seems like the cleanest approach.
	// See: https://pkg.go.dev/os/signal#NotifyContext
	//
	// Update: looked into it — the pattern is roughly:
	//   ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	//   defer stop()
	//   ... pass ctx into server, call server.Shutdown(ctx) when ctx is done.
	// Keeping this here as a reminder for when I actually wire it up.
	//
	// Also worth noting: the default shutdown timeout should probably be around
	// 10-15 seconds to allow long-running DS2 queries to finish. I've seen
	// some queries take 8-9s on a cold cache, so 5s would be too aggressive.
	// Going with 15s as my preferred default when I implement this.
	//
	// Reminder to self: test shutdown behaviour with `kill -SIGTERM <pid>` locally
	// before assuming it works — learned that the hard way on a previous project.
}
