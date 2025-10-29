<template>
  <AppLayout>
    <div class="flex flex-col h-full">
      <!-- Header -->
      <header class="bg-white border-b border-gray-200 px-6 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-800">ChatBot AI</h1>
            <p class="text-sm text-gray-500">Powered by OpenAI</p>
          </div>
          <div class="flex items-center gap-3">
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
  </AppLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import MessageList from '@/components/chat/MessageList.vue'
import MessageInput from '@/components/chat/MessageInput.vue'
import { useWebSocket } from '@/composables/useWebSocket'

const { isConnected } = useWebSocket()

onMounted(() => {
  console.log('ChatView mounted')
})
</script>
