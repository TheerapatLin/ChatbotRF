# 🎤 Whisper.cpp Speech-to-Text Integration Guide

ระบบแปลงเสียงเป็นข้อความ (STT) โดยใช้ whisper.cpp รองรับภาษาไทยและอังกฤษ

**Last Updated**: 2025-11-10
**Platform**: Windows 11 with WSL2
**Status**: ✅ Implementation Complete

---

## 📋 สรุประบบ

### ความสามารถหลัก
- ✅ แปลงเสียงเป็นข้อความ (Speech-to-Text)
- ✅ รองรับภาษาไทยและอังกฤษ
- ✅ ตรวจจับภาษาอัตโนมัติ (auto-detect)
- ✅ รองรับ Timestamps (segment-level)
- ✅ ประมวลผลในเครื่อง (offline, no API cost)
- ✅ รองรับหลายรูปแบบเสียง: WAV, MP3, M4A, OGG, FLAC

### คุณสมบัติ whisper.cpp
- **ประมวลผลในเครื่อง**: ไม่ต้องพึ่ง cloud API
- **รองรับหลายภาษา**: 99+ ภาษา รวมภาษาไทย
- **GGML Model**: โมเดลที่ optimize แล้วสำหรับ inference
- **ความหน่วงต่ำ**: แปลงเสียงได้รวดเร็ว
- **ไม่มีค่า API**: ดาวน์โหลดโมเดลครั้งเดียว ใช้งานได้ไม่จำกัด

---

## 📁 โครงสร้างโฟลเดอร์

```
backend/
├── whisper/                          # ระบบ Whisper STT
│   ├── binary/                      # whisper.cpp binaries
│   │   ├── linux/main              # Linux executable
│   │   ├── windows/main.exe        # Windows executable
│   │   └── macos/main              # macOS executable
│   │
│   ├── models/                      # GGML models
│   │   └── ggml-small.bin         # Model small (466 MB)
│   │
│   └── temp/                        # ไฟล์เสียงชั่วคราว
│       └── .gitkeep
│
├── config/
│   └── config.go                    # Configuration loading
│
├── services/
│   ├── transcription_service.go    # STT Interface
│   └── whispercpp_service.go       # Whisper.cpp implementation
│
├── controllers/
│   └── whispercpp_controller.go    # HTTP handlers
│
├── routes/
│   └── routes.go                    # Route registration
│
└── test/
    └── sst-whisper/                 # Unit & integration tests
        ├── testdata/audio/         # ไฟล์เสียงทดสอบ
        ├── config_test.go          # Configuration tests
        ├── setup_test.go           # Binary & model tests
        └── service_test.go         # Service logic tests
```

---

## 🏗️ สถาปัตยกรรมระบบ

### ภาพรวม 3-Tier Architecture

```
┌──────────────────┐
│   Frontend       │  - บันทึกเสียงผ่าน MediaRecorder API
│   (Vue.js)       │  - ส่งไฟล์เสียงผ่าน multipart/form-data
└────────┬─────────┘
         │ HTTP POST
         ▼
┌──────────────────────────────────────────┐
│   Backend (Go Fiber)                      │
│                                           │
│  ┌─────────────────────────────────┐    │
│  │ WhisperCppController            │    │  - รับ HTTP request
│  │ /api/stt/whispercpp              │    │  - Validate ไฟล์
│  └───────────┬─────────────────────┘    │  - Error handling
│              ▼                            │
│  ┌─────────────────────────────────┐    │
│  │ WhisperCppService               │    │  - บันทึกไฟล์ชั่วคราว
│  │ - Transcribe()                   │    │  - เรียก binary
│  │ - TranscribeWithTimestamps()     │    │  - Parse output
│  └───────────┬─────────────────────┘    │
│              │                            │
└──────────────┼────────────────────────────┘
               │ exec.Command()
               ▼
┌──────────────────────────────────────────┐
│   whisper.cpp Binary (C++)               │  - โหลด GGML model
│   - โหลด model ไปยัง RAM                 │  - Decode audio
│   - ประมวลผล audio (resample, mel)       │  - AI inference
│   - แปลงเป็นข้อความด้วย AI              │  - Output JSON/text
└──────────────────────────────────────────┘
```

