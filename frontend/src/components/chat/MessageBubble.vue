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
          <div class="whitespace-pre-wrap break-words">
            {{ message.content }}
          </div>
          <div
            :class="[
              'text-xs mt-2',
              message.role === 'user' ? 'text-blue-100' : 'text-gray-500'
            ]"
          >
            {{ formatTime(message.created_at) }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  message: {
    type: Object,
    required: true
  }
})

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleTimeString('th-TH', {
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>
