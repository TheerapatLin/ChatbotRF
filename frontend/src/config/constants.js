/**
 * Application-wide constants
 */

// API Configuration
export const API_TIMEOUT = 30000 // 30 seconds
export const MESSAGE_TIMEOUT = 30000 // 30 seconds

// WebSocket Configuration
export const WS_RECONNECT_DELAY = 3000 // 3 seconds
export const WS_MAX_RECONNECT_ATTEMPTS = 5

// File Upload Configuration
export const MAX_FILE_SIZE = 20 * 1024 * 1024 // 20MB
export const MAX_FILES_PER_UPLOAD = 5
export const ACCEPTED_FILE_TYPES = '.pdf,.docx,.xlsx,.txt,.png,.jpg,.jpeg,.gif,.webp'
export const ACCEPTED_MIME_TYPES = [
  'application/pdf',
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  'text/plain',
  'image/png',
  'image/jpeg',
  'image/gif',
  'image/webp'
]

// Chat Configuration
export const CHAT_HISTORY_LIMIT = 50
export const CHAT_HISTORY_OFFSET = 0
export const DEFAULT_PERSONA_ID = 1

// Audio Configuration
export const AUDIO_SAMPLE_RATE = 16000
export const AUDIO_CHANNELS = 1
export const AUDIO_BIT_RATE = 128000

// UI Configuration
export const SCROLL_THRESHOLD = 100 // pixels from bottom to trigger auto-scroll
export const DEBOUNCE_DELAY = 300 // milliseconds

// Session Configuration
export const SESSION_STORAGE_KEY = 'current_session_id'
export const PERSONA_STORAGE_KEY = 'current_persona_id'
