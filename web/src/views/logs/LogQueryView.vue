<script setup lang="ts">
import { ref, h, computed } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NSpace, NTag, NInput, NSelect, NCard, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { queryLogs, cleanupLogs } from '@/api/logs'
import type { SyslogLog } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'

const message = useMessage()
const router = useRouter()
const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const confirmDialogShow = ref(false)
const confirmLoading = ref(false)

const filters = ref({
  keyword: '',
  deviceName: '',
  sourceIp: '',
  destIp: '',
  eventType: '',
  severity: '',
  pushStatus: '',
  logId: '',
  startTime: '',
  endTime: '',
})

const columns: DataTableColumns<SyslogLog> = [
  { title: '时间', key: 'receivedAt', width: 160 },
  { title: '设备', key: 'deviceName', width: 120, ellipsis: { tooltip: true } },
  { title: '源IP', key: 'sourceIp', width: 130 },
  { title: '事件类型', key: 'eventType', width: 100 },
  {
    title: '严重程度', key: 'severity', width: 80,
    render(row) { return h(StatusBadge, { status: row.severity, type: 'severity' }) },
  },
  {
    title: '推送状态', key: 'pushStatus', width: 80,
    render(row) { return h(StatusBadge, { status: row.pushStatus, type: 'push' }) },
  },
  { title: '消息', key: 'message', ellipsis: { tooltip: true } },
  {
    title: '操作', key: 'actions', width: 100,
    render(row) {
      return h(NButton, { size: 'small', text: true, type: 'primary', onClick: () => router.push(`/logs/trace/${row.id}`) }, { default: () => '追踪' })
    },
  },
]

const extraParams = computed(() => ({ ...filters.value }))

async function fetchData(params: any) {
  const merged = { ...params, ...filters.value }
  const res = await queryLogs(merged)
  return res
}

function handleSearch() {
  tableRef.value?.loadData()
}

function handleReset() {
  filters.value = { keyword: '', deviceName: '', sourceIp: '', destIp: '', eventType: '', severity: '', pushStatus: '', logId: '', startTime: '', endTime: '' }
  handleSearch()
}

function handleCleanup() {
  confirmDialogShow.value = true
}

async function doCleanup() {
  confirmLoading.value = true
  try {
    const res = await cleanupLogs()
    message.success(`清理完成，共删除 ${res.data?.deleted || 0} 条日志`)
    tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '清理失败') }
  finally { confirmLoading.value = false; confirmDialogShow.value = false }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="日志查询" description="查询与检索系统日志">
      <n-button type="warning" ghost @click="handleCleanup">清理日志</n-button>
    </PageHeader>

    <n-card size="small" style="margin-bottom: 16px">
      <n-space wrap>
        <n-input v-model:value="filters.keyword" placeholder="关键字" clearable style="width: 160px" @keyup.enter="handleSearch" />
        <n-input v-model:value="filters.deviceName" placeholder="设备名称" clearable style="width: 140px" />
        <n-input v-model:value="filters.sourceIp" placeholder="源IP" clearable style="width: 140px" />
        <n-input v-model:value="filters.destIp" placeholder="目标IP" clearable style="width: 140px" />
        <n-select v-model:value="filters.severity" placeholder="严重程度" clearable style="width: 120px"
          :options="[
            { label: '严重', value: 'critical' }, { label: '高', value: 'high' },
            { label: '中', value: 'medium' }, { label: '低', value: 'low' }, { label: '信息', value: 'info' },
          ]" />
        <n-select v-model:value="filters.pushStatus" placeholder="推送状态" clearable style="width: 120px"
          :options="[
            { label: '等待中', value: 'pending' }, { label: '成功', value: 'success' },
            { label: '失败', value: 'failed' }, { label: '跳过', value: 'skipped' },
          ]" />
        <n-input v-model:value="filters.logId" placeholder="日志ID" clearable style="width: 200px" />
        <n-input v-model:value="filters.startTime" placeholder="开始时间" clearable style="width: 160px" />
        <n-input v-model:value="filters.endTime" placeholder="结束时间" clearable style="width: 160px" />
        <n-button type="primary" @click="handleSearch">搜索</n-button>
        <n-button @click="handleReset">重置</n-button>
      </n-space>
    </n-card>

    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :show-search="false" :extra-params="extraParams" />

    <ConfirmDialog
      v-model:show="confirmDialogShow"
      title="清理日志"
      content="确定要清理旧日志吗？此操作不可恢复！"
      :loading="confirmLoading"
      @confirm="doCleanup"
    />
  </div>
</template>