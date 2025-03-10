# agentAI

agentAI is a modular Go application that provides an API service with dynamic tool execution. It supports both internal tools (built as part of the project) and external tools (dropped in by users). Internal tools are enabled by default unless explicitly disabled via configuration, and external tools must be fully configured in a YAML file.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Build and Run](#build-and-run)
- [Configuration](#configuration)

## Features

- **Dynamic Tool Execution:**  
  Execute tools via HTTP POST requests. Tools can be internal (built as part of the project) or external (configured via a YAML file).

- **Plugin Architecture:**  
  New tools can be added by simply dropping in a new tool configuration in the YAML file (external) or by updating the internal registry (internal).

- **Configurable via YAML:**  
  Use a configuration file (`config.yaml` or `config.yml`) to define server settings, models, and tool overrides.

- **Swagger Integration:**  
  Interactive API documentation is available at `/swagger/index.html`.

- **CLI Tool:**  
  A command-line interface is provided to test API endpoints.

## Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/yourusername/agentAI.git
   cd agentAI
   ```

2. Install Dependencies:

Ensure you have Go installed (version 1.16+ is recommended). Then run:

```bash
go mod tidy
```

3 .Install Swagger Tools (Optional):

If you want to generate Swagger docs, install the swag CLI tool:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger
```

## Build and Run
### Build with Makefile
The project uses a Makefile to simplify building, testing, and docker image creation.

Build All:

```bash
make
```

Build API Only:

```bash
make api
```
Build CLI Tool:

```bash
make cli
```

Build All Tools:

```bash
make tools
```

Run Tests:

```bash
make test
```

Run Tests with Coverage:

```bash
make coverage
```

### Run API
After building, run the API server:

```bash
./bin/agentAI-api
```

The API will start on the port specified in your configuration (default is 8080).

### Swagger Documentation
Once the API is running, open your browser and navigate to:

http://localhost:8080/swagger/index.html

This provides an interactive interface for testing API endpoints.

### Configuration
The API uses a YAML configuration file. By default, the API looks for config.yaml (or config.yml) in the same directory as the API binary. Here is an example config:

Disable a tool globally
```yaml
tools:
- id: ccc
  enabled: false
```

Adding a tool to a model enables (if not globally disabled) it for that model.
```yaml
model:
    tools:
    - someTool
```


```yaml
version: 1.0
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

  - id: anotherModel
    name: anotherModel
    endpoint: http://127.0.0.1:8080/
    tools_supported: true
    tool_tag_start: "<tool>"
    tool_tag_end: "</tool>"
    tools:
      - fstool
      - externaltool

tools:
  - id: fstool
    enabled: false
    
  # Example external tool (must be fully configured)
  - id: externaltool
    name: External Tool
    enabled: true
    description: An external tool example
    command_key: externaltool
    command_args: '{"arg": "value"}'
    example: '{"tool": "externaltool", "args": {"arg": "value"}}'
    example_response: '{"output": "external output"}'
```