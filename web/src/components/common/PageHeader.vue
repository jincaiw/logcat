<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<template>
  <div class="page-header">
    <div>
      <h2>{{ title }}</h2>
      <p v-if="description" style="color: var(--text-color-secondary); font-size: 14px; margin-top: 4px">
        {{ description }}
      </p>
    </div>
    <div class="page-header-actions">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.page-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>