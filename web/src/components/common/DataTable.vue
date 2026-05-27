<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import {
  NDataTable, NInput, NButton, NSpace, NPagination, NIcon, NSpin,
} from 'naive-ui'
import { SearchOutline, RefreshOutline } from '@vicons/ionicons5'
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import type { PageResponse } from '@/types'

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

const paginationReactive = ref(false)

async function loadData() {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      pageSize: pageSize.value,
      ...props.extraParams,
    }
    if (searchKeyword.value && props.searchFields.length > 0) {
      props.searchFields.forEach((field) => {
        params[field] = searchKeyword.value
      })
    }
    const res = await props.fetchApi(params)
    if (res && res.data) {
      data.value = res.data.list || []
      total.value = res.data.total || 0
      emit('data-loaded', res.data)
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
        <n-input
          v-if="showSearch"
          v-model:value="searchKeyword"
          :placeholder="searchPlaceholder"
          clearable
          style="width: 240px"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <n-icon :component="SearchOutline" />
          </template>
        </n-input>
        <n-button v-if="showSearch" size="small" type="primary" @click="handleSearch">
          搜索
        </n-button>
        <n-button v-if="showRefresh" size="small" quaternary @click="handleRefresh">
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
        :columns="columns"
        :data="data"
        :row-key="(row: any) => row[props.rowKey]"
        :loading="loading"
        :bordered="false"
        :single-line="false"
        size="small"
        striped
        @update:checked-row-keys="(keys: DataTableRowKey[]) => checkedRowKeys = keys"
        @row-click="handleRowClick"
      />
    </n-spin>

    <div v-if="showPagination" style="display: flex; justify-content: flex-end; margin-top: 16px">
      <n-pagination
        v-model:page="page"
        :page-size="pageSize"
        :item-count="total"
        :page-sizes="[10, 20, 50, 100]"
        show-size-picker
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
  border-radius: 8px;
  padding: 16px;
  border: 1px solid var(--border-color);
}
</style>