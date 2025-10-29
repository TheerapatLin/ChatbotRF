import { defineStore } from 'pinia'

export const usePersonasStore = defineStore('personas', {
  state: () => ({
    personas: [
      {
        id: 1,
        name: 'AI Assistant',
        icon: '🤖',
        description: 'General purpose AI assistant'
      }
    ],
    currentPersonaId: 1,
    isLoading: false
  }),

  getters: {
    getCurrentPersona: (state) => {
      return state.personas.find(p => p.id === state.currentPersonaId)
    },

    getPersonaById: (state) => (id) => {
      return state.personas.find(p => p.id === id)
    }
  },

  actions: {
    setCurrentPersona(id) {
      this.currentPersonaId = id
    }
  }
})