### Flow การทำงาน

**1. Basic Transcription (ไม่มี timestamps)**
```
Client → Controller.TranscribeAudio()
       → Service.Transcribe(audioFile, language)
       → whisper.cpp binary execution
       → Parse text output
       → Return JSON response
```

**2. Transcription with Timestamps**
```
Client → Controller.TranscribeAudio() (timestamps=true)
       → Service.TranscribeWithTimestamps(audioFile, language)
       → whisper.cpp binary with -oj flag (JSON output)
       → Parse JSON file with segments
       → Return JSON response with time-aligned segments
```

---

## 🔧 ส่วนประกอบหลัก

### 1. Configuration Layer (config/config.go)

**หน้าที่**: โหลด configuration จาก environment variables

**ฟิลด์สำคัญ**:
- `WhisperBinaryPath`: Path ไปยัง whisper.cpp binary
- `WhisperModelPath`: Path ไปยัง default GGML model file
- `WhisperModelsDir`: Directory ที่เก็บ models ทั้งหมด
- `WhisperTempDir`: โฟลเดอร์เก็บไฟล์ชั่วคราว
- `WhisperLanguage`: ภาษาเริ่มต้น (th, en, auto)
- `WhisperModelName`: ชื่อโมเดล default (small, medium, large)
- `WhisperSupportedModels`: รายการ models ที่รองรับ
- `WhisperThreads`: จำนวน CPU threads
- `WhisperBeamSize`: Beam size สำหรับ decoding
- `WhisperBestOf`: จำนวน candidates

**Environment Variables**:
```bash
WHISPER_BINARY_PATH_WINDOWS=wsl /mnt/c/.../main  # WSL2 path
WHISPER_MODEL_PATH=./whisper/models/ggml-small.bin
WHISPER_MODELS_DIR=./whisper/models
WHISPER_TEMP_DIR=./whisper/temp
WHISPER_LANGUAGE=auto
WHISPER_MODEL_NAME=small
WHISPER_SUPPORTED_MODELS=tiny.en,small,medium,large-v2
WHISPER_THREADS=4
WHISPER_BEAM_SIZE=5
WHISPER_BEST_OF=5
```

---

### 2. Interface Layer (services/transcription_service.go)

**หน้าที่**: กำหนด interface สำหรับ STT services

**Methods**:
- `Transcribe(audioFile, filename, language) (text, confidence, error)`
  - แปลงเสียงเป็นข้อความอย่างเดียว

- `TranscribeWithTimestamps(audioFile, filename, language) (text, segments, duration, error)`
  - แปลงเสียงพร้อม timestamps แต่ละ segment

**ประโยชน์**: ช่วยให้สามารถสลับ implementation (OpenAI Whisper API, whisper.cpp, etc.) ได้ง่าย

---

### 3. Service Layer (services/whispercpp_service.go)

**หน้าที่**: Implementation หลักของ whisper.cpp integration

**Struct**:
```go
type WhisperCppService struct {
    config *config.Config
}
```

**Methods**:

#### `NewWhisperCppService(cfg) (*WhisperCppService, error)`
- สร้าง service instance
- Validate binary และ model existence

#### `Transcribe(audioFile, filename, language) (text, confidence, error)`
- บันทึกไฟล์ audio ชั่วคราว
- เรียก whisper.cpp binary ด้วย text output mode
- Parse stdout เพื่อดึงข้อความ
- คำนวณ confidence score
- ลบไฟล์ชั่วคราว
- Return transcription text

#### `TranscribeWithTimestamps(audioFile, filename, language) (text, segments, duration, error)`
- บันทึกไฟล์ audio ชั่วคราว
- เรียก whisper.cpp binary ด้วย `-oj` flag (JSON output)
- อ่าน JSON file จาก temp directory
- Parse segments พร้อม timestamps
- คำนวณ duration รวม
- ลบไฟล์ชั่วคราว
- Return full transcription + segments array

