# 🚀 Development Workflow - ChatBot Project

> **ขั้นตอนการพัฒนาโปรเจกต์ ChatBot แบบละเอียด**
>
> สำหรับ Developer ทุกระดับ - มีทั้ง Checklist, Timeline และ Commands ที่ใช้

---

## 📅 Timeline Overview (ภาพรวมระยะเวลา)

| Phase | ระยะเวลา | สถานะ | หัวข้อหลัก |
|-------|---------|-------|------------|
| **Phase 1** | Week 1 (5 วัน) | ✅ **Complete** | Setup & Database Connection |
| **Phase 2** | Week 2-3 (10 วัน) | ⏳ In Progress | Backend API Development |
| **Phase 3** | Week 3-4 (10 วัน) | ⏹️ Pending | Frontend Development |
| **Phase 4** | Week 5 (5 วัน) | ⏹️ Pending | Testing & Deployment |

**รวมทั้งหมด**: 5 สัปดาห์ (25 วันทำงาน)

---

## 📊 Progress Tracking

```
[████████░░░░░░░░░░░░░░░░░░] 30% Complete

✅ Phase 1: Setup                    (100%)
⏳ Phase 2: Backend Development      (20%)
⏹️ Phase 3: Frontend Development     (0%)
⏹️ Phase 4: Testing & Deployment     (0%)
```

---

## 🎯 Phase 1: Setup & Database Connection

**ระยะเวลา**: 5 วันทำงาน (Week 1)
**สถานะ**: ✅ **Complete**
**ผู้รับผิดชอบ**: Backend Developer

### Day 1: Project Initialization

#### Task 1.1: ติดตั้ง Development Environment
**เวลา**: 2-3 ชั่วโมง

- [ ] ติดตั้ง Docker Desktop
- [ ] ติดตั้ง Go 1.21+
- [ ] ติดตั้ง Node.js 18+
- [ ] ติดตั้ง PostgreSQL client (psql หรือ pgAdmin)
- [ ] ติดตั้ง VS Code + Extensions:
  - Go extension
  - Vetur (Vue)
  - PostgreSQL extension

**Commands:**
```bash
# ตรวจสอบ versions
docker --version
go version
node --version
psql --version
```

#### Task 1.2: Clone และ Setup Project Structure
**เวลา**: 1-2 ชั่วโมง

- [x] สร้างโครงสร้างโฟลเดอร์
- [x] สร้าง `.gitignore`
- [x] สร้าง `README.md`
- [x] สร้าง `.env.example`

**Commands:**
```bash
cd ChatBotProject

# สร้างโครงสร้าง
mkdir -p backend/{config,models,repositories,services,controllers,middleware,routes,utils}
mkdir -p frontend/src/{components,views,services,store,composables,assets}
mkdir -p database/init

# Initialize Git
git init
git add .
git commit -m "Initial project structure"
```

---

### Day 2: Database Setup

#### Task 2.1: Setup PostgreSQL with Docker
**เวลา**: 1-2 ชั่วโมง

- [x] สร้าง `docker-compose.yml`
- [x] เริ่ม PostgreSQL container
- [x] ตรวจสอบ connection
- [x] เพิ่ม UUID extension

**Commands:**
```bash
# เริ่ม PostgreSQL
cd backend
docker-compose up -d postgres

# ตรวจสอบสถานะ
docker-compose ps

# เข้าสู่ PostgreSQL
docker exec -it chatbot_postgres psql -U chatbot_user -d chatbot_db

# เพิ่ม extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
\q
```

**Verification:**
```bash
# ทดสอบ connection
docker exec chatbot_postgres psql -U chatbot_user -d chatbot_db -c "SELECT 1;"
```

#### Task 2.2: สร้าง Database Schema
**เวลา**: 2-3 ชั่วโมง

- [x] สร้าง `database/init/01_init_extensions.sql`
- [x] สร้าง `database/init/02_create_tables.sql`
- [x] สร้าง `database/init/03_seed_data.sql`
- [ ] ทดสอบ schema

**Verification:**
```bash
# ตรวจสอบ tables
docker exec chatbot_postgres psql -U chatbot_user -d chatbot_db -c "\dt"

# ตรวจสอบ personas
docker exec chatbot_postgres psql -U chatbot_user -d chatbot_db -c "SELECT * FROM personas;"
```

