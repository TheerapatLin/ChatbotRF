# คู่มือการติดตั้ง Whisper.cpp

## ภาพรวม

เอกสารนี้อธิบายการติดตั้งระบบ Speech-to-Text (STT) โดยใช้ whisper.cpp สำหรับโปรเจค ChatBot **พร้อมรองรับทั้งภาษาไทยและภาษาอังกฤษ**

### whisper.cpp คืออะไร?

whisper.cpp เป็น implementation ของ OpenAI Whisper ASR model ที่เขียนด้วย C++ เพื่อประสิทธิภาพสูง มีคุณสมบัติ:

- **ประมวลผลในเครื่อง**: ไม่ต้องพึ่งพา cloud สามารถประมวลผล audio ในเครื่องได้
- **รองรับหลายภาษา**: รองรับ 99+ ภาษา รวมถึง**ภาษาไทย**และ**ภาษาอังกฤษ**
- **รูปแบบ GGML**: โมเดลที่ optimize แล้วเพื่อ inference ที่มีประสิทธิภาพ
- **ความหน่วงต่ำ**: แปลงคำพูดเป็นข้อความได้รวดเร็ว ใช้ทรัพยากรน้อย
- **ไม่มีค่า API**: ดาวน์โหลดโมเดลครั้งเดียว ใช้งานได้ไม่จำกัด

## โครงสร้างโฟลเดอร์

ระบบ Speech-to-Text ทั้งหมดจะถูกสร้างภายใต้โฟลเดอร์ `backend/whisper`:

```
backend/
├── whisper/                       # โฟลเดอร์หลักของระบบ Whisper STT
│   ├── binary/                    # whisper.cpp binaries สำหรับแต่ละ OS
│   │   ├── linux/                # Binary สำหรับ Linux
│   │   │   └── main              # Executable binary สำหรับ Linux
│   │   ├── windows/              # Binary สำหรับ Windows
│   │   │   └── main.exe          # Executable binary สำหรับ Windows
│   │   └── macos/                # Binary สำหรับ macOS
│   │       └── main              # Executable binary สำหรับ macOS
│   │
│   ├── models/                   # GGML models (โมเดล AI สำหรับ transcription)
│   │   ├── ggml-small.bin       # โมเดล small (461 MB) - แนะนำสำหรับใช้งานทั่วไป
│   │   ├── ggml-medium.bin      # โมเดล medium (1.5 GB) - ความแม่นยำสูงขึ้น
│   │   └── ggml-large-v2.bin    # โมเดล large-v2 (2.9 GB) - ความแม่นยำสูงสุด (ตัวเลือก)
│   │
│   ├── temp/                     # โฟลเดอร์เก็บไฟล์ audio ชั่วคราว
│   │   └── .gitkeep             # ไฟล์ placeholder เพื่อให้ Git track โฟลเดอร์ว่าง
│   │
│   └── README.md                # คู่มือการใช้งานระบบ Whisper STT
│
├── services/                     # Go service layer (รวม services ทั้งหมด)
│   ├── transcription_service.go      # Interface definition สำหรับ STT services
│   ├── whispercpp_service.go         # Implementation ของ whisper.cpp service
│   └── ... (services อื่นๆ ที่มีอยู่แล้ว)
│
├── controllers/                  # Go controllers (รวม controllers ทั้งหมด)
│   ├── whispercpp_controller.go      # Controller สำหรับ whisper.cpp API
│   └── ... (controllers อื่นๆ ที่มีอยู่แล้ว)
│
└── test/
    └── sst-whisper/              # Unit tests สำหรับระบบ Whisper STT
        ├── testdata/            # ไฟล์ทดสอบ audio samples
        │   ├── audio/          # ตัวอย่างไฟล์เสียง
        │   │   ├── thai_short.wav        # ไฟล์เสียงภาษาไทยสั้น (3-5 วินาที)
        │   │   ├── thai_long.wav         # ไฟล์เสียงภาษาไทยยาว (10-30 วินาที)
        │   │   ├── english_short.wav     # ไฟล์เสียงภาษาอังกฤษสั้น
        │   │   ├── english_long.wav      # ไฟล์เสียงภาษาอังกฤษยาว
        │   │   └── mixed_language.wav    # ไฟล์เสียงผสมภาษาไทย-อังกฤษ
        │   └── .gitkeep
        │
        ├── setup_test.go        # Unit tests สำหรับการติดตั้งและ binary
        ├── config_test.go       # Unit tests สำหรับ configuration
        ├── service_test.go      # Unit tests สำหรับ WhisperCppService
        └── integration_test.go  # Integration tests แบบ end-to-end
```

### คำอธิบายโครงสร้างโฟลเดอร์และไฟล์

**โฟลเดอร์หลัก `backend/whisper/`**
- เป็นโฟลเดอร์รากของระบบ Speech-to-Text
- เก็บ binaries และ models ของ Whisper
- **ไม่มี** services และ controllers (ถูกย้ายไปรวมกับโฟลเดอร์หลักของ backend)

**`backend/whisper/binary/`**
- เก็บ whisper.cpp binary executables สำหรับแต่ละระบบปฏิบัติการ
- ระบบจะเลือก binary ที่เหมาะสมโดยอัตโนมัติตาม OS ที่รันอยู่
- Binary เหล่านี้คอมไพล์มาจาก whisper.cpp (C++ implementation)

**`backend/whisper/models/`**
- เก็บ GGML model files ที่ใช้สำหรับ transcription
- โมเดลเหล่านี้ download มาจาก Hugging Face
- ยิ่งโมเดลใหญ่ ความแม่นยำยิ่งสูง แต่ช้ากว่า

**`backend/whisper/temp/`**
- เก็บไฟล์ audio ชั่วคราวที่ upload จาก client
- ไฟล์จะถูกลบหลังจาก transcription เสร็จสิ้น
- ช่วยลดการใช้พื้นที่ disk

**`backend/services/`** (โฟลเดอร์หลักของ backend - รวม services ทั้งหมด)
- `transcription_service.go`: Interface สำหรับ STT services (รองรับการสลับ implementation)
- `whispercpp_service.go`: Implementation หลักที่รัน whisper.cpp binary และประมวลผล output
- ไฟล์อื่นๆ: Services ที่มีอยู่แล้วในโปรเจ็ค (เช่น chat service, bedrock service, etc.)

**`backend/controllers/`** (โฟลเดอร์หลักของ backend - รวม controllers ทั้งหมด)
- `whispercpp_controller.go`: HTTP handler สำหรับ REST API endpoints ของ Whisper STT
- จัดการ request/response, validation, และ error handling
- ไฟล์อื่นๆ: Controllers ที่มีอยู่แล้วในโปรเจ็ค (เช่น chat controller, message controller, etc.)

**`backend/test/sst-whisper/`**
- เก็บ unit tests และ integration tests ทั้งหมดของระบบ Whisper
- `testdata/audio/`: ไฟล์เสียงตัวอย่างสำหรับทดสอบในหลายภาษา
- `setup_test.go`: ทดสอบการมีอยู่และทำงานของ binary และ model
- `config_test.go`: ทดสอบการ load configuration และ environment variables
- `service_test.go`: ทดสอบ business logic ของ WhisperCppService
- `integration_test.go`: ทดสอบการทำงานแบบ end-to-end ด้วยไฟล์เสียงจริง

### สถาปัตยกรรมการผสานระบบ

```
┌─────────────────┐
│   Frontend      │
│   (Vue.js)      │
└────────┬────────┘
         │ HTTP POST
         │ multipart/form-data
         ▼
┌─────────────────────────────────────────┐
│   Backend (Go Fiber)                     │
│                                          │
│  ┌────────────────────────────────┐    │
│  │  WhisperCppController          │    │
│  │  POST /api/stt/whispercpp      │    │
│  └──────────┬─────────────────────┘    │
│             │                            │
│             ▼                            │
│  ┌────────────────────────────────┐    │
│  │  WhisperCppService             │    │
│  │  - ประมวลผลไฟล์ audio          │    │
│  │  - เรียก whisper.cpp binary    │    │
│  │  - จัดการภาษาไทย               │    │
│  │  - ส่งคืนผล transcription      │    │
│  └──────────┬─────────────────────┘    │
│             │                            │
└─────────────┼────────────────────────────┘
              │ exec.Command()
              ▼
┌─────────────────────────────────────────┐
│   whisper.cpp Binary                     │
│   - โหลด GGML model                     │
│   - แปลง audio เป็นข้อความ              │
│   - ส่งผลทาง stdout                     │
└─────────────────────────────────────────┘
```

## หลักการทำงานของ Whisper Speech-to-Text

### 1. ภาพรวมการทำงาน (High-Level Overview)

ระบบ Whisper Speech-to-Text ในโปรเจ็คนี้ทำงานผ่าน 3 ชั้น (3-tier architecture):

**Layer 1: Presentation Layer (Frontend - Vue.js)**
- ผู้ใช้บันทึกเสียงผ่าน Web browser (MediaRecorder API)
- แปลงข้อมูลเสียงเป็น Blob หรือ File object
- ส่งไฟล์เสียงไปยัง Backend ผ่าน HTTP POST (multipart/form-data)

**Layer 2: Application Layer (Backend - Go Fiber)**
- รับ HTTP request พร้อมไฟล์เสียง
- Validate ไฟล์ (ขนาด, ประเภท, รูปแบบ)
- บันทึกไฟล์ชั่วคราวใน `backend/whisper/temp/`
- เรียกใช้ whisper.cpp binary ผ่าน Go `exec.Command()`
- รอรับผลลัพธ์และประมวลผล output
- ส่งผลลัพธ์กลับไปยัง Frontend เป็น JSON

**Layer 3: Processing Layer (whisper.cpp Binary - C++)**
- โหลด GGML model ไปยัง memory
- ประมวลผลไฟล์เสียง (decode, resample, normalize)
- แปลงเสียงเป็นข้อความด้วย AI model
- Output ผลลัพธ์ทาง stdout

---

### 2. Data Flow (ขั้นตอนการทำงานแบบละเอียด)

#### ขั้นตอนที่ 1: การบันทึกเสียงจาก Frontend

```
[ผู้ใช้พูด] → [Microphone] → [MediaRecorder API] → [Audio Blob]
                                                        ↓
                                             [แปลงเป็น WAV/MP3]
                                                        ↓
                                             [FormData object]
```

**การทำงาน:**
- Frontend ใช้ `navigator.mediaDevices.getUserMedia()` เพื่อเข้าถึงไมโครโฟน
- ใช้ `MediaRecorder` API บันทึกเสียงเป็น chunks
- รวม chunks เป็น Blob object
- สร้าง FormData และแนบไฟล์เสียง

