# 🔧 การแก้ปัญหา Error: exit status 0xc0000135

## 📋 สรุปปัญหา

จากการรัน Unit Tests บน Windows พบ error ดังนี้:

```bash
=== RUN   TestWhisperBinaryExists
    ✓ พบ whisper.cpp binary ที่ ..\..\whisper\binary\windows\main.exe (OS: windows)
--- PASS: TestWhisperBinaryExists (0.00s)

=== RUN   TestWhisperModelExists
    ✓ พบโมเดล whisper ที่ ..\..\whisper\models\ggml-small.bin (ขนาด: 465 MB)
--- PASS: TestWhisperModelExists (0.00s)

=== RUN   TestWhisperVersion
    ไม่สามารถรัน whisper.cpp: exit status 0xc0000135 (OS: windows)
--- FAIL: TestWhisperVersion (0.01s)

=== RUN   TestWhisperTranscribeThaiAudio
    ไม่สามารถแปลงเสียงภาษาไทย: exit status 0xc0000135
--- FAIL: TestWhisperTranscribeThaiAudio (0.01s)

=== RUN   TestWhisperTranscribeEnglishAudio
    ไม่สามารถแปลงเสียงภาษาอังกฤษ: exit status 0xc0000135
--- FAIL: TestWhisperTranscribeEnglishAudio (0.01s)
```

## 🔍 การวิเคราะห์ปัญหา

### Error Code: 0xc0000135 คืออะไร?

**Error Code `0xc0000135`** = `STATUS_DLL_NOT_FOUND`

หมายความว่า:
- Windows ไม่สามารถหาไฟล์ DLL ที่โปรแกรม `main.exe` ต้องการ
- Binary file `main.exe` ถูกคอมไพล์บนระบบอื่น และต้องการ runtime libraries
- ส่วนใหญ่มักขาด **Microsoft Visual C++ Redistributable**

### ทำไม Binary และ Model Tests ผ่าน?

✅ **TestWhisperBinaryExists** - PASS
- ตรวจสอบแค่ว่าไฟล์ `main.exe` มีอยู่หรือไม่
- ไม่ได้พยายามรันไฟล์

✅ **TestWhisperModelExists** - PASS
- ตรวจสอบแค่ว่าไฟล์ `ggml-small.bin` มีอยู่และมีขนาดถูกต้อง
- ไม่ได้พยายามใช้งาน model

❌ **TestWhisperVersion** - FAIL
- พยายาม**รัน** `main.exe --help`
- ตรงนี้ Windows ต้องโหลด DLL dependencies
- พบว่าขาด DLL → Error 0xc0000135

❌ **Transcription Tests** - FAIL
- พยายาม**รัน** `main.exe` กับ audio files
- เจอปัญหาเดียวกัน

### Binary มาจากไหน?

จากการตรวจสอบ:
```bash
backend/whisper/binary/windows/main.exe (111 KB)
```

Binary นี้:
- ดาวน์โหลดจาก **whisper.cpp releases v1.5.4** (pre-compiled)
- คอมไพล์บน Linux หรือระบบอื่นด้วย MinGW
- ต้องการ **Visual C++ Runtime Libraries** ที่ไม่มีใน Windows ของคุณ

### DLL ที่อาจขาดหายไป

Binary อาจต้องการ DLL เหล่านี้:
- `vcruntime140.dll` - Visual C++ Runtime
- `msvcp140.dll` - C++ Standard Library
- `ucrtbase.dll` - Universal C Runtime
- `libgcc_s_seh-1.dll` - GCC Runtime (ถ้าคอมไพล์ด้วย MinGW)
- `libstdc++-6.dll` - C++ Standard Library (MinGW)
- `libwinpthread-1.dll` - POSIX Threads (MinGW)

## ✅ วิธีแก้ปัญหา (3 วิธี)

---

## วิธีที่ 1: ติดตั้ง Visual C++ Redistributable (แนะนำสำหรับ Windows)

### ขั้นตอนที่ 1.1: ดาวน์โหลดและติดตั้ง

