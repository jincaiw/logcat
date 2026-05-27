<script setup lang="ts">
import { ref, h } from 'vue'
import { NTag } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getAuditLogs } from '@/api/auditLogs'
import type { AuditLog } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import PageHeader from '@/components/common/PageHeader.vue'

const columns: DataTableColumns<AuditLog> = [
  { title: '时间', key: 'createdAt', width: 160 },
  { title: '用户', key: 'username', width: 100 },
  { title: '操作', key: 'action', width: 120 },
  { title: '资源', key: 'resource', width: 120 },
  { title: '资源ID', key: 'resourceId', width: 100 },
  {
    title: '结果', key: 'result', width: 80,
    render(row) { return h(NTag, { type: row.result === 'success' ? 'success' : 'error', size: 'small', bordered: false }, { default: () => row.result === 'success' ? '成功' : '失败' }) },
  },
  { title: '详情', key: 'detail', ellipsis: { tooltip: true } },
  { title: 'IP', key: 'ip', width: 130 },
]

async function fetchData(params: any) { return getAuditLogs(params) }
</script>

<template>
  <div class="page-container">
    <PageHeader title="审计日志" description="系统操作审计记录" />
    <DataTable :columns="columns" :fetch-api="fetchData" :search-fields="['username', 'action', 'resource']" search-placeholder="搜索用户、操作或资源" />
  </div>
</template>