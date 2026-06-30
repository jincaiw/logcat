// 统一 API 层：完全匹配后端 internal/api/router.go 路由定义。
// 移除所有 Wails 回退逻辑，修复签名不一致问题。
import { http } from './client'
import type {
  Device,
  DeviceGroup,
  ParseTemplate,
  OutputTemplate,
  FieldMappingDoc,
  FilterPolicy,
  AlertPolicy,
  AlertRule,
  Robot,
  SystemConfig,
  SystemStats,
  ServiceStatus,
  LogQueryParams,
  LogQueryResult,
  FieldStatsRequest,
  FieldStatsResult,
  StatsField,
  LogTraceInfo,
  ParseTestRequest,
  ParseTestResult,
  TestSyslogRequest,
  TestSyslogResult,
  TestSyslogForwardRequest,
  ImportResult,
  ConfigExport,
  ApiSuccess,
  UnmatchedCount,
  User,
  LoginRequest,
  LoginResponse,
  UpdateProfileRequest,
  ChangePasswordRequest,
} from '@/types'

// ==================== 认证 ====================

export const authApi = {
  login: (data: LoginRequest) =>
    http.post<LoginResponse>('auth/login', data, true),
  logout: () => http.post<ApiSuccess>('auth/logout'),
  getProfile: () => http.get<User>('auth/profile'),
  updateProfile: (data: UpdateProfileRequest) =>
    http.put<User>('auth/profile', data),
  changePassword: (data: ChangePasswordRequest) =>
    http.put<ApiSuccess>('auth/password', data),
}

// ==================== 设备管理 ====================

export const deviceApi = {
  list: () => http.get<Device[]>('devices'),
  get: (id: number) => http.get<Device>(`devices/${id}`),
  create: (data: Omit<Device, 'id' | 'createdAt' | 'updatedAt'>) =>
    http.post<Device>('devices', data),
  update: (id: number, data: Partial<Device>) =>
    http.put<Device>(`devices/${id}`, data),
  delete: (id: number) => http.delete<ApiSuccess>(`devices/${id}`),
}

// ==================== 设备分组 ====================

export const deviceGroupApi = {
  list: () => http.get<DeviceGroup[]>('device-groups'),
  get: (id: number) => http.get<DeviceGroup>(`device-groups/${id}`),
  create: (data: Omit<DeviceGroup, 'id' | 'createdAt' | 'updatedAt'>) =>
    http.post<DeviceGroup>('device-groups', data),
  update: (id: number, data: Partial<DeviceGroup>) =>
    http.put<DeviceGroup>(`device-groups/${id}`, data),
  delete: (id: number) => http.delete<ApiSuccess>(`device-groups/${id}`),
}

// ==================== 解析模板 ====================

export const parseTemplateApi = {
  list: () => http.get<ParseTemplate[]>('parse-templates'),
  get: (id: number) => http.get<ParseTemplate>(`parse-templates/${id}`),
  create: (data: Omit<ParseTemplate, 'id' | 'createdAt' | 'updatedAt'>) =>
    http.post<ParseTemplate>('parse-templates', data),
  update: (id: number, data: Partial<ParseTemplate>) =>
    http.put<ParseTemplate>(`parse-templates/${id}`, data),
  delete: (id: number) => http.delete<ApiSuccess>(`parse-templates/${id}`),
}

// ==================== 输出模板 ====================

export const outputTemplateApi = {
  list: () => http.get<OutputTemplate[]>('output-templates'),
  get: (id: number) => http.get<OutputTemplate>(`output-templates/${id}`),
  create: (data: Omit<OutputTemplate, 'id' | 'createdAt' | 'updatedAt'>) =>
    http.post<OutputTemplate>('output-templates', data),
  update: (id: number, data: Partial<OutputTemplate>) =>
    http.put<OutputTemplate>(`output-templates/${id}`, data),
  delete: (id: number) => http.delete<ApiSuccess>(`output-templates/${id}`),
}

// ==================== 字段映射文档 ====================

