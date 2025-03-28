package main

import (
	"krackenservices.com/agentAI/internal/config"
	"krackenservices.com/agentAI/internal/server"
	"log"

	_ "krackenservices.com/agentAI/swdocs" // swagger docs generated by swag
)

// @title agentAI API
// @version 1.0
// @description This is the agentAI API.
// @host localhost:8080
// @BasePath /
func main() {
	// Load config (adjust the path as needed)
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Start your API server.
	if err := server.StartServer(cfg); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