---

### Day 3: Backend Configuration

#### Task 3.1: Setup Go Project
**เวลา**: 2-3 ชั่วโมง

- [x] สร้าง `go.mod`
- [x] ติดตั้ง dependencies
- [x] สร้าง `config/config.go` - โหลด .env
- [x] สร้าง `config/database.go` - Database connection

**Commands:**
```bash
cd backend

# Initialize Go module
go mod init chatbot

# ติดตั้ง dependencies
go get github.com/gofiber/fiber/v2
go get github.com/sashabaranov/go-openai
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/joho/godotenv
go get github.com/google/uuid
go get gorm.io/datatypes

# Tidy dependencies
go mod tidy
```

#### Task 3.2: สร้าง Models
**เวลา**: 1-2 ชั่วโมง

- [x] สร้าง `models/message.go`
- [x] สร้าง `models/persona.go`
- [ ] ทดสอบ models

**Verification:**
```bash
# Compile test
go build -o /dev/null .
```

---

### Day 4: Database Connection & Migration

#### Task 4.1: สร้าง Main Application
**เวลา**: 2-3 ชั่วโมง

- [x] สร้าง `main.go`
- [x] Setup Fiber app
- [x] เชื่อมต่อ Database
- [x] Auto-migrate models
- [x] Seed personas

**Commands:**
```bash
# รัน application
go run main.go

# ควรเห็น output:
# ✓ Configuration loaded successfully
# ✓ Database connected successfully
# ✓ Database migration completed
# ✓ Seeded 3 personas
# 🚀 Server starting on port 3000
```

#### Task 4.2: ทดสอบ Connection
**เวลา**: 1 ชั่วโมง

- [x] ทดสอบ health endpoint
- [x] ตรวจสอบ database
- [x] ตรวจสอบ logs

**Verification:**
```bash
# ทดสอบ health endpoint
curl http://localhost:3000/api/health

# Expected: {"status":"ok","message":"ChatBot API is running"}

# ตรวจสอบ database
docker exec chatbot_postgres psql -U chatbot_user -d chatbot_db -c "SELECT COUNT(*) FROM messages;"
docker exec chatbot_postgres psql -U chatbot_user -d chatbot_db -c "SELECT COUNT(*) FROM personas;"
```

---

### Day 5: Documentation & Review

#### Task 5.1: เขียน Documentation
**เวลา**: 2-3 ชั่วโมง

- [x] อัพเดท README.md
- [x] สร้าง BMAD.md
- [x] สร้าง WORKFLOW.md (ไฟล์นี้)
- [ ] สร้าง API documentation

#### Task 5.2: Code Review & Cleanup
**เวลา**: 1-2 ชั่วโมง

- [ ] Review code quality
- [ ] ตรวจสอบ error handling
- [ ] ลบ unused code
- [ ] Format code

**Commands:**
```bash
# Format Go code
go fmt ./...

# Run linter (ถ้าติดตั้ง golangci-lint)
golangci-lint run

# Git commit
git add .
git commit -m "feat: Phase 1 complete - Database setup and connection"
```

---

## 🔧 Phase 2: Backend API Development

**ระยะเวลา**: 10 วันทำงาน (Week 2-3)
**สถานะ**: ⏳ **In Progress (20%)**
**ผู้รับผิดชอบ**: Backend Developer

### Week 2: Core API Development

#### Day 6-7: Repositories & Services Layer

**Day 6 Tasks** (6-8 ชั่วโมง):

##### Task 6.1: สร้าง Repositories
- [ ] สร้าง `repositories/message_repository.go`
  - [ ] `Create(message)` - บันทึกข้อความ
  - [ ] `FindByID(id)` - หาข้อความจาก ID
  - [ ] `FindRecent(limit)` - ดึงข้อความล่าสุด N ข้อความ
  - [ ] `Delete(id)` - ลบข้อความ

- [ ] สร้าง `repositories/persona_repository.go`
  - [ ] `FindAll()` - ดึง personas ทั้งหมด
  - [ ] `FindByID(id)` - หา persona จาก ID
  - [ ] `FindActive()` - ดึงเฉพาะ active personas