export const fieldMappingDocApi = {
  list: () => http.get<FieldMappingDoc[]>('field-mapping-docs'),
  get: (id: number) => http.get<FieldMappingDoc>(`field-mapping-docs/${id}`),
  create: (data: Omit<FieldMappingDoc, 'id' | 'createdAt' | 'updatedAt'>) =>
    http.post<FieldMappingDoc>('field-mapping-docs', data),
  update: (id: number, data: Partial<FieldMappingDoc>) =>
    http.put<FieldMappingDoc>(`field-mapping-docs/${id}`, data),
  delete: (id: number) => http.delete<ApiSuccess>(`field-mapping-docs/${id}`),
}

// ==================== 筛选策略 ====================

export const filterPolicyApi = {
  list: () => http.get<FilterPolicy[]>('filter-policies'),
  get: (id: number) => http.get<FilterPolicy>(`filter-policies/${id}`),
  create: (data: Omit<FilterPolicy, 'id' | 'createdAt' | 'updatedAt'>) =>
    http.post<FilterPolicy>('filter-policies', data),
  update: (id: number, data: Partial<FilterPolicy>) =>
    http.put<FilterPolicy>(`filter-policies/${id}`, data),
  delete: (id: number) => http.delete<ApiSuccess>(`filter-policies/${id}`),
}

// ==================== 告警策略 ====================

export const alertPolicyApi = {
  list: () => http.get<AlertPolicy[]>('alert-policies'),
  get: (id: number) => http.get<AlertPolicy>(`alert-policies/${id}`),
  create: (data: Omit<AlertPolicy, 'id' | 'createdAt' | 'updatedAt'>) =>
    http.post<AlertPolicy>('alert-policies', data),
  update: (id: number, data: Partial<AlertPolicy>) =>
    http.put<AlertPolicy>(`alert-policies/${id}`, data),
  delete: (id: number) => http.delete<ApiSuccess>(`alert-policies/${id}`),
}

// ==================== 告警规则 ====================

export const alertRuleApi = {
  listByRobot: (robotId: number) =>
    http.get<AlertRule[]>(`alert-rules/robot/${robotId}`),
  get: (id: number) => http.get<AlertRule>(`alert-rules/${id}`),
  create: (data: Omit<AlertRule, 'id' | 'createdAt' | 'updatedAt'>) =>
    http.post<AlertRule>('alert-rules', data),
  update: (id: number, data: Partial<AlertRule>) =>
    http.put<AlertRule>(`alert-rules/${id}`, data),
  delete: (id: number) => http.delete<ApiSuccess>(`alert-rules/${id}`),
  deleteByRobot: (robotId: number) =>
    http.delete<ApiSuccess>(`alert-rules/robot/${robotId}`),
}

// ==================== 推送通道（机器人） ====================

export const robotApi = {
  list: () => http.get<Robot[]>('robots'),
  get: (id: number) => http.get<Robot>(`robots/${id}`),
  create: (data: Omit<Robot, 'id' | 'createdAt' | 'updatedAt'>) =>
    http.post<Robot>('robots', data),
  update: (id: number, data: Partial<Robot>) =>
    http.put<Robot>(`robots/${id}`, data),
  delete: (id: number) => http.delete<ApiSuccess>(`robots/${id}`),
  test: (data: Partial<Robot>) => http.post<ApiSuccess>('test-robot', data),
}

// ==================== 日志 ====================

export const logApi = {
  list: (params: LogQueryParams) =>
    http.get<LogQueryResult>('logs', {
      page: params.page,
      pageSize: params.pageSize,
      deviceId: params.deviceId,
      startTime: params.startTime,
      endTime: params.endTime,
      keyword: params.keyword,
    }),
  cleanup: (days: number) => http.post<ApiSuccess>('logs/cleanup', { days }),
  cleanupAll: () => http.delete<ApiSuccess>('logs/cleanup-all'),
  getUnmatchedCount: () => http.get<UnmatchedCount>('logs/unmatched-count'),
  cleanupUnmatched: (days: number) =>
    http.post<ApiSuccess>('logs/cleanup-unmatched', { days }),
}

// ==================== 服务管理 ====================

