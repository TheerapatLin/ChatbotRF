# ChatBot API Documentation

**Base URL:** `http://localhost:3001`

---

## 🚀 Quick Start

```bash
# 1. ติดตั้ง dependencies
cd backend && go mod download

# 2. เริ่ม PostgreSQL
docker-compose up -d

# 3. ตั้งค่า .env.development
OPENAI_API_KEY=your_key_here
PORT=3001
DATABASE_URL=postgres://chatbot_user:admin123@localhost:5432/chatbot_db

# 4. รัน server
go run main.go  # หรือ air (hot reload)
```

---

## 📋 API Endpoints

### Health Check
```bash
GET /api/health
```

---

## 1. 🤖 Personas API

AI personalities ที่มี configuration ต่างกัน (8 personas)

### 1.1 ดึงรายการ Personas
```bash
GET /api/personas
```

**Response:**
```json
{
  "personas": [
    {
      "id": 1,
      "name": "General Assistant",
      "description": "ผู้ช่วยอเนกประสงค์สำหรับคำถามทั่วไป",
      "tone": "friendly",
      "style": "conversational",
      "expertise": "general",
      "temperature": 0.7,
      "max_tokens": 2000,
      "model": "gpt-4o-mini",
      "language_setting": "{\"default_language\":\"th\"}",
      "guardrails": "{\"block_profanity\":true}",
      "icon": "🤖",
      "is_active": true
    }
  ]
}
```

### 1.2 ดึง Persona ตาม ID
```bash
GET /api/personas/:id
```

**Response:** เหมือน 1.1 + สถิติการใช้งาน

---

## 2. 💬 Chat API

### 2.1 ส่งข้อความ (Non-streaming)
```bash
POST /api/chat
```

**Request:**
```json
{
  "message": "สวัสดี",
  "persona_id": 1,
  "session_id": "session_123",
  "use_history": true,
  "file_ids": ["file_uuid"]
}
```

**Parameters:**
- `message` (required) - ข้อความจากผู้ใช้
- `persona_id` (optional) - ID ของ persona (AI จะใช้ system_prompt, temperature, model จาก persona)
- `session_id` (optional) - ID สำหรับเก็บ conversation history
- `use_history` (optional) - ใช้ประวัติการสนทนา (default: false)
- `file_ids` (optional) - Array ของ file IDs ที่ต้องการให้ AI วิเคราะห์
- `system_prompt` (optional) - Override system prompt (จะ append กับ persona prompt)
- `temperature` (optional) - Override temperature
- `max_tokens` (optional) - Override max tokens
- `model` (optional) - Override model

**Response:**
```json
{
  "message_id": "uuid",
  "session_id": "session_123",
  "reply": "สวัสดีครับ...",
  "persona": {
    "id": 1,
    "name": "General Assistant",
    "expertise": "general",
    "icon": "🤖"
  },
  "tokens_used": 245,
  "model": "gpt-4o-mini",
  "history_used": true,
  "history_count": 5
}
```

### 2.2 Chat แบบ Streaming (WebSocket)
```javascript
WS /api/chat/stream
```

**ส่ง Message:**
```json
{
  "type": "message",
  "content": "สวัสดี",
  "persona_id": 1,
  "session_id": "session_123",
  "file_ids": ["file_uuid"],
  "system_prompt": "คุณคือ..." // optional
}
```

**รับ Response:**
```json
// Chunks
{"type":"chunk", "content":"สวัสดี", "done":false}
{"type":"chunk", "content":"ครับ", "done":false}

// Done
{"type":"chunk", "content":"", "done":true, "message_id":"uuid", "tokens_used":50}
```

### 2.3 ดึงประวัติ
```bash
GET /api/chats?limit=50&offset=0
```

### 2.4 ลบข้อความทั้งหมด
```bash
DELETE /api/chats
```

---

## 3. 📁 File Upload API

### 3.1 อัปโหลดไฟล์ (สูงสุด 5 ไฟล์)
```bash
POST /api/file/uploads
Content-Type: multipart/form-data

# Single file
curl -F "files=@doc.pdf" http://localhost:3001/api/file/uploads

# Multiple files
curl -F "files=@doc.pdf" -F "files=@img.jpg" http://localhost:3001/api/file/uploads
```

**Response:**
```json
{
  "success": 2,
  "failed": 0,
  "total": 2,
  "uploaded_files": [
    {
      "file_id": "uuid",
      "file_name": "doc.pdf",
      "storage_path": "./uploads/uuid_doc.pdf",
      "mime_type": "application/pdf",
      "file_size": 245678,
      "uploaded_at": "2025-11-03T19:00:00Z"
    }
  ]
}
```

