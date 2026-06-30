<script setup lang="ts">
import { ref, computed, onMounted, reactive, h } from 'vue'
import dayjs from 'dayjs'
import { NDataTable, NInput, NSelect, NDatePicker, NButton, NTag, NPagination, NModal, NDescriptions, NDescriptionsItem, NEmpty } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { API } from '@/api'
import { useI18n } from '@/i18n'


const { t } = useI18n()

interface LogItem {
  id: number
  deviceName: string
  sourceIp: string
  rawMessage: string
  parsedData: string
  parsedFields: string
  priority: string
  severity: number
  receivedAt: string
  isAlerted: boolean
  filterStatus: string
}

const loading = ref(false)
const logs = ref<LogItem[]>([])
const total = ref(0)
const dialogVisible = ref(false)
const currentLog = ref<LogItem | null>(null)

const queryParams = reactive({
  page: 1,
  pageSize: 20,
  deviceId: undefined as number | undefined,
  startTime: null as number | null,
  endTime: null as number | null,
  keyword: '',
})

const devices = ref<any[]>([])

onMounted(async () => {
  await loadDevices()
  await loadLogs()
})

async function loadDevices() {
  try {
    devices.value = await API.GetDevices()
  } catch (e) {
    console.error(e)
  }
}

async function loadLogs() {
  loading.value = true
  try {
    const result = await API.QueryLogs({
      page: queryParams.page,
      pageSize: queryParams.pageSize,
      deviceId: queryParams.deviceId || 0,
      startTime: queryParams.startTime ? dayjs(queryParams.startTime).format('YYYY-MM-DD HH:mm:ss') : '',
      endTime: queryParams.endTime ? dayjs(queryParams.endTime).format('YYYY-MM-DD HH:mm:ss') : '',
      keyword: queryParams.keyword,
    })
    logs.value = result.logs || []
    total.value = result.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  queryParams.page = 1
  loadLogs()
}

function handleReset() {
  queryParams.page = 1
  queryParams.deviceId = undefined
  queryParams.startTime = null
  queryParams.endTime = null
  queryParams.keyword = ''
  loadLogs()
}

function handlePageChange(page: number) {
  queryParams.page = page
  loadLogs()
}

function handleSizeChange(size: number) {
  queryParams.pageSize = size
  queryParams.page = 1
  loadLogs()
}

function viewLogDetail(log: LogItem) {
  currentLog.value = log
  dialogVisible.value = true
}

function getLevelTag(log: LogItem): { type: 'default' | 'error' | 'warning' | 'info' | 'success'; text: string } {
  if (log.parsedFields) {
    try {
      const fields = JSON.parse(log.parsedFields)
      if (fields.severity) {
        return getLevelTagByText(fields.severity)
      }
      if (fields.levelDesc) {
        return getLevelTagByText(fields.levelDesc)
      }
    } catch (e) {
      console.error(e)
    }
  }
  return { type: 'info', text: t('log.levelInfo') }
}

function getLevelTagByText(level: string): { type: 'default' | 'error' | 'warning' | 'info' | 'success'; text: string } {
  const levelMap: Record<string, { type: 'default' | 'error' | 'warning' | 'info' | 'success'; text: string }> = {
    '危急': { type: 'error', text: t('log.levelCritical') },
    '高危': { type: 'error', text: t('log.levelHigh') },
    '中危': { type: 'warning', text: t('log.levelMedium') },
    '低危': { type: 'info', text: t('log.levelLow') },
    '信息': { type: 'default', text: t('log.levelInfo') },
    '严重': { type: 'error', text: t('log.levelSevere') },
  }
  return levelMap[level] || { type: 'info', text: level }
}

function getFilterStatusType(status: string): 'default' | 'error' | 'warning' | 'info' | 'success' {
  const statusMap: Record<string, 'default' | 'error' | 'warning' | 'info' | 'success'> = {
    pending: 'warning',
    matched: 'success',
    unmatched: 'info',
    whitelisted: 'warning',
    discarded: 'error',
    disabled: 'info',
  }
  return statusMap[status] || 'info'
}

function getFilterStatusText(status: string): string {
  const statusMap: Record<string, string> = {
    pending: t('log.statusPending'),
    matched: t('log.statusMatched'),
    unmatched: t('log.statusUnmatched'),
    whitelisted: t('log.statusWhitelisted'),
    discarded: t('log.statusDiscarded'),
    disabled: t('log.statusDisabled'),
  }
  return statusMap[status] || status
}

const columns: DataTableColumns<LogItem> = [
  { title: t('log.id'), key: 'id', width: 80 },
  { title: t('log.deviceName'), key: 'deviceName', width: 120, ellipsis: { tooltip: true } },
  { title: t('log.sourceIp'), key: 'sourceIp', width: 140 },
  {
    title: t('log.severity'),
    key: 'severity',
    width: 90,
    render(row) {
      const tag = getLevelTag(row)
      return h(NTag, { type: tag.type, size: 'small' }, { default: () => tag.text })
    },
  },
  { title: t('log.logContent'), key: 'rawMessage', minWidth: 200, ellipsis: { tooltip: true } },
  {
    title: t('log.receivedAt'),
    key: 'receivedAt',
    width: 170,
    render(row) {
      return dayjs(row.receivedAt).format('YYYY-MM-DD HH:mm:ss')
    },
  },
  {
    title: t('log.filterStatus'),
    key: 'filterStatus',
    width: 90,
    render(row) {
      return h(NTag, { type: getFilterStatusType(row.filterStatus), size: 'small' }, { default: () => getFilterStatusText(row.filterStatus) })
    },
  },
  {
    title: t('log.alertStatus'),
    key: 'isAlerted',
    width: 70,
    align: 'center',
    render(row) {
      return h('span', {
        style: { color: row.isAlerted ? 'var(--success)' : 'var(--text-muted)', fontSize: '16px' },
      }, row.isAlerted ? '✓' : '—')
    },
  },
  {
    title: t('log.action'),
    key: 'actions',
    width: 60,
    align: 'center',
    render(row) {
      return h(NButton, {
        text: true,
        type: 'primary',
        size: 'small',
        onClick: () => viewLogDetail(row),
      }, { default: () => t('log.detail') })
    },
  },
]

const deviceOptions = computed(() => [
  { label: t('log.allDevices'), value: 0 },
  ...devices.value.map((d: any) => ({ label: d.name, value: d.id })),
])

</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('log.title') }}</h1>
        <p class="page-subtitle text-muted">{{ t('log.subtitle') }}</p>
      </div>
    </div>

    <div class="search-toolbar">
      <NInput
        v-model:value="queryParams.keyword"
        :placeholder="t('log.searchPlaceholder')"
        clearable
        style="width: 220px"
        @keyup.enter="handleSearch"
      />
      <NSelect
        v-model:value="queryParams.deviceId"
        :placeholder="t('log.selectDevice')"
        :options="deviceOptions"
        clearable
        style="width: 180px"
      />
      <NDatePicker
        v-model:value="queryParams.startTime"
        type="datetime"
        :placeholder="t('log.startTime')"
        clearable
        style="width: 200px"
      />
      <NDatePicker
        v-model:value="queryParams.endTime"
        type="datetime"
        :placeholder="t('log.endTime')"
        clearable
        style="width: 200px"
      />
      <NButton type="primary" @click="handleSearch">{{ t('common.search') }}</NButton>
      <NButton @click="handleReset">{{ t('common.reset') }}</NButton>
    </div>

    <div class="data-table-wrap mt-4">
      <NDataTable
        :columns="columns"
        :data="logs"
        :loading="loading"
        :row-props="(row: LogItem) => ({ style: 'cursor: pointer', ondblclick: () => viewLogDetail(row) })"
        :bordered="false"
        striped
      />
      <div v-if="!logs.length && !loading" class="empty-state">
        <NEmpty :description="t('log.noLogsDesc')" />
      </div>
    </div>

    <div class="pagination-wrap mt-4">
      <NPagination
        v-model:page="queryParams.page"
        v-model:page-size="queryParams.pageSize"
        :item-count="total"
        :page-sizes="[10, 20, 50, 100]"
        show-size-picker
        @update:page="handlePageChange"
        @update:page-size="handleSizeChange"
      />
    </div>

    <!-- Detail Modal -->
    <NModal
      v-model:show="dialogVisible"
      :title="t('log.logDetail')"
      preset="card"
      style="width: 700px"
      :bordered="true"
    >
      <NDescriptions v-if="currentLog" :column="2" bordered label-placement="left">
        <NDescriptionsItem :label="t('log.id')">{{ currentLog.id }}</NDescriptionsItem>
        <NDescriptionsItem :label="t('log.deviceName')">{{ currentLog.deviceName }}</NDescriptionsItem>
        <NDescriptionsItem :label="t('log.sourceIp')">{{ currentLog.sourceIp }}</NDescriptionsItem>
        <NDescriptionsItem :label="t('log.priority')">{{ currentLog.priority }}</NDescriptionsItem>
        <NDescriptionsItem :label="t('log.severity')">
          <NTag :type="getLevelTag(currentLog).type" size="small">{{ getLevelTag(currentLog).text }}</NTag>
        </NDescriptionsItem>
        <NDescriptionsItem :label="t('log.receivedAt')">{{ dayjs(currentLog.receivedAt).format('YYYY-MM-DD HH:mm:ss') }}</NDescriptionsItem>
        <NDescriptionsItem :label="t('log.rawLog')" :span="2">
          <pre class="log-content">{{ currentLog.rawMessage }}</pre>
        </NDescriptionsItem>
        <NDescriptionsItem :label="t('log.parsedData')" :span="2">
          <pre class="log-content">{{ currentLog.parsedData || '-' }}</pre>
        </NDescriptionsItem>
      </NDescriptions>
    </NModal>
  </div>
</template>

<style scoped>
.log-content {
  background: var(--bg-sunken);
  padding: 8px 12px;
  border-radius: 8px;
  font-family: 'SF Mono', 'Menlo', 'Monaco', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 240px;
  overflow: auto;
  margin: 0;
  color: var(--text-secondary);
}
</style>
