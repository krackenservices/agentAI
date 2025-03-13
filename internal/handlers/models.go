package handlers

import (
	"encoding/json"
	"krackenservices.com/agentAI/internal/config"
	"net/http"
	"strings"
)

func maskModel(m config.ModelConfig) config.ModelConfig {
	// Create a new instance with simple fields copied.
	mCopy := config.ModelConfig{
		ID:                        m.ID,
		Name:                      m.Name,
		Endpoint:                  m.Endpoint,
		Enabled:                   m.Enabled,
		APIKey:                    "<masked>", // Will be set below
		AdditionalSystemPrompt:    m.AdditionalSystemPrompt,
		AdditionalUserPrompt:      m.AdditionalUserPrompt,
		AdditionalAssistantPrompt: m.AdditionalAssistantPrompt,
		ToolsSupported:            m.ToolsSupported,
		ToolTagStart:              m.ToolTagStart,
		ToolTagEnd:                m.ToolTagEnd,
	}

	// Deep copy Headers.
	if m.Headers != nil {
		mCopy.Headers = make(map[string]string)
		for k, v := range m.Headers {
			mCopy.Headers[k] = v
		}
	}

	// Deep copy Parameters.
	if m.Parameters != nil {
		mCopy.Parameters = make(map[string]interface{})
		for k, v := range m.Parameters {
			mCopy.Parameters[k] = v
		}
	}

	// Deep copy Tools slice.
	if m.Tools != nil {
		mCopy.Tools = make([]string, len(m.Tools))
		copy(mCopy.Tools, m.Tools)
	}

	return mCopy
}

// ListModels godoc
// @Summary List all models
// @Description Returns a list of all models with their details (API keys masked).
// @Tags models
// @Accept json
// @Produce json
// @Success 200 {array} config.ModelConfig
// @Router /api/v1/models [get]
func ListModels(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maskedModels := make([]config.ModelConfig, len(cfg.Models))
		for i, m := range cfg.Models {
			maskedModels[i] = maskModel(m)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(maskedModels)
	}
}

// GetModel godoc
// @Summary Get model details
// @Description Returns the details of a model by its ID (API key is masked).
// @Tags models
// @Accept json
// @Produce json
// @Param modelID path string true "Model ID"
// @Success 200 {object} config.ModelConfig
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /api/v1/model/{modelID} [get]
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
				json.NewEncoder(w).Encode(maskModel(model))
				return
			}
		}
		http.Error(w, "Model not found", http.StatusNotFound)
	}
}
