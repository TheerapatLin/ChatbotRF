import api from './api'

/**
 * Audio Service
 * Handles audio transcription and text-to-speech API calls
 */

/**
 * Transcribe audio to text using OpenAI Whisper API
 * @param {Blob} audioBlob - Audio file blob
 * @param {string} filename - Original filename
 * @returns {Promise<Object>} Transcription result
 */
export async function transcribeAudio(audioBlob, filename = 'recording.webm') {
  const formData = new FormData()
  formData.append('audio', audioBlob, filename)

  try {
    const response = await api.post('/audio/transcribe', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    return response
  } catch (error) {
    console.error('Transcription error:', error)
    throw new Error('ไม่สามารถแปลงเสียงเป็นข้อความได้: ' + error.message)
  }
}

/**
 * Convert text to speech using OpenAI TTS API
 * @param {string} text - Text to convert
 * @param {Object} options - TTS options
 * @returns {Promise<Object>} TTS result with audio data
 */
export async function textToSpeech(text, options = {}) {
  const {
    voice = 'nova',
    model = 'tts-1',
    responseFormat = 'mp3',
    speed = 1.0
  } = options

  try {
    const response = await api.post('/audio/tts', {
      text,
      voice,
      model,
      response_format: responseFormat,
      speed
    })

    return response
  } catch (error) {
    console.error('TTS error:', error)
    throw new Error('ไม่สามารถแปลงข้อความเป็นเสียงได้: ' + error.message)
  }
}

/**
 * Convert base64 audio to Blob
 * @param {string} base64Audio - Base64 encoded audio
 * @param {string} mimeType - MIME type of audio
 * @returns {Blob} Audio blob
 */
export function base64ToBlob(base64Audio, mimeType = 'audio/mpeg') {
  const binaryString = atob(base64Audio)
  const bytes = new Uint8Array(binaryString.length)

  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i)
  }

  return new Blob([bytes], { type: mimeType })
}

/**
 * Create object URL from audio blob
 * @param {Blob} blob - Audio blob
 * @returns {string} Object URL
 */
export function createAudioURL(blob) {
  return URL.createObjectURL(blob)
}

/**
 * Revoke object URL to free memory
 * @param {string} url - Object URL to revoke
 */
export function revokeAudioURL(url) {
  if (url) {
    URL.revokeObjectURL(url)
  }
}

export const audioService = {
  transcribeAudio,
  textToSpeech,
  base64ToBlob,
  createAudioURL,
  revokeAudioURL
}

export default audioService
