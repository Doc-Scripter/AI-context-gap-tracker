# AI Context Gap Tracker - MCP Server Setup Guide

This guide explains how to set up and use the AI Context Gap Tracker as an MCP (Model Context Protocol) server for Claude Desktop and other MCP-compliant applications.

## Overview

The AI Context Gap Tracker MCP Server exposes the following tools:
- **rewrite_prompt**: Enhance prompts with context and clarity flags
- **audit_response**: Audit AI responses for quality and assumptions
- **track_context**: Track conversation context and detect information gaps
- **evaluate_rules**: Evaluate logical rules against user input
- **get_session_context**: Retrieve stored context for a conversation session

## Installation Options

### Option 1: Using uvx (Recommended)

The easiest way to run the MCP server is using `uvx`, which handles Python dependencies automatically:

```bash
# Install and run directly from the project directory
uvx --from . ai-context-gap-tracker

# Or run the alternative command
uvx --from . context-tracker-mcp
```

### Option 2: Using pip

```bash
# Install from the project directory
pip install .

# Run the MCP server
ai-context-gap-tracker
```

### Option 3: Direct Python execution

```bash
# Navigate to the mcp_server directory
cd mcp_server

# Install dependencies
pip install -r requirements.txt

# Run directly
python main.py
```

## Prerequisites

Before running the MCP server, ensure the main AI Context Gap Tracker backend services are running:

### 1. Start Backend Services

```bash
# From the project root directory
docker-compose up -d

# Or just the core services (if you don't need the containerized MCP server)
docker-compose up -d postgres redis ai-context-tracker nlp-service
```

### 2. Verify Services

Check that the main API is accessible:

```bash
curl http://localhost:8080/api/v1/health
```

You should see a response indicating the service is healthy.

## Claude Desktop Configuration

To use the AI Context Gap Tracker with Claude Desktop, add the following configuration to your `claude_desktop_config.json` file:

### Configuration File Locations

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

### Configuration Example

```json
{
  "mcpServers": {
    "ai-context-gap-tracker": {
      "command": "uvx",
      "args": [
        "--from",
        "/path/to/ai-context-gap-tracker",
        "ai-context-gap-tracker"
      ],
      "env": {
        "TRACKER_API_ENDPOINT": "http://localhost:8080"
      }
    }
  }
}
```

**Note**: Replace `/path/to/ai-context-gap-tracker` with the actual path to your project directory.

### Alternative Configuration (using Python directly)

If you prefer not to use uvx:

```json
{
  "mcpServers": {
    "ai-context-gap-tracker": {
      "command": "python",
      "args": [
        "/path/to/ai-context-gap-tracker/mcp_server/main.py"
      ],
      "env": {
        "TRACKER_API_ENDPOINT": "http://localhost:8080",
        "PYTHONPATH": "/path/to/ai-context-gap-tracker"
      }
    }
  }
}
```

## Environment Variables

The MCP server supports the following environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `TRACKER_API_ENDPOINT` | `http://localhost:8080` | URL of the main AI Context Gap Tracker API |

## Usage Examples

Once configured, you can use the following tools in Claude Desktop:

### 1. Rewrite Prompt
Enhance a prompt with context and clarity flags:

```
Use the rewrite_prompt tool with this prompt: "Help me with that thing we discussed earlier"
```

### 2. Audit Response
Audit an AI response for quality issues:

```
Use the audit_response tool to analyze this response: "Based on my assumptions about your previous request, here's what I think..."
```

### 3. Track Context
Track conversation context:

```
Use the track_context tool with this input: "I want to visit that European country we mentioned"
```

### 4. Evaluate Rules
Check for logical inconsistencies:

```
Use the evaluate_rules tool with: "I always never go there, but sometimes I do"
```

### 5. Get Session Context
Retrieve stored conversation context:

```
Use the get_session_context tool to see what context has been tracked
```

## Troubleshooting

### Common Issues

1. **Connection Refused**: Ensure the backend services are running (`docker-compose up -d`)
2. **Module Not Found**: Make sure dependencies are installed (`pip install -r requirements.txt`)
3. **Path Issues**: Verify the paths in your Claude Desktop configuration are correct

### Debug Mode

To run with debug logging:

```bash
# Set log level to DEBUG
export LOG_LEVEL=DEBUG
uvx --from . ai-context-gap-tracker
```

### Offline Mode

The MCP server gracefully handles offline scenarios. If the backend services are not available, it will:
- Return the original prompt for rewrite requests
- Indicate that services are unavailable
- Provide basic fallback responses

## Development Setup

For development or testing:

```bash
# Clone and navigate to the project
git clone <repository-url>
cd ai-context-gap-tracker

# Start backend services
docker-compose up -d

# Install in development mode
pip install -e .

# Run the MCP server
ai-context-gap-tracker
```

## Universal MCP Server

This MCP server is designed to be universal and can work with any MCP-compliant host application, not just Claude Desktop. The server follows the MCP specification and can be integrated with:

- Claude Desktop
- VS Code extensions
- Custom IDE integrations
- Other AI-powered applications supporting MCP

## Support

If you encounter issues:

1. Check that backend services are healthy: `curl http://localhost:8080/api/v1/health`
2. Verify your configuration paths are correct
3. Check the MCP server logs for error messages
4. Ensure all dependencies are properly installed

For more information, see the main [README.md](../README.md) and [DEPLOYMENT.md](./DEPLOYMENT.md) files.