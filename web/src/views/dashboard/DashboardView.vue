<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  NCard, NGrid, NGi, NSpace, NStatistic, NIcon, NTag, NProgress, NDivider, NSpin,
} from 'naive-ui'
import {
  ServerOutline, TrendingUpOutline, DocumentTextOutline,
  WarningOutline, CheckmarkCircleOutline, LayersOutline,
  CloseCircleOutline,
} from '@vicons/ionicons5'
import { getDashboardStats } from '@/api/dashboard'
import type { DashboardStats } from '@/types'
import PageHeader from '@/components/common/PageHeader.vue'

const stats = ref<DashboardStats | null>(null)
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
</script>

<template>
  <div class="page-container">
    <PageHeader title="仪表盘" description="系统运行概览" />

    <n-spin :show="loading">
      <div class="stats-grid">
        <!-- Service Status -->
        <n-card size="small" hoverable>
          <n-space align="center" justify="center" vertical>
            <n-icon size="28" color="var(--primary-color)" :component="ServerOutline" />
            <n-statistic label="服务状态" :value="stats?.serviceStatus || '--'" />
            <n-tag :type="stats?.serviceStatus === 'running' ? 'success' : 'default'" size="small">
              {{ stats?.serviceStatus === 'running' ? '运行中' : '已停止' }}
            </n-tag>
          </n-space>
        </n-card>

        <!-- Receive Rate -->
        <n-card size="small" hoverable>
          <n-space align="center" justify="center" vertical>
            <n-icon size="28" color="#2080f0" :component="TrendingUpOutline" />
            <n-statistic label="接收速率">
              {{ stats?.receiveRate ? stats.receiveRate.toFixed(0) + ' /s' : '-- /s' }}
            </n-statistic>
          </n-space>
        </n-card>

        <!-- Today Total -->
        <n-card size="small" hoverable>
          <n-space align="center" justify="center" vertical>
            <n-icon size="28" color="#18a058" :component="DocumentTextOutline" />
            <n-statistic label="今日日志量">
              {{ stats ? formatNumber(stats.todayTotal) : '--' }}
            </n-statistic>
          </n-space>
        </n-card>

        <!-- Today Alerts -->
        <n-card size="small" hoverable>
          <n-space align="center" justify="center" vertical>
            <n-icon size="28" color="#f0a020" :component="WarningOutline" />
            <n-statistic label="今日告警">
              {{ stats?.todayAlerts ?? '--' }}
            </n-statistic>
          </n-space>
        </n-card>

        <!-- Push Success Rate -->
        <n-card size="small" hoverable>
          <n-space align="center" justify="center" vertical>
            <n-icon size="28" color="var(--primary-color)" :component="CheckmarkCircleOutline" />
            <n-statistic label="推送成功率">
              {{ stats?.pushSuccessRate ? (stats.pushSuccessRate * 100).toFixed(1) + '%' : '--' }}
            </n-statistic>
          </n-space>
        </n-card>
      </div>

      <!-- Queue Backlog -->
      <n-divider title-placement="left">队列积压</n-divider>
      <n-grid v-if="stats?.queueBacklog?.length" cols="1 s:2 m:3 l:4" x-gap="12" y-gap="12" style="margin-bottom: 24px">
        <n-gi v-for="queue in stats.queueBacklog" :key="queue.name">
          <n-card size="small">
            <n-space vertical>
              <span style="font-weight: 600">{{ queue.name }}</span>
              <n-progress
                type="line"
                :percentage="queue.capacity > 0 ? Math.round((queue.size / queue.capacity) * 100) : 0"
                :color="queue.size / queue.capacity > 0.8 ? '#d03050' : '#18a058'"
                :height="20"
              >
                {{ queue.size }} / {{ queue.capacity }}
              </n-progress>
              <div style="font-size: 12px; color: var(--text-color-secondary)">
                入队: {{ queue.enqueueRate?.toFixed(1) }}/s | 出队: {{ queue.dequeueRate?.toFixed(1) }}/s
              </div>
            </n-space>
          </n-card>
        </n-gi>
      </n-grid>
      <n-card v-else size="small">
        <span style="color: var(--text-color-secondary)">暂无队列数据</span>
      </n-card>

      <!-- Recent Push Failures -->
      <n-divider title-placement="left">近期推送失败</n-divider>
      <n-card v-if="stats?.recentPushFailures?.length" size="small" style="margin-bottom: 24px">
        <div
          v-for="(item, index) in stats.recentPushFailures"
          :key="index"
          style="padding: 8px 0; border-bottom: 1px solid var(--border-color); display: flex; justify-content: space-between"
        >
          <div>
            <n-icon :component="CloseCircleOutline" color="#d03050" size="14" />
            <span style="margin-left: 8px; font-size: 13px">{{ item.channel }}</span>
            <n-tag size="tiny" style="margin-left: 8px">{{ item.logId }}</n-tag>
          </div>
          <span style="font-size: 12px; color: var(--text-color-secondary)">{{ item.time }}</span>
        </div>
        <div v-if="stats.recentPushFailures.length === 0" style="color: var(--text-color-secondary); text-align: center; padding: 16px">
          暂无推送失败记录
        </div>
      </n-card>
      <n-card v-else size="small">
        <span style="color: var(--text-color-secondary)">暂无推送失败记录</span>
      </n-card>

      <!-- Health Status -->
      <n-divider title-placement="left">健康状态</n-divider>
      <n-grid v-if="stats?.healthStatus" cols="1 s:2 m:3 l:5" x-gap="12" y-gap="12">
        <n-gi>
          <n-card size="small">
            <n-statistic label="CPU 使用率" :value="`${stats.healthStatus.cpu?.toFixed(1)}%`" />
            <n-progress
              type="line"
              :percentage="stats.healthStatus.cpu || 0"
              :color="(stats.healthStatus.cpu || 0) > 80 ? '#d03050' : '#18a058'"
              :height="10"
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
              :color="(stats.healthStatus.memory || 0) > 80 ? '#d03050' : '#18a058'"
              :height="10"
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
              :color="(stats.healthStatus.diskUsage || 0) > 80 ? '#d03050' : '#18a058'"
              :height="10"
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