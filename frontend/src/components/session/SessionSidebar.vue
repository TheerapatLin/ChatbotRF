<template>
  <div class="session-sidebar">
    <!-- Header -->
    <div class="sidebar-header">
      <h2>Chats</h2>
      <button @click="handleNewChat" class="new-chat-btn">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 5v14M5 12h14"/>
        </svg>
        New Chat
      </button>
    </div>

    <!-- Sessions List -->
    <div class="sessions-list">
      <div v-if="sessionStore.isLoading" class="loading">
        <div class="spinner"></div>
        Loading...
      </div>

      <div
        v-for="session in sessionStore.sortedSessions"
        :key="session.session_id"
        :class="['session-item', {
          active: session.session_id === chatStore.sessionId
        }]"
        @click="handleSwitchSession(session.session_id)"
      >
        <div class="session-info">
          <h3>{{ sessionStore.getSessionTitle(session) }}</h3>
          <div class="session-meta">
            <span class="count">{{ session.message_count }} msgs</span>
            <span class="date">{{ formatDate(session.last_message_at) }}</span>
          </div>
        </div>

        <div class="session-actions">
          <button
            @click.stop="handleRename(session)"
            class="action-btn"
            title="Rename"
          >
            ✏️
          </button>
          <button
            @click.stop="handleDelete(session.session_id)"
            class="action-btn delete"
            title="Delete"
          >
            🗑️
          </button>
        </div>
      </div>

      <div v-if="!sessionStore.isLoading && sessionStore.sessions.length === 0" class="empty-state">
        <p>No chats yet</p>
        <p>Start a new conversation!</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useChatStore } from '@/stores/chat'
import { useSessionStore } from '@/stores/session'

const chatStore = useChatStore()
const sessionStore = useSessionStore()

const handleNewChat = () => {
  sessionStore.createNewChat()
}

const handleSwitchSession = async (sessionId) => {
  try {
    await sessionStore.switchToSession(sessionId)
  } catch (error) {
    alert('Failed to load chat')
  }
}

const handleRename = (session) => {
  const currentTitle = sessionStore.getSessionTitle(session)
  const newTitle = prompt('Enter new title:', currentTitle)

  if (newTitle && newTitle.trim() !== '' && newTitle !== currentTitle) {
    sessionStore.renameSession(session.session_id, newTitle.trim())
  }
}

const handleDelete = (sessionId) => {
  if (confirm('Delete this chat? This cannot be undone.')) {
    sessionStore.deleteSession(sessionId)
  }
}

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return 'Just now'
  if (diffMins < 60) return `${diffMins}m`
  if (diffHours < 24) return `${diffHours}h`
  if (diffDays < 7) return `${diffDays}d`

  return date.toLocaleDateString('th-TH', {
    month: 'short',
    day: 'numeric'
  })
}

onMounted(() => {
  sessionStore.loadSessions()
})
</script>

<style scoped>
.session-sidebar {
  width: 260px;
  background: #f9fafb;
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
  background: white;
}

.sidebar-header h2 {
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 12px 0;
  color: #111827;
}

.new-chat-btn {
  width: 100%;
  padding: 10px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  transition: transform 0.2s;
}

.new-chat-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.sessions-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.session-item {
  padding: 10px 12px;
  margin-bottom: 4px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.session-item:hover {
  border-color: #667eea;
  transform: translateX(2px);
}

.session-item.active {
  background: linear-gradient(135deg, #eff6ff 0%, #f3f0ff 100%);
  border-color: #667eea;
  border-width: 2px;
}

.session-info {
  flex: 1;
  min-width: 0;
}

.session-info h3 {
  font-size: 13px;
  font-weight: 600;
  margin: 0 0 4px 0;
  color: #111827;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.session-meta {
  display: flex;
  gap: 8px;
  font-size: 11px;
  color: #6b7280;
}

.session-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.2s;
}

.session-item:hover .session-actions {
  opacity: 1;
}

.action-btn {
  padding: 4px;
  background: none;
  border: none;
  cursor: pointer;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.15s;
}

.action-btn:hover {
  background: #f3f4f6;
  transform: scale(1.1);
}

.action-btn.delete:hover {
  background: #fee2e2;
}

.loading {
  padding: 40px 20px;
  text-align: center;
  color: #9ca3af;
  font-size: 13px;
}

.empty-state {
  padding: 60px 20px;
  text-align: center;
  color: #9ca3af;
  font-size: 13px;
}

.empty-state p {
  margin: 4px 0;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #e5e7eb;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  margin: 0 auto 8px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Scrollbar */
.sessions-list::-webkit-scrollbar {
  width: 6px;
}

.sessions-list::-webkit-scrollbar-track {
  background: transparent;
}

.sessions-list::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 3px;
}

.sessions-list::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}
</style>
