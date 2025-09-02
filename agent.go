// agent.go - EVE Generic Agent using LLM Provider Interface
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

// GenericAgent uses the LLMProvider interface for provider-agnostic operation
type GenericAgent struct {
	provider       LLMProvider
	getUserMessage func() (string, bool)
	tools          []ToolDefinition
	verbose        bool
}

func NewGenericAgent(
	provider LLMProvider,
	getUserMessage func() (string, bool),
	tools []ToolDefinition,
	verbose bool,
) *GenericAgent {
	return &GenericAgent{
		provider:       provider,
		getUserMessage: getUserMessage,
		tools:          tools,
		verbose:        verbose,
	}
}

func (a *GenericAgent) Run(ctx context.Context) error {
	conversation := []Message{}

	// EVE Welcome Banner
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                           🤖 EVE                           ║")
	fmt.Println("║              Your AI-Powered Coding Assistant              ║")
	fmt.Println("║                                                              ║")
	fmt.Println("║  🔧 Multi-Provider Support | 🛠️  Advanced Tools           ║")
	fmt.Println("║  📁 File Operations       | 🔍 Code Search                 ║")
	fmt.Println("║  💻 Terminal Commands     | ✏️  Code Editing               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	if a.verbose {
		log.Printf("EVE starting chat session with provider: %s", a.provider.Name())
		log.Printf("Available models: %v", a.provider.AvailableModels())
	}
	fmt.Printf("🤖 EVE with %s (use 'ctrl-c' to quit)\n", a.provider.Name())
	fmt.Println("💡 Try: 'Read the riddle.txt file and solve the puzzle'")
	fmt.Println()

	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		userInput, ok := a.getUserMessage()
		if !ok {
			if a.verbose {
				log.Println("User input ended, breaking from chat loop")
			}
			break
		}

		// Skip empty messages
		if userInput == "" {
			if a.verbose {
				log.Println("Skipping empty message")
			}
			continue
		}

		if a.verbose {
			log.Printf("User input received: %q", userInput)
		}

		// Add user message to conversation
		userMessage := Message{
			Role:    "user",
			Content: userInput,
		}
		conversation = append(conversation, userMessage)

		if a.verbose {
			log.Printf("Sending message to %s, conversation length: %d", a.provider.Name(), len(conversation))
		}

		// Get response from provider
		response, err := a.provider.SendMessage(ctx, conversation, a.tools)
		if err != nil {
			if a.verbose {
				log.Printf("Error during inference: %v", err)
			}
			return err
		}

		// Process the response
		var toolResults []ContentBlock
		var hasToolUse bool

		if a.verbose {
			log.Printf("Processing %d content blocks from %s", len(response.Content), a.provider.Name())
		}

		for _, content := range response.Content {
			switch content.Type {
			case "text":
				fmt.Printf("\u001b[93m%s\u001b[0m: %s\n", a.provider.Name(), content.Text)
			case "tool_use":
				hasToolUse = true
				toolUse := content.ToolUse
				if a.verbose {
					log.Printf("Tool use detected: %s with input: %s", toolUse.Name, string(toolUse.Input))
				}
				fmt.Printf("\u001b[96mtool\u001b[0m: %s(%s)\n", toolUse.Name, string(toolUse.Input))

				// Find and execute the tool
				var toolResult string
				var toolError error
				var toolFound bool
				for _, tool := range a.tools {
					if tool.Name == toolUse.Name {
						if a.verbose {
							log.Printf("Executing tool: %s", tool.Name)
						}
						toolResult, toolError = tool.Function(toolUse.Input)
						fmt.Printf("\u001b[92mresult\u001b[0m: %s\n", toolResult)
						if toolError != nil {
							fmt.Printf("\u001b[91merror\u001b[0m: %s\n", toolError.Error())
						}
						if a.verbose {
							if toolError != nil {
								log.Printf("Tool execution failed: %v", toolError)
							} else {
								log.Printf("Tool execution successful, result length: %d chars", len(toolResult))
							}
						}
						toolFound = true
						break
					}
				}

				if !toolFound {
					toolError = fmt.Errorf("tool '%s' not found", toolUse.Name)
					fmt.Printf("\u001b[91merror\u001b[0m: %s\n", toolError.Error())
				}

				// Add tool result to collection
				if toolError != nil {
					toolResults = append(toolResults, ContentBlock{
						Type: "tool_result",
						ToolResult: &ToolResult{
							ToolCallID: toolUse.ID,
							Content:    toolError.Error(),
							IsError:    true,
						},
					})
				} else {
					toolResults = append(toolResults, ContentBlock{
						Type: "tool_result",
						ToolResult: &ToolResult{
							ToolCallID: toolUse.ID,
							Content:    toolResult,
							IsError:    false,
						},
					})
				}
			}
		}

		// If there were no tool uses, add the assistant's response to conversation
		if !hasToolUse {
			assistantMessage := Message{
				Role:    "assistant",
				Content: response.Content, // Keep the full content blocks
			}
			conversation = append(conversation, assistantMessage)
		} else {
			// Send tool results back to the provider
			if a.verbose {
				log.Printf("Sending %d tool results back to %s", len(toolResults), a.provider.Name())
			}
			toolResultMessage := Message{
				Role:    "user",
				Content: toolResults,
			}
			conversation = append(conversation, toolResultMessage)

			// Get the provider's response after tool execution
			followupResponse, err := a.provider.SendMessage(ctx, conversation, a.tools)
			if err != nil {
				if a.verbose {
					log.Printf("Error during followup inference: %v", err)
				}
				return err
			}

			// Process followup response
			for _, content := range followupResponse.Content {
				if content.Type == "text" {
					fmt.Printf("\u001b[93m%s\u001b[0m: %s\n", a.provider.Name(), content.Text)
				}
			}

			// Add followup response to conversation
			assistantMessage := Message{
				Role:    "assistant",
				Content: followupResponse.Content,
			}
			conversation = append(conversation, assistantMessage)

			if a.verbose {
				log.Printf("Received followup response with %d content blocks", len(followupResponse.Content))
			}
		}
	}

	if a.verbose {
		log.Println("Chat session ended")
	}
	return nil
}

