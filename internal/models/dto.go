package models

import "time"

// LogTraceInfo 日志追踪信息（全链路）
type LogTraceInfo struct {
	LogID            uint             `json:"logId"`
	ReceivedAt       time.Time        `json:"receivedAt"`
	SourceIP         string           `json:"sourceIp"`
	RawMessage       string           `json:"rawMessage"`
	ReceiveStatus    string           `json:"receiveStatus"`
	ReceiveError     string           `json:"receiveError,omitempty"`
	ParseStatus      string           `json:"parseStatus"`
	ParseTemplate    string           `json:"parseTemplate,omitempty"`
	ParsedData       string           `json:"parsedData,omitempty"`
	ParseError       string           `json:"parseError,omitempty"`
	FilterStatus     string           `json:"filterStatus"`
	FilterEnabled    bool             `json:"filterEnabled"`
	MatchedPolicy    string           `json:"matchedPolicy,omitempty"`
	FilterConditions string           `json:"filterConditions,omitempty"`
	FilterResult     string           `json:"filterResult,omitempty"`
	AlertStatus      string           `json:"alertStatus"`
	AlertRecords     []AlertTraceInfo `json:"alertRecords,omitempty"`
}

// AlertTraceInfo 告警追踪信息
type AlertTraceInfo struct {
	RobotID   uint      `json:"robotId"`
	RobotName string    `json:"robotName"`
	Platform  string    `json:"platform"`
	Status    string    `json:"status"`
	ErrorMsg  string    `json:"errorMsg,omitempty"`
	SentAt    time.Time `json:"sentAt,omitempty"`
}

// FieldStatsRequest 字段统计请求
type FieldStatsRequest struct {
	DeviceID       uint   `json:"deviceId"`
	FilterPolicyID uint   `json:"filterPolicyId"`
	StartTime      string `json:"startTime"`
	EndTime        string `json:"endTime"`
	Field          string `json:"field"`
	TopN           int    `json:"topN"`
}

// FieldStatsResult 字段统计结果
type FieldStatsResult struct {
	Field       string      `json:"field"`
	TotalLogs   int64       `json:"totalLogs"`
	UniqueCount int64       `json:"uniqueCount"`
	Items       []StatsItem `json:"items"`
}

// StatsItem 统计条目
type StatsItem struct {
	Value    string `json:"value"`
	Location string `json:"location"`
	Count    int64  `json:"count"`
	Percent  string `json:"percent"`
	LastSeen string `json:"lastSeen"`
}

// StatsField 统计字段
type StatsField struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// LogQueryParams 日志查询参数
type LogQueryParams struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
	DeviceID  int    `json:"deviceId"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword   string `json:"keyword"`
}

// LogQueryResult 日志查询结果
type LogQueryResult struct {
	Logs  []SyslogLog `json:"logs"`
	Total int64       `json:"total"`
}

// SystemStats 系统统计
type SystemStats struct {
	TotalLogs      int64   `json:"totalLogs"`
	DeviceCount    int     `json:"deviceCount"`
	ServiceRunning bool    `json:"serviceRunning"`
	ListenPort     int     `json:"listenPort"`
	StartTime      string  `json:"startTime"`
	MemoryUsage    uint64  `json:"memoryUsage"`
	CPUUsage       float64 `json:"cpuUsage"`
	Connections    int     `json:"connections"`
	ReceiveRate    float64 `json:"receiveRate"`
	Protocol       string  `json:"protocol"`
	DatabaseSize   int64   `json:"databaseSize"`
	// 仪表盘扩展字段
	MatchedLogs          int64 `json:"matchedLogs"`
	AlertCount           int64 `json:"alertCount"`
	UnmatchedLogs        int64 `json:"unmatchedLogs"`
	ParseTemplateCount   int64 `json:"parseTemplateCount"`
	ActiveFilterPolicies int   `json:"activeFilterPolicies"`
	ActiveAlertPolicies  int   `json:"activeAlertPolicies"`
	ActiveRobots         int   `json:"activeRobots"`
	ActiveDevices        int   `json:"activeDevices"`
	GoroutineCount       int   `json:"goroutineCount"`
}

// ParseTestRequest 解析测试请求
type ParseTestRequest struct {
	ParseType      string `json:"parseType"`
	HeaderRegex    string `json:"headerRegex"`
	FieldMapping   string `json:"fieldMapping"`
	ValueTransform string `json:"valueTransform"`
	SampleLog      string `json:"sampleLog"`
}

// ParseTestResult 解析测试结果
type ParseTestResult struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Fields  []string               `json:"fields"`
	Data    map[string]interface{} `json:"data"`
}

// TestSyslogRequest 测试 Syslog 发送请求
type TestSyslogRequest struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	Message    string `json:"message"`
	Count      int    `json:"count"`
	IntervalMs int    `json:"intervalMs"`
}

// TestSyslogResult 测试 Syslog 发送结果
type TestSyslogResult struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	SentCount   int      `json:"sentCount"`
	FailedCount int      `json:"failedCount"`
	Errors      []string `json:"errors"`
}

// ImportResult 导入结果
type ImportResult struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Count   int      `json:"count"`
	Errors  []string `json:"errors"`
}

// ConfigExport 配置导出结构
type ConfigExport struct {
	Version        string          `json:"version"`
	ExportedAt     string          `json:"exportedAt"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	ParseTemplates []ParseTemplate `json:"parseTemplates,omitempty"`
	FilterPolicies []FilterPolicy  `json:"filterPolicies,omitempty"`
}
