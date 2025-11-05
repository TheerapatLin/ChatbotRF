import { defineStore } from 'pinia'
import { ref } from 'vue'
import { audioService } from '@/api/audioService'

export const useAudioStore = defineStore('audio', () => {
  // State
  const isRecording = ref(false)
  const isProcessing = ref(false)
  const isSpeaking = ref(false)
  const transcript = ref('')
  const audioMode = ref('text') // 'text' or 'speech'
  const mediaRecorder = ref(null)
  const audioChunks = ref([])
  const currentAudio = ref(null)
  const error = ref(null)

  // Voice settings
  const voiceSettings = ref({
    // OpenAI TTS settings
    openai: {
      voice: 'nova', // alloy, echo, fable, onyx, nova, shimmer
      model: 'tts-1', // tts-1, tts-1-hd
      speed: 1.0
    },
    // ElevenLabs TTS settings
    elevenlabs: {
      voice_id: '21m00Tcm4TlvDq8ikWAM', // Rachel
      model_id: 'eleven_multilingual_v2',
      stability: 0.5,
      similarity_boost: 0.75,
      style: 0.0,
      speed: 1.0,
      use_speaker_boost: true
    },
    provider: 'elevenlabs' // 'openai' or 'elevenlabs'
  })

  // Actions
  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      mediaRecorder.value = new MediaRecorder(stream)
      audioChunks.value = []

      mediaRecorder.value.ondataavailable = (event) => {
        audioChunks.value.push(event.data)
      }

      mediaRecorder.value.onstop = async () => {
        const audioBlob = new Blob(audioChunks.value, { type: 'audio/webm' })
        await transcribeAudio(audioBlob)
      }

      mediaRecorder.value.start()
      isRecording.value = true
      error.value = null
    } catch (err) {
      console.error('Failed to start recording:', err)
      error.value = 'ไม่สามารถเข้าถึงไมโครโฟนได้ กรุณาอนุญาตการใช้ไมโครโฟน'
      throw err
    }
  }

  const stopRecording = () => {
    if (mediaRecorder.value && isRecording.value) {
      mediaRecorder.value.stop()
      mediaRecorder.value.stream.getTracks().forEach((track) => track.stop())
      isRecording.value = false
    }
  }

  const transcribeAudio = async (audioBlob) => {
    try {
      isProcessing.value = true
      error.value = null

      const formData = new FormData()
      formData.append('file', audioBlob, 'audio.webm')

      const response = await audioService.transcribe(formData)
      transcript.value = response.data.text || ''

      return transcript.value
    } catch (err) {
      console.error('Failed to transcribe audio:', err)
      error.value = 'ไม่สามารถถอดเสียงได้ กรุณาลองใหม่อีกครั้ง'
      throw err
    } finally {
      isProcessing.value = false
    }
  }

  const textToSpeech = async (text, customSettings = null) => {
    try {
      isSpeaking.value = true
      error.value = null

      let audioBlob

      const provider = customSettings?.provider || voiceSettings.value.provider

      if (provider === 'elevenlabs') {
        // Use ElevenLabs TTS
        const settings = customSettings?.elevenlabs || voiceSettings.value.elevenlabs
        const response = await audioService.elevenLabsTTS({
          text: text,
          voice_id: settings.voice_id,
          model_id: settings.model_id,
          stability: settings.stability,
          similarity_boost: settings.similarity_boost,
          style: settings.style,
          speed: settings.speed,
          use_speaker_boost: settings.use_speaker_boost
        })
        audioBlob = response.data
      } else {
        // Use OpenAI TTS
        const settings = customSettings?.openai || voiceSettings.value.openai
        const response = await audioService.textToSpeech({
          text: text,
          voice: settings.voice,
          model: settings.model,
          speed: settings.speed
        })
        audioBlob = response.data
      }

      // Play audio
      await playAudio(audioBlob)

      return audioBlob
    } catch (err) {
      console.error('Failed to convert text to speech:', err)
      error.value = 'ไม่สามารถแปลงข้อความเป็นเสียงได้'
      throw err
    } finally {
      isSpeaking.value = false
    }
  }

  const playAudio = (audioBlob) => {
    return new Promise((resolve, reject) => {
      try {
        // Stop current audio if playing
        stopAudioPlayback()

        const audioUrl = URL.createObjectURL(audioBlob)
        currentAudio.value = new Audio(audioUrl)

        currentAudio.value.onended = () => {
          isSpeaking.value = false
          URL.revokeObjectURL(audioUrl)
          resolve()
        }

        currentAudio.value.onerror = (err) => {
          isSpeaking.value = false
          URL.revokeObjectURL(audioUrl)
          reject(err)
        }

        currentAudio.value.play()
      } catch (err) {
        console.error('Failed to play audio:', err)
        reject(err)
      }
    })
  }

  const stopAudioPlayback = () => {
    if (currentAudio.value) {
      currentAudio.value.pause()
      currentAudio.value.currentTime = 0
      currentAudio.value = null
      isSpeaking.value = false
    }
  }

  const toggleAudioMode = () => {
    audioMode.value = audioMode.value === 'text' ? 'speech' : 'text'
  }

  const setAudioMode = (mode) => {
    if (mode === 'text' || mode === 'speech') {
      audioMode.value = mode
    }
  }

  const updateVoiceSettings = (provider, settings) => {
    if (provider === 'openai') {
      voiceSettings.value.openai = { ...voiceSettings.value.openai, ...settings }
    } else if (provider === 'elevenlabs') {
      voiceSettings.value.elevenlabs = { ...voiceSettings.value.elevenlabs, ...settings }
    }
    voiceSettings.value.provider = provider
  }

  const clearTranscript = () => {
    transcript.value = ''
  }

  const clearError = () => {
    error.value = null
  }

  return {
    // State
    isRecording,
    isProcessing,
    isSpeaking,
    transcript,
    audioMode,
    voiceSettings,
    error,

    // Actions
    startRecording,
    stopRecording,
    transcribeAudio,
    textToSpeech,
    playAudio,
    stopAudioPlayback,
    toggleAudioMode,
    setAudioMode,
    updateVoiceSettings,
    clearTranscript,
    clearError
  }
})
