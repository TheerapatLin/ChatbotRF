# Text-to-Speech (TTS) Implementation Plan

## Overview

เพิ่มฟีเจอร์ Text-to-Speech เพื่อให้ AI สามารถตอบกลับเป็นเสียงได้ โดยใช้ OpenAI Text-to-Speech API

---

## Architecture

```
User Types/Speaks
    ↓
Frontend sends message
    ↓
Backend receives via WebSocket
    ↓
OpenAI generates text response (streaming)
    ↓
Backend converts text to speech (OpenAI TTS API)
    ↓
Backend sends audio data to Frontend
    ↓
Frontend plays audio
```

---

## OpenAI TTS API

### Endpoint
```
POST https://api.openai.com/v1/audio/speech
```

### Request Format
```json
{
  "model": "tts-1",
  "input": "Text to convert to speech",
  "voice": "alloy",
  "response_format": "mp3",
  "speed": 1.0
}
```

### Available Models
- `tts-1` - Standard quality, faster, cheaper
- `tts-1-hd` - High definition quality, slower, more expensive

### Available Voices
- `alloy` - Neutral
- `echo` - Male
- `fable` - British accent
- `onyx` - Deep male
- `nova` - Female
- `shimmer` - Soft female

### Response Format Options
- `mp3` (default)
- `opus` - Better for real-time streaming
- `aac` - Good compression
- `flac` - Lossless, larger file
- `wav` - Uncompressed, largest file
- `pcm` - Raw audio

### Pricing (as of 2024)
- `tts-1`: $0.015 / 1K characters
- `tts-1-hd`: $0.030 / 1K characters

---

## Implementation Options

### Option 1: Convert Full Response (Recommended)
**Flow:**
1. รับ text response จาก OpenAI chat completion
2. รอให้ streaming เสร็จ (ได้ full text)
3. ส่ง full text ไป OpenAI TTS API
4. ส่ง audio file กลับไป Frontend
5. Frontend เล่นเสียง

**Pros:**
- เสียงต่อเนื่อง ไม่ขาดๆ
- TTS quality ดีกว่า (ได้ context เต็ม)
- ง่ายต่อการ implement

**Cons:**
- ต้องรอ response เสร็จก่อน
- Latency สูงกว่า

### Option 2: Streaming TTS (Advanced)
**Flow:**
1. รับ text chunks จาก streaming
2. แบ่ง text เป็น sentences
3. ส่ง sentence ไป TTS API ทันที
4. Stream audio chunks กลับ Frontend
5. Frontend เล่นเสียงแบบ streaming

**Pros:**
- Latency ต่ำกว่า
- ประสบการณ์ real-time ดีกว่า

**Cons:**
- ซับซ้อนกว่ามาก
- ต้องจัดการ audio buffering
- อาจมีเสียงขาดๆ

---

## Implementation Steps

### Phase 1: Backend - Add TTS Service

#### 1.1 Create TTS Service

**File:** `backend/services/tts_service.go`

```go
package services

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
)

type TTSService struct {
	client *openai.Client
	model  string // "tts-1" or "tts-1-hd"
	voice  string // voice name
}

func NewTTSService(apiKey string) *TTSService {
	client := openai.NewClient(apiKey)
	return &TTSService{
		client: client,
		model:  "tts-1",
		voice:  "nova", // default voice
	}
}

// TextToSpeech converts text to speech audio
func (s *TTSService) TextToSpeech(ctx context.Context, text string) ([]byte, error) {
	req := openai.CreateSpeechRequest{
		Model: s.model,
		Input: text,
		Voice: openai.VoiceNova,
		ResponseFormat: openai.SpeechResponseFormatMp3,
		Speed: 1.0,
	}

	response, err := s.client.CreateSpeech(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("TTS API error: %w", err)
	}
	defer response.Close()

	// Read audio data
	audioData, err := io.ReadAll(response)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}

	return audioData, nil
}

// SetVoice changes the TTS voice
func (s *TTSService) SetVoice(voice string) {
	s.voice = voice
}

// SetModel changes the TTS model
func (s *TTSService) SetModel(model string) {
	s.model = model
}
```

#### 1.2 Update WebSocket Controller

**File:** `backend/controllers/websocket_controller.go`

```go
// Add to WSMessage struct
type WSMessage struct {
	Type         string `json:"type"`
	Content      string `json:"content"`
	PersonaID    *int   `json:"persona_id"`
	SystemPrompt string `json:"system_prompt"`
	EnableTTS    bool   `json:"enable_tts"`    // ✅ New: Enable voice response
	TTSVoice     string `json:"tts_voice"`     // ✅ New: Voice preference
}

// Add to WSResponse struct
type WSResponse struct {
	Type       string `json:"type"`
	Content    string `json:"content"`
	Done       bool   `json:"done"`
	MessageID  string `json:"message_id,omitempty"`
	TokensUsed int    `json:"tokens_used,omitempty"`
	AudioData  string `json:"audio_data,omitempty"` // ✅ New: Base64 encoded audio
	AudioURL   string `json:"audio_url,omitempty"`  // ✅ New: Audio file URL
}
```

