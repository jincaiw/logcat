<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NTag, NModal, NInput, NSelect, NCard } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getAlertRecords, disposeAlert } from '@/api/alerts'
import type { AlertRecord } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import { useAppMessage } from '@/composables/useMessage'
import { useIsMobile } from '@/composables/useIsMobile'

const message = useAppMessage()
const { isMobile } = useIsMobile()
const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const noteDialogShow = ref(false)
const noteLoading = ref(false)
const currentAlert = ref<AlertRecord | null>(null)
const noteValue = ref('')
const disposeAction = ref('confirm')

const columns: DataTableColumns<AlertRecord> = [
  { title: '时间', key: 'createdAt', width: 160 },
  { title: '日志ID', key: 'logId', width: 180 },
  { title: '通道类型', key: 'channelType', width: 80, render(row) { return h(NTag, { size: 'small', bordered: false }, { default: () => row.channelType?.toUpperCase() || '--' }) } },
  { title: '状态', key: 'status', width: 80, render(row) { return h(StatusBadge, { status: row.status, type: 'alert' }) } },
  { title: '请求摘要', key: 'requestSummary', ellipsis: { tooltip: true } },
  {
    title: '操作', key: 'actions', width: 240,
    render(row) {
      return h(NSpace, { size: 'small' }, { default: () => [
        h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => openNote(row, 'confirm') }, { default: () => '确认' }),
        h(NButton, { size: 'small', onClick: () => openNote(row, 'ignore') }, { default: () => '忽略' }),
        h(NButton, { size: 'small', type: 'warning', ghost: true, onClick: () => openNote(row, 'close') }, { default: () => '关闭' }),
      ]})
    },
  },
]

async function fetchData(params: any) { return getAlertRecords(params) }

function openNote(row: AlertRecord, action: string) {
  currentAlert.value = row
  disposeAction.value = action
  noteValue.value = ''
  noteDialogShow.value = true
}

async function handleDispose() {
  if (!currentAlert.value) return
  noteLoading.value = true
  try {
    await disposeAlert(currentAlert.value.id, { action: disposeAction.value, note: noteValue.value })
    message.success('处置成功')
    noteDialogShow.value = false
    tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '处置失败') }
  finally { noteLoading.value = false }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="告警记录" description="查看和处理告警记录" />
    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['channelType', 'requestSummary']" search-placeholder="搜索通道类型或请求摘要" />

    <!-- Disposition Note Dialog -->
    <n-modal v-model:show="noteDialogShow" title="处置告警" preset="card" :style="{ width: isMobile ? 'calc(100vw - 32px)' : '500px', maxWidth: 'calc(100vw - 32px)' }" :mask-closable="false">
      <n-space vertical>
        <div>操作: {{ disposeAction === 'confirm' ? '确认' : disposeAction === 'ignore' ? '忽略' : '关闭' }}</div>
        <n-input v-model:value="noteValue" type="textarea" placeholder="处置备注（可选）" :autosize="{ minRows: 3, maxRows: 6 }" />
      </n-space>
      <template #footer>
        <n-space justify="end">
          <n-button @click="noteDialogShow = false">取消</n-button>
          <n-button type="primary" :loading="noteLoading" @click="handleDispose">确认处置</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>