**ตัวอย่าง Code (Vue.js):**
```javascript
const mediaRecorder = new MediaRecorder(stream)
const audioChunks = []

mediaRecorder.ondataavailable = (event) => {
  audioChunks.push(event.data)
}

mediaRecorder.onstop = async () => {
  const audioBlob = new Blob(audioChunks, { type: 'audio/wav' })
  const formData = new FormData()
  formData.append('audio', audioBlob, 'recording.wav')
  formData.append('language', 'th')

  // ส่งไปยัง Backend
  await fetch('/api/stt/whispercpp', {
    method: 'POST',
    body: formData
  })
}
```

---

#### ขั้นตอนที่ 2: Controller รับ Request (WhisperCppController)

```
[HTTP POST Request] → [Fiber Router] → [WhisperCppController.TranscribeAudio()]
                                                        ↓
                                          [Validate Request & Extract File]
                                                        ↓
                                          [เรียก WhisperCppService.Transcribe()]
```

**การทำงาน:**
- Fiber framework route request ไปยัง `WhisperCppController`
- Controller ตรวจสอบ multipart/form-data
- Extract ไฟล์เสียงจาก form field "audio"
- ตรวจสอบพารามิเตอร์ (language, timestamps)
- Validate ไฟล์:
  - ขนาดไฟล์ไม่เกิน limit (เช่น 25 MB)
  - รูปแบบไฟล์ถูกต้อง (WAV, MP3, M4A, OGG, FLAC)
  - ไฟล์ไม่เสียหาย
- เรียก Service layer เพื่อประมวลผล

**ตัวอย่าง Code (Go):**
```go
func (ctrl *WhisperCppController) TranscribeAudio(c *fiber.Ctx) error {
    // รับไฟล์จาก form
    file, err := c.FormFile("audio")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "ไม่พบไฟล์เสียง",
        })
    }

    // รับพารามิเตอร์
    language := c.FormValue("language", "auto")

    // เปิดไฟล์
    audioFile, err := file.Open()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "ไม่สามารถเปิดไฟล์ได้",
        })
    }
    defer audioFile.Close()

    // เรียก Service
    transcription, confidence, err := ctrl.service.Transcribe(
        audioFile,
        file.Filename,
        language,
    )

    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "transcription": transcription,
        "confidence": confidence,
        "language": language,
    })
}
```

---

#### ขั้นตอนที่ 3: Service ประมวลผล (WhisperCppService)

```
[Transcribe() method called]
         ↓
[1. บันทึกไฟล์ชั่วคราวใน temp/]
         ↓
[2. กำหนด command arguments สำหรับ whisper.cpp]
         ↓
[3. เลือก binary path ตาม OS (Linux/Windows/macOS)]
         ↓
[4. รัน exec.Command() เพื่อเรียก whisper.cpp binary]
         ↓
[5. รอรับ output จาก stdout]
         ↓
[6. Parse output เพื่อดึง transcription text]
         ↓
[7. คำนวณ confidence score]
         ↓
[8. ลบไฟล์ชั่วคราว]
         ↓
[9. ส่งผลลัพธ์กลับ]
```

**รายละเอียดการทำงาน:**

**ขั้นตอน 1: บันทึกไฟล์ชั่วคราว**
```go
// สร้างชื่อไฟล์ชั่วคราว (unique)
tempFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), filename)
tempPath := filepath.Join(s.cfg.WhisperTempDir, tempFilename)

// เขียนข้อมูลลงไฟล์
tempFile, err := os.Create(tempPath)
if err != nil {
    return "", 0, err
}
defer os.Remove(tempPath) // ลบหลังเสร็จสิ้น

// Copy ข้อมูลจาก audioFile ไปยัง tempFile
_, err = io.Copy(tempFile, audioFile)
tempFile.Close()
```

**ขั้นตอน 2-4: กำหนด command และรัน binary**
```go
// เลือก binary ตาม OS
binaryPath := s.cfg.WhisperBinaryPath // ถูกกำหนดใน config

// กำหนด arguments
args := []string{
    "-m", s.cfg.WhisperModelPath,     // path ไปยัง GGML model
    "-f", tempPath,                    // path ไปยังไฟล์เสียงชั่วคราว
    "-l", language,                    // รหัสภาษา (th, en, auto)
    "-t", strconv.Itoa(s.cfg.WhisperThreads), // จำนวน CPU threads
    "-bs", strconv.Itoa(s.cfg.WhisperBeamSize), // beam size
    "-bo", strconv.Itoa(s.cfg.WhisperBestOf),   // best of N
}

// สร้าง command
cmd := exec.Command(binaryPath, args...)

// รัน command และรับ output
output, err := cmd.CombinedOutput()
if err != nil {
    return "", 0, fmt.Errorf("whisper.cpp error: %v", err)
}
```

**ขั้นตอน 5-7: Parse output และดึงข้อความ**
```go
// Output ตัวอย่างจาก whisper.cpp:
// [00:00:00.000 --> 00:00:03.000]  สวัสดีครับ
// [00:00:03.000 --> 00:00:06.000]  ยินดีต้อนรับสู่ระบบแชทบอท

outputStr := string(output)
lines := strings.Split(outputStr, "\n")

var transcriptionText string
for _, line := range lines {
    // ดึงเฉพาะบรรทัดที่มี timestamp
    if strings.Contains(line, "-->") {
        // ตัดส่วน timestamp ออก
        parts := strings.SplitN(line, "]", 2)
        if len(parts) == 2 {
            text := strings.TrimSpace(parts[1])
            transcriptionText += text + " "
        }
    }
}

transcriptionText = strings.TrimSpace(transcriptionText)

// คำนวณ confidence (ประมาณการจากความยาวข้อความ)
confidence := calculateConfidence(transcriptionText)

return transcriptionText, confidence, nil
```

---

#### ขั้นตอนที่ 4: whisper.cpp Binary ประมวลผล

```
[Binary เริ่มทำงาน]
         ↓
[1. โหลด GGML model ไปยัง RAM]
         ↓
[2. อ่านไฟล์เสียง และ decode format (WAV/MP3/etc.)]
         ↓
[3. Resample เสียงเป็น 16kHz mono]
         ↓
[4. แปลงเสียงเป็น mel-spectrogram features]
         ↓
[5. ป้อน features เข้า Whisper AI model]
         ↓
[6. Model ทำนาย text tokens (ทีละ token)]
         ↓
[7. Decode tokens เป็นข้อความ (support Unicode สำหรับภาษาไทย)]
         ↓
[8. แบ่ง segments ตาม timestamps]
         ↓
[9. Output ผลลัพธ์ทาง stdout]
```

**รายละเอียดภายใน whisper.cpp:**

**Audio Processing:**
- รองรับหลายรูปแบบ (WAV, MP3, M4A, OGG, FLAC)
- Automatic resampling เป็น 16kHz (ความถี่ที่ Whisper model ต้องการ)
- แปลง stereo → mono (ถ้าจำเป็น)
- Normalize audio levels

**Feature Extraction:**
- แปลง waveform เป็น mel-spectrogram
- Window size: 25ms
- Hop length: 10ms
- 80 mel filter banks

**AI Model Inference:**
- โหลด GGML quantized model (เพื่อลดการใช้ RAM)
- Encoder-Decoder transformer architecture
- Encoder: ประมวลผล audio features
- Decoder: สร้างข้อความทีละ token
- Beam search สำหรับหาลำดับคำที่ดีที่สุด

**Language Detection (ถ้าใช้ auto):**
- Model วิเคราะห์ audio features 30 วินาทีแรก
- คำนวณ probability สำหรับแต่ละภาษา
- เลือกภาษาที่มี probability สูงสุด

**Text Generation:**
- ใช้ tokenizer ที่รองรับ Unicode (ภาษาไทยใช้ได้)
- Beam search เพื่อหา N-best hypotheses
- Temperature sampling สำหรับ diversity
- Output พร้อม timestamps

---

### 3. การจัดการหลายภาษา (Multi-language Support)

**Whisper model รองรับ 99+ ภาษา** โดยใช้ multilingual training:

```
[Audio Input] → [Encoder] → [Language Detection]
                                    ↓
                    ┌───────────────┼───────────────┐
                    ↓               ↓               ↓
              [Thai Decoder]  [English Decoder]  [Auto Decoder]
                    ↓               ↓               ↓
              "สวัสดีครับ"     "Hello"         (ตรวจจับภาษาแล้วเลือก)
```

**การบังคับภาษา (language="th"):**
```go
// Service ส่ง parameter "-l th" ไปยัง whisper.cpp
args := []string{
    "-m", modelPath,
    "-f", audioPath,
    "-l", "th",  // บังคับให้ใช้ Thai decoder
}
```

**การตรวจจับอัตโนมัติ (language="auto"):**
```go
// ไม่ส่ง "-l" parameter หรือส่ง "auto"
args := []string{
    "-m", modelPath,
    "-f", audioPath,
    // whisper.cpp จะตรวจจับภาษาเอง
}
```

---

### 4. Performance Optimization

**การเลือก Model:**
- `small` (461 MB): รวดเร็ว, ความแม่นยำดี (แนะนำ)
- `medium` (1.5 GB): ช้ากว่า, ความแม่นยำสูงกว่า
- `large-v2` (2.9 GB): ช้าที่สุด, ความแม่นยำสูงสุด

**การใช้ Multi-threading:**
```env
WHISPER_THREADS=4  # ใช้ 4 CPU cores
```

**Beam Search Parameters:**
```env
WHISPER_BEAM_SIZE=5   # ยิ่งมาก ยิ่งแม่นยำ แต่ช้ากว่า
WHISPER_BEST_OF=5     # เลือกจาก 5 candidates ที่ดีที่สุด
```

**การ Cache Model:**
- Model ถูกโหลดครั้งแรกเมื่อ binary รัน
- ไม่ต้องโหลดใหม่ทุกครั้ง (ถ้าใช้ long-running process)

---

### 5. Error Handling & Recovery

**ประเภทของ Errors:**

1. **File Errors:**
   - ไฟล์ใหญ่เกินไป → Return 413 Payload Too Large
   - รูปแบบไฟล์ไม่รองรับ → Return 400 Bad Request
   - ไฟล์เสียหาย → Return 400 Bad Request

2. **Binary Errors:**
   - Binary ไม่พบ → Return 503 Service Unavailable
   - Binary ไม่มีสิทธิ์รัน → Return 500 Internal Server Error

3. **Model Errors:**
   - Model ไม่พบ → Return 503 Service Unavailable
   - Model format ผิด → Return 500 Internal Server Error

4. **Processing Errors:**
   - Audio decode ล้มเหลว → Return 400 Bad Request
   - Timeout (audio ยาวเกินไป) → Return 408 Request Timeout
   - Out of memory → Return 500 Internal Server Error