**Code Template:**
```go
// repositories/message_repository.go
package repositories

import (
    "chatbot/models"
    "gorm.io/gorm"
)

type MessageRepository struct {
    db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
    return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
    return r.db.Create(message).Error
}

func (r *MessageRepository) FindRecent(limit int) ([]models.Message, error) {
    var messages []models.Message
    err := r.db.Order("created_at DESC").Limit(limit).Find(&messages).Error
    return messages, err
}
```

**Day 7 Tasks** (6-8 ชั่วโมง):

##### Task 7.1: สร้าง OpenAI Service
- [ ] สร้าง `services/openai_service.go`
  - [ ] Setup OpenAI client
  - [ ] `SendChatRequest(messages)` - ส่งข้อความและรับ response
  - [ ] `BuildContext(recentMessages)` - สร้าง context array
  - [ ] Error handling และ retry logic

**Code Template:**
```go
// services/openai_service.go
package services

import (
    "context"
    "chatbot/config"
    "github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
    client *openai.Client
    config *config.Config
}

func NewOpenAIService(cfg *config.Config) *OpenAIService {
    client := openai.NewClient(cfg.OpenAIAPIKey)
    return &OpenAIService{
        client: client,
        config: cfg,
    }
}

func (s *OpenAIService) SendChatRequest(messages []openai.ChatCompletionMessage) (string, error) {
    resp, err := s.client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model:    s.config.OpenAIModel,
            Messages: messages,
        },
    )
    if err != nil {
        return "", err
    }
    return resp.Choices[0].Message.Content, nil
}
```

##### Task 7.2: ทดสอบ OpenAI Integration
- [ ] เขียน test function
- [ ] ทดสอบ API call
- [ ] ตรวจสอบ response

**Verification:**
```bash
# สร้าง test file
# services/openai_service_test.go
go test ./services -v
```

---

#### Day 8-9: Controllers & Routes

**Day 8 Tasks** (6-8 ชั่วโมง):

##### Task 8.1: สร้าง Chat Controller
- [ ] สร้าง `controllers/chat_controller.go`
  - [ ] `HandleChat(c *fiber.Ctx)` - รับข้อความและตอบกลับ
  - [ ] Validate input
  - [ ] Build context
  - [ ] Call OpenAI
  - [ ] Save to database
  - [ ] Return response

**Code Template:**
```go
// controllers/chat_controller.go
package controllers

import (
    "chatbot/models"
    "chatbot/repositories"
    "chatbot/services"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

type ChatController struct {
    messageRepo   *repositories.MessageRepository
    personaRepo   *repositories.PersonaRepository
    openaiService *services.OpenAIService
}

func NewChatController(
    msgRepo *repositories.MessageRepository,
    persRepo *repositories.PersonaRepository,
    aiSvc *services.OpenAIService,
) *ChatController {
    return &ChatController{
        messageRepo:   msgRepo,
        personaRepo:   persRepo,
        openaiService: aiSvc,
    }
}

type ChatRequest struct {
    Message   string `json:"message"`
    PersonaID *int   `json:"persona_id"`
}

type ChatResponse struct {
    MessageID uuid.UUID `json:"message_id"`
    Reply     string    `json:"reply"`
    Timestamp string    `json:"timestamp"`
}

func (ctrl *ChatController) HandleChat(c *fiber.Ctx) error {
    var req ChatRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    // TODO: Implement chat logic
    // 1. Get persona (if specified)
    // 2. Build context from recent messages
    // 3. Call OpenAI API
    // 4. Save user message and AI response
    // 5. Return response

    return c.JSON(ChatResponse{
        MessageID: uuid.New(),
        Reply:     "TODO: Implement chat logic",
        Timestamp: time.Now().Format(time.RFC3339),
    })
}
```

**Day 9 Tasks** (6-8 ชั่วโมง):

##### Task 9.1: สร้าง Persona Controller
- [ ] สร้าง `controllers/persona_controller.go`
  - [ ] `GetAllPersonas(c *fiber.Ctx)` - ดึง personas ทั้งหมด
  - [ ] `GetPersonaByID(c *fiber.Ctx)` - ดึง persona ตาม ID

##### Task 9.2: Setup Routes
- [ ] อัพเดท `routes/routes.go`
  - [ ] เพิ่ม `/api/chat` endpoint
  - [ ] เพิ่ม `/api/personas` endpoints
  - [ ] Setup dependency injection

