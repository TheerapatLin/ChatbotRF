<template>
  <div class="w-80 bg-white border-l border-gray-200 flex flex-col h-full overflow-y-auto">
    <!-- Header -->
    <div class="p-4 border-b border-gray-200">
      <h2 class="text-lg font-semibold text-gray-800">โหมดการสนทนา</h2>
    </div>

    <!-- Mode Selection -->
    <div class="p-4 space-y-3">
      <!-- Text to Text Mode -->
      <div>
        <button
          @click="selectMode('text')"
          :class="[
            'w-full p-4 rounded-lg border-2 transition-all text-left',
            chatMode === 'text'
              ? 'border-blue-500 bg-blue-50'
              : 'border-gray-200 hover:border-gray-300'
          ]"
        >
          <div class="flex items-center gap-3">
            <div class="text-3xl">💬</div>
            <div class="flex-1">
              <div class="font-semibold text-gray-800">Text to Text</div>
              <div class="text-sm text-gray-500">พิมพ์ข้อความเพื่อสนทนา</div>
            </div>
            <div v-if="chatMode === 'text'" class="text-blue-500">
              <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
              </svg>
            </div>
          </div>
        </button>

        <!-- Text Mode Settings -->
        <div v-if="chatMode === 'text'" class="mt-3 p-3 bg-blue-50 rounded-lg space-y-3">
          <!-- Persona Selector -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              เลือก Persona
            </label>
            <select
              v-model="selectedPersonaId"
              @change="handlePersonaChange"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500 text-sm"
            >
              <option v-if="personasLoading" disabled>กำลังโหลด...</option>
              <option
                v-for="persona in personas"
                :key="persona.id"
                :value="persona.id"
              >
                {{ persona.icon }} {{ persona.name }}
              </option>
            </select>
          </div>

          <!-- System Prompt -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              System Prompt (ถ้ามี)
            </label>
            <textarea
              v-model="systemPrompt"
              @input="handleSystemPromptChange"
              placeholder="ระบุ system prompt เพิ่มเติม..."
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500 text-sm resize-none"
              rows="3"
            />
          </div>

          <!-- Persona Info -->
          <div v-if="currentPersona" class="text-xs text-gray-600 bg-white p-2 rounded">
            <div class="font-medium mb-1">{{ currentPersona.name }}</div>
            <div class="text-gray-500">{{ currentPersona.description }}</div>
          </div>
        </div>
      </div>

      <!-- Speech to Speech Mode -->
      <button
        @click="selectMode('speech')"
        :class="[
          'w-full p-4 rounded-lg border-2 transition-all text-left',
          chatMode === 'speech'
            ? 'border-blue-500 bg-blue-50'
            : 'border-gray-200 hover:border-gray-300'
        ]"
      >
        <div class="flex items-center gap-3">
          <div class="text-3xl">🎤</div>
          <div class="flex-1">
            <div class="font-semibold text-gray-800">Speech to Speech</div>
            <div class="text-sm text-gray-500">พูดเพื่อสนทนา</div>
          </div>
          <div v-if="chatMode === 'speech'" class="text-blue-500">
            <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
            </svg>
          </div>
        </div>
      </button>
    </div>

    <!-- Footer Info -->
    <div class="mt-auto p-4 border-t border-gray-200 bg-gray-50">
      <div class="text-xs text-gray-500 space-y-1">
        <div class="flex justify-between">
          <span>ข้อความ:</span>
          <span class="font-semibold">{{ messageCount }}</span>
        </div>
        <div class="flex justify-between">
          <span>โหมด:</span>
          <span class="font-semibold">{{ chatMode === 'text' ? 'ข้อความ' : 'เสียง' }}</span>
        </div>
        <div v-if="chatMode === 'text' && currentPersona" class="flex justify-between">
          <span>Persona:</span>
          <span class="font-semibold">{{ currentPersona.name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useChatStore } from '@/store/chat'
import { usePersonasStore } from '@/store/personas'

const chatStore = useChatStore()
const personasStore = usePersonasStore()

const chatMode = computed(() => chatStore.chatMode)
const messageCount = computed(() => chatStore.messages.length)
const personas = computed(() => personasStore.personas)
const personasLoading = computed(() => personasStore.isLoading)
const currentPersona = computed(() => personasStore.getCurrentPersona)

const selectedPersonaId = ref(1)
const systemPrompt = ref('')

const selectMode = (mode) => {
  chatStore.setChatMode(mode)
}

const handlePersonaChange = () => {
  personasStore.setCurrentPersona(selectedPersonaId.value)
  chatStore.setCurrentPersona(selectedPersonaId.value)
}

const handleSystemPromptChange = () => {
  chatStore.setSystemPrompt(systemPrompt.value)
}

onMounted(async () => {
  try {
    await personasStore.fetchPersonas()
    selectedPersonaId.value = personasStore.currentPersonaId
  } catch (error) {
    console.error('Failed to fetch personas:', error)
  }
})
</script>
