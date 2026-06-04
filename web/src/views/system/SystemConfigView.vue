<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NAlert, NButton, NCard, NForm, NFormItem, NInput, NInputNumber, NSelect, NSpin, NSpace, NSwitch, NTag } from 'naive-ui'
import { getSystemConfigs, updateSystemConfigs } from '@/api/system'
import type { SystemConfig } from '@/types'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAppMessage } from '@/composables/useMessage'
import { useIsMobile } from '@/composables/useIsMobile'
import { useTimeFormat } from '@/composables/useTimeFormat'

type FieldType = 'text' | 'number' | 'select' | 'switch' | 'password' | 'timezone' | 'timeformat'

interface ConfigFieldMeta {
  key: string
  label: string
  type: FieldType
  placeholder?: string
  description: string
  effect: string
  defaultValue: string
  options?: { label: string; value: string }[]
  condition?: (configs: Record<string, string>) => boolean
}

interface ConfigGroupMeta {
  key: string
  title: string
  description: string
  fields: ConfigFieldMeta[]
}

const message = useAppMessage()
const { isMobile } = useIsMobile()
const { setTimezone, setTimeFormat, commonTimezones, timeFormatOptions } = useTimeFormat()
const loading = ref(false)
const saving = ref(false)
const configs = ref<Record<string, string>>({})