**Code Template:**
```go
// routes/routes.go
package routes

import (
    "chatbot/config"
    "chatbot/controllers"
    "chatbot/repositories"
    "chatbot/services"
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
    // Initialize repositories
    messageRepo := repositories.NewMessageRepository(db)
    personaRepo := repositories.NewPersonaRepository(db)

    // Initialize services
    openaiService := services.NewOpenAIService(cfg)

    // Initialize controllers
    chatCtrl := controllers.NewChatController(messageRepo, personaRepo, openaiService)
    personaCtrl := controllers.NewPersonaController(personaRepo)

    // API group
    api := app.Group("/api")

    // Health check
    api.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })

    // Chat routes
    api.Post("/chat", chatCtrl.HandleChat)

    // Persona routes
    api.Get("/personas", personaCtrl.GetAllPersonas)
    api.Get("/personas/:id", personaCtrl.GetPersonaByID)
}
```

---

#### Day 10: Testing & Integration

**Day 10 Tasks** (6-8 ชั่วโมง):

##### Task 10.1: Integration Testing
- [ ] ทดสอบ `/api/chat` endpoint
  - [ ] ส่งข้อความธรรมดา
  - [ ] ส่งข้อความกับ persona
  - [ ] ทดสอบ error cases

- [ ] ทดสอบ `/api/personas` endpoints
  - [ ] ดึง personas ทั้งหมด
  - [ ] ดึง persona ตาม ID
  - [ ] ทดสอบ ID ที่ไม่มี

**Test Commands:**
```bash
# ทดสอบ chat endpoint
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"สวัสดี","persona_id":1}'

# ทดสอบ personas endpoint
curl http://localhost:3000/api/personas

# ทดสอบ persona by ID
curl http://localhost:3000/api/personas/1
```

##### Task 10.2: Performance Testing
- [ ] ทดสอบ response time
- [ ] ทดสอบ concurrent requests
- [ ] Monitor database queries

**Verification:**
```bash
# Load testing (ถ้ามี hey installed)
hey -n 100 -c 10 http://localhost:3000/api/health
```

---

### Week 3: Advanced Features

#### Day 11-12: WebSocket Streaming

**Day 11 Tasks** (6-8 ชั่วโมง):

##### Task 11.1: Setup WebSocket
- [ ] สร้าง WebSocket endpoint `/api/chat/stream`
- [ ] Handle WebSocket connection
- [ ] Implement message handling

**Code Template:**
```go
// routes/routes.go
import "github.com/gofiber/websocket/v2"

func SetupRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
    // ... existing code ...

    // WebSocket route
    api.Get("/chat/stream", websocket.New(func(c *websocket.Conn) {
        // Handle WebSocket connection
        for {
            msgType, msg, err := c.ReadMessage()
            if err != nil {
                break
            }
            // TODO: Process message and stream response
            c.WriteMessage(msgType, []byte("Echo: "+string(msg)))
        }
    }))
}
```

**Day 12 Tasks** (6-8 ชั่วโมง):

##### Task 12.1: Implement Streaming Response
- [ ] เรียก OpenAI Streaming API
- [ ] ส่ง chunks กลับไป client
- [ ] Handle errors และ disconnections

##### Task 12.2: ทดสอบ WebSocket
- [ ] ทดสอบ connection
- [ ] ทดสอบ streaming
- [ ] ทดสอบ reconnection

---

#### Day 13-14: Audio Transcription

**Day 13 Tasks** (6-8 ชั่วโมง):

##### Task 13.1: สร้าง Audio Controller
- [ ] สร้าง `controllers/audio_controller.go`
- [ ] รับ multipart file upload
- [ ] Validate file type และ size

**Day 14 Tasks** (6-8 ชั่วโมง):

##### Task 14.1: สร้าง Whisper Service
- [ ] สร้าง `services/whisper_service.go`
- [ ] เรียก Whisper API
- [ ] Return transcribed text

##### Task 14.2: Integration
- [ ] เพิ่ม `/api/audio/transcribe` route
- [ ] ทดสอบ audio upload
- [ ] ทดสอบ transcription

**Test Command:**
```bash
# ทดสอบ audio transcription
curl -X POST http://localhost:3000/api/audio/transcribe \
  -F "audio=@test.mp3"
```

---

#### Day 15: Code Review & Optimization

