<template>
  <div class="loading-indicator">
    <div class="loading-avatar">
      <span class="avatar-icon">{{ personaIcon }}</span>
    </div>

    <div class="loading-content">
      <div class="loading-bubble">
        <div class="loading-dots">
          <span class="dot"></span>
          <span class="dot"></span>
          <span class="dot"></span>
        </div>
        <p v-if="message" class="loading-message">{{ message }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  message: {
    type: String,
    default: 'Bot กำลังประมวลผล...'
  },
  personaIcon: {
    type: String,
    default: '🤖'
  }
})
</script>

<style scoped>
.loading-indicator {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.loading-avatar {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.avatar-icon {
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.loading-content {
  flex: 1;
  max-width: 70%;
}

.loading-bubble {
  padding: 12px 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  border-bottom-left-radius: 4px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.loading-dots {
  display: flex;
  gap: 6px;
  align-items: center;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #9ca3af;
  animation: bounce 1.4s infinite ease-in-out;
}

.dot:nth-child(1) {
  animation-delay: -0.32s;
}

.dot:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes bounce {
  0%,
  80%,
  100% {
    transform: scale(0);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

.loading-message {
  margin: 0;
  font-size: 13px;
  color: #6b7280;
  font-style: italic;
}

@media (max-width: 768px) {
  .loading-content {
    max-width: 85%;
  }

  .loading-avatar {
    width: 32px;
    height: 32px;
    font-size: 18px;
  }
}
</style>
