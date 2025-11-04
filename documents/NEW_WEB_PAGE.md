# Frontend Development Guide - ChatBot Project

## ภาพรวมโปรเจ็ค

โปรเจ็คนี้เป็น **AI Chatbot Platform** ที่รองรับการสื่อสารแบบ Text-to-Text และ Speech-to-Speech พร้อมการอัปโหลดไฟล์เพื่อให้ AI วิเคราะห์ และระบบ Persona ที่หนากหลาย

**Stack แนะนำ:**
- **Frontend Framework:** Vue 3 (Composition API)
- **HTTP Client:** Axios
- **WebSocket Client:** Native WebSocket API
- **UI Framework:** Tailwind CSS / Vuetify / Element Plus
- **State Management:** Pinia (แนะนำสำหรับ Vue 3)
- **Audio Recording:** MediaRecorder API / RecordRTC

---

## Backend API Endpoints

**Base URL:** `http://localhost:3001`

### 1. Personas API
- `GET /api/personas` - ดึงรายการ AI personas ทั้งหมด
- `GET /api/personas/:id` - ดึงข้อมูล persona ตาม ID

### 2. Chat API
- `POST /api/chat` - ส่งข้อความแบบ non-streaming (HTTP)
- `WS /api/chat/stream` - ส่งข้อความแบบ streaming (WebSocket)
- `GET /api/chats` - ดึงประวัติการสนทนา
- `DELETE /api/chats` - ลบข้อความทั้งหมด

### 3. File Upload API
- `POST /api/file/uploads` - อัปโหลดไฟล์ (รองรับสูงสุด 5 ไฟล์/ครั้ง)
- `GET /api/file/history` - ดึงประวัติการอัปโหลด
- `DELETE /api/file/uploads` - ลบบันทึกไฟล์ทั้งหมด

### 4. Audio API
- `POST /api/audio/transcribe` - แปลงเสียงเป็นข้อความ (Speech-to-Text)
- `POST /api/audio/tts` - แปลงข้อความเป็นเสียง (Text-to-Speech)

---

## UI Components ที่ต้องพัฒนา

## 1. ส่วนแสดงประวัติการสนทนา (Chat Log / Message Area)
พื้นที่ขนาดใหญ่ที่แสดงข้อความโต้ตอบทั้งหมด รวมถึงผลลัพธ์จากการวิเคราะห์ไฟล์

**Components:**
- **ChatBubble.vue** - แสดงข้อความแยกระหว่าง User และ Bot
  - รองรับการแสดงผลแบบ WebSocket streaming
  - แสดง loading state เมื่อรอ Bot ประมวลผล
  - แสดงผลการวิเคราะห์ไฟล์ในรูปแบบที่อ่านง่าย

**Props:**
```typescript
interface ChatBubbleProps {
  role: 'user' | 'assistant' | 'system';
  content: string;
  timestamp: string;
  isStreaming?: boolean;
  fileAnalysis?: FileAnalysisResult;
}
```

**Features ที่ต้องมี:**
- Markdown rendering สำหรับ Bot response
- Syntax highlighting สำหรับโค้ด
- Auto-scroll to bottom เมื่อมีข้อความใหม่
- Infinite scroll สำหรับโหลดประวัติย้อนหลัง

---

## 2. ส่วนป้อนข้อมูลหลัก (Input Area)
ส่วนที่อยู่ด้านล่างสุดของหน้าจอ ใช้สำหรับการสื่อสารพื้นฐาน

**Components:**
- **MessageInput.vue** - ช่องพิมพ์ข้อความและปุ่มส่ง

**Features:**
- Text Input Field: รองรับ multi-line, กด Enter ส่ง, Shift+Enter ขึ้นบรรทัดใหม่
- Send Button: เปลี่ยนสีเมื่อมีข้อความ, disabled ขณะส่ง
- Attachment Button: เปิด file picker หรือ drag-drop area

