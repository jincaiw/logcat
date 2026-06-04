<script setup lang="ts">
import { NModal, NSpace, NButton } from 'naive-ui'

defineProps<{
  title?: string
  content?: string
  type?: 'warning' | 'error' | 'info'
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'confirm'): void
  (e: 'cancel'): void
  (e: 'update:show', value: boolean): void
}>()

const show = defineModel<boolean>('show', { default: false })

function handleConfirm() {
  emit('confirm')
}

function handleCancel() {
  show.value = false
  emit('cancel')
}
</script>

<template>
  <n-modal
    v-model:show="show"
    :title="title || '确认操作'"
    preset="dialog"
    type="warning"
    positive-text="确认"
    negative-text="取消"
    :loading="loading"
    @positive-click="handleConfirm"
    @negative-click="handleCancel"
    @close="handleCancel"
  >
    <div v-if="$slots.default"><slot /></div>
    <div v-else>{{ content || '确定要执行此操作吗？' }}</div>
  </n-modal>
</template>