<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createDeviceTemplate, updateDeviceTemplate, deleteDeviceTemplate, getDeviceTemplates } from '@/api/deviceTemplates'
import type { DeviceTemplate } from '@/types'
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
const editingRow = ref<DeviceTemplate | null>(null)

const formFields: FieldConfig[] = [
  { key: 'name', label: '模板名称', type: 'text', required: true },
  { key: 'description', label: '描述', type: 'textarea' },
  { key: 'config', label: '配置 (JSON)', type: 'code' },
]

const columns: DataTableColumns<DeviceTemplate> = [
  { title: '模板名称', key: 'name' },
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

async function fetchData(params: any) { return getDeviceTemplates(params) }
function handleAdd() { editingRow.value = null; formDialogRef.value?.open() }
function handleEdit(row: DeviceTemplate) {
  editingRow.value = row
  formDialogRef.value?.open({
    name: row.name,
    description: row.description,
    config: typeof row.config === 'string' ? row.config : JSON.stringify(row.config, null, 2),
  })
}

async function handleFormSubmit(data: Record<string, any>) {
  try {
    const payload = { ...data }
    if (payload.config && typeof payload.config === 'string') {
      try { payload.config = JSON.parse(payload.config) } catch { message.warning('配置 JSON 格式不正确'); return }
    }
    if (editingRow.value) { await updateDeviceTemplate(editingRow.value.id, payload); message.success('更新成功') }
    else { await createDeviceTemplate(payload); message.success('创建成功') }
    formDialogRef.value?.close(); tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '操作失败') }
}

function handleDelete(row: DeviceTemplate) {
  confirmTitle.value = '删除模板'; confirmContent.value = `确定要删除模板 "${row.name}" 吗？`
  confirmAction.value = async () => { await deleteDeviceTemplate(row.id); message.success('删除成功'); tableRef.value?.loadData() }
  confirmDialogShow.value = true
}

async function handleConfirm() {
  confirmLoading.value = true
  try { await confirmAction.value() } catch (err: any) { message.error(err?.message || '操作失败') }
  finally { confirmLoading.value = false; confirmDialogShow.value = false }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="设备模板" description="管理设备配置模板"><n-button type="primary" @click="handleAdd">添加模板</n-button></PageHeader>
    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['name']" search-placeholder="搜索模板名称" />
    <FormDialog ref="formDialogRef" :title="editingRow ? '编辑模板' : '添加模板'" :fields="formFields" @submit="handleFormSubmit" />
    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>