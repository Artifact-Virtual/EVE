


<div style="background: linear-gradient(135deg, #D291BC 0%, #B38CB4 50%, #A3A1FF 100%); padding: 20px; text-align: center; margin-bottom: 20px; border-radius: 10px;">
  <h1 style="color: white; margin: 0; font-size: 3em; text-shadow: 2px 2px 4px rgba(0,0,0,0.3);">🌸 EVE 🌸</h1>
  <p style="color: white; margin: 10px 0 0 0; font-size: 1.2em; opacity: 0.9;">Multi-Provider AI Coding Assistant with Advanced Tool Integration</p>
</div>

<div align="center">
<h1 style="color:#D291BC;">EVE</h1>
<p><strong style="color:#B38CB4;">Multi-Provider AI Coding Assistant with Tool Integration</strong></p>
<p style="color:#A3A1FF;">Build, debug, and deploy code with AI-powered assistance</p>
<br />

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.24.2+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-FF69B4?style=for-the-badge)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Artifact-Virtual/EVE?style=for-the-badge)](https://goreportcard.com/report/github.com/Artifact-Virtual/EVE)
[![GitHub Stars](https://img.shields.io/github/stars/Artifact-Virtual/EVE?style=for-the-badge)](https://github.com/Artifact-Virtual/EVE/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/Artifact-Virtual/EVE?style=for-the-badge)](https://github.com/Artifact-Virtual/EVE/issues)
[![Contributions Welcome](https://img.shields.io/badge/Contributions-Welcome-FFC0CB?style=for-the-badge)](CONTRIBUTING.md)
[![Build Status](https://img.shields.io/badge/Build-Passing-32CD32?style=for-the-badge)](Makefile)

[🚀 Quick Start](#-quick-start) • [📚 Documentation](#-documentation) • [🤝 Contributing](#-contributing) • [📄 License](#-license)

</div>## ✨ What is EVE?

🌷 **EVE** is a sophisticated, AI-powered coding assistant that seamlessly integrates with multiple AI providers (Anthropic Claude, Google Gemini, OpenAI) to help you:

- **🔍 Read & Analyze** code files and project structures
- **✏️ Edit & Modify** code with precision and safety
- **⚡ Execute Commands** safely in your terminal
- **🔎 Search Codebase** with powerful pattern matching
- **🌐 Make API Calls** to external services
- **🕷️ Scrape Webpages** for data extraction
- **🧠 Remember Conversations** across sessions
- **🔄 Switch Providers** on-the-fly for optimal performance

Built with a clean provider abstraction layer, EVE gives you the flexibility to choose the best AI model for your specific coding tasks while maintaining a consistent, powerful interface.

---

## 🚀 Key Features

### 🤖 Multi-Provider Support

- **🧠 Anthropic Claude** - Excellent for complex reasoning and code analysis
- **⚡ Google Gemini** - Fast and efficient for quick tasks
- **🎨 OpenAI GPT** - Versatile for creative coding solutions
- **🔄 Easy Switching** - Change providers with a single environment variable

### 🛠️ Powerful Tool Integration

- **📁 File Operations** - Read, edit, create, and manage files
- **💻 Terminal Commands** - Execute shell commands safely with output capture
- **🔍 Code Search** - Find patterns across your entire codebase using ripgrep
- **📂 Directory Exploration** - Navigate and understand project structure
- **🌐 API Caller** - Make HTTP requests to external services
- **🕷️ Web Scraper** - Extract data from webpages using CSS selectors
- **🧠 Memory Persistence** - Save and load conversation history

### 🏗️ Clean Architecture

- **🔌 Provider Abstraction** - Clean interface for adding new AI providers
- **🧰 Tool System** - Extensible tool framework with JSON schema validation
- **⚙️ Configuration Management** - Environment-based setup with validation
- **📊 Verbose Logging** - Detailed debugging and monitoring capabilities

---

## 🏁 Quick Start

### 📋 Prerequisites

- **Go 1.24.2+** or [devenv](https://devenv.sh/) (recommended for easy setup)
- **At least one AI provider API key**:
  - [Anthropic API Key](https://www.anthropic.com/product/claude) (recommended)
  - [Google Gemini API Key](https://makersuite.google.com/app/apikey) (optional)
  - [OpenAI API Key](https://platform.openai.com/api-keys) (optional)

### 🔧 Setup Environment

#### Option 1: Recommended (using devenv)

```bash
# Clone the repository
git clone https://github.com/ghuntley/how-to-build-a-coding-agent.git
cd how-to-build-a-coding-agent

# Use devenv for reproducible environment (recommended)
devenv shell
```

#### Option 2: Manual setup

```bash
# Clone the repository
git clone https://github.com/ghuntley/how-to-build-a-coding-agent.git
cd how-to-build-a-coding-agent

# Install dependencies manually
go mod tidy
```

### 🔐 Configure API Keys

```bash
# Choose your preferred provider
export LLM_PROVIDER="anthropic"  # Options: anthropic, gemini, openai

# Set your API keys (at least one required)
export ANTHROPIC_API_KEY="your-anthropic-key-here"
export GEMINI_API_KEY="your-gemini-key-here"
export OPENAI_API_KEY="your-openai-key-here"
```

### 🎯 Run EVE

```bash
# Run the full-featured multi-provider agent
go run agent.go

# Or explore individual tools
go run chat.go          # Basic chat interface
go run read.go          # With file reading capabilities
go run list_files.go    # With directory exploration
go run bash_tool.go     # With terminal command execution
go run edit_tool.go     # With file editing capabilities
go run code_search_tool.go  # With code search functionality
```

---

## 💬 Example Usage

### Basic Chat
```bash
go run chat.go
```

```text
You: Hello, can you help me write a function?
EVE: I'd be happy to help you write a function! What kind of function are you thinking of?
```

### File Analysis

```bash
go run read.go
```

```text
You: Read the main.go file and explain what it does
EVE: Let me read that file for you...

[Reads file content]
This appears to be the main entry point for a Go application that...
```

### Code Editing

```bash
go run edit_tool.go
```

```text
You: Add error handling to the database connection function
EVE: I'll help you add proper error handling. Let me first read the current implementation...

[Analyzes code and makes targeted edits]
I've added comprehensive error handling to your database connection function.
```

### Terminal Operations
```bash
go run bash_tool.go
```
```
You: Run the test suite and show me the results
EVE: I'll execute the test command for you.

[Running: go test ./...]
Test results: 15 passed, 2 failed...
```

---

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐
│     EVE         │    │  Provider        │
│   Agent         │◄──►│  Abstraction     │
│                 │    │  Layer           │
└─────────────────┘    └──────────────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌──────────────────┐
│   Tool System   │    │  Anthropic       │
│                 │    │  Claude          │
│ • File Reader   │    ├──────────────────┤
│ • Code Editor   │    │  Google Gemini   │
│ • Terminal      │    ├──────────────────┤
│ • Search        │    │  OpenAI GPT      │
└─────────────────┘    └──────────────────┘
```

### Core Components

- **`agent.go`** - Main multi-provider agent with tool orchestration
- **`llm.go`** - Provider interfaces and shared types with JSON schema generation
- **`config.go`** - Configuration management with environment variable validation
- **Provider implementations** - Anthropic, Gemini, OpenAI with unified interfaces
- **Tool implementations** - File operations, terminal commands, code search

### How It Works

Each agent follows this **event loop**:

1. **📥 Waits** for your input
2. **📤 Sends** it to your chosen AI provider (Claude, Gemini, or GPT)
3. **🤔 AI responds** directly or requests tool usage
4. **⚡ Executes** tools (read files, run commands, etc.)
5. **📨 Returns** results back to the AI
6. **💬 Provides** final answer to you

---

## 🔧 Development

### Building
```bash
make build    # Build all binaries
make fmt      # Format code
make check    # Run linting and checks
make clean    # Clean build artifacts
```

### Adding New Providers
1. Implement the `LLMProvider` interface in `llm.go`
2. Add provider type to `config.go`
3. Update environment variable handling
4. Add to the provider factory in `config.go`

### Adding New Tools
1. Define tool schema in `llm.go` using struct tags
2. Implement tool execution logic
3. Register tool with the agent
4. Update tool definitions and help text

### Testing
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestAnthropicProvider ./...
```

---

## � Roadmap

- [ ] 🌐 **Web UI interface** - Browser-based interface for EVE
- [ ] 🔌 **Plugin system** - Custom tools and extensions
- [ ] 🧠 **Memory features** - Conversation persistence and context
- [ ] 💻 **IDE integration** - VS Code, IntelliJ, and other editors
- [ ] 📊 **Advanced analytics** - Code analysis and insights
- [ ] 🌍 **Multi-language support** - Python, JavaScript, Rust, etc.
- [ ] 🔒 **Security enhancements** - Sandboxed execution
- [ ] 📈 **Performance optimization** - Caching and parallel processing

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### How to Contribute
1. **🍴 Fork** the repository
2. **🌿 Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **💻 Make** your changes
4. **✅ Add** tests if applicable
5. **📤 Submit** a pull request

### Development Setup
```bash
# Fork and clone
git clone https://github.com/yourusername/how-to-build-a-coding-agent.git
cd how-to-build-a-coding-agent

# Set up development environment
devenv shell

# Run tests
make check
```

---

## 📄 License

**EVE** is open source software licensed under the [MIT License](LICENSE).

---

## 🙏 Acknowledgments

- **Anthropic** for the Claude API
- **Google** for the Gemini API
- **OpenAI** for the GPT API
- **BurntSushi** for ripgrep (code search)
- **The Go Community** for excellent tooling

---

<div align="center">
  <p><strong>Built with ❤️ for developers, by developers</strong></p>
  <p>
    <a href="#-what-is-eve">What is EVE</a> •
    <a href="#-key-features">Features</a> •
    <a href="#-quick-start">Quick Start</a> •
    <a href="#-architecture">Architecture</a> •
    <a href="#-contributing">Contributing</a>
  </p>
  <br />
  <p><em>Transform your coding workflow with AI-powered assistance</em></p>
</div>
```

---

## 🏁 Start with the Basics

### 1. `chat.go` — Basic Chat

A simple chatbot that talks to your chosen AI provider (Anthropic Claude by default).

```bash
go run chat.go
```

- ➡️ Try: "Hello!"
- ➡️ Add `--verbose` to see detailed logs

---

## 🛠️ Add Tools (One Step at a Time)

### 2. `read.go` — Read Files

Now your AI can read files from your computer.

```bash
go run read.go
```

- ➡️ Try: "Read fizzbuzz.js"

---

### 3. `list_files.go` — Explore Folders

Lets Claude look around your directory.

```bash
go run list_files.go
```

* ➡️ Try: “List all files in this folder”
* ➡️ Try: “What’s in fizzbuzz.js?”

---

### 4. `bash_tool.go` — Run Shell Commands

Allows Claude to run safe terminal commands.

```bash
go run bash_tool.go
```

* ➡️ Try: “Run git status”
* ➡️ Try: “List all .go files using bash”

---

### 5. `edit_tool.go` — Edit Files

Claude can now **modify code**, create files, and make changes.

```bash
go run edit_tool.go
```

* ➡️ Try: “Create a Python hello world script”
* ➡️ Try: “Add a comment to the top of fizzbuzz.js”

---

### 6. `code_search_tool.go` — Search Code

Use pattern search (powered by [ripgrep](https://github.com/BurntSushi/ripgrep)).

```bash
go run code_search_tool.go
```

* ➡️ Try: “Find all function definitions in Go files”
* ➡️ Try: “Search for TODO comments”

---

## � 7. `agent.go` — Multi-Provider Generic Agent

The final version combines all tools with **multi-provider support**. Switch between Anthropic Claude, Google Gemini, and OpenAI seamlessly!

```bash
# Use Anthropic (default)
export LLM_PROVIDER="anthropic"
go run agent.go

# Switch to Gemini
export LLM_PROVIDER="gemini"
go run agent.go

# Switch to OpenAI
export LLM_PROVIDER="openai"
go run agent.go
```

- ➡️ Try: "Read the riddle.txt file and solve the puzzle"
- ➡️ Try: "Create a new Go file with a hello world function"
- ➡️ Try: "Search for all TODO comments in the codebase"

---

## �🧪 Sample Files (Already Included)

1. `fizzbuzz.js`: for file reading and editing
1. `riddle.txt`: a fun text file to explore
1. `AGENT.md`: info about the project environment

---

## 🐞 Troubleshooting

**API key not working?**

* Make sure it’s exported: `echo $ANTHROPIC_API_KEY`
* Check your quota on [Anthropic’s dashboard](https://www.anthropic.com)

**Go errors?**

* Run `go mod tidy`
* Make sure you’re using Go 1.24.2 or later

**Tool errors?**

* Use `--verbose` for full error logs
* Check file paths and permissions

**Environment issues?**

* Use `devenv shell` to avoid config problems

---

## 💡 How Tools Work (Under the Hood)

Tools are like plugins. You define:

* **Name** (e.g., `read_file`)
* **Input Schema** (what info it needs)
* **Function** (what it does)

Example tool definition in Go:

```go
var ToolDefinition = ToolDefinition{
    Name:        "read_file",
    Description: "Reads the contents of a file",
    InputSchema: GenerateSchema[ReadFileInput](),
    Function:    ReadFile,
}
```

Schema generation uses Go structs — so it’s easy to define and reuse.


---

## 🛠️ Developer Environment (Optional)

If you use [`devenv`](https://devenv.sh/), it gives you:

* Go, Node, Python, Rust, .NET
* Git and other dev tools

```bash
devenv shell   # Load everything
devenv test    # Run checks
hello          # Greeting script
```

---

## 🏗️ Advanced Features

### Database Integration

EVE now includes a sophisticated file-based database system for project management:

- **📁 Project Files**: Version-controlled file storage with metadata
- **📸 Checkpoints**: Create and restore project snapshots
- **🔌 MCP Integrations**: Connect to external Model Context Protocol servers
- **👥 Multiplayer**: Track collaborative development actions

### Built-in Tools

EVE comes with 7 powerful integrated tools:

1. **💾 SaveToDatabase** - Store files with automatic versioning
2. **📸 CreateCheckpoint** - Create project snapshots for backup
3. **📋 ListCheckpoints** - View all available checkpoints
4. **🔌 AddMCPIntegration** - Register external MCP servers
5. **📊 GetMCPIntegrations** - List MCP connections with status
6. **👥 RecordMultiplayerAction** - Track team collaboration
7. **📈 GetMultiplayerHistory** - View collaboration timeline

### Web Interface

EVE includes a modern web interface with:

- **🌐 Responsive Design**: Works on desktop and mobile
- **🌙 Dark Theme**: Beautiful purple/pink gradient UI
- **⚡ Interactive Terminal**: Web-based command execution
- **📊 Real-time Monitoring**: Live system status and logs

---

## 🚀 Getting Started with Advanced Features

### Database Setup

```bash
# Initialize project database
./eve --init-project --name "my-awesome-project"

# Save your first file
./eve --tool save-file --path "main.go" --content "package main\n\nfunc main() {\n\tfmt.Println(\"Hello, EVE!\")\n}"

# Create a checkpoint
./eve --tool create-checkpoint --name "initial-setup" --description "Project initialization"
```

### MCP Integration

```bash
# Add an MCP server
./eve --tool add-mcp --name "github-mcp" --server-type "http" --config '{"url":"https://api.github.com","auth":"bearer-token"}'

# List all integrations
./eve --tool get-mcp-integrations
```

### Multiplayer Collaboration

```bash
# Enable multiplayer mode
./eve --enable-multiplayer --project "team-project"

# Record collaborative actions
./eve --tool record-action --user "alice" --action "code-review" --data '{"file":"api.go","status":"approved"}'

# View team activity
./eve --tool get-multiplayer-history --limit 20
```

### Launch Web Interface

```bash
# Start the web server
./eve --web-server --port 8080

# Access at http://localhost:8080
```

---

## 📊 Performance & Benchmarks

### Supported Platforms

- **🐧 Linux** (x86_64, ARM64)
- **🍎 macOS** (Intel, Apple Silicon)
- **🪟 Windows** (x86_64, ARM64)
- **🐳 Docker** containers

### System Requirements

- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 500MB for installation, varies by project size
- **Network**: Internet connection for AI provider APIs

### Performance Tips

- Use Claude for complex coding tasks
- Gemini for fast, simple operations
- Enable caching for repeated operations
- Use checkpoints for large projects

---

## 🔧 Configuration Options

### Environment Variables

```bash
# Core Settings
export EVE_DEFAULT_PROVIDER="claude"        # anthropic, gemini, openai
export EVE_DATABASE_PATH="./eve_project_data"
export EVE_LOG_LEVEL="info"                 # debug, info, warn, error

# Performance
export EVE_MAX_WORKERS="10"
export EVE_CACHE_SIZE="100MB"
export EVE_TIMEOUT="30"

# Security
export EVE_ENCRYPTION_KEY="your-key"
export EVE_API_KEY_ROTATION="7d"

# Advanced Features
export EVE_ENABLE_MULTIPLAY="true"
export EVE_MCP_AUTO_DISCOVER="true"
export EVE_WEB_INTERFACE="true"
```

### Configuration File

Create `eve_config.yaml`:

```yaml
provider:
  default: "claude"
  fallback: ["gemini", "openai"]

database:
  type: "file"
  path: "./eve_project_data"
  compression: true
  encryption: true

performance:
  max_workers: 10
  cache_size: "100MB"
  timeout: 30

security:
  encryption_key: "your-key"
  api_key_rotation: "7d"
  access_control: true

features:
  multiplayer: true
  mcp_integration: true
  web_interface: true
```

---

## 🐛 Troubleshooting Guide

### Common Issues

### Build Failures

```bash
# Clean and rebuild
go clean
go mod tidy
go build -o eve .

# Check Go version
go version  # Should be 1.24.2+
```

### API Connection Issues

```bash
# Test API connectivity
curl -H "Authorization: Bearer $ANTHROPIC_API_KEY" \
     https://api.anthropic.com/v1/messages

# Check API key format
echo $ANTHROPIC_API_KEY | head -c 10  # Should start with sk-ant-
```

### Database Problems

```bash
# Reset database
rm -rf eve_project_data/
# Will be recreated on next run

# Check permissions
ls -la eve_project_data/
```

### Memory Issues

```bash
# Monitor memory usage
./eve --memory-profile

# Limit memory usage
export EVE_MAX_MEMORY="2GB"
```

### Debug Mode

```bash
# Run with verbose logging
./eve --verbose --debug

# Enable API request logging
export EVE_DEBUG_API="true"
```

---

## 🌟 Use Cases & Examples

### Web Development

```bash
# Create a React component
./eve --tool create-file --path "src/components/Button.jsx" --content "const Button = ({ children }) => <button>{children}</button>;"

# Set up Express server
./eve --tool save-file --path "server.js" --content "const express = require('express'); const app = express(); app.listen(3000);"
```

### Data Science

```bash
# Create Python analysis script
./eve --tool save-file --path "analysis.py" --content "import pandas as pd\nimport matplotlib.pyplot as plt\n\ndf = pd.read_csv('data.csv')\ndf.plot()\nplt.show()"

# Run Jupyter notebook
./eve --bash "jupyter notebook"
```

### DevOps & Infrastructure

```bash
# Create Dockerfile
./eve --tool save-file --path "Dockerfile" --content "FROM golang:1.21\nCOPY . .\nRUN go build -o app .\nCMD [\"./app\"]"

# Set up CI/CD pipeline
./eve --tool save-file --path ".github/workflows/ci.yml" --content "name: CI\non: push\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n    - uses: actions/checkout@v3\n    - run: go test ./..."
```

---

## 🤝 Community & Support

### Getting Help

- **📖 Documentation**: Check `/docs` directory and [BUILD_GUIDE.md](BUILD_GUIDE.md)
- **🐛 Issues**: [GitHub Issues](https://github.com/Artifact-Virtual/EVE/issues) for bugs
- **💬 Discussions**: [GitHub Discussions](https://github.com/Artifact-Virtual/EVE/discussions) for questions
- **💬 Discord**: Join our community server for real-time help

### Contributing

We welcome contributions of all kinds!

**Ways to Contribute:**

- 🐛 Report bugs and issues
- ✨ Suggest new features
- 📝 Improve documentation
- 🔧 Submit code changes
- 🧪 Add tests
- 🌍 Translate documentation

**Development Workflow:**

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Run `make check` to ensure quality
6. Submit a pull request

### Code of Conduct

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) to understand our community standards.

---

## 📈 Roadmap & Future Plans

### Q1 2024

- [ ] 🌐 Enhanced web interface with collaborative editing
- [ ] 🔌 Plugin system for custom tools
- [ ] 🧠 Advanced memory and context management

### Q2 2024

- [ ] 💻 IDE integrations (VS Code, IntelliJ)
- [ ] 📊 Code analytics and insights
- [ ] 🌍 Multi-language support expansion

### Q3 2024

- [ ] 🔒 Enterprise security features
- [ ] 📈 Performance optimizations
- [ ] 🏢 Team management features

### Q4 2024

- [ ] 📱 Mobile companion app
- [ ] ☁️ Cloud deployment options
- [ ] 🤖 AI model fine-tuning

### Long-term Vision

- **🧠 AGI Integration**: Connect with advanced AI systems
- **🌐 Decentralized**: P2P collaboration without central servers
- **🔮 Predictive**: AI-powered code suggestions and bug detection
- **🎨 Creative**: AI-assisted design and architecture planning

---

## 📄 License & Legal

**EVE** is open source software licensed under the [MIT License](LICENSE).

### Third-party Licenses

- **Anthropic Claude API**: Subject to Anthropic's Terms of Service
- **Google Gemini API**: Subject to Google's Terms of Service
- **OpenAI API**: Subject to OpenAI's Terms of Service
- **ripgrep**: MIT License (for code search functionality)

### Privacy & Data

- API keys are stored locally and never transmitted except to AI providers
- Project data is stored locally in the `eve_project_data/` directory
- No telemetry or usage data is collected without explicit consent

---

## 🙏 Acknowledgments & Thanks

### Core Technologies

- **🤖 Anthropic** for Claude API - Exceptional reasoning capabilities
- **⚡ Google** for Gemini API - Fast and efficient AI processing
- **🎨 OpenAI** for GPT API - Versatile and creative AI solutions
- **🔍 BurntSushi** for ripgrep - Blazingly fast code search
- **🐹 The Go Team** for an amazing programming language

### Community Contributors

Special thanks to all our contributors who help make EVE better every day!

### Inspiration

EVE was inspired by the need for flexible, powerful AI-assisted development tools that work seamlessly across different AI providers and development environments.

---

## 🎯 Quick Reference

### Essential Commands

```bash
# Start EVE
go run agent.go

# Build for production
go build -ldflags="-s -w" -o eve .

# Run with web interface
./eve --web-server --port 8080

# Initialize project
./eve --init-project --name "my-project"
```

### Tool Shortcuts

```bash
# Save file
./eve --tool save-file --path "file.go" --content "code"

# Create checkpoint
./eve --tool create-checkpoint --name "v1.0"

# Add MCP server
./eve --tool add-mcp --name "server" --config '{"url":"..."}'

# View history
./eve --tool get-multiplayer-history --limit 10
```

### Environment Setup

```bash
# Required
export ANTHROPIC_API_KEY="your-key"

# Optional
export EVE_DEFAULT_PROVIDER="claude"
export EVE_WEB_INTERFACE="true"
```

---

<div align="center">

## 🌸 **Happy Coding with EVE!** 🌸

<p><em>Transform your development workflow with the power of AI</em></p>

### 🚀 **Ready to Get Started?**

```bash
git clone https://github.com/Artifact-Virtual/EVE.git
cd EVE
devenv shell
go run agent.go
```

### 📞 **Need Help?**
- 📖 [Full Documentation](BUILD_GUIDE.md)
- 🐛 [Report Issues](https://github.com/Artifact-Virtual/EVE/issues)
- 💬 [Community Discussions](https://github.com/Artifact-Virtual/EVE/discussions)

---

---

<div align="center">

## 🌸 **Happy Coding with EVE!** 🌸

_Transform your development workflow with the power of AI_

### 🚀 **Ready to Get Started?**

```bash
git clone https://github.com/Artifact-Virtual/EVE.git
cd EVE
devenv shell
go run agent.go
```

### 📞 **Need Help?**

- 📖 [Full Documentation](BUILD_GUIDE.md)
- 🐛 [Report Issues](https://github.com/Artifact-Virtual/EVE/issues)
- 💬 [Community Discussions](https://github.com/Artifact-Virtual/EVE/discussions)

---

**Built with ❤️ by the EVE development team**  
*Empowering developers with AI-assisted coding since 2024*

</div>
