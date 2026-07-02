<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSelect, NTag, NPopconfirm, NEmpty, NTabs, NTabPane, NSpace, NSwitch, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { API } from '@/api'
import type { ParseTemplate, FieldMappingDoc } from '@/types'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const message = useMessage()
const activeTab = ref('templates')

// ==================== Parse Templates ====================
const templates = ref<ParseTemplate[]>([])
const templatesLoading = ref(false)
const templateDialogVisible = ref(false)
const templateDialogTitle = ref(t('parseTemplate.addTemplate'))
const templateForm = ref<Partial<ParseTemplate>>({
  name: '',
  description: '',
  parseType: 'regex',
  headerRegex: '',
  fieldMapping: '',
  valueTransform: '',
  sampleLog: '',
  deviceType: '',
  delimiter: '|!',
  typeField: 0,
  subTemplates: '',
  isActive: true,
})
const parseTestResult = ref<any>(null)
const parseTesting = ref(false)

const parseTypeOptions = [
  { label: t('parseTemplate.typeSyslogJson'), value: 'syslog_json' },
  { label: t('parseTemplate.typeJson'), value: 'json' },
  { label: t('parseTemplate.typeDelimiter'), value: 'delimiter' },
  { label: t('parseTemplate.typeSmartDelimiter'), value: 'smart_delimiter' },
  { label: t('parseTemplate.typeKeyvalue'), value: 'keyvalue' },
  { label: t('parseTemplate.typeRegex'), value: 'regex' },
  { label: t('parseTemplate.typeKv'), value: 'kv' },
]

async function loadTemplates() {
  templatesLoading.value = true
  try {
    templates.value = await API.GetParseTemplates()
  } catch (e) {
    console.error(e)
  } finally {
    templatesLoading.value = false
  }
}

function handleAddTemplate() {
  templateDialogTitle.value = t('parseTemplate.addTemplate')
  templateForm.value = {
    name: '', description: '', parseType: 'regex', headerRegex: '',
    fieldMapping: '', valueTransform: '', sampleLog: '',
    deviceType: '', delimiter: '|!', typeField: 0, subTemplates: '', isActive: true,
  }
  parseTestResult.value = null
  templateDialogVisible.value = true
}

function handleEditTemplate(row: ParseTemplate) {
  templateDialogTitle.value = t('parseTemplate.editTemplate')
  templateForm.value = { ...row }
  parseTestResult.value = null
  templateDialogVisible.value = true
}

async function handleDeleteTemplate(row: ParseTemplate) {
  try {
    await API.DeleteParseTemplate(row.id!)
    message.success(t('message.deleteSuccess'))
    loadTemplates()
  } catch (e) {
    message.error(t('message.deleteFailed'))
  }
}

async function handleSubmitTemplate() {
  if (!templateForm.value.name) {
    message.warning(t('parseTemplate.pleaseInputName'))
    return
  }
  try {
    if (templateForm.value.id) {
      await API.UpdateParseTemplate({ ...templateForm.value, id: templateForm.value.id! } as any)
      message.success(t('message.updateSuccess'))
    } else {
      await API.AddParseTemplate(templateForm.value as any)
      message.success(t('message.addSuccess'))
    }
    templateDialogVisible.value = false
    loadTemplates()
  } catch (e) {
    message.error(t('message.operationFailed'))
  }
}

async function handleTestParse() {
  if (!templateForm.value.sampleLog) {
    message.warning(t('parseTemplate.pleaseInputLogSample'))
    return
  }
  parseTesting.value = true
  parseTestResult.value = null
  try {
    const result = await API.TestParseTemplate({
      parseType: templateForm.value.parseType || 'regex',
      headerRegex: templateForm.value.headerRegex || '',
      fieldMapping: templateForm.value.fieldMapping || '',
      valueTransform: templateForm.value.valueTransform || '',
      sampleLog: templateForm.value.sampleLog,
    })
    parseTestResult.value = result
  } catch (e: any) {
    message.error(t('message.testFailed') + (e.message || e))
  } finally {
    parseTesting.value = false
  }
}

