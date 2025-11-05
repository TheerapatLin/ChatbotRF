# Frontend Development Workflow

**Tech Stack:** Vue 3 + Axios + WebSocket
**Backend API:** http://localhost:3001
**Target Port:** http://localhost:5173

---

## Task 1: Setup Project & Dependencies ✅ COMPLETED

**วัตถุประสงค์:** ติดตั้ง Vue 3 project และ dependencies ที่จำเป็น

### สิ่งที่ทำเสร็จแล้ว:

✅ **Project Setup**
- สร้าง Vue 3 project ด้วย Vite
- ติดตั้ง dependencies: axios, pinia
- Setup Pinia ใน main.js

✅ **API Service Layer** (`src/api/`)
- `axios.js` - Axios instance config พร้อม interceptors
- `personaService.js` - CRUD operations สำหรับ Personas
- `chatService.js` - Chat operations
- `fileService.js` - File upload/management
- `audioService.js` - Speech-to-Text & Text-to-Speech

✅ **Pinia Stores** (`src/stores/`)
- `chat.js` - จัดการ messages, WebSocket connection, streaming
- `persona.js` - จัดการ personas (CRUD, selection)
- `audio.js` - จัดการ audio recording, transcription, TTS
- `file.js` - จัดการ file uploads

✅ **WebSocket Service** (`src/services/`)
- `websocketService.js` - WebSocket connection manager พร้อม auto-reconnect

✅ **Configuration**
- `vite.config.js` - Proxy setup, path alias (@)
- `.env`, `.env.development`, `.env.production` - Environment variables
- `README.md` - Documentation

### โครงสร้างโปรเจค:
```
frontend/
├── src/
│   ├── api/
│   │   ├── axios.js
│   │   ├── personaService.js
│   │   ├── chatService.js
│   │   ├── fileService.js
│   │   └── audioService.js
│   ├── stores/
│   │   ├── chat.js
│   │   ├── persona.js
│   │   ├── audio.js
│   │   └── file.js
│   ├── services/
│   │   └── websocketService.js
│   └── main.js
├── .env
├── .env.development
├── .env.production
└── vite.config.js
```

### วิธีรัน:
```bash
cd frontend
npm run dev
```

**Status:** ✅ Task 1 เสร็จสมบูรณ์ - พร้อมสำหรับ Task 2

---

## Task 2: Create Persona Management Components ✅ COMPLETED

**วัตถุประสงค์:** สร้าง Sidebar สำหรับจัดการ Personas

### สิ่งที่ทำเสร็จแล้ว:

✅ **Common UI Components** (`src/components/common/`)
- `BaseButton.vue` - ปุ่มพร้อม variants, sizes, loading state
- `BaseModal.vue` - Modal dialog พร้อม transitions
- `BaseInput.vue` - Input field พร้อม validation
- `BaseSelect.vue` - Dropdown select
- `BaseTextarea.vue` - Textarea พร้อม auto-resize

✅ **PersonaSidebar.vue** (`src/components/persona/`)
**Features:**
- ✅ Dropdown เลือก Persona (ใช้ personaStore)
- ✅ แสดงรายละเอียด Persona ที่เลือก (Icon, Name, Description, Model, Tone, Temperature, MaxTokens, SystemPrompt)
- ✅ ปุ่ม "New" สร้าง Persona ใหม่ (เปิด Modal)
- ✅ ปุ่ม "Edit" แก้ไข Persona
- ✅ ปุ่ม "Delete" ลบ Persona (พร้อม confirmation)
- ✅ Loading state ขณะโหลดข้อมูล
- ✅ Error handling และ retry
- ✅ Auto-load personas เมื่อ mount
- ✅ Responsive design

✅ **PersonaModal.vue** (`src/components/persona/`)
**Features:**
- ✅ Form สร้าง/แก้ไข Persona (2 modes: create, edit)
- ✅ ฟิลด์ Basic Info: Name, Icon, Description, System Prompt
- ✅ ฟิลด์ Personality: Tone, Style, Expertise
- ✅ ฟิลด์ Model Config: Model (dropdown), Temperature (0-2), MaxTokens
- ✅ Language Settings: Default Language, Response Style, Language Code
- ✅ Guardrails: BlockProfanity, BlockSensitive, AllowedTopics, BlockedTopics, MaxResponseLength, RequireModeration
- ✅ Form validation
- ✅ Loading state ขณะ submit
- ✅ Auto-populate data เมื่อ edit mode
- ✅ Model options: gpt-4o-mini ⭐, gpt-4o, gpt-4-turbo, gpt-4, gpt-3.5-turbo

