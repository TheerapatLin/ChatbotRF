<template>
  <AppLayout>
    <!-- Text Mode: Traditional Chat Interface -->
    <div v-if="chatMode === 'text'" class="flex flex-col h-full">
      <!-- Header -->
      <header class="bg-white border-b border-gray-200 px-6 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-800">ChatBot AI</h1>
            <p class="text-sm text-gray-500">Powered by OpenAI</p>
          </div>
          <div class="flex items-center gap-3">
            <!-- New Chat Button -->
            <button
              @click="handleNewChat"
              class="flex items-center gap-2 px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-all"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
              </svg>
              <span class="text-sm font-medium">แชทใหม่</span>
            </button>

            <!-- Connection Status -->
            <div class="flex items-center gap-2 px-3 py-2 bg-gray-100 rounded-lg">
              <div
                :class="[
                  'w-2 h-2 rounded-full',
                  isConnected ? 'bg-green-500' : 'bg-red-500'
                ]"
              ></div>
              <span class="text-sm text-gray-600">
                {{ isConnected ? 'เชื่อมต่อแล้ว' : 'ไม่ได้เชื่อมต่อ' }}
              </span>
            </div>
          </div>
        </div>
      </header>

      <!-- Chat Messages -->
      <MessageList />

      <!-- Input Area -->
      <MessageInput />
    </div>

    <!-- Speech Mode: Speech-to-Speech Interface -->
    <SpeechView v-else-if="chatMode === 'speech'" />
  </AppLayout>
</template>

<script setup>
import { provide, onMounted, computed } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import MessageList from '@/components/chat/MessageList.vue'
import MessageInput from '@/components/chat/MessageInput.vue'
import SpeechView from './SpeechView.vue'
import { useWebSocket } from '@/composables/useWebSocket'
import { useChatStore } from '@/store/chat'

const chatStore = useChatStore()
const wsConnection = useWebSocket()
const { isConnected } = wsConnection

// Get chat mode from store
const chatMode = computed(() => chatStore.chatMode)

// Provide WebSocket connection to child components
provide('wsConnection', wsConnection)

const handleNewChat = () => {
  if (confirm('คุณต้องการเริ่มแชทใหม่ใช่หรือไม่? ข้อความทั้งหมดจะถูกลบ')) {
    chatStore.clearChat()
  }
}

const loadChatHistory = async () => {
  try {
    await chatStore.fetchHistory()
    console.log('Chat history loaded:', chatStore.messages.length, 'messages')
  } catch (error) {
    console.error('Failed to load chat history:', error)
  }
}

onMounted(async () => {
  console.log('ChatView mounted, mode:', chatMode.value)

  // Load chat history from API only in text mode
  if (chatMode.value === 'text') {
    await loadChatHistory()
  }
})
</script>
