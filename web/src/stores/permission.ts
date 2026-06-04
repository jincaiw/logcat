import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Permission, Role } from '@/types'
import { getAllPermissions, getRolePermissions, getAllRoles } from '@/api/roles'
import { getCurrentUser } from '@/api/users'

export interface MenuItem {
  label: string
  key: string
  icon?: string
  permission?: string
  children?: MenuItem[]
}

export const usePermissionStore = defineStore('permission', () => {
  const permissions = ref<Permission[]>([])
  const roles = ref<Role[]>([])
  const loaded = ref(false)
  const userPermissionCodes = ref<Set<string>>(new Set())

  const permissionCodes = computed(() => userPermissionCodes.value)

  function hasPermission(code: string): boolean {
    if (!code) return true
    return userPermissionCodes.value.has(code)
  }

  function setUserPermissions(codes: string[]) {
    userPermissionCodes.value = new Set(codes)
  }

  const menuItems: MenuItem[] = [
    { label: '仪表盘', key: '/', icon: 'grid-outline' },
    {
      label: '日志中心',
      key: 'logs-center',
      icon: 'reader-outline',
      children: [
        { label: '日志查询', key: '/logs', permission: 'logs:list' },
        { label: '告警记录', key: '/alerts', permission: 'alerts:list' },
        { label: '告警处置', key: '/alerts/disposition', permission: 'alerts:disposition:list' },
        { label: '聚合告警', key: '/aggregated-alerts', permission: 'aggregated-alerts:list' },
        { label: '高频IP', key: '/high-freq-ips', permission: 'high-freq-ips:list' },
      ],
    },
    {
      label: '设备与模板',
      key: 'devices-templates',
      icon: 'hardware-chip-outline',
      children: [
        { label: '设备管理', key: '/devices', permission: 'devices:list' },
        { label: '设备分组', key: '/device-groups', permission: 'device-groups:list' },
        { label: '设备模板', key: '/device-templates', permission: 'device-templates:list' },
        { label: '解析模板', key: '/parse-templates', permission: 'parse-templates:list' },
        { label: '字段映射', key: '/field-mappings', permission: 'field-mappings:list' },
        { label: '过滤策略', key: '/filter-policies', permission: 'filter-policies:list' },
        { label: '输出模板', key: '/output-templates', permission: 'output-templates:list' },
      ],
    },
    {
      label: '告警与推送',
      key: 'alert-push',
      icon: 'notifications-outline',
      children: [
        { label: '告警规则', key: '/alert-rules', permission: 'alert-rules:list' },
        { label: '推送配置', key: '/push-configs', permission: 'push-configs:list' },
        { label: '脱敏规则', key: '/desensitize', permission: 'desensitize-rules:list' },
      ],
    },
    {
      label: '系统管理',
      key: 'system',
      icon: 'settings-outline',
      children: [
        { label: '用户管理', key: '/users', permission: 'users:list' },
        { label: '角色管理', key: '/roles', permission: 'roles:list' },
        { label: '系统配置', key: '/system/config', permission: 'system:config:read' },
        { label: '服务状态', key: '/system/status', permission: 'system:status' },
        { label: '指标监控', key: '/system/metrics', permission: 'system:status' },
        { label: '数据统计', key: '/stats', permission: 'stats:fields' },
        { label: '导入导出', key: '/import-export', permission: 'export:config' },
        { label: '审计日志', key: '/audit-logs', permission: 'audit-logs:list' },
      ],
    },
  ]

  function filterVisibleMenus(menus: MenuItem[]): MenuItem[] {
    return menus
      .filter((m) => {
        if (m.permission && !hasPermission(m.permission)) return false
        if (m.children) {
          const filteredChildren = filterVisibleMenus(m.children)
          return filteredChildren.length > 0
        }
        return true
      })
      .map((m) => ({
        ...m,
        children: m.children ? filterVisibleMenus(m.children) : undefined,
      }))
  }

  const visibleMenus = computed(() => filterVisibleMenus(JSON.parse(JSON.stringify(menuItems))))

  async function fetchPermissions() {
    try {
      const roleRes = await getAllRoles()
      roles.value = roleRes.data || []
      const permRes = await getAllPermissions()
      permissions.value = permRes.data || []
      try {
        const userRes = await getCurrentUser()
        if (userRes.data?.permissions) {
          setUserPermissions(userRes.data.permissions)
        }
      } catch {
        // ignore - user may not be logged in
      }
      loaded.value = true
    } catch {
      loaded.value = true
    }
  }

  function reset() {
    permissions.value = []
    roles.value = []
    loaded.value = false
    userPermissionCodes.value = new Set()
  }

  return {
    permissions,
    roles,
    loaded,
    permissionCodes,
    hasPermission,
    setUserPermissions,
    menuItems,
    visibleMenus,
    fetchPermissions,
    reset,
  }
})