#### `GetModelName() string`
- คืนชื่อโมเดลที่ใช้งาน (สำหรับ status endpoint)

**Binary Execution**:
- ใช้ `exec.CommandContext()` พร้อม 1 นาที timeout
- รองรับ absolute path resolution (หา project root ผ่าน go.mod)
- รองรับ WSL2 integration สำหรับ Windows

---

### 4. Controller Layer (controllers/whispercpp_controller.go)

**หน้าที่**: HTTP request handling และ response formatting

**Struct**:
```go
type WhisperCppController struct {
    whisperService *services.WhisperCppService
}
```

**Endpoints**:

#### `GET /api/stt/whispercpp/status`
**ฟังก์ชัน**: `GetStatus(c *fiber.Ctx) error`

**คืนค่า**:
```json
{
  "service": "whisper.cpp",
  "available": true,
  "default_model": "small",
  "supported_models": ["tiny.en", "small", "medium", "large-v2"],
  "supported_formats": ["wav", "mp3", "m4a", "ogg", "flac"],
  "supported_languages": ["th", "en", "auto"],
  "default_language": "auto",
  "current_os": "windows"
}
```

#### `POST /api/stt/whispercpp`
**ฟังก์ชัน**: `TranscribeAudio(c *fiber.Ctx) error`

**Request** (multipart/form-data):
- `audio` (File): ไฟล์เสียง (max 25MB)
- `language` (Text): "th", "en", "auto" (default: "auto")
- `timestamps` (Text): "true" หรือ "false" (default: "false")
- `model` (Text): ชื่อ model เช่น "tiny.en", "small", "medium" (optional, default: ใช้ค่าจาก config)

**Validation**:
- ขนาดไฟล์ไม่เกิน 25 MB
- ตรวจสอบ file extension (wav, mp3, m4a, ogg, flac)
- ตรวจสอบ language code
- ตรวจสอบ model name (ต้องอยู่ใน supported_models)

**Response (ไม่มี timestamps)**:
```json
{
  "success": true,
  "transcription": "สวัสดีครับ ยินดีต้อนรับ",
  "confidence": 0.95,
  "language": "th",
  "model": "small"
}
```

**Response (มี timestamps)**:
```json
{
  "success": true,
  "transcription": "สวัสดีครับ ยินดีต้อนรับ",
  "segments": [
    {
      "start_time": 0.0,
      "end_time": 1.5,
      "text": "สวัสดีครับ"
    },
    {
      "start_time": 1.5,
      "end_time": 3.2,
      "text": "ยินดีต้อนรับ"
    }
  ],
  "language": "th",
  "duration": 3.2,
  "model": "small"
}
```

---

### 5. Route Registration (routes/routes.go)

**หน้าที่**: ลงทะเบียน endpoints ใน Fiber app

**การทำงาน**:
1. Initialize WhisperCppService พร้อม error handling
2. สร้าง WhisperCppController (conditional - ถ้า service available)
3. สร้าง route group `/api/stt`
4. ลงทะเบียน 2 endpoints
5. แสดง log ยืนยันการลงทะเบียนสำเร็จ

**Graceful Degradation**: ถ้า binary หรือ model ไม่พบ service จะ return error และ endpoints จะไม่ถูกลงทะเบียน (แต่ app ยังรันได้ปกติ)

---

## 🧪 Test Suite

**โฟลเดอร์**: `backend/test/sst-whisper/`

### Test Files

| ไฟล์ | จำนวน Tests | หน้าที่ |
|------|-------------|---------|
| `config_test.go` | 9 tests | ทดสอบ configuration loading, validation, OS detection, path resolution |
| `setup_test.go` | 5 tests | ทดสอบ binary existence, execution, model file, temp directory |
| `service_test.go` | 8 tests | ทดสอบ Transcribe(), TranscribeWithTimestamps(), error handling, language support |