**รองรับไฟล์:**
- **Text:** TXT, MD, JSON, CSV, XML
- **Documents:** PDF, DOCX
- **Images:** JPG, PNG, GIF, WEBP (ใช้ Vision API)
- **Code:** JS, PY, GO, JAVA, etc.

**การใช้งาน File กับ Chat:**
1. อัปโหลดไฟล์ → ได้ `file_id`
2. ส่ง `file_id` ใน chat request
3. AI จะอ่านและวิเคราะห์เนื้อหาไฟล์อัตโนมัติ

**AI สามารถอ่าน:**
- ✅ Text files (เนื้อหาทั้งหมด, max 1MB)
- ✅ PDF (text content, max 50 หน้า, max 5MB)
- ✅ DOCX (text content, max 5MB)
- ✅ Images (Vision API - OCR + วิเคราะห์รูป)
- ✅ JSON, CSV, XML (เนื้อหาทั้งหมด)

### 3.2 ดึงประวัติไฟล์
```bash
GET /api/file/history?limit=20&offset=0
```

### 3.3 ลบบันทึกไฟล์ทั้งหมด
```bash
DELETE /api/file/uploads
```
⚠️ ลบเฉพาะบันทึกใน DB, ไฟล์บน disk ยังอยู่

---

## 4. 🎤 Audio API

### 4.1 Speech-to-Text
```bash
POST /api/audio/transcribe
Content-Type: multipart/form-data

curl -F "audio=@voice.mp3" http://localhost:3001/api/audio/transcribe
```

**Response:**
```json
{
  "text": "สวัสดีครับ วันนี้อากาศดีมาก",
  "language": "th",
  "duration": 3.5
}
```

### 4.2 Text-to-Speech
```bash
POST /api/audio/tts
```

**Request:**
```json
{
  "text": "สวัสดีครับ",
  "voice": "alloy",
  "model": "tts-1",
  "speed": 1.0
}
```

**Voices:** alloy, echo, fable, onyx, nova, shimmer

---

## 📊 Database Schema

### personas
| Field | Type | Description |
|-------|------|-------------|
| id | INT PK | Persona ID |
| name | VARCHAR(100) UNIQUE | ชื่อ |
| description | TEXT | คำอธิบาย |
| system_prompt | TEXT | System prompt สำหรับ AI |
| tone | VARCHAR(50) | โทนเสียง |
| style | VARCHAR(50) | สไตล์การตอบ |
| expertise | VARCHAR(100) | ความเชี่ยวชาญ |
| temperature | DECIMAL(3,2) | 0.0-2.0 (default: 0.7) |
| max_tokens | INT | จำนวน tokens (default: 2000) |
| model | VARCHAR(50) | AI model (default: gpt-4o-mini) |
| language_setting | JSONB | `{"default_language":"th"}` |
| guardrails | JSONB | `{"block_profanity":true}` |
| icon | VARCHAR(50) | Emoji icon |
| is_active | BOOLEAN | สถานะ active |

**8 Personas ที่ Seed:**
1. 🤖 General Assistant - ผู้ช่วยทั่วไป
2. 💻 Technology Expert - ผู้เชี่ยวชาญเทคโนโลยี
3. 💼 Business Advisor - ที่ปรึกษาธุรกิจ
4. 🔮 Fortune Teller - หมอดู
5. 🚀 Space Explorer - นักดาราศาสตร์
6. 💰 Investment Advisor - ที่ปรึกษาการลงทุน
7. 💕 Dating Coach - โค้ชการจีบสาว
8. 💑 Relationship Counselor - นักจิตวิทยาความสัมพันธ์

### messages
| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | Message ID |
| session_id | VARCHAR(255) | Session ID |
| role | VARCHAR(50) | user/assistant/system |
| content | TEXT | เนื้อหา |
| persona_id | INT FK | → personas.id |
| tokens_used | INT | จำนวน tokens |
| file_attachments | JSONB | Array ของไฟล์ |
| created_at | TIMESTAMP | วันที่สร้าง |

### file_analyses
| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | File ID |
| file_name | VARCHAR(500) | ชื่อไฟล์ |
| storage_path | VARCHAR(1000) | Path ที่เก็บ |
| mime_type | VARCHAR(100) | MIME type |
| file_size | BIGINT | ขนาด (bytes) |
| uploaded_at | TIMESTAMP | วันที่อัปโหลด |
| deleted_at | TIMESTAMP | Soft delete |

