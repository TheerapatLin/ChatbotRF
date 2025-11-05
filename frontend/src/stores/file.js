import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fileService } from '@/api/fileService'

export const useFileStore = defineStore('file', () => {
  // State
  const uploadedFiles = ref([])
  const isUploading = ref(false)
  const uploadProgress = ref(0)
  const error = ref(null)

  // Actions
  const uploadFiles = async (files, options = {}) => {
    try {
      isUploading.value = true
      error.value = null
      uploadProgress.value = 0

      const formData = new FormData()

      // Add files
      files.forEach((file) => {
        formData.append('files', file)
      })

      // Add optional fields
      if (options.prompt) {
        formData.append('prompt', options.prompt)
      }
      if (options.session_id) {
        formData.append('session_id', options.session_id)
      }
      if (options.use_history !== undefined) {
        formData.append('use_history', options.use_history)
      }
      if (options.system_prompt) {
        formData.append('system_prompt', options.system_prompt)
      }

      const response = await fileService.uploadFiles(formData)

      // Store uploaded files info
      const uploadedFilesData = response.data.uploaded_files || []
      uploadedFiles.value.push(...uploadedFilesData)

      uploadProgress.value = 100

      return {
        success: response.data.success,
        failed: response.data.failed,
        total: response.data.total,
        file_ids: uploadedFilesData.map((f) => f.file_id),
        files: uploadedFilesData
      }
    } catch (err) {
      console.error('Failed to upload files:', err)
      error.value = err.response?.data?.error || 'Failed to upload files'
      throw err
    } finally {
      isUploading.value = false
    }
  }

  const getFileHistory = async (params = {}) => {
    try {
      const response = await fileService.getFileHistory(params)
      return response.data
    } catch (err) {
      console.error('Failed to fetch file history:', err)
      error.value = err.response?.data?.error || 'Failed to fetch file history'
      throw err
    }
  }

  const deleteAllFiles = async () => {
    try {
      await fileService.deleteAllFiles()
      uploadedFiles.value = []
    } catch (err) {
      console.error('Failed to delete files:', err)
      error.value = err.response?.data?.error || 'Failed to delete files'
      throw err
    }
  }

  const clearUploadedFiles = () => {
    uploadedFiles.value = []
  }

  const removeFile = (fileId) => {
    uploadedFiles.value = uploadedFiles.value.filter((f) => f.file_id !== fileId)
  }

  const clearError = () => {
    error.value = null
  }

  return {
    // State
    uploadedFiles,
    isUploading,
    uploadProgress,
    error,

    // Actions
    uploadFiles,
    getFileHistory,
    deleteAllFiles,
    clearUploadedFiles,
    removeFile,
    clearError
  }
})
