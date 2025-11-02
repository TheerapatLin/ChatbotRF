<template>
  <div class="file-analysis-card bg-gray-50 border border-gray-200 rounded-lg p-3 hover:bg-gray-100 transition-colors">
    <div class="flex items-start gap-3">
      <!-- File Icon -->
      <div class="flex-shrink-0">
        <div :class="[
          'w-10 h-10 rounded-lg flex items-center justify-center text-white font-semibold text-sm',
          getFileTypeColor(file.fileType || file.file_type)
        ]">
          {{ getFileTypeIcon(file.fileType || file.file_type) }}
        </div>
      </div>

      <!-- File Info -->
      <div class="flex-1 min-w-0">
        <div class="flex items-start justify-between gap-2">
          <div class="flex-1 min-w-0">
            <h4 class="text-sm font-medium text-gray-900 truncate" :title="file.filename">
              {{ file.filename }}
            </h4>
            <div class="flex items-center gap-2 mt-1">
              <span class="text-xs text-gray-500">
                {{ formatFileType(file.fileType || file.file_type) }}
              </span>
              <span v-if="file.fileSize || file.file_size" class="text-xs text-gray-400">
                {{ formatFileSize(file.fileSize || file.file_size) }}
              </span>
            </div>
          </div>

          <!-- Delete Button -->
          <button
            v-if="showDelete"
            @click="handleDelete"
            class="flex-shrink-0 w-6 h-6 flex items-center justify-center text-gray-400 hover:text-red-600 hover:bg-red-50 rounded transition-colors"
            title="ลบไฟล์"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <!-- Analysis Summary (Expandable) -->
        <div v-if="file.summary || file.analysisSummary || file.analysis_summary" class="mt-2">
          <div
            class="text-xs text-gray-600 line-clamp-2"
            :class="{ 'line-clamp-none': expanded }"
          >
            {{ file.summary || file.analysisSummary || file.analysis_summary }}
          </div>
          <button
            v-if="hasLongSummary"
            @click="expanded = !expanded"
            class="text-xs text-blue-600 hover:text-blue-700 font-medium mt-1"
          >
            {{ expanded ? 'แสดงน้อยลง' : 'แสดงเพิ่มเติม' }}
          </button>
        </div>

        <!-- Status Badge (if applicable) -->
        <div v-if="file.status" class="mt-2">
          <span :class="[
            'inline-flex items-center px-2 py-1 rounded text-xs font-medium',
            getStatusClass(file.status)
          ]">
            {{ getStatusText(file.status) }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { formatFileSize } from '@/utils/formatters'
import { getFileTypeIcon, getFileTypeColor } from '@/utils/fileHelpers'

const props = defineProps({
  file: {
    type: Object,
    required: true
  },
  showDelete: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['delete'])

const expanded = ref(false)

const hasLongSummary = computed(() => {
  const summary = props.file.summary || props.file.analysisSummary || props.file.analysis_summary || ''
  return summary.length > 100
})

const handleDelete = () => {
  emit('delete')
}

const formatFileType = (fileType) => {
  if (!fileType) return 'Unknown'

  const type = fileType.toLowerCase()
  if (type.includes('pdf')) return 'PDF'
  if (type.includes('word') || type.includes('docx')) return 'Word Document'
  if (type.includes('excel') || type.includes('xlsx')) return 'Excel Spreadsheet'
  if (type.includes('image')) return 'Image'
  if (type.includes('text')) return 'Text File'
  return fileType
}

const getStatusClass = (status) => {
  switch (status) {
    case 'processing':
      return 'bg-yellow-100 text-yellow-800'
    case 'completed':
      return 'bg-green-100 text-green-800'
    case 'error':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'processing':
      return 'กำลังประมวลผล'
    case 'completed':
      return 'เสร็จสิ้น'
    case 'error':
      return 'ผิดพลาด'
    default:
      return status
  }
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-none {
  display: block;
}
</style>
