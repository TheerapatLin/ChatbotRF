# ChatBot API Documentation

## โครงสร้างโปรเจ็ค

```
backend/
├── config/              # การตั้งค่า
│   ├── config.go        # โหลด environment variables
│   └── database.go      # เชื่อมต่อ PostgreSQL
├── controllers/         # จัดการ HTTP requests
│   ├── audio_controller.go
│   ├── chat_controller.go
│   ├── file_controller.go
│   ├── persona_controller.go
│   └── websocket_controller.go
├── middleware/          # HTTP middleware
│   └── logger.go
├── models/              # โครงสร้างข้อมูล Database
│   ├── file_analysis.go
│   ├── message.go
│   └── persona.go
├── repositories/        # เข้าถึงข้อมูล Database
│   ├── file_analysis_repository.go
│   ├── message_repository.go
│   └── persona_repository.go
├── routes/              # กำหนดเส้นทาง API
│   └── routes.go
├── services/            # Business logic
│   ├── context_service.go    # สร้าง context สำหรับ AI
│   ├── file_service.go       # วิเคราะห์ไฟล์
│   ├── openai_service.go     # เชื่อมต่อ OpenAI API
│   ├── tts_service.go        # Text-to-Speech
│   └── whisper_service.go    # Speech-to-Text
├── main.go              # Entry point
├── .env.development     # ตั้งค่า Environment
└── go.mod               # Go dependencies
```

**Architecture Pattern:** Repository-Service-Controller
- **Repository:** เข้าถึงข้อมูลจาก Database
- **Service:** Business logic และการเชื่อมต่อ External APIs
- **Controller:** รับ HTTP request และส่ง response

---

## วิธีเริ่ม Server

### 1. ติดตั้ง Dependencies
```bash
cd backend
go mod download
```

### 2. เริ่ม Database (PostgreSQL)
```bash
docker-compose up -d
```

### 3. ตั้งค่า Environment Variables
แก้ไขไฟล์ `.env.development`:
```env
OPENAI_API_KEY=your_openai_api_key_here
PORT=3001
DATABASE_URL=postgres://chatbot_user:admin123@localhost:5432/chatbot_db?sslmode=disable
```

### 4. รัน Server
```bash
# รันแบบปกติ
go run main.go

# รันแบบ Hot Reload (สำหรับ Development)
air
```

Server จะเปิดที่: `http://localhost:3001`

---

## API Endpoints

### Health Check
**GET** `/api/health`

ตรวจสอบสถานะของ API

**Response:**
```json
{
  "status": "healthy",
  "environment": "development",
  "timestamp": "2025-11-03T19:00:00Z"
}
```

---

## 1. Personas API

### 1.1 ดึงรายการ Personas ทั้งหมด
**GET** `/api/personas`

ดึงรายการ AI personas ที่ active ทั้งหมดพร้อม configuration ครบถ้วน

**Response:**
```json
{
  "personas": [
    {
      "id": 1,
      "name": "ผู้ช่วยทั่วไป",
      "description": "ผู้ช่วย AI อเนกประสงค์ที่ตอบคำถามได้หลากหลาย",
      "system_prompt": "คุณคือผู้ช่วย AI ที่เป็นมิตรและให้ข้อมูลที่ถูกต้อง...",
      "tone": "friendly",
      "style": "conversational",
      "expertise": "general",
      "temperature": 0.7,
      "max_tokens": 2000,
      "model": "gpt-4o-mini",
      "language_setting": {
        "default_language": "th",
        "response_style": "casual",
        "language_code": "th-TH"
      },
      "guardrails": {
        "block_profanity": true,
        "block_sensitive": true,
        "allowed_topics": ["general", "technology", "education"],
        "blocked_topics": ["politics", "religion"],
        "max_response_length": 4000,
        "require_moderation": false
      },
      "icon": "🤖",
      "is_active": true,
      "created_at": "2025-11-03T19:00:00Z",
      "updated_at": "2025-11-03T19:00:00Z"
    }
  ]
}
```

**วิธีทดสอบ:**
```bash
curl http://localhost:3001/api/personas
```

---

### 1.2 ดึงข้อมูล Persona ตามไอดี
**GET** `/api/personas/:id`

ดึงข้อมูลละเอียดของ persona พร้อมสถิติการใช้งาน

