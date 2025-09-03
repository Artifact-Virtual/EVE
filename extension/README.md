The EVE VS Code extension integrates the EVE Go CLI for chat, planning, and applying edits directly in your editor. It supports both stdio and HTTP modes for the CLI daemon.

## Installation

- Install the generated `.vsix` file in VS Code via `Extensions > Install from VSIX...`.

## Configuration

Configure the extension in VS Code settings (`settings.json` or UI):

- `eve.cliPath`: Path to the EVE CLI binary (default: `"eve"`). On Windows, append `.exe` if needed.
- `eve.args`: Additional arguments to pass to the CLI (default: `[]`).
- `eve.logLevel`: Log level for output (options: `"error"`, `"warn"`, `"info"`, `"debug"`, `"trace"`; default: `"info"`).
- `eve.serverMode`: Mode to run the CLI in (options: `"stdio"`, `"http"`; default: `"stdio"`).
- `eve.httpPort`: Port for HTTP mode (default: `0` for auto-assignment).

## Usage

- **EVE: Chat** (`Ctrl+Alt+E` on Windows/Linux, `Cmd+Alt+E` on Mac): Open a prompt to chat with EVE using the active editor context.
- **EVE: Plan Edits**: Plan edits for the workspace.
- **EVE: Apply Edits**: Apply planned edits to files.
- **EVE: Open WebUI**: Open a WebView panel for interactive chat, planning, and applying.
- **EVE: Start CLI Daemon**: Start the CLI daemon manually.
- **EVE: Stop CLI Daemon**: Stop the running daemon.

Output is logged to the "EVE" output channel in VS Code.

## Modes

- **Stdio Mode** (default): Runs the CLI as `eve --daemon` and communicates via stdio.
- **HTTP Mode**: Runs the CLI as `eve --serve --port=<port>` and communicates via HTTP POST to `/rpc`.

## Development

- Watch for changes: `npm run watch`
- Compile: `npm run compile`
- Package: `npm run package`

## Changelog

### 0.0.1

- Initial release with basic chat, plan, apply, and WebUI commands.
