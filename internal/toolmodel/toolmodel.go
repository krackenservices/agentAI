package toolmodel

// ToolConfig represents the configuration for a tool.
// swagger:model ToolConfig
type ToolConfig struct {
	ID              string                 `yaml:"id"`
	Name            string                 `yaml:"name"`
	Description     string                 `yaml:"description"`
	CommandKey      string                 `yaml:"command_key"`
	CommandArgs     map[string]interface{} `yaml:"command_args"`
	Example         map[string]interface{} `yaml:"example"`
	ExampleResponse map[string]interface{} `yaml:"example_response"`
	Enabled         *bool                  `yaml:"enabled,omitempty"`
}
