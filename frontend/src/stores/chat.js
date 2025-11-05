import { defineStore } from 'pinia'
import { ref } from 'vue'
import { chatService } from '@/api/chatService'

export const useChatStore = defineStore('chat', () => {
  // State
  const messages = ref([])
  const sessionId = ref(`session_${Date.now()}`)
  const isLoading = ref(false)
  const webSocket = ref(null)
  const isConnected = ref(false)
  const currentStreamingMessage = ref(null)

  // Actions
  const connectWebSocket = () => {
    return new Promise((resolve, reject) => {
      try {
        webSocket.value = new WebSocket('ws://localhost:3001/api/chat/stream')

        webSocket.value.onopen = () => {
          console.log('✅ WebSocket connected')
          isConnected.value = true
          resolve()
        }

        webSocket.value.onmessage = (event) => {
          const data = JSON.parse(event.data)

          if (data.type === 'error') {
            console.error('WebSocket error:', data.error)
            return
          }

          if (data.type === 'chunk') {
            if (!data.done) {
              // Streaming in progress
              if (!currentStreamingMessage.value) {
                currentStreamingMessage.value = {
                  id: Date.now(),
                  role: 'assistant',
                  content: data.content,
                  timestamp: new Date().toISOString(),
                  isStreaming: true
                }
                messages.value.push(currentStreamingMessage.value)
              } else {
                currentStreamingMessage.value.content += data.content
              }
            } else {
              // Streaming complete
              if (currentStreamingMessage.value) {
                currentStreamingMessage.value.isStreaming = false
                currentStreamingMessage.value.message_id = data.message_id
                currentStreamingMessage.value.tokens_used = data.tokens_used
              }
              currentStreamingMessage.value = null
              isLoading.value = false
            }
          }
        }

        webSocket.value.onerror = (error) => {
          console.error('❌ WebSocket error:', error)
          isConnected.value = false
          reject(error)
        }

        webSocket.value.onclose = () => {
          console.log('🔌 WebSocket disconnected')
          isConnected.value = false
        }
      } catch (error) {
        console.error('Failed to connect WebSocket:', error)
        reject(error)
      }
    })
  }

  const disconnectWebSocket = () => {
    if (webSocket.value) {
      webSocket.value.close()
      webSocket.value = null
      isConnected.value = false
    }
  }

  const sendMessage = async (content, personaId, fileIds = []) => {
    if (!content.trim() && fileIds.length === 0) {
      return
    }

    // Add user message to chat
    const userMessage = {
      id: Date.now(),
      role: 'user',
      content: content,
      timestamp: new Date().toISOString(),
      file_ids: fileIds
    }
    messages.value.push(userMessage)

    // Send via WebSocket
    if (!isConnected.value) {
      await connectWebSocket()
    }

    isLoading.value = true

    const messageData = {
      type: 'message',
      content: content,
      persona_id: personaId,
      session_id: sessionId.value,
      file_ids: fileIds
    }

    webSocket.value.send(JSON.stringify(messageData))
  }

  const loadChatHistory = async (limit = 50, offset = 0) => {
    try {
      isLoading.value = true
      const response = await chatService.getChatHistory({ limit, offset })
      messages.value = response.data.messages || []
      return response.data
    } catch (error) {
      console.error('Failed to load chat history:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const clearChatHistory = async () => {
    try {
      await chatService.deleteAllMessages()
      messages.value = []
      sessionId.value = `session_${Date.now()}`
    } catch (error) {
      console.error('Failed to clear chat history:', error)
      throw error
    }
  }

  const resetSession = () => {
    messages.value = []
    sessionId.value = `session_${Date.now()}`
    currentStreamingMessage.value = null
    isLoading.value = false
  }

  return {
    // State
    messages,
    sessionId,
    isLoading,
    webSocket,
    isConnected,
    currentStreamingMessage,

    // Actions
    connectWebSocket,
    disconnectWebSocket,
    sendMessage,
    loadChatHistory,
    clearChatHistory,
    resetSession
  }
})
