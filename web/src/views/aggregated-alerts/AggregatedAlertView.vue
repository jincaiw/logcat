<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NModal, NDataTable } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getAggregatedAlerts, getAggregatedAlertLogs, acknowledgeAggregatedAlert, resolveAggregatedAlert } from '@/api/aggregatedAlerts'
import type { AggregatedAlert, SyslogLog } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import { useAppMessage } from '@/composables/useMessage'
import { useIsMobile } from '@/composables/useIsMobile'
import { useTimeFormat } from '@/composables/useTimeFormat'

const message = useAppMessage()
const { isMobile } = useIsMobile()
const { formatTime } = useTimeFormat()
const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const detailShow = ref(false)
const detailLogs = ref<SyslogLog[]>([])
const detailLoading = ref(false)
const currentAgg = ref<AggregatedAlert | null>(null)

function canAcknowledge(row: AggregatedAlert) {
  return !['acknowledged', 'resolved', 'closed'].includes(row.status)
}

function canResolve(row: AggregatedAlert) {
  return !['resolved', 'closed'].includes(row.status)
}

const columns: DataTableColumns<AggregatedAlert> = [
  { title: '时间', key: 'firstSeenAt', width: 160 },
  { title: '类型', key: 'eventType' },
  { title: '严重程度', key: 'severity', width: 80, render(row) { return h(StatusBadge, { status: row.severity, type: 'severity' }) } },
  { title: '状态', key: 'status', width: 90, render(row) { return h(StatusBadge, { status: row.status, type: 'alert' }) } },
  { title: '聚合键', key: 'aggregateKey', ellipsis: { tooltip: true } },
  { title: '数量', key: 'count', width: 60 },
  { title: '源IP', key: 'sourceIp' },
  {
    title: '操作', key: 'actions', width: 260,
    render(row) {
      return h(NSpace, { size: 'small' }, { default: () => [
        h(NButton, { size: 'small', text: true, type: 'primary', onClick: () => showLogs(row) }, { default: () => '关联日志' }),
        h(
          NButton,
          { size: 'small', text: true, disabled: !canAcknowledge(row), onClick: () => handleAck(row) },
          { default: () => '确认' },
        ),
        h(
          NButton,
          { size: 'small', text: true, type: 'warning', disabled: !canResolve(row), onClick: () => handleResolve(row) },
          { default: () => '解决' },
        ),
      ]})
    },
  },
]

async function fetchData(params: any) { return getAggregatedAlerts(params) }

async function showLogs(row: AggregatedAlert) {
  currentAgg.value = row
  detailShow.value = true
  detailLoading.value = true
  try {
    const res = await getAggregatedAlertLogs(row.id, { page: 1, pageSize: 50 })
    detailLogs.value = res.data?.list || []
  } catch { message.error('获取日志失败') }
  finally { detailLoading.value = false }
}

async function handleAck(row: AggregatedAlert) {
  try { await acknowledgeAggregatedAlert(row.id); message.success('已确认'); tableRef.value?.loadData() }
  catch (err: any) { message.error(err?.message || '操作失败') }
}

async function handleResolve(row: AggregatedAlert) {
  try { await resolveAggregatedAlert(row.id); message.success('已解决'); tableRef.value?.loadData() }
  catch (err: any) { message.error(err?.message || '操作失败') }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="聚合告警" description="查看聚合后的告警信息" />
    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['eventType', 'sourceIp']" search-placeholder="搜索事件类型或源IP" />

    <n-modal v-model:show="detailShow" title="关联日志" preset="card" :style="{ width: isMobile ? 'calc(100vw - 32px)' : '900px', maxWidth: 'calc(100vw - 32px)' }">
      <n-data-table v-if="!detailLoading" :columns="[
        { title: '时间', key: 'receivedAt', width: 160, render: (row: any) => formatTime(row.receivedAt) },
        { title: '设备', key: 'deviceName', width: 120 },
        { title: '源IP', key: 'sourceIp', width: 130 },
        { title: '消息', key: 'rawMessage', ellipsis: { tooltip: true } },
      ]" :data="detailLogs" :single-line="false" size="small" :pagination="false" />
      <div v-else style="text-align: center; padding: 40px">加载中...</div>
      <template #footer>
        <n-space justify="end"><n-button @click="detailShow = false">关闭</n-button></n-space>
      </template>
    </n-modal>
  </div>
</template>
