<script setup>
import { computed } from 'vue'
import { useChatStore } from '@/stores/chat'
import {
  AI_PROVIDERS,
  OPENAI_MODELS,
  BEDROCK_MODELS,
  MODEL_CATEGORIES
} from '@/config/aiModels'

const chatStore = useChatStore()

const selectedProvider = computed({
  get: () => chatStore.selectedProvider,
  set: (value) => chatStore.setProvider(value)
})

const selectedModel = computed({
  get: () => chatStore.selectedModel,
  set: (value) => chatStore.setModel(value)
})

const temperature = computed({
  get: () => chatStore.temperature,
  set: (value) => chatStore.setTemperature(value)
})

const maxTokens = computed({
  get: () => chatStore.maxTokens,
  set: (value) => chatStore.setMaxTokens(value)
})

const availableModels = computed(() => {
  return selectedProvider.value === AI_PROVIDERS.OPENAI
    ? OPENAI_MODELS
    : BEDROCK_MODELS
})

const groupedModels = computed(() => {
  if (selectedProvider.value === AI_PROVIDERS.OPENAI) {
    return { openai: OPENAI_MODELS }
  }

  // Group Bedrock models by category
  const grouped = {}
  BEDROCK_MODELS.forEach(model => {
    const category = model.category || 'other'
    if (!grouped[category]) {
      grouped[category] = []
    }
    grouped[category].push(model)
  })
  return grouped
})

const currentModel = computed(() => {
  return availableModels.value.find(m => m.id === selectedModel.value)
})
</script>

<template>
  <div class="model-selector">
    <div class="selector-header">
      <h3>🤖 AI Configuration</h3>
    </div>

    <!-- Provider Selection -->
    <div class="form-group">
      <label class="form-label">
        <span class="label-icon">🔌</span>
        AI Provider
      </label>
      <div class="provider-tabs">
        <button
          class="provider-tab"
          :class="{ active: selectedProvider === AI_PROVIDERS.BEDROCK }"
          @click="selectedProvider = AI_PROVIDERS.BEDROCK"
        >
          <span class="tab-icon">🔵</span>
          AWS Bedrock
        </button>
        <button
          class="provider-tab"
          :class="{ active: selectedProvider === AI_PROVIDERS.OPENAI }"
          @click="selectedProvider = AI_PROVIDERS.OPENAI"
        >
          <span class="tab-icon">🟢</span>
          OpenAI
        </button>
      </div>
    </div>

    <!-- Model Selection -->
    <div class="form-group">
      <label class="form-label">
        <span class="label-icon">🎯</span>
        Model
      </label>
      <select v-model="selectedModel" class="form-select">
        <template v-for="(models, category) in groupedModels" :key="category">
          <optgroup
            v-if="selectedProvider === AI_PROVIDERS.BEDROCK"
            :label="MODEL_CATEGORIES[category]?.name || category"
          >
            <option
              v-for="model in models"
              :key="model.id"
              :value="model.id"
            >
              {{ model.name }}
              {{ model.recommended ? '⭐' : '' }}
            </option>
          </optgroup>
          <option
            v-else
            v-for="model in models"
            :key="model.id"
            :value="model.id"
          >
            {{ model.name }}
            {{ model.recommended ? '⭐' : '' }}
          </option>
        </template>
      </select>
    </div>

    <!-- Model Info Display -->
    <div v-if="currentModel" class="model-info">
      <div class="info-row">
        <span class="info-label">📝 Description:</span>
        <span class="info-value">{{ currentModel.description }}</span>
      </div>
      <div class="info-row">
        <span class="info-label">📊 Context:</span>
        <span class="info-value">{{ currentModel.contextWindow }}</span>
      </div>
      <div class="info-row">
        <span class="info-label">💰 Pricing:</span>
        <span class="info-value">{{ currentModel.pricing }}</span>
      </div>
      <div class="info-row">
        <span class="info-label">✨ Capabilities:</span>
        <span class="info-value">
          <span
            v-for="cap in currentModel.capabilities"
            :key="cap"
            class="capability-badge"
          >
            {{ cap }}
          </span>
        </span>
      </div>
    </div>

    <!-- Advanced Parameters -->
    <details class="advanced-section">
      <summary>⚙️ Advanced Parameters</summary>

      <div class="form-group">
        <label class="form-label">
          🌡️ Temperature: {{ temperature.toFixed(2) }}
        </label>
        <input
          type="range"
          v-model.number="temperature"
          min="0"
          max="1"
          step="0.1"
          class="form-range"
        >
        <div class="range-labels">
          <span>Precise (0.0)</span>
          <span>Creative (1.0)</span>
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">
          📏 Max Tokens: {{ maxTokens }}
        </label>
        <input
          type="range"
          v-model.number="maxTokens"
          min="500"
          max="4096"
          step="100"
          class="form-range"
        >
        <div class="range-labels">
          <span>Short (500)</span>
          <span>Long (4096)</span>
        </div>
      </div>
    </details>
  </div>
</template>

<style scoped>
.model-selector {
  background: #f9fafb;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  border: 1px solid #e5e7eb;
}

.selector-header h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.form-group {
  margin-bottom: 16px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
}

.label-icon {
  margin-right: 6px;
}

.provider-tabs {
  display: flex;
  gap: 8px;
}

.provider-tab {
  flex: 1;
  padding: 10px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  transition: all 0.2s;
}

.provider-tab:hover {
  border-color: #d1d5db;
}

.provider-tab.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.tab-icon {
  margin-right: 6px;
}

.form-select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}

.form-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.model-info {
  background: white;
  border-radius: 8px;
  padding: 12px;
  font-size: 13px;
  border: 1px solid #e5e7eb;
}

.info-row {
  display: flex;
  margin-bottom: 8px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-label {
  font-weight: 500;
  color: #6b7280;
  min-width: 110px;
}

.info-value {
  color: #111827;
  flex: 1;
}

.capability-badge {
  display: inline-block;
  padding: 2px 8px;
  background: #ede9fe;
  color: #7c3aed;
  border-radius: 4px;
  font-size: 12px;
  margin-right: 4px;
}

.advanced-section {
  margin-top: 16px;
  border-top: 1px solid #e5e7eb;
  padding-top: 16px;
}

.advanced-section summary {
  cursor: pointer;
  font-weight: 500;
  color: #374151;
  margin-bottom: 16px;
}

.form-range {
  width: 100%;
  height: 6px;
  border-radius: 3px;
  background: #e5e7eb;
  outline: none;
}

.range-labels {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #6b7280;
  margin-top: 4px;
}
</style>
