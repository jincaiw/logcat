<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NSpace, NSelect, NInput, NCard, NTable, NInputNumber, useMessage } from 'naive-ui'
import { getStats, exportStatsCsv, copyIpList } from '@/api/stats'
import type { StatsResult } from '@/types'
import PageHeader from '@/components/common/PageHeader.vue'

const message = useMessage()
const loading = ref(false)
const results = ref<StatsResult[]>([])
const csvLoading = ref(false)
const copyLoading = ref(false)

const formData = ref({
  policy: '',
  startTime: '',
  endTime: '',
  field: '',
  topN: 20,
})

async function handleQuery() {
  if (!formData.value.field) { message.warning('请选择统计字段'); return }
  loading.value = true
  try {
    const res = await getStats(formData.value as any)
    results.value = res.data || []
  } catch (err: any) { message.error(err?.message || '查询失败') }
  finally { loading.value = false }
}

async function handleExportCsv() {
  csvLoading.value = true
  try {
    const res = await exportStatsCsv(formData.value as any)
    if (res.data?.url) window.open(res.data.url)
    else message.success('导出成功')
  } catch (err: any) { message.error(err?.message || '导出失败') }
  finally { csvLoading.value = false }
}

async function handleCopyIpList() {
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
</script>

<template>
  <div class="page-container">
    <PageHeader title="数据统计" description="日志数据统计与分析" />

    <n-card size="small" style="margin-bottom: 16px">
      <n-space wrap>
        <n-select v-model:value="formData.policy" placeholder="策略 (可选)" clearable style="width: 150px"
          :options="[{ label: '默认', value: 'default' }]" />
        <n-input v-model:value="formData.startTime" placeholder="开始时间" clearable style="width: 160px" />
        <n-input v-model:value="formData.endTime" placeholder="结束时间" clearable style="width: 160px" />
        <n-select v-model:value="formData.field" placeholder="统计字段" style="width: 150px"
          :options="[
            { label: '源IP', value: 'sourceIp' }, { label: '设备', value: 'deviceName' },
            { label: '事件类型', value: 'eventType' }, { label: '严重程度', value: 'severity' },
          ]" />
        <div style="font-size: 13px; color: var(--text-color-secondary); margin-bottom: 4px">
          Top N:
          <n-input-number v-model:value="formData.topN" :min="1" :max="1000" style="width: 100px" />
        </div>
      </n-space>
      <n-space style="margin-top: 12px">
        <n-button type="primary" :loading="loading" @click="handleQuery">查询统计</n-button>
        <n-button :loading="csvLoading" @click="handleExportCsv">导出 CSV</n-button>
        <n-button :loading="copyLoading" @click="handleCopyIpList">复制 IP 列表</n-button>
      </n-space>
    </n-card>

    <n-card size="small" title="统计结果">
      <n-table v-if="results.length" :single-line="false" size="small">
        <thead>
          <tr>
            <th>排名</th>
            <th>{{ formData.field || '字段' }}</th>
            <th>数量</th>
            <th>占比</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in results" :key="index">
            <td>{{ index + 1 }}</td>
            <td>{{ item.value }}</td>
            <td>{{ item.count }}</td>
            <td>{{ (item.percentage * 100).toFixed(2) }}%</td>
          </tr>
        </tbody>
      </n-table>
      <div v-else style="text-align: center; padding: 40px; color: var(--text-color-secondary)">
        暂无数据，请先查询
      </div>
    </n-card>
  </div>
</template>