// main function for the generic agent
func main() {
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	flag.Parse()

	if *verbose {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Verbose logging enabled")
	} else {
		log.SetOutput(os.Stdout)
		log.SetFlags(0)
		log.SetPrefix("")
	}

	// Load configuration
	config, err := NewConfigFromEnv()
	if err != nil {
		fmt.Printf("Configuration error: %s\n", err.Error())
		fmt.Println("Please set the appropriate environment variables:")
		fmt.Println("  For Anthropic: ANTHROPIC_API_KEY")
		fmt.Println("  For OpenAI: OPENAI_API_KEY")
		fmt.Println("  For Gemini: GEMINI_API_KEY")
		fmt.Println("  Optional: LLM_PROVIDER (anthropic, openai, gemini)")
		fmt.Println("  Optional: LLM_MODEL (specific model name)")
		os.Exit(1)
	}

	// Create provider
	provider, err := config.CreateProvider()
	if err != nil {
		fmt.Printf("Provider creation error: %s\n", err.Error())
		os.Exit(1)
	}

	if *verbose {
		log.Printf("Initialized provider: %s with model: %s", provider.Name(), config.Model)
	}

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	// Define tools (same as before)
	tools := []ToolDefinition{
		ReadFileDefinition,
		ListFilesDefinition,
		BashDefinition,
		EditFileDefinition,
		CodeSearchDefinition,
	}

	if *verbose {
		log.Printf("Initialized %d tools", len(tools))
	}

	agent := NewGenericAgent(provider, getUserMessage, tools, *verbose)
	err = agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
