import apiClient from './axios'

export const audioService = {
  /**
   * Transcribe audio (Speech-to-Text)
   * @param {FormData} audioFormData - Form data with audio file
   */
  transcribe: (audioFormData) =>
    apiClient.post('/api/audio/transcribe', audioFormData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    }),

  /**
   * Text-to-Speech (OpenAI)
   * @param {Object} data - TTS data
   * @param {string} data.text - Text to convert
   * @param {string} data.voice - Voice name (alloy, echo, fable, onyx, nova, shimmer)
   * @param {string} data.model - Model (tts-1, tts-1-hd)
   * @param {string} data.response_format - Format (mp3, opus, aac, flac, wav, pcm)
   * @param {number} data.speed - Speed (0.25-4.0)
   */
  textToSpeech: (data) =>
    apiClient.post('/api/audio/tts', data, {
      responseType: 'blob'
    }),

  /**
   * Text-to-Speech (ElevenLabs)
   * @param {Object} data - TTS data
   * @param {string} data.text - Text to convert
   * @param {string} data.voice_id - ElevenLabs voice ID
   * @param {string} data.model_id - Model ID
   * @param {number} data.stability - Stability (0.0-1.0)
   * @param {number} data.similarity_boost - Similarity boost (0.0-1.0)
   * @param {number} data.style - Style (0.0-1.0)
   * @param {number} data.speed - Speed (0.7-1.2)
   * @param {boolean} data.use_speaker_boost - Use speaker boost
   */
  elevenLabsTTS: (data) =>
    apiClient.post('/audio/elevenlabs/tts', data, {
      headers: { Accept: 'audio/mpeg' },
      responseType: 'blob'
    }),

  /**
   * Get ElevenLabs voices
   */
  getElevenLabsVoices: () => apiClient.get('/audio/elevenlabs/voices')
}
