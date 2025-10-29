# 📡 ChatBot API Documentation

> **REST API และ WebSocket API สำหรับ ChatBot Application**
>
> Base URL: `http://localhost:3000/api`
>
> Version: 1.0.0

---

## 📋 Table of Contents

- [Authentication](#authentication)
- [Response Format](#response-format)
- [Error Handling](#error-handling)
- [Health Check](#health-check)
- [Chat Endpoints](#chat-endpoints)
- [Persona Endpoints](#persona-endpoints)
- [Audio Endpoints](#audio-endpoints)
- [WebSocket API](#websocket-api)
- [Rate Limiting](#rate-limiting)
- [Examples](#examples)

---

## 🔐 Authentication

**สถานะปัจจุบัน**: ❌ **No Authentication Required**

> ⚠️ **หมายเหตุ**: โปรเจกต์นี้สร้างขึ้นเพื่อการเรียนรู้ จึงไม่มีระบบ Authentication
>
> สำหรับ Production ควรเพิ่ม:
> - JWT Authentication
> - API Key validation
> - Rate limiting per user

---

## 📦 Response Format

### Success Response

```json
{
  "status": "success",
  "data": {
    // Response data here
  },
  "timestamp": "2025-10-28T10:30:00Z"
}
```

### Error Response

```json
{
  "status": "error",
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {}  // Optional
  },
  "timestamp": "2025-10-28T10:30:00Z"
}
```

---

## ⚠️ Error Handling

### HTTP Status Codes

| Status Code | Meaning | Description |
|-------------|---------|-------------|
| `200` | OK | Request successful |
| `201` | Created | Resource created successfully |
| `400` | Bad Request | Invalid request format or parameters |
| `404` | Not Found | Resource not found |
| `422` | Unprocessable Entity | Validation error |
| `429` | Too Many Requests | Rate limit exceeded |
| `500` | Internal Server Error | Server error |
| `503` | Service Unavailable | Service temporarily unavailable |

### Common Error Codes

| Error Code | Description |
|------------|-------------|
| `INVALID_INPUT` | Request body validation failed |
| `PERSONA_NOT_FOUND` | Specified persona ID doesn't exist |
| `OPENAI_ERROR` | OpenAI API error |
| `DATABASE_ERROR` | Database operation failed |
| `FILE_TOO_LARGE` | Uploaded file exceeds size limit |
| `INVALID_FILE_TYPE` | Uploaded file type not supported |

---

## 🏥 Health Check

### `GET /api/health`

ตรวจสอบสถานะของ API server

#### Request

```bash
GET /api/health
```

#### Response

**Status**: `200 OK`

```json
{
  "status": "ok",
  "message": "ChatBot API is running",
  "env": "development",
  "timestamp": "2025-10-28T10:30:00Z"
}
```

#### Example

```bash
curl http://localhost:3000/api/health
```

---

## 💬 Chat Endpoints

### 1. Send Chat Message (Sync)

`POST /api/chat`

ส่งข้อความไปยัง AI และรอรับคำตอบทั้งหมดกลับมา

#### Request

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "message": "สวัสดีครับ ช่วยแนะนำเกี่ยวกับ Go programming หน่อย",
  "persona_id": 2
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `message` | string | ✅ Yes | ข้อความจากผู้ใช้ (max 4000 characters) |
| `persona_id` | integer | ❌ No | ID ของ AI persona (default: 1) |

#### Response

**Status**: `200 OK`

```json
{
  "message_id": "550e8400-e29b-41d4-a716-446655440000",
  "reply": "สวัสดีครับ! ผมยินดีให้คำแนะนำเกี่ยวกับ Go programming ครับ\n\nGo เป็นภาษาโปรแกรมที่...",
  "persona": {
    "id": 2,
    "name": "Technology Expert",
    "icon": "💻"
  },
  "tokens_used": 245,
  "timestamp": "2025-10-28T10:30:15Z"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `message_id` | UUID | ID ของข้อความที่บันทึกใน database |
| `reply` | string | คำตอบจาก AI |
| `persona` | object | ข้อมูล persona ที่ใช้ตอบ |
| `tokens_used` | integer | จำนวน tokens ที่ใช้ (สำหรับคำนวณค่าใช้จ่าย) |
| `timestamp` | string (ISO 8601) | เวลาที่ตอบกลับ |

#### Example

```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "อธิบาย BMAD method หน่อย",
    "persona_id": 1
  }'
```

#### Error Responses

**400 Bad Request** - Invalid input
```json
{
  "error": "message is required and cannot be empty"
}
```

**404 Not Found** - Persona not found
```json
{
  "error": "persona with id 99 not found"
}
```

**500 Internal Server Error** - OpenAI API error
```json
{
  "error": "failed to get response from OpenAI API"
}
```

---

### 2. Get Chat History

`GET /api/chat/history`

ดึงประวัติการสนทนาล่าสุด

#### Request

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `limit` | integer | ❌ No | 50 | จำนวนข้อความที่ต้องการ (max 100) |
| `offset` | integer | ❌ No | 0 | เริ่มจากข้อความที่ (สำหรับ pagination) |

#### Response

**Status**: `200 OK`

```json
{
  "messages": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "role": "user",
      "content": "สวัสดีครับ",
      "persona_id": 1,
      "created_at": "2025-10-28T10:30:00Z"
    },
    {
      "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "role": "assistant",
      "content": "สวัสดีครับ! ผมยินดีให้ความช่วยเหลือครับ",
      "persona_id": 1,
      "created_at": "2025-10-28T10:30:05Z"
    }
  ],
  "total": 2,
  "limit": 50,
  "offset": 0
}
```

#### Example

```bash
# ดึง 10 ข้อความล่าสุด
curl "http://localhost:3000/api/chat/history?limit=10"

# Pagination - ข้อความที่ 11-20
curl "http://localhost:3000/api/chat/history?limit=10&offset=10"
```

---

## 👤 Persona Endpoints

### 1. Get All Personas

`GET /api/personas`

ดึงรายการ AI Personas ทั้งหมดที่ active

#### Request

```bash
GET /api/personas
```

#### Response

**Status**: `200 OK`

```json
{
  "personas": [
    {
      "id": 1,
      "name": "General Assistant",
      "system_prompt": "You are a helpful, friendly, and knowledgeable AI assistant...",
      "expertise": "general",
      "description": "A versatile AI assistant for general questions and conversations",
      "icon": "🤖",
      "is_active": true,
      "created_at": "2025-10-28T08:00:00Z"
    },
    {
      "id": 2,
      "name": "Technology Expert",
      "system_prompt": "You are a technology expert with deep knowledge...",
      "expertise": "technology",
      "description": "Specialized in programming, software architecture, and tech solutions",
      "icon": "💻",
      "is_active": true,
      "created_at": "2025-10-28T08:00:00Z"
    },
    {
      "id": 3,
      "name": "Business Advisor",
      "system_prompt": "You are a professional business consultant...",
      "expertise": "business",
      "description": "Expert in business strategy, entrepreneurship, and market analysis",
      "icon": "💼",
      "is_active": true,
      "created_at": "2025-10-28T08:00:00Z"
    }
  ],
  "total": 3
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Persona ID |
| `name` | string | ชื่อของ persona |
| `system_prompt` | string | System prompt ที่ส่งไปยัง OpenAI |
| `expertise` | string | สาขาความเชี่ยวชาญ |
| `description` | string | คำอธิบาย persona |
| `icon` | string | Emoji icon |
| `is_active` | boolean | สถานะเปิดใช้งาน |
| `created_at` | string | วันที่สร้าง |

#### Example

```bash
curl http://localhost:3000/api/personas
```

---

### 2. Get Persona by ID

`GET /api/personas/:id`

ดึงรายละเอียด persona ตาม ID

#### Request

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Persona ID |

#### Response

**Status**: `200 OK`

```json
{
  "id": 2,
  "name": "Technology Expert",
  "system_prompt": "You are a technology expert with deep knowledge in software development, programming, cloud computing, and IT infrastructure. Provide technical and accurate responses.",
  "expertise": "technology",
  "description": "Specialized in programming, software architecture, and tech solutions",
  "icon": "💻",
  "is_active": true,
  "created_at": "2025-10-28T08:00:00Z",
  "stats": {
    "total_messages": 1542,
    "avg_response_time": "2.3s"
  }
}
```

#### Example

```bash
curl http://localhost:3000/api/personas/2
```

#### Error Responses

**404 Not Found**
```json
{
  "error": "persona with id 99 not found"
}
```

---

## 🎤 Audio Endpoints

### 1. Transcribe Audio

`POST /api/audio/transcribe`

อัพโหลดไฟล์เสียงและแปลงเป็นข้อความด้วย OpenAI Whisper API

#### Request

**Headers:**
```
Content-Type: multipart/form-data
```

**Body (Form Data):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `audio` | file | ✅ Yes | ไฟล์เสียง (mp3, mp4, wav, webm, m4a) |

**Constraints:**
- Max file size: **25 MB**
- Max duration: **30 minutes**
- Supported formats: mp3, mp4, mpeg, mpga, m4a, wav, webm

#### Response

**Status**: `200 OK`

```json
{
  "text": "สวัสดีครับ ช่วยแนะนำเกี่ยวกับการเขียนโปรแกรมด้วย Go หน่อยครับ",
  "language": "th",
  "duration": 5.2,
  "confidence": 0.95,
  "timestamp": "2025-10-28T10:30:00Z"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `text` | string | ข้อความที่แปลงได้ |
| `language` | string | ภาษาที่ตรวจพบ (ISO 639-1 code) |
| `duration` | float | ความยาวของเสียง (seconds) |
| `confidence` | float | ความมั่นใจในการแปลง (0.0-1.0) |
| `timestamp` | string | เวลาที่ประมวลผล |

#### Example

```bash
# ด้วย curl
curl -X POST http://localhost:3000/api/audio/transcribe \
  -F "audio=@voice_message.mp3"

# ด้วย JavaScript (Browser)
const formData = new FormData()
formData.append('audio', audioBlob, 'recording.webm')

fetch('http://localhost:3000/api/audio/transcribe', {
  method: 'POST',
  body: formData
})
  .then(res => res.json())
  .then(data => console.log(data.text))
```

#### Error Responses

**400 Bad Request** - No file uploaded
```json
{
  "error": "audio file is required"
}
```

**413 Payload Too Large** - File too large
```json
{
  "error": "file size exceeds maximum allowed (25MB)"
}
```

**415 Unsupported Media Type** - Invalid file type
```json
{
  "error": "unsupported file format. Allowed: mp3, mp4, wav, webm, m4a"
}
```

**500 Internal Server Error** - Whisper API error
```json
{
  "error": "failed to transcribe audio"
}
```

---

## 🔌 WebSocket API

### Connect to WebSocket

`WS /api/chat/stream`

เชื่อมต่อ WebSocket สำหรับรับคำตอบแบบ streaming (real-time typing effect)

#### Connection

```javascript
const ws = new WebSocket('ws://localhost:3000/api/chat/stream')

ws.onopen = () => {
  console.log('WebSocket connected')
}

ws.onmessage = (event) => {
  const data = JSON.parse(event.data)
  console.log('Received:', data)
}

ws.onerror = (error) => {
  console.error('WebSocket error:', error)
}

ws.onclose = () => {
  console.log('WebSocket disconnected')
}
```

---

### Send Message

ส่งข้อความผ่าน WebSocket

#### Request Message Format

```json
{
  "type": "message",
  "content": "เล่านิทานให้ฟังหน่อย",
  "persona_id": 1
}
```

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | ✅ Yes | Message type (always "message") |
| `content` | string | ✅ Yes | ข้อความจากผู้ใช้ |
| `persona_id` | integer | ❌ No | Persona ID (default: 1) |

#### Response Message Format

**Streaming Chunks** (multiple messages):

```json
{
  "type": "chunk",
  "content": "กาล",
  "done": false
}
```

```json
{
  "type": "chunk",
  "content": "คราว",
  "done": false
}
```

```json
{
  "type": "chunk",
  "content": "หนึ่ง",
  "done": false
}
```

**Final Message**:

```json
{
  "type": "chunk",
  "content": "",
  "done": true,
  "message_id": "550e8400-e29b-41d4-a716-446655440000",
  "tokens_used": 156
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Message type ("chunk") |
| `content` | string | เนื้อหาของ chunk (คำหรือวลี) |
| `done` | boolean | `true` = streaming เสร็จสิ้น |
| `message_id` | UUID | ID ของข้อความ (มีเฉพาะตอนจบ) |
| `tokens_used` | integer | Tokens ที่ใช้ (มีเฉพาะตอนจบ) |

---

### Example: Full WebSocket Flow

```javascript
const ws = new WebSocket('ws://localhost:3000/api/chat/stream')

ws.onopen = () => {
  // ส่งข้อความ
  ws.send(JSON.stringify({
    type: 'message',
    content: 'เล่าเรื่องสั้นให้ฟังหน่อย',
    persona_id: 1
  }))
}

let fullResponse = ''

ws.onmessage = (event) => {
  const data = JSON.parse(event.data)

  if (data.type === 'chunk') {
    if (data.done) {
      // Streaming เสร็จสิ้น
      console.log('Full response:', fullResponse)
      console.log('Message ID:', data.message_id)
      console.log('Tokens used:', data.tokens_used)
    } else {
      // รับ chunk
      fullResponse += data.content
      console.log('Chunk:', data.content)
      // แสดงใน UI แบบ real-time
      updateChatUI(data.content)
    }
  }
}

ws.onerror = (error) => {
  console.error('Error:', error)
}
```

---

### Error Messages

WebSocket อาจส่ง error message:

```json
{
  "type": "error",
  "error": {
    "code": "OPENAI_ERROR",
    "message": "Failed to connect to OpenAI API"
  }
}
```

---

## ⏱️ Rate Limiting

### Current Limits

**สถานะปัจจุบัน**: ❌ **No Rate Limiting**

**แนะนำสำหรับ Production**:

| Endpoint | Limit | Window |
|----------|-------|--------|
| `POST /api/chat` | 20 requests | per minute |
| `POST /api/audio/transcribe` | 10 requests | per minute |
| `GET /api/personas` | 60 requests | per minute |
| WebSocket connections | 5 connections | per user |

### Rate Limit Headers (Future)

```
X-RateLimit-Limit: 20
X-RateLimit-Remaining: 15
X-RateLimit-Reset: 1635724800
```

### 429 Response

```json
{
  "error": "Rate limit exceeded. Try again in 45 seconds",
  "retry_after": 45
}
```

---

## 📝 Examples

### Example 1: Simple Chat Flow

```javascript
// 1. เลือก persona
const personasResponse = await fetch('http://localhost:3000/api/personas')
const personas = await personasResponse.json()

// 2. ส่งข้อความ
const chatResponse = await fetch('http://localhost:3000/api/chat', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    message: 'สวัสดีครับ',
    persona_id: personas.personas[0].id
  })
})

const chatResult = await chatResponse.json()
console.log('AI:', chatResult.reply)
```

---

### Example 2: Voice Message Flow

```javascript
// 1. บันทึกเสียง (MediaRecorder API)
const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
const mediaRecorder = new MediaRecorder(stream)
const audioChunks = []

mediaRecorder.ondataavailable = (event) => {
  audioChunks.push(event.data)
}

mediaRecorder.onstop = async () => {
  // 2. สร้าง audio blob
  const audioBlob = new Blob(audioChunks, { type: 'audio/webm' })

  // 3. Upload และ transcribe
  const formData = new FormData()
  formData.append('audio', audioBlob, 'recording.webm')

  const response = await fetch('http://localhost:3000/api/audio/transcribe', {
    method: 'POST',
    body: formData
  })

  const result = await response.json()
  console.log('Transcribed text:', result.text)

  // 4. ส่งเป็นข้อความ
  const chatResponse = await fetch('http://localhost:3000/api/chat', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      message: result.text,
      persona_id: 1
    })
  })

  const chatResult = await chatResponse.json()
  console.log('AI:', chatResult.reply)
}

// เริ่มบันทึก
mediaRecorder.start()

// หยุดบันทึก (หลังจาก 5 วินาที)
setTimeout(() => {
  mediaRecorder.stop()
  stream.getTracks().forEach(track => track.stop())
}, 5000)
```

---

### Example 3: Streaming Chat with WebSocket

```javascript
class ChatClient {
  constructor() {
    this.ws = null
    this.currentResponse = ''
  }

  connect() {
    this.ws = new WebSocket('ws://localhost:3000/api/chat/stream')

    this.ws.onopen = () => {
      console.log('Connected to chat server')
    }

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      this.handleMessage(data)
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    this.ws.onclose = () => {
      console.log('Disconnected')
      // Auto-reconnect
      setTimeout(() => this.connect(), 3000)
    }
  }

  sendMessage(message, personaId = 1) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.currentResponse = ''
      this.ws.send(JSON.stringify({
        type: 'message',
        content: message,
        persona_id: personaId
      }))
    }
  }

  handleMessage(data) {
    if (data.type === 'chunk') {
      if (data.done) {
        console.log('Complete response:', this.currentResponse)
        console.log('Tokens used:', data.tokens_used)
      } else {
        this.currentResponse += data.content
        // Update UI
        this.updateChatUI(data.content)
      }
    } else if (data.type === 'error') {
      console.error('Error:', data.error.message)
    }
  }

  updateChatUI(chunk) {
    // Append chunk to chat interface
    const chatBox = document.getElementById('chat-messages')
    const lastMessage = chatBox.lastElementChild
    if (lastMessage && lastMessage.classList.contains('ai-message')) {
      lastMessage.textContent += chunk
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
    }
  }
}

// Usage
const client = new ChatClient()
client.connect()

// ส่งข้อความ
client.sendMessage('เล่าเรื่องสั้นให้ฟังหน่อย', 1)
```

---

### Example 4: Error Handling

```javascript
async function sendChatMessage(message, personaId) {
  try {
    const response = await fetch('http://localhost:3000/api/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        message,
        persona_id: personaId
      })
    })

    // ตรวจสอบ HTTP status
    if (!response.ok) {
      const error = await response.json()

      switch (response.status) {
        case 400:
          throw new Error(`Invalid input: ${error.error}`)
        case 404:
          throw new Error(`Persona not found: ${error.error}`)
        case 429:
          throw new Error(`Rate limit exceeded. Try again later.`)
        case 500:
          throw new Error(`Server error: ${error.error}`)
        default:
          throw new Error(`Unexpected error: ${error.error}`)
      }
    }

    const data = await response.json()
    return data

  } catch (error) {
    console.error('Failed to send message:', error.message)

    // แสดง error ให้ user
    showErrorNotification(error.message)

    // Retry logic (optional)
    if (shouldRetry(error)) {
      console.log('Retrying in 3 seconds...')
      await delay(3000)
      return sendChatMessage(message, personaId)
    }

    throw error
  }
}

function shouldRetry(error) {
  // Retry on network errors or 500 errors
  return error.message.includes('network') ||
         error.message.includes('Server error')
}

function delay(ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}
```

---

## 🧪 Testing the API

### Using curl

```bash
# Health check
curl http://localhost:3000/api/health

# Get all personas
curl http://localhost:3000/api/personas

# Get persona by ID
curl http://localhost:3000/api/personas/1

# Send chat message
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"สวัสดี","persona_id":1}'

# Get chat history
curl "http://localhost:3000/api/chat/history?limit=10"

# Upload audio file
curl -X POST http://localhost:3000/api/audio/transcribe \
  -F "audio=@test.mp3"
```

### Using Postman

1. **Import Collection**: Create new collection "ChatBot API"
2. **Add Requests**:
   - GET Health Check
   - GET Personas
   - POST Send Message
   - POST Audio Transcribe
3. **Set Environment Variables**:
   - `base_url`: `http://localhost:3000/api`

### Using JavaScript (Node.js)

```javascript
const axios = require('axios')

const API_BASE_URL = 'http://localhost:3000/api'

// Test health
async function testHealth() {
  const response = await axios.get(`${API_BASE_URL}/health`)
  console.log('Health:', response.data)
}

// Test chat
async function testChat() {
  const response = await axios.post(`${API_BASE_URL}/chat`, {
    message: 'สวัสดีครับ',
    persona_id: 1
  })
  console.log('Chat response:', response.data.reply)
}

testHealth()
testChat()
```

---

## 📚 Additional Resources

### Related Documentation
- [BMAD.md](BMAD.md) - Architecture และ Design
- [WORKFLOW.md](WORKFLOW.md) - Development Workflow
- [README.md](README.md) - Project Overview

### External APIs
- [OpenAI Chat Completions API](https://platform.openai.com/docs/api-reference/chat)
- [OpenAI Whisper API](https://platform.openai.com/docs/api-reference/audio)

### Tools
- [Postman](https://www.postman.com/) - API Testing
- [WebSocket Test Client](https://www.websocket.org/echo.html) - WebSocket Testing
- [curl](https://curl.se/) - Command-line HTTP client

---

## 🔄 Changelog

### Version 1.0.0 (2025-10-28)
- ✅ Initial API documentation
- ✅ Health check endpoint
- ✅ Chat endpoints (sync)
- ✅ Persona endpoints
- ⏳ Audio transcription endpoint (planned)
- ⏳ WebSocket streaming (planned)
- ⏳ Chat history endpoint (planned)

---

## 📞 Support

หากพบปัญหาหรือมีคำถาม:

1. ตรวจสอบ [Error Handling](#error-handling) section
2. ดู [Examples](#examples) สำหรับ use cases ทั่วไป
3. เปิด issue ใน GitHub repository
4. อ่าน [BMAD.md](BMAD.md) สำหรับ architecture details

---

**📅 Last Updated**: 2025-10-28
**📝 Version**: 1.0.0
**🔗 Base URL**: http://localhost:3000/api
**📧 Contact**: [Your Email]
