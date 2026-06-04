import type { ApiResponse } from '@/types'

const mockUser = {
  id: 1,
  username: 'admin',
  displayName: '管理员',
  email: 'admin@logcat.local',
  status: 'enabled',
  isAdmin: true,
  mustChangePassword: false,
  lastLoginAt: '2026-05-30T10:00:00Z',
  createdAt: '2026-01-01T00:00:00Z',
  updatedAt: '2026-05-30T10:00:00Z',
  roles: ['admin'],
  permissions: [
    'dashboard:view', 'system:status', 'system:config:read',
    'users:list', 'roles:list', 'devices:list', 'device-groups:list',
    'device-templates:list', 'field-mappings:list', 'parse-templates:list',
    'filter-policies:list', 'output-templates:list', 'push-configs:list',
    'alert-rules:list', 'logs:list', 'logs:trace', 'alerts:list',
    'alerts:disposition:list', 'aggregated-alerts:list', 'high-freq-ips:list',
    'desensitize-rules:list', 'stats:fields', 'export:config', 'audit-logs:list',
  ],
}

const mockDevices = Array.from({ length: 12 }, (_, i) => ({
  id: i + 1,
  name: `防火墙-${String(i + 1).padStart(3, '0')}`,
  host: `10.0.${Math.floor(i / 10)}.${i + 1}`,
  port: 514,
  protocol: i % 3 === 0 ? 'tcp' : i % 3 === 1 ? 'udp' : 'tls',
  deviceGroupId: (i % 3) + 1,
  deviceGroupName: ['核心交换区', 'DMZ区域', '办公区域'][(i % 3)],
  deviceTemplateId: (i % 4) + 1,
  deviceTemplateName: ['FortiGate模板', 'PaloAlto模板', 'Cisco ASA模板', '通用Syslog模板'][(i % 4)],
  status: i === 3 ? 0 : 1,
  description: `${['核心防火墙', '边界防火墙', '内网防火墙', 'VPN网关'][i % 4]} - 机房${Math.floor(i / 4) + 1}`,
  createdAt: '2026-01-15T08:00:00Z',
  updatedAt: '2026-05-30T10:00:00Z',
}))

