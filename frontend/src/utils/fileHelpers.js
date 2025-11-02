/**
 * Utility functions for file handling
 */

/**
 * Get emoji icon based on file type
 * @param {string} fileType - MIME type of the file
 * @returns {string} Emoji representing the file type
 */
export function getFileTypeIcon(fileType) {
  if (!fileType) return '📄'
  const type = fileType.toLowerCase()

  if (type.includes('pdf')) return '📕'
  if (type.includes('word') || type.includes('docx')) return '📘'
  if (type.includes('excel') || type.includes('xlsx')) return '📗'
  if (type.includes('powerpoint') || type.includes('pptx')) return '📙'
  if (type.includes('image')) return '🖼️'
  if (type.includes('text')) return '📄'
  if (type.includes('video')) return '🎥'
  if (type.includes('audio')) return '🎵'
  if (type.includes('zip') || type.includes('rar') || type.includes('7z')) return '📦'

  return '📄'
}

/**
 * Get Tailwind CSS color class based on file type
 * @param {string} fileType - MIME type of the file
 * @returns {string} Tailwind background color class
 */
export function getFileTypeColor(fileType) {
  if (!fileType) return 'bg-gray-500'
  const type = fileType.toLowerCase()

  if (type.includes('pdf')) return 'bg-red-500'
  if (type.includes('word') || type.includes('docx')) return 'bg-blue-500'
  if (type.includes('excel') || type.includes('xlsx')) return 'bg-green-500'
  if (type.includes('powerpoint') || type.includes('pptx')) return 'bg-orange-500'
  if (type.includes('image')) return 'bg-purple-500'
  if (type.includes('text')) return 'bg-gray-500'
  if (type.includes('video')) return 'bg-pink-500'
  if (type.includes('audio')) return 'bg-indigo-500'
  if (type.includes('zip') || type.includes('rar') || type.includes('7z')) return 'bg-yellow-500'

  return 'bg-gray-500'
}

/**
 * Validate file type against allowed types
 * @param {File} file - File object to validate
 * @param {string[]} allowedTypes - Array of allowed MIME types
 * @returns {boolean} True if file type is allowed
 */
export function isFileTypeAllowed(file, allowedTypes) {
  if (!file || !allowedTypes || allowedTypes.length === 0) return false
  return allowedTypes.some(type => file.type.includes(type))
}

/**
 * Check if file size is within limit
 * @param {File} file - File object to check
 * @param {number} maxSize - Maximum file size in bytes
 * @returns {boolean} True if file size is within limit
 */
export function isFileSizeValid(file, maxSize) {
  if (!file || !maxSize) return false
  return file.size <= maxSize
}