```bash
# ดาวน์โหลด Visual C++ Redistributable (x64)
# URL: https://aka.ms/vs/17/release/vc_redist.x64.exe
```

หรือใช้ PowerShell:
```powershell
# Download
Invoke-WebRequest -Uri "https://aka.ms/vs/17/release/vc_redist.x64.exe" -OutFile "$env:TEMP\vc_redist.x64.exe"

# Install
Start-Process -FilePath "$env:TEMP\vc_redist.x64.exe" -ArgumentList "/quiet", "/norestart" -Wait

# Verify
Write-Host "Installation completed. Please restart your terminal."
```

### ขั้นตอนที่ 1.2: Restart Terminal

```bash
# ปิด terminal ปัจจุบัน
exit

# เปิดใหม่แล้วรัน tests อีกครั้ง
cd C:\Users\boatr\MyBoat\RealFactory\ChatBotProject\backend\test\sst-whisper
go test -v
```

### ผลลัพธ์ที่คาดหวัง

```bash
=== RUN   TestWhisperBinaryExists
    ✓ พบ whisper.cpp binary ที่ ..\..\whisper\binary\windows\main.exe (OS: windows)
--- PASS: TestWhisperBinaryExists (0.00s)

=== RUN   TestWhisperModelExists
    ✓ พบโมเดล whisper ที่ ..\..\whisper\models\ggml-small.bin (ขนาด: 465 MB)
--- PASS: TestWhisperModelExists (0.00s)

=== RUN   TestWhisperVersion
    ✓ whisper.cpp ทำงานได้ปกติ (OS: windows)
--- PASS: TestWhisperVersion (0.15s)

=== RUN   TestWhisperTranscribeThaiAudio
    ✓ แปลงเสียงภาษาไทยสำเร็จ
--- PASS: TestWhisperTranscribeThaiAudio (4.32s)

=== RUN   TestWhisperTranscribeEnglishAudio
    ✓ แปลงเสียงภาษาอังกฤษสำเร็จ
--- PASS: TestWhisperTranscribeEnglishAudio (5.87s)

PASS
ok      chatbot/test/sst-whisper        10.352s
```

### ข้อดี
- ✅ แก้ปัญหาได้ทันที (ใช้เวลา 2-3 นาที)
- ✅ ไม่ต้องคอมไพล์ binary ใหม่
- ✅ ใช้งานได้กับ pre-compiled binary

### ข้อเสีย
- ⚠️ ต้องติดตั้ง software เพิ่ม (~25 MB)
- ⚠️ ทุกคนใน team ต้องติดตั้ง

---

## วิธีที่ 2: คอมไพล์ Binary ใหม่บน Windows (เหมาะสำหรับ Advanced Users)

### ขั้นตอนที่ 2.1: ติดตั้ง MSYS2

```bash
# ดาวน์โหลด MSYS2 Installer
# URL: https://www.msys2.org/
# File: msys2-x86_64-YYYYMMDD.exe

# หรือใช้ winget (Windows 11)
winget install MSYS2.MSYS2
```

### ขั้นตอนที่ 2.2: ติดตั้ง Build Tools

เปิด **MSYS2 MinGW 64-bit** terminal:

```bash
# Update packages
pacman -Syu

# Install build tools
pacman -S --needed base-devel mingw-w64-x86_64-toolchain
pacman -S mingw-w64-x86_64-cmake
pacman -S git
```

### ขั้นตอนที่ 2.3: คอมไพล์ whisper.cpp

```bash
# Navigate to whisper source (ใน MSYS2 terminal)
cd /c/Users/boatr/MyBoat/RealFactory/ChatBotProject/backend/whisper/whisper-source

# Clean previous build
rm -rf build

# Build with CMake
cmake -B build -G "MinGW Makefiles" -DCMAKE_BUILD_TYPE=Release
cmake --build build --config Release

# Copy binary
cp build/bin/whisper-cli.exe ../binary/windows/main.exe
```

### ขั้นตอนที่ 2.4: Copy Required DLLs

