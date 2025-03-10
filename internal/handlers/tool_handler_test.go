package handlers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	"krackenservices.com/agentAI/internal/handlers"
	"krackenservices.com/agentAI/internal/toolmodel"
)

// fakeExecCommand simulates exec.Command by calling the test binary with a special flag.
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	// Set an environment variable to signal the helper process.
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess is a helper that is invoked by fakeExecCommand.
// It is not a real test.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Simply output a fixed string.
	os.Stdout.WriteString("fake output")
	os.Exit(0)
}

func TestDynamicToolHandler_DefaultArgs(t *testing.T) {
	// Override the command executor.
	origExecCommand := handlers.ExecCommand
	handlers.ExecCommand = fakeExecCommand
	defer func() { handlers.ExecCommand = origExecCommand }()

	// Define an internal tool configuration (using toolmodel.ToolConfig).
	toolCfg := toolmodel.ToolConfig{
		ID:          "fstool",
		Name:        "fstool",
		Description: "Internal tool to list files on the local filesystem",
		CommandKey:  "fstool",
		CommandArgs: map[string]interface{}{"path": "."},
	}

	// Create a POST request with an empty JSON body (to use default args).
	req := httptest.NewRequest(http.MethodPost, "/tool/fstool", bytes.NewBuffer([]byte(`{}`)))
	rr := httptest.NewRecorder()

	// Call the dynamic tool handler.
	handler := handlers.DynamicToolHandler(toolCfg)
	handler(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status OK; got %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading response body: %v", err)
	}

	var respMap map[string]string
	if err := json.Unmarshal(body, &respMap); err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

	// Remove any extra warning messages if present.
	output := strings.ReplaceAll(respMap["output"], "warning: GOCOVERDIR not set, no coverage data emitted\n", "")
	expected := "fake output"
	if output != expected {
		t.Errorf("expected output %q; got %q", expected, output)
	}
}

func TestDynamicToolHandler_OverrideArgs(t *testing.T) {
	// Override the command executor.
	origExecCommand := handlers.ExecCommand
	handlers.ExecCommand = fakeExecCommand
	defer func() { handlers.ExecCommand = origExecCommand }()

	toolCfg := toolmodel.ToolConfig{
		ID:          "fstool",
		Name:        "fstool",
		Description: "Internal tool to list files on the local filesystem",
		CommandKey:  "fstool",
		CommandArgs: map[string]interface{}{"path": "."},
	}

	// Create a request that overrides the "path" argument and adds an extra key ("extra") that should be ignored.
	reqBody := []byte(`{"args": {"path": "/tmp", "extra": "ignored"}}`)
	req := httptest.NewRequest(http.MethodPost, "/tool/fstool", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	handler := handlers.DynamicToolHandler(toolCfg)
	handler(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status OK; got %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading response body: %v", err)
	}

	var respMap map[string]string
	if err := json.Unmarshal(body, &respMap); err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

	output := strings.ReplaceAll(respMap["output"], "warning: GOCOVERDIR not set, no coverage data emitted\n", "")
	expected := "fake output"
	if output != expected {
		t.Errorf("expected output %q; got %q", expected, output)
	}
}
