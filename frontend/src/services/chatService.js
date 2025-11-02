import api from './api'
import { CHAT_HISTORY_LIMIT, CHAT_HISTORY_OFFSET } from '../config/constants'

export const chatService = {
  async sendMessage(message, sessionId = null, personaId = null, options = {}) {
    const payload = {
      message,
      session_id: sessionId,
      persona_id: personaId,
      use_history: options.useHistory !== undefined ? options.useHistory : true,
      system_prompt: options.systemPrompt || '',
      temperature: options.temperature || 0.7,
      max_tokens: options.maxTokens || 0,
      model: options.model || 'gpt-4o-mini'
    }

    return await api.post('/chat', payload)
  },

  async getChatHistory(sessionId = null, limit = CHAT_HISTORY_LIMIT, offset = CHAT_HISTORY_OFFSET) {
    const params = { limit, offset }
    if (sessionId) {
      params.session_id = sessionId
    }

    return await api.get('/chat/history', { params })
  }
}