**State Management:**
```javascript
const state = {
  message: '',
  selectedFiles: [],
  isSending: false,
}
```

**Methods:**
```javascript
async function sendMessage() {
  // 1. ตรวจสอบว่ามีไฟล์แนบหรือไม่
  if (selectedFiles.length > 0) {
    // 2. Upload files ก่อน
    const fileIds = await uploadFiles(selectedFiles)

    // 3. ส่งข้อความพร้อม file_ids ผ่าน WebSocket
    sendWebSocketMessage({
      type: 'message',
      message: message,
      file_ids: fileIds,
      persona_id: selectedPersona.id
    })
  } else {
    // ส่งข้อความปกติผ่าน WebSocket
    sendWebSocketMessage({
      type: 'message',
      message: message,
      persona_id: selectedPersona.id
    })
  }
}
```

---

## 3. ส่วนควบคุม Speech-to-Speech
องค์ประกอบสำคัญที่เปิดใช้งานฟังก์ชันเสียง

**Components:**
- **MicrophoneButton.vue** - ปุ่มบันทึกเสียง
- **AudioStatusIndicator.vue** - แสดงสถานะการประมวลผลเสียง
- **TranscriptDisplay.vue** - แสดงข้อความที่ถอดจากเสียง

**Workflow:**
```
1. User กดปุ่ม Microphone → เริ่มบันทึกเสียง (MediaRecorder)
2. User ปล่อยปุ่ม → หยุดบันทึก
3. Frontend ส่ง audio file ไปยัง POST /api/audio/transcribe
4. แสดงข้อความที่ถอดได้ใน TranscriptDisplay
5. ส่งข้อความไปยัง POST /api/chat (หรือ WebSocket)
6. รับ response จาก AI
7. ส่ง response text ไปยัง POST /api/audio/tts
8. รับ audio file กลับมา
9. เล่นเสียงให้ user ฟังทันที (Audio API)
10. แสดง stop button เพื่อให้หยุดเสียงกลางคันได้
```

**State:**
```javascript
const audioState = {
  isRecording: false,
  isProcessing: false,
  isSpeaking: false,
  transcript: '',
  audioPlayer: null,
}
```

**Key Methods:**
```javascript
async function startRecording() {
  const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
  mediaRecorder = new MediaRecorder(stream)
  // ... handle recording
}

async function processAudio(audioBlob) {
  // 1. Speech-to-Text
  const formData = new FormData()
  formData.append('audio', audioBlob, 'recording.webm')
  const { text } = await axios.post('/api/audio/transcribe', formData)

  // 2. Send to AI
  const { response } = await axios.post('/api/chat', {
    message: text,
    persona_id: selectedPersona.id
  })

  // 3. Text-to-Speech
  const audioResponse = await axios.post('/api/audio/tts',
    { text: response },
    { responseType: 'blob' }
  )

  // 4. Play audio
  const audioUrl = URL.createObjectURL(audioResponse.data)
  audioPlayer.src = audioUrl
  audioPlayer.play()
}

function stopAudio() {
  audioPlayer.pause()
  audioPlayer.currentTime = 0
}
```

---

## 4. ส่วนจัดการไฟล์อัปโหลด (File Upload Modal/Widget)
ส่วนนี้จะปรากฏขึ้นเมื่อผู้ใช้คลิกปุ่มอัปโหลดไฟล์

**Components:**
- **FileUploadModal.vue** - Modal สำหรับอัปโหลดไฟล์
- **FileDragDropArea.vue** - พื้นที่ drag-and-drop
- **FileList.vue** - รายการไฟล์ที่เลือก

**Features:**
- Drag-and-Drop Area: ลากไฟล์มาวางได้
- File List: แสดงชื่อไฟล์, ขนาด, สถานะ (pending/uploading/uploaded/error)
- File Type Restriction: แสดงประเภทไฟล์ที่รองรับ
- Progress Bar: แสดง upload progress แต่ละไฟล์
- Remove Button: ลบไฟล์ออกจากรายการก่อนส่ง

