package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"krackenservices.com/agentAI/internal/toolregistry"
)

// internalToolArgOrder maps internal tool IDs to an ordered list of argument keys.
var internalToolArgOrder = map[string][]string{
	"fstool": {"path"},
	// Add additional internal tools here with their positional argument keys.
}

// printUsage displays CLI usage instructions.
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  agentai tool <tool_id> [arguments...]")
	fmt.Println("")
	fmt.Println("For internal tools:")
	fmt.Println("  If a single argument is provided, the default key is assumed based on the tool's argument order.")
	fmt.Println("  If more than one argument is provided, each argument must be in key=value format.")
	fmt.Println("")
	fmt.Println("For external tools, provide arguments as key=value pairs.")
}

func main() {
	flag.Parse()
	args := flag.Args()

	// We expect at least two arguments: "tool" and <tool_id>
	if len(args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// First argument should be "tool".
	if args[0] != "tool" {
		fmt.Printf("Unknown command: %s\n", args[0])
		printUsage()
		os.Exit(1)
	}

	toolID := args[1]
	overrideArgs := make(map[string]interface{})

	// Determine if this is an internal tool.
	_, isInternal := toolregistry.InternalTools[toolID]
	if isInternal {
		argOrder, exists := internalToolArgOrder[toolID]
		if !exists {
			fmt.Printf("No argument order defined for internal tool %q\n", toolID)
			os.Exit(1)
		}
		// For internal tools:
		// If exactly one parameter is provided (after toolID), assume it's for the first key.
		if len(args) == 3 {
			overrideArgs[argOrder[0]] = args[2]
		} else if len(args) > 3 {
			// Expect each argument to be in key=value format.
			for _, arg := range args[2:] {
				if !strings.Contains(arg, "=") {
					fmt.Printf("Error: for internal tool %q, when providing more than one argument, use key=value format. Invalid argument: %q\n", toolID, arg)
					os.Exit(1)
				}
				parts := strings.SplitN(arg, "=", 2)
				overrideArgs[parts[0]] = parts[1]
			}
		}
	} else {
		// External tool: require key=value pairs.
		if len(args) < 3 {
			fmt.Printf("For external tools, provide arguments in key=value format.\n")
			printUsage()
			os.Exit(1)
		}
		for _, arg := range args[2:] {
			if !strings.Contains(arg, "=") {
				fmt.Printf("Error: argument '%s' is not in key=value format.\n", arg)
				os.Exit(1)
			}
			parts := strings.SplitN(arg, "=", 2)
			overrideArgs[parts[0]] = parts[1]
		}
	}

	// Build the JSON payload.
	payload := map[string]interface{}{
		"args": overrideArgs,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling JSON payload: %v", err)
	}

	// Build the API URL.
	url := fmt.Sprintf("http://localhost:8080/tool/%s", toolID)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status %s with message: %s", resp.Status, string(body))
	}

	fmt.Printf("Response:\n%s\n", string(body))
}
