package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"krackenservices.com/agentAI/internal/config"
	"krackenservices.com/agentAI/internal/routes"
	"krackenservices.com/agentAI/internal/toolmodel"
)

var apiv1 = "/api/v1"

// boolPtr is a helper to return a pointer to a bool.
func boolPtr(b bool) *bool {
	return &b
}

func TestRouter_ToolEndpoints(t *testing.T) {
	// Create a dummy configuration with:
	// - An internal tool "fstool" that is disabled via config.
	// - An external tool "externaltool" with a complete configuration.
	cfg := &config.Config{
		Version: "1.0",
		Models: []config.ModelConfig{
			{
				ID:             "local",
				Name:           "mymodel",
				Endpoint:       "http://127.0.0.1:8080/",
				ToolsSupported: true,
			},
		},
		Tools: []toolmodel.ToolConfig{
			// Override for internal tool "fstool" to disable it.
			{
				ID:      "fstool",
				Enabled: boolPtr(false),
			},
			// External tool configuration.
			{
				ID:          "externaltool",
				Name:        "External Tool",
				Description: "An external tool",
				CommandKey:  "externaltool",
				CommandArgs: map[string]interface{}{"arg": "default"},
				Example:     map[string]interface{}{"tool": "externaltool"},
				ExampleResponse: map[string]interface{}{
					"output": "external output",
				},
				Enabled: boolPtr(true),
			},
		},
	}

	router := routes.NewRouter(cfg)

	// Check that the /hello endpoint is registered.
	reqHello := httptest.NewRequest(http.MethodGet, apiv1+"/hello", nil)
	rrHello := httptest.NewRecorder()
	router.ServeHTTP(rrHello, reqHello)
	if rrHello.Code == http.StatusNotFound {
		t.Errorf("expected %s/hello endpoint, got 404", apiv1)
	}

	// Check that the internal tool "fstool" is disabled, so /tool/fstool should not be registered.
	reqFstool := httptest.NewRequest(http.MethodPost, apiv1+"/tool/fstool", nil)
	rrFstool := httptest.NewRecorder()
	router.ServeHTTP(rrFstool, reqFstool)
	if rrFstool.Code != http.StatusNotFound {
		t.Errorf("expected %s/tool/fstool to be unregistered (disabled), got %d", apiv1, rrFstool.Code)
	}

	// Check that the external tool "externaltool" is registered.
	reqExternal := httptest.NewRequest(http.MethodPost, apiv1+"/tool/externaltool", nil)
	rrExternal := httptest.NewRecorder()
	router.ServeHTTP(rrExternal, reqExternal)
	if rrExternal.Code == http.StatusNotFound {
		t.Errorf("expected %s/tool/externaltool to be registered, got 404", apiv1)
	}
}
