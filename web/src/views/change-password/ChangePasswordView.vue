<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NAlert, NSpace } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { changePassword as apiChangePassword } from '@/api/auth'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const formData = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const rules: FormRules = {
  oldPassword: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule, value) => value === formData.newPassword,
      message: '两次输入的密码不一致',
      trigger: ['blur', 'change'],
    },
  ],
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    errorMsg.value = ''
    successMsg.value = ''
    loading.value = true

    await apiChangePassword(formData.oldPassword, formData.newPassword)
    authStore.mustChangePassword = false
    successMsg.value = '密码修改成功！'

    setTimeout(() => {
      router.push('/')
    }, 1500)
  } catch (err: any) {
    errorMsg.value = err?.message || '密码修改失败'
  } finally {
    loading.value = false
  }
}

function handleCancel() {
  if (!authStore.mustChangePassword) {
    router.push('/')
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-logo">
        <h1>GoLog</h1>
        <p v-if="authStore.mustChangePassword">首次登录，请修改密码</p>
        <p v-else>修改密码</p>
      </div>

      <n-alert
        v-if="successMsg"
        type="success"
        :bordered="false"
        style="margin-bottom: 16px"
      >
        {{ successMsg }}
      </n-alert>

      <n-alert
        v-if="errorMsg"
        type="error"
        :bordered="false"
        style="margin-bottom: 16px"
        closable
        @update:show="() => errorMsg = ''"
      >
        {{ errorMsg }}
      </n-alert>

      <n-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-placement="left"
        label-width="0"
        size="large"
      >
        <n-form-item path="oldPassword">
          <n-input
            v-model:value="formData.oldPassword"
            type="password"
            show-password-on="click"
            placeholder="原密码"
          />
        </n-form-item>

        <n-form-item path="newPassword">
          <n-input
            v-model:value="formData.newPassword"
            type="password"
            show-password-on="click"
            placeholder="新密码"
          />
        </n-form-item>

        <n-form-item path="confirmPassword">
          <n-input
            v-model:value="formData.confirmPassword"
            type="password"
            show-password-on="click"
            placeholder="确认新密码"
          />
        </n-form-item>

        <n-form-item>
          <n-space style="width: 100%" :justify="'end'">
            <n-button
              v-if="!authStore.mustChangePassword"
              @click="handleCancel"
            >
              取消
            </n-button>
            <n-button
              type="primary"
              :loading="loading"
              :disabled="!!successMsg"
              @click="handleSubmit"
            >
              确认修改
            </n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </div>
  </div>
</template>