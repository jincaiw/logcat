<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { NInput, NButton, NForm, NFormItem, NAlert, useMessage } from 'naive-ui'
import { useI18n } from '@/i18n'
import { useAuthStore } from '@/stores/auth'
import { useRouter, useRoute } from 'vue-router'

const { t, locale, setLocale } = useI18n()
const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const message = useMessage()

const loading = ref(false)
const errorMsg = ref('')

const form = reactive({
  username: 'admin',
  password: 'admin',
})

onMounted(() => {
  if (authStore.isAuthenticated) {
    router.push('/dashboard')
  }
})

async function handleLogin() {
  if (!form.username || !form.password) {
    errorMsg.value = t('auth.loginFailed')
    return
  }
  loading.value = true
  errorMsg.value = ''
  try {
    await authStore.login(form.username, form.password)
    message.success(t('auth.loginSuccess'))
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch (e: any) {
    errorMsg.value = e.message || t('auth.loginFailed')
  } finally {
    loading.value = false
  }
}

function toggleLang() {
  setLocale(locale.value === 'zh-CN' ? 'en-US' : 'zh-CN')
}
</script>

<template>
  <div class="login-page">
    <button class="btn-lang-float" @click="toggleLang">
      {{ locale === 'zh-CN' ? 'EN' : '中' }}
    </button>

    <div class="login-card">
      <div class="login-logo">L</div>
      <h1 class="login-title">logcat</h1>
      <p class="login-subtitle">{{ t('auth.subtitle') }}</p>

      <NAlert v-if="errorMsg" type="error" :bordered="false" style="margin-bottom: 16px">
        {{ errorMsg }}
      </NAlert>

      <NForm @submit.prevent="handleLogin">
        <NFormItem>
          <NInput
            v-model:value="form.username"
            size="large"
            :placeholder="t('auth.usernamePlaceholder')"
            @keyup.enter="handleLogin"
          />
        </NFormItem>
        <NFormItem>
          <NInput
            v-model:value="form.password"
            type="password"
            size="large"
            show-password-on="click"
            :placeholder="t('auth.passwordPlaceholder')"
            @keyup.enter="handleLogin"
          />
        </NFormItem>
        <NButton
          type="primary"
          size="large"
          block
          :loading="loading"
          attr-type="submit"
          style="margin-top: 8px"
        >
          {{ t('auth.login') }}
        </NButton>
      </NForm>
    </div>
  </div>
</template>

<style scoped>
.btn-lang-float {
  position: fixed;
  top: 20px;
  right: 20px;
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  color: var(--text-secondary);
  padding: 6px 14px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s;
  z-index: 10;
}
.btn-lang-float:hover {
  background: var(--accent-subtle-hover);
  color: var(--text-primary);
}
</style>
