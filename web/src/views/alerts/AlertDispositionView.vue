<script setup lang="ts">
import { h } from 'vue'
import { NTag } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getAlertDispositions } from '@/api/alerts'
import type { AlertDisposition } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import PageHeader from '@/components/common/PageHeader.vue'

const columns: DataTableColumns<AlertDisposition> = [
  { title: '时间', key: 'createdAt', width: 160 },
  {
    title: '告警ID',
    key: 'alertRecordId',
    width: 100,
    render(row) { return row.alertRecordId || row.aggregatedAlertId || '-' },
  },
  {
    title: '状态', key: 'status', width: 100,
    render(row) {
      const map: Record<string, 'default' | 'success' | 'warning' | 'info'> = {
        confirmed: 'success',
        ignored: 'default',
        closed: 'warning',
        acknowledged: 'info',
        resolved: 'success',
      }
      const labels: Record<string, string> = {
        confirmed: '已确认',
        ignored: '已忽略',
        closed: '已关闭',
        acknowledged: '已确认',
        resolved: '已解决',
      }
      return h(NTag, { type: map[row.status] || 'default', size: 'small', bordered: false }, { default: () => labels[row.status] || row.status })
    },
  },
  { title: '操作人', key: 'operatorName', width: 100 },
  { title: '备注', key: 'note', ellipsis: { tooltip: true } },
]

async function fetchData(params: any) { return getAlertDispositions(params) }
</script>

<template>
  <div class="page-container">
    <PageHeader title="告警处置记录" description="查看告警处置历史" />
    <DataTable :columns="columns" :fetch-api="fetchData" :search-fields="['operatorName']" search-placeholder="搜索操作人" />
  </div>
</template>
