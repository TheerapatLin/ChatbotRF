import apiClient from './axios'

export const chatService = {
  /**
   * Send message (non-streaming)
   * @param {Object} data - Message data
   * @param {string} data.message - Message content
   * @param {number} data.persona_id - Persona ID
   * @param {string} data.session_id - Session ID
   * @param {boolean} data.use_history - Include history
   * @param {Array} data.file_ids - File IDs
   */
  sendMessage: (data) => apiClient.post('/api/chat', data),

  /**
   * Get chat history
   * @param {Object} params - Query parameters
   * @param {number} params.limit - Limit
   * @param {number} params.offset - Offset
   */
  getChatHistory: (params = { limit: 50, offset: 0 }) =>
    apiClient.get('/api/chats', { params }),

  /**
   * Delete all messages
   */
  deleteAllMessages: () => apiClient.delete('/api/chats')
}
