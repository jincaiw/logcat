<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<template>
  <div class="page-header">
    <div class="page-header-info">
      <h2>{{ title }}</h2>
      <p v-if="description" class="page-header-desc">{{ description }}</p>
    </div>
    <div class="page-header-actions">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 12px;
}

.page-header-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.page-header h2 {
  font-size: 22px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--text-color);
}

.page-header-desc {
  color: var(--text-color-secondary);
  font-size: 14px;
}

.page-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
