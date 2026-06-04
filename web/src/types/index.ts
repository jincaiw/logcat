// ==================== Common ====================

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PageResponse<T = any> {
  items?: T[]
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
  phone?: string
  status: string
  isAdmin?: boolean
  mustChangePassword: boolean
  lastLoginAt?: string | null
  createdAt: string
  updatedAt: string
  roles?: string[]
  permissions?: string[]
}

export interface Role {
  id: number
  name: string
  code: string
  description: string
  status?: number | string
  builtIn?: boolean
  createdAt: string
  updatedAt: string
  permissions?: Permission[]
}

export interface Permission {
  id: number
  name: string
  code: string
  type?: string
  resource?: string
  action?: string
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
  ipAddress: string
  groupId: number | null
  templateId: number | null
  parseTemplateId: number | null
  deviceType: string
  description: string
  enabled: boolean
  createdAt: string
  updatedAt: string
  group?: DeviceGroup | null
  template?: DeviceTemplate | null
  parseTemplate?: ParseTemplate | null
}

export interface DeviceGroup {
  id: number
  name: string
  description: string
  color: string
  sortOrder: number
  deviceCount?: number
  createdAt: string
  updatedAt: string
}

// ==================== Device Template ====================

export interface DeviceTemplate {
  id: number
  name: string
  deviceType: string
  parseTemplateId: number | null
  fieldMappingDocId: number | null
  recommendedPolicy: string
  enabled: boolean
  createdAt: string
  updatedAt: string
}

// ==================== Field Mapping ====================

export interface FieldMappingDoc {
  id: number
  deviceType: string
  standardField: string
  originalField: string
  description: string
  fieldType: string
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
  deviceType?: string
  parseType: string
  headerRegex?: string
  delimiter?: string
  fieldMapping?: string
  valueTransform?: string
  sampleLog?: string
  subTemplates?: string
  enabled?: boolean
  createdAt: string
  updatedAt: string
}

export interface ParseTestResult {
  success: boolean
  fields: Record<string, any>
  error?: string
}

// ==================== Filter Policy ====================

export interface FilterPolicy {
  id: number
  name: string
  deviceId: number | null
  deviceGroupId: number | null
  parseTemplateId: number | null
  conditions: string
  conditionLogic: string
  whitelistEnabled: boolean
  whitelistField: string
  whitelistValues: string
  action: string
  priority: number
  dedupEnabled: boolean
  dedupWindow: number
  enabled: boolean
  createdAt: string
  updatedAt: string
  deviceGroup?: DeviceGroup | null
  parseTemplate?: ParseTemplate | null
}

export interface FilterCondition {
  field: string
  operator: string // eq, ne, contains, regex, gt, lt, gte, lte, in, notIn
  value: string
}

export interface FilterTestResult {
  matched: boolean
  action: string
  message?: string
  whitelistResult?: string
  policy?: Partial<FilterPolicy>
}

// ==================== Output Template ====================

export interface OutputTemplate {
  id: number
  name: string
  channelType: string
  content: string
  fields: string
  deviceType: string
  enabled: boolean
  createdAt: string
  updatedAt: string
}

// ==================== Push Config ====================

export interface PushConfig {
  id: number
  name: string
  type: string // http, email, syslog
  enabled: boolean
  url?: string
  method?: string
  timeout?: number
  retryCount?: number
  retryDelay?: number
  notesIds?: string
  headers?: string
  bodyTemplate?: string
  successStatusCodes?: string
  successBodyKeyword?: string
  authType?: string
  token?: string
  contentType?: string
  retryOnStatusCodes?: string
  maxResponseLogSize?: number
  smtpHost?: string
  smtpPort?: number
  smtpUsername?: string
  smtpPassword?: string
  fromAddress?: string
  toAddresses?: string
  subjectTemplate?: string
  emailBodyTemplate?: string
  syslogHost?: string
  syslogPort?: number
  syslogProtocol?: string
  syslogFormat?: string
  syslogFields?: string
  createdAt: string
  updatedAt: string
}

export interface PushTestResult {
  success: boolean
  channel: string
  statusCode: number
  responseBody: string
  errorMessage?: string
  summary?: string
}

// ==================== Alert Rule ====================

export interface AlertRule {
  id: number
  name: string
  filterPolicyId: number | null
  pushConfigId: number | null
  outputTemplateId: number | null
  channelType: string
  enabled: boolean
  createdAt: string
  updatedAt: string
  filterPolicy?: FilterPolicy | null
  pushConfig?: PushConfig | null
  outputTemplate?: OutputTemplate | null
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
  id: number
  logId: string
  sourceIp: string
  destinationIp: string
  eventType: string
  severity: string
  facility: string
  deviceId: number | null
  deviceName: string
  rawMessage: string
  parsedData: string
  filterStatus: string
  matchedFilterPolicyId: number | null
  alertStatus: string
  alertRuleId: number | null
  aggregatedAlertId: number | null
  receivedAt: string
  createdAt: string
}

