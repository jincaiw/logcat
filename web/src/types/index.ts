// ==================== Common ====================

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PageResponse<T = any> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

export interface PageParams {
  page?: number
  pageSize?: number
  sortField?: string
  sortOrder?: 'ascend' | 'descend'
}

// ==================== Auth ====================

export interface User {
  id: number
  username: string
  displayName: string
  email: string
  phone: string
  status: number // 1=active, 0=disabled
  isAdmin: boolean
  mustChangePassword: boolean
  lastLoginAt: string
  createdAt: string
  updatedAt: string
}

export interface Role {
  id: number
  name: string
  code: string
  description: string
  status: number
  createdAt: string
  updatedAt: string
  permissions?: Permission[]
}

export interface Permission {
  id: number
  name: string
  code: string
  group: string
  description: string
}

export interface UserRole {
  userId: number
  roleId: number
}

// ==================== Device ====================

export interface Device {
  id: number
  name: string
  host: string
  port: number
  protocol: string // tcp, udp, tls
  deviceGroupId: number
  deviceGroupName?: string
  deviceTemplateId: number
  deviceTemplateName?: string
  status: number // 1=online, 0=offline
  description: string
  createdAt: string
  updatedAt: string
}

export interface DeviceGroup {
  id: number
  name: string
  description: string
  deviceCount?: number
  createdAt: string
  updatedAt: string
}

// ==================== Device Template ====================

export interface DeviceTemplate {
  id: number
  name: string
  description: string
  config: Record<string, any>
  createdAt: string
  updatedAt: string
}

// ==================== Field Mapping ====================

export interface FieldMappingDoc {
  id: number
  name: string
  description: string
  mappings: FieldMappingItem[]
  createdAt: string
  updatedAt: string
}

export interface FieldMappingItem {
  sourceField: string
  targetField: string
  transformRule: string
}

// ==================== Parse Template ====================

export interface ParseTemplate {
  id: number
  name: string
  description: string
  type: string // regex, json, kv, grok, custom
  pattern: string
  sample: string
  fieldMappings: string
  createdAt: string
  updatedAt: string
}

export interface ParseTestResult {
  success: boolean
  parsed: Record<string, any>
  errors: string[]
}

// ==================== Filter Policy ====================

export interface FilterPolicy {
  id: number
  name: string
  description: string
  priority: number
  status: number // 1=enabled, 0=disabled
  action: string // accept, drop
  conditions: FilterCondition[]
  createdAt: string
  updatedAt: string
}

export interface FilterCondition {
  field: string
  operator: string // eq, ne, contains, regex, gt, lt, gte, lte, in, notIn
  value: string
}

export interface FilterTestResult {
  matched: boolean
  action: string
  conditions: FilterCondition[]
  results: boolean[]
}

// ==================== Output Template ====================

export interface OutputTemplate {
  id: number
  name: string
  description: string
  format: string // json, csv, syslog, raw, custom
  template: string
  createdAt: string
  updatedAt: string
}

// ==================== Push Config ====================

export interface PushConfig {
  id: number
  name: string
  type: string // http, email, syslog
  status: number // 1=enabled, 0=disabled
  config: HttpPushConfig | EmailPushConfig | SyslogPushConfig
  createdAt: string
  updatedAt: string
}

export interface HttpPushConfig {
  url: string
  method: string // POST, PUT
  headers: Record<string, string>
  timeout: number
  retryCount: number
  retryInterval: number
}

export interface EmailPushConfig {
  smtpHost: string
  smtpPort: number
  username: string
  password: string
  from: string
  to: string[]
  subject: string
  useTLS: boolean
}

export interface SyslogPushConfig {
  host: string
  port: number
  protocol: string // tcp, udp
  facility: string
  severity: string
  tag: string
}

// ==================== Alert Rule ====================

export interface AlertRule {
  id: number
  name: string
  description: string
  status: number // 1=enabled, 0=disabled
  severity: string // critical, high, medium, low, info
  condition: AlertCondition
  pushConfigIds: number[]
  cooldownSeconds: number
  createdAt: string
  updatedAt: string
}