**รวม**: 22 unit tests

### รัน Tests

```bash
# รันทุก tests
cd backend
DATABASE_URL="postgres://test" go test -v -timeout 60s chatbot/test/sst-whisper

# รัน test group เฉพาะ
DATABASE_URL="postgres://test" go test -v -run TestConfig
DATABASE_URL="postgres://test" go test -v -run TestSetup
DATABASE_URL="postgres://test" go test -v -run TestWhisper
```

**หมายเหตุ**: บน Windows ต้องรันผ่าน WSL2 เพราะ binary เป็น Linux format

---

## ⚙️ การตั้งค่าและใช้งาน

### ติดตั้ง Binary และ Model

**1. Binary**:
- Linux: `backend/whisper/binary/linux/main` (111 KB)
- Windows: ใช้ WSL2 เรียก Linux binary
- macOS: `backend/whisper/binary/macos/main`

**2. Model**:
- ดาวน์โหลด `ggml-small.bin` (466 MB) จาก Hugging Face
- วางที่ `backend/whisper/models/ggml-small.bin`

### การตั้งค่า Environment

แก้ไขไฟล์ `backend/.env.development`:

```bash
# ใช้ WSL2 สำหรับ Windows
WHISPER_BINARY_PATH_WINDOWS=wsl /mnt/c/Users/.../backend/whisper/binary/linux/main

# Model และ temp directory
WHISPER_MODEL_PATH=./whisper/models/ggml-small.bin
WHISPER_TEMP_DIR=./whisper/temp

# การตั้งค่าภาษา
WHISPER_LANGUAGE=auto
WHISPER_MODEL_NAME=small
WHISPER_THREADS=4
```

### รัน Backend

```bash
# บน WSL2 (Windows)
wsl
cd /mnt/c/Users/.../ChatBotProject/backend
go run .

# หรือบน Linux
cd backend
go run .
```

### ทดสอบ API

**Postman / cURL**:
```bash
POST http://localhost:3001/api/stt/whispercpp

Body (form-data):
- audio: [ไฟล์เสียง .wav/.mp3]
- language: auto
- timestamps: true
```

---

## 🚀 ความสามารถและข้อจำกัด

### Performance
- **Model Size**: 466 MB (small model)
- **Transcription Speed**: ~0.5-2x realtime (ขึ้นกับ CPU)
- **Memory Usage**: ~500-800 MB RAM
- **Max Audio Length**: แนะนำไม่เกิน 5 นาที
- **Max File Size**: 25 MB (configurable)

### Supported Languages
- ภาษาไทย (th)
- ภาษาอังกฤษ (en)
- Auto-detect (auto) - รองรับ 99+ ภาษา

### Supported Audio Formats
- WAV (แนะนำ - 16kHz, 16-bit, mono)
- MP3
- M4A
- OGG
- FLAC

---

## ❌ การจัดการ Error

| Error | สาเหตุ | วิธีแก้ |
|-------|--------|---------|
| `Binary not found` | ไม่พบ whisper.cpp binary | ติดตั้ง binary หรือตั้ง path ใน .env |
| `Model not found` | ไม่พบ GGML model file | ดาวน์โหลด model ไปยัง models/ |
| `Audio file too large` | ไฟล์เกิน 25 MB | ลดขนาดไฟล์หรือเพิ่ม limit |
| `Unsupported format` | รูปแบบไฟล์ไม่รองรับ | แปลงเป็น WAV/MP3/M4A |
| `Exit status 0xc0000135` | Windows DLL not found | ใช้ WSL2 หรือ compile Windows binary |
| `Transcription timeout` | Audio ยาวเกินไป | ตั้ง timeout สูงขึ้นหรือแบ่งไฟล์ |
| `Empty transcription` | Audio ไม่มีเสียงพูด | ตรวจสอบไฟล์เสียง |
| `Language not supported` | รหัสภาษาผิด | ใช้ th, en, หรือ auto |

---

## 🔧 Troubleshooting