**Recovery Mechanisms:**
```go
// Cleanup ไฟล์ชั่วคราวเสมอ (ไม่ว่าจะสำเร็จหรือไม่)
defer os.Remove(tempPath)

// Timeout protection
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
cmd := exec.CommandContext(ctx, binaryPath, args...)

// Retry logic (สำหรับ transient errors)
for i := 0; i < maxRetries; i++ {
    output, err := cmd.CombinedOutput()
    if err == nil {
        break
    }
    time.Sleep(retryDelay)
}
```

---

### 6. Integration กับ Chatbot System

**Use Case 1: Voice Message ใน Chat**
```
[User records voice] → [STT converts to text] → [Send to LLM] → [Get response]
                                                                        ↓
                                                          [TTS converts to speech]
                                                                        ↓
                                                            [Play audio response]
```

**Use Case 2: Voice Commands**
```
[User: "สร้างรายการสินค้า"] → [STT] → [Parse command] → [Execute action]
```

**Use Case 3: Meeting Transcription**
```
[Long audio file] → [STT with timestamps] → [Save to database] → [Display transcript]
```

---

### 7. Security Considerations

**Input Validation:**
- ตรวจสอบขนาดไฟล์สูงสุด (ป้องกัน DoS)
- ตรวจสอบ MIME type และ file extension
- Sanitize filename (ป้องกัน path traversal)

**Resource Limits:**
- จำกัดจำนวน concurrent requests
- ตั้ง timeout สำหรับ processing
- จำกัดขนาด temp directory

**Data Privacy:**
- ลบไฟล์เสียงชั่วคราวทันทีหลัง transcription
- ไม่เก็บ log ของเนื้อหาเสียง (ถ้าเป็นข้อมูลละเอียดอ่อน)
- ใช้ HTTPS สำหรับการส่งข้อมูล

**Code Example:**
```go
// Validate file size
maxSize := int64(25 * 1024 * 1024) // 25 MB
if file.Size > maxSize {
    return c.Status(413).JSON(fiber.Map{
        "error": "ไฟล์ใหญ่เกินไป",
    })
}

// Validate file extension
allowedExts := []string{".wav", ".mp3", ".m4a", ".ogg", ".flac"}
ext := filepath.Ext(file.Filename)
if !contains(allowedExts, ext) {
    return c.Status(400).JSON(fiber.Map{
        "error": "รูปแบบไฟล์ไม่รองรับ",
    })
}

// Sanitize filename (ป้องกัน ../../../etc/passwd)
filename := filepath.Base(file.Filename)
filename = strings.ReplaceAll(filename, "..", "")
```

---

### 8. Monitoring & Logging

**Metrics ที่ควร Track:**
- จำนวน requests ต่อวินาที
- Average processing time
- Success rate vs Error rate
- Model accuracy (confidence scores)
- Audio duration distribution

**Logging Strategy:**
```go
log.Printf("[STT] Starting transcription: file=%s, size=%d, language=%s",
    filename, fileSize, language)

startTime := time.Now()
// ... process ...
duration := time.Since(startTime)

log.Printf("[STT] Completed: file=%s, duration=%dms, length=%d, confidence=%.2f",
    filename, duration.Milliseconds(), len(transcription), confidence)
```

---

### สรุปหลักการทำงานทั้งหมด

1. **Frontend** บันทึกเสียงและส่งไปยัง Backend ผ่าน HTTP POST
2. **Controller** รับ request, validate input, เรียก Service
3. **Service** บันทึกไฟล์ชั่วคราว, รัน whisper.cpp binary, parse output
4. **whisper.cpp** โหลด model, ประมวลผลเสียง, แปลงเป็นข้อความ, output ผลลัพธ์
5. **Response** ส่งข้อความที่แปลงได้กลับไปยัง Frontend
6. **Cleanup** ลบไฟล์ชั่วคราวทันที

ระบบนี้ออกแบบมาเพื่อ:
- **Scalability**: รองรับ multiple concurrent requests
- **Reliability**: Error handling และ cleanup ที่ดี
- **Performance**: ใช้ binary ที่ optimize แล้ว และ multi-threading
- **Maintainability**: แยก layer ชัดเจน, มี unit tests
- **Security**: Validation, resource limits, data cleanup

---

### การรองรับหลายภาษา (ภาษาไทย และ ภาษาอังกฤษ)

whisper.cpp รองรับการแปลงคำพูด**ทั้งภาษาไทยและภาษาอังกฤษ**ได้อย่างสมบูรณ์ พร้อมการตรวจจับภาษาอัตโนมัติ

**โมเดลแนะนำสำหรับภาษาไทยและอังกฤษ:**
- `small` (461 MB) - สมดุลที่ดีระหว่างความเร็วและความแม่นยำสำหรับทั้งสองภาษา
- `medium` (1.5 GB) - ความแม่นยำสูงขึ้นสำหรับคำพูดที่ซับซ้อนและผสมภาษา
- `large-v2` (2.9 GB) - ความแม่นยำสูงสุด (ช้ากว่า)

**URL สำหรับดาวน์โหลดโมเดล:**
```bash
# โมเดล Small (แนะนำสำหรับการใช้งานทั่วไป)
https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin

# โมเดล Medium (ความแม่นยำดีกว่า)
https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-medium.bin
```

**การตั้งค่าสำหรับแต่ละภาษา:**

```env
# สำหรับภาษาไทย
WHISPER_LANGUAGE=th           # บังคับตรวจจับภาษาไทย
WHISPER_MODEL_NAME=small
WHISPER_BEAM_SIZE=5           # เพิ่มความแม่นยำสำหรับเสียงวรรณยุกต์ไทย
WHISPER_BEST_OF=5

# สำหรับภาษาอังกฤษ
WHISPER_LANGUAGE=en           # บังคับตรวจจับภาษาอังกฤษ
WHISPER_MODEL_NAME=small
WHISPER_BEAM_SIZE=5
WHISPER_BEST_OF=5

# สำหรับการตรวจจับภาษาอัตโนมัติ (Auto-detect)
WHISPER_LANGUAGE=auto         # ตรวจจับภาษาอัตโนมัติ
WHISPER_MODEL_NAME=small
WHISPER_BEAM_SIZE=5
WHISPER_BEST_OF=5
```

---

## งานที่ต้องทำ (Implementation Tasks)

### งานที่ 1: ติดตั้งและตั้งค่า whisper.cpp

**วัตถุประสงค์**: ดาวน์โหลด คอมไพล์ และติดตั้ง whisper.cpp binary และโมเดล พร้อมรองรับ Linux, Windows, และ macOS

**ขั้นตอน**:

1. **สร้างโครงสร้างโฟลเดอร์**:
```bash
# สร้างโฟลเดอร์ whisper จาก root ของโปรเจค
mkdir -p backend/whisper/binary/linux backend/whisper/binary/windows backend/whisper/binary/macos
mkdir -p backend/whisper/models backend/whisper/temp
mkdir -p backend/test/sst-whisper/testdata/audio
```

2. **Clone repository whisper.cpp**:
```bash
cd backend/whisper
git clone https://github.com/ggerganov/whisper.cpp.git whisper-source
cd whisper-source
```

3. **คอมไพล์ whisper.cpp สำหรับแต่ละระบบปฏิบัติการ**:

#### สำหรับ Linux:
```bash
# ติดตั้ง dependencies (Ubuntu/Debian)
sudo apt-get update
sudo apt-get install build-essential

# คอมไพล์
make clean
make

# คัดลอก binary ไปยังโฟลเดอร์ที่กำหนด
cp main ../binary/linux/main
chmod +x ../binary/linux/main
```

#### สำหรับ Windows:
```bash
# ใช้ MinGW หรือ MSYS2
# ติดตั้ง MinGW-w64 ก่อน จาก https://www.msys2.org/
# จาก MSYS2 terminal:
pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-cmake make

# คอมไพล์
make clean
make

# หรือใช้ CMake
cmake -B build -G "MinGW Makefiles"
cmake --build build --config Release

# คัดลอก binary
cp main.exe ../binary/windows/main.exe
# หรือจาก build folder: cp build/bin/Release/main.exe ../binary/windows/main.exe
```

#### สำหรับ macOS:
```bash
# ติดตั้ง Xcode Command Line Tools
xcode-select --install

# คอมไพล์
make clean
make

# คัดลอก binary
cp main ../binary/macos/main
chmod +x ../binary/macos/main

# สำหรับ Apple Silicon (M1/M2) - Metal acceleration
# WHISPER_METAL=1 make
```

4. **ดาวน์โหลดโมเดล GGML**:
```bash
cd ../models

# Linux/macOS
curl -L -o ggml-small.bin https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin

# Windows PowerShell
# Invoke-WebRequest -Uri "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin" -OutFile "ggml-small.bin"
```

5. **ตรวจสอบการติดตั้งในแต่ละระบบ**:

#### Linux:
```bash
# ทดสอบภาษาไทย
../binary/linux/main -m ggml-small.bin -f ../test/audio/thai_short.wav -l th
# ทดสอบภาษาอังกฤษ
../binary/linux/main -m ggml-small.bin -f ../test/audio/english_short.wav -l en
```

#### Windows:
```bash
# ทดสอบภาษาไทย
..\binary\windows\main.exe -m ggml-small.bin -f ..\test\audio\thai_short.wav -l th
# ทดสอบภาษาอังกฤษ
..\binary\windows\main.exe -m ggml-small.bin -f ..\test\audio\english_short.wav -l en
```

#### macOS:
```bash
# ทดสอบภาษาไทย
../binary/macos/main -m ggml-small.bin -f ../test/audio/thai_short.wav -l th
# ทดสอบภาษาอังกฤษ
../binary/macos/main -m ggml-small.bin -f ../test/audio/english_short.wav -l en
```

**ผลลัพธ์ที่คาดหวัง**:
```
whisper_init_from_file: loaded model
whisper_model_load: n_vocab = 51864
whisper_model_load: n_audio_ctx = 1500
whisper_model_load: n_text_ctx = 448
...
[00:00:00.000 --> 00:00:11.000]  (ข้อความที่แปลงได้ตามภาษาที่ทดสอบ)
```

**การทดสอบด้วย Postman**: ไม่มี (การตั้งค่าผ่าน command-line เท่านั้น)

**Unit Test**: สร้างไฟล์ `backend/test/sst-whisper/setup_test.go`

**รันคำสั่ง Test**:
```bash
# Linux/macOS
cd backend/test/sst-whisper
go test -v -run TestWhisperBinaryExists
go test -v -run TestWhisperModelExists
go test -v -run TestWhisperVersion

# Windows
cd backend\test\sst-whisper
go test -v -run TestWhisperBinaryExists
go test -v -run TestWhisperModelExists
go test -v -run TestWhisperVersion

# รัน test ทั้งหมด
go test -v
```

