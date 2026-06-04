<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NTag, NSwitch } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createAlertRule, updateAlertRule, deleteAlertRule, getAlertRules, toggleAlertRule } from '@/api/alertRules'
import { getFilterPolicies } from '@/api/filterPolicies'
import { getPushConfigs } from '@/api/pushConfigs'
import { getOutputTemplates } from '@/api/outputTemplates'
import type { AlertRule } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import FormDialog, { type FieldConfig } from '@/components/common/FormDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAppMessage } from '@/composables/useMessage'

const message = useAppMessage(); const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const formDialogRef = ref<InstanceType<typeof FormDialog> | null>(null)
const confirmDialogShow = ref(false); const confirmTitle = ref(''); const confirmContent = ref('')
const confirmAction = ref<() => Promise<void>>(() => Promise.resolve()); const confirmLoading = ref(false)
const editingRow = ref<AlertRule | null>(null)

const channelTypeOptions = [
  { label: 'HTTP', value: 'http' }, { label: 'Email', value: 'email' }, { label: 'Syslog', value: 'syslog' },
]

const formFields: FieldConfig[] = [
  { key: 'name', label: '规则名称', type: 'text', required: true },
  { key: 'filterPolicyId', label: '过滤策略', type: 'select', required: true, options: [], placeholder: '请选择过滤策略' },
  { key: 'pushConfigId', label: '推送配置', type: 'select', required: true, options: [], placeholder: '请选择推送配置' },
  { key: 'outputTemplateId', label: '输出模板', type: 'select', options: [], placeholder: '请选择输出模板' },
  { key: 'channelType', label: '通道类型', type: 'select', required: true, options: channelTypeOptions, defaultValue: 'http' },
  { key: 'enabled', label: '启用状态', type: 'select', options: [{ label: '启用', value: true }, { label: '禁用', value: false }], defaultValue: true },
]

const columns: DataTableColumns<AlertRule> = [
  { title: '名称', key: 'name' },
  { title: '通道类型', key: 'channelType', render(row) { return h(NTag, { size: 'small', bordered: false }, { default: () => row.channelType?.toUpperCase() || '--' }) } },
  { title: '启用', key: 'enabled', render(row) { return h(NSwitch, { size: 'small', value: row.enabled, onUpdateValue: (v: boolean) => handleToggle(row, v) }) } },
  {
    title: '操作', key: 'actions',
    render(row) { return h(NSpace, { size: 'small' }, { default: () => [
      h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', type: 'error', ghost: true, onClick: () => handleDelete(row) }, { default: () => '删除' }),
    ]})},
  },
]

async function loadOptions() {
  try {
    const [fpRes, pcRes, otRes] = await Promise.all([
      getFilterPolicies({ page: 1, pageSize: 1000 }),
      getPushConfigs({ page: 1, pageSize: 1000 }),
      getOutputTemplates({ page: 1, pageSize: 1000 }),
    ])
    const fpField = formFields.find((f) => f.key === 'filterPolicyId')
    const pcField = formFields.find((f) => f.key === 'pushConfigId')
    const otField = formFields.find((f) => f.key === 'outputTemplateId')
    if (fpField) fpField.options = (fpRes.data?.list || fpRes.data?.items || []).map((i: any) => ({ label: i.name, value: i.id }))
    if (pcField) pcField.options = (pcRes.data?.list || pcRes.data?.items || []).map((i: any) => ({ label: i.name, value: i.id }))
    if (otField) otField.options = (otRes.data?.list || otRes.data?.items || []).map((i: any) => ({ label: i.name, value: i.id }))
  } catch { /* ignore */ }
}

async function fetchData(params: any) { return getAlertRules(params) }
async function handleAdd() { editingRow.value = null; await loadOptions(); formDialogRef.value?.open({ enabled: true, channelType: 'http' }) }
async function handleEdit(row: AlertRule) {
  editingRow.value = row
  await loadOptions()
  formDialogRef.value?.open({
    name: row.name, filterPolicyId: row.filterPolicyId, pushConfigId: row.pushConfigId,
    outputTemplateId: row.outputTemplateId, channelType: row.channelType, enabled: row.enabled,
  })
}

async function handleFormSubmit(data: Record<string, any>) {
  try {
    const payload = { ...data }
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
