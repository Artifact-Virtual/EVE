// agent.go - EVE Generic Agent using LLM Provider Interface
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GenericAgent uses the LLMProvider interface for provider-agnostic operation
type GenericAgent struct {
	provider       LLMProvider
	getUserMessage func() (string, bool)
	tools          []ToolDefinition
	verbose        bool
	database       *ProjectDatabase
}

// Global database instance
var globalDB *ProjectDatabase

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
		database:       globalDB,
	}
}

type APICallInput struct {
	URL     string            `json:"url" jsonschema_description:"The URL to make the HTTP request to"`
	Method  string            `json:"method" jsonschema_description:"HTTP method (GET, POST, PUT, DELETE, etc.)"`
	Headers map[string]string `json:"headers,omitempty" jsonschema_description:"Optional headers as key-value pairs"`
	Body    string            `json:"body,omitempty" jsonschema_description:"Optional request body"`
}

var APICallInputSchema = GenerateSchema[APICallInput]()

func APICall(input json.RawMessage) (string, error) {
	apiInput := APICallInput{}
	err := json.Unmarshal(input, &apiInput)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(apiInput.Method, apiInput.URL, strings.NewReader(apiInput.Body))
	if err != nil {
		return "", err
	}
	for k, v := range apiInput.Headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Status: %s\nBody: %s", resp.Status, string(body)), nil
}

var APICallDefinition = ToolDefinition{
	Name:        "api_call",
	Description: "Make an HTTP request to a given URL. Supports GET, POST, PUT, DELETE, etc.",
	InputSchema: APICallInputSchema,
	Function:    APICall,
}

type WebScraperInput struct {
	URL      string `json:"url" jsonschema_description:"The URL to scrape"`
	Selector string `json:"selector" jsonschema_description:"CSS selector to extract text from"`
}

var WebScraperInputSchema = GenerateSchema[WebScraperInput]()

