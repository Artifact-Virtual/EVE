# EVE Agent Blueprint

## 🌟 System Overview

EVE is a sophisticated multi-provider AI coding assistant built with Go 1.24.2+, featuring advanced tool integration, file-based database management, and collaborative development capabilities.

## 🏗️ Architecture

### Core Components

```
EVE/
├── agent.go              # Main agent with 7 advanced tools
├── database.go           # File-based project management system
├── llm.go               # Provider abstraction layer
├── providers/           # AI provider implementations
│   ├── anthropic_provider.go (Claude)
│   ├── gemini_provider.go    (Google Gemini)
│   └── openai_provider.go     (GPT models)
├── tools/               # Standalone tool implementations
│   ├── bash_tool.go      # Terminal command execution
│   ├── edit_tool.go      # File editing operations
│   ├── code_search_tool.go # Codebase search
│   ├── list_files.go     # File system navigation
│   ├── read.go          # File reading capabilities
│   └── ...
└── eve_project_data/    # Auto-generated database directory
```

### Database Architecture

EVE uses a comprehensive file-based database system with JSON storage:

- **Project Files**: Version-controlled file storage with metadata
- **Checkpoints**: Project snapshots for backup/restore functionality
- **MCP Integrations**: External service connections and configurations
- **Multiplayer Actions**: Collaborative development activity tracking
- **Edit History**: Complete audit trail of file modifications

## 🛠️ Advanced Tool System

### Core Tools (7 Built-in)

1. **SaveToDatabase** - Store project files with automatic versioning
2. **CreateCheckpoint** - Create project snapshots with descriptions
3. **ListCheckpoints** - View available project checkpoints
4. **AddMCPIntegration** - Register Model Context Protocol servers
5. **GetMCPIntegrations** - List configured MCP integrations
6. **RecordMultiplayerAction** - Track collaborative development actions
7. **GetMultiplayerHistory** - View multiplayer activity history

### Tool Categories

- **File Operations**: Read, write, edit, and manage project files
- **System Integration**: Terminal commands, file system navigation
- **Code Analysis**: Search, analyze, and understand codebases
- **Database Management**: Project data persistence and retrieval
- **Collaboration**: Multiplayer features and activity tracking
- **External Services**: MCP server integration and API connectivity

## 🔧 Development Environment

### Prerequisites
- **Go 1.24.2+** - Core language runtime
- **Git** - Version control system
- **Nix** (optional) - Reproducible environments with devenv
- **API Keys** - Anthropic Claude, Google Gemini, or OpenAI

### Environment Setup

```bash
# Using devenv (recommended)
devenv shell

# Manual setup
go mod tidy
go build -o eve .
```

### Configuration

```bash
# API Keys (choose your provider)
export ANTHROPIC_API_KEY="your-claude-key"
export GOOGLE_AI_API_KEY="your-gemini-key"
export OPENAI_API_KEY="your-gpt-key"

# Optional: Database location
export EVE_DB_PATH="./eve_project_data"
```

## 🚀 Usage Patterns

### Basic Operation

```bash
# Start EVE
./eve

# With verbose logging
./eve --verbose

# Using specific provider
./eve --provider claude
```

### Tool Execution Examples

```bash
# File operations
./eve --tool save-file --path "main.go" --content "package main..."

# Database management
./eve --tool create-checkpoint --name "Version 1.0"
./eve --tool list-checkpoints

# MCP integration
./eve --tool add-mcp --name "my-server" --url "http://localhost:3000"
./eve --tool get-mcp-integrations

# Collaboration
./eve --tool record-action --user "alice" --action "edit" --data "Updated API"
./eve --tool get-multiplayer-history --limit 10
```

## 📊 Database Schema

### ProjectFile Structure
```go
type ProjectFile struct {
    ID         int       `json:"id"`
    Path       string    `json:"path"`
    Content    string    `json:"content"`
    Hash       string    `json:"hash"`
    Version    int       `json:"version"`
    CreatedAt  time.Time `json:"created_at"`
    ModifiedAt time.Time `json:"modified_at"`
    IsActive   bool      `json:"is_active"`
}
```

### Checkpoint Structure
```go
type Checkpoint struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Timestamp   time.Time `json:"timestamp"`
    FileCount   int       `json:"file_count"`
}
```

### MCP Integration Structure
```go
type MCPIntegration struct {
    ID        int                    `json:"id"`
    Name      string                 `json:"name"`
    Type      string                 `json:"type"`
    Config    map[string]interface{} `json:"config"`
    CreatedAt time.Time              `json:"created_at"`
    IsActive  bool                   `json:"is_active"`
}
```

## 🔄 Provider Abstraction

### Interface Definition

```go
type LLMProvider interface {
    Name() string
    GenerateResponse(prompt string, tools []ToolDefinition) (string, error)
    SupportsTools() bool
    MaxTokens() int
}
```

### Supported Providers

- **Anthropic Claude**: Excellent for complex reasoning and code analysis
- **Google Gemini**: Fast and efficient for quick tasks
- **OpenAI GPT**: Versatile for creative coding solutions

