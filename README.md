# 🤖 ChatBot Project

โปรเจกต์ AI Chatbot ที่สร้างด้วย **Go Fiber**, **PostgreSQL**, **OpenAI API** และ **WebSocket Streaming**

---

## 📖 โปรเจ็คนี้คืออะไร

**ChatBot Project** เป็นระบบแชทบอท AI แบบ Full-Stack ที่ผสมผสานเทคโนโลยี:

- 🤖 **Multi-Persona AI Assistant** - ตัวช่วย AI หลายบุคลิก (General Assistant, Technology Expert, Business Advisor)
- 🎤 **Voice Input** - บันทึกเสียงและแปลงเป็นข้อความด้วย OpenAI Whisper API
- ⚡ **Real-time Streaming** - รับคำตอบแบบ streaming ผ่าน WebSocket (typing effect แบบ real-time)
- 💾 **Conversation History** - เก็บบันทึกการสนทนาทั้งหมดใน PostgreSQL
- 🚀 **High-Performance API** - Backend ที่สร้างด้วย Go Fiber รองรับการเชื่อมต่อพร้อมกันได้หลายพัน connections

**Use Cases:**
- ระบบ Customer Support แบบ AI
- Virtual Assistant สำหรับธุรกิจ
- Educational Chatbot
- Personal AI Assistant

---

## ✨ Feature มีอะไรบ้าง

### 🎯 Core Features (ที่ implement แล้ว)

#### 1. **Chat System**
- ✅ **Synchronous Chat** - ส่งข้อความและรับคำตอบแบบปกติ
- ✅ **WebSocket Streaming** - รับคำตอบแบบ real-time typing effect
- ✅ **Message History** - ดูประวัติการสนทนาพร้อม pagination
- ✅ **Context Awareness** - AI จำบริบทการสนทนา

#### 2. **Multi-Persona System**
- ✅ **3 Personas ในระบบ:**
  - 🤖 **General Assistant** - ผู้ช่วยทั่วไปตอบคำถามหลากหลาย
  - 💻 **Technology Expert** - ผู้เชี่ยวชาญด้านเทคโนโลยีและการเขียนโปรแกรม
  - 💼 **Business Advisor** - ที่ปรึกษาธุรกิจและการตลาด
- ✅ **Custom System Prompts** - แต่ละ persona มี behavior ต่างกัน
- ✅ **Persona Statistics** - ดูจำนวนข้อความและประสิทธิภาพของแต่ละ persona

#### 3. **Voice Input (Audio Transcription)**
- ✅ **Upload Audio Files** - รองรับไฟล์เสียงหลายฟอร์แมต (mp3, wav, webm, m4a, mp4)
- ✅ **OpenAI Whisper Integration** - แปลงเสียงเป็นข้อความด้วย AI
- ✅ **File Size Validation** - จำกัดขนาดไฟล์ไม่เกิน 25MB
- ✅ **Language Detection** - ตรวจจับภาษาที่พูดอัตโนมัติ

#### 4. **Real-time Streaming**
- ✅ **WebSocket Connection** - เชื่อมต่อแบบ bidirectional
- ✅ **Chunk-based Streaming** - รับข้อความทีละ chunk แบบ real-time
- ✅ **Typing Effect** - แสดงผลแบบพิมพ์ทีละคำเหมือนคนจริง
- ✅ **Auto-reconnect** - เชื่อมต่อใหม่อัตโนมัติเมื่อการเชื่อมต่อขาด

#### 5. **Database & Persistence**
- ✅ **PostgreSQL 15** - ฐานข้อมูลที่มั่นคงและรวดเร็ว
- ✅ **UUID Primary Keys** - รองรับระบบ distributed
- ✅ **JSONB Metadata** - เก็บข้อมูลเพิ่มเติมแบบ flexible
- ✅ **Auto-Migration** - สร้างตารางอัตโนมัติเมื่อเริ่มโปรแกรม
- ✅ **Connection Pooling** - จัดการ connections อย่างมีประสิทธิภาพ

#### 6. **API Documentation**
- ✅ **Comprehensive API Docs** - เอกสารครบถ้วนใน [API_DOCUMENTATION.md](API_DOCUMENTATION.md)
- ✅ **Request/Response Examples** - ตัวอย่างการใช้งานทุก endpoint
- ✅ **Error Handling** - จัดการ errors อย่างถูกต้อง

