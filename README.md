


<div align="center">
<h1 style="color:#D291BC;">EVE</h1>
<p><strong style="color:#B38CB4;">Multi-Provider AI Coding Assistant with Tool Integration</strong></p>
<p style="color:#A3A1FF;">Build, debug, and deploy code with AI-powered assistance</p>
<br />

![License](https://img.shields.io/badge/License-MIT-FF69B4?style=for-the-badge)
![Contributions Welcome](https://img.shields.io/badge/Contributions-Welcome-FFC0CB?style=for-the-badge)
![Build Status](https://img.shields.io/badge/Build-Passing-32CD32?style=for-the-badge)

---

## ✨ What is EVE?

🌷 **EVE** is a sophisticated, AI-powered coding assistant that seamlessly integrates with multiple AI providers (Anthropic Claude, Google Gemini, OpenAI) to help you:

- **🔍 Read & Analyze** code files and project structures
- **✏️ Edit & Modify** code with precision and safety
- **⚡ Execute Commands** safely in your terminal
- **🔎 Search Codebase** with powerful pattern matching
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
```
You: Hello, can you help me write a function?
EVE: I'd be happy to help you write a function! What kind of function are you thinking of?
```

### File Analysis
```bash
go run read.go
```
```
You: Read the main.go file and explain what it does
EVE: Let me read that file for you...

[Reads file content]
This appears to be the main entry point for a Go application that...
```

### Code Editing
```bash
go run edit_tool.go
```
```
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

## 🧭 Workshop Path: Learn by Building

| Phase | What to Focus On                                 |
| ----- | ------------------------------------------------ |
| **1** | `chat.go`: API integration and response handling |
| **2** | `read.go`: Tool system, schema generation        |
| **3** | `list_files.go`: Multiple tools, file system     |
| **4** | `bash_tool.go`: Shell execution, error capture   |
| **5** | `edit_tool.go`: File editing, safety checks      |
| **6** | `code_search_tool.go`: Pattern search, ripgrep   |

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

## 🚀 What's Next?

Once you complete the workshop, try building:

* Custom tools (e.g., API caller, web scraper)
* Tool chains (run tools in a sequence)
* Memory features (remember things across sessions)
* A web UI for your agent
* Integration with other AI models

---

## 📦 Summary

This workshop helps you:

* Understand agent architecture
* Learn to build smart assistants
* Grow capabilities step-by-step
* Practice using Claude and Go together

---

Have fun exploring and building your own AI-powered tools! 💻✨

If you have questions or ideas, feel free to fork the repo, open issues, or connect with the community!
