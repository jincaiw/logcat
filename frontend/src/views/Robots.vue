<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSelect, NSwitch, NTag, NPopconfirm, NEmpty, NTabs, NTabPane, NSpace, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { API } from '@/api'
import type { Robot, OutputTemplate } from '@/types'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const message = useMessage()
const activeTab = ref('robots')

// ==================== Robots ====================
const robots = ref<Robot[]>([])
const robotsLoading = ref(false)
const robotDialogVisible = ref(false)
const robotDialogTitle = ref(t('robot.addRobot'))
const robotForm = ref<Partial<Robot>>({
  name: '',
  platform: 'feishu',
  webhookUrl: '',
  secret: '',
  isActive: true,
  description: '',
})

const robotTypeOptions = computed(() => [
  { label: t('robot.typeFeishu'), value: 'feishu' },
  { label: t('robot.typeWebhook'), value: 'webhook' },
  { label: t('robot.typeEmail'), value: 'email' },
])

const needsWebhook = computed(() => robotForm.value.platform !== 'email')

async function loadRobots() {
  robotsLoading.value = true
  try {
    robots.value = await API.GetRobots()
  } catch (e) { console.error(e) }
  finally { robotsLoading.value = false }
}

function handleAddRobot() {
  robotDialogTitle.value = t('robot.addRobot')
  robotForm.value = { name: '', platform: 'feishu', webhookUrl: '', secret: '', isActive: true, description: '' }
  robotDialogVisible.value = true
}

function handleEditRobot(row: Robot) {
  robotDialogTitle.value = t('robot.editRobot')
  robotForm.value = { ...row }
  robotDialogVisible.value = true
}

async function handleDeleteRobot(row: Robot) {
  try {
    await API.DeleteRobot(row.id!)
    message.success(t('message.deleteSuccess'))
    loadRobots()
  } catch (e) { message.error(t('message.deleteFailed')) }
}

async function handleSubmitRobot() {
  if (!robotForm.value.name) {
    message.warning(t('common.requiredFields'))
    return
  }
  if (needsWebhook.value && !robotForm.value.webhookUrl) {
    message.warning(t('common.requiredFields'))
    return
  }
  try {
    if (robotForm.value.id) {
      await API.UpdateRobot({ ...robotForm.value, id: robotForm.value.id! } as any)
      message.success(t('message.updateSuccess'))
    } else {
      await API.AddRobot(robotForm.value as any)
      message.success(t('message.addSuccess'))
    }
    robotDialogVisible.value = false
    loadRobots()
  } catch (e) { message.error(t('message.operationFailed')) }
}

async function handleTestRobot(row: Robot) {
  try {
    const result = await API.TestRobot(row)
    if (result.success) {
      message.success(t('robot.testSuccess'))
    } else {
      message.error(t('robot.testFailed'))
    }
  } catch (e: any) {
    message.error(t('robot.testFailed') + (e.message || e))
  }
}

const robotColumns: DataTableColumns<Robot> = [
  { title: t('common.id'), key: 'id', width: 70 },
  { title: t('robot.name'), key: 'name', width: 160, ellipsis: { tooltip: true } },
  {
    title: t('robot.type'), key: 'platform', width: 100,
    render(row) {
      return h(NTag, {
        type: row.platform === 'feishu' ? 'info' : row.platform === 'webhook' ? 'success' : 'default',
        size: 'small',
      }, { default: () => robotTypeOptions.value.find((o: any) => o.value === row.platform)?.label || row.platform })
  },
  },
  {
    title: t('robot.webhookUrl'), key: 'webhookUrl', ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'mono text-muted', style: { fontSize: '13px' } }, row.webhookUrl || '-')
    },
  },
  { title: t('common.description'), key: 'description', ellipsis: { tooltip: true } },
  {
    title: t('common.status'), key: 'isActive', width: 80, align: 'center',
    render(row) {
      return h(NTag, { type: row.isActive ? 'success' : 'error', size: 'small' }, {
        default: () => row.isActive ? t('common.enable') : t('common.disable'),
      })
    },
  },
  {
    title: t('common.action'), key: 'actions', width: 200,
    render(row) {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { text: true, type: 'primary', size: 'small', onClick: () => handleEditRobot(row) }, { default: () => t('common.edit') }),
          h(NButton, { text: true, type: 'info', size: 'small', onClick: () => handleTestRobot(row) }, { default: () => t('robot.testSend') }),
          h(NPopconfirm, { onPositiveClick: () => handleDeleteRobot(row) }, {
            trigger: () => h(NButton, { text: true, type: 'error', size: 'small' }, { default: () => t('common.delete') }),
            default: () => t('robot.deleteConfirm'),
          }),
        ],
      })
    },
  },
]

