<script setup lang="ts">
import { ref, computed } from 'vue'
import { NButton, NSpace, NCard, NGrid, NGi, NUpload, NTag, NModal, NList, NListItem, NText } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { exportData, importData, type ExportType } from '@/api/importExport'
import { useAppMessage } from '@/composables/useMessage'
import { useIsMobile } from '@/composables/useIsMobile'
import PageHeader from '@/components/common/PageHeader.vue'
import type { ImportResult } from '@/types'

const message = useAppMessage()
const { isMobile } = useIsMobile()
const exportLoading = ref('')
const importLoading = ref('')
const importResult = ref<ImportResult | null>(null)
const importResultText = ref('')

const previewVisible = ref(false)
const previewContent = ref('')
const previewType = ref<ExportType>('device-templates')
const previewFileContent = ref('')

const types: { type: ExportType; label: string }[] = [
  { type: 'device-templates', label: '设备模板' },
  { type: 'parse-templates', label: '解析模板' },
  { type: 'filter-policies', label: '过滤策略' },
  { type: 'push-configs', label: '推送配置' },
]

const previewLabel = computed(() => types.find(t => t.type === previewType.value)?.label || '')

async function handleExport(type: ExportType) {
  exportLoading.value = type
  try {
    const res = await exportData(type)
    if (res.data?.url) window.open(res.data.url)
    message.success(`导出成功，共 ${res.data?.count || 0} 条`)
  } catch (err: any) { message.error(err?.message || '导出失败') }
  finally { exportLoading.value = '' }
}

function handleFileSelect(type: ExportType, data: { file: UploadFileInfo; fileList: UploadFileInfo[] }) {
  if (!data.file.file) return
  const file = data.file.file

  if (!file.name.endsWith('.json')) {
    message.error('仅支持 .json 格式文件')
    return
  }

  file.text().then((content) => {
    let parsed: any
    try {
      parsed = JSON.parse(content)
    } catch {
      message.error('文件内容不是有效的 JSON 格式')
      return
    }

    if (!parsed.version) {
      message.error('导入文件缺少 version 字段，格式不正确')
      return
    }

    previewFileContent.value = content
    previewType.value = type
    try {
      previewContent.value = JSON.stringify(parsed, null, 2)
    } catch {
      previewContent.value = content
    }
    previewVisible.value = true
  }).catch(() => {
    message.error('读取文件失败')
  })
}

function cancelPreview() {
  previewVisible.value = false
  previewContent.value = ''
  previewFileContent.value = ''
}

async function confirmImport() {
  previewVisible.value = false
  const type = previewType.value
  const content = previewFileContent.value

  importLoading.value = type
  importResult.value = null
  importResultText.value = ''
  try {
    const res = await importData(type, content)
    importResult.value = res.data
    importResultText.value = `导入完成：新增 ${res.data?.created || 0}，更新 ${res.data?.updated || 0}，失败 ${res.data?.failed || 0}`
    if (res.data?.failed > 0) {
      message.warning(importResultText.value)
    } else {
      message.success(importResultText.value)
    }
  } catch (err: any) { message.error(err?.message || '导入失败') }
  finally { importLoading.value = '' }
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="导入导出" description="数据导入与导出" />

    <n-grid cols="1 s:2 m:3 l:4" x-gap="12" y-gap="12">
      <n-gi v-for="item in types" :key="item.type">
        <n-card size="small" :title="item.label">
          <n-space vertical>
            <n-button
              size="small"
              :loading="exportLoading === item.type"
              @click="handleExport(item.type)"
            >
              导出 {{ item.label }}
            </n-button>
            <n-upload
              :show-file-list="false"
              accept=".json"
              @change="(data: any) => handleFileSelect(item.type, data)"
            >
              <n-button
                size="small"
                :loading="importLoading === item.type"
              >
                导入 {{ item.label }}
              </n-button>
            </n-upload>
          </n-space>
        </n-card>
      </n-gi>
    </n-grid>

    <n-tag v-if="importResultText" :type="importResult && importResult.failed > 0 ? 'warning' : 'success'" style="margin-top: 16px">{{ importResultText }}</n-tag>

    <n-list v-if="importResult && importResult.errors && importResult.errors.length > 0" bordered style="margin-top: 12px">
      <n-list-item v-for="(err, idx) in importResult.errors" :key="idx">
        <n-text type="error">{{ err }}</n-text>
      </n-list-item>
    </n-list>

    <n-modal
      v-model:show="previewVisible"
      preset="dialog"
      :title="`确认导入 ${previewLabel}`"
      positive-text="确认导入"
      negative-text="取消"
      :style="{ width: isMobile ? 'calc(100vw - 32px)' : '600px', maxWidth: 'calc(100vw - 32px)' }"
      @positive-click="confirmImport"
      @negative-click="cancelPreview"
    >
      <n-text>以下是将要导入的文件内容预览：</n-text>
      <pre style="max-height: 300px; overflow: auto; background: var(--bg-color-embedded); padding: 12px; border-radius: 4px; margin-top: 8px; font-size: 12px; white-space: pre-wrap; word-break: break-all;">{{ previewContent }}</pre>
    </n-modal>
  </div>
</template>
