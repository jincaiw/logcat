<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NAlert, NIcon } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { LockClosedOutline } from '@vicons/ionicons5'
import { initAdmin } from '@/api/auth'

const router = useRouter()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const formData = reactive({
  password: '',
  confirmPassword: '',
})

const rules: FormRules = {
  password: [
    { required: true, message: '请输入管理员密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (_rule, value) => value === formData.password,
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

    await initAdmin(formData.password)
    successMsg.value = '管理员密码设置成功！即将跳转到登录页面...'

    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (err: any) {
    errorMsg.value = err?.message || '初始化失败'
  } finally {
    loading.value = false
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
        <p>初始化管理员密码</p>
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
        <n-form-item path="password">
          <n-input
            v-model:value="formData.password"
            type="password"
            show-password-on="click"
            placeholder="管理员密码"
          >
            <template #prefix>
              <n-icon :component="LockClosedOutline" style="color: var(--text-color-tertiary)" />
            </template>
          </n-input>
        </n-form-item>

        <n-form-item path="confirmPassword">
          <n-input
            v-model:value="formData.confirmPassword"
            type="password"
            show-password-on="click"
            placeholder="确认密码"
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
            :disabled="!!successMsg"
            @click="handleSubmit"
            style="height: 42px; font-size: 15px; border-radius: 8px"
          >
            设置密码
          </n-button>
        </n-form-item>
      </n-form>
    </div>
  </div>
</template>