const templateColumns: DataTableColumns<ParseTemplate> = [
  { title: 'ID', key: 'id', width: 70 },
  { title: t('parseTemplate.name'), key: 'name', width: 160, ellipsis: { tooltip: true } },
  {
    title: t('parseTemplate.parseType'),
    key: 'parseType',
    width: 130,
    render(row) {
      return h(NTag, { size: 'small' }, { default: () => row.parseType })
    },
  },
  { title: t('common.description'), key: 'description', ellipsis: { tooltip: true } },
  {
    title: t('common.status'),
    key: 'isActive',
    width: 80,
    render(row) {
      return h(NTag, { type: row.isActive ? 'success' : 'error', size: 'small' }, {
        default: () => row.isActive ? t('common.enable') : t('common.disable'),
      })
    },
  },
  {
    title: t('common.action'),
    key: 'actions',
    width: 150,
    render(row) {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { text: true, type: 'primary', size: 'small', onClick: () => handleEditTemplate(row) }, { default: () => t('common.edit') }),
          h(NPopconfirm, { onPositiveClick: () => handleDeleteTemplate(row) }, {
            trigger: () => h(NButton, { text: true, type: 'error', size: 'small' }, { default: () => t('common.delete') }),
            default: () => t('parseTemplate.deleteConfirm'),
          }),
        ],
      })
    },
  },
]

// ==================== Field Mapping Docs ====================
const docs = ref<FieldMappingDoc[]>([])
const docsLoading = ref(false)
const docDialogVisible = ref(false)
const docDialogTitle = ref(t('fieldMappingDoc.addDocTitle'))
const isEditDoc = ref(false)
const docForm = ref<Partial<FieldMappingDoc>>({
  name: '',
  deviceType: '',
  description: '',
  fieldMappings: '',
  isActive: true,
})

async function loadDocs() {
  docsLoading.value = true
  try {
    docs.value = await API.GetFieldMappingDocs()
  } catch (e) {
    console.error(e)
  } finally {
    docsLoading.value = false
  }
}

function handleAddDoc() {
  isEditDoc.value = false
  docDialogTitle.value = t('fieldMappingDoc.addDocTitle')
  docForm.value = { name: '', deviceType: '', description: '', fieldMappings: '', isActive: true }
  docDialogVisible.value = true
}

function handleEditDoc(row: FieldMappingDoc) {
  isEditDoc.value = true
  docDialogTitle.value = t('fieldMappingDoc.editDoc')
  docForm.value = { ...row }
  docDialogVisible.value = true
}

async function handleDeleteDoc(row: FieldMappingDoc) {
  try {
    await API.DeleteFieldMappingDoc(row.id!)
    message.success(t('message.deleteSuccess'))
    loadDocs()
  } catch (e) {
    message.error(t('message.deleteFailed'))
  }
}

async function handleSubmitDoc() {
  if (!docForm.value.name || !docForm.value.deviceType) {
    message.warning(t('common.requiredFields'))
    return
  }
  try {
    if (isEditDoc.value) {
      await API.UpdateFieldMappingDoc({ ...docForm.value, id: docForm.value.id! } as any)
      message.success(t('message.updateSuccess'))
    } else {
      await API.AddFieldMappingDoc(docForm.value as any)
      message.success(t('message.addSuccess'))
    }
    docDialogVisible.value = false
    loadDocs()
  } catch (e) {
    message.error(t('message.operationFailed'))
  }
}

function getFieldCount(fieldMappings: string): number {
  if (!fieldMappings) return 0
  try {
    return Object.keys(JSON.parse(fieldMappings)).length
  } catch {
    return 0
  }
}

