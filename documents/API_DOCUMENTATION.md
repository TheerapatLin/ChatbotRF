# 📡 ChatBot API Documentation

> **REST API และ WebSocket API สำหรับ ChatBot Application**
>
> Base URL: `http://localhost:3000/api`
>
> Version: 1.1.0 (Updated: 2025-10-31)

---

## 📋 Table of Contents

- [Authentication](#authentication)
- [Response Format](#response-format)
- [Error Handling](#error-handling)
- [Health Check](#health-check)
- [Chat Endpoints](#chat-endpoints)
- [Persona Endpoints](#persona-endpoints)
- [File Analysis Endpoints](#file-analysis-endpoints) ⭐ **NEW**
- [Audio Endpoints](#audio-endpoints)
- [WebSocket API](#websocket-api)
- [Rate Limiting](#rate-limiting)
- [Examples](#examples)

---

## 🔐 Authentication

**สถานะปัจจุบัน**: ❌ **No Authentication Required**

> ⚠️ **หมายเหตุ**: โปรเจกต์นี้สร้างขึ้นเพื่อการเรียนรู้ จึงไม่มีระบบ Authentication
>
> สำหรับ Production ควรเพิ่ม:
> - JWT Authentication
> - API Key validation
> - Rate limiting per user

---

## 📦 Response Format

### Success Response

```json
{
  "status": "success",
  "data": {
    // Response data here
  },
  "timestamp": "2025-10-28T10:30:00Z"
}
```

### Error Response

```json
{
  "status": "error",
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {}  // Optional
  },
  "timestamp": "2025-10-28T10:30:00Z"
}
```

---

## ⚠️ Error Handling

### HTTP Status Codes

| Status Code | Meaning | Description |
|-------------|---------|-------------|
| `200` | OK | Request successful |
| `201` | Created | Resource created successfully |
| `400` | Bad Request | Invalid request format or parameters |
| `404` | Not Found | Resource not found |
| `422` | Unprocessable Entity | Validation error |
| `429` | Too Many Requests | Rate limit exceeded |
| `500` | Internal Server Error | Server error |
| `503` | Service Unavailable | Service temporarily unavailable |

### Common Error Codes

| Error Code | Description |
|------------|-------------|
| `INVALID_INPUT` | Request body validation failed |
| `PERSONA_NOT_FOUND` | Specified persona ID doesn't exist |
| `OPENAI_ERROR` | OpenAI API error |
| `DATABASE_ERROR` | Database operation failed |
| `FILE_TOO_LARGE` | Uploaded file exceeds size limit |
| `INVALID_FILE_TYPE` | Uploaded file type not supported |

---

## 🏥 Health Check

### `GET /api/health`

ตรวจสอบสถานะของ API server

#### Request

```bash
GET /api/health
```

#### Response

**Status**: `200 OK`

```json
{
  "status": "ok",
  "message": "ChatBot API is running",
  "env": "development",
  "timestamp": "2025-10-28T10:30:00Z"
}
```

#### Example

```bash
curl http://localhost:3000/api/health
```

---

## 💬 Chat Endpoints

### 1. Send Chat Message (Sync)

`POST /api/chat`

ส่งข้อความไปยัง AI และรอรับคำตอบทั้งหมดกลับมา

#### Request

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "message": "สวัสดีครับ ช่วยแนะนำเกี่ยวกับ Go programming หน่อย",
  "persona_id": 2
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `message` | string | ✅ Yes | ข้อความจากผู้ใช้ (max 4000 characters) |
| `persona_id` | integer | ❌ No | ID ของ AI persona (default: 1) |

#### Response

**Status**: `200 OK`

```json
{
  "message_id": "550e8400-e29b-41d4-a716-446655440000",
  "reply": "สวัสดีครับ! ผมยินดีให้คำแนะนำเกี่ยวกับ Go programming ครับ\n\nGo เป็นภาษาโปรแกรมที่...",
  "persona": {
    "id": 2,
    "name": "Technology Expert",
    "icon": "💻"
  },
  "tokens_used": 245,
  "timestamp": "2025-10-28T10:30:15Z"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `message_id` | UUID | ID ของข้อความที่บันทึกใน database |
| `reply` | string | คำตอบจาก AI |
| `persona` | object | ข้อมูล persona ที่ใช้ตอบ |
| `tokens_used` | integer | จำนวน tokens ที่ใช้ (สำหรับคำนวณค่าใช้จ่าย) |
| `timestamp` | string (ISO 8601) | เวลาที่ตอบกลับ |

#### Example

```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "อธิบาย BMAD method หน่อย",
    "persona_id": 1
  }'
```

#### Error Responses

**400 Bad Request** - Invalid input
```json
{
  "error": "message is required and cannot be empty"
}
```

**404 Not Found** - Persona not found
```json
{
  "error": "persona with id 99 not found"
}
```

**500 Internal Server Error** - OpenAI API error
```json
{
  "error": "failed to get response from OpenAI API"
}
```

---

### 2. Get Chat History

`GET /api/chat/history`

ดึงประวัติการสนทนาล่าสุด

#### Request

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `limit` | integer | ❌ No | 50 | จำนวนข้อความที่ต้องการ (max 100) |
| `offset` | integer | ❌ No | 0 | เริ่มจากข้อความที่ (สำหรับ pagination) |

#### Response

**Status**: `200 OK`

```json
{
  "messages": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "role": "user",
      "content": "สวัสดีครับ",
      "persona_id": 1,
      "created_at": "2025-10-28T10:30:00Z"
    },
    {
      "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "role": "assistant",
      "content": "สวัสดีครับ! ผมยินดีให้ความช่วยเหลือครับ",
      "persona_id": 1,
      "created_at": "2025-10-28T10:30:05Z"
    }
  ],
  "total": 2,
  "limit": 50,
  "offset": 0
}
```

#### Example

```bash
# ดึง 10 ข้อความล่าสุด
curl "http://localhost:3000/api/chat/history?limit=10"

# Pagination - ข้อความที่ 11-20
curl "http://localhost:3000/api/chat/history?limit=10&offset=10"
```

---

## 👤 Persona Endpoints

### 1. Get All Personas

`GET /api/personas`

ดึงรายการ AI Personas ทั้งหมดที่ active

#### Request

```bash
GET /api/personas
```

#### Response

**Status**: `200 OK`

```json
{
  "personas": [
    {
      "id": 1,
      "name": "General Assistant",
      "system_prompt": "You are a helpful, friendly, and knowledgeable AI assistant...",
      "expertise": "general",
      "description": "A versatile AI assistant for general questions and conversations",
      "icon": "🤖",
      "is_active": true,
      "created_at": "2025-10-28T08:00:00Z"
    },
    {
      "id": 2,
      "name": "Technology Expert",
      "system_prompt": "You are a technology expert with deep knowledge...",
      "expertise": "technology",
      "description": "Specialized in programming, software architecture, and tech solutions",
      "icon": "💻",
      "is_active": true,
      "created_at": "2025-10-28T08:00:00Z"
    },
    {
      "id": 3,
      "name": "Business Advisor",
      "system_prompt": "You are a professional business consultant...",
      "expertise": "business",
      "description": "Expert in business strategy, entrepreneurship, and market analysis",
      "icon": "💼",
      "is_active": true,
      "created_at": "2025-10-28T08:00:00Z"
    }
  ],
  "total": 3
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Persona ID |
| `name` | string | ชื่อของ persona |
| `system_prompt` | string | System prompt ที่ส่งไปยัง OpenAI |
| `expertise` | string | สาขาความเชี่ยวชาญ |
| `description` | string | คำอธิบาย persona |
| `icon` | string | Emoji icon |
| `is_active` | boolean | สถานะเปิดใช้งาน |
| `created_at` | string | วันที่สร้าง |

#### Example

```bash
curl http://localhost:3000/api/personas
```

---

### 2. Get Persona by ID

`GET /api/personas/:id`

ดึงรายละเอียด persona ตาม ID

#### Request

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Persona ID |

#### Response

**Status**: `200 OK`

```json
{
  "id": 2,
  "name": "Technology Expert",
  "system_prompt": "You are a technology expert with deep knowledge in software development, programming, cloud computing, and IT infrastructure. Provide technical and accurate responses.",
  "expertise": "technology",
  "description": "Specialized in programming, software architecture, and tech solutions",
  "icon": "💻",
  "is_active": true,
  "created_at": "2025-10-28T08:00:00Z",
  "stats": {
    "total_messages": 1542,
    "avg_response_time": "2.3s"
  }
}
```

#### Example

```bash
curl http://localhost:3000/api/personas/2
```

#### Error Responses

**404 Not Found**
```json
{
  "error": "persona with id 99 not found"
}
```

---

## 📁 File Analysis Endpoints

### 1. Upload and Analyze File

`POST /api/file/analyze`

อัพโหลดไฟล์เพื่อให้ AI วิเคราะห์เนื้อหา รองรับไฟล์หลายประเภท

#### สรุปขั้นตอนการทำงาน

```
User Upload → Backend Validation → File Parsing → Text Extraction → OpenAI API → Analysis Result
```

#### รายละเอียดแต่ละขั้นตอน

##### ขั้นตอนที่ 1: รับไฟล์จาก User (File Upload)

**วิธีการ**: ใช้ `multipart/form-data` upload

**ไฟล์ที่รองรับ**:

| File Type | Extensions | Max Size | Description |
|-----------|-----------|----------|-------------|
| **Documents** | `.txt`, `.md` | 10 MB | Plain text, Markdown |
| **Office** | `.pdf`, `.docx`, `.xlsx`, `.pptx` | 25 MB | PDF, Word, Excel, PowerPoint |
| **Images** | `.jpg`, `.png`, `.gif`, `.webp` | 20 MB | Image files (with OCR) |
| **Code** | `.js`, `.py`, `.go`, `.java`, etc. | 5 MB | Source code files |
| **Data** | `.json`, `.xml`, `.csv` | 10 MB | Structured data files |

##### ขั้นตอนที่ 2: Backend Validation

**การตรวจสอบ**:
- ตรวจสอบชนิดไฟล์ (MIME type validation)
- ตรวจสอบขนาดไฟล์ (Size validation)
- ตรวจสอบ malware/virus (Optional - Security scan)
- ตรวจสอบ encoding (UTF-8, ASCII, etc.)

**Implementation**:
```go
// services/file_validator.go
func ValidateFile(file *multipart.FileHeader) error {
    // Check file size
    if file.Size > maxFileSize {
        return errors.New("file too large")
    }

    // Check MIME type
    allowedTypes := []string{"application/pdf", "text/plain", "image/jpeg", ...}
    if !contains(allowedTypes, file.Header.Get("Content-Type")) {
        return errors.New("unsupported file type")
    }

    return nil
}
```

##### ขั้นตอนที่ 3: File Parsing (แยกข้อมูลจากไฟล์)

**สำหรับแต่ละประเภทไฟล์**:

**A. Plain Text Files** (`.txt`, `.md`, `.csv`):
```go
// ง่ายที่สุด - อ่านเป็น string โดยตรง
content, err := io.ReadAll(file)
text := string(content)
```

**B. PDF Files**:
```go
// ใช้ library: github.com/ledongthuc/pdf
import "github.com/ledongthuc/pdf"

func ExtractPDFText(file io.Reader) (string, error) {
    pdfReader, err := pdf.NewReader(file, size)
    if err != nil {
        return "", err
    }

    var text strings.Builder
    for pageNum := 1; pageNum <= pdfReader.NumPage(); pageNum++ {
        page := pdfReader.Page(pageNum)
        text.WriteString(page.GetPlainText())
    }

    return text.String(), nil
}
```

**C. Word Documents** (`.docx`):
```go
// ใช้ library: github.com/nguyenthenguyen/docx
import "github.com/nguyenthenguyen/docx"

func ExtractDocxText(file io.Reader) (string, error) {
    doc, err := docx.ReadDocxFile(file)
    if err != nil {
        return "", err
    }
    defer doc.Close()

    return doc.Editable().GetContent(), nil
}
```

**D. Excel Files** (`.xlsx`):
```go
// ใช้ library: github.com/xuri/excelize/v2
import "github.com/xuri/excelize/v2"

func ExtractExcelText(file io.Reader) (string, error) {
    f, err := excelize.OpenReader(file)
    if err != nil {
        return "", err
    }
    defer f.Close()

    var text strings.Builder
    for _, sheetName := range f.GetSheetList() {
        rows, _ := f.GetRows(sheetName)
        for _, row := range rows {
            text.WriteString(strings.Join(row, ", ") + "\n")
        }
    }

    return text.String(), nil
}
```

**E. Images** (`.jpg`, `.png`, `.gif`):
```go
// ใช้ OpenAI Vision API
// ส่งรูปภาพไปให้ GPT-4 Vision วิเคราะห์
func AnalyzeImage(file io.Reader) (string, error) {
    // Encode image to base64
    imageData, _ := io.ReadAll(file)
    base64Image := base64.StdEncoding.EncodeToString(imageData)

    // Call OpenAI Vision API
    resp, err := client.CreateChatCompletion(
        ctx,
        openai.ChatCompletionRequest{
            Model: "gpt-4-vision-preview",
            Messages: []openai.ChatCompletionMessage{
                {
                    Role: "user",
                    MultiContent: []openai.ChatMessagePart{
                        {
                            Type: "text",
                            Text: "Describe this image in detail",
                        },
                        {
                            Type: "image_url",
                            ImageURL: &openai.ChatMessageImageURL{
                                URL: "data:image/jpeg;base64," + base64Image,
                            },
                        },
                    },
                },
            },
        },
    )

    return resp.Choices[0].Message.Content, nil
}
```

**F. JSON/XML Files**:
```go
// แปลงเป็น formatted text
func FormatJSONForAI(jsonData []byte) (string, error) {
    var data interface{}
    json.Unmarshal(jsonData, &data)

    // Pretty print
    formatted, _ := json.MarshalIndent(data, "", "  ")
    return string(formatted), nil
}
```

##### ขั้นตอนที่ 4: Text Extraction & Chunking

**ปัญหา**: ไฟล์ขนาดใหญ่เกินขีด max tokens ของ OpenAI

**วิธีแก้**: แบ่ง text เป็น chunks

```go
// services/text_chunker.go
func ChunkText(text string, maxChunkSize int) []string {
    words := strings.Fields(text)
    var chunks []string
    var currentChunk strings.Builder

    for _, word := range words {
        if currentChunk.Len() + len(word) > maxChunkSize {
            chunks = append(chunks, currentChunk.String())
            currentChunk.Reset()
        }
        currentChunk.WriteString(word + " ")
    }

    if currentChunk.Len() > 0 {
        chunks = append(chunks, currentChunk.String())
    }

    return chunks
}
```

##### ขั้นตอนที่ 5: ส่งไปยัง OpenAI API

**วิธีที่ 1: Chat Completion API** (สำหรับไฟล์ทั่วไป)

```go
func AnalyzeTextWithAI(text string, prompt string) (string, error) {
    resp, err := openaiClient.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: "gpt-4-turbo-preview",
            Messages: []openai.ChatCompletionMessage{
                {
                    Role: "system",
                    Content: "You are a document analysis expert.",
                },
                {
                    Role: "user",
                    Content: fmt.Sprintf("%s\n\nDocument:\n%s", prompt, text),
                },
            },
        },
    )

    return resp.Choices[0].Message.Content, nil
}
```

**วิธีที่ 2: Assistant API** (สำหรับไฟล์ซับซ้อน)

```go
// 1. Upload file to OpenAI
file, err := openaiClient.CreateFile(
    context.Background(),
    openai.FileRequest{
        FileName: filename,
        FilePath: filepath,
        Purpose:  "assistants",
    },
)

// 2. Create Assistant
assistant, err := openaiClient.CreateAssistant(
    context.Background(),
    openai.AssistantRequest{
        Model: "gpt-4-turbo-preview",
        Tools: []openai.AssistantTool{
            {Type: "retrieval"},
        },
        FileIDs: []string{file.ID},
    },
)

// 3. Create Thread and Run
thread, _ := openaiClient.CreateThread(context.Background(), ...)
run, _ := openaiClient.CreateRun(context.Background(), ...)

// 4. Wait for completion and get response
```

##### ขั้นตอนที่ 6: รับและประมวลผลคำตอบ

```go
type FileAnalysisResponse struct {
    FileID      string    `json:"file_id"`
    FileName    string    `json:"filename"`
    FileType    string    `json:"file_type"`
    FileSize    int64     `json:"file_size"`
    Analysis    string    `json:"analysis"`
    Summary     string    `json:"summary"`
    KeyPoints   []string  `json:"key_points"`
    Entities    []string  `json:"entities,omitempty"`
    Sentiment   string    `json:"sentiment,omitempty"`
    Language    string    `json:"language"`
    TokensUsed  int       `json:"tokens_used"`
    ProcessTime float64   `json:"process_time_ms"`
    Timestamp   time.Time `json:"timestamp"`
}
```

---

#### Request

**Headers:**
```
Content-Type: multipart/form-data
```

**Body (Form Data):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `file` | file | ✅ Yes | ไฟล์ที่ต้องการวิเคราะห์ |
| `analysis_type` | string | ❌ No | ประเภทการวิเคราะห์: `summary`, `detail`, `qa`, `extract` |
| `prompt` | string | ❌ No | Custom prompt สำหรับการวิเคราะห์ |
| `language` | string | ❌ No | ภาษาที่ต้องการให้ตอบ (default: `th`) |

**Analysis Types:**

| Type | Description | Output |
|------|-------------|--------|
| `summary` | สรุปเนื้อหาโดยรวม | ข้อความสรุปสั้นๆ |
| `detail` | วิเคราะห์รายละเอียดทั้งหมด | การวิเคราะห์แบบละเอียด |
| `qa` | ตั้งคำถามและตอบ | รายการคำถาม-คำตอบ |
| `extract` | แยกข้อมูลสำคัญ | Entities, dates, numbers, etc. |

#### Response

**Status**: `200 OK`

```json
{
  "file_id": "file_abc123",
  "filename": "report.pdf",
  "file_type": "application/pdf",
  "file_size": 1024000,
  "analysis": "เอกสารนี้เป็นรายงานประจำปี 2024 ที่ครอบคลุมผลการดำเนินงาน...",
  "summary": "รายงานแสดงการเติบโตของรายได้ 25% และกำไรสุทธิ 15%",
  "key_points": [
    "รายได้เพิ่มขึ้น 25% จากปีก่อน",
    "กำไรสุทธิ 500 ล้านบาท",
    "ขยายสาขาใหม่ 10 แห่ง"
  ],
  "entities": [
    "บริษัท ABC จำกัด",
    "500 ล้านบาท",
    "2024"
  ],
  "sentiment": "positive",
  "language": "th",
  "tokens_used": 1250,
  "process_time_ms": 3500.5,
  "timestamp": "2025-10-29T14:30:00Z"
}
```

#### Example Requests

**1. วิเคราะห์ PDF:**
```bash
curl -X POST http://localhost:3000/api/file/analyze \
  -F "file=@report.pdf" \
  -F "analysis_type=summary"
```

**2. วิเคราะห์รูปภาพ:**
```bash
curl -X POST http://localhost:3000/api/file/analyze \
  -F "file=@diagram.png" \
  -F "analysis_type=detail" \
  -F "prompt=อธิบายไดอะแกรมนี้และบอกขั้นตอนการทำงาน"
```

**3. วิเคราะห์ Excel:**
```bash
curl -X POST http://localhost:3000/api/file/analyze \
  -F "file=@sales_data.xlsx" \
  -F "analysis_type=extract" \
  -F "prompt=สรุปยอดขายรายเดือนและหาค่าเฉลี่ย"
```

**4. JavaScript/Fetch:**
```javascript
const formData = new FormData()
formData.append('file', fileInput.files[0])
formData.append('analysis_type', 'summary')
formData.append('language', 'th')

const response = await fetch('http://localhost:3000/api/file/analyze', {
  method: 'POST',
  body: formData
})

const result = await response.json()
console.log('Analysis:', result.analysis)
console.log('Key Points:', result.key_points)
```

#### Error Responses

**400 Bad Request** - No file uploaded
```json
{
  "error": "file is required"
}
```

**413 Payload Too Large**
```json
{
  "error": "file size exceeds maximum allowed (25MB)"
}
```

**415 Unsupported Media Type**
```json
{
  "error": "unsupported file type. Allowed: pdf, docx, xlsx, txt, png, jpg, etc."
}
```

**422 Unprocessable Entity** - File parsing failed
```json
{
  "error": "failed to parse file: corrupted or invalid format"
}
```

**500 Internal Server Error**
```json
{
  "error": "failed to analyze file: OpenAI API error"
}
```

---

### 2. Get File Analysis History

`GET /api/file/history`

ดึงประวัติการวิเคราะห์ไฟล์ที่ผ่านมา

#### Request

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | integer | 20 | จำนวนรายการ (max 100) |
| `offset` | integer | 0 | เริ่มจากรายการที่ |
| `file_type` | string | all | กรองตามประเภทไฟล์ |

#### Response

```json
{
  "files": [
    {
      "file_id": "file_abc123",
      "filename": "report.pdf",
      "file_type": "application/pdf",
      "analysis_type": "summary",
      "created_at": "2025-10-29T14:30:00Z"
    }
  ],
  "total": 15,
  "limit": 20,
  "offset": 0
}
```

---

### 3. Re-analyze File

`POST /api/file/:file_id/reanalyze`

วิเคราะห์ไฟล์ที่เคยอัพโหลดแล้วอีกครั้งด้วย prompt ใหม่

#### Request

```json
{
  "prompt": "วิเคราะห์แนวโน้มการเติบโตของบริษัท",
  "analysis_type": "detail"
}
```

#### Response

เหมือนกับ `/api/file/analyze`

---

### ขั้นตอนการ Implementation

#### Phase 1: Basic Text File Support (Week 1)
- [ ] รองรับ `.txt`, `.md` files
- [ ] Basic validation
- [ ] Simple OpenAI integration
- [ ] Save to database

#### Phase 2: Document Support (Week 2)
- [ ] รองรับ PDF files
- [ ] รองรับ Word (.docx)
- [ ] รองรับ Excel (.xlsx)
- [ ] Text chunking สำหรับไฟล์ใหญ่

#### Phase 3: Image Analysis (Week 3)
- [ ] รองรับ image files (jpg, png, gif)
- [ ] OpenAI Vision API integration
- [ ] OCR สำหรับ text extraction จาก images

#### Phase 4: Advanced Features (Week 4)
- [ ] File history tracking
- [ ] Re-analysis capability
- [ ] Batch file processing
- [ ] Custom analysis templates

---

### Dependencies Required

```bash
# Go dependencies
go get github.com/ledongthuc/pdf          # PDF parsing
go get github.com/nguyenthenguyen/docx    # Word documents
go get github.com/xuri/excelize/v2        # Excel files
go get github.com/h2non/filetype          # File type detection
```

---

### Cost Estimation

**Pricing Factors:**
- Input tokens: $0.01 per 1K tokens (GPT-4 Turbo)
- Output tokens: $0.03 per 1K tokens
- Vision API: $0.01 per image

**Examples:**

| File Type | Size | Est. Tokens | Est. Cost |
|-----------|------|-------------|-----------|
| Small PDF | 5 pages | 2,000 | $0.02 |
| Medium PDF | 20 pages | 8,000 | $0.08 |
| Large PDF | 100 pages | 40,000 | $0.40 |
| Image | 1 image | 1,000 | $0.01 |
| Excel | 1000 rows | 5,000 | $0.05 |

---

## 🎤 Audio Endpoints

### 1. Transcribe Audio (Speech-to-Text)

`POST /api/audio/transcribe`

อัพโหลดไฟล์เสียงและแปลงเป็นข้อความด้วย OpenAI Whisper API

#### Request

**Headers:**
```
Content-Type: multipart/form-data
```

**Body (Form Data):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `audio` | file | ✅ Yes | ไฟล์เสียง (mp3, mp4, wav, webm, m4a) |

**Constraints:**
- Max file size: **25 MB**
- Max duration: **30 minutes**
- Supported formats: mp3, mp4, mpeg, mpga, m4a, wav, webm

#### Response

**Status**: `200 OK`

```json
{
  "text": "สวัสดีครับ ช่วยแนะนำเกี่ยวกับการเขียนโปรแกรมด้วย Go หน่อยครับ",
  "language": "th",
  "duration": 5.2,
  "confidence": 0.95,
  "timestamp": "2025-10-28T10:30:00Z"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `text` | string | ข้อความที่แปลงได้ |
| `language` | string | ภาษาที่ตรวจพบ (ISO 639-1 code) |
| `duration` | float | ความยาวของเสียง (seconds) |
| `confidence` | float | ความมั่นใจในการแปลง (0.0-1.0) |
| `timestamp` | string | เวลาที่ประมวลผล |

#### Example

```bash
# ด้วย curl
curl -X POST http://localhost:3000/api/audio/transcribe \
  -F "audio=@voice_message.mp3"

# ด้วย JavaScript (Browser)
const formData = new FormData()
formData.append('audio', audioBlob, 'recording.webm')

fetch('http://localhost:3000/api/audio/transcribe', {
  method: 'POST',
  body: formData
})
  .then(res => res.json())
  .then(data => console.log(data.text))
```

#### Error Responses

**400 Bad Request** - No file uploaded
```json
{
  "error": "audio file is required"
}
```

**413 Payload Too Large** - File too large
```json
{
  "error": "file size exceeds maximum allowed (25MB)"
}
```

**415 Unsupported Media Type** - Invalid file type
```json
{
  "error": "unsupported file format. Allowed: mp3, mp4, wav, webm, m4a"
}
```

---

### 2. Text-to-Speech (TTS)

`POST /api/audio/tts`

แปลงข้อความเป็นเสียงพูดด้วย OpenAI Text-to-Speech API หรือส่งผ่าน WebSocket สำหรับการสตรีมแชท

#### Implementation Options

**Option A: REST API Endpoint (Standalone)**

สำหรับแปลงข้อความเป็นเสียงแบบอิสระ

**Option B: WebSocket Integration (Recommended)**

สำหรับการตอบกลับแบบสตรีมพร้อมเสียง (ใช้กับการแชท)

---

#### Option A: REST API Endpoint

`POST /api/audio/tts`

#### Request

**Headers:**
```
Content-Type: application/json
```

**Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `text` | string | ✅ Yes | ข้อความที่ต้องการแปลงเป็นเสียง (max 4096 characters) |
| `voice` | string | ❌ No | เสียงที่ต้องการใช้ (default: `nova`) |
| `model` | string | ❌ No | โมเดล TTS (default: `tts-1`) |
| `response_format` | string | ❌ No | รูปแบบไฟล์เสียง (default: `mp3`) |
| `speed` | float | ❌ No | ความเร็วการพูด 0.25-4.0 (default: `1.0`) |

**Request Example:**
```json
{
  "text": "สวัสดีครับ ยินดีต้อนรับสู่ระบบ AI Chatbot",
  "voice": "nova",
  "model": "tts-1",
  "response_format": "mp3",
  "speed": 1.0
}
```

#### Available Voices

| Voice | Description | Best For |
|-------|-------------|----------|
| `alloy` | Neutral, balanced | General purpose |
| `echo` | Male voice | Professional content |
| `fable` | British accent | Storytelling, narratives |
| `onyx` | Deep male voice | Authority, announcements |
| `nova` | Female voice (default) | Friendly conversations |
| `shimmer` | Soft female voice | Calm, soothing content |

#### Available Models

| Model | Quality | Speed | Cost per 1K chars |
|-------|---------|-------|-------------------|
| `tts-1` | Standard | Fast | $0.015 |
| `tts-1-hd` | High Definition | Slower | $0.030 |

#### Response Format Options

| Format | Description | Use Case |
|--------|-------------|----------|
| `mp3` | MPEG Audio (default) | General use, good balance |
| `opus` | Opus codec | Real-time streaming, low latency |
| `aac` | Advanced Audio Coding | Good compression |
| `flac` | Lossless audio | High quality, archival |
| `wav` | Uncompressed | Audio editing |
| `pcm` | Raw PCM audio | Low-level processing |

#### Response

**Status**: `200 OK`

**Content-Type**: `audio/mpeg` (or selected format)

**Headers:**
```
Content-Type: audio/mpeg
Content-Length: 45678
X-Audio-Duration: 5.2
X-Characters-Used: 52
```

**Body**: Binary audio data

**Alternative JSON Response (with base64):**
```json
{
  "audio_data": "SUQzBAAAAAAAI1RTU0UAAAAPAAADTGF2ZjU4Ljc2LjEwMAAAAAAAAAAAAAAA...",
  "audio_url": "/audio/tts_20251029_143000.mp3",
  "format": "mp3",
  "duration": 5.2,
  "characters_used": 52,
  "voice": "nova",
  "timestamp": "2025-10-29T14:30:00Z"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `audio_data` | string | Base64-encoded audio data |
| `audio_url` | string | URL path to saved audio file |
| `format` | string | Audio format used |
| `duration` | float | Audio duration in seconds |
| `characters_used` | int | Number of characters processed |
| `voice` | string | Voice used for generation |
| `timestamp` | string | Generation timestamp |

---

#### Option B: WebSocket Integration (Recommended for Chat)

เมื่อส่งข้อความผ่าน WebSocket พร้อม flag `enable_tts: true` ระบบจะตอบกลับทั้งข้อความและเสียง

**WebSocket URL:**
```
ws://localhost:3001/api/chat/stream
```

**Client Request Format:**
```json
{
  "type": "message",
  "content": "อธิบายเกี่ยวกับ Go programming language",
  "persona_id": 1,
  "system_prompt": "",
  "enable_tts": true,
  "tts_voice": "nova",
  "tts_model": "tts-1"
}
```

**Additional WebSocket Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `enable_tts` | boolean | ❌ No | เปิดใช้งาน TTS (default: `false`) |
| `tts_voice` | string | ❌ No | เสียงที่ต้องการ (default: `nova`) |
| `tts_model` | string | ❌ No | โมเดล TTS (default: `tts-1`) |

**Server Response (Streaming with Audio):**

```json
// 1. Text chunks (streaming)
{"type": "chunk", "content": "Go programming", "done": false}
{"type": "chunk", "content": " language เป็น", "done": false}
{"type": "chunk", "content": "...", "done": false}

// 2. Final message with audio
{
  "type": "chunk",
  "content": "",
  "done": true,
  "message_id": "550e8400-e29b-41d4-a716-446655440000",
  "tokens_used": 150,
  "audio_data": "SUQzBAAAAAAAI1RTU0UAAAA...",
  "audio_url": "/audio/msg_550e8400.mp3"
}
```

**Enhanced WSResponse Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `audio_data` | string | Base64-encoded audio (if `enable_tts: true`) |
| `audio_url` | string | URL to audio file (alternative to audio_data) |

---

#### Testing Methods

##### 1. Testing with curl (REST API)

**Basic TTS Request:**
```bash
curl -X POST http://localhost:3001/api/audio/tts \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, this is a test message",
    "voice": "nova",
    "model": "tts-1"
  }' \
  --output test_audio.mp3
```

**With Different Voice:**
```bash
curl -X POST http://localhost:3001/api/audio/tts \
  -H "Content-Type: application/json" \
  -d '{
    "text": "สวัสดีครับ ทดสอบเสียงภาษาไทย",
    "voice": "onyx",
    "speed": 0.9
  }' \
  --output thai_voice.mp3
```

**High Quality TTS:**
```bash
curl -X POST http://localhost:3001/api/audio/tts \
  -H "Content-Type: application/json" \
  -d '{
    "text": "High quality audio test",
    "model": "tts-1-hd",
    "voice": "alloy",
    "response_format": "flac"
  }' \
  --output hq_audio.flac
```

##### 2. Testing with JavaScript/Fetch (REST API)

```javascript
async function testTTS() {
  const response = await fetch('http://localhost:3001/api/audio/tts', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      text: 'Testing text-to-speech API',
      voice: 'nova',
      model: 'tts-1'
    })
  })

  // Option 1: Get binary audio data
  const audioBlob = await response.blob()
  const audioUrl = URL.createObjectURL(audioBlob)

  const audio = new Audio(audioUrl)
  audio.play()

  // Option 2: Get JSON response with base64
  const data = await response.json()
  const audioData = atob(data.audio_data) // decode base64
  const bytes = new Uint8Array(audioData.length)
  for (let i = 0; i < audioData.length; i++) {
    bytes[i] = audioData.charCodeAt(i)
  }
  const blob = new Blob([bytes], { type: 'audio/mpeg' })
  const url = URL.createObjectURL(blob)

  const audioElement = new Audio(url)
  audioElement.play()
}
```

##### 3. Testing with WebSocket (Integrated Chat)

```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:3001/api/chat/stream')

ws.onopen = () => {
  // Send message with TTS enabled
  ws.send(JSON.stringify({
    type: 'message',
    content: 'Tell me about artificial intelligence',
    persona_id: 1,
    enable_tts: true,
    tts_voice: 'nova'
  }))
}

ws.onmessage = (event) => {
  const data = JSON.parse(event.data)

  if (data.type === 'chunk' && data.done) {
    // Final message with audio
    console.log('Message ID:', data.message_id)
    console.log('Tokens used:', data.tokens_used)

    if (data.audio_data) {
      // Play audio
      const audioData = atob(data.audio_data)
      const bytes = new Uint8Array(audioData.length)
      for (let i = 0; i < audioData.length; i++) {
        bytes[i] = audioData.charCodeAt(i)
      }
      const blob = new Blob([bytes], { type: 'audio/mpeg' })
      const audioUrl = URL.createObjectURL(blob)

      const audio = new Audio(audioUrl)
      audio.play()
    }
  } else if (data.type === 'chunk') {
    // Text streaming
    console.log('Chunk:', data.content)
  }
}
```

##### 4. Testing with Postman

**Setup:**
1. Create new POST request to `http://localhost:3001/api/audio/tts`
2. Set Headers: `Content-Type: application/json`
3. Set Body (raw JSON):
```json
{
  "text": "Testing with Postman",
  "voice": "nova",
  "model": "tts-1"
}
```
4. Click "Send"
5. In response, click "Save to file" to download audio
6. Play the downloaded file

**Testing Different Voices:**
- Create collection with 6 requests (one for each voice)
- Change `voice` field: alloy, echo, fable, onyx, nova, shimmer
- Compare audio quality and characteristics

##### 5. Testing with Vue Frontend Component

```vue
<template>
  <div>
    <button @click="testTTS">Test TTS</button>
    <audio ref="audioPlayer" controls></audio>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const audioPlayer = ref(null)

async function testTTS() {
  try {
    const response = await fetch('http://localhost:3001/api/audio/tts', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        text: 'Vue component test',
        voice: 'nova'
      })
    })

    const audioBlob = await response.blob()
    const audioUrl = URL.createObjectURL(audioBlob)

    if (audioPlayer.value) {
      audioPlayer.value.src = audioUrl
      audioPlayer.value.play()
    }
  } catch (error) {
    console.error('TTS test failed:', error)
  }
}
</script>
```

---

#### Error Responses

**400 Bad Request** - Missing text
```json
{
  "error": "text is required"
}
```

**400 Bad Request** - Text too long
```json
{
  "error": "text exceeds maximum length of 4096 characters"
}
```

**400 Bad Request** - Invalid voice
```json
{
  "error": "invalid voice. Allowed: alloy, echo, fable, onyx, nova, shimmer"
}
```

**400 Bad Request** - Invalid model
```json
{
  "error": "invalid model. Allowed: tts-1, tts-1-hd"
}
```

**400 Bad Request** - Invalid speed
```json
{
  "error": "speed must be between 0.25 and 4.0"
}
```

**413 Payload Too Large** - Text too long
```json
{
  "error": "text is too long (max 4096 characters)"
}
```

**500 Internal Server Error** - TTS API failure
```json
{
  "error": "failed to generate speech: API error"
}
```

**503 Service Unavailable** - OpenAI API unavailable
```json
{
  "error": "TTS service temporarily unavailable"
}
```

---

#### Cost Estimation

**Pricing (as of 2024):**
- `tts-1`: $0.015 per 1,000 characters
- `tts-1-hd`: $0.030 per 1,000 characters

**Examples:**

| Text Length | Model | Cost |
|-------------|-------|------|
| 100 chars | tts-1 | $0.0015 |
| 500 chars | tts-1 | $0.0075 |
| 1,000 chars | tts-1 | $0.015 |
| 500 chars | tts-1-hd | $0.015 |
| 1,000 chars | tts-1-hd | $0.030 |

**Average AI Response:** ~500 characters = **$0.0075 per response** (tts-1)

---

#### Rate Limiting

**Recommended Limits:**
- 10 requests per minute per user
- 1,000 requests per day per user
- Maximum 4,096 characters per request
- Maximum 100 concurrent TTS generations

**Headers for Rate Limiting:**
```
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 7
X-RateLimit-Reset: 1698765432
```

---

#### Best Practices

1. **Use `tts-1` for most cases** - Good quality, lower cost
2. **Cache common responses** - Save audio files for repeated text
3. **Let users opt-in** - Don't enable TTS by default
4. **Use `opus` format for streaming** - Better for real-time applications
5. **Implement retry logic** - Handle temporary API failures
6. **Monitor costs** - Track character usage per user
7. **Validate text length** - Prevent abuse with length limits
8. **Use appropriate voice** - Match voice to content type and audience

---

#### Integration Example (Full Flow)

**Frontend: Enable TTS in Chat**
```javascript
// In chat store
const enableTTS = ref(false)
const selectedVoice = ref('nova')

// Send message with TTS
function sendMessage(content) {
  wsConnection.send(JSON.stringify({
    type: 'message',
    content: content,
    persona_id: currentPersonaId.value,
    enable_tts: enableTTS.value,
    tts_voice: selectedVoice.value
  }))
}
```

**Backend: Process and Generate Audio**
```go
// In WebSocket controller
if msg.EnableTTS {
    audioData, err := ctrl.ttsService.TextToSpeech(ctx, fullContent)
    if err != nil {
        log.Printf("TTS failed: %v", err)
    } else {
        audioBase64 := base64.StdEncoding.EncodeToString(audioData)
        response.AudioData = audioBase64
    }
}
```

**Frontend: Play Audio Response**
```javascript
// In WebSocket handler
if (data.done && data.audio_data) {
  const audioBlob = base64ToBlob(data.audio_data, 'audio/mpeg')
  const audioUrl = URL.createObjectURL(audioBlob)
  const audio = new Audio(audioUrl)
  audio.play()
}
```

**500 Internal Server Error** - Whisper API error
```json
{
  "error": "failed to transcribe audio"
}
```

---

## 🔌 WebSocket API

### Connect to WebSocket

`WS /api/chat/stream`

เชื่อมต่อ WebSocket สำหรับรับคำตอบแบบ streaming (real-time typing effect)

#### Connection

```javascript
const ws = new WebSocket('ws://localhost:3000/api/chat/stream')

ws.onopen = () => {
  console.log('WebSocket connected')
}

ws.onmessage = (event) => {
  const data = JSON.parse(event.data)
  console.log('Received:', data)
}

ws.onerror = (error) => {
  console.error('WebSocket error:', error)
}

ws.onclose = () => {
  console.log('WebSocket disconnected')
}
```

---

### Send Message

ส่งข้อความผ่าน WebSocket

#### Request Message Format

```json
{
  "type": "message",
  "content": "เล่านิทานให้ฟังหน่อย",
  "persona_id": 1
}
```

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | ✅ Yes | Message type (always "message") |
| `content` | string | ✅ Yes | ข้อความจากผู้ใช้ |
| `persona_id` | integer | ❌ No | Persona ID (default: 1) |

#### Response Message Format

**Streaming Chunks** (multiple messages):

```json
{
  "type": "chunk",
  "content": "กาล",
  "done": false
}
```

```json
{
  "type": "chunk",
  "content": "คราว",
  "done": false
}
```

```json
{
  "type": "chunk",
  "content": "หนึ่ง",
  "done": false
}
```

**Final Message**:

```json
{
  "type": "chunk",
  "content": "",
  "done": true,
  "message_id": "550e8400-e29b-41d4-a716-446655440000",
  "tokens_used": 156
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Message type ("chunk") |
| `content` | string | เนื้อหาของ chunk (คำหรือวลี) |
| `done` | boolean | `true` = streaming เสร็จสิ้น |
| `message_id` | UUID | ID ของข้อความ (มีเฉพาะตอนจบ) |
| `tokens_used` | integer | Tokens ที่ใช้ (มีเฉพาะตอนจบ) |

---

### Example: Full WebSocket Flow

```javascript
const ws = new WebSocket('ws://localhost:3000/api/chat/stream')

ws.onopen = () => {
  // ส่งข้อความ
  ws.send(JSON.stringify({
    type: 'message',
    content: 'เล่าเรื่องสั้นให้ฟังหน่อย',
    persona_id: 1
  }))
}

let fullResponse = ''

ws.onmessage = (event) => {
  const data = JSON.parse(event.data)

  if (data.type === 'chunk') {
    if (data.done) {
      // Streaming เสร็จสิ้น
      console.log('Full response:', fullResponse)
      console.log('Message ID:', data.message_id)
      console.log('Tokens used:', data.tokens_used)
    } else {
      // รับ chunk
      fullResponse += data.content
      console.log('Chunk:', data.content)
      // แสดงใน UI แบบ real-time
      updateChatUI(data.content)
    }
  }
}

ws.onerror = (error) => {
  console.error('Error:', error)
}
```

---

### Error Messages

WebSocket อาจส่ง error message:

```json
{
  "type": "error",
  "error": {
    "code": "OPENAI_ERROR",
    "message": "Failed to connect to OpenAI API"
  }
}
```

---

## ⏱️ Rate Limiting

### Current Limits

**สถานะปัจจุบัน**: ❌ **No Rate Limiting**

**แนะนำสำหรับ Production**:

| Endpoint | Limit | Window |
|----------|-------|--------|
| `POST /api/chat` | 20 requests | per minute |
| `POST /api/audio/transcribe` | 10 requests | per minute |
| `GET /api/personas` | 60 requests | per minute |
| WebSocket connections | 5 connections | per user |

### Rate Limit Headers (Future)

```
X-RateLimit-Limit: 20
X-RateLimit-Remaining: 15
X-RateLimit-Reset: 1635724800
```

### 429 Response

```json
{
  "error": "Rate limit exceeded. Try again in 45 seconds",
  "retry_after": 45
}
```

---

## 📝 Examples

### Example 1: Simple Chat Flow

```javascript
// 1. เลือก persona
const personasResponse = await fetch('http://localhost:3000/api/personas')
const personas = await personasResponse.json()

// 2. ส่งข้อความ
const chatResponse = await fetch('http://localhost:3000/api/chat', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    message: 'สวัสดีครับ',
    persona_id: personas.personas[0].id
  })
})

const chatResult = await chatResponse.json()
console.log('AI:', chatResult.reply)
```

---

### Example 2: Voice Message Flow

```javascript
// 1. บันทึกเสียง (MediaRecorder API)
const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
const mediaRecorder = new MediaRecorder(stream)
const audioChunks = []

mediaRecorder.ondataavailable = (event) => {
  audioChunks.push(event.data)
}

mediaRecorder.onstop = async () => {
  // 2. สร้าง audio blob
  const audioBlob = new Blob(audioChunks, { type: 'audio/webm' })

  // 3. Upload และ transcribe
  const formData = new FormData()
  formData.append('audio', audioBlob, 'recording.webm')

  const response = await fetch('http://localhost:3000/api/audio/transcribe', {
    method: 'POST',
    body: formData
  })

  const result = await response.json()
  console.log('Transcribed text:', result.text)

  // 4. ส่งเป็นข้อความ
  const chatResponse = await fetch('http://localhost:3000/api/chat', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      message: result.text,
      persona_id: 1
    })
  })

  const chatResult = await chatResponse.json()
  console.log('AI:', chatResult.reply)
}

// เริ่มบันทึก
mediaRecorder.start()

// หยุดบันทึก (หลังจาก 5 วินาที)
setTimeout(() => {
  mediaRecorder.stop()
  stream.getTracks().forEach(track => track.stop())
}, 5000)
```

---

### Example 3: Streaming Chat with WebSocket

```javascript
class ChatClient {
  constructor() {
    this.ws = null
    this.currentResponse = ''
  }

  connect() {
    this.ws = new WebSocket('ws://localhost:3000/api/chat/stream')

    this.ws.onopen = () => {
      console.log('Connected to chat server')
    }

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      this.handleMessage(data)
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    this.ws.onclose = () => {
      console.log('Disconnected')
      // Auto-reconnect
      setTimeout(() => this.connect(), 3000)
    }
  }

  sendMessage(message, personaId = 1) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.currentResponse = ''
      this.ws.send(JSON.stringify({
        type: 'message',
        content: message,
        persona_id: personaId
      }))
    }
  }

  handleMessage(data) {
    if (data.type === 'chunk') {
      if (data.done) {
        console.log('Complete response:', this.currentResponse)
        console.log('Tokens used:', data.tokens_used)
      } else {
        this.currentResponse += data.content
        // Update UI
        this.updateChatUI(data.content)
      }
    } else if (data.type === 'error') {
      console.error('Error:', data.error.message)
    }
  }

  updateChatUI(chunk) {
    // Append chunk to chat interface
    const chatBox = document.getElementById('chat-messages')
    const lastMessage = chatBox.lastElementChild
    if (lastMessage && lastMessage.classList.contains('ai-message')) {
      lastMessage.textContent += chunk
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
    }
  }
}

// Usage
const client = new ChatClient()
client.connect()

// ส่งข้อความ
client.sendMessage('เล่าเรื่องสั้นให้ฟังหน่อย', 1)
```

---

### Example 4: Error Handling

```javascript
async function sendChatMessage(message, personaId) {
  try {
    const response = await fetch('http://localhost:3000/api/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        message,
        persona_id: personaId
      })
    })

    // ตรวจสอบ HTTP status
    if (!response.ok) {
      const error = await response.json()

      switch (response.status) {
        case 400:
          throw new Error(`Invalid input: ${error.error}`)
        case 404:
          throw new Error(`Persona not found: ${error.error}`)
        case 429:
          throw new Error(`Rate limit exceeded. Try again later.`)
        case 500:
          throw new Error(`Server error: ${error.error}`)
        default:
          throw new Error(`Unexpected error: ${error.error}`)
      }
    }

    const data = await response.json()
    return data

  } catch (error) {
    console.error('Failed to send message:', error.message)

    // แสดง error ให้ user
    showErrorNotification(error.message)

    // Retry logic (optional)
    if (shouldRetry(error)) {
      console.log('Retrying in 3 seconds...')
      await delay(3000)
      return sendChatMessage(message, personaId)
    }

    throw error
  }
}

