import { defineStore } from 'pinia'
import { chatService } from '@/services/chatService'
import { fileService } from '@/services/fileService'

export const useChatStore = defineStore('chat', {
  state: () => ({
    messages: [],
    sessionId: null, // เพิ่ม session tracking
    uploadedFiles: [], // ไฟล์ที่อัปโหลดใน session ปัจจุบัน
    isLoading: false,
    isStreaming: false,
    streamingContent: '',
    currentPersonaId: 1,
    chatMode: 'text', // 'text' or 'speech'
    systemPrompt: '',
    pagination: {
      total: 0,
      limit: 50,
      offset: 0
    }
  }),

  getters: {
    userMessages: (state) => state.messages.filter(m => m.role === 'user'),
    assistantMessages: (state) => state.messages.filter(m => m.role === 'assistant'),
    totalTokensUsed: (state) => {
      return state.messages
        .filter(m => m.tokens_used)
        .reduce((sum, m) => sum + m.tokens_used, 0)
    }
  },

  actions: {
    // สร้าง session ใหม่
    createNewSession() {
      this.sessionId = `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
      this.messages = []
      this.uploadedFiles = []
      localStorage.setItem('current_session_id', this.sessionId)
      console.log('Created new session:', this.sessionId)
      return this.sessionId
    },

    // Initialize chat session
    initializeChat() {
      // โหลด session_id จาก localStorage (ถ้ามี)
      const savedSessionId = localStorage.getItem('current_session_id')
      if (savedSessionId) {
        this.sessionId = savedSessionId
        console.log('Loaded saved session:', this.sessionId)
      } else {
        this.createNewSession()
      }

      // โหลดประวัติ (ถ้ามี)
      this.fetchHistory()
    },

    async fetchHistory(limit = 50, offset = 0) {
      this.isLoading = true
      try {
        const response = await chatService.getChatHistory(this.sessionId, limit, offset)
        if (response.messages && response.messages.length > 0) {
          this.messages = response.messages.reverse() // Newest at bottom
        }
        this.pagination = {
          total: response.total || 0,
          limit: response.limit || limit,
          offset: response.offset || offset
        }
      } catch (error) {
        console.error('Failed to fetch chat history:', error)
      } finally {
        this.isLoading = false
      }
    },

    // Upload ไฟล์
    async uploadFiles(files, analysisType = 'summary') {
      const uploadPromises = files.map(async (file) => {
        try {
          const result = await fileService.analyzeFile(file, {
            sessionId: this.sessionId,
            analysisType,
            language: 'th'
          })

          // Response ตอนนี้เป็นรูปแบบเดียวกับ ChatResponse
          // มี file_info ซ้อนอยู่ข้างใน
          const fileInfo = result.file_info || {}

          this.uploadedFiles.push({
            fileId: fileInfo.file_id || result.message_id,
            filename: fileInfo.filename,
            fileType: fileInfo.file_type,
            fileSize: fileInfo.file_size
          })

          console.log('File uploaded:', fileInfo.filename)
          return result
        } catch (error) {
          console.error('Failed to upload file:', file.name, error)
          throw error
        }
      })

      await Promise.all(uploadPromises)
    },

    addMessage(message) {
      this.messages.push(message)
    },

    updateStreamingContent(chunk) {
      this.streamingContent += chunk
    },

    startStreaming() {
      this.isStreaming = true
      this.streamingContent = ''
    },

    finishStreaming(data) {
      const assistantMessage = {
        id: data.messageId,
        role: 'assistant',
        content: this.streamingContent,
        persona_id: this.currentPersonaId,
        tokens_used: data.tokensUsed,
        created_at: new Date().toISOString()
      }

      this.addMessage(assistantMessage)
      this.isStreaming = false
      this.streamingContent = ''
    },

    setCurrentPersona(personaId) {
      this.currentPersonaId = personaId
    },

    setChatMode(mode) {
      this.chatMode = mode
    },

    setSystemPrompt(prompt) {
      this.systemPrompt = prompt
    },

    clearChat() {
      this.messages = []
      this.streamingContent = ''
      this.isStreaming = false
      this.uploadedFiles = []
      // สร้าง session ใหม่เมื่อ clear chat
      this.createNewSession()
    },

    // Remove a specific uploaded file
    removeUploadedFile(fileId) {
      const index = this.uploadedFiles.findIndex(f => f.fileId === fileId)
      if (index !== -1) {
        this.uploadedFiles.splice(index, 1)
        console.log('Removed file:', fileId)
      }
    },

    // Clear all uploaded files
    clearUploadedFiles() {
      this.uploadedFiles = []
      console.log('Cleared all uploaded files')
    }
  }
})
