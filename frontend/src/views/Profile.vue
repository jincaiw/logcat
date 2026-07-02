<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { NForm, NFormItem, NInput, NButton, NAlert, useMessage } from 'naive-ui'
import { useI18n } from '@/i18n'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const message = useMessage()
const authStore = useAuthStore()

const profileForm = reactive({
  username: '',
  nickname: '',
  email: '',
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const profileLoading = ref(false)
const passwordLoading = ref(false)

onMounted(async () => {
  if (authStore.user) {
    profileForm.username = authStore.user.username || ''
    profileForm.nickname = authStore.user.nickname || ''
    profileForm.email = authStore.user.email || ''
  } else {
    await loadProfile()
  }
})

async function loadProfile() {
  try {
    const user = await authApi.getProfile()
    authStore.user = user
    profileForm.username = user.username || ''
    profileForm.nickname = user.nickname || ''
    profileForm.email = user.email || ''
  } catch (e) {
    console.error(e)
  }
}

async function handleSaveProfile() {
  if (!profileForm.nickname) {
    message.warning(t('profile.nicknameRequired'))
    return
  }
  profileLoading.value = true
  try {
    const updated = await authApi.updateProfile({
      nickname: profileForm.nickname,
      email: profileForm.email,
    })
    authStore.user = updated
    message.success(t('profile.profileUpdateSuccess'))
  } catch (e: any) {
    message.error(e.message || t('message.operationFailed'))
  } finally {
    profileLoading.value = false
  }
}

async function handleChangePassword() {
  if (!passwordForm.oldPassword || !passwordForm.newPassword || !passwordForm.confirmPassword) {
    message.warning(t('profile.allFieldsRequired'))
    return
  }
  if (passwordForm.newPassword.length < 6) {
    message.warning(t('profile.passwordTooShort'))
    return
  }
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    message.warning(t('profile.passwordNotMatch'))
    return
  }
  passwordLoading.value = true
  try {
    await authStore.changePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword,
    })
    message.success(t('profile.passwordChangeSuccess'))
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    await authStore.logout()
    window.location.hash = '#/login'
  } catch (e: any) {
    message.error(e.message || t('message.operationFailed'))
  } finally {
    passwordLoading.value = false
  }
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('profile.title') }}</h1>
        <p class="page-subtitle text-muted">{{ t('profile.subtitle') }}</p>
      </div>
    </div>

    <div class="profile-layout mt-4">
      <div class="section-card profile-section">
        <div class="card-header">
          <h3 class="card-title text-accent">{{ t('profile.title') }}</h3>
        </div>
        <NForm label-placement="left" :label-width="100">
          <NFormItem :label="t('profile.username')">
            <NInput v-model:value="profileForm.username" disabled />
          </NFormItem>
          <NFormItem :label="t('profile.nickname')">
            <NInput v-model:value="profileForm.nickname" :placeholder="t('profile.nickname')" />
          </NFormItem>
          <NFormItem label="Email">
            <NInput v-model:value="profileForm.email" placeholder="Email" />
          </NFormItem>
          <NFormItem>
            <NButton type="primary" :loading="profileLoading" @click="handleSaveProfile">
              {{ t('common.save') }}
            </NButton>
          </NFormItem>
        </NForm>
      </div>

      <div class="section-card profile-section">
        <div class="card-header">
          <h3 class="card-title text-accent">{{ t('profile.changePassword') }}</h3>
        </div>
        <NAlert type="info" :bordered="false" style="margin-bottom: 16px">
          {{ t('profile.subtitle') }}
        </NAlert>
        <NForm label-placement="left" :label-width="100">
          <NFormItem :label="t('profile.oldPassword')">
            <NInput v-model:value="passwordForm.oldPassword" type="password" show-password-on="click" :placeholder="t('profile.oldPasswordPlaceholder')" />
          </NFormItem>
          <NFormItem :label="t('profile.newPassword')">
            <NInput v-model:value="passwordForm.newPassword" type="password" show-password-on="click" :placeholder="t('profile.newPasswordPlaceholder')" />
          </NFormItem>
          <NFormItem :label="t('profile.confirmPassword')">
            <NInput v-model:value="passwordForm.confirmPassword" type="password" show-password-on="click" :placeholder="t('profile.confirmPasswordPlaceholder')" />
          </NFormItem>
          <NFormItem>
            <NButton type="primary" :loading="passwordLoading" @click="handleChangePassword">
              {{ t('profile.changePassword') }}
            </NButton>
          </NFormItem>
        </NForm>
      </div>
    </div>
  </div>
</template>
