<template>
  <div class="microphone-button-container">
    <button
      @click="handleClick"
      @mousedown="handleMouseDown"
      @mouseup="handleMouseUp"
      @touchstart="handleMouseDown"
      @touchend="handleMouseUp"
      :disabled="disabled || state === 'processing'"
      :class="[
        'mic-button',
        `mic-button-${state}`,
        { 'mic-button-disabled': disabled }
      ]"
    >
      <!-- Microphone Icon -->
      <svg
        v-if="state === 'idle'"
        class="mic-icon"
        fill="currentColor"
        viewBox="0 0 20 20"
      >
        <path
          fill-rule="evenodd"
          d="M7 4a3 3 0 016 0v4a3 3 0 11-6 0V4zm4 10.93A7.001 7.001 0 0017 8a1 1 0 10-2 0A5 5 0 015 8a1 1 0 00-2 0 7.001 7.001 0 006 6.93V17H6a1 1 0 100 2h8a1 1 0 100-2h-3v-2.07z"
          clip-rule="evenodd"
        />
      </svg>

      <!-- Recording Icon (animated) -->
      <div v-else-if="state === 'recording'" class="recording-indicator">
        <div class="recording-pulse"></div>
        <svg class="mic-icon" fill="currentColor" viewBox="0 0 20 20">
          <path
            fill-rule="evenodd"
            d="M7 4a3 3 0 016 0v4a3 3 0 11-6 0V4zm4 10.93A7.001 7.001 0 0017 8a1 1 0 10-2 0A5 5 0 015 8a1 1 0 00-2 0 7.001 7.001 0 006 6.93V17H6a1 1 0 100 2h8a1 1 0 100-2h-3v-2.07z"
            clip-rule="evenodd"
          />
        </svg>
      </div>

      <!-- Processing Icon (spinner) -->
      <svg
        v-else-if="state === 'processing'"
        class="mic-icon spinner"
        fill="none"
        viewBox="0 0 24 24"
      >
        <circle
          class="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          stroke-width="4"
        ></circle>
        <path
          class="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        ></path>
      </svg>

      <!-- Playing Icon -->
      <svg
        v-else-if="state === 'playing'"
        class="mic-icon"
        fill="currentColor"
        viewBox="0 0 20 20"
      >
        <path
          fill-rule="evenodd"
          d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zM7 8a1 1 0 012 0v4a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v4a1 1 0 102 0V8a1 1 0 00-1-1z"
          clip-rule="evenodd"
        />
      </svg>
    </button>

    <!-- Button label -->
    <div class="mic-label">
      {{ buttonLabel }}
    </div>

    <!-- Hint text -->
    <div v-if="showHint && state === 'idle'" class="mic-hint">
      {{ hintText }}
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  // Button state: idle, recording, processing, playing
  state: {
    type: String,
    default: 'idle',
    validator: (value) => ['idle', 'recording', 'processing', 'playing'].includes(value)
  },

  // Disabled state
  disabled: {
    type: Boolean,
    default: false
  },

  // Show hint text
  showHint: {
    type: Boolean,
    default: true
  },

  // Custom hint text
  hintText: {
    type: String,
    default: 'คลิกเพื่อเริ่มพูด'
  }
})

const emit = defineEmits(['click', 'start-recording', 'stop-recording'])

const buttonLabel = computed(() => {
  switch (props.state) {
    case 'idle':
      return 'เริ่มพูด'
    case 'recording':
      return 'กำลังบันทึก...'
    case 'processing':
      return 'กำลังประมวลผล...'
    case 'playing':
      return 'กำลังเล่นเสียง...'
    default:
      return 'เริ่มพูด'
  }
})

const handleClick = () => {
  if (props.disabled || props.state === 'processing' || props.state === 'playing') {
    return
  }

  emit('click')
}

const handleMouseDown = (event) => {
  if (props.disabled || props.state !== 'idle') {
    return
  }

  event.preventDefault()
  emit('start-recording')
}

const handleMouseUp = (event) => {
  if (props.disabled || props.state !== 'recording') {
    return
  }

  event.preventDefault()
  emit('stop-recording')
}
</script>

<style scoped>
.microphone-button-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.mic-button {
  position: relative;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Idle state */
.mic-button-idle {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
}

.mic-button-idle:hover:not(:disabled) {
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
  transform: scale(1.05);
  box-shadow: 0 20px 25px -5px rgba(59, 130, 246, 0.3), 0 10px 10px -5px rgba(59, 130, 246, 0.2);
}

.mic-button-idle:active:not(:disabled) {
  transform: scale(0.95);
}

/* Recording state */
.mic-button-recording {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
  animation: pulse-button 2s ease-in-out infinite;
}

/* Processing state */
.mic-button-processing {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
  cursor: not-allowed;
}

/* Playing state */
.mic-button-playing {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  cursor: not-allowed;
}

/* Disabled state */
.mic-button-disabled {
  background: linear-gradient(135deg, #9ca3af 0%, #6b7280 100%);
  cursor: not-allowed;
  opacity: 0.5;
}

.mic-icon {
  width: 48px;
  height: 48px;
}

/* Recording indicator */
.recording-indicator {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.recording-pulse {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  animation: pulse-wave 2s ease-out infinite;
}

/* Spinner animation */
.spinner {
  animation: spin 1s linear infinite;
}

/* Labels */
.mic-label {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1f2937;
  text-align: center;
}

.mic-hint {
  font-size: 0.875rem;
  color: #6b7280;
  text-align: center;
  padding: 0.5rem 1rem;
  background: #f3f4f6;
  border-radius: 9999px;
}

/* Animations */
@keyframes pulse-button {
  0%, 100% {
    box-shadow: 0 10px 25px -5px rgba(239, 68, 68, 0.4), 0 10px 10px -5px rgba(239, 68, 68, 0.3);
  }
  50% {
    box-shadow: 0 20px 40px -5px rgba(239, 68, 68, 0.6), 0 15px 20px -5px rgba(239, 68, 68, 0.5);
  }
}

@keyframes pulse-wave {
  0% {
    transform: scale(1);
    opacity: 0.7;
  }
  100% {
    transform: scale(2);
    opacity: 0;
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* Touch device optimizations */
@media (hover: none) {
  .mic-button {
    -webkit-tap-highlight-color: transparent;
  }

  .mic-button-idle:active:not(:disabled) {
    transform: scale(0.95);
  }
}
</style>
