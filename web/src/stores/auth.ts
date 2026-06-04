import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import * as authApi from '@/api/auth'
import { usePermissionStore } from './permission'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isAuthenticated = ref(false)
  const isInitialized = ref(false)
  const mustChangePassword = ref(false)

  const username = computed(() => user.value?.username || '')
  const isAdmin = computed(() => {
    if (!user.value?.roles) return false
    return user.value.roles.includes('admin')
  })

  async function fetchCurrentUser() {
    try {
      const res = await authApi.getCurrentUser()
      user.value = res.data
      isAuthenticated.value = true
      mustChangePassword.value = res.data.mustChangePassword || false
      isInitialized.value = true

      const permissionStore = usePermissionStore()
      permissionStore.setUserPermissions(res.data.permissions || [])
    } catch {
      isAuthenticated.value = false
      user.value = null
      mustChangePassword.value = false
      isInitialized.value = true
    }
  }

  async function login(uname: string, password: string) {
    const res = await authApi.login(uname, password)
    user.value = res.data.user
    isAuthenticated.value = true
    isInitialized.value = true
    mustChangePassword.value = res.data.user.mustChangePassword || false

    const permissionStore = usePermissionStore()
    permissionStore.setUserPermissions(res.data.user.permissions || [])

    return res.data.user
  }

  async function logout() {
    try {
      await authApi.logout()
    } finally {
      user.value = null
      isAuthenticated.value = false
      mustChangePassword.value = false

      const permissionStore = usePermissionStore()
      permissionStore.reset()
    }
  }

  async function changePassword(oldPwd: string, newPwd: string) {
    await authApi.changePassword(oldPwd, newPwd)
    mustChangePassword.value = false
  }

  function setUser(u: User) {
    user.value = u
    isAuthenticated.value = true

    const permissionStore = usePermissionStore()
    permissionStore.setUserPermissions(u.permissions || [])
  }

  return {
    user,
    isAuthenticated,
    isInitialized,
    mustChangePassword,
    username,
    isAdmin,
    fetchCurrentUser,
    login,
    logout,
    changePassword,
    setUser,
  }
})