const docColumns: DataTableColumns<FieldMappingDoc> = [
  { title: 'ID', key: 'id', width: 70 },
  { title: t('fieldMappingDoc.docName'), key: 'name', width: 180 },
  {
    title: t('fieldMappingDoc.deviceType'),
    key: 'deviceType',
    width: 120,
    render(row) {
      return h(NTag, { type: 'info', size: 'small' }, { default: () => row.deviceType })
    },
  },
  { title: t('common.description'), key: 'description', ellipsis: { tooltip: true } },
  {
    title: t('fieldMappingDoc.fieldCount'),
    key: 'fieldMappings',
    width: 100,
    render(row) {
      return h(NTag, { size: 'small' }, { default: () => `${getFieldCount(row.fieldMappings)} ${t('fieldMappingDoc.fieldCountUnit')}` })
    },
  },
  {
    title: t('common.status'),
    key: 'isActive',
    width: 80,
    render(row) {
      return h(NTag, { type: row.isActive ? 'success' : 'error', size: 'small' }, {
        default: () => row.isActive ? t('common.enable') : t('common.disable'),
      })
    },
  },
  {
    title: t('common.action'),
    key: 'actions',
    width: 150,
    render(row) {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { text: true, type: 'primary', size: 'small', onClick: () => handleEditDoc(row) }, { default: () => t('common.edit') }),
          h(NPopconfirm, { onPositiveClick: () => handleDeleteDoc(row) }, {
            trigger: () => h(NButton, { text: true, type: 'error', size: 'small' }, { default: () => t('common.delete') }),
            default: () => t('fieldMappingDoc.deleteConfirm'),
          }),
        ],
      })
    },
  },
]