```bash
# Find required DLLs
ldd build/bin/whisper-cli.exe

# Copy DLLs ที่จำเป็น (ถ้ามี)
# ตัวอย่าง:
cp /mingw64/bin/libgcc_s_seh-1.dll ../binary/windows/
cp /mingw64/bin/libstdc++-6.dll ../binary/windows/
cp /mingw64/bin/libwinpthread-1.dll ../binary/windows/
```

### ขั้นตอนที่ 2.5: ทดสอบ

```bash
# ใน PowerShell หรือ CMD
cd C:\Users\boatr\MyBoat\RealFactory\ChatBotProject\backend\test\sst-whisper
go test -v
```

### ข้อดี
- ✅ Binary รัน native บน Windows
- ✅ ไม่ต้องพึ่ง Visual C++ Redistributable
- ✅ Performance อาจดีกว่า

### ข้อเสีย
- ⚠️ ใช้เวลานาน (10-15 นาที ครั้งแรก)
- ⚠️ ต้องติดตั้ง MSYS2 และ build tools (~1 GB)
- ⚠️ ซับซ้อนกว่า

---

## วิธีที่ 3: ใช้ WSL2 (Windows Subsystem for Linux)

### ขั้นตอนที่ 3.1: ติดตั้ง WSL2

```powershell
# ใน PowerShell (Admin)
wsl --install

# Restart เครื่อง
Restart-Computer
```

### ขั้นตอนที่ 3.2: Setup WSL2 Ubuntu

```bash
# เปิด WSL2
wsl

# Update packages
sudo apt update
sudo apt upgrade -y

# Install build tools
sudo apt install -y build-essential cmake git golang-go
```

### ขั้นตอนที่ 3.3: คอมไพล์ whisper.cpp สำหรับ Linux

```bash
# Navigate to project (WSL can access Windows files)
cd /mnt/c/Users/boatr/MyBoat/RealFactory/ChatBotProject/backend/whisper/whisper-source

# ⚠️ หมายเหตุ: whisper.cpp ใช้ CMake ไม่ใช่ Makefile โดยตรง
# ต้อง generate build files ด้วย CMake ก่อน

# วิธีที่ 1: ใช้ CMake (แนะนำ)
cmake -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build --config Release

# วิธีที่ 2: ใช้ make (จะเรียก cmake อัตโนมัติ)
make

# Copy binary for Linux
cp build/bin/whisper-cli ../binary/linux/main
chmod +x ../binary/linux/main

# Verify binary
ls -lh ../binary/linux/main
file ../binary/linux/main
```

**💡 คำอธิบาย:**
- `cmake -B build` = สร้างโฟลเดอร์ build และ generate build files
- `cmake --build build` = คอมไพล์โปรเจค
- `make` = wrapper ที่จะเรียก cmake อัตโนมัติ (ถ้ามี Makefile)
- Binary จะอยู่ที่ `build/bin/whisper-cli` (ไม่ใช่ `bin/main`)

**🔧 แก้ปัญหา Error ที่พบ:**

```bash
# ถ้าเจอ: make: *** No rule to make target 'clean'. Stop.
# แก้: ใช้ CMake แทน
rm -rf build  # ลบโฟลเดอร์ build (แทน make clean)
cmake -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build --config Release

# ถ้าเจอ: cp: cannot stat 'build/bin/whisper-cli'
# แก้: ตรวจสอบว่าคอมไพล์สำเร็จหรือไม่
ls -la build/bin/  # ดูว่ามี binary อะไรบ้าง
# Binary อาจชื่อ: whisper-cli, main, หรืออื่นๆ

# ถ้า binary ชื่ออื่น เช่น main
cp build/bin/main ../binary/linux/main

# หรือถ้ามีหลายตัว
cp build/bin/whisper-cli ../binary/linux/main  # CLI version
# หรือ
cp build/main ../binary/linux/main  # ตัวเก่า
```

### ขั้นตอนที่ 3.4: รัน Tests ใน WSL2

```bash
# Navigate to test directory
cd /mnt/c/Users/boatr/MyBoat/RealFactory/ChatBotProject/backend/test/sst-whisper

# Run tests (จะใช้ Linux binary)
go test -v
```

