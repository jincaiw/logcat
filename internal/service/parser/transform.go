package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"syslog-alert/pkg/constants"
)

// SubTemplateConfig 智能分隔符解析的子模板配置。
// 按告警类型（alertType）映射到不同的字段索引方案。
type SubTemplateConfig struct {
	AlertNameField    int                 `json:"alertNameField"`
	AttackIPField     int                 `json:"attackIPField"`
	VictimIPField     int                 `json:"victimIPField"`
	AlertTimeField    int                 `json:"alertTimeField"`
	SeverityField     int                 `json:"severityField"`
	AttackResultField int                 `json:"attackResultField"`
	CustomFields      []CustomFieldConfig `json:"customFields,omitempty"`
}

// CustomFieldConfig 自定义字段配置，将分隔符字段索引映射到字段名。
type CustomFieldConfig struct {
	Name       string `json:"name"`
	FieldIndex int    `json:"fieldIndex"`
}

// SmartDelimiterConfig 智能分隔符解析的整体配置。
type SmartDelimiterConfig struct {
	Delimiter    string                       `json:"delimiter"`
	TypeField    int                          `json:"typeField"`
	SkipHeader   bool                         `json:"skipHeader"`
	HeaderRegex  string                       `json:"headerRegex"`
	SubTemplates map[string]SubTemplateConfig `json:"subTemplates"`
}

// DelimiterConfig 普通分隔符解析配置。
type DelimiterConfig struct {
	Delimiter   string              `json:"delimiter"`
	Fields      []string            `json:"fields"`
	TypeField   string              `json:"type_field"`
	TypeMapping map[string][]string `json:"type_mapping"`
}

// KVDelimiterConfig 键值分隔符解析配置。
type KVDelimiterConfig struct {
	Delimiter   string `json:"delimiter"`
	KVSeparator string `json:"kv_separator"`
}

// ApplyFieldMapping 根据字段映射配置重命名字段。
//
// 支持两种映射格式：
//  1. 简单映射 {"oldField": "newField"}
//  2. 复杂映射 {"targetField": {"source": "header|json|field", "group": 1, "path": "a.b.c"}}
//
// 参数 regex 用于复杂映射中 source="header" 时按分组索引取命名捕获。
// 该函数合并了原 parser.go 的 applyFieldMapping 与 syslog_forward.go 的 applyFieldMapping 逻辑。
func ApplyFieldMapping(data map[string]interface{}, fieldMappingJSON string, regex *regexp.Regexp) map[string]interface{} {
	if fieldMappingJSON == "" {
		return data
	}

	// 尝试复杂映射格式
	var complexMapping map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(fieldMappingJSON), &complexMapping); err == nil {
		return applyComplexFieldMapping(data, complexMapping, regex)
	}

	// 回退到简单映射格式
	var simpleMapping map[string]string
	if err := json.Unmarshal([]byte(fieldMappingJSON), &simpleMapping); err != nil {
		return data
	}
	return applySimpleFieldMapping(data, simpleMapping)
}

// applySimpleFieldMapping 处理 {"oldField": "newField"} 格式的简单映射。
func applySimpleFieldMapping(data map[string]interface{}, mapping map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	// 保留原始数据
	for k, v := range data {
		result[k] = v
	}
	// 应用映射（添加新键，不删除旧键，保持兼容）
	for oldField, newField := range mapping {
		if v, exists := data[oldField]; exists {
			result[newField] = v
		}
	}
	return result
}

// applyComplexFieldMapping 处理 {"targetField": {"source": "...", ...}} 格式的复杂映射。
func applyComplexFieldMapping(data map[string]interface{}, mapping map[string]map[string]interface{}, regex *regexp.Regexp) map[string]interface{} {
	result := make(map[string]interface{})

	for targetField, sourceConfig := range mapping {
		source, ok := sourceConfig["source"].(string)
		if !ok {
			continue
		}

		switch source {
		case "header":
			// 从正则命名分组取值
			if group, ok := sourceConfig["group"].(float64); ok && regex != nil {
				groupIndex := int(group)
				if groupIndex > 0 && groupIndex <= len(regex.SubexpNames()) {
					name := regex.SubexpNames()[groupIndex]
					if val, exists := data[name]; exists {
						result[targetField] = val
					}
				}
			}
		case "json":
			// 按点分路径从嵌套数据取值
			if path, ok := sourceConfig["path"].(string); ok {
				if value := getNestedValue(data, path); value != nil {
					result[targetField] = value
				}
			}
		default:
			// 直接按字段名取值
			if val, exists := data[source]; exists {
				result[targetField] = val
			}
		}
	}

	// 保留未被映射覆盖的原始字段
	for k, v := range data {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}
	return result
}

// ApplyValueTransform 根据值转换规则替换字段值。
//
// 转换规则格式：{"field": {"原始值": "新值", ...}}
// 同时处理 alertTime 字段的 Unix 时间戳（秒/毫秒）转 "2006-01-02 15:04:05"。
//
// 该函数合并了原 parser.go 的 applyValueTransform 与 syslog_forward.go 的 applyValueTransform 逻辑。
// 注意：原 syslog_forward 版本不处理 alertTime 时间戳转换，此处统一处理以保持一致。
// 若调用方不需要 alertTime 转换，应在调用后自行覆盖。
func ApplyValueTransform(data map[string]interface{}, valueTransformJSON string) map[string]interface{} {
	if valueTransformJSON == "" {
		return data
	}

	var transforms map[string]map[string]string
	if err := json.Unmarshal([]byte(valueTransformJSON), &transforms); err != nil {
		return data
	}

	for field, transformMap := range transforms {
		if value, exists := data[field]; exists {
			strValue := fmt.Sprintf("%v", value)
			data[field+"Raw"] = strValue
			if newValue, ok := transformMap[strValue]; ok {
				data[field] = newValue
			}
		}
	}

	// alertTime 时间戳转换（毫秒/秒 → 可读格式）
	convertAlertTimeTimestamp(data)

	return data
}

