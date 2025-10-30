<template>
  <div class="waveform-container">
    <canvas
      ref="canvasRef"
      :width="width"
      :height="height"
      class="waveform-canvas"
    ></canvas>

    <!-- Status indicator -->
    <div v-if="mode !== 'idle'" class="status-indicator">
      <div :class="['status-dot', `status-${mode}`]"></div>
      <span class="status-text">
        {{ statusText }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'

const props = defineProps({
  // Audio source: MediaStream (recording) or HTMLAudioElement (playback)
  audioSource: {
    type: [MediaStream, HTMLAudioElement],
    default: null
  },

  // Visualization mode
  mode: {
    type: String,
    default: 'idle', // 'idle' | 'recording' | 'playing'
    validator: (value) => ['idle', 'recording', 'playing'].includes(value)
  },

  // Canvas dimensions
  width: {
    type: Number,
    default: 600
  },

  height: {
    type: Number,
    default: 200
  },

  // Color scheme
  barColor: {
    type: String,
    default: null // Will use mode-based colors if null
  }
})

const canvasRef = ref(null)
const audioContext = ref(null)
const analyser = ref(null)
const dataArray = ref(null)
const animationId = ref(null)
const sourceNode = ref(null)
const connectedAudioElement = ref(null) // Track which audio element is connected

const statusText = computed(() => {
  switch (props.mode) {
    case 'recording':
      return 'กำลังบันทึกเสียง...'
    case 'playing':
      return 'กำลังเล่นเสียง...'
    default:
      return ''
  }
})

// Initialize audio context and analyser
const initAudioContext = () => {
  try {
    audioContext.value = new (window.AudioContext || window.webkitAudioContext)()
    analyser.value = audioContext.value.createAnalyser()
    analyser.value.fftSize = 256

    const bufferLength = analyser.value.frequencyBinCount
    dataArray.value = new Uint8Array(bufferLength)

    console.log('Audio context initialized', { bufferLength })
  } catch (error) {
    console.error('Failed to initialize audio context:', error)
  }
}

// Connect audio source to analyser
const connectAudioSource = () => {
  if (!audioContext.value || !analyser.value || !props.audioSource) {
    return
  }

  try {
    // Disconnect previous source if exists
    if (sourceNode.value) {
      sourceNode.value.disconnect()
      sourceNode.value = null
    }

    // Connect based on source type
    if (props.audioSource instanceof MediaStream) {
      // For recording: MediaStream
      sourceNode.value = audioContext.value.createMediaStreamSource(props.audioSource)
      sourceNode.value.connect(analyser.value)
      console.log('MediaStream connected to analyser')
    } else if (props.audioSource instanceof HTMLAudioElement) {
      // For playback: HTMLAudioElement
      // Check if this audio element was already connected
      if (connectedAudioElement.value === props.audioSource && sourceNode.value) {
        console.log('HTMLAudioElement already connected, reusing existing connection')
        // Make sure it's still connected properly
        if (!sourceNode.value.context || sourceNode.value.context.state === 'closed') {
          console.log('Previous connection was closed, reconnecting...')
          connectedAudioElement.value = null
        } else {
          return
        }
      }

      // Try to create new media element source
      try {
        sourceNode.value = audioContext.value.createMediaElementSource(props.audioSource)
        sourceNode.value.connect(analyser.value)
        analyser.value.connect(audioContext.value.destination)
        connectedAudioElement.value = props.audioSource
        console.log('HTMLAudioElement connected to analyser and destination')
      } catch (err) {
        // If element is already connected to another source, we can't visualize it
        // This is a Web Audio API limitation - just play without visualization
        console.warn('Cannot create MediaElementSource (already connected), playing without visualization:', err.message)
        // Audio will still play, just without waveform
      }
    }
  } catch (error) {
    console.error('Failed to connect audio source:', error)
  }
}

// Get color based on mode
const getBarColor = () => {
  if (props.barColor) {
    return props.barColor
  }

  switch (props.mode) {
    case 'recording':
      return '#ef4444' // Red
    case 'playing':
      return '#3b82f6' // Blue
    default:
      return '#6b7280' // Gray
  }
}

// Draw waveform visualization
const draw = () => {
  if (!canvasRef.value || !analyser.value || !dataArray.value) {
    return
  }

  const canvas = canvasRef.value
  const canvasCtx = canvas.getContext('2d')
  const width = canvas.width
  const height = canvas.height

  // Get frequency data
  analyser.value.getByteFrequencyData(dataArray.value)

  // Clear canvas
  canvasCtx.fillStyle = 'rgba(255, 255, 255, 0.1)'
  canvasCtx.fillRect(0, 0, width, height)

  // Draw bars
  const barWidth = (width / dataArray.value.length) * 2.5
  let barHeight
  let x = 0

  const barColor = getBarColor()

  for (let i = 0; i < dataArray.value.length; i++) {
    barHeight = (dataArray.value[i] / 255) * height * 0.8

    // Create gradient
    const gradient = canvasCtx.createLinearGradient(0, height, 0, height - barHeight)
    gradient.addColorStop(0, barColor)
    gradient.addColorStop(1, barColor + 'aa') // Add transparency

    canvasCtx.fillStyle = gradient
    canvasCtx.fillRect(x, height - barHeight, barWidth, barHeight)

    x += barWidth + 1
  }

  // Continue animation
  if (props.mode !== 'idle') {
    animationId.value = requestAnimationFrame(draw)
  }
}

// Draw idle state (flat line)
const drawIdle = () => {
  if (!canvasRef.value) {
    return
  }

  const canvas = canvasRef.value
  const canvasCtx = canvas.getContext('2d')
  const width = canvas.width
  const height = canvas.height

  // Clear canvas
  canvasCtx.fillStyle = '#f9fafb'
  canvasCtx.fillRect(0, 0, width, height)

  // Draw center line
  canvasCtx.strokeStyle = '#d1d5db'
  canvasCtx.lineWidth = 2
  canvasCtx.beginPath()
  canvasCtx.moveTo(0, height / 2)
  canvasCtx.lineTo(width, height / 2)
  canvasCtx.stroke()
}

// Start visualization
const startVisualization = () => {
  if (props.mode === 'idle') {
    drawIdle()
    return
  }

  if (!audioContext.value) {
    initAudioContext()
  }

  // Resume audio context if suspended
  if (audioContext.value.state === 'suspended') {
    audioContext.value.resume()
  }

  connectAudioSource()
  draw()
}

// Stop visualization
const stopVisualization = () => {
  if (animationId.value) {
    cancelAnimationFrame(animationId.value)
    animationId.value = null
  }

  if (sourceNode.value) {
    sourceNode.value.disconnect()
    sourceNode.value = null
  }
}

// Watch for mode changes
watch(() => props.mode, (newMode, oldMode) => {
  console.log('Mode changed:', oldMode, '->', newMode)
  stopVisualization()

  if (newMode === 'idle') {
    drawIdle()
  } else if (props.audioSource) {
    // Only start if audioSource is already available
    startVisualization()
  } else {
    console.log('Mode is', newMode, 'but audioSource not yet available, waiting...')
  }
})

// Watch for audioSource changes (for when mode changes before stream is ready)
watch(() => props.audioSource, (newSource, oldSource) => {
  console.log('AudioSource changed:', oldSource ? 'exists' : 'null', '->', newSource ? 'exists' : 'null')

  // Only react if we're in a non-idle mode and source is now available
  if (newSource && props.mode !== 'idle') {
    console.log('AudioSource now available, starting visualization')
    startVisualization()
  }
})

onMounted(() => {
  drawIdle()
})

onUnmounted(() => {
  stopVisualization()

  if (audioContext.value) {
    audioContext.value.close()
  }
})
</script>

<style scoped>
.waveform-container {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.waveform-canvas {
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  background: linear-gradient(135deg, #f9fafb 0%, #f3f4f6 100%);
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: white;
  border-radius: 9999px;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

.status-dot.status-recording {
  background-color: #ef4444;
}

.status-dot.status-playing {
  background-color: #3b82f6;
}

.status-text {
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>
