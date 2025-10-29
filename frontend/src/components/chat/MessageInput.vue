<template>
  <div class="border-t bg-white p-4">
    <div class="flex items-end gap-3 max-w-4xl mx-auto">
      <!-- Microphone Button (for Speech mode) -->
      <button
        v-if="chatMode === 'speech'"
        @click="toggleRecording"
        :class="[
          'p-3 rounded-full transition-all',
          isRecording
            ? 'bg-red-500 hover:bg-red-600 animate-pulse'
            : 'bg-gray-200 hover:bg-gray-300'
        ]"
        :disabled="isStreaming"
      >
        <svg
          class="w-5 h-5"
          :class="isRecording ? 'text-white' : 'text-gray-600'"
          fill="currentColor"
          viewBox="0 0 20 20"
        >
          <path d="M7 4a3 3 0 016 0v6a3 3 0 11-6 0V4zm4 10.93A7.001 7.001 0 0017 8a1 1 0 10-2 0A5 5 0 015 8a1 1 0 00-2 0 7.001 7.001 0 006 6.93V17H6a1 1 0 100 2h8a1 1 0 100-2h-3v-2.07z"/>
        </svg>
      </button>

      <!-- Text Input -->
      <div class="flex-1">
        <textarea
          ref="textareaRef"
          v-model="message"
          @keydown.enter.exact.prevent="sendMessage"
          placeholder="พิมพ์ข้อความ... (Enter เพื่อส่ง, Shift+Enter ขึ้นบรรทัดใหม่)"
          class="w-full px-4 py-3 border border-gray-300 rounded-2xl focus:outline-none focus:border-blue-500 resize-none"
          rows="1"
          :disabled="isStreaming"
        />
      </div>

      <!-- Send Button -->
      <button
        @click="sendMessage"
        :disabled="!message.trim() || isStreaming"
        class="p-3 rounded-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition-all"
      >
        <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
          <path d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z"/>
        </svg>
      </button>
    </div>

    <!-- Recording Indicator -->
    <div v-if="isRecording" class="mt-2 text-sm text-red-500 flex items-center justify-center gap-2">
      <span class="w-2 h-2 bg-red-500 rounded-full animate-pulse"></span>
      กำลังบันทึกเสียง...
    </div>
  </div>
</template>

<script setup>
import { ref, computed, inject, nextTick } from 'vue'
import { useChatStore } from '@/store/chat'
import { useAudioRecorder } from '@/composables/useAudioRecorder'

const chatStore = useChatStore()
const { isRecording, startRecording, stopRecording } = useAudioRecorder()

// Get WebSocket from parent component
const wsConnection = inject('wsConnection', null)

const message = ref('')
const textareaRef = ref(null)
const isStreaming = computed(() => chatStore.isStreaming)
const chatMode = computed(() => chatStore.chatMode)

const sendMessage = async () => {
  if (!message.value.trim() || isStreaming.value) return

  const userMessage = {
    role: 'user',
    content: message.value,
    created_at: new Date().toISOString()
  }

  chatStore.addMessage(userMessage)

  // Send via WebSocket
  if (wsConnection && wsConnection.sendMessage) {
    chatStore.startStreaming()
    const messagePayload = {
      content: message.value,
      persona_id: chatStore.currentPersonaId,
      system_prompt: chatStore.systemPrompt || ''
    }
    wsConnection.sendMessage(messagePayload)
  } else {
    console.error('WebSocket not connected')
  }

  message.value = ''

  // Focus back to textarea after sending
  await nextTick()
  if (textareaRef.value) {
    textareaRef.value.focus()
  }
}

const toggleRecording = async () => {
  if (isRecording.value) {
    const audioBlob = await stopRecording()
    if (audioBlob) {
      console.log('Audio recorded:', audioBlob.size, 'bytes')
      console.log('Transcription integration needed')
    }
  } else {
    await startRecording()
  }
}
</script>