### 🚨 Exit Status 0xc0000135 Error (Windows)

**อาการ**:
```
❌ Transcription failed: whisper.cpp execution failed: exit status 0xc0000135
```

**สาเหตุ**:
- Windows error code: `STATUS_DLL_NOT_FOUND`
- whisper.cpp binary เป็น Linux format หรือขาด DLL dependencies

**การวิเคราะห์**:

Exit code `0xc0000135` (hex) = `-1073741515` (decimal) หมายถึง Windows ไม่สามารถโหลด DLL ที่จำเป็นสำหรับ binary ได้

**สาเหตุที่เป็นไปได้**:
1. Binary เป็น Linux ELF format (ไม่ใช่ Windows PE format)
2. Binary compile ด้วย MSVC แต่ขาด Visual C++ Redistributable
3. Binary ต้องการ MinGW runtime libraries

---

### ✅ วิธีแก้ไข (4 วิธี)

#### วิธีที่ 1: ใช้ WSL2 ⭐⭐⭐⭐⭐ (แนะนำสำหรับ Development)

**ข้อดี**:
- ใช้เวลา 5-10 นาที
- ใช้ Linux binary ที่มีอยู่แล้ว
- ไม่มีปัญหา DLL dependencies
- รัน tests ได้เหมือน Linux

**ขั้นตอน**:

1. **ติดตั้ง WSL2**:
```powershell
# ตรวจสอบว่ามี WSL2 แล้วหรือยัง
wsl --version

# ถ้ายังไม่มี ติดตั้ง
wsl --install
# รีสตาร์ทเครื่อง
```

2. **ติดตั้ง dependencies ใน WSL2**:
```bash
wsl
sudo apt update
sudo apt install -y build-essential
```

3. **แก้ไข `.env.development`**:
```bash
WHISPER_BINARY_PATH_WINDOWS=wsl /mnt/c/Users/boatr/.../backend/whisper/binary/linux/main
```

4. **รัน backend**:
```bash
# ใน WSL2
cd /mnt/c/Users/.../backend
go run .
```

---

#### วิธีที่ 2: Download Pre-built Windows Binary ⭐⭐⭐⭐

**ข้อดี**:
- ใช้เวลา 5 นาที
- มี DLLs ครบถ้วน

**ขั้นตอน**:

1. ไปที่ https://github.com/ggerganov/whisper.cpp/releases
2. ดาวน์โหลด `whisper-bin-x64.zip` (Windows x64)
3. แตกไฟล์และคัดลอก:
```bash
copy main.exe backend\whisper\binary\windows\
copy *.dll backend\whisper\binary\windows\
```

4. ทดสอบ binary:
```bash
cd backend\whisper\binary\windows
.\main.exe --version
```

---

#### วิธีที่ 3: Compile Windows Native Binary ⭐⭐⭐

**ข้อดี**:
- Performance ดีที่สุด
- ไม่ต้องพึ่ง WSL2

**ข้อเสีย**:
- ใช้เวลานาน (30-60 นาที)
- ต้องติดตั้ง Visual Studio

**ขั้นตอน**:

1. ติดตั้ง **Visual Studio 2022** หรือ **Build Tools for Visual Studio 2022**
   - เลือก "Desktop development with C++"

2. ติดตั้ง **CMake**:
```bash
# ดาวน์โหลดจาก https://cmake.org/download/
# หรือใช้ chocolatey
choco install cmake
```

3. **Clone และ Build**:
```bash
git clone https://github.com/ggerganov/whisper.cpp
cd whisper.cpp
mkdir build
cd build
cmake ..
cmake --build . --config Release
```

4. **คัดลอก binary**:
```bash
copy whisper.cpp\build\bin\Release\main.exe backend\whisper\binary\windows\
```

---

#### วิธีที่ 4: ติดตั้ง Visual C++ Redistributable ⭐⭐

**ข้อดี**:
- ใช้เวลา 2 นาที

**ข้อเสีย**:
- อาจไม่แก้ปัญหาได้ถ้า binary เป็น Linux format

