// Package filter 提供日志筛选与白名单匹配能力。
//
// 合并了原 filter.go 与 syslog_service.go 中重复的筛选逻辑，统一使用：
//   - 真正的正则匹配（修复原 syslog_service.go 中 regexpMatch 用 strings.Contains 的 Bug）
//   - 结构化日志替代 stdlog.Printf("[DEBUG] ...")
//   - CIDR 白名单匹配
//   - AND/OR 条件逻辑
//
// 核心入口：
//   - NewEngine(policy) 创建筛选引擎
//   - Engine.Match(log) 匹配单条日志
//   - ProcessLogWithPolicies(log, device) 按设备策略链匹配
package filter

import (
	"encoding/json"
	"fmt"
	"sync"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	"syslog-alert/internal/service/cache"
	"syslog-alert/internal/service/parser"
	"syslog-alert/pkg/constants"
	applogger "syslog-alert/pkg/logger"
)

// Engine 筛选引擎，绑定一个筛选策略与可选的解析器。
type Engine struct {
	policy *models.FilterPolicy
	parser parser.Parser
}

// NewEngine 创建筛选引擎。若策略配置了解析模板，会自动创建对应的解析器。
func NewEngine(policy *models.FilterPolicy) (*Engine, error) {
	engine := &Engine{policy: policy}

	if policy.ParseTemplateID > 0 {
		p, err := cache.GetParserByTemplateID(repository.GetParseTemplates, policy.ParseTemplateID)
		if err != nil {
			return nil, fmt.Errorf("failed to create parser: %v", err)
		}
		engine.parser = p
	}

	return engine, nil
}

// Match 匹配日志，返回是否命中策略、解析数据、错误。
//
// 流程：
//  1. 解析日志（使用策略绑定的解析器，或回退到日志已有的 ParsedData）
//  2. 检查白名单（命中白名单返回 false，不推送告警）
//  3. 检查筛选条件（按 AND/OR 逻辑组合）
func (e *Engine) Match(log *models.SyslogLog) (bool, map[string]interface{}, error) {
	parsedData, err := e.parseLog(log)
	if err != nil {
		return false, nil, err
	}

	// 白名单检查：命中则跳过（返回 false）
	if e.policy.Whitelist != "" && e.policy.WhitelistField != "" {
		matched, err := MatchWhitelist(parsedData, e.policy.WhitelistField, e.policy.Whitelist)
		if err != nil {
			return false, nil, err
		}
		if matched {
			return false, parsedData, nil
		}
	}

	// 无条件 = 匹配
	if e.policy.Conditions == "" {
		return true, parsedData, nil
	}

	conditions, err := ParseConditions(e.policy.Conditions)
	if err != nil {
		return false, nil, fmt.Errorf("invalid conditions: %v", err)
	}

	matched := EvaluateConditions(conditions, parsedData, e.policy.ConditionLogic)
	return matched, parsedData, nil
}

// parseLog 解析日志。优先使用策略绑定的解析器，回退到日志已有的 ParsedData。
func (e *Engine) parseLog(log *models.SyslogLog) (map[string]interface{}, error) {
	if e.parser != nil {
		return e.parser.Parse(log.RawMessage)
	}

	// 回退：使用已存储的解析数据
	if log.ParsedData == "" {
		return make(map[string]interface{}), nil
	}
	var parsedData map[string]interface{}
	if err := json.Unmarshal([]byte(log.ParsedData), &parsedData); err != nil {
		return make(map[string]interface{}), nil
	}
	return parsedData, nil
}

// parseConditions 解析条件 JSON。
var conditionsCache sync.Map

func ParseConditions(conditionsJSON string) ([]models.FilterCondition, error) {
	if cached, ok := conditionsCache.Load(conditionsJSON); ok {
		return cached.([]models.FilterCondition), nil
	}
	var conditions []models.FilterCondition
	if err := json.Unmarshal([]byte(conditionsJSON), &conditions); err != nil {
		return nil, err
	}
	conditionsCache.Store(conditionsJSON, conditions)
	return conditions, nil
}

// ProcessLogWithPolicies 按设备策略链匹配日志。
//
// 策略优先级：
//  1. 设备绑定的策略
//  2. 设备分组绑定的策略
//  3. 全局策略
//
// 返回匹配的策略与解析数据。若策略动作为 discard，返回 error。
func ProcessLogWithPolicies(log *models.SyslogLog, device *models.Device) (*models.FilterPolicy, map[string]interface{}, error) {
	policies := getPoliciesForDevice(device)

	applogger.Debug("ProcessLogWithPolicies: LogID=%d, PoliciesCount=%d", log.ID, len(policies))

	for i := range policies {
		policy := &policies[i]
		applogger.Debug("Checking policy: ID=%d, Name=%s, IsActive=%v", policy.ID, policy.Name, policy.IsActive)

		if !policy.IsActive {
			continue
		}

		// 设备/分组匹配过滤
		if !policyMatchesDevice(policy, device) {
			continue
		}

		engine, err := NewEngine(policy)
		if err != nil {
			applogger.Warn("Failed to create filter engine for policy %s: %v", policy.Name, err)
			continue
		}

		matched, parsedData, err := engine.Match(log)
		if err != nil {
			applogger.Debug("Policy %s match error: %v", policy.Name, err)
			continue
		}

		if matched {
			if policy.Action == constants.ActionKeep {
				return policy, parsedData, nil
			}
			return nil, nil, fmt.Errorf("discarded by policy: %s", policy.Name)
		}
	}

	return nil, nil, nil
}

// getPoliciesForDevice 按优先级获取设备相关的策略列表。
func getPoliciesForDevice(device *models.Device) []models.FilterPolicy {
	if device != nil && device.ID > 0 {
		policies := cache.GetFilterPoliciesByDeviceID(repository.GetFilterPolicies, device.ID)
		if len(policies) > 0 {
			return policies
		}
		if device.GroupID > 0 {
			policies = cache.GetFilterPoliciesByDeviceGroupID(repository.GetFilterPolicies, device.GroupID)
			if len(policies) > 0 {
				return policies
			}
		}
	}
	return cache.GetFilterPolicies(repository.GetFilterPolicies)
}

// policyMatchesDevice 检查策略是否适用于指定设备。
func policyMatchesDevice(policy *models.FilterPolicy, device *models.Device) bool {
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