✅ **DebugPanel.vue** (`src/components/persona/`)
**Features:**
- ✅ แสดง Stats: Message Count, Session ID, Total Tokens, WebSocket Status
- ✅ แสดง Last Response Time (Latency)
- ✅ ปุ่ม "Clear History" (DELETE `/api/chats`)
- ✅ ปุ่ม "Refresh" อัพเดตข้อมูล
- ✅ Auto-refresh ทุก 5 วินาที
- ✅ Collapsible panel
- ✅ Color-coded WebSocket status (Connected/Connecting/Disconnected)

### โครงสร้าง Components:
```
frontend/src/components/
├── common/
│   ├── BaseButton.vue         ✅
│   ├── BaseModal.vue          ✅
│   ├── BaseInput.vue          ✅
│   ├── BaseSelect.vue         ✅
│   ├── BaseTextarea.vue       ✅
│   └── index.js               ✅
└── persona/
    ├── PersonaSidebar.vue     ✅
    ├── PersonaModal.vue       ✅
    ├── DebugPanel.vue         ✅
    └── index.js               ✅
```

### Integration:
- ✅ อัพเดต App.vue ให้ใช้ PersonaSidebar
- ✅ เชื่อมต่อกับ Pinia stores (personaStore, chatStore)
- ✅ ใช้ API services ที่สร้างไว้ใน Task 1

### Build Test:
```bash
✓ 91 modules transformed
✓ built in 1.34s
```

**Status:** ✅ Task 2 เสร็จสมบูรณ์ - พร้อมสำหรับ Task 3

---

## Task 3: Create Chat Interface Components ✅ COMPLETED

**วัตถุประสงค์:** สร้างส่วนแสดงและส่งข้อความ

### สิ่งที่ทำเสร็จแล้ว:

✅ **ChatBubble.vue** (`src/components/chat/`)
**Features:**
- ✅ แยก UI สำหรับ User/Bot/System messages
- ✅ แสดง Avatar พร้อม persona icon
- ✅ แสดง File Attachments ("แนบไฟล์ ${Filename}")
- ✅ แสดง Message content พร้อม markdown formatting (bold, italic, code)
- ✅ แสดง Timestamp (relative time: "Just now", "5m ago")
- ✅ แสดง Tokens used
- ✅ Streaming indicator (blinking cursor)
- ✅ Smooth animations

✅ **LoadingIndicator.vue** (`src/components/chat/`)
**Features:**
- ✅ Animated dots (3 bouncing dots)
- ✅ Loading message ("Bot กำลังประมวลผล...")
- ✅ Persona avatar
- ✅ Fade-in animation

✅ **ChatLog.vue** (`src/components/chat/`)
**Features:**
- ✅ แสดงประวัติการสนทนาทั้งหมด
- ✅ Empty state ("เริ่มการสนทนา")
- ✅ Auto-scroll to bottom เมื่อมีข้อความใหม่
- ✅ Scroll to bottom button (แสดงเมื่อ scroll ขึ้น)
- ✅ Loading indicator ขณะ Bot ตอบกลับ
- ✅ WebSocket integration (connect on mount)
- ✅ Smooth scroll behavior
- ✅ Custom scrollbar styling

✅ **ChatInput.vue** (`src/components/chat/`)
**Features:**
- ✅ Textarea auto-resize (max 150px)
- ✅ Enter to send, Shift+Enter for new line
- ✅ File attachment button
- ✅ File preview panel (แสดงไฟล์ที่เลือก)
- ✅ Remove individual files
- ✅ Send button (disabled เมื่อ input ว่าง)
- ✅ Loading state (spinner)
- ✅ File upload workflow:
  - มีไฟล์: Upload → ได้ file_ids → Send via WebSocket
  - ไม่มีไฟล์: Send via WebSocket โดยตรง
- ✅ Supported file types: TXT, PDF, DOCX, XLSX, Images
- ✅ File size display
- ✅ Upload progress indicator

### โครงสร้าง Components:
```
frontend/src/components/chat/
├── ChatBubble.vue         ✅ (270 lines)
├── LoadingIndicator.vue   ✅ (120 lines)
├── ChatLog.vue            ✅ (230 lines)
├── ChatInput.vue          ✅ (420 lines)
└── index.js               ✅
```

### Integration:
- ✅ อัพเดต App.vue ให้ใช้ ChatLog + ChatInput
- ✅ เชื่อมต่อกับ Pinia stores (chatStore, personaStore, fileStore)
- ✅ เชื่อมต่อกับ WebSocket service
- ✅ File upload integration

### Build Test:
```bash
✓ 102 modules transformed
✓ built in 1.56s
```

