<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NTag, NSpace, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createDevice, updateDevice, deleteDevice, getDevices, getDeviceGroups, getAllDeviceGroups } from '@/api/devices'
import { getDeviceTemplates, getAllDeviceTemplates } from '@/api/deviceTemplates'
import type { Device, DeviceGroup, DeviceTemplate } from '@/types'
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

const editingRow = ref<Device | null>(null)

const formFields: FieldConfig[] = [
  { key: 'name', label: '设备名称', type: 'text', required: true },
  { key: 'host', label: '主机地址', type: 'text', required: true },
  { key: 'port', label: '端口', type: 'number', required: true, min: 1, max: 65535, defaultValue: 514 },
  { key: 'protocol', label: '协议', type: 'select', required: true, options: [{ label: 'TCP', value: 'tcp' }, { label: 'UDP', value: 'udp' }, { label: 'TLS', value: 'tls' }], defaultValue: 'udp' },
  { key: 'deviceGroupId', label: '设备分组', type: 'select', options: [], placeholder: '请选择分组' },
  { key: 'deviceTemplateId', label: '设备模板', type: 'select', options: [], placeholder: '请选择模板' },
  { key: 'description', label: '描述', type: 'textarea' },
]

const columns: DataTableColumns<Device> = [
  { title: '名称', key: 'name' },
  { title: '主机', key: 'host' },
  { title: '端口', key: 'port' },
  { title: '协议', key: 'protocol' },
  { title: '分组', key: 'deviceGroupName' },
  {
    title: '状态', key: 'status',
    render(row) {
      return h(NTag, { type: row.status === 1 ? 'success' : 'default', size: 'small', bordered: false },
        { default: () => row.status === 1 ? '在线' : '离线' })
    },
  },
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

async function loadOptions() {
  try {
    const [groupsRes, templatesRes] = await Promise.all([
      getAllDeviceGroups(),
      getAllDeviceTemplates(),
    ])
    const groups = groupsRes.data || []
    const templates = templatesRes.data || []
    const groupField = formFields.find((f) => f.key === 'deviceGroupId')
    const templateField = formFields.find((f) => f.key === 'deviceTemplateId')
    if (groupField) groupField.options = groups.map((g) => ({ label: g.name, value: g.id }))
    if (templateField) templateField.options = templates.map((t) => ({ label: t.name, value: t.id }))
  } catch { /* ignore */ }
}

async function fetchData(params: any) {
  const res = await getDevices(params)
  return res
}

function handleAdd() {
  editingRow.value = null
  loadOptions()
  formDialogRef.value?.open()
}

function handleEdit(row: Device) {
  editingRow.value = row
  loadOptions()
  formDialogRef.value?.open(row)
}

async function handleFormSubmit(data: Record<string, any>) {
  try {
    if (editingRow.value) {
      await updateDevice(editingRow.value.id, data)
      message.success('设备更新成功')
    } else {
      await createDevice(data)
      message.success('设备创建成功')
    }
    formDialogRef.value?.close()
    tableRef.value?.loadData()
  } catch (err: any) {
    message.error(err?.message || '操作失败')
  }
}

function handleDelete(row: Device) {
  confirmTitle.value = '删除设备'
  confirmContent.value = `确定要删除设备 "${row.name}" 吗？`
  confirmAction.value = async () => {
    await deleteDevice(row.id)
    message.success('删除成功')
    tableRef.value?.loadData()
  }
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
    <PageHeader title="设备管理" description="管理日志采集设备">
      <n-button type="primary" @click="handleAdd">添加设备</n-button>
    </PageHeader>
    <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['name', 'host']" search-placeholder="搜索设备名称或主机" />
    <FormDialog ref="formDialogRef" :title="editingRow ? '编辑设备' : '添加设备'" :fields="formFields" @submit="handleFormSubmit" />
    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>