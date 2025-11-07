package tts_el_test

import (
	"reflect"
	"testing"

	"chatbot/utils"
)

// TestChunkText - ทดสอบฟังก์ชันการแบ่ง text เป็น chunks
func TestChunkText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "Simple sentence with comma",
			input: "Hello, world!",
			expected: []string{
				"Hello,",
				"world!",
			},
		},
		{
			name:  "Multiple sentences",
			input: "Hello! How are you? I am fine.",
			expected: []string{
				"Hello!",
				"How are you?",
				"I am fine.",
			},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:  "Text with semicolon",
			input: "First part; second part.",
			expected: []string{
				"First part;",
				"second part.",
			},
		},
		{
			name:  "Text without punctuation",
			input: "Hello world",
			expected: []string{
				"Hello world",
			},
		},
		{
			name:  "Multiple commas",
			input: "One, two, three, four.",
			expected: []string{
				"One,",
				"two,",
				"three,",
				"four.",
			},
		},
		{
			name:  "Mixed punctuation",
			input: "What? Yes! Okay, done.",
			expected: []string{
				"What?",
				"Yes!",
				"Okay,",
				"done.",
			},
		},
		{
			name:  "Thai text with punctuation",
			input: "สวัสดีครับ, คุณสบายดีไหม? ผมสบายดีครับ.",
			expected: []string{
				"สวัสดีครับ,",
				"คุณสบายดีไหม?",
				"ผมสบายดีครับ.",
			},
		},
		{
			name:  "Long sentence with multiple clauses",
			input: "This is the first clause, and this is the second; finally this is the last.",
			expected: []string{
				"This is the first clause,",
				"and this is the second;",
				"finally this is the last.",
			},
		},
		{
			name:  "Single word with exclamation",
			input: "Wow!",
			expected: []string{
				"Wow!",
			},
		},
	}

	// วนลูปทดสอบแต่ละ test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// เรียกใช้ฟังก์ชัน ChunkText และเปรียบเทียบผลลัพธ์
			result := utils.ChunkText(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ChunkText() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestChunkTextEdgeCases - ทดสอบ edge cases ต่างๆ
func TestChunkTextEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Only spaces",
			input:    "   ",
			expected: []string{}, // Empty string after trimming
		},
		{
			name:  "Only punctuation",
			input: ".,!?;",
			expected: []string{
				".,!?;",
			},
		},
		{
			name:  "Repeated punctuation",
			input: "Hello!!! How are you???",
			expected: []string{
				"Hello!!!",
				"How are you???",
			},
		},
		{
			name:  "Multiple spaces between words",
			input: "Hello,    world!",
			expected: []string{
				"Hello,",
				"world!",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ChunkText(tt.input)
			// Handle nil vs empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ChunkText() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// BenchmarkChunkText - ทดสอบประสิทธิภาพของฟังก์ชัน ChunkText
func BenchmarkChunkText(b *testing.B) {
	text := "Hello, world! How are you? I am testing the performance of text chunking. This is a longer sentence with multiple clauses, semicolons; and various punctuation marks!"

	// รัน benchmark
	for i := 0; i < b.N; i++ {
		utils.ChunkText(text)
	}
}

// BenchmarkChunkTextLong - ทดสอบประสิทธิภาพกับข้อความยาว
func BenchmarkChunkTextLong(b *testing.B) {
	// สร้างข้อความยาวมาก
	text := ""
	for i := 0; i < 100; i++ {
		text += "This is sentence number " + string(rune(i)) + ". "
	}

	for i := 0; i < b.N; i++ {
		utils.ChunkText(text)
	}
}
