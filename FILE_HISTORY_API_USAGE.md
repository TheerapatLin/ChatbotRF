# 📂 File History API - Usage Guide

## API Endpoint

### `GET /api/file/history`

ดึงประวัติการวิเคราะห์ไฟล์ที่ผ่านมา พร้อมรองรับการกรองและ pagination

---

## Query Parameters

| Parameter | Type | Default | Required | Description |
|-----------|------|---------|----------|-------------|
| `limit` | integer | 20 | ❌ | จำนวนรายการต่อหน้า (สูงสุด 100) |
| `offset` | integer | 0 | ❌ | เริ่มต้นจากรายการที่ (สำหรับ pagination) |
| `file_type` | string | "all" | ❌ | กรองตามประเภทไฟล์ (MIME type) |

---

## Response Format

### Success Response (200 OK)

```json
{
  "files": [
    {
      "file_id": "550e8400-e29b-41d4-a716-446655440000",
      "filename": "report.pdf",
      "file_type": "application/pdf",
      "file_size": 156234,
      "analysis_type": "summary",
      "language": "th",
      "tokens_used": 450,
      "created_at": "2025-10-31T14:30:00Z"
    },
    {
      "file_id": "660e9500-e29b-41d4-a716-446655440001",
      "filename": "data.xlsx",
      "file_type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
      "file_size": 89456,
      "analysis_type": "detail",
      "language": "en",
      "tokens_used": 780,
      "created_at": "2025-10-31T13:15:00Z"
    }
  ],
  "total": 15,
  "limit": 20,
  "offset": 0
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `files` | array | รายการไฟล์ที่วิเคราะห์แล้ว |
| `files[].file_id` | string (UUID) | ID ของไฟล์จาก database |
| `files[].filename` | string | ชื่อไฟล์ต้นฉบับ |
| `files[].file_type` | string | MIME type ของไฟล์ |
| `files[].file_size` | integer | ขนาดไฟล์ (bytes) |
| `files[].analysis_type` | string | ประเภทการวิเคราะห์: summary, detail, qa, extract |
| `files[].language` | string | ภาษาที่ใช้: th, en |
| `files[].tokens_used` | integer | จำนวน tokens ที่ใช้จาก OpenAI |
| `files[].created_at` | string (ISO 8601) | วันเวลาที่วิเคราะห์ |
| `total` | integer | จำนวนรายการทั้งหมดใน database |
| `limit` | integer | จำนวนรายการต่อหน้า |
| `offset` | integer | ตำแหน่งเริ่มต้น |

### Error Response (500 Internal Server Error)

```json
{
  "error": "failed to fetch file history"
}
```

---

## Usage Examples

### 1. ดึงประวัติทั้งหมด (default: 20 รายการแรก)

```bash
curl -X GET "http://localhost:3000/api/file/history"
```

### 2. ดึงแบบ pagination (10 รายการ, เริ่มจากรายการที่ 20)

```bash
curl -X GET "http://localhost:3000/api/file/history?limit=10&offset=20"
```

### 3. กรองเฉพาะไฟล์ PDF

```bash
curl -X GET "http://localhost:3000/api/file/history?file_type=application/pdf"
```

### 4. กรองเฉพาะไฟล์ Excel

```bash
curl -X GET "http://localhost:3000/api/file/history?file_type=application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
```

### 5. กรองเฉพาะไฟล์รูปภาพ (JPEG)

```bash
curl -X GET "http://localhost:3000/api/file/history?file_type=image/jpeg"
```

### 6. ดึงรายการจำนวนมาก (max 100)

```bash
curl -X GET "http://localhost:3000/api/file/history?limit=100"
```

---

## Common File Types

### Document Files
- **PDF**: `application/pdf`
- **Word**: `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
- **Excel**: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- **PowerPoint**: `application/vnd.openxmlformats-officedocument.presentationml.presentation`

### Text Files
- **Plain Text**: `text/plain`
- **Markdown**: `text/markdown`
- **CSV**: `text/csv`
- **JSON**: `application/json`
- **XML**: `application/xml`

### Image Files
- **JPEG**: `image/jpeg`
- **PNG**: `image/png`
- **GIF**: `image/gif`
- **WebP**: `image/webp`

### Code Files
- **JavaScript**: `application/javascript`
- **Python**: `text/x-python`
- **Go**: `text/x-go`
- **Java**: `text/x-java`

---

## Pagination Example

สมมติมีข้อมูล 45 รายการใน database:

### หน้าที่ 1 (รายการ 1-20)
```bash
GET /api/file/history?limit=20&offset=0
```
Response: `"total": 45, "limit": 20, "offset": 0`

### หน้าที่ 2 (รายการ 21-40)
```bash
GET /api/file/history?limit=20&offset=20
```
Response: `"total": 45, "limit": 20, "offset": 20`

### หน้าที่ 3 (รายการ 41-45)
```bash
GET /api/file/history?limit=20&offset=40
```
Response: `"total": 45, "limit": 20, "offset": 40` (จะมีแค่ 5 รายการ)

---

## Implementation Details

