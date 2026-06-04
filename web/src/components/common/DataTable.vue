<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import {
  NDataTable, NInput, NButton, NPagination, NIcon, NSpin,
} from 'naive-ui'
import { SearchOutline, RefreshOutline } from '@vicons/ionicons5'
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import type { PageResponse } from '@/types'
import { useIsMobile } from '@/composables/useIsMobile'
import { useTimeFormat } from '@/composables/useTimeFormat'

const TIME_COLUMN_KEYS = new Set([
  'createdAt', 'updatedAt', 'lastLoginAt', 'receivedAt',
  'lastSeenAt', 'firstSeenAt', 'lastReceivedAt', 'firstAt',
  'occurredAt', 'lockedUntil', 'expiredAt',
  'firstSeen', 'lastSeen', 'sentAt',
])

const props = withDefaults(defineProps<{
  columns: DataTableColumns<any>
  fetchApi: (params: any) => Promise<any>
  searchPlaceholder?: string
  searchFields?: string[]
  defaultPageSize?: number
  rowKey?: string
  showSearch?: boolean
  showRefresh?: boolean
  showPagination?: boolean
  extraParams?: Record<string, any>
}>(), {
  searchPlaceholder: '搜索...',
  searchFields: () => ['name'],
  defaultPageSize: 20,
  rowKey: 'id',
  showSearch: true,
  showRefresh: true,
  showPagination: true,
  extraParams: () => ({}),
})

const emit = defineEmits<{
  (e: 'row-click', row: any): void
  (e: 'data-loaded', data: PageResponse): void
}>()

const data = ref<any[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(props.defaultPageSize)
const searchKeyword = ref('')
const checkedRowKeys = ref<DataTableRowKey[]>([])
const { isMobile } = useIsMobile()
const { formatTime } = useTimeFormat()

const processedColumns = computed(() => {
  return props.columns.map((col: any) => {
    if (col.render || col.type) return col
    const key = col.key as string
    if (key && TIME_COLUMN_KEYS.has(key)) {
      return {
        ...col,
        render: (row: any) => formatTime(row[key]),
      }
    }
    return col
  })
})

const scrollX = computed(() => {
  let totalWidth = 0
  props.columns.forEach((col: any) => {
    if (col.width) {
      totalWidth += Number(col.width)
    }
  })
  return totalWidth > 0 ? totalWidth : undefined
})

const paginationReactive = ref(false)

function normalizeTablePayload(payload: any): PageResponse | { list: any[]; total: number; page: number; pageSize: number } {
  if (Array.isArray(payload)) {
    return {
      list: payload,
      total: payload.length,
      page: 1,
      pageSize: payload.length,
    }
  }

  const list = payload?.items ?? payload?.list ?? payload?.rows ?? []
  return {
    ...payload,
    list: Array.isArray(list) ? list : [],
    total: Number(payload?.total ?? payload?.count ?? 0),
    page: Number(payload?.page ?? 1),
    pageSize: Number(payload?.pageSize ?? list.length ?? 0),
  }
}

async function loadData() {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      pageSize: pageSize.value,
      ...props.extraParams,
    }
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    const res = await props.fetchApi(params)
    if (res && res.data) {
      const pageData = normalizeTablePayload(res.data)
      data.value = pageData.list || []
      total.value = pageData.total || 0
      emit('data-loaded', pageData as PageResponse)
    }
  } catch (err) {
    console.error('Failed to load data:', err)
  } finally {
    loading.value = false
    paginationReactive.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadData()
}

function handlePageChange(p: number) {
  page.value = p
  loadData()
}

function handlePageSizeChange(ps: number) {
  pageSize.value = ps
  page.value = 1
  loadData()
}

function handleRefresh() {
  loadData()
}

function handleRowClick(row: any) {
  emit('row-click', row)
}

onMounted(() => {
  loadData()
})

watch(() => props.extraParams, () => {
  page.value = 1
  loadData()
}, { deep: true })

defineExpose({
  loadData,
  refresh: handleRefresh,
  data,
  loading,
})
</script>

<template>
  <div class="data-table-wrapper">
    <div class="table-toolbar">
      <div class="table-toolbar-left">
        <slot name="toolbar-left" />
        <div v-if="showSearch" class="search-group">
          <n-input
            v-model:value="searchKeyword"
            :placeholder="searchPlaceholder"
            clearable
            size="small"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <n-icon :component="SearchOutline" style="color: var(--text-color-tertiary)" />
            </template>
          </n-input>
          <n-button size="small" type="primary" @click="handleSearch">
            搜索
          </n-button>
        </div>
        <n-button v-if="showRefresh" size="small" quaternary class="refresh-btn" @click="handleRefresh">
          <template #icon>
            <n-icon :component="RefreshOutline" />
          </template>
        </n-button>
      </div>
      <div class="table-toolbar-right">
        <slot name="toolbar-right" />
      </div>
    </div>

    <n-spin :show="loading">
      <n-data-table
        :columns="processedColumns"
        :data="data"
        :row-key="(row: any) => row[props.rowKey]"
        :loading="loading"
        :bordered="false"
        :single-line="false"
        :scroll-x="scrollX"
        size="small"
        striped
        @update:checked-row-keys="(keys: DataTableRowKey[]) => checkedRowKeys = keys"
        @row-click="handleRowClick"
      />
    </n-spin>

    <div v-if="showPagination" class="pagination-wrapper">
      <n-pagination
        v-model:page="page"
        :page-size="pageSize"
        :item-count="total"
        :page-sizes="isMobile ? [10, 20] : [10, 20, 50, 100]"
        :show-size-picker="!isMobile"
        :disabled="paginationReactive"
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
      />
    </div>
  </div>
</template>

<style scoped>
.data-table-wrapper {
  background: var(--bg-color-card);
  border-radius: var(--radius-md);
  padding: 16px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-card);
  transition: background-color 0.3s ease, border-color 0.3s ease;
}

.search-group {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  max-width: 320px;
  min-width: 200px;
}

.refresh-btn {
  color: var(--text-color-secondary);
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