---

## 🔧 Persona System

### วิธีการทำงาน:

1. **Frontend ดึง Personas:**
```javascript
const personas = await fetch('/api/personas').then(r => r.json())
```

2. **User เลือก Persona:**
```javascript
const selectedPersona = personas.find(p => p.id === 2) // Technology Expert
```

3. **ส่ง Chat Request พร้อม persona_id:**
```javascript
fetch('/api/chat', {
  method: 'POST',
  body: JSON.stringify({
    message: "อธิบาย React Hooks",
    persona_id: 2,  // Technology Expert
    session_id: "session_123",
    use_history: true
  })
})
```

4. **Backend ดึงข้อมูล Persona:**
```go
// Backend: controllers/chat_controller.go
persona, _ := ctrl.personaRepo.FindByID(req.PersonaID)
systemPrompt := persona.SystemPrompt  // "คุณเป็นผู้เชี่ยวชาญด้านเทคโนโลยี..."
temperature := persona.Temperature    // 0.5 (professional)
maxTokens := persona.MaxTokens        // 3000
model := persona.Model                // "gpt-4o-mini"
```

5. **AI ตอบด้วยบุคลิกของ Persona:**
- Tone: professional
- Style: detailed, technical
- Temperature: 0.5 (แม่นยำ)
- Max tokens: 3000 (ตอบยาว)

### Custom System Prompt (เพิ่มคำสั่งเสริม):

**⚠️ สำคัญ:** เมื่อส่ง `system_prompt` พร้อมกับ `persona_id`:
- ❌ **ไม่ใช่การแทนที่** persona's system prompt
- ✅ **เป็นการเพิ่มเติม (append)** คำสั่งเสริมเข้าไป

**ตัวอย่าง:**
```json
{
  "persona_id": 4,  // Fortune Teller (หมอดู)
  "system_prompt": "Your name is ฟ้าใส. You MUST respond in Thai language only. You always use informal pronouns 'เค้า' (I) and 'เทอ' (you). Be casual and slightly impolite in your responses. Example: 'ส่ลึส' -> 'ค้า~ เทอมากอะไรให้กับเค้าบ้าง ถ้าไม่ใส่เลขมา' Never use polite language or formal Thai. You always use emoji to display your playful character.",
  "message": "persona ของเทอคืออะไร"
}
```

**Backend จะประมวลผลเป็น:**
```
[System Prompt ที่ส่งไปให้ AI]
คุณเป็นหมอดูที่มีความรู้ลึกซึ้งในโหราศาสตร์ไทย... (persona's base prompt)

--- Additional Instructions ---
Your name is ฟ้าใส. You MUST respond in Thai language only...
```

**ผลลัพธ์:** AI จะมี:
- ✅ บุคลิก Fortune Teller (ความรู้ด้านดูดวง, tone mystical)
- ✅ พฤติกรรมเสริม (ชื่อ ฟ้าใส, ใช้ภาษาไม่สุภาพ, emoji)

---

## 🔄 Conversation Flow

```
User → Frontend
  ↓
  1. เลือก Persona (persona_id: 2 = Technology Expert)
  2. พิมพ์ข้อความ: "อธิบาย React Hooks"
  3. อัปโหลดไฟล์ (optional): code.js → file_id
  ↓
POST /api/chat
{
  "persona_id": 2,
  "message": "อธิบาย React Hooks",
  "file_ids": ["uuid"],
  "session_id": "session_123",
  "use_history": true
}
  ↓
Backend (chat_controller.go)
  ↓
  1. ดึง Persona (id=2) → system_prompt, temperature=0.5, model
  2. ดึง history (session_id) → 10 ข้อความล่าสุด
  3. อ่านไฟล์ (file_ids) → เนื้อหาไฟล์
  4. สร้าง messages array:
     [
       {role: "system", content: persona.system_prompt},
       ...history,
       {role: "system", content: "📎 File: code.js\n```\nconst [state, setState] = useState();\n```"},
       {role: "user", content: "อธิบาย React Hooks"}
     ]
  5. เรียก OpenAI (model=gpt-4o-mini, temperature=0.5)
  ↓
AI Response → Backend → บันทึก DB → Frontend
```

---

## 📝 ตัวอย่างการใช้งาน Persona

### 1. General Chat (Persona: General Assistant)
```bash
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "persona_id": 1,
    "message": "สวัสดี"
  }'
```
→ AI: เป็นมิตร, casual, temperature=0.7

