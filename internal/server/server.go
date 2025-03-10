package server

import (
	"log"
	"net/http"

	"krackenservices.com/agentAI/internal/config"
	"krackenservices.com/agentAI/internal/routes"
)

// StartServer initializes and starts the HTTP server.
func StartServer(cfg *config.Config) error {
	router := routes.NewRouter(cfg)
	address := ":" + cfg.Server.Port
	log.Printf("Server starting on port %s", cfg.Server.Port)
	return http.ListenAndServe(address, router)
}