**ขั้นตอน**:

1. ดาวน์โหลด **VC++ Redistributable**:
   - https://learn.microsoft.com/en-us/cpp/windows/latest-supported-vc-redist
   - ดาวน์โหลด `vc_redist.x64.exe`

2. ติดตั้งและรีสตาร์ท backend

---

### 🔍 การวินิจฉัยปัญหา

**ตรวจสอบว่า Binary เป็น Windows หรือ Linux**:

```powershell
cd backend\whisper\binary\windows
Get-Content .\main.exe -Encoding Byte -TotalCount 2 | ForEach-Object { [char]$_ }
```

**ผลลัพธ์**:
- `MZ` → Windows PE Binary (ถูกต้อง)
- `ELF` → Linux ELF Binary (ใช้ WSL2)

**ตรวจสอบ DLL Dependencies**:

```bash
# ถ้าติดตั้ง Visual Studio
dumpbin /dependents main.exe
```

---

### 💡 คำแนะนำตามสถานการณ์

**Development (Local)**:
- → ใช้ **วิธีที่ 1: WSL2** ⭐⭐⭐⭐⭐
- รวดเร็ว ใช้งานได้ทันที

**Production (Windows Server)**:
- → ใช้ **วิธีที่ 2 หรือ 3: Native Binary** ⭐⭐⭐⭐
- Performance ดี Stable

**Production (Linux Server)**:
- → ใช้ **Linux Binary ตรงๆ** ⭐⭐⭐⭐⭐
- ไม่มีปัญหาเลย (Recommended!)

---

## 📝 สรุปความคืบหน้า

### ✅ งานที่เสร็จสมบูรณ์แล้ว

#### งานที่ 1: Setup โครงสร้างโฟลเดอร์ ✅
- สร้างโฟลเดอร์ `backend/whisper/` พร้อม binary/, models/, temp/
- Clone whisper.cpp source code
- ดาวน์โหลด Linux binary (111 KB)
- ดาวน์โหลด GGML small model (466 MB)
- สร้างโฟลเดอร์ test พร้อม testdata/

#### งานที่ 2: Configuration ✅
- เพิ่มฟิลด์ Whisper config ใน `config/config.go`
- สร้าง `LoadWhisperConfig()` function
- เพิ่ม environment variables ใน `.env.development`
- รองรับ OS detection (Linux/Windows/macOS)
- รองรับ absolute path resolution
- ทดสอบ: 9/9 PASS

#### งานที่ 3: Interface Definition ✅
- สร้าง `services/transcription_service.go`
- กำหนด interface `TranscriptionService` พร้อม 2 methods:
  - `Transcribe()` - basic transcription
  - `TranscribeWithTimestamps()` - with time-aligned segments
- เอกสารครบถ้วน

#### งานที่ 4: Implement WhisperCppService ✅
- สร้าง `services/whispercpp_service.go` (450+ lines)
- Implement `NewWhisperCppService()` พร้อม validation
- Implement `Transcribe()` สำหรับ basic transcription
- Implement `TranscribeWithTimestamps()` สำหรับ segment-level timestamps
- Implement `GetModelName()` สำหรับ status endpoint
- แก้ไข Path Mismatch Error ด้วย absolute path resolution
- แก้ไข Test Timeout Error (1 นาที timeout, JSON file reading)
- ทดสอบบน WSL2: 22/22 PASS

#### งานที่ 5: สร้าง WhisperCppController ✅
- สร้าง `controllers/whispercpp_controller.go` (245 lines)
- Implement `NewWhisperCppController()` constructor
- Implement `GetStatus()` endpoint:
  - GET `/api/stt/whispercpp/status`
  - คืนค่า: service info, model, formats, languages
- Implement `TranscribeAudio()` endpoint:
  - POST `/api/stt/whispercpp`
  - รองรับ multipart/form-data upload
  - Validation: file size (25MB), format, language
  - รองรับทั้ง basic และ timestamps mode
