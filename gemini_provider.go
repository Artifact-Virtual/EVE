// gemini_provider.go - Google Gemini Provider Implementation
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiProvider implements the LLMProvider interface for Google Gemini
type GeminiProvider struct {
	client *genai.Client
	model  string
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey string, model string) *GeminiProvider {
	if model == "" {
		model = "gemini-1.5-pro"
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		// In a real implementation, you'd want to handle this error properly
		panic(fmt.Sprintf("Failed to create Gemini client: %v", err))
	}

	return &GeminiProvider{
		client: client,
		model:  model,
	}
}

// Name returns the provider name
func (p *GeminiProvider) Name() string {
	return "Google Gemini"
}

// AvailableModels returns available Gemini models
func (p *GeminiProvider) AvailableModels() []string {
	return []string{
		"gemini-1.5-pro",
		"gemini-1.5-flash",
		"gemini-1.0-pro",
	}
}

// SendMessage sends a message to Gemini and returns the response
func (p *GeminiProvider) SendMessage(ctx context.Context, conversation []Message, tools []ToolDefinition) (*LLMResponse, error) {
	// Get the Gemini model
	model := p.client.GenerativeModel(p.model)

	// Configure the model for function calling if tools are provided
	if len(tools) > 0 {
		var geminiTools []*genai.Tool
		for _, tool := range tools {
			// Create function declaration
			function := &genai.FunctionDeclaration{
				Name:        tool.Name,
				Description: tool.Description,
			}

			// For now, create a simple schema - in a full implementation,
			// you'd want to properly convert from the anthropic schema format
			function.Parameters = &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"input": {
						Type:        genai.TypeString,
						Description: "Input parameter for the tool",
					},
				},
			}

			geminiTools = append(geminiTools, &genai.Tool{
				FunctionDeclarations: []*genai.FunctionDeclaration{function},
			})
		}
		model.Tools = geminiTools
	}

	// Start a chat session
	chat := model.StartChat()

	// Convert our conversation to Gemini format
	for _, msg := range conversation {
		var parts []genai.Part

		switch msg.Role {
		case "user":
			if content, ok := msg.Content.(string); ok {
				parts = append(parts, genai.Text(content))
			} else if contentBlocks, ok := msg.Content.([]ContentBlock); ok {
				for _, block := range contentBlocks {
					switch block.Type {
					case "text":
						parts = append(parts, genai.Text(block.Text))
					case "tool_result":
						if block.ToolResult != nil {
							// For tool results, we need to send them as function responses
							parts = append(parts, genai.FunctionResponse{
								Name: block.ToolResult.ToolCallID,
								Response: map[string]interface{}{
									"result": block.ToolResult.Content,
									"error":  block.ToolResult.IsError,
								},
							})
						}
					}
				}
			}
			chat.SendMessage(ctx, parts...)

		case "assistant":
			if content, ok := msg.Content.(string); ok {
				parts = append(parts, genai.Text(content))
			} else if contentBlocks, ok := msg.Content.([]ContentBlock); ok {
				for _, block := range contentBlocks {
					switch block.Type {
					case "text":
						parts = append(parts, genai.Text(block.Text))
					case "tool_use":
						if block.ToolUse != nil {
							// For tool calls, we need to send them as function calls
							var input map[string]interface{}
							json.Unmarshal(block.ToolUse.Input, &input)
							parts = append(parts, genai.FunctionCall{
								Name: block.ToolUse.Name,
								Args: input,
							})
						}
					}
				}
			}
			chat.SendMessage(ctx, parts...)
		}
	}

	// Send the last user message and get response
	var lastUserParts []genai.Part
	if len(conversation) > 0 {
		lastMsg := conversation[len(conversation)-1]
		if lastMsg.Role == "user" {
			if content, ok := lastMsg.Content.(string); ok {
				lastUserParts = append(lastUserParts, genai.Text(content))
			} else if contentBlocks, ok := lastMsg.Content.([]ContentBlock); ok {
				for _, block := range contentBlocks {
					switch block.Type {
					case "text":
						lastUserParts = append(lastUserParts, genai.Text(block.Text))
					}
				}
			}
		}
	}

	if len(lastUserParts) == 0 {
		lastUserParts = append(lastUserParts, genai.Text("Hello"))
	}

	resp, err := chat.SendMessage(ctx, lastUserParts...)
	if err != nil {
		return nil, fmt.Errorf("gemini API error: %w", err)
	}

	// Convert Gemini response to our generic format
	response := &LLMResponse{
		Content:      make([]ContentBlock, len(resp.Candidates[0].Content.Parts)),
		FinishReason: string(resp.Candidates[0].FinishReason),
	}

	if resp.UsageMetadata != nil {
		response.Usage = &Usage{
			PromptTokens:     int(resp.UsageMetadata.PromptTokenCount),
			CompletionTokens: int(resp.UsageMetadata.CandidatesTokenCount),
			TotalTokens:      int(resp.UsageMetadata.TotalTokenCount),
		}
	}

	for i, part := range resp.Candidates[0].Content.Parts {
		switch p := part.(type) {
		case genai.Text:
			response.Content[i] = ContentBlock{
				Type: "text",
				Text: string(p),
			}
		case genai.FunctionCall:
			input, _ := json.Marshal(p.Args)
			response.Content[i] = ContentBlock{
				Type: "tool_use",
				ToolUse: &ToolUse{
					ID:    fmt.Sprintf("call_%d", i), // Gemini doesn't provide IDs, so we generate one
					Name:  p.Name,
					Input: input,
				},
			}
		}
	}

	return response, nil
}
