# 🤖 ChatBot Project

โปรเจกต์ AI Chatbot ที่สร้างด้วย **Vue.js 3**, **Go Fiber**, **PostgreSQL** และ **OpenAI API**

## 📋 สารบัญ

- [ฟีเจอร์หลัก](#ฟีเจอร์หลัก)
- [เทคโนโลยีที่ใช้](#เทคโนโลยีที่ใช้)
- [การติดตั้ง](#การติดตั้ง)
- [การใช้งาน Docker](#การใช้งาน-docker)
- [การตั้งค่า Environment Variables](#การตั้งค่า-environment-variables)
- [การรัน Backend และ Frontend](#การรัน-backend-และ-frontend)
- [API Endpoints](#api-endpoints)
- [Database Schema](#database-schema)

---

## ✨ ฟีเจอร์หลัก

- 💬 **Chat Interface**: ส่งข้อความและรับคำตอบจาก AI แบบ real-time
- 🎯 **Persona System**: เลือกบุคลิกของ AI (General, Technology, Business, etc.)
- 🔊 **Audio Transcription**: บันทึกเสียงและแปลงเป็นข้อความด้วย Whisper API
- 🌊 **WebSocket Streaming**: รับคำตอบแบบ streaming (typing effect)
- 💾 **Session Management**: บันทึกประวัติการสนทนา
- 🔐 **Authentication**: ระบบ login/register (coming soon)

---

## 🛠 เทคโนโลยีที่ใช้

### Frontend
- **Vue.js 3** - Progressive JavaScript Framework
- **Vite** - Fast build tool
- **Pinia** - State management
- **Axios** - HTTP client

### Backend
- **Go 1.21+** - Programming language
- **Fiber v2** - Web framework (Express-like for Go)
- **GORM** - ORM for database operations
- **OpenAI Go SDK** - API integration

### Database
- **PostgreSQL 15+** - Relational database
- **Redis** (optional) - Caching and session storage

### DevOps
- **Docker & Docker Compose** - Containerization
- **pgAdmin** - Database management UI

---

## 📦 การติดตั้ง

### Prerequisites

ติดตั้งโปรแกรมเหล่านี้ก่อน:

```bash
# ตรวจสอบ Docker
docker --version
docker-compose --version

# ตรวจสอบ Go (สำหรับ local development)
go version

# ตรวจสอบ Node.js (สำหรับ frontend)
node --version
npm --version
```

### Clone Repository

```bash
git clone <repository-url>
cd ChatBotProject
```

---

## 🐳 การใช้งาน Docker

### 1. สร้างไฟล์ `.env`

```bash
# คัดลอกจาก .env.example
cp .env.example .env

# แก้ไขไฟล์ .env และใส่ค่าที่ต้องการ
# โดยเฉพาะ OPENAI_API_KEY
```

### 2. เริ่มต้น PostgreSQL และ Services

```bash
# เริ่ม PostgreSQL, pgAdmin, และ Redis
docker-compose up -d

# ดู logs
docker-compose logs -f

# ตรวจสอบสถานะ
docker-compose ps
```

### 3. เข้าถึง Services

| Service | URL | Credentials |
|---------|-----|-------------|
| **PostgreSQL** | `localhost:5432` | User: `chatbot_user`<br>Password: `chatbot_password`<br>Database: `chatbot_db` |
| **pgAdmin** | [http://localhost:5050](http://localhost:5050) | Email: `admin@chatbot.local`<br>Password: `admin123` |
| **Redis** | `localhost:6379` | Password: `redis_password` |

### 4. หยุด Services

```bash
# หยุด containers
docker-compose stop

# หยุดและลบ containers (แต่ data ยังอยู่)
docker-compose down

# ลบทั้ง containers และ volumes (ข้อมูลจะหายหมด!)
docker-compose down -v
```

---

## 🔧 การตั้งค่า Environment Variables

แก้ไขไฟล์ `.env` ตามความต้องการ:

```bash
# Database
POSTGRES_USER=chatbot_user
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=chatbot_db
DATABASE_URL=postgres://chatbot_user:your_secure_password@localhost:5432/chatbot_db?sslmode=disable

# OpenAI (สำคัญ! ต้องมี API Key)
OPENAI_API_KEY=sk-your-actual-api-key-here
OPENAI_MODEL=gpt-3.5-turbo  # หรือ gpt-4

# Backend
PORT=3000
CORS_ORIGIN=http://localhost:5173

# JWT (ถ้าใช้ authentication)
JWT_SECRET=your-super-secret-key-min-32-chars
```

---

## 🚀 การรัน Backend และ Frontend

### Backend (Go Fiber)

```bash
cd backend

# ติดตั้ง dependencies
go mod download

# รัน development mode (with hot reload)
# ติดตั้ง air ก่อน: go install github.com/cosmtrek/air@latest
air

# หรือรันแบบปกติ
go run main.go
```

Backend จะรันที่: [http://localhost:3000](http://localhost:3000)

### Frontend (Vue.js)

```bash
cd frontend

# ติดตั้ง dependencies
npm install

# รัน development server
npm run dev
```

Frontend จะรันที่: [http://localhost:5173](http://localhost:5173)

---

## 📡 API Endpoints

### Health Check
```bash
GET /api/health
```

### Chat
```bash
# Send message (sync)
POST /api/chat
Body: {
  "session_id": "uuid",
  "message": "Hello",
  "persona_id": 1
}

# WebSocket streaming
WS /api/chat/stream
```

### Audio
```bash
POST /api/audio/transcribe
Content-Type: multipart/form-data
Body: { audio: File }
```

### Personas
```bash
GET /api/personas              # List all personas
GET /api/personas/:id          # Get persona details
```

### Sessions (coming soon)
```bash
GET    /api/sessions           # List user's sessions
POST   /api/sessions           # Create new session
PUT    /api/sessions/:id       # Update session
DELETE /api/sessions/:id       # Archive session
```

---

## 🗄 Database Schema

### Tables

#### `users`
- `id` (UUID, PK)
- `email` (VARCHAR, UNIQUE)
- `username` (VARCHAR, UNIQUE)
- `password_hash` (VARCHAR)
- `created_at`, `updated_at`, `last_login_at` (TIMESTAMP)
- `is_active` (BOOLEAN)

#### `sessions`
- `id` (UUID, PK)
- `user_id` (UUID, FK → users)
- `title` (VARCHAR)
- `persona_id` (INTEGER, FK → personas)
- `created_at`, `updated_at` (TIMESTAMP)
- `is_archived` (BOOLEAN)

#### `messages`
- `id` (UUID, PK)
- `session_id` (UUID, FK → sessions)
- `user_id` (UUID, FK → users)
- `role` (VARCHAR: 'user', 'assistant', 'system')
- `content` (TEXT)
- `tokens_used` (INTEGER)
- `metadata` (JSONB)
- `created_at` (TIMESTAMP)

#### `personas`
- `id` (SERIAL, PK)
- `name` (VARCHAR)
- `system_prompt` (TEXT)
- `expertise` (VARCHAR)
- `description` (TEXT)
- `icon` (VARCHAR)
- `is_active` (BOOLEAN)
- `created_at` (TIMESTAMP)

#### `audit_logs`
- `id` (SERIAL, PK)
- `user_id` (UUID, FK → users)
- `action`, `resource` (VARCHAR)
- `details` (JSONB)
- `ip_address` (INET)
- `created_at` (TIMESTAMP)

---

## 🔍 การตรวจสอบ Database

### ใช้ psql (Command Line)

```bash
# เข้าสู่ PostgreSQL container
docker exec -it chatbot_postgres psql -U chatbot_user -d chatbot_db

# คำสั่ง psql ที่มีประโยชน์
\dt                    # แสดงทุก tables
\d messages           # แสดง schema ของ messages table
SELECT * FROM personas;           # ดู personas ทั้งหมด
SELECT COUNT(*) FROM messages;    # นับจำนวน messages
```

### ใช้ pgAdmin (GUI)

1. เปิด [http://localhost:5050](http://localhost:5050)
2. Login: `admin@chatbot.local` / `admin123`
3. Add Server:
   - Name: `ChatBot DB`
   - Host: `postgres` (ชื่อ service ใน docker-compose)
   - Port: `5432`
   - Username: `chatbot_user`
   - Password: `chatbot_password`
   - Database: `chatbot_db`

---

## 🧪 การทดสอบ

### ทดสอบ API ด้วย curl

```bash
# Health check
curl http://localhost:3000/api/health

# List personas
curl http://localhost:3000/api/personas

# Send chat message
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "สวัสดีครับ",
    "persona_id": 1
  }'
```

---

## 📝 TODO List

- [ ] เพิ่ม JWT Authentication
- [ ] เพิ่ม Rate Limiting
- [ ] เพิ่ม Unit Tests
- [ ] เพิ่ม CI/CD Pipeline
- [ ] เพิ่ม Monitoring (Prometheus, Grafana)
- [ ] เพิ่ม Error Tracking (Sentry)
- [ ] เพิ่ม Docker Compose สำหรับ Backend + Frontend

---

## 🐛 Troubleshooting

### Database connection failed

```bash
# ตรวจสอบว่า PostgreSQL รันอยู่
docker-compose ps

# ดู logs
docker-compose logs postgres

# Restart service
docker-compose restart postgres
```

### OpenAI API error

```bash
# ตรวจสอบว่ามี OPENAI_API_KEY ใน .env
cat .env | grep OPENAI_API_KEY

# ตรวจสอบว่า API key ใช้งานได้
curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY"
```

### Port already in use

```bash
# หา process ที่ใช้ port 5432
# Windows
netstat -ano | findstr :5432
taskkill /PID <PID> /F

# Linux/Mac
lsof -i :5432
kill -9 <PID>
```

---

## 📄 License

MIT License

---

## 👥 Contributors

- Your Name - Initial work

---

## 📞 Support

หากมีปัญหาหรือคำถาม กรุณาเปิด [Issue](https://github.com/your-repo/issues) ในโปรเจกต์
