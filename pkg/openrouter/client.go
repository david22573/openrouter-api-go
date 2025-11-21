package openrouter

import (
	"net/http"
	"time"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apiKey, model string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://api.openrouter.ai",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
