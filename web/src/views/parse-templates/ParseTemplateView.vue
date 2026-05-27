<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NTabs, NTabPane, NInput, NCard, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createParseTemplate, updateParseTemplate, deleteParseTemplate, getParseTemplates, testParseTemplate } from '@/api/parseTemplates'
import type { ParseTemplate, ParseTestResult } from '@/types'
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
const editingRow = ref<ParseTemplate | null>(null)
const testTabActive = ref('list')
const testSample = ref('')
const testResult = ref<ParseTestResult | null>(null)
const testLoading = ref(false)

const formFields: FieldConfig[] = [
  { key: 'name', label: '模板名称', type: 'text', required: true },
  { key: 'type', label: '类型', type: 'select', required: true, options: [
    { label: 'Regex', value: 'regex' }, { label: 'JSON', value: 'json' },
    { label: 'KV', value: 'kv' }, { label: 'Grok', value: 'grok' }, { label: 'Custom', value: 'custom' },
  ], defaultValue: 'regex' },
  { key: 'pattern', label: '解析模式', type: 'code', required: true },
  { key: 'description', label: '描述', type: 'textarea' },
]

const columns: DataTableColumns<ParseTemplate> = [
  { title: '名称', key: 'name' },
  {
    title: '类型', key: 'type',
    render(row) { return h(NTag, { size: 'small', bordered: false }, { default: () => row.type.toUpperCase() }) },
  },
  { title: '描述', key: 'description' },
  { title: '创建时间', key: 'createdAt' },
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

async function fetchData(params: any) { return getParseTemplates(params) }
function handleAdd() { editingRow.value = null; formDialogRef.value?.open() }
function handleEdit(row: ParseTemplate) { editingRow.value = row; formDialogRef.value?.open(row) }

async function handleFormSubmit(data: Record<string, any>) {
  try {
    if (editingRow.value) { await updateParseTemplate(editingRow.value.id, data); message.success('更新成功') }
    else { await createParseTemplate(data); message.success('创建成功') }
    formDialogRef.value?.close(); tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '操作失败') }
}

function handleDelete(row: ParseTemplate) {
  confirmTitle.value = '删除'; confirmContent.value = `确定要删除 "${row.name}" 吗？`
  confirmAction.value = async () => { await deleteParseTemplate(row.id); message.success('删除成功'); tableRef.value?.loadData() }
  confirmDialogShow.value = true
}

function handleTest(row: ParseTemplate) {
  editingRow.value = row
  testSample.value = ''
  testResult.value = null
  testTabActive.value = 'test'
}

async function handleRunTest() {
  if (!editingRow.value || !testSample.value) { message.warning('请输入测试数据'); return }
  testLoading.value = true
  try {
    const res = await testParseTemplate(editingRow.value.id, testSample.value)
    testResult.value = res.data
  } catch (err: any) { message.error(err?.message || '测试失败') }
  finally { testLoading.value = false }
}

async function handleConfirm() {
  confirmLoading.value = true
  try { await confirmAction.value() } catch (err: any) { message.error(err?.message || '操作失败') } finally { confirmLoading.value = false; confirmDialogShow.value = false }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="解析模板" description="管理日志解析模板"><n-button type="primary" @click="handleAdd">添加模板</n-button></PageHeader>

    <n-tabs v-model:value="testTabActive" type="line">
      <n-tab-pane name="list" tab="模板列表">
        <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['name']" search-placeholder="搜索模板名称" />
      </n-tab-pane>
      <n-tab-pane name="test" tab="解析测试" :disabled="!editingRow">
        <n-card v-if="editingRow && testTabActive" size="small">
          <h3>测试模板: {{ editingRow.name }}</h3>
          <n-input v-model:value="testSample" type="textarea" placeholder="输入测试日志..." :autosize="{ minRows: 4, maxRows: 10 }" style="margin-bottom: 12px" />
          <n-button type="primary" :loading="testLoading" @click="handleRunTest">执行测试</n-button>
          <div v-if="testResult" style="margin-top: 16px">
            <n-tag :type="testResult.success ? 'success' : 'error'" size="small">{{ testResult.success ? '解析成功' : '解析失败' }}</n-tag>
            <div v-if="testResult.errors?.length" style="margin-top: 8px; color: var(--error-color)">
              <div v-for="(err, i) in testResult.errors" :key="i">{{ err }}</div>
            </div>
            <pre v-if="testResult.parsed" style="margin-top: 8px; padding: 12px; background: var(--bg-color); border-radius: 4px; overflow-x: auto">{{ JSON.stringify(testResult.parsed, null, 2) }}</pre>
          </div>
        </n-card>
      </n-tab-pane>
    </n-tabs>

    <FormDialog ref="formDialogRef" :title="editingRow ? '编辑模板' : '添加模板'" :fields="formFields" :width="700" @submit="handleFormSubmit" />
    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>