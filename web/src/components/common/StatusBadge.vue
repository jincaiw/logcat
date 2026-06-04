<script setup lang="ts">
import { computed } from 'vue'
import { NTag } from 'naive-ui'

const props = defineProps<{
  status: number | string
  type?: 'status' | 'severity' | 'push' | 'alert'
}>()

const config = computed(() => {
  if (props.type === 'severity') {
    const map: Record<string, any> = {
      critical: { type: 'error', label: '严重' },
      high: { type: 'warning', label: '高' },
      medium: { type: 'info', label: '中' },
      low: { type: 'success', label: '低' },
      info: { type: 'default', label: '信息' },
    }
    return map[String(props.status)] || { type: 'default', label: props.status }
  }
  if (props.type === 'push') {
    const map: Record<string, any> = {
      pending: { type: 'default', label: '等待中' },
      success: { type: 'success', label: '成功' },
      failed: { type: 'error', label: '失败' },
      skipped: { type: 'warning', label: '跳过' },
    }
    return map[String(props.status)] || { type: 'default', label: props.status }
  }
  if (props.type === 'alert') {
    const map: Record<string, any> = {
      pending: { type: 'default', label: '待处理' },
      success: { type: 'success', label: '成功' },
      sent: { type: 'warning', label: '已发送' },
      confirmed: { type: 'info', label: '已确认' },
      ignored: { type: 'default', label: '已忽略' },
      closed: { type: 'warning', label: '已关闭' },
      acknowledged: { type: 'info', label: '已确认' },
      resolved: { type: 'success', label: '已解决' },
      failed: { type: 'error', label: '失败' },
      active: { type: 'warning', label: '活跃' },
    }
    return map[String(props.status)] || { type: 'default', label: props.status }
  }
  // status type
  const map: Record<string, any> = {
    '1': { type: 'success', label: '启用' },
    '0': { type: 'default', label: '禁用' },
    online: { type: 'success', label: '在线' },
    offline: { type: 'default', label: '离线' },
    running: { type: 'success', label: '运行中' },
    stopped: { type: 'default', label: '已停止' },
  }
  return map[String(props.status)] || { type: 'default', label: props.status }
})
</script>

<template>
  <n-tag :type="config.type" size="small" :bordered="false">
    {{ config.label }}
  </n-tag>
</template>
