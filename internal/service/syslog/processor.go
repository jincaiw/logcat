package syslog

import (
	"encoding/json"
	"fmt"
	"strings"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	"syslog-alert/internal/service/filter"
	"syslog-alert/internal/service/parser"
	"syslog-alert/pkg/constants"
	applogger "syslog-alert/pkg/logger"
)

// handleMessage 处理单条 Syslog 消息：创建日志记录、更新统计、触发过滤与告警。
func (s *Server) handleMessage(msg Message) {
	device, _ := repository.GetDeviceByIP(msg.SourceIP)

	deviceName := "Unknown"
	deviceID := uint(0)
	if device != nil {
		deviceName = device.Name
		deviceID = device.ID
	}

	priority, facility, severity := parsePriority(msg.Message)
	isForwarded := checkForwardedMark(msg.Message)

	syslogLog := &models.SyslogLog{
		DeviceID:     deviceID,
		DeviceName:   deviceName,
		SourceIP:     msg.SourceIP,
		RawMessage:   msg.Message,
		Priority:     fmt.Sprintf("%d", priority),
		Facility:     facility,
		Severity:     severity,
		Timestamp:    msg.ReceivedAt,
		ReceivedAt:   msg.ReceivedAt,
		FilterStatus: constants.FilterStatusPending,
		AlertStatus:  constants.AlertStatusNone,
	}

	repository.CreateLog(syslogLog)
	s.incrementReceiveCount()
	s.createTrace(syslogLog.ID, msg.SourceIP, msg.Message)

	if s.statsUpdate != nil {
		s.statsUpdate.UpdateStats(repository.GetLogCount(), int(repository.GetDeviceCount()), true)
	}

	config := repository.GetSystemConfig()
	applogger.Debug("AlertEnabled: %v, LogID: %d, IsForwarded: %v", config.AlertEnabled, syslogLog.ID, isForwarded)
	if config.AlertEnabled && !isForwarded {
		s.processLogWithPolicies(syslogLog, device)
	}
}

// processLogWithPolicies 按筛选策略链处理日志。
//
// 流程：
//  1. 遍历所有策略，按设备/分组过滤
//  2. 对每个策略：解析日志 → 检查白名单 → 检查条件
//  3. 匹配成功：更新日志状态，发送告警
//  4. 匹配失败：检查是否丢弃未匹配日志
func (s *Server) processLogWithPolicies(syslogLog *models.SyslogLog, device *models.Device) {
	policies := repository.GetFilterPolicies()
	applogger.Debug("processLogWithPolicies: LogID=%d, PoliciesCount=%d", syslogLog.ID, len(policies))

	var matchedPolicy *models.FilterPolicy
	var parsedData map[string]interface{}
	var hasActivePolicy bool
	var hasAnyPolicy bool
	var whitelistedByPolicy string

	for i := range policies {
		hasAnyPolicy = true
		policy := &policies[i]
		applogger.Debug("Checking policy: ID=%d, Name=%s, IsActive=%v", policy.ID, policy.Name, policy.IsActive)

		if !policy.IsActive {
			continue
		}
		hasActivePolicy = true

		if !devicePolicyMatch(policy, device) {
			continue
		}

		data, templateName, parseErr := s.parseLogForPolicy(syslogLog, policy)
		if parseErr != nil {
			applogger.Debug("Parse failed: %v", parseErr)
			s.updateTraceParse(syslogLog.ID, "failed", templateName, "", parseErr.Error())
			continue
		}
		applogger.Debug("Parsed data: %+v", data)
		s.updateTraceParse(syslogLog.ID, "success", templateName, fmt.Sprintf("%+v", data), "")

		// 白名单检查
		if policy.Whitelist != "" && policy.WhitelistField != "" {
			whitelisted, whitelistErr := filter.MatchWhitelist(data, policy.WhitelistField, policy.Whitelist)
			if whitelistErr != nil {
				applogger.Debug("Whitelist check error: %v", whitelistErr)
			}
			if whitelisted {
				applogger.Debug("Log matched whitelist, skipping")
				whitelistedByPolicy = policy.Name
				continue
			}
		}

		// 条件匹配
		if s.matchConditions(data, policy) {
			matchedPolicy = policy
			parsedData = data
			applogger.Debug("Policy %s matched!", policy.Name)
			break
		}
		applogger.Debug("Policy %s did not match", policy.Name)
	}

	if matchedPolicy != nil {
		s.handleMatchedPolicy(syslogLog, device, matchedPolicy, parsedData)
	} else {
		s.handleUnmatchedPolicy(syslogLog, hasAnyPolicy, hasActivePolicy, whitelistedByPolicy)
	}
}

// devicePolicyMatch 检查策略是否适用于设备。
func devicePolicyMatch(policy *models.FilterPolicy, device *models.Device) bool {
	if policy.DeviceID > 0 && (device == nil || policy.DeviceID != device.ID) {
		applogger.Debug("Policy %s DeviceID mismatch, skipping", policy.Name)
		return false
	}
	if policy.DeviceGroupID > 0 && (device == nil || policy.DeviceGroupID != device.GroupID) {
		applogger.Debug("Policy %s DeviceGroupID mismatch, skipping", policy.Name)
		return false
	}
	return true
}

