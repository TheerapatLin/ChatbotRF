# Frontend Development Guide - ChatBot Project

## 📋 Overview

**AI Chatbot Platform** with Text-to-Text, Speech-to-Speech, File Analysis, and Dynamic AI Personas

**Stack:** Vue 3 (Composition API) + Pinia + Axios + WebSocket + Tailwind CSS

---

## 🔌 Backend API Endpoints

**Base URL:** `http://localhost:3001`

### Personas API
- `GET /api/personas` - List all personas
- `GET /api/personas/:id` - Get persona by ID
- `POST /api/persona` - ✅ Create new persona
- `DELETE /api/persona/:id` - ✅ Delete persona (with cascade delete messages)

### Chat API
- `POST /api/chat` - Send message (non-streaming)
- `WS /api/chat/stream` - Send message (WebSocket streaming)
- `GET /api/chats` - Get chat history
- `DELETE /api/chats` - Delete all messages

### File Upload API
- `POST /api/file/uploads` - Upload files (max 5 files)
- `GET /api/file/history` - Get file history
- `DELETE /api/file/uploads` - Delete all files

### Audio API
- `POST /api/audio/transcribe` - Speech-to-Text (OpenAI Whisper)
- `POST /api/audio/tts` - Text-to-Speech (OpenAI TTS)

---

## 🎨 Main Pages

### 1. 💬 Chat Page
**Components:**
- `MessageList.vue` - Display conversation history
- `MessageInput.vue` - Input box + send button + file attachment
- `MicrophoneButton.vue` - Voice recording
- `PersonaSelector.vue` - Select persona
- `ConfigAICharactor.vue` - ✅ Create new AI persona with full configuration

**Features:**
- WebSocket streaming
- File upload + analysis (PDF, DOCX, Images)
- Speech-to-Speech
- Auto-scroll
- Markdown rendering

### 2. 🎭 Persona Management Page
**Components:**
- `PersonaList.vue` - List of personas
- `PersonaCard.vue` - Display persona details
- `PersonaCreateModal.vue` - ✅ Create new persona
- `PersonaEditModal.vue` - Edit persona (future)

**Features:**
- ✅ Create Persona with custom settings
- ✅ Delete Persona (cascade delete all messages)
- View Persona statistics
- Configure tone, style, temperature, guardrails

---

## ✅ Feature: Create Persona

### API: `POST /api/persona`

**Request Body:**
```typescript
interface PersonaCreateRequest {
  name: string              // Required, max 100 chars
  description: string       // Required
  system_prompt: string     // Required
  tone?: string            // Max 200 chars (default: "friendly")
  style?: string           // Max 500 chars (default: "conversational")
  expertise?: string       // Max 500 chars (default: "general")
  temperature?: number     // 0.0-2.0 (default: 0.7)
  max_tokens?: number      // Default: 2000
  model?: string          // Default: "gpt-4o-mini"
  icon?: string           // Max 10 chars (default: "🤖")
  language_setting?: {
    default_language: string
    response_style: string
    language_code: string
  }
  guardrails?: {
    block_profanity: boolean
    block_sensitive: boolean
    allowed_topics: string[]
    blocked_topics: string[]
    max_response_length: number
    require_moderation: boolean
  }
}
```

**Valid Models:**
- OpenAI: `gpt-4o-mini`, `gpt-4o`, `gpt-4`, `gpt-3.5-turbo`
- Claude: `claude-sonnet-4`, `claude-3-opus`, `claude-3-sonnet`
- AWS Bedrock: `apac.anthropic.claude-sonnet-4-20250514-v1:0`

**Response:**
```json
{
  "id": 9,
  "name": "Marketing Expert",
  "description": "AI ผู้เชี่ยวชาญด้านการตลาด",
  "temperature": 0.6,
  "model": "gpt-4o-mini",
  "icon": "📈",
  "is_active": true
}
```

### Example: Create Persona
```javascript
async function createPersona() {
  const response = await axios.post('/api/persona', {
    name: "Marketing Expert",
    description: "ผู้เชี่ยวชาญด้านการตลาดดิจิทัล",
    system_prompt: "You are a marketing expert...",
    tone: "professional",
    style: "detailed",
    expertise: "marketing",
    temperature: 0.6,
    max_tokens: 2500,
    model: "gpt-4o-mini",
    icon: "📈",
    language_setting: {
      default_language: "th",
      response_style: "formal",
      language_code: "th-TH"
    },
    guardrails: {
      block_profanity: true,
      block_sensitive: false,
      allowed_topics: ["marketing", "business"],
      blocked_topics: ["politics"],
      max_response_length: 3000,
      require_moderation: false
    }
  })
  return response.data
}
```

