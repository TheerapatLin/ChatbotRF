<template>
  <div
    :class="[
      'flex mb-4',
      message.role === 'user' ? 'justify-end' : 'justify-start'
    ]"
  >
    <div
      :class="[
        'max-w-[70%] rounded-2xl px-4 py-3 shadow-sm',
        message.role === 'user'
          ? 'bg-gradient-to-r from-blue-500 to-blue-600 text-white'
          : 'bg-white text-gray-800 border border-gray-200'
      ]"
    >
      <div class="flex items-start gap-2">
        <div v-if="message.role === 'assistant'" class="text-2xl">
          🤖
        </div>
        <div class="flex-1">
          <!-- File Attachments (if any) -->
          <div v-if="hasAttachments" class="mb-3">
            <div class="flex items-center gap-1 mb-2">
              <svg class="w-4 h-4" :class="message.role === 'user' ? 'text-blue-100' : 'text-gray-500'"
                   fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"/>
              </svg>
              <span :class="['text-xs font-medium', message.role === 'user' ? 'text-blue-100' : 'text-gray-500']">
                ไฟล์แนบ ({{ attachments.length }})
              </span>
            </div>
            <div class="space-y-1">
              <div
                v-for="(file, index) in attachments"
                :key="index"
                :class="[
                  'text-xs px-2 py-1.5 rounded-lg flex items-center gap-2',
                  message.role === 'user'
                    ? 'bg-blue-400 bg-opacity-30 text-white'
                    : 'bg-gray-100 text-gray-700'
                ]"
              >
                <span class="font-medium">{{ getFileTypeIcon(file.file_type) }}</span>
                <span class="truncate flex-1" :title="file.filename">{{ file.filename }}</span>
                <span class="text-xs opacity-75">{{ formatFileSize(file.file_size) }}</span>
              </div>
            </div>
          </div>

          <!-- Message Content -->
          <div class="whitespace-pre-wrap break-words">
            {{ message.content }}
          </div>

          <!-- Timestamp and Token Info -->
          <div
            :class="[
              'text-xs mt-2 flex items-center gap-2',
              message.role === 'user' ? 'text-blue-100' : 'text-gray-500'
            ]"
          >
            <span>{{ formatTime(message.created_at) }}</span>
            <span v-if="message.tokens_used" class="flex items-center gap-1">
              <span>•</span>
              <span>{{ message.tokens_used }} tokens</span>
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { formatTime, formatFileSize } from '@/utils/formatters'
import { getFileTypeIcon } from '@/utils/fileHelpers'

const props = defineProps({
  message: {
    type: Object,
    required: true
  }
})

// Check if message has file attachments
const hasAttachments = computed(() => {
  return props.message.file_attachments && props.message.file_attachments.length > 0
})

// Parse file attachments from JSONB string or array
const attachments = computed(() => {
  if (!hasAttachments.value) return []

  const fileData = props.message.file_attachments

  // If it's a string (JSONB from database), parse it
  if (typeof fileData === 'string') {
    try {
      return JSON.parse(fileData)
    } catch (error) {
      console.error('Failed to parse file attachments:', error)
      return []
    }
  }

  // If it's already an array, return it
  return Array.isArray(fileData) ? fileData : []
})
</script>
