// 类型定义：与后端 internal/models 结构对齐，保持 JSON 字段名一致。
// 仅导出接口，不包含运行时逻辑。

// ==================== 基础实体 ====================

export interface DeviceGroup {
  id: number
  name: string
  description: string
  color: string
  sortOrder: number
  createdAt: string
  updatedAt: string
}

export interface Device {
  id: number
  name: string
  ipAddress: string
  groupId: number
  groupName: string
  description: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface ParseTemplate {
  id: number
  name: string
  description: string
  parseType: ParseType
  headerRegex: string
  fieldMapping: string
  valueTransform: string
  sampleLog: string
  deviceType: string
  delimiter: string
  typeField: number
  subTemplates: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface OutputTemplate {
  id: number
  name: string
  platform: string
  description: string
  content: string
  fields: string
  deviceType: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface FieldMappingDoc {
  id: number
  name: string
  deviceType: string
  description: string
  fieldMappings: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface FilterPolicy {
  id: number
  name: string
  description: string
  deviceId: number
  deviceGroupId: number
  parseTemplateId: number
  conditions: string
  conditionLogic: 'AND' | 'OR'
  whitelist: string
  whitelistField: string
  action: 'keep' | 'discard'
  priority: number
  isActive: boolean
  dedupEnabled: boolean
  dedupWindow: number
  dropUnmatched: boolean
  createdAt: string
  updatedAt: string
}

export interface FilterCondition {
  field: string
  operator: string
  value: string
}

export interface WhitelistItem {
  cidr: string
  description: string
  enabled: boolean
}

export interface AlertPolicy {
  id: number
  name: string
  description: string
  filterPolicyId: number
  robotId: number
  outputTemplateId: number
  deviceId: number
  deviceGroupId: number
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface AlertRule {
  id: number
  robotId: number
  filterPolicyId: number
  outputTemplateId: number
  outputFormat: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface Robot {
  id: number
  name: string
  platform: string
  description: string
  isActive: boolean
  feishuWebhookUrl: string
  feishuSecret: string
  smtpHost: string
  smtpPort: number
  smtpUsername: string
  smtpPassword: string
  smtpFrom: string
  smtpTo: string
  syslogHost: string
  syslogPort: number
  syslogProtocol: string
  syslogFormat: string
  createdAt: string
  updatedAt: string
}

export interface AlertRecord {
  id: number
  logId: number
  robotId: number
  alertPolicyId: number
  deviceName: string
  message: string
  status: string
  errorMsg: string
  sentAt: string
}

export interface SyslogLog {
  id: number
  deviceId: number
  deviceName: string
  sourceIp: string
  rawMessage: string
  parsedData: string
  parsedFields: string
  filterStatus: string
  matchedPolicyId: number
  alertStatus: string
  alertPolicyId: number
  priority: string
  facility: number
  severity: number
  timestamp: string
  receivedAt: string
  isProcessed: boolean
  isAlerted: boolean
}

// ==================== 枚举类型 ====================

export type ParseType =
  | 'json'
  | 'regex'
  | 'kv'
  | 'syslog_json'
  | 'smart_delimiter'
  | 'delimiter'
  | 'keyvalue'

export type FilterOperator =
  | 'equals'
  | 'not_equals'
  | 'contains'
  | 'not_contains'
  | 'in'
  | 'not_in'
  | 'starts_with'
  | 'ends_with'
  | 'regex'
  | 'not_regex'
  | 'exists'
  | 'not_exists'
  | 'gt'
  | 'gte'
  | 'lt'
  | 'lte'
  // 别名（后端兼容）
  | '=='
  | '!='
  | '=~'
  | '!~'
  | '>'
  | '>='
  | '<'
  | '<='

export type Platform = 'feishu' | 'email' | 'syslog'

export type Protocol = 'udp' | 'tcp'

// ==================== DTO ====================

export interface SystemConfig {
  id: number
  listenPort: number
  protocol: Protocol
  logRetention: number
  maxLogSize: number
  autoStart: boolean
  minimizeToTray: boolean
  alertEnabled: boolean
  alertInterval: number
  unmatchedLogRetention: number
  unmatchedLogAlert: boolean
  defaultFilterAction: 'keep' | 'discard'
  theme: 'dark' | 'light'
  language: 'zh-CN' | 'en-US'
  dataDir: string
  configDir: string
}

export interface SystemStats {
  totalLogs: number
  deviceCount: number
  serviceRunning: boolean
  listenPort: number
  startTime: string
  memoryUsage: number
  cpuUsage: number
  connections: number
  receiveRate: number
  protocol: Protocol
  databaseSize: number
  // 仪表盘扩展字段（后端可能不返回，默认 0）
  matchedLogs?: number
  alertCount?: number
  unmatchedLogs?: number
  activeRobots?: number
  activeFilterPolicies?: number
  activeAlertPolicies?: number
  parseTemplateCount?: number
  activeDevices?: number
  goroutineCount?: number
}

export interface ServiceStatus {
  serviceRunning: boolean
  listenPort: number
  receiveCount: number
  receiveRate: number
  connections: number
}

export interface LogQueryParams {
  page?: number
  pageSize?: number
  deviceId?: number
  startTime?: string
  endTime?: string
  keyword?: string
}

export interface LogQueryResult {
  logs: SyslogLog[]
  total: number
}

export interface FieldStatsRequest {
  deviceId?: number
  filterPolicyId?: number
  startTime: string
  endTime: string
  field: string
  topN: number
}

export interface StatsItem {
  value: string
  location: string
  count: number
  percent: string
  lastSeen: string
}

export interface FieldStatsResult {
  field: string
  totalLogs: number
  uniqueCount: number
  items: StatsItem[]
}

export interface StatsField {
  name: string
  displayName: string
}

export interface LogTraceInfo {
  logId: number
  receivedAt: string
  sourceIp: string
  rawMessage: string
  receiveStatus: string
  receiveError?: string
  parseStatus: string
  parseTemplate?: string
  parsedData?: string
  parseError?: string
  filterStatus: string
  filterEnabled: boolean
  matchedPolicy?: string
  filterConditions?: string
  filterResult?: string
  alertStatus: string
  alertRecords?: AlertTraceInfo[]
}

export interface AlertTraceInfo {
  robotId: number
  robotName: string
  platform: string
  status: string
  errorMsg?: string
  sentAt?: string
}

export interface ParseTestRequest {
  templateId?: number
  parseType: ParseType
  headerRegex?: string
  fieldMapping?: string
  valueTransform?: string
  sampleLog: string
  delimiter?: string
}

export interface ParseTestResult {
  success: boolean
  error: string
  fields: string[]
  data: Record<string, unknown>
}

export interface TestSyslogRequest {
  host: string
  port: number
  protocol: Protocol
  message: string
  count: number
  intervalMs: number
}

export interface TestSyslogResult {
  success: boolean
  message: string
  sentCount: number
  failedCount: number
  errors: string[]
}

export interface TestSyslogForwardRequest {
  host: string
  port: number
  protocol: Protocol
  format: string
}

export interface ImportResult {
  success: boolean
  message: string
  count: number
  errors: string[]
}

export interface ConfigExport {
  version: string
  exportedAt: string
  name: string
  description?: string
  parseTemplates?: ParseTemplate[]
  filterPolicies?: FilterPolicy[]
}

// ==================== 通用响应 ====================

export interface ApiSuccess {
  success: boolean
}

export interface UnmatchedCount {
  count: number
}

// ==================== 认证 ====================

export interface User {
  id: number
  username: string
  nickname: string
  email: string
  avatar: string
  createdAt: string
  updatedAt: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface UpdateProfileRequest {
  nickname?: string
  email?: string
  avatar?: string
}

export interface ChangePasswordRequest {
  oldPassword: string
  newPassword: string
}
