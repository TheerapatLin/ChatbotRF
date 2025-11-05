<template>
  <div :class="bubbleClasses">
    <div v-if="role === 'assistant'" class="avatar">
      <span class="avatar-icon">{{ personaIcon }}</span>
    </div>

    <div class="bubble-content">
      <!-- File Attachments -->
      <div v-if="fileAttachments && fileAttachments.length > 0" class="file-attachments">
        <div v-for="file in fileAttachments" :key="file.file_id" class="file-item">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
            <path d="M13 2v7h7" />
          </svg>
          <span>แนบไฟล์ {{ file.file_name || 'Unknown' }}</span>
        </div>
      </div>

      <!-- Message Content -->
      <div class="message-text">
        <p v-if="content" v-html="formattedContent"></p>
        <span v-if="isStreaming" class="streaming-indicator">▊</span>
      </div>

      <!-- Metadata -->
      <div class="message-meta">
        <span class="message-time">{{ formattedTime }}</span>
        <span v-if="tokensUsed" class="message-tokens">{{ tokensUsed }} tokens</span>
      </div>
    </div>

    <div v-if="role === 'user'" class="avatar">
      <span class="avatar-icon">👤</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  role: {
    type: String,
    required: true,
    validator: (value) => ['user', 'assistant', 'system'].includes(value)
  },
  content: {
    type: String,
    default: ''
  },
  timestamp: {
    type: String,
    default: null
  },
  tokensUsed: {
    type: Number,
    default: null
  },
  fileAttachments: {
    type: Array,
    default: () => []
  },
  personaIcon: {
    type: String,
    default: '🤖'
  },
  isStreaming: {
    type: Boolean,
    default: false
  }
})

const bubbleClasses = computed(() => ({
  'chat-bubble': true,
  'bubble-user': props.role === 'user',
  'bubble-assistant': props.role === 'assistant',
  'bubble-system': props.role === 'system',
  'bubble-streaming': props.isStreaming
}))

const formattedContent = computed(() => {
  if (!props.content) return ''

  // Convert newlines to <br>
  let formatted = props.content.replace(/\n/g, '<br>')

  // Simple markdown-like formatting
  formatted = formatted.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
  formatted = formatted.replace(/\*(.*?)\*/g, '<em>$1</em>')
  formatted = formatted.replace(/`(.*?)`/g, '<code>$1</code>')

  return formatted
})

const formattedTime = computed(() => {
  if (!props.timestamp) return ''

  const date = new Date(props.timestamp)
  const now = new Date()
  const diff = now - date

  // Less than 1 minute
  if (diff < 60000) {
    return 'Just now'
  }

  // Less than 1 hour
  if (diff < 3600000) {
    const minutes = Math.floor(diff / 60000)
    return `${minutes}m ago`
  }

  // Less than 24 hours
  if (diff < 86400000) {
    const hours = Math.floor(diff / 3600000)
    return `${hours}h ago`
  }

  // More than 24 hours
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  })
})
</script>

<style scoped>
.chat-bubble {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.bubble-user {
  flex-direction: row-reverse;
}

.bubble-assistant {
  flex-direction: row;
}

.avatar {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.bubble-user .avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.bubble-assistant .avatar {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.avatar-icon {
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.bubble-content {
  flex: 1;
  max-width: 70%;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.bubble-user .bubble-content {
  align-items: flex-end;
}

.bubble-assistant .bubble-content {
  align-items: flex-start;
}

.file-attachments {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 8px;
  font-size: 13px;
  color: #667eea;
}

.file-item svg {
  flex-shrink: 0;
}

.message-text {
  padding: 12px 16px;
  border-radius: 16px;
  word-wrap: break-word;
  position: relative;
}

.bubble-user .message-text {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-bottom-right-radius: 4px;
}

.bubble-assistant .message-text {
  background: white;
  color: #1f2937;
  border: 1px solid #e5e7eb;
  border-bottom-left-radius: 4px;
}

.bubble-system .message-text {
  background: #f3f4f6;
  color: #6b7280;
  text-align: center;
  font-size: 13px;
  font-style: italic;
}

.message-text p {
  margin: 0;
  line-height: 1.5;
}

.message-text :deep(strong) {
  font-weight: 600;
}

.message-text :deep(em) {
  font-style: italic;
}

.message-text :deep(code) {
  background: rgba(0, 0, 0, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 0.9em;
}

.bubble-user .message-text :deep(code) {
  background: rgba(255, 255, 255, 0.2);
}

.streaming-indicator {
  display: inline-block;
  animation: blink 1s infinite;
  margin-left: 2px;
}

@keyframes blink {
  0%, 50% {
    opacity: 1;
  }
  51%, 100% {
    opacity: 0;
  }
}

.message-meta {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: #9ca3af;
  padding: 0 4px;
}

.bubble-user .message-meta {
  justify-content: flex-end;
}

.bubble-assistant .message-meta {
  justify-content: flex-start;
}

.message-time,
.message-tokens {
  white-space: nowrap;
}

@media (max-width: 768px) {
  .bubble-content {
    max-width: 85%;
  }

  .avatar {
    width: 32px;
    height: 32px;
    font-size: 18px;
  }
}
</style>
