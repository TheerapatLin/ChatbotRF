# Whisper.cpp Speech-to-Text System

ระบบแปลงเสียงเป็นข้อความโดยใช้ whisper.cpp รองรับภาษาไทยและภาษาอังกฤษ

## โครงสร้างโฟลเดอร์

```
whisper/
├── binary/              # whisper.cpp binaries
│   ├── linux/          # Binary สำหรับ Linux
│   ├── windows/        # Binary สำหรับ Windows (main.exe)
│   └── macos/          # Binary สำหรับ macOS
├── models/             # GGML models
│   └── ggml-small.bin  # Model small (466 MB)
├── temp/               # ไฟล์เสียงชั่วคราว
├── whisper-source/     # Source code จาก GitHub
└── README.md           # ไฟล์นี้
```

## สถานะการติดตั้ง

### ✅ ติดตั้งเรียบร้อยแล้ว

- [x] โครงสร้างโฟลเดอร์
- [x] Clone whisper.cpp source (whisper-source/)
- [x] ดาวน์โหลด binary สำหรับ Windows (main.exe)
- [x] ดาวน์โหลด GGML model small (466 MB)

### Binary Information

- **Windows**: `binary/windows/main.exe` (111 KB)
  - Downloaded from: whisper.cpp releases v1.5.4
  - Pre-compiled binary for x64

### Model Information

- **ggml-small.bin**: 466 MB
  - รองรับ 99+ ภาษา รวมภาษาไทยและอังกฤษ
  - ความแม่นยำดี และประมวลผลเร็ว
  - Downloaded from: Hugging Face (ggerganov/whisper.cpp)

## การใช้งาน

### ทดสอบ Binary

```bash
# Windows
cmd /c "backend\whisper\binary\windows\main.exe --help"

# หรือใช้ผ่าน Go service (แนะนำ)
```

### ทดสอบการแปลงเสียง

```bash
# แปลงไฟล์เสียงภาษาไทย
cmd /c "backend\whisper\binary\windows\main.exe -m backend\whisper\models\ggml-small.bin -f audio.wav -l th"

# แปลงไฟล์เสียงภาษาอังกฤษ
cmd /c "backend\whisper\binary\windows\main.exe -m backend\whisper\models\ggml-small.bin -f audio.wav -l en"

# ตรวจจับภาษาอัตโนมัติ
cmd /c "backend\whisper\binary\windows\main.exe -m backend\whisper\models\ggml-small.bin -f audio.wav"
```

## ขั้นตอนถัดไป

1. สร้าง Unit Tests (backend/test/sst-whisper/)
2. สร้าง WhisperCppService (backend/services/)
3. สร้าง WhisperCppController (backend/controllers/)
4. ตั้งค่า Environment Variables
5. ลงทะเบียน Routes

## หมายเหตุ

- Binary สำหรับ Linux และ macOS ยังไม่ได้ติดตั้ง (ต้องคอมไพล์บนระบบนั้นๆ)
- ไฟล์ในโฟลเดอร์ temp/ จะถูกลบอัตโนมัติหลังการประมวลผล
- สำหรับการใช้งานจริง ควรใช้ผ่าน Go service แทนการเรียก binary โดยตรง

## ข้อมูลเพิ่มเติม

- GitHub: https://github.com/ggerganov/whisper.cpp
- Models: https://huggingface.co/ggerganov/whisper.cpp
- Documentation: documents/WHISPER_START.md
