<template>
  <div class="w-64 bg-white border-l border-gray-200 flex flex-col h-full">
    <!-- Header -->
    <div class="p-4 border-b border-gray-200">
      <h2 class="text-lg font-semibold text-gray-800">โหมดการสนทนา</h2>
    </div>

    <!-- Mode Selection -->
    <div class="flex-1 p-4 space-y-3">
      <button
        @click="selectMode('text')"
        :class="[
          'w-full p-4 rounded-lg border-2 transition-all text-left',
          chatMode === 'text'
            ? 'border-blue-500 bg-blue-50'
            : 'border-gray-200 hover:border-gray-300'
        ]"
      >
        <div class="flex items-center gap-3">
          <div class="text-3xl">💬</div>
          <div class="flex-1">
            <div class="font-semibold text-gray-800">Text to Text</div>
            <div class="text-sm text-gray-500">พิมพ์ข้อความเพื่อสนทนา</div>
          </div>
          <div v-if="chatMode === 'text'" class="text-blue-500">
            <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
            </svg>
          </div>
        </div>
      </button>

      <button
        @click="selectMode('speech')"
        :class="[
          'w-full p-4 rounded-lg border-2 transition-all text-left',
          chatMode === 'speech'
            ? 'border-blue-500 bg-blue-50'
            : 'border-gray-200 hover:border-gray-300'
        ]"
      >
        <div class="flex items-center gap-3">
          <div class="text-3xl">🎤</div>
          <div class="flex-1">
            <div class="font-semibold text-gray-800">Speech to Speech</div>
            <div class="text-sm text-gray-500">พูดเพื่อสนทนา</div>
          </div>
          <div v-if="chatMode === 'speech'" class="text-blue-500">
            <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
            </svg>
          </div>
        </div>
      </button>
    </div>

    <!-- Footer Info -->
    <div class="p-4 border-t border-gray-200 bg-gray-50">
      <div class="text-xs text-gray-500 space-y-1">
        <div class="flex justify-between">
          <span>ข้อความ:</span>
          <span class="font-semibold">{{ messageCount }}</span>
        </div>
        <div class="flex justify-between">
          <span>โหมด:</span>
          <span class="font-semibold">{{ chatMode === 'text' ? 'ข้อความ' : 'เสียง' }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useChatStore } from '@/store/chat'

const chatStore = useChatStore()

const chatMode = computed(() => chatStore.chatMode)
const messageCount = computed(() => chatStore.messages.length)

const selectMode = (mode) => {
  chatStore.setChatMode(mode)
}
</script>
