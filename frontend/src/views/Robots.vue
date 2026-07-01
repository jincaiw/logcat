<script setup lang="ts">
import { ref, computed, onMounted, h, watch } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, NSwitch, NTag, NPopconfirm, NEmpty, NTabs, NTabPane, NSpace, NRadioGroup, NRadio, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { API } from '@/api'
import type { Robot, OutputTemplate } from '@/types'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const message = useMessage()
const activeTab = ref('robots')

// ==================== Platform Options ====================
const platformOptions = computed(() => [
  { label: t('robot.typeFeishu'), value: 'feishu' },
  { label: t('robot.typeEmail'), value: 'email' },
  { label: t('robot.typeSyslog'), value: 'syslog' },
  { label: t('robot.typeHttp'), value: 'http' },
])

const syslogFormatOptions = [
  { label: 'JSON', value: 'json' },
  { label: 'RFC 3164', value: 'rfc3164' },
  { label: 'RFC 5424', value: 'rfc5424' },
]

function getPlatformLabel(platform: string): string {
  const opt = platformOptions.value.find(o => o.value === platform)
  return opt ? opt.label : platform
}

function getPlatformTagType(platform: string): 'default' | 'primary' | 'success' | 'warning' | 'info' | 'error' {
  const map: Record<string, 'default' | 'primary' | 'success' | 'warning' | 'info' | 'error'> = {
    feishu: 'success',
    email: 'info',
    syslog: 'default',
    http: 'warning',
  }
  return map[platform] || 'default'
}

function isSupportedPlatform(platform?: string): boolean {
  return !!platform && platformOptions.value.some(o => o.value === platform)
}

// ==================== Robots ====================
const robots = ref<Robot[]>([])
const robotsLoading = ref(false)
const robotDialogVisible = ref(false)
const robotDialogTitle = ref(t('robot.addRobot'))
const robotForm = ref<Partial<Robot>>({
  name: '',
  platform: 'feishu',
  feishuWebhookUrl: '',
  feishuSecret: '',
  smtpHost: '',
  smtpPort: 25,
  smtpUsername: '',
  smtpPassword: '',
  smtpFrom: '',
  smtpTo: '',
  syslogHost: '',
  syslogPort: 514,
  syslogProtocol: 'udp',
  syslogFormat: 'json',
  httpUrl: 'http://168.63.6.81:8080/cib-message/public/service/sendwbg.do',
  httpTimeout: 3,
  httpRetryCount: 3,
  httpRetryDelay: 2,
  httpNotesIds: '420102,420809',
  isActive: true,
  description: '',
})

watch(() => robotForm.value.platform, (newPlatform) => {
  if (newPlatform === 'syslog' && !robotForm.value.syslogPort) {
    robotForm.value.syslogPort = 514
  }
  if (newPlatform === 'syslog' && !robotForm.value.syslogProtocol) {
    robotForm.value.syslogProtocol = 'udp'
  }
  if (newPlatform === 'syslog' && !robotForm.value.syslogFormat) {
    robotForm.value.syslogFormat = 'json'
  }
  if (newPlatform === 'email' && !robotForm.value.smtpPort) {
    robotForm.value.smtpPort = 25
  }
  if (newPlatform === 'http') {
    if (!robotForm.value.httpUrl) robotForm.value.httpUrl = 'http://168.63.6.81:8080/cib-message/public/service/sendwbg.do'
    if (!robotForm.value.httpTimeout) robotForm.value.httpTimeout = 3
    if (robotForm.value.httpRetryCount === undefined) robotForm.value.httpRetryCount = 3
    if (robotForm.value.httpRetryDelay === undefined) robotForm.value.httpRetryDelay = 2
    if (!robotForm.value.httpNotesIds) robotForm.value.httpNotesIds = '420102,420809'
  }
})

const isFeishu = computed(() => robotForm.value.platform === 'feishu')
const isEmail = computed(() => robotForm.value.platform === 'email')
const isSyslog = computed(() => robotForm.value.platform === 'syslog')
const isHttp = computed(() => robotForm.value.platform === 'http')

async function loadRobots() {
  robotsLoading.value = true
  try {
    robots.value = await API.GetRobots()
  } catch (e) { console.error(e) }
  finally { robotsLoading.value = false }
}

function createDefaultRobotForm(): Partial<Robot> {
  return {
    name: '',
    platform: 'feishu',
    feishuWebhookUrl: '',
    feishuSecret: '',
    smtpHost: '',
    smtpPort: 25,
    smtpUsername: '',
    smtpPassword: '',
    smtpFrom: '',
    smtpTo: '',
    syslogHost: '',
    syslogPort: 514,
    syslogProtocol: 'udp',
    syslogFormat: 'json',
    httpUrl: 'http://168.63.6.81:8080/cib-message/public/service/sendwbg.do',
    httpTimeout: 3,
    httpRetryCount: 3,
    httpRetryDelay: 2,
    httpNotesIds: '420102,420809',
    isActive: true,
    description: '',
  }
}

