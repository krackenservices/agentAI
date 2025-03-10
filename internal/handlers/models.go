package handlers

import (
	"encoding/json"
	"krackenservices.com/agentAI/internal/config"
	"net/http"
	"strings"
)

// GetModel godoc
// @Summary Get model details
// @Description Returns details of the model identified by modelID.
// @Tags models
// @Accept json
// @Produce json
// @Param modelID path string true "Model ID"
// @Success 200 {object} config.ModelConfig
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /model/{modelID} [get]
func GetModel(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Expect URL: /model/<modelID>
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 || parts[2] == "" {
			http.Error(w, "Model ID not provided", http.StatusBadRequest)
			return
		}
		modelID := parts[2]
		for _, model := range cfg.Models {
			if model.ID == modelID {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(model)
				return
			}
		}
		http.Error(w, "Model not found", http.StatusNotFound)
	}
}

// ListModels godoc
// @Summary List all models
// @Description Returns a list of all models with their details.
// @Tags models
// @Accept json
// @Produce json
// @Success 200 {array} config.ModelConfig
// @Router /models [get]
func ListModels(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg.Models)
	}
}
