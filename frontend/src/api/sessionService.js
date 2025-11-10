import apiClient from './axios'

export const sessionService = {
  /**
   * ดูรายการ sessions จาก messages (group by session_id)
   * @param {number} limit - จำนวน sessions ที่ต้องการ
   * @returns {Promise<Array>} - รายการ sessions
   */
  async getSessionsList(limit = 50) {
    try {
      const response = await apiClient.get('/api/chats', {
        params: { limit: 1000 } // ดึงข้อความมาเยอะๆ เพื่อ group
      })

      // Group messages by session_id
      const sessionMap = new Map()

      response.data.messages.forEach(msg => {
        if (!msg.session_id) return

        if (!sessionMap.has(msg.session_id)) {
          sessionMap.set(msg.session_id, {
            session_id: msg.session_id,
            first_message: msg.content,
            created_at: msg.created_at,
            messages: []
          })
        }
        sessionMap.get(msg.session_id).messages.push(msg)
      })

      // Convert to array and calculate stats
      const sessions = Array.from(sessionMap.values()).map(session => ({
        session_id: session.session_id,
        title: session.first_message.substring(0, 30) +
               (session.first_message.length > 30 ? '...' : ''),
        message_count: session.messages.length,
        last_message_at: session.messages[session.messages.length - 1].created_at,
        created_at: session.created_at
      }))

      // Sort by last message time (newest first)
      return sessions.sort((a, b) =>
        new Date(b.last_message_at) - new Date(a.last_message_at)
      ).slice(0, limit)
    } catch (error) {
      console.error('Failed to get sessions list:', error)
      return []
    }
  },

  /**
   * ดู messages ของ session เฉพาะ
   * @param {string} sessionId - Session ID
   * @returns {Promise} - Response ที่มี messages
   */
  getSessionMessages(sessionId) {
    return apiClient.get(`/api/chats/session/${sessionId}`)
  },

  /**
   * ลบ session (ลบ messages ทั้งหมดใน session)
   * @param {string} sessionId - Session ID ที่จะลบ
   * @returns {Promise}
   */
  deleteSession(sessionId) {
    return apiClient.delete(`/api/chats/session/${sessionId}`)
  }
}
