package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
	"krackenservices.com/agentAI/internal/toolmodel"
	"krackenservices.com/agentAI/internal/toolregistry"
)

// Config represents the entire configuration file.
type Config struct {
	Version string                 `yaml:"version"`
	Server  ServerConfig           `yaml:"server,omitempty"`
	Models  []ModelConfig          `yaml:"models"`
	Tools   []toolmodel.ToolConfig `yaml:"tools,omitempty"`
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Port      string `yaml:"port,omitempty"`
	Env       string `yaml:"env,omitempty"`
	Interface string `yaml:"interface,omitempty"`
}

// ModelConfig represents a model configuration.
type ModelConfig struct {
	ID             string   `yaml:"id"`
	Name           string   `yaml:"name"`
	Endpoint       string   `yaml:"endpoint"`
	ToolsSupported bool     `yaml:"tools_supported"`
	ToolTagStart   string   `yaml:"tool_tag_start,omitempty"`
	ToolTagEnd     string   `yaml:"tool_tag_end,omitempty"`
	Tools          []string `yaml:"tools,omitempty"`
}

// LoadConfig loads the configuration from the given YAML file path,
// applies sensible defaults, and validates required fields.
func LoadConfig(path string) (*Config, error) {
	// If no path is provided, look for config.yaml or config.yml in the executable's directory.
	if path == "" {
		binaryPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("could not get executable path: %w", err)
		}
		baseDir := filepath.Dir(binaryPath)
		path = filepath.Join(baseDir, "config")
		validExtensions := []string{".yaml", ".yml"}
		found := false
		for _, ext := range validExtensions {
			if _, err := os.Stat(path + ext); err == nil {
				path = path + ext
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("no config file found in %s with extensions: %v", baseDir, validExtensions)
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	// Set defaults for Server fields.
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.Env == "" {
		cfg.Server.Env = "development"
	}
	if cfg.Server.Interface == "" {
		cfg.Server.Interface = "0.0.0.0"
	}
	if cfg.Version == "" {
		cfg.Version = "1.0"
	}

	// Validate that at least one model is provided.
	if len(cfg.Models) == 0 {
		return nil, fmt.Errorf("config must define at least one model")
	}

	// Validate tool configurations.
	// For tools that are internal, allow a minimal config (e.g. only 'id' and 'enabled').
	// For external tools, require complete configuration.
	for _, tool := range cfg.Tools {
		if _, isInternal := toolregistry.InternalTools[tool.ID]; isInternal {
			// Internal tool override: allow minimal configuration.
			continue
		}
		// External tool: require all fields.
		if tool.ID == "" || tool.Name == "" || tool.Description == "" || tool.CommandKey == "" {
			return nil, fmt.Errorf("incomplete tool configuration for tool with id '%s'", tool.ID)
		}
	}

	// Default the Enabled flag to true for all tools if not provided.
	for i, tool := range cfg.Tools {
		if tool.Enabled == nil {
			enabled := true
			cfg.Tools[i].Enabled = &enabled
		}
	}

	return &cfg, nil
}
