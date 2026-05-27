<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NTag, NSwitch, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createDesensitizeRule, updateDesensitizeRule, deleteDesensitizeRule, getDesensitizeRules, toggleDesensitizeRule } from '@/api/desensitizeRules'
import type { DesensitizeRule } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import FormDialog, { type FieldConfig } from '@/components/common/FormDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'

const message = useMessage(); const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const formDialogRef = ref<InstanceType<typeof FormDialog> | null>(null)
const confirmDialogShow = ref(false); const confirmTitle = ref(''); const confirmContent = ref('')
const confirmAction = ref<() => Promise<void>>(() => Promise.resolve()); const confirmLoading = ref(false)
const editingRow = ref<DesensitizeRule | null>(null)

const formFields: FieldConfig[] = [
  { key: 'name', label: '规则名称', type: 'text', required: true },
  { key: 'field', label: '目标字段', type: 'text', required: true },
  { key: 'pattern', label: '匹配模式 (Regex)', type: 'text', required: true },
  { key: 'replacement', label: '替换文本', type: 'text', defaultValue: '***' },
  { key: 'description', label: '描述', type: 'textarea' },
  { key: 'status', label: '状态', type: 'select', options: [{ label: '启用', value: 1 }, { label: '禁用', value: 0 }], defaultValue: 1 },
]

const columns: DataTableColumns<DesensitizeRule> = [
  { title: '名称', key: 'name' },
  { title: '字段', key: 'field' },
  { title: '模式', key: 'pattern' },
  { title: '替换', key: 'replacement' },
  { title: '状态', key: 'status', render(row) { return h(NSwitch, { size: 'small', value: row.status === 1, onUpdateValue: (v: boolean) => handleToggle(row, v) }) } },
  {
    title: '操作', key: 'actions',
    render(row) { return h(NSpace, { size: 'small' }, { default: () => [
      h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', type: 'error', ghost: true, onClick: () => handleDelete(row) }, { default: () => '删除' }),
    ]})},
  },
]

async function fetchData(params: any) { return getDesensitizeRules(params) }
function handleAdd() { editingRow.value = null; formDialogRef.value?.open({ status: 1, replacement: '***' }) }
function handleEdit(row: DesensitizeRule) { editingRow.value = row; formDialogRef.value?.open(row) }

async function handleFormSubmit(data: Record<string, any>) {
  try {
    if (editingRow.value) { await updateDesensitizeRule(editingRow.value.id, data); message.success('更新成功') }
    else { await createDesensitizeRule(data); message.success('创建成功') }
    formDialogRef.value?.close(); tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '操作失败') }
}

async function handleToggle(row: DesensitizeRule, value: boolean) {
  try { await toggleDesensitizeRule(row.id, value ? 1 : 0); message.success(value ? '已启用' : '已禁用'); tableRef.value?.loadData() }
  catch (err: any) { message.error(err?.message || '操作失败') }
}

function handleDelete(row: DesensitizeRule) {
  confirmTitle.value = '删除'; confirmContent.value = `确定要删除 "${row.name}" 吗？`
  confirmAction.value = async () => { await deleteDesensitizeRule(row.id); message.success('删除成功'); tableRef.value?.loadData() }
  confirmDialogShow.value = true
}

async function handleConfirm() {
  confirmLoading.value = true
  try { await confirmAction.value() } catch (err: any) { message.error(err?.message || '操作失败') } finally { confirmLoading.value = false; confirmDialogShow.value = false }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="脱敏规则" description="管理日志敏感信息脱敏规则"><n-button type="primary" @click="handleAdd">添加规则</n-button></PageHeader>
    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['name', 'field']" search-placeholder="搜索规则名称或字段" />
    <FormDialog ref="formDialogRef" :title="editingRow ? '编辑脱敏规则' : '添加脱敏规则'" :fields="formFields" @submit="handleFormSubmit" />
    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>