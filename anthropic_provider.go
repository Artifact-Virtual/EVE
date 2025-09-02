// anthropic_provider.go - Anthropic Claude Provider Implementation
package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// AnthropicProvider implements the LLMProvider interface for Claude
type AnthropicProvider struct {
	client *anthropic.Client
	model  string
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(apiKey string, model string) *AnthropicProvider {
	if model == "" {
		model = string(anthropic.ModelClaude3_7SonnetLatest)
	}

	// Set the API key in environment if provided
	if apiKey != "" {
		// Note: In production, you'd want to handle this more securely
		// For now, we'll assume it's set in the environment
	}

	client := anthropic.NewClient()
	return &AnthropicProvider{
		client: &client,
		model:  model,
	}
}

// Name returns the provider name
func (p *AnthropicProvider) Name() string {
	return "Anthropic Claude"
}

// AvailableModels returns available Claude models
func (p *AnthropicProvider) AvailableModels() []string {
	return []string{
		string(anthropic.ModelClaude3_7SonnetLatest),
		string(anthropic.ModelClaude3_5SonnetLatest),
		"claude-3-haiku-20240307", // Using string literal for Haiku
	}
}

// SendMessage sends a message to Claude and returns the response
func (p *AnthropicProvider) SendMessage(ctx context.Context, conversation []Message, tools []ToolDefinition) (*LLMResponse, error) {
	// Convert our generic messages to Anthropic format
	var anthropicMessages []anthropic.MessageParam

	for _, msg := range conversation {
		switch msg.Role {
		case "user":
			if content, ok := msg.Content.(string); ok {
				anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(content)))
			} else if contentBlocks, ok := msg.Content.([]ContentBlock); ok {
				var blocks []anthropic.ContentBlockParamUnion
				for _, block := range contentBlocks {
					switch block.Type {
					case "text":
						blocks = append(blocks, anthropic.NewTextBlock(block.Text))
					case "tool_result":
						if block.ToolResult != nil {
							blocks = append(blocks, anthropic.NewToolResultBlock(
								block.ToolResult.ToolCallID,
								block.ToolResult.Content,
								block.ToolResult.IsError,
							))
						}
					}
				}
				anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(blocks...))
			}
		case "assistant":
			if content, ok := msg.Content.(string); ok {
				anthropicMessages = append(anthropicMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(content)))
			} else if contentBlocks, ok := msg.Content.([]ContentBlock); ok {
				var blocks []anthropic.ContentBlockParamUnion
				for _, block := range contentBlocks {
					switch block.Type {
					case "text":
						blocks = append(blocks, anthropic.NewTextBlock(block.Text))
					case "tool_use":
						if block.ToolUse != nil {
							blocks = append(blocks, anthropic.NewToolUseBlock(
								block.ToolUse.ID,
								block.ToolUse.Name,
								string(block.ToolUse.Input), // Convert RawMessage to string
							))
						}
					}
				}
				anthropicMessages = append(anthropicMessages, anthropic.NewAssistantMessage(blocks...))
			}
		}
	}

	// Convert tools to Anthropic format
	var anthropicTools []anthropic.ToolUnionParam
	for _, tool := range tools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}

	// Make the API call
	message, err := p.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(p.model),
		MaxTokens: int64(1024),
		Messages:  anthropicMessages,
		Tools:     anthropicTools,
	})

	if err != nil {
		return nil, fmt.Errorf("anthropic API error: %w", err)
	}

	// Convert Anthropic response to our generic format
	response := &LLMResponse{
		Content:      make([]ContentBlock, len(message.Content)),
		FinishReason: string(message.StopReason),
	}

	if message.Usage.InputTokens > 0 || message.Usage.OutputTokens > 0 {
		response.Usage = &Usage{
			PromptTokens:     int(message.Usage.InputTokens),
			CompletionTokens: int(message.Usage.OutputTokens),
			TotalTokens:      int(message.Usage.InputTokens + message.Usage.OutputTokens),
		}
	}

	for i, content := range message.Content {
		switch content.Type {
		case "text":
			response.Content[i] = ContentBlock{
				Type: "text",
				Text: content.Text,
			}
		case "tool_use":
			toolUse := content.AsToolUse()
			response.Content[i] = ContentBlock{
				Type: "tool_use",
				ToolUse: &ToolUse{
					ID:    toolUse.ID,
					Name:  toolUse.Name,
					Input: toolUse.Input,
				},
			}
		}
	}

	return response, nil
}