### Dynamic Switching

```go
// Runtime provider switching
agent := NewGenericAgent("claude")  // Switch to Claude
agent.SetProvider("gemini")         // Switch to Gemini
agent.SetProvider("openai")         // Switch to OpenAI
```

## 🌐 Web Interface

EVE includes a modern web interface (`index.html`) featuring:

- **Responsive Design**: Desktop and mobile optimized
- **Dark Theme**: Purple/pink gradient aesthetic
- **Interactive Elements**: Animated terminal simulation
- **Feature Showcase**: Comprehensive capability overview
- **Modern UI**: Tailwind CSS with custom animations

## 🔒 Security & Privacy

### Data Protection
- **API Key Security**: Environment variable storage only
- **Local Processing**: Optional offline operation mode
- **Encryption**: End-to-end encryption for sensitive data
- **Access Control**: Configurable multiplayer permissions

### Best Practices
- Never commit API keys to version control
- Use environment variables for configuration
- Enable verbose logging for debugging only
- Regularly rotate API keys

## 🚨 Error Handling

### Graceful Degradation
- Automatic fallback between providers
- Database corruption recovery
- Network timeout handling
- Tool execution error recovery

### Logging Levels
- **Info**: Normal operation logging
- **Verbose**: Detailed execution tracing
- **Error**: Critical failure reporting
- **Debug**: Development troubleshooting

## 📈 Performance Optimization

### Memory Management
- Efficient file-based database with caching
- Lazy loading for large projects
- Automatic cleanup of old checkpoints
- Memory-bounded operation queues

### Network Optimization
- Connection pooling for API calls
- Intelligent retry mechanisms
- Request batching for bulk operations
- Rate limiting and backoff strategies

## 🔧 Development Workflow

### Building
```bash
# Standard build
go build -o eve .

# Optimized build
go build -ldflags="-s -w" -o eve .

# Cross-platform
GOOS=windows GOARCH=amd64 go build -o eve.exe .
```

### Testing
```bash
# Unit tests
go test ./...

# Integration tests
go test -tags=integration ./...

# Benchmark tests
go test -bench=. ./...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Vet for issues
go vet ./...

# Security scan
go mod audit
```

## 🎯 Advanced Features

### MCP Server Integration
- **Dynamic Registration**: Runtime MCP server addition
- **Configuration Management**: Flexible server configuration
- **Health Monitoring**: Automatic server health checks
- **Load Balancing**: Multi-server request distribution

### Collaborative Development
- **Real-time Tracking**: Live activity monitoring
- **User Attribution**: Action ownership tracking
- **Conflict Resolution**: Collaborative editing support
- **Audit Trails**: Complete activity history

### Intelligent Caching
- **Response Caching**: API response memoization
- **File Hashing**: Content change detection
- **Database Indexing**: Fast query optimization
- **Memory Pooling**: Efficient resource utilization

## 🚀 Deployment Options

### Standalone Binary
```bash
# Build for deployment
go build -ldflags="-s -w" -o eve .
# Copy to target system
scp eve user@server:/usr/local/bin/
```

### Docker Container
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o eve .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/eve /usr/local/bin/
CMD ["eve"]
```

### System Service
```systemd
[Unit]
Description=EVE AI Assistant
After=network.target

[Service]
Type=simple
User=eve
ExecStart=/usr/local/bin/eve
Restart=always

[Install]
WantedBy=multi-user.target
```

## 📞 Support & Community

### Documentation
- **README.md**: Comprehensive usage guide
- **BUILD_GUIDE.md**: Detailed build instructions
- **API Reference**: Tool and provider documentation

### Troubleshooting
- **Verbose Logging**: Enable with `--verbose` flag
- **Database Reset**: Remove `eve_project_data/` directory
- **Provider Testing**: Individual provider connectivity tests

### Contributing
- **Code Standards**: Go formatting and best practices
- **Testing**: Comprehensive test coverage required
- **Documentation**: Update docs for all changes
- **Review Process**: Pull request review workflow

---

**EVE Agent Blueprint - Version 2.0** 🌸✨
*Last updated: September 2025*
- File operations (reading, writing, listing files with sizes/counts)
- Bash command execution (commands run, output, errors)
- Conversation flow (message processing, content blocks)
- Error details with stack traces

**Log output locations:**
- **Verbose mode**: Detailed logs go to stderr with timestamps and file locations
- **Normal mode**: Only essential output goes to stdout

**Common troubleshooting scenarios:**
- **API failures**: Check verbose logs for authentication errors or rate limits
- **Tool failures**: See exactly which tool failed and why (file not found, permission errors)
- **Unexpected responses**: View full conversation flow and Claude's reasoning
- **Performance issues**: See API call timing and response sizes

### Environment Issues
- Ensure `ANTHROPIC_API_KEY` environment variable is set
- Run `devenv shell` to ensure proper development environment
- Use `go mod tidy` to ensure dependencies are installed

## Notes
- Requires ANTHROPIC_API_KEY environment variable to be set
- Chat application provides a simple terminal interface to Claude
- Use ctrl-c to quit the chat session