**Upload Logic:**
```javascript
async function uploadFiles(files) {
  const formData = new FormData()

  // เพิ่มไฟล์ทั้งหมด (สูงสุด 5 ไฟล์)
  files.forEach(file => {
    formData.append('files', file)
  })

  try {
    const response = await axios.post('/api/file/uploads', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (progressEvent) => {
        // อัพเดต progress bar
        const percentCompleted = Math.round(
          (progressEvent.loaded * 100) / progressEvent.total
        )
        updateProgress(percentCompleted)
      }
    })

    // ส่ง file_ids กลับไป
    return response.data.uploaded_files.map(f => f.file_id)
  } catch (error) {
    // จัดการ error
    handleUploadError(error)
  }
}
```

**Supported File Types:**
- Documents: PDF, DOCX, XLSX, PPTX, TXT, MD, CSV, JSON, XML
- Images: JPG, PNG, GIF, WEBP, BMP
- Code: JS, PY, GO, JAVA, CPP, etc.
- Others: ZIP, RAR, MP3, MP4, etc.

**Validation:**
```javascript
function validateFiles(files) {
  // 1. ตรวจสอบจำนวน (สูงสุด 5 ไฟล์)
  if (files.length > 5) {
    throw new Error('Maximum 5 files allowed')
  }

  // 2. ตรวจสอบขนาดแต่ละไฟล์ (สูงสุด 10MB)
  const maxSize = 10 * 1024 * 1024 // 10MB
  files.forEach(file => {
    if (file.size > maxSize) {
      throw new Error(`File ${file.name} is too large`)
    }
  })

  return true
}
```

---

## 5. ส่วนตั้งค่า/เครื่องมือทดสอบ (Settings Panel)
ส่วนนี้สำคัญมากสำหรับโปรเจกต์ทดสอบ อาจเป็น Sidebar หรือ Modal

**Components:**
- **SettingsPanel.vue** - Panel หลัก
- **PersonaSelector.vue** - Dropdown เลือก Persona
- **LanguageSelector.vue** - เลือกภาษา
- **DebugStats.vue** - แสดงสถิติการใช้งาน

**Features:**

### 5.1 Persona Selector
```vue
<select v-model="selectedPersona" @change="onPersonaChange">
  <option v-for="persona in personas" :key="persona.id" :value="persona">
    {{ persona.icon }} {{ persona.name }} - {{ persona.expertise }}
  </option>
</select>
```

**Persona Data:**
```javascript
// ดึงจาก GET /api/personas
const personas = ref([
  {
    id: 1,
    name: "General Assistant",
    expertise: "general",
    icon: "🤖",
    tone: "friendly",
    temperature: 0.7,
    model: "gpt-4o-mini"
  },
  {
    id: 2,
    name: "Technology Expert",
    expertise: "technology",
    icon: "💻",
    tone: "professional",
    temperature: 0.5
  },
  // ... 6 personas อื่นๆ
])
```

### 5.2 Language Selector
```javascript
const languages = [
  { code: 'th', name: 'ไทย', flag: '🇹🇭' },
  { code: 'en', name: 'English', flag: '🇬🇧' },
]
```

### 5.3 Debug/Stats Display
แสดงข้อมูลเชิง technical:
- **Response Latency:** เวลาตอบกลับจาก AI (ms)
- **API Usage:** จำนวนครั้งที่เรียก API
- **Session ID:** ID ของ session ปัจจุบัน
- **Messages Count:** จำนวนข้อความใน session
- **WebSocket Status:** สถานะการเชื่อมต่อ (Connected/Disconnected)

```html
<div class="debug-stats">
  <div>Latency: {{ latency }}ms</div>
  <div>Messages: {{ messageCount }}</div>
  <div>Session: {{ sessionId }}</div>
  <div>WS: <span :class="wsStatus">{{ wsStatus }}</span></div>
  <button @click="clearHistory">Clear History</button>
</div>
```

