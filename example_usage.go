// example_usage.go - Example of how to use the provider abstraction
package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

func exampleUsage() {
	// Example 1: Using Anthropic
	fmt.Println("=== Example 1: Anthropic Claude ===")
	anthropicProvider := NewAnthropicProvider(os.Getenv("ANTHROPIC_API_KEY"), "")
	fmt.Printf("Provider: %s\n", anthropicProvider.Name())
	fmt.Printf("Models: %v\n", anthropicProvider.AvailableModels())

	// Example 2: Using OpenAI (when implemented)
	fmt.Println("\n=== Example 2: OpenAI (Conceptual) ===")
	openaiProvider := NewOpenAIProvider(os.Getenv("OPENAI_API_KEY"), "")
	fmt.Printf("Provider: %s\n", openaiProvider.Name())
	fmt.Printf("Models: %v\n", openaiProvider.AvailableModels())

	// Example 3: Configuration-based provider selection
	fmt.Println("\n=== Example 3: Configuration-based Selection ===")

	// Set environment variables for different providers
	os.Setenv("LLM_PROVIDER", "anthropic")
	os.Setenv("ANTHROPIC_API_KEY", "your-anthropic-key")
	os.Setenv("LLM_MODEL", "claude-3-5-sonnet-20241022")

	config, err := NewConfigFromEnv()
	if err != nil {
		log.Printf("Config error: %v", err)
		return
	}

	provider, err := config.CreateProvider()
	if err != nil {
		log.Printf("Provider creation error: %v", err)
		return
	}

	fmt.Printf("Selected provider: %s\n", provider.Name())
	fmt.Printf("Using model: %s\n", config.Model)

	// Example conversation
	conversation := []Message{
		{
			Role:    "user",
			Content: "Hello! Can you help me with some coding tasks?",
		},
	}

	// This would work once the provider is properly implemented
	_, err = provider.SendMessage(context.Background(), conversation, []ToolDefinition{})
	if err != nil {
		fmt.Printf("Expected error (provider not fully implemented): %v\n", err)
	}
}

func mainExample() {
	exampleUsage()
}
