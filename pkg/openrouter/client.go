package openrouter

import (
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

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key cannot be empty")
	}

	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://api.openrouter.ai",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *Client) do(req *http.Request, v any) error {
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Non-200 / non-2xx handling
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, string(body))
	}

	// decode JSON into v
	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
