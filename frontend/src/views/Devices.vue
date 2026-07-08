<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSwitch, NTag, NPopconfirm, NEmpty, NSpace, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { API } from '@/api'
import type { Device } from '@/types'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const message = useMessage()

const loading = ref(false)
const devices = ref<Device[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref(t('device.addDevice'))
const formData = ref<Partial<Device>>({
  name: '',
  ipAddress: '',
  groupId: 0,
  description: '',
  groupName: 'default',
  isActive: true,
})

onMounted(() => {
  loadDevices()
})

async function loadDevices() {
  loading.value = true
  try {
    devices.value = await API.GetDevices()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  dialogTitle.value = t('device.addDevice')
  formData.value = {
    name: '',
    ipAddress: '',
    groupId: 0,
    description: '',
    groupName: 'default',
    isActive: true,
  }
  dialogVisible.value = true
}

function handleEdit(row: Device) {
  dialogTitle.value = t('device.editDevice')
  formData.value = { ...row }
  dialogVisible.value = true
}

async function handleDelete(row: Device) {
  try {
    await API.DeleteDevice(row.id!)
    message.success(t('message.deleteSuccess'))
    loadDevices()
  } catch (e: any) {
    message.error(t('message.deleteFailed'))
  }
}

async function handleSubmit() {
  if (!formData.value.name || !formData.value.ipAddress) {
    message.warning(t('common.requiredFields'))
    return
  }
  try {
    if (formData.value.id) {
      await API.UpdateDevice({ ...formData.value, id: formData.value.id! } as any)
      message.success(t('message.updateSuccess'))
    } else {
      await API.AddDevice(formData.value as any)
      message.success(t('message.addSuccess'))
    }
    dialogVisible.value = false
    loadDevices()
  } catch (e) {
    message.error(t('message.operationFailed'))
  }
}

const columns: DataTableColumns<Device> = [
  { title: t('common.id'), key: 'id', width: 80 },
  { title: t('device.name'), key: 'name', width: 150 },
  { title: t('device.ipAddress'), key: 'ipAddress', width: 140 },
  {
    title: t('device.group'),
    key: 'groupName',
    width: 100,
    render(row) {
      return h(NTag, { size: 'small' }, { default: () => row.groupName || 'default' })
    },
  },
  { title: t('common.description'), key: 'description', ellipsis: { tooltip: true } },
  {
    title: t('common.status'),
    key: 'isActive',
    width: 80,
    render(row) {
      return h(NTag, { type: row.isActive ? 'success' : 'error', size: 'small' }, {
        default: () => row.isActive ? t('common.enable') : t('common.disable'),
      })
    },
  },
  {
    title: t('common.action'),
    key: 'actions',
    width: 150,
    render(row) {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { text: true, type: 'primary', size: 'small', onClick: () => handleEdit(row) }, { default: () => t('common.edit') }),
          h(NPopconfirm, { onPositiveClick: () => handleDelete(row) }, {
            trigger: () => h(NButton, { text: true, type: 'error', size: 'small' }, { default: () => t('common.delete') }),
            default: () => t('device.deleteConfirm'),
          }),
        ],
      })
    },
  },
]
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('device.list') }}</h1>
        <p class="page-subtitle text-muted">{{ t('device.subtitle') }}</p>
      </div>
      <NSpace class="page-actions">
        <NButton type="primary" @click="handleAdd">
          {{ t('device.addDevice') }}
        </NButton>
      </NSpace>
    </div>

    <div class="table-card mt-4">
      <NDataTable
        :columns="columns"
        :data="devices"
        :loading="loading"
        :bordered="false"
        striped
      />
      <div v-if="!devices.length && !loading" class="empty-state">
        <NEmpty :description="t('common.noDataDesc')" />
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <NModal
      v-model:show="dialogVisible"
      :title="dialogTitle"
      preset="card"
      class="modal-520"
      :bordered="true"
    >
      <NForm :model="formData" label-placement="left" :label-width="100">
        <NFormItem :label="t('device.name')" required>
          <NInput v-model:value="formData.name" :placeholder="t('device.pleaseInputName')" />
        </NFormItem>
        <NFormItem :label="t('device.ipAddress')" required>
          <NInput v-model:value="formData.ipAddress" :placeholder="t('device.pleaseInputIp')" />
        </NFormItem>
        <NFormItem :label="t('device.group')">
          <NInput v-model:value="formData.groupName" :placeholder="t('device.defaultGroupPlaceholder')" />
        </NFormItem>
        <NFormItem :label="t('common.description')">
          <NInput v-model:value="formData.description" type="textarea" :rows="3" :placeholder="t('common.pleaseInputDescription')" />
        </NFormItem>
        <NFormItem :label="t('common.status')">
          <NSwitch v-model:value="formData.isActive" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="dialogVisible = false">{{ t('common.cancel') }}</NButton>
          <NButton type="primary" @click="handleSubmit">{{ t('common.confirmButtonText') }}</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>