### 🚧 Features ที่วางแผนไว้ (Coming Soon)

- ⏳ **JWT Authentication** - ระบบ login/register
- ⏳ **User Management** - จัดการผู้ใช้งาน
- ⏳ **Rate Limiting** - จำกัดจำนวนการเรียก API
- ⏳ **Frontend UI (Vue.js 3)** - หน้าตาผู้ใช้งาน (โครงสร้างพร้อมแล้ว)
- ⏳ **Unit Tests** - ทดสอบระบบ
- ⏳ **CI/CD Pipeline** - Deploy อัตโนมัติ
- ⏳ **Monitoring** - ติดตามสถานะระบบ

---

## 🛠 เทคโนโลยีที่ใช้

### Backend Stack

| Technology | Version | Purpose |
|-----------|---------|---------|
| **Go** | 1.21+ | ภาษาโปรแกรมหลัก (high-performance, concurrent) |
| **Fiber** | v2.52.6 | Web framework (Express-like for Go) |
| **GORM** | v1.25.5 | ORM สำหรับจัดการฐานข้อมูล |
| **PostgreSQL Driver** | v1.5.4 | Database driver |
| **go-openai** | v1.41.2 | OpenAI API SDK (Chat + Whisper) |
| **Fiber WebSocket** | v1.3.4 | WebSocket support สำหรับ streaming |
| **google/uuid** | v1.6.0 | UUID generation |
| **godotenv** | v1.5.1 | Environment variables loader |

### Database

- **PostgreSQL 15+** - RDBMS พร้อม JSONB และ UUID extensions
- **uuid-ossp Extension** - สำหรับสร้าง UUID ในฐานข้อมูล
- **Connection Pooling** - Max 100 open connections, 10 idle

### External Services

- **OpenAI API**
  - gpt-4o-mini / GPT-4 - Chat completions (sync & streaming)
  - Whisper API - Audio transcription
- **Docker & Docker Compose** - Container orchestration

### Frontend (Placeholder)

- **Vue.js 3** - Progressive JavaScript framework
- **Vite** - Build tool
- **Pinia** - State management (planned)
- **Axios** - HTTP client (planned)

---

## 📦 การติดตั้ง

### ✅ Prerequisites

ตรวจสอบว่าติดตั้งโปรแกรมเหล่านี้แล้ว:

```bash
# ตรวจสอบ Docker
docker --version              # ควรเป็น 20.0 ขึ้นไป
docker-compose --version      # ควรเป็น 1.29 ขึ้นไป

# ตรวจสอบ Go (สำหรับ local development)
go version                    # ควรเป็น 1.21 ขึ้นไป
```

### 📥 Clone Repository

```bash
git clone <repository-url>
cd ChatBotProject
```

---

## 🐳 การใช้งาน Docker

### 1. ตั้งค่า Environment Variables

สร้างไฟล์ `.env.development` ในโฟลเดอร์ `backend/`:

```bash
cd backend

# คัดลอกจาก template (ถ้ามี)
cp .env.example .env.development

# หรือสร้างใหม่ด้วย text editor
nano .env.development
```

**ไฟล์ `.env.development` ต้องมี:**

```bash
# Server Configuration
PORT=3000
APP_ENV=development
APP_NAME=ChatBotAPI

# OpenAI API (⚠️ จำเป็นต้องมี!)
OPENAI_API_KEY=sk-proj-YOUR_ACTUAL_API_KEY_HERE
OPENAI_MODEL=gpt-4o-mini
OPENAI_MAX_TOKENS=2000
OPENAI_TEMPERATURE=0.7

# PostgreSQL Database
DATABASE_URL=postgres://chatbot_user:admin123@localhost:5432/chatbot_db?sslmode=disable
POSTGRES_USER=chatbot_user
POSTGRES_PASSWORD=admin123
POSTGRES_DB=chatbot_db
POSTGRES_PORT=5432

# CORS (สำหรับ Frontend)
CORS_ORIGIN=http://localhost:5173
```

⚠️ **สำคัญ!** ต้องใส่ OpenAI API Key ของคุณที่ `OPENAI_API_KEY`

### 2. เริ่มต้น PostgreSQL Database

