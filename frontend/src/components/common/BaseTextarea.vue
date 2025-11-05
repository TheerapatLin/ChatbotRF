<template>
  <div class="textarea-wrapper">
    <label v-if="label" :for="id" class="textarea-label">
      {{ label }}
      <span v-if="required" class="required">*</span>
    </label>
    <textarea
      :id="id"
      v-model="localValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :required="required"
      :rows="rows"
      :class="textareaClasses"
      @input="handleInput"
      @blur="handleBlur"
      @focus="handleFocus"
    ></textarea>
    <p v-if="error" class="textarea-error">{{ error }}</p>
    <p v-else-if="hint" class="textarea-hint">{{ hint }}</p>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
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
  rows: {
    type: Number,
    default: 4
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
    default: () => `textarea-${Math.random().toString(36).substr(2, 9)}`
  }
})

const emit = defineEmits(['update:modelValue', 'blur', 'focus'])

const localValue = ref(props.modelValue)
const isFocused = ref(false)

watch(
  () => props.modelValue,
  (newValue) => {
    localValue.value = newValue
  }
)

const textareaClasses = computed(() => ({
  'textarea-field': true,
  'textarea-error-state': !!props.error,
  'textarea-disabled': props.disabled,
  'textarea-focused': isFocused.value
}))

const handleInput = (event) => {
  emit('update:modelValue', event.target.value)
}

const handleBlur = (event) => {
  isFocused.value = false
  emit('blur', event)
}

const handleFocus = (event) => {
  isFocused.value = true
  emit('focus', event)
}
</script>

<style scoped>
.textarea-wrapper {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.textarea-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.required {
  color: #dc3545;
  margin-left: 2px;
}

.textarea-field {
  padding: 10px 14px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
  color: #1f2937;
  background: white;
  transition: all 0.2s;
  outline: none;
  resize: vertical;
  min-height: 80px;
}

.textarea-field:hover:not(.textarea-disabled) {
  border-color: #9ca3af;
}

.textarea-field:focus,
.textarea-focused {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.textarea-field::placeholder {
  color: #9ca3af;
}

.textarea-error-state {
  border-color: #dc3545;
}

.textarea-error-state:focus {
  border-color: #dc3545;
  box-shadow: 0 0 0 3px rgba(220, 53, 69, 0.1);
}

.textarea-disabled {
  background: #f3f4f6;
  color: #9ca3af;
  cursor: not-allowed;
  resize: none;
}

.textarea-error {
  font-size: 13px;
  color: #dc3545;
  margin: 0;
}

.textarea-hint {
  font-size: 13px;
  color: #6b7280;
  margin: 0;
}
</style>
