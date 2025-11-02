# ChatBot API Documentation

**Base URL**: `http://localhost:3000/api`
**WebSocket**: `ws://localhost:3000/api/chat/stream`
**Version**: 2.1.0
**Last Updated**: 2025-11-02

---

## 📋 Table of Contents

1. [Quick Start](#quick-start)
2. [API Endpoints](#api-endpoints)
3. [Data Models](#data-models)
4. [Error Handling](#error-handling)
5. [Best Practices](#best-practices)

---

## 🚀 Quick Start

```bash
# Health check
curl http://localhost:3000/api/health

# List personas
curl http://localhost:3000/api/personas

# Send chat message
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "สวัสดี", "session_id": "test_001"}'

# Analyze file
curl -X POST http://localhost:3000/api/file/analyze \
  -F "file=@document.pdf" \
  -F "session_id=test_001" \
  -F "prompt=สรุปเอกสารนี้"
```

---

## 📡 API Endpoints

### 1. Health Check

#### `GET /health`
ตรวจสอบสถานะ API

**Response:**
```json
{
  "status": "ok",
  "message": "ChatBot API is running",
  "env": "development"
}
```

---

### 2. Personas (บุคลิก AI)

#### `GET /personas`
ดึงรายการบุคลิก AI ทั้งหมด

**Response:**
```json
{
  "personas": [
    {
      "id": 1,
      "name": "General Assistant",
      "system_prompt": "You are a helpful assistant...",
      "expertise": "general",
      "description": "ผู้ช่วย AI ทั่วไป",
      "icon": "🤖",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

#### `GET /personas/:id`
ดึงรายละเอียดบุคลิกพร้อมสถิติ

**Response:**
```json
{
  "id": 1,
  "name": "General Assistant",
  "stats": {
    "total_messages": 42,
    "avg_response_time": "2.3s"
  }
}
```

---

### 3. Chat (แชท)

#### `POST /chat`
ส่งข้อความและรับคำตอบจาก AI

**Request:**
```json
{
  "message": "อธิบาย Machine Learning",
  "session_id": "session_123",
  "persona_id": 2,
  "system_prompt": "คุณเป็นครูสอน AI",
  "use_history": true,
  "file_ids": ["uuid1", "uuid2"]
}
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `message` | string | ✅ | - | ข้อความ |
| `session_id` | string | No | auto | Session ID |
| `persona_id` | int | No | null | ID ของบุคลิก (1-3) |
| `system_prompt` | string | No | "" | คำสั่งกำหนดพฤติกรรม AI |
| `use_history` | bool | No | false | ใช้ประวัติการสนทนา |
| `file_ids` | array | No | [] | UUID ของไฟล์ (สูงสุด 5 ไฟล์) |

**Response:**
```json
{
  "message_id": "uuid",
  "session_id": "session_123",
  "reply": "Machine Learning คือ...",
  "persona": {
    "id": 2,
    "name": "Technology Expert",
    "icon": "💻"
  },
  "tokens_used": 145,
  "model": "gpt-4o-mini",
  "timestamp": "2024-11-02T12:30:00Z",
  "history_used": true,
  "history_count": 8
}
```

**หมายเหตุ:**
- ส่งไฟล์ได้สูงสุด 5 ไฟล์ต่อข้อความ
- ประวัติการสนทนา: 10 ข้อความล่าสุด
- Backend จะดึงเนื้อหาไฟล์มารวมให้อัตโนมัติ

#### `GET /chat/history`
ดึงประวัติการสนทนา

**Query Parameters:**
- `limit` (default: 50, max: 100)
- `offset` (default: 0)

**Response:**
```json
{
  "messages": [
    {
      "id": "uuid",
      "session_id": "session_123",
      "role": "user",
      "content": "สวัสดี",
      "file_attachments": [
        {
          "file_id": "uuid",
          "filename": "doc.pdf",
          "file_type": "application/pdf",
          "file_size": 102400
        }
      ],
      "created_at": "2024-11-02T12:00:00Z"
    }
  ],
  "total": 247,
  "limit": 50,
  "offset": 0
}
```

---

### 4. File Analysis (วิเคราะห์ไฟล์)

#### `POST /file/analyze`
อัปโหลดและวิเคราะห์ไฟล์ด้วย AI

**Content-Type:** `multipart/form-data`

**Form Fields:**

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `file` | file | ✅ | - | ไฟล์ที่ต้องการวิเคราะห์ |
| `analysis_type` | string | No | summary | `summary`, `detail`, `qa`, `extract` |
| `session_id` | string | No | null | เชื่อมกับการสนทนา |
| `prompt` | string | No | "" | คำสั่งวิเคราะห์ |
| `language` | string | No | th | `th` หรือ `en` |
| `system_prompt` | string | No | "" | กำหนดพฤติกรรม AI |
| `use_history` | bool | No | false | ใช้ประวัติการสนทนา |

**Supported File Types:**

| Type | Extensions | Max Size |
|------|-----------|----------|
| Text | .txt, .md | 10 MB |
| PDF | .pdf | 25 MB |
| Office | .docx, .xlsx, .pptx | 25 MB |
| Images | .jpg, .png, .gif, .webp | 20 MB |
| Code | .js, .py, .go, .java | 5 MB |
| Data | .json, .xml, .csv | 10 MB |

**Analysis Types:**
- `summary` - สรุปย่อ
- `detail` - วิเคราะห์แบบละเอียด พร้อมจุดสำคัญ
- `qa` - รูปแบบคำถาม-คำตอบ
- `extract` - ดึงข้อมูลสำคัญ (ชื่อ, วันที่, ตัวเลข)

**Response:**
```json
{
  "message_id": "uuid",
  "session_id": "session_123",
  "reply": "เอกสารนี้กล่าวถึงผลประกอบการรายไตรมาส...",
  "tokens_used": 450,
  "model": "gpt-4o-mini",
  "timestamp": "2024-11-02T12:30:00Z",
  "file_info": {
    "file_id": "uuid",
    "filename": "document.pdf",
    "file_type": "application/pdf",
    "file_size": 102400
  }
}
```

**Example:**
```bash
# วิเคราะห์ PDF พร้อม prompt
curl -X POST http://localhost:3000/api/file/analyze \
  -F "file=@contract.pdf" \
  -F "analysis_type=detail" \
  -F "prompt=สรุปสาระสำคัญและภาระผูกพัน" \
  -F "session_id=session_001"

# วิเคราะห์ไฟล์พร้อมใช้ประวัติการสนทนา
curl -X POST http://localhost:3000/api/file/analyze \
  -F "file=@data.xlsx" \
  -F "session_id=session_001" \
  -F "use_history=true" \
  -F "prompt=วิเคราะห์ตามบริบทการสนทนาก่อนหน้า"
```

**การบันทึก Messages:**
เมื่ออัปโหลดไฟล์ระบุ `session_id` ระบบจะบันทึก 2 messages:
1. **User message:**
   ```
   อัปโหลดไฟล์: document.pdf
   คำสั่ง: สรุปเอกสารนี้
   ```
2. **AI response:** ผลการวิเคราะห์

#### `GET /file/history`
ดึงประวัติการวิเคราะห์ไฟล์

**Query Parameters:**
- `limit` (default: 20, max: 100)
- `offset` (default: 0)
- `file_type` (optional)

**Response:**
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
      "created_at": "2024-11-02T12:30:00Z"
    }
  ],
  "total": 23,
  "limit": 20,
  "offset": 0
}
```

---

### 5. Audio (เสียง)

#### `POST /audio/transcribe`
แปลงเสียงเป็นข้อความ (Speech-to-Text)

**Content-Type:** `multipart/form-data`

**Form Fields:**
- `audio` (required) - ไฟล์เสียง

**Supported:** MP3, MP4, WAV, WEBM, M4A (max 25 MB)

**Response:**
```json
{
  "text": "สวัสดีครับ นี่คือข้อความจากไฟล์เสียง",
  "language": "th",
  "duration": 5.2,
  "confidence": 0.95,
  "timestamp": "2024-11-02T12:30:00Z"
}
```

#### `POST /audio/tts`
แปลงข้อความเป็นเสียง (Text-to-Speech)

**Request:**
```json
{
  "text": "สวัสดีครับ",
  "voice": "nova",
  "model": "tts-1",
  "speed": 1.0
}
```

| Field | Default | Options |
|-------|---------|---------|
| `text` | - | สูงสุด 4096 ตัวอักษร |
| `voice` | nova | alloy, echo, fable, onyx, nova, shimmer |
| `model` | tts-1 | tts-1, tts-1-hd |
| `speed` | 1.0 | 0.25 - 4.0 |

**Response (JSON):**
```json
{
  "audio_data": "base64_encoded_audio",
  "format": "mp3",
  "duration": 1.5,
  "characters_used": 11,
  "voice": "nova"
}
```

**Response (Binary):** หาก header `Accept: audio/mpeg`
- ส่งคืน binary audio stream
- Headers: `Content-Type: audio/mpeg`, `X-Audio-Duration`

---

### 6. WebSocket Streaming

#### `WebSocket /api/chat/stream`
รับคำตอบแบบ real-time streaming

**Connection:**
```javascript
const ws = new WebSocket('ws://localhost:3000/api/chat/stream');
```

**Client → Server:**
```json
{
  "type": "message",
  "content": "AI คืออะไร",
  "session_id": "session_123",
  "persona_id": 1,
  "system_prompt": "คุณเป็นครู",
  "file_ids": ["uuid1", "uuid2"]
}
```

**Server → Client:**

**Chunk (กำลังส่ง):**
```json
{
  "type": "chunk",
  "content": "AI คือ",
  "done": false
}
```

**Done (เสร็จสิ้น):**
```json
{
  "type": "chunk",
  "content": "",
  "done": true,
  "message_id": "uuid",
  "tokens_used": 156
}
```

**Error:**
```json
{
  "type": "error",
  "error": "คำอธิบายข้อผิดพลาด"
}
```

**Features:**
- Streaming แบบ token-by-token
- บันทึก messages อัตโนมัติ
- รองรับประวัติการสนทนา (10 ข้อความ)
- รองรับไฟล์แนบ (สูงสุด 5 ไฟล์)

---

## 🗄️ Data Models

### Message
```
Table: messages
- id (uuid, PK)
- session_id (varchar, indexed)
- role (user|assistant|system)
- content (text)
- persona_id (int, FK, nullable)
- tokens_used (int, nullable)
- file_attachments (jsonb, default '[]')
- created_at (timestamp)
- metadata (jsonb)
```

**file_attachments Structure:**
```json
[
  {
    "file_id": "uuid",
    "filename": "document.pdf",
    "file_type": "application/pdf",
    "file_size": 102400
  }
]
```

### FileAnalysis
```
Table: file_analyses
- id (uuid, PK)
- session_id (varchar, indexed, nullable)
- filename (varchar)
- file_type (varchar)
- file_size (bigint)
- analysis_type (summary|detail|qa|extract)
- custom_prompt (text)
- language (th|en)
- analysis (text)
- key_points (text[])
- entities (text[])
- sentiment (varchar)
- tokens_used (int)
- process_time_ms (float)
- created_at (timestamp)
- updated_at (timestamp)
- deleted_at (timestamp, soft delete)
- reanalysis_count (int)
```

### Persona
```
Table: personas
- id (int, PK, auto_increment)
- name (varchar)
- system_prompt (text)
- expertise (varchar)
- description (text)
- icon (varchar)
- is_active (boolean)
- created_at (timestamp)
```

---

## ⚠️ Error Handling

### HTTP Status Codes

| Code | Meaning | Common Causes |
|------|---------|---------------|
| 200 | Success | - |
| 400 | Bad Request | ข้อมูลไม่ครบ, รูปแบบไม่ถูกต้อง |
| 404 | Not Found | ไม่พบ persona หรือไฟล์ |
| 413 | Payload Too Large | ไฟล์ใหญ่เกินกำหนด |
| 415 | Unsupported Media | ประเภทไฟล์ไม่รองรับ |
| 422 | Unprocessable Entity | ไม่สามารถแยกข้อมูลจากไฟล์ |
| 500 | Internal Server Error | ข้อผิดพลาดระบบ |
| 503 | Service Unavailable | OpenAI API ไม่พร้อมใช้งาน |

### Error Response Format
```json
{
  "error": "คำอธิบายข้อผิดพลาด"
}
```

---

## 💡 Best Practices

### 1. Conversation History
✅ **ควรทำ:**
- ใช้ `session_id` เดียวกันตลอดการสนทนา
- ตั้งค่า `use_history: true` เพื่อความต่อเนื่อง
- Backend จำกัดประวัติที่ 10 ข้อความล่าสุด

❌ **ไม่ควรทำ:**
- เปลี่ยน `session_id` กลางคัน
- ส่งประวัติทั้งหมดในข้อความ

### 2. File Attachments
✅ **ควรทำ:**
1. อัปโหลดไฟล์ผ่าน `/file/analyze` ก่อน
2. เก็บ `file_id` ที่ได้รับ
3. ส่ง `file_ids` array ใน chat request
4. สูงสุด 5 ไฟล์ต่อข้อความ

❌ **ไม่ควรทำ:**
- ส่งไฟล์ตรงไปที่ chat endpoint
- อัปโหลดไฟล์ใหญ่เกินกำหนด

### 3. System Prompts
✅ **ควรทำ:**
- ใช้ภาษาอังกฤษในการเขียน system prompt
- ระบุรูปแบบผลลัพธ์อย่างชัดเจน
- ระบุภาษาที่ต้องการในการตอบ

**ตัวอย่างที่ดี:**
```json
{
  "system_prompt": "You MUST respond in Thai language. You are a professional consultant. Structure your answers: 1) สรุป 2) รายละเอียด 3) คำแนะนำ"
}
```

❌ **ไม่ควรทำ:**
- ใช้คำสั่งคลุมเครือ
- สมมติว่า AI รู้ภาษาที่ต้องการตอบ

### 4. Performance
- ใช้ WebSocket สำหรับคำตอบยาว (streaming)
- ใช้ REST API สำหรับคำถาม-คำตอบสั้น
- Cache personas ที่ client
- ใช้ pagination สำหรับ history

---

## 📊 Architecture Overview

```
┌─────────────┐
│   Persona   │
│  (3 types)  │
└──────┬──────┘
       │
       │ 1:N
       ▼
┌──────────────────────┐        ┌──────────────────┐
│      Message         │        │   FileAnalysis   │
│  - session_id        │◄───────│  - session_id    │
│  - file_attachments  │ refs   │  - analysis      │
│    (JSONB)           │        │  - key_points    │
└──────────────────────┘        └──────────────────┘
```

**Data Flow:**
1. อัปโหลดไฟล์ → FileAnalysis (เก็บผลวิเคราะห์)
2. ส่งข้อความ + file_ids → Message (เก็บการสนทนา)
3. Backend ดึง FileAnalysis ตาม file_ids → สร้าง context
4. ส่งไปยัง OpenAI → รับคำตอบ
5. บันทึก assistant message → ส่งกลับ client

---

## 🔧 Environment Variables

```env
# Server
PORT=3000
APP_ENV=development

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

### Version 2.1.0 (2025-11-02) - **Current**
- ❌ **ลบ `summary` field** จาก FileAnalysisResponse และ FileAnalysis model
- ✅ **เพิ่ม prompt ใน user message** เมื่ออัปโหลดไฟล์:
  ```
  อัปโหลดไฟล์: document.pdf
  คำสั่ง: สรุปเอกสารนี้
  ```
- ✅ ลบ `analysis_summary` จาก FileAttachment
- ✅ ปรับปรุง API response ให้กระชับขึ้น
- ✅ อัปเดตเอกสารให้สั้นและเข้าใจง่าย

### Version 2.0.0 (2025-11-01)
- ✅ เพิ่มการแนบไฟล์ใน chat (`file_ids`)
- ✅ Backend สร้าง context จากไฟล์อัตโนมัติ
- ✅ เพิ่ม `file_attachments` (JSONB) ใน Message
- ✅ รองรับ WebSocket พร้อมไฟล์
- ✅ เชื่อม session_id ระหว่าง chat และ file

### Version 1.2.0 (2025-10-31)
- เพิ่มประวัติการสนทนา
- เพิ่ม file analysis endpoints
- รองรับ Vision API

### Version 1.0.0 (2025-10-28)
- เปิดตัวครั้งแรก
- Chat, Persona, Audio endpoints
- WebSocket streaming

---

**API Base URL**: `http://localhost:3000/api`
**WebSocket URL**: `ws://localhost:3000/api/chat/stream`
**Frontend**: `http://localhost:5173`
