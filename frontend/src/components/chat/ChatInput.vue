<template>
  <div class="chat-input">
    <!-- Model Selector -->
    <ModelSelector />

    <!-- File Preview -->
    <Transition name="slide-up">
      <div v-if="selectedFiles.length > 0" class="file-preview">
        <div class="file-preview-header">
          <span class="file-count">{{ selectedFiles.length }} file(s) selected</span>
          <button class="clear-files-btn" @click="clearFiles">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="file-list">
          <div v-for="(file, index) in selectedFiles" :key="index" class="file-item">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
              <path d="M13 2v7h7" />
            </svg>
            <span class="file-name">{{ file.name }}</span>
            <span class="file-size">{{ formatFileSize(file.size) }}</span>
            <button class="remove-file-btn" @click="removeFile(index)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Input Area -->
    <div class="input-container">
      <!-- File Upload Button -->
      <button
        class="icon-button"
        :disabled="isLoading || isUploading"
        @click="triggerFileInput"
        title="Attach files"
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48" />
        </svg>
      </button>

      <!-- Hidden File Input -->
      <input
        ref="fileInput"
        type="file"
        multiple
        accept=".txt,.md,.json,.csv,.xml,.pdf,.docx,.xlsx,.jpg,.jpeg,.png,.gif,.webp"
        style="display: none"
        @change="handleFileSelect"
      />

      <!-- Text Input -->
      <textarea
        ref="textInput"
        v-model="message"
        class="message-input"
        placeholder="Type your message..."
        :disabled="isLoading || isUploading"
        @keydown.enter.exact.prevent="handleSend"
        @keydown.shift.enter="handleNewLine"
        @input="handleInput"
      ></textarea>

      <!-- Send Button -->
      <button
        class="send-button"
        :disabled="!canSend"
        @click="handleSend"
      >
        <span v-if="isLoading || isUploading" class="spinner"></span>
        <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M22 2L11 13M22 2l-7 20-4-9-9-4 20-7z" />
        </svg>
      </button>
    </div>

    <!-- Status Bar -->
    <div v-if="isUploading" class="status-bar">
      <span class="status-text">Uploading files...</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import { useChatStore } from '@/stores/chat'
import { usePersonaStore } from '@/stores/persona'
import { useFileStore } from '@/stores/file'
import ModelSelector from './ModelSelector.vue'

const chatStore = useChatStore()
const personaStore = usePersonaStore()
const fileStore = useFileStore()

// Refs
const message = ref('')
const selectedFiles = ref([])
const fileInput = ref(null)
const textInput = ref(null)
const isUploading = ref(false)

// Computed
const isLoading = computed(() => chatStore.isLoading)
const canSend = computed(() => {
  return (message.value.trim() || selectedFiles.value.length > 0) && !isLoading.value && !isUploading.value
})

// Methods
const handleInput = () => {
  // Auto-resize textarea
  if (textInput.value) {
    textInput.value.style.height = 'auto'
    textInput.value.style.height = Math.min(textInput.value.scrollHeight, 150) + 'px'
  }
}

const handleNewLine = () => {
  // Shift+Enter adds new line (default behavior)
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event) => {
  const files = Array.from(event.target.files)
  selectedFiles.value.push(...files)

  // Clear input so same file can be selected again
  event.target.value = ''
}

const removeFile = (index) => {
  selectedFiles.value.splice(index, 1)
}

const clearFiles = () => {
  selectedFiles.value = []
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'

  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

const handleSend = async () => {
  if (!canSend.value) return

  const messageText = message.value.trim()
  const personaId = personaStore.selectedPersona?.id

  if (!personaId) {
    alert('Please select a persona first')
    return
  }

  let fileIds = []

  try {
    // Upload files if any
    if (selectedFiles.value.length > 0) {
      isUploading.value = true

      const result = await fileStore.uploadFiles(selectedFiles.value, {
        session_id: chatStore.sessionId
      })

      fileIds = result.file_ids || []
    }

    // Send message via WebSocket
    await chatStore.sendMessage(messageText, personaId, fileIds)

    // Clear input
    message.value = ''
    selectedFiles.value = []

    // Reset textarea height
    nextTick(() => {
      if (textInput.value) {
        textInput.value.style.height = 'auto'
      }
    })
  } catch (error) {
    console.error('Failed to send message:', error)
    alert(`Failed to send message: ${error.message}`)
  } finally {
    isUploading.value = false
  }
}

// Auto-focus on mount
nextTick(() => {
  textInput.value?.focus()
})
</script>

<style scoped>
.chat-input {
  border-top: 1px solid #e5e7eb;
  background: white;
  display: flex;
  flex-direction: column;
}

.file-preview {
  padding: 12px 16px;
  background: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
}

.file-preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.file-count {
  font-size: 13px;
  font-weight: 500;
  color: #6b7280;
}

.clear-files-btn {
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  border-radius: 4px;
  transition: all 0.2s;
}

.clear-files-btn:hover {
  background: #e5e7eb;
  color: #374151;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 13px;
}

.file-item svg {
  flex-shrink: 0;
  color: #667eea;
}

.file-name {
  flex: 1;
  color: #374151;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  color: #9ca3af;
  font-size: 12px;
}

.remove-file-btn {
  background: none;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  border-radius: 4px;
  transition: all 0.2s;
}

.remove-file-btn:hover {
  background: #fee2e2;
  color: #dc2626;
}

.input-container {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  padding: 16px;
}

.icon-button {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  border: none;
  background: #f3f4f6;
  color: #6b7280;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.icon-button:hover:not(:disabled) {
  background: #e5e7eb;
  color: #374151;
}

.icon-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.message-input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
  color: #1f2937;
  background: white;
  resize: none;
  min-height: 40px;
  max-height: 150px;
  outline: none;
  transition: all 0.2s;
}

.message-input:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.message-input:disabled {
  background: #f3f4f6;
  color: #9ca3af;
  cursor: not-allowed;
}

.message-input::placeholder {
  color: #9ca3af;
}

.send-button {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  border: none;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.send-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.send-button:active:not(:disabled) {
  transform: translateY(0);
}

.send-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.status-bar {
  padding: 8px 16px;
  background: #fef3c7;
  border-top: 1px solid #fde68a;
}

.status-text {
  font-size: 13px;
  color: #92400e;
}

/* Transitions */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
  max-height: 300px;
  overflow: hidden;
}

.slide-up-enter-from,
.slide-up-leave-to {
  max-height: 0;
  opacity: 0;
}

@media (max-width: 768px) {
  .input-container {
    padding: 12px;
  }
}
</style>