onMounted(async () => {
  await Promise.all([loadTemplates(), loadDocs()])
})
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('menu.logParser') }}</h1>
        <p class="page-subtitle text-muted">{{ t('parseTemplate.subtitle') }}</p>
      </div>
      <NSpace class="page-actions">
        <NButton v-if="activeTab === 'templates'" type="primary" @click="handleAddTemplate">
          {{ t('parseTemplate.addTemplate') }}
        </NButton>
        <NButton v-if="activeTab === 'mappings'" type="primary" @click="handleAddDoc">
          {{ t('fieldMappingDoc.addDoc') }}
        </NButton>
      </NSpace>
    </div>

    <div class="section-card mt-4 tabbed-panel">
      <NTabs v-model:value="activeTab" type="line" animated>
        <NTabPane name="templates" :tab="t('parseTemplate.title')">
          <div class="table-card">
            <NDataTable
              :columns="templateColumns"
              :data="templates"
              :loading="templatesLoading"
              :bordered="false"
              striped
            />
            <div v-if="!templates.length && !templatesLoading" class="empty-state">
              <NEmpty :description="t('common.noDataDesc')" />
            </div>
          </div>
        </NTabPane>
        <NTabPane name="mappings" :tab="t('fieldMappingDoc.title')">
          <div class="table-card">
            <NDataTable
              :columns="docColumns"
              :data="docs"
              :loading="docsLoading"
              :bordered="false"
              striped
            />
            <div v-if="!docs.length && !docsLoading" class="empty-state">
              <NEmpty :description="t('common.noDataDesc')" />
            </div>
          </div>
        </NTabPane>
      </NTabs>
    </div>

    <!-- Template Modal -->
    <NModal
      v-model:show="templateDialogVisible"
      :title="templateDialogTitle"
      preset="card"
      style="width: min(860px, 92vw)"
      :bordered="true"
    >
      <NForm :model="templateForm" label-placement="left" :label-width="120">
        <NFormItem :label="t('parseTemplate.name')" required>
          <NInput v-model:value="templateForm.name" :placeholder="t('parseTemplate.pleaseInputName')" />
        </NFormItem>
        <NFormItem :label="t('parseTemplate.parseType')">
          <NSelect v-model:value="templateForm.parseType" :options="parseTypeOptions" style="width: 100%" />
        </NFormItem>
        <NFormItem :label="t('parseTemplate.headerRegex')">
          <NInput v-model:value="templateForm.headerRegex" :placeholder="t('parseTemplate.headerRegexTip')" />
        </NFormItem>
        <NFormItem :label="t('parseTemplate.fieldMapping')">
          <NInput v-model:value="templateForm.fieldMapping" type="textarea" :rows="4" />
        </NFormItem>
        <NFormItem :label="t('parseTemplate.valueTransform')">
          <NInput v-model:value="templateForm.valueTransform" type="textarea" :rows="3" />
        </NFormItem>
        <NFormItem :label="t('parseTemplate.sampleLog')">
          <NInput v-model:value="templateForm.sampleLog" type="textarea" :rows="4" :placeholder="t('parseTemplate.pleaseInputSampleLog')" />
        </NFormItem>
        <NFormItem :label="t('parseTemplate.deviceType')">
          <NInput v-model:value="templateForm.deviceType" />
        </NFormItem>
        <NFormItem :label="t('parseTemplate.delimiter')">
          <NInput v-model:value="templateForm.delimiter" style="width: 200px" />
        </NFormItem>
        <NFormItem :label="t('common.description')">
          <NInput v-model:value="templateForm.description" type="textarea" :rows="2" :placeholder="t('common.pleaseInputDescription')" />
        </NFormItem>
        <NFormItem :label="t('common.status')">
          <NSwitch v-model:value="templateForm.isActive" />
        </NFormItem>
      </NForm>

      <!-- Test Parse Section -->
      <div class="mt-4">
        <NButton :loading="parseTesting" type="primary" secondary @click="handleTestParse">
          {{ t('parseTemplate.testParse') }}
        </NButton>
        <div v-if="parseTestResult" class="parse-result mt-4">
          <h4 class="text-secondary">{{ t('parseTemplate.parseResult') }}</h4>
          <NTag v-if="parseTestResult.success" type="success" size="small" class="mb-4">
            {{ t('parseTemplate.parseSuccess', { count: parseTestResult.fields?.length || 0 }) }}
          </NTag>
          <NTag v-else type="error" size="small" class="mb-4">
            {{ parseTestResult.error || t('parseTemplate.parseFailed') }}
          </NTag>
          <pre v-if="parseTestResult.data" class="log-entry mono">{{ JSON.stringify(parseTestResult.data, null, 2) }}</pre>
        </div>
      </div>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="templateDialogVisible = false">{{ t('common.cancel') }}</NButton>
          <NButton type="primary" @click="handleSubmitTemplate">{{ t('common.confirmButtonText') }}</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- Field Mapping Doc Modal -->
    <NModal
      v-model:show="docDialogVisible"
      :title="docDialogTitle"
      preset="card"
      style="width: min(640px, 92vw)"
      :bordered="true"
    >
      <NForm :model="docForm" label-placement="left" :label-width="100">
        <NFormItem :label="t('fieldMappingDoc.docName')" required>
          <NInput v-model:value="docForm.name" :placeholder="t('fieldMappingDoc.pleaseInputDocName')" />
        </NFormItem>
        <NFormItem :label="t('fieldMappingDoc.deviceType')" required>
          <NInput v-model:value="docForm.deviceType" :placeholder="t('fieldMappingDoc.deviceTypePlaceholder')" />
        </NFormItem>
        <NFormItem :label="t('common.description')">
          <NInput v-model:value="docForm.description" type="textarea" :rows="2" :placeholder="t('common.pleaseInputDescription')" />
        </NFormItem>
        <NFormItem :label="t('fieldMappingDoc.fieldMappings')">
          <NInput v-model:value="docForm.fieldMappings" type="textarea" :rows="8" placeholder="JSON format field mappings" />
        </NFormItem>
        <NFormItem :label="t('common.status')">
          <NSwitch v-model:value="docForm.isActive" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="docDialogVisible = false">{{ t('common.cancel') }}</NButton>
          <NButton type="primary" @click="handleSubmitDoc">{{ t('common.confirmButtonText') }}</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
.parse-result {
  padding: 16px;
  background: var(--bg-sunken);
  border-radius: 8px;
}
</style>