```go
package whisper_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// getBinaryPath คืนค่า path ของ binary ตามระบบปฏิบัติการ
func getBinaryPath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join("..", "..", "whisper", "binary", "windows", "main.exe")
	case "darwin":
		return filepath.Join("..", "..", "whisper", "binary", "macos", "main")
	case "linux":
		return filepath.Join("..", "..", "whisper", "binary", "linux", "main")
	default:
		return filepath.Join("..", "..", "whisper", "binary", "linux", "main")
	}
}

// getModelPath คืนค่า path ของโมเดล
func getModelPath() string {
	return filepath.Join("..", "..", "whisper", "models", "ggml-small.bin")
}

// TestWhisperBinaryExists ตรวจสอบว่ามี whisper.cpp binary สำหรับระบบปฏิบัติการปัจจุบัน
func TestWhisperBinaryExists(t *testing.T) {
	binaryPath := getBinaryPath()

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatalf("ไม่พบ whisper.cpp binary ที่ %s (OS: %s)", binaryPath, runtime.GOOS)
	}

	t.Logf("✓ พบ whisper.cpp binary ที่ %s (OS: %s)", binaryPath, runtime.GOOS)
}

// TestWhisperModelExists ตรวจสอบว่ามีไฟล์โมเดล GGML
func TestWhisperModelExists(t *testing.T) {
	modelPath := getModelPath()

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Fatalf("ไม่พบโมเดล whisper ที่ %s", modelPath)
	}

	info, _ := os.Stat(modelPath)
	t.Logf("✓ พบโมเดล whisper ที่ %s (ขนาด: %d MB)", modelPath, info.Size()/1024/1024)
}

// TestWhisperVersion ตรวจสอบว่า whisper.cpp สามารถทำงานได้
func TestWhisperVersion(t *testing.T) {
	binaryPath := getBinaryPath()

	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("ไม่สามารถรัน whisper.cpp: %v (OS: %s)", err, runtime.GOOS)
	}

	if len(output) == 0 {
		t.Fatal("whisper.cpp ไม่มี output")
	}

	t.Logf("✓ whisper.cpp ทำงานได้ปกติ (OS: %s)", runtime.GOOS)
}

// TestWhisperTranscribeThaiAudio ทดสอบการแปลงเสียงภาษาไทย
func TestWhisperTranscribeThaiAudio(t *testing.T) {
	binaryPath := getBinaryPath()
	modelPath := getModelPath()
	audioPath := filepath.Join("testdata", "audio", "thai_short.wav")

	// ข้าม test ถ้าไม่มีไฟล์เสียงทดสอบ
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		t.Skipf("ข้าม test: ไม่พบไฟล์เสียงทดสอบ %s", audioPath)
	}

	cmd := exec.Command(binaryPath, "-m", modelPath, "-f", audioPath, "-l", "th")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("ไม่สามารถแปลงเสียงภาษาไทย: %v\nOutput: %s", err, string(output))
	}

	t.Logf("✓ แปลงเสียงภาษาไทยสำเร็จ\nOutput:\n%s", string(output))
}

// TestWhisperTranscribeEnglishAudio ทดสอบการแปลงเสียงภาษาอังกฤษ
func TestWhisperTranscribeEnglishAudio(t *testing.T) {
	binaryPath := getBinaryPath()
	modelPath := getModelPath()
	audioPath := filepath.Join("testdata", "audio", "english_short.wav")

	// ข้าม test ถ้าไม่มีไฟล์เสียงทดสอบ
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		t.Skipf("ข้าม test: ไม่พบไฟล์เสียงทดสอบ %s", audioPath)
	}

	cmd := exec.Command(binaryPath, "-m", modelPath, "-f", audioPath, "-l", "en")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("ไม่สามารถแปลงเสียงภาษาอังกฤษ: %v\nOutput: %s", err, string(output))
	}

	t.Logf("✓ แปลงเสียงภาษาอังกฤษสำเร็จ\nOutput:\n%s", string(output))
}
```

**✅ สถานะการดำเนินงาน (Updated: 2025-11-10)**

งานที่ 1 ดำเนินการเสร็จสิ้นแล้ว:

✅ **สำเร็จแล้ว:**
1. สร้างโครงสร้างโฟลเดอร์สำหรับ Whisper ครบถ้วน
   - `backend/whisper/binary/` (linux, windows, macos)
   - `backend/whisper/models/`
   - `backend/whisper/temp/`
   - `backend/test/sst-whisper/testdata/audio/`

2. Clone whisper.cpp source code
   - Repository: `backend/whisper/whisper-source/`
   - Version: Latest from GitHub

3. ดาวน์โหลด Binary สำหรับ Windows
   - Path: `backend/whisper/binary/windows/main.exe` (111 KB)
   - Source: Pre-compiled binary from whisper.cpp releases v1.5.4

4. ดาวน์โหลด GGML Model
   - Model: `backend/whisper/models/ggml-small.bin` (466 MB)
   - Source: Hugging Face (ggerganov/whisper.cpp)
   - รองรับ 99+ ภาษา รวมภาษาไทยและอังกฤษ

5. สร้าง README.md
   - Path: `backend/whisper/README.md`
   - รวมคำแนะนำการใช้งานและสถานะการติดตั้ง

6. **✅ สร้าง Unit Tests**
   - Path: `backend/test/sst-whisper/setup_test.go`
   - รวม 5 test functions:
     - `TestWhisperBinaryExists` - ตรวจสอบว่ามี binary ✅ PASS
     - `TestWhisperModelExists` - ตรวจสอบว่ามี model ✅ PASS
     - `TestWhisperVersion` - ตรวจสอบ binary ทำงานได้ ✅ PASS (WSL2)
     - `TestWhisperTranscribeThaiAudio` - ทดสอบแปลงเสียงภาษาไทย ✅ PASS (WSL2)
     - `TestWhisperTranscribeEnglishAudio` - ทดสอบแปลงเสียงภาษาอังกฤษ ✅ PASS (WSL2)

7. **✅ สร้างไฟล์เสียงทดสอบ**
   - `backend/test/sst-whisper/testdata/audio/th_audio.wav` (70 KB)
     - เสียงภาษาไทยจริง (~4.4 วินาที)
     - เนื้อหา: "ฉันเดินทางไปเที่ยวที่จังหวัดเชียงใหม่ในช่วงฤดูหนาวเพื่อสัมผัสอากาศเย็นสบาย"
   - `backend/test/sst-whisper/testdata/audio/en_audio.mp3` (276 KB)
     - เสียงภาษาอังกฤษ (~17.3 วินาที)
     - เนื้อหา: "What is this? This is a coffee. This is a cup..." (คำถาม-คำตอบเกี่ยวกับของใช้)

8. **✅ คอมไพล์ Binary สำหรับ Linux**
   - Path: `backend/whisper/binary/linux/main` (938 KB)
   - คอมไพล์ด้วย CMake บน WSL2
   - Version: whisper.cpp latest

9. **✅ สร้าง README.md และ FIX_ERROR.md**
   - Path: `backend/test/sst-whisper/README.md` - คำแนะนำการรัน tests
   - Path: `backend/test/sst-whisper/FIX_ERROR.md` - วิธีแก้ปัญหา DLL error และ troubleshooting

**📊 ผลการทดสอบ Unit Tests (Windows - DLL Error):**

```bash
# รัน test บน Windows
cd backend/test/sst-whisper
go test -v

# ผลลัพธ์:
✓ TestWhisperBinaryExists - PASS (0.00s)
✓ TestWhisperModelExists - PASS (0.00s)
✗ TestWhisperVersion - FAIL (exit status 0xc0000135 - ต้องการ Visual C++ Redistributable)
✗ TestWhisperTranscribeThaiAudio - FAIL (exit status 0xc0000135)
✗ TestWhisperTranscribeEnglishAudio - FAIL (exit status 0xc0000135)
```

**📊 ผลการทดสอบ Unit Tests (WSL2 - ✅ สำเร็จ):**

```bash
# รัน test บน WSL2
wsl bash -c "cd backend/test/sst-whisper && go test -v"

# ผลลัพธ์:
=== RUN   TestWhisperBinaryExists
    ✓ พบ whisper.cpp binary ที่ ../../whisper/binary/linux/main (OS: linux)
--- PASS: TestWhisperBinaryExists (0.00s)

=== RUN   TestWhisperModelExists
    ✓ พบโมเดล whisper ที่ ../../whisper/models/ggml-small.bin (ขนาด: 465 MB)
--- PASS: TestWhisperModelExists (0.00s)

=== RUN   TestWhisperVersion
    ✓ whisper.cpp ทำงานได้ปกติ (OS: linux)
--- PASS: TestWhisperVersion (0.11s)

=== RUN   TestWhisperTranscribeThaiAudio
    ✓ แปลงเสียงภาษาไทยสำเร็จ
    Output: "ฉันเดินทางไปเที่ยวที่จังวัดเชียงใหม่ในช่วงรูดูหนาว เพื่อสัมพัดอากาศเย็นสบาย"
--- PASS: TestWhisperTranscribeThaiAudio (12.73s)

=== RUN   TestWhisperTranscribeEnglishAudio
    ✓ แปลงเสียงภาษาอังกฤษสำเร็จ
    Output: "What is this? This is a coffee. This is a cup. This is a chair..."
--- PASS: TestWhisperTranscribeEnglishAudio (12.98s)

PASS
ok      chatbot/test/sst-whisper        25.821s
```

**🎯 วิธีการรัน Unit Tests:**

**วิธีที่ 1: ใช้ WSL2 (แนะนำ - ไม่มีปัญหา DLL)**
```bash
# รัน test ทั้งหมด
wsl bash -c "cd backend/test/sst-whisper && go test -v"

# รัน test เฉพาะเจาะจง
wsl bash -c "cd backend/test/sst-whisper && go test -v -run TestWhisperVersion"
wsl bash -c "cd backend/test/sst-whisper && go test -v -run TestWhisperTranscribeThaiAudio"
wsl bash -c "cd backend/test/sst-whisper && go test -v -run TestWhisperTranscribeEnglishAudio"

# รัน test แบบกลุ่ม (Version + Transcription tests)
wsl bash -c "cd backend/test/sst-whisper && go test -v -run 'TestWhisper(Version|Transcribe)'"

# รัน test พร้อม timeout (2 นาที)
wsl bash -c "cd backend/test/sst-whisper && go test -v -timeout 2m"
```

