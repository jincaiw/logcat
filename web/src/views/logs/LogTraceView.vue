<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NCard, NTag, NTimeline, NTimelineItem, NSpin, NSpace, NButton, NIcon } from 'naive-ui'
import { ArrowBackOutline } from '@vicons/ionicons5'
import { getLogTrace } from '@/api/logs'
import type { SyslogLog } from '@/types'
import StatusBadge from '@/components/common/StatusBadge.vue'
import PageHeader from '@/components/common/PageHeader.vue'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const log = ref<SyslogLog | null>(null)
const trace = ref<any[]>([])

const logId = computed(() => route.params.id as string)

async function loadTrace() {
  if (!logId.value) return
  loading.value = true
  try {
    const res = await getLogTrace(logId.value)
    log.value = res.data?.log || null
    trace.value = res.data?.trace || buildDefaultTrace(log.value)
  } catch {
    log.value = null
    trace.value = []
  } finally {
    loading.value = false
  }
}

function buildDefaultTrace(l: SyslogLog | null): any[] {
  if (!l) return []
  return [
    { stage: '接收', status: 'success' as const, detail: `${l.deviceName} (${l.deviceHost})`, time: l.receivedAt },
    { stage: '解析', status: 'success' as const, detail: l.parsedData ? '解析完成' : '原始消息', time: l.receivedAt },
    { stage: '过滤', status: 'success' as const, detail: '通过过滤策略', time: l.receivedAt },
    { stage: '去重', status: 'success' as const, detail: '无重复', time: l.receivedAt },
    { stage: '聚合', status: 'success' as const, detail: '聚合处理', time: l.receivedAt },
    { stage: '推送', status: l.pushStatus === 'success' ? 'success' as const : l.pushStatus === 'failed' ? 'error' as const : 'warning' as const, detail: l.pushStatus, time: l.receivedAt },
  ]
}

onMounted(() => {
  loadTrace()
})
</script>

<template>
  <div class="page-container">
    <PageHeader title="日志追踪">
      <n-button text @click="router.push('/logs')">
        <template #icon><n-icon :component="ArrowBackOutline" /></template>
        返回
      </n-button>
    </PageHeader>

    <n-spin :show="loading">
      <n-card v-if="log" size="small" style="margin-bottom: 16px">
        <n-space vertical>
          <div><strong>日志ID:</strong> {{ log.id }}</div>
          <div><strong>时间:</strong> {{ log.receivedAt }}</div>
          <div><strong>设备:</strong> {{ log.deviceName }} ({{ log.deviceHost }})</div>
          <div><strong>源IP:</strong> {{ log.sourceIp }}</div>
          <div><strong>严重程度:</strong> <StatusBadge v-if="log.severity" :status="log.severity" type="severity" /></div>
          <div>
            <strong>原始消息:</strong>
            <pre style="margin-top: 4px; padding: 8px; background: var(--bg-color); border-radius: 4px; font-size: 12px; overflow-x: auto">{{ log.rawMessage }}</pre>
          </div>
        </n-space>
      </n-card>

      <n-card title="处理链路" size="small">
        <n-timeline>
          <n-timeline-item
            v-for="(item, index) in trace"
            :key="index"
            :type="item.status === 'error' ? 'error' : item.status === 'warning' ? 'warning' : 'success'"
            :title="item.stage"
            :time="item.time"
          >
            <div>{{ item.detail || '-' }}</div>
          </n-timeline-item>
        </n-timeline>
      </n-card>
    </n-spin>
  </div>
</template>