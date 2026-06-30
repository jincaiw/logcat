package parser

import "strings"

// kvParser 解析 "key=value key=value" 格式的键值对日志。
// 值可以用双引号包裹，解析时会去除引号。
type kvParser struct {
	*baseParser
}

func (p *kvParser) Parse(rawLog string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	pairs := strings.Fields(rawLog)
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, `"`)
			result[key] = value
		}
	}

	if p.template.ValueTransform != "" {
		result = ApplyValueTransform(result, p.template.ValueTransform)
	}
	return result, nil
}