**Response:**
```json
{
  "id": 1,
  "name": "ผู้ช่วยทั่วไป",
  "description": "ผู้ช่วย AI อเนกประสงค์ที่ตอบคำถามได้หลากหลาย",
  "system_prompt": "คุณคือผู้ช่วย AI ที่เป็นมิตร...",
  "tone": "friendly",
  "style": "conversational",
  "expertise": "general",
  "temperature": 0.7,
  "max_tokens": 2000,
  "model": "gpt-4o-mini",
  "language_setting": {
    "default_language": "th",
    "response_style": "casual",
    "language_code": "th-TH"
  },
  "guardrails": {
    "block_profanity": true,
    "block_sensitive": true,
    "allowed_topics": ["general", "technology"],
    "blocked_topics": ["politics"],
    "max_response_length": 4000,
    "require_moderation": false
  },
  "icon": "🤖",
  "is_active": true,
  "created_at": "2025-11-03T19:00:00Z",
  "updated_at": "2025-11-03T19:00:00Z",
  "stats": {
    "total_messages": 150,
    "avg_response_time": "2.3s"
  }
}
```

**วิธีทดสอบ:**
```bash
curl http://localhost:3001/api/personas/1
```

---

## 2. Chat API

### 2.1 ส่งข้อความแบบไม่ Streaming
**POST** `/api/chat`

ส่งข้อความและรับคำตอบจาก AI (Non-streaming)

**Request Body:**
```json
{
  "message": "สวัสดี บอกเกี่ยวกับตัวคุณหน่อย",
  "session_id": "session_123",
  "persona_id": 1,
  "system_prompt": "คุณคือผู้เชี่ยวชาญด้านการเงิน",
  "use_history": true,
  "max_tokens": 2000,
  "temperature": 0.7,
  "file_ids": ["file_uuid_123"]
}
```

**Parameters:**
- `message` (required): ข้อความจากผู้ใช้
- `session_id` (optional): ID สำหรับเก็บ conversation history
- `persona_id` (optional): ID ของ persona ที่ต้องการใช้
- `system_prompt` (optional): Custom system prompt
- `use_history` (optional): ใช้ประวัติการสนทนา (default: false)
- `file_ids` (optional): Array ของ file IDs ที่ต้องการอ้างอิง

**Response:**
```json
{
  "message_id": "msg_uuid_123",
  "session_id": "session_123",
  "reply": "สวัสดีครับ ผมคือ AI ที่พร้อมช่วยเหลือคุณ...",
  "persona": {
    "id": 1,
    "name": "ผู้ช่วยทั่วไป",
    "expertise": "ทั่วไป",
    "icon": "🤖",
    "description": "ผู้ช่วยอเนกประสงค์"
  },
  "tokens_used": 245,
  "model": "gpt-4o-mini",
  "timestamp": "2025-11-03T19:00:00Z",
  "history_used": true,
  "history_count": 5
}
```

**วิธีทดสอบ:**
```bash
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "สวัสดี",
    "session_id": "test_session",
    "use_history": false
  }'
```

---

### 2.2 ลบข้อความทั้งหมด
**DELETE** `/api/chat`

ลบข้อความทั้งหมดในฐานข้อมูล (ใช้สำหรับล้างข้อมูลทั้งหมด)

**Response:**
```json
{
  "message": "All messages deleted successfully"
}
```

**วิธีทดสอบ:**
```bash
curl -X DELETE http://localhost:3001/api/chat
```

**⚠️ คำเตือน:** คำสั่งนี้จะลบข้อความทั้งหมดในฐานข้อมูล ใช้ด้วยความระมัดระวัง!

---

### 2.3 ดึงประวัติการสนทนา
**GET** `/api/chat/history?limit=50&offset=0`

ดึงประวัติข้อความทั้งหมด พร้อม pagination

**Query Parameters:**
- `limit` (optional): จำนวนข้อความต่อหน้า (default: 50, max: 100)
- `offset` (optional): เริ่มที่ข้อความลำดับที่ (default: 0)

**Response:**
```json
{
  "messages": [
    {
      "id": "msg_uuid_123",
      "role": "user",
      "content": "สวัสดี",
      "persona_id": 1,
      "created_at": "2025-11-03T19:00:00Z"
    },
    {
      "id": "msg_uuid_124",
      "role": "assistant",
      "content": "สวัสดีครับ",
      "persona_id": 1,
      "created_at": "2025-11-03T19:00:01Z"
    }
  ],
  "total": 150,
  "limit": 50,
  "offset": 0
}
```

**วิธีทดสอบ:**
```bash
curl "http://localhost:3001/api/chat/history?limit=10&offset=0"
```

---

### 2.4 Chat แบบ Streaming (WebSocket)
**WebSocket** `/api/chat/stream`

