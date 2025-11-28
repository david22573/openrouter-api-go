package openrouter

// ChatRequest matches the OpenRouter chat endpoint payload.
type ChatRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// ChatResponse is your simplified output type.
type ChatResponse struct {
	Message string `json:"message"`
}
