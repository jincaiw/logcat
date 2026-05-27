<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import {
  NModal, NForm, NFormItem, NInput, NSelect, NInputNumber,
  NSwitch, NButton, NSpace, useMessage,
} from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'

export interface FieldConfig {
  key: string
  label: string
  type: 'text' | 'number' | 'select' | 'switch' | 'textarea' | 'password' | 'code'
  required?: boolean
  placeholder?: string
  options?: { label: string; value: any }[]
  disabled?: boolean
  defaultValue?: any
  span?: number
  rules?: any[]
  min?: number
  max?: number
}

const props = defineProps<{
  title: string
  fields: FieldConfig[]
  initialData?: Record<string, any>
  loading?: boolean
  width?: number
  labelWidth?: number
}>()

const emit = defineEmits<{
  (e: 'submit', data: Record<string, any>): void
  (e: 'cancel'): void
  (e: 'update:show', value: boolean): void
}>()

const show = ref(false)
const formRef = ref<FormInst | null>(null)
const formData = ref<Record<string, any>>({})
const message = useMessage()

const rules = computed<FormRules>(() => {
  const r: FormRules = {}
  props.fields.forEach((field) => {
    if (field.required) {
      r[field.key] = [
        ...(field.rules || []),
        { required: true, message: `请输入${field.label}`, trigger: ['blur', 'change'] },
      ]
    } else if (field.rules) {
      r[field.key] = field.rules
    }
  })
  return r
})

function open(data?: Record<string, any>) {
  show.value = true
  initForm(data || props.initialData || {})
}

function close() {
  show.value = false
  emit('cancel')
}

function initForm(data: Record<string, any>) {
  const d: Record<string, any> = {}
  props.fields.forEach((field) => {
    d[field.key] = data[field.key] !== undefined ? data[field.key] : (field.defaultValue ?? '')
  })
  formData.value = d
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    emit('submit', { ...formData.value })
  } catch {
    message.warning('请检查表单填写')
  }
}

watch(() => props.initialData, (val) => {
  if (val && show.value) {
    initForm(val)
  }
}, { deep: true })

defineExpose({ open, close })
</script>

<template>
  <n-modal
    v-model:show="show"
    :title="title"
    preset="card"
    :style="{ width: `${width || 600}px` }"
    :mask-closable="false"
    @close="close"
  >
    <n-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      :label-width="labelWidth || 100"
      label-placement="left"
    >
      <n-form-item
        v-for="field in fields"
        :key="field.key"
        :label="field.label"
        :path="field.key"
      >
        <n-input
          v-if="field.type === 'text'"
          v-model:value="formData[field.key]"
          :placeholder="field.placeholder || `请输入${field.label}`"
          :disabled="field.disabled"
        />
        <n-input
          v-else-if="field.type === 'password'"
          v-model:value="formData[field.key]"
          type="password"
          show-password-on="click"
          :placeholder="field.placeholder || `请输入${field.label}`"
        />
        <n-input
          v-else-if="field.type === 'textarea'"
          v-model:value="formData[field.key]"
          type="textarea"
          :placeholder="field.placeholder || `请输入${field.label}`"
          :autosize="{ minRows: 3, maxRows: 8 }"
        />
        <n-input
          v-else-if="field.type === 'code'"
          v-model:value="formData[field.key]"
          type="textarea"
          :placeholder="field.placeholder || `请输入${field.label}`"
          :autosize="{ minRows: 5, maxRows: 15 }"
          style="font-family: monospace"
        />
        <n-input-number
          v-else-if="field.type === 'number'"
          v-model:value="formData[field.key]"
          :placeholder="field.placeholder || `请输入${field.label}`"
          :min="field.min"
          :max="field.max"
        />
        <n-select
          v-else-if="field.type === 'select'"
          v-model:value="formData[field.key]"
          :options="field.options || []"
          :placeholder="field.placeholder || `请选择${field.label}`"
        />
        <n-switch
          v-else-if="field.type === 'switch'"
          v-model:value="formData[field.key]"
        />
      </n-form-item>
    </n-form>

    <template #footer>
      <n-space justify="end">
        <n-button @click="close">取消</n-button>
        <n-button type="primary" :loading="loading" @click="handleSubmit">
          确定
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>