const configGroups: ConfigGroupMeta[] = [
  {
    key: 'listener',
    title: '服务监听',
    description: '管理 Web/API 与 Syslog 接收端口。',
    fields: [
      { key: 'serverHost', label: '监听地址', type: 'text', defaultValue: '0.0.0.0', description: 'Web/API 服务监听地址。', effect: '重启生效' },
      { key: 'serverPort', label: 'Web 端口', type: 'number', defaultValue: '5080', description: '管理端和 API 默认访问端口。', effect: '重启生效' },
      {
        key: 'syslogEnabled',
        label: '启用 Syslog',
        type: 'select',
        defaultValue: 'false',
        description: '控制 Syslog 接收器默认启用状态。',
        effect: '重启生效',
        options: [
          { label: '启用', value: 'true' },
          { label: '禁用', value: 'false' },
        ],
      },
      { key: 'syslogUdpPort', label: 'Syslog UDP 端口', type: 'number', defaultValue: '5140', description: 'UDP 日志接收端口。', effect: '重启生效' },
      { key: 'syslogTcpPort', label: 'Syslog TCP 端口', type: 'number', defaultValue: '5140', description: 'TCP 日志接收端口。', effect: '重启生效' },
    ],
  },
  {
    key: 'security',
    title: '会话与登录安全',
    description: '控制登录失败锁定策略与会话存活时间。',
    fields: [
      { key: 'sessionExpireHours', label: '会话过期(小时)', type: 'number', defaultValue: '24', description: 'HttpOnly Session 默认过期时间。', effect: '重启生效' },
      { key: 'maxFailedLogin', label: '失败阈值', type: 'number', defaultValue: '5', description: '连续失败次数达到阈值后锁定账号。', effect: '重启生效' },
      { key: 'lockDurationMinutes', label: '锁定时长(分钟)', type: 'number', defaultValue: '30', description: '登录失败达到阈值后的锁定时长。', effect: '重启生效' },
    ],
  },
  {
    key: 'retention',
    title: '保留与清理',
    description: '控制日志保留周期和容量阈值。',
    fields: [
      { key: 'retentionDays', label: '日志保留天数', type: 'number', defaultValue: '90', description: '已匹配日志默认保留天数。', effect: '重启生效' },
      { key: 'unmatchedRetentionDays', label: '未匹配日志保留', type: 'number', defaultValue: '30', description: '未命中模板日志默认保留天数。', effect: '重启生效' },
      { key: 'maxLogSize', label: '最大日志大小(MB)', type: 'number', defaultValue: '10000', description: '日志文件或数据量的上限控制。', effect: '重启生效' },
      {
        key: 'defaultFilterAction',
        label: '默认过滤动作',
        type: 'select',
        defaultValue: 'keep',
        description: '日志未命中任何过滤规则时的默认处理动作。(FR-231)',
        effect: '重启生效',
        options: [
          { label: '保留 (keep)', value: 'keep' },
          { label: '丢弃 (drop)', value: 'drop' },
        ],
      },
    ],
  },
  {
    key: 'worker',
    title: 'Worker 与队列',
    description: '控制解析、筛选、推送并发和队列策略。',
    fields: [
      { key: 'parseWorkers', label: '解析 Worker', type: 'number', defaultValue: '4', description: '解析阶段 worker 数量。', effect: '重启生效' },
      { key: 'filterWorkers', label: '筛选 Worker', type: 'number', defaultValue: '4', description: '筛选阶段 worker 数量。', effect: '重启生效' },
      { key: 'pushWorkers', label: '推送 Worker', type: 'number', defaultValue: '4', description: '推送阶段 worker 数量。', effect: '重启生效' },
      { key: 'queueCapacity', label: '队列容量', type: 'number', defaultValue: '10000', description: '流水线队列默认容量。', effect: '重启生效' },
      {
        key: 'queueFullPolicy',
        label: '满队列策略',
        type: 'select',
        defaultValue: 'block_drop',
        description: '队列满时的处理策略。',
        effect: '重启生效',
        options: [
          { label: 'block_drop', value: 'block_drop' },
          { label: 'drop', value: 'drop' },
          { label: 'block', value: 'block' },
        ],
      },
    ],
  },
  {
    key: 'data',
    title: '数据与路径',
    description: '管理默认数据库类型、连接配置和存储路径。',
    fields: [
      {
        key: 'databaseType',
        label: '数据库类型',
        type: 'select',
        defaultValue: 'sqlite',
        description: '当前支持 sqlite 和 mysql。',
        effect: '重启生效',
        options: [
          { label: 'SQLite', value: 'sqlite' },
          { label: 'MySQL', value: 'mysql' },
        ],
      },
      { key: 'sqlitePath', label: 'SQLite 路径', type: 'text', defaultValue: 'data/logcat.db', description: 'SQLite 数据文件路径。', effect: '重启生效', condition: (c) => c.databaseType !== 'mysql' },
      { key: 'configDir', label: '配置目录', type: 'text', defaultValue: 'data/config', description: '系统配置文件存储目录路径。', effect: '重启生效' },
      { key: 'logDir', label: '日志目录', type: 'text', defaultValue: 'data/logs', description: '系统运行日志存储目录路径。', effect: '重启生效' },
      { key: 'mysqlHost', label: 'MySQL 主机', type: 'text', defaultValue: '127.0.0.1', description: 'MySQL 服务连接地址。', effect: '重启生效', condition: (c) => c.databaseType === 'mysql' },
      { key: 'mysqlPort', label: 'MySQL 端口', type: 'number', defaultValue: '3306', description: 'MySQL 服务连接端口。', effect: '重启生效', condition: (c) => c.databaseType === 'mysql' },
      { key: 'mysqlDatabase', label: 'MySQL 数据库', type: 'text', defaultValue: 'logcat', description: 'MySQL 数据库名称。', effect: '重启生效', condition: (c) => c.databaseType === 'mysql' },
      { key: 'mysqlUsername', label: 'MySQL 用户名', type: 'text', defaultValue: 'root', description: 'MySQL 连接用户名。', effect: '重启生效', condition: (c) => c.databaseType === 'mysql' },
      { key: 'mysqlPassword', label: 'MySQL 密码', type: 'password', defaultValue: '', description: 'MySQL 连接密码。', effect: '重启生效', condition: (c) => c.databaseType === 'mysql' },
      { key: 'mysqlCharset', label: 'MySQL 字符集', type: 'text', defaultValue: 'utf8mb4', description: 'MySQL 连接字符集。', effect: '重启生效', condition: (c) => c.databaseType === 'mysql' },
      { key: 'mysqlTimezone', label: 'MySQL 时区', type: 'text', defaultValue: 'Asia/Shanghai', description: 'MySQL 连接时区。', effect: '重启生效', condition: (c) => c.databaseType === 'mysql' },
      { key: 'mysqlMaxOpenConns', label: '最大打开连接数', type: 'number', defaultValue: '50', description: 'MySQL 最大打开连接数。', effect: '重启生效', condition: (c) => c.databaseType === 'mysql' },
      { key: 'mysqlMaxIdleConns', label: '最大空闲连接数', type: 'number', defaultValue: '10', description: 'MySQL 最大空闲连接数。', effect: '重启生效', condition: (c) => c.databaseType === 'mysql' },
    ],
  },
  {
    key: 'alert',
    title: '告警配置',
    description: '控制告警开关与告警频率。',
    fields: [
      { key: 'alertEnabled', label: '告警开关', type: 'switch', defaultValue: 'true', description: '全局告警启用/禁用。(FR-227)', effect: '即时生效' },
      { key: 'alertInterval', label: '告警间隔(秒)', type: 'number', defaultValue: '60', description: '同一告警规则的最小触发间隔，单位秒。(FR-228)', effect: '即时生效' },
      { key: 'unmatchedAlertEnabled', label: '未匹配日志告警', type: 'switch', defaultValue: 'false', description: '当日志未命中任何模板时是否触发告警。(FR-230)', effect: '重启生效' },
    ],
  },
  {
    key: 'theme',
    title: '界面配置',
    description: '控制前端界面主题与显示偏好。',
    fields: [
      {
        key: 'theme',
        label: '主题',
        type: 'select',
        defaultValue: 'light',
        description: '界面主题切换。(FR-232)',
        effect: '即时生效',
        options: [
          { label: '浅色 (light)', value: 'light' },
          { label: '深色 (dark)', value: 'dark' },
        ],
      },
      {
        key: 'timezone',
        label: '时区',
        type: 'timezone',
        defaultValue: 'Asia/Shanghai',
        description: '系统全局时间显示时区，影响所有页面的时间展示。',
        effect: '即时生效',
      },
      {
        key: 'timeFormat',
        label: '时间格式',
        type: 'timeformat',
        defaultValue: 'YYYY-MM-DD HH:mm:ss',
        description: '系统全局时间显示格式，影响所有页面的时间展示。',
        effect: '即时生效',
      },
    ],
  },
]