export const serviceApi = {
  getStatus: () => http.get<ServiceStatus>('service/status'),
  start: (port: number, protocol: string) =>
    http.post<ApiSuccess>('service/start', { port, protocol }),
  stop: () => http.post<ApiSuccess>('service/stop'),
}

// ==================== 配置 ====================

export const configApi = {
  get: () => http.get<SystemConfig>('config'),
  save: (config: Partial<SystemConfig>) => http.put<SystemConfig>('config', config),
}

// ==================== 统计 ====================

export const statsApi = {
  getSystemStats: () => http.get<SystemStats>('stats'),
  getFieldStats: (req: FieldStatsRequest) =>
    http.post<FieldStatsResult>('field-stats', req),
  getAvailableFields: (policyId: number) =>
    http.get<StatsField[]>(`available-stats-fields/${policyId}`),
}

// ==================== 测试工具 ====================

export const testApi = {
  sendTestSyslog: (req: TestSyslogRequest) =>
    http.post<TestSyslogResult>('test-syslog', req),
  testSyslogForward: (req: TestSyslogForwardRequest) =>
    http.post<ApiSuccess>('test-syslog-forward', req),
  testParse: (req: ParseTestRequest) =>
    http.post<ParseTestResult>('test-parse', req),
}

// ==================== 日志追踪 ====================

export const traceApi = {
  get: (logId: number) => http.get<LogTraceInfo>(`log-trace/${logId}`),
}

// ==================== 网络 ====================

export const networkApi = {
  getLocalIPs: () => http.get<string[]>('local-ips'),
  getServerIP: () => http.get<string>('server-ip'),
}

// ==================== 导入导出 ====================

export const importExportApi = {
  exportParseTemplates: (ids?: number[]) =>
    http.get<ConfigExport>('export/parse-templates', {
      ids: ids?.length ? ids.join(',') : undefined,
    }),
  exportFilterPolicies: (ids?: number[]) =>
    http.get<ConfigExport>('export/filter-policies', {
      ids: ids?.length ? ids.join(',') : undefined,
    }),
  importParseTemplates: (data: ConfigExport) =>
    http.post<ImportResult>('import/parse-templates', data),
  importFilterPolicies: (data: ConfigExport) =>
    http.post<ImportResult>('import/filter-policies', data),
}

// ==================== 预设模板 ====================

export const presetApi = {
  getTemplates: () => http.get<unknown[]>('preset-templates'),
}

// ==================== 兼容性导出 ====================
// 保留旧 API 对象以渐进迁移视图文件，新代码应使用上述分模块 API。

