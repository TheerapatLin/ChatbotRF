/**
 * Utility functions for formatting data
 */

/**
 * Format file size from bytes to human-readable format
 * @param {number} bytes - File size in bytes
 * @returns {string} Formatted file size
 */
export function formatFileSize(bytes) {
  if (!bytes || bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

/**
 * Format timestamp to localized time string
 * @param {string|number|Date} timestamp - Timestamp to format
 * @param {string} locale - Locale for formatting (default: 'th-TH')
 * @returns {string} Formatted time string
 */
export function formatTime(timestamp, locale = 'th-TH') {
  if (!timestamp) return ''
  return new Date(timestamp).toLocaleTimeString(locale, {
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * Format timestamp to localized date and time string
 * @param {string|number|Date} timestamp - Timestamp to format
 * @param {string} locale - Locale for formatting (default: 'th-TH')
 * @returns {string} Formatted date and time string
 */
export function formatDateTime(timestamp, locale = 'th-TH') {
  if (!timestamp) return ''
  return new Date(timestamp).toLocaleString(locale, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
