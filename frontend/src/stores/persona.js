import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { personaService } from '@/api/personaService'

export const usePersonaStore = defineStore('persona', () => {
  // State
  const personas = ref([])
  const selectedPersona = ref(null)
  const isLoadingPersonas = ref(false)
  const error = ref(null)

  // Computed
  const hasPersonas = computed(() => personas.value.length > 0)
  const selectedPersonaId = computed(() => selectedPersona.value?.id || null)

  // Actions
  const fetchPersonas = async () => {
    try {
      isLoadingPersonas.value = true
      error.value = null
      const response = await personaService.getAll()
      personas.value = response.data.personas || []

      // Auto-select first persona if none selected
      if (!selectedPersona.value && personas.value.length > 0) {
        selectedPersona.value = personas.value[0]
      }

      return personas.value
    } catch (err) {
      console.error('Failed to fetch personas:', err)
      error.value = err.response?.data?.error || 'Failed to fetch personas'
      throw err
    } finally {
      isLoadingPersonas.value = false
    }
  }

  const fetchPersonaById = async (id) => {
    try {
      isLoadingPersonas.value = true
      error.value = null
      const response = await personaService.getById(id)
      return response.data
    } catch (err) {
      console.error('Failed to fetch persona:', err)
      error.value = err.response?.data?.error || 'Failed to fetch persona'
      throw err
    } finally {
      isLoadingPersonas.value = false
    }
  }

  const createPersona = async (personaData) => {
    try {
      isLoadingPersonas.value = true
      error.value = null
      const response = await personaService.create(personaData)
      const newPersona = response.data

      // Add to list
      personas.value.push(newPersona)

      return newPersona
    } catch (err) {
      console.error('Failed to create persona:', err)
      error.value = err.response?.data?.error || 'Failed to create persona'
      throw err
    } finally {
      isLoadingPersonas.value = false
    }
  }

  const updatePersona = async (id, personaData) => {
    try {
      isLoadingPersonas.value = true
      error.value = null
      const response = await personaService.update(id, personaData)
      const updatedPersona = response.data

      // Update in list
      const index = personas.value.findIndex((p) => p.id === id)
      if (index !== -1) {
        personas.value[index] = updatedPersona
      }

      // Update selected if it's the one being updated
      if (selectedPersona.value?.id === id) {
        selectedPersona.value = updatedPersona
      }

      return updatedPersona
    } catch (err) {
      console.error('Failed to update persona:', err)
      error.value = err.response?.data?.error || 'Failed to update persona'
      throw err
    } finally {
      isLoadingPersonas.value = false
    }
  }

  const deletePersona = async (id) => {
    try {
      isLoadingPersonas.value = true
      error.value = null
      await personaService.delete(id)

      // Remove from list
      personas.value = personas.value.filter((p) => p.id !== id)

      // Clear selection if deleted
      if (selectedPersona.value?.id === id) {
        selectedPersona.value = personas.value[0] || null
      }
    } catch (err) {
      console.error('Failed to delete persona:', err)
      error.value = err.response?.data?.error || 'Failed to delete persona'
      throw err
    } finally {
      isLoadingPersonas.value = false
    }
  }

  const selectPersona = (persona) => {
    selectedPersona.value = persona
  }

  const clearError = () => {
    error.value = null
  }

  return {
    // State
    personas,
    selectedPersona,
    isLoadingPersonas,
    error,

    // Computed
    hasPersonas,
    selectedPersonaId,

    // Actions
    fetchPersonas,
    fetchPersonaById,
    createPersona,
    updatePersona,
    deletePersona,
    selectPersona,
    clearError
  }
})
