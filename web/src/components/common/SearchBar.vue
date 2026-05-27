<script setup lang="ts">
import { ref } from 'vue'
import { NInput, NButton, NIcon, NSpace, NSelect } from 'naive-ui'
import { SearchOutline } from '@vicons/ionicons5'

export interface SearchField {
  key: string
  label: string
  type: 'text' | 'select' | 'date'
  placeholder?: string
  options?: { label: string; value: any }[]
}

const props = defineProps<{
  fields: SearchField[]
  showReset?: boolean
}>()

const emit = defineEmits<{
  (e: 'search', values: Record<string, any>): void
  (e: 'reset'): void
}>()

const searchValues = ref<Record<string, any>>({})

function handleSearch() {
  emit('search', { ...searchValues.value })
}

function handleReset() {
  searchValues.value = {}
  emit('reset')
}

function handleKeyup(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    handleSearch()
  }
}
</script>

<template>
  <div class="search-bar">
    <n-space align="center" wrap>
      <template v-for="field in fields" :key="field.key">
        <n-select
          v-if="field.type === 'select'"
          v-model:value="searchValues[field.key]"
          :options="field.options"
          :placeholder="field.placeholder || field.label"
          clearable
          style="width: 160px"
        />
        <n-input
          v-else
          v-model:value="searchValues[field.key]"
          :placeholder="field.placeholder || field.label"
          clearable
          style="width: 180px"
          @keyup="handleKeyup"
        />
      </template>
      <n-button type="primary" size="small" @click="handleSearch">
        <template #icon>
          <n-icon :component="SearchOutline" />
        </template>
        搜索
      </n-button>
      <n-button v-if="showReset" size="small" @click="handleReset">
        重置
      </n-button>
    </n-space>
  </div>
</template>

<style scoped>
.search-bar {
  margin-bottom: 16px;
}
</style>