### 5.4 Clear History Button
```javascript
async function clearHistory() {
  if (confirm('Delete all messages?')) {
    await axios.delete('/api/chats')
    messages.value = []
    sessionId.value = generateNewSessionId()
  }
}
```

---

## WebSocket Integration

### Connection Setup
```javascript
import { ref, onMounted, onUnmounted } from 'vue'

const ws = ref(null)
const messages = ref([])
const currentStreamingMessage = ref('')

function connectWebSocket() {
  ws.value = new WebSocket('ws://localhost:3001/api/chat/stream')

  ws.value.onopen = () => {
    console.log('WebSocket connected')
  }

  ws.value.onmessage = (event) => {
    const data = JSON.parse(event.data)

    if (data.type === 'start') {
      // เริ่ม streaming ข้อความใหม่
      currentStreamingMessage.value = ''
      messages.value.push({
        role: 'assistant',
        content: '',
        isStreaming: true
      })
    } else if (data.type === 'content') {
      // รับ chunk ของข้อความ
      currentStreamingMessage.value += data.delta
      messages.value[messages.value.length - 1].content = currentStreamingMessage.value
    } else if (data.type === 'done') {
      // จบการ streaming
      messages.value[messages.value.length - 1].isStreaming = false
    } else if (data.type === 'error') {
      // จัดการ error
      console.error('WebSocket error:', data.error)
    }
  }

  ws.value.onerror = (error) => {
    console.error('WebSocket error:', error)
  }

  ws.value.onclose = () => {
    console.log('WebSocket disconnected')
    // Auto-reconnect logic
    setTimeout(connectWebSocket, 3000)
  }
}

function sendMessage(message, fileIds = []) {
  const payload = {
    type: 'message',
    persona_id: selectedPersona.value.id,
    message: message,
    session_id: sessionId.value,
    use_history: true
  }

  if (fileIds.length > 0) {
    payload.file_ids = fileIds
  }

  ws.value.send(JSON.stringify(payload))

  // เพิ่มข้อความของ user ใน UI
  messages.value.push({
    role: 'user',
    content: message,
    timestamp: new Date().toISOString()
  })
}

onMounted(() => {
  connectWebSocket()
})

onUnmounted(() => {
  if (ws.value) {
    ws.value.close()
  }
})
```

---

## Axios Configuration

### Setup Axios Instance
```javascript
// src/services/api.js
import axios from 'axios'

const apiClient = axios.create({
  baseURL: 'http://localhost:3001/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
apiClient.interceptors.request.use(
  (config) => {
    // เพิ่ม loading state
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // Global error handling
    if (error.response) {
      console.error('API Error:', error.response.data)
    }
    return Promise.reject(error)
  }
)

export default apiClient
```

### API Service Functions
```javascript
// src/services/chatService.js
import apiClient from './api'

export const chatService = {
  // Get all personas
  async getPersonas() {
    const response = await apiClient.get('/personas')
    return response.data.personas
  },

  // Get persona by ID
  async getPersonaById(id) {
    const response = await apiClient.get(`/personas/${id}`)
    return response.data
  },

  // Send chat message (non-streaming)
  async sendMessage(personaId, message, sessionId, useHistory = true) {
    const response = await apiClient.post('/chat', {
      persona_id: personaId,
      message,
      session_id: sessionId,
      use_history: useHistory
    })
    return response.data
  },

  // Get chat history
  async getChatHistory(limit = 50, offset = 0) {
    const response = await apiClient.get('/chats', {
      params: { limit, offset }
    })
    return response.data
  },

  // Delete all messages
  async deleteAllMessages() {
    const response = await apiClient.delete('/chats')
    return response.data
  }
}

// src/services/fileService.js
export const fileService = {
  // Upload files
  async uploadFiles(files) {
    const formData = new FormData()
    files.forEach(file => {
      formData.append('files', file)
    })

    const response = await apiClient.post('/file/uploads', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    return response.data
  },

  // Get file history
  async getFileHistory(limit = 20, offset = 0) {
    const response = await apiClient.get('/file/history', {
      params: { limit, offset }
    })
    return response.data
  },

  // Delete all files
  async deleteAllFiles() {
    const response = await apiClient.delete('/file/uploads')
    return response.data
  }
}

// src/services/audioService.js
export const audioService = {
  // Speech-to-Text
  async transcribeAudio(audioBlob) {
    const formData = new FormData()
    formData.append('audio', audioBlob, 'recording.webm')

    const response = await apiClient.post('/audio/transcribe', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    return response.data
  },

  // Text-to-Speech
  async textToSpeech(text, voice = 'alloy') {
    const response = await apiClient.post('/audio/tts',
      { text, voice },
      { responseType: 'blob' }
    )
    return response.data
  }
}
```

