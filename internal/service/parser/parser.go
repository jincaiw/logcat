// Package parser 提供日志解析能力，采用策略模式支持多种解析类型。
//
// 解析类型包括：syslog_json、json、regex、kv、delimiter、keyvalue、smart_delimiter。
// 每种类型实现 Parser 接口，由 New 工厂函数根据 ParseTemplate.ParseType 分发。
//
// 设计要点：
//   - 策略模式：新增解析类型只需实现 Parser 接口并在 factory 中注册
//   - 公共能力复用：extractJSON/flattenJSON/getNestedValue 等辅助函数集中管理
//   - 字段映射与值转换统一收敛到 transform.go，消除 parser.go 与 syslog_forward.go 的重复
package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/pkg/constants"
)

// Parser 日志解析器接口，每种解析类型实现该接口。
type Parser interface {
	// Parse 解析原始日志，返回结构化字段映射。
	Parse(rawLog string) (map[string]interface{}, error)
}

// New 根据解析模板创建对应的解析器实例。
// 若模板配置了 HeaderRegex，会预编译正则以提升性能。
func New(template *models.ParseTemplate) (Parser, error) {
	regex, err := compileHeaderRegex(template.HeaderRegex)
	if err != nil {
		return nil, fmt.Errorf("invalid header regex: %v", err)
	}

	base := &baseParser{
		template: template,
		regex:    regex,
	}

	switch template.ParseType {
	case constants.ParseTypeSyslogJSON:
		return &syslogJSONParser{base}, nil
	case constants.ParseTypeJSON:
		return &jsonParser{base}, nil
	case constants.ParseTypeRegex:
		return &regexParser{base}, nil
	case constants.ParseTypeKV:
		return &kvParser{base}, nil
	case constants.ParseTypeDelimiter:
		return &delimiterParser{base}, nil
	case constants.ParseTypeKeyValue:
		return &kvDelimiterParser{base}, nil
	case constants.ParseTypeSmartDelimiter:
		return &smartDelimiterParser{base}, nil
	default:
		// 未知类型默认使用 syslog_json，保持与原行为兼容
		return &syslogJSONParser{base}, nil
	}
}

// baseParser 所有解析器共享的基础结构，持有模板与预编译正则。
type baseParser struct {
	template *models.ParseTemplate
	regex    *regexp.Regexp
}

// compileHeaderRegex 编译头部正则，空字符串返回 nil。
func compileHeaderRegex(pattern string) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, nil
	}
	return regexp.Compile(pattern)
}

// extractHeaderFields 使用预编译正则提取头部命名分组，并返回匹配结束位置。
// 若未匹配返回 (-1, nil)。
func (p *baseParser) extractHeaderFields(rawLog string) (int, map[string]interface{}) {
	if p.regex == nil {
		return -1, nil
	}
	matches := p.regex.FindStringSubmatch(rawLog)
	if matches == nil {
		return -1, nil
	}
	result := make(map[string]interface{})
	subexpNames := p.regex.SubexpNames()
	for i, name := range subexpNames {
		if name != "" && i < len(matches) {
			result[name] = matches[i]
		}
	}
	loc := p.regex.FindStringIndex(rawLog)
	if loc == nil {
		return -1, result
	}
	return loc[1], result
}

// applyTransforms 统一应用字段映射与值转换，消除各解析器中的重复调用。
func (p *baseParser) applyTransforms(result map[string]interface{}) map[string]interface{} {
	if p.template.FieldMapping != "" {
		result = ApplyFieldMapping(result, p.template.FieldMapping, p.regex)
	}
	if p.template.ValueTransform != "" {
		result = ApplyValueTransform(result, p.template.ValueTransform)
	}
	return result
}

// ---- 公共辅助函数 ----

// extractJSON 从字符串中提取第一个完整的 JSON 对象或数组。
// 处理 JSON 后混入非 JSON 文本的情况，通过括号深度匹配截取。
func extractJSON(str string) string {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return str
	}
	if str[0] != '{' && str[0] != '[' {
		return str
	}

	depth := 0
	inString := false
	escape := false

	for i, c := range str {
		if escape {
			escape = false
			continue
		}
		switch c {
		case '\\':
			if inString {
				escape = true
			}
		case '"':
			inString = !inString
		case '{', '[':
			if !inString {
				depth++
			}
		case '}', ']':
			if !inString {
				depth--
				if depth == 0 {
					return str[:i+1]
				}
			}
		}
	}
	return str
}

