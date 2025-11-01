# Frontend Documentation - ChatBot Project

**เวอร์ชั่น**: 1.0
**อัพเดตล่าสุด**: 2025-11-01
**สถานะ**: พร้อมพัฒนา (Backend สมบูรณ์แล้ว)

---

## สารบัญ

1. [ภาพรวมและสถานะ](#ภาพรวมและสถานะ)
2. [เทคโนโลยีที่ใช้](#เทคโนโลยีที่ใช้)
3. [โครงสร้างโปรเจ็ค](#โครงสร้างโปรเจ็ค)
4. [การติดตั้งและเริ่มต้น](#การติดตั้งและเริ่มต้น)
5. [Features ที่ต้องพัฒนา](#features-ที่ต้องพัฒนา)
6. [API Endpoints ที่พร้อมใช้](#api-endpoints-ที่พร้อมใช้)
7. [คู่มือการพัฒนา](#คู่มือการพัฒนา)

---

## ภาพรวมและสถานะ

### Backend Status: ✅ สมบูรณ์ (95%)

Backend พร้อมใช้งานแล้ว มี API endpoints ครบถ้วนสำหรับ:
- ✅ Chat System (พร้อม WebSocket streaming)
- ✅ Multi-Persona System (3 personas พร้อมใช้)
- ✅ File Upload & Analysis (PDF, DOCX, XLSX, Images)
- ✅ Audio (Speech-to-Text, Text-to-Speech)
- ✅ Database (PostgreSQL พร้อม schema)

### Frontend Status: ⏳ รอการพัฒนา

Frontend มี **โครงสร้างพื้นฐาน** อยู่แล้ว แต่ต้องปรับให้สอดคล้องกับ Backend:
- ⚠️ มี Vue 3 components พื้นฐาน
- ⚠️ ต้องอัพเดต API integration
- ⚠️ ต้องเพิ่ม File Upload UI
- ⚠️ ต้องปรับ session management

---

## เทคโนโลยีที่ใช้

| เทคโนโลยี | เวอร์ชั่น | หมายเหตุ |
|-----------|----------|----------|
| **Vue.js** | 3.4+ | Composition API |
| **Vite** | 5.0+ | Build tool |
| **Pinia** | 2.1+ | State management |
| **Vue Router** | 4.0+ | Routing |
| **Axios** | 1.6+ | HTTP client |
| **Tailwind CSS** | 3.4+ | Styling |
| **WebSocket** | Native API | Real-time streaming |

---

## โครงสร้างโปรเจ็ค

```
frontend/
├── src/
│   ├── components/
│   │   ├── chat/
│   │   │   ├── MessageBubble.vue       ✅ มีอยู่แล้ว
│   │   │   ├── MessageInput.vue        ⚠️ ต้องเพิ่ม File Upload
│   │   │   └── MessageList.vue         ✅ มีอยู่แล้ว
│   │   ├── file/                       ❌ ยังไม่มี - ต้องสร้าง
│   │   │   ├── FileUpload.vue
│   │   │   ├── FileList.vue
│   │   │   └── FileAnalysisCard.vue
│   │   └── layout/
│   │       └── AppLayout.vue           ✅ มีอยู่แล้ว
│   ├── composables/
│   │   ├── useWebSocket.js             ✅ มีอยู่แล้ว
│   │   ├── useAudioRecorder.js         ✅ มีอยู่แล้ว
│   │   └── useFileUpload.js            ❌ ยังไม่มี - ต้องสร้าง
│   ├── services/
│   │   ├── api.js                      ✅ มีอยู่แล้ว
│   │   ├── chatService.js              ⚠️ ต้องเพิ่ม session_id, file_attachments
│   │   ├── audioService.js             ✅ มีอยู่แล้ว
│   │   └── fileService.js              ❌ ยังไม่มี - ต้องสร้าง
│   ├── store/
│   │   ├── chat.js                     ⚠️ ต้องเพิ่ม session management
│   │   ├── personas.js                 ✅ มีอยู่แล้ว
│   │   └── files.js                    ❌ ยังไม่มี - ต้องสร้าง
│   └── views/
│       ├── ChatView.vue                ⚠️ ต้องเพิ่ม File Upload UI
│       └── SpeechView.vue              ✅ มีอยู่แล้ว
├── .env.development                    ⚠️ ต้องตรวจสอบ API URL
└── package.json
```

---

## การติดตั้งและเริ่มต้น

### 1. ติดตั้ง Dependencies

```bash
cd frontend
npm install
```

### 2. ตั้งค่า Environment Variables

สร้างไฟล์ `.env.development`:

```env
VITE_API_BASE_URL=http://localhost:3000/api
VITE_WS_URL=ws://localhost:3000/api/chat/stream
```

### 3. รัน Development Server

```bash
npm run dev
```

Frontend จะรันที่ `http://localhost:5173`

---

## Features ที่ต้องพัฒนา

### Priority 1: Chat System พร้อม File Upload ⚡

**ความสำคัญ**: สูงสุด
**Backend Status**: ✅ พร้อมใช้

**สิ่งที่ต้องทำ**:

1. **อัพเดต MessageInput Component**
   - เพิ่มปุ่ม Upload File
   - รองรับ Multiple Files (สูงสุด 5 ไฟล์/ข้อความ)
   - แสดง Preview ไฟล์ก่อนส่งในรูปแบบ card
   - ตรวจสอบ File Type และ Size
   - เพิ่มปุ่มลบ Preview ไฟล์ก่อนส่ง ให้อยู่มุมบนขวาของ card

2. **สร้าง File Upload Service**

   **ไฟล์**: `src/services/fileService.js`

   ```javascript
   import api from './api'

   export const fileService = {
     // Upload และวิเคราะห์ไฟล์
     async analyzeFile(file, options = {}) {
       const formData = new FormData()
       formData.append('file', file)
       formData.append('analysis_type', options.analysisType || 'summary')
       formData.append('language', options.language || 'th')
       formData.append('session_id', options.sessionId || '')

       if (options.prompt) {
         formData.append('prompt', options.prompt)
       }

       return await api.post('/file/analyze', formData, {
         headers: { 'Content-Type': 'multipart/form-data' }
       })
     },

     // ดึงประวัติไฟล์
     async getFileHistory(sessionId, limit = 20) {
       return await api.get('/file/history', {
         params: { session_id: sessionId, limit }
       })
     }
   }
   ```

3. **อัพเดต Chat Service เพิ่ม session_id และ use_history**

   **ไฟล์**: `src/services/chatService.js` (เพิ่ม)

   ```javascript
   async sendMessage(message, sessionId, personaId = null, options = {}) {
     const payload = {
       message,
       session_id: sessionId,      // ⚡ เพิ่ม session_id
       persona_id: personaId,
       use_history: true,          // ⚡ เปิด conversation history
       temperature: options.temperature || 0.7,
       max_tokens: options.maxTokens || 0,
       model: options.model || 'gpt-4o-mini'
     }

     return await api.post('/chat', payload)
   }
   ```

4. **อัพเดต Chat Store เพิ่ม Session Management**

   **ไฟล์**: `src/store/chat.js` (เพิ่ม)

   ```javascript
   state: () => ({
     messages: [],
     sessionId: null,              // ⚡ เพิ่ม session tracking
     uploadedFiles: [],            // ⚡ ไฟล์ที่อัปโหลดใน session ปัจจุบัน
     isLoading: false,
     isStreaming: false,
     currentPersonaId: 1
   }),

   actions: {
     // สร้าง session ใหม่
     createNewSession() {
       this.sessionId = `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
       this.messages = []
       this.uploadedFiles = []
       return this.sessionId
     },

     // Upload ไฟล์
     async uploadFiles(files, analysisType = 'summary') {
       for (const file of files) {
         const result = await fileService.analyzeFile(file, {
           sessionId: this.sessionId,
           analysisType,
           language: 'th'
         })
         this.uploadedFiles.push({
           fileId: result.file_id,
           filename: result.filename,
           summary: result.summary
         })
       }
     }
   }
   ```

---

### Priority 2: File Upload UI Components 📁

**Backend Endpoints พร้อมใช้**:
- `POST /api/file/analyze` - อัปโหลดและวิเคราะห์ไฟล์
- `GET /api/file/history` - ดูประวัติไฟล์

**Components ที่ต้องสร้าง**:

#### 1. FileUpload.vue
```vue
<template>
  <div class="file-upload">
    <input
      ref="fileInput"
      type="file"
      multiple
      accept=".pdf,.docx,.xlsx,.txt,.png,.jpg,.jpeg"
      @change="handleFileSelect"
      class="hidden"
    />

    <button @click="$refs.fileInput.click()" class="upload-btn">
      📎 แนบไฟล์ (สูงสุด 5 ไฟล์)
    </button>

    <!-- Preview ไฟล์ที่เลือก -->
    <div v-if="selectedFiles.length > 0" class="file-preview">
      <div v-for="file in selectedFiles" :key="file.name" class="file-item">
        {{ file.name }} ({{ formatFileSize(file.size) }})
        <button @click="removeFile(file)">✕</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const selectedFiles = ref([])
const emit = defineEmits(['files-selected'])

const handleFileSelect = (event) => {
  const files = Array.from(event.target.files)

  // จำกัด 5 ไฟล์
  if (selectedFiles.value.length + files.length > 5) {
    alert('สามารถเลือกได้สูงสุด 5 ไฟล์')
    return
  }

  selectedFiles.value.push(...files)
  emit('files-selected', selectedFiles.value)
}

const removeFile = (file) => {
  selectedFiles.value = selectedFiles.value.filter(f => f !== file)
  emit('files-selected', selectedFiles.value)
}

const formatFileSize = (bytes) => {
  return (bytes / 1024 / 1024).toFixed(2) + ' MB'
}
</script>
```

#### 2. ปรับ MessageInput.vue เพิ่ม File Upload

```vue
<template>
  <div class="message-input">
    <!-- File Upload Component -->
    <FileUpload @files-selected="handleFilesSelected" />

    <textarea
      v-model="message"
      @keydown.enter.prevent="sendMessage"
      placeholder="พิมพ์ข้อความ..."
    />

    <button
      @click="sendMessage"
      :disabled="!canSend"
    >
      ส่ง
    </button>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useChatStore } from '@/store/chat'
import FileUpload from '@/components/file/FileUpload.vue'

const chatStore = useChatStore()
const message = ref('')
const selectedFiles = ref([])

const canSend = computed(() => {
  return message.value.trim() || selectedFiles.value.length > 0
})

const handleFilesSelected = (files) => {
  selectedFiles.value = files
}

const sendMessage = async () => {
  if (!canSend.value) return

  // อัปโหลดไฟล์ก่อน (ถ้ามี)
  if (selectedFiles.value.length > 0) {
    await chatStore.uploadFiles(selectedFiles.value, 'summary')
  }

  // ส่งข้อความ
  await chatStore.sendMessage(message.value)

  // รีเซ็ต
  message.value = ''
  selectedFiles.value = []
}
</script>
```

---

### Priority 3: WebSocket Streaming ⚡

**Backend**: ✅ WebSocket endpoint พร้อมใช้ที่ `ws://localhost:3000/api/chat/stream`

**สิ่งที่ต้องทำ**:

1. **ปรับ useWebSocket.js ให้รองรับ session_id และ file context**

```javascript
function sendMessage(content, sessionId, personaId = 1) {
  if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
    console.error('WebSocket is not connected')
    return false
  }

  const message = {
    message: content,
    session_id: sessionId,     // ⚡ ส่ง session_id
    persona_id: personaId,
    use_history: true          // ⚡ เปิด history
  }

  ws.value.send(JSON.stringify(message))
  return true
}
```

2. **ปรับ Chat Store ให้ใช้ WebSocket แทน HTTP**

```javascript
async sendMessage(message) {
  // ใช้ WebSocket streaming แทน HTTP
  this.sendStreamingMessage(message)
}

sendStreamingMessage(message) {
  const { sendMessage } = useWebSocket()

  const userMessage = {
    role: 'user',
    content: message,
    created_at: new Date().toISOString()
  }
  this.addMessage(userMessage)

  this.isStreaming = true
  this.streamingContent = ''

  // ส่งผ่าน WebSocket พร้อม session_id
  sendMessage(message, this.sessionId, this.currentPersonaId)
}
```

---

### Priority 4: Session Management & History 📝

**Backend**: ✅ รองรับ `session_id` และ `use_history`

**สิ่งที่ต้องทำ**:

1. **สร้าง Session ใหม่เมื่อเริ่มต้น**

```javascript
// src/store/chat.js
actions: {
  initializeChat() {
    // สร้าง session ใหม่
    if (!this.sessionId) {
      this.createNewSession()
    }

    // โหลดประวัติ (ถ้ามี)
    this.fetchHistory()
  },

  createNewSession() {
    this.sessionId = `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    this.messages = []
    this.uploadedFiles = []
    localStorage.setItem('current_session_id', this.sessionId)
  },

  async fetchHistory() {
    const response = await chatService.getChatHistory(this.sessionId, 50)
    this.messages = response.messages
  }
}
```

2. **ปรับ ChatView.vue เพื่อ initialize session**

```javascript
// src/views/ChatView.vue
import { onMounted } from 'vue'
import { useChatStore } from '@/store/chat'

const chatStore = useChatStore()

onMounted(() => {
  chatStore.initializeChat()
})
```

---

## API Endpoints ที่พร้อมใช้

### 1. Chat Endpoints ✅

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/chat` | ส่งข้อความ (รองรับ `session_id`, `use_history`) |
| GET | `/api/chat/history` | ดึงประวัติการสนทนา (ใช้ `session_id`) |
| WS | `/api/chat/stream` | WebSocket streaming |

**ตัวอย่างการใช้งาน**:

```javascript
// ส่งข้อความพร้อม session และ history
await chatService.sendMessage('สวัสดี', 'session_123', 1, {
  use_history: true
})

// ดึงประวัติ
const history = await chatService.getChatHistory('session_123', 50)
```

### 2. File Endpoints ✅

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/file/analyze` | อัปโหลดและวิเคราะห์ไฟล์ (รองรับ `session_id`) |
| GET | `/api/file/history` | ดึงประวัติไฟล์ (filter ด้วย `session_id`) |

**รองรับไฟล์**:
- เอกสาร: PDF, DOCX, XLSX, TXT (max 10MB)
- รูปภาพ: PNG, JPG, JPEG (max 20MB)

**Analysis Types**:
- `summary` - สรุปเนื้อหา
- `detail` - วิเคราะห์ละเอียด
- `qa` - ถาม-ตอบจากไฟล์
- `extract` - ดึงข้อมูลสำคัญ

### 3. Persona Endpoints ✅

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/personas` | รายชื่อ personas ทั้งหมด (3 personas) |
| GET | `/api/personas/:id` | รายละเอียด persona |

**Personas พร้อมใช้**:
1. General Assistant (ID: 1)
2. Technology Expert (ID: 2)
3. Business Advisor (ID: 3)

### 4. Audio Endpoints ✅

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/audio/transcribe` | แปลงเสียงเป็นข้อความ |
| POST | `/api/audio/tts` | แปลงข้อความเป็นเสียง |

---

## คู่มือการพัฒนา

### ขั้นตอนการพัฒนาแบบละเอียด

#### Week 1: File Upload & Session Management

**Day 1-2**: สร้าง File Upload Components
- [ ] สร้าง `FileUpload.vue`
- [ ] สร้าง `FileList.vue` (แสดงไฟล์ที่อัปโหลดแล้ว)
- [ ] สร้าง `fileService.js`
- [ ] ทดสอบ upload ไฟล์ผ่าน API

**Day 3-4**: Session Management
- [ ] เพิ่ม session tracking ใน Chat Store
- [ ] สร้าง `createNewSession()` action
- [ ] เพิ่ม `session_id` ใน chat requests
- [ ] ทดสอบ conversation history

**Day 5**: Integration
- [ ] เชื่อม File Upload กับ MessageInput
- [ ] ทดสอบส่งข้อความพร้อมไฟล์
- [ ] ทดสอบ AI ตอบคำถามจากไฟล์

#### Week 2: WebSocket & Real-time Features

**Day 1-2**: WebSocket Integration
- [ ] ปรับ `useWebSocket.js` รองรับ session
- [ ] เพิ่ม file context ใน WebSocket request
- [ ] ทดสอบ streaming responses

**Day 3**: UI Polish
- [ ] เพิ่ม loading states
- [ ] แสดงไฟล์ที่แนบใน message bubbles
- [ ] เพิ่ม error handling

**Day 4-5**: Testing & Bug Fixes
- [ ] ทดสอบ flow ทั้งหมด
- [ ] แก้ไข bugs
- [ ] Optimize performance

---

## ตัวอย่าง Flow การใช้งาน

### Flow 1: Chat ปกติ (พร้อม History)

```
1. User เปิดแอพ
   → Frontend สร้าง session_id ใหม่

2. User พิมพ์ข้อความ "สวัสดี"
   → ส่ง POST /api/chat {
       message: "สวัสดี",
       session_id: "session_123",
       use_history: true
     }

3. AI ตอบ "สวัสดีครับ"
   → บันทึกใน messages table พร้อม session_id

4. User พิมพ์ "ชื่ออะไร"
   → Backend โหลด history จาก session_123
   → AI จำบริบทได้
```

### Flow 2: Chat พร้อม File Upload

```
1. User อัปโหลดไฟล์ PDF
   → POST /api/file/analyze {
       file: report.pdf,
       session_id: "session_123",
       analysis_type: "summary"
     }
   → Backend วิเคราะห์และบันทึกใน file_analyses table

2. User ถามคำถาม "สรุปเอกสารนี้หน่อย"
   → POST /api/chat {
       message: "สรุปเอกสารนี้หน่อย",
       session_id: "session_123",
       use_history: true
     }
   → Backend โหลดไฟล์จาก session_123 (สูงสุด 5 ไฟล์)
   → AI ตอบโดยอ้างอิงไฟล์
```

### Flow 3: WebSocket Streaming

```
1. User ส่งข้อความ
   → WebSocket.send({
       message: "เล่าเรื่องยาวๆ",
       session_id: "session_123",
       use_history: true
     })

2. Backend stream response กลับมา
   → { type: "start", session_id: "session_123" }
   → { type: "chunk", content: "กา" }
   → { type: "chunk", content: "ลค" }
   → { type: "chunk", content: "รั้" }
   → { type: "done", tokens_used: 250 }

3. Frontend แสดงข้อความทีละตัวอักษร
```

---

## การทดสอบ

### 1. ทดสอบ API Integration

```bash
# ทดสอบ Chat API
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "สวัสดี",
    "session_id": "test_session",
    "use_history": true
  }'

# ทดสอบ File Upload
curl -X POST http://localhost:3000/api/file/analyze \
  -F "file=@test.pdf" \
  -F "session_id=test_session" \
  -F "analysis_type=summary"
```

### 2. ทดสอบ WebSocket

```javascript
const ws = new WebSocket('ws://localhost:3000/api/chat/stream')

ws.onopen = () => {
  ws.send(JSON.stringify({
    message: 'ทดสอบ WebSocket',
    session_id: 'test_session',
    use_history: true
  }))
}

ws.onmessage = (event) => {
  console.log('Response:', JSON.parse(event.data))
}
```

---

## Troubleshooting

### ปัญหาที่อาจพบ

**1. WebSocket Connection Failed**
```
Solution:
- ตรวจสอบ Backend รันที่ port 3000
- ตรวจสอบ VITE_WS_URL ใน .env
- เช็ค CORS settings
```

**2. File Upload Failed**
```
Solution:
- ตรวจสอบขนาดไฟล์ (PDF max 10MB, Image max 20MB)
- ตรวจสอบ file type
- เช็ค multipart/form-data header
```

**3. Session History ไม่ทำงาน**
```
Solution:
- ตรวจสอบ session_id ถูกส่งไปหรือไม่
- ตรวจสอบ use_history: true
- เช็ค Database มี messages ใน session นั้นหรือไม่
```

---

## Checklist สำหรับการพัฒนา

### Phase 1: Core Features ⚡
- [ ] Session Management (สร้าง/จัดการ session_id)
- [ ] File Upload UI (FileUpload.vue)
- [ ] อัพเดต MessageInput รองรับไฟล์
- [ ] อัพเดต chatService เพิ่ม session_id, use_history
- [ ] สร้าง fileService.js
- [ ] ทดสอบส่งข้อความพร้อมไฟล์

### Phase 2: WebSocket Streaming
- [ ] ปรับ useWebSocket รองรับ session
- [ ] ทดสอบ streaming responses
- [ ] แสดง typing indicator

### Phase 3: UI/UX Polish
- [ ] แสดงไฟล์ที่แนบใน message bubbles
- [ ] Loading states
- [ ] Error handling
- [ ] Responsive design

### Phase 4: Testing
- [ ] ทดสอบ File Upload flow
- [ ] ทดสอบ Session History
- [ ] ทดสอบ WebSocket Streaming
- [ ] ทดสอบบน Mobile

---

## API Base URL

**Development**:
```
HTTP: http://localhost:3000/api
WebSocket: ws://localhost:3000/api/chat/stream
```

**Production** (เมื่อ deploy):
```
HTTP: https://your-domain.com/api
WebSocket: wss://your-domain.com/api/chat/stream
```

---

## สรุป

Backend **พร้อมใช้งาน 95%** แล้ว มี API endpoints ครบถ้วน
Frontend ต้อง**พัฒนาเพิ่ม** 4 ส่วนหลัก:

1. ✅ **File Upload UI** - รองรับหลายไฟล์
2. ✅ **Session Management** - จัดการ conversation history
3. ✅ **WebSocket Integration** - Real-time streaming
4. ✅ **UI Polish** - แสดงไฟล์แนบ, loading states

ใช้เวลาพัฒนาประมาณ **2 สัปดาห์** หากทำเต็มเวลา

---

**อัพเดตล่าสุด**: 2025-11-01
**Next Update**: เมื่อเริ่มพัฒนา Frontend
