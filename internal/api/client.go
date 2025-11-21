package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://openrouter.ai/api/v1",
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// ChatMessage represents a single message in the conversation
type ChatMessage struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"` // message text
}

// ExtraBody for advanced features like reasoning
type ExtraBody struct {
	Reasoning struct {
		Enabled bool `json:"enabled"`
	} `json:"reasoning"`
}

// ChatRequest mirrors OpenRouter /chat/completions body
type ChatRequest struct {
	Model     string        `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	ExtraBody *ExtraBody    `json:"extra_body,omitempty"`
}

// ChatResponse minimal fields
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Output  []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"output"`
}

// Chat sends a prompt to OpenRouter and returns the assistant's reply
func (c *Client) Chat(model string, messages []ChatMessage, extra *ExtraBody) (string, error) {
	reqBody := ChatRequest{
		Model:     model,
		Messages:  messages,
		ExtraBody: extra,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Output) == 0 {
		return "", fmt.Errorf("empty response from OpenRouter")
	}

	return chatResp.Output[0].Content, nil
}
