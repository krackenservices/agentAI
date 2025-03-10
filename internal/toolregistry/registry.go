package toolregistry

import "krackenservices.com/agentAI/internal/toolmodel"

// InternalTools holds the built‑in tools that are enabled by default.
var InternalTools = map[string]toolmodel.ToolConfig{
	"fstool": {
		ID:          "fstool",
		Name:        "fstool",
		Description: "Internal tool to list files on the local filesystem",
		CommandKey:  "fstool",
		CommandArgs: map[string]interface{}{"path": "."},
	},
	// Add other internal tools as needed.
}