export interface AlertCondition {
  field: string
  operator: string // eq, ne, gt, lt, gte, lte, contains, regex
  value: string
  threshold: number
  windowSeconds: number
}

// ==================== Syslog Log ====================

export interface SyslogLog {
  id: string
  receivedAt: string
  deviceName: string
  deviceHost: string
  sourceIp: string
  destIp: string
  facility: string
  severity: string
  eventType: string
  message: string
  rawMessage: string
  parsedData: Record<string, any>
  pushStatus: string // pending, success, failed, skipped
  pushError: string
  traceId: string
}

// ==================== Alert ====================

export interface AlertRecord {
  id: number
  ruleId: number
  ruleName: string
  severity: string
  message: string
  deviceName: string
  sourceIp: string
  logId: string
  channel: string
  status: string // sent, failed, acknowledged, resolved
  error: string
  disposition: AlertDisposition | null
  createdAt: string
}

export interface AlertDisposition {
  id: number
  alertId: number
  action: string // confirm, ignore, close
  note: string
  operatorId: number
  operatorName: string
  createdAt: string
}

export interface AggregatedAlert {
  id: number
  ruleId: number
  ruleName: string
  severity: string
  summary: string
  count: number
  firstAt: string
  lastAt: string
  sourceIps: string[]
  deviceNames: string[]
  status: string
  logIds: string[]
  createdAt: string
  updatedAt: string
}

// ==================== High Freq IP ====================

export interface HighFreqIp {
  id: number
  ip: string
  count: number
  firstSeen: string
  lastSeen: string
  deviceNames: string[]
}

export interface HighFreqIpConfig {
  timeWindowSeconds: number
  threshold: number
}

// ==================== Desensitize ====================

export interface DesensitizeRule {
  id: number
  name: string
  description: string
  field: string
  pattern: string
  replacement: string
  status: number
  createdAt: string
  updatedAt: string
}

// ==================== Audit Log ====================

export interface AuditLog {
  id: number
  userId: number
  username: string
  action: string
  resource: string
  resourceId: string
  detail: string
  result: string // success, failure
  ip: string
  createdAt: string
}

// ==================== System ====================

export interface SystemConfig {
  id: number
  configKey: string
  configValue: string
  description: string
  updatedAt: string
}

export interface SystemStatus {
  serviceRunning: boolean
  startedAt: string
  uptime: string
  listeners: ListenerInfo[]
  connections: ConnectionStatus
  queue: QueueStatus
  workers: WorkerStatus[]
}

export interface ListenerInfo {
  protocol: string
  port: number
  address: string
  status: string
}

export interface ConnectionStatus {
  total: number
  active: number
  idle: number
  closed: number
}

export interface QueueStatus {
  name: string
  size: number
  capacity: number
  enqueueRate: number
  dequeueRate: number
}

export interface WorkerStatus {
  id: number
  status: string
  processedCount: number
  errorCount: number
  lastActiveAt: string
}

// ==================== Dashboard ====================

export interface DashboardStats {
  serviceStatus: string
  receiveRate: number
  todayTotal: number
  todayAlerts: number
  pushSuccessRate: number
  queueBacklog: QueueStatus[]
  recentPushFailures: PushFailureItem[]
  healthStatus: HealthStatus
}

export interface PushFailureItem {
  time: string
  channel: string
  error: string
  logId: string
}

export interface HealthStatus {
  service: string
  cpu: number
  memory: number
  diskUsage: number
  networkIn: number
  networkOut: number
}

// ==================== Stats ====================

export interface StatsQueryParams {
  policy: string
  startTime: string
  endTime: string
  field: string
  topN: number
}

export interface StatsResult {
  field: string
  value: string
  count: number
  percentage: number
}

// ==================== Import/Export ====================

export interface ImportResult {
  success: boolean
  imported: number
  failed: number
  errors: string[]
}