**Day 15 Tasks** (6-8 ชั่วโมง):

##### Task 15.1: Code Review
- [ ] Review error handling
- [ ] Review logging
- [ ] Check for memory leaks
- [ ] Optimize database queries

##### Task 15.2: Documentation
- [ ] เขียน API documentation
- [ ] Update README
- [ ] เขียน inline comments

**Commands:**
```bash
# Format code
go fmt ./...

# Run tests
go test ./... -v

# Git commit
git add .
git commit -m "feat: Phase 2 complete - Backend API development"
```

---

## 🎨 Phase 3: Frontend Development

**ระยะเวลา**: 10 วันทำงาน (Week 3-4)
**สถานะ**: ⏹️ **Pending**
**ผู้รับผิดชอบ**: Frontend Developer

### Week 3: Vue.js Setup & Core Components

#### Day 16-17: Project Setup

**Day 16 Tasks** (6-8 ชั่วโมง):

##### Task 16.1: Initialize Vue Project
- [ ] สร้าง Vue 3 project ด้วย Vite
- [ ] ติดตั้ง dependencies
- [ ] Setup project structure

**Commands:**
```bash
cd frontend

# Create Vue project
npm create vue@latest

# Options:
# - Project name: chatbot-frontend
# - TypeScript: No
# - JSX: No
# - Vue Router: Yes
# - Pinia: Yes
# - Vitest: No
# - ESLint: Yes

# Install dependencies
npm install axios

# Run dev server
npm run dev
```

##### Task 16.2: Setup API Service
- [ ] สร้าง `services/api.js` - Axios instance
- [ ] สร้าง `services/chatService.js`
- [ ] สร้าง `services/personaService.js`

**Code Template:**
```javascript
// services/api.js
import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

export default api
```

**Day 17 Tasks** (6-8 ชั่วโมง):

##### Task 17.1: Setup Pinia Store
- [ ] สร้าง `store/chat.js`
- [ ] สร้าง `store/personas.js`
- [ ] Define state, getters, actions

**Code Template:**
```javascript
// store/chat.js
import { defineStore } from 'pinia'
import chatService from '@/services/chatService'

export const useChatStore = defineStore('chat', {
  state: () => ({
    messages: [],
    isLoading: false,
    error: null
  }),

  actions: {
    async sendMessage(message, personaId) {
      this.isLoading = true
      try {
        const response = await chatService.sendMessage(message, personaId)
        this.messages.push({
          role: 'user',
          content: message
        })
        this.messages.push({
          role: 'assistant',
          content: response.reply
        })
      } catch (error) {
        this.error = error.message
      } finally {
        this.isLoading = false
      }
    }
  }
})
```

---

#### Day 18-20: Core Components

**Day 18 Tasks** (6-8 ชั่วโมง):

##### Task 18.1: สร้าง Chat Components
- [ ] สร้าง `components/ChatInterface.vue` - Main chat container
- [ ] สร้าง `components/MessageList.vue` - แสดงข้อความ
- [ ] Basic styling

**Day 19 Tasks** (6-8 ชั่วโมง):

##### Task 19.1: สร้าง Input Components
- [ ] สร้าง `components/MessageInput.vue` - ช่องพิมพ์ข้อความ
- [ ] เพิ่ม Send button
- [ ] Handle Enter key

##### Task 19.2: สร้าง Persona Selector
- [ ] สร้าง `components/PersonaSelector.vue`
- [ ] Fetch personas from API
- [ ] Display persona options

**Day 20 Tasks** (6-8 ชั่วโมง):

##### Task 20.1: Integration
- [ ] เชื่อมต่อ components เข้าด้วยกัน
- [ ] ทดสอบการส่งข้อความ
- [ ] ทดสอบการเลือก persona

---

### Week 4: Advanced Features & Polish

#### Day 21-22: WebSocket Integration

**Day 21-22 Tasks** (12-16 ชั่วโมง):

##### Task 21.1: สร้าง WebSocket Composable
- [ ] สร้าง `composables/useWebSocket.js`
- [ ] Handle connection, reconnection
- [ ] Handle incoming messages