---

## State Management (Pinia)

### Chat Store
```javascript
// src/stores/chatStore.js
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { chatService } from '@/services/chatService'

export const useChatStore = defineStore('chat', () => {
  const messages = ref([])
  const selectedPersona = ref(null)
  const personas = ref([])
  const sessionId = ref(generateSessionId())
  const isLoading = ref(false)

  async function loadPersonas() {
    isLoading.value = true
    try {
      personas.value = await chatService.getPersonas()
      selectedPersona.value = personas.value[0] // เลือก persona แรกเป็นค่าเริ่มต้น
    } catch (error) {
      console.error('Failed to load personas:', error)
    } finally {
      isLoading.value = false
    }
  }

  function addMessage(role, content) {
    messages.value.push({
      role,
      content,
      timestamp: new Date().toISOString()
    })
  }

  async function clearAllMessages() {
    try {
      await chatService.deleteAllMessages()
      messages.value = []
      sessionId.value = generateSessionId()
    } catch (error) {
      console.error('Failed to clear messages:', error)
    }
  }

  function generateSessionId() {
    return `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  }

  return {
    messages,
    selectedPersona,
    personas,
    sessionId,
    isLoading,
    loadPersonas,
    addMessage,
    clearAllMessages
  }
})
```

### File Store
```javascript
// src/stores/fileStore.js
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fileService } from '@/services/fileService'

