<template>
  <BaseModal
    v-model="isOpen"
    :title="modalTitle"
    size="lg"
    @close="handleClose"
  >
    <form @submit.prevent="handleSubmit">
      <div class="form-content">
        <!-- Basic Info -->
        <div class="form-section">
          <h4 class="section-title">Basic Information</h4>

          <div class="form-row">
            <BaseInput
              v-model="formData.name"
              label="Name"
              placeholder="e.g., Marketing Expert"
              required
              :error="errors.name"
            />
          </div>

          <div class="form-row">
            <BaseInput
              v-model="formData.icon"
              label="Icon"
              placeholder="e.g., 🤖"
              :hint="'Choose an emoji (max 10 characters)'"
              :error="errors.icon"
            />
          </div>

          <div class="form-row">
            <BaseTextarea
              v-model="formData.description"
              label="Description"
              placeholder="Describe this persona..."
              required
              :rows="3"
              :error="errors.description"
            />
          </div>

          <div class="form-row">
            <BaseTextarea
              v-model="formData.system_prompt"
              label="System Prompt"
              placeholder="You are a helpful assistant..."
              required
              :rows="4"
              :error="errors.system_prompt"
            />
          </div>
        </div>

        <!-- Personality Settings -->
        <div class="form-section">
          <h4 class="section-title">Personality</h4>

          <div class="form-grid">
            <BaseSelect
              v-model="formData.tone"
              label="Tone"
              :options="toneOptions"
              placeholder="Select tone..."
              :error="errors.tone"
            />

            <BaseSelect
              v-model="formData.style"
              label="Style"
              :options="styleOptions"
              placeholder="Select style..."
              :error="errors.style"
            />
          </div>

          <div class="form-row">
            <BaseInput
              v-model="formData.expertise"
              label="Expertise"
              placeholder="e.g., marketing, technology, general"
              :error="errors.expertise"
            />
          </div>
        </div>

        <!-- Model Settings -->
        <div class="form-section">
          <h4 class="section-title">Model Configuration</h4>

          <div class="form-row">
            <BaseSelect
              v-model="formData.model"
              label="Model"
              :options="modelOptions"
              placeholder="Select model..."
              required
              :error="errors.model"
            />
          </div>

          <div class="form-grid">
            <BaseInput
              v-model.number="formData.temperature"
              label="Temperature"
              type="number"
              step="0.1"
              min="0"
              max="2"
              :hint="'0.0 = Focused, 2.0 = Creative'"
              :error="errors.temperature"
            />

            <BaseInput
              v-model.number="formData.max_tokens"
              label="Max Tokens"
              type="number"
              step="100"
              min="100"
              :hint="'Maximum response length'"
              :error="errors.max_tokens"
            />
          </div>
        </div>

        <!-- Language Settings -->
        <div class="form-section">
          <h4 class="section-title">Language Settings</h4>

          <div class="form-grid">
            <BaseSelect
              v-model="languageSettings.default_language"
              label="Default Language"
              :options="languageOptions"
            />

            <BaseSelect
              v-model="languageSettings.response_style"
              label="Response Style"
              :options="responseStyleOptions"
            />
          </div>

          <div class="form-row">
            <BaseInput
              v-model="languageSettings.language_code"
              label="Language Code"
              placeholder="e.g., en-US, th-TH"
            />
          </div>
        </div>

        <!-- Guardrails -->
        <div class="form-section">
          <h4 class="section-title">Content Guardrails</h4>

          <div class="form-grid">
            <div class="checkbox-group">
              <label class="checkbox-label">
                <input v-model="guardrails.block_profanity" type="checkbox" />
                Block Profanity
              </label>
            </div>

            <div class="checkbox-group">
              <label class="checkbox-label">
                <input v-model="guardrails.block_sensitive" type="checkbox" />
                Block Sensitive Content
              </label>
            </div>
          </div>

          <div class="form-row">
            <BaseInput
              v-model="guardrails.allowed_topics"
              label="Allowed Topics"
              placeholder="marketing, business, technology (comma-separated)"
              :hint="'Leave empty to allow all topics'"
            />
          </div>

          <div class="form-row">
            <BaseInput
              v-model="guardrails.blocked_topics"
              label="Blocked Topics"
              placeholder="politics, religion (comma-separated)"
              :hint="'Topics to avoid'"
            />
          </div>

          <div class="form-row">
            <BaseInput
              v-model.number="guardrails.max_response_length"
              label="Max Response Length"
              type="number"
              placeholder="3000"
              :hint="'Maximum characters in response'"
            />
          </div>

          <div class="checkbox-group">
            <label class="checkbox-label">
              <input v-model="guardrails.require_moderation" type="checkbox" />
              Require Moderation
            </label>
          </div>
        </div>
      </div>
    </form>

    <template #footer>
      <BaseButton variant="secondary" @click="handleClose"> Cancel </BaseButton>
      <BaseButton :loading="isSubmitting" @click="handleSubmit">
        {{ mode === 'create' ? 'Create Persona' : 'Update Persona' }}
      </BaseButton>
    </template>
  </BaseModal>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { usePersonaStore } from '@/stores/persona'
