<script setup lang="ts">
import { ref, h, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NSpace, NTag, NInput, NSelect, NCard, NDatePicker, NDrawer, NDrawerContent, NDescriptions, NDescriptionsItem, NInputNumber } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { queryLogs, cleanupLogs, getUnmatchedLogCount } from '@/api/logs'
import type { SyslogLog } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import { useAppMessage } from '@/composables/useMessage'
import { useTimeFormat } from '@/composables/useTimeFormat'
import { useIsMobile } from '@/composables/useIsMobile'

const message = useAppMessage()
const { formatTime } = useTimeFormat()
const { isMobile } = useIsMobile()
const router = useRouter()

const drawerWidth = computed(() => isMobile.value ? window.innerWidth * 0.9 : 520)
const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const confirmDialogShow = ref(false)
const confirmLoading = ref(false)
const cleanupDays = ref<number | null>(30)

// Unmatched log count (FR-167)
const unmatchedCount = ref<number | null>(null)
const unmatchedLoading = ref(false)

// Log detail drawer (FR-151-157)
const drawerShow = ref(false)
const selectedLog = ref<SyslogLog | null>(null)

const parsedDataObj = computed(() => {
  if (!selectedLog.value?.parsedData) return null
  try { return JSON.parse(selectedLog.value.parsedData) } catch { return null }
})

// Date range picker (replace startTime/endTime)
const dateRange = ref<[number, number] | null>(null)

const filters = ref({
  keyword: '',
  deviceName: '',
  sourceIp: '',
  destinationIp: '',
  eventType: '',
  severity: '',
  alertStatus: '',
  logId: '',
  startTime: '',
  endTime: '',
  parsedFieldKey: '',
  parsedFieldValue: '',
  filterStatus: '',
})

const columns: DataTableColumns<SyslogLog> = [
  { title: '时间', key: 'receivedAt', width: 160 },
  { title: '设备', key: 'deviceName', width: 120, ellipsis: { tooltip: true } },
  { title: '源IP', key: 'sourceIp', width: 130 },
  { title: '事件类型', key: 'eventType', width: 100 },
  { title: 'Facility', key: 'facility', width: 90 }, // FR-160
  {
    title: '严重程度', key: 'severity', width: 80,
    render(row) { return h(StatusBadge, { status: row.severity, type: 'severity' }) },
  },
  {
    title: '告警状态', key: 'alertStatus', width: 80,
    render(row) { return h(StatusBadge, { status: row.alertStatus, type: 'push' }) },
  },
  { title: '消息', key: 'rawMessage', ellipsis: { tooltip: true } },
  {
    title: '操作', key: 'actions', width: 140,
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'small', text: true, type: 'primary', onClick: (e: Event) => { e.stopPropagation(); showLogDetail(row) } }, { default: () => '详情' }),
          h(NButton, { size: 'small', text: true, type: 'primary', onClick: (e: Event) => { e.stopPropagation(); router.push(`/logs/trace/${row.id}`) } }, { default: () => '追踪' }),
        ],
      })
    },
  },
]

const extraParams = computed(() => {
  const params: Record<string, any> = { ...filters.value }
  // Convert date range to startTime/endTime strings
  if (dateRange.value) {
    params.startTime = new Date(dateRange.value[0]).toISOString()
    params.endTime = new Date(dateRange.value[1]).toISOString()
  } else {
    params.startTime = ''
    params.endTime = ''
  }
  return params
})

async function fetchData(params: any) {
  const merged = { ...params, ...extraParams.value }
  const res = await queryLogs(merged)
  return res
}

function handleSearch() {
  tableRef.value?.loadData()
}

function handleReset() {
  filters.value = {
    keyword: '', deviceName: '', sourceIp: '', destinationIp: '', eventType: '',
    severity: '', alertStatus: '', logId: '', startTime: '', endTime: '',
    parsedFieldKey: '', parsedFieldValue: '', filterStatus: '',
  }
  dateRange.value = null
  handleSearch()
}

// FR-151-157: Log detail drawer
function showLogDetail(row: SyslogLog) {
  selectedLog.value = row
  drawerShow.value = true
}

// FR-167: Load unmatched log count
async function loadUnmatchedCount() {
  unmatchedLoading.value = true
  try {
    const res = await getUnmatchedLogCount()
    unmatchedCount.value = res.data?.count ?? 0
  } catch {
    unmatchedCount.value = null
  } finally {
    unmatchedLoading.value = false
  }
}

