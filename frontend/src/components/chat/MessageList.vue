<template>
  <div
    ref="messagesContainer"
    class="flex-1 overflow-y-auto p-6 space-y-4"
  >
    <!-- Empty State -->
    <div v-if="messages.length === 0" class="flex items-center justify-center h-full">
      <div class="text-center">
        <div class="text-6xl mb-4">💬</div>
        <h3 class="text-xl font-semibold text-gray-700 mb-2">
          เริ่มต้นการสนทนา
        </h3>
        <p class="text-gray-500">
          พิมพ์ข้อความหรือใช้ไมโครโฟนเพื่อเริ่มต้น
        </p>
      </div>
    </div>

    <!-- Messages -->
    <MessageBubble
      v-for="message in messages"
      :key="message.id || message.created_at"
      :message="message"
    />

    <!-- Streaming Indicator -->
    <div v-if="isStreaming" class="flex justify-start">
      <div class="max-w-[70%] rounded-2xl px-4 py-3 bg-white border border-gray-200">
        <div class="flex items-center gap-2">
          <div class="text-2xl">🤖</div>
          <div class="flex-1">
            <div class="whitespace-pre-wrap break-words">
              {{ streamingContent }}
              <span class="inline-block w-2 h-4 bg-gray-400 animate-pulse ml-1"></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading Indicator -->
    <div v-if="isLoading" class="flex justify-center py-4">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick, computed } from 'vue'
import { useChatStore } from '@/store/chat'
import MessageBubble from './MessageBubble.vue'

const chatStore = useChatStore()
const messagesContainer = ref(null)

const messages = computed(() => chatStore.messages)
const isLoading = computed(() => chatStore.isLoading)
const isStreaming = computed(() => chatStore.isStreaming)
const streamingContent = computed(() => chatStore.streamingContent)

watch([messages, streamingContent], () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
})
</script>
