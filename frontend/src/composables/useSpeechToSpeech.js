import { ref, computed } from 'vue'
import { useAudioRecorder } from './useAudioRecorder'
import { audioService } from '@/services/audioService'

/**
 * Speech-to-Speech Composable
 * Orchestrates the full workflow: Record → Transcribe → Chat → TTS → Play
 */
export function useSpeechToSpeech(wsConnection) {
  // State
  const state = ref('idle') // idle | recording | processing | playing
  const error = ref(null)
  const transcribedText = ref('')
  const aiResponseText = ref('')
  const audioElement = ref(null)
  const currentAudioURL = ref(null)

  // Audio recorder
  const audioRecorder = useAudioRecorder()
  // Use the stream from audioRecorder for waveform visualization
  const mediaStream = audioRecorder.stream

  // Computed
  const isIdle = computed(() => state.value === 'idle')
  const isRecording = computed(() => state.value === 'recording')
  const isProcessing = computed(() => state.value === 'processing')
  const isPlaying = computed(() => state.value === 'playing')
  const hasError = computed(() => error.value !== null)

  /**
   * Start recording audio
   */
  const startRecording = async () => {
    try {
      error.value = null
      state.value = 'recording'

      // Start recording - this will also set the mediaStream
      await audioRecorder.startRecording()

      console.log('Recording started, stream available:', !!mediaStream.value)
    } catch (err) {
      console.error('Failed to start recording:', err)
      error.value = 'ไม่สามารถเข้าถึงไมโครโฟนได้ กรุณาอนุญาตการใช้งานไมโครโฟน'
      state.value = 'idle'
      throw err
    }
  }

  /**
   * Stop recording and process audio through the pipeline
   */
  const stopRecordingAndProcess = async () => {
    try {
      state.value = 'processing'
      console.log('Stopping recording and processing...')

      // Stop recording and get audio blob
      // (stopRecording also stops the media stream)
      const audioBlob = await audioRecorder.stopRecording()
      console.log('Audio blob created:', audioBlob.size, 'bytes')

      // Step 1: Transcribe audio to text
      console.log('Step 1: Transcribing audio...')
      const transcribeResult = await audioService.transcribeAudio(audioBlob)
      transcribedText.value = transcribeResult.text
      console.log('Transcribed text:', transcribedText.value)

      if (!transcribedText.value) {
        throw new Error('ไม่พบข้อความจากเสียงที่บันทึก')
      }

      // Step 2: Send to chat API via WebSocket and wait for response
      console.log('Step 2: Sending to chat API...')
      const aiResponse = await sendToChatAndWaitForResponse(transcribedText.value)
      aiResponseText.value = aiResponse
      console.log('AI response:', aiResponseText.value)

      if (!aiResponseText.value) {
        throw new Error('ไม่ได้รับการตอบกลับจาก AI')
      }

      // Step 3: Convert AI response to speech
      console.log('Step 3: Converting to speech...')
      const ttsResult = await audioService.textToSpeech(aiResponseText.value, {
        voice: 'nova',
        model: 'tts-1',
        responseFormat: 'mp3',
        speed: 1.0
      })
      console.log('TTS result:', ttsResult)

      // Step 4: Play audio response
      console.log('Step 4: Playing audio...')
      await playAudioResponse(ttsResult.audio_data)

    } catch (err) {
      console.error('Processing error:', err)
      error.value = err.message || 'เกิดข้อผิดพลาดในการประมวลผล'
      state.value = 'idle'

      // Clean up media stream on error
      if (mediaStream.value) {
        mediaStream.value.getTracks().forEach(track => track.stop())
        mediaStream.value = null
      }

      throw err
    }
  }

  /**
   * Send message to chat API via WebSocket and wait for full response
   */
  const sendToChatAndWaitForResponse = (message) => {
    return new Promise((resolve, reject) => {
      // Validate WebSocket connection
      if (!wsConnection) {
        reject(new Error('WebSocket connection object not available'))
        return
      }

      if (!wsConnection.ws || !wsConnection.ws.value) {
        reject(new Error('WebSocket not initialized'))
        return
      }

      if (wsConnection.ws.value.readyState !== WebSocket.OPEN) {
        reject(new Error('WebSocket is not connected (state: ' + wsConnection.ws.value.readyState + ')'))
        return
      }

      console.log('WebSocket is connected, sending message...')

      let fullResponse = ''
      let messageHandler = null

      // Create temporary message handler
      messageHandler = (event) => {
        try {
          const data = JSON.parse(event.data)

          if (data.type === 'chunk') {
            if (data.done) {
              // Streaming complete
              console.log('Chat streaming complete. Full response:', fullResponse)

              // Remove listener
              wsConnection.ws.value.removeEventListener('message', messageHandler)

              resolve(fullResponse)
            } else {
              // Accumulate chunks
              fullResponse += data.content
              console.log('Received chunk:', data.content)
            }
          } else if (data.type === 'error') {
            // Remove listener
            wsConnection.ws.value.removeEventListener('message', messageHandler)

            reject(new Error(data.error || 'Chat API error'))
          }
        } catch (err) {
          console.error('Message handler error:', err)
          wsConnection.ws.value.removeEventListener('message', messageHandler)
          reject(err)
        }
      }

      // Add message listener
      wsConnection.ws.value.addEventListener('message', messageHandler)

      // Send message with system prompt for speech mode
      const success = wsConnection.sendMessage({
        content: message,
        persona_id: 1,
        system_prompt: 'คุณกำลังพูดอยู่กับผู้ใช้งาน ตอบกลับด้วยประโยคสั้นๆ และกระชับ'
      })

      if (!success) {
        wsConnection.ws.value.removeEventListener('message', messageHandler)
        reject(new Error('Failed to send message to WebSocket'))
        return
      }

      // Set timeout (30 seconds)
      setTimeout(() => {
        if (wsConnection.ws && wsConnection.ws.value) {
          wsConnection.ws.value.removeEventListener('message', messageHandler)
        }
        reject(new Error('Chat response timeout (30s)'))
      }, 30000)
    })
  }

  /**
   * Play audio response and update state
   */
  const playAudioResponse = async (base64Audio) => {
    return new Promise((resolve, reject) => {
      try {
        // Clean up previous audio
        if (currentAudioURL.value) {
          audioService.revokeAudioURL(currentAudioURL.value)
          currentAudioURL.value = null
        }

        // Clean up previous audio element
        if (audioElement.value) {
          audioElement.value.pause()
          audioElement.value.src = ''
          audioElement.value.onended = null
          audioElement.value.onerror = null
          audioElement.value.oncanplaythrough = null
          audioElement.value = null
        }

        // Convert base64 to blob
        const audioBlob = audioService.base64ToBlob(base64Audio, 'audio/mpeg')
        currentAudioURL.value = audioService.createAudioURL(audioBlob)

        // Create NEW audio element
        audioElement.value = new Audio()

        // Setup event handlers BEFORE setting src
        audioElement.value.oncanplaythrough = () => {
          console.log('Audio can play through, starting playback...')

          // Set state to playing
          state.value = 'playing'

          // Start playback
          audioElement.value.play()
            .then(() => {
              console.log('Audio play() promise resolved')
            })
            .catch(err => {
              console.error('Play error:', err)
              state.value = 'idle'
              reject(new Error('ไม่สามารถเล่นเสียงได้: ' + err.message))
            })
        }

        audioElement.value.onended = () => {
          console.log('Audio playback ended')
          state.value = 'idle'

          // Clean up
          if (currentAudioURL.value) {
            audioService.revokeAudioURL(currentAudioURL.value)
            currentAudioURL.value = null
          }

          resolve()
        }

        audioElement.value.onerror = (err) => {
          // Only reject if there's a real error (check audioElement.value.error)
          if (audioElement.value.error) {
            console.error('Audio playback error:', err, audioElement.value.error)
            state.value = 'idle'
            reject(new Error('ไม่สามารถเล่นเสียงได้: ' + audioElement.value.error.message))
          } else {
            // False alarm - ignore error event without actual error
            console.warn('Audio error event without error object, ignoring')
          }
        }

        // Set src AFTER event handlers are installed
        audioElement.value.src = currentAudioURL.value
        audioElement.value.load() // Explicitly start loading

        console.log('Audio element created and loading:', {
          src: currentAudioURL.value.substring(0, 50) + '...',
          readyState: audioElement.value.readyState
        })

      } catch (err) {
        console.error('Play audio error:', err)
        state.value = 'idle'
        reject(err)
      }
    })
  }

  /**
   * Cancel current operation
   */
  const cancel = () => {
    if (isRecording.value) {
      audioRecorder.cancelRecording()

      if (mediaStream.value) {
        mediaStream.value.getTracks().forEach(track => track.stop())
        mediaStream.value = null
      }
    }

    if (isPlaying.value && audioElement.value) {
      audioElement.value.pause()
      audioElement.value.currentTime = 0
    }

    if (currentAudioURL.value) {
      audioService.revokeAudioURL(currentAudioURL.value)
      currentAudioURL.value = null
    }

    state.value = 'idle'
    error.value = null
  }

  /**
   * Clear error
   */
  const clearError = () => {
    error.value = null
  }

  /**
   * Reset to initial state
   */
  const reset = () => {
    cancel()
    transcribedText.value = ''
    aiResponseText.value = ''
  }

  return {
    // State
    state,
    error,
    transcribedText,
    aiResponseText,
    audioElement,
    mediaStream,

    // Computed
    isIdle,
    isRecording,
    isProcessing,
    isPlaying,
    hasError,

    // Methods
    startRecording,
    stopRecordingAndProcess,
    cancel,
    clearError,
    reset
  }
}
