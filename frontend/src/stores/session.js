import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { sessionService } from '@/api/sessionService'
import { useChatStore } from './chat'

export const useSessionStore = defineStore('sessions', () => {
  const chatStore = useChatStore()

  // State
  const sessions = ref([])
  const isLoading = ref(false)
  const sessionTitles = ref({}) // เก็บ custom titles ใน memory

  // Computed
  const sortedSessions = computed(() => sessions.value)

  // Load custom titles from localStorage
  const loadCustomTitles = () => {
    try {
      const saved = localStorage.getItem('session_titles')
      if (saved) {
        sessionTitles.value = JSON.parse(saved)
      }
    } catch (error) {
      console.error('Failed to load custom titles:', error)
      sessionTitles.value = {}
    }
  }

  // Save custom titles to localStorage
  const saveCustomTitles = () => {
    try {
      localStorage.setItem('session_titles', JSON.stringify(sessionTitles.value))
    } catch (error) {
      console.error('Failed to save custom titles:', error)
    }
  }

  // Get title for session (custom or auto-generated)
  const getSessionTitle = (session) => {
    return sessionTitles.value[session.session_id] || session.title
  }

  // Actions
  const loadSessions = async () => {
    try {
      isLoading.value = true
      loadCustomTitles()
      sessions.value = await sessionService.getSessionsList(50)
    } catch (error) {
      console.error('Failed to load sessions:', error)
      sessions.value = []
    } finally {
      isLoading.value = false
    }
  }

  const switchToSession = async (sessionId) => {
    try {
      const response = await sessionService.getSessionMessages(sessionId)
      chatStore.sessionId = sessionId
      chatStore.messages = response.data.messages || []
    } catch (error) {
      console.error('Failed to switch session:', error)
      throw error
    }
  }

  const createNewChat = () => {
    chatStore.newChat()
    // Reload sessions list after user sends first message
  }

  const renameSession = (sessionId, newTitle) => {
    sessionTitles.value[sessionId] = newTitle
    saveCustomTitles()

    // Update in sessions array
    const session = sessions.value.find(s => s.session_id === sessionId)
    if (session) {
      session.title = newTitle
    }
  }

  const deleteSession = async (sessionId) => {
    try {
      await sessionService.deleteSession(sessionId)
      sessions.value = sessions.value.filter(s => s.session_id !== sessionId)

      // Remove custom title
      delete sessionTitles.value[sessionId]
      saveCustomTitles()

      // If deleted active session, create new
      if (chatStore.sessionId === sessionId) {
        createNewChat()
      }
    } catch (error) {
      console.error('Failed to delete session:', error)
      throw error
    }
  }

  return {
    sessions,
    sortedSessions,
    isLoading,
    loadSessions,
    switchToSession,
    createNewChat,
    renameSession,
    deleteSession,
    getSessionTitle
  }
})
