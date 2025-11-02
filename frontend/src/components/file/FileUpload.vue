<template>
  <div class="file-upload">
    <input
      ref="fileInput"
      type="file"
      multiple
      accept=".pdf,.docx,.xlsx,.txt,.png,.jpg,.jpeg,.gif,.webp"
      @change="handleFileSelect"
      class="hidden"
    />

    <button
      @click="$refs.fileInput.click()"
      type="button"
      class="upload-btn px-3 py-2 text-sm text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded-lg transition-colors flex items-center gap-2"
      :disabled="disabled"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-5 h-5"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="m18.375 12.739-7.693 7.693a4.5 4.5 0 0 1-6.364-6.364l10.94-10.94A3 3 0 1 1 19.5 7.372L8.552 18.32m.009-.01-.01.01m5.699-9.941-7.81 7.81a1.5 1.5 0 0 0 2.112 2.13"
        />
      </svg>
      <span>แนบไฟล์</span>
    </button>

    <!-- Preview ไฟล์ที่เลือกในรูปแบบ card -->
    <div v-if="selectedFiles.length > 0" class="file-preview mt-3 flex flex-wrap gap-2">
      <div
        v-for="(file, index) in selectedFiles"
        :key="index"
        class="file-card relative bg-white border border-gray-300 rounded-lg px-3 py-2 pr-8 max-w-[250px]"
      >
        <!-- ปุ่มลบที่มุมบนขวา -->
        <button
          @click="removeFile(index)"
          type="button"
          class="absolute top-1 right-1 w-5 h-5 flex items-center justify-center text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-full transition-colors"
          title="ลบไฟล์"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="2"
            stroke="currentColor"
            class="w-4 h-4"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
          </svg>
        </button>

        <!-- ไอคอนไฟล์ -->
        <div class="flex items-start gap-2">
          <div class="file-icon flex-shrink-0 mt-0.5">
            <svg
              v-if="getFileType(file) === 'image'"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-5 h-5 text-blue-500"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
              />
            </svg>
            <svg
              v-else
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-5 h-5 text-gray-500"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
              />
            </svg>
          </div>

          <!-- ชื่อไฟล์และขนาด -->
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-gray-800 truncate" :title="file.name">
              {{ file.name }}
            </p>
            <p class="text-xs text-gray-500">{{ formatFileSize(file.size) }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Error message -->
    <div v-if="errorMessage" class="mt-2 text-sm text-red-500">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { formatFileSize } from '@/utils/formatters'
import { MAX_FILES_PER_UPLOAD, MAX_FILE_SIZE } from '@/config/constants'

const props = defineProps({
  disabled: {
    type: Boolean,
    default: false
  },
  maxFiles: {
    type: Number,
    default: MAX_FILES_PER_UPLOAD
  },
  maxFileSize: {
    type: Number,
    default: MAX_FILE_SIZE
  }
})

const emit = defineEmits(['files-selected'])

const fileInput = ref(null)
const selectedFiles = ref([])
const errorMessage = ref('')

const handleFileSelect = (event) => {
  const files = Array.from(event.target.files)

  // ตรวจสอบจำนวนไฟล์
  if (selectedFiles.value.length + files.length > props.maxFiles) {
    errorMessage.value = `สามารถเลือกได้สูงสุด ${props.maxFiles} ไฟล์`
    setTimeout(() => {
      errorMessage.value = ''
    }, 3000)
    return
  }

  // ตรวจสอบขนาดไฟล์
  const oversizedFiles = files.filter((f) => f.size > props.maxFileSize)
  if (oversizedFiles.length > 0) {
    errorMessage.value = `ไฟล์บางไฟล์มีขนาดเกิน ${formatFileSize(props.maxFileSize)}`
    setTimeout(() => {
      errorMessage.value = ''
    }, 3000)
    return
  }

  selectedFiles.value.push(...files)
  emit('files-selected', selectedFiles.value)

  // Reset input
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const removeFile = (index) => {
  selectedFiles.value.splice(index, 1)
  emit('files-selected', selectedFiles.value)
}

const getFileType = (file) => {
  if (file.type.startsWith('image/')) {
    return 'image'
  }
  return 'document'
}

// Method to clear files (exposed to parent)
const clearFiles = () => {
  selectedFiles.value = []
  emit('files-selected', [])
}

defineExpose({
  clearFiles
})
</script>

<style scoped>
.file-card {
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.file-card:hover {
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
}
</style>
