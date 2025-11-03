# Frontend Documentation - ChatBot AI

**Framework**: Vue 3 + Vite
**Version**: 2.1.0
**Last Updated**: 2025-11-02

---

## 📋 สารบัญ

1. [ภาพรวมโปรเจ็ค](#ภาพรวมโปรเจ็ค)
2. [คุณสมบัติหลัก](#คุณสมบัติหลัก)
3. [โครงสร้างโฟลเดอร์](#โครงสร้างโฟลเดอร์)
4. [Components](#components)
5. [State Management](#state-management)
6. [Services](#services)
7. [วิธีใช้งาน](#วิธีใช้งาน)
8. [การติดตั้งและรัน](#การติดตั้งและรัน)

---

## 🎯 ภาพรวมโปรเจ็ค

### เทคโนโลยีที่ใช้

| Technology | Version | Purpose |
|------------|---------|---------|
| Vue 3 | Latest | Framework หลัก (Composition API) |
| Vite | 5.2.0 | Build tool & Dev server |
| Pinia | 2.1.7 | State management |
| Vue Router | 4.3.0 | Routing |
| Axios | 1.6.8 | HTTP client |
| Tailwind CSS | 3.4.3 | Styling |
| Headless UI | 1.7.22 | UI components |

### Port & URL

- **Dev Server**: `http://localhost:5173`
- **API**: `http://localhost:3001/api`
- **WebSocket**: `ws://localhost:3001/api/chat/stream`

---

## ⭐ คุณสมบัติหลัก

### 1. Text Chat (แชทข้อความ)
- ส่งข้อความแบบ real-time ผ่าน WebSocket
- แนบไฟล์ได้สูงสุด 5 ไฟล์ต่อข้อความ
- รองรับไฟล์: PDF, Word, Excel, รูปภาพ, Code, JSON
- Streaming แบบ token-by-token
- บันทึกประวัติการสนทนา
- แสดง token usage

### 2. Speech-to-Speech (แชทด้วยเสียง)
- กดค้างเพื่อบันทึกเสียง
- แปลงเสียงเป็นข้อความอัตโนมัติ
- AI ตอบกลับเป็นเสียง
- Waveform visualization แสดงคลื่นเสียงแบบ real-time
- เล่นเสียงตอบกลับอัตโนมัติ

### 3. File Analysis (วิเคราะห์ไฟล์)
- อัปโหลดไฟล์เพื่อวิเคราะห์
- ประเภทการวิเคราะห์: summary, detail, qa, extract
- แสดง preview ไฟล์พร้อมไอคอน
- บันทึกไฟล์ลงในการสนทนา
- รองรับรูปภาพผ่าน Vision API

### 4. AI Personas (บุคลิก AI)
- เลือกบุคลิก AI ได้ 3 แบบ
- กำหนด system prompt เอง
- บันทึกบุคลิกที่เลือกใช้

### 5. UI Features
- ภาษาไทยทั้งระบบ
- Responsive design (มือถือ + Desktop)
- Auto-scroll ในการสนทนา
- แสดงสถานะการเชื่อมต่อ
- File attachment preview
- Error handling พร้อมข้อความแจ้งเตือน

---

## 📁 โครงสร้างโฟลเดอร์

```
frontend/
├── src/
│   ├── views/               # หน้าหลัก
│   │   ├── ChatView.vue    # แชทข้อความ
│   │   └── SpeechView.vue  # แชทเสียง
│   │
│   ├── components/
│   │   ├── layout/         # Layout components
│   │   ├── chat/           # แชทและข้อความ
│   │   ├── file/           # อัปโหลดและแสดงไฟล์
│   │   └── speech/         # บันทึกและแสดงเสียง
│   │
│   ├── store/              # Pinia stores
│   │   ├── chat.js         # สถานะแชท
│   │   ├── personas.js     # บุคลิก AI
│   │   └── ui.js           # UI state
│   │
│   ├── services/           # API services
│   │   ├── api.js          # Axios setup
│   │   ├── chatService.js  # Chat API
│   │   ├── fileService.js  # File API
│   │   ├── audioService.js # Audio API
│   │   └── personaService.js
│   │
│   ├── composables/        # Vue composables
│   │   ├── useWebSocket.js
│   │   ├── useAudioRecorder.js
│   │   └── useSpeechToSpeech.js
│   │
│   ├── utils/              # Utilities
│   │   ├── formatters.js   # Format time, file size
│   │   └── fileHelpers.js  # File icons, colors
│   │
│   ├── config/
│   │   └── constants.js    # Constants
│   │
│   └── router/
│       └── index.js        # Routes
│
├── vite.config.js
├── tailwind.config.js
└── package.json
```

---

## 🧩 Components

### Layout

**AppLayout.vue**
- Layout หลักของแอป
- มี Sidebar และพื้นที่เนื้อหา

**AppSidebar.vue**
- เลือกโหมด (Text / Speech)
- เลือกบุคลิก AI
- ตั้งค่า System Prompt
- ดูจำนวนข้อความ

### Chat Components

**MessageList.vue**
- แสดงรายการข้อความทั้งหมด
- Auto-scroll ไปข้อความล่าสุด
- แสดง loading และ streaming indicator

**MessageInput.vue**
- พิมพ์ข้อความ (Enter = ส่ง, Shift+Enter = บรรทัดใหม่)
- ปุ่มแนบไฟล์
- แสดง uploading indicator
- ส่งข้อความผ่าน WebSocket

**MessageBubble.vue**
- แสดงข้อความแต่ละข้อความ
- User: ฟองสีน้ำเงิน (ชิดขวา)
- AI: ฟองสีขาว (ชิดซ้าย) พร้อมอิโมจิ 🤖
- แสดงไฟล์แนบพร้อมไอคอน
- แสดง timestamp และ token count

### File Components

**FileUpload.vue**
- เลือกไฟล์ (สูงสุด 5 ไฟล์)
- Preview ไฟล์พร้อมขนาด
- ปุ่มลบไฟล์
- Validation: ขนาดไม่เกิน 20MB

**FileAnalysisCard.vue**
- แสดงข้อมูลไฟล์ที่วิเคราะห์
- ไอคอนตามประเภทไฟล์
- สรุปการวิเคราะห์ (expandable)
- Status badge

### Speech Components

**MicrophoneButton.vue**
- ปุ่มบันทึกเสียง 120x120px
- แสดงสถานะ: idle, recording, processing, playing
- Animation: pulse ขณะบันทึก, spinning ขณะประมวลผล

**WaveformVisualizer.vue**
- แสดงคลื่นเสียงแบบ real-time
- Canvas 600x200px
- สีแดงขณะบันทึก, สีน้ำเงินขณะเล่น
- ใช้ Web Audio API

---

## 💾 State Management (Pinia)

### Chat Store (`chat.js`)

**State:**
```javascript
{
  messages: [],              // ข้อความทั้งหมด
  sessionId: null,           // Session ID ปัจจุบัน
  uploadedFiles: [],         // ไฟล์ที่อัปโหลดชั่วคราว
  isLoading: false,
  isStreaming: false,
  streamingContent: '',      // ข้อความที่กำลัง stream
  currentPersonaId: 1,
  chatMode: 'text',          // 'text' | 'speech'
  systemPrompt: '',
  pagination: {...}
}
```

**Actions สำคัญ:**
- `createNewSession()` - สร้าง session ใหม่
- `initializeChat()` - โหลดประวัติ
- `uploadFiles(files)` - อัปโหลดไฟล์
- `addMessage(message)` - เพิ่มข้อความ
- `updateStreamingContent(chunk)` - อัปเดตขณะ stream
- `finishStreaming(data)` - เสร็จสิ้นการ stream
- `clearChat()` - ล้างแชท

### Personas Store (`personas.js`)

**State:**
```javascript
{
  personas: [],              // รายการบุคลิก AI
  currentPersonaId: 1,
  isLoading: false
}
```

**Actions:**
- `fetchPersonas()` - ดึงข้อมูลบุคลิก
- `setCurrentPersona(id)` - เลือกบุคลิก

### UI Store (`ui.js`)

**State:**
```javascript
{
  sidebarOpen: true
}
```

---

## 🔌 Services

### API Service (`api.js`)
- Axios instance พร้อม interceptors
- Timeout: 30 วินาที
- Auto-retry logic

### Chat Service (`chatService.js`)
```javascript
sendMessage(message, sessionId, personaId, options)
  // POST /chat

getChatHistory(sessionId, limit, offset)
  // GET /chat/history
```

### File Service (`fileService.js`)
```javascript
analyzeFile(file, options)
  // POST /file/analyze

getFileHistory(sessionId, limit, offset)
  // GET /file/history
```

### Audio Service (`audioService.js`)
```javascript
transcribeAudio(audioBlob)
  // POST /audio/transcribe
  // Returns: {text, language, ...}

textToSpeech(text, options)
  // POST /audio/tts
  // Returns: {audio_data (base64), ...}

base64ToBlob(base64)
createAudioURL(blob)
```

### Persona Service (`personaService.js`)
```javascript
getAllPersonas()
  // GET /personas

getPersonaById(id)
  // GET /personas/:id
```

---

## 🎮 วิธีใช้งาน

### Text Chat Mode

1. **เริ่มการสนทนา:**
   - พิมพ์ข้อความใน textarea
   - กด Enter เพื่อส่ง (Shift+Enter = บรรทัดใหม่)

2. **แนบไฟล์:**
   - คลิกปุ่ม 📎 (แนบไฟล์)
   - เลือกไฟล์ (สูงสุด 5 ไฟล์)
   - ไฟล์จะแสดง preview ด้านบน
   - ส่งข้อความพร้อมไฟล์

3. **เลือกบุคลิก AI:**
   - เปิด Sidebar (☰)
   - เลือกบุคลิกจาก dropdown
   - บุคลิกจะถูกบันทึกอัตโนมัติ

4. **กำหนด System Prompt:**
   - เปิด Sidebar
   - พิมพ์ในช่อง "System Prompt (Optional)"
   - Prompt จะถูกส่งไปกับทุกข้อความ

5. **เริ่มแชทใหม่:**
   - คลิกปุ่ม "แชทใหม่"
   - Session ID และข้อความจะถูกล้าง

### Speech-to-Speech Mode

1. **สลับโหมด:**
   - เปิด Sidebar
   - คลิก "Speech Mode"

2. **บันทึกเสียง:**
   - กดค้างปุ่มไมโครโฟน 🎤
   - พูดข้อความ
   - จะเห็นคลื่นเสียงเคลื่อนไหว
   - ปล่อยปุ่มเพื่อหยุดบันทึก

3. **ระบบจะทำงานอัตโนมัติ:**
   - แปลงเสียงเป็นข้อความ
   - ส่งไปยัง AI
   - แปลงคำตอบเป็นเสียง
   - เล่นเสียงอัตโนมัติ

4. **สถานะต่างๆ:**
   - **Idle** (พร้อมใช้งาน): ปุ่มสีเทา
   - **Recording** (กำลังบันทึก): ปุ่มสีแดง pulsing
   - **Processing** (กำลังประมวลผล): ปุ่มหมุน
   - **Playing** (กำลังเล่นเสียง): ปุ่มสีฟ้า

5. **รีเซ็ต:**
   - คลิกปุ่ม "Reset" เพื่อล้างสถานะ

---

## 🎨 Data Flow

### Text Chat Flow
```
1. User พิมพ์ข้อความ + แนบไฟล์ (ถ้ามี)
   ↓
2. FileUpload → uploadFiles() → POST /file/analyze
   ↓
3. MessageInput → addMessage() → สร้าง user message พร้อม file_attachments
   ↓
4. WebSocket.sendMessage() → ส่งข้อความ
   ↓
5. Backend → Streaming chunks
   ↓
6. useWebSocket → updateStreamingContent() → แสดงแบบ real-time
   ↓
7. Done chunk → finishStreaming() → บันทึก assistant message
```

### Speech-to-Speech Flow
```
1. กดค้างปุ่มไมโครโฟน
   ↓
2. useAudioRecorder → MediaStream → WaveformVisualizer
   ↓
3. ปล่อยปุ่ม → stopRecordingAndProcess()
   ↓
4. POST /audio/transcribe → {text: "..."}
   ↓
5. WebSocket → ส่งข้อความ → รอคำตอบ
   ↓
6. POST /audio/tts → {audio_data: base64}
   ↓
7. base64 → Blob → ObjectURL → Audio element
   ↓
8. playAudioResponse() → เล่นเสียงอัตโนมัติ
```

---

## 🛠️ การติดตั้งและรัน

### ติดตั้ง Dependencies

```bash
cd frontend
npm install
```

### รันในโหมด Development

```bash
npm run dev
```

เปิดเบราว์เซอร์ที่ `http://localhost:5173`

### Build Production

```bash
npm run build
```

ไฟล์ output จะอยู่ในโฟลเดอร์ `dist/`

### Preview Production Build

```bash
npm run preview
```

---

## ⚙️ Configuration

### Environment Variables

สร้างไฟล์ `.env` หรือ `.env.development`:

```env
VITE_API_BASE_URL=http://localhost:3001/api
VITE_WS_URL=ws://localhost:3001/api/chat/stream
```

### Constants (`src/config/constants.js`)

```javascript
// Timeouts
API_TIMEOUT = 30000ms          // API request timeout
WS_RECONNECT_DELAY = 3000ms    // WebSocket reconnect delay

// File Upload
MAX_FILE_SIZE = 20MB           // ขนาดไฟล์สูงสุด
MAX_FILES_PER_UPLOAD = 5       // จำนวนไฟล์สูงสุดต่อครั้ง

// Chat
CHAT_HISTORY_LIMIT = 50        // จำนวนข้อความที่โหลด
DEFAULT_PERSONA_ID = 1         // บุคลิกเริ่มต้น
```

---

## 🔑 Key Features Summary

### ✅ Completed Features

| Feature | Status | Description |
|---------|--------|-------------|
| Text Chat | ✅ | แชทข้อความแบบ real-time streaming |
| File Upload | ✅ | อัปโหลดและวิเคราะห์ไฟล์ (สูงสุด 5 ไฟล์) |
| Speech-to-Speech | ✅ | บันทึกเสียง → ตอบเป็นเสียง |
| AI Personas | ✅ | เลือกบุคลิก AI 3 แบบ |
| System Prompt | ✅ | กำหนด prompt เอง |
| Session Persistence | ✅ | บันทึก session ใน localStorage |
| Chat History | ✅ | โหลดประวัติการสนทนา |
| Waveform Visualization | ✅ | แสดงคลื่นเสียงแบบ real-time |
| Token Tracking | ✅ | นับ token ที่ใช้ |
| File Attachments | ✅ | แสดงไฟล์แนบพร้อมไอคอน |
| Error Handling | ✅ | แจ้งเตือนข้อผิดพลาด |
| WebSocket Reconnect | ✅ | เชื่อมต่อใหม่อัตโนมัติ |
| Responsive Design | ✅ | รองรับมือถือและ Desktop |
| Thai Language | ✅ | ภาษาไทยทั้งระบบ |

---

## 📊 Technical Highlights

1. **Real-time Communication**: WebSocket streaming แบบ token-by-token
2. **Audio Processing**: MediaRecorder API + Web Audio API
3. **State Persistence**: localStorage สำหรับ session & persona
4. **Reactive UI**: Vue 3 Composition API
5. **Component Reusability**: Stateless และ stateful components
6. **Clean Architecture**: แยก services, composables, และ utils
7. **Type Safety**: Props validation ทุก component
8. **Error Handling**: Try-catch พร้อม user feedback
9. **Performance**: Auto-scroll optimization, debouncing
10. **Accessibility**: Semantic HTML, ARIA labels

---

## 📈 Metrics

- **Total Components**: 10
- **Total Stores**: 3 (chat, personas, ui)
- **Total Services**: 5
- **Total Composables**: 3
- **Total Utility Functions**: 6+
- **Lines of Code**: ~3000+

---

## 🔄 Changelog

### Version 2.1.0 (2025-11-02) - **Current**
- ✅ ลบ `analysis_summary` จาก FileAttachment
- ✅ แสดง prompt ใน user message เมื่ออัปโหลดไฟล์
- ✅ ปรับปรุง File preview component
- ✅ Refactor code: ลด duplication, สร้าง utils
- ✅ อัปเดตเอกสารให้สั้นและชัดเจน

### Version 2.0.0 (2025-11-01)
- ✅ เพิ่มการแนบไฟล์ใน chat
- ✅ FileAttachment JSONB structure
- ✅ WebSocket streaming พร้อมไฟล์
- ✅ Session linking

### Version 1.2.0 (2025-10-31)
- Speech-to-Speech mode
- Waveform visualization
- Audio transcription & TTS

### Version 1.0.0 (2025-10-28)
- เปิดตัวครั้งแรก
- Text chat
- File upload
- AI personas

---

**Frontend URL**: `http://localhost:5173`
**Backend API**: `http://localhost:3001/api`
