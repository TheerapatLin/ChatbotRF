- personas config

- whisper.cpp: v5e-1 TPU, .wave

- brief: การร่างรายละเอียดเกี่ยวกับโปรเจ็ค

- คำว่า POC เป็นคำย่อที่ใช้กันอย่างแพร่หลายในแวดวงธุรกิจ เทคโนโลยี และการพัฒนาผลิตภัณฑ์ ย่อมาจาก Proof of Concept
- Proof of Concept (POC) หรือ การพิสูจน์แนวคิด คือ การดำเนินการในระดับเล็กและจำกัดขอบเขต เพื่อแสดงให้เห็นถึงความเป็นไปได้ หรือความมีชีวิต (Feasibility) ของแนวคิด วิธีการ หรือเทคโนโลยีที่เสนอขึ้นมา ก่อนที่จะทุ่มเททรัพยากรจำนวนมากในการพัฒนาจริง

- mini kube

- OpenAI TTS API ไม่รองรับ การตั้งค่าอารมณ์และ style

- SSML (Speech Synthesis Markup Language) คือ ภาษามาร์กอัปที่ใช้ XML (Extensible Markup Language) ซึ่งถูกออกแบบมาเพื่อให้ผู้พัฒนาสามารถ ควบคุม ลักษณะการสังเคราะห์เสียง (Text-to-Speech หรือ TTS) ที่สร้างโดยระบบคอมพิวเตอร์ได้อย่างละเอียดและแม่นยำยิ่งขึ้น เช่น 
```xml
<break time="500ms"/> (หยุด 500 มิลลิวินาที)

<prosody rate="slow" pitch="high">สวัสดี</prosody> (การปรับ ระดับเสียง (Pitch), ความเร็ว (Rate/Speed), และ ระดับความดัง (Volume))

ฉัน <emphasis level="strong">ต้องการ</emphasis> สิ่งนี้ (การกำหนด การเน้นคำ (Emphasis))

<say-as interpret-as="date" format="mdy">11/06/2025</say-as> (การระบุ รูปแบบการอ่าน สำหรับตัวเลข, วันที่, เวลา, หรือตัวย่อ)
```

- unit test, APi test, Integration test, E2E

- hey gen