function handleAddRobot() {
  robotDialogTitle.value = t('robot.addRobot')
  robotForm.value = createDefaultRobotForm()
  robotDialogVisible.value = true
}

function handleEditRobot(row: Robot) {
  robotDialogTitle.value = t('robot.editRobot')
  robotForm.value = { ...createDefaultRobotForm(), ...row }
  robotDialogVisible.value = true
}

async function handleDeleteRobot(row: Robot) {
  try {
    await API.DeleteRobot(row.id!)
    message.success(t('message.deleteSuccess'))
    loadRobots()
  } catch (e) { message.error(t('message.deleteFailed')) }
}

function validateRobotForm(): boolean {
  if (!robotForm.value.name) {
    message.warning(t('common.requiredFields'))
    return false
  }
  switch (robotForm.value.platform) {
    case 'feishu':
      if (!robotForm.value.feishuWebhookUrl) {
        message.warning(t('robot.webhookRequired'))
        return false
      }
      break
    case 'email':
      if (!robotForm.value.smtpHost || !robotForm.value.smtpFrom || !robotForm.value.smtpTo) {
        message.warning(t('robot.emailRequired'))
        return false
      }
      break
    case 'syslog':
      if (!robotForm.value.syslogHost || !robotForm.value.syslogPort) {
        message.warning(t('robot.syslogRequired'))
        return false
      }
      break
    case 'http':
      if (!robotForm.value.httpUrl) {
        message.warning(t('robot.httpRequired'))
        return false
      }
      break
    default:
      message.warning(t('common.requiredFields'))
      return false
  }
  return true
}

async function handleSubmitRobot() {
  if (!validateRobotForm()) return
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
  } catch (e: any) { message.error(e?.message || t('message.operationFailed')) }
}

async function handleTestRobot(row: Robot) {
  try {
    const result = await API.TestRobot(row)
    message.success(result.message || t('robot.testSuccess'))
  } catch (e: any) {
    message.error(t('robot.testFailed') + (e?.message ? `: ${e.message}` : ''))
  }
}

async function handleTestRobotForm() {
  if (!validateRobotForm()) return
  try {
    const result = await API.TestRobot(robotForm.value as any)
    message.success(result.message || t('robot.testSuccess'))
  } catch (e: any) {
    message.error(t('robot.testFailed') + (e?.message ? `: ${e.message}` : ''))
  }
}

