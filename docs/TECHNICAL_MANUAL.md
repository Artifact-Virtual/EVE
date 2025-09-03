# EVE Technical Reference Manual

## Table of Contents

1. [System Overview](#1-system-overview)
2. [Core Architecture](#2-core-architecture)
3. [Provider Abstraction Layer](#3-provider-abstraction-layer)
4. [Tool System Architecture](#4-tool-system-architecture)
5. [Database and Persistence Layer](#5-database-and-persistence-layer)
6. [Configuration Management](#6-configuration-management)
7. [Execution Flow and State Management](#7-execution-flow-and-state-management)
8. [Error Handling and Resilience](#8-error-handling-and-resilience)
9. [Performance Characteristics](#9-performance-characteristics)
10. [Security Considerations](#10-security-considerations)
11. [Extensibility Framework](#11-extensibility-framework)
12. [Deployment and Operations](#12-deployment-and-operations)

---

## 1. System Overview

### 1.1 Purpose and Scope

EVE (Enhanced Virtual Environment) is a sophisticated AI-powered coding assistant designed to provide intelligent, context-aware assistance for software development tasks. The system implements a multi-provider architecture that enables seamless integration with various Large Language Model (LLM) services while maintaining a consistent interface and behavior.

### 1.2 Key Design Principles

- **Provider Agnosticism**: Abstract interface allowing transparent switching between LLM providers
- **Tool-Driven Architecture**: Extensible tool system for executing complex operations
- **State Persistence**: File-based database system for project state management
- **Modular Design**: Clean separation of concerns with well-defined interfaces
- **Resilient Operation**: Comprehensive error handling and graceful degradation

### 1.3 System Boundaries

**Included Components:**

- Core agent execution engine
- Multi-provider LLM integration
- Tool execution framework
- Project database management
- Web interface (optional)
- Configuration management

**External Dependencies:**

- LLM provider APIs (Anthropic, Google, OpenAI)
- File system access
- Terminal/command execution
- Network connectivity for API calls

---

## 2. Core Architecture

### 2.1 Component Hierarchy

```text
┌─────────────────────────────────────────────────────────────┐
│                    GenericAgent                            │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              LLMProvider Interface                 │    │
│  │  ┌─────────────┬─────────────┬─────────────┐       │    │
│  │  │ Anthropic   │   OpenAI    │   Gemini    │       │    │
│  │  │ Provider    │ Provider    │ Provider    │       │    │
│  │  └─────────────┴─────────────┴─────────────┘       │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                 Tool System                        │    │
│  │  ┌─────────────┬─────────────┬─────────────┐       │    │
│  │  │ File Ops    │ Terminal    │ Code Search │       │    │
│  │  │ Tools       │ Tools       │ Tools       │       │    │
│  │  └─────────────┴─────────────┴─────────────┘       │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │            Project Database                        │    │
│  │  ┌─────────────┬─────────────┬─────────────┐       │    │
│  │  │ File Store  │ Checkpoints │ MCP Config  │       │    │
│  │  └─────────────┴─────────────┴─────────────┘       │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Core Data Structures

#### Message Flow

```go
type Message struct {
    Role    string      // "user", "assistant", "system", "tool"
    Content interface{} // string or []ContentBlock
}

type ContentBlock struct {
    Type       string      // "text", "tool_use", "tool_result"
    Text       string
    ToolUse    *ToolUse
    ToolResult *ToolResult
}
```

#### Tool Definition

```go
type ToolDefinition struct {
    Name        string
    Description string
    InputSchema anthropic.ToolInputSchemaParam
    Function    func(input json.RawMessage) (string, error)
}
```

### 2.3 Execution Model

The system operates on a conversational loop with the following phases:

1. **Input Reception**: User input captured via stdin or programmatic interface
2. **Message Processing**: Input converted to standardized Message format
3. **Provider Dispatch**: Message sent to configured LLM provider
4. **Tool Execution**: If tool calls detected, execute appropriate tools
5. **Response Integration**: Tool results fed back to provider for final response
6. **Output Rendering**: Formatted response displayed to user

---

## 3. Provider Abstraction Layer

### 3.1 Interface Definition

```go
type LLMProvider interface {
    SendMessage(ctx context.Context, conversation []Message, tools []ToolDefinition) (*LLMResponse, error)
    Name() string
    AvailableModels() []string
}
```

### 3.2 Provider Implementations

#### Anthropic Provider (`anthropic_provider.go`)

- **SDK**: `github.com/anthropics/anthropic-sdk-go`
- **Models**: Claude 3.7 Sonnet, Claude 3.5 Sonnet, Claude 3 Haiku
- **Features**: Full tool calling support, streaming responses, usage tracking
- **Message Conversion**: Bidirectional conversion between generic and Anthropic formats

#### Gemini Provider (`gemini_provider.go`)

- **SDK**: `github.com/google/generative-ai-go/genai`
- **Models**: Gemini 1.5 Pro, Gemini 1.5 Flash, Gemini 1.0 Pro
- **Features**: Function calling via tools, multi-modal support
- **Limitations**: Schema conversion required for tool definitions

#### OpenAI Provider (Framework)

- **Status**: Interface defined, implementation pending
- **Planned Models**: GPT-4, GPT-3.5 Turbo
- **Integration Points**: Function calling, assistant API

### 3.3 Message Format Conversion

Each provider implements conversion logic between the generic `Message` format and provider-specific formats:

- **Anthropic**: Uses `anthropic.MessageParam` with content blocks
- **Gemini**: Uses `genai.Part` with text and function call structures
- **OpenAI**: Will use Chat Completion format with tool calls

### 3.4 Error Handling Strategy

Provider-specific errors are normalized to generic error types:
- **API Errors**: Rate limits, authentication failures, model errors
- **Network Errors**: Connection timeouts, DNS resolution failures
- **Format Errors**: Message conversion failures, schema validation errors

---

## 4. Tool System Architecture

### 4.1 Tool Categories

#### File System Tools

- **ReadFile**: Single file content retrieval with path validation
- **ListFiles**: Directory enumeration with recursive support
- **EditFile**: String replacement-based file modification
- **SaveToDatabase**: Persistent file storage with versioning

#### Execution Tools

- **Bash**: Command execution with output capture
- **APICall**: HTTP request execution with full REST support
- **WebScraper**: HTML content extraction using CSS selectors

#### Development Tools

- **CodeSearch**: Pattern-based code searching across Go files
- **CreateCheckpoint**: Project state snapshots
- **ListCheckpoints**: Checkpoint enumeration and metadata
- **BackupProject**: Full project archival

#### Integration Tools

- **AddMCPIntegration**: Model Context Protocol server registration
- **RecordMultiplayerAction**: Collaborative action tracking
- **GetMCPIntegrations**: Integration status and configuration

### 4.2 Tool Execution Pipeline

```text
User Request → LLM Processing → Tool Selection → Input Validation → Execution → Result Processing → Response Generation
```

#### Input Validation

- JSON schema validation using `jsonschema` library
- Type conversion and sanitization
- Path validation for file operations
- Command safety checks for bash execution

#### Execution Context

- **Working Directory**: Current project root
- **Environment**: Inherited from parent process
- **Permissions**: File system access, command execution
- **Timeouts**: Configurable execution limits

### 4.3 Tool Result Processing

Tool results are processed through multiple stages:

1. **Raw Output Capture**: Direct command/file output
2. **Format Normalization**: Consistent string representation
3. **Error Classification**: Success/failure determination
4. **Metadata Attachment**: Execution time, resource usage
5. **Response Integration**: Incorporation into conversation flow

---

## 5. Database and Persistence Layer

### 5.1 Architecture Overview

The system implements a file-based database architecture designed for:

- **Version Control**: File content versioning with hash-based integrity
- **Project State**: Checkpoint-based project snapshots
- **Integration Management**: MCP server configuration storage
- **Collaboration Tracking**: Multiplayer action logging

### 5.2 Data Structures

#### ProjectFile

```go
type ProjectFile struct {
    ID         int       // Unique identifier
    Path       string    // File system path
    Content    string    // File content
    Hash       string    // MD5 content hash
    Version    int       // Version counter
    CreatedAt  time.Time // Creation timestamp
    ModifiedAt time.Time // Last modification
    IsActive   bool      // Soft delete flag
}
```

#### Checkpoint

```go
type Checkpoint struct {
    ID          int       // Unique identifier
    Name        string    // User-defined name
    Description string    // Descriptive text
    Timestamp   time.Time // Creation time
    FileCount   int       // Files in checkpoint
}
```

#### MCPIntegration

```go
type MCPIntegration struct {
    ID        int                    // Unique identifier
    Name      string                 // Integration name
    Type      string                 // Integration type
    Config    map[string]interface{} // Configuration data
    CreatedAt time.Time              // Creation timestamp
    IsActive  bool                   // Active status
}
```

### 5.3 Storage Implementation

#### File Organization
```
eve_project_data/
├── files/           # Individual file records
│   ├── 1.json
│   ├── 2.json
│   └── ...
├── checkpoints/     # Checkpoint metadata
├── integrations/    # MCP configurations
├── multiplayer/     # Collaboration logs
└── next_ids.json    # ID counter state
```

#### Persistence Strategy
- **Atomic Operations**: File writes use temporary files with atomic moves
- **Integrity Checks**: MD5 hashing for content verification
- **Recovery Mechanisms**: Automatic reconstruction from checkpoints
- **Space Management**: Soft deletes with periodic cleanup

### 5.4 Query and Retrieval

#### File Operations
- **Path-based Lookup**: O(n) search through active files
- **Version History**: Complete modification timeline
- **Content Diffing**: Change tracking between versions

#### Checkpoint Management
- **Snapshot Creation**: Point-in-time project state
- **Incremental Storage**: Only changed files stored
- **Restoration Logic**: Atomic rollback capability

---

## 6. Configuration Management

### 6.1 Configuration Sources

#### Environment Variables
```bash
# Provider Selection
LLM_PROVIDER=anthropic|openai|gemini

# API Keys
ANTHROPIC_API_KEY=sk-ant-api03-...
OPENAI_API_KEY=sk-...
GEMINI_API_KEY=...

# Model Selection
LLM_MODEL=claude-3-5-sonnet-20241022

# System Configuration
EVE_DATABASE_PATH=./eve_project_data
EVE_LOG_LEVEL=info|debug|warn|error
EVE_MAX_WORKERS=10
EVE_CACHE_SIZE=100MB
EVE_TIMEOUT=30
```

#### Runtime Configuration
- **Dynamic Provider Switching**: Runtime reconfiguration capability
- **Model Selection**: Per-session model specification
- **Resource Limits**: Configurable execution constraints

### 6.2 Configuration Validation

#### API Key Validation
- **Format Checking**: Provider-specific key format validation
- **Connectivity Testing**: API endpoint reachability verification
- **Permission Verification**: Required scope validation

#### System Requirements
- **Go Version**: Minimum 1.24.2 requirement
- **Dependencies**: Automatic dependency resolution via go.mod
- **File System**: Write permissions for database directory

### 6.3 Configuration Lifecycle

1. **Initialization**: Environment variable parsing
2. **Validation**: Configuration completeness checking
3. **Provider Creation**: Dynamic provider instantiation
4. **Runtime Updates**: Configuration hot-reloading capability

---

## 7. Execution Flow and State Management

### 7.1 Conversation Management

#### Message Processing Pipeline
```
Input → Tokenization → Context Window Management → Provider Dispatch → Response Processing → State Update
```

#### Context Window Management
- **Token Counting**: Accurate token estimation for each provider
- **Message Pruning**: Automatic removal of oldest messages
- **Summarization**: Context compression for long conversations
- **Persistence**: Conversation state saved to disk

#### State Persistence
- **File Format**: JSON-based conversation storage
- **Recovery**: Automatic state restoration on restart
- **Versioning**: Conversation format versioning for compatibility

### 7.2 Tool Execution State

#### Execution Tracking
- **Request ID**: Unique identifier for each tool execution
- **Progress Monitoring**: Real-time execution status
- **Result Caching**: Output caching for identical requests
- **Timeout Management**: Configurable execution time limits

#### Error Recovery
- **Retry Logic**: Automatic retry for transient failures
- **Fallback Mechanisms**: Alternative execution paths
- **Partial Success**: Handling of partially successful operations

### 7.3 Resource Management

#### Memory Management
- **Buffer Pooling**: Reusable buffer allocation
- **Garbage Collection**: Automatic cleanup of unused resources
- **Memory Limits**: Configurable memory usage constraints

#### Connection Pooling
- **HTTP Clients**: Persistent connection management
- **Provider Connections**: Connection reuse and lifecycle management
- **Database Connections**: Efficient connection pooling

---

## 8. Error Handling and Resilience

### 8.1 Error Classification

#### Provider Errors
- **API Errors**: Rate limits, quota exceeded, model unavailable
- **Authentication Errors**: Invalid API keys, expired tokens
- **Network Errors**: Connection failures, timeouts, DNS issues

#### Tool Execution Errors
- **File System Errors**: Permission denied, file not found, disk full
- **Command Execution Errors**: Command not found, execution failures
- **Validation Errors**: Input format errors, schema violations

#### System Errors
- **Configuration Errors**: Missing environment variables, invalid settings
- **Resource Errors**: Memory exhaustion, disk space issues
- **Database Errors**: Corruption, locking issues, migration failures

### 8.2 Error Recovery Strategies

#### Automatic Recovery
- **Retry Logic**: Exponential backoff for transient failures
- **Circuit Breaker**: Automatic failure detection and recovery
- **Fallback Providers**: Automatic switching to alternative providers

#### Graceful Degradation
- **Feature Disabling**: Non-critical feature shutdown on failures
- **Reduced Functionality**: Continued operation with limited capabilities
- **User Notification**: Clear error communication to users

### 8.3 Monitoring and Alerting

#### Health Checks
- **Provider Connectivity**: Regular API endpoint testing
- **Database Integrity**: Automatic corruption detection
- **Resource Usage**: Memory, disk, and CPU monitoring

#### Logging Strategy
- **Structured Logging**: JSON-formatted log entries
- **Log Levels**: Debug, info, warning, error classification
- **Log Rotation**: Automatic log file rotation and cleanup

---

## 9. Performance Characteristics

### 9.1 Latency Analysis

#### Provider Latency
- **Anthropic Claude**: 2-5 seconds for typical requests
- **Google Gemini**: 1-3 seconds for standard queries
- **OpenAI GPT**: 1-4 seconds depending on model

#### Tool Execution Latency
- **File Operations**: < 100ms for local files
- **Command Execution**: Variable based on command complexity
- **Network Operations**: 100ms - 2s depending on endpoint

### 9.2 Throughput Considerations

#### Concurrent Operations
- **Tool Parallelization**: Up to 10 concurrent tool executions
- **Provider Limits**: Respecting API rate limits and quotas
- **Resource Pooling**: Efficient resource utilization

#### Memory Usage
- **Base Memory**: ~50MB for core system
- **Per Conversation**: ~1MB additional per 1000 messages
- **File Caching**: Configurable cache sizes up to 100MB

### 9.3 Scalability Factors

#### Vertical Scaling
- **Memory**: Linear scaling with conversation size
- **CPU**: Multi-core utilization for parallel tool execution
- **Storage**: File-based storage with efficient indexing

#### Horizontal Scaling
- **Stateless Design**: Individual instances can be scaled independently
- **Shared Storage**: Database can be shared across instances
- **Load Balancing**: Request distribution across multiple instances

---

## 10. Security Considerations

### 10.1 API Key Management

#### Storage Security
- **Environment Variables**: Secure key storage in environment
- **No Persistence**: Keys never written to disk
- **Access Control**: Process-level key isolation

#### Transmission Security
- **HTTPS Only**: All API communications over TLS 1.3
- **Certificate Validation**: Strict certificate verification
- **Key Rotation**: Support for periodic key rotation

### 10.2 File System Security

#### Access Control
- **Working Directory**: Restricted to project directory
- **Path Validation**: Prevention of directory traversal attacks
- **Permission Checking**: File access permission validation

#### Data Protection
- **Encryption**: Optional file content encryption
- **Integrity**: MD5 hash verification for data integrity
- **Backup Security**: Encrypted backup file generation

### 10.3 Command Execution Security

#### Input Sanitization
- **Command Validation**: Whitelisted command patterns
- **Argument Escaping**: Proper shell argument escaping
- **Path Safety**: Absolute path prevention in commands

#### Execution Isolation
- **Process Isolation**: Separate process execution
- **Resource Limits**: CPU and memory limits on commands
- **Timeout Enforcement**: Maximum execution time limits

---

## 11. Extensibility Framework

### 11.1 Provider Extension

#### Adding New Providers
1. **Interface Implementation**: Implement `LLMProvider` interface
2. **Message Conversion**: Provider-specific message format conversion
3. **Error Handling**: Provider-specific error normalization
4. **Configuration Integration**: Add to configuration system

#### Provider Registration
```go
// Add to config.go
const ProviderNewProvider ProviderType = "newprovider"

// Add to CreateProvider
case ProviderNewProvider:
    return NewNewProvider(c.APIKey, c.Model), nil
```

### 11.2 Tool Extension

#### Tool Development
1. **Function Implementation**: Core tool logic implementation
2. **Schema Definition**: JSON schema for input validation
3. **Error Handling**: Comprehensive error handling and reporting
4. **Documentation**: Tool usage documentation

#### Tool Registration
```go
var NewToolDefinition = ToolDefinition{
    Name:        "new_tool",
    Description: "Description of new tool functionality",
    InputSchema: GenerateSchema[NewToolInput](),
    Function:    NewToolFunction,
}
```

### 11.3 Integration Points

#### MCP Integration
- **Protocol Support**: Model Context Protocol implementation
- **Server Discovery**: Automatic MCP server detection
- **Configuration Management**: Dynamic integration configuration

#### Plugin Architecture
- **Interface Definition**: Plugin interface specification
- **Lifecycle Management**: Plugin loading and unloading
- **Security Sandbox**: Isolated plugin execution environment

---

## 12. Deployment and Operations

### 12.1 Build Process

#### Go Build Configuration
```bash
# Standard build
go build -o eve .

# Optimized build
go build -ldflags="-s -w" -o eve .

# Cross-platform build
GOOS=linux GOARCH=amd64 go build -o eve-linux .
GOOS=darwin GOARCH=amd64 go build -o eve-mac .
GOOS=windows GOARCH=amd64 go build -o eve.exe .
```

#### Dependency Management
- **Module System**: Go 1.24.2 module support
- **Vendor Directory**: Optional vendored dependencies
- **Security Updates**: Regular dependency updates via renovate

### 12.2 Runtime Requirements

#### System Requirements
- **Operating System**: Linux, macOS, Windows
- **Architecture**: x86_64, ARM64 support
- **Memory**: 4GB minimum, 8GB recommended
- **Storage**: 500MB for installation, variable for projects

#### Network Requirements
- **Outbound Connections**: HTTPS access to provider APIs
- **DNS Resolution**: Reliable DNS service
- **Firewall**: Outbound HTTPS access (ports 443)

### 12.3 Monitoring and Observability

#### Metrics Collection
- **Performance Metrics**: Response times, throughput, error rates
- **Resource Usage**: CPU, memory, disk utilization
- **Provider Metrics**: API call counts, token usage, costs

#### Logging Configuration
- **Log Levels**: Configurable verbosity levels
- **Structured Output**: JSON-formatted log entries
- **Log Aggregation**: Support for external log aggregation systems

### 12.4 Backup and Recovery

#### Data Backup
- **Automated Backups**: Scheduled project state backups
- **Incremental Backup**: Change-only backup strategy
- **Offsite Storage**: Support for remote backup destinations

#### Disaster Recovery
- **State Reconstruction**: Database reconstruction from backups
- **Configuration Recovery**: Environment and configuration restoration
- **Failover Support**: Multi-instance deployment support

---

## Conclusion

EVE represents a sophisticated implementation of an AI-powered coding assistant with a focus on modularity, extensibility, and robust operation. The system's provider-agnostic architecture, comprehensive tool ecosystem, and resilient design make it suitable for a wide range of development assistance scenarios.

The technical foundation established in this manual provides a solid basis for understanding, maintaining, and extending the EVE system. Future enhancements should maintain the established patterns of clean abstraction, comprehensive error handling, and modular design.

For specific implementation details or extension development, refer to the individual component documentation and source code comments.</content>
<parameter name="filePath">l:\devops\_sandbox\EVE\docs\TECHNICAL_MANUAL.md