// parseLogForPolicy 使用策略绑定的解析模板解析日志。
// 返回解析数据、模板名称、错误。
func (s *Server) parseLogForPolicy(syslogLog *models.SyslogLog, policy *models.FilterPolicy) (map[string]interface{}, string, error) {
	if policy.ParseTemplateID <= 0 {
		data := parseSyslogToMap(syslogLog.RawMessage)
		applogger.Debug("Using syslog map: %+v", data)
		return data, "syslog", nil
	}

	template, err := repository.GetParseTemplateByID(policy.ParseTemplateID)
	if err != nil {
		applogger.Debug("Failed to get parse template: %v", err)
		return nil, "", err
	}

	p, err := parser.New(template)
	if err != nil {
		return nil, template.Name, err
	}

	data, err := p.Parse(syslogLog.RawMessage)
	if err != nil {
		return nil, template.Name, err
	}
	return data, template.Name, nil
}

// matchConditions 评估策略条件是否匹配。
func (s *Server) matchConditions(data map[string]interface{}, policy *models.FilterPolicy) bool {
	if policy.Conditions == "" {
		return true
	}

	var conditions []models.FilterCondition
	if err := json.Unmarshal([]byte(policy.Conditions), &conditions); err != nil {
		return false
	}

	return filter.EvaluateConditions(conditions, data, policy.ConditionLogic)
}

// handleMatchedPolicy 处理匹配到策略的日志。
func (s *Server) handleMatchedPolicy(syslogLog *models.SyslogLog, device *models.Device, matchedPolicy *models.FilterPolicy, parsedData map[string]interface{}) {
	syslogLog.FilterStatus = constants.FilterStatusMatched
	syslogLog.MatchedPolicyID = matchedPolicy.ID

	if parsedData != nil {
		parsedBytes, _ := json.Marshal(parsedData)
		syslogLog.ParsedData = string(parsedBytes)
		syslogLog.ParsedFields = parser.ExtractKeyFields(parsedData)
	}

	repository.UpdateLogFilterStatus(syslogLog.ID, constants.FilterStatusMatched, matchedPolicy.ID)
	if syslogLog.ParsedData != "" {
		repository.UpdateLogParsedFields(syslogLog.ID, syslogLog.ParsedData, syslogLog.ParsedFields)
	}

	s.updateTraceFilter(syslogLog.ID, constants.FilterStatusMatched, true, matchedPolicy.Name, matchedPolicy.Conditions, constants.ActionKeep)

	if matchedPolicy.Action == constants.ActionKeep {
		s.updateTraceAlert(syslogLog.ID, constants.AlertStatusPending)
		s.sendAlertWithPolicy(syslogLog, device, matchedPolicy, parsedData)
	} else if matchedPolicy.Action == constants.ActionDiscard {
		s.updateTraceFilter(syslogLog.ID, constants.FilterStatusMatched, true, matchedPolicy.Name, matchedPolicy.Conditions, constants.ActionDiscard)
		repository.DeleteLog(syslogLog.ID)
	}
}

// handleUnmatchedPolicy 处理未匹配任何策略的日志。
func (s *Server) handleUnmatchedPolicy(syslogLog *models.SyslogLog, hasAnyPolicy, hasActivePolicy bool, whitelistedByPolicy string) {
	syslogLog.FilterStatus = constants.FilterStatusUnmatched
	repository.UpdateLogFilterStatus(syslogLog.ID, constants.FilterStatusUnmatched, 0)

	// 白名单匹配但无条件匹配
	if whitelistedByPolicy != "" {
		applogger.Debug("Log matched whitelist by policy %s but no filter conditions matched", whitelistedByPolicy)
		s.updateTraceFilter(syslogLog.ID, constants.FilterStatusWhitelisted, true, whitelistedByPolicy, "", "whitelist")
		repository.UpdateLogFilterStatus(syslogLog.ID, constants.FilterStatusWhitelisted, 0)
		return
	}

	// 检查是否需要丢弃未匹配日志
	if hasActivePolicy {
		for _, p := range repository.GetFilterPolicies() {
			if p.IsActive && p.DropUnmatched {
				applogger.Debug("Dropping unmatched log by policy: %s", p.Name)
				s.updateTraceFilter(syslogLog.ID, constants.FilterStatusDropped, true, "", "", "drop_unmatched")
				repository.DeleteLog(syslogLog.ID)
				return
			}
		}
	}

	if !hasAnyPolicy || !hasActivePolicy {
		applogger.Debug("No active policy found")
		s.updateTraceFilter(syslogLog.ID, constants.FilterStatusDisabled, false, "", "", "no active policy")
	} else {
		applogger.Debug("No policy matched")
		s.updateTraceFilter(syslogLog.ID, constants.FilterStatusUnmatched, true, "", "", "no policy matched")
	}
}

// parseSyslogToMap 将原始 syslog 消息解析为简单 map（无解析模板时的回退方案）。
func parseSyslogToMap(msg string) map[string]interface{} {
	result := make(map[string]interface{})
	if len(msg) == 0 {
		return result
	}

	start := strings.Index(msg, ">")
	if start == -1 {
		result["message"] = msg
		return result
	}

	content := msg[start+1:]
	result["message"] = content

	// 尝试提取 JSON 部分
	jsonStart := strings.Index(content, "{")
	if jsonStart != -1 {
		jsonStr := content[jsonStart:]
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &jsonData); err == nil {
			for k, v := range jsonData {
				result[k] = v
			}
		}
	}
	return result
}
