# ChatBot API Documentation

**Base URL**: `http://localhost:3000/api`
**Version**: 1.2.0
**Last Updated**: 2025-10-31

---

## Table of Contents

- [Response Format & Error Handling](#response-format--error-handling)
- [Health Check](#health-check)
- [Chat Endpoints](#chat-endpoints)
- [Persona Endpoints](#persona-endpoints)
- [File Analysis Endpoints](#file-analysis-endpoints)
- [Audio Endpoints](#audio-endpoints)
- [WebSocket API](#websocket-api)

---

## Response Format & Error Handling

### Success Response
```json
{
  "data": { ... }
}
```

### Error Response
```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes

| Status | Meaning |
|--------|---------|
| `200` | OK - Success |
| `400` | Bad Request - Invalid input |
| `404` | Not Found - Resource not found |
| `413` | Payload Too Large - File too large |
| `415` | Unsupported Media Type - Invalid file type |
| `422` | Unprocessable Entity - Cannot process request |
| `500` | Internal Server Error - Server error |
| `503` | Service Unavailable - External service error |

---

## Health Check

### Check API Status

**Endpoint**: `GET /health`

**Response**:
```json
{
  "status": "ok",
  "timestamp": "2025-10-31T10:30:00Z"
}
```

---

## Chat Endpoints

### 1. Send Chat Message

**Endpoint**: `POST /api/chat`

**Request Body**:
```json
{
  "message": "Hello, how are you?",
  "session_id": "user_session_123",
  "persona_id": 1,
  "system_prompt": "You are a helpful assistant",
  "use_history": true,
  "temperature": 0.7,
  "max_tokens": 1000,
  "model": "gpt-4o-mini"
}
```

**Parameters**:

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `message` | string | Yes | - | User's message |
| `session_id` | string | No | auto-generated | Session ID for conversation history |
| `persona_id` | integer | No | null | Persona ID to use |
| `system_prompt` | string | No | "" | Custom system prompt (overrides persona) |
| `use_history` | boolean | No | false | Include conversation history in context |
| `temperature` | float | No | 0.7 | Response creativity (0.0-2.0) |
| `max_tokens` | integer | No | - | Maximum response length |
| `model` | string | No | "gpt-4o-mini" | OpenAI model to use |

**Response**:
```json
{
  "message_id": "uuid",
  "session_id": "user_session_123",
  "reply": "I'm doing well, thank you!",
  "persona": {
    "id": 1,
    "name": "General Assistant",
    "expertise": "General knowledge",
    "icon": "🤖",
    "description": "Helpful AI assistant"
  },
  "tokens_used": 45,
  "model": "gpt-4o-mini",
  "timestamp": "2025-10-31T10:30:00Z",
  "history_used": true,
  "history_count": 4
}
```

**Notes**:
- When `use_history: true`, the API maintains conversation context using `session_id`
- `history_count` shows the number of previous messages included in the context
- For better AI compliance with Thai language prompts, use English with clear instructions:

```json
{
  "system_prompt": "You MUST respond in Thai language. You are a professional consultant. Always provide structured answers with: 1) Summary 2) Details 3) Recommendations."
}
```

**Example**:
```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "What is AI?",
    "session_id": "session_001",
    "use_history": true
  }'
```

### 2. Get Chat History

**Endpoint**: `GET /api/chat/history`

**Query Parameters**:
- `limit` (default: 50, max: 100) - Number of messages
- `offset` (default: 0) - Pagination offset

**Response**:
```json
{
  "messages": [
    {
      "id": "uuid",
      "role": "user",
      "content": "Hello",
      "persona_id": 1,
      "created_at": "2025-10-31T10:30:00Z"
    }
  ],
  "total": 150,
  "limit": 50,
  "offset": 0
}
```

---

## Persona Endpoints

### 1. List All Personas

**Endpoint**: `GET /api/personas`

**Response**:
```json
{
  "personas": [
    {
      "id": 1,
      "name": "General Assistant",
      "expertise": "General knowledge",
      "icon": "🤖",
      "description": "Helpful AI assistant for general queries"
    }
  ],
  "total": 3
}
```

### 2. Get Persona Details

**Endpoint**: `GET /api/personas/:id`

**Response**:
```json
{
  "id": 1,
  "name": "General Assistant",
  "expertise": "General knowledge",
  "icon": "🤖",
  "description": "Helpful AI assistant",
  "system_prompt": "You are a helpful assistant...",
  "created_at": "2025-01-01T00:00:00Z"
}
```

---

## File Analysis Endpoints

### 1. Analyze File

**Endpoint**: `POST /api/file/analyze`

**Content-Type**: `multipart/form-data`

**Form Fields**:

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `file` | file | Yes | - | File to analyze |
| `analysis_type` | string | No | "summary" | Analysis type: `summary`, `detail`, `qa`, `extract` |
| `prompt` | string | No | "" | Custom analysis instructions |
| `language` | string | No | "th" | Response language: `th`, `en` |

**Supported File Types**:
- **Documents**: PDF, DOCX, TXT, MD, CSV (max 10MB)
- **Images**: PNG, JPG, JPEG, GIF, WEBP (max 20MB, uses Vision API)

**Analysis Types**:
- `summary` - Brief overview of content
- `detail` - Comprehensive analysis with key points
- `qa` - Question-answering format
- `extract` - Extract structured data and entities

**Response**:
```json
{
  "file_id": "uuid",
  "filename": "document.pdf",
  "file_type": "application/pdf",
  "file_size": 245678,
  "analysis": "This document discusses...",
  "summary": "Brief overview",
  "key_points": ["Point 1", "Point 2"],
  "entities": ["Name", "Company"],
  "sentiment": "neutral",
  "language": "th",
  "tokens_used": 1234,
  "process_time": 2500
}
```

**Example**:
```bash
# Analyze PDF
curl -X POST http://localhost:3000/api/file/analyze \
  -F "file=@report.pdf" \
  -F "analysis_type=summary" \
  -F "language=th"

# Analyze Image
curl -X POST http://localhost:3000/api/file/analyze \
  -F "file=@diagram.png" \
  -F "analysis_type=detail" \
  -F "prompt=Explain this diagram"
```

### 2. Get File Analysis History

**Endpoint**: `GET /api/file/history`

**Query Parameters**:
- `limit` (default: 20, max: 100)
- `offset` (default: 0)
- `file_type` (optional) - Filter by file type or "all"

**Response**:
```json
{
  "files": [
    {
      "file_id": "uuid",
      "filename": "document.pdf",
      "file_type": "application/pdf",
      "file_size": 245678,
      "analysis_type": "summary",
      "language": "th",
      "tokens_used": 1234,
      "created_at": "2025-10-31T10:30:00Z"
    }
  ],
  "total": 50,
  "limit": 20,
  "offset": 0
}
```

---

## Audio Endpoints

### 1. Transcribe Audio (Speech-to-Text)

**Endpoint**: `POST /api/audio/transcribe`

**Content-Type**: `multipart/form-data`

**Form Fields**:
- `audio` (required) - Audio file

**Supported Formats**: MP3, MP4, WAV, WEBM, M4A (max 25MB, max 30 minutes)

**Response**:
```json
{
  "text": "Transcribed text from audio",
  "language": "th",
  "duration": 45.5,
  "timestamp": "2025-10-31T10:30:00Z"
}
```

**Example**:
```bash
curl -X POST http://localhost:3000/api/audio/transcribe \
  -F "audio=@voice.mp3"
```

### 2. Text-to-Speech (TTS)

**Endpoint**: `POST /api/audio/tts`

**Request Body**:
```json
{
  "text": "Hello world",
  "voice": "nova",
  "model": "tts-1",
  "speed": 1.0,
  "response_format": "mp3"
}
```

**Parameters**:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `text` | string | - | Text to convert (max 4096 characters) |
| `voice` | string | "nova" | Voice: `alloy`, `echo`, `fable`, `onyx`, `nova`, `shimmer` |
| `model` | string | "tts-1" | Model: `tts-1` (faster), `tts-1-hd` (higher quality) |
| `speed` | float | 1.0 | Speed: 0.25 to 4.0 |
| `response_format` | string | "mp3" | Format: `mp3`, `opus`, `aac`, `flac`, `wav`, `pcm` |

**Response**: Binary audio file

**Example**:
```bash
curl -X POST http://localhost:3000/api/audio/tts \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello world",
    "voice": "nova"
  }' \
  --output speech.mp3
```

**Voice Options**:

| Voice | Description |
|-------|-------------|
| `alloy` | Neutral, balanced |
| `echo` | Male voice |
| `fable` | British accent |
| `onyx` | Deep male voice |
| `nova` | Female voice (default) |
| `shimmer` | Soft female voice |

---

## WebSocket API

### Stream Chat Responses

**Endpoint**: `WS /api/chat/stream`

**Client Request**:
```javascript
const ws = new WebSocket('ws://localhost:3000/api/chat/stream');

ws.onopen = () => {
  ws.send(JSON.stringify({
    message: "Tell me a story",
    session_id: "test_session",
    use_history: true,
    temperature: 0.8
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);

  if (data.type === 'start') {
    console.log('Started:', data);
  } else if (data.type === 'chunk') {
    process.stdout.write(data.content); // Stream content
  } else if (data.type === 'done') {
    console.log('\nCompleted:', data);
  } else if (data.type === 'error') {
    console.error('Error:', data.error);
  }
};
```

**Message Types**:

1. **Start Message**:
```json
{
  "type": "start",
  "session_id": "test_session",
  "model": "gpt-4o-mini",
  "timestamp": "2025-10-31T10:30:00Z"
}
```

2. **Chunk Message** (streaming content):
```json
{
  "type": "chunk",
  "content": "partial response text"
}
```

3. **Done Message**:
```json
{
  "type": "done",
  "session_id": "test_session",
  "tokens_used": 150,
  "finish_reason": "stop"
}
```

4. **Error Message**:
```json
{
  "type": "error",
  "error": "Error description"
}
```

---

## Common Use Cases

### 1. Simple Chat
```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "What is AI?"}'
```

### 2. Conversational Chat with History
```bash
# First message
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "My name is John",
    "session_id": "session_001",
    "use_history": true
  }'

# Follow-up (AI remembers context)
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "What is my name?",
    "session_id": "session_001",
    "use_history": true
  }'
```

### 3. Custom Persona Chat
```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Explain async/await",
    "persona_id": 2,
    "use_history": true
  }'
```

### 4. Analyze PDF Document
```bash
curl -X POST http://localhost:3000/api/file/analyze \
  -F "file=@contract.pdf" \
  -F "analysis_type=detail" \
  -F "prompt=Extract key terms and obligations"
```

### 5. Voice Message Flow
```bash
# Step 1: Transcribe audio
TRANSCRIPT=$(curl -X POST http://localhost:3000/api/audio/transcribe \
  -F "audio=@question.mp3" | jq -r '.text')

# Step 2: Send to chat
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d "{\"message\": \"$TRANSCRIPT\"}"
```

---

## Best Practices

### Conversation History
- Use the same `session_id` for related conversations
- Set `use_history: true` to maintain context
- Consider limiting history to recent messages for cost efficiency

### System Prompts
- **Use English for system prompts** for better AI compliance
- Be specific and provide examples for complex behaviors
- Example:
```json
{
  "system_prompt": "You MUST respond in Thai language. You are a professional consultant. Always provide structured answers with: 1) Summary 2) Details 3) Recommendations. Example: 'สรุป: ..., รายละเอียด: ..., คำแนะนำ: ...'"
}
```

### File Analysis
- Use `summary` for quick overview
- Use `detail` for comprehensive analysis
- Provide specific `prompt` for targeted extraction
- Choose appropriate `language` for response

### Error Handling
- Always check HTTP status codes
- Parse error messages for user-friendly display
- Implement retry logic for 503 errors

---

## Rate Limiting (Recommended for Production)

| Endpoint | Recommended Limit |
|----------|-------------------|
| `POST /api/chat` | 20 requests/minute |
| `POST /api/audio/transcribe` | 10 requests/minute |
| `POST /api/file/analyze` | 20 requests/minute |
| `GET /api/personas` | 60 requests/minute |
| WebSocket connections | 5 connections/user |

---

## Changelog

### Version 1.2.0 (2025-10-31) - **Current**
- Added conversation history feature (`use_history`, `session_id`)
- Added `history_count` to chat responses
- Refactored controllers with utility helpers
- Improved error handling and response consistency
- Added debug capabilities for system prompt tracing
- Updated documentation to be more concise

### Version 1.1.0 (2025-10-31)
- Added file analysis endpoints (PDF, DOCX, images)
- Implemented Vision API for image analysis
- Added file history tracking

### Version 1.0.0 (2025-10-28)
- Initial API release
- Chat, persona, and audio endpoints
- WebSocket streaming support

---

## Additional Resources

- **OpenAI API Docs**: https://platform.openai.com/docs
- **WebSocket Guide**: https://developer.mozilla.org/en-US/docs/Web/API/WebSocket

**Base URL**: http://localhost:3000/api
**Version**: 1.2.0
**Last Updated**: 2025-10-31