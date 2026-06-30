package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// delimiterParser 使用固定分隔符拆分日志，按字段名列表或类型映射赋值。
// 配置格式（FieldMapping 字段）：
//
//	{
//	  "delimiter": "|!",
//	  "fields": ["field1", "field2"],
//	  "type_field": "alertType",
//	  "type_mapping": {"typeA": ["f1", "f2"], "typeB": ["g1", "g2"]}
//	}
//
// 也支持简单映射格式 {"oldField": "newField"}。
type delimiterParser struct {
	*baseParser
}

func (p *delimiterParser) Parse(rawLog string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	config, simpleMapping, err := parseDelimiterConfigWithFallback(p.template.FieldMapping)
	if err != nil {
		return nil, err
	}

	content := rawLog

	// 提取头部字段（若配置了 HeaderRegex）
	if p.regex != nil {
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

	// 按类型映射或字段列表赋值
	if config.TypeField != "" && len(values) > 0 {
		alertType := values[0]
		result[config.TypeField] = alertType

		if fields, ok := config.TypeMapping[alertType]; ok {
			assignFieldsByIndex(result, fields, values)
		} else if len(config.Fields) > 0 {
			assignFieldsByIndex(result, config.Fields, values)
		} else {
			assignFieldsByIndexDefault(result, values)
		}
	} else if len(config.Fields) > 0 {
		assignFieldsByIndex(result, config.Fields, values)
	} else {
		assignFieldsByIndexDefault(result, values)
	}

	// 应用简单字段映射（重命名）
	if len(simpleMapping) > 0 {
		applySimpleMappingToResult(result, simpleMapping)
	}

	if p.template.ValueTransform != "" {
		result = ApplyValueTransform(result, p.template.ValueTransform)
	}
	return result, nil
}

// parseDelimiterConfigWithFallback 解析分隔符配置，支持复杂格式与简单映射格式的回退。
func parseDelimiterConfigWithFallback(fieldMappingJSON string) (*DelimiterConfig, map[string]string, error) {
	if fieldMappingJSON == "" {
		return &DelimiterConfig{Delimiter: "|!"}, nil, nil
	}

	// 先尝试复杂格式
	var config DelimiterConfig
	if err := json.Unmarshal([]byte(fieldMappingJSON), &config); err == nil {
		// 检查是否为有效的复杂配置（至少有一个有意义的字段）
		// 否则任意 JSON 对象都会"成功"解析为零值 DelimiterConfig，导致简单映射回退永不触发
		if config.Delimiter != "" || len(config.Fields) > 0 || config.TypeField != "" || len(config.TypeMapping) > 0 {
			if config.Delimiter == "" {
				config.Delimiter = "|!"
			}
			return &config, nil, nil
		}
	}

	// 回退到简单映射格式
	var simpleMapping map[string]string
	if err := json.Unmarshal([]byte(fieldMappingJSON), &simpleMapping); err != nil {
		return nil, nil, fmt.Errorf("invalid delimiter config: %v", err)
	}
	return &DelimiterConfig{Delimiter: "|!"}, simpleMapping, nil
}

// applySimpleMappingToResult 对已解析结果应用简单字段映射（支持点分路径取值）。
func applySimpleMappingToResult(result map[string]interface{}, mapping map[string]string) {
	for oldField, newField := range mapping {
		var value interface{}
		if strings.Contains(oldField, ".") {
			value = getNestedValue(result, oldField)
		} else if v, exists := result[oldField]; exists {
			value = v
		}
		if value != nil {
			result[newField] = value
		}
	}
}

// kvDelimiterParser 解析 "key:val|!key:val" 格式的键值对日志。
// 分隔符与键值分隔符均可配置。
type kvDelimiterParser struct {
	*baseParser
}

func (p *kvDelimiterParser) Parse(rawLog string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	config, err := ParseKVDelimiterConfig(p.template.FieldMapping)
	if err != nil {
		return nil, err
	}

	pairs := strings.Split(rawLog, config.Delimiter)
	for _, pair := range pairs {
		idx := strings.Index(pair, config.KVSeparator)
		if idx > 0 {
			key := strings.TrimSpace(pair[:idx])
			value := strings.TrimSpace(pair[idx+1:])
			result[key] = value
		}
	}

	if p.template.ValueTransform != "" {
		result = ApplyValueTransform(result, p.template.ValueTransform)
	}
	return result, nil
}
