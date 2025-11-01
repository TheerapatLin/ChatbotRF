# ChatBot API Documentation

**Base URL**: `http://localhost:3001/api`
**Version**: 2.0.0
**Last Updated**: 2025-11-01
**Database**: PostgreSQL with UUID-based IDs
**AI Models**: OpenAI GPT-4o-mini, Whisper-1, TTS-1, Vision

---

## 📑 Table of Contents

1. [Quick Start](#quick-start)
2. [Authentication & Headers](#authentication--headers)
3. [Response Format](#response-format)
4. [API Endpoints](#api-endpoints)
   - [Health Check](#health-check)
   - [Personas](#personas)
   - [Chat](#chat)
   - [File Analysis](#file-analysis)
   - [Audio](#audio)
   - [WebSocket Streaming](#websocket-streaming)
5. [Database Models](#database-models)
6. [Error Handling](#error-handling)
7. [Best Practices](#best-practices)

---

## 🚀 Quick Start

```bash
# 1. Check API health
curl http://localhost:3001/api/health

# 2. List personas
curl http://localhost:3001/api/personas

# 3. Send a chat message
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello!",
    "session_id": "session_001",
    "use_history": true
  }'

# 4. Analyze a file
curl -X POST http://localhost:3001/api/file/analyze \
  -F "file=@document.pdf" \
  -F "analysis_type=summary" \
  -F "session_id=session_001"
```

---

## 🔐 Authentication & Headers

**Current Version**: No authentication required
**CORS Enabled**: `http://localhost:5173` (Frontend)

**Recommended Headers**:
```
Content-Type: application/json
Accept: application/json
```

---

## 📦 Response Format

### Success Response
```json
{
  "field1": "value1",
  "field2": "value2"
}
```

### Error Response
```json
{
  "error": "Descriptive error message"
}
```

### HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200  | Success |
| 400  | Bad Request - Invalid input |
| 404  | Not Found - Resource doesn't exist |
| 413  | Payload Too Large - File exceeds limit |
| 415  | Unsupported Media Type - Invalid file type |
| 422  | Unprocessable Entity - Cannot parse file |
| 500  | Internal Server Error |
| 503  | Service Unavailable - OpenAI API issue |

---

## 📡 API Endpoints

### Health Check

#### `GET /health`
Check if API is running.

**Response**:
```json
{
  "status": "ok",
  "message": "ChatBot API is running",
  "env": "development"
}
```

---

### Personas

AI personalities with predefined system prompts.

#### `GET /personas`
List all active personas.

**Response**:
```json
{
  "personas": [
    {
      "id": 1,
      "name": "General Assistant",
      "system_prompt": "You are a helpful, friendly AI assistant...",
      "expertise": "general",
      "description": "A versatile AI assistant for general questions",
      "icon": "🤖",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": 2,
      "name": "Technology Expert",
      "expertise": "technology",
      "icon": "💻"
    },
    {
      "id": 3,
      "name": "Business Advisor",
      "expertise": "business",
      "icon": "💼"
    }
  ]
}
```

#### `GET /personas/:id`
Get persona details with usage statistics.

**Response**:
```json
{
  "id": 1,
  "name": "General Assistant",
  "system_prompt": "You are a helpful...",
  "expertise": "general",
  "description": "...",
  "icon": "🤖",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "stats": {
    "total_messages": 42,
    "avg_response_time": "2.3s"
  }
}
```

---

### Chat

Send messages and receive AI responses.

#### `POST /chat`
Send chat message with optional conversation history and file attachments.

**Request**:
```json
{
  "message": "What is machine learning?",
  "session_id": "session_123",
  "persona_id": 2,
  "system_prompt": "You are a teacher",
  "use_history": true,
  "temperature": 0.7,
  "max_tokens": 2000,
  "model": "gpt-4o-mini",
  "file_ids": ["uuid1", "uuid2"]
}
```

**Parameters**:

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `message` | string | ✅ Yes | - | User's message content |
| `session_id` | string | No | auto-generated | Session ID for conversation tracking |
| `persona_id` | int | No | null | Persona to use (1, 2, or 3) |
| `system_prompt` | string | No | "" | Custom prompt (overrides persona) |
| `use_history` | bool | No | false | Include conversation history |
| `temperature` | float | No | 0.7 | Creativity (0.0-2.0) |
| `max_tokens` | int | No | 2000 | Max response length |
| `model` | string | No | gpt-4o-mini | OpenAI model |
| `file_ids` | array | No | [] | File UUIDs to include (max 5) |

**Response**:
```json
{
  "message_id": "uuid-here",
  "session_id": "session_123",
  "reply": "Machine learning is a subset of AI...",
  "persona": {
    "id": 2,
    "name": "Technology Expert",
    "expertise": "technology",
    "icon": "💻",
    "description": "..."
  },
  "tokens_used": 145,
  "model": "gpt-4o-mini",
  "timestamp": "2024-11-01T12:30:00Z",
  "history_used": true,
  "history_count": 8
}
```

**Notes**:
- **Max 5 files** per message
- File context automatically included from `file_ids`
- History limit: 10 recent messages
- Backend will fetch file analysis content and include in context

#### `GET /chat/history`
Retrieve chat history with pagination.

**Query Parameters**:
- `limit` (default: 50, max: 100)
- `offset` (default: 0)

**Response**:
```json
{
  "messages": [
    {
      "id": "uuid",
      "session_id": "session_123",
      "role": "user",
      "content": "Hello",
      "persona_id": 1,
      "tokens_used": null,
      "file_attachments": [
        {
          "file_id": "uuid",
          "filename": "doc.pdf",
          "file_type": "application/pdf",
          "file_size": 102400,
          "analysis_summary": "Summary of the file..."
        }
      ],
      "created_at": "2024-11-01T12:00:00Z"
    },
    {
      "id": "uuid",
      "role": "assistant",
      "content": "Hi there!",
      "persona_id": 1,
      "tokens_used": 15,
      "file_attachments": null,
      "created_at": "2024-11-01T12:01:00Z"
    }
  ],
  "total": 247,
  "limit": 50,
  "offset": 0
}
```

---

### File Analysis

Upload and analyze documents using AI.

#### `POST /file/analyze`
Analyze uploaded file with AI.

**Content-Type**: `multipart/form-data`

**Form Fields**:

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `file` | file | ✅ Yes | - | File to analyze |
| `analysis_type` | string | No | summary | `summary`, `detail`, `qa`, `extract` |
| `session_id` | string | No | null | Link to conversation session |
| `prompt` | string | No | "" | Custom analysis instructions |
| `language` | string | No | th | `th` or `en` |

**Supported File Types**:

| Type | Extensions | Max Size |
|------|-----------|----------|
| Text | .txt, .md | 10 MB |
| Documents | .pdf | 25 MB |
| Office | .docx, .xlsx, .pptx | 25 MB |
| Images | .jpg, .png, .gif, .webp | 20 MB |
| Code | .js, .py, .go, .java | 5 MB |
| Data | .json, .xml, .csv | 10 MB |

**Analysis Types**:
- `summary` - Brief overview
- `detail` - Comprehensive analysis with key points
- `qa` - Question-answer format
- `extract` - Extract entities and structured data

**Response**:
```json
{
  "file_id": "uuid-here",
  "filename": "document.pdf",
  "file_type": "application/pdf",
  "file_size": 102400,
  "analysis": "This document discusses the quarterly financial results...",
  "summary": "Q3 financial report showing 15% revenue growth...",
  "key_points": [
    "Revenue increased 15% YoY",
    "Operating margin improved to 28%",
    "Customer base grew by 10,000 users"
  ],
  "entities": ["Q3 2024", "ABC Corporation", "$5.2M"],
  "sentiment": "positive",
  "language": "th",
  "tokens_used": 450,
  "process_time_ms": 2345.67,
  "timestamp": "2024-11-01T12:30:00Z"
}
```

**Example**:
```bash
# Analyze PDF with custom prompt
curl -X POST http://localhost:3001/api/file/analyze \
  -F "file=@contract.pdf" \
  -F "analysis_type=detail" \
  -F "prompt=Extract key terms, obligations, and dates" \
  -F "language=th" \
  -F "session_id=session_001"

# Analyze image
curl -X POST http://localhost:3001/api/file/analyze \
  -F "file=@chart.png" \
  -F "analysis_type=summary" \
  -F "prompt=Explain this chart"
```

#### `GET /file/history`
Retrieve file analysis history.

**Query Parameters**:
- `limit` (default: 20, max: 100)
- `offset` (default: 0)
- `file_type` (optional) - Filter by MIME type

**Response**:
```json
{
  "files": [
    {
      "file_id": "uuid",
      "filename": "report.pdf",
      "file_type": "application/pdf",
      "file_size": 102400,
      "analysis_type": "summary",
      "language": "th",
      "tokens_used": 450,
      "created_at": "2024-11-01T12:30:00Z"
    }
  ],
  "total": 23,
  "limit": 20,
  "offset": 0
}
```

---

### Audio

Speech-to-text and text-to-speech.

#### `POST /audio/transcribe`
Convert audio to text using Whisper API.

**Content-Type**: `multipart/form-data`

**Form Fields**:
- `audio` (required) - Audio file

**Supported Formats**: MP3, MP4, WAV, WEBM, M4A, OGG, FLAC
**Max Size**: 25 MB
**Max Duration**: ~30 minutes

**Response**:
```json
{
  "text": "Hello, this is a transcription of the audio file",
  "language": "en",
  "duration": 0.0,
  "confidence": 0.95,
  "timestamp": "2024-11-01T12:30:00Z"
}
```

**Example**:
```bash
curl -X POST http://localhost:3001/api/audio/transcribe \
  -F "audio=@voice.mp3"
```

#### `POST /audio/tts`
Convert text to speech.

**Request**:
```json
{
  "text": "สวัสดีครับ",
  "voice": "nova",
  "model": "tts-1",
  "response_format": "mp3",
  "speed": 1.0
}
```

**Parameters**:

| Field | Type | Default | Options |
|-------|------|---------|---------|
| `text` | string | - | Max 4096 characters |
| `voice` | string | nova | alloy, echo, fable, onyx, nova, shimmer |
| `model` | string | tts-1 | tts-1, tts-1-hd |
| `response_format` | string | mp3 | mp3, opus, aac, flac, wav, pcm |
| `speed` | float | 1.0 | 0.25 - 4.0 |

**Response** (Default - JSON with base64):
```json
{
  "audio_data": "//NExAA...",
  "format": "mp3",
  "duration": 1.5,
  "characters_used": 11,
  "voice": "nova",
  "timestamp": "2024-11-01T12:30:00Z"
}
```

**Alternative Response** (If `Accept: audio/mpeg` header sent):
- Returns binary audio stream
- Headers: `Content-Type: audio/mpeg`, `X-Audio-Duration`, `X-Characters-Used`

**Voice Characteristics**:

| Voice | Description |
|-------|-------------|
| alloy | Neutral, balanced |
| echo | Male voice |
| fable | British accent |
| onyx | Deep male voice |
| nova | Female voice (recommended) |
| shimmer | Soft female voice |

---

### WebSocket Streaming

Real-time streaming chat responses.

#### `WebSocket /api/chat/stream`

**Connection**:
```javascript
const ws = new WebSocket('ws://localhost:3001/api/chat/stream');
```

**Client → Server Message**:
```json
{
  "type": "message",
  "content": "What is AI?",
  "session_id": "session_123",
  "persona_id": 1,
  "system_prompt": "Optional custom system prompt",
  "file_ids": ["uuid1", "uuid2"]
}
```

**Server → Client Messages**:

**Chunk (streaming response)**:
```json
{
  "type": "chunk",
  "content": "AI is",
  "done": false
}
```

**Done (stream complete)**:
```json
{
  "type": "chunk",
  "content": "",
  "done": true,
  "message_id": "uuid-here",
  "tokens_used": 156
}
```

**Error**:
```json
{
  "type": "error",
  "error": "Error description"
}
```

**Features**:
- Real-time token-by-token streaming
- Auto-saves messages to database
- Supports conversation history (last 10 messages)
- Supports file attachments (max 5 files)
- Graceful disconnect handling

**Example**:
```javascript
const ws = new WebSocket('ws://localhost:3001/api/chat/stream');

ws.onopen = () => {
  ws.send(JSON.stringify({
    type: 'message',
    content: 'Tell me a story',
    session_id: 'session_123',
    persona_id: 1,
    file_ids: []
  }));
};

let fullResponse = '';

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);

  if (data.type === 'chunk') {
    if (data.done) {
      console.log('\n✅ Complete!');
      console.log('Message ID:', data.message_id);
      console.log('Tokens:', data.tokens_used);
      console.log('Full Response:', fullResponse);
    } else {
      process.stdout.write(data.content);
      fullResponse += data.content;
    }
  } else if (data.type === 'error') {
    console.error('❌ Error:', data.error);
  }
};
```

---

## 🗄️ Database Models

### Persona
```
Table: personas
- id (int, PK, AUTO_INCREMENT)
- name (varchar(100))
- system_prompt (text)
- expertise (varchar(100))
- description (text)
- icon (varchar(50))
- is_active (boolean, default: true)
- created_at (timestamp)
```

### Message
```
Table: messages
- id (uuid, PK, DEFAULT uuid_generate_v4())
- session_id (varchar(100), INDEXED)
- role (varchar(20): 'user'|'assistant'|'system')
- content (text)
- persona_id (int, FK → personas.id, NULLABLE)
- tokens_used (int, NULLABLE)
- file_attachments (jsonb, DEFAULT '[]')
- metadata (jsonb, DEFAULT '{}')
- created_at (timestamp)
```

**file_attachments JSONB structure**:
```json
[
  {
    "file_id": "uuid",
    "filename": "document.pdf",
    "file_type": "application/pdf",
    "file_size": 102400,
    "analysis_summary": "First 200 chars..."
  }
]
```

### FileAnalysis
```
Table: file_analyses
- id (uuid, PK, DEFAULT uuid_generate_v4())
- session_id (varchar(100), INDEXED, NULLABLE)
- file_name (varchar(255))
- file_type (varchar(100))
- file_size (bigint)
- file_path (varchar(500), NULLABLE)
- analysis_type (varchar(50): summary|detail|qa|extract)
- custom_prompt (text)
- language (varchar(10): th|en)
- analysis (text) - Full analysis
- summary (text) - Brief summary
- key_points (text[], ARRAY)
- entities (text[], ARRAY)
- sentiment (varchar(50))
- tokens_used (int)
- process_time_ms (float)
- created_at (timestamp)
- updated_at (timestamp)
- deleted_at (timestamp, NULLABLE) - Soft delete
- reanalysis_count (int, DEFAULT 0)
```

### Relationships

```
Persona (1) ──────→ Message (Many)
                    └─→ file_attachments (JSONB array)
                        └─→ References FileAnalysis by file_id

Message.session_id ←→ FileAnalysis.session_id (Linked by session)
```

---

## ⚠️ Error Handling

### Common Errors

| Status | Error | Cause | Solution |
|--------|-------|-------|----------|
| 400 | "message is required" | Empty message | Provide message content |
| 400 | "maximum 5 files allowed" | Too many files | Limit to 5 files |
| 404 | "persona with ID X not found" | Invalid persona_id | Use valid ID (1, 2, or 3) |
| 413 | File too large | Exceeds limit | Reduce file size |
| 415 | Unsupported file type | Invalid format | Use supported formats |
| 422 | Failed to parse file | Corrupted file | Check file integrity |
| 503 | OpenAI API unavailable | API down | Retry after delay |

### Error Response Format
```json
{
  "error": "Descriptive error message"
}
```

### Best Practices
- Always check HTTP status code
- Parse error messages for user display
- Implement retry logic for 503 errors (exponential backoff)
- Validate file size before upload
- Validate session_id format on client

---

## 💡 Best Practices

### Conversation History
✅ **DO**:
- Use consistent `session_id` per conversation
- Set `use_history: true` for context
- Limit history to recent messages (backend limits to 10)

❌ **DON'T**:
- Change `session_id` mid-conversation
- Include entire history manually in message

### File Attachments
✅ **DO**:
- Upload files first via `/file/analyze`
- Store returned `file_id`
- Send `file_ids` array in chat request
- Max 5 files per message

❌ **DON'T**:
- Send files directly in chat endpoint
- Exceed file size limits

### System Prompts
✅ **DO**:
- Use English for system prompts (better AI compliance)
- Be specific and provide examples
- Specify output language explicitly

**Example**:
```json
{
  "system_prompt": "You MUST respond in Thai language. You are a professional consultant. Always provide structured answers with: 1) Summary 2) Details 3) Recommendations. Example: 'สรุป: ..., รายละเอียด: ..., คำแนะนำ: ...'"
}
```

❌ **DON'T**:
- Use vague prompts
- Assume AI knows output language

### Performance
- Use WebSocket for long responses (streaming)
- Use REST API for simple Q&A
- Implement client-side caching for personas
- Paginate history requests

---

## 📊 Database Schema Visual

```
┌─────────────────┐
│    Persona      │
│  (id: 1,2,3)    │
└────────┬────────┘
         │
         │ 1:N (optional FK)
         ▼
┌─────────────────────────────────┐
│         Message                 │
│  - id (uuid)                    │
│  - session_id                   │
│  - role (user/assistant)        │
│  - content                      │
│  - persona_id (FK, nullable)    │
│  - file_attachments (jsonb) ────┐
└─────────────────────────────────┘
                                   │
                                   │ References by file_id
                                   ▼
                         ┌─────────────────────┐
                         │   FileAnalysis      │
                         │  - id (uuid)        │
                         │  - session_id       │
                         │  - file_name        │
                         │  - analysis         │
                         │  - summary          │
                         └─────────────────────┘
```

---

## 🔧 Environment Variables

Required in `.env.development`:

```env
# Server
PORT=3001
APP_ENV=development
APP_NAME=ChatBotAPI

# Database
DATABASE_URL=postgresql://user:password@localhost:5432/chatbot

# OpenAI
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o-mini
OPENAI_MAX_TOKENS=2000
OPENAI_TEMPERATURE=0.7

# CORS
CORS_ORIGIN=http://localhost:5173
```

---

## 📝 Changelog

### Version 2.0.0 (2025-11-01) - **Current**
- ✅ Added file attachment support in chat (`file_ids`)
- ✅ Backend automatically builds file context from UUIDs
- ✅ Updated base URL to port 3001
- ✅ Added JSONB `file_attachments` to Message model
- ✅ Improved WebSocket with file support
- ✅ Added session linking for files
- ✅ Comprehensive backend analysis
- ✅ Simplified and clarified documentation

### Version 1.2.0 (2025-10-31)
- Added conversation history (`use_history`, `session_id`)
- Added file analysis endpoints
- Implemented Vision API for images

### Version 1.0.0 (2025-10-28)
- Initial release
- Chat, persona, audio endpoints
- WebSocket streaming

---

## 🌐 Resources

- **OpenAI API**: https://platform.openai.com/docs
- **WebSocket Guide**: https://developer.mozilla.org/en-US/docs/Web/API/WebSocket
- **PostgreSQL JSONB**: https://www.postgresql.org/docs/current/datatype-json.html

---

**API Base URL**: `http://localhost:3001/api`
**WebSocket URL**: `ws://localhost:3001/api/chat/stream`
**Frontend**: `http://localhost:5173`
