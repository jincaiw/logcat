<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSelect, NInputNumber, NSwitch, NTag, NPopconfirm, NEmpty, NSpace, NRadioGroup, NRadio, NAlert, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { API } from '@/api'
import type { FilterPolicy, FilterCondition, WhitelistItem } from '@/types'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const message = useMessage()

const loading = ref(false)
const policies = ref<FilterPolicy[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref(t('filterPolicy.addDialogTitle'))
const parseTemplates = ref<any[]>([])
const devices = ref<any[]>([])

const formData = ref<Partial<FilterPolicy>>({
  id: 0, name: '', description: '', deviceId: 0, deviceGroupId: 0,
  parseTemplateId: 0, conditions: '', conditionLogic: 'AND' as const,
  whitelist: '', whitelistField: '', action: 'keep' as const,
  priority: 0, isActive: true, dedupEnabled: true, dedupWindow: 60,
  dropUnmatched: false, createdAt: '', updatedAt: '',
})

const conditions = ref<FilterCondition[]>([])
const whitelist = ref<WhitelistItem[]>([])
const newCondition = ref({ field: '', operator: 'equals', value: '' })
const newWhitelistItem = ref({ cidr: '', description: '', enabled: true })

const operators = computed(() => [
  { value: 'equals', label: t('filterPolicy.opEquals') },
  { value: 'not_equals', label: t('filterPolicy.opNotEquals') },
  { value: 'contains', label: t('filterPolicy.opContains') },
  { value: 'not_contains', label: t('filterPolicy.opNotContains') },
  { value: 'in', label: t('filterPolicy.opIn') },
  { value: 'not_in', label: t('filterPolicy.opNotIn') },
  { value: 'starts_with', label: t('filterPolicy.opStartsWith') },
  { value: 'ends_with', label: t('filterPolicy.opEndsWith') },
  { value: 'regex', label: t('filterPolicy.opRegex') },
  { value: 'exists', label: t('filterPolicy.opExists') },
  { value: 'not_exists', label: t('filterPolicy.opNotExists') },
  { value: 'gt', label: t('filterPolicy.opGt') },
  { value: 'gte', label: t('filterPolicy.opGte') },
  { value: 'lt', label: t('filterPolicy.opLt') },
  { value: 'lte', label: t('filterPolicy.opLte') },
])

const actionOptions = [
  { value: 'keep', label: t('filterPolicy.actionKeep') },
  { value: 'discard', label: t('filterPolicy.actionDiscard') },
]

const availableFields = computed(() => {
  if (!formData.value.parseTemplateId) return []
  const template = parseTemplates.value.find((tm: any) => tm.id === formData.value.parseTemplateId)
  if (!template) return []
  if (template.fieldMapping) {
    try {
      const mapping = JSON.parse(template.fieldMapping)
      return Object.entries(mapping).map(([value, label]) => ({ value, label: `${label} (${value})` }))
    } catch { /* ignore */ }
  }
  return []
})

onMounted(async () => {
  await Promise.all([loadPolicies(), loadParseTemplates(), loadDevices()])
})

async function loadPolicies() {
  loading.value = true
  try {
    policies.value = await API.GetFilterPolicies()
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

async function loadParseTemplates() {
  try { parseTemplates.value = await API.GetParseTemplates() } catch (e) { console.error(e) }
}

async function loadDevices() {
  try { devices.value = await API.GetDevices() } catch (e) { console.error(e) }
}

function handleAdd() {
  dialogTitle.value = t('filterPolicy.addDialogTitle')
  const maxPriority = policies.value.length > 0 ? Math.max(...policies.value.map(p => p.priority)) : 0
  formData.value = {
    id: 0, name: '', description: '', deviceId: 0, deviceGroupId: 0,
    parseTemplateId: 0, conditions: '', conditionLogic: 'AND',
    whitelist: '', whitelistField: '', action: 'keep',
    priority: maxPriority + 1, isActive: true, dedupEnabled: true,
    dedupWindow: 60, dropUnmatched: false, createdAt: '', updatedAt: '',
  }
  conditions.value = []
  whitelist.value = []
  dialogVisible.value = true
}

function handleEdit(row: FilterPolicy) {
  dialogTitle.value = t('filterPolicy.editPolicy')
  formData.value = { ...row }
  try { conditions.value = row.conditions ? JSON.parse(row.conditions) : [] } catch { conditions.value = [] }
  try { whitelist.value = row.whitelist ? JSON.parse(row.whitelist) : [] } catch { whitelist.value = [] }
  dialogVisible.value = true
}

async function handleDelete(row: FilterPolicy) {
  try {
    await API.DeleteFilterPolicy(row.id!)
    message.success(t('message.deleteSuccess'))
    loadPolicies()
  } catch (e) { message.error(t('message.deleteFailed')) }
}

function addCondition() {
  if (!newCondition.value.field) { message.warning(t('filterPolicy.pleaseInputFieldName')); return }
  conditions.value.push({ ...newCondition.value })
  newCondition.value = { field: '', operator: 'equals', value: '' }
}

function removeCondition(index: number) { conditions.value.splice(index, 1) }

function addWhitelistItem() {
  if (!newWhitelistItem.value.cidr) { message.warning(t('filterPolicy.pleaseInputIpOrCidr')); return }
  whitelist.value.push({ ...newWhitelistItem.value })
  newWhitelistItem.value = { cidr: '', description: '', enabled: true }
}

function removeWhitelistItem(index: number) { whitelist.value.splice(index, 1) }

async function handleSubmit() {
  if (!formData.value.name) { message.warning(t('filterPolicy.pleaseInputName')); return }
  formData.value.conditions = JSON.stringify(conditions.value)
  formData.value.whitelist = JSON.stringify(whitelist.value)
  try {
    if (formData.value.id) {
      await API.UpdateFilterPolicy(formData.value as any)
      message.success(t('message.updateSuccess'))
    } else {
      await API.AddFilterPolicy(formData.value as any)
      message.success(t('message.addSuccess'))
    }
    dialogVisible.value = false
    loadPolicies()
  } catch (e) { message.error(t('message.operationFailed')) }
}

function getParseTemplateName(id: number): string {
  return parseTemplates.value.find((tm: any) => tm.id === id)?.name || '-'
}

function getDeviceName(id: number): string {
  if (id === 0) return t('common.allDevices')
  return devices.value.find((d: any) => d.id === id)?.name || '-'
}

function getActionText(action: string): string {
  return actionOptions.find(a => a.value === action)?.label || action
}

const columns: DataTableColumns<FilterPolicy> = [
  { title: t('common.id'), key: 'id', width: 70 },
  { title: t('filterPolicy.policyName'), key: 'name', width: 160, ellipsis: { tooltip: true } },
  { title: t('filterPolicy.parseTemplate'), key: 'parseTemplateId', width: 140, render(row) { return getParseTemplateName(row.parseTemplateId) } },
  { title: t('filterPolicy.device'), key: 'deviceId', width: 100, render(row) { return getDeviceName(row.deviceId) } },
  {
    title: t('filterPolicy.action'), key: 'action', width: 80, align: 'center',
    render(row) { return h(NTag, { type: row.action === 'keep' ? 'success' : 'error', size: 'small' }, { default: () => getActionText(row.action) }) },
  },
  { title: t('filterPolicy.priority'), key: 'priority', width: 70, align: 'center' },
  { title: t('common.description'), key: 'description', ellipsis: { tooltip: true } },
  {
    title: t('common.status'), key: 'isActive', width: 70, align: 'center',
    render(row) { return h(NTag, { type: row.isActive ? 'success' : 'error', size: 'small' }, { default: () => row.isActive ? t('common.enable') : t('common.disable') }) },
  },
  {
    title: t('common.action'), key: 'actions', width: 120,
    render(row) {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { text: true, type: 'primary', size: 'small', onClick: () => handleEdit(row) }, { default: () => t('common.edit') }),
          h(NPopconfirm, { onPositiveClick: () => handleDelete(row) }, {
            trigger: () => h(NButton, { text: true, type: 'error', size: 'small' }, { default: () => t('common.delete') }),
            default: () => t('filterPolicy.deleteConfirm'),
          }),
        ],
      })
    },
  },
]
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('filterPolicy.title') }}</h1>
        <p class="page-subtitle text-muted">{{ t('filterPolicy.subtitle') }}</p>
      </div>
      <NSpace class="page-actions">
        <NButton type="primary" @click="handleAdd">
          {{ t('filterPolicy.addPolicy') }}
        </NButton>
      </NSpace>
    </div>

    <div class="table-card mt-4">
      <NDataTable
        :columns="columns"
        :data="policies"
        :loading="loading"
        :bordered="false"
        striped
      />
      <div v-if="!policies.length && !loading" class="empty-state">
        <NEmpty :description="t('common.noDataDesc')" />
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <NModal
      v-model:show="dialogVisible"
      :title="dialogTitle"
      preset="card"
      class="modal-860"
      :bordered="true"
    >
      <NForm :model="formData" label-placement="left" :label-width="110">
        <NFormItem :label="t('filterPolicy.policyName')" required>
          <NInput v-model:value="formData.name" :placeholder="t('filterPolicy.pleaseInputName')" />
        </NFormItem>
        <NFormItem :label="t('filterPolicy.parseTemplate')">
          <NSelect v-model:value="formData.parseTemplateId" :placeholder="t('filterPolicy.selectParseTemplate')" clearable class="field-full"
            :options="parseTemplates.map((tm: any) => ({ label: tm.name, value: tm.id }))" />
        </NFormItem>
        <NFormItem :label="t('filterPolicy.device')">
          <NSelect v-model:value="formData.deviceId" :placeholder="t('filterPolicy.selectDevice')" clearable class="field-full"
            :options="[{ label: t('common.allDevices'), value: 0 }, ...devices.map((d: any) => ({ label: d.name, value: d.id }))]" />
        </NFormItem>

        <!-- Conditions -->
        <NFormItem :label="t('filterPolicy.conditions')">
          <div class="field-full">
            <div v-if="!formData.parseTemplateId" class="mb-4">
              <NAlert type="info" :bordered="false">{{ t('filterPolicy.pleaseSelectParseTemplateFirst') }}</NAlert>
            </div>
            <div class="condition-input flex gap-3 mb-4">
              <NSelect v-model:value="newCondition.field" :placeholder="t('filterPolicy.selectField')" :options="availableFields" class="field-control-200" filterable clearable :disabled="!formData.parseTemplateId" />
              <NSelect v-model:value="newCondition.operator" :options="operators" class="field-control-130" />
              <NInput v-model:value="newCondition.value" :placeholder="t('filterPolicy.value')" class="field-control-200" />
              <NButton type="primary" @click="addCondition">{{ t('common.add') }}</NButton>
            </div>
            <NRadioGroup v-if="conditions.length > 0" v-model:value="formData.conditionLogic" size="small" class="mb-4">
              <NRadio value="AND">{{ t('filterPolicy.conditionLogicAll') }}</NRadio>
              <NRadio value="OR">{{ t('filterPolicy.conditionLogicAny') }}</NRadio>
            </NRadioGroup>
            <div v-for="(cond, idx) in conditions" :key="idx" class="condition-item flex items-center justify-between gap-3 mb-4">
              <NSpace align="center">
                <NTag type="info" size="small">{{ cond.field }}</NTag>
                <span class="text-muted">{{ operators.find(o => o.value === cond.operator)?.label }}</span>
                <span class="mono">{{ cond.value || '-' }}</span>
              </NSpace>
              <NButton text type="error" size="small" @click="removeCondition(idx)">{{ t('common.delete') }}</NButton>
            </div>
          </div>
        </NFormItem>

        <!-- Whitelist -->
        <NFormItem :label="t('filterPolicy.whitelistConfig')">
          <div class="field-full">
            <p class="field-note mb-4">{{ t('filterPolicy.whitelistTip') }}</p>
            <div class="flex gap-3 mb-4">
              <NInput v-model:value="newWhitelistItem.cidr" :placeholder="t('filterPolicy.ipCidrPlaceholder')" class="field-control-260" />
              <NInput v-model:value="newWhitelistItem.description" :placeholder="t('filterPolicy.descriptionPlaceholder')" class="field-control-150" />
              <NSwitch v-model:value="newWhitelistItem.enabled" />
              <NButton type="primary" size="small" @click="addWhitelistItem">{{ t('common.add') }}</NButton>
            </div>
            <div v-for="(item, idx) in whitelist" :key="idx" class="condition-item flex items-center justify-between gap-3 mb-4">
              <NSpace align="center">
                <span class="mono text-accent">{{ item.cidr }}</span>
                <span class="text-muted">{{ item.description || '-' }}</span>
                <NSwitch v-model:value="item.enabled" size="small" />
              </NSpace>
              <NButton text type="error" size="small" @click="removeWhitelistItem(idx)">{{ t('common.delete') }}</NButton>
            </div>
            <p v-if="!whitelist.length" class="field-note">{{ t('filterPolicy.noWhitelist') }}</p>
          </div>
        </NFormItem>

        <!-- Action / Priority / DropUnmatched -->
        <NFormItem :label="t('filterPolicy.matchAction')">
          <NSelect v-model:value="formData.action" :options="actionOptions" class="field-control-200" />
        </NFormItem>
        <NFormItem :label="t('filterPolicy.dropUnmatched')">
          <NSwitch v-model:value="formData.dropUnmatched" />
        </NFormItem>
        <NFormItem :label="t('filterPolicy.priority')">
          <NInputNumber v-model:value="formData.priority" :min="0" :max="100" class="field-control-200" />
        </NFormItem>
        <NFormItem :label="t('common.description')">
          <NInput v-model:value="formData.description" type="textarea" :rows="2" :placeholder="t('common.pleaseInputDescription')" />
        </NFormItem>
        <NFormItem :label="t('common.status')">
          <NSwitch v-model:value="formData.isActive" />
        </NFormItem>
        <NFormItem :label="t('filterPolicy.dedupEnabled')">
          <NSpace align="center">
            <NSwitch v-model:value="formData.dedupEnabled" />
            <span class="text-muted">{{ formData.dedupEnabled ? t('common.enabled') : t('common.disabled') }}</span>
            <NInputNumber v-if="formData.dedupEnabled" v-model:value="formData.dedupWindow" :min="10" :max="3600" :step="10" size="small" class="field-control-120" />
            <span v-if="formData.dedupEnabled" class="text-muted">{{ t('common.second') }}</span>
          </NSpace>
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="dialogVisible = false">{{ t('common.cancel') }}</NButton>
          <NButton type="primary" @click="handleSubmit">{{ t('common.confirmButtonText') }}</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>
