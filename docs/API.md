# AI Context Gap Tracker - API Documentation

## Overview

The AI Context Gap Tracker provides RESTful APIs for managing conversational context, evaluating logic rules, auditing responses, and rewriting prompts. This document describes all available endpoints and their usage.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, the API does not require authentication. In production, consider implementing API keys or JWT tokens.

## Table of Contents

1. [Health Check](#health-check)
2. [Context Tracking](#context-tracking)
3. [Logic Rules](#logic-rules)
4. [Response Auditing](#response-auditing)
5. [Prompt Rewriting](#prompt-rewriting)
6. [Pipeline Processing](#pipeline-processing)
7. [Error Handling](#error-handling)
8. [Rate Limiting](#rate-limiting)

## Health Check

### Get System Health

```http
GET /health
```

**Response:**

```json
{
  "status": "healthy",
  "service": "ai-context-gap-tracker"
}
```

## Context Tracking

### Track Context

Store and analyze conversational context for a specific turn.

```http
POST /context/track
```

**Request Body:**

```json
{
  "session_id": "demo-session-001",
  "turn_number": 1,
  "user_input": "I want to get a visa for a place I mentioned earlier."
}
```

**Response:**

```json
{
  "id": 1,
  "session_id": "demo-session-001",
  "turn_number": 1,
  "user_input": "I want to get a visa for a place I mentioned earlier.",
  "entities": {},
  "topics": [],
  "timeline": [],
  "assertions": [],
  "ambiguities": [],
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### Get Session Context

Retrieve all context turns for a session.

```http
GET /context/session/{sessionId}
```

**Parameters:**
- `sessionId` (path): Session identifier

**Response:**

```json
[
  {
    "id": 1,
    "session_id": "demo-session-001",
    "turn_number": 1,
    "user_input": "I want to get a visa for a place I mentioned earlier.",
    "entities": {},
    "topics": [],
    "timeline": [],
    "assertions": [],
    "ambiguities": [],
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
]
```

### Get Context for Specific Turn

Retrieve context for a specific turn in a session.

```http
GET /context/session/{sessionId}/turn/{turnNumber}
```

**Parameters:**
- `sessionId` (path): Session identifier
- `turnNumber` (path): Turn number

**Response:**

```json
{
  "id": 1,
  "session_id": "demo-session-001",
  "turn_number": 1,
  "user_input": "I want to get a visa for a place I mentioned earlier.",
  "entities": {},
  "topics": [],
  "timeline": [],
  "assertions": [],
  "ambiguities": [],
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### Get Memory Graph

Retrieve the memory graph for a session.

```http
GET /context/session/{sessionId}/memory
```

**Parameters:**
- `sessionId` (path): Session identifier

**Response:**

```json
{
  "session_id": "demo-session-001",
  "nodes": {},
  "edges": [],
  "updated_at": "2024-01-01T12:00:00Z"
}
```

## Logic Rules

### Get Active Rules

Retrieve all active logic rules.

```http
GET /rules
```

**Response:**

```json
[
  {
    "id": 1,
    "name": "Temporal Consistency Check",
    "description": "Checks for temporal consistency in user input",
    "rule_type": "temporal_consistency",
    "conditions": {"enabled": true},
    "actions": {"type": "temporal_check"},
    "priority": 100,
    "is_active": true,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
]
```

### Create Rule

Create a new logic rule.

```http
POST /rules
```

**Request Body:**

```json
{
  "name": "Custom Rule",
  "description": "A custom logic rule",
  "rule_type": "custom",
  "conditions": {"enabled": true},
  "actions": {"type": "custom_action"},
  "priority": 50,
  "is_active": true
}
```

**Response:**

```json
{
  "id": 2,
  "name": "Custom Rule",
  "description": "A custom logic rule",
  "rule_type": "custom",
  "conditions": {"enabled": true},
  "actions": {"type": "custom_action"},
  "priority": 50,
  "is_active": true,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### Update Rule

Update an existing rule.

```http
PUT /rules/{id}
```

**Parameters:**
- `id` (path): Rule identifier

**Request Body:**

```json
{
  "name": "Updated Rule",
  "description": "An updated logic rule",
  "rule_type": "custom",
  "conditions": {"enabled": true},
  "actions": {"type": "updated_action"},
  "priority": 75,
  "is_active": true
}
```

### Delete Rule

Delete a rule.

```http
DELETE /rules/{id}
```

**Parameters:**
- `id` (path): Rule identifier

**Response:**

```json
{
  "message": "rule deleted successfully"
}
```

### Evaluate Rules

Evaluate all active rules against a given context.

```http
POST /rules/evaluate
```

**Request Body:**

```json
{
  "session_id": "demo-session-001",
  "turn_number": 1,
  "user_input": "I want to go there soon, but I'm not sure about the requirements.",
  "entities": {},
  "topics": [],
  "timeline": [],
  "assertions": [],
  "ambiguities": [],
  "history": []
}
```

**Response:**

```json
[
  {
    "rule_id": 1,
    "rule_name": "Temporal Consistency Check",
    "matched": true,
    "confidence": 0.8,
    "actions": [
      {
        "type": "temporal_check",
        "parameters": {
          "keyword": "soon",
          "context": "temporal_reference_detected"
        }
      }
    ],
    "violations": [],
    "suggestions": ["Specify a more precise time than 'soon'"]
  }
]
```

### Initialize Default Rules

Initialize the system with default rules.

```http
POST /rules/initialize
```

**Response:**

```json
{
  "message": "default rules initialized successfully"
}
```

## Response Auditing

### Audit Response

Audit a response for quality, assumptions, and contradictions.

```http
POST /audit/response
```

**Request Body:**

```json
{
  "session_id": "demo-session-001",
  "turn_number": 1,
  "response_text": "Based on your previous mention of Canada, I assume you're asking about Canadian visa requirements.",
  "context": {}
}
```

**Response:**

```json
{
  "id": 1,
  "session_id": "demo-session-001",
  "turn_number": 1,
  "response_text": "Based on your previous mention of Canada, I assume you're asking about Canadian visa requirements.",
  "certainty_level": "assumed",
  "flags": {
    "contains_assumptions": true,
    "contains_contradictions": false,
    "response_length": 89
  },
  "assumptions": [
    {
      "text": "assume you're asking about Canadian visa requirements",
      "confidence": 0.8,
      "source": "keyword_detection",
      "critical": true
    }
  ],
  "contradictions": [],
  "retry_count": 0,
  "recommendations": [
    "Consider explicitly stating assumptions or seeking confirmation"
  ],
  "quality_score": 0.7,
  "created_at": "2024-01-01T12:00:00Z"
}
```

### Get Audit History

Retrieve audit history for a session.

```http
GET /audit/session/{sessionId}/history
```

**Parameters:**
- `sessionId` (path): Session identifier

**Response:**

```json
[
  {
    "id": 1,
    "session_id": "demo-session-001",
    "turn_number": 1,
    "response_text": "Based on your previous mention of Canada, I assume you're asking about Canadian visa requirements.",
    "certainty_level": "assumed",
    "flags": {
      "contains_assumptions": true,
      "contains_contradictions": false,
      "response_length": 89
    },
    "assumptions": [
      {
        "text": "assume you're asking about Canadian visa requirements",
        "confidence": 0.8,
        "source": "keyword_detection",
        "critical": true
      }
    ],
    "contradictions": [],
    "retry_count": 0,
    "created_at": "2024-01-01T12:00:00Z"
  }
]
```

## Prompt Rewriting

### Rewrite Prompt (Advanced)

Rewrite a prompt with full configuration options.

```http
POST /prompt/rewrite
```

**Request Body:**

```json
{
  "session_id": "demo-session-001",
  "turn_number": 1,
  "user_input": "Can you help me with that thing we discussed?",
  "system_prompt": "You are a helpful assistant.",
  "context": {},
  "options": {
    "include_context": true,
    "include_ambiguities": true,
    "include_assumptions": true,
    "include_history": true,
    "add_disambiguation": true,
    "add_clarity_flags": true,
    "max_context_length": 2000,
    "max_history_turns": 5,
    "optimize_for_clarity": true,
    "optimize_for_accuracy": true
  }
}
```

**Response:**

```json
{
  "original_prompt": "Can you help me with that thing we discussed?",
  "rewritten_prompt": "You are a helpful assistant.\n\nCONTEXT INFORMATION:\n...\n\nDISAMBIGUATION REQUIRED:\n- Clarify what 'that thing' refers to\n\nUSER INPUT:\nCan you help me with that thing we discussed?\n\nRESPONSE REQUIREMENTS:\n- Provide clear, unambiguous responses\n- State assumptions explicitly\n- Ask for clarification when needed\n",
  "context": {},
  "ambiguities": ["Clarify what 'that thing' refers to"],
  "assumptions": [],
  "clarity_flags": ["Clarification needed: Missing Information Detection"],
  "disambiguation_flags": ["Disambiguation needed: Ambiguity Resolution"],
  "quality_score": 0.6,
  "recommendations": ["Consider resolving ambiguities before processing"],
  "processing_time": 0.012
}
```

### Simple Prompt Rewrite

Rewrite a prompt with default settings.

```http
POST /prompt/simple-rewrite
```

**Request Body:**

```json
{
  "session_id": "demo-session-001",
  "turn_number": 1,
  "user_input": "Can you help me with that thing we discussed?"
}
```

**Response:**

```json
{
  "rewritten_prompt": "CLARITY CONSIDERATIONS:\n- Clarification needed: Missing Information Detection\n\nIDENTIFIED AMBIGUITIES:\n- Clarify what 'that' refers to\n\nUSER INPUT:\nCan you help me with that thing we discussed?\n\nRESPONSE REQUIREMENTS:\n- Provide clear, unambiguous responses\n- State assumptions explicitly\n- Ask for clarification when needed\n"
}
```

## Pipeline Processing

### Process Complete Pipeline

Process input through the complete AI Context Gap Tracker pipeline.

```http
POST /pipeline/process
```

**Request Body:**

```json
{
  "session_id": "demo-session-001",
  "turn_number": 1,
  "user_input": "I need information about the process for that European country.",
  "system_prompt": "You are a helpful assistant for travel planning.",
  "options": {}
}
```

**Response:**

```json
{
  "session_id": "demo-session-001",
  "turn_number": 1,
  "context": {
    "id": 1,
    "session_id": "demo-session-001",
    "turn_number": 1,
    "user_input": "I need information about the process for that European country.",
    "entities": {},
    "topics": [],
    "timeline": [],
    "assertions": [],
    "ambiguities": [],
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  },
  "rule_results": [
    {
      "rule_id": 1,
      "rule_name": "Missing Information Detection",
      "matched": true,
      "confidence": 0.9,
      "actions": [],
      "violations": [
        {
          "type": "vague_reference",
          "description": "Vague reference detected: 'that'",
          "severity": "medium",
          "confidence": 0.8
        }
      ],
      "suggestions": ["Clarify what 'that' refers to"]
    }
  ],
  "prompt_result": {
    "original_prompt": "I need information about the process for that European country.",
    "rewritten_prompt": "You are a helpful assistant for travel planning.\n\nCLARITY CONSIDERATIONS:\n- Clarification needed: Missing Information Detection\n\nIDENTIFIED AMBIGUITIES:\n- Clarify what 'that' refers to\n\nUSER INPUT:\nI need information about the process for that European country.\n\nRESPONSE REQUIREMENTS:\n- Provide clear, unambiguous responses\n- State assumptions explicitly\n- Ask for clarification when needed\n",
    "context": {},
    "ambiguities": ["Clarify what 'that' refers to"],
    "assumptions": [],
    "clarity_flags": ["Clarification needed: Missing Information Detection"],
    "disambiguation_flags": [],
    "quality_score": 0.7,
    "recommendations": [],
    "processing_time": 0.015
  },
  "pipeline_stage": "completed"
}
```

## NLP Service API

The NLP service provides natural language processing capabilities and is available at `http://localhost:5000`.

### Entity Extraction

```http
POST /entities
```

**Request Body:**

```json
{
  "text": "I want to visit Paris, France next month for vacation.",
  "session_id": "demo-session-001",
  "turn_number": 1
}
```

### Sentiment Analysis

```http
POST /sentiment
```

**Request Body:**

```json
{
  "text": "I am very excited about this trip!",
  "session_id": "demo-session-001",
  "turn_number": 1
}
```

### Complete NLP Analysis

```http
POST /analyze
```

**Request Body:**

```json
{
  "text": "I want to visit Paris, France next month for vacation.",
  "session_id": "demo-session-001",
  "turn_number": 1
}
```

## Error Handling

### Error Response Format

All API errors follow this format:

```json
{
  "error": "Error message describing what went wrong",
  "code": "ERROR_CODE",
  "details": {
    "field": "Additional error details"
  }
}
```

### Common HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request format or parameters
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Service temporarily unavailable

## Rate Limiting

Currently, no rate limiting is implemented. In production, consider implementing:

- Rate limiting per IP address
- Rate limiting per API key
- Different limits for different endpoints

## SDK and Client Libraries

### cURL Examples

```bash
# Track context
curl -X POST http://localhost:8080/api/v1/context/track \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "demo-session-001",
    "turn_number": 1,
    "user_input": "I want to get a visa for a place I mentioned earlier."
  }'

# Evaluate rules
curl -X POST http://localhost:8080/api/v1/rules/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "demo-session-001",
    "turn_number": 1,
    "user_input": "I want to go there soon.",
    "entities": {},
    "topics": [],
    "timeline": [],
    "assertions": [],
    "ambiguities": [],
    "history": []
  }'
```

### JavaScript Example

```javascript
const response = await fetch('http://localhost:8080/api/v1/context/track', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    session_id: 'demo-session-001',
    turn_number: 1,
    user_input: 'I want to get a visa for a place I mentioned earlier.'
  })
});

const data = await response.json();
console.log(data);
```

### Python Example

```python
import requests

response = requests.post('http://localhost:8080/api/v1/context/track', 
  json={
    'session_id': 'demo-session-001',
    'turn_number': 1,
    'user_input': 'I want to get a visa for a place I mentioned earlier.'
  }
)

data = response.json()
print(data)
```

## Changelog

### v1.0.0 (Current)
- Initial release
- Context tracking functionality
- Logic rule engine
- Response auditing
- Prompt rewriting
- Pipeline processing
- NLP integration

For more detailed information, see the [README](README.md) and [DEPLOYMENT](DEPLOYMENT.md) guides.