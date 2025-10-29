import api from './api'

export const personaService = {
  async getAllPersonas() {
    return await api.get('/personas')
  },

  async getPersonaById(id) {
    return await api.get(`/personas/${id}`)
  }
}
