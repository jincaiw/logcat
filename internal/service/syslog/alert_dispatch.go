package syslog

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	"syslog-alert/internal/service/alert"
	"syslog-alert/internal/service/cache"
	"syslog-alert/internal/service/parser"
	"syslog-alert/pkg/constants"
	applogger "syslog-alert/pkg/logger"
)

// sendAlertWithPolicy 按告警规则分发告警到各推送通道。
//
// 流程：
//  1. 获取策略关联的告警规则
//  2. 对每条规则：获取机器人 → 去重检查 → 渲染模板 → 发送 → 记录结果
//  3. 更新日志告警状态与追踪信息
func (s *Server) sendAlertWithPolicy(log *models.SyslogLog, device *models.Device, filterPolicy *models.FilterPolicy, parsedData map[string]interface{}) {
	applogger.Debug("sendAlertWithPolicy: LogID=%d, PolicyID=%d, PolicyName=%s", log.ID, filterPolicy.ID, filterPolicy.Name)

	rules := cache.GetAlertRulesByFilterPolicyID(repository.GetAlertRules, filterPolicy.ID)
	applogger.Debug("Found %d alert rules for policy %d", len(rules), filterPolicy.ID)

	for _, rule := range rules {
		robot, err := cache.GetRobotByID(repository.GetRobots, rule.RobotID)
		if err != nil {
			applogger.Debug("Robot not found for rule %d: %v", rule.ID, err)
			continue
		}

		applogger.Debug("Processing robot: %s (ID: %d, Platform: %s)", robot.Name, robot.ID, robot.Platform)

		if !robot.IsActive || !rule.IsActive {
			applogger.Debug("Robot %s or rule is not active, skipping", robot.Name)
			continue
		}

		platform, ok := constants.NormalizeNotificationPlatform(robot.Platform)
		if !ok {
			applogger.Debug("Robot %s uses unsupported platform %s, skipping", robot.Name, robot.Platform)
			continue
		}

		// 去重检查
		alertKey := s.generateAlertKey(log, filterPolicy, parsedData)
		if filterPolicy.DedupEnabled && s.isDuplicateAlert(alertKey, filterPolicy.DedupWindow) {
			applogger.Debug("Duplicate alert for robot %s, skipping", robot.Name)
			continue
		}

		// 渲染消息
		message, outputTemplate := s.resolveMessage(robot, rule, platform, parsedData, device, log)

		applogger.Debug("Sending to platform: %s", platform)

		// 发送告警
		sendErr := s.dispatchAlert(robot, rule, platform, message, parsedData, log, filterPolicy, outputTemplate)

		// 记录结果
		s.recordAlertResult(log, robot, platform, message, sendErr, alertKey)
	}
}

// resolveMessage 解析输出模板并渲染告警消息。
func (s *Server) resolveMessage(robot *models.Robot, rule models.AlertRule, platform string, parsedData map[string]interface{}, device *models.Device, log *models.SyslogLog) (string, *models.OutputTemplate) {
	var message string
	var outputTemplate *models.OutputTemplate

	// 优先按规则绑定的模板
	if rule.OutputTemplateID > 0 {
		outputTemplate, _ = cache.GetOutputTemplateByID(repository.GetOutputTemplates, rule.OutputTemplateID)
		if outputTemplate != nil {
			applogger.Debug("Got outputTemplate by ID: name=%s, platform=%s", outputTemplate.Name, outputTemplate.Platform)
		}
	}

	// 模板平台不匹配时回退到平台默认模板
	if outputTemplate == nil || (outputTemplate.Platform != platform && platform != constants.PlatformSyslog) {
		outputTemplate, _ = cache.GetOutputTemplateByPlatform(repository.GetOutputTemplates, platform)
		if outputTemplate != nil {
			applogger.Debug("Got outputTemplate by platform: name=%s, platform=%s", outputTemplate.Name, outputTemplate.Platform)
		}
	}

	// 渲染消息（syslog 平台由 dispatchAlert 内部处理）
	if outputTemplate != nil && platform != constants.PlatformSyslog {
		message = renderOutputTemplate(outputTemplate, parsedData, device, log)
		applogger.Debug("Using template %s for robot %s", outputTemplate.Name, robot.Name)
	} else if platform != constants.PlatformSyslog {
		message = defaultAlertMessage(log, device)
		applogger.Debug("No template found, using default message")
	}

	return message, outputTemplate
}

// dispatchAlert 按平台分发告警。
// syslog 平台需要额外的字段映射参数，其他平台通过 alert.Send 统一发送。
func (s *Server) dispatchAlert(robot *models.Robot, rule models.AlertRule, platform, message string, parsedData map[string]interface{}, log *models.SyslogLog, filterPolicy *models.FilterPolicy, outputTemplate *models.OutputTemplate) error {
	if platform == constants.PlatformSyslog {
		return s.sendSyslogAlert(robot, rule, message, parsedData, log, filterPolicy, outputTemplate)
	}
	return alert.Send(robot, message, parsedData, log)
}

