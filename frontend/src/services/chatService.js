import api from './api'

export const chatService = {
  async sendMessage(message, personaId = null, options = {}) {
    const payload = {
      message,
      persona_id: personaId,
      system_prompt: options.systemPrompt || '',
      temperature: options.temperature || 0,
      max_tokens: options.maxTokens || 0,
      model: options.model || ''
    }

    return await api.post('/chat', payload)
  },

  async getChatHistory(limit = 50, offset = 0) {
    return await api.get('/chat/history', {
      params: { limit, offset }
    })
  },

  async transcribeAudio(audioBlob) {
    const formData = new FormData()
    formData.append('audio', audioBlob, 'recording.webm')

    return await api.post('/audio/transcribe', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}
