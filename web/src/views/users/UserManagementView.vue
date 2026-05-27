<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NTag, NSpace, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createUser, updateUser, deleteUser, getUsers, resetPassword, unlockUser, forcePasswordChange, assignRoles, getUserRoles } from '@/api/users'
import { getAllRoles } from '@/api/roles'
import type { User } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import FormDialog, { type FieldConfig } from '@/components/common/FormDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import PageHeader from '@/components/common/PageHeader.vue'

const message = useMessage()

const tableRef = ref<InstanceType<typeof DataTable> | null>(null)
const formDialogRef = ref<InstanceType<typeof FormDialog> | null>(null)
const roleDialogRef = ref<InstanceType<typeof FormDialog> | null>(null)
const confirmDialogShow = ref(false)
const confirmTitle = ref('')
const confirmContent = ref('')
const confirmAction = ref<() => Promise<void>>(() => Promise.resolve())
const confirmLoading = ref(false)

const editingUser = ref<User | null>(null)
const formFields: FieldConfig[] = [
  { key: 'username', label: '用户名', type: 'text', required: true },
  { key: 'displayName', label: '显示名称', type: 'text', required: true },
  { key: 'email', label: '邮箱', type: 'text', placeholder: 'example@mail.com' },
  { key: 'phone', label: '手机号', type: 'text' },
  { key: 'status', label: '状态', type: 'select', options: [{ label: '启用', value: 1 }, { label: '禁用', value: 0 }], defaultValue: 1 },
]

const roleFields: FieldConfig[] = [
  { key: 'roleIds', label: '角色', type: 'select', placeholder: '请选择角色' },
]

const columns: DataTableColumns<User> = [
  { title: '用户名', key: 'username', sorter: true },
  { title: '显示名称', key: 'displayName' },
  { title: '邮箱', key: 'email' },
  {
    title: '状态', key: 'status',
    render(row) {
      return h(NTag, { type: row.status === 1 ? 'success' : 'default', size: 'small', bordered: false }, { default: () => row.status === 1 ? '启用' : '禁用' })
    },
  },
  { title: '最后登录', key: 'lastLoginAt', width: 160 },
  {
    title: '操作', key: 'actions', width: 320,
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'small', type: 'primary', ghost: true, onClick: () => handleEdit(row) }, { default: () => '编辑' }),
          h(NButton, { size: 'small', onClick: () => handleAssignRoles(row) }, { default: () => '分配角色' }),
          h(NButton, { size: 'small', onClick: () => handleResetPassword(row) }, { default: () => '重置密码' }),
          h(NButton, { size: 'small', onClick: () => handleUnlock(row) }, { default: () => '解锁' }),
          h(NButton, { size: 'small', type: 'error', ghost: true, onClick: () => handleDelete(row) }, { default: () => '删除' }),
        ],
      })
    },
  },
]

async function fetchData(params: any) {
  const res = await getUsers(params)
  return res
}

function handleAdd() {
  editingUser.value = null
  formDialogRef.value?.open()
}

function handleEdit(row: User) {
  editingUser.value = row
  formDialogRef.value?.open({
    username: row.username,
    displayName: row.displayName,
    email: row.email,
    phone: row.phone,
    status: row.status,
  })
}

async function handleFormSubmit(data: Record<string, any>) {
  try {
    if (editingUser.value) {
      await updateUser(editingUser.value.id, data)
      message.success('用户更新成功')
    } else {
      await createUser(data)
      message.success('用户创建成功')
    }
    formDialogRef.value?.close()
    tableRef.value?.loadData()
  } catch (err: any) {
    message.error(err?.message || '操作失败')
  }
}

function handleDelete(row: User) {
  confirmTitle.value = '删除用户'
  confirmContent.value = `确定要删除用户 "${row.username}" 吗？`
  confirmAction.value = async () => {
    await deleteUser(row.id)
    message.success('删除成功')
    tableRef.value?.loadData()
  }
  confirmDialogShow.value = true
}

function handleResetPassword(row: User) {
  confirmTitle.value = '重置密码'
  confirmContent.value = `确定要将用户 "${row.username}" 的密码重置为默认密码吗？`
  confirmAction.value = async () => {
    await resetPassword(row.id, 'default123')
    message.success('密码已重置')
  }
  confirmDialogShow.value = true
}

function handleUnlock(row: User) {
  confirmTitle.value = '解锁用户'
  confirmContent.value = `确定要解锁用户 "${row.username}" 吗？`
  confirmAction.value = async () => {
    await unlockUser(row.id)
    message.success('用户已解锁')
  }
  confirmDialogShow.value = true
}

async function handleAssignRoles(row: User) {
  try {
    const [roleRes, userRoleRes] = await Promise.all([getAllRoles(), getUserRoles(row.id)])
    const roleOptions = (roleRes.data || []).map((r) => ({ label: r.name, value: r.id }))
    const assignedRoleIds = userRoleRes.data?.roles?.map((r) => r.id) || []

    roleFields[0].options = roleOptions
    editingUser.value = row
    roleDialogRef.value?.open({ roleIds: assignedRoleIds })
  } catch {
    message.error('获取角色信息失败')
  }
}

async function handleRoleSubmit(data: Record<string, any>) {
  if (!editingUser.value) return
  try {
    const roleIds = Array.isArray(data.roleIds) ? data.roleIds : [data.roleIds].filter(Boolean)
    await assignRoles(editingUser.value.id, roleIds)
    message.success('角色分配成功')
    roleDialogRef.value?.close()
  } catch (err: any) {
    message.error(err?.message || '分配失败')
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
</script>

<template>
  <div class="page-container">
    <PageHeader title="用户管理" description="管理系统用户">
      <n-button type="primary" @click="handleAdd">添加用户</n-button>
    </PageHeader>

    <DataTable
      ref="tableRef"
      :columns="columns"
      :fetch-api="fetchData"
      :search-fields="['username', 'displayName']"
      search-placeholder="搜索用户名或显示名称"
    />

    <FormDialog
      ref="formDialogRef"
      :title="editingUser ? '编辑用户' : '添加用户'"
      :fields="formFields"
      @submit="handleFormSubmit"
    />

    <FormDialog
      ref="roleDialogRef"
      title="分配角色"
      :fields="roleFields"
      @submit="handleRoleSubmit"
    />

    <ConfirmDialog
      v-model:show="confirmDialogShow"
      :title="confirmTitle"
      :content="confirmContent"
      :loading="confirmLoading"
      @confirm="handleConfirm"
    />
  </div>
</template>