<template>
  <aside class="persona-sidebar">
    <div class="sidebar-header">
      <h2 class="sidebar-title">AI Personas</h2>
      <BaseButton size="sm" @click="showCreateModal">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 5v14M5 12h14" />
        </svg>
        New
      </BaseButton>
    </div>

    <!-- Persona Selector -->
    <div class="persona-selector">
      <BaseSelect
        v-model="selectedPersonaId"
        label="Select Persona"
        placeholder="Choose a persona..."
        :options="personaOptions"
        value-key="id"
        label-key="name"
        :disabled="isLoadingPersonas"
        @change="handlePersonaChange"
      />
    </div>

    <!-- Loading State -->
    <div v-if="isLoadingPersonas" class="loading-state">
      <div class="spinner"></div>
      <p>Loading personas...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <p class="error-message">{{ error }}</p>
      <BaseButton size="sm" variant="secondary" @click="retryLoad">Retry</BaseButton>
    </div>

    <!-- Persona Details -->
    <div v-else-if="selectedPersona" class="persona-details">
      <div class="detail-header">
        <span class="persona-icon">{{ selectedPersona.icon }}</span>
        <h3 class="persona-name">{{ selectedPersona.name }}</h3>
      </div>

      <div class="detail-section">
        <label class="detail-label">Description</label>
        <p class="detail-value">{{ selectedPersona.description }}</p>
      </div>

      <div class="detail-grid">
        <div class="detail-section">
          <label class="detail-label">Model</label>
          <p class="detail-value">{{ selectedPersona.model }}</p>
        </div>

        <div class="detail-section">
          <label class="detail-label">Tone</label>
          <p class="detail-value">{{ selectedPersona.tone }}</p>
        </div>
      </div>

      <div class="detail-grid">
        <div class="detail-section">
          <label class="detail-label">Temperature</label>
          <p class="detail-value">{{ selectedPersona.temperature }}</p>
        </div>

        <div class="detail-section">
          <label class="detail-label">Max Tokens</label>
          <p class="detail-value">{{ selectedPersona.max_tokens }}</p>
        </div>
      </div>

      <div class="detail-section">
        <label class="detail-label">System Prompt</label>
        <p class="detail-value detail-prompt">{{ selectedPersona.system_prompt }}</p>
      </div>

      <!-- Action Buttons -->
      <div class="action-buttons">
        <BaseButton size="sm" variant="outline" full-width @click="handleEdit">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
          </svg>
          Edit
        </BaseButton>
        <BaseButton size="sm" variant="danger" full-width @click="handleDelete">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
          </svg>
          Delete
        </BaseButton>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="empty-state">
      <p>No persona selected</p>
      <BaseButton size="sm" @click="showCreateModal">Create Persona</BaseButton>
    </div>

    <!-- Debug Panel -->
    <DebugPanel />

    <!-- Persona Modal -->
    <PersonaModal
      v-model="isModalOpen"
      :persona="editingPersona"
      :mode="modalMode"
      @saved="handlePersonaSaved"
    />
  </aside>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { usePersonaStore } from '@/stores/persona'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseSelect from '@/components/common/BaseSelect.vue'
import DebugPanel from './DebugPanel.vue'
import PersonaModal from './PersonaModal.vue'

const personaStore = usePersonaStore()

// State
const isModalOpen = ref(false)
const modalMode = ref('create') // 'create' or 'edit'
const editingPersona = ref(null)

// Computed
const selectedPersonaId = computed({
  get: () => personaStore.selectedPersona?.id || '',
  set: (value) => {
    const persona = personaStore.personas.find((p) => p.id === parseInt(value))
    if (persona) {
      personaStore.selectPersona(persona)
    }
  }
})

const personaOptions = computed(() => personaStore.personas)
const selectedPersona = computed(() => personaStore.selectedPersona)
const isLoadingPersonas = computed(() => personaStore.isLoadingPersonas)
const error = computed(() => personaStore.error)

// Methods
const showCreateModal = () => {
  modalMode.value = 'create'
  editingPersona.value = null
  isModalOpen.value = true
}

const handleEdit = () => {
  modalMode.value = 'edit'
  editingPersona.value = { ...selectedPersona.value }
  isModalOpen.value = true
}

const handleDelete = async () => {
  if (!selectedPersona.value) return

  const confirmed = confirm(
    `Are you sure you want to delete "${selectedPersona.value.name}"?\nThis will also delete all associated messages.`
  )

  if (confirmed) {
    try {
      await personaStore.deletePersona(selectedPersona.value.id)
      alert('Persona deleted successfully!')
    } catch (error) {
      alert(`Failed to delete persona: ${error.message}`)
    }
  }
}

const handlePersonaChange = (personaId) => {
  const persona = personaStore.personas.find((p) => p.id === parseInt(personaId))
  if (persona) {
    personaStore.selectPersona(persona)
  }
}

const handlePersonaSaved = () => {
  isModalOpen.value = false
  editingPersona.value = null
}

const retryLoad = async () => {
  try {
    await personaStore.fetchPersonas()
  } catch (error) {
    console.error('Failed to retry loading personas:', error)
  }
}

// Lifecycle
onMounted(async () => {
  try {
    await personaStore.fetchPersonas()
  } catch (error) {
    console.error('Failed to load personas:', error)
  }
})
</script>

<style scoped>
.persona-sidebar {
  width: 320px;
  height: 100vh;
  background: linear-gradient(180deg, #f8f9fa 0%, #ffffff 100%);
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  padding: 20px;
  gap: 20px;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sidebar-title {
  font-size: 20px;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}

.persona-selector {
  margin-bottom: 10px;
}

/* Loading State */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  gap: 12px;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e5e7eb;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Error State */
.error-state {
  padding: 20px;
  text-align: center;
}

.error-message {
  color: #dc3545;
  margin-bottom: 12px;
}

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
  gap: 12px;
}

.empty-state p {
  color: #6b7280;
  margin: 0;
}

/* Persona Details */
.persona-details {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e5e7eb;
}

.persona-icon {
  font-size: 32px;
}

.persona-name {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.detail-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.detail-label {
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.detail-value {
  font-size: 14px;
  color: #374151;
  margin: 0;
}

.detail-prompt {
  font-size: 13px;
  line-height: 1.5;
  max-height: 120px;
  overflow-y: auto;
  padding: 8px;
  background: #f9fafb;
  border-radius: 4px;
  border: 1px solid #e5e7eb;
}

.action-buttons {
  display: flex;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid #e5e7eb;
}
</style>
