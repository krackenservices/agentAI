package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"krackenservices.com/agentAI/internal/config"
	"krackenservices.com/agentAI/internal/toolregistry"
	"net/http"
	"regexp"
	"strings"
)

// Step1 Recieve Message from the user
// Step2 Construct the payload to send to the LLM (Tool context etc)
// Step3 Send the message to the LLM
// Step4 Recieve the response from the LLM
// Step5 Check for tool commands
// Step6 Call tools, send the response to LLM
// Step7 Loop 4 - 6 until no tool comamnds found
// Step8 Send the final response to the user

// ChatRequest defines the payload to send to the LLM.
type ChatRequest struct {
	Model   string                 `json:"model"`
	Message string                 `json:"message"`
	Params  map[string]interface{} `json:"params"`
}

func ChatHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w, "Only POST requests are allowed")
			return
		}

		var payload ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Step 2: Construct the payload to send to the LLM.
		// 1. Determine Model from payload
		var selectedModel *config.ModelConfig
		for i, m := range cfg.Models {
			if m.ID == payload.Model {
				selectedModel = &cfg.Models[i]
				break
			}
		}
		if selectedModel == nil {
			http.Error(w, fmt.Sprintf("Model %q not found", payload.Model), http.StatusBadRequest)
			return
		}
		toolContext := buildToolContext(*selectedModel)
		//fmt.Printf("Tool Context: %s\n", toolContext)

		llmResponse, err := callLLM(selectedModel, payload, toolContext)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error calling LLM: %v", err), http.StatusInternalServerError)
			return
		}

		// Loop until no tool commands are found.
		for {
			command, found := extractToolCommand(selectedModel, llmResponse)
			if !found {
				break
			}

			toolResult, err := callTool(command)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error calling tool: %v", err), http.StatusInternalServerError)
				return
			}

			// Step 7: Append the tool result to the conversation and send it back to the LLM.
			// For demonstration, we simply append the tool result to the current message.
			payload.Message = llmResponse + "\nTool result: " + toolResult
			llmResponse, err = callLLM(selectedModel, payload, toolContext)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error calling LLM after tool execution: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// Step 8: Send the final LLM response to the user.
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, llmResponse)
	}
}

func callLLM(model *config.ModelConfig, payload ChatRequest, context string) (string, error) {
	message := context + "\n" + payload.Message

	fmt.Println("Sending message:\n\n" + message)

	body, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	//model.Endpoint
	var sb strings.Builder

	if !strings.Contains(payload.Message, "<tool>") {
		sb.WriteString("I need the listing of the current directory")
		sb.WriteString("<tool>{ \"tool\": \"fstool\", \"args\": { \"path\": \".\" } }</tool>")
	} else {
		sb.WriteString("I have the listing of the current directory")
		sb.WriteString("{ \"output\": \"Listing contents of directory: .\\nfile1\\ndir1\\ndir1/subdir1\\n\"}")
	}
	return sb.String(), nil
}

func buildToolContext(model config.ModelConfig) string {
	var sb strings.Builder
	sb.WriteString("You are an assistant that can call external tools when needed.\n")
	sb.WriteString("Available Tools:\n")
	for _, toolID := range model.Tools {
		if tool, ok := toolregistry.InternalTools[toolID]; ok {
			// You can format this context as needed. For example, include ID and description.
			cmdArgs, _ := json.Marshal(tool.CommandArgs)
			sb.WriteString(fmt.Sprintf("- %s: %s\n%v", tool.ID, tool.Description, string(cmdArgs)))
		} else {
			// Optionally, you can add external tools if available.
			sb.WriteString(fmt.Sprintf("- %s: (external tool)\n", toolID))
		}
	}
	sb.WriteString("If you need to fetch external data or perform a task, return a tool call using the following format:\n")
	sb.WriteString(fmt.Sprintf("%s\n\"{\"name\": \"<tool_name>\", \"arguments\": {\"arg1\": \"value1\"}}\"\n%s\n", model.ToolTagStart, model.ToolTagEnd))

	// Just get the 1st internal tool for demonstration.
	var key string
	for k := range toolregistry.InternalTools {
		key = k
		break
	}

	cmdExample, _ := json.Marshal(toolregistry.InternalTools[key].Example)
	sb.WriteString(fmt.Sprintf("For example %s\n", string(cmdExample)))
	return sb.String()
}

// extractToolCommand searches for a tool command pattern in the response.
// We assume tool commands are enclosed in <tool>...</tool>.
func extractToolCommand(model *config.ModelConfig, response string) (string, bool) {
	re := regexp.MustCompile(model.ToolTagStart + `(.*?)` + model.ToolTagEnd)
	matches := re.FindStringSubmatch(response)
	if len(matches) > 1 {
		fmt.Println("**** MATCH *****")
		return matches[1], true
	}
	return "", false
}

// callTool simulates calling a tool given a command.
// Replace this stub with your actual tool-calling logic.
func callTool(command string) (string, error) {
	// For demonstration, simply return a string that echoes the command.
	return "Tool response for command: " + command, nil
}