import BaseModal from '@/components/common/BaseModal.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import BaseTextarea from '@/components/common/BaseTextarea.vue'
import BaseSelect from '@/components/common/BaseSelect.vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true
  },
  persona: {
    type: Object,
    default: null
  },
  mode: {
    type: String,
    default: 'create',
    validator: (value) => ['create', 'edit'].includes(value)
  }
})

const emit = defineEmits(['update:modelValue', 'saved'])

const personaStore = usePersonaStore()

// State
const isSubmitting = ref(false)
const errors = ref({})

// Form Data
const formData = ref({
  name: '',
  description: '',
  system_prompt: '',
  tone: 'friendly',
  style: 'conversational',
  expertise: 'general',
  temperature: 0.7,
  max_tokens: 2000,
  model: 'gpt-4o-mini',
  icon: '🤖',
  is_active: true
})

const languageSettings = ref({
  default_language: 'en',
  response_style: 'balanced',
  language_code: 'en-US'
})

const guardrails = ref({
  block_profanity: false,
  block_sensitive: false,
  allowed_topics: '',
  blocked_topics: '',
  max_response_length: 3000,
  require_moderation: false
})

// Options
const toneOptions = ['friendly', 'professional', 'empathetic', 'mystical', 'enthusiastic']
const styleOptions = ['conversational', 'detailed', 'concise']
const modelOptions = [
  { value: 'gpt-4o-mini', label: 'GPT-4o Mini ⭐ (Fast & Affordable)' },
  { value: 'gpt-4o', label: 'GPT-4o (Latest & Fast)' },
  { value: 'gpt-4-turbo', label: 'GPT-4 Turbo (Fast)' },
  { value: 'gpt-4', label: 'GPT-4 (Most Accurate)' },
  { value: 'gpt-3.5-turbo', label: 'GPT-3.5 Turbo (Economical)' }
]
const languageOptions = ['en', 'th', 'ja', 'zh', 'es', 'fr', 'de']
const responseStyleOptions = ['formal', 'casual', 'balanced']

// Computed
const isOpen = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const modalTitle = computed(() => {
  return props.mode === 'create' ? 'Create New Persona' : 'Edit Persona'
})

// Methods
const resetForm = () => {
  formData.value = {
    name: '',
    description: '',
    system_prompt: '',
    tone: 'friendly',
    style: 'conversational',
    expertise: 'general',
    temperature: 0.7,
    max_tokens: 2000,
    model: 'gpt-4o-mini',
    icon: '🤖',
    is_active: true
  }

  languageSettings.value = {
    default_language: 'en',
    response_style: 'balanced',
    language_code: 'en-US'
  }

  guardrails.value = {
    block_profanity: false,
    block_sensitive: false,
    allowed_topics: '',
    blocked_topics: '',
    max_response_length: 3000,
    require_moderation: false
  }

  errors.value = {}
}