---

## ✅ Feature: Delete Persona

### API: `DELETE /api/persona/:id`

**Behavior:**
- Deletes all messages associated with the persona first (cascade delete)
- Then deletes the persona itself
- Returns count of deleted messages

**Response:**
```json
{
  "message": "Persona deleted successfully",
  "id": 9,
  "messages_deleted": 42
}
```

**Example:**
```javascript
async function deletePersona(personaId) {
  const response = await axios.delete(`/api/persona/${personaId}`)
  console.log(`Deleted ${response.data.messages_deleted} messages`)
  return response.data
}
```

---

## 🔌 WebSocket Integration

### Connect WebSocket
```javascript
const ws = new WebSocket('ws://localhost:3001/api/chat/stream')

ws.onmessage = (event) => {
  const data = JSON.parse(event.data)

  if (data.type === 'chunk') {
    // Streaming chunk
    if (!data.done) {
      currentMessage.value += data.content
    } else {
      currentMessage.value = ''
    }
  } else if (data.type === 'error') {
    console.error('Error:', data.error)
  }
}
```

### Send Message
```javascript
function sendMessage(message, personaId, fileIds = []) {
  ws.send(JSON.stringify({
    type: 'message',
    content: message,
    persona_id: personaId,
    session_id: sessionId,
    file_ids: fileIds,
    system_prompt: '' // Optional
  }))
}
```

---

## 📁 Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── chat/
│   │   │   ├── MessageList.vue
│   │   │   ├── MessageInput.vue
│   │   │   └── MicrophoneButton.vue
│   │   ├── persona/
│   │   │   ├── PersonaSelector.vue
│   │   │   ├── ConfigAICharactor.vue     ✅ Create Persona Form
│   │   │   ├── PersonaList.vue
│   │   │   └── PersonaCard.vue
│   │   └── file/
│   │       └── FileDragDropArea.vue
│   ├── stores/
│   │   ├── chatStore.js
│   │   ├── personaStore.js              ✅ Persona CRUD
│   │   └── fileStore.js
│   ├── services/
│   │   ├── chatService.js
│   │   ├── personaService.js            ✅ API calls
│   │   └── audioService.js
│   └── views/
│       ├── ChatView.vue
│       └── PersonaManagementView.vue    ✅ Persona CRUD UI
```

---

## 📦 Pinia Store Example

### personaStore.js
```javascript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export const usePersonaStore = defineStore('persona', () => {
  const personas = ref([])
  const selectedPersona = ref(null)

  async function fetchPersonas() {
    const response = await axios.get('/api/personas')
    personas.value = response.data.personas
  }

  async function createPersona(data) {
    const response = await axios.post('/api/persona', data)
    personas.value.push(response.data)
    return response.data
  }

  async function deletePersona(id) {
    const response = await axios.delete(`/api/persona/${id}`)
    personas.value = personas.value.filter(p => p.id !== id)
    return response.data
  }

  return { personas, selectedPersona, fetchPersonas, createPersona, deletePersona }
})
```

---

## ✅ Testing Checklist

### Chat Features
- [ ] Send message via WebSocket streaming
- [ ] Switch persona → tone/style changes
- [ ] Upload files (PDF, DOCX, Images)
- [ ] Speech-to-Text → Chat → Text-to-Speech

### Persona Management
- [x] **List all personas**
- [x] **Create new persona with full configuration**
- [x] **Delete persona (cascade delete messages)**
- [ ] Edit persona (future feature)
- [ ] View persona statistics

---

## 🎯 Summary

### Implemented Features (v6.2)
1. ✅ Chat (Text + Speech)
2. ✅ File Analysis (PDF, DOCX, Images)
3. ✅ 8 Pre-configured AI Personas
4. ✅ WebSocket Streaming
5. ✅ **Create Custom Persona**
6. ✅ **Delete Persona (with cascade delete)**
7. ✅ Full Persona Configuration (tone, style, temperature, guardrails, language settings)

### User Flow
1. **Create Persona** → Configure name, description, system prompt, tone, style, expertise, temperature, model, guardrails
2. **Select Persona** → Choose from persona list
3. **Chat** → Send messages via WebSocket, upload files, use voice
4. **Delete Persona** → Remove persona and all related messages

### Backend API Status
- ✅ `POST /api/persona` - Create persona
- ✅ `DELETE /api/persona/:id` - Delete persona (cascade delete messages)
- ✅ `GET /api/personas` - List personas
- ✅ `GET /api/personas/:id` - Get persona by ID
- ⏳ `PUT /api/persona/:id` - Update persona (future)

---

**Version:** 6.2 (2025-11-04)
**Last Updated:** Persona Management API - Create & Delete with cascade operations