### ผลลัพธ์ที่คาดหวัง

```bash
=== RUN   TestWhisperBinaryExists
    ✓ พบ whisper.cpp binary ที่ ../../whisper/binary/linux/main (OS: linux)
--- PASS: TestWhisperBinaryExists (0.00s)

=== RUN   TestWhisperModelExists
    ✓ พบโมเดล whisper ที่ ../../whisper/models/ggml-small.bin (ขนาด: 465 MB)
--- PASS: TestWhisperModelExists (0.00s)

=== RUN   TestWhisperVersion
    ✓ whisper.cpp ทำงานได้ปกติ (OS: linux)
--- PASS: TestWhisperVersion (0.15s)

=== RUN   TestWhisperTranscribeThaiAudio
    ✓ แปลงเสียงภาษาไทยสำเร็จ
--- PASS: TestWhisperTranscribeThaiAudio (4.32s)

=== RUN   TestWhisperTranscribeEnglishAudio
    ✓ แปลงเสียงภาษาอังกฤษสำเร็จ
--- PASS: TestWhisperTranscribeEnglishAudio (5.87s)

PASS
ok      chatbot/test/sst-whisper        10.352s
```

### ข้อดี
- ✅ Environment เหมือน production (Linux)
- ✅ รัน Docker ได้เร็วกว่า
- ✅ เหมาะสำหรับ development ในอนาคต

### ข้อเสีย
- ⚠️ ต้อง restart เครื่องครั้งแรก
- ⚠️ ใช้ disk space เพิ่ม (~3-5 GB)
- ⚠️ ต้องเรียนรู้ Linux commands

---

## 📊 เปรียบเทียบวิธีแก้ปัญหา

| วิธี | เวลา | ความยาก | Disk Space | แนะนำสำหรับ |
|------|------|---------|------------|-------------|
| **1. Visual C++ Redistributable** | 2-3 นาที | ⭐ ง่าย | 25 MB | ทุกคน (Quick Fix) |
| **2. คอมไพล์ด้วย MSYS2** | 15-20 นาที | ⭐⭐⭐ ยาก | 1 GB | Advanced Users |
| **3. ใช้ WSL2** | 10-15 นาที | ⭐⭐ ปานกลาง | 3-5 GB | Developers |

---

## 🎯 คำแนะนำ

### สำหรับการแก้ปัญหาด่วน:
👉 **ใช้วิธีที่ 1: ติดตั้ง Visual C++ Redistributable**

### สำหรับการใช้งานระยะยาว:
👉 **ใช้วิธีที่ 3: WSL2** (เพราะ production จะรันบน Linux อยู่แล้ว)

---

## 🔍 วิธีตรวจสอบว่าแก้ปัญหาสำเร็จ

### ทดสอบ Binary

```bash
# ทดสอบรัน binary โดยตรง
cd C:\Users\boatr\MyBoat\RealFactory\ChatBotProject\backend\whisper\binary\windows
.\main.exe --help

# ควรเห็น help message แทน error
```

### รัน Unit Tests

```bash
cd C:\Users\boatr\MyBoat\RealFactory\ChatBotProject\backend\test\sst-whisper
go test -v

# ควรผ่านทั้ง 5 tests
```

### ทดสอบ Transcription จริง

```bash
# ทดสอบแปลงเสียงภาษาอังกฤษ
cd C:\Users\boatr\MyBoat\RealFactory\ChatBotProject\backend\whisper\binary\windows
.\main.exe -m ..\..\models\ggml-small.bin -f ..\..\..\..\test\sst-whisper\testdata\audio\en_audio.mp3 -l en

# ควรเห็นผลลัพธ์การแปลงเสียง
```

---

## 📚 เอกสารเพิ่มเติม

- **Error Code Reference**: https://learn.microsoft.com/en-us/windows/win32/debug/system-error-codes
- **Visual C++ Redistributable**: https://learn.microsoft.com/en-us/cpp/windows/latest-supported-vc-redist
- **MSYS2**: https://www.msys2.org/
- **WSL2**: https://learn.microsoft.com/en-us/windows/wsl/install
- **whisper.cpp**: https://github.com/ggerganov/whisper.cpp