// FR-168: Cleanup unmatched logs
const cleanupUnmatchedShow = ref(false)
const cleanupUnmatchedLoading = ref(false)

function handleCleanupUnmatched() {
  cleanupUnmatchedShow.value = true
}

async function doCleanupUnmatched() {
  cleanupUnmatchedLoading.value = true
  try {
    const res = await cleanupLogs(undefined, 0) // days=0 for unmatched only
    message.success(`清理完成，共删除 ${res.data?.deleted || 0} 条未匹配日志`)
    loadUnmatchedCount()
    tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '清理失败') }
  finally { cleanupUnmatchedLoading.value = false; cleanupUnmatchedShow.value = false }
}

// FR-165: Cleanup by days
function handleCleanup() {
  confirmDialogShow.value = true
}

async function doCleanup() {
  confirmLoading.value = true
  try {
    const days = cleanupDays.value ?? 30
    const res = await cleanupLogs(undefined, days)
    message.success(`清理完成，共删除 ${res.data?.deleted || 0} 条日志`)
    tableRef.value?.loadData()
    loadUnmatchedCount()
  } catch (err: any) { message.error(err?.message || '清理失败') }
  finally { confirmLoading.value = false; confirmDialogShow.value = false }
}

onMounted(() => {
  loadUnmatchedCount()
})
</script>