// ==================== Alert ====================

export interface AlertRecord {
  id: number
  logId: string
  alertRuleId: number | null
  pushConfigId: number | null
  channelType: string
  status: string
  retryCount: number
  requestSummary: string
  responseStatusCode: number
  responseSummary: string
  errorMessage: string
  dispositionStatus: string
  sentAt: string | null
  createdAt: string
  alertRule?: AlertRule | null
  pushConfig?: PushConfig | null
}

export interface AlertDisposition {
  id: number
  alertRecordId?: number
  aggregatedAlertId?: number
  status: string // confirmed, ignored, closed, acknowledged, resolved
  note: string
  operatorId?: number
  operatorName: string
  operatedAt?: string
  createdAt: string
}

export interface AggregatedAlert {
  id: number
  aggregateKey: string
  aggregateType: string
  sourceIp: string
  destinationIp: string
  eventType: string
  deviceId: number | null
  severity: string
  count: number
  firstSeenAt: string
  lastSeenAt: string
  status: string
  createdAt: string
  updatedAt: string
}

// ==================== High Freq IP ====================

export interface HighFreqIp {
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
  fieldName: string
  ruleType: string
  ruleConfig: string
  enabled: boolean
  createdAt: string
  updatedAt: string
}

// ==================== Audit Log ====================

export interface AuditLog {
  id: number
  userId?: number
  username: string
  action: string
  resource?: string
  resourceType?: string
  resourceId: string
  detail: string
  result: string // success, failure
  ip?: string
  clientIp?: string
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
  uptimeSeconds?: number
  database?: string
  databaseStatus?: string
  goVersion?: string
  listeners: ListenerInfo[]
  connections: ConnectionStatus
  queue: QueueStatus
  workers: WorkerStatus[]
  receiverMetrics?: ReceiverMetricsSnapshot
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
  enqueueRate?: number
  dequeueRate?: number
}

export interface WorkerStatus {
  id: number
  stage?: string
  count?: number
  status: string
  processedCount: number
  errorCount: number
  lastActiveAt: string
}

export interface ReceiverMetricsSnapshot {
  udpReceived: number
  tcpReceived: number
  udpErrors: number
  tcpErrors: number
  parseErrors: number
  channelDropped: number
  tcpConnections: number
  lastReceiveAt?: string
}

// ==================== Dashboard ====================

export interface DashboardStats {
  serviceStatus: string
  receiveRate: number
  todayTotal: number
  todayAlerts: number
  pushSuccessRate: number
  tcpConnections: number
  lastReceivedAt: string
  parseSuccessRate: number
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
  field?: string
  value: string
  count: number
  percentage: number
  lastSeenAt?: string
}

export interface StatsResponse {
  results: StatsResult[]
  totalLogs: number
  uniqueValues: number
}

export interface AvailableField {
  value: string
  label: string
}

export interface SystemHealthCheck {
  status: string
  time?: string
  reason?: string
}

export interface DatabaseMetricsSnapshot {
  open_connections?: number
  in_use?: number
  idle?: number
  wait_count?: number
  max_open?: number
}

export interface PipelineMetricsSnapshot {
  rawQueueDepth: number
  parsedQueueDepth: number
  dbQueueDepth: number
  pushQueueDepth: number
  parseProcessed: number
  parseErrors: number
  filterProcessed: number
  filterDropped: number
  dbWritten: number
  dbErrors: number
  pushProcessed: number
  pushErrors: number
  rawDropped: number
  dbDropped: number
  pushDropped: number
}

export interface RuntimeMetricsSnapshot {
  uptime_seconds: number
  status: string
  goroutines: number
  heap_alloc_mb: number
  heap_sys_mb: number
  heap_total_mb: number
  num_gc: number
  receiver: Partial<ReceiverMetricsSnapshot>
  pipeline: Partial<PipelineMetricsSnapshot>
}

export interface SystemMetricsSnapshot {
  status: string
  database: string
  db_stats?: DatabaseMetricsSnapshot
  uptime: number
  started_at: string
  receiver_metrics?: Record<string, number>
  pipeline_metrics?: Record<string, number>
}

// ==================== Import/Export ====================

export interface ImportResult {
  resourceType: string
  version: string
  created: number
  updated: number
  failed: number
  errors: string[]
}