#### 1.3 Add TTS Logic to handleMessage

```go
streamDone:
	// After streaming is complete, convert to speech if requested
	if msg.EnableTTS {
		// Convert response to speech
		audioData, err := ctrl.ttsService.TextToSpeech(ctx, fullContent)
		if err != nil {
			log.Printf("TTS conversion failed: %v", err)
		} else {
			// Encode audio to base64
			audioBase64 := base64.StdEncoding.EncodeToString(audioData)

			// Send audio data
			if err := ctrl.sendAudio(c, audioBase64); err != nil {
				log.Printf("Failed to send audio: %v", err)
			}
		}
	}

	// ... rest of the code
```

---

### Phase 2: Frontend - Audio Playback

#### 2.1 Create Audio Player Component

**File:** `frontend/src/components/audio/AudioPlayer.vue`

```vue
<template>
  <div class="audio-player">
    <button
      @click="togglePlay"
      :disabled="!audioUrl"
      class="p-2 rounded-full bg-blue-500 hover:bg-blue-600 text-white disabled:bg-gray-300"
    >
      <svg v-if="!isPlaying" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
        <!-- Play icon -->
        <path d="M6.3 2.841A1.5 1.5 0 004 4.11V15.89a1.5 1.5 0 002.3 1.269l9.344-5.89a1.5 1.5 0 000-2.538L6.3 2.84z"/>
      </svg>
      <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
        <!-- Pause icon -->
        <path d="M5.75 3a.75.75 0 00-.75.75v12.5c0 .414.336.75.75.75h1.5a.75.75 0 00.75-.75V3.75A.75.75 0 007.25 3h-1.5zM12.75 3a.75.75 0 00-.75.75v12.5c0 .414.336.75.75.75h1.5a.75.75 0 00.75-.75V3.75a.75.75 0 00-.75-.75h-1.5z"/>
      </svg>
    </button>

    <audio
      ref="audioElement"
      :src="audioUrl"
      @ended="onEnded"
      @play="isPlaying = true"
      @pause="isPlaying = false"
    />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  audioData: {
    type: String, // Base64 encoded audio
    default: null
  }
})

const audioElement = ref(null)
const audioUrl = ref(null)
const isPlaying = ref(false)

// Convert base64 to blob URL
watch(() => props.audioData, (newData) => {
  if (newData) {
    // Decode base64
    const binaryString = atob(newData)
    const bytes = new Uint8Array(binaryString.length)
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i)
    }

    // Create blob and URL
    const blob = new Blob([bytes], { type: 'audio/mpeg' })
    audioUrl.value = URL.createObjectURL(blob)

    // Auto-play
    setTimeout(() => {
      if (audioElement.value) {
        audioElement.value.play()
      }
    }, 100)
  }
})

const togglePlay = () => {
  if (audioElement.value) {
    if (isPlaying.value) {
      audioElement.value.pause()
    } else {
      audioElement.value.play()
    }
  }
}

const onEnded = () => {
  isPlaying.value = false
}
</script>
```

#### 2.2 Update MessageBubble to Show Audio Player

```vue
<template>
  <div class="message-bubble">
    <!-- Text content -->
    <div class="message-text">{{ message.content }}</div>

    <!-- Audio player (if has audio) -->
    <AudioPlayer
      v-if="message.audioData"
      :audioData="message.audioData"
    />
  </div>
</template>
```

#### 2.3 Update Chat Store to Handle Audio

```javascript
// In chat store
finishStreaming(data) {
  const assistantMessage = {
    id: data.messageId,
    role: 'assistant',
    content: this.streamingContent,
    persona_id: this.currentPersonaId,
    tokens_used: data.tokensUsed,
    audioData: data.audioData, // ✅ Store audio data
    created_at: new Date().toISOString()
  }

  this.addMessage(assistantMessage)
  this.isStreaming = false
  this.streamingContent = ''
}
```

#### 2.4 Update WebSocket Composable

```javascript
ws.value.onmessage = (event) => {
  const message = JSON.parse(event.data)

  if (message.type === 'chunk') {
    if (message.done) {
      // Include audio data if available
      chatStore.finishStreaming({
        messageId: message.message_id,
        tokensUsed: message.tokens_used,
        audioData: message.audio_data // ✅ Pass audio data
      })
    } else {
      chatStore.updateStreamingContent(message.content)
    }
  }
}
```

---

### Phase 3: UI Controls

#### 3.1 Add TTS Toggle in AppSidebar