---

## 🐛 ปัญหาที่พบระหว่างการแก้ไข (Troubleshooting Log)

### ปัญหาที่ 1: make: No rule to make target 'clean' (WSL2)

**Error:**
```bash
┌──(ikai㉿TheerapatLin)-[/mnt/c/.../whisper-source]
└─$ make clean
make: *** No rule to make target 'clean'.  Stop.
```

**สาเหตุ:**
- whisper.cpp ใช้ CMake เป็นหลัก ไม่ใช่ Makefile แบบดั้งเดิม
- โฟลเดอร์ยังไม่มี build files ที่ generate จาก CMake
- Makefile อาจไม่มีหรือไม่สมบูรณ์

**วิธีแก้:**
```bash
# แทนที่จะใช้ make clean
rm -rf build  # ลบโฟลเดอร์ build ทั้งหมด

# Generate build files ด้วย CMake
cmake -B build -DCMAKE_BUILD_TYPE=Release

# คอมไพล์
cmake --build build --config Release

# หรือใช้ make (จะเรียก cmake อัตโนมัติ)
make
```

**ผลลัพธ์:**
- ✅ สร้างโฟลเดอร์ `build/` พร้อม build files
- ✅ คอมไพล์สำเร็จ binary อยู่ใน `build/bin/`

---

### ปัญหาที่ 2: cp: cannot stat 'build/bin/whisper-cli'

**Error:**
```bash
┌──(ikai㉿TheerapatLin)-[/mnt/c/.../whisper-source]
└─$ cp build/bin/whisper-cli ../binary/linux/main
cp: cannot stat 'build/bin/whisper-cli': No such file or directory
```

**สาเหตุ:**
- โฟลเดอร์ `build/` ยังไม่ได้ถูกสร้าง
- คอมไพล์ยังไม่เสร็จหรือล้มเหลว
- Binary อาจชื่ออื่นหรืออยู่ตำแหน่งอื่น

**วิธีตรวจสอบ:**
```bash
# ตรวจสอบว่ามีโฟลเดอร์ build หรือไม่
ls -la build/

# ดูว่ามี binary อะไรบ้าง
find build -name "*whisper*" -type f
find build -name "main" -type f
ls -la build/bin/ 2>/dev/null || echo "ไม่มีโฟลเดอร์ build/bin/"
ls -la bin/ 2>/dev/null || echo "ไม่มีโฟลเดอร์ bin/"
```

**วิธีแก้:**
```bash
# ขั้นตอนที่ 1: คอมไพล์ให้เสร็จก่อน
cmake -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build --config Release

# ขั้นตอนที่ 2: หา binary ที่ถูกสร้าง
find build -type f -executable | grep -E "(whisper|main)"

# ขั้นตอนที่ 3: Copy binary (อาจอยู่ตำแหน่งใดตำแหน่งหนึ่ง)
# ตัวอย่างตำแหน่งที่เป็นไปได้:
cp build/bin/whisper-cli ../binary/linux/main         # ตำแหน่งใหม่
# หรือ
cp build/bin/main ../binary/linux/main                # ตำแหน่งเก่า
# หรือ
cp build/examples/cli/whisper-cli ../binary/linux/main  # ตำแหน่งใน examples

# Verify
chmod +x ../binary/linux/main
ls -lh ../binary/linux/main
file ../binary/linux/main
../binary/linux/main --help  # ทดสอบรัน
```

**ผลลัพธ์:**
- ✅ พบ binary ที่ถูกต้อง
- ✅ Copy สำเร็จ
- ✅ Binary รันได้

---

### ปัญหาที่ 3: Binary Path ไม่ตรงกับที่คาดหวัง

**Whisper.cpp Version ต่างๆ อาจวาง binary คนละที่:**

