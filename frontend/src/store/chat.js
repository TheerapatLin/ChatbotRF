import { defineStore } from 'pinia'

export const useChatStore = defineStore('chat', {
  state: () => ({
    messages: [],
    isLoading: false,
    isStreaming: false,
    streamingContent: '',
    currentPersonaId: 1,
    chatMode: 'text' // 'text' or 'speech'
  }),

  getters: {
    userMessages: (state) => state.messages.filter(m => m.role === 'user'),
    assistantMessages: (state) => state.messages.filter(m => m.role === 'assistant')
  },

  actions: {
    addMessage(message) {
      this.messages.push(message)
    },

    updateStreamingContent(chunk) {
      this.streamingContent += chunk
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

    clearChat() {
      this.messages = []
      this.streamingContent = ''
    }
  }
})
