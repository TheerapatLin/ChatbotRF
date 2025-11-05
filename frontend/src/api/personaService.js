import apiClient from './axios'

export const personaService = {
  /**
   * Get all personas
   */
  getAll: () => apiClient.get('/api/personas'),

  /**
   * Get persona by ID
   * @param {number} id - Persona ID
   */
  getById: (id) => apiClient.get(`/api/personas/${id}`),

  /**
   * Create new persona
   * @param {Object} data - Persona data
   */
  create: (data) => apiClient.post('/api/persona', data),

  /**
   * Update persona
   * @param {number} id - Persona ID
   * @param {Object} data - Updated persona data
   */
  update: (id, data) => apiClient.patch(`/api/persona/${id}`, data),

  /**
   * Delete persona
   * @param {number} id - Persona ID
   */
  delete: (id) => apiClient.delete(`/api/persona/${id}`)
}