function shouldRetry(error) {
  // Retry on network errors or 500 errors
  return error.message.includes('network') ||
         error.message.includes('Server error')
}

function delay(ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}
```

---

## 🧪 Testing the API

### Using curl

```bash
# Health check
curl http://localhost:3000/api/health

# Get all personas
curl http://localhost:3000/api/personas

# Get persona by ID
curl http://localhost:3000/api/personas/1

# Send chat message
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"สวัสดี","persona_id":1}'

# Get chat history
curl "http://localhost:3000/api/chat/history?limit=10"

# Upload audio file
curl -X POST http://localhost:3000/api/audio/transcribe \
  -F "audio=@test.mp3"
```

### Using Postman

1. **Import Collection**: Create new collection "ChatBot API"
2. **Add Requests**:
   - GET Health Check
   - GET Personas
   - POST Send Message
   - POST Audio Transcribe
3. **Set Environment Variables**:
   - `base_url`: `http://localhost:3000/api`

### Using JavaScript (Node.js)

```javascript
const axios = require('axios')

const API_BASE_URL = 'http://localhost:3000/api'

// Test health
async function testHealth() {
  const response = await axios.get(`${API_BASE_URL}/health`)
  console.log('Health:', response.data)
}

// Test chat
async function testChat() {
  const response = await axios.post(`${API_BASE_URL}/chat`, {
    message: 'สวัสดีครับ',
    persona_id: 1
  })
  console.log('Chat response:', response.data.reply)
}

testHealth()
testChat()
```

