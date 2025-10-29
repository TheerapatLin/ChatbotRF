# Frontend Operation Documentation

## Table of Contents
1. [Technology Stack](#technology-stack)
2. [Project Structure](#project-structure)
3. [Setup Instructions](#setup-instructions)
4. [Core Components](#core-components)
5. [State Management](#state-management)
6. [API Integration](#api-integration)
7. [WebSocket Integration](#websocket-integration)
8. [Audio Recording](#audio-recording)
9. [UI/UX Guidelines](#uiux-guidelines)
10. [Testing Strategy](#testing-strategy)
11. [Deployment](#deployment)

---

## Technology Stack

| Category | Technology | Version | Purpose |
|----------|-----------|---------|---------|
| **Framework** | Vue.js | 3.4+ | Progressive JavaScript framework |
| **Build Tool** | Vite | 5.0+ | Fast build tool and dev server |
| **State Management** | Pinia | 2.1+ | Modern Vue store |
| **Routing** | Vue Router | 4.0+ | Official router for Vue.js |
| **HTTP Client** | Axios | 1.6+ | Promise-based HTTP client |
| **CSS Framework** | Tailwind CSS | 3.4+ | Utility-first CSS framework |
| **UI Components** | Headless UI | 2.0+ | Unstyled, accessible components |
| **Icons** | Heroicons | 2.0+ | Beautiful hand-crafted SVG icons |
| **WebSocket** | Native API | - | Real-time communication |
| **Audio** | MediaRecorder API | - | Audio recording |

---

## Project Structure

```
frontend/
├── public/
│   └── favicon.ico
├── src/
│   ├── assets/              # Static assets (images, fonts)
│   ├── components/          # Vue components
│   │   ├── chat/
│   │   │   ├── MessageBubble.vue
│   │   │   ├── MessageList.vue
│   │   │   ├── MessageInput.vue
│   │   │   └── ChatContainer.vue
│   │   ├── personas/
│   │   │   ├── PersonaCard.vue
│   │   │   ├── PersonaSelector.vue
│   │   │   └── PersonaStats.vue
│   │   ├── audio/
│   │   │   ├── AudioRecorder.vue
│   │   │   └── AudioPlayer.vue
│   │   └── layout/
│   │       ├── AppHeader.vue
│   │       ├── AppSidebar.vue
│   │       └── AppLayout.vue
│   ├── composables/         # Reusable composition functions
│   │   ├── useWebSocket.js
│   │   ├── useAudioRecorder.js
│   │   └── useChat.js
│   ├── router/              # Vue Router configuration
│   │   └── index.js
│   ├── services/            # API services
│   │   ├── api.js
│   │   ├── chatService.js
│   │   ├── personaService.js
│   │   └── audioService.js
│   ├── store/               # Pinia stores
│   │   ├── chat.js
│   │   ├── personas.js
│   │   └── ui.js
│   ├── utils/               # Utility functions
│   │   ├── formatters.js
│   │   └── validators.js
│   ├── views/               # Page components
│   │   ├── HomeView.vue
│   │   ├── ChatView.vue
│   │   └── PersonasView.vue
│   ├── App.vue              # Root component
│   └── main.js              # Application entry point
├── .env.development         # Development environment variables
├── .env.production          # Production environment variables
├── index.html
├── package.json
├── tailwind.config.js
├── vite.config.js
└── README.md
```

---

## Setup Instructions

### 1. Initialize Vue 3 Project with Vite

```bash
# Navigate to project root
cd ChatBotProject

# Create Vue 3 project with Vite
npm create vite@latest frontend -- --template vue

# Navigate to frontend directory
cd frontend

# Install dependencies
npm install
```

### 2. Install Required Packages

```bash
# Core dependencies
npm install vue-router@4 pinia axios

# UI dependencies
npm install -D tailwindcss@latest postcss autoprefixer
npm install @headlessui/vue @heroicons/vue

# Development dependencies
npm install -D @vitejs/plugin-vue eslint prettier
```

### 3. Initialize Tailwind CSS

```bash
# Generate Tailwind config
npx tailwindcss init -p
```

**tailwind.config.js:**
```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f0f9ff',
          100: '#e0f2fe',
          500: '#0ea5e9',
          600: '#0284c7',
          700: '#0369a1',
        }
      }
    },
  },
  plugins: [],
}
```

### 4. Environment Configuration

**.env.development:**
```env
VITE_API_BASE_URL=http://localhost:3001/api
VITE_WS_URL=ws://localhost:3001/api/chat/stream
VITE_APP_TITLE=ChatBot AI
```

**.env.production:**
```env
VITE_API_BASE_URL=https://api.yourdomain.com/api
VITE_WS_URL=wss://api.yourdomain.com/api/chat/stream
VITE_APP_TITLE=ChatBot AI
```

---

## Core Components

### 1. MessageBubble Component

**`src/components/chat/MessageBubble.vue`:**
```vue
<template>
  <div
    :class="[
      'flex mb-4',
      message.role === 'user' ? 'justify-end' : 'justify-start'
    ]"
  >
    <div
      :class="[
        'max-w-[70%] rounded-2xl px-4 py-3 shadow-sm',
        message.role === 'user'
          ? 'bg-gradient-to-r from-blue-500 to-blue-600 text-white'
          : 'bg-white text-gray-800 border border-gray-200'
      ]"
    >
      <div class="flex items-start gap-2">
        <div v-if="message.role === 'assistant'" class="text-2xl">
          {{ personaIcon }}
        </div>
        <div class="flex-1">
          <div class="text-sm font-medium mb-1" v-if="message.role === 'assistant'">
            {{ personaName }}
          </div>
          <div class="whitespace-pre-wrap break-words">
            {{ message.content }}
          </div>
          <div
            :class="[
              'text-xs mt-2',
              message.role === 'user' ? 'text-blue-100' : 'text-gray-500'
            ]"
          >
            {{ formatTime(message.created_at) }}
          </div>
          <div
            v-if="message.tokens_used"
            class="text-xs mt-1 opacity-70"
          >
            {{ message.tokens_used }} tokens
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { usePersonasStore } from '@/store/personas'

const props = defineProps({
  message: {
    type: Object,
    required: true
  }
})

const personasStore = usePersonasStore()

const personaIcon = computed(() => {
  if (props.message.persona_id) {
    const persona = personasStore.getPersonaById(props.message.persona_id)
    return persona?.icon || '🤖'
  }
  return '🤖'
})

const personaName = computed(() => {
  if (props.message.persona_id) {
    const persona = personasStore.getPersonaById(props.message.persona_id)
    return persona?.name || 'AI Assistant'
  }
  return 'AI Assistant'
})

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleTimeString('th-TH', {
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>
```

### 2. MessageList Component

**`src/components/chat/MessageList.vue`:**
```vue
<template>
  <div
    ref="messagesContainer"
    class="flex-1 overflow-y-auto p-4 space-y-4 bg-gray-50"
  >
    <div v-if="isLoading" class="flex justify-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
    </div>

    <MessageBubble
      v-for="message in messages"
      :key="message.id"
      :message="message"
    />

    <div v-if="isStreaming" class="flex justify-start">
      <div class="max-w-[70%] rounded-2xl px-4 py-3 bg-white border border-gray-200">
        <div class="flex items-center gap-2">
          <div class="text-2xl">{{ currentPersonaIcon }}</div>
          <div class="flex-1">
            <div class="text-sm font-medium mb-1">{{ currentPersonaName }}</div>
            <div class="whitespace-pre-wrap break-words">
              {{ streamingContent }}
              <span class="inline-block w-2 h-4 bg-gray-400 animate-pulse ml-1"></span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick, computed } from 'vue'
import { useChatStore } from '@/store/chat'
import { usePersonasStore } from '@/store/personas'
import MessageBubble from './MessageBubble.vue'

const chatStore = useChatStore()
const personasStore = usePersonasStore()
const messagesContainer = ref(null)

const messages = computed(() => chatStore.messages)
const isLoading = computed(() => chatStore.isLoading)
const isStreaming = computed(() => chatStore.isStreaming)
const streamingContent = computed(() => chatStore.streamingContent)

const currentPersonaIcon = computed(() => {
  const persona = personasStore.getCurrentPersona()
  return persona?.icon || '🤖'
})

const currentPersonaName = computed(() => {
  const persona = personasStore.getCurrentPersona()
  return persona?.name || 'AI Assistant'
})

// Auto-scroll to bottom when new messages arrive
watch([messages, streamingContent], () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
})
</script>
```

### 3. MessageInput Component

**`src/components/chat/MessageInput.vue`:**
```vue
<template>
  <div class="border-t bg-white p-4">
    <div class="flex items-end gap-3">
      <button
        @click="toggleRecording"
        :class="[
          'p-3 rounded-full transition-all',
          isRecording
            ? 'bg-red-500 hover:bg-red-600 animate-pulse'
            : 'bg-gray-200 hover:bg-gray-300'
        ]"
        :disabled="isStreaming"
      >
        <MicrophoneIcon class="w-5 h-5" :class="isRecording ? 'text-white' : 'text-gray-600'" />
      </button>

      <div class="flex-1">
        <textarea
          v-model="message"
          @keydown.enter.exact.prevent="sendMessage"
          placeholder="พิมพ์ข้อความ... (Enter เพื่อส่ง, Shift+Enter เพื่อขึ้นบรรทัดใหม่)"
          class="w-full px-4 py-3 border border-gray-300 rounded-2xl focus:outline-none focus:border-blue-500 resize-none"
          rows="1"
          :disabled="isStreaming"
        />
      </div>

      <button
        @click="sendMessage"
        :disabled="!message.trim() || isStreaming"
        class="p-3 rounded-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition-all"
      >
        <PaperAirplaneIcon class="w-5 h-5 text-white" />
      </button>
    </div>

    <div v-if="isRecording" class="mt-2 text-sm text-red-500 flex items-center gap-2">
      <span class="w-2 h-2 bg-red-500 rounded-full animate-pulse"></span>
      กำลังบันทึกเสียง...
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { MicrophoneIcon, PaperAirplaneIcon } from '@heroicons/vue/24/solid'
import { useChatStore } from '@/store/chat'
import { useAudioRecorder } from '@/composables/useAudioRecorder'

const chatStore = useChatStore()
const { isRecording, startRecording, stopRecording } = useAudioRecorder()

const message = ref('')
const isStreaming = computed(() => chatStore.isStreaming)

const sendMessage = async () => {
  if (!message.value.trim() || isStreaming.value) return

  await chatStore.sendMessage(message.value)
  message.value = ''
}

const toggleRecording = async () => {
  if (isRecording.value) {
    const audioBlob = await stopRecording()
    if (audioBlob) {
      await chatStore.transcribeAndSend(audioBlob)
    }
  } else {
    await startRecording()
  }
}
</script>
```

---

## State Management

### Pinia Store - Chat

**`src/store/chat.js`:**
```javascript
import { defineStore } from 'pinia'
import { chatService } from '@/services/chatService'
import { useWebSocket } from '@/composables/useWebSocket'

export const useChatStore = defineStore('chat', {
  state: () => ({
    messages: [],
    isLoading: false,
    isStreaming: false,
    streamingContent: '',
    currentPersonaId: 1,
    pagination: {
      total: 0,
      limit: 50,
      offset: 0
    }
  }),

  getters: {
    userMessages: (state) => state.messages.filter(m => m.role === 'user'),
    assistantMessages: (state) => state.messages.filter(m => m.role === 'assistant'),
    totalTokensUsed: (state) => {
      return state.messages
        .filter(m => m.tokens_used)
        .reduce((sum, m) => sum + m.tokens_used, 0)
    }
  },

  actions: {
    // Fetch chat history
    async fetchHistory(limit = 50, offset = 0) {
      this.isLoading = true
      try {
        const response = await chatService.getChatHistory(limit, offset)
        this.messages = response.messages.reverse() // Newest at bottom
        this.pagination = {
          total: response.total,
          limit: response.limit,
          offset: response.offset
        }
      } catch (error) {
        console.error('Failed to fetch chat history:', error)
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Send message (non-streaming)
    async sendMessage(message) {
      const userMessage = {
        role: 'user',
        content: message,
        created_at: new Date().toISOString()
      }
      this.addMessage(userMessage)

      try {
        const response = await chatService.sendMessage(message, this.currentPersonaId)

        const assistantMessage = {
          id: response.message_id,
          role: 'assistant',
          content: response.reply,
          persona_id: this.currentPersonaId,
          tokens_used: response.tokens_used,
          created_at: response.timestamp
        }
        this.addMessage(assistantMessage)

        return response
      } catch (error) {
        console.error('Failed to send message:', error)
        throw error
      }
    },

    // Send streaming message via WebSocket
    sendStreamingMessage(message) {
      const { sendMessage } = useWebSocket()

      const userMessage = {
        role: 'user',
        content: message,
        created_at: new Date().toISOString()
      }
      this.addMessage(userMessage)

      this.isStreaming = true
      this.streamingContent = ''

      sendMessage(message, this.currentPersonaId)
    },

    // Update streaming content
    updateStreamingContent(chunk) {
      this.streamingContent += chunk
    },

    // Finish streaming
    finishStreaming(data) {
      const assistantMessage = {
        id: data.messageId,
        role: 'assistant',
        content: this.streamingContent,
        persona_id: this.currentPersonaId,
        tokens_used: data.tokensUsed,
        created_at: new Date().toISOString()
      }

      this.addMessage(assistantMessage)
      this.isStreaming = false
      this.streamingContent = ''
    },

    // Add message to list
    addMessage(message) {
      this.messages.push(message)
    },

    // Set current persona
    setCurrentPersona(personaId) {
      this.currentPersonaId = personaId
    },

    // Transcribe audio and send
    async transcribeAndSend(audioBlob) {
      try {
        const text = await chatService.transcribeAudio(audioBlob)
        await this.sendStreamingMessage(text)
      } catch (error) {
        console.error('Failed to transcribe audio:', error)
        throw error
      }
    },

    // Clear chat
    clearChat() {
      this.messages = []
      this.streamingContent = ''
    }
  }
})
```

### Pinia Store - Personas

**`src/store/personas.js`:**
```javascript
import { defineStore } from 'pinia'
import { personaService } from '@/services/personaService'

export const usePersonasStore = defineStore('personas', {
  state: () => ({
    personas: [],
    currentPersonaId: 1,
    isLoading: false
  }),

  getters: {
    getCurrentPersona: (state) => {
      return state.personas.find(p => p.id === state.currentPersonaId)
    },

    getPersonaById: (state) => (id) => {
      return state.personas.find(p => p.id === id)
    }
  },

  actions: {
    async fetchPersonas() {
      this.isLoading = true
      try {
        const response = await personaService.getAllPersonas()
        this.personas = response.personas
      } catch (error) {
        console.error('Failed to fetch personas:', error)
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async fetchPersonaById(id) {
      try {
        const persona = await personaService.getPersonaById(id)

        // Update in personas array
        const index = this.personas.findIndex(p => p.id === id)
        if (index !== -1) {
          this.personas[index] = persona
        } else {
          this.personas.push(persona)
        }

        return persona
      } catch (error) {
        console.error(`Failed to fetch persona ${id}:`, error)
        throw error
      }
    },

    setCurrentPersona(id) {
      this.currentPersonaId = id
    }
  }
})
```

---

## API Integration

### Axios Configuration

**`src/services/api.js`:**
```javascript
import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:3001/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    // Handle errors globally
    if (error.response) {
      switch (error.response.status) {
        case 401:
          console.error('Unauthorized - please login')
          // Redirect to login page
          break
        case 404:
          console.error('Resource not found')
          break
        case 500:
          console.error('Server error')
          break
        default:
          console.error('API error:', error.response.data)
      }
    } else if (error.request) {
      console.error('Network error - no response received')
    } else {
      console.error('Request error:', error.message)
    }
    return Promise.reject(error)
  }
)

export default api
```

### Chat Service

**`src/services/chatService.js`:**
```javascript
import api from './api'

export const chatService = {
  // Send message (non-streaming)
  async sendMessage(message, personaId = null, options = {}) {
    const payload = {
      message,
      persona_id: personaId,
      system_prompt: options.systemPrompt || '',
      temperature: options.temperature || 0,
      max_tokens: options.maxTokens || 0,
      model: options.model || ''
    }

    return await api.post('/chat', payload)
  },

  // Get chat history
  async getChatHistory(limit = 50, offset = 0) {
    return await api.get('/chat/history', {
      params: { limit, offset }
    })
  },

  // Transcribe audio
  async transcribeAudio(audioBlob) {
    const formData = new FormData()
    formData.append('audio', audioBlob, 'recording.webm')

    return await api.post('/audio/transcribe', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}
```

### Persona Service

**`src/services/personaService.js`:**
```javascript
import api from './api'

export const personaService = {
  // Get all personas
  async getAllPersonas() {
    return await api.get('/personas')
  },

  // Get persona by ID
  async getPersonaById(id) {
    return await api.get(`/personas/${id}`)
  }
}
```

---

## WebSocket Integration

### WebSocket Composable

**`src/composables/useWebSocket.js`:**
```javascript
import { ref, onUnmounted } from 'vue'
import { useChatStore } from '@/store/chat'

export function useWebSocket() {
  const ws = ref(null)
  const isConnected = ref(false)
  const reconnectTimer = ref(null)
  const chatStore = useChatStore()

  const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:3001/api/chat/stream'

  function connect() {
    try {
      ws.value = new WebSocket(wsUrl)

      ws.value.onopen = () => {
        console.log('WebSocket connected')
        isConnected.value = true

        // Clear reconnect timer if exists
        if (reconnectTimer.value) {
          clearTimeout(reconnectTimer.value)
          reconnectTimer.value = null
        }
      }

      ws.value.onmessage = (event) => {
        const message = JSON.parse(event.data)
        console.log('WebSocket message received:', message)

        if (message.type === 'chunk') {
          if (message.done) {
            // Streaming complete
            chatStore.finishStreaming({
              messageId: message.message_id,
              tokensUsed: message.tokens_used
            })
          } else {
            // Streaming chunk
            chatStore.updateStreamingContent(message.content)
          }
        } else if (message.type === 'error') {
          console.error('WebSocket error message:', message.error)
          chatStore.isStreaming = false
        }
      }

      ws.value.onerror = (error) => {
        console.error('WebSocket error:', error)
        isConnected.value = false
      }

      ws.value.onclose = () => {
        console.log('WebSocket disconnected')
        isConnected.value = false
        chatStore.isStreaming = false

        // Auto-reconnect after 3 seconds
        reconnectTimer.value = setTimeout(() => {
          console.log('Attempting to reconnect...')
          connect()
        }, 3000)
      }
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
    }
  }

  function sendMessage(content, personaId = 1) {
    if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
      console.error('WebSocket is not connected')
      return false
    }

    const message = {
      type: 'message',
      content,
      persona_id: personaId
    }

    ws.value.send(JSON.stringify(message))
    return true
  }

  function disconnect() {
    if (reconnectTimer.value) {
      clearTimeout(reconnectTimer.value)
      reconnectTimer.value = null
    }

    if (ws.value) {
      ws.value.close()
      ws.value = null
    }

    isConnected.value = false
  }

  // Auto-connect on composable initialization
  connect()

  // Auto-disconnect on component unmount
  onUnmounted(() => {
    disconnect()
  })

  return {
    isConnected,
    sendMessage,
    disconnect,
    reconnect: connect
  }
}
```

---

## Audio Recording

### Audio Recorder Composable

**`src/composables/useAudioRecorder.js`:**
```javascript
import { ref } from 'vue'

export function useAudioRecorder() {
  const isRecording = ref(false)
  const mediaRecorder = ref(null)
  const audioChunks = ref([])
  const stream = ref(null)

  async function startRecording() {
    try {
      // Request microphone access
      stream.value = await navigator.mediaDevices.getUserMedia({
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          sampleRate: 44100
        }
      })

      // Create MediaRecorder
      mediaRecorder.value = new MediaRecorder(stream.value, {
        mimeType: 'audio/webm'
      })

      audioChunks.value = []

      mediaRecorder.value.ondataavailable = (event) => {
        if (event.data.size > 0) {
          audioChunks.value.push(event.data)
        }
      }

      mediaRecorder.value.start()
      isRecording.value = true

      console.log('Recording started')
    } catch (error) {
      console.error('Failed to start recording:', error)
      throw error
    }
  }

  async function stopRecording() {
    return new Promise((resolve, reject) => {
      if (!mediaRecorder.value || !isRecording.value) {
        reject(new Error('No active recording'))
        return
      }

      mediaRecorder.value.onstop = () => {
        const audioBlob = new Blob(audioChunks.value, { type: 'audio/webm' })

        // Stop all tracks
        if (stream.value) {
          stream.value.getTracks().forEach(track => track.stop())
        }

        isRecording.value = false
        console.log('Recording stopped', audioBlob.size, 'bytes')

        resolve(audioBlob)
      }

      mediaRecorder.value.onerror = (error) => {
        reject(error)
      }

      mediaRecorder.value.stop()
    })
  }

  function cancelRecording() {
    if (mediaRecorder.value && isRecording.value) {
      mediaRecorder.value.stop()

      // Stop all tracks
      if (stream.value) {
        stream.value.getTracks().forEach(track => track.stop())
      }

      audioChunks.value = []
      isRecording.value = false

      console.log('Recording cancelled')
    }
  }

  return {
    isRecording,
    startRecording,
    stopRecording,
    cancelRecording
  }
}
```

---

## UI/UX Guidelines

### Design Principles

1. **Conversational Interface**: Chat-like interface similar to popular messaging apps
2. **Real-time Feedback**: Show typing indicators, streaming text, and loading states
3. **Responsive Design**: Mobile-first approach, works on all screen sizes
4. **Accessibility**: Keyboard navigation, screen reader support, ARIA labels
5. **Performance**: Lazy loading, virtual scrolling for long chat histories

### Color Scheme

```css
/* Primary Colors */
--primary-50: #f0f9ff;
--primary-100: #e0f2fe;
--primary-500: #0ea5e9;
--primary-600: #0284c7;
--primary-700: #0369a1;

/* Gray Scale */
--gray-50: #f9fafb;
--gray-100: #f3f4f6;
--gray-500: #6b7280;
--gray-800: #1f2937;

/* Semantic Colors */
--success: #10b981;
--error: #ef4444;
--warning: #f59e0b;
```

### Component Guidelines

**Message Bubbles:**
- User messages: Right-aligned, blue gradient background
- AI messages: Left-aligned, white background with border
- Max width: 70% of container
- Border radius: 18px for rounded appearance
- Include timestamp and persona info

**Input Area:**
- Fixed at bottom of screen
- Auto-resize textarea (max 5 lines)
- Send button disabled when empty or streaming
- Show recording indicator when active

**Persona Selector:**
- Display as cards with icon, name, and description
- Highlight currently selected persona
- Show message count and token usage stats

---

## Testing Strategy

### Unit Tests (Vitest)

```bash
npm install -D vitest @vue/test-utils happy-dom
```

**Example test for MessageBubble:**
```javascript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import MessageBubble from '@/components/chat/MessageBubble.vue'

describe('MessageBubble', () => {
  it('renders user message correctly', () => {
    const message = {
      role: 'user',
      content: 'Hello!',
      created_at: new Date().toISOString()
    }

    const wrapper = mount(MessageBubble, {
      props: { message }
    })

    expect(wrapper.text()).toContain('Hello!')
    expect(wrapper.classes()).toContain('justify-end')
  })

  it('renders assistant message with persona', () => {
    const message = {
      role: 'assistant',
      content: 'Hi there!',
      persona_id: 1,
      created_at: new Date().toISOString()
    }

    const wrapper = mount(MessageBubble, {
      props: { message }
    })

    expect(wrapper.text()).toContain('Hi there!')
    expect(wrapper.classes()).toContain('justify-start')
  })
})
```

### E2E Tests (Playwright)

```bash
npm install -D @playwright/test
```

**Example E2E test:**
```javascript
import { test, expect } from '@playwright/test'

test('send message and receive response', async ({ page }) => {
  await page.goto('http://localhost:5173')

  // Wait for connection
  await expect(page.locator('.status.connected')).toBeVisible()

  // Type message
  await page.fill('textarea', 'Hello AI!')

  // Send message
  await page.click('button[type="submit"]')

  // Wait for response
  await expect(page.locator('.assistant-message')).toBeVisible({ timeout: 10000 })
})
```

---

## Deployment

### Build for Production

```bash
# Build frontend
npm run build

# Preview production build
npm run preview
```

### Nginx Configuration

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    root /var/www/chatbot/dist;
    index index.html;

    # Frontend - Vue Router history mode
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Backend API proxy
    location /api/ {
        proxy_pass http://localhost:3001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # WebSocket proxy
    location /api/chat/stream {
        proxy_pass http://localhost:3001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_read_timeout 86400;
    }

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
}
```

### Docker Deployment

**Dockerfile:**
```dockerfile
# Build stage
FROM node:20-alpine AS build

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine

COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### Environment Variables for Production

Create `.env.production`:
```env
VITE_API_BASE_URL=https://api.yourdomain.com/api
VITE_WS_URL=wss://api.yourdomain.com/api/chat/stream
VITE_APP_TITLE=ChatBot AI
```

---

## Development Checklist

### Day 1: Project Setup
- [ ] Initialize Vue 3 project with Vite
- [ ] Install all dependencies
- [ ] Configure Tailwind CSS
- [ ] Set up environment variables
- [ ] Create project structure

### Day 2: API Integration
- [ ] Create Axios instance with interceptors
- [ ] Implement chat service
- [ ] Implement persona service
- [ ] Implement audio service
- [ ] Set up Pinia stores

### Day 3: Core Components
- [ ] Create MessageBubble component
- [ ] Create MessageList component
- [ ] Create MessageInput component
- [ ] Create PersonaSelector component
- [ ] Implement basic styling

### Day 4: WebSocket Integration
- [ ] Implement useWebSocket composable
- [ ] Connect WebSocket to chat store
- [ ] Handle streaming messages
- [ ] Add error handling and reconnection

### Day 5: Audio Recording
- [ ] Implement useAudioRecorder composable
- [ ] Create AudioRecorder component
- [ ] Integrate with transcription API
- [ ] Test audio recording flow

### Day 6: Polish & Testing
- [ ] Add loading states
- [ ] Implement error handling
- [ ] Write unit tests
- [ ] Write E2E tests
- [ ] Test on different browsers

### Day 7: Deployment
- [ ] Build for production
- [ ] Configure Nginx
- [ ] Set up SSL certificate
- [ ] Deploy to production server
- [ ] Monitor and optimize

---

## Best Practices

1. **Component Composition**: Break down complex components into smaller, reusable pieces
2. **State Management**: Use Pinia stores for global state, props/emits for local communication
3. **Error Handling**: Always wrap API calls in try-catch blocks and show user-friendly error messages
4. **Performance**: Use `v-memo`, `computed` properties, and lazy loading for optimal performance
5. **Accessibility**: Include proper ARIA labels, keyboard navigation, and focus management
6. **Code Quality**: Use ESLint, Prettier, and TypeScript for better code quality
7. **Git Workflow**: Follow conventional commits and create feature branches

---

## Troubleshooting

### WebSocket Connection Issues

**Problem**: WebSocket fails to connect
**Solution**:
- Check if backend server is running
- Verify VITE_WS_URL in .env file
- Check browser console for CORS errors
- Ensure WebSocket endpoint is not blocked by firewall

### Audio Recording Not Working

**Problem**: MediaRecorder API fails
**Solution**:
- Check browser compatibility (Chrome/Edge recommended)
- Ensure HTTPS is enabled (required for getUserMedia)
- Grant microphone permissions
- Check if another app is using the microphone

### Build Errors

**Problem**: Build fails with module errors
**Solution**:
- Delete node_modules and package-lock.json
- Run `npm install` again
- Check Node.js version (use 18+ or 20+)
- Clear Vite cache: `rm -rf node_modules/.vite`

---

## Additional Resources

- [Vue 3 Documentation](https://vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [Pinia Documentation](https://pinia.vuejs.org/)
- [Tailwind CSS Documentation](https://tailwindcss.com/)
- [WebSocket API Documentation](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)
- [MediaRecorder API Documentation](https://developer.mozilla.org/en-US/docs/Web/API/MediaRecorder)

---

**Last Updated**: {{ Current Date }}
**Document Version**: 1.0.0
**Maintained By**: Development Team