const loadPersonaData = (persona) => {
  if (!persona) return

  formData.value = {
    name: persona.name || '',
    description: persona.description || '',
    system_prompt: persona.system_prompt || '',
    tone: persona.tone || 'friendly',
    style: persona.style || 'conversational',
    expertise: persona.expertise || 'general',
    temperature: persona.temperature || 0.7,
    max_tokens: persona.max_tokens || 2000,
    model: persona.model || 'gpt-4o-mini',
    icon: persona.icon || '🤖',
    is_active: persona.is_active !== undefined ? persona.is_active : true
  }

  // Parse language_setting
  if (persona.language_setting) {
    try {
      const parsed =
        typeof persona.language_setting === 'string'
          ? JSON.parse(persona.language_setting)
          : persona.language_setting
      languageSettings.value = {
        default_language: parsed.default_language || 'en',
        response_style: parsed.response_style || 'balanced',
        language_code: parsed.language_code || 'en-US'
      }
    } catch (e) {
      console.error('Failed to parse language_setting:', e)
    }
  }

  // Parse guardrails
  if (persona.guardrails) {
    try {
      const parsed =
        typeof persona.guardrails === 'string' ? JSON.parse(persona.guardrails) : persona.guardrails
      guardrails.value = {
        block_profanity: parsed.block_profanity || false,
        block_sensitive: parsed.block_sensitive || false,
        allowed_topics: Array.isArray(parsed.allowed_topics)
          ? parsed.allowed_topics.join(', ')
          : '',
        blocked_topics: Array.isArray(parsed.blocked_topics)
          ? parsed.blocked_topics.join(', ')
          : '',
        max_response_length: parsed.max_response_length || 3000,
        require_moderation: parsed.require_moderation || false
      }
    } catch (e) {
      console.error('Failed to parse guardrails:', e)
    }
  }
}

const validateForm = () => {
  errors.value = {}

  if (!formData.value.name) {
    errors.value.name = 'Name is required'
  }

  if (!formData.value.description) {
    errors.value.description = 'Description is required'
  }

  if (!formData.value.system_prompt) {
    errors.value.system_prompt = 'System prompt is required'
  }

  if (formData.value.temperature < 0 || formData.value.temperature > 2) {
    errors.value.temperature = 'Temperature must be between 0 and 2'
  }

  return Object.keys(errors.value).length === 0
}

const handleSubmit = async () => {
  if (!validateForm()) {
    alert('Please fix the errors in the form')
    return
  }

  isSubmitting.value = true

  try {
    // Prepare data
    const personaData = {
      ...formData.value,
      language_setting: languageSettings.value,
      guardrails: {
        ...guardrails.value,
        allowed_topics: guardrails.value.allowed_topics
          ? guardrails.value.allowed_topics.split(',').map((t) => t.trim())
          : [],
        blocked_topics: guardrails.value.blocked_topics
          ? guardrails.value.blocked_topics.split(',').map((t) => t.trim())
          : []
      }
    }

    if (props.mode === 'create') {
      await personaStore.createPersona(personaData)
      alert('Persona created successfully!')
    } else {
      await personaStore.updatePersona(props.persona.id, personaData)
      alert('Persona updated successfully!')
    }

    emit('saved')
    handleClose()
  } catch (error) {
    alert(`Failed to ${props.mode} persona: ${error.message}`)
  } finally {
    isSubmitting.value = false
  }
}

const handleClose = () => {
  resetForm()
  emit('update:modelValue', false)
}

// Watch for persona changes
watch(
  () => props.persona,
  (newPersona) => {
    if (newPersona && props.mode === 'edit') {
      loadPersonaData(newPersona)
    } else {
      resetForm()
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.form-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
  max-height: 70vh;
  overflow-y: auto;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 20px;
  background: #f9fafb;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.section-title {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
  color: #374151;
}

.form-row {
  display: flex;
  flex-direction: column;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.checkbox-group {
  display: flex;
  align-items: center;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #374151;
  cursor: pointer;
  user-select: none;
}

.checkbox-label input[type='checkbox'] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
