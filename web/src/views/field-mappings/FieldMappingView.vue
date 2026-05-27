<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createFieldMapping, updateFieldMapping, deleteFieldMapping, getFieldMappings } from '@/api/fieldMappings'
import type { FieldMappingDoc } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import FormDialog, { type FieldConfig } from '@/components/common/FormDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'

const message = useMessage()
const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const formDialogRef = ref<InstanceType<typeof FormDialog> | null>(null)
const confirmDialogShow = ref(false)
const confirmTitle = ref('')
const confirmContent = ref('')
const confirmAction = ref<() => Promise<void>>(() => Promise.resolve())
const confirmLoading = ref(false)
const editingRow = ref<FieldMappingDoc | null>(null)

const formFields: FieldConfig[] = [
  { key: 'name', label: '映射名称', type: 'text', required: true },
  { key: 'description', label: '描述', type: 'textarea' },
  { key: 'mappings', label: '映射规则 (JSON)', type: 'code' },
]

const columns: DataTableColumns<FieldMappingDoc> = [
  { title: '名称', key: 'name' },
  { title: '描述', key: 'description' },
  { title: '创建时间', key: 'createdAt' },
  {
    title: '操作', key: 'actions',
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
          h(NButton, { size: 'small', type: 'error', ghost: true, onClick: () => handleDelete(row) }, { default: () => '删除' }),
        ],
      })
    },
  },
]

async function fetchData(params: any) { return getFieldMappings(params) }
function handleAdd() { editingRow.value = null; formDialogRef.value?.open() }
function handleEdit(row: FieldMappingDoc) {
  editingRow.value = row
  formDialogRef.value?.open({
    name: row.name,
    description: row.description,
    mappings: JSON.stringify(row.mappings || [], null, 2),
  })
}

async function handleFormSubmit(data: Record<string, any>) {
  try {
    const payload = { ...data }
    if (payload.mappings && typeof payload.mappings === 'string') {
      try { payload.mappings = JSON.parse(payload.mappings) } catch { message.warning('JSON 格式不正确'); return }
    }
    if (editingRow.value) { await updateFieldMapping(editingRow.value.id, payload); message.success('更新成功') }
    else { await createFieldMapping(payload); message.success('创建成功') }
    formDialogRef.value?.close(); tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '操作失败') }
}

function handleDelete(row: FieldMappingDoc) {
  confirmTitle.value = '删除'; confirmContent.value = `确定要删除 "${row.name}" 吗？`
  confirmAction.value = async () => { await deleteFieldMapping(row.id); message.success('删除成功'); tableRef.value?.loadData() }
  confirmDialogShow.value = true
}

async function handleConfirm() {
  confirmLoading.value = true
  try { await confirmAction.value() } catch (err: any) { message.error(err?.message || '操作失败') } finally { confirmLoading.value = false; confirmDialogShow.value = false }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="字段映射" description="管理日志字段映射规则"><n-button type="primary" @click="handleAdd">添加映射</n-button></PageHeader>
    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['name']" search-placeholder="搜索映射名称" />
    <FormDialog ref="formDialogRef" :title="editingRow ? '编辑映射' : '添加映射'" :fields="formFields" @submit="handleFormSubmit" />
    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>