"""
AI Context Gap Tracker - MCP Server Package

Universal MCP server that exposes AI Context Gap Tracker functionality as tools
for Claude Desktop and other MCP-compliant applications.
"""

__version__ = "1.0.0"
__author__ = "AI Context Gap Tracker Team"
__email__ = "team@aicontextgap.com"

from .main import MCPServer

__all__ = ["MCPServer"]