// ==================== Output Templates ====================
const templates = ref<OutputTemplate[]>([])
const templatesLoading = ref(false)
const templateDialogVisible = ref(false)
const templateDialogTitle = ref(t('outputTemplate.addTitle'))
const templateForm = ref<Partial<OutputTemplate>>({
  name: '',
  platform: 'feishu',
  content: '',
  description: '',
  isActive: true,
})

async function loadTemplates() {
  templatesLoading.value = true
  try {
    templates.value = await API.GetOutputTemplates()
  } catch (e) { console.error(e) }
  finally { templatesLoading.value = false }
}

function handleAddTemplate() {
  templateDialogTitle.value = t('outputTemplate.addTitle')
  templateForm.value = { name: '', platform: 'feishu', content: '', description: '', isActive: true }
  templateDialogVisible.value = true
}

function handleEditTemplate(row: OutputTemplate) {
  templateDialogTitle.value = t('outputTemplate.editTitle')
  templateForm.value = { ...row }
  templateDialogVisible.value = true
}

async function handleDeleteTemplate(row: OutputTemplate) {
  try {
    await API.DeleteOutputTemplate(row.id!)
    message.success(t('message.deleteSuccess'))
    loadTemplates()
  } catch (e) { message.error(t('message.deleteFailed')) }
}

async function handleSubmitTemplate() {
  if (!templateForm.value.name) {
    message.warning(t('outputTemplate.pleaseInputName'))
    return
  }
  try {
    if (templateForm.value.id) {
      await API.UpdateOutputTemplate({ ...templateForm.value, id: templateForm.value.id! } as any)
      message.success(t('message.updateSuccess'))
    } else {
      await API.AddOutputTemplate(templateForm.value as any)
      message.success(t('message.addSuccess'))
    }
    templateDialogVisible.value = false
    loadTemplates()
  } catch (e) { message.error(t('message.operationFailed')) }
}

const templateColumns: DataTableColumns<OutputTemplate> = [
  { title: t('common.id'), key: 'id', width: 70 },
  { title: t('outputTemplate.name'), key: 'name', width: 160, ellipsis: { tooltip: true } },
  {
    title: t('outputTemplate.type'), key: 'platform', width: 100,
    render(row) {
      return h(NTag, { type: 'info', size: 'small' }, {
        default: () => robotTypeOptions.value.find((o: any) => o.value === row.platform)?.label || row.platform,
      })
    },
  },
  { title: t('common.description'), key: 'description', ellipsis: { tooltip: true } },
  {
    title: t('common.status'), key: 'isActive', width: 80, align: 'center',
    render(row) {
      return h(NTag, { type: row.isActive ? 'success' : 'error', size: 'small' }, {
        default: () => row.isActive ? t('common.enable') : t('common.disable'),
      })
    },
  },
  {
    title: t('common.action'), key: 'actions', width: 150,
    render(row) {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { text: true, type: 'primary', size: 'small', onClick: () => handleEditTemplate(row) }, { default: () => t('common.edit') }),
          h(NPopconfirm, { onPositiveClick: () => handleDeleteTemplate(row) }, {
            trigger: () => h(NButton, { text: true, type: 'error', size: 'small' }, { default: () => t('common.delete') }),
            default: () => t('outputTemplate.deleteConfirm'),
          }),
        ],
      })
    },
  },
]

