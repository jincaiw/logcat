<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NTabs, NTabPane, NInput, NCard, NTag, NInputNumber, NSelect, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createFilterPolicy, updateFilterPolicy, deleteFilterPolicy, getFilterPolicies, testFilterPolicy } from '@/api/filterPolicies'
import type { FilterPolicy, FilterCondition } from '@/types'
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
const editingRow = ref<FilterPolicy | null>(null)
const testSample = ref('')
const testResult = ref<any>(null)
const testLoading = ref(false)

const formFields: FieldConfig[] = [
  { key: 'name', label: '策略名称', type: 'text', required: true },
  { key: 'description', label: '描述', type: 'textarea' },
  { key: 'priority', label: '优先级', type: 'number', required: true, defaultValue: 100 },
  { key: 'status', label: '状态', type: 'select', options: [{ label: '启用', value: 1 }, { label: '禁用', value: 0 }], defaultValue: 1 },
  { key: 'action', label: '动作', type: 'select', required: true, options: [{ label: '接受', value: 'accept' }, { label: '丢弃', value: 'drop' }], defaultValue: 'accept' },
  { key: 'conditions', label: '条件 (JSON)', type: 'code' },
]

const columns: DataTableColumns<FilterPolicy> = [
  { title: '名称', key: 'name' },
  { title: '优先级', key: 'priority' },
  { title: '动作', key: 'action', render(row) { return h(NTag, { type: row.action === 'accept' ? 'success' : 'error', size: 'small', bordered: false }, { default: () => row.action === 'accept' ? '接受' : '丢弃' }) } },
  { title: '状态', key: 'status', render(row) { return h(NTag, { type: row.status === 1 ? 'success' : 'default', size: 'small', bordered: false }, { default: () => row.status === 1 ? '启用' : '禁用' }) } },
  {
    title: '操作', key: 'actions',
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
          h(NButton, { size: 'small', onClick: () => handleTest(row) }, { default: () => '测试' }),
          h(NButton, { size: 'small', type: 'error', ghost: true, onClick: () => handleDelete(row) }, { default: () => '删除' }),
        ],
      })
    },
  },
]

async function fetchData(params: any) { return getFilterPolicies(params) }
function handleAdd() { editingRow.value = null; formDialogRef.value?.open() }
function handleEdit(row: FilterPolicy) {
  editingRow.value = row
  formDialogRef.value?.open({
    name: row.name, description: row.description, priority: row.priority,
    status: row.status, action: row.action,
    conditions: JSON.stringify(row.conditions || [], null, 2),
  })
}

async function handleFormSubmit(data: Record<string, any>) {
  try {
    const payload = { ...data }
    if (payload.conditions && typeof payload.conditions === 'string') {
      try { payload.conditions = JSON.parse(payload.conditions) } catch { message.warning('JSON 格式不正确'); return }
    }
    if (editingRow.value) { await updateFilterPolicy(editingRow.value.id, payload); message.success('更新成功') }
    else { await createFilterPolicy(payload); message.success('创建成功') }
    formDialogRef.value?.close(); tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '操作失败') }
}

function handleDelete(row: FilterPolicy) {
  confirmTitle.value = '删除'; confirmContent.value = `确定要删除 "${row.name}" 吗？`
  confirmAction.value = async () => { await deleteFilterPolicy(row.id); message.success('删除成功'); tableRef.value?.loadData() }
  confirmDialogShow.value = true
}

function handleTest(row: FilterPolicy) { editingRow.value = row; testSample.value = ''; testResult.value = null }
async function handleRunTest() {
  if (!editingRow.value || !testSample.value) { message.warning('请输入测试数据'); return }
  testLoading.value = true
  try {
    let logData: Record<string, any>
    try { logData = JSON.parse(testSample.value) } catch { logData = { message: testSample.value } }
    const res = await testFilterPolicy(editingRow.value.id, logData)
    testResult.value = res.data
  } catch (err: any) { message.error(err?.message || '测试失败') } finally { testLoading.value = false }
}

async function handleConfirm() {
  confirmLoading.value = true
  try { await confirmAction.value() } catch (err: any) { message.error(err?.message || '操作失败') } finally { confirmLoading.value = false; confirmDialogShow.value = false }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="过滤策略" description="管理日志过滤规则"><n-button type="primary" @click="handleAdd">添加策略</n-button></PageHeader>

    <n-tabs type="line">
      <n-tab-pane name="list" tab="策略列表">
        <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['name']" search-placeholder="搜索策略名称" />
      </n-tab-pane>
      <n-tab-pane name="test" tab="策略测试" :disabled="!editingRow">
        <n-card v-if="editingRow" size="small">
          <h3>测试策略: {{ editingRow.name }}</h3>
          <n-input v-model:value="testSample" type="textarea" placeholder="输入测试数据 (JSON 格式)..." :autosize="{ minRows: 4, maxRows: 10 }" style="margin-bottom: 12px" />
          <n-button type="primary" :loading="testLoading" @click="handleRunTest">执行测试</n-button>
          <div v-if="testResult" style="margin-top: 16px">
            <n-tag :type="testResult.matched ? 'success' : 'warning'" size="small">
              匹配操作: {{ testResult.action === 'accept' ? '接受' : '丢弃' }} | 匹配结果: {{ testResult.matched ? '匹配' : '未匹配' }}
            </n-tag>
            <pre style="margin-top: 8px; padding: 12px; background: var(--bg-color); border-radius: 4px; overflow-x: auto">{{ JSON.stringify(testResult, null, 2) }}</pre>
          </div>
        </n-card>
      </n-tab-pane>
    </n-tabs>

    <FormDialog ref="formDialogRef" :title="editingRow ? '编辑策略' : '添加策略'" :fields="formFields" :width="700" @submit="handleFormSubmit" />
    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>