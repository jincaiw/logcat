import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: { public: true },
  },
  {
    path: '/init',
    name: 'Init',
    component: () => import('@/views/init/InitAdminView.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/DashboardView.vue'),
        meta: { title: '仪表盘', permission: 'dashboard:view' },
      },
      {
        path: 'system/status',
        name: 'SystemStatus',
        component: () => import('@/views/system/SystemStatusView.vue'),
        meta: { title: '服务状态', permission: 'system:status' },
      },
      {
        path: 'system/config',
        name: 'SystemConfig',
        component: () => import('@/views/system/SystemConfigView.vue'),
        meta: { title: '系统配置', permission: 'system:config' },
      },
      {
        path: 'users',
        name: 'UserManagement',
        component: () => import('@/views/users/UserManagementView.vue'),
        meta: { title: '用户管理', permission: 'user:view' },
      },
      {
        path: 'roles',
        name: 'RoleManagement',
        component: () => import('@/views/roles/RoleManagementView.vue'),
        meta: { title: '角色管理', permission: 'role:view' },
      },
      {
        path: 'devices',
        name: 'DeviceManagement',
        component: () => import('@/views/devices/DeviceManagementView.vue'),
        meta: { title: '设备管理', permission: 'device:view' },
      },
      {
        path: 'device-groups',
        name: 'DeviceGroupManagement',
        component: () => import('@/views/devices/DeviceGroupManagementView.vue'),
        meta: { title: '设备分组', permission: 'deviceGroup:view' },
      },
      {
        path: 'device-templates',
        name: 'DeviceTemplate',
        component: () => import('@/views/device-templates/DeviceTemplateView.vue'),
        meta: { title: '设备模板', permission: 'deviceTemplate:view' },
      },
      {
        path: 'field-mappings',
        name: 'FieldMapping',
        component: () => import('@/views/field-mappings/FieldMappingView.vue'),
        meta: { title: '字段映射', permission: 'fieldMapping:view' },
      },
      {
        path: 'parse-templates',
        name: 'ParseTemplate',
        component: () => import('@/views/parse-templates/ParseTemplateView.vue'),
        meta: { title: '解析模板', permission: 'parseTemplate:view' },
      },
      {
        path: 'filter-policies',
        name: 'FilterPolicy',
        component: () => import('@/views/filter-policies/FilterPolicyView.vue'),
        meta: { title: '过滤策略', permission: 'filterPolicy:view' },
      },
      {
        path: 'output-templates',
        name: 'OutputTemplate',
        component: () => import('@/views/output-templates/OutputTemplateView.vue'),
        meta: { title: '输出模板', permission: 'outputTemplate:view' },
      },
      {
        path: 'push-configs',
        name: 'PushConfig',
        component: () => import('@/views/push-configs/PushConfigView.vue'),
        meta: { title: '推送配置', permission: 'pushConfig:view' },
      },
      {
        path: 'alert-rules',
        name: 'AlertRule',
        component: () => import('@/views/alert-rules/AlertRuleView.vue'),
        meta: { title: '告警规则', permission: 'alertRule:view' },
      },
      {
        path: 'logs',
        name: 'LogQuery',
        component: () => import('@/views/logs/LogQueryView.vue'),
        meta: { title: '日志查询', permission: 'log:view' },
      },
      {
        path: 'logs/trace/:id',
        name: 'LogTrace',
        component: () => import('@/views/logs/LogTraceView.vue'),
        meta: { title: '日志追踪', permission: 'log:view' },
      },
      {
        path: 'alerts',
        name: 'AlertRecord',
        component: () => import('@/views/alerts/AlertRecordView.vue'),
        meta: { title: '告警记录', permission: 'alert:view' },
      },
      {
        path: 'alerts/disposition',
        name: 'AlertDisposition',
        component: () => import('@/views/alerts/AlertDispositionView.vue'),
        meta: { title: '告警处置', permission: 'alert:dispose' },
      },
      {
        path: 'aggregated-alerts',
        name: 'AggregatedAlert',
        component: () => import('@/views/aggregated-alerts/AggregatedAlertView.vue'),
        meta: { title: '聚合告警', permission: 'aggregatedAlert:view' },
      },
      {
        path: 'high-freq-ips',
        name: 'HighFreqIp',
        component: () => import('@/views/high-freq-ips/HighFreqIpView.vue'),
        meta: { title: '高频IP', permission: 'highFreqIp:view' },
      },
      {
        path: 'desensitize',
        name: 'DesensitizeRule',
        component: () => import('@/views/desensitize/DesensitizeRuleView.vue'),
        meta: { title: '脱敏规则', permission: 'desensitizeRule:view' },
      },
      {
        path: 'stats',
        name: 'Stats',
        component: () => import('@/views/stats/StatsView.vue'),
        meta: { title: '数据统计', permission: 'stats:view' },
      },
      {
        path: 'import-export',
        name: 'ImportExport',
        component: () => import('@/views/import-export/ImportExportView.vue'),
        meta: { title: '导入导出', permission: 'importExport:view' },
      },
      {
        path: 'audit-logs',
        name: 'AuditLog',
        component: () => import('@/views/audit-logs/AuditLogView.vue'),
        meta: { title: '审计日志', permission: 'auditLog:view' },
      },
      {
        path: 'change-password',
        name: 'ChangePassword',
        component: () => import('@/views/change-password/ChangePasswordView.vue'),
        meta: { title: '修改密码' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()
  const permissionStore = usePermissionStore()

  // Public routes
  if (to.meta.public) {
    if (authStore.isAuthenticated && to.path === '/login') {
      return next('/')
    }
    return next()
  }

  // Check auth
  if (!authStore.isInitialized) {
    await authStore.fetchCurrentUser()
  }

  if (!authStore.isAuthenticated) {
    return next('/login')
  }

  // Check must change password
  if (authStore.mustChangePassword && to.path !== '/change-password') {
    return next('/change-password')
  }

  // Load permissions if not loaded
  if (!permissionStore.loaded) {
    await permissionStore.fetchPermissions()
  }

  // Check permission
  if (to.meta.permission && !permissionStore.hasPermission(to.meta.permission as string)) {
    return next('/')
  }

  next()
})

export default router