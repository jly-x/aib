# AIB - AI-powered Command Line Assistant

AIB (AI Buddy) is a lightweight CLI tool that I built after getting tired of alt-tabbing between my terminal and browser. It provides AI assistance right in your terminal, especially useful for bash command generation.

## Features

- **Quick answers** - Get concise responses to questions right in your terminal
- **Command generation** - Generate executable shell commands with the `-s` flag
- **Clipboard integration** - Commands are automatically copied to your clipboard for easy use

## Installation

### Prerequisites

- Go 1.19 or higher
- A Google API key for Gemini

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/aib.git
   cd aib
   ```

2. Create a `.env` file in the project directory with your Gemini API Key from Google AI Studio:
   ```
   GOOGLE_API_KEY=your_api_key_here
   ```

3. Build and install:
   ```bash
   go install
   ```

   Make sure your `$GOPATH/bin` is in your PATH.

## Usage

### General Questions

```bash
aib how do I check if a file exists in Go
```

### Generate Shell Commands

```bash
aib -s find all files modified in the last 7 days
```
The command will be displayed and automatically copied to your clipboard.

## Future Features if I need them

- **Choose your model**: Select between different AI models.
- **Disable clipboard**: Disable automatic copying to clipboard.
- **Customize prompt**: Customize the prompt as a config option.

## License

MIT
