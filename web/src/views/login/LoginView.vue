<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NAlert, NIcon } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { PersonOutline, LockClosedOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const errorMsg = ref('')

const formData = reactive({
  username: '',
  password: '',
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleLogin() {
  try {
    await formRef.value?.validate()
    errorMsg.value = ''
    loading.value = true

    await authStore.login(formData.username, formData.password)

    if (authStore.mustChangePassword) {
      router.push('/change-password')
    } else {
      router.push('/')
    }
  } catch (err: any) {
    if (err?.message === 'Network Error' || err?.code === 'ECONNABORTED' || err?.code === 'ERR_NETWORK') {
      errorMsg.value = '无法连接到服务器，请检查网络连接或确认服务是否正常运行'
    } else {
      errorMsg.value = err?.message || '登录失败，请检查用户名和密码'
    }
  } finally {
    loading.value = false
  }
}

function handleKeyup(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    handleLogin()
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-logo">
        <div class="login-logo-icon">
          <svg width="36" height="36" viewBox="0 0 36 36" fill="none">
            <rect width="36" height="36" rx="10" fill="var(--primary-color)" />
            <path d="M10 12h4l4 8 4-8h4v14h-4v-8l-4 8-4-8v8h-4V12z" fill="white" />
          </svg>
        </div>
        <h1>logcat</h1>
        <p>日志采集与管理平台</p>
      </div>

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
        <n-form-item path="username">
          <n-input
            v-model:value="formData.username"
            placeholder="用户名"
            autocomplete="username"
            @keyup="handleKeyup"
          >
            <template #prefix>
              <n-icon :component="PersonOutline" style="color: var(--text-color-tertiary)" />
            </template>
          </n-input>
        </n-form-item>

        <n-form-item path="password">
          <n-input
            v-model:value="formData.password"
            type="password"
            show-password-on="click"
            placeholder="密码"
            autocomplete="current-password"
            @keyup="handleKeyup"
          >
            <template #prefix>
              <n-icon :component="LockClosedOutline" style="color: var(--text-color-tertiary)" />
            </template>
          </n-input>
        </n-form-item>

        <n-form-item>
          <n-button
            type="primary"
            block
            :loading="loading"
            @click="handleLogin"
            style="height: 42px; font-size: 15px; border-radius: 8px"
          >
            登 录
          </n-button>
        </n-form-item>
      </n-form>
    </div>
  </div>
</template>
