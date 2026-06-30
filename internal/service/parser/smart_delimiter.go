package parser

import (
	"fmt"
	"regexp"
	"strings"

	"syslog-alert/pkg/constants"
)

// smartDelimiterParser 智能分隔符解析器。
// 根据日志中的类型字段（TypeField）选择对应的子模板配置，按字段索引提取告警字段。
// 支持自动跳过 syslog 头部（SkipHeader）。
type smartDelimiterParser struct {
	*baseParser
}

func (p *smartDelimiterParser) Parse(rawLog string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	config, err := ParseSmartDelimiterConfig(p.template.FieldMapping)
	if err != nil {
		return nil, err
	}

	content := rawLog

	// 处理头部
	if config.SkipHeader {
		content = p.extractSmartHeader(rawLog, config.HeaderRegex, result)
	} else if p.regex != nil {
		end, headerFields := p.extractHeaderFields(rawLog)
		if headerFields == nil {
			return nil, fmt.Errorf("log does not match header pattern")
		}
		for k, v := range headerFields {
			result[k] = v
		}
		if end > 0 && end < len(rawLog) {
			content = strings.TrimSpace(rawLog[end:])
		}
	}

	values := strings.Split(content, config.Delimiter)

	if len(values) <= config.TypeField {
		return nil, fmt.Errorf("log does not have enough fields")
	}

	alertType := values[config.TypeField]
	result["alertType"] = alertType

	// 保存所有字段索引
	for i, v := range values {
		result[fmt.Sprintf("field_%d", i)] = v
	}

	// 按子模板配置提取命名字段
	if subConfig, ok := config.SubTemplates[alertType]; ok {
		applySubTemplate(result, values, &subConfig)
	}

	// IOC 告警默认攻击结果为"失陷"
	if alertType == constants.AlertTypeIOC {
		if _, exists := result["attackResult"]; !exists {
			result["attackResult"] = constants.AttackResultCompromised
		}
	}

	if p.template.ValueTransform != "" {
		result = ApplyValueTransform(result, p.template.ValueTransform)
	}
	return result, nil
}

// extractSmartHeader 处理 SkipHeader 模式下的头部提取。
// 返回去除头部后的日志内容，头部字段写入 result。
func (p *smartDelimiterParser) extractSmartHeader(rawLog, headerRegex string, result map[string]interface{}) string {
	pattern := headerRegex
	if pattern == "" {
		pattern = constants.DefaultSyslogHeaderRegex
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return rawLog
	}
	matches := re.FindStringSubmatch(rawLog)
	if matches == nil {
		return rawLog
	}
	subexpNames := re.SubexpNames()
	for i, name := range subexpNames {
		if name != "" && i < len(matches) {
			result[name] = matches[i]
		}
	}
	loc := re.FindStringIndex(rawLog)
	if loc != nil && loc[1] < len(rawLog) {
		return strings.TrimSpace(rawLog[loc[1]:])
	}
	return rawLog
}

// applySubTemplate 按子模板配置从值列表中提取命名字段。
func applySubTemplate(result map[string]interface{}, values []string, sub *SubTemplateConfig) {
	fieldMappings := []struct {
		index int
		name  string
	}{
		{sub.AlertNameField, "alertName"},
		{sub.AttackIPField, "attackIP"},
		{sub.VictimIPField, "victimIP"},
		{sub.AlertTimeField, "alertTime"},
		{sub.SeverityField, "severity"},
		{sub.AttackResultField, "attackResult"},
	}
	for _, fm := range fieldMappings {
		if fm.index >= 0 && fm.index < len(values) {
			result[fm.name] = values[fm.index]
		}
	}
	// 处理自定义字段
	for _, cf := range sub.CustomFields {
		if cf.FieldIndex >= 0 && cf.FieldIndex < len(values) && cf.Name != "" {
			result[cf.Name] = values[cf.FieldIndex]
		}
	}
}
