<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NTag, NSwitch } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createDesensitizeRule, updateDesensitizeRule, deleteDesensitizeRule, getDesensitizeRules, toggleDesensitizeRule } from '@/api/desensitizeRules'
import type { DesensitizeRule } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import FormDialog, { type FieldConfig } from '@/components/common/FormDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAppMessage } from '@/composables/useMessage'

const message = useAppMessage(); const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const formDialogRef = ref<InstanceType<typeof FormDialog> | null>(null)
const confirmDialogShow = ref(false); const confirmTitle = ref(''); const confirmContent = ref('')
const confirmAction = ref<() => Promise<void>>(() => Promise.resolve()); const confirmLoading = ref(false)
const editingRow = ref<DesensitizeRule | null>(null)

const formFields: FieldConfig[] = [
  { key: 'fieldName', label: '字段名称', type: 'text', required: true },
  { key: 'ruleType', label: '规则类型', type: 'select', required: true, options: [
    { label: '正则替换', value: 'regex' }, { label: '掩码', value: 'mask' }, { label: '哈希', value: 'hash' }, { label: '删除', value: 'remove' },
  ], defaultValue: 'regex' },
  { key: 'ruleConfig', label: '规则配置 (JSON)', type: 'code' },
  { key: 'enabled', label: '启用状态', type: 'select', options: [{ label: '启用', value: true }, { label: '禁用', value: false }], defaultValue: true },
]

const columns: DataTableColumns<DesensitizeRule> = [
  { title: '字段', key: 'fieldName' },
  { title: '规则类型', key: 'ruleType' },
  { title: '启用', key: 'enabled', render(row) { return h(NSwitch, { size: 'small', value: row.enabled, onUpdateValue: (v: boolean) => handleToggle(row, v) }) } },
  {
    title: '操作', key: 'actions',
    render(row) { return h(NSpace, { size: 'small' }, { default: () => [
      h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', type: 'error', ghost: true, onClick: () => handleDelete(row) }, { default: () => '删除' }),
    ]})},
  },
]

async function fetchData(params: any) { return getDesensitizeRules(params) }
function handleAdd() { editingRow.value = null; formDialogRef.value?.open({ enabled: true, ruleType: 'regex' }) }
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
  confirmTitle.value = '删除'; confirmContent.value = `确定要删除字段 "${row.fieldName}" 的脱敏规则吗？`
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
    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['fieldName']" search-placeholder="搜索字段名称" />
    <FormDialog ref="formDialogRef" :title="editingRow ? '编辑脱敏规则' : '添加脱敏规则'" :fields="formFields" @submit="handleFormSubmit" />
    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>