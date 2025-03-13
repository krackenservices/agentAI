package handlers

import (
	"encoding/json"
	"fmt"
	"krackenservices.com/agentAI/internal/config"
	"krackenservices.com/agentAI/internal/toolregistry"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"krackenservices.com/agentAI/internal/toolmodel"
)

// ExecCommand allows overriding exec.Command in tests.
var ExecCommand = exec.Command

// ToolRequest represents the expected JSON request body for dynamic tools.
type ToolRequest struct {
	// Args represents key-value pairs that override the default command arguments.
	Args map[string]interface{} `json:"args"`
}

// mergeArgsOnlyExisting returns a new map containing only keys from defaultArgs,
// replacing values with those provided in requestArgs. Extra keys are ignored.
func mergeArgsOnlyExisting(defaultArgs, requestArgs map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for key, defaultValue := range defaultArgs {
		if newVal, ok := requestArgs[key]; ok {
			merged[key] = newVal
		} else {
			merged[key] = defaultValue
		}
	}
	return merged
}

// buildCommandArgs converts a map into a slice of command-line arguments.
func buildCommandArgs(args map[string]interface{}) []string {
	var cmdArgs []string
	for k, v := range args {
		cmdArgs = append(cmdArgs, "-"+k, fmt.Sprintf("%v", v))
	}
	return cmdArgs
}

// DynamicToolHandler godoc
// @Summary Executes a dynamic tool
// @Description Executes the specified tool using default command arguments overridden by provided values.
// @Tags tool
// @Accept json
// @Produce json
// @Param tool body ToolRequest true "Tool Request"
// @Success 200 {object} map[string]string "Output of the tool"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Error"
// @Router /api/v1/tool/{tool_id} [post]
func DynamicToolHandler(toolConfig toolmodel.ToolConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ToolRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		argsMap := mergeArgsOnlyExisting(toolConfig.CommandArgs, req.Args)
		cmdArgs := buildCommandArgs(argsMap)

		exePath, err := os.Executable()
		if err != nil {
			http.Error(w, "Error determining executable path: "+err.Error(), http.StatusInternalServerError)
			return
		}
		baseDir := filepath.Dir(exePath)
		toolBinary := filepath.Join(baseDir, "tools", "agentAI-"+toolConfig.ID)

		cmd := ExecCommand(toolBinary, cmdArgs...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, "Error executing tool: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"output": string(output)})
	}
}

// ListTools godoc
// @Summary List all tools
// @Description Returns a list of all tools (both internal and external) including details from the configuration.
// @Tags tools
// @Accept json
// @Produce json
// @Success 200 {array} interface{}
// @Router /api/v1/tools [get]
func ListTools(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var toolsList []interface{}

		// Add internal tools from registry; merge any override from config.
		for id, tool := range toolregistry.InternalTools {
			// Look for a config override.
			for _, cfgTool := range cfg.Tools {
				if cfgTool.ID == id {
					tool.Enabled = cfgTool.Enabled // Only overriding the Enabled flag.
					break
				}
			}
			toolsList = append(toolsList, tool)
		}

		// Add external tools (those not in internal registry).
		for _, tool := range cfg.Tools {
			if _, exists := toolregistry.InternalTools[tool.ID]; exists {
				continue
			}
			toolsList = append(toolsList, tool)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toolsList)
	}
}

// ListInternalTools godoc
// @Summary List internal tools
// @Description Returns a list of internal tools with their configuration details.
// @Tags tools
// @Accept json
// @Produce json
// @Success 200 {array} interface{}
// @Router /api/v1/tools/internal [get]
func ListInternalTools(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var toolsList []interface{}
		for id, tool := range toolregistry.InternalTools {
			// Merge any override.
			for _, cfgTool := range cfg.Tools {
				if cfgTool.ID == id {
					tool.Enabled = cfgTool.Enabled
					break
				}
			}
			toolsList = append(toolsList, tool)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toolsList)
	}
}

// ListExternalTools godoc
// @Summary List external tools
// @Description Returns a list of external tools with their configuration details.
// @Tags tools
// @Accept json
// @Produce json
// @Success 200 {array} interface{}
// @Router /api/v1/tools/external [get]
func ListExternalTools(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var toolsList []interface{}
		for _, tool := range cfg.Tools {
			if _, exists := toolregistry.InternalTools[tool.ID]; exists {
				continue
			}
			toolsList = append(toolsList, tool)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toolsList)
	}
}