const robotColumns: DataTableColumns<Robot> = [
  { title: t('common.id'), key: 'id', width: 70 },
  { title: t('robot.name'), key: 'name', width: 160, ellipsis: { tooltip: true } },
  {
    title: t('robot.type'), key: 'platform', width: 110,
    render(row) {
      return h(NTag, {
        type: getPlatformTagType(row.platform),
        size: 'small',
      }, { default: () => getPlatformLabel(row.platform) })
    },
  },
  {
    title: t('robot.configInfo'), key: 'configInfo', ellipsis: { tooltip: true },
    render(row) {
      let info = '-'
      switch (row.platform) {
        case 'feishu': info = row.feishuWebhookUrl || '-'; break
        case 'email': info = row.smtpHost ? `${row.smtpHost}:${row.smtpPort || 25}` : '-'; break
        case 'syslog': info = row.syslogHost ? `${row.syslogHost}:${row.syslogPort || 514} (${(row.syslogProtocol || 'udp').toUpperCase()}/${row.syslogFormat || 'json'})` : '-'; break
        case 'http': info = row.httpUrl || '-'; break
        default: info = '-'
      }
      return h('span', { class: 'mono text-muted', style: { fontSize: '13px' } }, info)
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
          h(NButton, { text: true, type: 'info', size: 'small', onClick: () => handleTestRobot(row), disabled: !row.isActive }, { default: () => t('robot.testSend') }),
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
  fields: '',
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
  templateForm.value = { name: '', platform: 'feishu', content: '', fields: '', description: '', isActive: true }
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
  if (!isSupportedPlatform(templateForm.value.platform)) {
    message.warning(t('common.requiredFields'))
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
  } catch (e: any) { message.error(e?.message || t('message.operationFailed')) }
}

const templateColumns: DataTableColumns<OutputTemplate> = [
  { title: t('common.id'), key: 'id', width: 70 },
  { title: t('outputTemplate.name'), key: 'name', width: 160, ellipsis: { tooltip: true } },
  {
    title: t('outputTemplate.type'), key: 'platform', width: 110,
    render(row) {
      return h(NTag, { type: getPlatformTagType(row.platform), size: 'small' }, {
        default: () => getPlatformLabel(row.platform),
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
      style="width: 640px"
      :bordered="true"
    >
      <NForm :model="robotForm" label-placement="left" :label-width="120">
        <NFormItem :label="t('robot.name')" required>
          <NInput v-model:value="robotForm.name" :placeholder="t('robot.pleaseInputName')" />
        </NFormItem>
        <NFormItem :label="t('robot.type')">
          <NSelect v-model:value="robotForm.platform" :options="platformOptions" style="width: 100%" />
        </NFormItem>

        <!-- Feishu -->
        <template v-if="isFeishu">
          <NFormItem :label="t('robot.webhookUrl')" required>
            <NInput v-model:value="robotForm.feishuWebhookUrl" :placeholder="t('robot.feishuWebhookPlaceholder')" />
          </NFormItem>
          <NFormItem :label="t('robot.secret')">
            <NInput v-model:value="robotForm.feishuSecret" type="password" show-password-on="click" :placeholder="t('robot.feishuSecretPlaceholder')" />
          </NFormItem>
        </template>

        <!-- Email -->
        <template v-if="isEmail">
          <NFormItem :label="t('robot.smtpServer')" required>
            <NSpace align="center" :wrap="false" style="width: 100%">
              <NInput v-model:value="robotForm.smtpHost" :placeholder="t('robot.smtpHostPlaceholder')" style="flex: 2" />
              <NInputNumber v-model:value="robotForm.smtpPort" :min="1" :max="65535" placeholder="Port" style="width: 100px; flex-shrink: 0" />
            </NSpace>
          </NFormItem>
          <NFormItem :label="t('robot.smtpUsername')">
            <NInput v-model:value="robotForm.smtpUsername" :placeholder="t('robot.smtpUsernamePlaceholder')" />
          </NFormItem>
          <NFormItem :label="t('robot.smtpPassword')">
            <NInput v-model:value="robotForm.smtpPassword" type="password" show-password-on="click" :placeholder="t('robot.smtpPasswordPlaceholder')" />
          </NFormItem>
          <NFormItem :label="t('robot.smtpFrom')" required>
            <NInput v-model:value="robotForm.smtpFrom" :placeholder="t('robot.smtpFromPlaceholder')" />
          </NFormItem>
          <NFormItem :label="t('robot.smtpTo')" required>
            <NInput v-model:value="robotForm.smtpTo" :placeholder="t('robot.smtpToPlaceholder')" />
          </NFormItem>
        </template>

        <!-- Syslog Forward -->
        <template v-if="isSyslog">
          <NFormItem :label="t('robot.syslogTarget')" required>
            <NSpace align="center" :wrap="false" style="width: 100%">
              <NInput v-model:value="robotForm.syslogHost" :placeholder="t('robot.syslogHostPlaceholder')" style="flex: 2" />
              <NInputNumber v-model:value="robotForm.syslogPort" :min="1" :max="65535" placeholder="Port" style="width: 100px; flex-shrink: 0" />
            </NSpace>
          </NFormItem>
          <NFormItem :label="t('robot.syslogProtocol')">
            <NRadioGroup v-model:value="robotForm.syslogProtocol">
              <NRadio value="udp">UDP</NRadio>
              <NRadio value="tcp">TCP</NRadio>
            </NRadioGroup>
          </NFormItem>
          <NFormItem :label="t('robot.syslogFormat')">
            <NSelect v-model:value="robotForm.syslogFormat" :options="syslogFormatOptions" style="width: 200px" />
          </NFormItem>
        </template>

        <!-- HTTP Interface -->
        <template v-if="isHttp">
          <NFormItem :label="t('robot.httpUrl')" required>
            <NInput v-model:value="robotForm.httpUrl" :placeholder="t('robot.httpUrlPlaceholder')" />
          </NFormItem>
          <NFormItem :label="t('robot.httpNotesIds')">
            <NInput v-model:value="robotForm.httpNotesIds" :placeholder="t('robot.httpNotesIdsPlaceholder')" />
          </NFormItem>
          <NFormItem :label="t('robot.httpRetrySettings')">
            <NSpace align="center" :wrap="false" style="width: 100%">
              <NInputNumber v-model:value="robotForm.httpTimeout" :min="1" :max="60" :placeholder="t('robot.httpTimeout')" style="width: 150px" />
              <NInputNumber v-model:value="robotForm.httpRetryCount" :min="0" :max="10" :placeholder="t('robot.httpRetryCount')" style="width: 150px" />
              <NInputNumber v-model:value="robotForm.httpRetryDelay" :min="0" :max="60" :placeholder="t('robot.httpRetryDelay')" style="width: 150px" />
            </NSpace>
          </NFormItem>
        </template>

        <NFormItem :label="t('common.description')">
          <NInput v-model:value="robotForm.description" type="textarea" :rows="2" :placeholder="t('common.pleaseInputDescription')" />
        </NFormItem>
        <NFormItem :label="t('common.status')">
          <NSwitch v-model:value="robotForm.isActive" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="robotDialogVisible = false">{{ t('common.cancel') }}</NButton>
          <NButton @click="handleTestRobotForm">{{ t('robot.testSend') }}</NButton>
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
          <NSelect v-model:value="templateForm.platform" :options="platformOptions" style="width: 100%" />
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
