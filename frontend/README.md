# ChatBot AI - Frontend

โปรเจ็ค Frontend สำหรับระบบ ChatBot AI ที่สร้างด้วย Vue 3 + Vite

## คุณสมบัติ

- ✅ Vue 3 Composition API
- ✅ Vite สำหรับ Development Server ที่รวดเร็ว
- ✅ Pinia สำหรับ State Management
- ✅ Vue Router สำหรับการจัดการเส้นทาง
- ✅ Tailwind CSS สำหรับการออกแบบ UI แบบมินิมอล
- ✅ WebSocket สำหรับ Real-time Streaming
- ✅ Audio Recording สำหรับ Speech to Speech
- ✅ แถบเมนูด้านขวาสำหรับเลือกโหมดการสนทนา

## โครงสร้างโปรเจ็ค

```
frontend/
├── src/
│   ├── components/
│   │   ├── chat/           # Chat components
│   │   ├── layout/         # Layout components
│   │   └── audio/          # Audio components
│   ├── composables/        # Vue composables
│   ├── router/             # Vue Router
│   ├── services/           # API services
│   ├── store/              # Pinia stores
│   ├── views/              # Page components
│   ├── App.vue
│   └── main.js
├── public/
├── index.html
└── package.json
```

## การติดตั้ง

```bash
# ติดตั้ง dependencies
npm install
```

## การรันโปรเจ็ค

```bash
# รัน development server
npm run dev

# เปิดเบราว์เซอร์ไปที่ http://localhost:5173
```

## การ Build

```bash
# Build สำหรับ production
npm run build

# Preview production build
npm run preview
```

## โหมดการใช้งาน

### Text to Text Mode
- พิมพ์ข้อความในช่องข้อความด้านล่าง
- กด Enter เพื่อส่งข้อความ
- รอรับคำตอบจาก AI

### Speech to Speech Mode
- คลิกปุ่มไมโครโฟน 🎤 เพื่อเริ่มบันทึกเสียง
- พูดข้อความที่ต้องการ
- คลิกปุ่มอีกครั้งเพื่อหยุดและส่งเสียง
- ระบบจะแปลงเสียงเป็นข้อความและตอบกลับ

## Environment Variables

สร้างไฟล์ `.env.development` สำหรับการพัฒนา:

```
VITE_API_BASE_URL=http://localhost:3001/api
VITE_WS_URL=ws://localhost:3001/api/chat/stream
VITE_APP_TITLE=ChatBot AI
```

## เทคโนโลยีที่ใช้

- **Vue 3** - Progressive JavaScript Framework
- **Vite** - Fast Build Tool
- **Pinia** - State Management
- **Vue Router** - Routing
- **Axios** - HTTP Client
- **Tailwind CSS** - Utility-first CSS Framework
- **WebSocket API** - Real-time Communication
- **MediaRecorder API** - Audio Recording

## ขั้นตอนต่อไป

1. รัน Backend Server ที่พอร์ต 3001
2. เริ่มใช้งาน Frontend
3. เลือกโหมดการสนทนา (Text หรือ Speech)
4. เริ่มสนทนากับ AI

## หมายเหตุ

- ต้องการ Backend API ที่รันอยู่ที่ http://localhost:3001
- สำหรับ Speech Mode ต้องอนุญาตให้เข้าถึงไมโครโฟน
- รองรับเฉพาะเบราว์เซอร์ที่รองรับ MediaRecorder API (Chrome, Edge แนะนำ)