// fixMalformedJSON 修复包含非法嵌套 JSON 字符串的字段（如 fullTree 字段）。
// 某些设备会将 JSON 数组作为字符串值嵌入，导致外层 JSON 解析失败，
// 此函数移除这类问题字段以恢复外层 JSON 的合法性。
func fixMalformedJSON(jsonStr string) string {
	result := jsonStr
	for {
		idx := strings.Index(result, `"fullTree":"`)
		if idx == -1 {
			break
		}
		valueStart := idx + len(`"fullTree":"`)
		if valueStart >= len(result) || result[valueStart] != '[' {
			break
		}

		depth := 0
		valueEnd := -1
		for i := valueStart; i < len(result); i++ {
			c := result[i]
			if c == '[' {
				depth++
			} else if c == ']' {
				depth--
				if depth == 0 {
					// 匹配 ] 后跟 " 或 ]\n"
					if i+1 < len(result) && result[i+1] == '"' {
						valueEnd = i + 2
						break
					}
					if i+3 < len(result) && result[i+1] == '\\' && result[i+2] == 'n' && result[i+3] == '"' {
						valueEnd = i + 4
						break
					}
				}
			}
		}
		if valueEnd == -1 {
			break
		}

		before := result[:idx]
		after := result[valueEnd:]
		// 移除字段前后的逗号
		if len(before) > 0 && before[len(before)-1] == ',' {
			before = before[:len(before)-1]
		} else if len(after) > 0 && after[0] == ',' {
			after = after[1:]
		}
		result = before + after
	}
	return result
}

// flattenJSON 将嵌套的 JSON map 扁平化，键用 "." 连接。
// 例如 {"a": {"b": 1}} -> {"a.b": 1}
func flattenJSON(data map[string]interface{}, prefix string, result map[string]interface{}) {
	for k, v := range data {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		if nested, ok := v.(map[string]interface{}); ok {
			flattenJSON(nested, key, result)
		} else {
			result[key] = v
		}
	}
}

// getNestedValue 按点分路径从嵌套 map 中取值。
// 例如 getNestedValue(data, "a.b.c") 等价于 data["a"]["b"]["c"]。
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

// ParseTimestamp 将多种格式的时间戳解析为 time.Time。
// 支持 float64（Unix 秒/毫秒）、字符串（多种布局、Unix 秒/毫秒）。
// 解析失败返回当前时间，保持与原行为兼容。
func ParseTimestamp(ts interface{}) time.Time {
	switch v := ts.(type) {
	case float64:
		if v > 1e12 {
			return time.UnixMilli(int64(v))
		}
		return time.Unix(int64(v), 0)
	case string:
		layouts := []string{
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.999Z",
			"Jan 02 15:04:05",
			time.RFC3339,
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, v); err == nil {
				if layout == "Jan 02 15:04:05" {
					t = t.AddDate(time.Now().Year(), 0, 0)
				}
				return t
			}
		}
		if milli, err := strconv.ParseInt(v, 10, 64); err == nil {
			if milli > 1e12 {
				return time.UnixMilli(milli)
			}
			return time.Unix(milli, 0)
		}
	}
	return time.Now()
}

// convertSyslogTimestamp 将 syslog 时间戳（如 "Jan  2 15:04:05"）转为 "2006-01-02 15:04:05"。
// syslog 时间戳不含年份，此处补当前年份。
func convertSyslogTimestamp(ts string) string {
	layouts := []string{
		"Jan _2 15:04:05",
		"Jan 02 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, ts); err == nil {
			t = t.AddDate(time.Now().Year(), 0, 0)
			return t.Format("2006-01-02 15:04:05")
		}
	}
	return ts
}

// parseJSONToMap 解析 JSON 字符串为 map，失败时尝试修复后重试。
func parseJSONToMap(jsonStr string) (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		fixed := fixMalformedJSON(jsonStr)
		if fixedErr := json.Unmarshal([]byte(fixed), &jsonData); fixedErr != nil {
			return nil, fmt.Errorf("failed to parse JSON: %v", err)
		}
	}
	return jsonData, nil
}
