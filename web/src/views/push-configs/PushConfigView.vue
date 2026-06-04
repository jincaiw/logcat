<script setup lang="ts">
import { computed, h, ref } from 'vue'
import {
  NAlert,
  NButton,
  NCode,
  NDescriptions,
  NDescriptionsItem,
  NDynamicTags,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NSpace,
  NSwitch,
  NTabPane,
  NTabs,
  NTag,
  NSelect,
} from 'naive-ui'
import type { DataTableColumns, FormInst } from 'naive-ui'
import { createPushConfig, updatePushConfig, deletePushConfig, getPushConfigs, testPushConfig } from '@/api/pushConfigs'
import type { PushConfig, PushTestResult } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAppMessage } from '@/composables/useMessage'
import { useIsMobile } from '@/composables/useIsMobile'

const message = useAppMessage()
const { isMobile } = useIsMobile()
const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const formRef = ref<FormInst | null>(null)
const formShow = ref(false)
const formLoading = ref(false)
const confirmDialogShow = ref(false)
const confirmTitle = ref('')
const confirmContent = ref('')
const confirmAction = ref<() => Promise<void>>(() => Promise.resolve())
const confirmLoading = ref(false)
const editingRow = ref<PushConfig | null>(null)
const testResultShow = ref(false)
const testResult = ref<PushTestResult | null>(null)
const testingId = ref<number | null>(null)

// 动态标签：notesIds（接收标识列表）和 retryOnStatusCodes（需重试状态码）
const notesIdsTags = ref<string[]>([])
const retryOnStatusCodesTags = ref<string[]>([])

/** 将 JSON 数组字符串解析为 string[] */
function parseJsonArray(val: string | undefined | null): string[] {
  if (!val) return []
  try {
    const arr = JSON.parse(val)
    return Array.isArray(arr) ? arr.map(String) : []
  } catch {
    return []
  }
}

/** 将 string[] 编码为 JSON 数组字符串 */
function stringifyArray(arr: string[]): string {
  return JSON.stringify(arr)
}

/** 掩码 URL 中的敏感查询参数 */
function maskUrl(url: string | undefined): string {
  if (!url) return '-'
  try {
    const u = new URL(url)
    const sensitiveKeys = ['token', 'key', 'secret', 'password', 'access_token', 'api_key', 'apikey', 'auth']
    u.searchParams.forEach((_, key) => {
      if (sensitiveKeys.some(sk => key.toLowerCase().includes(sk))) {
        u.searchParams.set(key, '****')
      }
    })
    return u.toString()
  } catch {
    return url
  }
}

const typeOptions = [
  { label: 'HTTP', value: 'http' },
  { label: 'Email', value: 'email' },
  { label: 'Syslog', value: 'syslog' },
]

const authTypeOptions = [
  { label: '无认证', value: 'none' },
  { label: 'Bearer', value: 'bearer' },
  { label: 'Basic', value: 'basic' },
  { label: '自定义 Header', value: 'custom_header' },
]

const syslogProtocolOptions = [
  { label: 'UDP', value: 'udp' },
  { label: 'TCP', value: 'tcp' },
]

const syslogFormatOptions = [
  { label: 'RFC3164', value: 'rfc3164' },
  { label: 'RFC5424', value: 'rfc5424' },
  { label: 'JSON', value: 'json' },
]

function createDefaultForm(): Partial<PushConfig> {
  return {
    name: '',
    type: 'http',
    enabled: true,
    method: 'POST',
    timeout: 30,
    retryCount: 0,
    retryDelay: 5,
    headers: '{\n  "X-Source": "logcat"\n}',
    bodyTemplate: '{\n  "message": "${message}",\n  "sourceIp": "${source_ip}",\n  "timestamp": "${timestamp}"\n}',
    successStatusCodes: '[200,201,202]',
    authType: 'none',
    token: '',
    contentType: 'application/json',
    maxResponseLogSize: 4096,
    notesIds: '[]',
    retryOnStatusCodes: '[]',
    smtpPort: 587,
    smtpUsername: '',
    smtpPassword: '',
    fromAddress: '',
    toAddresses: '["security@example.com"]',
    subjectTemplate: 'logcat 告警通知 ${event_type}',
    emailBodyTemplate: '来源 IP: ${source_ip}\n消息内容: ${message}\n时间: ${timestamp}',
    syslogPort: 514,
    syslogProtocol: 'udp',
    syslogFormat: 'rfc3164',
    syslogFields: '{\n  "message": "${message}",\n  "source_ip": "${source_ip}",\n  "event_type": "${event_type}"\n}',
  }
}