export const useFileStore = defineStore('file', () => {
  const selectedFiles = ref([])
  const uploadedFiles = ref([])
  const isUploading = ref(false)
  const uploadProgress = ref(0)

  function addFiles(files) {
    const validFiles = Array.from(files).filter(file => {
      return file.size <= 10 * 1024 * 1024 // 10MB max
    })

    if (selectedFiles.value.length + validFiles.length > 5) {
      throw new Error('Maximum 5 files allowed')
    }

    selectedFiles.value.push(...validFiles)
  }

  function removeFile(index) {
    selectedFiles.value.splice(index, 1)
  }

  async function uploadFiles() {
    if (selectedFiles.value.length === 0) return []

    isUploading.value = true
    try {
      const response = await fileService.uploadFiles(selectedFiles.value)
      uploadedFiles.value = response.uploaded_files

      // Reset selected files
      selectedFiles.value = []

      // Return file IDs
      return response.uploaded_files.map(f => f.file_id)
    } catch (error) {
      console.error('Upload failed:', error)
      throw error
    } finally {
      isUploading.value = false
      uploadProgress.value = 0
    }
  }

  return {
    selectedFiles,
    uploadedFiles,
    isUploading,
    uploadProgress,
    addFiles,
    removeFile,
    uploadFiles
  }
})
```

---

## Project Structure (Recommended)

```
frontend/
├── src/
│   ├── components/
│   │   ├── chat/
│   │   │   ├── ChatBubble.vue
│   │   │   ├── MessageList.vue
│   │   │   └── MessageInput.vue
│   │   ├── file/
│   │   │   ├── FileUploadModal.vue
│   │   │   ├── FileDragDropArea.vue
│   │   │   └── FileList.vue
│   │   ├── audio/
│   │   │   ├── MicrophoneButton.vue
│   │   │   ├── AudioStatusIndicator.vue
│   │   │   └── TranscriptDisplay.vue
│   │   ├── settings/
│   │   │   ├── SettingsPanel.vue
│   │   │   ├── PersonaSelector.vue
│   │   │   └── DebugStats.vue
│   │   └── layout/
│   │       ├── AppLayout.vue
│   │       └── AppSidebar.vue
│   ├── stores/
│   │   ├── chatStore.js
│   │   ├── fileStore.js
│   │   └── audioStore.js
│   ├── services/
│   │   ├── api.js
│   │   ├── chatService.js
│   │   ├── fileService.js
│   │   ├── audioService.js
│   │   └── websocketService.js
│   ├── composables/
│   │   ├── useWebSocket.js
│   │   ├── useAudioRecorder.js
│   │   └── useFileUpload.js
│   ├── views/
│   │   ├── ChatView.vue
│   │   └── SettingsView.vue
│   ├── router/
│   │   └── index.js
│   ├── App.vue
│   └── main.js
├── package.json
└── vite.config.js
```

---

## Key Dependencies

```json
{
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.2.0",
    "pinia": "^2.1.0",
    "axios": "^1.6.0",
    "tailwindcss": "^3.4.0",
    "@vueuse/core": "^10.7.0",
    "marked": "^11.0.0",
    "highlight.js": "^11.9.0"
  }
}
```

---

## Critical Implementation Notes

### 1. WebSocket Reconnection
ต้องมี auto-reconnect logic เมื่อ connection ขาด:
```javascript
function setupWebSocket() {
  let reconnectAttempts = 0
  const maxReconnectAttempts = 5

  function connect() {
    ws = new WebSocket('ws://localhost:3001/api/chat/stream')

    ws.onclose = () => {
      if (reconnectAttempts < maxReconnectAttempts) {
        reconnectAttempts++
        setTimeout(connect, 2000 * reconnectAttempts)
      }
    }
  }

  connect()
}
```

### 2. Audio Cleanup
ต้อง cleanup audio resources เมื่อ component unmount:
```javascript
onUnmounted(() => {
  if (audioPlayer) {
    audioPlayer.pause()
    audioPlayer.src = ''
  }
  if (mediaRecorder && mediaRecorder.state !== 'inactive') {
    mediaRecorder.stop()
  }
})
```

### 3. File Upload Error Handling
จัดการ partial upload success:
```javascript
const response = await fileService.uploadFiles(files)

if (response.failed > 0) {
  // แสดง warning ว่าบางไฟล์ upload ไม่สำเร็จ
  showWarning(`${response.failed} files failed to upload`)
}

