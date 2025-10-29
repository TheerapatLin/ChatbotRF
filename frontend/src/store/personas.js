import { defineStore } from 'pinia'
import { personaService } from '@/services/personaService'

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
    async fetchPersonas() {
      this.isLoading = true
      try {
        const response = await personaService.getAllPersonas()
        this.personas = response.personas || response
        console.log('Personas loaded:', this.personas)
      } catch (error) {
        console.error('Failed to fetch personas:', error)
        // Keep default personas if API fails
      } finally {
        this.isLoading = false
      }
    },

    async fetchPersonaById(id) {
      try {
        const persona = await personaService.getPersonaById(id)

        const index = this.personas.findIndex(p => p.id === id)
        if (index !== -1) {
          this.personas[index] = persona
        } else {
          this.personas.push(persona)
        }

        return persona
      } catch (error) {
        console.error(`Failed to fetch persona ${id}:`, error)
        throw error
      }
    },

    setCurrentPersona(id) {
      this.currentPersonaId = id
    }
  }
})