// sendSyslogAlert 发送 Syslog 转发告警（需要字段映射等高级配置）。
func (s *Server) sendSyslogAlert(robot *models.Robot, rule models.AlertRule, message string, parsedData map[string]interface{}, log *models.SyslogLog, filterPolicy *models.FilterPolicy, outputTemplate *models.OutputTemplate) error {
	outputFormat := rule.OutputFormat
	if outputFormat == "" {
		outputFormat = robot.SyslogFormat
	}

	// 解析选定字段
	var selectedFields []string
	if outputTemplate != nil && outputTemplate.Fields != "" {
		if err := json.Unmarshal([]byte(outputTemplate.Fields), &selectedFields); err != nil {
			applogger.Debug("Failed to parse template fields: %v", err)
		}
	}

	// 提取字段映射配置
	fieldMapping, fieldNameMapping, valueTransform := s.extractSyslogFieldMapping(filterPolicy, parsedData)

	return alert.SendSyslogForward(
		robot.SyslogHost, robot.SyslogPort, robot.SyslogProtocol, outputFormat,
		message, parsedData, log, fieldMapping, fieldNameMapping, selectedFields, valueTransform,
	)
}

// extractSyslogFieldMapping 从筛选策略的解析模板中提取字段映射配置。
// 针对 smart_delimiter 类型使用子模板的字段索引映射，其他类型使用 FieldMapping。
func (s *Server) extractSyslogFieldMapping(filterPolicy *models.FilterPolicy, parsedData map[string]interface{}) (string, map[string]string, string) {
	if filterPolicy.ParseTemplateID <= 0 {
		return "", nil, ""
	}

	template, err := cache.GetParseTemplateByID(repository.GetParseTemplates, filterPolicy.ParseTemplateID)
	if err != nil || template == nil {
		return "", nil, ""
	}

	valueTransform := template.ValueTransform

	if template.ParseType == constants.ParseTypeSmartDelimiter && template.SubTemplates != "" {
		fieldNameMapping := s.buildSmartDelimiterFieldNameMapping(template, parsedData)
		return "", fieldNameMapping, valueTransform
	}

	return template.FieldMapping, nil, valueTransform
}

// buildSmartDelimiterFieldNameMapping 为 smart_delimiter 类型构建字段名映射。
func (s *Server) buildSmartDelimiterFieldNameMapping(template *models.ParseTemplate, parsedData map[string]interface{}) map[string]string {
	var config struct {
		SubTemplates map[string]parser.SubTemplateConfig `json:"subTemplates"`
	}
	if err := json.Unmarshal([]byte(template.SubTemplates), &config); err != nil {
		return nil
	}

	alertType, ok := parsedData["alertType"].(string)
	if !ok {
		return nil
	}

	subTemplate, ok := config.SubTemplates[alertType]
	if !ok {
		return nil
	}

	fieldNameMapping := make(map[string]string)
	fieldMappings := []struct {
		index int
		name  string
	}{
		{subTemplate.AlertNameField, "告警名称"},
		{subTemplate.AttackIPField, "攻击IP"},
		{subTemplate.VictimIPField, "受害IP"},
		{subTemplate.AlertTimeField, "告警时间"},
		{subTemplate.SeverityField, "威胁等级"},
		{subTemplate.AttackResultField, "攻击结果"},
	}
	for _, fm := range fieldMappings {
		if fm.index >= 0 {
			fieldNameMapping[fmt.Sprintf("field_%d", fm.index)] = fm.name
		}
	}
	return fieldNameMapping
}

// recordAlertResult 记录告警发送结果并更新追踪。
func (s *Server) recordAlertResult(log *models.SyslogLog, robot *models.Robot, platform, message string, sendErr error, alertKey string) {
	record := &models.AlertRecord{
		LogID:      log.ID,
		RobotID:    robot.ID,
		DeviceName: log.DeviceName,
		Message:    message,
		SentAt:     time.Now(),
	}

	if sendErr != nil {
		record.Status = "failed"
		record.ErrorMsg = sendErr.Error()
		log.AlertStatus = constants.AlertStatusFailed
		applogger.Debug("Failed to send to %s: %v", robot.Name, sendErr)
	} else {
		record.Status = "sent"
		log.AlertStatus = constants.AlertStatusSent
		s.markAlertSent(alertKey)
		applogger.Debug("Successfully sent to %s", robot.Name)
	}

	repository.CreateAlertRecord(record)
	repository.UpdateLogAlertStatus(log.ID, log.AlertStatus, 0)

	s.addTraceAlertRecord(log.ID, models.AlertTraceInfo{
		RobotID:   robot.ID,
		RobotName: robot.Name,
		Platform:  platform,
		Status:    record.Status,
		ErrorMsg:  record.ErrorMsg,
		SentAt:    record.SentAt,
	})

	if record.Status == "sent" {
		s.updateTraceAlert(log.ID, constants.AlertStatusSent)
	} else {
		s.updateTraceAlert(log.ID, constants.AlertStatusFailed)
	}
}