const formData = ref<Partial<PushConfig>>(createDefaultForm())
const currentType = computed(() => formData.value.type || 'http')

const columns: DataTableColumns<PushConfig> = [
  { title: '名称', key: 'name' },
  {
    title: '类型',
    key: 'type',
    render(row) {
      return h(NTag, { size: 'small', bordered: false }, { default: () => row.type?.toUpperCase() })
    },
  },
  {
    title: '目标',
    key: 'target',
    render(row) {
      if (row.type === 'http') {
        return maskUrl(row.url)
      }
      if (row.type === 'email') {
        return row.smtpHost || '-'
      }
      return row.syslogHost || '-'
    },
  },
  {
    title: 'Token',
    key: 'token',
    render(row) {
      if (row.type !== 'http' || !row.token) return '-'
      return '****'
    },
  },
  {
    title: 'SMTP密码',
    key: 'smtpPassword',
    render(row) {
      if (row.type !== 'email' || !row.smtpPassword) return '-'
      return '****'
    },
  },
  {
    title: '状态',
    key: 'enabled',
    render(row) {
      return h(NSwitch, {
        size: 'small',
        value: !!row.enabled,
        onUpdateValue: (val: boolean) => handleToggleEnabled(row, val),
      })
    },
  },
  {
    title: '操作',
    key: 'actions',
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
          h(
            NButton,
            { size: 'small', loading: testingId.value === row.id, onClick: () => handleTest(row) },
            { default: () => '测试' }
          ),
          h(NButton, { size: 'small', type: 'error', ghost: true, onClick: () => handleDelete(row) }, { default: () => '删除' }),
        ],
      })
    },
  },
]

async function fetchData(params: any) {
  return getPushConfigs(params)
}

function openForm(data?: Partial<PushConfig>) {
  formData.value = {
    ...createDefaultForm(),
    ...(data || {}),
  }
  // 同步动态标签状态
  notesIdsTags.value = parseJsonArray(formData.value.notesIds)
  retryOnStatusCodesTags.value = parseJsonArray(formData.value.retryOnStatusCodes)
  formShow.value = true
}

function handleAdd() {
  editingRow.value = null
  openForm()
}

function handleEdit(row: PushConfig) {
  editingRow.value = row
  openForm({
    ...row,
    headers: row.headers || createDefaultForm().headers,
    successStatusCodes: row.successStatusCodes || createDefaultForm().successStatusCodes,
    toAddresses: row.toAddresses || createDefaultForm().toAddresses,
    syslogFields: row.syslogFields || createDefaultForm().syslogFields,
    notesIds: row.notesIds || createDefaultForm().notesIds,
    retryOnStatusCodes: row.retryOnStatusCodes || createDefaultForm().retryOnStatusCodes,
  })
}

function closeForm() {
  formShow.value = false
}

function validateForm() {
  if (!formData.value.name?.trim()) {
    message.warning('请输入配置名称')
    return false
  }
  if (!formData.value.type) {
    message.warning('请选择推送类型')
    return false
  }

  if (currentType.value === 'http') {
    if (!formData.value.url?.trim()) {
      message.warning('请输入 HTTP 地址')
      return false
    }
  }

  if (currentType.value === 'email') {
    if (!formData.value.smtpHost?.trim()) {
      message.warning('请输入 SMTP 主机')
      return false
    }
    if (!formData.value.fromAddress?.trim()) {
      message.warning('请输入发件人地址')
      return false
    }
    if (!formData.value.toAddresses?.trim()) {
      message.warning('请输入收件人地址')
      return false
    }
  }

  if (currentType.value === 'syslog') {
    if (!formData.value.syslogHost?.trim()) {
      message.warning('请输入 Syslog 主机')
      return false
    }
    if (!formData.value.syslogPort) {
      message.warning('请输入 Syslog 端口')
      return false
    }
  }

  return true
}

