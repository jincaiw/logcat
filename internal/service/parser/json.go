package parser

// jsonParser 解析纯 JSON 日志，无 syslog 头部。
type jsonParser struct {
	*baseParser
}

func (p *jsonParser) Parse(rawLog string) (map[string]interface{}, error) {
	jsonData, err := parseJSONToMap(rawLog)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	flattenJSON(jsonData, "", result)

	return p.applyTransforms(result), nil
}
