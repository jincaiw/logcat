<script setup lang="ts">
import { ref, h } from 'vue'
import { NTag } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getAlertDispositions } from '@/api/alerts'
import type { AlertDisposition } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import PageHeader from '@/components/common/PageHeader.vue'

const columns: DataTableColumns<AlertDisposition> = [
  { title: '时间', key: 'createdAt', width: 160 },
  { title: '告警ID', key: 'alertId', width: 80 },
  {
    title: '操作', key: 'action', width: 80,
    render(row) { const m: Record<string, any> = { confirm: 'success', ignore: 'default', close: 'warning' }; const l: Record<string, string> = { confirm: '确认', ignore: '忽略', close: '关闭' }; return h(NTag, { type: m[row.action], size: 'small', bordered: false }, { default: () => l[row.action] || row.action }) },
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