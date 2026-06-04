<script setup lang="ts">
import { ref, h, onMounted } from 'vue'
import { NButton, NSpace, NInputNumber, NCard, NGrid, NGi } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getHighFreqIps, getHighFreqIpConfig, updateHighFreqIpConfig } from '@/api/highFreqIps'
import type { HighFreqIp } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAppMessage } from '@/composables/useMessage'

const message = useAppMessage()
const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const configLoading = ref(false)
const timeWindow = ref(300)
const threshold = ref(100)

const columns: DataTableColumns<HighFreqIp> = [
  { title: 'IP 地址', key: 'ip', width: 160 },
  { title: '访问次数', key: 'count', width: 100 },
  { title: '首次出现', key: 'firstSeen', width: 160 },
  { title: '最后出现', key: 'lastSeen', width: 160 },
  { title: '关联设备', key: 'deviceNames', render(row: any) { return h('span', null, (row.deviceNames || []).join(', ')) } },
]

async function fetchData(params: any) { return getHighFreqIps(params) }

async function loadConfig() {
  try {
    const res = await getHighFreqIpConfig()
    if (res.data) {
      timeWindow.value = res.data.timeWindowSeconds || 300
      threshold.value = res.data.threshold || 100
    }
  } catch { /* ignore */ }
}

async function saveConfig() {
  configLoading.value = true
  try {
    await updateHighFreqIpConfig({ timeWindowSeconds: timeWindow.value, threshold: threshold.value })
    message.success('配置保存成功')
  } catch (err: any) { message.error(err?.message || '保存失败') }
  finally { configLoading.value = false }
}

onMounted(() => { loadConfig() })
</script>

<template>
  <div class="page-container">
    <PageHeader title="高频IP" description="识别和分析高频访问IP" />

    <n-card size="small" style="margin-bottom: 16px" title="检测配置">
      <n-grid cols="1 s:3" x-gap="12" y-gap="12" responsive="screen">
        <n-gi>
          <div style="font-size: 13px; margin-bottom: 4px; color: var(--text-color-secondary)">时间窗口 (秒)</div>
          <n-input-number v-model:value="timeWindow" :min="10" :max="3600" />
        </n-gi>
        <n-gi>
          <div style="font-size: 13px; margin-bottom: 4px; color: var(--text-color-secondary)">阈值 (次)</div>
          <n-input-number v-model:value="threshold" :min="1" :max="100000" />
        </n-gi>
        <n-gi style="display: flex; align-items: flex-end">
          <n-button type="primary" :loading="configLoading" @click="saveConfig">保存配置</n-button>
        </n-gi>
      </n-grid>
    </n-card>

    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" search-placeholder="搜索IP" />
  </div>
</template>