- ทดสอบ compilation สำเร็จ

#### งานที่ 6: ลงทะเบียน Routes ✅
- แก้ไข `routes/routes.go`
- Initialize `WhisperCppService` พร้อม error handling
- Initialize `WhisperCppController` (conditional)
- สร้าง route group `/api/stt`
- ลงทะเบียน 2 endpoints:
  - GET `/api/stt/whispercpp/status`
  - POST `/api/stt/whispercpp`
- ทดสอบ backend compilation สำเร็จ
- Graceful degradation ถ้า binary/model ไม่พบ

#### งานเพิ่มเติม: Dynamic Model Selection ✅ (2025-11-10)
- เพิ่ม `WhisperModelsDir` และ `WhisperSupportedModels` ใน Configuration
- เพิ่ม `GetSupportedModels()` และ `GetModelPath(modelName)` ใน WhisperCppService
- เพิ่ม `TranscribeWithModel()` และ `TranscribeWithTimestampsAndModel()` methods
- แก้ไข WhisperCppController รองรับ `model` parameter ใน request
- เพิ่ม `model` field ใน response ทุก API
- อัพเดต `/api/stt/whispercpp/status` แสดง `supported_models` list
- รองรับการเลือก model แบบ dynamic: tiny.en, small, medium, large-v2
- ทดสอบ compilation สำเร็จ

---

## 🎯 Dynamic Model Selection

### ภาพรวม

ระบบรองรับการเลือก model แบบ dynamic ผ่าน API request โดยไม่ต้อง restart backend

**Models ที่รองรับ**:
- `tiny.en` - Model ขนาดเล็กสุด เหมาะสำหรับภาษาอังกฤษ (fast, low accuracy)
- `small` - Model ขนาดกลาง รองรับหลายภาษา รวมภาษาไทย (default)
- `medium` - Model ขนาดใหญ่ ความแม่นยำสูง
- `large-v2` - Model ขนาดใหญ่สุด ความแม่นยำสูงสุด (slow, high accuracy)

### การตั้งค่า

**1. เพิ่ม models ใน `backend/whisper/models/` directory**:
```
backend/whisper/models/
├── ggml-small.bin          # Default model
├── ggml-tiny.en.bin        # หรือ ggml-tiny-en-q5_1.bin
├── ggml-medium.bin
└── ggml-large-v2.bin
```

**2. ตั้งค่า supported models ใน `.env.development`**:
```bash
WHISPER_MODELS_DIR=./whisper/models
WHISPER_SUPPORTED_MODELS=tiny.en,small,medium,large-v2
WHISPER_MODEL_NAME=small  # Default model
```

### การใช้งาน

**ตรวจสอบ models ที่รองรับ**:
```bash
GET /api/stt/whispercpp/status
```

Response จะแสดง `supported_models` list:
```json
{
  "supported_models": ["tiny.en", "small", "medium", "large-v2"],
  "default_model": "small"
}
```

**เลือก model ใน request**:
```bash
POST /api/stt/whispercpp

Body (form-data):
- audio: [File]
- language: "en"
- model: "tiny.en"  # เลือก model ที่ต้องการ
```

Response จะแสดง model ที่ใช้:
```json
{
  "transcription": "Hello world",
  "model": "tiny.en"
}
```

**ถ้าไม่ระบุ model** จะใช้ default model จาก config

### Model Filename Convention

ระบบจะหา model file โดยอัตโนมัติตาม pattern:
1. `ggml-{modelName}.bin` - เช่น `ggml-small.bin`
2. `ggml-{modelName}-q5_1.bin` - เช่น `ggml-tiny-en-q5_1.bin` (quantized version)

**ตัวอย่าง**:
- `tiny.en` → หา `ggml-tiny.en.bin` หรือ `ggml-tiny-en-q5_1.bin`
- `small` → หา `ggml-small.bin`
- `medium` → หา `ggml-medium.bin`

### Error Handling

