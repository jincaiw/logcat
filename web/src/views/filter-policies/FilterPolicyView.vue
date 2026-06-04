<script setup lang="ts">
import { computed, ref, h, watch } from 'vue'
import {
  NButton, NSpace, NTabs, NTabPane, NInput, NCard, NTag, NTable,
  NModal, NForm, NFormItem, NInputNumber, NSelect, NSwitch, NEmpty,
} from 'naive-ui'
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import { createFilterPolicy, updateFilterPolicy, deleteFilterPolicy, getFilterPolicies, testFilterPolicy } from '@/api/filterPolicies'
import { getAllDevices, getAllDeviceGroups } from '@/api/devices'
import { getParseTemplates } from '@/api/parseTemplates'
import type { FilterPolicy, FilterTestResult, FilterCondition, Device, DeviceGroup, ParseTemplate } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAppMessage } from '@/composables/useMessage'
import { useIsMobile } from '@/composables/useIsMobile'

const message = useAppMessage()
const { isMobile } = useIsMobile()
const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const confirmDialogShow = ref(false)
const confirmTitle = ref('')
const confirmContent = ref('')
const confirmAction = ref<() => Promise<void>>(() => Promise.resolve())
const confirmLoading = ref(false)
const editingRow = ref<FilterPolicy | null>(null)
const activeTab = ref('list')
const testSample = ref('')
const testResult = ref<FilterTestResult | null>(null)
const testLoading = ref(false)
const conditionExample = '[{"field":"severity","operator":"equals","value":"high"},{"field":"count","operator":"gte","value":"10"}]'
const testInputPlaceholder = '{"message":"Blocked By Policy","severity":"high","count":42,"source_ip":"10.0.0.5"}'

// ==================== Form State ====================
const formShow = ref(false)
const formRef = ref<FormInst | null>(null)
const formLoading = ref(false)
const formData = ref<Record<string, any>>({})
const conditions = ref<FilterCondition[]>([])

// ==================== Select Options ====================
const deviceOptions = ref<{ label: string; value: number }[]>([])
const deviceGroupOptions = ref<{ label: string; value: number }[]>([])
const parseTemplateOptions = ref<{ label: string; value: number }[]>([])

const operatorOptions = [
  { label: '等于 (equals)', value: 'equals' },
  { label: '不等于 (not_equals)', value: 'not_equals' },
  { label: '包含 (contains)', value: 'contains' },
  { label: '不包含 (not_contains)', value: 'not_contains' },
  { label: '开头为 (starts_with)', value: 'starts_with' },
  { label: '结尾为 (ends_with)', value: 'ends_with' },
  { label: '在列表中 (in)', value: 'in' },
  { label: '不在列表中 (not_in)', value: 'not_in' },
  { label: '正则匹配 (regex)', value: 'regex' },
  { label: '大于 (gt)', value: 'gt' },
  { label: '小于 (lt)', value: 'lt' },
  { label: '大于等于 (gte)', value: 'gte' },
  { label: '小于等于 (lte)', value: 'lte' },
  { label: '存在 (exists)', value: 'exists' },
  { label: '不存在 (not_exists)', value: 'not_exists' },
]

const actionOptions = [
  { label: '保留 (keep)', value: 'keep' },
  { label: '丢弃 (discard)', value: 'drop' },
]

const conditionLogicOptions = [
  { label: 'AND', value: 'AND' },
  { label: 'OR', value: 'OR' },
]

const formRules: FormRules = {
  name: [{ required: true, message: '请输入策略名称', trigger: ['blur', 'change'] }],
  priority: [{ type: 'number', required: true, message: '请输入优先级', trigger: ['blur', 'change'] }],
  action: [{ required: true, message: '请选择动作', trigger: ['blur', 'change'] }],
}

// ==================== Computed ====================
const resultRows = computed(() => {
  if (!testResult.value) return []
  return [
    { label: '命中结果', value: testResult.value.matched ? '命中' : '未命中' },
    { label: '执行动作', value: testResult.value.action === 'keep' ? '保留' : '丢弃' },
    { label: '白名单结果', value: testResult.value.whitelistResult || '--' },
    { label: '说明', value: testResult.value.message || '--' },
    { label: '命中策略', value: testResult.value.policy?.name || editingRow.value?.name || '--' },
  ]
})

