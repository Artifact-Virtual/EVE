# Multi-Provider LLM Abstraction

This enhanced version of the coding agent now supports **multiple LLM providers** through a clean abstraction layer. You can easily switch between Anthropic Claude, OpenAI GPT, Google Gemini, and other providers.

## 🏗️ Architecture Overview

The abstraction consists of several key components:

### Core Interfaces (`llm.go`)
- **`LLMProvider`**: Main interface that all providers must implement
- **`Message`**: Generic message format
- **`ContentBlock`**: Flexible content representation
- **`ToolUse`/`ToolResult`**: Tool interaction structures

### Provider Implementations
- **`anthropic_provider.go`**: Full Anthropic Claude implementation
- **`openai_provider.go`**: OpenAI GPT implementation (framework ready)
- **`gemini_provider.go`**: Google Gemini (to be implemented)

### Configuration System (`config.go`)
- Environment-based provider selection
- Automatic API key detection
- Model selection support

## 🚀 Quick Start

### 1. Set Environment Variables

```bash
# Choose your provider
export LLM_PROVIDER=anthropic  # or openai, gemini

# Set API keys
export ANTHROPIC_API_KEY="your-anthropic-key"
export OPENAI_API_KEY="your-openai-key"
export GEMINI_API_KEY="your-gemini-key"

# Optional: specify model
export LLM_MODEL="claude-3-5-sonnet-20241022"
```

### 2. Run the Generic Agent

```bash
go run agent.go
```

The agent will automatically detect your configuration and use the appropriate provider!

## 🔧 Adding New Providers

### Step 1: Implement the LLMProvider Interface

```go
type YourProvider struct {
    client *your_sdk.Client
    model  string
}

func (p *YourProvider) SendMessage(ctx context.Context, conversation []Message, tools []ToolDefinition) (*LLMResponse, error) {
    // Convert generic messages to provider-specific format
    // Make API call
    // Convert response back to generic format
    // Return LLMResponse
}

func (p *YourProvider) Name() string {
    return "Your Provider Name"
}

func (p *YourProvider) AvailableModels() []string {
    return []string{"model1", "model2"}
}
```

### Step 2: Add to Configuration

```go
// In config.go
const (
    ProviderYourProvider ProviderType = "yourprovider"
)

// In CreateProvider method
case ProviderYourProvider:
    return NewYourProvider(c.APIKey, c.Model), nil
```

### Step 3: Handle API Key Environment Variable

```go
// In NewConfigFromEnv
case ProviderYourProvider:
    apiKey = os.Getenv("YOURPROVIDER_API_KEY")
```

## 📋 Supported Providers

| Provider | Status | Environment Variable | Dependencies |
|----------|--------|---------------------|--------------|
| Anthropic Claude | ✅ Complete | `ANTHROPIC_API_KEY` | Built-in |
| OpenAI GPT | 🔄 Framework Ready | `OPENAI_API_KEY` | `github.com/sashabaranov/go-openai` |
| Google Gemini | 📝 Planned | `GEMINI_API_KEY` | Google AI SDK |

## 🔄 Message Format Conversion

Each provider needs to convert between the generic format and provider-specific formats:

### Anthropic Claude
- Generic `Message` → `anthropic.MessageParam`
- `anthropic.Message` → Generic `LLMResponse`
- Tool calls use `anthropic.ToolUseBlock`

### OpenAI GPT
- Generic `Message` → `openai.ChatCompletionMessage`
- `openai.ChatCompletion` → Generic `LLMResponse`
- Tool calls use `openai.ToolCall`

### Google Gemini
- Generic `Message` → Gemini `Content`
- Gemini response → Generic `LLMResponse`
- Tool calls use Gemini function calling

## 🛠️ Tool Support

The abstraction maintains full tool support across providers:

- **Function Calling**: All providers support tool/function calling
- **Schema Generation**: Uses existing `GenerateSchema[T any]()` function
- **Error Handling**: Consistent error handling across providers
- **Result Processing**: Unified tool result processing

## 📊 Usage Tracking

All providers return standardized usage information:

```go
type Usage struct {
    PromptTokens     int `json:"prompt_tokens"`
    CompletionTokens int `json:"completion_tokens"`
    TotalTokens      int `json:"total_tokens"`
}
```

## 🎯 Example Usage

```go
// Automatic provider selection
config, _ := NewConfigFromEnv()
provider, _ := config.CreateProvider()

// Use with existing tools
agent := NewGenericAgent(provider, getUserInput, tools, verbose)
agent.Run(context.Background())
```

## 🔍 Debugging

Enable verbose logging to see provider-specific details:

```bash
go run agent.go --verbose
```

This will show:
- Provider initialization
- Message conversion details
- API call timing
- Tool execution flow
- Error details with context

## 🚀 Benefits of This Abstraction

1. **Provider Agnostic**: Switch providers without changing agent logic
2. **Easy Extension**: Add new providers by implementing one interface
3. **Unified API**: Same agent code works with any supported provider
4. **Tool Compatibility**: All tools work across providers
5. **Configuration Driven**: Environment-based provider selection
6. **Future Proof**: Easy to add new LLM providers as they emerge

## 📈 Next Steps

1. **Complete OpenAI Implementation**: Add the OpenAI SDK dependency and finish the implementation
2. **Add Gemini Support**: Implement Google Gemini provider
3. **Add More Providers**: Support for Grok, Mistral, etc.
4. **Provider-specific Features**: Model-specific optimizations
5. **Load Balancing**: Distribute requests across multiple providers
6. **Fallback Support**: Automatic fallback to alternative providers

This abstraction makes your coding agent truly provider-agnostic and future-ready! 🎉
