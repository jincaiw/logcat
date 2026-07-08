<script setup lang="ts">
import { computed, ref, reactive, onMounted } from 'vue'
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
  password: '',
})

const highlights = computed(() => [
  t('auth.loginPoint1'),
  t('auth.loginPoint2'),
  t('auth.loginPoint3'),
])

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
    <button class="btn-lang-float" :aria-label="t('common.language')" :title="t('common.language')" @click="toggleLang">
      {{ locale === 'zh-CN' ? 'EN' : '中' }}
    </button>

    <div class="login-shell">
      <section class="login-marketing">
        <div class="login-brand">
          <div class="login-logo">L</div>
          <div class="login-brand-copy">
            <p class="login-kicker">logcat</p>
            <h1 class="login-title">{{ t('auth.subtitle') }}</h1>
          </div>
        </div>

        <p class="login-pitch">{{ t('auth.loginPitch') }}</p>

        <div class="login-points">
          <div v-for="(item, idx) in highlights" :key="idx" class="login-point">
            <span class="login-point-mark"></span>
            <p>{{ item }}</p>
          </div>
        </div>

        <div class="login-footer-note">
          {{ t('auth.loginMeta') }}
        </div>
      </section>

      <section class="login-card">
        <div class="login-card-header">
          <h2>{{ t('auth.login') }}</h2>
          <p>{{ t('auth.consoleIntro') }}</p>
        </div>

        <NAlert v-if="errorMsg" type="error" :bordered="false" class="login-alert">
          {{ errorMsg }}
        </NAlert>

        <NForm @submit.prevent="handleLogin" class="login-form">
          <NFormItem :label="t('auth.username')">
            <NInput
              v-model:value="form.username"
              size="large"
              :placeholder="t('auth.usernamePlaceholder')"
              autocomplete="username"
              spellcheck="false"
              @keyup.enter="handleLogin"
            />
          </NFormItem>
          <NFormItem :label="t('auth.password')">
            <NInput
              v-model:value="form.password"
              type="password"
              size="large"
              show-password-on="click"
              :placeholder="t('auth.passwordPlaceholder')"
              autocomplete="current-password"
              spellcheck="false"
              @keyup.enter="handleLogin"
            />
          </NFormItem>
          <NButton
            type="primary"
            size="large"
            block
            :loading="loading"
            attr-type="submit"
            class="login-submit"
          >
            {{ t('auth.login') }}
          </NButton>
        </NForm>
      </section>
    </div>
  </div>
</template>

