<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NTabs, NTabPane, NTag, NSwitch, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createPushConfig, updatePushConfig, deletePushConfig, getPushConfigs, testPushConfig } from '@/api/pushConfigs'
import type { PushConfig } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import FormDialog, { type FieldConfig } from '@/components/common/FormDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'

const message = useMessage()
const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const formDialogRef = ref<InstanceType<typeof FormDialog> | null>(null)
const confirmDialogShow = ref(false)
const confirmTitle = ref(''); const confirmContent = ref('')
const confirmAction = ref<() => Promise<void>>(() => Promise.resolve())
const confirmLoading = ref(false)
const editingRow = ref<PushConfig | null>(null)
const pushType = ref('http')

const getFormFields = (): FieldConfig[] => {
  const base: FieldConfig[] = [
    { key: 'name', label: '配置名称', type: 'text', required: true },
    { key: 'type', label: '类型', type: 'select', required: true, options: [
      { label: 'HTTP', value: 'http' }, { label: 'Email', value: 'email' }, { label: 'Syslog', value: 'syslog' },
    ], defaultValue: 'http' },
    { key: 'status', label: '状态', type: 'select', options: [{ label: '启用', value: 1 }, { label: '禁用', value: 0 }], defaultValue: 1 },
    { key: 'config', label: '配置 (JSON)', type: 'code', required: true },
  ]
  return base
}

const columns: DataTableColumns<PushConfig> = [
  { title: '名称', key: 'name' },
  { title: '类型', key: 'type', render(row: any) { return h(NTag, { size: 'small', bordered: false }, { default: () => row.type?.toUpperCase() }) } },
  { title: '状态', key: 'status', render(row) { return h(NSwitch, { size: 'small', value: row.status === 1, onUpdateValue: () => {} }) } },
  {
    title: '操作', key: 'actions',
    render(row) {
      return h(NSpace, { size: 'small' }, { default: () => [
        h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'small', onClick: () => handleTest(row) }, { default: () => '测试' }),
        h(NButton, { size: 'small', type: 'error', ghost: true, onClick: () => handleDelete(row) }, { default: () => '删除' }),
      ]})
    },
  },
]

async function fetchData(params: any) { return getPushConfigs(params) }
function handleAdd() { editingRow.value = null; formDialogRef.value?.open({ type: 'http', status: 1 }) }
function handleEdit(row: PushConfig) {
  editingRow.value = row
  const data: any = { name: row.name, type: row.type, status: row.status }
  data.config = typeof row.config === 'string' ? row.config : JSON.stringify(row.config, null, 2)
  formDialogRef.value?.open(data)
}

async function handleFormSubmit(data: Record<string, any>) {
  try {
    const payload = { ...data }
    if (payload.config && typeof payload.config === 'string') {
      try { payload.config = JSON.parse(payload.config) } catch { message.warning('JSON 格式不正确'); return }
    }
    if (editingRow.value) { await updatePushConfig(editingRow.value.id, payload); message.success('更新成功') }
    else { await createPushConfig(payload); message.success('创建成功') }
    formDialogRef.value?.close(); tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '操作失败') }
}

async function handleTest(row: PushConfig) {
  try {
    const res = await testPushConfig(row.id)
    message.success(res.data?.message || '测试成功')
  } catch (err: any) { message.error(err?.message || '测试失败') }
}

function handleDelete(row: PushConfig) {
  confirmTitle.value = '删除'; confirmContent.value = `确定要删除 "${row.name}" 吗？`
  confirmAction.value = async () => { await deletePushConfig(row.id); message.success('删除成功'); tableRef.value?.loadData() }
  confirmDialogShow.value = true
}

async function handleConfirm() {
  confirmLoading.value = true
  try { await confirmAction.value() } catch (err: any) { message.error(err?.message || '操作失败') } finally { confirmLoading.value = false; confirmDialogShow.value = false }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="推送配置" description="管理日志推送通道 (HTTP/Email/Syslog)"><n-button type="primary" @click="handleAdd">添加配置</n-button></PageHeader>

    <n-tabs type="line">
      <n-tab-pane name="all" tab="全部"><DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['name']" :extra-params="{}" search-placeholder="搜索配置名称" /></n-tab-pane>
    </n-tabs>

    <FormDialog ref="formDialogRef" :title="editingRow ? '编辑推送配置' : '添加推送配置'" :fields="getFormFields()" :width="700" @submit="handleFormSubmit" />
    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>