เชื่อมต่อ WebSocket สำหรับรับคำตอบแบบ real-time streaming

**ส่งข้อความ (WebSocket Message):**
```json
{
  "type": "message",
  "content": "เขียนเรื่องสั้นให้หน่อย",
  "persona_id": 1,
  "session_id": "session_123",
  "system_prompt": "คุณคือนักเขียน",
  "file_ids": ["file_uuid_123"]
}
```

**รับข้อความ (WebSocket Response):**
```json
// Streaming chunks
{
  "type": "chunk",
  "content": "กาลครั้งหนึ่ง",
  "done": false
}

// Final message
{
  "type": "chunk",
  "content": "",
  "done": true,
  "message_id": "msg_uuid_125",
  "tokens_used": 320
}
```

**วิธีทดสอบ (JavaScript):**
```javascript
const ws = new WebSocket('ws://localhost:3001/api/chat/stream');

ws.onopen = () => {
  ws.send(JSON.stringify({
    type: 'message',
    content: 'สวัสดี',
    session_id: 'test_session'
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log(data);
};
```

---

## 3. File Upload API

### 3.1 อัปโหลดไฟล์ (รองรับหลายไฟล์)
**POST** `/api/file/upload`

อัปโหลดไฟล์ไปยัง server โดยไม่มีการวิเคราะห์ด้วย AI (รองรับสูงสุด 5 ไฟล์ต่อครั้ง)

**Content-Type:** `multipart/form-data`

**Form Parameters:**
- `files` (required): ไฟล์ที่ต้องการอัปโหลด (สูงสุด 5 ไฟล์)

**รองรับไฟล์ทุกประเภท:**
- เอกสาร: PDF, DOCX, XLSX, PPTX, TXT, MD, CSV, JSON, XML
- รูปภาพ: JPG, PNG, GIF, WEBP, BMP
- โค้ด: JS, PY, GO, JAVA, CPP, etc.
- อื่นๆ: ZIP, RAR, MP3, MP4, etc.

**Response (อัปโหลดสำเร็จทั้งหมด):**
```json
{
  "success": 3,
  "failed": 0,
  "total": 3,
  "uploaded_files": [
    {
      "file_id": "uuid-1",
      "file_name": "report.pdf",
      "storage_path": "./uploads/uuid-1_report.pdf",
      "mime_type": "application/pdf",
      "file_size": 245678,
      "uploaded_at": "2025-11-03T19:00:00Z"
    },
    {
      "file_id": "uuid-2",
      "file_name": "image.jpg",
      "storage_path": "./uploads/uuid-2_image.jpg",
      "mime_type": "image/jpeg",
      "file_size": 156789,
      "uploaded_at": "2025-11-03T19:00:00Z"
    },
    {
      "file_id": "uuid-3",
      "file_name": "data.xlsx",
      "storage_path": "./uploads/uuid-3_data.xlsx",
      "mime_type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
      "file_size": 89012,
      "uploaded_at": "2025-11-03T19:00:00Z"
    }
  ]
}
```

**Response (อัปโหลดบางไฟล์สำเร็จ):**
```json
{
  "success": 2,
  "failed": 1,
  "total": 3,
  "uploaded_files": [
    {
      "file_id": "uuid-1",
      "file_name": "report.pdf",
      "storage_path": "./uploads/uuid-1_report.pdf",
      "mime_type": "application/pdf",
      "file_size": 245678,
      "uploaded_at": "2025-11-03T19:00:00Z"
    },
    {
      "file_id": "uuid-2",
      "file_name": "image.jpg",
      "storage_path": "./uploads/uuid-2_image.jpg",
      "mime_type": "image/jpeg",
      "file_size": 156789,
      "uploaded_at": "2025-11-03T19:00:00Z"
    }
  ],
  "failed_files": [
    {
      "file_name": "corrupted.docx",
      "error": "failed to save file to disk"
    }
  ]
}
```

**Error Response (เกินจำนวนไฟล์ที่อนุญาต):**
```json
{
  "error": "maximum 5 files allowed per upload"
}
```

**วิธีทดสอบ (อัปโหลดไฟล์เดียว):**
```bash
curl -X POST http://localhost:3001/api/file/upload \
  -F "files=@document.pdf"
```

**วิธีทดสอบ (อัปโหลดหลายไฟล์):**
```bash
curl -X POST http://localhost:3001/api/file/upload \
  -F "files=@document.pdf" \
  -F "files=@image.jpg" \
  -F "files=@data.xlsx"
```

