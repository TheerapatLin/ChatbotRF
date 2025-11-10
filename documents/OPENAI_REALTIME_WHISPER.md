# OpenAI Realtime Whisper Implementation Guide

**เวอร์ชัน**: 1.0
**วันที่สร้าง**: 2025-11-10
**สถานะ**: 📝 Planning Phase

---

## 📋 สารบัญ

1. [ภาพรวม](#ภาพรวม)
2. [หลักการทำงานของ OpenAI Realtime API](#หลักการทำงานของ-openai-realtime-api)
3. [เปรียบเทียบวิธีการ Transcription](#เปรียบเทียบวิธีการ-transcription)
4. [สถาปัตยกรรมระบบ](#สถาปัตยกรรมระบบ)
5. [รายละเอียดทางเทคนิค](#รายละเอียดทางเทคนิค)
6. [Implementation Tasks](#implementation-tasks)
7. [โครงสร้างไฟล์](#โครงสร้างไฟล์)
8. [การทดสอบ](#การทดสอบ)
9. [อ้างอิง](#อ้างอิง)

---

## ภาพรวม

OpenAI Realtime API เป็น WebSocket-based API ที่ออกแบบมาเพื่อการถอดความเสียง (Speech-to-Text Transcription) แบบเรียลไทม์ด้วยความหน่วงต่ำ (Low Latency: 300-800ms) เหมาะสำหรับ:

- **Live Captioning**: การสร้างคำบรรยายสดในการประชุม, งานสัมมนา
- **Voice Assistants**: ผู้ช่วยเสียงที่ตอบสนองทันที
- **Meeting Transcription**: บันทึกการประชุมแบบเรียลไทม์
- **Call Center Analytics**: วิเคราะห์การสนทนาลูกค้า
- **Interactive Voice Applications**: แอปพลิเคชันที่ต้องการปฏิสัมพันธ์ด้วยเสียงแบบทันที

### โมเดลที่รองรับ

| โมเดล | คุณสมบัติ | Latency | ความแม่นยำ | ราคา | Use Case |
|-------|-----------|---------|-----------|------|----------|
| **gpt-4o-transcribe** | Full-featured, state-of-the-art | ~500-800ms | สูงมาก (99%+) | สูงกว่า | Production, Meeting transcription, สำคัญ |
| **gpt-4o-mini-transcribe** | Lightweight, faster | ~300-500ms | สูง (95-98%) | ต่ำกว่า | Testing, Voice commands, Real-time captions |

---

## หลักการทำงานของ OpenAI Realtime API

### 1. Architecture Flow

```
┌─────────────┐         WebSocket          ┌─────────────────┐
│   Client    │ ←───────────────────────→  │  OpenAI API     │
│  (Browser/  │   wss://api.openai.com    │  Realtime       │
│   Mobile)   │   /v1/realtime            │  Transcription  │
└─────────────┘                            └─────────────────┘
      │                                              │
      │ 1. Audio Capture (Mic)                      │
      │ ─────────────────────────────────────────→ │
      │                                              │
      │ 2. Encode to PCM16 24kHz                    │
      │ ─────────────────────────────────────────→ │
      │                                              │
      │ 3. Base64 Encode & Send via WebSocket      │
      │ ─────────────────────────────────────────→ │
      │                                              │
      │                          4. Process Audio   │
      │                          ←──────────────── │
      │                                              │
      │ 5. Partial Transcription (Streaming)        │
      │ ←───────────────────────────────────────── │
      │   {"type": "partial", "text": "hello..."}   │
      │                                              │
      │ 6. Final Transcription                      │
      │ ←───────────────────────────────────────── │
      │   {"type": "final", "text": "hello world"}  │
      │                                              │
      │ 7. Display/Process Result                   │
      └─────────────────────────────────────────────┘
```

### 2. WebSocket Connection Lifecycle

#### Phase 1: Connection & Session Setup

```javascript
// 1. สร้าง WebSocket connection
const ws = new WebSocket(
  'wss://api.openai.com/v1/realtime?intent=transcription',
  {
    headers: {
      'Authorization': 'Bearer YOUR_OPENAI_API_KEY',
      'OpenAI-Beta': 'realtime=v1'
    }
  }
);

// 2. เมื่อเชื่อมต่อสำเร็จ ส่ง session configuration
ws.on('open', () => {
  ws.send(JSON.stringify({
    type: 'session.update',
    session: {
      model: 'gpt-4o-mini-transcribe',
      intent: 'transcription',
      input_audio_format: 'pcm16',
      output_audio_format: 'pcm16',
      input_audio_transcription: {
        model: 'gpt-4o-mini-transcribe',
        language: 'en',  // หรือ 'th' สำหรับภาษาไทย
      },
      turn_detection: {
        type: 'server_vad',
        threshold: 0.5,
        prefix_padding_ms: 300,
        silence_duration_ms: 500
      },
      input_audio_noise_reduction: 'near_field'  // 'near_field' หรือ 'far_field'
    }
  }));
});
```

**คำอธิบาย Parameters:**
- `intent: 'transcription'` - ระบุว่าต้องการใช้สำหรับถอดความเสียงเท่านั้น
- `input_audio_format: 'pcm16'` - รูปแบบเสียงที่รับ (PCM 16-bit)
- `language: 'en'` - ภาษาที่ต้องการถอดความ (auto-detect ถ้าไม่ระบุ)
- `turn_detection.type: 'server_vad'` - ใช้ Voice Activity Detection ฝั่ง server
- `turn_detection.silence_duration_ms: 500` - หยุดรับเสียงหลังเงียบ 500ms
- `input_audio_noise_reduction: 'near_field'` - ลดเสียงรบกวนสำหรับไมค์ใกล้

#### Phase 2: Audio Streaming

```javascript
// 3. จับเสียงจาก microphone
navigator.mediaDevices.getUserMedia({ audio: true })
  .then(stream => {
    const audioContext = new AudioContext({ sampleRate: 24000 });
    const source = audioContext.createMediaStreamSource(stream);
    const processor = audioContext.createScriptProcessor(4096, 1, 1);

    processor.onaudioprocess = (e) => {
      // Convert Float32Array to Int16Array (PCM16)
      const float32 = e.inputBuffer.getChannelData(0);
      const int16 = new Int16Array(float32.length);

      for (let i = 0; i < float32.length; i++) {
        const s = Math.max(-1, Math.min(1, float32[i]));
        int16[i] = s < 0 ? s * 0x8000 : s * 0x7FFF;
      }

      // Base64 encode และส่งผ่าน WebSocket
      const base64Audio = btoa(
        String.fromCharCode.apply(null, new Uint8Array(int16.buffer))
      );

      ws.send(JSON.stringify({
        type: 'input_audio_buffer.append',
        audio: base64Audio
      }));
    };

    source.connect(processor);
    processor.connect(audioContext.destination);
  });
```

#### Phase 3: Receiving Transcriptions

```javascript
// 4. รับผลการถอดความแบบ streaming
ws.on('message', (message) => {
  const event = JSON.parse(message);

  switch(event.type) {
    case 'input_audio_buffer.speech_started':
      console.log('🎤 User started speaking');
      break;

    case 'input_audio_buffer.speech_stopped':
      console.log('🔇 User stopped speaking');
      break;

    case 'conversation.item.input_audio_transcription.partial':
      // Partial transcription (อัพเดตระหว่างพูด)
      console.log('📝 Partial:', event.transcript);
      updateTranscript(event.transcript, false);
      break;

    case 'conversation.item.input_audio_transcription.final':
      // Final transcription (เสร็จสมบูรณ์)
      console.log('✅ Final:', event.transcript);
      updateTranscript(event.transcript, true);

      // Optional: Log probabilities for confidence
      if (event.logprobs) {
        const avgConfidence = calculateConfidence(event.logprobs);
        console.log('📊 Confidence:', avgConfidence);
      }
      break;

    case 'error':
      console.error('❌ Error:', event.error);
      break;
  }
});

function updateTranscript(text, isFinal) {
  const transcriptDiv = document.getElementById('transcript');
  if (isFinal) {
    transcriptDiv.innerHTML += `<p class="final">${text}</p>`;
  } else {
    // อัพเดต partial transcript (ตัวเอียง)
    const partial = transcriptDiv.querySelector('.partial');
    if (partial) {
      partial.textContent = text;
    } else {
      transcriptDiv.innerHTML += `<p class="partial"><em>${text}</em></p>`;
    }
  }
}
```

### 3. Event Types & Message Format

#### Client → Server Events

| Event Type | Description | Payload |
|-----------|-------------|---------|
| `session.update` | อัพเดต session configuration | `{ type, session: {...} }` |
| `input_audio_buffer.append` | ส่ง audio chunk | `{ type, audio: "base64..." }` |
| `input_audio_buffer.commit` | ยืนยันจบการส่งเสียง | `{ type }` |
| `input_audio_buffer.clear` | ล้าง buffer | `{ type }` |

#### Server → Client Events

| Event Type | Description | Data |
|-----------|-------------|------|
| `session.created` | Session สร้างสำเร็จ | Session ID, config |
| `input_audio_buffer.speech_started` | ตรวจจับเสียงพูดเริ่มต้น | Timestamp |
| `input_audio_buffer.speech_stopped` | ตรวจจับเสียงพูดหยุด | Timestamp, audio duration |
| `conversation.item.input_audio_transcription.partial` | ผล transcription แบบ partial | `{ transcript, item_id }` |
| `conversation.item.input_audio_transcription.final` | ผล transcription สุดท้าย | `{ transcript, item_id, logprobs }` |
| `error` | ข้อผิดพลาด | `{ error: {type, message} }` |

---

## เปรียบเทียบวิธีการ Transcription

### 1. OpenAI Whisper API (Current Implementation)

**ตำแหน่ง**: `backend/services/openai_service.go` → `TranscribeAudio()`

```go
// POST /api/audio/transcribe
func (s *OpenAIService) TranscribeAudio(file io.Reader, filename string) (*TranscriptionResponse, error) {
    req := openai.AudioRequest{
        Model:    openai.Whisper1,
        FilePath: filename,
        Reader:   file,
    }
    resp, err := s.client.CreateTranscription(ctx, req)
    // ...
}
```

**ลักษณะการทำงาน:**
- 📁 Upload ไฟล์เสียงทั้งไฟล์ (HTTP POST)
- ⏱️ รอจนเสียงบันทึกเสร็จ → อัพโหลด → ได้ผลทั้งหมดพร้อมกัน
- 🔄 Latency: 3-10 วินาที (ขึ้นอยู่กับความยาวไฟล์)

**Use Cases:**
- ✅ บันทึกเสียงเสร็จแล้ว (Recorded audio)
- ✅ ไฟล์เสียงที่มีอยู่แล้ว
- ✅ ไม่ต้องการ real-time feedback

### 2. OpenAI Realtime API (New - To Be Implemented)

```
WebSocket Connection → Streaming Audio → Partial + Final Results
```

**ลักษณะการทำงาน:**
- 🎤 Stream เสียงแบบ real-time ผ่าน WebSocket
- 📝 ได้ผลทันทีขณะกำลังพูด (Partial transcription)
- ⚡ Latency: 300-800ms

**Use Cases:**
- ✅ Live transcription (คำบรรยายสด)
- ✅ Voice commands ที่ต้องตอบสนองทันที
- ✅ Interactive conversations

### 3. Whisper.cpp (Local Implementation)

**ตำแหน่ง**: `backend/whisper/` (งานที่ 1-2 เสร็จแล้ว)

```bash
./backend/whisper/binary/linux/main -m model.bin -f audio.wav
```

**ลักษณะการทำงาน:**
- 💻 รันบน server เอง (ไม่เรียก API)
- 🔒 ข้อมูลไม่ออกนอก server
- ⚡ เร็วถ้ามี GPU

**Use Cases:**
- ✅ Privacy-sensitive data
- ✅ Offline deployment
- ✅ ลดค่าใช้จ่าย API

### สรุปเปรียบเทียบ

| Feature | Whisper API | Realtime API | Whisper.cpp |
|---------|-------------|--------------|-------------|
| **Connection** | HTTP POST | WebSocket | Local Process |
| **Latency** | 3-10s | 300-800ms | 2-5s (CPU), <1s (GPU) |
| **Streaming** | ❌ | ✅ | ❌ |
| **Partial Results** | ❌ | ✅ | ❌ |
| **File Size Limit** | 25 MB | Unlimited (stream) | Unlimited |
| **Cost** | $0.006/min | $0.10/min (4o), $0.01/min (mini) | ฟรี (infrastructure only) |
| **Privacy** | ❌ ส่งไป OpenAI | ❌ ส่งไป OpenAI | ✅ อยู่บน server |
| **Languages** | 99+ | 99+ | 99+ |
| **Setup Complexity** | ⭐ ง่าย | ⭐⭐ ปานกลาง | ⭐⭐⭐ ซับซ้อน |
| **Use Case** | Batch processing | Real-time UI | Privacy/Offline |

---

## สถาปัตยกรรมระบบ

### Current Architecture (Whisper API)

```
┌──────────┐   HTTP POST    ┌──────────┐   HTTP POST    ┌─────────────┐
│ Frontend │ ─────────────→ │ Backend  │ ─────────────→ │ OpenAI      │
│ (Vue.js) │                │ (Go)     │                │ Whisper API │
└──────────┘                └──────────┘                └─────────────┘
     │                            │                             │
     │ 1. Record audio            │                             │
     │ 2. Upload file             │                             │
     │                            │ 3. Forward to OpenAI        │
     │                            │                             │
     │                            │ 4. Wait for response        │
     │                            │ ←───────────────────────── │
     │ 5. Return transcript       │                             │
     │ ←───────────────────────── │                             │

⏱️ Total Latency: 3-10 seconds
```

### Proposed Architecture (Realtime API)

```
┌──────────┐   WebSocket    ┌──────────┐   WebSocket    ┌─────────────┐
│ Frontend │ ←────────────→ │ Backend  │ ←────────────→ │ OpenAI      │
│ (Vue.js) │  (ws://...)    │ (Go)     │  (wss://...)   │ Realtime    │
└──────────┘                └──────────┘                └─────────────┘
     │                            │                             │
     │ 1. Capture audio (mic)     │                             │
     │ ────────────────────────→ │                             │
     │ 2. Stream PCM chunks       │                             │
     │ ────────────────────────→ │ 3. Forward to OpenAI       │
     │                            │ ──────────────────────────→│
     │                            │                             │
     │                            │ 4. Partial transcript       │
     │ 5. Update UI (streaming)   │ ←───────────────────────── │
     │ ←───────────────────────── │                             │
     │ (ทันทีขณะกำลังพูด)         │                             │
     │                            │                             │
     │ 6. Final transcript         │                             │
     │ ←───────────────────────── │ ←───────────────────────── │

⚡ Latency per chunk: 300-800ms
```

### Hybrid Architecture (แนะนำ)

```
                        ┌──────────────────────────┐
                        │      Frontend (Vue.js)    │
                        │                          │
                        │  ┌─────────────────────┐ │
                        │  │ Audio Recorder      │ │
                        │  │ - MediaRecorder API │ │
                        │  └─────────────────────┘ │
                        │           │              │
                        │           ├──────────────┼──────┐
                        │           │              │      │
                        └───────────┼──────────────┘      │
                                    │                     │
                            (Switch Mode)                 │
                                    │                     │
                   ┌────────────────┴─────────┐          │
                   │                          │          │
        ┌──────────▼─────────┐    ┌──────────▼─────────┐│
        │  Batch Mode        │    │  Realtime Mode     ││
        │  (HTTP POST)       │    │  (WebSocket)       ││
        │                    │    │                    ││
        │  /api/audio/       │    │  /ws/realtime/     ││
        │  transcribe        │    │  transcribe        ││
        └────────┬───────────┘    └────────┬───────────┘│
                 │                         │            │
                 │  Backend (Go Fiber)     │            │
                 │                         │            │
        ┌────────▼─────────┐    ┌──────────▼───────────▼┐
        │ OpenAI Service   │    │ Realtime Service       │
        │ (openai_service) │    │ (realtime_service)     │
        └────────┬───────────┘    └────────┬───────────┘
                 │                         │
                 │  HTTP                   │  WebSocket
                 │                         │
        ┌────────▼─────────────────────────▼────────────┐
        │          OpenAI API                           │
        │                                               │
        │  ┌─────────────────┐  ┌──────────────────┐  │
        │  │ Whisper API     │  │ Realtime API     │  │
        │  │ (whisper-1)     │  │ (gpt-4o-         │  │
        │  │                 │  │  transcribe)     │  │
        │  └─────────────────┘  └──────────────────┘  │
        └────────────────────────────────────────────┘

Use Case Selection:
├─ Batch Mode: Recorded audio files, non-urgent transcription
└─ Realtime Mode: Live captions, voice commands, meetings
```

---

## รายละเอียดทางเทคนิค

### 1. Audio Format Requirements

#### Input Audio (Client → Server)

```
Format:      PCM16 (16-bit Pulse Code Modulation)
Sample Rate: 24000 Hz (24 kHz) - แนะนำ
             16000 Hz (16 kHz) - รองรับ
Channels:    1 (Mono)
Bit Depth:   16-bit
Byte Order:  Little-endian
Encoding:    Base64 (สำหรับส่งผ่าน WebSocket)
Chunk Size:  4096 samples (~170ms @ 24kHz)
             2048 samples (~128ms @ 16kHz)
```

#### Audio Processing Pipeline

```
[Microphone]
      ↓
[AudioContext (Browser API)]
  - Sample Rate: 24000 Hz
  - Format: Float32 (-1.0 to 1.0)
      ↓
[ScriptProcessor / AudioWorklet]
  - Convert Float32 → Int16
  - Formula: int16 = float32 * 32767
      ↓
[Base64 Encoding]
  - Convert Int16Array to Base64 string
  - Why Base64? WebSocket supports text/binary, Base64 ง่ายกว่า
      ↓
[WebSocket Send]
  - Type: 'input_audio_buffer.append'
  - Audio: "base64_encoded_pcm_data"
```

### 2. Go Implementation Details

#### WebSocket Libraries

```bash
# ติดตั้ง dependencies
go get github.com/gorilla/websocket      # WebSocket server
go get github.com/sashabaranov/go-openai # OpenAI client (มีอยู่แล้ว)
```

#### Backend Service Structure

```
backend/
├── services/
│   ├── openai_service.go           # มีอยู่แล้ว
│   └── realtime_service.go         # ใหม่ - WebSocket to OpenAI Realtime
│
├── controllers/
│   ├── audio_controller.go         # มีอยู่แล้ว
│   └── realtime_controller.go      # ใหม่ - WebSocket endpoint
│
├── models/
│   └── transcription.go            # ใหม่ - Realtime event types
│
└── routes/
    └── routes.go                   # อัพเดต - เพิ่ม WebSocket route
```

### 3. Environment Variables

```env
# เพิ่มใน .env.development

# OpenAI Realtime API
OPENAI_REALTIME_ENDPOINT=wss://api.openai.com/v1/realtime
OPENAI_REALTIME_MODEL=gpt-4o-mini-transcribe  # หรือ gpt-4o-transcribe
OPENAI_REALTIME_LANGUAGE=en                   # en, th, auto
OPENAI_REALTIME_VAD_THRESHOLD=0.5             # Voice Activity Detection threshold
OPENAI_REALTIME_SILENCE_MS=500                # Silence duration (ms)
OPENAI_REALTIME_NOISE_REDUCTION=near_field    # near_field, far_field
```

### 4. Security Considerations

```go
// CORS Configuration
app.Use(cors.New(cors.Config{
    AllowOrigins:     "http://localhost:5173",
    AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
    AllowCredentials: true,
    AllowWebSockets:  true,  // ⚠️ สำคัญ!
}))

// API Key ไม่ควร expose ไปยัง Frontend
// Backend เป็นตัวกลางเชื่อม OpenAI
```

### 5. Error Handling

| Error Code | Description | Solution |
|-----------|-------------|----------|
| `invalid_api_key` | API key ไม่ถูกต้อง | ตรวจสอบ OPENAI_API_KEY |
| `invalid_audio_format` | Format ไม่ถูกต้อง | ต้องเป็น PCM16 24kHz |
| `rate_limit_exceeded` | เกินโควต้า | Retry with exponential backoff |
| `connection_timeout` | WebSocket timeout | Reconnect with exponential backoff |
| `unsupported_language` | ภาษาไม่รองรับ | ใช้ 'auto' หรือเปลี่ยนภาษา |

---

## Implementation Tasks

### งานที่ 1: Backend - สร้าง Realtime Service ⏳

**วัตถุประสงค์**: สร้าง Go service สำหรับเชื่อมต่อ OpenAI Realtime API

**ไฟล์**: `backend/services/realtime_service.go`

**ขั้นตอน**:

1. **ติดตั้ง Dependencies**
```bash
cd backend
go get github.com/gorilla/websocket
```

2. **สร้าง RealtimeService struct**
```go
type RealtimeService struct {
    config        *config.Config
    openaiConn    *websocket.Conn
    mu            sync.Mutex
    sessionID     string
    isConnected   bool
}
```

3. **Implement Connection Methods**
- `ConnectToOpenAI()` - เชื่อมต่อไปยัง OpenAI Realtime API
- `ConfigureSession()` - ตั้งค่า session (model, language, VAD)
- `StreamAudio()` - ส่ง audio chunks ไปยัง OpenAI
- `ReceiveTranscription()` - รับผลการถอดความจาก OpenAI
- `Close()` - ปิดการเชื่อมต่อ

4. **Event Handling**
- จัดการ events ต่างๆ จาก OpenAI (partial, final, error)
- Implement retry logic สำหรับการเชื่อมต่อขาด

**Unit Test**: `backend/test/realtime/service_test.go`
```bash
go test -v ./test/realtime/
```

---

### งานที่ 2: Backend - สร้าง WebSocket Controller ⏳

**วัตถุประสงค์**: สร้าง WebSocket endpoint สำหรับ Frontend

**ไฟล์**: `backend/controllers/realtime_controller.go`

**ขั้นตอน**:

1. **สร้าง RealtimeController**
```go
type RealtimeController struct {
    realtimeService *services.RealtimeService
    upgrader        websocket.Upgrader
}
```

2. **Implement WebSocket Handler**
- `HandleRealtimeTranscription(c *fiber.Ctx)` - อัพเกรด HTTP → WebSocket
- อ่าน audio chunks จาก client
- ส่งต่อไปยัง RealtimeService
- Forward transcription results กลับไปยัง client

3. **Connection Management**
- จัดการหลาย WebSocket connections พร้อมกัน
- Cleanup เมื่อ client disconnect
- Implement heartbeat/ping-pong

**API Endpoint**:
```
WS /ws/realtime/transcribe
```

---

### งานที่ 3: Backend - อัพเดต Config และ Routes ⏳

**วัตถุประสงค์**: เพิ่ม configuration และ routing

**ไฟล์**:
- `backend/config/config.go`
- `backend/routes/routes.go`

**ขั้นตอน**:

1. **อัพเดต Config struct** (config.go)
```go
type Config struct {
    // ... existing fields ...

    // OpenAI Realtime
    OpenAIRealtimeEndpoint    string
    OpenAIRealtimeModel       string
    OpenAIRealtimeLanguage    string
    OpenAIRealtimeVADThreshold float64
    OpenAIRealtimeSilenceMs   int
    OpenAIRealtimeNoiseReduction string
}
```

2. **เพิ่ม WebSocket Route** (routes.go)
```go
func SetupRoutes(app *fiber.App, controllers *Controllers) {
    // ... existing routes ...

    // WebSocket Realtime Transcription
    app.Get("/ws/realtime/transcribe",
        controllers.Realtime.HandleRealtimeTranscription)
}
```

---

### งานที่ 4: Frontend - Audio Capture Component 🎤 ⏳

**วัตถุประสงค์**: สร้าง Vue component สำหรับจับเสียงและส่งผ่าน WebSocket

**ไฟล์**: `frontend/src/components/audio/RealtimeRecorder.vue`

**ขั้นตอน**:

1. **สร้าง Component Template**
```vue
<template>
  <div class="realtime-recorder">
    <button @click="startRecording" :disabled="isRecording">
      🎤 Start Recording
    </button>
    <button @click="stopRecording" :disabled="!isRecording">
      ⏹️ Stop
    </button>

    <div class="transcript">
      <div v-for="item in transcripts" :key="item.id">
        <p :class="{ partial: !item.isFinal }">
          {{ item.text }}
        </p>
      </div>
    </div>
  </div>
</template>
```

2. **Implement Audio Capture**
- ใช้ `navigator.mediaDevices.getUserMedia()`
- สร้าง `AudioContext` ที่ 24kHz
- Implement `ScriptProcessor` หรือ `AudioWorklet`
- Convert Float32 → Int16 → Base64

3. **WebSocket Integration**
- เชื่อมต่อกับ `ws://localhost:3000/ws/realtime/transcribe`
- ส่ง audio chunks แบบ streaming
- รับและแสดงผล partial/final transcripts

4. **State Management**
```javascript
const state = reactive({
  isRecording: false,
  isConnected: false,
  transcripts: [],
  ws: null,
  audioContext: null,
  mediaStream: null
});
```

---

### งานที่ 5: Frontend - Realtime Service/Composable 📡 ⏳

**วัตถุประสงค์**: สร้าง reusable service สำหรับ WebSocket communication

**ไฟล์**: `frontend/src/composables/useRealtimeTranscription.js`

**ขั้นตอน**:

1. **สร้าง Composable**
```javascript
export function useRealtimeTranscription() {
  const ws = ref(null);
  const transcripts = ref([]);
  const isConnected = ref(false);

  const connect = () => {
    ws.value = new WebSocket('ws://localhost:3000/ws/realtime/transcribe');
    // ...
  };

  const sendAudio = (audioData) => {
    if (ws.value && isConnected.value) {
      ws.value.send(JSON.stringify({
        type: 'audio_chunk',
        data: audioData  // Base64
      }));
    }
  };

  const disconnect = () => {
    ws.value?.close();
  };

  return { connect, sendAudio, disconnect, transcripts, isConnected };
}
```

2. **Event Handling**
- `onOpen`: ตั้งค่า isConnected = true
- `onMessage`: แยก partial/final transcripts
- `onError`: จัดการ error และ reconnect
- `onClose`: cleanup และ reconnect logic

---

### งานที่ 6: Integration Testing 🧪 ⏳

**วัตถุประสงค์**: ทดสอบ end-to-end flow

**ขั้นตอน**:

1. **Unit Tests** (Go)
```bash
# Test RealtimeService
go test -v ./services/realtime_service_test.go

# Test RealtimeController
go test -v ./controllers/realtime_controller_test.go
```

2. **Frontend Unit Tests** (Vue)
```bash
# Test RealtimeRecorder component
npm run test:unit

# Test useRealtimeTranscription composable
npm run test:unit -- RealtimeRecorder
```

3. **Integration Test**
- เปิด backend: `go run main.go`
- เปิด frontend: `npm run dev`
- ทดสอบการบันทึกเสียงและรับ transcription

4. **Performance Testing**
- วัด latency (time to first partial transcript)
- ทดสอบการเชื่อมต่อหลาย clients พร้อมกัน
- ทดสอบการ reconnect

---

### งานที่ 7: Documentation และ Deployment ✅ ⏳

**วัตถุประสงค์**: เอกสารและ deploy

**ขั้นตอน**:

1. **อัพเดตเอกสาร**
- README.md - เพิ่มคำแนะนำการใช้ Realtime API
- API.md - เอกสาร WebSocket endpoints
- TROUBLESHOOTING.md - วิธีแก้ปัญหา

2. **Environment Setup**
```bash
# Development
cp .env.development.example .env.development
# แก้ไข OPENAI_API_KEY

# Production
cp .env.production.example .env.production
```

3. **Docker Support** (Optional)
```dockerfile
# อัพเดต Dockerfile รองรับ WebSocket
# docker-compose.yml เพิ่ม environment variables
```

---

## โครงสร้างไฟล์

### Backend Files (New/Modified)

```
backend/
├── services/
│   ├── openai_service.go             # มีอยู่แล้ว - ไม่ต้องแก้
│   └── realtime_service.go           # ✨ ใหม่
│
├── controllers/
│   ├── audio_controller.go           # มีอยู่แล้ว - ไม่ต้องแก้
│   └── realtime_controller.go        # ✨ ใหม่
│
├── models/
│   └── realtime_events.go            # ✨ ใหม่ - Event types
│
├── config/
│   └── config.go                     # 🔧 แก้ไข - เพิ่ม Realtime config
│
├── routes/
│   └── routes.go                     # 🔧 แก้ไข - เพิ่ม WebSocket route
│
└── test/
    └── realtime/
        ├── service_test.go           # ✨ ใหม่
        └── controller_test.go        # ✨ ใหม่
```

### Frontend Files (New/Modified)

```
frontend/src/
├── components/
│   └── audio/
│       └── RealtimeRecorder.vue      # ✨ ใหม่
│
├── composables/
│   └── useRealtimeTranscription.js   # ✨ ใหม่
│
├── services/
│   └── realtimeService.js            # ✨ ใหม่ (Optional - ถ้าไม่ใช้ composable)
│
└── views/
    └── RealtimeTranscriptionView.vue # ✨ ใหม่ - Demo page
```

### Configuration Files

```
.env.development                       # 🔧 แก้ไข - เพิ่ม Realtime vars
.env.production                        # 🔧 แก้ไข - เพิ่ม Realtime vars
docker-compose.yml                     # 🔧 แก้ไข - เพิ่ม environment
```

---

## การทดสอบ

### 1. Manual Testing Checklist

**Backend Testing:**
- [ ] WebSocket server รันได้
- [ ] เชื่อมต่อกับ OpenAI Realtime API สำเร็จ
- [ ] ส่ง audio chunks ได้
- [ ] รับ partial transcriptions ได้
- [ ] รับ final transcriptions ได้
- [ ] จัดการ client disconnect ได้
- [ ] Retry mechanism ทำงาน
- [ ] Error handling ครบถ้วน

**Frontend Testing:**
- [ ] ขอ microphone permission ได้
- [ ] จับเสียงจาก mic ได้
- [ ] Audio conversion (Float32→Int16→Base64) ถูกต้อง
- [ ] WebSocket connection สำเร็จ
- [ ] ส่ง audio ได้แบบ streaming
- [ ] แสดง partial transcripts (real-time)
- [ ] แสดง final transcripts
- [ ] UI responsive และ smooth
- [ ] Handle errors gracefully
- [ ] Cleanup resources เมื่อหยุด

### 2. Unit Tests

**Go Backend:**
```bash
# Test Realtime Service
go test -v ./services/realtime_service_test.go -run TestConnectToOpenAI
go test -v ./services/realtime_service_test.go -run TestStreamAudio
go test -v ./services/realtime_service_test.go -run TestReceiveTranscription

# Test Controller
go test -v ./controllers/realtime_controller_test.go -run TestWebSocketUpgrade
go test -v ./controllers/realtime_controller_test.go -run TestHandleClientAudio

# Run all tests
go test -v ./...
```

**Vue Frontend:**
```bash
# Test component
npm run test:unit -- RealtimeRecorder.spec.js

# Test composable
npm run test:unit -- useRealtimeTranscription.spec.js

# Coverage report
npm run test:coverage
```

### 3. Integration Test

**Scenario 1: Single User Transcription**
```bash
# 1. Start backend
cd backend && go run main.go

# 2. Start frontend
cd frontend && npm run dev

# 3. Open browser http://localhost:5173
# 4. Navigate to /realtime-transcription
# 5. Click "Start Recording"
# 6. Speak: "Hello, this is a test"
# 7. Verify partial transcripts appear in real-time
# 8. Click "Stop Recording"
# 9. Verify final transcript is correct
```

**Scenario 2: Multiple Users**
```bash
# เปิดหลาย browser tabs/windows
# ทดสอบ concurrent connections (3-5 users)
# ตรวจสอบว่าไม่มี cross-talk ระหว่าง users
```

**Scenario 3: Connection Issues**
```bash
# 1. Start recording
# 2. ปิด backend ระหว่างบันทึก
# 3. เปิด backend กลับมา
# 4. ตรวจสอบว่า reconnect อัตโนมัติ
# 5. ตรวจสอบว่า transcription ทำงานต่อได้
```

### 4. Performance Testing

```bash
# ใช้ tool เช่น k6 หรือ Artillery
# Load test: 10-50 concurrent WebSocket connections
k6 run load-test-websocket.js

# Expected Results:
# - Latency: < 1s (95th percentile)
# - Connection success rate: > 99%
# - No memory leaks
# - CPU usage stable
```

---

## อ้างอิง

### Official Documentation

- [OpenAI Realtime API Documentation](https://platform.openai.com/docs/guides/realtime)
- [OpenAI Audio API (Whisper)](https://platform.openai.com/docs/guides/speech-to-text)
- [WebSocket RFC 6455](https://datatracker.ietf.org/doc/html/rfc6455)
- [MDN Web Audio API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Audio_API)

### Libraries & Tools

- [Gorilla WebSocket (Go)](https://github.com/gorilla/websocket)
- [go-openai (Go SDK)](https://github.com/sashabaranov/go-openai)
- [MediaRecorder API (Browser)](https://developer.mozilla.org/en-US/docs/Web/API/MediaRecorder)
- [Vue.js Composables](https://vuejs.org/guide/reusability/composables.html)

### Pricing (as of 2025)

| Model | Price | Input | Note |
|-------|-------|-------|------|
| gpt-4o-transcribe | $0.10/min | Audio stream | High accuracy |
| gpt-4o-mini-transcribe | $0.01/min | Audio stream | Fast, lower cost |
| whisper-1 (Batch) | $0.006/min | Audio file | Non-realtime |

**ตัวอย่างการคำนวณ:**
- Meeting 1 ชั่วโมง (60 นาที)
  - gpt-4o-transcribe: $6.00
  - gpt-4o-mini-transcribe: $0.60
  - whisper-1 (Batch): $0.36

### Related Documents

- [WHISPER_START.md](./WHISPER_START.md) - Whisper.cpp Implementation (Tasks 1-2 ✅)
- Backend API documentation: `/backend/README.md`
- Frontend components: `/frontend/src/components/audio/README.md`

---

## สรุป

### ✅ สิ่งที่มีอยู่แล้ว (Current State)

1. **OpenAI Whisper API Integration**
   - File: `backend/services/openai_service.go`
   - Endpoint: `POST /api/audio/transcribe`
   - ใช้สำหรับ batch transcription

2. **Whisper.cpp Local Implementation**
   - Binary: `backend/whisper/binary/`
   - Config: `.env.development` (WHISPER_* variables)
   - Tests: `backend/test/sst-whisper/`

3. **Frontend Audio Recording**
   - Component: `frontend/src/components/chat/ChatInput.vue` (มี audio recorder)
   - ยังไม่มี realtime streaming

### 🎯 สิ่งที่ต้องสร้าง (To Be Implemented)

1. **Backend Realtime Service** (`backend/services/realtime_service.go`)
2. **Backend Realtime Controller** (`backend/controllers/realtime_controller.go`)
3. **Frontend Realtime Recorder** (`frontend/src/components/audio/RealtimeRecorder.vue`)
4. **Frontend Composable** (`frontend/src/composables/useRealtimeTranscription.js`)
5. **Integration Tests**
6. **Documentation**

### 📊 Expected Timeline

| Task | Estimated Time | Priority |
|------|---------------|----------|
| งานที่ 1: Backend Service | 4-6 hours | 🔴 High |
| งานที่ 2: Backend Controller | 3-4 hours | 🔴 High |
| งานที่ 3: Config & Routes | 1-2 hours | 🔴 High |
| งานที่ 4: Frontend Component | 4-5 hours | 🟡 Medium |
| งานที่ 5: Frontend Composable | 2-3 hours | 🟡 Medium |
| งานที่ 6: Integration Tests | 3-4 hours | 🟢 Low |
| งานที่ 7: Documentation | 2-3 hours | 🟢 Low |
| **Total** | **19-27 hours** | |

---

**Created**: 2025-11-10
**Last Updated**: 2025-11-10
**Status**: 📝 Planning Phase - Ready for Implementation
