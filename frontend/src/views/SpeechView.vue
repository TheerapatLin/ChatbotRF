<template>
  <div class="speech-view">
    <!-- Header -->
    <header class="speech-header">
      <div class="flex items-center justify-between w-full max-w-4xl mx-auto">
        <div>
          <h1 class="text-2xl font-bold text-gray-800">Speech to Speech</h1>
          <p class="text-sm text-gray-500">พูดเพื่อสนทนากับ AI</p>
        </div>

        <div class="flex items-center gap-3">
          <!-- Connection Status -->
          <div class="flex items-center gap-2 px-3 py-2 bg-gray-100 rounded-lg">
            <div
              :class="[
                'w-2 h-2 rounded-full',
                isConnected ? 'bg-green-500' : 'bg-red-500'
              ]"
            ></div>
            <span class="text-sm text-gray-600">
              {{ isConnected ? 'เชื่อมต่อแล้ว' : 'ไม่ได้เชื่อมต่อ' }}
            </span>
          </div>

          <!-- Reset Button -->
          <button
            v-if="!isIdle"
            @click="handleReset"
            class="px-4 py-2 bg-gray-500 hover:bg-gray-600 text-white rounded-lg transition-all text-sm font-medium"
          >
            รีเซ็ต
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <div class="speech-content">
      <div class="content-wrapper">
        <!-- Error Message -->
        <div v-if="hasError" class="error-alert">
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path
              fill-rule="evenodd"
              d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
              clip-rule="evenodd"
            />
          </svg>
          <span>{{ error }}</span>
          <button @click="clearError" class="error-close">
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path
                fill-rule="evenodd"
                d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                clip-rule="evenodd"
              />
            </svg>
          </button>
        </div>

        <!-- Waveform Visualizer -->
        <WaveformVisualizer
          :audio-source="audioSource"
          :mode="visualizerMode"
          :width="600"
          :height="200"
        />

        <!-- Microphone Button -->
        <MicrophoneButton
          :state="state"
          :disabled="!isConnected"
          @start-recording="handleStartRecording"
          @stop-recording="handleStopRecording"
        />

        <!-- Status Messages (for debugging, can be hidden) -->
        <div v-if="false" class="debug-info">
          <div v-if="transcribedText" class="debug-item">
            <strong>Transcribed:</strong> {{ transcribedText }}
          </div>
          <div v-if="aiResponseText" class="debug-item">
            <strong>AI Response:</strong> {{ aiResponseText }}
          </div>
        </div>

        <!-- Instructions -->
        <div class="instructions">
          <div class="instruction-item">
            <div class="instruction-icon">1</div>
            <div class="instruction-text">กดค้างปุ่มไมโครโฟนเพื่อเริ่มพูด</div>
          </div>
          <div class="instruction-item">
            <div class="instruction-icon">2</div>
            <div class="instruction-text">ปล่อยปุ่มเมื่อพูดเสร็จ</div>
          </div>
          <div class="instruction-item">
            <div class="instruction-icon">3</div>
            <div class="instruction-text">รอ AI ตอบกลับเป็นเสียง</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import WaveformVisualizer from '@/components/speech/WaveformVisualizer.vue'
import MicrophoneButton from '@/components/speech/MicrophoneButton.vue'
import { useSpeechToSpeech } from '@/composables/useSpeechToSpeech'

// Inject WebSocket connection from parent
const wsConnection = inject('wsConnection')
const { isConnected } = wsConnection

// Speech-to-Speech composable
const speechToSpeech = useSpeechToSpeech(wsConnection)

const {
  state,
  error,
  transcribedText,
  aiResponseText,
  audioElement,
  mediaStream,
  isIdle,
  isRecording,
  isProcessing,
  isPlaying,
  hasError,
  startRecording,
  stopRecordingAndProcess,
  cancel,
  clearError,
  reset
} = speechToSpeech

// Computed for waveform visualizer
const audioSource = computed(() => {
  // Only use waveform visualization for recording
  // Playback visualization is disabled to avoid Web Audio API connection issues
  if (isRecording.value && mediaStream.value) {
    return mediaStream.value
  }
  return null
})

const visualizerMode = computed(() => {
  if (isRecording.value) return 'recording'
  if (isPlaying.value) return 'playing'
  return 'idle'
})

// Event handlers
const handleStartRecording = async () => {
  try {
    await startRecording()
  } catch (err) {
    console.error('Start recording error:', err)
  }
}

const handleStopRecording = async () => {
  try {
    await stopRecordingAndProcess()
  } catch (err) {
    console.error('Stop recording error:', err)
    // Error is already set in composable
  }
}

const handleReset = () => {
  if (confirm('คุณต้องการรีเซ็ตการสนทนาใช่หรือไม่?')) {
    reset()
  }
}
</script>

<style scoped>
.speech-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.speech-header {
  padding: 1.5rem 2rem;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.speech-content {
  flex: 1;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 3rem 2rem 2rem 2rem;
  overflow: auto;
}

.content-wrapper {
  width: 100%;
  max-width: 800px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2rem;
  margin-top: 1rem;
}

/* Error Alert */
.error-alert {
  width: 100%;
  max-width: 600px;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  background: #fee2e2;
  border: 1px solid #fecaca;
  border-radius: 12px;
  color: #991b1b;
  font-size: 0.875rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.error-alert svg:first-child {
  flex-shrink: 0;
}

.error-close {
  margin-left: auto;
  padding: 0.25rem;
  border: none;
  background: transparent;
  color: #991b1b;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.2s;
}

.error-close:hover {
  background: rgba(0, 0, 0, 0.1);
}

/* Debug Info */
.debug-info {
  width: 100%;
  max-width: 600px;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  backdrop-filter: blur(10px);
}

.debug-item {
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  color: white;
}

.debug-item strong {
  font-weight: 600;
}

/* Instructions */
.instructions {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  width: 100%;
  max-width: 600px;
  padding: 1.5rem;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.2);
}

.instruction-item {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.instruction-icon {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  font-weight: 600;
  font-size: 0.875rem;
}

.instruction-text {
  color: #374151;
  font-size: 0.9375rem;
}

/* Responsive */
@media (max-width: 768px) {
  .speech-header {
    padding: 1rem 1.5rem;
  }

  .speech-header h1 {
    font-size: 1.5rem;
  }

  .speech-content {
    padding: 1rem;
  }

  .instructions {
    padding: 1rem;
  }
}
</style>
