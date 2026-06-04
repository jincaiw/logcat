<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  NCard, NGrid, NGi, NSpace, NStatistic, NIcon, NTag, NProgress, NDivider, NSpin,
} from 'naive-ui'
import {
  ServerOutline, TrendingUpOutline, DocumentTextOutline,
  WarningOutline, CheckmarkCircleOutline, LayersOutline,
  CloseCircleOutline, GitNetworkOutline, TimeOutline, PieChartOutline,
} from '@vicons/ionicons5'
import { getDashboardStats } from '@/api/dashboard'
import type { DashboardStats } from '@/types'
import PageHeader from '@/components/common/PageHeader.vue'
import { useTimeFormat } from '@/composables/useTimeFormat'

const stats = ref<DashboardStats | null>(null)
const { formatTime } = useTimeFormat()
const loading = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

async function loadStats() {
  loading.value = true
  try {
    const res = await getDashboardStats()
    stats.value = res.data
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStats()
  timer = setInterval(loadStats, 10000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

function formatNumber(num: number): string {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return String(num)
}

function getProgressColor(ratio: number): string {
  return ratio > 0.8 ? 'var(--error-color)' : 'var(--success-color)'
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="仪表盘" description="系统运行概览" />

    <n-spin :show="loading">
      <div class="stats-grid">
        <n-card size="small" hoverable class="stat-item">
          <n-space align="center" justify="center" vertical>
            <div class="stat-icon stat-icon--primary">
              <n-icon size="24" :component="ServerOutline" />
            </div>
            <n-statistic label="服务状态" :value="stats?.serviceStatus || '--'" />
            <n-tag :type="stats?.serviceStatus === 'running' ? 'success' : 'default'" size="small" round>
              {{ stats?.serviceStatus === 'running' ? '运行中' : '已停止' }}
            </n-tag>
          </n-space>
        </n-card>

        <n-card size="small" hoverable class="stat-item">
          <n-space align="center" justify="center" vertical>
            <div class="stat-icon stat-icon--info">
              <n-icon size="24" :component="TrendingUpOutline" />
            </div>
            <n-statistic label="接收速率">
              {{ stats?.receiveRate != null ? stats.receiveRate.toFixed(0) + ' /s' : '-- /s' }}
            </n-statistic>
          </n-space>
        </n-card>

        <n-card size="small" hoverable class="stat-item">
          <n-space align="center" justify="center" vertical>
            <div class="stat-icon stat-icon--success">
              <n-icon size="24" :component="DocumentTextOutline" />
            </div>
            <n-statistic label="今日日志量">
              {{ stats ? formatNumber(stats.todayTotal) : '--' }}
            </n-statistic>
          </n-space>
        </n-card>

        <n-card size="small" hoverable class="stat-item">
          <n-space align="center" justify="center" vertical>
            <div class="stat-icon stat-icon--warning">
              <n-icon size="24" :component="WarningOutline" />
            </div>
            <n-statistic label="今日告警">
              {{ stats?.todayAlerts ?? '--' }}
            </n-statistic>
          </n-space>
        </n-card>

        <n-card size="small" hoverable class="stat-item">
          <n-space align="center" justify="center" vertical>
            <div class="stat-icon stat-icon--primary">
              <n-icon size="24" :component="CheckmarkCircleOutline" />
            </div>
            <n-statistic label="推送成功率">
              {{ stats?.pushSuccessRate != null ? (stats.pushSuccessRate * 100).toFixed(1) + '%' : '--' }}
            </n-statistic>
          </n-space>
        </n-card>

        <n-card size="small" hoverable class="stat-item">
          <n-space align="center" justify="center" vertical>
            <div class="stat-icon stat-icon--success">
              <n-icon size="24" :component="GitNetworkOutline" />
            </div>
            <n-statistic label="TCP 连接数">
              {{ stats?.tcpConnections ?? '--' }}
            </n-statistic>
          </n-space>
        </n-card>

        <n-card size="small" hoverable class="stat-item">
          <n-space align="center" justify="center" vertical>
            <div class="stat-icon stat-icon--warning">
              <n-icon size="24" :component="TimeOutline" />
            </div>
            <n-statistic label="最后接收日志">
              {{ formatTime(stats?.lastReceivedAt) }}
            </n-statistic>
          </n-space>
        </n-card>

        <n-card size="small" hoverable class="stat-item">
          <n-space align="center" justify="center" vertical>
            <div class="stat-icon stat-icon--info">
              <n-icon size="24" :component="PieChartOutline" />
            </div>
            <n-statistic label="解析成功率">
              {{ stats?.parseSuccessRate != null ? (stats.parseSuccessRate * 100).toFixed(1) + '%' : '--' }}
            </n-statistic>
          </n-space>
        </n-card>
      </div>

      <n-divider title-placement="left">队列积压</n-divider>
      <n-grid v-if="stats?.queueBacklog?.length" cols="1 s:2 m:3 l:4" x-gap="12" y-gap="12" style="margin-bottom: 24px">
        <n-gi v-for="queue in stats.queueBacklog" :key="queue.name">
          <n-card size="small">
            <n-space vertical>
              <span style="font-weight: 600">{{ queue.name }}</span>
              <n-progress
                type="line"
                :percentage="queue.capacity > 0 ? Math.round((queue.size / queue.capacity) * 100) : 0"
                :color="getProgressColor(queue.capacity > 0 ? queue.size / queue.capacity : 0)"
                :height="20"
                :border-radius="4"
              >
                {{ queue.size }} / {{ queue.capacity }}
              </n-progress>
              <div class="stat-sub-text">
                入队: {{ queue.enqueueRate?.toFixed(1) }}/s | 出队: {{ queue.dequeueRate?.toFixed(1) }}/s
              </div>
            </n-space>
          </n-card>
        </n-gi>
      </n-grid>
      <n-card v-else size="small">
        <span class="stat-sub-text">暂无队列数据</span>
      </n-card>

      <n-divider title-placement="left">近期推送失败</n-divider>
      <n-card v-if="stats?.recentPushFailures?.length" size="small" style="margin-bottom: 24px">
        <div
          v-for="(item, index) in stats.recentPushFailures"
          :key="index"
          class="push-failure-item"
        >
          <div>
            <n-icon :component="CloseCircleOutline" color="var(--error-color)" size="14" />
            <span style="margin-left: 8px; font-size: 13px">{{ item.channel }}</span>
            <n-tag size="tiny" style="margin-left: 8px">{{ item.logId }}</n-tag>
          </div>
          <span class="stat-sub-text">{{ formatTime(item.time) }}</span>
        </div>
      </n-card>
      <n-card v-else size="small">
        <span class="stat-sub-text">暂无推送失败记录</span>
      </n-card>

      <n-divider title-placement="left">健康状态</n-divider>
      <n-grid v-if="stats?.healthStatus" cols="1 s:2 m:3 l:4" x-gap="12" y-gap="12">
        <n-gi>
          <n-card size="small">
            <n-statistic label="CPU 使用率" :value="`${stats.healthStatus.cpu?.toFixed(1)}%`" />
            <n-progress
              type="line"
              :percentage="stats.healthStatus.cpu || 0"
              :color="getProgressColor((stats.healthStatus.cpu || 0) / 100)"
              :height="10"
              :border-radius="4"
              style="margin-top: 8px"
            />
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic label="内存使用率" :value="`${stats.healthStatus.memory?.toFixed(1)}%`" />
            <n-progress
              type="line"
              :percentage="stats.healthStatus.memory || 0"
              :color="getProgressColor((stats.healthStatus.memory || 0) / 100)"
              :height="10"
              :border-radius="4"
              style="margin-top: 8px"
            />
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic label="磁盘使用率" :value="`${stats.healthStatus.diskUsage?.toFixed(1)}%`" />
            <n-progress
              type="line"
              :percentage="stats.healthStatus.diskUsage || 0"
              :color="getProgressColor((stats.healthStatus.diskUsage || 0) / 100)"
              :height="10"
              :border-radius="4"
              style="margin-top: 8px"
            />
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic label="网络入">
              {{ stats.healthStatus.networkIn ? (stats.healthStatus.networkIn / 1024).toFixed(1) + ' KB/s' : '--' }}
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic label="网络出">
              {{ stats.healthStatus.networkOut ? (stats.healthStatus.networkOut / 1024).toFixed(1) + ' KB/s' : '--' }}
            </n-statistic>
          </n-card>
        </n-gi>
      </n-grid>
    </n-spin>
  </div>
</template>

<style scoped>
.stat-item {
  text-align: center;
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 4px;
}

.stat-icon--primary {
  background: rgba(24, 160, 88, 0.1);
  color: var(--primary-color);
}

html.dark .stat-icon--primary {
  background: rgba(99, 226, 183, 0.12);
}

.stat-icon--info {
  background: rgba(32, 128, 240, 0.1);
  color: var(--info-color);
}

html.dark .stat-icon--info {
  background: rgba(112, 192, 232, 0.12);
}

.stat-icon--success {
  background: rgba(24, 160, 88, 0.1);
  color: var(--success-color);
}

html.dark .stat-icon--success {
  background: rgba(99, 226, 183, 0.12);
}

.stat-icon--warning {
  background: rgba(240, 160, 32, 0.1);
  color: var(--warning-color);
}

html.dark .stat-icon--warning {
  background: rgba(242, 201, 125, 0.12);
}

.stat-sub-text {
  font-size: 12px;
  color: var(--text-color-secondary);
}

.push-failure-item {
  padding: 8px 0;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 4px;
}

.push-failure-item:last-child {
  border-bottom: none;
}
</style>
