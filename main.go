package main

import (
	"zion/internal/config"
	"zion/internal/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	// Create a new server
	srv := server.CreateServer(cfg)
	// Start the server
	srv.Start()
}
