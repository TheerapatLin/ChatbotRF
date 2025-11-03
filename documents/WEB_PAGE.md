# Chatbot: Chat Interface

## 1. ส่วนแสดงประวัติการสนทนา (Chat Log / Message Area)
พื้นที่ขนาดใหญ่ที่แสดงข้อความโต้ตอบทั้งหมด รวมถึงผลลัพธ์จากการวิเคราะห์ไฟล์
- Chat Bubble Component: แสดงข้อความแยกกันระหว่างผู้ใช้ (User) และแชทบอท (Bot) อย่างชัดเจน
    - ต้องรองรับการแสดงผลแบบ WebSocket สำหรับข้อความที่ตอบกลับ ในขณะรอการตอบกลับจาก Bot จะมี UX/UI แสดงว่า Bot กำลังประมวลผลอยู่
    - สำหรับ File Analysis: ควรแสดงสรุปการวิเคราะห์หรือผลลัพธ์ที่ได้จากการประมวลผลไฟล์ในรูปแบบที่อ่านง่าย 

## 2. ส่วนป้อนข้อมูลหลัก (Input Area)
ส่วนที่อยู่ด้านล่างสุดของหน้าจอ ใช้สำหรับการสื่อสารพื้นฐาน
- Text Input Field: ช่องสำหรับพิมพ์ข้อความ (Text-to-Text) และรองรับการกด Enter เพื่อส่งข้อความ
- Send Button: ปุ่มสำหรับส่งข้อความที่พิมพ์
- Attachment/Upload Button: ปุ่มสำคัญสำหรับฟังก์ชัน File Upload โดยจะเปิด Modal หรือ File Explorer เพื่อให้ผู้ใช้อัปโหลดไฟล์

## 3. ส่วนควบคุม Speech-to-Speech
องค์ประกอบสำคัญที่เปิดใช้งานฟังก์ชันเสียง
- Mode Switch/Toggle Button: ปุ่มสลับ/เลือกโหมดการสื่อสารหลัก (Text/Speech)
- Microphone Button: ปุ่มหลักสำหรับ เริ่ม/หยุด การบันทึกเสียง
- Status Indicator: แสดงสถานะที่ชัดเจน:
    - "Listening..." (กำลังรอรับเสียงผู้ใช้)
    - "Processing..." (กำลังถอดเสียง/วิเคราะห์)
    - "Bot Speaking..." (กำลังเล่นเสียงคำตอบ)
- Transcript Display: แสดงข้อความที่ถูกถอดเสียงจากคำพูดของผู้ใช้ (Speech-to-Text) ก่อนส่งไปยังแชทบอท
- หลักการทำงานคือ เมื่อ user อัดเสียงเสร็จจะส่ง API POST `/api/audio/transcribe` เพื่อแปลงเสียงเป็นข้อความ จากนั้นส่งข้อความไปให้ AI ด้วย POST `/api/chat` จากนั้นนำ response ที่ได้จาก AI ไป POST `/api/audio/tts` เพื่อแปลงข้อความเป็นเสียง จากนั้นก็ส่งให้ user ฟัง โดยจะเปิดเสียงให้ฟังทันที โดยสามารถหยุดเสียง AI ก่อนที่เสียงจะเล่นจบได้

## 4. ส่วนจัดการไฟล์อัปโหลด (File Upload Modal/Widget)
ส่วนนี้จะปรากฏขึ้นเมื่อผู้ใช้คลิกปุ่มอัปโหลดไฟล์ (Attachment Button)
- Drag-and-Drop Area: พื้นที่ที่สามารถลากไฟล์มาวางได้เพื่อความสะดวกในการใช้งาน
- File List: แสดงชื่อไฟล์และสถานะของไฟล์ที่กำลังรอการประมวลผล (เช่น "Pending", "Uploaded", "Analyzing")
- File Type Restriction Info: แสดงข้อมูลประเภทไฟล์ที่รองรับ (เช่น PDF, DOCX, TXT, CSV, JPEG)
- โดยเมื่อมีการส่งไฟล์ไป frontend จะตรวจสอบว่า user ได้แนบไฟล์มาด้วยมั้ย หากไม่มี ให้ใช้ WebSocket ตามปกติ แต่ถ้าหาก user แนบไฟล์มาด้วย ให้ส่ง API POST `/api/file/upload` เพื่อ upload ไฟล์ไปเก็บไว้ใน Database ก่อน จากนั้นนำ response: file_id เก็บมาเป็น array เพื่อใส่ใน message ของ Websocket ("file_ids": ["e69acc73-6af1-4e44-b5c7-9d6273b048a4"],) ก่อนที่จะส่งไป

## 5. ส่วนตั้งค่า/เครื่องมือทดสอบ (Settings/Configuration Panel)
ส่วนนี้สำคัญมากสำหรับโปรเจกต์ ทดสอบ (Testing) โดยอาจเป็น Sidebar หรือ Modal ที่ซ่อนอยู่
- Model Selector: Dropdown เพื่อเลือก Chatbot Model/Engine ที่จะใช้ (เช่น "Model LLM-A", "Vision Model-B")
- Language Selector: Dropdown สำหรับเลือกภาษาสำหรับ Speech-to-Speech และ Text-to-Text
- System Prompt/Persona Input: ช่องข้อความสำหรับกำหนด System Prompt เพื่อควบคุมบุคลิกภาพหรือบริบทของแชทบอท
- Debug/Stats Display: พื้นที่แสดงข้อมูลทางเทคนิคเพื่อประเมินประสิทธิภาพ:
    - Latency: เวลาตอบกลับ (Response Time)
    - API Usage: ข้อมูลการเรียกใช้ API (ถ้ามี)
    - Clear History Button: ปุ่มสำหรับล้างประวัติการสนทนาบนหน้าจอ