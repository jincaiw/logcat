package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// syslogJSONParser 先提取 syslog 头部，再解析剩余部分的 JSON。
// 典型场景：syslog 消息体为 JSON 格式的安全告警。
type syslogJSONParser struct {
	*baseParser
}

// defaultSyslogTimeRegex 无 HeaderRegex 时的回退方案，匹配 syslog 时间戳。
var defaultSyslogTimeRegex = regexp.MustCompile(`^<\d+>(\w{3}\s+\d{1,2}\s+[\d:]+)`)

func (p *syslogJSONParser) Parse(rawLog string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var jsonStart int

	if p.regex != nil {
		// 使用配置的 HeaderRegex 提取头部字段
		end, headerFields := p.extractHeaderFields(rawLog)
		if headerFields == nil {
			return nil, fmt.Errorf("log does not match header pattern")
		}
		for k, v := range headerFields {
			result[k] = v
		}
		jsonStart = end
	} else {
		// 无 HeaderRegex 时，尝试匹配默认 syslog 时间戳并定位 JSON 起始
		if matches := defaultSyslogTimeRegex.FindStringSubmatch(rawLog); matches != nil {
			result["timestamp"] = matches[1]
		}
		jsonStart = strings.Index(rawLog, "{")
	}

	// 截取 JSON 部分
	jsonStr := rawLog
	if jsonStart > 0 && jsonStart < len(rawLog) {
		jsonStr = strings.TrimSpace(rawLog[jsonStart:])
	}
	jsonStr = extractJSON(jsonStr)

	// 解析 JSON
	jsonData, err := parseJSONToMap(jsonStr)
	if err != nil {
		return nil, err
	}
	flattenJSON(jsonData, "", result)

	// 转换 syslog 时间戳为标准格式
	if ts, ok := result["timestamp"].(string); ok && ts != "" {
		result["alertTime"] = convertSyslogTimestamp(ts)
	}

	return p.applyTransforms(result), nil
}
