import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api'
import type { User } from '@/types'

const TOKEN_KEY = 'logcat-token'

type AuthUser = User

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY))
  const user = ref<AuthUser | null>(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  function setToken(t: string) {
    token.value = t
    localStorage.setItem(TOKEN_KEY, t)
  }

  function clearToken() {
    token.value = null
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
  }

  async function login(username: string, password: string) {
    loading.value = true
    try {
      const result = await authApi.login({ username, password })
      setToken(result.token)
      if (result.user) {
        user.value = result.user
      }
      return result
    } finally {
      loading.value = false
    }
  }

  async function fetchProfile() {
    if (!token.value) return null
    try {
      const profile = await authApi.getProfile()
      user.value = profile
      return profile
    } catch (e) {
      clearToken()
      throw e
    }
  }

  async function updateProfile(data: Partial<AuthUser>) {
    const result = await authApi.updateProfile(data)
    user.value = result
    return result
  }

  async function changePassword(data: { oldPassword: string; newPassword: string }) {
    await authApi.changePassword(data)
  }

  async function logout() {
    try {
      if (token.value) {
        await authApi.logout()
      }
    } catch {
      // ignore
    } finally {
      clearToken()
    }
  }

  return {
    token,
    user,
    loading,
    isAuthenticated,
    login,
    logout,
    fetchProfile,
    updateProfile,
    changePassword,
    setToken,
    clearToken,
  }
})