### UI/UX Improvements:
- ✅ **Full-screen Layout** - หน้าเว็บครอบคลุมพื้นที่ทั้งหมด (ไม่มีพื้นที่สีเทา)
- ✅ **Chat Container** - ช่องแชทขยายเต็มพื้นที่ main content
- ✅ **Fixed Positioning** - ใช้ position: fixed เพื่อให้แน่ใจว่าครอบคลุม viewport
- ✅ **Proper Height/Width** - ตั้งค่า 100% height/width ทุก level (html, body, #app, .app)
- ✅ **No Scrollbars** - ซ่อน scrollbar ของ body, เหลือแค่ใน messages-container

**Status:** ✅ Task 3 เสร็จสมบูรณ์ - Chat Interface พร้อมใช้งาน 100%!

---

## Task 4: Create File Upload Components

**วัตถุประสงค์:** จัดการการอัปโหลดไฟล์

### Component: `FileUploadModal.vue`

**Features:**
- Drag-and-Drop Area
- File List แสดงไฟล์ที่เลือก (ชื่อ, ขนาด, สถานะ)
- แสดงประเภทไฟล์ที่รองรับ: PDF, DOCX, TXT, CSV, JPEG, PNG
- ปุ่ม Upload

**API Call:**
```javascript
// Upload files
const formData = new FormData()
files.forEach(file => formData.append('files', file))

axios.post('/api/file/uploads', formData, {
  headers: { 'Content-Type': 'multipart/form-data' }
})
.then(response => {
  // เก็บ file_ids จาก response.uploaded_files
  const fileIds = response.data.uploaded_files.map(f => f.file_id)
})
```

**Supported File Types:**
- Text: TXT, MD, JSON, CSV, XML (max 10 MB)
- Documents: PDF, DOCX, XLSX (max 25 MB)
- Images: JPG, PNG, GIF, WebP (max 20 MB)

---

## Task 5: Create Speech-to-Speech Components

**วัตถุประสงค์:** เพิ่ม Speech-to-Speech mode

### Component: `SpeechControl.vue`

**Features:**
- Toggle Switch: โหมด Text-to-Text / Speech-to-Speech
- Microphone Button (เริ่ม/หยุด บันทึกเสียง)
- Status Indicator:
  - "Listening..." (กำลังอัดเสียง)
  - "Processing..." (กำลังถอดเสียง)
  - "Bot Speaking..." (กำลังเล่นเสียง)
- Transcript Display (แสดงข้อความที่ถอดเสียงได้)

**Workflow:**
1. User อัดเสียง → ส่ง POST `/api/audio/transcribe`
2. ได้ข้อความ → ส่ง POST `/api/chat` หรือ WebSocket
3. ได้ response → ส่ง POST `/api/audio/tts` (OpenAI) หรือ `/audio/elevenlabs/tts`
4. เล่นเสียงให้ User ฟัง (สามารถหยุดเสียงก่อนจบได้)

**API Calls:**
```javascript
// Speech-to-Text (Whisper)
const audioFormData = new FormData()
audioFormData.append('file', audioBlob)
axios.post('/api/audio/transcribe', audioFormData)

// Text-to-Speech (OpenAI)
axios.post('/api/audio/tts', {
  text: botResponse,
  voice: "nova",
  model: "tts-1"
})

// Text-to-Speech (ElevenLabs - better quality)
axios.post('/audio/elevenlabs/tts', {
  text: botResponse,
  voice_id: "21m00Tcm4TlvDq8ikWAM",
  model_id: "eleven_multilingual_v2"
}, {
  headers: { 'Accept': 'audio/mpeg' },
  responseType: 'blob'
})
```

**Audio Playback:**
```javascript
const audio = new Audio(URL.createObjectURL(audioBlob))
audio.play()
// สามารถ audio.pause() ได้ตลอดเวลา
```

---

## Task 6: State Management (Pinia Stores)

**วัตถุประสงค์:** จัดการ State ของแอพ

### Store: `stores/chat.js`

**State:**
- `messages: []` - ประวัติการสนทนา
- `sessionId: string` - Session ID ปัจจุบัน
- `isLoading: boolean` - สถานะกำลังโหลด
- `webSocket: WebSocket` - WebSocket connection

**Actions:**
- `sendMessage(content, personaId, fileIds)`
- `loadChatHistory()`
- `clearChatHistory()` (DELETE `/api/chats`)
- `connectWebSocket()`
- `disconnectWebSocket()`

### Store: `stores/persona.js`

**State:**
- `personas: []` - รายการ personas ทั้งหมด
- `selectedPersona: object` - Persona ที่เลือก
- `isLoadingPersonas: boolean`

**Actions:**
- `fetchPersonas()` (GET `/api/personas`)
- `fetchPersonaById(id)` (GET `/api/personas/:id`)
- `createPersona(data)` (POST `/api/persona`)
- `updatePersona(id, data)` (PATCH `/api/persona/:id`)
- `deletePersona(id)` (DELETE `/api/persona/:id`)
- `selectPersona(persona)`

### Store: `stores/audio.js`

**State:**
- `isRecording: boolean`
- `isProcessing: boolean`
- `isSpeaking: boolean`
- `transcript: string`
- `audioMode: 'text' | 'speech'`

**Actions:**
- `startRecording()`
- `stopRecording()`
- `transcribeAudio(audioBlob)` (POST `/api/audio/transcribe`)
- `textToSpeech(text, voiceOptions)` (POST `/api/audio/tts` หรือ `/audio/elevenlabs/tts`)
- `stopAudioPlayback()`

---

## Task 7: API Service Layer

**วัตถุประสงค์:** สร้าง API service functions

### File: `src/api/axios.js`

```javascript
import axios from 'axios'

const apiClient = axios.create({
  baseURL: 'http://localhost:3001',
  headers: {
    'Content-Type': 'application/json'
  }
})

export default apiClient
```

### File: `src/api/personaService.js`

```javascript
import apiClient from './axios'

export const personaService = {
  getAll: () => apiClient.get('/api/personas'),
  getById: (id) => apiClient.get(`/api/personas/${id}`),
  create: (data) => apiClient.post('/api/persona', data),
  update: (id, data) => apiClient.patch(`/api/persona/${id}`, data),
  delete: (id) => apiClient.delete(`/api/persona/${id}`)
}
```

### File: `src/api/chatService.js`

```javascript
import apiClient from './axios'

export const chatService = {
  sendMessage: (data) => apiClient.post('/api/chat', data),
  getChatHistory: (params) => apiClient.get('/api/chats', { params }),
  deleteAllMessages: () => apiClient.delete('/api/chats')
}
```

### File: `src/api/fileService.js`

```javascript
import apiClient from './axios'

export const fileService = {
  uploadFiles: (formData) =>
    apiClient.post('/api/file/uploads', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    }),
  getFileHistory: (params) => apiClient.get('/api/file/history', { params }),
  deleteAllFiles: () => apiClient.delete('/api/file/uploads')
}
```

### File: `src/api/audioService.js`

```javascript
import apiClient from './axios'

export const audioService = {
  transcribe: (audioFormData) =>
    apiClient.post('/api/audio/transcribe', audioFormData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    }),
  textToSpeech: (data) =>
    apiClient.post('/api/audio/tts', data, {
      responseType: 'blob'
    }),
  elevenLabsTTS: (data) =>
    apiClient.post('/audio/elevenlabs/tts', data, {
      headers: { 'Accept': 'audio/mpeg' },
      responseType: 'blob'
    }),
  getElevenLabsVoices: () =>
    apiClient.get('/audio/elevenlabs/voices')
}
```

---

## Task 8: WebSocket Integration

**วัตถุประสงค์:** จัดการ WebSocket connection

### File: `src/services/websocketService.js`

```javascript
class WebSocketService {
  constructor() {
    this.ws = null
    this.messageHandlers = []
    this.errorHandlers = []
  }

  connect() {
    this.ws = new WebSocket('ws://localhost:3001/api/chat/stream')

    this.ws.onopen = () => {
      console.log('WebSocket connected')
    }

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      this.messageHandlers.forEach(handler => handler(data))
    }

    this.ws.onerror = (error) => {
      this.errorHandlers.forEach(handler => handler(error))
    }

    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
    }
  }

  send(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    }
  }

  onMessage(handler) {
    this.messageHandlers.push(handler)
  }

  onError(handler) {
    this.errorHandlers.push(handler)
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }
}

export default new WebSocketService()
```

---

## Task 9: UI/UX Components

**วัตถุประสงค์:** สร้าง UI components ที่สวยงามและใช้งานง่าย

### Recommended Component Library:
- **Option 1:** Vuetify (Material Design)
- **Option 2:** Element Plus (Enterprise UI)
- **Option 3:** PrimeVue (Rich UI components)
- **Option 4:** Custom CSS + Tailwind CSS

### Components to Create:
1. `LoadingSpinner.vue` - แสดง loading animation
2. `ChatBubble.vue` - แสดงข้อความ User/Bot
3. `Button.vue` - ปุ่มพื้นฐาน
4. `Dropdown.vue` - Dropdown menu
5. `Modal.vue` - Modal dialog
6. `FileItem.vue` - แสดงข้อมูลไฟล์
7. `AudioVisualizer.vue` - แสดง waveform ขณะอัดเสียง

---

## Task 10: Integration & Testing

**วัตถุประสงค์:** รวม components และทดสอบการทำงาน

### Integration Steps:
1. รวม PersonaSidebar + ChatLog + ChatInput ใน Main Layout
2. เชื่อม FileUploadModal กับ ChatInput
3. เชื่อม SpeechControl กับ Audio API
4. ทดสอบ WebSocket streaming
5. ทดสอบ File Upload workflow
6. ทดสอบ Speech-to-Speech workflow

### Test Scenarios:
1. **Text Chat:**
   - เลือก Persona → พิมพ์ข้อความ → ได้รับ response แบบ streaming

2. **File Upload + Chat:**
   - อัปโหลดไฟล์ PDF → ส่งข้อความถาม → Bot ตอบจากเนื้อหาไฟล์

3. **Speech-to-Speech:**
   - อัดเสียง → ถอดเสียงเป็นข้อความ → ส่งให้ Bot → รับเสียงตอบกลับ

4. **Persona Management:**
   - สร้าง Persona ใหม่ → เลือกใช้งาน → แก้ไข → ลบ

5. **Clear History:**
   - ส่งข้อความหลายข้อ → Clear History → ตรวจสอบว่าข้อความหายหมด

---

## Task 11: Error Handling & Edge Cases

**วัตถุประสงค์:** จัดการ errors และ edge cases

### Error Scenarios:
1. **WebSocket Disconnection:**
   - แสดงข้อความ "กำลังเชื่อมต่อใหม่..."
   - Auto-reconnect mechanism

2. **API Errors:**
   - แสดง error message จาก backend
   - Fallback UI สำหรับ failed requests

3. **File Upload Errors:**
   - ไฟล์ใหญ่เกินไป → แสดง error
   - ประเภทไฟล์ไม่รองรับ → แสดง error

4. **Audio Recording Errors:**
   - ไม่มี microphone permission → ขอ permission
   - Audio format ไม่รองรับ → แสดง error

5. **Empty Messages:**
   - ป้องกันการส่งข้อความว่าง
   - Disable send button เมื่อ input ว่าง

---

## Summary Checklist

- [ ] **Task 1:** Setup Project (Vue 3 + Axios + Pinia)
- [ ] **Task 2:** Persona Management (Sidebar + Modal + CRUD)
- [ ] **Task 3:** Chat Interface (ChatLog + ChatInput + WebSocket)
- [ ] **Task 4:** File Upload (Modal + Drag-and-Drop)
- [ ] **Task 5:** Speech-to-Speech (Recording + Transcription + TTS)
- [ ] **Task 6:** State Management (Pinia Stores: chat, persona, audio)
- [ ] **Task 7:** API Services (personaService, chatService, fileService, audioService)
- [ ] **Task 8:** WebSocket Service (websocketService.js)
- [ ] **Task 9:** UI Components (LoadingSpinner, ChatBubble, Modal, etc.)
- [ ] **Task 10:** Integration & Testing (ทดสอบ workflows ทั้งหมด)
- [ ] **Task 11:** Error Handling (Disconnection, API errors, validation)

---

## Key API Endpoints Summary

| Feature | Method | Endpoint | Purpose |
|---------|--------|----------|---------|
| Personas | GET | `/api/personas` | ดึงรายการ personas |
| Persona Detail | GET | `/api/personas/:id` | ดึงข้อมูล persona |
| Create Persona | POST | `/api/persona` | สร้าง persona ใหม่ |
| Update Persona | PATCH | `/api/persona/:id` | แก้ไข persona |
| Delete Persona | DELETE | `/api/persona/:id` | ลบ persona |
| Chat | POST | `/api/chat` | ส่งข้อความ (non-streaming) |
| Chat Streaming | WS | `/api/chat/stream` | ส่งข้อความ (streaming) |
| Chat History | GET | `/api/chats` | ดึงประวัติการสนทนา |
| Clear History | DELETE | `/api/chats` | ลบประวัติทั้งหมด |
| Upload Files | POST | `/api/file/uploads` | อัปโหลดไฟล์ |
| Transcribe | POST | `/api/audio/transcribe` | Speech-to-Text |
| TTS (OpenAI) | POST | `/api/audio/tts` | Text-to-Speech |
| TTS (ElevenLabs) | POST | `/audio/elevenlabs/tts` | Text-to-Speech (high quality) |
| Voices (ElevenLabs) | GET | `/audio/elevenlabs/voices` | ดึง voice list |

---

**End of Workflow**