function buildPayload(): Partial<PushConfig> {
  const payload: Partial<PushConfig> = {
    name: formData.value.name?.trim(),
    type: formData.value.type,
    enabled: !!formData.value.enabled,
    url: formData.value.url?.trim(),
    method: formData.value.method?.trim() || 'POST',
    timeout: Number(formData.value.timeout || 30),
    retryCount: Number(formData.value.retryCount || 0),
    retryDelay: Number(formData.value.retryDelay || 0),
    headers: formData.value.headers?.trim(),
    bodyTemplate: formData.value.bodyTemplate || '',
    successStatusCodes: formData.value.successStatusCodes?.trim(),
    successBodyKeyword: formData.value.successBodyKeyword?.trim(),
    authType: formData.value.authType || 'none',
    contentType: formData.value.contentType?.trim() || 'application/json',
    maxResponseLogSize: Number(formData.value.maxResponseLogSize || 4096),
    notesIds: stringifyArray(notesIdsTags.value),
    retryOnStatusCodes: stringifyArray(retryOnStatusCodesTags.value),
    smtpHost: formData.value.smtpHost?.trim(),
    smtpPort: Number(formData.value.smtpPort || 587),
    smtpUsername: formData.value.smtpUsername?.trim(),
    fromAddress: formData.value.fromAddress?.trim(),
    toAddresses: formData.value.toAddresses?.trim(),
    subjectTemplate: formData.value.subjectTemplate || '',
    emailBodyTemplate: formData.value.emailBodyTemplate || '',
    syslogHost: formData.value.syslogHost?.trim(),
    syslogPort: Number(formData.value.syslogPort || 514),
    syslogProtocol: formData.value.syslogProtocol || 'udp',
    syslogFormat: formData.value.syslogFormat || 'rfc3164',
    syslogFields: formData.value.syslogFields?.trim(),
  }
  if (formData.value.token) {
    payload.token = formData.value.token
  }
  if (formData.value.smtpPassword) {
    payload.smtpPassword = formData.value.smtpPassword
  }
  return payload
}

async function handleSubmit() {
  await formRef.value?.validate().catch(() => undefined)
  if (!validateForm()) {
    return
  }

  formLoading.value = true
  try {
    const payload = buildPayload()
    if (editingRow.value) {
      await updatePushConfig(editingRow.value.id, payload)
      message.success('更新成功')
    } else {
      await createPushConfig(payload)
      message.success('创建成功')
    }
    closeForm()
    tableRef.value?.loadData()
  } catch (err: any) {
    message.error(err?.message || '操作失败')
  } finally {
    formLoading.value = false
  }
}

async function handleToggleEnabled(row: PushConfig, val: boolean) {
  try {
    await updatePushConfig(row.id, { enabled: val })
    message.success(val ? '已启用' : '已禁用')
    tableRef.value?.loadData()
  } catch (err: any) {
    message.error(err?.message || '操作失败')
  }
}

async function handleTest(row: PushConfig) {
  testingId.value = row.id
  try {
    const res = await testPushConfig(row.id)
    testResult.value = res.data
    testResultShow.value = true
    if (res.data?.success) {
      message.success('测试完成')
    } else {
      message.warning('测试失败')
    }
  } catch (err: any) {
    testResult.value = {
      success: false,
      channel: row.type,
      statusCode: 500,
      responseBody: '',
      errorMessage: err?.message || '测试失败',
      summary: err?.message || '测试失败',
    }
    testResultShow.value = true
    message.error(err?.message || '测试失败')
  } finally {
    testingId.value = null
  }
}

function handleDelete(row: PushConfig) {
  confirmTitle.value = '删除'
  confirmContent.value = `确定要删除 "${row.name}" 吗？`
  confirmAction.value = async () => {
    await deletePushConfig(row.id)
    message.success('删除成功')
    tableRef.value?.loadData()
  }
  confirmDialogShow.value = true
}

