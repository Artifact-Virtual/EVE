// llm.go - Abstract LLM Provider Interface
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

// ToolDefinition represents a tool that can be used by the LLM
type ToolDefinition struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	InputSchema anthropic.ToolInputSchemaParam `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

// Input structs for tools
type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

type ListFilesInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a directory in the working directory."`
}

type BashInput struct {
	Command string `json:"command" jsonschema_description:"The bash command to execute."`
}

type EditFileInput struct {
	Path   string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
	OldStr string `json:"old_string" jsonschema_description:"The text to replace in the file."`
	NewStr string `json:"new_string" jsonschema_description:"The new text to replace the old text with."`
}

type CodeSearchInput struct {
	Query string `json:"query" jsonschema_description:"The search query to find in the codebase."`
}

// GenerateSchema generates a JSON schema for a given type
func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T

	schema := reflector.Reflect(v)

	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}

// Tool schemas
var ReadFileInputSchema = GenerateSchema[ReadFileInput]()
var ListFilesInputSchema = GenerateSchema[ListFilesInput]()
var BashInputSchema = GenerateSchema[BashInput]()
var EditFileInputSchema = GenerateSchema[EditFileInput]()
var CodeSearchInputSchema = GenerateSchema[CodeSearchInput]()

// Tool function implementations
func ReadFile(input json.RawMessage) (string, error) {
	readFileInput := ReadFileInput{}
	err := json.Unmarshal(input, &readFileInput)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal ReadFile input: %w", err)
	}

	log.Printf("Reading file: %s", readFileInput.Path)
	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		log.Printf("Failed to read file %s: %v", readFileInput.Path, err)
		return "", err
	}
	log.Printf("Successfully read file %s (%d bytes)", readFileInput.Path, len(content))
	return string(content), nil
}

func ListFiles(input json.RawMessage) (string, error) {
	listFilesInput := ListFilesInput{}
	err := json.Unmarshal(input, &listFilesInput)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal ListFiles input: %w", err)
	}

	dir := "."
	if listFilesInput.Path != "" {
		dir = listFilesInput.Path
	}

	log.Printf("Listing files in directory: %s", dir)

	var files []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		// Skip .devenv directory and its contents
		if info.IsDir() && (relPath == ".devenv" || strings.HasPrefix(relPath, ".devenv/")) {
			return filepath.SkipDir
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Failed to list files in %s: %v", dir, err)
		return "", err
	}

	log.Printf("Successfully listed %d items in %s", len(files), dir)

	result, err := json.Marshal(files)
	if err != nil {
		return "", fmt.Errorf("failed to marshal file list: %w", err)
	}

	return string(result), nil
}

func Bash(input json.RawMessage) (string, error) {
	bashInput := BashInput{}
	err := json.Unmarshal(input, &bashInput)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal Bash input: %w", err)
	}

	log.Printf("Executing bash command: %s", bashInput.Command)

	cmd := exec.Command("bash", "-c", bashInput.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command failed: %v, output: %s", err, string(output))
		return "", fmt.Errorf("command failed: %w, output: %s", err, string(output))
	}

	log.Printf("Command executed successfully, output length: %d chars", len(output))
	return string(output), nil
}

func EditFile(input json.RawMessage) (string, error) {
	editFileInput := EditFileInput{}
	err := json.Unmarshal(input, &editFileInput)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal EditFile input: %w", err)
	}

	log.Printf("Editing file: %s", editFileInput.Path)

	content, err := os.ReadFile(editFileInput.Path)
	if err != nil {
		log.Printf("Failed to read file %s: %v", editFileInput.Path, err)
		return "", err
	}

	newContent := strings.ReplaceAll(string(content), editFileInput.OldStr, editFileInput.NewStr)

	err = os.WriteFile(editFileInput.Path, []byte(newContent), 0644)
	if err != nil {
		log.Printf("Failed to write file %s: %v", editFileInput.Path, err)
		return "", err
	}

	log.Printf("Successfully edited file %s", editFileInput.Path)
	return "File edited successfully", nil
}

func CodeSearch(input json.RawMessage) (string, error) {
	codeSearchInput := CodeSearchInput{}
	err := json.Unmarshal(input, &codeSearchInput)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal CodeSearch input: %w", err)
	}

	log.Printf("Searching for code pattern: %s", codeSearchInput.Query)

	// Simple implementation - search for the query in all .go files
	var results []string
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if strings.Contains(string(content), codeSearchInput.Query) {
				results = append(results, path)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Failed to search code: %v", err)
		return "", err
	}

	log.Printf("Found %d files containing the search pattern", len(results))

	result, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to marshal search results: %w", err)
	}

	return string(result), nil
}

// Tool definitions with function implementations
var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function:    ReadFile,
}

var ListFilesDefinition = ToolDefinition{
	Name:        "list_files",
	Description: "List the contents of a given relative directory path. Use this when you want to see what files and directories are in a directory.",
	InputSchema: ListFilesInputSchema,
	Function:    ListFiles,
}

var BashDefinition = ToolDefinition{
	Name:        "bash",
	Description: "Execute a bash command. Use this when you need to run shell commands.",
	InputSchema: BashInputSchema,
	Function:    Bash,
}

var EditFileDefinition = ToolDefinition{
	Name:        "edit_file",
	Description: "Edit a file by replacing old text with new text. Use this when you need to modify file contents.",
	InputSchema: EditFileInputSchema,
	Function:    EditFile,
}

var CodeSearchDefinition = ToolDefinition{
	Name:        "code_search",
	Description: "Search for code patterns in the codebase. Use this when you need to find specific code or patterns.",
	InputSchema: CodeSearchInputSchema,
	Function:    CodeSearch,
}

// LLMProvider defines the interface for different LLM providers
type LLMProvider interface {
	// Send a conversation and get a response
	SendMessage(ctx context.Context, conversation []Message, tools []ToolDefinition) (*LLMResponse, error)

	// Get the provider name
	Name() string

	// Get available models
	AvailableModels() []string
}

// Message represents a chat message
type Message struct {
	Role    string      `json:"role"`    // "user", "assistant", "system", "tool"
	Content interface{} `json:"content"` // string for text, or []ContentBlock for complex content
}

// ContentBlock represents different types of content in a message
type ContentBlock struct {
	Type       string      `json:"type"` // "text", "tool_use", "tool_result"
	Text       string      `json:"text,omitempty"`
	ToolUse    *ToolUse    `json:"tool_use,omitempty"`
	ToolResult *ToolResult `json:"tool_result,omitempty"`
}

// ToolUse represents a tool call from the LLM
type ToolUse struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
	IsError    bool   `json:"is_error"`
}

// LLMResponse represents the response from an LLM
type LLMResponse struct {
	Content      []ContentBlock `json:"content"`
	Usage        *Usage         `json:"usage,omitempty"`
	FinishReason string         `json:"finish_reason,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
