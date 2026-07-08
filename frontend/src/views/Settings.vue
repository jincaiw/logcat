<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import { NForm, NFormItem, NInputNumber, NSelect, NSwitch, NButton, NSpin, NDivider, useMessage } from 'naive-ui'
import { API } from '@/api'
import type { SystemConfig } from '@/types'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const message = useMessage()
const loading = ref(false)
const saving = ref(false)

const settings = reactive<Partial<SystemConfig>>({
  listenPort: 5140,
  protocol: 'udp',
  logRetention: 30,
  maxLogSize: 100,
  autoStart: false,
  minimizeToTray: false,
  alertEnabled: true,
  alertInterval: 300,
  unmatchedLogRetention: 7,
  unmatchedLogAlert: false,
  defaultFilterAction: 'keep',
  theme: 'dark',
  language: 'zh-CN',
})

const protocolOptions = [
  { label: 'UDP', value: 'udp' },
  { label: 'TCP', value: 'tcp' },
]

const themeOptions = computed(() => [
  { label: t('settings.themeDark'), value: 'dark' },
  { label: t('settings.themeLight'), value: 'light' },
])

const languageOptions = computed(() => [
  { label: t('settings.langZhCN'), value: 'zh-CN' },
  { label: t('settings.langEnUS'), value: 'en-US' },
])

onMounted(async () => {
  await loadSettings()
})

async function loadSettings() {
  loading.value = true
  try {
    const data = await API.GetSystemConfig()
    Object.assign(settings, data)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    await API.SaveSystemConfig(settings)
    message.success(t('message.updateSuccess'))
  } catch (e) {
    message.error(t('message.operationFailed'))
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('settings.title') }}</h1>
        <p class="page-subtitle text-muted">{{ t('settings.description') }}</p>
      </div>
      <NSpace class="page-actions">
        <NButton type="primary" :loading="saving" @click="handleSave">
          {{ t('settings.save') }}
        </NButton>
      </NSpace>
    </div>

    <NSpin :show="loading">
      <div class="settings-layout mt-4">
        <!-- Basic Settings -->
        <div class="section-card settings-section">
          <div class="card-header">
            <h3 class="card-title text-accent">{{ t('settings.basicSettings') }}</h3>
            <p class="settings-desc text-muted">{{ t('settings.basicSettingsDesc') }}</p>
          </div>
          <NDivider class="settings-divider" />
          <NForm :model="settings" label-placement="left" :label-width="180">
            <NFormItem :label="t('settings.syslogListenPort')">
              <NInputNumber v-model:value="settings.listenPort" :min="1" :max="65535" class="field-control-200" />
            </NFormItem>
            <NFormItem :label="t('service.protocol')">
              <NSelect v-model:value="settings.protocol" :options="protocolOptions" class="field-control-200" />
            </NFormItem>
            <NFormItem :label="t('settings.logRetentionDays')">
              <NInputNumber v-model:value="settings.logRetention" :min="1" :max="365" class="field-control-200" />
              <span class="field-unit text-muted">{{ t('settings.days') }}</span>
            </NFormItem>
            <NFormItem :label="t('settings.maxLogSize')">
              <NInputNumber v-model:value="settings.maxLogSize" :min="1" :max="10000" :step="10" class="field-control-200" />
              <span class="field-unit text-muted">{{ t('common.mb') }}</span>
            </NFormItem>
            <NFormItem :label="t('settings.autoStart')">
              <NSwitch v-model:value="settings.autoStart" />
            </NFormItem>
            <NFormItem :label="t('settings.minimizeToTray')">
              <NSwitch v-model:value="settings.minimizeToTray" />
            </NFormItem>
          </NForm>
        </div>

        <!-- Notification Settings -->
        <div class="section-card settings-section">
          <div class="card-header">
            <h3 class="card-title text-accent">{{ t('settings.notificationSettings') }}</h3>
            <p class="settings-desc text-muted">{{ t('settings.notificationSettingsDesc') }}</p>
          </div>
          <NDivider class="settings-divider" />
          <NForm :model="settings" label-placement="left" :label-width="180">
            <NFormItem :label="t('settings.alertEnabled')">
              <NSwitch v-model:value="settings.alertEnabled" />
            </NFormItem>
            <NFormItem :label="t('settings.alertInterval')">
              <NInputNumber v-model:value="settings.alertInterval" :min="0" :max="3600" :step="10" class="field-control-200" />
              <span class="field-unit text-muted">{{ t('common.second') }}</span>
            </NFormItem>
            <NFormItem :label="t('settings.unmatchedLogAlert')">
              <NSwitch v-model:value="settings.unmatchedLogAlert" />
            </NFormItem>
            <NFormItem :label="t('settings.unmatchedLogRetention')">
              <NInputNumber v-model:value="settings.unmatchedLogRetention" :min="1" :max="365" class="field-control-200" />
              <span class="field-unit text-muted">{{ t('settings.days') }}</span>
            </NFormItem>
            <NFormItem :label="t('settings.defaultFilterAction')">
              <NSelect
                v-model:value="settings.defaultFilterAction"
                :options="[{ label: t('filterPolicy.actionKeep'), value: 'keep' }, { label: t('filterPolicy.actionDiscard'), value: 'discard' }]"
                class="field-control-200"
              />
            </NFormItem>
          </NForm>
        </div>

        <!-- Appearance Settings -->
        <div class="section-card settings-section">
          <div class="card-header">
            <h3 class="card-title text-accent">{{ t('settings.advancedSettings') }}</h3>
            <p class="settings-desc text-muted">{{ t('settings.advancedSettingsDesc') }}</p>
          </div>
          <NDivider class="settings-divider" />
          <NForm :model="settings" label-placement="left" :label-width="180">
            <NFormItem :label="t('settings.theme')">
              <NSelect v-model:value="settings.theme" :options="themeOptions" class="field-control-200" />
            </NFormItem>
            <NFormItem :label="t('settings.language')">
              <NSelect v-model:value="settings.language" :options="languageOptions" class="field-control-200" />
            </NFormItem>
          </NForm>
        </div>
      </div>
    </NSpin>
  </div>
</template>