```bash
# เริ่ม PostgreSQL container
cd backend
docker-compose up -d

# ตรวจสอบสถานะ
docker-compose ps

# ดู logs
docker-compose logs -f postgres
```

**ผลลัพธ์ที่ควรเห็น:**
```
✓ Container chatbot_postgres  Started
```

### 3. เข้าถึง PostgreSQL

| วิธี | URL/Command | Credentials |
|------|-------------|-------------|
| **psql (CLI)** | `docker exec -it chatbot_postgres psql -U chatbot_user -d chatbot_db` | User: `chatbot_user`<br>Password: `admin123` |
| **Connection String** | `postgres://chatbot_user:admin123@localhost:5432/chatbot_db` | - |
| **pgAdmin** | (ยังไม่ได้ติดตั้ง - ใช้ GUI tool อื่นได้) | - |

### 4. หยุด Services

```bash
# หยุด containers
docker-compose stop

# หยุดและลบ containers (แต่ data ยังอยู่)
docker-compose down

# ลบทั้ง containers และ volumes (⚠️ ข้อมูลจะหายหมด!)
docker-compose down -v
```

---

## 🚀 วิธีรันโปรเจ็ค

### 🔹 Option 1: รัน Backend อย่างเดียว (แนะนำสำหรับเริ่มต้น)

#### **Step 1: เริ่ม PostgreSQL**

```bash
cd backend
docker-compose up -d
```

#### **Step 2: ติดตั้ง Go Dependencies**

```bash
# อยู่ในโฟลเดอร์ backend
go mod download
go mod tidy
```

#### **Step 3: รัน Backend Server**

```bash
# รันแบบปกติ
go run main.go

# หรือรันพร้อม hot-reload (ติดตั้ง air ก่อน)
# go install github.com/cosmtrek/air@latest
air
```

**ผลลัพธ์ที่ควรเห็น:**

```
2025/10/29 11:27:31 ✓ Configuration loaded successfully (env: development)
2025/10/29 11:27:31 ✓ Database connected successfully
2025/10/29 11:27:31 ✓ PostgreSQL extensions initialized
2025/10/29 11:27:31 ✓ Database migration completed
2025/10/29 11:27:31 ✓ Personas already exist, skipping seed
2025/10/29 11:27:31 🚀 Server starting on port 3000 (env: development)

 ┌───────────────────────────────────────────────────┐
 │                    ChatBotAPI                     │
 │                   Fiber v2.52.6                   │
 │               http://127.0.0.1:3000               │
 │       (bound on host 0.0.0.0 and port 3000)       │
 │                                                   │
 │ Handlers ............ 16  Processes ........... 1 │
 │ Prefork ....... Disabled  PID ............. 10040 │
 └───────────────────────────────────────────────────┘
```

Backend จะรันที่: **http://localhost:3000**

#### **Step 4: ทดสอบ API**

```bash
# Health check
curl http://localhost:3000/api/health

# ผลลัพธ์:
# {"env":"development","message":"ChatBot API is running","status":"ok"}
```

### 🔹 Option 2: รัน Backend + Frontend (เมื่อ Frontend พร้อม)

```bash
# Terminal 1: Backend
cd backend
go run main.go

# Terminal 2: Frontend (ยังไม่พร้อมใช้งาน)
cd frontend
npm install
npm run dev
```

---

## 📡 API Endpoints

### 🏥 Health Check

```bash
GET /api/health

# Response:
{
  "status": "ok",
  "message": "ChatBot API is running",
  "env": "development"
}
```

### 💬 Chat Endpoints

#### 1. Send Message (Synchronous)

```bash
POST /api/chat
Content-Type: application/json

{
  "message": "สวัสดีครับ แนะนำตัวหน่อย",
  "persona_id": 1,
  "temperature": 0.7,
  "max_tokens": 1000
}

# Response:
{
  "message_id": "550e8400-e29b-41d4-a716-446655440000",
  "reply": "สวัสดีครับ! ผม...",
  "persona": {
    "id": 1,
    "name": "General Assistant",
    "expertise": "general",
    "icon": "🤖"
  },
  "tokens_used": 156,
  "model": "gpt-4o-mini",
  "timestamp": "2025-10-29T10:30:00Z"
}
```

#### 2. WebSocket Streaming