---

### 3.2 ดึงประวัติการอัปโหลดไฟล์
**GET** `/api/file/history?limit=20&offset=0&file_type=all`

ดึงประวัติการอัปโหลดไฟล์ทั้งหมด

**Query Parameters:**
- `limit` (optional): จำนวนรายการต่อหน้า (default: 20, max: 100)
- `offset` (optional): เริ่มที่รายการลำดับที่ (default: 0)
- `file_type` (optional): กรองตามประเภทไฟล์ เช่น `application/pdf`

**Response:**
```json
{
  "files": [
    {
      "file_id": "uuid-123",
      "file_name": "report.pdf",
      "storage_path": "./uploads/uuid-123_report.pdf",
      "mime_type": "application/pdf",
      "file_size": 245678,
      "uploaded_at": "2025-11-03T19:00:00Z"
    }
  ],
  "total": 25,
  "limit": 20,
  "offset": 0
}
```

**วิธีทดสอบ:**
```bash
curl "http://localhost:3001/api/file/history?limit=10"
```

---

### 3.3 ลบไฟล์ทั้งหมด
**DELETE** `/api/file/uploads`

ลบข้อมูลไฟล์ทั้งหมดในฐานข้อมูล (ใช้สำหรับล้างข้อมูลทั้งหมด)

**Response:**
```json
{
  "message": "All file records deleted successfully"
}
```

**วิธีทดสอบ:**
```bash
curl -X DELETE http://localhost:3001/api/file/uploads
```

**⚠️ คำเตือน:**
- คำสั่งนี้จะลบ**บันทึกไฟล์ทั้งหมด**ในฐานข้อมูล
- **ไฟล์จริงบน disk จะยังคงอยู่** ที่ `./uploads/`
- ใช้ด้วยความระมัดระวัง!

---

## 4. Audio API

### 4.1 แปลงเสียงเป็นข้อความ (Speech-to-Text)
**POST** `/api/audio/transcribe`

อัปโหลดไฟล์เสียงและแปลงเป็นข้อความด้วย Whisper API

**Content-Type:** `multipart/form-data`

**Form Parameters:**
- `audio` (required): ไฟล์เสียง (max 25MB)

**รองรับไฟล์:** mp3, mp4, mpeg, mpga, m4a, wav, webm

**Response:**
```json
{
  "text": "สวัสดีครับ วันนี้อากาศดีมาก",
  "language": "th",
  "duration": 3.5,
  "confidence": 0.95,
  "timestamp": "2025-11-03T19:00:00Z"
}
```

**วิธีทดสอบ:**
```bash
curl -X POST http://localhost:3001/api/audio/transcribe \
  -F "audio=@voice.mp3"
```

---

### 4.2 แปลงข้อความเป็นเสียง (Text-to-Speech)
**POST** `/api/audio/tts`

แปลงข้อความเป็นไฟล์เสียงด้วย OpenAI TTS API

**Request Body:**
```json
{
  "text": "สวัสดีครับ ยินดีต้อนรับ",
  "voice": "alloy",
  "model": "tts-1",
  "response_format": "mp3",
  "speed": 1.0
}
```

**Parameters:**
- `text` (required): ข้อความที่ต้องการแปลง (max 4096 characters)
- `voice` (optional): เสียง - `alloy`, `echo`, `fable`, `onyx`, `nova`, `shimmer` (default: alloy)
- `model` (optional): โมเดล - `tts-1`, `tts-1-hd` (default: tts-1)
- `response_format` (optional): รูปแบบ - `mp3`, `opus`, `aac`, `flac`, `wav`, `pcm` (default: mp3)
- `speed` (optional): ความเร็ว 0.25-4.0 (default: 1.0)

**Response (JSON):**
```json
{
  "audio_data": "base64_encoded_audio_data...",
  "format": "mp3",
  "duration": 2.5,
  "characters_used": 25,
  "voice": "alloy",
  "timestamp": "2025-11-03T19:00:00Z"
}
```

**Response (Binary):**
หากส่ง Header `Accept: audio/mpeg` จะได้ไฟล์เสียงโดยตรง

**วิธีทดสอบ:**
```bash
# รับ JSON response
curl -X POST http://localhost:3001/api/audio/tts \
  -H "Content-Type: application/json" \
  -d '{
    "text": "สวัสดีครับ",
    "voice": "nova"
  }'

# รับไฟล์เสียงโดยตรง
curl -X POST http://localhost:3001/api/audio/tts \
  -H "Content-Type: application/json" \
  -H "Accept: audio/mpeg" \
  -d '{
    "text": "สวัสดีครับ"
  }' \
  --output audio.mp3
```

