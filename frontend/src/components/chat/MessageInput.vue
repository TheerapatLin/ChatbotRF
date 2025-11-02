<template>
  <div class="border-t bg-white p-4">
    <div class="max-w-4xl mx-auto">
      <!-- File Upload Component -->
      <div class="mb-3">
        <FileUpload
          ref="fileUploadRef"
          @files-selected="handleFilesSelected"
          :disabled="isStreaming"
        />
      </div>

      <div class="flex items-end gap-3">
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
          :disabled="!canSend || isStreaming"
          class="p-3 rounded-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition-all"
        >
          <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
            <path d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z"/>
          </svg>
        </button>
      </div>

      <!-- Uploading Indicator -->
      <div v-if="isUploading" class="mt-2 text-sm text-blue-500 flex items-center justify-center gap-2">
        <svg class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        กำลังอัปโหลดไฟล์...
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, inject, nextTick } from 'vue'
import { useChatStore } from '@/store/chat'
import FileUpload from '@/components/file/FileUpload.vue'

const chatStore = useChatStore()

// Get WebSocket from parent component
const wsConnection = inject('wsConnection', null)

const message = ref('')
const textareaRef = ref(null)
const fileUploadRef = ref(null)
const selectedFiles = ref([])
const isUploading = ref(false)

const isStreaming = computed(() => chatStore.isStreaming)
const canSend = computed(() => {
  return message.value.trim() || selectedFiles.value.length > 0
})

const handleFilesSelected = (files) => {
  selectedFiles.value = files
}

const sendMessage = async () => {
  if (!canSend.value || isStreaming.value) return

  // อัปโหลดไฟล์ก่อน (ถ้ามี)
  if (selectedFiles.value.length > 0) {
    isUploading.value = true
    try {
      await chatStore.uploadFiles(selectedFiles.value, 'summary')
    } catch (error) {
      console.error('Failed to upload files:', error)
      isUploading.value = false
      return
    }
    isUploading.value = false
  }

  // Build user message with file attachments
  const userMessage = {
    role: 'user',
    content: message.value || '(แนบไฟล์)',
    created_at: new Date().toISOString(),
    file_attachments: chatStore.uploadedFiles.map(file => ({
      file_id: file.fileId,
      filename: file.filename,
      file_type: file.fileType,
      file_size: file.fileSize
    }))
  }

  chatStore.addMessage(userMessage)

  // Send via WebSocket
  if (wsConnection && wsConnection.sendMessage) {
    chatStore.startStreaming()
    const messagePayload = {
      content: message.value || 'วิเคราะห์ไฟล์ที่แนบมา',
      persona_id: chatStore.currentPersonaId,
      system_prompt: chatStore.systemPrompt || ''
    }
    wsConnection.sendMessage(messagePayload)
  } else {
    console.error('WebSocket not connected')
  }

  message.value = ''
  selectedFiles.value = []

  // Clear uploaded files from store
  chatStore.clearUploadedFiles()

  // Clear file upload component
  if (fileUploadRef.value) {
    fileUploadRef.value.clearFiles()
  }

  // Focus back to textarea after sending
  await nextTick()
  if (textareaRef.value) {
    textareaRef.value.focus()
  }
}
</script>
