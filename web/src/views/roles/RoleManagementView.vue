<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NSpace, NCheckbox, NCheckboxGroup } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createRole, updateRole, deleteRole, getRoles, getRolePermissions, assignPermissions, getAllPermissions } from '@/api/roles'
import type { Role, Permission } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import FormDialog, { type FieldConfig } from '@/components/common/FormDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAppMessage } from '@/composables/useMessage'
import { useIsMobile } from '@/composables/useIsMobile'

const message = useAppMessage()
const { isMobile } = useIsMobile()

const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const formDialogRef = ref<InstanceType<typeof FormDialog> | null>(null)
const permDialogShow = ref(false)
const confirmDialogShow = ref(false)
const confirmTitle = ref('')
const confirmContent = ref('')
const confirmAction = ref<() => Promise<void>>(() => Promise.resolve())
const confirmLoading = ref(false)

const editingRole = ref<Role | null>(null)
const allPermissions = ref<Permission[]>([])
const selectedPermIds = ref<string[]>([])
const permLoading = ref(false)

const formFields: FieldConfig[] = [
  { key: 'name', label: '角色名称', type: 'text', required: true },
  { key: 'code', label: '角色编码', type: 'text', required: true },
  { key: 'description', label: '描述', type: 'textarea' },
]

const columns: DataTableColumns<Role> = [
  { title: '角色名称', key: 'name' },
  { title: '角色编码', key: 'code' },
  { title: '描述', key: 'description' },
  { title: '创建时间', key: 'createdAt' },
  {
    title: '操作', key: 'actions',
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
          h(NButton, { size: 'small', onClick: () => handlePermission(row) }, { default: () => '分配权限' }),
          h(NButton, { size: 'small', type: 'error', ghost: true, onClick: () => handleDelete(row) }, { default: () => '删除' }),
        ],
      })
    },
  },
]

async function fetchData(params: any) {
  const res = await getRoles(params)
  return res
}

function handleAdd() {
  editingRole.value = null
  formDialogRef.value?.open()
}

function handleEdit(row: Role) {
  editingRole.value = row
  formDialogRef.value?.open({
    name: row.name,
    code: row.code,
    description: row.description,
  })
}

async function handleFormSubmit(data: Record<string, any>) {
  try {
    if (editingRole.value) {
      await updateRole(editingRole.value.id, data)
      message.success('角色更新成功')
    } else {
      await createRole(data)
      message.success('角色创建成功')
    }
    formDialogRef.value?.close()
    tableRef.value?.loadData()
  } catch (err: any) {
    message.error(err?.message || '操作失败')
  }
}

function handleDelete(row: Role) {
  confirmTitle.value = '删除角色'
  confirmContent.value = `确定要删除角色 "${row.name}" 吗？`
  confirmAction.value = async () => {
    await deleteRole(row.id)
    message.success('删除成功')
    tableRef.value?.loadData()
  }
  confirmDialogShow.value = true
}

async function handlePermission(row: Role) {
  editingRole.value = row
  permLoading.value = true
  try {
    const [permRes, rolePermRes] = await Promise.all([
      getAllPermissions(),
      getRolePermissions(row.id),
    ])
    allPermissions.value = permRes.data || []
    selectedPermIds.value = (rolePermRes.data || []).map((p) => String(p.id))
    permDialogShow.value = true
  } catch {
    message.error('加载权限失败')
  } finally {
    permLoading.value = false
  }
}

async function handlePermSave() {
  if (!editingRole.value) return
  try {
    await assignPermissions(editingRole.value.id, selectedPermIds.value.map(Number))
    message.success('权限分配成功')
    permDialogShow.value = false
  } catch {
    message.error('权限分配失败')
  }
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

function getPermGrouped(): Record<string, Permission[]> {
  const grouped: Record<string, Permission[]> = {}
  for (const p of allPermissions.value) {
    if (!grouped[p.group]) grouped[p.group] = []
    grouped[p.group].push(p)
  }
  return grouped
}
</script>

<template>
  <div class="page-container">
    <PageHeader title="角色管理" description="管理角色与权限">
      <n-button type="primary" @click="handleAdd">添加角色</n-button>
    </PageHeader>

    <DataTable
      ref="tableRef"
      :columns="columns"
      :fetch-api="fetchData"
      :search-fields="['name', 'code']"
      search-placeholder="搜索角色名称或编码"
    />

    <FormDialog
      ref="formDialogRef"
      :title="editingRole ? '编辑角色' : '添加角色'"
      :fields="formFields"
      @submit="handleFormSubmit"
    />

    <!-- Permission Dialog -->
    <n-modal
      v-model:show="permDialogShow"
      :title="`分配权限 - ${editingRole?.name}`"
      preset="card"
      :style="{ width: isMobile ? 'calc(100vw - 32px)' : '700px', maxWidth: 'calc(100vw - 32px)' }"
      :mask-closable="false"
    >
      <div v-if="permLoading" style="text-align: center; padding: 40px">
        加载中...
      </div>
      <div v-else>
        <n-checkbox-group v-model:value="selectedPermIds">
          <div v-for="(perms, group) in getPermGrouped()" :key="group" style="margin-bottom: 16px">
            <div style="font-weight: 600; margin-bottom: 8px; color: var(--text-color)">{{ group }}</div>
            <n-space>
              <n-checkbox
                v-for="perm in perms"
                :key="perm.id"
                :value="String(perm.id)"
                :label="perm.name"
              />
            </n-space>
          </div>
        </n-checkbox-group>
      </div>
      <template #footer>
        <n-space justify="end">
          <n-button @click="permDialogShow = false">取消</n-button>
          <n-button type="primary" @click="handlePermSave">保存</n-button>
        </n-space>
      </template>
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