---

## โครงสร้าง Database (PostgreSQL)

### Table: `personas`
เก็บข้อมูล AI personalities พร้อม configuration ครบถ้วน

| Column | Type | Description |
|--------|------|-------------|
| id | SERIAL PRIMARY KEY | ID |
| name | VARCHAR(100) UNIQUE | ชื่อ persona (ไม่ซ้ำ) |
| description | TEXT | คำอธิบายละเอียด |
| system_prompt | TEXT NOT NULL | System prompt สำหรับ AI |
| tone | VARCHAR(50) | โทนเสียง (e.g., friendly, professional) |
| style | VARCHAR(50) | สไตล์การตอบ (e.g., concise, detailed) |
| expertise | VARCHAR(100) | ความเชี่ยวชาญ (e.g., technology, healthcare) |
| temperature | DECIMAL(3,2) | ค่า temperature (0.0-2.0, default: 0.7) |
| max_tokens | INTEGER | จำนวน tokens สูงสุด (default: 2000) |
| model | VARCHAR(50) | โมเดล AI (default: gpt-4o-mini) |
| language_setting | JSONB | การตั้งค่าภาษา (JSON) |
| guardrails | JSONB | กฎการกรองเนื้อหา (JSON) |
| icon | VARCHAR(50) | Emoji icon |
| is_active | BOOLEAN | สถานะ active (default: true) |
| created_at | TIMESTAMP | วันที่สร้าง |
| updated_at | TIMESTAMP | วันที่อัปเดตล่าสุด |

**JSON Schema - language_setting:**
```json
{
  "default_language": "th",      // ภาษาหลัก (th, en, etc.)
  "response_style": "casual",    // สไตล์การตอบ (formal, casual, professional)
  "language_code": "th-TH"       // รหัสภาษา ISO 639-1
}
```

**JSON Schema - guardrails:**
```json
{
  "block_profanity": true,        // บล็อกคำหยาบ
  "block_sensitive": true,        // บล็อกเนื้อหาละเอียดอ่อน
  "allowed_topics": [],           // หัวข้อที่อนุญาต
  "blocked_topics": [],           // หัวข้อที่ห้าม
  "max_response_length": 4000,    // ความยาวคำตอบสูงสุด
  "require_moderation": false     // ต้องการ moderation
}
```

---

### Table: `messages`
เก็บข้อความการสนทนา

| Column | Type | Description |
|--------|------|-------------|
| id | UUID PRIMARY KEY | Message ID |
| session_id | VARCHAR(255) | Session ID สำหรับจัดกลุ่มการสนทนา |
| role | VARCHAR(50) | บทบาท: user, assistant, system |
| content | TEXT | เนื้อหาข้อความ |
| persona_id | INTEGER | Foreign key → personas.id |
| tokens_used | INTEGER | จำนวน tokens ที่ใช้ |
| file_attachments | JSONB | Array ของไฟล์แนบ |
| created_at | TIMESTAMP | วันที่สร้าง |
| metadata | JSONB | ข้อมูลเพิ่มเติม |

**Indexes:**
- `session_id` - ค้นหาประวัติการสนทนา
- `created_at` - เรียงลำดับตามเวลา

---

### Table: `file_analyses`
เก็บผลการวิเคราะห์ไฟล์

| Column | Type | Description |
|--------|------|-------------|
| id | UUID PRIMARY KEY | File ID |
| session_id | VARCHAR(255) | เชื่อมโยงกับการสนทนา |
| file_name | VARCHAR(500) | ชื่อไฟล์ |
| file_type | VARCHAR(100) | MIME type |
| file_size | BIGINT | ขนาดไฟล์ (bytes) |
| file_path | VARCHAR(1000) | Path ที่เก็บไฟล์ |
| analysis_type | VARCHAR(50) | ประเภทการวิเคราะห์ |
| custom_prompt | TEXT | คำสั่งจากผู้ใช้ |
| language | VARCHAR(10) | ภาษา |
| analysis | TEXT | ผลการวิเคราะห์ |
| key_points | TEXT[] | Array ของประเด็นสำคัญ |
| entities | TEXT[] | Array ของ entities ที่พบ |
| sentiment | VARCHAR(50) | Sentiment analysis |
| tokens_used | INTEGER | จำนวน tokens ที่ใช้ |
| process_time_ms | FLOAT | เวลาประมวลผล (ms) |
| reanalysis_count | INTEGER | จำนวนครั้งที่วิเคราะห์ซ้ำ |
| created_at | TIMESTAMP | วันที่สร้าง |
| updated_at | TIMESTAMP | วันที่แก้ไขล่าสุด |
| deleted_at | TIMESTAMP | Soft delete timestamp |

