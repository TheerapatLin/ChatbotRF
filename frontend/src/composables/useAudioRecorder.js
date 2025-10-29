import { ref } from 'vue'

export function useAudioRecorder() {
  const isRecording = ref(false)
  const mediaRecorder = ref(null)
  const audioChunks = ref([])
  const stream = ref(null)

  async function startRecording() {
    try {
      stream.value = await navigator.mediaDevices.getUserMedia({
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          sampleRate: 44100
        }
      })

      mediaRecorder.value = new MediaRecorder(stream.value, {
        mimeType: 'audio/webm'
      })

      audioChunks.value = []

      mediaRecorder.value.ondataavailable = (event) => {
        if (event.data.size > 0) {
          audioChunks.value.push(event.data)
        }
      }

      mediaRecorder.value.start()
      isRecording.value = true

      console.log('Recording started')
    } catch (error) {
      console.error('Failed to start recording:', error)
      throw error
    }
  }

  async function stopRecording() {
    return new Promise((resolve, reject) => {
      if (!mediaRecorder.value || !isRecording.value) {
        reject(new Error('No active recording'))
        return
      }

      mediaRecorder.value.onstop = () => {
        const audioBlob = new Blob(audioChunks.value, { type: 'audio/webm' })

        if (stream.value) {
          stream.value.getTracks().forEach(track => track.stop())
        }

        isRecording.value = false
        console.log('Recording stopped', audioBlob.size, 'bytes')

        resolve(audioBlob)
      }

      mediaRecorder.value.onerror = (error) => {
        reject(error)
      }

      mediaRecorder.value.stop()
    })
  }

  function cancelRecording() {
    if (mediaRecorder.value && isRecording.value) {
      mediaRecorder.value.stop()

      if (stream.value) {
        stream.value.getTracks().forEach(track => track.stop())
      }

      audioChunks.value = []
      isRecording.value = false

      console.log('Recording cancelled')
    }
  }

  return {
    isRecording,
    startRecording,
    stopRecording,
    cancelRecording
  }
}
