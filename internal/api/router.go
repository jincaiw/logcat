package api

import (
	"net/http"
)

// NewRouter 创建并注册所有 API 路由。
// 使用 Go 1.22 增强的 http.ServeMux，支持方法路由与路径参数 {id}。
// 所有 /api/ 路由（除 /api/auth/login）均经过 AuthMiddleware 认证保护。
func NewRouter(ws *WebServer) *http.ServeMux {
	mux := http.NewServeMux()

	// 公开路由：无需登录
	mux.HandleFunc("GET /healthz", ws.Health)
	mux.HandleFunc("GET /api/health", ws.Health)
	mux.HandleFunc("POST /api/auth/login", ws.Login)

	// 受保护路由：需登录后访问
	protected := http.NewServeMux()

	// ---- 认证相关 ----
	protected.HandleFunc("POST /api/auth/logout", ws.Logout)
	protected.HandleFunc("GET /api/auth/profile", ws.GetProfile)
	protected.HandleFunc("PUT /api/auth/profile", ws.UpdateProfile)
	protected.HandleFunc("PUT /api/auth/password", ws.ChangePassword)

	// ---- 设备管理 ----
	protected.HandleFunc("GET /api/devices", ws.ListDevices)
	protected.HandleFunc("POST /api/devices", ws.CreateDevice)
	protected.HandleFunc("GET /api/devices/{id}", ws.GetDevice)
	protected.HandleFunc("PUT /api/devices/{id}", ws.UpdateDevice)
	protected.HandleFunc("DELETE /api/devices/{id}", ws.DeleteDevice)

	// ---- 设备分组 ----
	protected.HandleFunc("GET /api/device-groups", ws.ListDeviceGroups)
	protected.HandleFunc("POST /api/device-groups", ws.CreateDeviceGroup)
	protected.HandleFunc("GET /api/device-groups/{id}", ws.GetDeviceGroup)
	protected.HandleFunc("PUT /api/device-groups/{id}", ws.UpdateDeviceGroup)
	protected.HandleFunc("DELETE /api/device-groups/{id}", ws.DeleteDeviceGroup)

	// ---- 解析模板 ----
	protected.HandleFunc("GET /api/parse-templates", ws.ListParseTemplates)
	protected.HandleFunc("POST /api/parse-templates", ws.CreateParseTemplate)
	protected.HandleFunc("GET /api/parse-templates/{id}", ws.GetParseTemplate)
	protected.HandleFunc("PUT /api/parse-templates/{id}", ws.UpdateParseTemplate)
	protected.HandleFunc("DELETE /api/parse-templates/{id}", ws.DeleteParseTemplate)

	// ---- 输出模板 ----
	protected.HandleFunc("GET /api/output-templates", ws.ListOutputTemplates)
	protected.HandleFunc("POST /api/output-templates", ws.CreateOutputTemplate)
	protected.HandleFunc("GET /api/output-templates/{id}", ws.GetOutputTemplate)
	protected.HandleFunc("PUT /api/output-templates/{id}", ws.UpdateOutputTemplate)
	protected.HandleFunc("DELETE /api/output-templates/{id}", ws.DeleteOutputTemplate)

	// ---- 筛选策略 ----
	protected.HandleFunc("GET /api/filter-policies", ws.ListFilterPolicies)
	protected.HandleFunc("POST /api/filter-policies", ws.CreateFilterPolicy)
	protected.HandleFunc("GET /api/filter-policies/{id}", ws.GetFilterPolicy)
	protected.HandleFunc("PUT /api/filter-policies/{id}", ws.UpdateFilterPolicy)
	protected.HandleFunc("DELETE /api/filter-policies/{id}", ws.DeleteFilterPolicy)

	// ---- 告警策略 ----
	protected.HandleFunc("GET /api/alert-policies", ws.ListAlertPolicies)
	protected.HandleFunc("POST /api/alert-policies", ws.CreateAlertPolicy)
	protected.HandleFunc("GET /api/alert-policies/{id}", ws.GetAlertPolicy)
	protected.HandleFunc("PUT /api/alert-policies/{id}", ws.UpdateAlertPolicy)
	protected.HandleFunc("DELETE /api/alert-policies/{id}", ws.DeleteAlertPolicy)

	// ---- 机器人/推送通道 ----
	protected.HandleFunc("GET /api/robots", ws.ListRobots)
	protected.HandleFunc("POST /api/robots", ws.CreateRobot)
	protected.HandleFunc("GET /api/robots/{id}", ws.GetRobot)
	protected.HandleFunc("PUT /api/robots/{id}", ws.UpdateRobot)
	protected.HandleFunc("DELETE /api/robots/{id}", ws.DeleteRobot)
	protected.HandleFunc("POST /api/test-robot", ws.TestRobot)

	// ---- 告警规则 ----
	protected.HandleFunc("GET /api/alert-rules/robot/{robotId}", ws.ListAlertRulesByRobot)
	protected.HandleFunc("POST /api/alert-rules", ws.CreateAlertRule)
	protected.HandleFunc("GET /api/alert-rules/{id}", ws.GetAlertRule)
	protected.HandleFunc("PUT /api/alert-rules/{id}", ws.UpdateAlertRule)
	protected.HandleFunc("DELETE /api/alert-rules/{id}", ws.DeleteAlertRule)
	protected.HandleFunc("DELETE /api/alert-rules/robot/{robotId}", ws.DeleteAlertRulesByRobot)

	// ---- 字段映射文档 ----
	protected.HandleFunc("GET /api/field-mapping-docs", ws.ListFieldMappingDocs)
	protected.HandleFunc("POST /api/field-mapping-docs", ws.CreateFieldMappingDoc)
	protected.HandleFunc("GET /api/field-mapping-docs/{id}", ws.GetFieldMappingDoc)
	protected.HandleFunc("PUT /api/field-mapping-docs/{id}", ws.UpdateFieldMappingDoc)
	protected.HandleFunc("DELETE /api/field-mapping-docs/{id}", ws.DeleteFieldMappingDoc)

	// ---- 日志 ----
	protected.HandleFunc("GET /api/logs", ws.ListLogs)
	protected.HandleFunc("POST /api/logs/cleanup", ws.CleanupLogs)
	protected.HandleFunc("DELETE /api/logs/cleanup-all", ws.CleanupAllLogs)
	protected.HandleFunc("GET /api/logs/unmatched-count", ws.GetUnmatchedLogsCount)
	protected.HandleFunc("POST /api/logs/cleanup-unmatched", ws.CleanupUnmatchedLogs)

	// ---- 服务管理 ----
	protected.HandleFunc("GET /api/service/status", ws.GetServiceStatus)
	protected.HandleFunc("POST /api/service/start", ws.StartService)
	protected.HandleFunc("POST /api/service/stop", ws.StopService)

	// ---- 配置 ----
	protected.HandleFunc("GET /api/config", ws.GetConfig)
	protected.HandleFunc("PUT /api/config", ws.SaveConfig)

	// ---- 统计 ----
	protected.HandleFunc("GET /api/stats", ws.GetStats)
	protected.HandleFunc("POST /api/field-stats", ws.GetFieldStats)
	protected.HandleFunc("GET /api/available-stats-fields/{policyId}", ws.GetAvailableStatsFields)

	// ---- 测试工具 ----
	protected.HandleFunc("POST /api/test-syslog", ws.SendTestSyslog)
	protected.HandleFunc("POST /api/test-syslog-forward", ws.TestSyslogForward)
	protected.HandleFunc("POST /api/test-parse", ws.TestParse)

	// ---- 追踪 ----
	protected.HandleFunc("GET /api/log-trace/{logId}", ws.GetLogTrace)

	// ---- 网络 ----
	protected.HandleFunc("GET /api/local-ips", ws.GetLocalIPs)
	protected.HandleFunc("GET /api/server-ip", ws.GetServerIP)

	// ---- 导入导出 ----
	protected.HandleFunc("GET /api/export/parse-templates", ws.ExportParseTemplates)
	protected.HandleFunc("GET /api/export/filter-policies", ws.ExportFilterPolicies)
	protected.HandleFunc("POST /api/import/parse-templates", ws.ImportParseTemplates)
	protected.HandleFunc("POST /api/import/filter-policies", ws.ImportFilterPolicies)

	// ---- 预设模板 ----
	protected.HandleFunc("GET /api/preset-templates", ws.GetPresetTemplates)

	// 将受保护路由挂载到 /api/ 前缀，外层包裹认证中间件
	mux.Handle("/api/", AuthMiddleware(protected))

	return mux
}
