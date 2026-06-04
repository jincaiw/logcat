<script setup lang="ts">
import { computed, ref, h, watch } from 'vue'
import { NButton, NSpace, NTabs, NTabPane, NInput, NCard, NTag, NTable, NSwitch, NSelect, NEmpty } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createParseTemplate, updateParseTemplate, deleteParseTemplate, getParseTemplates, testParseTemplate } from '@/api/parseTemplates'
import type { ParseTemplate, ParseTestResult } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import FormDialog, { type FieldConfig } from '@/components/common/FormDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAppMessage } from '@/composables/useMessage'

const message = useAppMessage()
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

const parseTypeExamples: Record<string, string> = {
  regex: '{"src_ip":"source_ip","severity":"severity"}',
  json: '{"event":"event_type","severity":"severity"}',
  kv: '["source_ip","severity","count"]',
  delimiter: '["source_ip","severity","count"]',
  syslog_json: '{"event":"event_type"}',
  sub_template: '[{"templateId":1,"matchType":"contains","matchField":"raw_message","matchValue":"event="}]',
}

const parseFields = computed(() => Object.entries(testResult.value?.fields || {}))
const testInputPlaceholder = computed(() => {
  if (!editingRow.value) return '输入测试日志...'
  if (editingRow.value.sampleLog) return editingRow.value.sampleLog
  switch (editingRow.value.parseType) {
    case 'json':
      return '{"event":"login","severity":"high"}'
    case 'syslog_json':
      return '<34>1 2026-05-28T10:00:00Z host app - - - {"event":"login","severity":"high"}'
    case 'kv':
      return 'src_ip=10.0.0.8 severity=high count=42'
    case 'delimiter':
      return '10.0.0.8|high|42'
    case 'sub_template':
      return '{"event":"login","result":"success"}'
    default:
      return 'src=10.0.0.8 severity=high count=42'
  }
})

const formFields: FieldConfig[] = [
  { key: 'name', label: '模板名称', type: 'text', required: true },
  { key: 'deviceType', label: '设备类型', type: 'select', options: [
    { label: '防火墙', value: 'firewall' }, { label: 'WAF', value: 'waf' },
    { label: 'EDR', value: 'edr' }, { label: 'IDS', value: 'ids' },
    { label: 'IPS', value: 'ips' }, { label: 'SIEM', value: 'siem' },
    { label: '服务器', value: 'server' }, { label: '应用', value: 'application' },
    { label: '其他', value: 'other' },
  ] },
  { key: 'parseType', label: '类型', type: 'select', required: true, options: [
    { label: 'Regex', value: 'regex' }, { label: 'JSON', value: 'json' },
    { label: 'KV', value: 'kv' }, { label: '分隔符', value: 'delimiter' },
    { label: 'Syslog+JSON', value: 'syslog_json' }, { label: '子模板路由', value: 'sub_template' },
  ], defaultValue: 'regex' },
  { key: 'enabled', label: '启用', type: 'switch', defaultValue: true },
  { key: 'headerRegex', label: 'Regex/Header', type: 'code', visible: (data) => ['regex', 'syslog_json'].includes(data.parseType) },
  { key: 'delimiter', label: '分隔符', type: 'text', visible: (data) => data.parseType === 'delimiter' },
  { key: 'subTemplates', label: '子模板路由(JSON)', type: 'code', visible: (data) => data.parseType === 'sub_template' },
  { key: 'fieldMapping', label: '字段映射(JSON)', type: 'code' },
  { key: 'valueTransform', label: '值转换(JSON)', type: 'code' },
  { key: 'sampleLog', label: '示例日志', type: 'textarea' },
]