| Version | Binary Path |
|---------|-------------|
| ≤ v1.5.x | `build/main` หรือ `build/bin/main` |
| ≥ v1.6.x | `build/bin/whisper-cli` |
| Latest | `build/bin/whisper-cli` หรือ `build/examples/cli/whisper-cli` |

**วิธีแก้แบบยืดหยุ่น:**
```bash
# หา binary อัตโนมัติ
BINARY_PATH=$(find build -type f -name "whisper-cli" -o -name "main" | grep -v ".o$" | head -1)

if [ -z "$BINARY_PATH" ]; then
    echo "❌ ไม่พบ binary"
    echo "💡 ลอง list ไฟล์ทั้งหมดใน build:"
    find build -type f -executable
else
    echo "✅ พบ binary: $BINARY_PATH"
    cp "$BINARY_PATH" ../binary/linux/main
    chmod +x ../binary/linux/main
    echo "✅ Copy สำเร็จ"
fi
```

---

## 🆘 ยังแก้ไม่ได้?

### ตรวจสอบ DLL ที่ขาดหาย

ใช้ **Dependency Walker** หรือ **Dependencies.exe**:

```bash
# ดาวน์โหลด Dependencies.exe
# URL: https://github.com/lucasg/Dependencies/releases

# เปิดไฟล์ main.exe ดู DLL ที่ขาดหาย
```

### ดู Detailed Error

```bash
# ใช้ Process Monitor (Procmon)
# URL: https://learn.microsoft.com/en-us/sysinternals/downloads/procmon

# Filter: Process Name is "main.exe"
# ดู DLL loading failures
```

---

## ✅ Checklist การแก้ปัญหา

- [ ] ตรวจสอบ error code (0xc0000135 = DLL not found)
- [ ] เลือกวิธีแก้ (แนะนำ: Visual C++ Redistributable)
- [ ] ติดตั้ง dependencies ที่จำเป็น
- [ ] Restart terminal/เครื่อง (ถ้าจำเป็น)
- [ ] ทดสอบ binary: `.\main.exe --help`
- [ ] รัน Unit Tests: `go test -v`
- [ ] ตรวจสอบผลลัพธ์: PASS ทั้ง 5 tests
- [ ] Document solution ใน team

---

## 📊 สรุปการอัพเดต

### เวอร์ชัน 1.1 (2025-11-10)

**เพิ่มเติม:**
- ✅ แก้ไขขั้นตอนวิธีที่ 3 (WSL2) ให้ใช้ CMake แทน make
- ✅ เพิ่มส่วน Troubleshooting Log สำหรับ error ที่พบจริง
- ✅ เพิ่มวิธีตรวจสอบและหา binary path อัตโนมัติ
- ✅ เพิ่มตารางเปรียบเทียบ binary path ในแต่ละเวอร์ชัน

**ปัญหาที่แก้ไข:**
1. ❌ `make: No rule to make target 'clean'` → ✅ ใช้ `rm -rf build` แทน
2. ❌ `cp: cannot stat 'build/bin/whisper-cli'` → ✅ ใช้ CMake compile ก่อน + หา binary path
3. ❌ Binary path ไม่ตรง → ✅ เพิ่ม script หา binary อัตโนมัติ

**คำสั่งที่ถูกต้องสำหรับ WSL2:**
```bash
# Navigate
cd /mnt/c/Users/boatr/MyBoat/RealFactory/ChatBotProject/backend/whisper/whisper-source

# Clean (ถ้าต้องการ)
rm -rf build

# Build with CMake
cmake -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build --config Release

# Find and copy binary
BINARY_PATH=$(find build -type f -name "whisper-cli" -o -name "main" | grep -v ".o$" | head -1)
cp "$BINARY_PATH" ../binary/linux/main
chmod +x ../binary/linux/main

# Verify
../binary/linux/main --help

# Run tests
cd /mnt/c/Users/boatr/MyBoat/RealFactory/ChatBotProject/backend/test/sst-whisper
go test -v
```

---

**Created**: 2025-11-10
**Last Updated**: 2025-11-10 (v1.1)
**Status**: ✅ Ready to Use - Updated with WSL2 fixes