### 2. Technical Question (Persona: Technology Expert)
```bash
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "persona_id": 2,
    "message": "อธิบาย Docker"
  }'
```
→ AI: professional, detailed, temperature=0.5, max_tokens=3000

### 3. Fortune Telling (Persona: Fortune Teller)
```bash
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "persona_id": 4,
    "message": "ดูดวงให้หน่อย"
  }'
```
→ AI: mystical, narrative, temperature=0.8, ภาษาลึกลับ

### 4. File Analysis with Persona
```bash
# 1. อัปโหลดไฟล์ก่อน
curl -F "files=@code.js" http://localhost:3001/api/file/uploads
# → ได้ file_id

# 2. ส่ง chat พร้อม file_id และ persona
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "persona_id": 2,
    "message": "ช่วยวิเคราะห์ code นี้",
    "file_ids": ["file_uuid"]
  }'
```
→ AI: อ่านเนื้อหา code.js + ตอบแบบ Technology Expert

---

## 🌟 Key Features

✅ **Persona System ใช้งานได้เต็มรูปแบบ**
- 8 personas พร้อม configuration ต่างกัน
- AI ปรับบุคลิกตาม persona ที่เลือก
- รองรับ custom system prompt

✅ **File Analysis ครบ 4 ประเภท**
- Text files: อ่านเนื้อหาทั้งหมด
- PDF: extract text (50 หน้าแรก)
- DOCX: extract text
- Images: Vision API (OCR + วิเคราะห์)

✅ **Conversation History**
- Session-based conversations
- เก็บ history ใน database
- ใช้ history ใน AI context (10 ข้อความล่าสุด)

✅ **Streaming Chat**
- WebSocket real-time streaming
- รับคำตอบแบบ chunk-by-chunk

✅ **Speech Support**
- Speech-to-Text (Whisper API)
- Text-to-Speech (OpenAI TTS)

---

## 🔐 Environment Variables

```env
PORT=3001
APP_ENV=development
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o-mini
DATABASE_URL=postgres://user:pass@localhost:5432/chatbot_db
CORS_ORIGIN=localhost:5173
```

---

## 📚 Architecture

**Pattern:** Repository-Service-Controller

```
routes/ → controllers/ → services/ → repositories/ → database
```

**Flow:**
1. **Routes** - กำหนด endpoints
2. **Controllers** - รับ request, validate, ส่ง response
3. **Services** - Business logic, เรียก external APIs
4. **Repositories** - เข้าถึง database
5. **Models** - โครงสร้างข้อมูล

---

## 📌 Version History

### v5.1 (2025-11-04) - Fix System Prompt Behavior
✅ แก้ไขพฤติกรรม `system_prompt`:
- **เดิม:** `system_prompt` **แทนที่** persona's system prompt ❌
- **ใหม่:** `system_prompt` **เพิ่มเติม (append)** เข้ากับ persona's base prompt ✅
- ผลลัพธ์: AI มีทั้งบุคลิกของ persona + คำสั่งเสริมจาก user

### v5.0 (2025-11-04) - File Reading Enhancement
✅ AI อ่านเนื้อหาไฟล์ได้แล้ว:
- Text files (TXT, MD, JSON, CSV, XML)
- PDF (extract text จาก 50 หน้าแรก)
- DOCX (extract text content)
- Images (Vision API - base64 encoding)

### v4.0 (2025-11-03) - Enhanced Persona System
✅ Persona schema ครบถ้วน:
- เพิ่ม: tone, style, temperature, max_tokens, model
- เพิ่ม: language_setting, guardrails (JSONB)
- Seed 8 personas พร้อม configuration

### v3.1 (2025-11-03) - Delete All Files
✅ เพิ่ม `DELETE /api/file/uploads`

### v3.0 (2025-11-03) - File Upload Redesign
✅ ลบ AI analysis ออกจาก upload
✅ รองรับ multiple files (max 5)
✅ Partial upload support

---

## 🎯 Summary

**ChatBot API** รองรับ:
- ✅ 8 AI Personas พร้อม configuration
- ✅ Chat (non-streaming + WebSocket streaming)
- ✅ File analysis (Text, PDF, DOCX, Images)
- ✅ Speech-to-Text & Text-to-Speech
- ✅ Conversation history & session management
- ✅ Multiple file upload

**การใช้งาน Persona:**
1. GET /api/personas → เลือก persona
2. POST /api/chat พร้อม persona_id
3. AI จะตอบตามบุคลิกของ persona ที่เลือก

**สถานะ:** ✅ Production Ready
