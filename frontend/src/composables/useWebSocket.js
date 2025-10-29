import { ref, onUnmounted } from 'vue'
import { useChatStore } from '@/store/chat'

export function useWebSocket() {
  const ws = ref(null)
  const isConnected = ref(false)
  const reconnectTimer = ref(null)
  const chatStore = useChatStore()

  const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:3001/api/chat/stream'

  function connect() {
    try {
      ws.value = new WebSocket(wsUrl)

      ws.value.onopen = () => {
        console.log('WebSocket connected')
        isConnected.value = true

        if (reconnectTimer.value) {
          clearTimeout(reconnectTimer.value)
          reconnectTimer.value = null
        }
      }

      ws.value.onmessage = (event) => {
        const message = JSON.parse(event.data)

        if (message.type === 'chunk') {
          if (message.done) {
            chatStore.finishStreaming({
              messageId: message.message_id,
              tokensUsed: message.tokens_used
            })
          } else {
            chatStore.updateStreamingContent(message.content)
          }
        } else if (message.type === 'error') {
          console.error('WebSocket error message:', message.error)
          chatStore.isStreaming = false
        }
      }

      ws.value.onerror = (error) => {
        console.error('WebSocket error:', error)
        isConnected.value = false
      }

      ws.value.onclose = () => {
        console.log('WebSocket disconnected')
        isConnected.value = false
        chatStore.isStreaming = false

        reconnectTimer.value = setTimeout(() => {
          console.log('Attempting to reconnect...')
          connect()
        }, 3000)
      }
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
    }
  }

  function sendMessage(payload) {
    if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
      console.error('WebSocket is not connected')
      return false
    }

    // Support both old format (string) and new format (object)
    let message
    if (typeof payload === 'string') {
      message = {
        type: 'message',
        content: payload,
        persona_id: 1
      }
    } else {
      message = {
        type: 'message',
        content: payload.content,
        persona_id: payload.persona_id || 1,
        system_prompt: payload.system_prompt || ''
      }
    }

    console.log('Sending WebSocket message:', message)
    ws.value.send(JSON.stringify(message))
    return true
  }

  function disconnect() {
    if (reconnectTimer.value) {
      clearTimeout(reconnectTimer.value)
      reconnectTimer.value = null
    }

    if (ws.value) {
      ws.value.close()
      ws.value = null
    }

    isConnected.value = false
  }

  connect()

  onUnmounted(() => {
    disconnect()
  })

  return {
    isConnected,
    sendMessage,
    disconnect,
    reconnect: connect
  }
}
