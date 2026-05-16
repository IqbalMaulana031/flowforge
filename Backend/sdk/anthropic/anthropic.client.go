package anthropic

import "flowforge-api/config"

type Client struct {
	apiKey string
}

func NewClient(cfg *config.Config) *Client {
	return &Client{apiKey: cfg.AI.AnthropicAPIKey}
}