func WebScraper(input json.RawMessage) (string, error) {
	scraperInput := WebScraperInput{}
	err := json.Unmarshal(input, &scraperInput)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(scraperInput.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	doc.Find(scraperInput.Selector).Each(func(i int, s *goquery.Selection) {
		result.WriteString(s.Text() + "\n")
	})
	return result.String(), nil
}

var WebScraperDefinition = ToolDefinition{
	Name:        "web_scraper",
	Description: "Scrape text from a webpage using a CSS selector.",
	InputSchema: WebScraperInputSchema,
	Function:    WebScraper,
}

// Database tool definitions
type SaveToDatabaseInput struct {
	Path    string `json:"path" jsonschema_description:"File path to save"`
	Content string `json:"content" jsonschema_description:"File content to save"`
}

var SaveToDatabaseInputSchema = GenerateSchema[SaveToDatabaseInput]()

func SaveToDatabase(input json.RawMessage) (string, error) {
	var dbInput SaveToDatabaseInput
	if err := json.Unmarshal(input, &dbInput); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	// Calculate simple hash for content
	hash := fmt.Sprintf("%x", len(dbInput.Content))

	// Save to database
	if globalDB != nil {
		if err := globalDB.SaveFile(dbInput.Path, dbInput.Content, hash); err != nil {
			return "", fmt.Errorf("failed to save to database: %w", err)
		}
		return fmt.Sprintf("Successfully saved file '%s' to database", dbInput.Path), nil
	}

	return "", fmt.Errorf("database not initialized")
}

var SaveToDatabaseDefinition = ToolDefinition{
	Name:        "save_to_database",
	Description: "Save a file to the project database for version control and backup.",
	InputSchema: SaveToDatabaseInputSchema,
	Function:    SaveToDatabase,
}

type CreateCheckpointInput struct {
	Name        string `json:"name" jsonschema_description:"Name of the checkpoint"`
	Description string `json:"description" jsonschema_description:"Description of the checkpoint"`
}

var CreateCheckpointInputSchema = GenerateSchema[CreateCheckpointInput]()

func CreateCheckpoint(input json.RawMessage) (string, error) {
	var cpInput CreateCheckpointInput
	if err := json.Unmarshal(input, &cpInput); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	if globalDB != nil {
		if err := globalDB.CreateCheckpoint(cpInput.Name, cpInput.Description); err != nil {
			return "", fmt.Errorf("failed to create checkpoint: %w", err)
		}
		return fmt.Sprintf("Successfully created checkpoint '%s'", cpInput.Name), nil
	}

	return "", fmt.Errorf("database not initialized")
}

var CreateCheckpointDefinition = ToolDefinition{
	Name:        "create_checkpoint",
	Description: "Create a checkpoint (snapshot) of the current project state.",
	InputSchema: CreateCheckpointInputSchema,
	Function:    CreateCheckpoint,
}

type RestoreCheckpointInput struct {
	CheckpointID int `json:"checkpoint_id" jsonschema_description:"ID of the checkpoint to restore"`
}

var RestoreCheckpointInputSchema = GenerateSchema[RestoreCheckpointInput]()

func RestoreCheckpoint(input json.RawMessage) (string, error) {
	var cpInput RestoreCheckpointInput
	if err := json.Unmarshal(input, &cpInput); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	if globalDB != nil {
		if err := globalDB.RestoreCheckpoint(cpInput.CheckpointID); err != nil {
			return "", fmt.Errorf("failed to restore checkpoint: %w", err)
		}
		return fmt.Sprintf("Successfully restored checkpoint %d", cpInput.CheckpointID), nil
	}

	return "", fmt.Errorf("database not initialized")
}

var RestoreCheckpointDefinition = ToolDefinition{
	Name:        "restore_checkpoint",
	Description: "Restore project to a previous checkpoint state.",
	InputSchema: RestoreCheckpointInputSchema,
	Function:    RestoreCheckpoint,
}

type ListCheckpointsInput struct{}

var ListCheckpointsInputSchema = GenerateSchema[ListCheckpointsInput]()

func ListCheckpoints(input json.RawMessage) (string, error) {
	if globalDB != nil {
		checkpoints, err := globalDB.ListCheckpoints()
		if err != nil {
			return "", fmt.Errorf("failed to list checkpoints: %w", err)
		}

		result := "Available Checkpoints:\n"
		for _, cp := range checkpoints {
			result += fmt.Sprintf("- ID: %d, Name: %s, Description: %s, Files: %d, Time: %s\n",
				cp.ID, cp.Name, cp.Description, cp.FileCount, cp.Timestamp.Format("2006-01-02 15:04:05"))
		}
		return result, nil
	}

	return "", fmt.Errorf("database not initialized")
}

var ListCheckpointsDefinition = ToolDefinition{
	Name:        "list_checkpoints",
	Description: "List all available project checkpoints.",
	InputSchema: ListCheckpointsInputSchema,
	Function:    ListCheckpoints,
}

type MCPIntegrationInput struct {
	Name      string `json:"name" jsonschema_description:"Name of the MCP integration"`
	Endpoint  string `json:"endpoint" jsonschema_description:"MCP server endpoint URL"`
	AuthToken string `json:"auth_token,omitempty" jsonschema_description:"Authentication token for MCP server"`
	Config    string `json:"config,omitempty" jsonschema_description:"Additional configuration JSON"`
}

var MCPIntegrationInputSchema = GenerateSchema[MCPIntegrationInput]()

func AddMCPIntegration(input json.RawMessage) (string, error) {
	var mcpInput MCPIntegrationInput
	if err := json.Unmarshal(input, &mcpInput); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	if globalDB != nil {
		if err := globalDB.AddMCPIntegration(mcpInput.Name, mcpInput.Endpoint, mcpInput.AuthToken, mcpInput.Config); err != nil {
			return "", fmt.Errorf("failed to add MCP integration: %w", err)
		}
		return fmt.Sprintf("Successfully added MCP integration '%s'", mcpInput.Name), nil
	}

	return "", fmt.Errorf("database not initialized")
}

var AddMCPIntegrationDefinition = ToolDefinition{
	Name:        "add_mcp_integration",
	Description: "Add a new MCP (Model Context Protocol) server integration.",
	InputSchema: MCPIntegrationInputSchema,
	Function:    AddMCPIntegration,
}

type MultiplayerActionInput struct {
	SessionID string `json:"session_id" jsonschema_description:"Multiplayer session ID"`
	UserID    string `json:"user_id" jsonschema_description:"User identifier"`
	Action    string `json:"action" jsonschema_description:"Action performed"`
	Data      string `json:"data" jsonschema_description:"Action data"`
}

var MultiplayerActionInputSchema = GenerateSchema[MultiplayerActionInput]()

func RecordMultiplayerAction(input json.RawMessage) (string, error) {
	var mpInput MultiplayerActionInput
	if err := json.Unmarshal(input, &mpInput); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	if globalDB != nil {
		if err := globalDB.RecordMultiplayerAction(mpInput.SessionID, mpInput.UserID, mpInput.Action, mpInput.Data); err != nil {
			return "", fmt.Errorf("failed to record multiplayer action: %w", err)
		}
		return fmt.Sprintf("Recorded multiplayer action: %s by %s", mpInput.Action, mpInput.UserID), nil
	}

	return "", fmt.Errorf("database not initialized")
}

var RecordMultiplayerActionDefinition = ToolDefinition{
	Name:        "record_multiplayer_action",
	Description: "Record a multiplayer session action for collaboration tracking.",
	InputSchema: MultiplayerActionInputSchema,
	Function:    RecordMultiplayerAction,
}

type BackupProjectInput struct {
	Path string `json:"path" jsonschema_description:"Path where to save the backup file"`
}

var BackupProjectInputSchema = GenerateSchema[BackupProjectInput]()

func BackupProject(input json.RawMessage) (string, error) {
	var backupInput BackupProjectInput
	if err := json.Unmarshal(input, &backupInput); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	if globalDB != nil {
		if err := globalDB.BackupProject(backupInput.Path); err != nil {
			return "", fmt.Errorf("failed to backup project: %w", err)
		}
		return fmt.Sprintf("Successfully backed up project to '%s'", backupInput.Path), nil
	}

	return "", fmt.Errorf("database not initialized")
}

var BackupProjectDefinition = ToolDefinition{
	Name:        "backup_project",
	Description: "Create a full backup of the project including all files and checkpoints.",
	InputSchema: BackupProjectInputSchema,
	Function:    BackupProject,
}

func (a *GenericAgent) Run(ctx context.Context) error {
	conversation := []Message{}

	// Load conversation from file
	if data, err := os.ReadFile("conversation.json"); err == nil {
		json.Unmarshal(data, &conversation)
	}

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

	// Save conversation to file
	data, _ := json.Marshal(conversation)
	os.WriteFile("conversation.json", data, 0644)

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

	// Initialize database
	dbPath := "eve_project.db"
	db, err := NewProjectDatabase(dbPath)
	if err != nil {
		fmt.Printf("Database initialization error: %s\n", err.Error())
		if *verbose {
			log.Printf("Failed to initialize database at %s: %v", dbPath, err)
		}
		// Continue without database - tools will handle gracefully
		globalDB = nil
	} else {
		globalDB = db
		if *verbose {
			log.Printf("Database initialized successfully at %s", dbPath)
		}
		defer db.Close()
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
		if globalDB != nil {
			globalDB.Close()
		}
		os.Exit(1)
	}

	// Create provider
	provider, err := config.CreateProvider()
	if err != nil {
		fmt.Printf("Provider creation error: %s\n", err.Error())
		if globalDB != nil {
			globalDB.Close()
		}
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

	// Define tools with database integration
	tools := []ToolDefinition{
		ReadFileDefinition,
		ListFilesDefinition,
		BashDefinition,
		EditFileDefinition,
		CodeSearchDefinition,
		APICallDefinition,
		WebScraperDefinition,
		SaveToDatabaseDefinition,
		CreateCheckpointDefinition,
		RestoreCheckpointDefinition,
		ListCheckpointsDefinition,
		AddMCPIntegrationDefinition,
		RecordMultiplayerActionDefinition,
		BackupProjectDefinition,
	}

	if *verbose {
		log.Printf("Initialized %d tools", len(tools))
	}

	agent := NewGenericAgent(provider, getUserMessage, tools, *verbose)
	err = agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		if globalDB != nil {
			globalDB.Close()
		}
		os.Exit(1)
	}

	if globalDB != nil {
		globalDB.Close()
	}
}