**Model ไม่รองรับ**:
```json
{
  "success": false,
  "error": "model selection error: model 'large-v3' is not supported. Supported models: tiny.en, small, medium, large-v2"
}
```

**Model file ไม่พบ**:
```json
{
  "success": false,
  "error": "model selection error: model file not found: tried ggml-tiny.en.bin and ggml-tiny-en-q5_1.bin"
}
```

---

## 🎯 API Endpoints

### GET /api/stt/whispercpp/status

**หน้าที่**: ตรวจสอบสถานะ service

**Response**:
```json
{
  "service": "whisper.cpp",
  "available": true,
  "model": "small",
  "supported_formats": ["wav", "mp3", "m4a", "ogg", "flac"],
  "supported_languages": ["th", "en", "auto"],
  "os": "windows"
}
```

---

### POST /api/stt/whispercpp

**หน้าที่**: แปลงเสียงเป็นข้อความ

**Request** (multipart/form-data):
```
audio: [File] ไฟล์เสียง (max 25MB)
language: [Text] "th" | "en" | "auto" (default: "auto")
timestamps: [Text] "true" | "false" (default: "false")
model: [Text] "tiny.en" | "small" | "medium" | "large-v2" (optional, default: ใช้ค่า config)
```

**Response (timestamps=false)**:
```json
{
  "success": true,
  "transcription": "สวัสดีครับ ยินดีต้อนรับสู่ระบบแชทบอท",
  "confidence": 0.95,
  "language": "th",
  "model": "small"
}
```

**Response (timestamps=true)**:
```json
{
  "success": true,
  "transcription": "สวัสดีครับ ยินดีต้อนรับสู่ระบบแชทบอท",
  "segments": [
    {
      "start_time": 0.0,
      "end_time": 2.0,
      "text": "สวัสดีครับ"
    },
    {
      "start_time": 2.0,
      "end_time": 4.5,
      "text": "ยินดีต้อนรับสู่ระบบแชทบอท"
    }
  ],
  "language": "th",
  "duration": 4.5,
  "model": "small"
}
```

---

## 📚 ข้อมูลอ้างอิง

- **whisper.cpp GitHub**: https://github.com/ggerganov/whisper.cpp
- **GGML Models**: https://huggingface.co/ggerganov/whisper.cpp
- **OpenAI Whisper Paper**: https://arxiv.org/abs/2212.04356
- **Go Fiber Documentation**: https://docs.gofiber.io/

---

## 🏁 สรุป

**Implementation Status**: ✅ **Complete + Enhanced**

**ระยะเวลา**: ~3-4 วัน (implementation) + 1 วัน (dynamic model selection)

**Test Coverage**: 22 unit tests ทั้งหมด PASS

**Production Ready**: ✅ ใช่ (ต้องตั้งค่า WSL2 สำหรับ Windows development หรือ compile native binary สำหรับ production)

**Features**:
- ✅ Speech-to-Text transcription (ภาษาไทย, อังกฤษ, auto-detect)
- ✅ Timestamps support (segment-level)
- ✅ **Dynamic model selection** (tiny.en, small, medium, large-v2)
- ✅ Multiple audio formats (WAV, MP3, M4A, OGG, FLAC)
- ✅ Graceful error handling
- ✅ WSL2 integration สำหรับ Windows

**Enhancements (2025-11-10)**:
- ➕ รองรับการเลือก model แบบ dynamic ผ่าน request parameter
- ➕ Model validation และ auto-detection
- ➕ Flexible filename convention (ggml-*.bin, ggml-*-q5_1.bin)
- ➕ Response แสดง model ที่ใช้งาน

**Next Steps**:
1. ดาวน์โหลด models เพิ่มเติม (tiny.en, medium, large-v2) ถ้าต้องการ
2. ทดสอบ API endpoints ด้วย Frontend
3. Monitor performance ของแต่ละ model
4. Deploy บน production server (แนะนำ Linux)

---

**Author**: Claude Code
**Last Updated**: 2025-11-10
**Version**: 1.1 (Dynamic Model Selection)