```vue
<template>
  <div class="sidebar">
    <!-- Existing content -->

    <!-- TTS Settings (for Speech mode) -->
    <div v-if="chatMode === 'speech'" class="mt-3 p-3 bg-gray-50 rounded-lg">
      <div class="flex items-center justify-between mb-2">
        <label class="text-sm font-medium">Voice Response</label>
        <input
          v-model="enableTTS"
          type="checkbox"
          class="toggle"
        />
      </div>

      <!-- Voice Selector -->
      <select
        v-if="enableTTS"
        v-model="selectedVoice"
        class="w-full px-3 py-2 border rounded-lg text-sm"
      >
        <option value="alloy">Alloy (Neutral)</option>
        <option value="echo">Echo (Male)</option>
        <option value="fable">Fable (British)</option>
        <option value="onyx">Onyx (Deep Male)</option>
        <option value="nova">Nova (Female)</option>
        <option value="shimmer">Shimmer (Soft Female)</option>
      </select>
    </div>
  </div>
</template>
```

#### 3.2 Update MessageInput to Send TTS Flag

```javascript
const sendMessage = async () => {
  const messagePayload = {
    content: message.value,
    persona_id: chatStore.currentPersonaId,
    system_prompt: chatStore.systemPrompt || '',
    enable_tts: chatStore.enableTTS, // ✅ Add TTS flag
    tts_voice: chatStore.selectedVoice // ✅ Add voice preference
  }
  wsConnection.sendMessage(messagePayload)
}
```

---

## File Storage Option (Alternative)

แทนที่จะส่ง base64 audio ผ่าน WebSocket, สามารถบันทึกไฟล์และส่ง URL:

### Backend - Save Audio Files

```go
// Save audio to file system
audioFilename := fmt.Sprintf("audio_%s.mp3", assistantMessage.ID.String())
audioPath := filepath.Join("public", "audio", audioFilename)

if err := os.WriteFile(audioPath, audioData, 0644); err != nil {
    log.Printf("Failed to save audio file: %v", err)
}

// Send audio URL instead of data
audioURL := fmt.Sprintf("/audio/%s", audioFilename)
ctrl.sendDone(c, assistantMessage.ID.String(), tokensUsed, audioURL)
```

### Frontend - Play from URL

```javascript
// Simpler - just use URL directly
<audio :src="message.audioUrl" autoplay />
```

---

## Environment Configuration

### Backend `.env.development`

```env
# Existing configs...

# TTS Configuration
TTS_MODEL=tts-1          # or tts-1-hd
TTS_DEFAULT_VOICE=nova   # default voice
TTS_ENABLED=true         # enable/disable TTS feature
```

### Frontend `.env.development`

```env
# Existing configs...

# Audio Configuration
VITE_AUDIO_AUTOPLAY=true
VITE_AUDIO_FORMAT=mp3
```

---

## Cost Considerations

### Example Costs (tts-1)

- 100 characters = $0.0015
- 1,000 characters = $0.015
- 10,000 characters = $0.15

**Average AI response:** ~500 characters = **$0.0075 per response**

### Cost Optimization

1. **Use tts-1 instead of tts-1-hd**
   - 50% cheaper
   - Good enough quality for most use cases

2. **Cache common responses**
   - Store frequently used audio files
   - Reuse instead of regenerating

3. **Let users opt-in**
   - Don't enable TTS by default
   - Only generate when requested

4. **Rate limiting**
   - Limit TTS requests per user
   - Prevent abuse

---

## Testing Plan

### Unit Tests

```go
func TestTTSService_TextToSpeech(t *testing.T) {
    service := NewTTSService(apiKey)

    audio, err := service.TextToSpeech(context.Background(), "Hello world")
    assert.NoError(t, err)
    assert.NotNil(t, audio)
    assert.Greater(t, len(audio), 0)
}
```

### Integration Tests

1. ✅ Send text message with `enable_tts: true`
2. ✅ Verify audio data received
3. ✅ Verify audio plays in browser
4. ✅ Test different voices
5. ✅ Test error handling

---

## Implementation Priority

### Must Have (MVP)
- ✅ Basic TTS conversion (full response)
- ✅ Audio playback in Frontend
- ✅ Toggle TTS on/off

### Nice to Have
- Voice selection
- Audio file caching
- Download audio option

### Future Enhancements
- Streaming TTS
- Multiple languages
- Custom voice fine-tuning
- Audio speed control

---

## Summary

### Approach: Convert Full Response (Recommended for MVP)

**Steps:**
1. Add `TTSService` to Backend
2. Update WebSocket to accept `enable_tts` flag
3. Convert text response to audio using OpenAI TTS API
4. Send base64 audio via WebSocket
5. Play audio in Frontend component

**Timeline:**
- Backend implementation: 2-3 hours
- Frontend implementation: 2-3 hours
- Testing: 1-2 hours
- **Total: 5-8 hours**

---

**Next Steps:**
1. Confirm approach with user
2. Implement Backend TTS service
3. Implement Frontend audio player
4. Test end-to-end flow