**Code Template:**
```javascript
// composables/useWebSocket.js
import { ref, onMounted, onUnmounted } from 'vue'

export function useWebSocket() {
  const ws = ref(null)
  const isConnected = ref(false)
  const messages = ref([])

  const connect = () => {
    ws.value = new WebSocket('ws://localhost:3000/api/chat/stream')

    ws.value.onopen = () => {
      isConnected.value = true
    }

    ws.value.onmessage = (event) => {
      const data = JSON.parse(event.data)
      messages.value.push(data)
    }

    ws.value.onclose = () => {
      isConnected.value = false
      // Reconnect after 3 seconds
      setTimeout(connect, 3000)
    }
  }

  const sendMessage = (message) => {
    if (ws.value && isConnected.value) {
      ws.value.send(JSON.stringify(message))
    }
  }

  onMounted(connect)
  onUnmounted(() => ws.value?.close())

  return { isConnected, messages, sendMessage }
}
```

##### Task 21.2: Implement Streaming Display
- [ ] แสดงข้อความแบบ typing effect
- [ ] Handle streaming chunks

---

#### Day 23-24: Audio Recording

**Day 23-24 Tasks** (12-16 ชั่วโมง):

##### Task 23.1: สร้าง Audio Recorder
- [ ] สร้าง `components/AudioRecorder.vue`
- [ ] สร้าง `composables/useAudioRecorder.js`
- [ ] ใช้ MediaRecorder API

##### Task 23.2: Integration
- [ ] Upload audio file
- [ ] แสดง transcribed text
- [ ] ส่งเป็นข้อความ

---

#### Day 25: UI/UX Polish

**Day 25 Tasks** (6-8 ชั่วโมง):

##### Task 25.1: Styling & Animations
- [ ] เพิ่ม CSS transitions
- [ ] ปรับปรุง responsive design
- [ ] เพิ่ม loading states

##### Task 25.2: Error Handling
- [ ] แสดง error messages
- [ ] Retry mechanisms
- [ ] Offline detection

**Commands:**
```bash
# Build for production
npm run build

# Preview production build
npm run preview

# Git commit
git add .
git commit -m "feat: Phase 3 complete - Frontend development"
```

---

## 🧪 Phase 4: Testing & Deployment

**ระยะเวลา**: 5 วันทำงาน (Week 5)
**สถานะ**: ⏹️ **Pending**
**ผู้รับผิดชอบ**: Full Team

### Day 26-27: Testing

#### Day 26: Backend Testing

**Day 26 Tasks** (6-8 ชั่วโมง):

##### Task 26.1: Unit Tests
- [ ] เขียน tests สำหรับ repositories
- [ ] เขียน tests สำหรับ services
- [ ] เขียน tests สำหรับ controllers

**Commands:**
```bash
cd backend

# Run tests
go test ./... -v -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

##### Task 26.2: Integration Tests
- [ ] ทดสอบ API endpoints ทั้งหมด
- [ ] ทดสอบ error scenarios
- [ ] ทดสอบ database transactions

#### Day 27: Frontend Testing

**Day 27 Tasks** (6-8 ชั่วโมง):

##### Task 27.1: Component Tests
- [ ] ทดสอบ components แต่ละตัว
- [ ] ทดสอบ user interactions
- [ ] ทดสอบ API calls

##### Task 27.2: E2E Tests
- [ ] ทดสอบ user flows ทั้งหมด
- [ ] ทดสอบ cross-browser

---

### Day 28-29: Bug Fixes & Optimization

**Day 28-29 Tasks** (12-16 ชั่วโมง):

##### Task 28.1: Bug Fixes
- [ ] แก้ไข bugs ที่พบจากการ testing
- [ ] Refactor code ที่จำเป็น
- [ ] Update documentation

##### Task 28.2: Performance Optimization
- [ ] Optimize database queries
- [ ] Optimize bundle size
- [ ] Optimize API response time

**Performance Checklist:**
- [ ] Database query time < 100ms
- [ ] API response time < 500ms
- [ ] Frontend bundle size < 500KB
- [ ] Time to Interactive < 3s

---

### Day 30: Deployment

**Day 30 Tasks** (6-8 ชั่วโมง):

##### Task 30.1: Prepare for Deployment
- [ ] Setup environment variables
- [ ] Configure production database
- [ ] Setup SSL certificates

##### Task 30.2: Deploy
- [ ] Deploy database (PostgreSQL)
- [ ] Deploy backend (Go Fiber)
- [ ] Deploy frontend (Static hosting)

**Deployment Checklist:**
- [ ] Environment variables configured
- [ ] Database migrated
- [ ] SSL enabled
- [ ] Monitoring setup
- [ ] Backup configured

**Commands:**
```bash
# Build backend
cd backend
go build -o chatbot-api

