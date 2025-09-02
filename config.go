// config.go - Provider Configuration
package main

import (
	"fmt"
	"os"
)

// ProviderType represents different LLM providers
type ProviderType string

const (
	ProviderAnthropic ProviderType = "anthropic"
	ProviderOpenAI    ProviderType = "openai"
	ProviderGemini    ProviderType = "gemini"
)

// Config holds the configuration for the LLM provider
type Config struct {
	Provider ProviderType `json:"provider"`
	APIKey   string       `json:"api_key"`
	Model    string       `json:"model"`
}

// NewConfigFromEnv creates a config from environment variables
func NewConfigFromEnv() (*Config, error) {
	provider := ProviderType(os.Getenv("LLM_PROVIDER"))
	if provider == "" {
		provider = ProviderAnthropic // Default to Anthropic
	}

	var apiKey string
	switch provider {
	case ProviderAnthropic:
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	case ProviderOpenAI:
		apiKey = os.Getenv("OPENAI_API_KEY")
	case ProviderGemini:
		apiKey = os.Getenv("GEMINI_API_KEY")
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API key not found for provider %s", provider)
	}

	return &Config{
		Provider: provider,
		APIKey:   apiKey,
		Model:    os.Getenv("LLM_MODEL"),
	}, nil
}

// CreateProvider creates the appropriate LLM provider based on config
func (c *Config) CreateProvider() (LLMProvider, error) {
	switch c.Provider {
	case ProviderAnthropic:
		return NewAnthropicProvider(c.APIKey, c.Model), nil
	case ProviderOpenAI:
		// return NewOpenAIProvider(c.APIKey, c.Model), nil
		return nil, fmt.Errorf("OpenAI provider not implemented yet")
	case ProviderGemini:
		return NewGeminiProvider(c.APIKey, c.Model), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", c.Provider)
	}
}
