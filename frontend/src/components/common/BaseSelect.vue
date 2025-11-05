<template>
  <div class="select-wrapper">
    <label v-if="label" :for="id" class="select-label">
      {{ label }}
      <span v-if="required" class="required">*</span>
    </label>
    <select
      :id="id"
      v-model="localValue"
      :disabled="disabled"
      :required="required"
      :class="selectClasses"
      @change="handleChange"
    >
      <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
      <option
        v-for="option in options"
        :key="getOptionValue(option)"
        :value="getOptionValue(option)"
      >
        {{ getOptionLabel(option) }}
      </option>
    </select>
    <p v-if="error" class="select-error">{{ error }}</p>
    <p v-else-if="hint" class="select-hint">{{ hint }}</p>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  options: {
    type: Array,
    required: true
  },
  valueKey: {
    type: String,
    default: 'value'
  },
  labelKey: {
    type: String,
    default: 'label'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  required: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  },
  hint: {
    type: String,
    default: ''
  },
  id: {
    type: String,
    default: () => `select-${Math.random().toString(36).substr(2, 9)}`
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const localValue = ref(props.modelValue)

watch(
  () => props.modelValue,
  (newValue) => {
    localValue.value = newValue
  }
)

const selectClasses = computed(() => ({
  'select-field': true,
  'select-error-state': !!props.error,
  'select-disabled': props.disabled
}))

const getOptionValue = (option) => {
  return typeof option === 'object' ? option[props.valueKey] : option
}

const getOptionLabel = (option) => {
  return typeof option === 'object' ? option[props.labelKey] : option
}

const handleChange = (event) => {
  const value = event.target.value
  emit('update:modelValue', value)
  emit('change', value)
}
</script>

<style scoped>
.select-wrapper {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.select-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.required {
  color: #dc3545;
  margin-left: 2px;
}

.select-field {
  padding: 10px 14px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
  color: #1f2937;
  background: white;
  transition: all 0.2s;
  outline: none;
  cursor: pointer;
}

.select-field:hover:not(.select-disabled) {
  border-color: #9ca3af;
}

.select-field:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.select-error-state {
  border-color: #dc3545;
}

.select-error-state:focus {
  border-color: #dc3545;
  box-shadow: 0 0 0 3px rgba(220, 53, 69, 0.1);
}

.select-disabled {
  background: #f3f4f6;
  color: #9ca3af;
  cursor: not-allowed;
}

.select-error {
  font-size: 13px;
  color: #dc3545;
  margin: 0;
}

.select-hint {
  font-size: 13px;
  color: #6b7280;
  margin: 0;
}
</style>
