// Package constants 定义全局常量，消除散落在代码中的魔法值。
package constants

// 系统版本信息
const (
	AppName    = "logcat"
	AppVersion = "0.2.0"
	AppAuthor  = "迷人安全"
)

// 默认端口与协议
const (
	DefaultListenPort = 5140
	DefaultWebPort    = 8080
	ProtocolUDP       = "udp"
	ProtocolTCP       = "tcp"
)

// 推送平台类型
const (
	PlatformFeishu = "feishu"
	PlatformEmail  = "email"
	PlatformSyslog = "syslog"
)

// Syslog 输出格式
const (
	FormatJSON    = "json"
	FormatRFC3164 = "rfc3164"
	FormatRFC5424 = "rfc5424"
)

// 日志过滤状态
const (
	FilterStatusPending     = "pending"
	FilterStatusMatched     = "matched"
	FilterStatusUnmatched   = "unmatched"
	FilterStatusWhitelisted = "whitelisted"
	FilterStatusDropped     = "dropped"
	FilterStatusDisabled    = "disabled"
)

// 告警状态
const (
	AlertStatusNone    = "none"
	AlertStatusPending = "pending"
	AlertStatusSent    = "sent"
	AlertStatusFailed  = "failed"
)

// 筛选动作
const (
	ActionKeep    = "keep"
	ActionDiscard = "discard"
)

// 条件逻辑
const (
	LogicAnd = "AND"
	LogicOR  = "OR"
)

// 默认配置值
const (
	DefaultLogRetention          = 7
	DefaultMaxLogSize            = 524288000 // 500MB
	DefaultAlertInterval         = 60
	DefaultUnmatchedLogRetention = 7
	DefaultDedupWindow           = 60
	DefaultTheme                 = "dark"
	DefaultLanguage              = "zh-CN"
)

// 默认分隔符
const (
	DefaultDelimiter   = "|!"
	DefaultKVSeparator = ":"
)

// UDP 数据包大小上限
const MaxUDPPacketSize = 65507

// 日志清理相关阈值
const (
	LogCleanupThreshold    = 10000
	AlertCleanupThreshold  = 10000
	LogVacuumThreshold     = 50000
	LogCleanupIntervalMins = 10
)

// 告警缓存上限
const AlertCacheMaxSize = 10000

// 环境变量名
const (
	EnvDataDir       = "SYSLG_ALERT_DATA_DIR"
	EnvTemplatesDir  = "SYSLG_ALERT_TEMPLATES_DIR"
	EnvConfigDir     = "SYSLG_ALERT_CONFIG_DIR"
	EnvOpenBrowser   = "LOGCAT_OPEN_BROWSER"
	EnvAdminUsername = "LOGCAT_ADMIN_USERNAME"
	EnvAdminPassword = "LOGCAT_ADMIN_PASSWORD"
)

// 文件/目录名
const (
	DataDirName      = "data"
	DatabaseFile     = "syslog.db"
	TemplatesDirName = "templates"
)

// 转发标记
const (
	ForwardedJSONMark = `"forwarded":true`
	ForwardedTextMark = "[FORWARDED]"
)

// 默认分组
const DefaultDeviceGroupName = "默认分组"

// 支持的语言
const (
	LangZhCN = "zh-CN"
	LangEnUS = "en-US"
)

// 解析类型
const (
	ParseTypeSyslogJSON     = "syslog_json"
	ParseTypeJSON           = "json"
	ParseTypeRegex          = "regex"
	ParseTypeKV             = "kv"
	ParseTypeDelimiter      = "delimiter"
	ParseTypeKeyValue       = "keyvalue"
	ParseTypeSmartDelimiter = "smart_delimiter"
)

// 筛选条件操作符
const (
	OpEquals      = "equals"
	OpNotEquals   = "not_equals"
	OpContains    = "contains"
	OpNotContains = "not_contains"
	OpIn          = "in"
	OpNotIn       = "not_in"
	OpStartsWith  = "starts_with"
	OpEndsWith    = "ends_with"
	OpRegex       = "regex"
	OpNotRegex    = "not_regex"
	OpExists      = "exists"
	OpNotExists   = "not_exists"
	OpGT          = "gt"
	OpGTE         = "gte"
	OpLT          = "lt"
	OpLTE         = "lte"
)

// 默认 syslog 头部正则（用于 smart_delimiter 等场景的 SkipHeader）
const DefaultSyslogHeaderRegex = `<(?P<priority>[0-9]+)>(?P<timestamp>[A-Za-z]+[ ]+[0-9]+ [0-9:]+) (?P<hostname>[^ ]+) (?P<program>[^:]+):`

// IOC 告警类型与默认攻击结果
const (
	AlertTypeIOC            = "ioc_alert"
	AttackResultCompromised = "失陷"
)

// 认证相关
const (
	AuthHeaderName      = "Authorization"
	AuthTokenType       = "Bearer"
	AuthTokenBytes      = 32           // token 随机字节数（hex 编码后 64 字符）
	AuthSessionTTL      = 24 * 60 * 60 // 会话有效期（秒）：24 小时
	AuthDefaultUsername = "admin"
	AuthDefaultPassword = "admin123"
	AuthDefaultNickname = "管理员"
	AuthBCryptCost      = 10
)
