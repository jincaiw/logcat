<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NForm, NFormItem, NInput, NButton, NSpace, NSpin, NCard, useMessage } from 'naive-ui'
import { getSystemConfigs, updateSystemConfigs } from '@/api/system'
import type { SystemConfig } from '@/types'
import PageHeader from '@/components/common/PageHeader.vue'

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const configs = ref<Record<string, string>>({})

const configFields = [
  { key: 'logRetentionDays', label: '日志保留天数' },
  { key: 'maxLogSize', label: '最大日志大小 (MB)' },
  { key: 'pushRetryCount', label: '推送重试次数' },
  { key: 'pushRetryInterval', label: '推送重试间隔 (秒)' },
  { key: 'queueCapacity', label: '队列容量' },
  { key: 'workerCount', label: '工作器数量' },
  { key: 'alertCooldown', label: '告警冷却时间 (秒)' },
  { key: 'maxConnections', label: '最大连接数' },
  { key: 'connectionTimeout', label: '连接超时 (秒)' },
  { key: 'syslogPort', label: 'Syslog 端口' },
  { key: 'tlsCert', label: 'TLS 证书路径' },
  { key: 'tlsKey', label: 'TLS 密钥路径' },
  { key: 'logLevel', label: '日志级别' },
]

async function loadConfigs() {
  loading.value = true
  try {
    const res = await getSystemConfigs()
    const map: Record<string, string> = {}
    if (res.data) {
      for (const cfg of res.data) {
        map[cfg.configKey] = cfg.configValue
      }
    }
    configs.value = map
  } catch {
    message.error('加载配置失败')
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    await updateSystemConfigs(configs.value)
    message.success('配置保存成功')
  } catch {
    message.error('配置保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadConfigs()
})
</script>

<template>
  <div class="page-container">
    <PageHeader title="系统配置" description="系统全局配置管理">
      <n-button type="primary" :loading="saving" @click="handleSave">
        保存配置
      </n-button>
    </PageHeader>

    <n-spin :show="loading">
      <n-card size="small">
        <n-form label-placement="left" label-width="160" require-mark-placement="left">
          <n-form-item
            v-for="field in configFields"
            :key="field.key"
            :label="field.label"
          >
            <n-input
              v-model:value="configs[field.key]"
              :placeholder="field.label"
              style="max-width: 400px"
            />
          </n-form-item>
        </n-form>
      </n-card>
    </n-spin>
  </div>
</template>