**วิธีที่ 2: ใช้ Windows (ต้องติดตั้ง Visual C++ Redistributable ก่อน)**
```bash
# รัน test ทั้งหมด
cd backend/test/sst-whisper
go test -v

# รัน test เฉพาะเจาะจง
go test -v -run TestWhisperBinaryExists
go test -v -run TestWhisperModelExists

# ดู test coverage
go test -v -cover
```

**📝 หมายเหตุ:**
- WSL2 tests ใช้ binary: `backend/whisper/binary/linux/main` (938 KB)
- Windows tests ใช้ binary: `backend/whisper/binary/windows/main.exe` (111 KB)
- Transcription tests ใช้เวลา ~12-13 วินาทีต่อไฟล์ (CPU mode, ไม่มี GPU acceleration)
- Model loading ใช้เวลา ~2 วินาที (487 MB)

**📝 คำอธิบาย Unit Tests:**

1. **TestWhisperBinaryExists**: ตรวจสอบว่าไฟล์ binary อยู่ในตำแหน่งที่ถูกต้องสำหรับ OS ปัจจุบัน
   - Windows: `backend/whisper/binary/windows/main.exe`
   - Linux: `backend/whisper/binary/linux/main`
   - macOS: `backend/whisper/binary/macos/main`

2. **TestWhisperModelExists**: ตรวจสอบว่ามีไฟล์ model และแสดงขนาด (ควรเป็น 465-466 MB)

3. **TestWhisperVersion**: พยายามรัน binary ด้วย `--help` เพื่อตรวจสอบว่าทำงานได้
   - ⚠️ บน Windows อาจ fail ถ้าไม่มี Visual C++ Redistributable

4. **TestWhisperTranscribeThaiAudio/English**: ทดสอบการแปลงเสียงจริง
   - ต้องมีไฟล์เสียงใน `testdata/audio/`
   - จะ skip ถ้าไม่มีไฟล์ทดสอบ

**⚠️ การแก้ปัญหา Windows DLL Error (0xc0000135):**

ปัญหา: Binary ต้องการ Visual C++ Redistributable ซึ่งมักจะไม่มีในระบบ Windows ใหม่

