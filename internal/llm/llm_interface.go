package llm

// LLM defines the common interface for language model providers.
type LLM interface {
	Call(request Request) (Response, error)
}

// Request is the input structure for LLM calls.
type Request struct {
	Model    string                 `json:"model"`
	Messages []Message              `json:"messages"` // Note: "messages" contains an array of Message objects.
	Params   map[string]interface{} `json:"params"`
}

// Message represents an individual message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response represents a generic response from an LLM.
type Response struct {
	Output string `json:"output"`
}