**Indexes:**
- `session_id` - ค้นหาไฟล์ในการสนทนา
- `file_type` - กรองตามประเภทไฟล์
- `created_at` - เรียงตามเวลา

---

## การทำงานของ Functions สำคัญ

### 1. Chat Flow (controllers/chat_controller.go)
```
HandleChat()
  ├─→ parseRequest() - ตรวจสอบข้อมูลที่รับเข้ามา
  ├─→ getPersonaInfo() - ดึงข้อมูล persona และ system prompt
  ├─→ getOrGenerateSessionID() - สร้าง session ID ถ้ายังไม่มี
  ├─→ buildMessages() - สร้าง messages array
  │    ├─→ ContextService.BuildContextWithHistory() - ดึงประวัติการสนทนา (ถ้า use_history=true)
  │    └─→ ContextService.BuildFileContext() - เพิ่ม context จากไฟล์
  ├─→ callOpenAI() - เรียก OpenAI API
  ├─→ saveMessages() - บันทึกข้อความลง database
  └─→ buildResponse() - สร้าง response ส่งกลับ
```

### 2. File Analysis Flow (controllers/file_controller.go)
```
AnalyzeFile()
  ├─→ parseFileRequest() - ตรวจสอบไฟล์และพารามิเตอร์
  ├─→ ตรวจสอบประเภทไฟล์
  │    ├─→ รูปภาพ: analyzeImageFile() → FileService.AnalyzeImage()
  │    └─→ เอกสาร: FileService.AnalyzeFile()
  │         ├─→ ParseFileContent() - แปลงไฟล์เป็นข้อความ
  │         ├─→ ContextService.BuildContextWithHistory() - เพิ่ม history (ถ้า use_history=true)
  │         └─→ OpenAIService - ส่งไปวิเคราะห์
  ├─→ saveFileAnalysis() - บันทึกผลลง file_analyses table
  ├─→ saveFileAnalysisMessages() - บันทึกข้อความลง messages table
  └─→ ส่ง response กลับ
```

### 3. WebSocket Streaming Flow (controllers/websocket_controller.go)
```
HandleStreamingChat()
  └─→ Message loop
       └─→ handleMessage()
            ├─→ ตรวจสอบข้อความ
            ├─→ ดึงข้อมูล persona
            ├─→ สร้าง system prompt
            ├─→ BuildContextWithHistory() - สร้าง context พร้อม history
            ├─→ OpenAIService.CreateStreamingChatCompletion() - เริ่ม streaming
            ├─→ วนลูปรับ chunks และส่งกลับ client
            │    └─→ sendChunk() - ส่ง chunk แต่ละส่วน
            ├─→ บันทึกข้อความลง database
            └─→ sendDone() - ส่ง completion message
```

### 4. Context Service (services/context_service.go)
```
BuildContextWithHistory()
  ├─→ MessageRepository.FindBySessionID() - ดึงประวัติการสนทนา
  ├─→ สร้าง system message (ถ้ามี)
  ├─→ แปลง messages เป็น OpenAI format
  ├─→ BuildFileContext() - เพิ่ม context จากไฟล์ (ถ้ามี)
  └─→ เพิ่ม current message
```

---

## สิ่งที่ควรเพิ่มในอนาคต

### 1. Authentication & Authorization
- JWT token authentication
- User management system
- Role-based access control (RBAC)
- API rate limiting per user

### 2. Advanced Features
- ✅ File re-analysis endpoint (ยังไม่ implement)
- Message editing และ deletion
- Conversation branching (fork conversations)
- Export conversation เป็น PDF/DOCX
- Search functionality ใน chat history
- Tags และ categories สำหรับ conversations

### 3. Performance Optimization
- Redis caching สำหรับ frequently accessed data
- Database query optimization
- Connection pooling
- Response compression (gzip)
- CDN สำหรับ static files

### 4. Monitoring & Analytics
- Request/response logging
- Error tracking (Sentry)
- Performance metrics (Prometheus)
- Usage analytics dashboard
- Token usage tracking และ cost estimation

### 5. File Handling
- File storage service (S3/MinIO)
- Larger file support (chunking)
- File format conversion
- OCR สำหรับรูปภาพ
- Audio/Video transcription improvements