### Database Query
- ใช้ `FileAnalysisRepository.GetAll()` สำหรับดึงข้อมูลทั้งหมด
- ใช้ `FileAnalysisRepository.GetByFileType()` สำหรับกรองตามประเภทไฟล์
- เรียงลำดับจากใหม่ไปเก่า (`ORDER BY created_at DESC`)
- รองรับ soft delete (ไม่แสดงรายการที่ถูกลบ)

### Validation Rules
- `limit`: ต้องอยู่ระหว่าง 1-100 (default: 20)
- `offset`: ต้องไม่ติดลบ (default: 0)
- `file_type`: รับค่าใดก็ได้ หรือ "all" สำหรับทั้งหมด

### Performance
- Query มี index บน `created_at` และ `file_type`
- Pagination ช่วยลดปริมาณข้อมูลที่ส่งกลับ
- ไม่ return field `analysis`, `summary`, `key_points` ใน history เพื่อลด payload size

---

## Frontend Integration Example

### JavaScript (Fetch API)

```javascript
async function getFileHistory(page = 1, limit = 20, fileType = 'all') {
  const offset = (page - 1) * limit;

  let url = `http://localhost:3000/api/file/history?limit=${limit}&offset=${offset}`;
  if (fileType !== 'all') {
    url += `&file_type=${encodeURIComponent(fileType)}`;
  }

  try {
    const response = await fetch(url);
    const data = await response.json();

    console.log(`Showing ${data.files.length} of ${data.total} files`);
    console.log('Files:', data.files);

    return data;
  } catch (error) {
    console.error('Failed to fetch file history:', error);
  }
}

// Usage
getFileHistory(1, 20, 'application/pdf');  // Page 1, PDF files only
```

### Vue.js Example

```vue
<template>
  <div>
    <h2>File History</h2>

    <!-- Filter -->
    <select v-model="selectedFileType" @change="loadHistory">
      <option value="all">All Files</option>
      <option value="application/pdf">PDF Files</option>
      <option value="image/jpeg">Images (JPEG)</option>
      <option value="text/plain">Text Files</option>
    </select>

    <!-- File List -->
    <ul>
      <li v-for="file in files" :key="file.file_id">
        <strong>{{ file.filename }}</strong>
        <span>{{ file.file_type }}</span>
        <span>{{ new Date(file.created_at).toLocaleString() }}</span>
        <span>Tokens: {{ file.tokens_used }}</span>
      </li>
    </ul>

    <!-- Pagination -->
    <div>
      <button @click="prevPage" :disabled="currentPage === 1">Previous</button>
      <span>Page {{ currentPage }} of {{ totalPages }}</span>
      <button @click="nextPage" :disabled="currentPage === totalPages">Next</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      files: [],
      total: 0,
      currentPage: 1,
      limit: 20,
      selectedFileType: 'all'
    };
  },
  computed: {
    totalPages() {
      return Math.ceil(this.total / this.limit);
    }
  },
  methods: {
    async loadHistory() {
      const offset = (this.currentPage - 1) * this.limit;
      let url = `http://localhost:3000/api/file/history?limit=${this.limit}&offset=${offset}`;

      if (this.selectedFileType !== 'all') {
        url += `&file_type=${encodeURIComponent(this.selectedFileType)}`;
      }

      try {
        const response = await fetch(url);
        const data = await response.json();

        this.files = data.files;
        this.total = data.total;
      } catch (error) {
        console.error('Failed to load file history:', error);
      }
    },
    prevPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
        this.loadHistory();
      }
    },
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++;
        this.loadHistory();
      }
    }
  },
  mounted() {
    this.loadHistory();
  }
};
</script>
```

---

## Testing

### 1. ทดสอบดึงข้อมูลทั้งหมด

```bash
curl -X GET "http://localhost:3000/api/file/history" | jq
```

### 2. ทดสอบ pagination

```bash
# Page 1
curl -X GET "http://localhost:3000/api/file/history?limit=5&offset=0" | jq

# Page 2
curl -X GET "http://localhost:3000/api/file/history?limit=5&offset=5" | jq
```

### 3. ทดสอบ filter

```bash
# PDF files only
curl -X GET "http://localhost:3000/api/file/history?file_type=application/pdf" | jq

# Image files only
curl -X GET "http://localhost:3000/api/file/history?file_type=image/jpeg" | jq
```

### 4. ทดสอบ edge cases

```bash
# Limit > 100 (should return max 100)
curl -X GET "http://localhost:3000/api/file/history?limit=200" | jq

# Negative offset (should default to 0)
curl -X GET "http://localhost:3000/api/file/history?offset=-10" | jq

# Invalid limit (should default to 20)
curl -X GET "http://localhost:3000/api/file/history?limit=0" | jq
```

---

## Related APIs

- **[POST /api/file/analyze](./FILE_ANALYSIS_API_USAGE.md)** - วิเคราะห์ไฟล์ใหม่
- **POST /api/file/:file_id/reanalyze** - วิเคราะห์ไฟล์เดิมอีกครั้ง (ยังไม่ implement)

---

## Status

✅ **Implemented** - Ready to use

**Created**: 2025-10-31
**Last Updated**: 2025-10-31
**Version**: 1.0.0