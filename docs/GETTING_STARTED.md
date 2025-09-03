# 🚀 EVE Getting Started Handbook

Welcome to **EVE** - your AI-powered coding assistant! This handbook will guide you through setting up and running EVE on your system. Whether you're a developer looking to enhance your workflow or someone exploring AI-assisted coding, this guide will get you up and running quickly.

## 📋 Table of Contents

- Quick Overview
- Prerequisites
- Installation Methods
- Configuration
- Building EVE
- Running EVE
- Basic Usage
- Troubleshooting
- Next Steps

---

## 🌟 Quick Overview

EVE is a multi-provider AI coding assistant that integrates with:

- **Anthropic Claude** - Excellent for complex reasoning and code analysis
- **Google Gemini** - Fast and efficient for quick tasks
- **OpenAI GPT** - Versatile for creative coding solutions

EVE provides powerful tools for:

- Reading and analyzing code files
- Editing code with precision
- Executing terminal commands safely
- Searching codebases with pattern matching
- Making API calls to external services

---

## 📋 Prerequisites

Before installing EVE, ensure you have the following:

### Required Software

- **Go 1.24.2 or later** - The core runtime for EVE
- **Git** - For cloning the repository
- **At least one AI provider API key** (see Configuration section)

### Optional but Recommended

- **[devenv](https://devenv.sh/)** - For reproducible development environment (highly recommended)
- **ripgrep** - For advanced code searching capabilities

### System Requirements

- **Operating System**: Linux, macOS, or Windows (with WSL)
- **RAM**: Minimum 4GB, recommended 8GB+
- **Storage**: 500MB free space for installation and dependencies

---

## 🛠️ Installation Methods

Choose one of the following installation methods based on your preference:

### Method 1: Recommended - Using devenv (Easiest)

devenv provides a reproducible development environment with all dependencies pre-configured.

#### Step 1: Install devenv

```bash
# On macOS (using Homebrew)
brew install devenv

# On Linux
curl -fsSL https://devenv.sh/install | bash

# On Windows (using Chocolatey)
choco install devenv
```

#### Step 2: Clone and Setup EVE

```bash
# Clone the repository
git clone https://github.com/Artifact-Virtual/EVE.git
cd EVE

# Enter the development environment
devenv shell
```

That's it! devenv will automatically:

- Install Go 1.24.2+
- Set up all required dependencies
- Configure your environment variables
- Install additional tools (ripgrep, etc.)

### Method 2: Manual Installation

If you prefer manual setup or can't use devenv:

#### Step 1: Install Go

Download and install Go 1.24.2+ from [golang.org](https://golang.org/dl/)

Verify installation:

```bash
go version
# Should show: go version go1.24.2 or later
```

#### Step 2: Clone the Repository

```bash
git clone https://github.com/Artifact-Virtual/EVE.git
cd EVE
```

#### Step 3: Install Dependencies

```bash
# Download and install Go modules
go mod tidy

# Verify dependencies
go mod verify
```

#### Step 4: Install Additional Tools (Optional)

```bash
# Install ripgrep for advanced code search
# On Ubuntu/Debian:
sudo apt-get install ripgrep

# On macOS:
brew install ripgrep

# On Windows:
winget install BurntSushi.ripgrep.MSVC
```

---

## ⚙️ Configuration

EVE requires configuration through environment variables. You need at least one AI provider API key.

### Step 1: Choose Your AI Provider

EVE supports three AI providers. Choose based on your needs:

| Provider | Best For | Setup Difficulty | Cost |
|----------|----------|------------------|------|
| **Anthropic Claude** | Complex reasoning, code analysis | Easy | Moderate |
| **Google Gemini** | Fast tasks, general coding | Easy | Low |
| **OpenAI GPT** | Creative solutions, versatility | Easy | Variable |

### Step 2: Get API Keys

#### Anthropic Claude (Recommended)

1. Visit [Anthropic Console](https://console.anthropic.com/)
2. Create an account or sign in
3. Navigate to API Keys section
4. Create a new API key
5. Copy the key (format: `sk-ant-api03-...`)

#### Google Gemini

1. Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Sign in with Google account
3. Create a new API key
4. Copy the key

#### OpenAI GPT

1. Visit [OpenAI Platform](https://platform.openai.com/api-keys)
2. Create an account or sign in
3. Navigate to API Keys
4. Create a new secret key
5. Copy the key (format: `sk-...`)

### Step 3: Set Environment Variables

Create a `.env` file in the EVE directory or set environment variables directly:

```bash
# Choose your provider (anthropic, gemini, or openai)
export LLM_PROVIDER="anthropic"

# Set your API key (replace with your actual key)
export ANTHROPIC_API_KEY="your-anthropic-key-here"
export GEMINI_API_KEY="your-gemini-key-here"
export OPENAI_API_KEY="your-openai-key-here"

# Optional: Specify model (defaults will be used if not set)
export LLM_MODEL="claude-3-5-sonnet-20241022"  # For Anthropic
# export LLM_MODEL="gemini-pro"                 # For Gemini
# export LLM_MODEL="gpt-4"                      # For OpenAI
```

### Step 4: Verify Configuration

Test your configuration:

```bash
# Check if variables are set
echo $LLM_PROVIDER
echo $ANTHROPIC_API_KEY  # (or your chosen provider's key)

# Should show your configured values
```

---

## 🔨 Building EVE

Once configured, build EVE using the provided Makefile:

### Build All Components

```bash
# Build all EVE binaries
make build

# Or build everything with formatting and checks
make all
```

This will create several executable files:

- `agent` - Full-featured multi-provider agent
- `chat` - Basic chat interface
- `read` - File reading capabilities
- `list_files` - Directory exploration
- `bash_tool` - Terminal command execution
- `edit_tool` - File editing capabilities

### Individual Build Commands

```bash
# Build specific components
go build -o agent agent.go
go build -o chat chat.go
go build -o read read.go
go build -o list_files list_files.go
go build -o bash_tool bash_tool.go
go build -o edit_tool edit_tool.go
```

### Verify Build

```bash
# Check if binaries were created
ls -la agent chat read list_files bash_tool edit_tool

# Should show executable files (green on Linux/macOS)
```

---

## ▶️ Running EVE

Now you're ready to run EVE! Start with the basic chat interface to get familiar.

### Basic Chat Interface

```bash
# Start the basic chat interface
./chat

# Or using go run
go run chat.go
```

You should see:

```text
🌸 EVE - AI Coding Assistant 🌸
Provider: anthropic
Model: claude-3-5-sonnet-20241022

You: [Your message here]
```

### Full-Featured Agent

For the complete EVE experience with all tools:

```bash
# Run the full agent
./agent

# Or using go run
go run agent.go
```

### Individual Tools

Test specific capabilities:

```bash
# File reading and analysis
./read

# Directory exploration
./list_files

# Terminal command execution
./bash_tool

# File editing
./edit_tool
```

---

## 💬 Basic Usage

Let's walk through some basic interactions with EVE:

### Example 1: Simple Chat

```text
You: Hello EVE, can you help me write a Go function?

EVE: Hello! I'd be happy to help you write a Go function. What kind of function are you looking for? Please provide some details about what it should do.
```

### Example 2: File Analysis

```text
You: Can you read the README.md file and summarize what EVE does?

EVE: I'll read the README.md file for you...

[Reading file contents...]

Based on the README.md, EVE is a multi-provider AI coding assistant that integrates with Anthropic Claude, Google Gemini, and OpenAI. It provides powerful tools for code analysis, editing, terminal operations, and more.
```

### Example 3: Code Editing

```text
You: Add error handling to this function: func connect() { db, _ := sql.Open("sqlite3", "test.db") }

EVE: I'll help you add proper error handling to your database connection function. Let me analyze the current code and make the necessary improvements...

[Making targeted edits to add error handling]
```

### Example 4: Terminal Commands

```text
You: Run the Go tests and show me the results

EVE: I'll execute the test command for you.

[Running: go test ./... -v]

Test Results:
PASS: 15 tests passed
FAIL: 2 tests failed (details below)
...
```

---

## 🔧 Troubleshooting

### Common Issues and Solutions

#### 1. "API key not found" Error

**Problem**: EVE can't find your API key
**Solution**:

```bash
# Check if environment variable is set
echo $ANTHROPIC_API_KEY

# If empty, set it again
export ANTHROPIC_API_KEY="your-key-here"
```

#### 2. "Unsupported provider" Error

**Problem**: Invalid provider specified
**Solution**:

```bash
# Check valid options
export LLM_PROVIDER="anthropic"  # Options: anthropic, gemini, openai
```

#### 3. "Command not found" Error

**Problem**: EVE binaries not built or not in PATH
**Solution**:

```bash
# Build the binaries
make build

# Run with full path
./agent
```

#### 4. "Go version too old" Error

**Problem**: Go version is below 1.24.2
**Solution**:

```bash
# Check current version
go version

# Update Go to 1.24.2+ from https://golang.org/dl/
```

#### 5. Network/Proxy Issues

**Problem**: Can't connect to AI provider APIs
**Solution**:

```bash
# Check internet connection
ping google.com

# If behind proxy, set proxy environment variables
export HTTP_PROXY="http://proxy.company.com:8080"
export HTTPS_PROXY="http://proxy.company.com:8080"
```

#### 6. Permission Denied

**Problem**: Can't execute binaries
**Solution**:

```bash
# Make binaries executable
chmod +x agent chat read list_files bash_tool edit_tool
```

### Getting Help

If you encounter issues not covered here:

1. **Check the logs** - EVE provides detailed logging
2. **Verify your setup** - Run `make check` to validate your installation
3. **Test with simple commands** - Start with basic chat before full agent
4. **Check API key validity** - Ensure your API keys are active and have credits

---

## 🎯 Next Steps

Now that you have EVE running, here are some recommended next steps:

### 1. Explore Advanced Features

- Try the full agent with all tools enabled
- Experiment with different AI providers
- Use EVE for real coding tasks

### 2. Customize Your Setup

- Configure your preferred AI model
- Set up multiple providers for comparison
- Create aliases for common commands

### 3. Learn More

- Read the [Technical Manual](TECHNICAL_MANUAL.md) for in-depth understanding
- Check out example usage in the `examples/` directory
- Join the community discussions

### 4. Contribute

- Report bugs or request features
- Submit pull requests for improvements
- Help improve documentation

---

## 📞 Support

- **Documentation**: [Technical Manual](TECHNICAL_MANUAL.md)
- **Issues**: [GitHub Issues](https://github.com/Artifact-Virtual/EVE/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Artifact-Virtual/EVE/discussions)

Happy coding with EVE! 🚀