### 6. AI Enhancements
- Multiple AI providers (Anthropic, Google Gemini)
- Custom fine-tuned models
- RAG (Retrieval-Augmented Generation)
- Vector database integration
- Conversation memory optimization

### 7. Testing
- Unit tests coverage > 80%
- Integration tests
- E2E tests
- Load testing
- Security testing (OWASP)

### 8. DevOps
- CI/CD pipeline
- Docker containerization
- Kubernetes deployment
- Auto-scaling configuration
- Backup และ disaster recovery

### 9. Security
- Input validation และ sanitization
- SQL injection prevention (GORM ป้องกันอยู่แล้ว)
- XSS protection
- CSRF protection
- Rate limiting
- API key management
- Secrets management (Vault)

### 10. Documentation
- OpenAPI/Swagger specification
- Postman collection
- Architecture diagrams
- Deployment guide
- API versioning strategy

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 3001 |
| APP_ENV | Environment (development/production) | development |
| APP_NAME | Application name | ChatBotAPI |
| OPENAI_API_KEY | OpenAI API key | - |
| OPENAI_MODEL | OpenAI model name | gpt-4o-mini |
| OPENAI_MAX_TOKENS | Max tokens per request | 2000 |
| OPENAI_TEMPERATURE | AI temperature (0-2) | 0.7 |
| DATABASE_URL | PostgreSQL connection string | - |
| CORS_ORIGIN | Allowed CORS origins | localhost:5173,... |

---

## Error Codes

| HTTP Status | Description |
|-------------|-------------|
| 200 | Success |
| 206 | Partial Content - Some files uploaded successfully |
| 400 | Bad Request - Invalid input |
| 404 | Not Found - Resource not found |
| 413 | Payload Too Large - File too large |
| 415 | Unsupported Media Type - Invalid file type |
| 422 | Unprocessable Entity - Cannot process file |
| 500 | Internal Server Error - Server error |
| 503 | Service Unavailable - External service error |

---

## Dependencies (go.mod)

**หลัก:**
- `fiber/v2` - Web framework
- `gorm` - ORM
- `go-openai` - OpenAI SDK
- `websocket` - WebSocket support
- `postgres` - PostgreSQL driver

**File Processing:**
- `pdf` - PDF parsing
- `docx` - DOCX parsing
- `excelize` - Excel parsing
- `etree` - XML parsing

**อื่นๆ:**
- `godotenv` - Environment variables
- `uuid` - UUID generation

---

## การเปลี่ยนแปลงล่าสุด

### Version 4.0 (2025-11-03)
**Breaking Changes - Enhanced Persona System:**
- ✅ **ขยาย Persona schema ให้มีการตั้งค่าครบถ้วน**
- ✅ เพิ่มฟิลด์ใหม่:
  - `tone` - โทนเสียง (friendly, professional, empathetic)
  - `style` - สไตล์การตอบ (concise, detailed, conversational)
  - `temperature` - ค่า temperature (0.0-2.0)
  - `max_tokens` - จำนวน tokens สูงสุด
  - `model` - โมเดล AI ที่ใช้
  - `language_setting` (JSONB) - การตั้งค่าภาษาและสไตล์การตอบ
  - `guardrails` (JSONB) - กฎการกรองเนื้อหาและความปลอดภัย
  - `updated_at` - เวลาอัปเดตล่าสุด
- ✅ เพิ่ม struct `LanguageSetting` และ `Guardrails` สำหรับ configuration
- ✅ เปลี่ยน `name` เป็น UNIQUE constraint

**คุณสมบัติใหม่:**
- ✅ รองรับ content filtering ผ่าน Guardrails
- ✅ ตั้งค่าภาษาและสไตล์การตอบแบบละเอียด
- ✅ กำหนด AI model และ parameters แยกต่างหากแต่ละ persona
- ✅ จัดการหัวข้อที่อนุญาตและห้าม
- ✅ ควบคุมความยาวคำตอบสูงสุด

**ตัวอย่าง Persona Configuration:**
```json
{
  "name": "Professional Assistant",
  "tone": "professional",
  "style": "detailed",
  "temperature": 0.5,
  "model": "gpt-4",
  "language_setting": {
    "default_language": "en",
    "response_style": "formal",
    "language_code": "en-US"
  },
  "guardrails": {
    "block_profanity": true,
    "allowed_topics": ["business", "technology"],
    "max_response_length": 3000
  }
}
```