// ใช้เฉพาะ file_ids ที่ upload สำเร็จ
return response.uploaded_files.map(f => f.file_id)
```

### 4. Session Management
สร้าง session ID ใหม่เมื่อ:
- User clear history
- User refresh page (optional - ขึ้นอยู่กับ requirement)
- Error ใน WebSocket connection ที่ไม่สามารถ recover ได้

### 5. Persona Configuration
ใช้ configuration จาก persona ที่เลือก:
```javascript
const currentConfig = computed(() => ({
  temperature: selectedPersona.value.temperature,
  maxTokens: selectedPersona.value.max_tokens,
  model: selectedPersona.value.model,
  languageSetting: JSON.parse(selectedPersona.value.language_setting),
  guardrails: JSON.parse(selectedPersona.value.guardrails)
}))
```

---

## Testing Checklist

### Chat Features
- [ ] ส่งข้อความผ่าน WebSocket ได้
- [ ] รับข้อความแบบ streaming ได้
- [ ] เปลี่ยน persona แล้วคำตอบเปลี่ยนตาม tone/style
- [ ] Auto-scroll เมื่อมีข้อความใหม่
- [ ] แสดง loading state เมื่อรอ Bot ตอบ

### File Upload
- [ ] Upload ไฟล์เดียวได้
- [ ] Upload หลายไฟล์ (สูงสุด 5) ได้พร้อมกัน
- [ ] Drag-and-drop ไฟล์ได้
- [ ] แสดง upload progress
- [ ] จัดการ error เมื่อ upload fail
- [ ] ส่งข้อความพร้อมไฟล์ได้

### Speech-to-Speech
- [ ] บันทึกเสียงได้
- [ ] แปลงเสียงเป็นข้อความได้
- [ ] ส่งข้อความไปยัง AI ได้
- [ ] แปลงคำตอบเป็นเสียงได้
- [ ] เล่นเสียงให้ user ฟังได้
- [ ] หยุดเสียงกลางคันได้

### Settings & Debug
- [ ] เลือก persona ได้และใช้งานได้จริง
- [ ] แสดง latency ถูกต้อง
- [ ] Clear history ได้และสร้าง session ใหม่
- [ ] WebSocket reconnect อัตโนมัติเมื่อ disconnect

---

## Performance Optimization

1. **Lazy Loading Components:**
```javascript
const FileUploadModal = defineAsyncComponent(() =>
  import('./components/file/FileUploadModal.vue')
)
```

2. **Debounce User Input:**
```javascript
import { useDebounceFn } from '@vueuse/core'

const debouncedSend = useDebounceFn(() => {
  sendMessage()
}, 500)
```

3. **Virtual Scrolling สำหรับ Chat History:**
ใช้ library เช่น `vue-virtual-scroller` สำหรับ chat ที่มีข้อความเยอะ

4. **Audio Preloading:**
Preload TTS audio สำหรับ common responses

---

## Security Considerations

1. **Input Sanitization:**
```javascript
import DOMPurify from 'dompurify'

function sanitizeMessage(message) {
  return DOMPurify.sanitize(message)
}
```

2. **File Validation:**
- ตรวจสอบ file type ก่อน upload
- จำกัดขนาดไฟล์
- Scan for malicious content (ฝั่ง backend)

3. **WebSocket Authentication:**
ถ้าจำเป็น ส่ง token ผ่าน WebSocket connection

---

## Deployment

### Build for Production
```bash
npm run build
```

### Environment Variables
```env
VITE_API_BASE_URL=https://api.yourdomain.com
VITE_WS_URL=wss://api.yourdomain.com/api/chat/stream
```

### Nginx Configuration
```nginx
location /api {
    proxy_pass http://backend:3001;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
}
```

---

## Summary

**สิ่งที่ต้องมี:**
1. ✅ Vue 3 + Composition API
2. ✅ Axios สำหรับ HTTP requests
3. ✅ WebSocket client สำหรับ streaming chat
4. ✅ Pinia สำหรับ state management
5. ✅ MediaRecorder API สำหรับบันทึกเสียง
6. ✅ File upload with drag-and-drop
7. ✅ Persona selector (8 personas)
8. ✅ Debug panel แสดง stats
9. ✅ Audio player with stop control
10. ✅ Markdown rendering สำหรับ Bot response

**Flow หลัก:**
- **Text Chat:** User → WebSocket → AI → WebSocket → User
- **Speech:** User (voice) → STT → AI → TTS → User (audio)
- **File:** User → Upload → Get file_ids → Send with message → AI analyzes → Response

**Backend Endpoints ที่ใช้:**
- `WS /api/chat/stream` - Chat streaming
- `POST /api/file/uploads` - File upload
- `POST /api/audio/transcribe` - Speech-to-Text
- `POST /api/audio/tts` - Text-to-Speech
- `GET /api/personas` - Get personas list
- `DELETE /api/chats` - Clear history
