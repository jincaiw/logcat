<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NSpace, NCard, NGrid, NGi, NUpload, NTag, useMessage } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { exportData, importData, type ExportType } from '@/api/importExport'
import PageHeader from '@/components/common/PageHeader.vue'

const message = useMessage()
const exportLoading = ref('')
const importLoading = ref('')
const importResult = ref('')

const types: { type: ExportType; label: string }[] = [
  { type: 'devices', label: '设备' },
  { type: 'deviceGroups', label: '设备分组' },
  { type: 'deviceTemplates', label: '设备模板' },
  { type: 'parseTemplates', label: '解析模板' },
  { type: 'filterPolicies', label: '过滤策略' },
  { type: 'outputTemplates', label: '输出模板' },
  { type: 'pushConfigs', label: '推送配置' },
  { type: 'alertRules', label: '告警规则' },
  { type: 'desensitizeRules', label: '脱敏规则' },
  { type: 'roles', label: '角色' },
]

async function handleExport(type: ExportType) {
  exportLoading.value = type
  try {
    const res = await exportData(type)
    if (res.data?.url) window.open(res.data.url)
    else message.success('导出成功')
  } catch (err: any) { message.error(err?.message || '导出失败') }
  finally { exportLoading.value = '' }
}

async function handleImport(type: ExportType, data: { file: UploadFileInfo; fileList: UploadFileInfo[] }) {
  if (!data.file.file) return
  importLoading.value = type
  importResult.value = ''
  try {
    const res = await importData(type, data.file.file)
    importResult.value = `导入: ${res.data?.imported || 0} 成功, ${res.data?.failed || 0} 失败`
    message.success(importResult.value)
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
              accept=".json,.yaml,.yml,.csv"
              @change="(data: any) => handleImport(item.type, data)"
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

    <n-tag v-if="importResult" type="success" style="margin-top: 16px">{{ importResult }}</n-tag>
  </div>
</template>