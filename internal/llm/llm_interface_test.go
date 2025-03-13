package llm

import "testing"

// TestOpenAICall verifies that the OpenAI implementation returns the expected stubbed response.
func TestOpenAICall(t *testing.T) {
	openai := &OpenAI{}
	req := Request{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "user", Content: "Hello, OpenAI!"},
		},
		Params: map[string]interface{}{"temperature": 0.7},
	}

	resp, err := openai.Call(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Response from OpenAI"
	if resp.Output != expected {
		t.Errorf("expected %q, got %q", expected, resp.Output)
	}
}

// TestOllamaCall verifies that the Ollama implementation returns the expected stubbed response.
func TestOllamaCall(t *testing.T) {
	ollama := &Ollama{}
	req := Request{
		Model: "ollama-model",
		Messages: []Message{
			{Role: "user", Content: "Hello, Ollama!"},
		},
		Params: map[string]interface{}{"param": "value"},
	}

	resp, err := ollama.Call(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Response from Ollama"
	if resp.Output != expected {
		t.Errorf("expected %q, got %q", expected, resp.Output)
	}
}
