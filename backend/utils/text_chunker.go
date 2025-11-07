package utils

import (
	"regexp"
	"strings"
)

// ChunkText แบ่งข้อความออกเป็น chunks ตามตัวคั่น: space, comma, !, ?, ;, .
// เพื่อให้ระบบสามารถส่ง audio กลับมาทีละส่วนได้เร็วขึ้น
func ChunkText(text string) []string {
	// ถ้าข้อความว่าง ให้คืน array ว่าง
	if text == "" {
		return []string{}
	}

	// ใช้ Regular expression ในการหาคำและ punctuation
	// Pattern นี้จะหาข้อความที่ไม่ใช่ punctuation ตามด้วย punctuation หรือคำเดี่ยว
	re := regexp.MustCompile(`([^.!?,;]+[.!?,;]+|\S+)`)
	matches := re.FindAllString(text, -1)

	var chunks []string      // เก็บ chunks ที่แบ่งแล้ว
	var currentChunk string  // chunk ที่กำลังสร้าง

	// วนลูปผ่านทุก match
	for _, match := range matches {
		match = strings.TrimSpace(match)
		if match == "" {
			continue
		}

		// เพิ่มคำเข้าไปใน chunk ปัจจุบัน
		if currentChunk == "" {
			currentChunk = match
		} else {
			currentChunk += " " + match
		}

		// ตรวจสอบว่าจบด้วย punctuation หรือไม่
		// ถ้าใช่ ให้แบ่ง chunk และเริ่ม chunk ใหม่
		if strings.HasSuffix(match, ".") ||
			strings.HasSuffix(match, "!") ||
			strings.HasSuffix(match, "?") ||
			strings.HasSuffix(match, ",") ||
			strings.HasSuffix(match, ";") {
			chunks = append(chunks, currentChunk)
			currentChunk = ""
		}
	}

	// เพิ่ม chunk สุดท้ายเข้าไป (ถ้ามี)
	if currentChunk != "" {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}
