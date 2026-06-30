package parser

import "fmt"

// regexParser 使用正则表达式命名分组提取字段。
// 要求 HeaderRegex 配置了命名捕获组（如 (?P<field>...)）。
type regexParser struct {
	*baseParser
}

func (p *regexParser) Parse(rawLog string) (map[string]interface{}, error) {
	if p.regex == nil {
		return nil, fmt.Errorf("no regex pattern configured")
	}

	matches := p.regex.FindStringSubmatch(rawLog)
	if matches == nil {
		return nil, fmt.Errorf("log does not match pattern")
	}

	result := make(map[string]interface{})
	subexpNames := p.regex.SubexpNames()
	for i, name := range subexpNames {
		if name != "" && i < len(matches) {
			result[name] = matches[i]
		}
	}

	// regex 类型只应用值转换，不应用字段映射（字段已由命名分组直接命名）
	if p.template.ValueTransform != "" {
		result = ApplyValueTransform(result, p.template.ValueTransform)
	}
	return result, nil
}