**วิธีแก้:**
1. ดาวน์โหลดและติดตั้ง [Visual C++ Redistributable (x64)](https://aka.ms/vs/17/release/vc_redist.x64.exe)
2. Restart เครื่อง (หรือ restart terminal)
3. รัน tests อีกครั้ง: `go test -v`

**ทางเลือกอื่น:**
- คอมไพล์ binary บน Windows โดยตรงโดยใช้ MinGW/MSYS2 (ตามขั้นตอนใน section Windows ด้านบน)
- ใช้ WSL2 และคอมไพล์สำหรับ Linux แทน
- ใช้ Docker container สำหรับการทดสอบ

**📝 หมายเหตุสำคัญ:**
- ✅ มีไฟล์เสียงภาษาไทยจริง (`th_audio.wav`) และภาษาอังกฤษ (`en_audio.mp3`)
- ✅ Tests ผ่านทั้งหมดบน WSL2 (Linux environment)
- ⚠️ Windows tests ยังคงต้องการ Visual C++ Redistributable
- แนะนำใช้ WAV format 16kHz mono สำหรับประสิทธิภาพที่ดีที่สุด

**🔗 เอกสารเพิ่มเติม:**
- รายละเอียด Unit Tests: [backend/test/sst-whisper/README.md](../backend/test/sst-whisper/README.md)
- วิธีแก้ปัญหา Error: [backend/test/sst-whisper/FIX_ERROR.md](../backend/test/sst-whisper/FIX_ERROR.md)
- Whisper Setup Guide: [backend/whisper/README.md](../backend/whisper/README.md)
- GitHub: https://github.com/ggerganov/whisper.cpp
- Models: https://huggingface.co/ggerganov/whisper.cpp

✅ **เสร็จสมบูรณ์:**
- ✅ Binary สำหรับ Linux (คอมไพล์บน WSL2) - 938 KB
- ✅ Binary สำหรับ Windows (pre-compiled) - 111 KB
- ✅ Model GGML Small - 465 MB
- ✅ ไฟล์เสียงทดสอบ (ภาษาไทยและอังกฤษ)
- ✅ Unit Tests ทั้งหมด (5 tests) - PASS บน WSL2
- ✅ Documentation และ Troubleshooting Guide

⚠️ **ยังไม่เสร็จ:**
- Binary สำหรับ macOS (ต้องคอมไพล์บน macOS)
- การรัน Binary บน Windows (ต้องติดตั้ง Visual C++ Redistributable หรือคอมไพล์ใหม่)

**📊 Performance Metrics (WSL2 - CPU Mode):**
- Model Loading Time: ~2 วินาที
- Thai Audio Transcription (4.4s): ~12.7 วินาที
- English Audio Transcription (17.3s): ~13.0 วินาที
- Total Test Suite: ~26 วินาที

**📝 หมายเหตุ:**
- Binary สำหรับ Windows ต้องการ **Visual C++ Redistributable** เพื่อรัน
- สำหรับ production ควรคอมไพล์ binary บนแต่ละระบบปฏิบัติการโดยตรง
- ขั้นตอนถัดไป: ตั้งค่า Environment Configuration และสร้าง Services/Controllers

---

### งานที่ 2: เพิ่ม Environment Configuration ✅

**วัตถุประสงค์**: ตั้งค่า path และ parameter ของ whisper.cpp ใน environment files

**สถานะ**: ✅ เสร็จสมบูรณ์

**ความคืบหน้า**:
- ✅ สร้างไฟล์ `.env.development` พร้อม Whisper configuration ครบถ้วน
- ✅ อัพเดต `backend/config/config.go` เพิ่ม struct fields สำหรับ Whisper
- ✅ เพิ่ม helper functions: `getWhisperBinaryPath()`, `getEnvAsInt()`, `getEnvAsBool()`
- ✅ อัพเดต `LoadConfig()` function เพื่อโหลดการตั้งค่า Whisper
- ✅ ทดสอบการโหลด configuration สำเร็จด้วย `backend/test/test-config.go`

**📊 ผลการทดสอบ:**
```bash
=== Testing Whisper Configuration ===
Current OS: windows

Whisper.cpp Configuration:
  Binary Path:        ./backend/whisper/binary/windows/main.exe
  Model Path:         ./backend/whisper/models/ggml-small.bin
  Temp Directory:     ./backend/whisper/temp
  Language:           th
  Model Name:         small
  Threads:            4
  Processors:         1
  Max Length:         0
  Beam Size:          5
  Best Of:            5
  Word Timestamps:    false
  Supported Languages: th,en,auto

✓ Configuration loaded successfully (env: development)
✓ Configuration test completed successfully
```

**ขั้นตอน**:

1. **อัปเดต `.env.development`**:
```env
# ========================================
# การตั้งค่า Whisper.cpp
# ========================================
# Binary paths สำหรับแต่ละ OS (ระบบจะเลือกอัตโนมัติตาม OS)
WHISPER_BINARY_PATH_LINUX=./backend/whisper/binary/linux/main
WHISPER_BINARY_PATH_WINDOWS=./backend/whisper/binary/windows/main.exe
WHISPER_BINARY_PATH_MACOS=./backend/whisper/binary/macos/main

# Model และการตั้งค่าทั่วไป
WHISPER_MODEL_PATH=./backend/whisper/models/ggml-small.bin
WHISPER_TEMP_DIR=./backend/whisper/temp

# การตั้งค่าภาษา (th=ไทย, en=อังกฤษ, auto=ตรวจจับอัตโนมัติ)
WHISPER_LANGUAGE=auto                  # auto, th, en
WHISPER_MODEL_NAME=small               # small, medium, large-v2
WHISPER_THREADS=4                      # จำนวน thread ของ CPU ที่ใช้
WHISPER_PROCESSORS=1                   # จำนวน processor
WHISPER_MAX_LEN=0                      # ความยาว segment สูงสุด (0 = ไม่จำกัด)
WHISPER_BEAM_SIZE=5                    # Beam size สำหรับความแม่นยำ
WHISPER_BEST_OF=5                      # Best of N candidates
WHISPER_WORD_TIMESTAMPS=false          # Timestamps ระดับคำ

# รองรับภาษา
WHISPER_SUPPORTED_LANGUAGES=th,en,auto # รายการภาษาที่รองรับ
```

2. **อัปเดต `backend/config/config.go`**:
```go
package config

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// ... ฟิลด์อื่นๆ ที่มีอยู่ ...

	// การตั้งค่า Whisper.cpp
	WhisperBinaryPath        string
	WhisperModelPath         string
	WhisperTempDir           string
	WhisperLanguage          string
	WhisperModelName         string
	WhisperThreads           int
	WhisperProcessors        int
	WhisperMaxLen            int
	WhisperBeamSize          int
	WhisperBestOf            int
	WhisperWordTimestamps    bool
	WhisperSupportedLanguages []string
}

// getWhisperBinaryPath คืนค่า binary path ตามระบบปฏิบัติการ
func getWhisperBinaryPath() string {
	// ตรวจสอบ OS และคืนค่า path ที่เหมาะสม
	switch runtime.GOOS {
	case "windows":
		return getEnv("WHISPER_BINARY_PATH_WINDOWS", "./backend/whisper/binary/windows/main.exe")
	case "darwin":
		return getEnv("WHISPER_BINARY_PATH_MACOS", "./backend/whisper/binary/macos/main")
	case "linux":
		return getEnv("WHISPER_BINARY_PATH_LINUX", "./backend/whisper/binary/linux/main")
	default:
		return getEnv("WHISPER_BINARY_PATH_LINUX", "./backend/whisper/binary/linux/main")
	}
}

// LoadConfig โหลดการตั้งค่าจาก environment variables
func LoadConfig() *Config {
	// โหลดไฟล์ .env
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	envFile := ".env." + appEnv
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("คำเตือน: ไม่พบไฟล์ %s กำลังใช้ environment variables ของระบบ", envFile)
	}

	// รายการภาษาที่รองรับ
	supportedLangs := strings.Split(
		getEnv("WHISPER_SUPPORTED_LANGUAGES", "th,en,auto"),
		",",
	)

	return &Config{
		// ... ฟิลด์อื่นๆ ที่มีอยู่ ...

		// การตั้งค่า Whisper.cpp
		WhisperBinaryPath:        getWhisperBinaryPath(),
		WhisperModelPath:         getEnv("WHISPER_MODEL_PATH", "./backend/whisper/models/ggml-small.bin"),
		WhisperTempDir:           getEnv("WHISPER_TEMP_DIR", "./backend/whisper/temp"),
		WhisperLanguage:          getEnv("WHISPER_LANGUAGE", "auto"),
		WhisperModelName:         getEnv("WHISPER_MODEL_NAME", "small"),
		WhisperThreads:           getEnvAsInt("WHISPER_THREADS", 4),
		WhisperProcessors:        getEnvAsInt("WHISPER_PROCESSORS", 1),
		WhisperMaxLen:            getEnvAsInt("WHISPER_MAX_LEN", 0),
		WhisperBeamSize:          getEnvAsInt("WHISPER_BEAM_SIZE", 5),
		WhisperBestOf:            getEnvAsInt("WHISPER_BEST_OF", 5),
		WhisperWordTimestamps:    getEnvAsBool("WHISPER_WORD_TIMESTAMPS", false),
		WhisperSupportedLanguages: supportedLangs,
	}
}

// getEnv ดึงค่า environment variable หรือคืนค่า default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt ดึงค่า environment variable เป็น integer หรือคืนค่า default
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("คำเตือน: ค่า integer ไม่ถูกต้องสำหรับ %s: %s ใช้ค่า default: %d", key, valueStr, defaultValue)
		return defaultValue
	}

	return value
}

// getEnvAsBool ดึงค่า environment variable เป็น boolean หรือคืนค่า default
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("คำเตือน: ค่า boolean ไม่ถูกต้องสำหรับ %s: %s ใช้ค่า default: %t", key, valueStr, defaultValue)
		return defaultValue
	}

	return value
}
```

**การทดสอบด้วย Postman**: ไม่มี (เป็นการตั้งค่าเท่านั้น)

**Unit Test**: สร้างไฟล์ `backend/test/sst-whisper/config_test.go`

**รันคำสั่ง Test**:
```bash
# Linux/macOS
cd backend/test/sst-whisper
go test -v -run TestWhisperConfigDefaults
go test -v -run TestWhisperConfigOverride
go test -v -run TestWhisperBinaryPathByOS
go test -v -run TestWhisperSupportedLanguages

# Windows
cd backend\test\sst-whisper
go test -v -run TestWhisperConfigDefaults
go test -v -run TestWhisperConfigOverride
go test -v -run TestWhisperBinaryPathByOS
go test -v -run TestWhisperSupportedLanguages

# รัน test ทั้งหมด
go test -v
```

```go
package whisper_test

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/yourusername/chatbot/config"
)

// TestWhisperConfigDefaults ทดสอบค่า default ของการตั้งค่า whisper
func TestWhisperConfigDefaults(t *testing.T) {
	// ล้างค่า env vars ของ whisper เพื่อทดสอบค่า default
	os.Unsetenv("WHISPER_BINARY_PATH_LINUX")
	os.Unsetenv("WHISPER_BINARY_PATH_WINDOWS")
	os.Unsetenv("WHISPER_BINARY_PATH_MACOS")
	os.Unsetenv("WHISPER_MODEL_PATH")
	os.Unsetenv("WHISPER_LANGUAGE")

	cfg := config.LoadConfig()

	// ตรวจสอบค่า default ของภาษา (ควรเป็น auto)
	if cfg.WhisperLanguage != "auto" {
		t.Errorf("คาดหวัง default language 'auto' แต่ได้ '%s'", cfg.WhisperLanguage)
	}

	if cfg.WhisperThreads != 4 {
		t.Errorf("คาดหวัง default threads 4 แต่ได้ %d", cfg.WhisperThreads)
	}

	if cfg.WhisperBeamSize != 5 {
		t.Errorf("คาดหวัง default beam size 5 แต่ได้ %d", cfg.WhisperBeamSize)
	}

	if cfg.WhisperTempDir != "./backend/whisper/temp" {
		t.Errorf("คาดหวัง default temp dir './backend/whisper/temp' แต่ได้ '%s'", cfg.WhisperTempDir)
	}

	t.Log("✓ ค่า default ทั้งหมดของ whisper config ถูกต้อง")
}

// TestWhisperBinaryPathByOS ทดสอบการเลือก binary path ตามระบบปฏิบัติการ
func TestWhisperBinaryPathByOS(t *testing.T) {
	cfg := config.LoadConfig()

	switch runtime.GOOS {
	case "windows":
		if !strings.Contains(cfg.WhisperBinaryPath, "windows") || !strings.HasSuffix(cfg.WhisperBinaryPath, ".exe") {
			t.Errorf("คาดหวัง Windows binary path แต่ได้ '%s'", cfg.WhisperBinaryPath)
		}
	case "darwin":
		if !strings.Contains(cfg.WhisperBinaryPath, "macos") {
			t.Errorf("คาดหวัง macOS binary path แต่ได้ '%s'", cfg.WhisperBinaryPath)
		}
	case "linux":
		if !strings.Contains(cfg.WhisperBinaryPath, "linux") {
			t.Errorf("คาดหวัง Linux binary path แต่ได้ '%s'", cfg.WhisperBinaryPath)
		}
	}

	t.Logf("✓ Binary path สำหรับ %s: %s", runtime.GOOS, cfg.WhisperBinaryPath)
}

// TestWhisperConfigOverride ทดสอบการ override ด้วย environment variables
func TestWhisperConfigOverride(t *testing.T) {
	// ตั้งค่าแบบกำหนดเอง
	os.Setenv("WHISPER_LANGUAGE", "en")
	os.Setenv("WHISPER_THREADS", "8")
	os.Setenv("WHISPER_BEAM_SIZE", "10")

	cfg := config.LoadConfig()

	if cfg.WhisperLanguage != "en" {
		t.Errorf("คาดหวัง language override 'en' แต่ได้ '%s'", cfg.WhisperLanguage)
	}

	if cfg.WhisperThreads != 8 {
		t.Errorf("คาดหวัง threads override 8 แต่ได้ %d", cfg.WhisperThreads)
	}

	if cfg.WhisperBeamSize != 10 {
		t.Errorf("คาดหวัง beam size override 10 แต่ได้ %d", cfg.WhisperBeamSize)
	}

	// รีเซ็ตค่ากลับ
	os.Unsetenv("WHISPER_LANGUAGE")
	os.Unsetenv("WHISPER_THREADS")
	os.Unsetenv("WHISPER_BEAM_SIZE")

	t.Log("✓ การ override ด้วย environment variables ทำงานถูกต้อง")
}

// TestWhisperSupportedLanguages ทดสอบการโหลดรายการภาษาที่รองรับ
func TestWhisperSupportedLanguages(t *testing.T) {
	cfg := config.LoadConfig()

	// ตรวจสอบว่ามีภาษาไทยและอังกฤษ
	foundTh := false
	foundEn := false
	foundAuto := false

	for _, lang := range cfg.WhisperSupportedLanguages {
		lang = strings.TrimSpace(lang)
		if lang == "th" {
			foundTh = true
		}
		if lang == "en" {
			foundEn = true
		}
		if lang == "auto" {
			foundAuto = true
		}
	}

	if !foundTh {
		t.Error("ไม่พบภาษาไทย (th) ในรายการภาษาที่รองรับ")
	}

	if !foundEn {
		t.Error("ไม่พบภาษาอังกฤษ (en) ในรายการภาษาที่รองรับ")
	}

	if !foundAuto {
		t.Error("ไม่พบการตรวจจับอัตโนมัติ (auto) ในรายการภาษาที่รองรับ")
	}

	t.Logf("✓ รองรับภาษา: %v", cfg.WhisperSupportedLanguages)
}
```

---

### งานที่ 3: สร้าง WhisperCppService Interface

**วัตถุประสงค์**: กำหนด interface สำหรับ service การแปลงคำพูดของ whisper.cpp

**ขั้นตอน**:

1. **สร้างไฟล์ `backend/services/transcription_service.go`** (ในโฟลเดอร์ services หลักของ backend):

```go
package services

import (
	"io"
)

// TranscriptionService กำหนด interface สำหรับ speech-to-text services
// ช่วยให้สามารถใช้งาน STT implementation หลายแบบ (whisper.cpp, Google STT, etc.)
type TranscriptionService interface {
	// Transcribe แปลงไฟล์ audio เป็นข้อความ
	// Parameters:
	//   - audioFile: io.Reader ที่มีข้อมูล audio
	//   - filename: ชื่อไฟล์เดิม (ใช้เพื่อกำหนดรูปแบบ)
	//   - language: รหัสภาษา (เช่น "th", "en", "auto")
	// Returns:
	//   - transcription: ข้อความที่แปลงได้
	//   - confidence: คะแนนความมั่นใจ (0.0 - 1.0)
	//   - error: ข้อผิดพลาดที่เกิดขึ้น
	Transcribe(audioFile io.Reader, filename string, language string) (transcription string, confidence float64, err error)

	// TranscribeWithTimestamps คืนค่าการแปลงพร้อม timestamps ระดับคำ
	// Parameters:
	//   - audioFile: io.Reader ที่มีข้อมูล audio
	//   - filename: ชื่อไฟล์เดิม
	//   - language: รหัสภาษา
	// Returns:
	//   - segments: array ของ segments การแปลงพร้อม timestamps
	//   - error: ข้อผิดพลาดที่เกิดขึ้น
	TranscribeWithTimestamps(audioFile io.Reader, filename string, language string) (segments []TranscriptionSegment, err error)

	// IsAvailable ตรวจสอบว่า service ตั้งค่าถูกต้องและพร้อมใช้งาน
	IsAvailable() bool

	// GetSupportedFormats คืนรายการรูปแบบ audio ที่รองรับ
	GetSupportedFormats() []string
}

// TranscriptionSegment แทน segment ของ audio ที่แปลงแล้วพร้อม timestamps
type TranscriptionSegment struct {
	StartTime float64 `json:"start_time"` // เวลาเริ่มต้นเป็นวินาที
	EndTime   float64 `json:"end_time"`   // เวลาสิ้นสุดเป็นวินาที
	Text      string  `json:"text"`       // ข้อความที่แปลงได้สำหรับ segment นี้
}

// TranscriptionResponse แทน API response สำหรับคำขอการแปลง
type TranscriptionResponse struct {
	Success      bool                   `json:"success"`
	Transcription string                 `json:"transcription"`
	Confidence   float64                `json:"confidence,omitempty"`
	Segments     []TranscriptionSegment `json:"segments,omitempty"`
	Language     string                 `json:"language"`
	Duration     float64                `json:"duration,omitempty"` // ระยะเวลา audio เป็นวินาที
	Error        string                 `json:"error,omitempty"`
}
```

**การทดสอบด้วย Postman**: ไม่มี (เป็นการกำหนด interface เท่านั้น)

---

### งานที่ 4: Implement WhisperCppService

**วัตถุประสงค์**: สร้าง service ที่รัน whisper.cpp binary และประมวลผล output

**ขั้นตอน**:

1. **สร้างไฟล์ `backend/services/whispercpp_service.go`** (ในโฟลเดอร์ services หลักของ backend): (ดูโค้ดตัวอย่างจากเอกสารภาษาอังกฤษ - มีประมาณ 250 บรรทัด พร้อม comments ภาษาไทย)

---

### งานที่ 5: สร้าง WhisperCppController (Endpoint เฉพาะ)

**วัตถุประสงค์**: สร้าง API endpoint เฉพาะสำหรับ whisper.cpp STT ที่ `POST /api/stt/whispercpp`

**ขั้นตอน**:

1. **สร้างไฟล์ `backend/controllers/whispercpp_controller.go`** (ในโฟลเดอร์ controllers หลักของ backend)

**คำอธิบาย API**:

**Endpoint**: `POST /api/stt/whispercpp`

**Request Body** (multipart/form-data):
- `audio`: ไฟล์ audio (จำเป็น) - WAV, MP3, M4A, OGG, FLAC
- `language`: รหัสภาษา (ตัวเลือก) - "th", "en", "auto" (ค่าเริ่มต้น: "th")
- `timestamps`: boolean (ตัวเลือก) - คืนค่า segments พร้อม timestamps (ค่าเริ่มต้น: false)

**Response สำหรับการแปลงแบบธรรมดา** (200):
```json
{
  "success": true,
  "transcription": "สวัสดีครับ ยินดีต้อนรับสู่ระบบแชทบอท",
  "confidence": 0.92,
  "language": "th"
}
```

**Response สำหรับการแปลงพร้อม Timestamps** (200):
```json
{
  "success": true,
  "transcription": "สวัสดีครับ ยินดีต้อนรับสู่ระบบแชทบอท",
  "segments": [
    {
      "start_time": 0.0,
      "end_time": 2.5,
      "text": "สวัสดีครับ"
    },
    {
      "start_time": 2.5,
      "end_time": 5.8,
      "text": "ยินดีต้อนรับสู่ระบบแชทบอท"
    }
  ],
  "language": "th",
  "duration": 5.8
}
```

**การทดสอบด้วย Postman**:

**Test 1: การแปลงภาษาไทยแบบง่าย**
```
POST http://localhost:3001/api/stt/whispercpp
Content-Type: multipart/form-data

Body:
- audio: [อัปโหลดไฟล์เสียงภาษาไทย - thai_sample.wav]
- language: th

ผลลัพธ์ที่คาดหวัง (200):
{
  "success": true,
  "transcription": "สวัสดีครับ ยินดีต้อนรับสู่ระบบแชทบอท",
  "confidence": 0.92,
  "language": "th"
}
```

**Test 2: การแปลงภาษาอังกฤษแบบง่าย**
```
POST http://localhost:3001/api/stt/whispercpp
Content-Type: multipart/form-data

Body:
- audio: [อัปโหลดไฟล์เสียงภาษาอังกฤษ - english_sample.wav]
- language: en

ผลลัพธ์ที่คาดหวัง (200):
{
  "success": true,
  "transcription": "Hello, welcome to our chatbot system",
  "confidence": 0.95,
  "language": "en"
}
```

**Test 3: การตรวจจับภาษาอัตโนมัติ**
```
POST http://localhost:3001/api/stt/whispercpp
Content-Type: multipart/form-data

Body:
- audio: [อัปโหลดไฟล์เสียงภาษาใดก็ได้]
- language: auto

ผลลัพธ์ที่คาดหวัง (200):
{
  "success": true,
  "transcription": "...",
  "confidence": 0.90,
  "language": "th"  // หรือ "en" ขึ้นอยู่กับภาษาที่ตรวจจับได้
}
```

**Test 4: การแปลงพร้อม Timestamps (ภาษาไทย)**
```
POST http://localhost:3001/api/stt/whispercpp
Content-Type: multipart/form-data

Body:
- audio: [อัปโหลดไฟล์เสียงภาษาไทย]
- language: th
- timestamps: true

ผลลัพธ์ที่คาดหวัง (200):
{
  "success": true,
  "transcription": "สวัสดีครับ ยินดีต้อนรับสู่ระบบแชทบอท",
  "segments": [
    {
      "start_time": 0.0,
      "end_time": 2.5,
      "text": "สวัสดีครับ"
    },
    {
      "start_time": 2.5,
      "end_time": 5.8,
      "text": "ยินดีต้อนรับสู่ระบบแชทบอท"
    }
  ],
  "language": "th",
  "duration": 5.8
}
```

**Test 5: ตรวจสอบสถานะ Service**
```
GET http://localhost:3001/api/stt/whispercpp/status

ผลลัพธ์ที่คาดหวัง (200):
{
  "service": "whisper.cpp",
  "available": true,
  "supported_formats": ["wav", "mp3", "m4a", "ogg", "flac"],
  "supported_languages": ["th", "en", "auto"],
  "default_language": "auto",
  "current_os": "windows",
  "model": "small"
}
```

---

### งานที่ 6: ลงทะเบียน Routes

**วัตถุประสงค์**: ลงทะเบียน endpoints ของ whisper.cpp ในการตั้งค่า routing

**ขั้นตอน**:

1. **อัปเดต `backend/routes/routes.go`**:

```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/chatbot/config"
	"github.com/yourusername/chatbot/controllers"  // Controllers อยู่ในโฟลเดอร์หลัก
	"github.com/yourusername/chatbot/middleware"
	"github.com/yourusername/chatbot/services"     // Services อยู่ในโฟลเดอร์หลัก
)

// SetupRoutes ตั้งค่า routes ทั้งหมดของแอปพลิเคชัน
func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// ... โค้ดเดิม ...

	// เริ่มต้น Whisper.cpp service และ controller
	// Services และ Controllers อยู่ในโฟลเดอร์หลักของ backend แล้ว
	whisperService := services.NewWhisperCppService(cfg)
	whisperController := controllers.NewWhisperCppController(whisperService)

	// ========================================
	// Speech-to-Text (STT) Routes
	// ========================================
	sttGroup := api.Group("/stt")

	// Endpoints เฉพาะของ Whisper.cpp
	sttGroup.Get("/whispercpp/status", whisperController.GetStatus)
	sttGroup.Post("/whispercpp",
		middleware.RateLimiter(cfg),
		whisperController.TranscribeAudio,
	)

	// ... routes อื่นๆ ที่เหลือ ...
}
```

**หมายเหตุ**:
- ไฟล์ `transcription_service.go` และ `whispercpp_service.go` อยู่ใน `backend/services/`
- ไฟล์ `whispercpp_controller.go` อยู่ใน `backend/controllers/`
- ทั้งสองโฟลเดอร์นี้เป็นโฟลเดอร์หลักของ backend ที่รวม services และ controllers ทั้งหมดของโปรเจ็ค

---

### งานที่ 7: Integration Testing

**วัตถุประสงค์**: ทดสอบแบบ end-to-end ด้วยไฟล์ audio จริง

**ขั้นตอน**:

1. **เตรียมไฟล์ audio สำหรับทดสอบ** ใน `backend/test/sst-whisper/testdata/audio/`:
   - `thai_short.wav` - ประโยคภาษาไทยสั้นๆ (3-5 วินาที)
   - `thai_long.wav` - คำพูดภาษาไทยยาวขึ้น (10-30 วินาที)
   - `english_short.wav` - คำพูดภาษาอังกฤษสั้น
   - `english_long.wav` - คำพูดภาษาอังกฤษยาว
   - `mixed_language.wav` - คำพูดผสมภาษาไทย-อังกฤษ

2. **สร้างไฟล์ `backend/test/sst-whisper/integration_test.go`**: (ดูโค้ดตัวอย่างจากเอกสารภาษาอังกฤษ)

3. **รัน integration tests**:
```bash
cd backend/test/sst-whisper
go test -v -run Integration
```

---

## เคล็ดลับการปรับแต่งสำหรับภาษาไทย

### การเลือกโมเดล

- **โมเดล Small** (ggml-small.bin): เหมาะที่สุดสำหรับการใช้งานภาษาไทยทั่วไป
  - ประมวลผลเร็ว
  - ความแม่นยำดีสำหรับคำพูดภาษาไทยที่ชัดเจน
  - ใช้พื้นที่ 461 MB

- **โมเดล Medium** (ggml-medium.bin): เพื่อความแม่นยำสูงขึ้น
  - ดีกว่าสำหรับภาษาไทยที่มีสำเนียง
  - จัดการเสียงรบกวนพื้นหลังได้ดีกว่า
  - ใช้พื้นที่ 1.5 GB

### การตั้งค่าสำหรับภาษาไทย

```env
WHISPER_LANGUAGE=th           # บังคับตรวจจับภาษาไทย
WHISPER_BEAM_SIZE=5           # เพิ่มความแม่นยำ
WHISPER_BEST_OF=5             # ผลลัพธ์ที่ดีขึ้นสำหรับเสียงวรรณยุกต์
WHISPER_THREADS=4             # ปรับตาม CPU
```

### คำแนะนำเกี่ยวกับคุณภาพ Audio

- **Sample rate**: ขั้นต่ำ 16 kHz (whisper.cpp จะทำ resample)
- **Format**: WAV (ไม่บีบอัด) เพื่อผลลัพธ์ที่ดีที่สุด
- **Channels**: Mono แนะนำ (stereo จะถูกแปลง)
- **เสียงรบกวนพื้นหลัง**: ลดให้น้อยที่สุดเพื่อความแม่นยำที่ดีขึ้น

### ปัญหาทั่วไปในการแปลงภาษาไทย

1. **ความสับสนเรื่องวรรณยุกต์**: ภาษาไทยมี 5 เสียง ต้องออกเสียงให้ชัดเจน
2. **คำยืม**: คำภาษาอังกฤษในคำพูดภาษาไทยอาจถูกแปลงเป็นภาษาอังกฤษ
3. **ภาษาทางการ vs ภาษาพูด**: โมเดลจัดการได้ทั้งสองแบบ แต่ภาษาทางการแปลงได้ดีกว่า
4. **ตัวเลข**: ตัวเลขภาษาไทยจะถูกแปลงเป็นคำภาษาไทย

---

## Checklist การ Deploy

### การติดตั้งและตั้งค่า
- [ ] สร้างโครงสร้างโฟลเดอร์ `backend/whisper/` แล้ว
- [ ] คอมไพล์ whisper.cpp สำหรับ OS ที่ใช้งาน (Linux/Windows/macOS)
- [ ] ดาวน์โหลดโมเดล GGML แล้ว (small หรือ medium)
- [ ] คัดลอก binary ไปยังโฟลเดอร์ที่ถูกต้อง
- [ ] ตั้งค่า environment variables ตาม OS
- [ ] สร้างไดเรกทอรี temp แล้ว (./backend/whisper/temp)
- [ ] สร้างโฟลเดอร์ทดสอบ (./backend/test/sst-whisper/)

### การตั้งค่าระบบ
- [ ] อัพเดต backend/config/config.go ให้รองรับหลาย OS
- [ ] ตั้งค่า binary path สำหรับ Linux, Windows, macOS
- [ ] ตั้งค่าการรองรับภาษา (th, en, auto)
- [ ] ตั้งค่าขีดจำกัดขนาดไฟล์อัปโหลดแล้ว (Fiber)
- [ ] ตั้งค่า rate limiting สำหรับ STT endpoint แล้ว

### การทดสอบ
- [ ] Unit tests ผ่านแล้ว (setup_test.go)
- [ ] Config tests ผ่านแล้ว (config_test.go)
- [ ] ทดสอบ binary บน OS ที่ใช้งานแล้ว
- [ ] ทดสอบการแปลงเสียงภาษาไทยแล้ว
- [ ] ทดสอบการแปลงเสียงภาษาอังกฤษแล้ว
- [ ] ทดสอบการตรวจจับภาษาอัตโนมัติแล้ว
- [ ] Integration tests ผ่านแล้ว
- [ ] ทดสอบ Postman สำเร็จแล้ว (ทั้ง th และ en)

### การตรวจสอบคุณภาพ
- [ ] ตรวจสอบ error handling แล้ว
- [ ] ตั้งค่า logging แล้ว
- [ ] ทดสอบการทำงานบนทั้ง 3 OS (ถ้าเป็นไปได้)
- [ ] อัปเดตเอกสารแล้ว
- [ ] ตรวจสอบความปลอดภัยของ API แล้ว

---

## สรุป API

### Endpoint: POST /api/stt/whispercpp

**วัตถุประสงค์**: Endpoint เฉพาะสำหรับการแปลง speech-to-text ด้วย whisper.cpp พร้อมรองรับ**ภาษาไทย และ ภาษาอังกฤษ**

**Request**:
```
POST /api/stt/whispercpp
Content-Type: multipart/form-data

ฟิลด์ฟอร์ม:
- audio: [File] ไฟล์ audio (WAV, MP3, M4A, OGG, FLAC)
- language: [String] รหัสภาษา - ค่าเริ่มต้น: "auto"
  * "th" = ภาษาไทย
  * "en" = ภาษาอังกฤษ
  * "auto" = ตรวจจับภาษาอัตโนมัติ
- timestamps: [Boolean] รวม timestamped segments - ค่าเริ่มต้น: false
```

**Response (แบบธรรมดา)**:
```json
{
  "success": true,
  "transcription": "สวัสดีครับ ยินดีต้อนรับสู่ระบบแชทบอท",
  "confidence": 0.92,
  "language": "th"
}
```

**Response (พร้อม Timestamps)**:
```json
{
  "success": true,
  "transcription": "สวัสดีครับ ยินดีต้อนรับสู่ระบบแชทบอท",
  "segments": [
    {
      "start_time": 0.0,
      "end_time": 2.5,
      "text": "สวัสดีครับ"
    },
    {
      "start_time": 2.5,
      "end_time": 5.8,
      "text": "ยินดีต้อนรับสู่ระบบแชทบอท"
    }
  ],
  "language": "th",
  "duration": 5.8
}
```

### Endpoint: GET /api/stt/whispercpp/status

**วัตถุประสงค์**: ตรวจสอบสถานะความพร้อมใช้งานของ whisper.cpp service

**Response**:
```json
{
  "service": "whisper.cpp",
  "available": true,
  "supported_formats": ["wav", "mp3", "m4a", "ogg", "flac"],
  "supported_languages": ["th", "en", "auto"],
  "default_language": "auto",
  "current_os": "windows",
  "model": "small"
}
```

---

## แหล่งข้อมูลอ้างอิง

- [whisper.cpp GitHub](https://github.com/ggerganov/whisper.cpp)
- [OpenAI Whisper Paper](https://arxiv.org/abs/2212.04356)
- [โมเดล GGML บน Hugging Face](https://huggingface.co/ggerganov/whisper.cpp)
- [การรองรับภาษาไทยใน Whisper](https://github.com/openai/whisper#available-models-and-languages)

---

## สรุปคำสั่งทดสอบสำหรับแต่ละระบบปฏิบัติการ

### Linux

```bash
# 1. การติดตั้งและคอมไพล์
cd backend/whisper
git clone https://github.com/ggerganov/whisper.cpp.git whisper-source
cd whisper-source
sudo apt-get update && sudo apt-get install build-essential
make clean && make
cp main ../binary/linux/main
chmod +x ../binary/linux/main

# 2. ดาวน์โหลดโมเดล
cd ../models
curl -L -o ggml-small.bin https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin

# 3. ทดสอบ Binary
../binary/linux/main -m ggml-small.bin -f ../../test/sst-whisper/testdata/audio/thai_short.wav -l th
../binary/linux/main -m ggml-small.bin -f ../../test/sst-whisper/testdata/audio/english_short.wav -l en

# 4. รัน Unit Tests
cd ../../test/sst-whisper
go test -v
go test -v -run TestWhisperBinaryExists
go test -v -run TestWhisperTranscribeThaiAudio
go test -v -run TestWhisperTranscribeEnglishAudio
```

### Windows

```powershell
# 1. การติดตั้งและคอมไพล์ (ใช้ MSYS2)
# ติดตั้ง MSYS2 จาก https://www.msys2.org/ ก่อน
# จาก MSYS2 terminal:
pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-cmake make

cd backend/whisper
git clone https://github.com/ggerganov/whisper.cpp.git whisper-source
cd whisper-source
make clean
make
cp main.exe ..\binary\windows\main.exe

# 2. ดาวน์โหลดโมเดล (PowerShell)
cd ..\models
Invoke-WebRequest -Uri "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin" -OutFile "ggml-small.bin"

# 3. ทดสอบ Binary
..\binary\windows\main.exe -m ggml-small.bin -f ..\..\test\sst-whisper\testdata\audio\thai_short.wav -l th
..\binary\windows\main.exe -m ggml-small.bin -f ..\..\test\sst-whisper\testdata\audio\english_short.wav -l en

# 4. รัน Unit Tests
cd ..\..\test\sst-whisper
go test -v
go test -v -run TestWhisperBinaryExists
go test -v -run TestWhisperTranscribeThaiAudio
go test -v -run TestWhisperTranscribeEnglishAudio
```

### macOS

```bash
# 1. การติดตั้งและคอมไพล์
xcode-select --install  # ติดตั้ง Xcode Command Line Tools

cd backend/whisper
git clone https://github.com/ggerganov/whisper.cpp.git whisper-source
cd whisper-source
make clean && make
cp main ../binary/macos/main
chmod +x ../binary/macos/main

# สำหรับ Apple Silicon (M1/M2) - ใช้ Metal acceleration
# WHISPER_METAL=1 make

# 2. ดาวน์โหลดโมเดล
cd ../models
curl -L -o ggml-small.bin https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin

# 3. ทดสอบ Binary
../binary/macos/main -m ggml-small.bin -f ../../test/sst-whisper/testdata/audio/thai_short.wav -l th
../binary/macos/main -m ggml-small.bin -f ../../test/sst-whisper/testdata/audio/english_short.wav -l en

# 4. รัน Unit Tests
cd ../../test/sst-whisper
go test -v
go test -v -run TestWhisperBinaryExists
go test -v -run TestWhisperTranscribeThaiAudio
go test -v -run TestWhisperTranscribeEnglishAudio
```

### คำสั่งทดสอบทั้งหมด (Cross-platform)

```bash
# รัน test แบบ specific
go test -v -run TestWhisperBinaryExists
go test -v -run TestWhisperModelExists
go test -v -run TestWhisperVersion
go test -v -run TestWhisperConfigDefaults
go test -v -run TestWhisperBinaryPathByOS
go test -v -run TestWhisperSupportedLanguages
go test -v -run TestWhisperTranscribeThaiAudio
go test -v -run TestWhisperTranscribeEnglishAudio

# รัน test ทั้งหมด
cd backend/test/sst-whisper
go test -v

# รัน test แบบ verbose พร้อมแสดงเวลา
go test -v -timeout 30m

# รัน specific test file
go test -v setup_test.go
go test -v config_test.go
```

---

## ขั้นตอนถัดไป

หลังจากเสร็จสิ้นการ implement นี้:

1. **การผสานระบบ Frontend**: เพิ่มการบันทึกและอัปโหลด audio ใน Vue.js frontend
2. **STT แบบ Real-time**: ใช้งานการแปลง audio แบบ streaming
3. **การสลับโมเดล**: ให้ผู้ใช้เลือกขนาดโมเดลได้ (small/medium/large)
4. **การประมวลผลแบบ Batch**: รองรับการอัปโหลดหลายไฟล์
5. **คำสั่งเสียง**: ผสานเข้ากับ chatbot สำหรับการโต้ตอบด้วยเสียง
6. **การติดตามประสิทธิภาพ**: ติดตามความหน่วงและความแม่นยำของการแปลง
7. **การรองรับหลายภาษา**: ขยายไปยังภาษาอื่นนอกเหนือจากไทย/อังกฤษ

---

## สรุปคุณสมบัติหลัก

✅ **รองรับหลายภาษา**: ภาษาไทย, ภาษาอังกฤษ และการตรวจจับอัตโนมัติ
✅ **รองรับหลาย OS**: Linux, Windows, และ macOS
✅ **โครงสร้างโฟลเดอร์ที่เป็นระเบียบ**:
   - Binaries และ Models: `backend/whisper/`
   - Services: `backend/services/` (รวมกับ services อื่นๆ)
   - Controllers: `backend/controllers/` (รวมกับ controllers อื่นๆ)
   - Tests: `backend/test/sst-whisper/`
✅ **Unit Tests**: ครอบคลุมทุกระบบปฏิบัติการ
✅ **คำสั่งทดสอบ**: มีคำสั่งชัดเจนสำหรับแต่ละ test
✅ **Integration ที่ดี**: Services และ Controllers รวมกับโครงสร้างหลักของ backend

---

**เวอร์ชันเอกสาร**: 2.1
**อัปเดตล่าสุด**: 2025-11-10
**สถานะ**: พร้อมสำหรับการ Implementation (รองรับ Multi-language & Multi-platform)
**การเปลี่ยนแปลง v2.1**:
- ย้าย services ไปยัง `backend/services/` (รวมกับ services อื่นๆ)
- ย้าย controllers ไปยัง `backend/controllers/` (รวมกับ controllers อื่นๆ)
- โครงสร้างโฟลเดอร์สอดคล้องกับ architecture ของโปรเจ็คมากขึ้น