<template>
  <div class="page-container">
    <PageHeader title="日志查询" description="查询与检索系统日志">
      <n-space align="center">
        <n-tag v-if="unmatchedCount !== null" :type="unmatchedCount > 0 ? 'warning' : 'success'" size="small">
          未匹配日志: {{ unmatchedCount }}
        </n-tag>
        <n-button v-if="unmatchedCount !== null && unmatchedCount > 0" type="warning" size="small" ghost @click="handleCleanupUnmatched">
          清理未匹配
        </n-button>
        <n-button type="warning" ghost @click="handleCleanup">清理日志</n-button>
      </n-space>
    </PageHeader>

    <n-card size="small" style="margin-bottom: 16px">
      <n-space wrap>
        <n-input v-model:value="filters.keyword" placeholder="关键字" clearable style="min-width: 120px; flex: 1 1 160px" @keyup.enter="handleSearch" />
        <n-input v-model:value="filters.deviceName" placeholder="设备名称" clearable style="min-width: 120px; flex: 1 1 140px" />
        <n-input v-model:value="filters.sourceIp" placeholder="源IP" clearable style="min-width: 120px; flex: 1 1 140px" />
        <n-input v-model:value="filters.destinationIp" placeholder="目标IP" clearable style="min-width: 120px; flex: 1 1 140px" />
        <n-select v-model:value="filters.severity" placeholder="严重程度" clearable style="min-width: 100px; flex: 1 1 120px"
          :options="[
            { label: '严重', value: 'critical' }, { label: '高', value: 'high' },
            { label: '中', value: 'medium' }, { label: '低', value: 'low' }, { label: '信息', value: 'info' },
          ]" />
        <n-select v-model:value="filters.alertStatus" placeholder="告警状态" clearable style="min-width: 100px; flex: 1 1 120px"
          :options="[
            { label: '已触发', value: 'triggered' }, { label: '无', value: 'none' },
          ]" />
        <n-select v-model:value="filters.filterStatus" placeholder="过滤状态" clearable style="min-width: 100px; flex: 1 1 120px"
          :options="[
            { label: '已匹配', value: 'matched' }, { label: '未匹配', value: 'unmatched' },
          ]" />
        <n-input v-model:value="filters.logId" placeholder="日志ID" clearable style="min-width: 120px; flex: 1 1 200px" />
        <n-date-picker v-model:value="dateRange" type="datetimerange" clearable style="max-width: 360px; width: 100%" />
        <n-input v-model:value="filters.parsedFieldKey" placeholder="解析字段名" clearable style="min-width: 100px; flex: 1 1 120px" />
        <n-input v-model:value="filters.parsedFieldValue" placeholder="解析字段值" clearable style="min-width: 100px; flex: 1 1 120px" @keyup.enter="handleSearch" />
        <n-button type="primary" @click="handleSearch">搜索</n-button>
        <n-button @click="handleReset">重置</n-button>
      </n-space>
    </n-card>

    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :show-search="false" :extra-params="extraParams" @row-click="showLogDetail" />

    <!-- Log Detail Drawer (FR-151-157) -->
    <n-drawer v-model:show="drawerShow" :width="drawerWidth" placement="right">
      <n-drawer-content title="日志详情" closable>
        <template v-if="selectedLog">
          <n-descriptions label-placement="left" bordered :column="1" size="small">
            <n-descriptions-item label="ID">{{ selectedLog.id }}</n-descriptions-item>
            <n-descriptions-item label="接收时间">{{ formatTime(selectedLog.receivedAt) }}</n-descriptions-item>
            <n-descriptions-item label="设备名称">{{ selectedLog.deviceName }}</n-descriptions-item>
            <n-descriptions-item label="源IP">{{ selectedLog.sourceIp }}</n-descriptions-item>
            <n-descriptions-item label="目标IP">{{ selectedLog.destinationIp }}</n-descriptions-item>
            <n-descriptions-item label="Facility">{{ selectedLog.facility }}</n-descriptions-item>
            <n-descriptions-item label="严重程度">
              <StatusBadge :status="selectedLog.severity" type="severity" />
            </n-descriptions-item>
            <n-descriptions-item label="事件类型">{{ selectedLog.eventType }}</n-descriptions-item>
            <n-descriptions-item label="告警状态">
              <StatusBadge :status="selectedLog.alertStatus" type="push" />
            </n-descriptions-item>
            <n-descriptions-item label="过滤状态">
              <n-tag :type="selectedLog.filterStatus === 'matched' ? 'success' : 'warning'" size="small">
                {{ selectedLog.filterStatus === 'matched' ? '已匹配' : '未匹配' }}
              </n-tag>
            </n-descriptions-item>
            <n-descriptions-item label="匹配过滤策略ID">
              {{ selectedLog.matchedFilterPolicyId ?? '--' }}
            </n-descriptions-item>
            <n-descriptions-item label="告警状态">
              <n-tag :type="selectedLog.alertStatus === 'triggered' ? 'error' : 'default'" size="small">
                {{ selectedLog.alertStatus === 'triggered' ? '已触发' : '无' }}
              </n-tag>
            </n-descriptions-item>
            <n-descriptions-item label="告警规则ID">
              {{ selectedLog.alertRuleId ?? '--' }}
            </n-descriptions-item>
            <n-descriptions-item label="原始消息">
              <div style="white-space: pre-wrap; word-break: break-all; max-height: 200px; overflow-y: auto; font-size: 12px; font-family: monospace">
                {{ selectedLog.rawMessage || '--' }}
              </div>
            </n-descriptions-item>
            <n-descriptions-item label="解析数据">
              <div v-if="parsedDataObj && Object.keys(parsedDataObj).length > 0"
                style="max-height: 300px; overflow-y: auto; font-size: 12px; font-family: monospace">
                <div v-for="(value, key) in parsedDataObj" :key="key"
                  style="padding: 2px 0; border-bottom: 1px solid var(--border-color)">
                  <span style="color: var(--primary-color); font-weight: 600">{{ key }}</span>: {{ value }}
                </div>
              </div>
              <span v-else style="color: var(--text-color-secondary)">--</span>
            </n-descriptions-item>
          </n-descriptions>
        </template>
      </n-drawer-content>
    </n-drawer>

    <!-- Cleanup by days dialog (FR-165) -->
    <ConfirmDialog
      v-model:show="confirmDialogShow"
      title="清理日志"
      :loading="confirmLoading"
      @confirm="doCleanup"
    >
      <div style="display: flex; flex-direction: column; gap: 8px">
        <span>确定要清理旧日志吗？此操作不可恢复！</span>
        <n-space align="center">
          <span>清理</span>
          <n-input-number v-model:value="cleanupDays" :min="1" :max="365" size="small" style="width: 100px" />
          <span>天前的日志</span>
        </n-space>
      </div>
    </ConfirmDialog>

    <!-- Cleanup unmatched dialog (FR-168) -->
    <ConfirmDialog
      v-model:show="cleanupUnmatchedShow"
      title="清理未匹配日志"
      content="确定要清理所有未匹配的日志吗？此操作不可恢复！"
      :loading="cleanupUnmatchedLoading"
      @confirm="doCleanupUnmatched"
    />
  </div>
</template>