const mockDeviceGroups = [
  { id: 1, name: '核心交换区', description: '核心网络交换设备组', deviceCount: 4, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 2, name: 'DMZ区域', description: 'DMZ区域安全设备组', deviceCount: 4, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 3, name: '办公区域', description: '办公网络设备组', deviceCount: 4, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
]

const mockDeviceTemplates = [
  { id: 1, name: 'FortiGate模板', description: 'FortiGate防火墙日志模板', config: { vendor: 'fortinet', format: 'syslog' }, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 2, name: 'PaloAlto模板', description: 'PaloAlto防火墙日志模板', config: { vendor: 'paloalto', format: 'syslog' }, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 3, name: 'Cisco ASA模板', description: 'Cisco ASA防火墙日志模板', config: { vendor: 'cisco', format: 'syslog' }, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 4, name: '通用Syslog模板', description: '通用Syslog日志模板', config: { vendor: 'generic', format: 'syslog' }, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
]

const mockParseTemplates = [
  { id: 1, name: 'FortiGate解析', deviceType: 'firewall', parseType: 'regex', headerRegex: 'src=(\\S+)', fieldMapping: '{"src":"source_ip"}', valueTransform: '', sampleLog: 'src=10.0.0.1 dst=10.0.0.2 action=allow', enabled: true, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 2, name: 'JSON格式解析', deviceType: 'siem', parseType: 'json', fieldMapping: '{"event":"event_type"}', valueTransform: '', sampleLog: '{"event":"login","severity":"high"}', enabled: true, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 3, name: 'KV格式解析', deviceType: 'firewall', parseType: 'kv', fieldMapping: '', valueTransform: '', sampleLog: 'src_ip=10.0.0.5 severity=high action=deny', enabled: true, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 4, name: '分隔符解析', deviceType: 'server', parseType: 'delimiter', delimiter: '|', fieldMapping: '', valueTransform: '', sampleLog: '10.0.0.8|high|42|denied', enabled: false, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
]

const mockFilterPolicies = [
  { id: 1, name: '丢弃调试日志', priority: 100, enabled: true, action: 'drop', conditionLogic: 'AND', conditions: JSON.stringify([{ field: 'severity', operator: 'equals', value: 'debug' }]), whitelistEnabled: false, dedupEnabled: false, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 2, name: '保留告警日志', priority: 50, enabled: true, action: 'keep', conditionLogic: 'OR', conditions: JSON.stringify([{ field: 'severity', operator: 'in', value: 'high,critical' }]), whitelistEnabled: false, dedupEnabled: false, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 3, name: '去重策略', priority: 200, enabled: true, action: 'keep', conditionLogic: 'AND', conditions: '[]', whitelistEnabled: false, dedupEnabled: true, dedupWindow: 300, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
]

const mockPushConfigs = [
  { id: 1, name: '企业微信推送', type: 'http', enabled: true, url: 'https://qyapi.weixin.qq.com/cgi-bin/webhook/send', method: 'POST', timeout: 10, retryCount: 3, retryDelay: 5, contentType: 'application/json', createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 2, name: '邮件告警', type: 'email', enabled: true, smtpHost: 'smtp.example.com', smtpPort: 587, fromAddress: 'alert@logcat.local', toAddresses: 'admin@logcat.local,ops@logcat.local', subjectTemplate: '[logcat告警]', createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 3, name: 'Syslog转发', type: 'syslog', enabled: false, syslogHost: '10.0.0.100', syslogPort: 514, syslogProtocol: 'udp', createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
]

const mockAlertRules = [
  { id: 1, name: '高危操作告警', description: '检测高危级别安全事件', status: 1, severity: 'critical', condition: { field: 'severity', operator: 'eq', value: 'critical', threshold: 1, windowSeconds: 60 }, pushConfigIds: [1, 2], cooldownSeconds: 300, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 2, name: '高频登录失败', description: '检测暴力破解攻击', status: 1, severity: 'high', condition: { field: 'event_type', operator: 'eq', value: 'login_failed', threshold: 10, windowSeconds: 300 }, pushConfigIds: [1], cooldownSeconds: 600, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 3, name: '异常流量告警', description: '检测异常网络流量', status: 0, severity: 'medium', condition: { field: 'traffic', operator: 'gte', value: '1000', threshold: 5, windowSeconds: 60 }, pushConfigIds: [2], cooldownSeconds: 180, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
]

const mockLogs = Array.from({ length: 20 }, (_, i) => ({
  id: i + 1,
  receivedAt: new Date(Date.now() - i * 600000).toISOString(),
  deviceName: mockDevices[i % mockDevices.length].name,
  deviceHost: mockDevices[i % mockDevices.length].host,
  sourceIp: `10.0.${Math.floor(i / 5)}.${(i * 7) % 256}`,
  destIp: `192.168.${Math.floor(i / 8)}.${(i * 3) % 256}`,
  facility: 'local0',
  severity: ['critical', 'high', 'medium', 'low', 'info'][i % 5],
  eventType: ['login', 'firewall', 'ids', 'vpn', 'system'][i % 5],
  message: `安全事件 #${i + 1}: ${['登录失败', '连接被拒绝', '入侵检测告警', 'VPN连接断开', '系统异常'][i % 5]}`,
  rawMessage: `<${134 + i % 5}>1 ${new Date(Date.now() - i * 600000).toISOString()} ${mockDevices[i % mockDevices.length].host} app - - - Raw message ${i + 1}`,
  parsedData: { source_ip: `10.0.${Math.floor(i / 5)}.${(i * 7) % 256}`, action: ['allow', 'deny', 'drop', 'reset', 'pass'][i % 5] },
  pushStatus: ['success', 'success', 'failed', 'pending', 'skipped'][i % 5],
  pushError: i % 5 === 2 ? 'Connection timeout' : '',
  filterStatus: i % 3 === 0 ? 'matched' : 'unmatched',
  matchedFilterPolicyId: i % 3 === 0 ? 1 : null,
  alertStatus: i % 4 === 0 ? 'triggered' : 'none',
  alertRuleId: i % 4 === 0 ? 1 : null,
  traceId: `trace-${String(i + 1).padStart(8, '0')}`,
}))

const mockAlertRecords = Array.from({ length: 10 }, (_, i) => ({
  id: i + 1,
  ruleId: mockAlertRules[i % mockAlertRules.length].id,
  ruleName: mockAlertRules[i % mockAlertRules.length].name,
  severity: mockAlertRules[i % mockAlertRules.length].severity,
  message: `告警 #${i + 1}: ${mockAlertRules[i % mockAlertRules.length].description}`,
  deviceName: mockDevices[i % mockDevices.length].name,
  sourceIp: `10.0.${Math.floor(i / 3)}.${(i * 7) % 256}`,
  logId: i + 1,
  channel: ['http', 'email', 'syslog'][i % 3],
  status: ['sent', 'failed', 'acknowledged', 'resolved'][i % 4],
  error: i % 4 === 1 ? '推送超时' : '',
  dispositionStatus: i % 4 === 2 ? 'confirmed' : i % 4 === 3 ? 'resolved' : '',
  disposition: null,
  createdAt: new Date(Date.now() - i * 3600000).toISOString(),
}))

const mockAggregatedAlerts = Array.from({ length: 5 }, (_, i) => ({
  id: i + 1,
  ruleId: mockAlertRules[i % mockAlertRules.length].id,
  ruleName: mockAlertRules[i % mockAlertRules.length].name,
  severity: mockAlertRules[i % mockAlertRules.length].severity,
  summary: `聚合告警 #${i + 1}: 过去1小时内${mockAlertRules[i % mockAlertRules.length].name}触发多次`,
  count: (i + 1) * 3,
  firstAt: new Date(Date.now() - (i + 1) * 3600000).toISOString(),
  lastAt: new Date(Date.now() - i * 600000).toISOString(),
  sourceIps: [`10.0.${i}.1`, `10.0.${i}.2`, `10.0.${i}.3`],
  deviceNames: [mockDevices[i].name, mockDevices[(i + 1) % mockDevices.length].name],
  status: i === 0 ? 'active' : i === 1 ? 'acknowledged' : 'resolved',
  logIds: mockLogs.slice(i * 2, i * 2 + 3).map(l => l.id),
  createdAt: new Date(Date.now() - (i + 1) * 3600000).toISOString(),
  updatedAt: new Date(Date.now() - i * 600000).toISOString(),
}))

const mockHighFreqIps = Array.from({ length: 8 }, (_, i) => ({
  id: i + 1,
  ip: `10.0.${Math.floor(i / 3)}.${(i * 37) % 256}`,
  count: (i + 1) * 150 + Math.floor(Math.random() * 100),
  firstSeen: new Date(Date.now() - (i + 1) * 7200000).toISOString(),
  lastSeen: new Date(Date.now() - i * 600000).toISOString(),
  deviceNames: [mockDevices[i % mockDevices.length].name],
}))

const mockDesensitizeRules = [
  { id: 1, name: 'IP地址脱敏', description: '对IP地址进行部分遮蔽', field: 'source_ip', pattern: '(\\d+\\.\\d+)\\.\\d+\\.\\d+', replacement: '$1.*.*', status: 1, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 2, name: '手机号脱敏', description: '对手机号中间4位脱敏', field: 'phone', pattern: '(\\d{3})\\d{4}(\\d{4})', replacement: '$1****$2', status: 1, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 3, name: '邮箱脱敏', description: '对邮箱地址脱敏', field: 'email', pattern: '(\\w{2})\\w+(@\\w+\\.\\w+)', replacement: '$1***$2', status: 0, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
]

const mockUsers = [
  { id: 1, username: 'admin', displayName: '管理员', email: 'admin@logcat.local', phone: '13800138000', status: 'enabled', isAdmin: true, mustChangePassword: false, lastLoginAt: '2026-05-30T10:00:00Z', createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z', roles: ['admin'] },
  { id: 2, username: 'operator', displayName: '运维人员', email: 'ops@logcat.local', phone: '13900139000', status: 'enabled', isAdmin: false, mustChangePassword: false, lastLoginAt: '2026-05-29T16:00:00Z', createdAt: '2026-02-01T00:00:00Z', updatedAt: '2026-05-29T16:00:00Z', roles: ['operator'] },
  { id: 3, username: 'viewer', displayName: '只读用户', email: 'viewer@logcat.local', phone: '', status: 'enabled', isAdmin: false, mustChangePassword: false, lastLoginAt: '2026-05-28T09:00:00Z', createdAt: '2026-03-01T00:00:00Z', updatedAt: '2026-05-28T09:00:00Z', roles: ['viewer'] },
]

const mockRoles = [
  { id: 1, name: '管理员', code: 'admin', description: '系统管理员，拥有所有权限', status: 1, builtIn: true, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z', permissions: [] },
  { id: 2, name: '运维人员', code: 'operator', description: '运维人员，拥有设备管理和日志查看权限', status: 1, builtIn: false, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z', permissions: [] },
  { id: 3, name: '只读用户', code: 'viewer', description: '只读用户，只能查看数据', status: 1, builtIn: false, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z', permissions: [] },
]

const mockPermissions = [
  { id: 1, name: '查看仪表盘', code: 'dashboard:view', type: 'menu', resource: 'dashboard', action: 'view', group: 'dashboard', description: '查看仪表盘页面' },
  { id: 2, name: '查看系统状态', code: 'system:status', type: 'menu', resource: 'system', action: 'status', group: 'system', description: '查看系统运行状态' },
  { id: 3, name: '读取系统配置', code: 'system:config:read', type: 'menu', resource: 'system', action: 'config:read', group: 'system', description: '读取系统配置信息' },
  { id: 4, name: '修改系统配置', code: 'system:config:write', type: 'button', resource: 'system', action: 'config:write', group: 'system', description: '修改系统配置' },
  { id: 5, name: '查看用户列表', code: 'users:list', type: 'menu', resource: 'users', action: 'list', group: 'users', description: '查看用户列表' },
  { id: 6, name: '创建用户', code: 'users:create', type: 'button', resource: 'users', action: 'create', group: 'users', description: '创建新用户' },
  { id: 7, name: '编辑用户', code: 'users:update', type: 'button', resource: 'users', action: 'update', group: 'users', description: '编辑用户信息' },
  { id: 8, name: '删除用户', code: 'users:delete', type: 'button', resource: 'users', action: 'delete', group: 'users', description: '删除用户' },
  { id: 9, name: '查看角色列表', code: 'roles:list', type: 'menu', resource: 'roles', action: 'list', group: 'roles', description: '查看角色列表' },
  { id: 10, name: '创建角色', code: 'roles:create', type: 'button', resource: 'roles', action: 'create', group: 'roles', description: '创建新角色' },
  { id: 11, name: '编辑角色', code: 'roles:update', type: 'button', resource: 'roles', action: 'update', group: 'roles', description: '编辑角色信息' },
  { id: 12, name: '删除角色', code: 'roles:delete', type: 'button', resource: 'roles', action: 'delete', group: 'roles', description: '删除角色' },
  { id: 13, name: '查看设备列表', code: 'devices:list', type: 'menu', resource: 'devices', action: 'list', group: 'devices', description: '查看设备列表' },
  { id: 14, name: '创建设备', code: 'devices:create', type: 'button', resource: 'devices', action: 'create', group: 'devices', description: '创建新设备' },
  { id: 15, name: '编辑设备', code: 'devices:update', type: 'button', resource: 'devices', action: 'update', group: 'devices', description: '编辑设备信息' },
  { id: 16, name: '删除设备', code: 'devices:delete', type: 'button', resource: 'devices', action: 'delete', group: 'devices', description: '删除设备' },
  { id: 17, name: '查看设备组列表', code: 'device-groups:list', type: 'menu', resource: 'device-groups', action: 'list', group: 'device-groups', description: '查看设备组列表' },
  { id: 18, name: '查看设备模板列表', code: 'device-templates:list', type: 'menu', resource: 'device-templates', action: 'list', group: 'device-templates', description: '查看设备模板列表' },
  { id: 19, name: '查看字段映射列表', code: 'field-mappings:list', type: 'menu', resource: 'field-mappings', action: 'list', group: 'field-mappings', description: '查看字段映射列表' },
  { id: 20, name: '查看解析模板列表', code: 'parse-templates:list', type: 'menu', resource: 'parse-templates', action: 'list', group: 'parse-templates', description: '查看解析模板列表' },
  { id: 21, name: '查看过滤策略列表', code: 'filter-policies:list', type: 'menu', resource: 'filter-policies', action: 'list', group: 'filter-policies', description: '查看过滤策略列表' },
  { id: 22, name: '查看输出模板列表', code: 'output-templates:list', type: 'menu', resource: 'output-templates', action: 'list', group: 'output-templates', description: '查看输出模板列表' },
  { id: 23, name: '查看推送配置列表', code: 'push-configs:list', type: 'menu', resource: 'push-configs', action: 'list', group: 'push-configs', description: '查看推送配置列表' },
  { id: 24, name: '查看告警规则列表', code: 'alert-rules:list', type: 'menu', resource: 'alert-rules', action: 'list', group: 'alert-rules', description: '查看告警规则列表' },
  { id: 25, name: '查看日志列表', code: 'logs:list', type: 'menu', resource: 'logs', action: 'list', group: 'logs', description: '查看日志列表' },
  { id: 26, name: '日志追踪', code: 'logs:trace', type: 'button', resource: 'logs', action: 'trace', group: 'logs', description: '查看日志追踪详情' },
  { id: 27, name: '导出日志', code: 'logs:export', type: 'button', resource: 'logs', action: 'export', group: 'logs', description: '导出日志数据' },
  { id: 28, name: '查看告警列表', code: 'alerts:list', type: 'menu', resource: 'alerts', action: 'list', group: 'alerts', description: '查看告警列表' },
  { id: 29, name: '查看告警处置', code: 'alerts:disposition:list', type: 'menu', resource: 'alerts', action: 'disposition:list', group: 'alerts', description: '查看告警处置记录' },
  { id: 30, name: '创建告警处置', code: 'alerts:disposition:create', type: 'button', resource: 'alerts', action: 'disposition:create', group: 'alerts', description: '处置告警' },
  { id: 31, name: '查看聚合告警列表', code: 'aggregated-alerts:list', type: 'menu', resource: 'aggregated-alerts', action: 'list', group: 'aggregated-alerts', description: '查看聚合告警列表' },
  { id: 32, name: '确认聚合告警', code: 'aggregated-alerts:acknowledge', type: 'button', resource: 'aggregated-alerts', action: 'acknowledge', group: 'aggregated-alerts', description: '确认聚合告警' },
  { id: 33, name: '解决聚合告警', code: 'aggregated-alerts:resolve', type: 'button', resource: 'aggregated-alerts', action: 'resolve', group: 'aggregated-alerts', description: '解决聚合告警' },
  { id: 34, name: '查看高频IP列表', code: 'high-freq-ips:list', type: 'menu', resource: 'high-freq-ips', action: 'list', group: 'high-freq-ips', description: '查看高频IP列表' },
  { id: 35, name: '配置高频IP', code: 'high-freq-ips:config', type: 'button', resource: 'high-freq-ips', action: 'config', group: 'high-freq-ips', description: '配置高频IP阈值' },
  { id: 36, name: '查看脱敏规则列表', code: 'desensitize-rules:list', type: 'menu', resource: 'desensitize-rules', action: 'list', group: 'desensitize-rules', description: '查看脱敏规则列表' },
  { id: 37, name: '查询统计', code: 'stats:query', type: 'menu', resource: 'stats', action: 'query', group: 'stats', description: '查询统计数据' },
  { id: 38, name: '查看统计字段', code: 'stats:fields', type: 'menu', resource: 'stats', action: 'fields', group: 'stats', description: '查看可用统计字段' },
  { id: 39, name: '导出统计', code: 'stats:export', type: 'button', resource: 'stats', action: 'export', group: 'stats', description: '导出统计数据' },
  { id: 40, name: '导出配置', code: 'export:config', type: 'button', resource: 'export', action: 'config', group: 'export', description: '导出系统配置' },
  { id: 41, name: '导入配置', code: 'import:config', type: 'button', resource: 'import', action: 'config', group: 'import', description: '导入系统配置' },
  { id: 42, name: '查看审计日志', code: 'audit-logs:list', type: 'menu', resource: 'audit-logs', action: 'list', group: 'audit-logs', description: '查看审计日志列表' },
]

const mockAuditLogs = Array.from({ length: 15 }, (_, i) => ({
  id: i + 1,
  userId: (i % 3) + 1,
  username: mockUsers[i % 3].username,
  action: ['login', 'create', 'update', 'delete', 'export'][i % 5],
  resource: ['device', 'filter_policy', 'push_config', 'alert_rule', 'user'][i % 5],
  resourceType: ['device', 'filter_policy', 'push_config', 'alert_rule', 'user'][i % 5],
  resourceId: String((i % 5) + 1),
  detail: `${['登录系统', '创建设备', '更新过滤策略', '删除推送配置', '导出告警规则'][i % 5]}`,
  result: i % 7 === 0 ? 'failure' : 'success',
  ip: `192.168.1.${100 + i}`,
  clientIp: `192.168.1.${100 + i}`,
  createdAt: new Date(Date.now() - i * 1800000).toISOString(),
}))

const mockFieldMappings = [
  { id: 1, name: 'FortiGate字段映射', description: 'FortiGate防火墙日志字段映射', mappings: [{ sourceField: 'src', targetField: 'source_ip', transformRule: '' }, { sourceField: 'dst', targetField: 'dest_ip', transformRule: '' }, { sourceField: 'action', targetField: 'action', transformRule: 'lowercase' }], createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 2, name: '通用Syslog映射', description: '通用Syslog日志字段映射', mappings: [{ sourceField: 'host', targetField: 'device_host', transformRule: '' }, { sourceField: 'msg', targetField: 'message', transformRule: '' }], createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
]

const mockOutputTemplates = [
  { id: 1, name: 'JSON标准输出', description: '标准JSON格式输出', format: 'json', template: '{"time":"{{.ReceivedAt}}","device":"{{.DeviceName}}","message":"{{.Message}}"}', createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
  { id: 2, name: 'CEF格式输出', description: 'Common Event Format', format: 'custom', template: 'CEF:0|logcat|Syslog|1.0|{{.EventType}}|{{.Message}}|{{.Severity}}', createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-05-30T10:00:00Z' },
]

const mockDashboardStats = {
  serviceStatus: 'running',
  receiveRate: 1250.5,
  todayTotal: 895432,
  todayAlerts: 47,
  pushSuccessRate: 0.985,
  tcpConnections: 156,
  lastReceivedAt: '2026-05-30T10:30:00Z',
  parseSuccessRate: 0.978,
  queueBacklog: [
    { name: 'raw', size: 234, capacity: 10000, enqueueRate: 1250.5, dequeueRate: 1248.2 },
    { name: 'parsed', size: 56, capacity: 5000, enqueueRate: 1248.2, dequeueRate: 1247.8 },
    { name: 'push', size: 12, capacity: 2000, enqueueRate: 45.2, dequeueRate: 45.0 },
  ],
  recentPushFailures: [
    { time: '2026-05-30T10:25:00Z', channel: 'http', error: 'Connection timeout', logId: 42 },
    { time: '2026-05-30T10:15:00Z', channel: 'email', error: 'SMTP auth failed', logId: 38 },
  ],
  healthStatus: {
    service: 'running',
    cpu: 35.2,
    memory: 48.7,
    diskUsage: 62.1,
    networkIn: 524288,
    networkOut: 262144,
  },
}

const mockSystemStatus = {
  serviceRunning: true,
  startedAt: '2026-05-29T00:00:00Z',
  uptime: '1d 10h 30m',
  uptimeSeconds: 123456,
  database: 'SQLite',
  databaseStatus: 'connected',
  goVersion: '1.22',
  listeners: [
    { protocol: 'TCP', port: 514, address: '0.0.0.0:514', status: 'running' },
    { protocol: 'UDP', port: 514, address: '0.0.0.0:514', status: 'running' },
    { protocol: 'TLS', port: 6514, address: '0.0.0.0:6514', status: 'running' },
  ],
  connections: { total: 156, active: 89, idle: 67, closed: 12 },
  queue: { name: 'main', size: 302, capacity: 17000, enqueueRate: 1250.5, dequeueRate: 1248.2 },
  workers: [
    { id: 1, stage: 'receive', count: 10, status: 'running', processedCount: 895432, errorCount: 23, lastActiveAt: '2026-05-30T10:30:00Z' },
    { id: 2, stage: 'parse', count: 5, status: 'running', processedCount: 892100, errorCount: 3332, lastActiveAt: '2026-05-30T10:30:00Z' },
    { id: 3, stage: 'push', count: 3, status: 'running', processedCount: 45200, errorCount: 678, lastActiveAt: '2026-05-30T10:30:00Z' },
  ],
  receiverMetrics: {
    udpReceived: 500000,
    tcpReceived: 400000,
    udpErrors: 12,
    tcpErrors: 5,
    parseErrors: 3332,
    channelDropped: 0,
    tcpConnections: 156,
    lastReceiveAt: '2026-05-30T10:30:00Z',
  },
}

const mockSystemConfigs = [
  { id: 1, configKey: 'serverHost', configValue: '0.0.0.0', description: 'Web/API 服务监听地址', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 2, configKey: 'serverPort', configValue: '5080', description: '管理端和 API 默认访问端口', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 3, configKey: 'syslogEnabled', configValue: 'false', description: '控制 Syslog 接收器默认启用状态', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 4, configKey: 'syslogUdpPort', configValue: '5140', description: 'UDP 日志接收端口', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 5, configKey: 'syslogTcpPort', configValue: '5140', description: 'TCP 日志接收端口', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 6, configKey: 'sessionExpireHours', configValue: '24', description: 'HttpOnly Session 默认过期时间', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 7, configKey: 'maxFailedLogin', configValue: '5', description: '连续失败次数达到阈值后锁定账号', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 8, configKey: 'lockDurationMinutes', configValue: '30', description: '登录失败达到阈值后的锁定时长', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 9, configKey: 'retentionDays', configValue: '90', description: '已匹配日志默认保留天数', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 10, configKey: 'unmatchedRetentionDays', configValue: '30', description: '未命中模板日志默认保留天数', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 11, configKey: 'maxLogSize', configValue: '10000', description: '日志文件或数据量的上限控制', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 12, configKey: 'defaultFilterAction', configValue: 'keep', description: '日志未命中任何过滤规则时的默认处理动作', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 13, configKey: 'parseWorkers', configValue: '4', description: '解析阶段 worker 数量', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 14, configKey: 'filterWorkers', configValue: '4', description: '筛选阶段 worker 数量', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 15, configKey: 'pushWorkers', configValue: '4', description: '推送阶段 worker 数量', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 16, configKey: 'queueCapacity', configValue: '10000', description: '流水线队列默认容量', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 17, configKey: 'queueFullPolicy', configValue: 'block_drop', description: '队列满时的处理策略', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 18, configKey: 'databaseType', configValue: 'sqlite', description: '当前支持 sqlite 和 mysql', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 19, configKey: 'sqlitePath', configValue: 'data/logcat.db', description: 'SQLite 数据文件路径', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 20, configKey: 'configDir', configValue: 'data/config', description: '系统配置文件存储目录路径', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 21, configKey: 'logDir', configValue: 'data/logs', description: '系统运行日志存储目录路径', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 22, configKey: 'mysqlHost', configValue: '127.0.0.1', description: 'MySQL 服务连接地址', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 23, configKey: 'mysqlPort', configValue: '3306', description: 'MySQL 服务连接端口', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 24, configKey: 'mysqlDatabase', configValue: 'logcat', description: 'MySQL 数据库名称', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 25, configKey: 'mysqlUsername', configValue: 'root', description: 'MySQL 连接用户名', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 26, configKey: 'mysqlPassword', configValue: '', description: 'MySQL 连接密码', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 27, configKey: 'mysqlCharset', configValue: 'utf8mb4', description: 'MySQL 连接字符集', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 28, configKey: 'mysqlTimezone', configValue: 'Asia/Shanghai', description: 'MySQL 连接时区', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 29, configKey: 'mysqlMaxOpenConns', configValue: '50', description: 'MySQL 最大打开连接数', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 30, configKey: 'mysqlMaxIdleConns', configValue: '10', description: 'MySQL 最大空闲连接数', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 31, configKey: 'alertEnabled', configValue: 'true', description: '全局告警启用/禁用', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 32, configKey: 'alertInterval', configValue: '60', description: '同一告警规则的最小触发间隔', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 33, configKey: 'unmatchedAlertEnabled', configValue: 'false', description: '当日志未命中任何模板时是否触发告警', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 34, configKey: 'theme', configValue: 'light', description: '界面主题切换', updatedAt: '2026-01-01T00:00:00Z' },
]

const mockRuntimeMetrics = {
  uptime_seconds: 123456,
  status: 'running',
  goroutines: 42,
  heap_alloc_mb: 128.5,
  heap_sys_mb: 256,
  heap_total_mb: 192,
  num_gc: 156,
  receiver: {
    udpReceived: 500000,
    tcpReceived: 400000,
    udpErrors: 12,
    tcpErrors: 5,
    parseErrors: 3332,
    channelDropped: 0,
    tcpConnections: 156,
  },
  pipeline: {
    rawQueueDepth: 234,
    parsedQueueDepth: 56,
    dbQueueDepth: 12,
    pushQueueDepth: 8,
    parseProcessed: 892100,
    parseErrors: 3332,
    filterProcessed: 890000,
    filterDropped: 5432,
    dbWritten: 895432,
    dbErrors: 0,
    pushProcessed: 45200,
    pushErrors: 678,
    rawDropped: 0,
    dbDropped: 0,
    pushDropped: 0,
  },
}

const mockMetricsSnapshot = {
  status: 'running',
  database: 'SQLite',
  uptime: 123456,
  started_at: '2026-05-29T00:00:00Z',
  db_stats: {
    open_connections: 5,
    in_use: 2,
    idle: 3,
    wait_count: 0,
    max_open: 10,
  },
  receiver_metrics: {
    udpReceived: 500000,
    tcpReceived: 400000,
    udpErrors: 12,
    tcpErrors: 5,
    parseErrors: 3332,
    channelDropped: 0,
    tcpConnections: 156,
  },
  pipeline_metrics: {
    rawQueueDepth: 234,
    parsedQueueDepth: 56,
    dbQueueDepth: 12,
    pushQueueDepth: 8,
    parseProcessed: 892100,
    filterProcessed: 890000,
    dbWritten: 895432,
    pushProcessed: 45200,
    pushErrors: 678,
  },
}

function paginate(list: any[], params: any) {
  const page = params?.page || 1
  const pageSize = params?.pageSize || 20
  const start = (page - 1) * pageSize
  const end = start + pageSize
  return { list: list.slice(start, end), total: list.length, page, pageSize }
}

function ok(data: any): ApiResponse<any> {
  return { code: 0, message: 'success', data }
}

type RouteHandler = (params: any, body: any) => any
type RouteEntry = [string | null, RegExp, RouteHandler]

const routes: RouteEntry[] = [
  // ==================== Auth ====================
  ['GET', /^\/api\/auth\/init-status$/, () => ok({ initialized: true })],
  ['POST', /^\/api\/auth\/login$/, () => ok({ user: mockUser })],
  ['POST', /^\/api\/auth\/logout$/, () => ok(null)],
  ['GET', /^\/api\/auth\/me$/, () => ok(mockUser)],
  ['POST', /^\/api\/auth\/change-password$/, () => ok(null)],
  ['POST', /^\/api\/auth\/init-admin$/, () => ok(null)],

  // ==================== Dashboard ====================
  ['GET', /^\/api\/dashboard\/stats$/, () => { mockDashboardStats.serviceStatus = mockSystemStatus.serviceRunning ? 'running' : 'stopped'; return ok(mockDashboardStats) }],

  // ==================== Monitoring ====================
  ['GET', /^\/api\/healthz$/, () => ok({ status: 'ok', time: new Date().toISOString() })],
  ['GET', /^\/api\/readyz$/, () => ok({ status: 'ready', time: new Date().toISOString() })],
  ['GET', /^\/api\/metrics\/runtime$/, () => ok(mockRuntimeMetrics)],
  ['GET', /^\/api\/metrics$/, () => ok(mockMetricsSnapshot)],

  // ==================== Permissions ====================
  ['GET', /^\/api\/permissions$/, () => ok(mockPermissions)],

  // ==================== System ====================
  ['GET', /^\/api\/system\/status$/, () => ok(mockSystemStatus)],
  ['GET', /^\/api\/system\/health$/, () => ok({ status: 'ok', time: new Date().toISOString() })],
  ['GET', /^\/api\/system\/config$/, () => ok(mockSystemConfigs)],
  ['PUT', /^\/api\/system\/config$/, () => ok(null)],
  ['POST', /^\/api\/system\/syslog\/start$/, () => { mockSystemStatus.serviceRunning = true; mockSystemStatus.listeners.forEach((l: any) => l.status = 'running'); mockSystemStatus.workers.forEach((w: any) => w.status = 'running'); return ok(null) }],
  ['POST', /^\/api\/system\/syslog\/stop$/, () => { mockSystemStatus.serviceRunning = false; mockSystemStatus.listeners.forEach((l: any) => l.status = 'stopped'); mockSystemStatus.workers.forEach((w: any) => w.status = 'stopped'); return ok(null) }],

  // ==================== Devices ====================
  ['GET', /^\/api\/devices\/all$/, () => ok(mockDevices)],
  ['GET', /^\/api\/devices\/\d+$/, () => ok(mockDevices[0])],
  ['POST', /^\/api\/devices$/, () => ok(mockDevices[0])],
  ['PUT', /^\/api\/devices\/\d+$/, () => ok(mockDevices[0])],
  ['DELETE', /^\/api\/devices\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/devices$/, (p) => ok(paginate(mockDevices, p))],

  // ==================== Device Groups ====================
  ['GET', /^\/api\/device-groups\/all$/, () => ok(mockDeviceGroups)],
  ['GET', /^\/api\/device-groups\/\d+$/, () => ok(mockDeviceGroups[0])],
  ['POST', /^\/api\/device-groups$/, () => ok(mockDeviceGroups[0])],
  ['PUT', /^\/api\/device-groups\/\d+$/, () => ok(mockDeviceGroups[0])],
  ['DELETE', /^\/api\/device-groups\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/device-groups$/, (p) => ok(paginate(mockDeviceGroups, p))],

  // ==================== Device Templates ====================
  ['GET', /^\/api\/device-templates\/all$/, () => ok(mockDeviceTemplates)],
  ['GET', /^\/api\/device-templates\/\d+$/, () => ok(mockDeviceTemplates[0])],
  ['POST', /^\/api\/device-templates$/, () => ok(mockDeviceTemplates[0])],
  ['PUT', /^\/api\/device-templates\/\d+$/, () => ok(mockDeviceTemplates[0])],
  ['DELETE', /^\/api\/device-templates\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/device-templates$/, (p) => ok(paginate(mockDeviceTemplates, p))],

  // ==================== Field Mappings ====================
  ['GET', /^\/api\/field-mappings\/\d+$/, () => ok(mockFieldMappings[0])],
  ['POST', /^\/api\/field-mappings$/, () => ok(mockFieldMappings[0])],
  ['PUT', /^\/api\/field-mappings\/\d+$/, () => ok(mockFieldMappings[0])],
  ['DELETE', /^\/api\/field-mappings\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/field-mappings$/, (p) => ok(paginate(mockFieldMappings, p))],

  // ==================== Parse Templates ====================
  ['POST', /^\/api\/parse-templates\/test$/, () => ok({ success: true, fields: { source_ip: '10.0.0.1', action: 'deny', severity: 'high' } })],
  ['GET', /^\/api\/parse-templates\/\d+$/, () => ok(mockParseTemplates[0])],
  ['POST', /^\/api\/parse-templates$/, () => ok(mockParseTemplates[0])],
  ['PUT', /^\/api\/parse-templates\/\d+$/, () => ok(mockParseTemplates[0])],
  ['DELETE', /^\/api\/parse-templates\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/parse-templates$/, (p) => ok(paginate(mockParseTemplates, p))],

  // ==================== Filter Policies ====================
  ['POST', /^\/api\/filter-policies\/test$/, () => ok({ matched: true, action: 'keep', message: '匹配成功', whitelistResult: '', policy: { name: '保留告警日志' } })],
  ['GET', /^\/api\/filter-policies\/\d+$/, () => ok(mockFilterPolicies[0])],
  ['POST', /^\/api\/filter-policies$/, () => ok(mockFilterPolicies[0])],
  ['PUT', /^\/api\/filter-policies\/\d+$/, () => ok(mockFilterPolicies[0])],
  ['DELETE', /^\/api\/filter-policies\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/filter-policies$/, (p) => ok(paginate(mockFilterPolicies, p))],

  // ==================== Output Templates ====================
  ['GET', /^\/api\/output-templates\/\d+$/, () => ok(mockOutputTemplates[0])],
  ['POST', /^\/api\/output-templates$/, () => ok(mockOutputTemplates[0])],
  ['PUT', /^\/api\/output-templates\/\d+$/, () => ok(mockOutputTemplates[0])],
  ['DELETE', /^\/api\/output-templates\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/output-templates$/, (p) => ok(paginate(mockOutputTemplates, p))],

  // ==================== Push Configs ====================
  ['POST', /^\/api\/push-configs\/\d+\/test$/, () => ok({ success: true, channel: 'http', statusCode: 200, responseBody: '{"errcode":0}', summary: '推送成功' })],
  ['GET', /^\/api\/push-configs\/\d+$/, () => ok(mockPushConfigs[0])],
  ['POST', /^\/api\/push-configs$/, () => ok(mockPushConfigs[0])],
  ['PUT', /^\/api\/push-configs\/\d+$/, () => ok(mockPushConfigs[0])],
  ['DELETE', /^\/api\/push-configs\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/push-configs$/, (p) => ok(paginate(mockPushConfigs, p))],

  // ==================== Alert Rules ====================
  ['PUT', /^\/api\/alert-rules\/\d+\/status$/, () => ok(null)],
  ['GET', /^\/api\/alert-rules\/\d+$/, () => ok(mockAlertRules[0])],
  ['POST', /^\/api\/alert-rules$/, () => ok(mockAlertRules[0])],
  ['PUT', /^\/api\/alert-rules\/\d+$/, () => ok(mockAlertRules[0])],
  ['DELETE', /^\/api\/alert-rules\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/alert-rules$/, (p) => ok(paginate(mockAlertRules, p))],

  // ==================== Logs ====================
  ['GET', /^\/api\/logs\/\d+\/trace$/, () => ok({ log: mockLogs[0], trace: [{ stage: '接收', status: 'success', detail: '防火墙-001 (10.0.0.1)', time: mockLogs[0].receivedAt }, { stage: '解析', status: 'success', detail: '解析完成', time: mockLogs[0].receivedAt }, { stage: '过滤', status: 'success', detail: '通过过滤策略', time: mockLogs[0].receivedAt }, { stage: '推送', status: 'success', detail: 'success', time: mockLogs[0].receivedAt }] })],
  ['GET', /^\/api\/logs\/unmatched-count$/, () => ok({ count: 42 })],
  ['DELETE', /^\/api\/logs\/cleanup$/, () => ok({ deleted: 0 })],
  ['GET', /^\/api\/logs\/export$/, () => ok({ url: '#' })],
  ['GET', /^\/api\/logs\/clean$/, () => ok(null)],
  ['GET', /^\/api\/logs$/, (p) => ok(paginate(mockLogs, p))],

  // ==================== Alerts ====================
  ['GET', /^\/api\/alerts\/\d+\/dispositions$/, () => ok([])],
  ['GET', /^\/api\/alerts\/\d+\/disposition$/, () => ok(null)],
  ['POST', /^\/api\/alerts\/\d+\/disposition$/, () => ok(null)],
  ['GET', /^\/api\/alerts\/\d+$/, () => ok(mockAlertRecords[0])],
  ['GET', /^\/api\/alerts$/, (p) => ok(paginate(mockAlertRecords, p))],

  // ==================== Alert Dispositions ====================
  ['GET', /^\/api\/alert-dispositions$/, (p) => ok(paginate(mockAlertRecords, p))],

  // ==================== Aggregated Alerts ====================
  ['POST', /^\/api\/aggregated-alerts\/\d+\/resolve$/, () => ok(null)],
  ['POST', /^\/api\/aggregated-alerts\/\d+\/acknowledge$/, () => ok(null)],
  ['GET', /^\/api\/aggregated-alerts\/\d+\/logs$/, () => ok(mockLogs.slice(0, 3))],
  ['GET', /^\/api\/aggregated-alerts\/\d+$/, () => ok(mockAggregatedAlerts[0])],
  ['GET', /^\/api\/aggregated-alerts$/, (p) => ok(paginate(mockAggregatedAlerts, p))],

  // ==================== High Frequency IPs ====================
  ['GET', /^\/api\/high-frequency-ips\/config$/, () => ok({ timeWindowSeconds: 300, threshold: 100 })],
  ['PUT', /^\/api\/high-frequency-ips\/config$/, () => ok(null)],
  ['POST', /^\/api\/high-frequency-ips\/refresh$/, () => ok(null)],
  ['GET', /^\/api\/high-frequency-ips$/, (p) => ok(paginate(mockHighFreqIps, p))],

  // ==================== Desensitize Rules ====================
  ['PUT', /^\/api\/desensitize-rules\/\d+\/status$/, () => ok(null)],
  ['GET', /^\/api\/desensitize-rules\/\d+$/, () => ok(mockDesensitizeRules[0])],
  ['POST', /^\/api\/desensitize-rules$/, () => ok(mockDesensitizeRules[0])],
  ['PUT', /^\/api\/desensitize-rules\/\d+$/, () => ok(mockDesensitizeRules[0])],
  ['DELETE', /^\/api\/desensitize-rules\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/desensitize-rules$/, (p) => ok(paginate(mockDesensitizeRules, p))],

  // ==================== Users ====================
  ['POST', /^\/api\/users\/\d+\/unlock$/, () => ok(null)],
  ['POST', /^\/api\/users\/\d+\/reset-password$/, () => ok(null)],
  ['POST', /^\/api\/users\/\d+\/force-password-change$/, () => ok(null)],
  ['GET', /^\/api\/users\/\d+\/roles$/, () => ok({ roles: [] })],
  ['POST', /^\/api\/users\/\d+\/roles$/, () => ok(null)],
  ['GET', /^\/api\/users\/\d+$/, () => ok(mockUsers[0])],
  ['POST', /^\/api\/users$/, () => ok(mockUsers[0])],
  ['PUT', /^\/api\/users\/\d+$/, () => ok(mockUsers[0])],
  ['DELETE', /^\/api\/users\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/users$/, (p) => ok(paginate(mockUsers, p))],

  // ==================== Roles ====================
  ['GET', /^\/api\/roles\/\d+\/permissions$/, () => ok([])],
  ['POST', /^\/api\/roles\/\d+\/permissions$/, () => ok(null)],
  ['GET', /^\/api\/roles\/\d+$/, () => ok(mockRoles[0])],
  ['POST', /^\/api\/roles$/, () => ok(mockRoles[0])],
  ['PUT', /^\/api\/roles\/\d+$/, () => ok(mockRoles[0])],
  ['DELETE', /^\/api\/roles\/\d+$/, () => ok(null)],
  ['GET', /^\/api\/roles$/, (p) => ok(paginate(mockRoles, p))],

  // ==================== Audit Logs ====================
  ['GET', /^\/api\/audit-logs$/, (p) => ok(paginate(mockAuditLogs, p))],

  // ==================== Stats ====================
  ['GET', /^\/api\/stats\/export-csv$/, () => ok({ url: '#', count: 10 })],
  ['GET', /^\/api\/stats\/ip-list$/, () => ok({ ips: ['10.0.1.5', '10.0.2.8', '10.0.3.12'] })],
  ['GET', /^\/api\/stats\/available-fields$/, () => ok([{ value: 'source_ip', label: '源IP' }, { value: 'dest_ip', label: '目标IP' }, { value: 'severity', label: '严重程度' }, { value: 'device_name', label: '设备名称' }, { value: 'event_type', label: '事件类型' }])],
  ['GET', /^\/api\/stats\/query$/, () => ok({ results: [{ value: '10.0.1.5', count: 1250, percentage: 25.0, lastSeenAt: '2026-05-30T10:30:00Z' }, { value: '10.0.2.8', count: 980, percentage: 19.6, lastSeenAt: '2026-05-30T10:28:00Z' }, { value: '10.0.3.12', count: 750, percentage: 15.0, lastSeenAt: '2026-05-30T10:25:00Z' }], totalLogs: 5000, uniqueValues: 128 })],

  // ==================== Import/Export ====================
  ['GET', /^\/api\/export\/history$/, (p) => ok(paginate([], p))],
  ['GET', /^\/api\/export\/[\w-]+$/, () => ok({ url: '#', version: '1.0', resourceType: 'device-templates', count: 10 })],
  ['GET', /^\/api\/import\/history$/, (p) => ok(paginate([], p))],
  ['POST', /^\/api\/import\/[\w-]+$/, () => ok({ resourceType: 'device-templates', version: '1.0', created: 2, updated: 1, failed: 0, errors: [] })],
]

function matchRoute(method: string, url: string): RouteHandler | null {
  for (const [m, pattern, handler] of routes) {
    if ((m === null || m === method) && pattern.test(url)) return handler
  }
  return null
}

export function setupMockAdapter(http: any) {
  const originalAdapter = http.defaults.adapter

  http.defaults.adapter = function mockAdapter(config: any) {
    const url = config.baseURL ? config.url?.replace(config.baseURL, '') : config.url
    const fullUrl = `${config.baseURL || ''}${url}`
    const method = (config.method || 'GET').toUpperCase()
    const handler = matchRoute(method, fullUrl)

    if (handler) {
      const params = config.params || {}
      const body = config.data ? (typeof config.data === 'string' ? JSON.parse(config.data) : config.data) : {}
      const responseData = handler(params, body)

      return new Promise((resolve) => {
        setTimeout(() => {
          resolve({
            data: responseData,
            status: 200,
            statusText: 'OK',
            headers: { 'content-type': 'application/json' },
            config,
          })
        }, 50 + Math.random() * 100)
      })
    }

    return originalAdapter(config)
  }
}
