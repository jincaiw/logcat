import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

export function useAuth() {
  const authStore = useAuthStore()
  const router = useRouter()

  const user = computed(() => authStore.user)
  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const isAdmin = computed(() => authStore.isAdmin)

  async function login(username: string, password: string) {
    await authStore.login(username, password)
    if (authStore.mustChangePassword) {
      router.push('/change-password')
    } else {
      router.push('/')
    }
  }

  async function logout() {
    await authStore.logout()
    router.push('/login')
  }

  async function changePassword(oldPwd: string, newPwd: string) {
    await authStore.changePassword(oldPwd, newPwd)
    router.push('/')
  }

  return {
    user,
    isAuthenticated,
    isAdmin,
    login,
    logout,
    changePassword,
  }
}