async function handleConfirm() {
  confirmLoading.value = true
  try {
    await confirmAction.value()
  } catch (err: any) {
    message.error(err?.message || '操作失败')
  } finally {
    confirmLoading.value = false
    confirmDialogShow.value = false
  }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="推送配置" description="管理日志推送通道 (HTTP / Email / Syslog)">
      <n-button type="primary" @click="handleAdd">添加配置</n-button>
    </PageHeader>

    <n-tabs type="line">
      <n-tab-pane name="all" tab="全部">
        <DataTable
          ref="tableRef"
          :columns="columns"
          :fetch-api="fetchData"
          :search-fields="['name']"
          :extra-params="{}"
          search-placeholder="搜索配置名称"
        />
      </n-tab-pane>
    </n-tabs>

    <n-modal
      v-model:show="formShow"
      preset="card"
      :title="editingRow ? '编辑推送配置' : '添加推送配置'"
      :style="{ width: isMobile ? 'calc(100vw - 32px)' : '860px', maxWidth: 'calc(100vw - 32px)' }"
      :mask-closable="false"
    >
      <div class="modal-scroll-content">
        <n-form ref="formRef" :model="formData" :label-placement="isMobile ? 'top' : 'left'" :label-width="isMobile ? undefined : 150">
        <n-form-item label="配置名称" path="name">
          <n-input v-model:value="formData.name" placeholder="例如：安全事件 HTTP 推送" />
        </n-form-item>
        <n-form-item label="推送类型" path="type">
          <n-select v-model:value="formData.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item label="启用状态" path="enabled">
          <n-switch v-model:value="formData.enabled" />
        </n-form-item>

        <template v-if="currentType === 'http'">
          <n-form-item label="请求地址" path="url">
            <n-input v-model:value="formData.url" placeholder="http://127.0.0.1:19080/mock" />
          </n-form-item>
          <n-form-item label="请求方法" path="method">
            <n-select
              v-model:value="formData.method"
              :options="[
                { label: 'POST', value: 'POST' },
                { label: 'PUT', value: 'PUT' },
                { label: 'PATCH', value: 'PATCH' },
              ]"
            />
          </n-form-item>
          <n-form-item label="请求头 JSON" path="headers">
            <n-input
              v-model:value="formData.headers"
              type="textarea"
              placeholder='{"Authorization":"Bearer ${token}"}'
              :autosize="{ minRows: 4, maxRows: 8 }"
            />
          </n-form-item>
          <n-form-item label="请求体模板" path="bodyTemplate">
            <n-input
              v-model:value="formData.bodyTemplate"
              type="textarea"
              :autosize="{ minRows: 5, maxRows: 10 }"
              placeholder='{"message":"${message}","sourceIp":"${source_ip}"}'
            />
          </n-form-item>
          <n-form-item label="认证方式" path="authType">
            <n-select v-model:value="formData.authType" :options="authTypeOptions" />
          </n-form-item>
          <n-form-item label="认证 Token" path="token">
            <n-input v-model:value="formData.token" type="password" show-password-on="click" placeholder="可选" />
          </n-form-item>
          <n-form-item label="内容类型" path="contentType">
            <n-input v-model:value="formData.contentType" placeholder="application/json" />
          </n-form-item>
          <n-form-item label="超时(秒)" path="timeout">
            <n-input-number v-model:value="formData.timeout" :min="1" :max="300" />
          </n-form-item>
          <n-form-item label="重试次数" path="retryCount">
            <n-input-number v-model:value="formData.retryCount" :min="0" :max="10" />
          </n-form-item>
          <n-form-item label="重试间隔(秒)" path="retryDelay">
            <n-input-number v-model:value="formData.retryDelay" :min="0" :max="300" />
          </n-form-item>
          <n-form-item label="成功状态码 JSON" path="successStatusCodes">
            <n-input v-model:value="formData.successStatusCodes" placeholder="[200,201,202]" />
          </n-form-item>
          <n-form-item label="成功关键字" path="successBodyKeyword">
            <n-input v-model:value="formData.successBodyKeyword" placeholder="可选，例如 success" />
          </n-form-item>
          <n-form-item label="响应摘要长度" path="maxResponseLogSize">
            <n-input-number v-model:value="formData.maxResponseLogSize" :min="128" :max="65535" />
          </n-form-item>
          <n-form-item label="接收标识列表" path="notesIds">
            <n-dynamic-tags v-model:value="notesIdsTags" />
          </n-form-item>
          <n-form-item label="需重试状态码" path="retryOnStatusCodes">
            <n-dynamic-tags v-model:value="retryOnStatusCodesTags" />
          </n-form-item>
        </template>

        <template v-else-if="currentType === 'email'">
          <n-form-item label="SMTP 主机" path="smtpHost">
            <n-input v-model:value="formData.smtpHost" placeholder="127.0.0.1" />
          </n-form-item>
          <n-form-item label="SMTP 端口" path="smtpPort">
            <n-input-number v-model:value="formData.smtpPort" :min="1" :max="65535" />
          </n-form-item>
          <n-form-item label="SMTP 用户名" path="smtpUsername">
            <n-input v-model:value="formData.smtpUsername" placeholder="可选" />
          </n-form-item>
          <n-form-item label="SMTP 密码" path="smtpPassword">
            <n-input v-model:value="formData.smtpPassword" type="password" show-password-on="click" placeholder="可选" />
          </n-form-item>
          <n-form-item label="发件人地址" path="fromAddress">
            <n-input v-model:value="formData.fromAddress" placeholder="sender@example.com" />
          </n-form-item>
          <n-form-item label="收件人 JSON" path="toAddresses">
            <n-input v-model:value="formData.toAddresses" placeholder='["security@example.com"]' />
          </n-form-item>
          <n-form-item label="邮件标题模板" path="subjectTemplate">
            <n-input v-model:value="formData.subjectTemplate" placeholder="logcat 告警 ${event_type}" />
          </n-form-item>
          <n-form-item label="邮件正文模板" path="emailBodyTemplate">
            <n-input
              v-model:value="formData.emailBodyTemplate"
              type="textarea"
              :autosize="{ minRows: 6, maxRows: 10 }"
              placeholder="来源 IP: ${source_ip}"
            />
          </n-form-item>
        </template>

        <template v-else>
          <n-form-item label="Syslog 主机" path="syslogHost">
            <n-input v-model:value="formData.syslogHost" placeholder="127.0.0.1" />
          </n-form-item>
          <n-form-item label="Syslog 端口" path="syslogPort">
            <n-input-number v-model:value="formData.syslogPort" :min="1" :max="65535" />
          </n-form-item>
          <n-form-item label="协议" path="syslogProtocol">
            <n-select v-model:value="formData.syslogProtocol" :options="syslogProtocolOptions" />
          </n-form-item>
          <n-form-item label="输出格式" path="syslogFormat">
            <n-select v-model:value="formData.syslogFormat" :options="syslogFormatOptions" />
          </n-form-item>
          <n-form-item label="字段模板 JSON" path="syslogFields">
            <n-input
              v-model:value="formData.syslogFields"
              type="textarea"
              :autosize="{ minRows: 5, maxRows: 10 }"
              placeholder='{"message":"${message}","source_ip":"${source_ip}"}'
            />
          </n-form-item>
        </template>
      </n-form>
      </div>

      <template #footer>
        <n-space justify="end">
          <n-button @click="closeForm">取消</n-button>
          <n-button type="primary" :loading="formLoading" @click="handleSubmit">保存</n-button>
        </n-space>
      </template>
    </n-modal>

    <n-modal
      v-model:show="testResultShow"
      preset="card"
      title="推送测试结果"
      :style="{ width: isMobile ? 'calc(100vw - 32px)' : '760px', maxWidth: 'calc(100vw - 32px)' }"
    >
      <n-space vertical :size="16">
        <n-alert :type="testResult?.success ? 'success' : 'error'" :show-icon="false">
          {{ testResult?.summary || (testResult?.success ? '测试成功' : '测试失败') }}
        </n-alert>
        <n-descriptions bordered label-placement="left" :column="isMobile ? 1 : 2">
          <n-descriptions-item label="通道">
            {{ testResult?.channel?.toUpperCase() || '-' }}
          </n-descriptions-item>
          <n-descriptions-item label="状态码">
            {{ testResult?.statusCode ?? '-' }}
          </n-descriptions-item>
          <n-descriptions-item label="执行结果">
            {{ testResult?.success ? '成功' : '失败' }}
          </n-descriptions-item>
          <n-descriptions-item label="摘要">
            {{ testResult?.summary || '-' }}
          </n-descriptions-item>
        </n-descriptions>
        <div>
          <div class="result-label">错误信息</div>
          <n-code :code="testResult?.errorMessage || '-'" language="text" word-wrap />
        </div>
        <div>
          <div class="result-label">响应 / 返回内容</div>
          <n-code :code="testResult?.responseBody || '-'" language="text" word-wrap />
        </div>
      </n-space>
    </n-modal>

    <ConfirmDialog
      v-model:show="confirmDialogShow"
      :title="confirmTitle"
      :content="confirmContent"
      :loading="confirmLoading"
      @confirm="handleConfirm"
    />
  </div>
</template>

<style scoped>
.result-label {
  margin-bottom: 8px;
  font-weight: 600;
}

.modal-scroll-content {
  max-height: calc(100vh - 200px);
  overflow-y: auto;
}
</style>
