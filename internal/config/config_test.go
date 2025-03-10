package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"krackenservices.com/agentAI/internal/config"
)

// writeTempConfig creates a temporary config file with the given content.
func writeTempConfig(t *testing.T, dir, filename, content string) string {
	t.Helper()
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}

// TestLoadConfig_ValidConfig_InternalToolMinimal verifies that minimal configuration for an internal tool is allowed.
func TestLoadConfig_ValidConfig_InternalToolMinimal(t *testing.T) {
	tmpDir := t.TempDir()
	// Minimal config for internal tool "fstool" (which is defined in toolregistry.InternalTools)
	yamlContent := `
version: "1.0"
server:
  port:
  env:
  interface:
models:
  - id: local
    name: mymodel
    endpoint: http://127.0.0.1:8080/
    tools_supported: true
    tool_tag_start: "<tool>"
    tool_tag_end: "</tool>"
    tools:
      - fstool
tools:
  - id: fstool
    enabled: true
`
	configPath := writeTempConfig(t, tmpDir, "config.yaml", yamlContent)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}

	if cfg.Version != "1.0" {
		t.Errorf("expected version '1.0', got %s", cfg.Version)
	}
	// Defaults for server settings should be applied.
	if cfg.Server.Port != "8080" {
		t.Errorf("expected default server port '8080', got %s", cfg.Server.Port)
	}
	if len(cfg.Models) != 1 {
		t.Errorf("expected 1 model, got %d", len(cfg.Models))
	}
	// For internal tool "fstool", a minimal config should be accepted.
	found := false
	for _, tool := range cfg.Tools {
		if tool.ID == "fstool" {
			found = true
			if tool.Enabled == nil || !*tool.Enabled {
				t.Errorf("expected fstool to be enabled")
			}
		}
	}
	if !found {
		t.Errorf("expected internal tool fstool to be found in config")
	}
}

// TestLoadConfig_MissingModels ensures an error is returned when models are missing.
func TestLoadConfig_MissingModels(t *testing.T) {
	tmpDir := t.TempDir()
	yamlContent := `
version: "1.0"
server:
  port: "9090"
`
	configPath := writeTempConfig(t, tmpDir, "config.yaml", yamlContent)
	_, err := config.LoadConfig(configPath)
	if err == nil {
		t.Fatal("expected error due to missing models, got nil")
	}
}

// TestLoadConfig_IncompleteTool_External verifies that an external tool entry must be complete.
func TestLoadConfig_IncompleteTool_External(t *testing.T) {
	tmpDir := t.TempDir()
	// "externaltool" is not in the internal registry, so its config must be complete.
	yamlContent := `
version: "1.0"
server:
  port: "8080"
models:
  - id: local
    name: mymodel
    endpoint: http://127.0.0.1:8080/
    tools_supported: true
tools:
  - id: externaltool
    enabled: true
`
	configPath := writeTempConfig(t, tmpDir, "config.yaml", yamlContent)
	_, err := config.LoadConfig(configPath)
	if err == nil {
		t.Fatal("expected error due to incomplete external tool configuration, got nil")
	}
}

// TestLoadConfig_FileDoesNotExist verifies that a non-existent file returns an error.
func TestLoadConfig_FileDoesNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	nonexistentPath := filepath.Join(tmpDir, "nonexistent.yaml")
	_, err := config.LoadConfig(nonexistentPath)
	if err == nil {
		t.Fatal("expected error due to non-existent config file, got nil")
	}
}

// TestLoadConfig_DefaultEnabled verifies that if 'enabled' is omitted, it defaults to true.
func TestLoadConfig_DefaultEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	// Minimal config for internal tool "fstool" without explicitly setting enabled.
	yamlContent := `
version: "1.0"
models:
  - id: local
    name: mymodel
    endpoint: http://127.0.0.1:8080/
    tools_supported: true
tools:
  - id: fstool
`
	configPath := writeTempConfig(t, tmpDir, "config.yaml", yamlContent)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
	for _, tool := range cfg.Tools {
		if tool.Enabled == nil || !*tool.Enabled {
			t.Errorf("expected tool %s to be enabled by default", tool.ID)
		}
	}
}

// TestLoadConfig_DisabledTool verifies that a tool can be explicitly disabled.
func TestLoadConfig_DisabledTool(t *testing.T) {
	tmpDir := t.TempDir()
	yamlContent := `
version: "1.0"
models:
  - id: local
    name: mymodel
    endpoint: http://127.0.0.1:8080/
    tools_supported: true
tools:
  - id: fstool
    enabled: false
`
	configPath := writeTempConfig(t, tmpDir, "config.yaml", yamlContent)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
	for _, tool := range cfg.Tools {
		if tool.ID == "fstool" {
			if tool.Enabled == nil || *tool.Enabled {
				t.Errorf("expected tool fstool to be disabled")
			}
		}
	}
}