// convertAlertTimeTimestamp 将 alertTime 中的 Unix 时间戳转为 "2006-01-02 15:04:05"。
// 超过 1e12 视为毫秒，否则视为秒。
func convertAlertTimeTimestamp(data map[string]interface{}) {
	alertTimeVal, exists := data["alertTime"]
	if !exists {
		return
	}
	strVal, ok := alertTimeVal.(string)
	if !ok {
		return
	}
	ts, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return
	}
	data["alertTimeRaw"] = strVal
	if ts > 1e12 {
		ts = ts / 1000
	}
	data["alertTime"] = time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

// ExtractKeyFields 从解析结果中提取关键告警字段，用于日志摘要展示。
// 返回 JSON 字符串，无关键字段时返回空字符串。
func ExtractKeyFields(data map[string]interface{}) string {
	keyFields := make(map[string]interface{})
	fieldNames := []string{
		"attackIp", "victimIp", "threatType", "attack_result", "result",
		"levelDesc", "description", "dealStatus", "threatSource",
		"timestamp", "localTimestamp", "machineName",
	}
	for _, name := range fieldNames {
		if value, exists := data[name]; exists {
			keyFields[name] = value
		}
	}
	if len(keyFields) == 0 {
		return ""
	}
	jsonBytes, _ := json.Marshal(keyFields)
	return string(jsonBytes)
}

// FormatAlertTime 从解析数据中格式化告警时间。
// 优先使用 localTimestamp（毫秒），其次 timestamp，最后返回当前时间。
func FormatAlertTime(data map[string]interface{}) string {
	if ts, ok := data["localTimestamp"]; ok {
		if milli, ok := ts.(float64); ok {
			if milli > 1e12 {
				return time.UnixMilli(int64(milli)).Format("2006-01-02 15:04:05")
			}
			return time.Unix(int64(milli), 0).Format("2006-01-02 15:04:05")
		}
	}
	if ts, ok := data["timestamp"]; ok {
		switch v := ts.(type) {
		case time.Time:
			return v.Format("2006-01-02 15:04:05")
		case string:
			return v
		}
	}
	return time.Now().Format("2006-01-02 15:04:05")
}

// ParseSmartDelimiterConfig 解析智能分隔符配置 JSON。
// 若 Delimiter 为空使用默认值 "|!"。
func ParseSmartDelimiterConfig(fieldMappingJSON string) (*SmartDelimiterConfig, error) {
	config := &SmartDelimiterConfig{
		Delimiter: constants.DefaultDelimiter,
	}
	if fieldMappingJSON == "" {
		return config, nil
	}
	if err := json.Unmarshal([]byte(fieldMappingJSON), &config); err != nil {
		return nil, fmt.Errorf("invalid smart delimiter config: %v", err)
	}
	if config.Delimiter == "" {
		config.Delimiter = constants.DefaultDelimiter
	}
	return config, nil
}

// ParseDelimiterConfig 解析普通分隔符配置 JSON。
// 若 Delimiter 为空使用默认值 "|!"。
func ParseDelimiterConfig(fieldMappingJSON string) (*DelimiterConfig, error) {
	config := &DelimiterConfig{
		Delimiter: constants.DefaultDelimiter,
	}
	if fieldMappingJSON == "" {
		return config, nil
	}
	if err := json.Unmarshal([]byte(fieldMappingJSON), &config); err != nil {
		return nil, fmt.Errorf("invalid delimiter config: %v", err)
	}
	if config.Delimiter == "" {
		config.Delimiter = constants.DefaultDelimiter
	}
	return config, nil
}

// ParseKVDelimiterConfig 解析键值分隔符配置 JSON。
// Delimiter 默认 "|!"，KVSeparator 默认 ":"。
func ParseKVDelimiterConfig(fieldMappingJSON string) (*KVDelimiterConfig, error) {
	config := &KVDelimiterConfig{
		Delimiter:   constants.DefaultDelimiter,
		KVSeparator: constants.DefaultKVSeparator,
	}
	if fieldMappingJSON == "" {
		return config, nil
	}
	if err := json.Unmarshal([]byte(fieldMappingJSON), &config); err != nil {
		return nil, fmt.Errorf("invalid keyvalue config: %v", err)
	}
	if config.Delimiter == "" {
		config.Delimiter = constants.DefaultDelimiter
	}
	if config.KVSeparator == "" {
		config.KVSeparator = constants.DefaultKVSeparator
	}
	return config, nil
}

// assignFieldsByIndex 按字段名列表与值列表按索引赋值。
func assignFieldsByIndex(result map[string]interface{}, fields []string, values []string) {
	for i, field := range fields {
		if i < len(values) {
			result[field] = values[i]
		}
	}
}

// assignFieldsByIndexDefault 无字段名时按 "field_N" 命名赋值。
func assignFieldsByIndexDefault(result map[string]interface{}, values []string) {
	for i, v := range values {
		result[fmt.Sprintf("field_%d", i)] = v
	}
}

// stringsContainsAny 检查字符串是否包含任意一个子串（辅助函数）。
func stringsContainsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
