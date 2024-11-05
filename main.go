package main

import (
	"log"
	"zion/internal/config"
	"zion/internal/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	// Create a new server
	s, err := server.InitializeZionServer(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Start the server
	s.Start()
}
