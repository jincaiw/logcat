<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NAlert } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
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
        <h1>GoLog</h1>
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
          />
        </n-form-item>

        <n-form-item path="confirmPassword">
          <n-input
            v-model:value="formData.confirmPassword"
            type="password"
            show-password-on="click"
            placeholder="确认密码"
          />
        </n-form-item>

        <n-form-item>
          <n-button
            type="primary"
            block
            :loading="loading"
            :disabled="!!successMsg"
            @click="handleSubmit"
          >
            设置密码
          </n-button>
        </n-form-item>
      </n-form>
    </div>
  </div>
</template>