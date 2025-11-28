package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) Chat(ctx context.Context, prompt string) (*ChatResponse, error) {
	reqBody := ChatRequest{
		Model:  "openrouter/auto", // could be configurable
		Prompt: prompt,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.BaseURL+"/v1/chat/completions",
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error: %d: %s", resp.StatusCode, string(body))
	}

	var out ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return &out, nil
}