---

## 📚 Additional Resources

### Related Documentation
- [BMAD.md](BMAD.md) - Architecture และ Design
- [WORKFLOW.md](WORKFLOW.md) - Development Workflow
- [README.md](README.md) - Project Overview

### External APIs
- [OpenAI Chat Completions API](https://platform.openai.com/docs/api-reference/chat)
- [OpenAI Whisper API](https://platform.openai.com/docs/api-reference/audio)

### Tools
- [Postman](https://www.postman.com/) - API Testing
- [WebSocket Test Client](https://www.websocket.org/echo.html) - WebSocket Testing
- [curl](https://curl.se/) - Command-line HTTP client

---

## 🔄 Changelog

### Version 1.1.0 (2025-10-31) - **File Analysis Update**
- ⭐ **NEW**: File Analysis Endpoints documentation
  - 📁 Upload and Analyze File endpoint (`POST /api/file/analyze`)
  - 📊 Detailed implementation steps (6 phases)
  - 🔧 Code examples for PDF, Word, Excel, Image parsing
  - 💡 OpenAI Vision API integration guide
  - 📈 Cost estimation and dependencies
- ✅ File history endpoint (`GET /api/file/history`)
- ✅ Re-analyze file endpoint (`POST /api/file/:file_id/reanalyze`)
- 📝 Implementation phases (Week 1-4 roadmap)

### Version 1.0.0 (2025-10-28)
- ✅ Initial API documentation
- ✅ Health check endpoint
- ✅ Chat endpoints (sync)
- ✅ Persona endpoints
- ✅ Audio transcription endpoint
- ✅ Text-to-Speech endpoint
- ✅ WebSocket streaming
- ✅ Chat history endpoint

---

## 📞 Support

หากพบปัญหาหรือมีคำถาม:

1. ตรวจสอบ [Error Handling](#error-handling) section
2. ดู [Examples](#examples) สำหรับ use cases ทั่วไป
3. เปิด issue ใน GitHub repository
4. อ่าน [BMAD.md](BMAD.md) สำหรับ architecture details

---

**📅 Last Updated**: 2025-10-31
**📝 Version**: 1.1.0
**🔗 Base URL**: http://localhost:3000/api
**✨ What's New**: File Analysis API - รองรับการวิเคราะห์ไฟล์ PDF, Word, Excel, Images และอื่นๆ
**📧 Contact**: [Your Email]