```javascript
// เชื่อมต่อ WebSocket
const ws = new WebSocket('ws://localhost:3000/api/chat/stream')

// ส่งข้อความ
ws.onopen = () => {
  ws.send(JSON.stringify({
    type: 'message',
    content: 'เล่านิทานให้ฟังหน่อย',
    persona_id: 1
  }))
}

// รับ chunks
ws.onmessage = (event) => {
  const data = JSON.parse(event.data)

  if (data.type === 'chunk' && !data.done) {
    // รับ chunk ของข้อความ
    console.log(data.content) // "กาล", "ครั้ง", "หนึ่ง", ...
  } else if (data.done) {
    // Streaming เสร็จสิ้น
    console.log('Message ID:', data.message_id)
    console.log('Tokens used:', data.tokens_used)
  }
}
```

#### 3. Get Chat History

```bash
GET /api/chat/history?limit=50&offset=0

# Response:
{
  "messages": [
    {
      "id": "uuid",
      "role": "user",
      "content": "สวัสดี",
      "persona_id": 1,
      "created_at": "2025-10-29T10:30:00Z"
    },
    {
      "id": "uuid",
      "role": "assistant",
      "content": "สวัสดีครับ...",
      "persona_id": 1,
      "created_at": "2025-10-29T10:30:05Z"
    }
  ],
  "total": 156,
  "limit": 50,
  "offset": 0
}
```

### 🎭 Persona Endpoints

#### 1. List All Personas

```bash
GET /api/personas

# Response:
{
  "personas": [
    {
      "id": 1,
      "name": "General Assistant",
      "system_prompt": "You are a helpful, friendly...",
      "expertise": "general",
      "description": "A versatile AI assistant...",
      "icon": "🤖",
      "is_active": true,
      "created_at": "2025-10-29T09:23:50Z"
    }
  ]
}
```

#### 2. Get Persona Details

```bash
GET /api/personas/1

# Response:
{
  "id": 1,
  "name": "General Assistant",
  "system_prompt": "You are...",
  "expertise": "general",
  "description": "A versatile AI...",
  "icon": "🤖",
  "is_active": true,
  "created_at": "2025-10-29T09:23:50Z",
  "stats": {
    "total_messages": 256,
    "avg_response_time": "2.3s"
  }
}
```

### 🎤 Audio Transcription

```bash
POST /api/audio/transcribe
Content-Type: multipart/form-data

# Form data:
audio: [audio file] (max 25MB)

# Response:
{
  "text": "สวัสดีครับ ผมต้องการถามเกี่ยวกับ...",
  "language": "th",
  "duration": 3.5,
  "confidence": 0.95,
  "timestamp": "2025-10-29T10:30:00Z"
}
```

**รองรับไฟล์:** mp3, wav, webm, m4a, mp4

---

## 🧪 วิธีทดสอบ WebSocket

### 🔹 Method 1: ใช้ HTML Test Page (แนะนำ)

โปรเจกต์มี test file ให้แล้วที่ `test_websocket.html`

#### **Step 1: เปิด Backend Server**

```bash
cd backend
go run main.go
# รอจนขึ้น "Server starting on port 3000"
```

#### **Step 2: เปิด Test Page**

```bash
# Windows
start test_websocket.html

# macOS
open test_websocket.html

# Linux
xdg-open test_websocket.html

# หรือเปิดไฟล์ด้วย browser โดยตรง
```

#### **Step 3: ทดสอบ**

1. ✅ เมื่อเปิดหน้าเว็บจะเชื่อมต่อ WebSocket อัตโนมัติ
2. ✅ สถานะจะเปลี่ยนจาก "○ Disconnected" → "● Connected"
3. พิมพ์ข้อความ เช่น **"เล่านิทานให้ฟังหน่อย"**
4. เลือก Persona (1: General, 2: Tech, 3: Business)
5. กดปุ่ม **"ส่ง"**
6. 🎉 คุณจะเห็น AI ตอบกลับแบบ **typing effect** (ทีละคำ real-time)

**ตัวอย่างการใช้งาน:**

```
You: เล่านิทานให้ฟังหน่อย
AI:  กาล┃                    <- chunk 1
AI:  กาลครั้ง┃               <- chunk 2
AI:  กาลครั้งหนึ่ง┃          <- chunk 3
AI:  กาลครั้งหนึ่ง มี...┃    <- chunk 4
...
```

