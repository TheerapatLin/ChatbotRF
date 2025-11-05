<template>
  <div class="debug-panel">
    <div class="panel-header">
      <h3 class="panel-title">Debug Info</h3>
      <button class="toggle-btn" @click="isExpanded = !isExpanded">
        <svg
          :style="{ transform: isExpanded ? 'rotate(180deg)' : 'rotate(0deg)' }"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M6 9l6 6 6-6" />
        </svg>
      </button>
    </div>

    <Transition name="slide">
      <div v-if="isExpanded" class="panel-content">
        <!-- Stats -->
        <div class="stats-grid">
          <div class="stat-item">
            <label class="stat-label">Messages</label>
            <p class="stat-value">{{ messageCount }}</p>
          </div>

          <div class="stat-item">
            <label class="stat-label">Session</label>
            <p class="stat-value">{{ truncateSessionId }}</p>
          </div>

          <div class="stat-item">
            <label class="stat-label">Tokens</label>
            <p class="stat-value">{{ totalTokens }}</p>
          </div>

          <div class="stat-item">
            <label class="stat-label">WebSocket</label>
            <p class="stat-value" :class="wsStatusClass">{{ wsStatus }}</p>
          </div>
        </div>

        <!-- Latency -->
        <div v-if="lastLatency" class="info-item">
          <label class="info-label">Last Response Time</label>
          <p class="info-value">{{ lastLatency }}ms</p>
        </div>

        <!-- API Usage (if available) -->
        <div v-if="apiUsage" class="info-item">
          <label class="info-label">API Usage</label>
          <p class="info-value">{{ apiUsage }}</p>
        </div>

        <!-- Actions -->
        <div class="action-section">
          <BaseButton
            size="sm"
            variant="danger"
            full-width
            :loading="isClearingHistory"
            @click="handleClearHistory"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
            </svg>
            Clear History
          </BaseButton>

          <BaseButton
            size="sm"
            variant="secondary"
            full-width
            @click="handleRefreshStats"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 2v6h-6M3 12a9 9 0 0 1 15-6.7L21 8M3 22v-6h6M21 12a9 9 0 0 1-15 6.7L3 16" />
            </svg>
            Refresh
          </BaseButton>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useChatStore } from '@/stores/chat'
import BaseButton from '@/components/common/BaseButton.vue'

const chatStore = useChatStore()

// State
const isExpanded = ref(true)
const isClearingHistory = ref(false)
const lastLatency = ref(null)
const apiUsage = ref(null)
const refreshInterval = ref(null)

// Computed
const messageCount = computed(() => chatStore.messages.length)

const truncateSessionId = computed(() => {
  const sessionId = chatStore.sessionId
  if (!sessionId) return 'N/A'
  return sessionId.length > 15 ? sessionId.substring(0, 15) + '...' : sessionId
})

const totalTokens = computed(() => {
  return chatStore.messages.reduce((sum, msg) => {
    return sum + (msg.tokens_used || 0)
  }, 0)
})

const wsStatus = computed(() => {
  if (chatStore.isConnected) return 'Connected'
  if (chatStore.isLoading) return 'Connecting...'
  return 'Disconnected'
})

const wsStatusClass = computed(() => ({
  'status-connected': chatStore.isConnected,
  'status-connecting': chatStore.isLoading,
  'status-disconnected': !chatStore.isConnected && !chatStore.isLoading
}))

// Methods
const handleClearHistory = async () => {
  const confirmed = confirm(
    'Are you sure you want to clear all chat history?\nThis action cannot be undone.'
  )

  if (!confirmed) return

  isClearingHistory.value = true

  try {
    await chatStore.clearChatHistory()
    alert('Chat history cleared successfully!')
  } catch (error) {
    alert(`Failed to clear history: ${error.message}`)
  } finally {
    isClearingHistory.value = false
  }
}

const handleRefreshStats = () => {
  // Calculate latency from last message if available
  if (chatStore.messages.length >= 2) {
    const lastMsg = chatStore.messages[chatStore.messages.length - 1]
    const prevMsg = chatStore.messages[chatStore.messages.length - 2]

    if (lastMsg.timestamp && prevMsg.timestamp) {
      const timeDiff = new Date(lastMsg.timestamp) - new Date(prevMsg.timestamp)
      lastLatency.value = Math.round(timeDiff)
    }
  }
}

const startAutoRefresh = () => {
  refreshInterval.value = setInterval(() => {
    handleRefreshStats()
  }, 5000) // Refresh every 5 seconds
}

const stopAutoRefresh = () => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
    refreshInterval.value = null
  }
}

// Lifecycle
onMounted(() => {
  handleRefreshStats()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.debug-panel {
  margin-top: auto;
  padding: 16px;
  background: white;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.panel-title {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin: 0;
}

.toggle-btn {
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s;
}

.toggle-btn:hover {
  background: #f3f4f6;
  color: #374151;
}

.toggle-btn svg {
  transition: transform 0.2s;
}

.panel-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 11px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-value {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.status-connected {
  color: #10b981;
}

.status-connecting {
  color: #f59e0b;
}

.status-disconnected {
  color: #ef4444;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
  background: #f9fafb;
  border-radius: 6px;
}

.info-label {
  font-size: 11px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-value {
  font-size: 13px;
  color: #374151;
  margin: 0;
}

.action-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-top: 8px;
  border-top: 1px solid #e5e7eb;
}

/* Transitions */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
  max-height: 500px;
  overflow: hidden;
}

.slide-enter-from,
.slide-leave-to {
  max-height: 0;
  opacity: 0;
}
</style>