const columns: DataTableColumns<ParseTemplate> = [
  { title: '名称', key: 'name' },
  {
    title: '类型', key: 'parseType',
    render(row) { return h(NTag, { size: 'small', bordered: false }, { default: () => row.parseType.toUpperCase() }) },
  },
  { title: '设备类型', key: 'deviceType' },
  {
    title: '启用', key: 'enabled',
    render(row) {
      return h(NSwitch, {
        value: row.enabled !== false,
        onUpdateValue: (val: boolean) => handleToggleEnabled(row, val),
      })
    },
  },
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
function handleEdit(row: ParseTemplate) {
  editingRow.value = row
  formDialogRef.value?.open({
    ...row,
    fieldMapping: row.fieldMapping || '',
    valueTransform: row.valueTransform || '',
    subTemplates: row.subTemplates || '',
    sampleLog: row.sampleLog || '',
    enabled: row.enabled !== false,
    deviceType: row.deviceType || '',
    headerRegex: row.headerRegex || '',
    delimiter: row.delimiter || '',
  })
}

function validateJSONField(value: unknown, label: string) {
  if (typeof value === 'string' && value.trim()) {
    try {
      JSON.parse(value)
    } catch {
      throw new Error(`${label} 不是合法 JSON`)
    }
  }
}

function validateParseTemplateForm(data: Record<string, any>) {
  if (!data.parseType) throw new Error('请选择解析类型')
  validateJSONField(data.fieldMapping, '字段映射')
  validateJSONField(data.valueTransform, '值转换')
  validateJSONField(data.subTemplates, '子模板路由')

  if (data.parseType === 'regex' && !String(data.headerRegex || '').trim()) {
    throw new Error('Regex 类型必须填写 Regex/Header')
  }
  if (data.parseType === 'delimiter' && !String(data.delimiter || '').trim()) {
    throw new Error('分隔符类型必须填写分隔符')
  }
  if (data.parseType === 'sub_template' && !String(data.subTemplates || '').trim()) {
    throw new Error('子模板路由类型必须填写子模板路由 JSON')
  }
}

async function handleFormSubmit(data: Record<string, any>) {
  try {
    validateParseTemplateForm(data)
    if (editingRow.value) { await updateParseTemplate(editingRow.value.id, data); message.success('更新成功') }
    else { await createParseTemplate(data); message.success('创建成功') }
    formDialogRef.value?.close(); tableRef.value?.loadData()
  } catch (err: any) { message.error(err?.message || '操作失败') }
}

async function handleToggleEnabled(row: ParseTemplate, enabled: boolean) {
  try {
    await updateParseTemplate(row.id, { ...row, enabled })
    message.success(enabled ? '已启用' : '已禁用')
    tableRef.value?.loadData()
  } catch (err: any) {
    message.error(err?.message || '操作失败')
  }
}

function handleDelete(row: ParseTemplate) {
  confirmTitle.value = '删除'; confirmContent.value = `确定要删除 "${row.name}" 吗？`
  confirmAction.value = async () => { await deleteParseTemplate(row.id); message.success('删除成功'); tableRef.value?.loadData() }
  confirmDialogShow.value = true
}

function handleTest(row: ParseTemplate) {
  editingRow.value = row
  testSelectedTemplateId.value = row.id
  testSample.value = row.sampleLog || ''
  testResult.value = null
  testTabActive.value = 'test'
}

const testSelectedTemplateId = ref<number | null>(null)
const templateOptions = ref<{ label: string; value: number }[]>([])

async function fetchTemplateOptions() {
  try {
    const res = await getParseTemplates({ page: 1, pageSize: 999 })
    const list = res.data?.list || res.data?.items || []
    templateOptions.value = list.map((t: ParseTemplate) => ({ label: t.name, value: t.id }))
  } catch { templateOptions.value = [] }
}

function onTestTemplateSelect(id: number) {
  const tmpl = templateOptions.value.find(t => t.value === id)
  if (tmpl) {
    editingRow.value = { id: tmpl.value, name: tmpl.label } as ParseTemplate
  }
  testResult.value = null
}

watch(testTabActive, (val) => {
  if (val === 'test') fetchTemplateOptions()
})

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
      <n-tab-pane name="test" tab="解析测试">
        <n-card size="small">
          <n-select
            v-model:value="testSelectedTemplateId"
            :options="templateOptions"
            placeholder="选择要测试的模板"
            filterable
            clearable
            style="margin-bottom: 12px"
            @update:value="onTestTemplateSelect"
          />
          <template v-if="editingRow">
            <h3>测试模板: {{ editingRow.name }}</h3>
            <div style="margin-bottom: 8px; color: var(--text-color-secondary)">
              配置示例：
              <code>{{ parseTypeExamples[editingRow.parseType] || '{}' }}</code>
            </div>
            <n-input v-model:value="testSample" type="textarea" :placeholder="testInputPlaceholder" :autosize="{ minRows: 4, maxRows: 10 }" style="margin-bottom: 12px" />
            <n-button type="primary" :loading="testLoading" :disabled="!testSelectedTemplateId" @click="handleRunTest">执行测试</n-button>
            <div v-if="testResult" style="margin-top: 16px">
              <n-tag :type="testResult.success ? 'success' : 'error'" size="small">{{ testResult.success ? '解析成功' : '解析失败' }}</n-tag>
              <div v-if="testResult.error" style="margin-top: 8px; color: var(--error-color)">
                {{ testResult.error }}
              </div>
              <n-table v-if="parseFields.length" :single-line="false" size="small" style="margin-top: 12px">
                <thead>
                  <tr>
                    <th style="width: 30%">字段</th>
                    <th>值</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="[key, value] in parseFields" :key="key">
                    <td>{{ key }}</td>
                    <td><code>{{ typeof value === 'string' ? value : JSON.stringify(value) }}</code></td>
                  </tr>
                </tbody>
              </n-table>
              <pre v-if="testResult.fields" style="margin-top: 8px; padding: 12px; background: var(--bg-color); border-radius: 4px; overflow-x: auto">{{ JSON.stringify(testResult.fields, null, 2) }}</pre>
            </div>
          </template>
          <n-empty v-else description="请选择一个模板进行测试" style="padding: 40px 0" />
        </n-card>
      </n-tab-pane>
    </n-tabs>

    <FormDialog ref="formDialogRef" :title="editingRow ? '编辑模板' : '添加模板'" :fields="formFields" :width="700" @submit="handleFormSubmit" />
    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>
