import api from './api'

export const chatService = {
  async sendMessage(message, sessionId = null, personaId = null, options = {}) {
    const payload = {
      message,
      session_id: sessionId, // เพิ่ม session_id
      persona_id: personaId,
      use_history: options.useHistory !== undefined ? options.useHistory : true, // เปิด conversation history
      system_prompt: options.systemPrompt || '',
      temperature: options.temperature || 0.7,
      max_tokens: options.maxTokens || 0,
      model: options.model || 'gpt-4o-mini'
    }

    return await api.post('/chat', payload)
  },

  async getChatHistory(sessionId = null, limit = 50, offset = 0) {
    const params = { limit, offset }
    if (sessionId) {
      params.session_id = sessionId
    }

    return await api.get('/chat/history', { params })
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