function applyDefaults() {
  for (const group of configGroups) {
    for (const field of group.fields) {
      if (configs.value[field.key] === undefined || configs.value[field.key] === '') {
        configs.value[field.key] = field.defaultValue
      }
    }
  }
}

function updateConfigValue(key: string, value: string | number | null) {
  configs.value[key] = value === null || value === undefined ? '' : String(value)
}

function numberValue(key: string) {
  const raw = configs.value[key]
  if (raw === undefined || raw === '') {
    return null
  }
  const value = Number(raw)
  return Number.isNaN(value) ? null : value
}

function switchValue(key: string): boolean {
  return configs.value[key] === 'true'
}

function updateSwitchValue(key: string, value: boolean) {
  configs.value[key] = String(value)
}

function isFieldVisible(field: ConfigFieldMeta): boolean {
  if (!field.condition) return true
  return field.condition(configs.value)
}

async function loadConfigs() {
  loading.value = true
  try {
    const res = await getSystemConfigs()
    const map: Record<string, string> = {}
    if (res.data) {
      for (const cfg of res.data as SystemConfig[]) {
        map[cfg.configKey] = cfg.configValue
      }
    }
    configs.value = map
    applyDefaults()
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
    if (configs.value.timezone) setTimezone(configs.value.timezone)
    if (configs.value.timeFormat) setTimeFormat(configs.value.timeFormat)
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
      <n-space vertical :size="16">
        <n-alert type="warning" :show-icon="false">
          当前页面保存的是系统配置记录，便于统一管理和审计。运行中的服务仍以配置文件和环境变量为准。标记为「重启生效」的参数需重启服务后生效，标记为「即时生效」的参数保存后立即生效。
        </n-alert>

        <n-card v-for="group in configGroups" :key="group.key" size="small" :title="group.title">
          <template #header-extra>
            <span>{{ group.description }}</span>
          </template>
          <n-form :label-placement="isMobile ? 'top' : 'left'" :label-width="isMobile ? undefined : 180" require-mark-placement="left">
            <n-form-item
              v-for="field in group.fields"
              v-show="isFieldVisible(field)"
              :key="field.key"
              :label="field.label"
            >
              <div class="config-field">
                <n-switch
                  v-if="field.type === 'switch'"
                  :value="switchValue(field.key)"
                  @update:value="(value) => updateSwitchValue(field.key, value)"
                />
                <n-input
                  v-else-if="field.type === 'password'"
                  type="password"
                  show-password-on="click"
                  :value="configs[field.key]"
                  :placeholder="field.placeholder || field.label"
                  @update:value="(value) => updateConfigValue(field.key, value)"
                />
                <n-input
                  v-else-if="field.type === 'text'"
                  :value="configs[field.key]"
                  :placeholder="field.placeholder || field.label"
                  @update:value="(value) => updateConfigValue(field.key, value)"
                />
                <n-input-number
                  v-else-if="field.type === 'number'"
                  :value="numberValue(field.key)"
                  :min="0"
                  style="max-width: 240px; width: 100%"
                  @update:value="(value) => updateConfigValue(field.key, value)"
                />
                <n-select
                  v-else-if="field.type === 'timezone'"
                  :value="configs[field.key]"
                  :options="commonTimezones"
                  filterable
                  style="max-width: 360px; width: 100%"
                  @update:value="(value) => updateConfigValue(field.key, value)"
                />
                <n-select
                  v-else-if="field.type === 'timeformat'"
                  :value="configs[field.key]"
                  :options="timeFormatOptions"
                  style="max-width: 360px; width: 100%"
                  @update:value="(value) => updateConfigValue(field.key, value)"
                />
                <n-select
                  v-else
                  :value="configs[field.key]"
                  :options="field.options || []"
                  style="max-width: 240px; width: 100%"
                  @update:value="(value) => updateConfigValue(field.key, value)"
                />
                <div class="field-meta">
                  <div class="field-description">{{ field.description }}</div>
                  <n-tag size="small" :type="field.effect === '即时生效' ? 'success' : 'warning'" :bordered="false">{{ field.effect }}</n-tag>
                </div>
              </div>
            </n-form-item>
          </n-form>
        </n-card>
      </n-space>
    </n-spin>
  </div>
</template>

<style scoped>
.config-field {
  width: 100%;
}

.field-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.field-description {
  color: var(--n-text-color-3);
  font-size: 12px;
}
</style>
