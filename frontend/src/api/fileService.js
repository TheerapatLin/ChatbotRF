import apiClient from './axios'

export const fileService = {
  /**
   * Upload files
   * @param {FormData} formData - Form data with files
   */
  uploadFiles: (formData) =>
    apiClient.post('/api/file/uploads', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    }),

  /**
   * Get file history
   * @param {Object} params - Query parameters
   * @param {number} params.limit - Limit
   * @param {number} params.offset - Offset
   * @param {string} params.file_type - File type filter
   */
  getFileHistory: (params = { limit: 50, offset: 0 }) =>
    apiClient.get('/api/file/history', { params }),

  /**
   * Delete all files
   */
  deleteAllFiles: () => apiClient.delete('/api/file/uploads')
}
