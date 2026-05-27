import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Permission, Role } from '@/types'
import { getAllPermissions, getRolePermissions, getAllRoles } from '@/api/roles'

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

  const permissionCodes = computed(() => new Set(permissions.value.map((p) => p.code)))

  function hasPermission(code: string): boolean {
    if (!code) return true
    return permissionCodes.value.has(code)
  }

  const menuItems: MenuItem[] = [
    { label: '仪表盘', key: '/', icon: 'grid-outline' },
    {
      label: '系统管理',
      key: 'system',
      icon: 'settings-outline',
      children: [
        { label: '服务状态', key: '/system/status', permission: 'system:status' },
        { label: '系统配置', key: '/system/config', permission: 'system:config' },
      ],
    },
    {
      label: '用户管理',
      key: 'users',
      icon: 'people-outline',
      children: [
        { label: '用户列表', key: '/users', permission: 'user:view' },
        { label: '角色管理', key: '/roles', permission: 'role:view' },
      ],
    },
    {
      label: '设备管理',
      key: 'devices',
      icon: 'hardware-chip-outline',
      children: [
        { label: '设备列表', key: '/devices', permission: 'device:view' },
        { label: '设备分组', key: '/device-groups', permission: 'deviceGroup:view' },
      ],
    },
    {
      label: '模板配置',
      key: 'templates',
      icon: 'document-text-outline',
      children: [
        { label: '设备模板', key: '/device-templates', permission: 'deviceTemplate:view' },
        { label: '字段映射', key: '/field-mappings', permission: 'fieldMapping:view' },
        { label: '解析模板', key: '/parse-templates', permission: 'parseTemplate:view' },
        { label: '过滤策略', key: '/filter-policies', permission: 'filterPolicy:view' },
        { label: '输出模板', key: '/output-templates', permission: 'outputTemplate:view' },
      ],
    },
    {
      label: '推送与告警',
      key: 'push-alert',
      icon: 'notifications-outline',
      children: [
        { label: '推送配置', key: '/push-configs', permission: 'pushConfig:view' },
        { label: '告警规则', key: '/alert-rules', permission: 'alertRule:view' },
      ],
    },
    {
      label: '日志管理',
      key: 'logs',
      icon: 'reader-outline',
      children: [
        { label: '日志查询', key: '/logs', permission: 'log:view' },
      ],
    },
    {
      label: '告警管理',
      key: 'alerts',
      icon: 'warning-outline',
      children: [
        { label: '告警记录', key: '/alerts', permission: 'alert:view' },
        { label: '告警处置', key: '/alerts/disposition', permission: 'alert:dispose' },
        { label: '聚合告警', key: '/aggregated-alerts', permission: 'aggregatedAlert:view' },
      ],
    },
    {
      label: '安全分析',
      key: 'security',
      icon: 'shield-checkmark-outline',
      children: [
        { label: '高频IP', key: '/high-freq-ips', permission: 'highFreqIp:view' },
        { label: '脱敏规则', key: '/desensitize', permission: 'desensitizeRule:view' },
      ],
    },
    {
      label: '数据分析',
      key: 'analysis',
      icon: 'stats-chart-outline',
      children: [
        { label: '数据统计', key: '/stats', permission: 'stats:view' },
        { label: '导入导出', key: '/import-export', permission: 'importExport:view' },
      ],
    },
    {
      label: '审计日志',
      key: '/audit-logs',
      icon: 'clipboard-outline',
      permission: 'auditLog:view',
    },
  ]

  function filterVisibleMenus(menus: MenuItem[]): MenuItem[] {
    return menus
      .filter((m) => {
        if (m.permission && !hasPermission(m.permission)) return false
        if (m.children) {
          m.children = filterVisibleMenus(m.children)
          return m.children.length > 0
        }
        return true
      })
      .map((m) => ({
        ...m,
        children: m.children && m.children.length > 0 ? m.children : undefined,
      }))
  }

  const visibleMenus = computed(() => filterVisibleMenus(JSON.parse(JSON.stringify(menuItems))))

  async function fetchPermissions() {
    try {
      const [permRes, roleRes] = await Promise.all([
        getAllPermissions(),
        getAllRoles(),
      ])
      permissions.value = permRes.data || []
      roles.value = roleRes.data || []
      loaded.value = true
    } catch {
      loaded.value = true
    }
  }

  return {
    permissions,
    roles,
    loaded,
    permissionCodes,
    hasPermission,
    menuItems,
    visibleMenus,
    fetchPermissions,
  }
})