# Build frontend
cd frontend
npm run build

# Deploy (example with Docker)
docker-compose -f docker-compose.prod.yml up -d
```

##### Task 30.3: Post-Deployment
- [ ] ทดสอบ production environment
- [ ] Monitor logs
- [ ] Setup alerts

---

## 📋 Daily Checklist Template

สำหรับทุกวันทำงาน:

### เริ่มวัน (Morning)
- [ ] Pull latest code: `git pull`
- [ ] Review TODO list
- [ ] เริ่ม development server
- [ ] Check dependencies

### ระหว่างวัน (During)
- [ ] Commit code บ่อยๆ
- [ ] เขียน meaningful commit messages
- [ ] ทดสอบ code ก่อน push
- [ ] Update documentation

### สิ้นวัน (Evening)
- [ ] Push code: `git push`
- [ ] Update progress tracking
- [ ] Plan tomorrow's tasks
- [ ] Backup important files

---

## 🎯 Definition of Done (DoD)

สำหรับแต่ละ task ต้อง:

- [ ] Code เขียนเสร็จและทำงานได้
- [ ] Tests passed ทั้งหมด
- [ ] Code review completed
- [ ] Documentation updated
- [ ] No critical bugs
- [ ] Performance acceptable
- [ ] Git committed and pushed

---

## 🚨 Emergency Procedures

### ถ้าเจอปัญหา:

#### Database Connection Issues
```bash
# Restart PostgreSQL
docker-compose restart postgres

# Check logs
docker-compose logs postgres

# Reset database (CAREFUL!)
docker-compose down -v
docker-compose up -d postgres
```

#### Backend Crashes
```bash
# Check logs
go run main.go 2>&1 | tee error.log

# Debug
dlv debug
```

#### Frontend Build Errors
```bash
# Clear cache
rm -rf node_modules
rm package-lock.json
npm install

# Clear Vite cache
rm -rf .vite
```

---

## 📊 Progress Tracking

### How to Track Progress

1. **Daily Standups** (15 minutes)
   - Yesterday: ทำอะไรไปบ้าง
   - Today: จะทำอะไร
   - Blockers: มีปัญหาอะไร

2. **Weekly Reviews** (1 hour)
   - Review completed tasks
   - Adjust timeline if needed
   - Plan next week

3. **Tools**
   - GitHub Issues
   - Trello/Jira
   - This WORKFLOW.md

---

## ✅ Phase Completion Criteria

### Phase 1 ✅
- [x] PostgreSQL running
- [x] Database connected
- [x] Models created
- [x] Health endpoint working

### Phase 2 ⏳
- [x] Repositories implemented (20%)
- [ ] OpenAI integration (0%)
- [ ] WebSocket streaming (0%)
- [ ] Audio transcription (0%)

### Phase 3 ⏹️
- [ ] Vue project setup
- [ ] Core components
- [ ] API integration
- [ ] UI/UX polish

### Phase 4 ⏹️
- [ ] Tests written
- [ ] Bugs fixed
- [ ] Performance optimized
- [ ] Deployed

---

## 🎓 Learning Resources

### Go & Fiber
- [Go by Example](https://gobyexample.com/)
- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Guide](https://gorm.io/docs/)

### Vue.js
- [Vue 3 Documentation](https://vuejs.org/)
- [Pinia Documentation](https://pinia.vuejs.org/)
- [Vite Guide](https://vitejs.dev/guide/)

### PostgreSQL
- [PostgreSQL Tutorial](https://www.postgresqltutorial.com/)
- [JSONB Guide](https://www.postgresql.org/docs/current/datatype-json.html)

### OpenAI
- [OpenAI API Documentation](https://platform.openai.com/docs)
- [Whisper API Guide](https://platform.openai.com/docs/guides/speech-to-text)

---

**📅 Created**: 2025-10-28
**👨‍💻 Last Updated**: 2025-10-28
**📍 Current Phase**: Phase 2 (Backend Development)
**⏱️ Estimated Completion**: Week 5 (5 สัปดาห์)
