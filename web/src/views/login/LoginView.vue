<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NSpace, NAlert } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { checkInitStatus } from '@/api/auth'

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

    // Check if system needs initialization
    const initRes = await checkInitStatus()
    if (initRes.data && !initRes.data.initialized) {
      router.push('/init')
      return
    }

    await authStore.login(formData.username, formData.password)

    if (authStore.mustChangePassword) {
      router.push('/change-password')
    } else {
      router.push('/')
    }
  } catch (err: any) {
    errorMsg.value = err?.message || err?.response?.data?.message || '登录失败，请检查用户名和密码'
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
        <h1>GoLog</h1>
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
            @keyup="handleKeyup"
          >
            <template #prefix>
              <span style="color: var(--text-color-secondary)">用户</span>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item path="password">
          <n-input
            v-model:value="formData.password"
            type="password"
            show-password-on="click"
            placeholder="密码"
            @keyup="handleKeyup"
          >
            <template #prefix>
              <span style="color: var(--text-color-secondary)">密码</span>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item>
          <n-button
            type="primary"
            block
            :loading="loading"
            @click="handleLogin"
          >
            登 录
          </n-button>
        </n-form-item>
      </n-form>
    </div>
  </div>
</template>