### 🔹 Method 2: ใช้ JavaScript Console

เปิด Browser Console (F12) แล้ววาง code นี้:

```javascript
// เชื่อมต่อ WebSocket
const ws = new WebSocket('ws://localhost:3000/api/chat/stream');

ws.onopen = () => {
  console.log('✅ Connected!');

  // ส่งข้อความ
  ws.send(JSON.stringify({
    type: 'message',
    content: 'สวัสดีครับ บอทช่วยอะไรได้บ้าง',
    persona_id: 1
  }));
};

let fullResponse = '';

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);

  if (data.type === 'chunk') {
    if (data.done) {
      console.log('✅ Done!');
      console.log('Full response:', fullResponse);
      console.log('Message ID:', data.message_id);
      console.log('Tokens used:', data.tokens_used);
    } else {
      fullResponse += data.content;
      console.log('Chunk:', data.content);
    }
  }
};

ws.onerror = (error) => {
  console.error('❌ Error:', error);
};

ws.onclose = () => {
  console.log('🔌 Disconnected');
};
```

### 🔹 Method 3: ใช้ Postman/Thunder Client

1. สร้าง WebSocket request
2. URL: `ws://localhost:3000/api/chat/stream`
3. เชื่อมต่อ
4. ส่งข้อความ JSON:
   ```json
   {
     "type": "message",
     "content": "Hello",
     "persona_id": 1
   }
   ```
5. ดู response streams

---

## 🗄 Database Schema

### Tables ที่ใช้งานจริง

#### `personas` - บุคลิกของ AI