### Version 3.1 (2025-11-03)
**New Features:**
- ✅ เพิ่ม API `DELETE /api/file/uploads` สำหรับลบบันทึกไฟล์ทั้งหมด
- ✅ เพิ่ม `DeleteAll()` method ใน FileAnalysisRepository
- ✅ ใช้ GORM `AllowGlobalUpdate` เพื่อป้องกันการลบข้อมูลโดยไม่ตั้งใจ

**⚠️ หมายเหตุ:** การลบจะลบเฉพาะบันทึกใน DB ไฟล์บน disk จะยังคงอยู่

### Version 3.0 (2025-11-03)
**Breaking Changes - File Upload System:**
- ✅ **ลบการวิเคราะห์ด้วย AI ออกจาก File Upload API**
- ✅ เปลี่ยน endpoint จาก `/api/file/analyze` เป็น `/api/file/uploads`
- ✅ รองรับการอัปโหลดหลายไฟล์พร้อมกัน (สูงสุด 5 ไฟล์)
- ✅ บันทึกเฉพาะ metadata (file_name, storage_path, mime_type, file_size, uploaded_at)
- ✅ ลบฟิลด์ที่เกี่ยวกับ AI analysis ออกจาก FileAnalysis model
- ✅ อัปเดต repository methods ให้ใช้ `mime_type` และ `uploaded_at`

**คุณสมบัติใหม่:**
- ✅ Form field เปลี่ยนจาก `file` เป็น `files` (รองรับ multiple files)
- ✅ Response แสดงสถิติการอัปโหลด (success, failed, total)
- ✅ Partial upload support - ถ้าบางไฟล์ล้มเหลว ไฟล์อื่นยังอัปโหลดได้
- ✅ Auto-cleanup - ลบไฟล์จาก disk ถ้าบันทึก DB ล้มเหลว
- ✅ HTTP Status 206 (Partial Content) สำหรับการอัปโหลดบางไฟล์สำเร็จ

**ตัวอย่างการใช้งาน:**
```bash
# อัปโหลดไฟล์เดียว
curl -X POST http://localhost:3001/api/file/uploads \
  -F "files=@document.pdf"

# อัปโหลดหลายไฟล์
curl -X POST http://localhost:3001/api/file/uploads \
  -F "files=@doc1.pdf" \
  -F "files=@image.jpg" \
  -F "files=@data.xlsx"
```

### Version 2.1 (2025-11-03)
**New Features:**
- ✅ เพิ่ม API `DELETE /api/chats` สำหรับลบข้อความทั้งหมด
- ✅ เพิ่ม `DeleteAll()` method ใน MessageRepository
- ✅ ใช้ GORM `AllowGlobalUpdate` เพื่อป้องกันการลบข้อมูลโดยไม่ตั้งใจ

### Version 2.0 (2025-11-03)
**Breaking Changes:**
- ❌ **ลบไฟล์ `backend/controllers/helpers.go` ออกแล้ว**
- ✅ เปลี่ยนการจัดการ HTTP responses ทั้งหมดให้ใช้ `c.Status().JSON()` โดยตรง
- ✅ ย้ายฟังก์ชัน pagination validation เข้าไปในแต่ละ controller
- ✅ ย้ายการ convert messages เข้าไปใน controller โดยตรง

**ข้อดี:**
- โค้ดชัดเจนขึ้น อ่านง่ายขึ้น
- ลด abstraction layer ที่ไม่จำเป็น
- ลด coupling ระหว่าง controllers
- ง่ายต่อการ debug และ maintain

**ตัวอย่างการเปลี่ยนแปลง:**
```go
// เดิม (Version 1.0)
return successJSON(c, response)
return badRequest(c, "Invalid input")
return internalError(c, "Server error")

// ใหม่ (Version 2.0)
return c.Status(fiber.StatusOK).JSON(response)
return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
```

---

## สรุป

API นี้เป็น chatbot backend ที่สมบูรณ์ รองรับ:
- ✅ Chat แบบ non-streaming และ streaming (WebSocket)
- ✅ AI personas ที่หลากหลาย
- ✅ การวิเคราะห์ไฟล์หลายรูปแบบ
- ✅ Speech-to-Text และ Text-to-Speech
- ✅ Conversation history และ context management
- ✅ Session-based conversations
- ✅ File integration ใน conversations

Architecture ที่ใช้: **Repository-Service-Controller** เหมาะสำหรับการขยายและบำรุงรักษา

**Code Style:** Direct HTTP responses ด้วย Fiber framework (ไม่ใช้ helper functions)
