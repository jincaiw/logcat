<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NButton, NSpace, NSelect, NCard, NTable, NInputNumber, NDatePicker, NStatistic, NForm, NFormItem, NGrid, NGi } from 'naive-ui'
import { getStats, exportStatsCsv, copyIpList, getAvailableFields } from '@/api/stats'
import type { StatsResult, AvailableField } from '@/types'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAppMessage } from '@/composables/useMessage'
import { useTimeFormat } from '@/composables/useTimeFormat'
import { useIsMobile } from '@/composables/useIsMobile'

const message = useAppMessage()
const { formatTime } = useTimeFormat()
const { isMobile } = useIsMobile()
const loading = ref(false)
const results = ref<StatsResult[]>([])
const totalLogs = ref(0)
const uniqueValues = ref(0)
const csvLoading = ref(false)
const copyLoading = ref(false)
const fieldOptions = ref<AvailableField[]>([])
const timeRange = ref<[number, number] | null>(null)

const FALLBACK_FIELDS: AvailableField[] = [
  { value: 'sourceIp', label: '源IP' },
  { value: 'deviceName', label: '设备' },
  { value: 'eventType', label: '事件类型' },
  { value: 'severity', label: '严重程度' },
]

const formData = ref({
  policy: '',
  startTime: '',
  endTime: '',
  field: '',
  topN: 20,
})

function toRfc3339(ts: number): string {
  return new Date(ts).toISOString()
}

function handleTimeRangeChange(value: [number, number] | null) {
  timeRange.value = value
  if (value) {
    formData.value.startTime = toRfc3339(value[0])
    formData.value.endTime = toRfc3339(value[1])
  } else {
    formData.value.startTime = ''
    formData.value.endTime = ''
  }
}

async function fetchFields() {
  try {
    const res = await getAvailableFields()
    fieldOptions.value = res.data?.length ? res.data : FALLBACK_FIELDS
  } catch {
    fieldOptions.value = FALLBACK_FIELDS
  }
}

async function handleQuery() {
  if (!formData.value.field) { message.warning('请选择统计字段'); return }
  if (!formData.value.startTime || !formData.value.endTime) {
    message.warning('请选择时间范围')
    return
  }
  loading.value = true
  try {
    const res = await getStats(formData.value as any)
    const data = res.data as any
    if (data && typeof data === 'object' && Array.isArray(data.results)) {
      results.value = data.results
      totalLogs.value = data.totalLogs ?? 0
      uniqueValues.value = data.uniqueValues ?? 0
    } else if (Array.isArray(data)) {
      results.value = data
      totalLogs.value = 0
      uniqueValues.value = 0
    } else {
      results.value = []
      totalLogs.value = 0
      uniqueValues.value = 0
    }
  } catch (err: any) { message.error(err?.message || '查询失败') }
  finally { loading.value = false }
}

async function handleExportCsv() {
  if (!formData.value.startTime || !formData.value.endTime) {
    message.warning('请选择时间范围')
    return
  }
  csvLoading.value = true
  try {
    const res = await exportStatsCsv(formData.value as any)
    if (res.data?.url) window.open(res.data.url)
    else message.success('导出成功')
  } catch (err: any) { message.error(err?.message || '导出失败') }
  finally { csvLoading.value = false }
}

async function handleCopyIpList() {
  if (!formData.value.startTime || !formData.value.endTime) {
    message.warning('请选择时间范围')
    return
  }
  copyLoading.value = true
  try {
    const res = await copyIpList({ startTime: formData.value.startTime, endTime: formData.value.endTime })
    const ips = res.data?.ips || []
    if (ips.length > 0) {
      await navigator.clipboard.writeText(ips.join('\n'))
      message.success(`已复制 ${ips.length} 个IP`)
    } else { message.warning('无IP数据') }
  } catch { message.error('操作失败') }
  finally { copyLoading.value = false }
}

onMounted(() => {
  fetchFields()
})
</script>

<template>
  <div class="page-container">
    <PageHeader title="数据统计" description="日志数据统计与分析" />

    <n-card size="small" style="margin-bottom: 16px">
      <n-form
        label-placement="left"
        :show-feedback="false"
        :label-width="isMobile ? 0 : 80"
      >
        <n-grid :cols="isMobile ? 1 : 4" :x-gap="12" :y-gap="12" responsive="screen">
          <n-gi>
            <n-form-item label="策略">
              <n-select
                v-model:value="formData.policy"
                placeholder="策略 (可选)"
                clearable
                :options="[{ label: '默认', value: 'default' }]"
              />
            </n-form-item>
          </n-gi>
          <n-gi :span="isMobile ? 1 : 2">
            <n-form-item label="时间范围">
              <n-date-picker
                v-model:value="timeRange"
                type="datetimerange"
                clearable
                placeholder="选择时间范围"
                style="width: 100%"
                @update:value="handleTimeRangeChange"
              />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="统计字段">
              <n-select
                v-model:value="formData.field"
                placeholder="请选择"
                :options="fieldOptions"
              />
            </n-form-item>
          </n-gi>
        </n-grid>

        <n-grid :cols="isMobile ? 1 : 4" :x-gap="12" :y-gap="12" responsive="screen" style="margin-top: 4px">
          <n-gi>
            <n-form-item label="Top N">
              <n-input-number
                v-model:value="formData.topN"
                :min="1"
                :max="1000"
                style="width: 100%"
              />
            </n-form-item>
          </n-gi>
          <n-gi :span="3">
            <n-space align="center" style="height: 100%; padding-top: 2px" :wrap="isMobile">
              <n-button type="primary" :loading="loading" @click="handleQuery">查询统计</n-button>
              <n-button :loading="csvLoading" @click="handleExportCsv">导出 CSV</n-button>
              <n-button :loading="copyLoading" @click="handleCopyIpList">复制 IP 列表</n-button>
            </n-space>
          </n-gi>
        </n-grid>
      </n-form>
    </n-card>

    <n-card size="small" title="统计结果">
      <n-space v-if="results.length" justify="start" wrap style="margin-bottom: 16px">
        <n-statistic label="总日志数" :value="totalLogs" />
        <n-statistic label="唯一值数量" :value="uniqueValues" />
      </n-space>
      <div v-if="results.length" style="overflow-x: auto">
        <n-table :single-line="false" size="small">
        <thead>
          <tr>
            <th>排名</th>
            <th>{{ formData.field || '字段' }}</th>
            <th>数量</th>
            <th>占比</th>
            <th>最后出现时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in results" :key="index">
            <td>{{ index + 1 }}</td>
            <td>{{ item.value }}</td>
            <td>{{ item.count }}</td>
            <td>{{ (item.percentage * 100).toFixed(2) }}%</td>
            <td>{{ formatTime(item.lastSeenAt) }}</td>
          </tr>
        </tbody>
      </n-table>
      </div>
      <div v-else style="text-align: center; padding: 40px; color: var(--text-color-secondary)">
        暂无数据，请先查询
      </div>
    </n-card>
  </div>
</template>
