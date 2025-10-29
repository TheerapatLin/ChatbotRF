<template>
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue'

onMounted(() => {
  console.log('App mounted')
  console.log('API Base URL:', import.meta.env.VITE_API_BASE_URL)
  console.log('WebSocket URL:', import.meta.env.VITE_WS_URL)
})
</script>