onMounted(async () => {
  await Promise.all([loadRobots(), loadTemplates()])
})
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('robot.title') }}</h1>
        <p class="page-subtitle text-muted">{{ t('robot.subtitle') }}</p>
      </div>
      <NSpace>
        <NButton v-if="activeTab === 'robots'" type="primary" @click="handleAddRobot">
          {{ t('robot.addRobot') }}
        </NButton>
        <NButton v-if="activeTab === 'templates'" type="primary" @click="handleAddTemplate">
          {{ t('outputTemplate.addTitle') }}
        </NButton>
      </NSpace>
    </div>

    <div class="glass-card mt-4">
      <NTabs v-model:value="activeTab" type="line" animated>
        <NTabPane name="robots" :tab="t('robot.robotList')">
          <div class="data-table-wrap">
            <NDataTable
              :columns="robotColumns"
              :data="robots"
              :loading="robotsLoading"
              :bordered="false"
              striped
            />
            <div v-if="!robots.length && !robotsLoading" class="empty-state">
              <NEmpty :description="t('common.noDataDesc')" />
            </div>
          </div>
        </NTabPane>
        <NTabPane name="templates" :tab="t('outputTemplate.title')">
          <div class="data-table-wrap">
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
      </NTabs>
    </div>

    <!-- Robot Modal -->
    <NModal
      v-model:show="robotDialogVisible"
      :title="robotDialogTitle"
      preset="card"
      style="width: 600px"
      :bordered="true"
    >
      <NForm :model="robotForm" label-placement="left" :label-width="100">
        <NFormItem :label="t('robot.name')" required>
          <NInput v-model:value="robotForm.name" :placeholder="t('robot.pleaseInputName')" />
        </NFormItem>
        <NFormItem :label="t('robot.type')">
          <NSelect v-model:value="robotForm.platform" :options="robotTypeOptions" style="width: 100%" />
        </NFormItem>
        <NFormItem :label="t('robot.webhookUrl')" required>
          <NInput v-model:value="robotForm.webhookUrl" :placeholder="t('robot.pleaseInputWebhookUrl')" />
        </NFormItem>
        <NFormItem :label="t('robot.secret')">
          <NInput v-model:value="robotForm.secret" type="password" show-password-on="click" :placeholder="t('robot.pleaseInputSecret')" />
        </NFormItem>
        <NFormItem :label="t('common.description')">
          <NInput v-model:value="robotForm.description" type="textarea" :rows="3" :placeholder="t('common.pleaseInputDescription')" />
        </NFormItem>
        <NFormItem :label="t('common.status')">
          <NSwitch v-model:value="robotForm.isActive" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="robotDialogVisible = false">{{ t('common.cancel') }}</NButton>
          <NButton type="primary" @click="handleSubmitRobot">{{ t('common.confirmButtonText') }}</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- Output Template Modal -->
    <NModal
      v-model:show="templateDialogVisible"
      :title="templateDialogTitle"
      preset="card"
      style="width: 700px"
      :bordered="true"
    >
      <NForm :model="templateForm" label-placement="left" :label-width="100">
        <NFormItem :label="t('outputTemplate.name')" required>
          <NInput v-model:value="templateForm.name" :placeholder="t('outputTemplate.pleaseInputName')" />
        </NFormItem>
        <NFormItem :label="t('outputTemplate.type')">
          <NSelect v-model:value="templateForm.platform" :options="robotTypeOptions" style="width: 100%" />
        </NFormItem>
        <NFormItem :label="t('outputTemplate.templateContent')">
          <NInput v-model:value="templateForm.content" type="textarea" :rows="10" :placeholder="t('outputTemplate.templateContentPlaceholder')" />
        </NFormItem>
        <NFormItem :label="t('common.description')">
          <NInput v-model:value="templateForm.description" type="textarea" :rows="2" :placeholder="t('common.pleaseInputDescription')" />
        </NFormItem>
        <NFormItem :label="t('common.status')">
          <NSwitch v-model:value="templateForm.isActive" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="templateDialogVisible = false">{{ t('common.cancel') }}</NButton>
          <NButton type="primary" @click="handleSubmitTemplate">{{ t('common.confirmButtonText') }}</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>
