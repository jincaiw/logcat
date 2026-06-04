<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NTag } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createOutputTemplate, updateOutputTemplate, deleteOutputTemplate, getOutputTemplates } from '@/api/outputTemplates'
import type { OutputTemplate } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import FormDialog, { type FieldConfig } from '@/components/common/FormDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAppMessage } from '@/composables/useMessage'

const message = useAppMessage(); const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const formDialogRef = ref<InstanceType<typeof FormDialog> | null>(null)
const confirmDialogShow = ref(false); const confirmTitle = ref(''); const confirmContent = ref('')
const confirmAction = ref<() => Promise<void>>(() => Promise.resolve()); const confirmLoading = ref(false)
const editingRow = ref<OutputTemplate | null>(null)

const formFields: FieldConfig[] = [
  { key: 'name', label: '模板名称', type: 'text', required: true },
  { key: 'channelType', label: '通道类型', type: 'select', required: true, options: [
    { label: 'HTTP', value: 'http' }, { label: 'Email', value: 'email' },
    { label: 'Syslog', value: 'syslog' },
  ], defaultValue: 'http' },
  { key: 'content', label: '模板内容', type: 'code', required: true },
  { key: 'fields', label: '字段 (JSON)', type: 'code' },
  { key: 'deviceType', label: '设备类型', type: 'text' },
  { key: 'enabled', label: '启用状态', type: 'select', options: [{ label: '启用', value: true }, { label: '禁用', value: false }], defaultValue: true },
]

const columns: DataTableColumns<OutputTemplate> = [
  { title: '名称', key: 'name' },
  { title: '通道类型', key: 'channelType', render(row) { return h(NTag, { size: 'small', bordered: false }, { default: () => row.channelType?.toUpperCase() || '--' }) } },
  {
    title: '操作', key: 'actions',
    render(row) { return h(NSpace, { size: 'small' }, { default: () => [
      h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', type: 'error', ghost: true, onClick: () => handleDelete(row) }, { default: () => '删除' }),
    ]})},
  },
]

async function fetchData(params: any) { return getOutputTemplates(params) }
function handleAdd() { editingRow.value = null; formDialogRef.value?.open() }
function handleEdit(row: OutputTemplate) { editingRow.value = row; formDialogRef.value?.open(row) }

async function handleFormSubmit(data: Record<string, any>) {
  try {
    if (editingRow.value) { await updateOutputTemplate(editingRow.value.id, data); message.success('更新成功') }
    else { await createOutputTemplate(data); message.success('创建成功') }
    formDialogRef.value?.close(); tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '操作失败') }
}

function handleDelete(row: OutputTemplate) {
  confirmTitle.value = '删除'; confirmContent.value = `确定要删除 "${row.name}" 吗？`
  confirmAction.value = async () => { await deleteOutputTemplate(row.id); message.success('删除成功'); tableRef.value?.loadData() }
  confirmDialogShow.value = true
}

async function handleConfirm() {
  confirmLoading.value = true
  try { await confirmAction.value() } catch (err: any) { message.error(err?.message || '操作失败') } finally { confirmLoading.value = false; confirmDialogShow.value = false }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="输出模板" description="管理日志输出格式"><n-button type="primary" @click="handleAdd">添加模板</n-button></PageHeader>
    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['name']" search-placeholder="搜索模板名称" />
    <FormDialog ref="formDialogRef" :title="editingRow ? '编辑模板' : '添加模板'" :fields="formFields" :width="700" @submit="handleFormSubmit" />
    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>