```sql
CREATE TABLE personas (
  id           SERIAL PRIMARY KEY,
  name         VARCHAR(255) NOT NULL,
  system_prompt TEXT NOT NULL,
  expertise    VARCHAR(100),
  description  TEXT,
  icon         VARCHAR(10),
  is_active    BOOLEAN DEFAULT true,
  created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Default Data:**
- General Assistant (🤖)
- Technology Expert (💻)
- Business Advisor (💼)

#### `messages` - ข้อความทั้งหมด

```sql
CREATE TABLE messages (
  id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  role         VARCHAR(20) NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
  content      TEXT NOT NULL,
  persona_id   INTEGER REFERENCES personas(id),
  tokens_used  INTEGER,
  metadata     JSONB DEFAULT '{}',
  created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes:**
- `idx_messages_created_at` - สำหรับ pagination
- `idx_messages_persona_id` - สำหรับ filter by persona

---

## 🔍 การตรวจสอบ Database

### ใช้ psql (Command Line)

```bash
# เข้าสู่ PostgreSQL container
docker exec -it chatbot_postgres psql -U chatbot_user -d chatbot_db

# คำสั่งที่มีประโยชน์
\dt                              # แสดงทุก tables
\d messages                      # แสดง schema ของ messages table
\d personas                      # แสดง schema ของ personas table

SELECT * FROM personas;                     # ดู personas ทั้งหมด
SELECT COUNT(*) FROM messages;              # นับจำนวน messages
SELECT * FROM messages ORDER BY created_at DESC LIMIT 10;  # ข้อความล่าสุด
SELECT role, COUNT(*) FROM messages GROUP BY role;  # นับตาม role
```

### ใช้ DBeaver / TablePlus / pgAdmin

**Connection Info:**
- Host: `localhost`
- Port: `5432`
- Database: `chatbot_db`
- Username: `chatbot_user`
- Password: `admin123`

---

## 🧪 การทดสอบด้วย cURL

### ทดสอบ Health Check

```bash
curl http://localhost:3000/api/health
```

### ทดสอบ List Personas

```bash
curl http://localhost:3000/api/personas
```

### ทดสอบ Send Chat Message

```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "สวัสดีครับ แนะนำตัวหน่อย",
    "persona_id": 1
  }'
```

### ทดสอบ Get Chat History

```bash
curl "http://localhost:3000/api/chat/history?limit=10&offset=0"
```

### ทดสอบ Audio Transcription

```bash
curl -X POST http://localhost:3000/api/audio/transcribe \
  -F "audio=@path/to/your/audio.mp3"
```

---

## 📝 Project Structure

```
ChatBotProject/
├── backend/                     # Go Fiber Backend
│   ├── main.go                 # Entry point
│   ├── go.mod                  # Dependencies
│   ├── .env.development        # Configuration
│   ├── docker-compose.yml      # PostgreSQL setup
│   ├── config/                 # Configuration management
│   │   ├── config.go          # Env loader
│   │   └── database.go        # DB connection
│   ├── models/                 # Database models (GORM)
│   │   ├── message.go         # Message entity
│   │   └── persona.go         # Persona entity
│   ├── repositories/           # Data access layer
│   │   ├── message_repository.go
│   │   └── persona_repository.go
│   ├── services/               # Business logic
│   │   └── openai_service.go  # OpenAI integration
│   ├── controllers/            # HTTP handlers
│   │   ├── chat_controller.go
│   │   ├── persona_controller.go
│   │   ├── audio_controller.go
│   │   └── websocket_controller.go
│   └── routes/
│       └── routes.go          # Route definitions
├── frontend/                   # Vue.js (placeholder)
│   └── .env.development
├── test_websocket.html        # WebSocket test page
├── API_DOCUMENTATION.md       # API documentation
├── BMAD.md                    # Architecture doc
├── WORKFLOW.md                # Development workflow
└── README.md                  # This file
```

---

## 🐛 Troubleshooting

### ❌ Database Connection Failed

**Symptoms:**
```
Failed to connect to database: dial tcp [::1]:5432: connectex: No connection could be made
```

**Solution:**
```bash
# ตรวจสอบว่า PostgreSQL รันอยู่
docker-compose ps

# ดู logs
docker-compose logs postgres

# Restart service
docker-compose restart postgres

# ตรวจสอบ connection string ใน .env.development
cat .env.development | grep DATABASE_URL
```

### ❌ OpenAI API Error

**Symptoms:**
```
Failed to get AI response: 401 Unauthorized
```

**Solution:**
```bash
# ตรวจสอบว่ามี OPENAI_API_KEY ใน .env
cat .env.development | grep OPENAI_API_KEY

# ทดสอบ API key
curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### ❌ Port Already in Use

**Symptoms:**
```
Failed to start server: listen tcp :3000: bind: Only one usage of each socket address
```

**Solution:**
```bash
# Windows
netstat -ano | findstr :3000
taskkill /F /PID <PID>

# Linux/Mac
lsof -i :3000
kill -9 <PID>
```

### ❌ WebSocket Connection Failed

**Symptoms:**
- หน้าเว็บแสดง "○ Disconnected"
- Console มี error "WebSocket connection failed"

**Solution:**
1. ตรวจสอบว่า Backend รันอยู่
2. ตรวจสอบ URL ใน test_websocket.html ว่าตรงกับ port ที่ใช้
3. ตรวจสอบ firewall ไม่ได้ block port

---

## 📚 Additional Documentation

- 📖 [API Documentation](API_DOCUMENTATION.md) - API reference แบบละเอียด
- 🏗 [Architecture Document](BMAD.md) - สถาปัตยกรรมระบบ
- 📋 [Development Workflow](WORKFLOW.md) - ขั้นตอนการพัฒนา

---

## 📄 License

MIT License

---

## 👥 Contributors

- **Your Name** - Initial development

---

## 📞 Support

หากมีปัญหาหรือคำถาม:
- 📧 Email: your.email@example.com
- 🐛 [Open an Issue](https://github.com/your-repo/issues)
- 💬 Discord: your-discord-server

---

## 🎯 Quick Start Checklist

ใช้ checklist นี้เพื่อเริ่มต้นใช้งาน:

- [ ] ✅ ติดตั้ง Docker และ Go แล้ว
- [ ] ✅ Clone repository
- [ ] ✅ สร้างไฟล์ `.env.development` พร้อม OpenAI API Key
- [ ] ✅ รัน `docker-compose up -d` (PostgreSQL)
- [ ] ✅ รัน `go run main.go` (Backend)
- [ ] ✅ เปิด http://localhost:3000/api/health เห็น `"status":"ok"`
- [ ] ✅ ทดสอบ WebSocket ด้วย test_websocket.html
- [ ] 🎉 พร้อมใช้งาน!

---

**สร้างด้วย ❤️ และ Go + OpenAI**