// ---- 告警去重 ----

// generateAlertKey 生成告警去重键，基于设备、策略和关键字段。
func (s *Server) generateAlertKey(log *models.SyslogLog, filterPolicy *models.FilterPolicy, parsedData map[string]interface{}) string {
	keyParts := []string{
		fmt.Sprintf("device:%d", log.DeviceID),
		fmt.Sprintf("policy:%d", filterPolicy.ID),
	}
	if attackIp, ok := parsedData["attackIp"]; ok {
		keyParts = append(keyParts, fmt.Sprintf("attackIp:%v", attackIp))
	}
	if threatType, ok := parsedData["threatType"]; ok {
		keyParts = append(keyParts, fmt.Sprintf("threatType:%v", threatType))
	}
	if description, ok := parsedData["description"]; ok {
		keyParts = append(keyParts, fmt.Sprintf("desc:%v", description))
	}
	return strings.Join(keyParts, "|")
}

// isDuplicateAlert 检查是否为重复告警（在去重时间窗口内）。
func (s *Server) isDuplicateAlert(key string, windowSeconds int) bool {
	s.alertMu.Lock()
	defer s.alertMu.Unlock()

	if windowSeconds <= 0 {
		windowSeconds = constants.DefaultDedupWindow
	}

	if lastSent, exists := s.alertCache[key]; exists {
		if time.Since(lastSent) < time.Duration(windowSeconds)*time.Second {
			return true
		}
	}
	return false
}

// markAlertSent 标记告警已发送，并清理过期缓存。
func (s *Server) markAlertSent(key string) {
	s.alertMu.Lock()
	defer s.alertMu.Unlock()

	s.alertCache[key] = time.Now()

	// 缓存超限时清理过期项
	if len(s.alertCache) > constants.AlertCacheMaxSize {
		cutoff := time.Now().Add(-5 * time.Minute)
		for k, v := range s.alertCache {
			if v.Before(cutoff) {
				delete(s.alertCache, k)
			}
		}
	}
}

// ---- 模板渲染 ----

// templateVarRegex 匹配 {{变量名}} 占位符。
var templateVarRegex = regexp.MustCompile(`\{\{([a-zA-Z0-9_.]+)\}\}`)

// renderOutputTemplate 渲染输出模板，将 {{变量}} 替换为实际值。
func renderOutputTemplate(template *models.OutputTemplate, data map[string]interface{}, device *models.Device, log *models.SyslogLog) string {
	if data == nil {
		data = make(map[string]interface{})
	}

	// 注入设备信息
	if device != nil {
		data["deviceName"] = device.Name
		data["deviceIP"] = device.IPAddress
	}

	// 注入日志元数据
	data["rawMessage"] = log.RawMessage
	data["receivedAt"] = log.ReceivedAt.Format("2006-01-02 15:04:05")

	// 处理 localTimestamp 毫秒时间戳
	if ts, ok := data["localTimestamp"]; ok {
		if milli, ok := ts.(float64); ok && milli > 1e12 {
			formatted := time.UnixMilli(int64(milli)).Format("2006-01-02 15:04:05")
			data["timestamp"] = formatted
			data["alertTime"] = formatted
		}
	}

	// 默认时间字段
	if _, ok := data["timestamp"]; !ok {
		data["timestamp"] = log.ReceivedAt.Format("2006-01-02 15:04:05")
	}
	if _, ok := data["alertTime"]; !ok {
		data["alertTime"] = data["timestamp"]
	}

	content := template.Content
	content = templateVarRegex.ReplaceAllStringFunc(content, func(match string) string {
		fieldName := strings.Trim(match, "{}")
		// 优先从扁平化数据中取值
		if value, exists := data[fieldName]; exists {
			return fmt.Sprintf("%v", value)
		}
		// 回退到嵌套取值
		if value := getNestedValue(data, fieldName); value != nil {
			return fmt.Sprintf("%v", value)
		}
		return ""
	})
	return content
}

// getNestedValue 按点分路径从嵌套 map 中取值。
func getNestedValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	var current interface{} = data
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			val, exists := v[part]
			if !exists {
				return nil
			}
			current = val
		default:
			return nil
		}
	}
	return current
}

// defaultAlertMessage 生成默认告警消息（无模板时使用）。
func defaultAlertMessage(log *models.SyslogLog, device *models.Device) string {
	deviceName := "Unknown"
	if device != nil {
		deviceName = device.Name
	}
	return fmt.Sprintf("### 🚨 安全告警\n\n"+
		"**设备名称**: %s\n\n"+
		"**来源IP**: %s\n\n"+
		"**告警时间**: %s\n\n"+
		"**告警内容**: %s",
		deviceName,
		log.SourceIP,
		log.ReceivedAt.Format("2006-01-02 15:04:05"),
		log.RawMessage,
	)
}
