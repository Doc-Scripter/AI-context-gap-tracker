#!/usr/bin/env python3
"""
AI Context Gap Tracker - MCP Server
Universal MCP server that exposes AI Context Gap Tracker functionality as tools.
"""

import asyncio
import json
import logging
import os
import sys
from typing import Any, Dict, List, Optional
import httpx

from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import (
    Resource,
    Tool,
    TextContent,
    CallToolResult,
    ListResourcesResult,
    ListToolsResult,
    ReadResourceResult,
)
from pydantic import BaseModel

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class MCPServer:
    def __init__(self):
        self.server = Server("ai-context-gap-tracker")
        self.tracker_api_endpoint = os.getenv(
            "TRACKER_API_ENDPOINT", 
            "http://localhost:8080"
        )
        self.http_client = httpx.AsyncClient(timeout=30.0)
        
        # Register tools and resources
        self._register_tools()
        self._register_resources()
        
    def _register_tools(self):
        """Register MCP tools that expose tracker functionality"""
        
        @self.server.list_tools()
        async def list_tools() -> List[Tool]:
            return [
                Tool(
                    name="rewrite_prompt",
                    description="Enhance a prompt with context and clarity flags to improve AI accuracy",
                    inputSchema={
                        "type": "object",
                        "properties": {
                            "prompt": {
                                "type": "string",
                                "description": "The original prompt to enhance"
                            },
                            "session_id": {
                                "type": "string",
                                "description": "Session ID for context tracking (optional)"
                            },
                            "context": {
                                "type": "object",
                                "description": "Additional context information (optional)"
                            }
                        },
                        "required": ["prompt"]
                    }
                ),
                Tool(
                    name="audit_response",
                    description="Audit an AI response for quality, assumptions, and potential issues",
                    inputSchema={
                        "type": "object",
                        "properties": {
                            "response": {
                                "type": "string",
                                "description": "The AI response to audit"
                            },
                            "original_prompt": {
                                "type": "string",
                                "description": "The original prompt that generated the response (optional)"
                            },
                            "context": {
                                "type": "object",
                                "description": "Context information for auditing (optional)"
                            }
                        },
                        "required": ["response"]
                    }
                ),
                Tool(
                    name="track_context",
                    description="Track conversation context and detect information gaps",
                    inputSchema={
                        "type": "object",
                        "properties": {
                            "user_input": {
                                "type": "string",
                                "description": "User input to analyze for context"
                            },
                            "session_id": {
                                "type": "string",
                                "description": "Session ID for context tracking"
                            },
                            "turn_number": {
                                "type": "integer",
                                "description": "Turn number in the conversation"
                            }
                        },
                        "required": ["user_input"]
                    }
                ),
                Tool(
                    name="evaluate_rules",
                    description="Evaluate logical rules against user input to detect inconsistencies",
                    inputSchema={
                        "type": "object",
                        "properties": {
                            "user_input": {
                                "type": "string",
                                "description": "User input to evaluate against rules"
                            },
                            "session_id": {
                                "type": "string",
                                "description": "Session ID for context"
                            },
                            "entities": {
                                "type": "object",
                                "description": "Extracted entities (optional)"
                            }
                        },
                        "required": ["user_input"]
                    }
                ),
                Tool(
                    name="get_session_context",
                    description="Retrieve stored context for a conversation session",
                    inputSchema={
                        "type": "object",
                        "properties": {
                            "session_id": {
                                "type": "string",
                                "description": "Session ID to retrieve context for"
                            }
                        },
                        "required": ["session_id"]
                    }
                )
            ]
        @self.server.call_tool()
        async def call_tool(name: str, arguments: Dict[str, Any]) -> CallToolResult:
            try:
                if name == "rewrite_prompt":
                    return await self._rewrite_prompt(arguments)
                elif name == "audit_response":
                    return await self._audit_response(arguments)
                elif name == "track_context":
                    return await self._track_context(arguments)
                elif name == "evaluate_rules":
                    return await self._evaluate_rules(arguments)
                elif name == "get_session_context":
                    return await self._get_session_context(arguments)
                else:
                    raise ValueError(f"Unknown tool: {name}")
                    
            except Exception as e:
                logger.error(f"Error calling tool {name}: {e}")
                return CallToolResult(
                    content=[
                        TextContent(
                            type="text",
                            text=f"Error executing {name}: {str(e)}"
                        )
                    ]
                )
    
    def _register_resources(self):
        """Register MCP resources"""
        
        @self.server.list_resources()
        async def list_resources() -> List[Resource]:
            return [
                Resource(
                    uri="context://session/current",
                    name="Current Session Context",
                    description="Current conversation context and tracked information",
                    mimeType="application/json"
                ),
                Resource(
                    uri="rules://active",
                    name="Active Rules",
                    description="Currently active logical rules for evaluation",
                    mimeType="application/json"
                ),
                Resource(
                    uri="metrics://performance",
                    name="Performance Metrics",
                    description="System performance and accuracy metrics",
                    mimeType="application/json"
                )
            ]
        
        @self.server.read_resource()
        async def read_resource(uri: str) -> ReadResourceResult:
            try:
                if uri == "context://session/current":
                    return await self._read_current_context()
                elif uri == "rules://active":
                    return await self._read_active_rules()
                elif uri == "metrics://performance":
                    return await self._read_performance_metrics()
                else:
                    raise ValueError(f"Unknown resource: {uri}")
            except Exception as e:
                logger.error(f"Error reading resource {uri}: {e}")
                return ReadResourceResult(
                    contents=[
                        TextContent(
                            type="text",
                            text=f"Error reading resource: {str(e)}"
                        )
                    ]
                )
    
    async def _rewrite_prompt(self, arguments: Dict[str, Any]) -> CallToolResult:
        """Rewrite a prompt for better clarity and context"""
        prompt = arguments["prompt"]
        session_id = arguments.get("session_id", "mcp-session")
        context = arguments.get("context", {})
        
        payload = {
            "prompt": prompt,
            "session_id": session_id,
            "context": context,
            "rewrite_type": "comprehensive"
        }
        
        try:
            response = await self.http_client.post(
                f"{self.tracker_api_endpoint}/api/v1/prompt/simple-rewrite",
                json=payload
            )
            response.raise_for_status()
            result = response.json()
            
            return CallToolResult(
                content=[
                    TextContent(
                        type="text",
                        text=f"**Enhanced Prompt:**\n{result.get('rewritten_prompt', prompt)}\n\n**Improvements:**\n{self._format_improvements(result.get('improvements', []))}\n\n**Clarity Flags:**\n{self._format_flags(result.get('clarity_flags', []))}"
                    )
                ]
            )
        
        except httpx.RequestError as e:
            logger.error(f"HTTP request failed: {e}")
            return CallToolResult(
                content=[
                    TextContent(
                        type="text",
                        text=f"Failed to connect to tracker service: {str(e)}\n\nOriginal prompt: {prompt}"
                    )
                ]
            )
    
    async def _audit_response(self, arguments: Dict[str, Any]) -> CallToolResult:
        """Audit an AI response for quality and assumptions"""
        response_text = arguments["response"]
        original_prompt = arguments.get("original_prompt", "")
        context = arguments.get("context", {})
        
        payload = {
            "response": response_text,
            "original_prompt": original_prompt,
            "context": context
        }
        
        try:
            response = await self.http_client.post(
                f"{self.tracker_api_endpoint}/api/v1/audit/response",
                json=payload
            )
            response.raise_for_status()
            result = response.json()
            
            audit_summary = self._format_audit_result(result)
            
            return CallToolResult(
                content=[
                    TextContent(
                        type="text",
                        text=audit_summary
                    )
                ]
            )
        
        except httpx.RequestError as e:
            logger.error(f"HTTP request failed: {e}")
            return CallToolResult(
                content=[
                    TextContent(
                        type="text",
                        text=f"Failed to connect to tracker service: {str(e)}\n\nResponse appears normal (offline mode)."
                    )
                ]
            )
    
    async def _track_context(self, arguments: Dict[str, Any]) -> CallToolResult:
        """Track conversation context"""
        user_input = arguments["user_input"]
        session_id = arguments.get("session_id", "mcp-session")
        turn_number = arguments.get("turn_number", 1)
        
        payload = {
            "session_id": session_id,
            "turn_number": turn_number,
            "user_input": user_input
        }
        
        try:
            response = await self.http_client.post(
                f"{self.tracker_api_endpoint}/api/v1/context/track",
                json=payload
            )
            response.raise_for_status()
            result = response.json()
            
            return CallToolResult(
                content=[
                    TextContent(
                        type="text",
                        text=f"**Context Tracked:**\n- Session: {session_id}\n- Turn: {turn_number}\n- Entities: {len(result.get('entities', {}))}\n- Topics: {', '.join(result.get('topics', []))}\n- Information Gaps: {len(result.get('gaps', []))}"
                    )
                ]
            )
        
        except httpx.RequestError as e:
            logger.error(f"HTTP request failed: {e}")
            return CallToolResult(
                content=[
                    TextContent(
                        type="text",
                        text=f"Failed to connect to tracker service: {str(e)}\n\nContext tracking unavailable (offline mode)."
                    )
                ]
            )
    
    async def _evaluate_rules(self, arguments: Dict[str, Any]) -> CallToolResult:
        """Evaluate logical rules against user input"""
        user_input = arguments["user_input"]
        session_id = arguments.get("session_id", "mcp-session")
        entities = arguments.get("entities", {})
        
        payload = {
            "session_id": session_id,
            "user_input": user_input,
            "entities": entities
        }
        
        try:
            response = await self.http_client.post(
                f"{self.tracker_api_endpoint}/api/v1/rules/evaluate",
                json=payload
            )
            response.raise_for_status()
            result = response.json()
            
            rules_summary = self._format_rules_result(result.get('results', []))
            
            return CallToolResult(
                content=[
                    TextContent(
                        type="text",
                        text=rules_summary
                    )
                ]
            )
        
        except httpx.RequestError as e:
            logger.error(f"HTTP request failed: {e}")
            return CallToolResult(
                content=[
                    TextContent(
                        type="text",
                        text=f"Failed to connect to tracker service: {str(e)}\n\nRule evaluation unavailable (offline mode)."
                    )
                ]
            )
    
    async def _get_session_context(self, arguments: Dict[str, Any]) -> CallToolResult:
        """Get stored session context"""
        session_id = arguments.get("session_id", "mcp-session")
        
        try:
            response = await self.http_client.get(
                f"{self.tracker_api_endpoint}/api/v1/context/session/{session_id}"
            )
            response.raise_for_status()
            result = response.json()
            
            context_summary = self._format_context_result(result)
            
            return CallToolResult(
                content=[
                    TextContent(
                        type="text",
                        text=context_summary
                    )
                ]
            )
        
        except httpx.RequestError as e:
            logger.error(f"HTTP request failed: {e}")
            return CallToolResult(
                content=[
                    TextContent(
                        type="text",
                        text=f"Failed to connect to tracker service: {str(e)}\n\nSession context unavailable (offline mode)."
                    )
                ]
            )
    
    async def _read_current_context(self) -> ReadResourceResult:
        """Read current session context resource"""
        try:
            response = await self.http_client.get(
                f"{self.tracker_api_endpoint}/api/v1/context/session/mcp-session"
            )
            response.raise_for_status()
            result = response.json()
            
            return ReadResourceResult(
                contents=[
                    TextContent(
                        type="text",
                        text=json.dumps(result, indent=2)
                    )
                ]
            )
        
        except Exception as e:
            return ReadResourceResult(
                contents=[
                    TextContent(
                        type="text",
                        text=f"Context unavailable: {str(e)}"
                    )
                ]
            )
    
    async def _read_active_rules(self) -> ReadResourceResult:
        """Read active rules resource"""
        return ReadResourceResult(
            contents=[
                TextContent(
                    type="text",
                    text=json.dumps({
                        "active_rules": [
                            "temporal_consistency",
                            "missing_information",
                            "contradiction_detection",
                            "ambiguity_resolution"
                        ],
                        "note": "Rules are evaluated dynamically based on input"
                    }, indent=2)
                )
            ]
        )
    
    async def _read_performance_metrics(self) -> ReadResourceResult:
        """Read performance metrics resource"""
        try:
            response = await self.http_client.get(
                f"{self.tracker_api_endpoint}/api/v1/health"
            )
            response.raise_for_status()
            result = response.json()
            
            return ReadResourceResult(
                contents=[
                    TextContent(
                        type="text",
                        text=json.dumps({
                            "service_status": result.get("status", "unknown"),
                            "uptime": result.get("uptime", "unknown"),
                            "last_check": result.get("timestamp", "unknown")
                        }, indent=2)
                    )
                ]
            )
        
        except Exception as e:
            return ReadResourceResult(
                contents=[
                    TextContent(
                        type="text",
                        text=json.dumps({
                            "service_status": "offline",
                            "error": str(e)
                        }, indent=2)
                    )
                ]
            )
    
    def _format_improvements(self, improvements: List[Dict]) -> str:
        """Format prompt improvements for display"""
        if not improvements:
            return "No specific improvements identified"
        
        formatted = []
        for improvement in improvements:
            formatted.append(f"• {improvement.get('description', 'Improvement applied')}")
        
        return "\n".join(formatted)
    
    def _format_flags(self, flags: List[Dict]) -> str:
        """Format clarity flags for display"""
        if not flags:
            return "No clarity issues detected"
        
        formatted = []
        for flag in flags:
            severity = flag.get('severity', 'medium')
            description = flag.get('description', 'Clarity issue detected')
            formatted.append(f"• [{severity.upper()}] {description}")
        
        return "\n".join(formatted)
    
    def _format_audit_result(self, result: Dict) -> str:
        """Format audit result for display"""
        quality_score = result.get('quality_score', 0)
        assumptions = result.get('assumptions', [])
        issues = result.get('issues', [])
        
        summary = f"**Quality Score:** {quality_score}/10\n\n"
        
        if assumptions:
            summary += "**Assumptions Detected:**\n"
            for assumption in assumptions:
                summary += f"• {assumption.get('description', 'Assumption made')}\n"
            summary += "\n"
        
        if issues:
            summary += "**Issues Identified:**\n"
            for issue in issues:
                summary += f"• [{issue.get('severity', 'medium').upper()}] {issue.get('description', 'Issue detected')}\n"
            summary += "\n"
        
        if not assumptions and not issues:
            summary += "**Status:** Response appears accurate with no major issues detected.\n"
        
        return summary
    
    def _format_rules_result(self, results: List[Dict]) -> str:
        """Format rule evaluation results for display"""
        if not results:
            return "No rule violations detected."
        
        summary = "**Rule Evaluation Results:**\n\n"
        
        for result in results:
            rule_name = result.get('rule_name', 'Unknown Rule')
            matched = result.get('matched', False)
            
            if matched:
                summary += f"**{rule_name}:** TRIGGERED\n"
                
                violations = result.get('violations', [])
                for violation in violations:
                    severity = violation.get('severity', 'medium')
                    description = violation.get('description', 'Violation detected')
                    summary += f"  • [{severity.upper()}] {description}\n"
                
                suggestions = result.get('suggestions', [])
                if suggestions:
                    summary += "  Suggestions:\n"
                    for suggestion in suggestions:
                        summary += f"    - {suggestion}\n"
                
                summary += "\n"
        
        return summary
    
    def _format_context_result(self, result: Dict) -> str:
        """Format context result for display"""
        session_id = result.get('session_id', 'unknown')
        turn_count = result.get('turn_count', 0)
        entities = result.get('entities', {})
        topics = result.get('topics', [])
        
        summary = f"**Session:** {session_id}\n"
        summary += f"**Turns:** {turn_count}\n"
        summary += f"**Entities:** {len(entities)} tracked\n"
        summary += f"**Topics:** {', '.join(topics) if topics else 'None identified'}\n"
        
        return summary

async def main():
    """Main entry point for the MCP server"""
    logger.info("Starting AI Context Gap Tracker MCP Server")
    
    # Create server instance
    mcp_server = MCPServer()
    
    # Run the server using stdio transport
    async with stdio_server() as (read_stream, write_stream):
        await mcp_server.server.run(
            read_stream,
            write_stream,
            mcp_server.server.create_initialization_options()
        )

if __name__ == "__main__":
    asyncio.run(main())