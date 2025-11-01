import api from './api'

export const fileService = {
  // Upload และวิเคราะห์ไฟล์
  async analyzeFile(file, options = {}) {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('analysis_type', options.analysisType || 'summary')
    formData.append('language', options.language || 'th')

    if (options.sessionId) {
      formData.append('session_id', options.sessionId)
    }

    if (options.prompt) {
      formData.append('prompt', options.prompt)
    }

    return await api.post('/file/analyze', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // ดึงประวัติไฟล์
  async getFileHistory(sessionId = null, limit = 20, offset = 0) {
    const params = { limit, offset }
    if (sessionId) {
      params.session_id = sessionId
    }

    return await api.get('/file/history', { params })
  }
}
