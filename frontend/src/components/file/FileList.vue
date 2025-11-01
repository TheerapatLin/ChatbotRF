<template>
  <div v-if="files.length > 0" class="file-list bg-white border border-gray-200 rounded-lg p-4 mb-4">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-sm font-semibold text-gray-700 flex items-center gap-2">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
        </svg>
        ไฟล์ที่อัปโหลด ({{ files.length }})
      </h3>
      <button
        v-if="showClearButton"
        @click="handleClearAll"
        class="text-xs text-red-600 hover:text-red-700 font-medium"
      >
        ลบทั้งหมด
      </button>
    </div>

    <div class="space-y-2">
      <FileAnalysisCard
        v-for="file in files"
        :key="file.fileId"
        :file="file"
        :show-delete="showDelete"
        @delete="handleDelete(file.fileId)"
      />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import FileAnalysisCard from './FileAnalysisCard.vue'
import { useChatStore } from '@/store/chat'

const props = defineProps({
  // If files prop is provided, use it; otherwise use store
  files: {
    type: Array,
    default: null
  },
  showDelete: {
    type: Boolean,
    default: true
  },
  showClearButton: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['delete', 'clear-all'])

const chatStore = useChatStore()

// Use provided files or get from store
const files = computed(() => {
  return props.files || chatStore.uploadedFiles
})

const handleDelete = (fileId) => {
  emit('delete', fileId)

  // If using store files, remove from store
  if (!props.files) {
    chatStore.removeUploadedFile(fileId)
  }
}

const handleClearAll = () => {
  emit('clear-all')

  // If using store files, clear all from store
  if (!props.files) {
    chatStore.clearUploadedFiles()
  }
}
</script>

<style scoped>
.file-list {
  max-height: 300px;
  overflow-y: auto;
}

.file-list::-webkit-scrollbar {
  width: 6px;
}

.file-list::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.file-list::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 3px;
}

.file-list::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}
</style>