export const API = {
  // Devices
  GetDevices: () => deviceApi.list(),
  GetDevice: (id: number) => deviceApi.get(id),
  AddDevice: (device: Parameters<typeof deviceApi.create>[0]) => deviceApi.create(device),
  UpdateDevice: (device: Partial<Device> & { id: number }) => deviceApi.update(device.id, device),
  DeleteDevice: (id: number) => deviceApi.delete(id),

  // Device Groups
  GetDeviceGroups: () => deviceGroupApi.list(),
  AddDeviceGroup: (group: Parameters<typeof deviceGroupApi.create>[0]) => deviceGroupApi.create(group),
  UpdateDeviceGroup: (group: Partial<DeviceGroup> & { id: number }) => deviceGroupApi.update(group.id, group),
  DeleteDeviceGroup: (id: number) => deviceGroupApi.delete(id),

  // Parse Templates
  GetParseTemplates: () => parseTemplateApi.list(),
  AddParseTemplate: (template: Parameters<typeof parseTemplateApi.create>[0]) => parseTemplateApi.create(template),
  UpdateParseTemplate: (template: Partial<ParseTemplate> & { id: number }) => parseTemplateApi.update(template.id, template),
  DeleteParseTemplate: (id: number) => parseTemplateApi.delete(id),

  // Output Templates
  GetOutputTemplates: () => outputTemplateApi.list(),
  AddOutputTemplate: (template: Parameters<typeof outputTemplateApi.create>[0]) => outputTemplateApi.create(template),
  UpdateOutputTemplate: (template: Partial<OutputTemplate> & { id: number }) => outputTemplateApi.update(template.id, template),
  DeleteOutputTemplate: (id: number) => outputTemplateApi.delete(id),

  // Filter Policies
  GetFilterPolicies: () => filterPolicyApi.list(),
  AddFilterPolicy: (policy: Parameters<typeof filterPolicyApi.create>[0]) => filterPolicyApi.create(policy),
  UpdateFilterPolicy: (policy: Partial<FilterPolicy> & { id: number }) => filterPolicyApi.update(policy.id, policy),
  DeleteFilterPolicy: (id: number) => filterPolicyApi.delete(id),

  // Alert Policies
  GetAlertPolicies: () => alertPolicyApi.list(),
  AddAlertPolicy: (policy: Parameters<typeof alertPolicyApi.create>[0]) => alertPolicyApi.create(policy),
  UpdateAlertPolicy: (policy: Partial<AlertPolicy> & { id: number }) => alertPolicyApi.update(policy.id, policy),
  DeleteAlertPolicy: (id: number) => alertPolicyApi.delete(id),

  // Alert Rules
  GetAlertRules: (robotId: number) => alertRuleApi.listByRobot(robotId),
  AddAlertRule: (rule: Parameters<typeof alertRuleApi.create>[0]) => alertRuleApi.create(rule),
  UpdateAlertRule: (rule: Partial<AlertRule> & { id: number }) => alertRuleApi.update(rule.id, rule),
  DeleteAlertRule: (id: number) => alertRuleApi.delete(id),
  DeleteAlertRulesByRobotID: (robotId: number) => alertRuleApi.deleteByRobot(robotId),

  // Field Mapping Docs
  GetFieldMappingDocs: () => fieldMappingDocApi.list(),
  AddFieldMappingDoc: (doc: Parameters<typeof fieldMappingDocApi.create>[0]) => fieldMappingDocApi.create(doc),
  UpdateFieldMappingDoc: (doc: Partial<FieldMappingDoc> & { id: number }) => fieldMappingDocApi.update(doc.id, doc),
  DeleteFieldMappingDoc: (id: number) => fieldMappingDocApi.delete(id),

  // Robots
  GetRobots: () => robotApi.list(),
  AddRobot: (robot: Parameters<typeof robotApi.create>[0]) => robotApi.create(robot),
  UpdateRobot: (robot: Partial<Robot> & { id: number }) => robotApi.update(robot.id, robot),
  DeleteRobot: (id: number) => robotApi.delete(id),
  TestRobot: (robot: Partial<Robot>) => robotApi.test(robot),

  // Logs
  QueryLogs: (params: LogQueryParams) => logApi.list(params),
  CleanupLogs: (days: number) => logApi.cleanup(days),

  // Service
  StartService: (port: number, protocol: string) => serviceApi.start(port, protocol),
  StopService: () => serviceApi.stop(),
  GetServiceStatus: () => serviceApi.getStatus(),

  // Config
  GetSystemConfig: () => configApi.get(),
  SaveSystemConfig: (config: Partial<SystemConfig>) => configApi.save(config),

  // Stats
  GetSystemStats: () => statsApi.getSystemStats(),
  GetFieldStats: (req: FieldStatsRequest) => statsApi.getFieldStats(req),
  GetAvailableStatsFields: (policyId: number) => statsApi.getAvailableFields(policyId),

  // Test
  SendTestSyslog: (req: TestSyslogRequest) => testApi.sendTestSyslog(req),
  TestParseTemplate: (req: ParseTestRequest) => testApi.testParse(req),
  TestSyslogForward: (req: TestSyslogForwardRequest) => testApi.testSyslogForward(req),
  GetLogTrace: (id: number) => traceApi.get(id),

  // Local IPs
  GetLocalIPs: () => networkApi.getLocalIPs(),

  // Export/Import
  ExportParseTemplates: (ids?: number[]) => importExportApi.exportParseTemplates(ids),
  ExportFilterPolicies: (ids?: number[]) => importExportApi.exportFilterPolicies(ids),
  ImportParseTemplates: (data: ConfigExport) => importExportApi.importParseTemplates(data),
  ImportFilterPolicies: (data: ConfigExport) => importExportApi.importFilterPolicies(data),
}
