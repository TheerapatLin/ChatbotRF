<template>
  <div class="chat-log">
    <!-- Empty State -->
    <div v-if="messages.length === 0 && !isLoading" class="empty-state">
      <div class="empty-icon">💬</div>
      <h3>เริ่มการสนทนา</h3>
      <p>ส่งข้อความเพื่อเริ่มคุยกับ AI</p>
    </div>

    <!-- Messages -->
    <div v-else class="messages-container" ref="messagesContainer">
      <ChatBubble
        v-for="message in messages"
        :key="message.id"
        :role="message.role"
        :content="message.content"
        :timestamp="message.timestamp"
        :tokens-used="message.tokens_used"
        :file-attachments="message.file_ids"
        :persona-icon="personaIcon"
        :is-streaming="message.isStreaming"
      />

      <!-- Loading Indicator -->
      <LoadingIndicator
        v-if="isLoading"
        :message="loadingMessage"
        :persona-icon="personaIcon"
      />
    </div>

    <!-- Scroll to Bottom Button -->
    <Transition name="fade">
      <button
        v-if="showScrollButton"
        class="scroll-to-bottom"
        @click="scrollToBottom"
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 5v14M19 12l-7 7-7-7" />
        </svg>
      </button>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useChatStore } from '@/stores/chat'
import { usePersonaStore } from '@/stores/persona'
import ChatBubble from './ChatBubble.vue'
import LoadingIndicator from './LoadingIndicator.vue'

const chatStore = useChatStore()
const personaStore = usePersonaStore()

// Refs
const messagesContainer = ref(null)
const showScrollButton = ref(false)
const userScrolled = ref(false)

// Computed
const messages = computed(() => chatStore.messages)
const isLoading = computed(() => chatStore.isLoading)
const personaIcon = computed(() => personaStore.selectedPersona?.icon || '🤖')

const loadingMessage = computed(() => {
  if (chatStore.currentStreamingMessage) {
    return 'กำลังตอบกลับ...'
  }
  return 'Bot กำลังประมวลผล...'
})

// Methods
const scrollToBottom = (smooth = true) => {
  if (!messagesContainer.value) return

  nextTick(() => {
    const container = messagesContainer.value
    const scrollOptions = {
      top: container.scrollHeight,
      behavior: smooth ? 'smooth' : 'auto'
    }
    container.scrollTo(scrollOptions)
    userScrolled.value = false
  })
}

const handleScroll = () => {
  if (!messagesContainer.value) return

  const container = messagesContainer.value
  const scrollTop = container.scrollTop
  const scrollHeight = container.scrollHeight
  const clientHeight = container.clientHeight

  // Check if user scrolled up
  const isAtBottom = scrollTop + clientHeight >= scrollHeight - 50
  userScrolled.value = !isAtBottom
  showScrollButton.value = !isAtBottom && messages.value.length > 0
}

// Watch for new messages and auto-scroll
watch(
  () => messages.value.length,
  () => {
    if (!userScrolled.value) {
      scrollToBottom()
    }
  }
)

// Watch for streaming updates
watch(
  () => chatStore.currentStreamingMessage?.content,
  () => {
    if (!userScrolled.value) {
      scrollToBottom()
    }
  }
)

// Lifecycle
onMounted(() => {
  // Connect WebSocket
  if (!chatStore.isConnected) {
    chatStore.connectWebSocket().catch((error) => {
      console.error('Failed to connect WebSocket:', error)
    })
  }

  // Add scroll listener
  if (messagesContainer.value) {
    messagesContainer.value.addEventListener('scroll', handleScroll)
  }

  // Initial scroll to bottom
  scrollToBottom(false)
})

onUnmounted(() => {
  // Remove scroll listener
  if (messagesContainer.value) {
    messagesContainer.value.removeEventListener('scroll', handleScroll)
  }
})
</script>

<style scoped>
.chat-log {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
  background: white;
  height: 100%;
  width: 100%;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 40px;
  color: #6b7280;
  height: 100%;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 24px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 8px 0;
}

.empty-state p {
  font-size: 16px;
  margin: 0;
  opacity: 0.8;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 20px;
  scroll-behavior: smooth;
  height: 100%;
  background: #f9fafb;
}

/* Scrollbar Styling */
.messages-container::-webkit-scrollbar {
  width: 8px;
}

.messages-container::-webkit-scrollbar-track {
  background: #f3f4f6;
}

.messages-container::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 4px;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}

.scroll-to-bottom {
  position: absolute;
  bottom: 20px;
  right: 20px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  z-index: 10;
}

.scroll-to-bottom:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.5);
}

.scroll-to-bottom:active {
  transform: scale(0.95);
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: scale(0.8);
}

@media (max-width: 768px) {
  .messages-container {
    padding: 12px;
  }

  .scroll-to-bottom {
    width: 40px;
    height: 40px;
    bottom: 16px;
    right: 16px;
  }
}
</style>
