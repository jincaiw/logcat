<script setup lang="ts">
import { ref, onMounted, reactive, computed, h } from 'vue'
import dayjs from 'dayjs'
import { NDataTable, NButton, NForm, NFormItem, NInput, NInputNumber, NSelect, NDatePicker, NTag, NSpace, NEmpty, NTabs, NTabPane, NAlert, NSpin, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { API } from '@/api'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const message = useMessage()
const activeTab = ref('statistics')

// ==================== Statistics ====================
const statsLoading = ref(false)
const devices = ref<any[]>([])
const fieldStats = ref<any[]>([])
const summaryStats = ref({
  totalLogs: 0,
  uniqueCount: 0,
  field: '',
})

const statsQuery = reactive({
  deviceId: 0 as number | undefined,
  startDate: null as number | null,
  endDate: null as number | null,
})

async function loadDevices() {
  try {
    devices.value = await API.GetDevices()
  } catch (e) { console.error(e) }
}

async function loadFieldStats() {
  statsLoading.value = true
  try {
    const params: any = {
      startTime: statsQuery.startDate ? dayjs(statsQuery.startDate).format('YYYY-MM-DD') : dayjs().subtract(7, 'day').format('YYYY-MM-DD'),
      endTime: statsQuery.endDate ? dayjs(statsQuery.endDate).format('YYYY-MM-DD') : dayjs().format('YYYY-MM-DD'),
      field: '',
      topN: 20,
    }
    if (statsQuery.deviceId && statsQuery.deviceId > 0) params.deviceId = statsQuery.deviceId

    const result = await API.GetFieldStats(params)
    fieldStats.value = result.items || []
    summaryStats.value = {
      totalLogs: result.totalLogs || 0,
      uniqueCount: result.uniqueCount || 0,
      field: result.field || '',
    }
  } catch (e) {
    console.error(e)
  } finally {
    statsLoading.value = false
  }
}

function handleStatsSearch() {
  loadFieldStats()
}

function handleStatsReset() {
  statsQuery.deviceId = 0
  statsQuery.startDate = null
  statsQuery.endDate = null
  loadFieldStats()
}

const summaryCards = computed(() => [
  { label: t('service.totalLogs'), value: summaryStats.value.totalLogs.toLocaleString(), color: '#3b82f6' },
  { label: t('stats.uniqueCount'), value: summaryStats.value.uniqueCount.toLocaleString(), color: '#22c55e' },
  { label: t('stats.fieldName'), value: summaryStats.value.field || '-', color: '#8b5cf6' },
])

const fieldColumns: DataTableColumns<any> = [
  { title: t('stats.fieldValue'), key: 'value', width: 200, ellipsis: { tooltip: true } },
  { title: t('stats.location'), key: 'location', width: 160, ellipsis: { tooltip: true } },
  { title: t('stats.count'), key: 'count', width: 120, align: 'right', render(row) { return h('span', { style: { color: '#3b82f6', fontWeight: 600 } }, row.count.toLocaleString()) } },
  { title: t('stats.percentage'), key: 'percent', width: 120, align: 'right',
    render(row) {
      const pct = row.percent || '-'
      const pctNum = typeof pct === 'string' ? parseFloat(pct.replace('%', '')) : pct
      return h(NTag, { type: pctNum > 50 ? 'success' : pctNum > 20 ? 'info' : 'default', size: 'small' }, { default: () => String(pct) })
    },
  },
  { title: t('stats.lastSeen'), key: 'lastSeen', width: 170 },
]

const deviceOptions = computed(() => [
  { label: t('common.allDevices'), value: 0 },
  ...devices.value.map((d: any) => ({ label: d.name, value: d.id })),
])

// ==================== Test Tools ====================
const testSyslogForm = reactive({
  host: '127.0.0.1',
  port: 514,
  protocol: 'udp' as 'udp' | 'tcp',
  message: t('testTools.defaultMessage'),
  count: 1,
  intervalMs: 1000,
})

const sendTestLoading = ref(false)
const testResult = ref<any>(null)

async function handleSendTestSyslog() {
  if (!testSyslogForm.message) {
    message.warning(t('testTools.pleaseInputMessage'))
    return
  }
  sendTestLoading.value = true
  testResult.value = null
  try {
    const result = await API.SendTestSyslog({
      host: testSyslogForm.host,
      port: testSyslogForm.port,
      protocol: testSyslogForm.protocol,
      message: testSyslogForm.message,
      count: testSyslogForm.count,
      intervalMs: testSyslogForm.intervalMs,
    })
    testResult.value = { success: true, data: result }
    message.success(t('testTools.sendSuccess'))
  } catch (e: any) {
    testResult.value = { success: false, error: e.message || String(e) }
    message.error(t('testTools.sendFailed'))
  } finally {
    sendTestLoading.value = false
  }
}

// Log Trace
const traceForm = reactive({
  logId: '',
})

// Syslog Forward Test
const forwardForm = reactive({
  host: '127.0.0.1',
  port: 514,
  protocol: 'udp' as 'udp' | 'tcp',
  format: 'BSD',
})
const forwardLoading = ref(false)
const forwardResult = ref<any>(null)

async function handleTestForward() {
  if (!forwardForm.host) {
    message.warning(t('testTools.pleaseInputForwardHost'))
    return
  }
  forwardLoading.value = true
  forwardResult.value = null
  try {
    const result = await API.TestSyslogForward({
      host: forwardForm.host,
      port: forwardForm.port,
      protocol: forwardForm.protocol,
      format: forwardForm.format,
    })
    forwardResult.value = { success: true, data: result }
    message.success(t('testTools.forwardSuccess'))
  } catch (e: any) {
    forwardResult.value = { success: false, error: e.message || String(e) }
    message.error(t('testTools.forwardFailed'))
  } finally {
    forwardLoading.value = false
  }
}

// Log Trace
const traceLoading = ref(false)
const traceResult = ref<any[]>([])

async function handleTrace() {
  if (!traceForm.logId) {
    message.warning(t('testTools.pleaseInputLogId'))
    return
  }
  traceLoading.value = true
  traceResult.value = []
  try {
    const logId = Number(traceForm.logId)
    if (isNaN(logId)) {
      message.warning(t('testTools.pleaseInputValidLogId'))
      traceLoading.value = false
      return
    }
    const result = await API.GetLogTrace(logId)
    traceResult.value = [result]
  } catch (e: any) {
    message.error(t('testTools.traceFailed') + (e.message || ''))
  } finally {
    traceLoading.value = false
  }
}

function handleTraceReset() {
  traceForm.logId = ''
  traceResult.value = []
}

onMounted(async () => {
  await loadDevices()
  await loadFieldStats()
})
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('menu.stats') }}</h1>
        <p class="page-subtitle text-muted">{{ t('stats.subtitle') }}</p>
      </div>
    </div>

    <div class="glass-card mt-4">
      <NTabs v-model:value="activeTab" type="line" animated>
        <!-- Statistics Tab -->
        <NTabPane name="statistics" :tab="t('stats.statistics')">
          <!-- Summary Cards -->
          <div class="stats-grid mt-4">
            <div v-for="card in summaryCards" :key="card.label" class="stat-card glass-card">
              <div class="stat-label text-muted">{{ card.label }}</div>
              <div class="stat-value" :style="{ color: card.color }">{{ card.value }}</div>
            </div>
          </div>

          <!-- Query Form -->
          <div class="search-toolbar mt-6">
            <NSpace>
              <NSelect
                v-model:value="statsQuery.deviceId"
                :placeholder="t('log.selectDevice')"
                :options="deviceOptions"
                clearable
                style="width: 180px"
              />
              <NDatePicker
                v-model:value="statsQuery.startDate"
                type="date"
                :placeholder="t('stats.startDate')"
                clearable
                style="width: 160px"
              />
              <NDatePicker
                v-model:value="statsQuery.endDate"
                type="date"
                :placeholder="t('stats.endDate')"
                clearable
                style="width: 160px"
              />
              <NButton type="primary" @click="handleStatsSearch">{{ t('common.search') }}</NButton>
              <NButton @click="handleStatsReset">{{ t('common.reset') }}</NButton>
            </NSpace>
          </div>

          <!-- Field Stats Table -->
          <div class="data-table-wrap mt-4">
            <NSpin :show="statsLoading">
              <NDataTable
                :columns="fieldColumns"
                :data="fieldStats"
                :bordered="false"
                striped
              />
              <div v-if="!fieldStats.length && !statsLoading" class="empty-state">
                <NEmpty :description="t('stats.noDataDesc')" />
              </div>
            </NSpin>
          </div>
        </NTabPane>

        <!-- Test Tools Tab -->
        <NTabPane name="tools" :tab="t('testTools.title')">
          <!-- Send Test Syslog -->
          <div class="glass-card test-section">
            <div class="card-header">
              <h3 class="card-title text-accent">{{ t('testTools.sendTestSyslog') }}</h3>
            </div>
            <NForm :model="testSyslogForm" label-placement="left" :label-width="100" class="mt-4">
              <NFormItem :label="t('testTools.host')">
                <NInput v-model:value="testSyslogForm.host" :placeholder="t('testTools.hostPlaceholder')" style="width: 300px" />
              </NFormItem>
              <NFormItem :label="t('testTools.port')">
                <NInputNumber v-model:value="testSyslogForm.port" :min="1" :max="65535" style="width: 300px" />
              </NFormItem>
              <NFormItem :label="t('testTools.protocol')">
                <NSelect v-model:value="testSyslogForm.protocol" :options="[{ label: 'UDP', value: 'udp' }, { label: 'TCP', value: 'tcp' }]" style="width: 300px" />
              </NFormItem>
              <NFormItem :label="t('testTools.message')">
                <NInput v-model:value="testSyslogForm.message" type="textarea" :rows="3" :placeholder="t('testTools.messagePlaceholder')" />
              </NFormItem>
              <NFormItem :label="t('testTools.count')">
                <NInputNumber v-model:value="testSyslogForm.count" :min="1" :max="100" style="width: 300px" />
              </NFormItem>
              <NFormItem>
                <NButton type="primary" :loading="sendTestLoading" @click="handleSendTestSyslog">
                  {{ t('testTools.send') }}
                </NButton>
              </NFormItem>
            </NForm>

            <div v-if="testResult" class="test-result mt-4">
              <NAlert v-if="testResult.success" type="success" :bordered="false">
                <pre class="mono">{{ JSON.stringify(testResult.data, null, 2) }}</pre>
              </NAlert>
              <NAlert v-else type="error" :bordered="false">
                {{ testResult.error }}
              </NAlert>
            </div>
          </div>

          <!-- Syslog Forward Test -->
          <div class="glass-card test-section mt-6">
            <div class="card-header">
              <h3 class="card-title text-accent">{{ t('testTools.testForward') }}</h3>
            </div>
            <NForm :model="forwardForm" label-placement="left" :label-width="130" class="mt-4">
              <NFormItem :label="t('testTools.forwardHost')">
                <NInput v-model:value="forwardForm.host" :placeholder="t('testTools.forwardHostPlaceholder')" style="width: 300px" />
              </NFormItem>
              <NFormItem :label="t('testTools.forwardPort')">
                <NInputNumber v-model:value="forwardForm.port" :min="1" :max="65535" style="width: 300px" />
              </NFormItem>
              <NFormItem :label="t('testTools.protocol')">
                <NSelect v-model:value="forwardForm.protocol" :options="[{ label: 'UDP', value: 'udp' }, { label: 'TCP', value: 'tcp' }]" style="width: 300px" />
              </NFormItem>
              <NFormItem :label="t('testTools.forwardFormat')">
                <NInput v-model:value="forwardForm.format" :placeholder="t('testTools.forwardFormatPlaceholder')" style="width: 300px" />
              </NFormItem>
              <NFormItem>
                <NButton type="primary" :loading="forwardLoading" @click="handleTestForward">
                  {{ t('common.confirmButtonText') }}
                </NButton>
              </NFormItem>
            </NForm>

            <div v-if="forwardResult" class="test-result mt-4">
              <NAlert v-if="forwardResult.success" type="success" :bordered="false">
                <pre class="mono">{{ JSON.stringify(forwardResult.data, null, 2) }}</pre>
              </NAlert>
              <NAlert v-else type="error" :bordered="false">
                {{ forwardResult.error }}
              </NAlert>
            </div>
          </div>

          <!-- Log Trace -->
          <div class="glass-card test-section mt-6">
            <div class="card-header">
              <h3 class="card-title text-accent">{{ t('testTools.logTrace') }}</h3>
            </div>
            <NForm :model="traceForm" label-placement="left" :label-width="100" class="mt-4">
              <NFormItem :label="t('testTools.logId')">
                <NInput v-model:value="traceForm.logId" :placeholder="t('testTools.logIdPlaceholder')" />
              </NFormItem>
              <NFormItem>
                <NSpace>
                  <NButton type="primary" :loading="traceLoading" @click="handleTrace">
                    {{ t('testTools.trace') }}
                  </NButton>
                  <NButton @click="handleTraceReset">{{ t('common.reset') }}</NButton>
                </NSpace>
              </NFormItem>
            </NForm>

            <div v-if="traceResult.length" class="trace-results mt-4">
              <div v-for="(item, idx) in traceResult" :key="idx" class="trace-item glass-card">
                <pre class="mono log-entry">{{ typeof item === 'string' ? item : JSON.stringify(item, null, 2) }}</pre>
              </div>
            </div>
          </div>
        </NTabPane>
      </NTabs>
    </div>
  </div>
</template>

<style scoped>
.test-section {
  padding: 24px;
  border-radius: 12px;
}

.test-result pre,
.trace-item pre {
  margin: 0;
  font-size: 12px;
  color: var(--text-secondary);
}

.trace-item {
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 8px;
}
</style>
