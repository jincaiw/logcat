<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NTag, NSwitch, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createAlertRule, updateAlertRule, deleteAlertRule, getAlertRules, toggleAlertRule } from '@/api/alertRules'
import type { AlertRule } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import FormDialog, { type FieldConfig } from '@/components/common/FormDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'

const message = useMessage(); const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const formDialogRef = ref<InstanceType<typeof FormDialog> | null>(null)
const confirmDialogShow = ref(false); const confirmTitle = ref(''); const confirmContent = ref('')
const confirmAction = ref<() => Promise<void>>(() => Promise.resolve()); const confirmLoading = ref(false)
const editingRow = ref<AlertRule | null>(null)

const formFields: FieldConfig[] = [
  { key: 'name', label: '规则名称', type: 'text', required: true },
  { key: 'description', label: '描述', type: 'textarea' },
  { key: 'severity', label: '严重程度', type: 'select', required: true, options: [
    { label: '严重', value: 'critical' }, { label: '高', value: 'high' },
    { label: '中', value: 'medium' }, { label: '低', value: 'low' }, { label: '信息', value: 'info' },
  ], defaultValue: 'medium' },
  { key: 'status', label: '状态', type: 'select', options: [{ label: '启用', value: 1 }, { label: '禁用', value: 0 }], defaultValue: 1 },
  { key: 'cooldownSeconds', label: '冷却时间 (秒)', type: 'number', defaultValue: 300, min: 0 },
  { key: 'condition', label: '告警条件 (JSON)', type: 'code' },
  { key: 'pushConfigIds', label: '推送配置ID (JSON数组)', type: 'text' },
]

const severityMap: Record<string, any> = { critical: 'error', high: 'warning', medium: 'info', low: 'success', info: 'default' }
const severityLabel: Record<string, string> = { critical: '严重', high: '高', medium: '中', low: '低', info: '信息' }

const columns: DataTableColumns<AlertRule> = [
  { title: '名称', key: 'name' },
  { title: '严重程度', key: 'severity', render(row) { return h(NTag, { type: severityMap[row.severity] || 'default', size: 'small', bordered: false }, { default: () => severityLabel[row.severity] || row.severity }) } },
  { title: '冷却(s)', key: 'cooldownSeconds' },
  { title: '状态', key: 'status', render(row) { return h(NSwitch, { size: 'small', value: row.status === 1, onUpdateValue: (v: boolean) => handleToggle(row, v) }) } },
  {
    title: '操作', key: 'actions',
    render(row) { return h(NSpace, { size: 'small' }, { default: () => [
      h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', type: 'error', ghost: true, onClick: () => handleDelete(row) }, { default: () => '删除' }),
    ]})},
  },
]

async function fetchData(params: any) { return getAlertRules(params) }
function handleAdd() { editingRow.value = null; formDialogRef.value?.open({ status: 1, severity: 'medium', cooldownSeconds: 300 }) }
function handleEdit(row: AlertRule) {
  editingRow.value = row
  formDialogRef.value?.open({
    name: row.name, description: row.description, severity: row.severity,
    status: row.status, cooldownSeconds: row.cooldownSeconds,
    condition: typeof row.condition === 'string' ? row.condition : JSON.stringify(row.condition, null, 2),
    pushConfigIds: JSON.stringify(row.pushConfigIds || []),
  })
}

async function handleFormSubmit(data: Record<string, any>) {
  try {
    const payload = { ...data }
    if (payload.condition && typeof payload.condition === 'string') {
      try { payload.condition = JSON.parse(payload.condition) } catch { message.warning('JSON 格式不正确'); return }
    }
    if (typeof payload.pushConfigIds === 'string') {
      try { payload.pushConfigIds = JSON.parse(payload.pushConfigIds) } catch { message.warning('推送配置ID 格式不正确'); return }
    }
    if (editingRow.value) { await updateAlertRule(editingRow.value.id, payload); message.success('更新成功') }
    else { await createAlertRule(payload); message.success('创建成功') }
    formDialogRef.value?.close(); tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '操作失败') }
}

async function handleToggle(row: AlertRule, value: boolean) {
  try { await toggleAlertRule(row.id, value ? 1 : 0); message.success(value ? '已启用' : '已禁用'); tableRef.value?.loadData() }
  catch (err: any) { message.error(err?.message || '操作失败') }
}

function handleDelete(row: AlertRule) {
  confirmTitle.value = '删除'; confirmContent.value = `确定要删除 "${row.name}" 吗？`
  confirmAction.value = async () => { await deleteAlertRule(row.id); message.success('删除成功'); tableRef.value?.loadData() }
  confirmDialogShow.value = true
}

async function handleConfirm() {
  confirmLoading.value = true
  try { await confirmAction.value() } catch (err: any) { message.error(err?.message || '操作失败') } finally { confirmLoading.value = false; confirmDialogShow.value = false }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="告警规则" description="管理日志告警规则"><n-button type="primary" @click="handleAdd">添加规则</n-button></PageHeader>
    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['name']" search-placeholder="搜索规则名称" />
    <FormDialog ref="formDialogRef" :title="editingRow ? '编辑告警规则' : '添加告警规则'" :fields="formFields" :width="700" @submit="handleFormSubmit" />
    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>