// ==================== Fetch Select Options ====================
async function fetchDeviceOptions() {
  try {
    const res = await getAllDevices()
    deviceOptions.value = (res.data || []).map((d: Device) => ({ label: d.name, value: d.id }))
  } catch { deviceOptions.value = [] }
}

async function fetchDeviceGroupOptions() {
  try {
    const res = await getAllDeviceGroups()
    deviceGroupOptions.value = (res.data || []).map((d: DeviceGroup) => ({ label: d.name, value: d.id }))
  } catch { deviceGroupOptions.value = [] }
}

async function fetchParseTemplateOptions() {
  try {
    const res = await getParseTemplates({ page: 1, pageSize: 999 })
    const list = res.data?.list || []
    parseTemplateOptions.value = list.map((t: ParseTemplate) => ({ label: t.name, value: t.id }))
  } catch { parseTemplateOptions.value = [] }
}

function fetchSelectOptions() {
  fetchDeviceOptions()
  fetchDeviceGroupOptions()
  fetchParseTemplateOptions()
}

// ==================== Condition Builder ====================
function addCondition() {
  conditions.value.push({ field: '', operator: 'equals', value: '' })
}

function removeCondition(index: number) {
  conditions.value.splice(index, 1)
}

// ==================== Table Columns ====================
const columns: DataTableColumns<FilterPolicy> = [
  { title: '名称', key: 'name' },
  { title: '优先级', key: 'priority' },
  { title: '动作', key: 'action', render(row) { return h(NTag, { type: row.action === 'keep' ? 'success' : 'error', size: 'small', bordered: false }, { default: () => row.action === 'keep' ? '保留' : '丢弃' }) } },
  { title: '状态', key: 'enabled', render(row) { return h(NTag, { type: row.enabled ? 'success' : 'default', size: 'small', bordered: false }, { default: () => row.enabled ? '启用' : '禁用' }) } },
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

// ==================== Data Operations ====================
async function fetchData(params: any) { return getFilterPolicies(params) }

function handleAdd() {
  editingRow.value = null
  formData.value = {
    name: '', priority: 100, enabled: true, action: 'keep',
    conditionLogic: 'AND', whitelistEnabled: false, whitelistField: '', whitelistValues: '',
    deviceId: null, deviceGroupId: null, parseTemplateId: null,
    dedupEnabled: false, dedupWindow: 300,
  }
  conditions.value = []
  formShow.value = true
  fetchSelectOptions()
}

function handleEdit(row: FilterPolicy) {
  editingRow.value = row
  formData.value = {
    name: row.name, priority: row.priority,
    enabled: row.enabled, action: row.action,
    conditionLogic: row.conditionLogic || 'AND',
    whitelistEnabled: row.whitelistEnabled || false,
    whitelistField: row.whitelistField || '',
    whitelistValues: row.whitelistValues || '',
    deviceId: row.deviceId || null,
    deviceGroupId: row.deviceGroupId || null,
    parseTemplateId: row.parseTemplateId || null,
    dedupEnabled: row.dedupEnabled || false,
    dedupWindow: row.dedupWindow || 300,
  }
  // Parse conditions from JSON string or array
  let parsedConditions: FilterCondition[] = []
  if (row.conditions) {
    if (typeof row.conditions === 'string') {
      try { parsedConditions = JSON.parse(row.conditions) } catch { parsedConditions = [] }
    } else if (Array.isArray(row.conditions)) {
      parsedConditions = row.conditions
    }
  }
  conditions.value = parsedConditions.length > 0 ? parsedConditions : []
  formShow.value = true
  fetchSelectOptions()
}

async function handleFormSubmit() {
  try {
    await formRef.value?.validate()
    const payload: Record<string, any> = { ...formData.value }

    // Validate whitelist
    if (payload.whitelistEnabled && !String(payload.whitelistField || '').trim()) {
      message.warning('启用白名单时必须填写白名单字段')
      return
    }

    // Serialize conditions - filter out rows with empty field
    const validConditions = conditions.value.filter(c => c.field?.trim() !== '')
    payload.conditions = JSON.stringify(validConditions)

    // Clean up null association fields
    if (!payload.deviceId) delete payload.deviceId
    if (!payload.deviceGroupId) delete payload.deviceGroupId
    if (!payload.parseTemplateId) delete payload.parseTemplateId

    // Clean up dedup fields when disabled
    if (!payload.dedupEnabled) {
      delete payload.dedupWindow
    }

    formLoading.value = true
    if (editingRow.value) {
      await updateFilterPolicy(editingRow.value.id, payload)
      message.success('更新成功')
    } else {
      await createFilterPolicy(payload)
      message.success('创建成功')
    }
    formShow.value = false
    tableRef.value?.loadData()
  } catch (err: any) {
    if (err?.message) message.error(err.message)
  } finally {
    formLoading.value = false
  }
}

function handleDelete(row: FilterPolicy) {
  confirmTitle.value = '删除'
  confirmContent.value = `确定要删除 "${row.name}" 吗？`
  confirmAction.value = async () => { await deleteFilterPolicy(row.id); message.success('删除成功'); tableRef.value?.loadData() }
  confirmDialogShow.value = true
}

// ==================== Test ====================
const testSelectedPolicyId = ref<number | null>(null)
const policyOptions = ref<{ label: string; value: number }[]>([])

async function fetchPolicyOptions() {
  try {
    const res = await getFilterPolicies({ page: 1, pageSize: 999 })
    const list = res.data?.list || res.data?.items || []
    policyOptions.value = list.map((p: FilterPolicy) => ({ label: p.name, value: p.id }))
  } catch { policyOptions.value = [] }
}

watch(activeTab, (val) => {
  if (val === 'test') fetchPolicyOptions()
})

function handleTest(row: FilterPolicy) {
  editingRow.value = row
  testSelectedPolicyId.value = row.id
  testSample.value = ''
  testResult.value = null
  activeTab.value = 'test'
}

function onTestPolicySelect(id: number) {
  const policy = policyOptions.value.find(p => p.value === id)
  if (policy) {
    editingRow.value = { id: policy.value, name: policy.label } as FilterPolicy
  }
  testResult.value = null
}
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

    <n-tabs v-model:value="activeTab" type="line">
      <n-tab-pane name="list" tab="策略列表">
        <DataTable ref="tableRef" :columns="columns" :fetch-api="fetchData" :search-fields="['name']" search-placeholder="搜索策略名称" />
      </n-tab-pane>
      <n-tab-pane name="test" tab="策略测试">
        <n-card size="small">
          <n-select
            v-model:value="testSelectedPolicyId"
            :options="policyOptions"
            placeholder="选择要测试的策略"
            filterable
            clearable
            style="margin-bottom: 12px"
            @update:value="onTestPolicySelect"
          />
          <template v-if="editingRow">
            <h3>测试策略: {{ editingRow.name }}</h3>
            <div style="margin-bottom: 8px; color: var(--text-color-secondary)">
              条件示例：
              <code>{{ conditionExample }}</code>
            </div>
            <n-input v-model:value="testSample" type="textarea" :placeholder="testInputPlaceholder" :autosize="{ minRows: 4, maxRows: 10 }" style="margin-bottom: 12px" />
            <n-button type="primary" :loading="testLoading" :disabled="!testSelectedPolicyId" @click="handleRunTest">执行测试</n-button>
            <div v-if="testResult" style="margin-top: 16px">
              <n-tag :type="testResult.matched ? 'success' : 'warning'" size="small">
                匹配操作: {{ testResult.action === 'keep' ? '保留' : '丢弃' }} | 匹配结果: {{ testResult.matched ? '匹配' : '未匹配' }}
              </n-tag>
              <n-table :single-line="false" size="small" style="margin-top: 12px">
                <tbody>
                  <tr v-for="item in resultRows" :key="item.label">
                    <td style="width: 30%">{{ item.label }}</td>
                    <td>{{ item.value }}</td>
                  </tr>
                </tbody>
              </n-table>
              <pre style="margin-top: 8px; padding: 12px; background: var(--bg-color); border-radius: 4px; overflow-x: auto">{{ JSON.stringify(testResult, null, 2) }}</pre>
            </div>
          </template>
          <n-empty v-else description="请选择一个策略进行测试" style="padding: 40px 0" />
        </n-card>
      </n-tab-pane>
    </n-tabs>

    <!-- Custom Form Modal -->
    <n-modal
      v-model:show="formShow"
      :title="editingRow ? '编辑策略' : '添加策略'"
      preset="card"
      :style="{ width: isMobile ? 'calc(100vw - 32px)' : '720px', maxWidth: 'calc(100vw - 32px)' }"
      :mask-closable="false"
    >
      <n-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        :label-width="isMobile ? undefined : 120"
        :label-placement="isMobile ? 'top' : 'left'"
      >
        <n-form-item label="策略名称" path="name">
          <n-input v-model:value="formData.name" placeholder="请输入策略名称" />
        </n-form-item>

        <n-form-item label="优先级" path="priority">
          <n-input-number v-model:value="formData.priority" placeholder="请输入优先级" style="width: 100%" />
        </n-form-item>

        <n-form-item label="状态">
          <n-switch v-model:value="formData.enabled" />
        </n-form-item>

        <n-form-item label="动作" path="action">
          <n-select v-model:value="formData.action" :options="actionOptions" placeholder="请选择动作" />
        </n-form-item>

        <n-form-item label="关联设备">
          <n-select v-model:value="formData.deviceId" :options="deviceOptions" placeholder="选择设备（可选）" clearable filterable />
        </n-form-item>

        <n-form-item label="关联设备组">
          <n-select v-model:value="formData.deviceGroupId" :options="deviceGroupOptions" placeholder="选择设备组（可选）" clearable filterable />
        </n-form-item>

        <n-form-item label="关联解析模板">
          <n-select v-model:value="formData.parseTemplateId" :options="parseTemplateOptions" placeholder="选择解析模板（可选）" clearable filterable />
        </n-form-item>

        <n-form-item label="条件逻辑">
          <n-select v-model:value="formData.conditionLogic" :options="conditionLogicOptions" />
        </n-form-item>

        <!-- Condition Builder -->
        <n-form-item label="过滤条件">
          <div style="width: 100%">
            <div v-for="(cond, index) in conditions" :key="index" style="display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 8px; align-items: center;">
              <n-input v-model:value="cond.field" placeholder="字段名" style="flex: 2" />
              <n-select v-model:value="cond.operator" :options="operatorOptions" placeholder="运算符" style="flex: 2" />
              <n-input v-model:value="cond.value" placeholder="值" style="flex: 2" />
              <n-button size="small" type="error" ghost @click="removeCondition(index)">删除</n-button>
            </div>
            <n-button size="small" dashed @click="addCondition">+ 添加条件</n-button>
          </div>
        </n-form-item>

        <n-form-item label="启用白名单">
          <n-switch v-model:value="formData.whitelistEnabled" />
        </n-form-item>

        <n-form-item v-if="formData.whitelistEnabled" label="白名单字段">
          <n-input v-model:value="formData.whitelistField" placeholder="请输入白名单字段" />
        </n-form-item>

        <n-form-item v-if="formData.whitelistEnabled" label="白名单值(逗号分隔)">
          <n-input v-model:value="formData.whitelistValues" type="textarea" placeholder="请输入白名单值" :autosize="{ minRows: 2, maxRows: 6 }" />
        </n-form-item>

        <n-form-item label="启用去重">
          <n-switch v-model:value="formData.dedupEnabled" />
        </n-form-item>

        <n-form-item v-if="formData.dedupEnabled" label="去重窗口(秒)">
          <n-input-number v-model:value="formData.dedupWindow" :min="1" placeholder="请输入去重窗口" style="width: 100%" />
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="formShow = false">取消</n-button>
          <n-button type="primary" :loading="formLoading" @click="handleFormSubmit">确定</n-button>
        </n-space>
      </template>
    </n-modal>

    <ConfirmDialog v-model:show="confirmDialogShow" :title="confirmTitle" :content="confirmContent" :loading="confirmLoading" @confirm="handleConfirm" />
  </div>
</template>
