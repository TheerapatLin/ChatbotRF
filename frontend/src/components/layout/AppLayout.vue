<template>
  <div class="flex h-screen bg-gray-50">
    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col">
      <slot></slot>
    </div>

    <!-- Right Sidebar -->
    <AppSidebar />
  </div>
</template>

<script setup>
import AppSidebar from './AppSidebar.vue'
</script>
