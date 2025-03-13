package routes

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	"krackenservices.com/agentAI/internal/config"
	"krackenservices.com/agentAI/internal/handlers"
	"krackenservices.com/agentAI/internal/toolregistry"
)

var apiv1 = "/api/v1"

// NewRouter returns an HTTP handler with routes for the API.
func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Serve Swagger docs at /swagger/index.html
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Register static endpoints.
	mux.HandleFunc(apiv1+"/hello", handlers.HelloHandler)

	// Register dynamic tool endpoints.
	for id, tool := range toolregistry.InternalTools {
		disabled := false
		for _, cfgTool := range cfg.Tools {
			if cfgTool.ID == id {
				if cfgTool.Enabled != nil && !*cfgTool.Enabled {
					disabled = true
				}
				break
			}
		}
		if !disabled {
			route := apiv1 + "/tool/" + id
			mux.HandleFunc(route, handlers.DynamicToolHandler(tool))
		}
	}
	for _, tool := range cfg.Tools {
		if _, exists := toolregistry.InternalTools[tool.ID]; exists {
			continue
		}
		if tool.Enabled != nil && !*tool.Enabled {
			continue
		}
		route := apiv1 + "/tool/" + tool.ID
		mux.HandleFunc(route, handlers.DynamicToolHandler(tool))
	}

	// Register info endpoints for tools.
	mux.HandleFunc(apiv1+"/tools", handlers.ListTools(cfg))
	mux.HandleFunc(apiv1+"/tools/internal", handlers.ListInternalTools(cfg))
	mux.HandleFunc(apiv1+"/tools/external", handlers.ListExternalTools(cfg))

	// Register endpoints for models.
	mux.HandleFunc(apiv1+"/models", handlers.ListModels(cfg))
	mux.HandleFunc(apiv1+"/model/", handlers.GetModel(cfg)) // expects /model/<modelID>

	// Register endpoint for chat
	// TODO: Create endpoints for ollama/openai to help ux
	mux.HandleFunc(apiv1+"/chat", handlers.ChatHandler(cfg))

	return mux
}
