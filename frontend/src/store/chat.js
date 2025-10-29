import { defineStore } from 'pinia'
import { chatService } from '@/services/chatService'

export const useChatStore = defineStore('chat', {
  state: () => ({
    messages: [],
    isLoading: false,
    isStreaming: false,
    streamingContent: '',
    currentPersonaId: 1,
    chatMode: 'text', // 'text' or 'speech'
    systemPrompt: '',
    pagination: {
      total: 0,
      limit: 50,
      offset: 0
    }
  }),

  getters: {
    userMessages: (state) => state.messages.filter(m => m.role === 'user'),
    assistantMessages: (state) => state.messages.filter(m => m.role === 'assistant'),
    totalTokensUsed: (state) => {
      return state.messages
        .filter(m => m.tokens_used)
        .reduce((sum, m) => sum + m.tokens_used, 0)
    }
  },

  actions: {
    async fetchHistory(limit = 50, offset = 0) {
      this.isLoading = true
      try {
        const response = await chatService.getChatHistory(limit, offset)
        this.messages = response.messages.reverse() // Newest at bottom
        this.pagination = {
          total: response.total,
          limit: response.limit,
          offset: response.offset
        }
      } catch (error) {
        console.error('Failed to fetch chat history:', error)
        throw error
      } finally {
        this.isLoading = false
      }
    },

    addMessage(message) {
      this.messages.push(message)
    },

    updateStreamingContent(chunk) {
      this.streamingContent += chunk
    },

    startStreaming() {
      this.isStreaming = true
      this.streamingContent = ''
    },

    finishStreaming(data) {
      const assistantMessage = {
        id: data.messageId,
        role: 'assistant',
        content: this.streamingContent,
        persona_id: this.currentPersonaId,
        tokens_used: data.tokensUsed,
        created_at: new Date().toISOString()
      }

      this.addMessage(assistantMessage)
      this.isStreaming = false
      this.streamingContent = ''
    },

    setCurrentPersona(personaId) {
      this.currentPersonaId = personaId
    },

    setChatMode(mode) {
      this.chatMode = mode
    },

    setSystemPrompt(prompt) {
      this.systemPrompt = prompt
    },

    clearChat() {
      this.messages = []
      this.streamingContent = ''
      this.isStreaming = false
    }
  }
})
