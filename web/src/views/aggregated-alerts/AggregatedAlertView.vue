<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NTag, NModal, NDataTable, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getAggregatedAlerts, getAggregatedAlertLogs, acknowledgeAggregatedAlert } from '@/api/aggregatedAlerts'
import type { AggregatedAlert, SyslogLog } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'

const message = useMessage()
const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const detailShow = ref(false)
const detailLogs = ref<SyslogLog[]>([])
const detailLoading = ref(false)
const currentAgg = ref<AggregatedAlert | null>(null)

const columns: DataTableColumns<AggregatedAlert> = [
  { title: '时间', key: 'firstAt', width: 160 },
  { title: '规则', key: 'ruleName' },
  { title: '严重程度', key: 'severity', width: 80, render(row) { return h(StatusBadge, { status: row.severity, type: 'severity' }) } },
  { title: '摘要', key: 'summary', ellipsis: { tooltip: true } },
  { title: '数量', key: 'count', width: 60 },
  { title: '源IP', key: 'sourceIps', render(row: any) { return h('span', null, (row.sourceIps || []).join(', ')) } },
  {
    title: '操作', key: 'actions', width: 200,
    render(row) {
      return h(NSpace, { size: 'small' }, { default: () => [
        h(NButton, { size: 'small', text: true, type: 'primary', onClick: () => showLogs(row) }, { default: () => '关联日志' }),
        h(NButton, { size: 'small', text: true, onClick: () => handleAck(row) }, { default: () => '确认' }),
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
</script>

<template>
  <div class="page-container">
    <PageHeader title="聚合告警" description="查看聚合后的告警信息" />
    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['ruleName']" search-placeholder="搜索规则名称" />

    <n-modal v-model:show="detailShow" title="关联日志" preset="card" style="width: 900px">
      <n-data-table v-if="!detailLoading" :columns="[
        { title: '时间', key: 'receivedAt', width: 160 },
        { title: '设备', key: 'deviceName', width: 120 },
        { title: '源IP', key: 'sourceIp', width: 130 },
        { title: '消息', key: 'message', ellipsis: { tooltip: true } },
      ]" :data="detailLogs" :single-line="false" size="small" :pagination="false" />
      <div v-else style="text-align: center; padding: 40px">加载中...</div>
      <template #footer>
        <n-space justify="end"><n-button @click="detailShow = false">关闭</n-button></n-space>
      </template>
    </n-modal>
  </div>
</template>