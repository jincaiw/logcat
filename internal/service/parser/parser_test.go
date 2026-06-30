package parser

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/pkg/constants"
)

// ==================== 工厂函数测试 ====================

// TestNew_FactoryDispatch 验证 New 工厂函数根据 ParseType 返回正确的解析器类型。
func TestNew_FactoryDispatch(t *testing.T) {
	tests := []struct {
		name      string
		parseType string
	}{
		{"syslog_json", constants.ParseTypeSyslogJSON},
		{"json", constants.ParseTypeJSON},
		{"regex", constants.ParseTypeRegex},
		{"kv", constants.ParseTypeKV},
		{"delimiter", constants.ParseTypeDelimiter},
		{"keyvalue", constants.ParseTypeKeyValue},
		{"smart_delimiter", constants.ParseTypeSmartDelimiter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := &models.ParseTemplate{ParseType: tt.parseType}
			p, err := New(tmpl)
			if err != nil {
				t.Fatalf("New(%s) error: %v", tt.name, err)
			}
			if p == nil {
				t.Fatalf("New(%s) returned nil parser", tt.name)
			}
		})
	}
}

// TestNew_UnknownTypeFallback 验证未知 ParseType 回退到 syslog_json。
func TestNew_UnknownTypeFallback(t *testing.T) {
	tmpl := &models.ParseTemplate{ParseType: "unknown_type"}
	p, err := New(tmpl)
	if err != nil {
		t.Fatalf("New(unknown) error: %v", err)
	}
	if p == nil {
		t.Fatal("New(unknown) returned nil")
	}
	// 应该能解析 syslog_json 格式
	result, err := p.Parse(`<134>Jan 02 15:04:05 host {"attackIp":"1.2.3.4"}`)
	if err != nil {
		t.Fatalf("fallback parser error: %v", err)
	}
	if result["attackIp"] != "1.2.3.4" {
		t.Errorf("expected attackIp=1.2.3.4, got %v", result["attackIp"])
	}
}

// TestNew_InvalidHeaderRegex 验证无效 HeaderRegex 返回错误。
func TestNew_InvalidHeaderRegex(t *testing.T) {
	tmpl := &models.ParseTemplate{
		ParseType:   constants.ParseTypeSyslogJSON,
		HeaderRegex: `[invalid(`, // 未闭合的括号
	}
	_, err := New(tmpl)
	if err == nil {
		t.Error("无效 HeaderRegex 应返回错误")
	}
}

// ==================== JSON 解析器测试 ====================

// TestJSONParser_Basic 验证纯 JSON 解析与扁平化。
func TestJSONParser_Basic(t *testing.T) {
	tmpl := &models.ParseTemplate{ParseType: constants.ParseTypeJSON}
	p, _ := New(tmpl)

	raw := `{"attackIp":"192.168.1.1","port":443,"nested":{"key":"value"}}`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result["attackIp"] != "192.168.1.1" {
		t.Errorf("attackIp = %v, want 192.168.1.1", result["attackIp"])
	}
	// 扁平化后嵌套字段应为 nested.key
	if result["nested.key"] != "value" {
		t.Errorf("nested.key = %v, want value", result["nested.key"])
	}
}

// TestJSONParser_InvalidJSON 验证无效 JSON 返回错误。
func TestJSONParser_InvalidJSON(t *testing.T) {
	tmpl := &models.ParseTemplate{ParseType: constants.ParseTypeJSON}
	p, _ := New(tmpl)

	_, err := p.Parse(`{invalid json}`)
	if err == nil {
		t.Error("无效 JSON 应返回错误")
	}
}

// TestJSONParser_EmptyInput 验证空输入处理。
func TestJSONParser_EmptyInput(t *testing.T) {
	tmpl := &models.ParseTemplate{ParseType: constants.ParseTypeJSON}
	p, _ := New(tmpl)

	_, err := p.Parse("")
	if err == nil {
		t.Error("空输入应返回错误")
	}
}

// ==================== Syslog JSON 解析器测试 ====================

// TestSyslogJSONParser_WithHeaderRegex 验证带 HeaderRegex 的 syslog_json 解析。
func TestSyslogJSONParser_WithHeaderRegex(t *testing.T) {
	tmpl := &models.ParseTemplate{
		ParseType:   constants.ParseTypeSyslogJSON,
		HeaderRegex: `<(?P<priority>[0-9]+)>(?P<timestamp>[A-Za-z]+[ ]+[0-9]+ [0-9:]+) (?P<hostname>[^ ]+) (?P<program>[^:]+):`,
	}
	p, _ := New(tmpl)

	raw := `<134>Jan 02 15:04:05 myhost myapp: {"attackIp":"10.0.0.1","alertName":"测试告警"}`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result["hostname"] != "myhost" {
		t.Errorf("hostname = %v, want myhost", result["hostname"])
	}
	if result["program"] != "myapp" {
		t.Errorf("program = %v, want myapp", result["program"])
	}
	if result["attackIp"] != "10.0.0.1" {
		t.Errorf("attackIp = %v, want 10.0.0.1", result["attackIp"])
	}
}

// TestSyslogJSONParser_DefaultRegex 验证无 HeaderRegex 时的默认 syslog 时间戳匹配。
func TestSyslogJSONParser_DefaultRegex(t *testing.T) {
	tmpl := &models.ParseTemplate{ParseType: constants.ParseTypeSyslogJSON}
	p, _ := New(tmpl)

	raw := `<134>Jan 02 15:04:05 {"attackIp":"172.16.0.1"}`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result["attackIp"] != "172.16.0.1" {
		t.Errorf("attackIp = %v, want 172.16.0.1", result["attackIp"])
	}
	// 默认正则应提取 timestamp
	if ts, ok := result["timestamp"].(string); !ok || ts == "" {
		t.Error("timestamp 应被默认正则提取")
	}
}

// TestSyslogJSONParser_HeaderMismatch 验证日志不匹配 HeaderRegex 时返回错误。
func TestSyslogJSONParser_HeaderMismatch(t *testing.T) {
	tmpl := &models.ParseTemplate{
		ParseType:   constants.ParseTypeSyslogJSON,
		HeaderRegex: `<(?P<priority>[0-9]+)>custom`,
	}
	p, _ := New(tmpl)

	raw := `no match here {"key":"value"}`
	_, err := p.Parse(raw)
	if err == nil {
		t.Error("不匹配 HeaderRegex 应返回错误")
	}
}

// TestSyslogJSONParser_TimestampConversion 验证 syslog 时间戳被转换为标准格式。
func TestSyslogJSONParser_TimestampConversion(t *testing.T) {
	tmpl := &models.ParseTemplate{ParseType: constants.ParseTypeSyslogJSON}
	p, _ := New(tmpl)

	raw := `<134>Jan 02 15:04:05 {"attackIp":"1.1.1.1"}`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	alertTime, ok := result["alertTime"].(string)
	if !ok {
		t.Fatal("alertTime 应为字符串")
	}
	// 应包含当前年份
	expectedPrefix := time.Now().Format("2006") + "-"
	if len(alertTime) < len(expectedPrefix) || alertTime[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("alertTime = %q, 应以 %q 开头", alertTime, expectedPrefix)
	}
}

// ==================== Regex 解析器测试 ====================

// TestRegexParser_NamedGroups 验证正则命名分组提取。
func TestRegexParser_NamedGroups(t *testing.T) {
	tmpl := &models.ParseTemplate{
		ParseType:   constants.ParseTypeRegex,
		HeaderRegex: `(?P<ip>\d+\.\d+\.\d+\.\d+) - - \[(?P<time>[^\]]+)\] "(?P<method>\w+) (?P<path>\S+)"`,
	}
	p, _ := New(tmpl)

	raw := `192.168.1.100 - - [10/Oct/2024:13:55:36 +0800] "GET /api/users"`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result["ip"] != "192.168.1.100" {
		t.Errorf("ip = %v, want 192.168.1.100", result["ip"])
	}
	if result["method"] != "GET" {
		t.Errorf("method = %v, want GET", result["method"])
	}
	if result["path"] != "/api/users" {
		t.Errorf("path = %v, want /api/users", result["path"])
	}
}

// TestRegexParser_NoPattern 验证未配置正则时返回错误。
func TestRegexParser_NoPattern(t *testing.T) {
	tmpl := &models.ParseTemplate{ParseType: constants.ParseTypeRegex}
	p, _ := New(tmpl)

	_, err := p.Parse("some log")
	if err == nil {
		t.Error("未配置正则应返回错误")
	}
}

// TestRegexParser_NoMatch 验证不匹配时返回错误。
func TestRegexParser_NoMatch(t *testing.T) {
	tmpl := &models.ParseTemplate{
		ParseType:   constants.ParseTypeRegex,
		HeaderRegex: `(?P<ip>\d+\.\d+\.\d+\.\d+)`,
	}
	p, _ := New(tmpl)

	_, err := p.Parse("no ip address here")
	if err == nil {
		t.Error("不匹配正则应返回错误")
	}
}

// ==================== KV 解析器测试 ====================

// TestKVParser_Basic 验证 key=value 格式解析。
func TestKVParser_Basic(t *testing.T) {
	tmpl := &models.ParseTemplate{ParseType: constants.ParseTypeKV}
	p, _ := New(tmpl)

	// 注意：KV 解析器使用 strings.Fields 按空格分割，不支持值中包含空格
	raw := `src=192.168.1.1 dst=10.0.0.1 action=block msg="attack_detected"`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result["src"] != "192.168.1.1" {
		t.Errorf("src = %v, want 192.168.1.1", result["src"])
	}
	if result["dst"] != "10.0.0.1" {
		t.Errorf("dst = %v, want 10.0.0.1", result["dst"])
	}
	if result["action"] != "block" {
		t.Errorf("action = %v, want block", result["action"])
	}
	// 引号应被去除
	if result["msg"] != "attack_detected" {
		t.Errorf("msg = %v, want attack_detected (引号应被去除)", result["msg"])
	}
}

// TestKVParser_EmptyInput 验证空输入返回空 map（不报错）。
func TestKVParser_EmptyInput(t *testing.T) {
	tmpl := &models.ParseTemplate{ParseType: constants.ParseTypeKV}
	p, _ := New(tmpl)

	result, err := p.Parse("")
	if err != nil {
		t.Fatalf("空输入不应报错: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("空输入应返回空 map, got %v", result)
	}
}

// ==================== Delimiter 解析器测试 ====================

// TestDelimiterParser_Fields 验证按字段名列表解析。
func TestDelimiterParser_Fields(t *testing.T) {
	tmpl := &models.ParseTemplate{
		ParseType:    constants.ParseTypeDelimiter,
		FieldMapping: `{"delimiter":"|!","fields":["alertName","attackIP","severity"]}`,
	}
	p, _ := New(tmpl)

	raw := `SQL注入|!192.168.1.100|!high`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result["alertName"] != "SQL注入" {
		t.Errorf("alertName = %v, want SQL注入", result["alertName"])
	}
	if result["attackIP"] != "192.168.1.100" {
		t.Errorf("attackIP = %v, want 192.168.1.100", result["attackIP"])
	}
	if result["severity"] != "high" {
		t.Errorf("severity = %v, want high", result["severity"])
	}
}

// TestDelimiterParser_TypeMapping 验证按类型映射解析。
// 注意：type_mapping 的字段列表从 index 0 开始赋值（包含类型字段本身），
// 因此第一个字段名应为类型字段名。
func TestDelimiterParser_TypeMapping(t *testing.T) {
	tmpl := &models.ParseTemplate{
		ParseType: constants.ParseTypeDelimiter,
		FieldMapping: `{
			"delimiter":"|!",
			"type_field":"alertType",
			"type_mapping": {
				"ioc": ["alertType", "alertName", "attackIP", "victimIP"]
			}
		}`,
	}
	p, _ := New(tmpl)

	// IOC 类型
	raw := `ioc|!恶意IP|!1.2.3.4|!5.6.7.8`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result["alertType"] != "ioc" {
		t.Errorf("alertType = %v, want ioc", result["alertType"])
	}
	if result["alertName"] != "恶意IP" {
		t.Errorf("alertName = %v, want 恶意IP", result["alertName"])
	}
	if result["attackIP"] != "1.2.3.4" {
		t.Errorf("attackIP = %v, want 1.2.3.4", result["attackIP"])
	}
	if result["victimIP"] != "5.6.7.8" {
		t.Errorf("victimIP = %v, want 5.6.7.8", result["victimIP"])
	}
}

// TestDelimiterParser_DefaultFieldNames 验证无字段名时使用 field_N 命名。
func TestDelimiterParser_DefaultFieldNames(t *testing.T) {
	tmpl := &models.ParseTemplate{
		ParseType:    constants.ParseTypeDelimiter,
		FieldMapping: `{"delimiter":","}`,
	}
	p, _ := New(tmpl)

	raw := `a,b,c`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result["field_0"] != "a" {
		t.Errorf("field_0 = %v, want a", result["field_0"])
	}
	if result["field_1"] != "b" {
		t.Errorf("field_1 = %v, want b", result["field_1"])
	}
	if result["field_2"] != "c" {
		t.Errorf("field_2 = %v, want c", result["field_2"])
	}
}

// TestDelimiterParser_SimpleMapping 验证简单字段映射格式。
func TestDelimiterParser_SimpleMapping(t *testing.T) {
	tmpl := &models.ParseTemplate{
		ParseType:    constants.ParseTypeDelimiter,
		FieldMapping: `{"field_0":"alertName","field_1":"attackIP"}`,
	}
	p, _ := New(tmpl)

	raw := `测试告警|!10.0.0.1`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result["alertName"] != "测试告警" {
		t.Errorf("alertName = %v, want 测试告警", result["alertName"])
	}
	if result["attackIP"] != "10.0.0.1" {
		t.Errorf("attackIP = %v, want 10.0.0.1", result["attackIP"])
	}
}

// ==================== KV Delimiter 解析器测试 ====================

// TestKVDelimiterParser_Basic 验证 key:val|!key:val 格式解析。
func TestKVDelimiterParser_Basic(t *testing.T) {
	tmpl := &models.ParseTemplate{
		ParseType:    constants.ParseTypeKeyValue,
		FieldMapping: `{"delimiter":"|!","kv_separator":":"}`,
	}
	p, _ := New(tmpl)

	raw := `alertName:测试告警|!attackIP:192.168.1.1|!severity:high`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result["alertName"] != "测试告警" {
		t.Errorf("alertName = %v, want 测试告警", result["alertName"])
	}
	if result["attackIP"] != "192.168.1.1" {
		t.Errorf("attackIP = %v, want 192.168.1.1", result["attackIP"])
	}
	if result["severity"] != "high" {
		t.Errorf("severity = %v, want high", result["severity"])
	}
}

// TestKVDelimiterParser_DefaultConfig 验证无配置时使用默认分隔符。
func TestKVDelimiterParser_DefaultConfig(t *testing.T) {
	tmpl := &models.ParseTemplate{ParseType: constants.ParseTypeKeyValue}
	p, _ := New(tmpl)

	// 默认 delimiter=|!, kv_separator=:
	raw := `key1:val1|!key2:val2`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result["key1"] != "val1" {
		t.Errorf("key1 = %v, want val1", result["key1"])
	}
	if result["key2"] != "val2" {
		t.Errorf("key2 = %v, want val2", result["key2"])
	}
}

// ==================== Smart Delimiter 解析器测试 ====================

// TestSmartDelimiterParser_WithSubTemplate 验证智能分隔符按子模板解析。
func TestSmartDelimiterParser_WithSubTemplate(t *testing.T) {
	config := `{
		"delimiter": "|!",
		"typeField": 0,
		"skipHeader": false,
		"subTemplates": {
			"ioc": {
				"alertNameField": 1,
				"attackIPField": 2,
				"victimIPField": 3,
				"alertTimeField": 4,
				"severityField": 5,
				"attackResultField": -1
			}
		}
	}`
	tmpl := &models.ParseTemplate{
		ParseType:    constants.ParseTypeSmartDelimiter,
		FieldMapping: config,
	}
	p, _ := New(tmpl)

	raw := `ioc|!恶意IP通信|!1.2.3.4|!5.6.7.8|!2024-01-01 10:00:00|!high`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result["alertType"] != "ioc" {
		t.Errorf("alertType = %v, want ioc", result["alertType"])
	}
	if result["alertName"] != "恶意IP通信" {
		t.Errorf("alertName = %v, want 恶意IP通信", result["alertName"])
	}
	if result["attackIP"] != "1.2.3.4" {
		t.Errorf("attackIP = %v, want 1.2.3.4", result["attackIP"])
	}
	if result["victimIP"] != "5.6.7.8" {
		t.Errorf("victimIP = %v, want 5.6.7.8", result["victimIP"])
	}
}

// TestSmartDelimiterParser_SkipHeader 验证 SkipHeader 跳过 syslog 头部。
func TestSmartDelimiterParser_SkipHeader(t *testing.T) {
	config := `{
		"delimiter": "|!",
		"typeField": 0,
		"skipHeader": true,
		"subTemplates": {
			"alert": {
				"alertNameField": 1,
				"attackIPField": 2,
				"alertTimeField": -1,
				"severityField": -1,
				"attackResultField": -1,
				"victimIPField": -1
			}
		}
	}`
	tmpl := &models.ParseTemplate{
		ParseType:    constants.ParseTypeSmartDelimiter,
		FieldMapping: config,
	}
	p, _ := New(tmpl)

	raw := `<134>Jan 02 15:04:05 myhost myapp: alert|!测试告警|!10.0.0.1`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result["alertType"] != "alert" {
		t.Errorf("alertType = %v, want alert", result["alertType"])
	}
	if result["alertName"] != "测试告警" {
		t.Errorf("alertName = %v, want 测试告警", result["alertName"])
	}
	if result["attackIP"] != "10.0.0.1" {
		t.Errorf("attackIP = %v, want 10.0.0.1", result["attackIP"])
	}
	// 头部字段应被提取
	if result["hostname"] != "myhost" {
		t.Errorf("hostname = %v, want myhost", result["hostname"])
	}
}

// TestSmartDelimiterParser_IOCDefaultAttackResult 验证 IOC 告警默认攻击结果为"失陷"。
func TestSmartDelimiterParser_IOCDefaultAttackResult(t *testing.T) {
	config := `{
		"delimiter": "|!",
		"typeField": 0,
		"skipHeader": false,
		"subTemplates": {
			"ioc_alert": {
				"alertNameField": 1,
				"attackIPField": 2,
				"victimIPField": -1,
				"alertTimeField": -1,
				"severityField": -1,
				"attackResultField": -1
			}
		}
	}`
	tmpl := &models.ParseTemplate{
		ParseType:    constants.ParseTypeSmartDelimiter,
		FieldMapping: config,
	}
	p, _ := New(tmpl)

	raw := `ioc_alert|!IOC告警|!1.1.1.1`
	result, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result["attackResult"] != constants.AttackResultCompromised {
		t.Errorf("attackResult = %v, want %v", result["attackResult"], constants.AttackResultCompromised)
	}
}

// TestSmartDelimiterParser_NotEnoughFields 验证字段不足时返回错误。
func TestSmartDelimiterParser_NotEnoughFields(t *testing.T) {
	config := `{
		"delimiter": "|!",
		"typeField": 5,
		"skipHeader": false
	}`
	tmpl := &models.ParseTemplate{
		ParseType:    constants.ParseTypeSmartDelimiter,
		FieldMapping: config,
	}
	p, _ := New(tmpl)

	raw := `only|!two|!fields`
	_, err := p.Parse(raw)
	if err == nil {
		t.Error("字段不足时应返回错误")
	}
}

// ==================== 辅助函数测试 ====================

// TestExtractJSON 验证从字符串中提取 JSON。
func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"纯JSON对象", `{"key":"value"}`, `{"key":"value"}`},
		{"JSON后有多余文本", `{"key":"value"} extra text`, `{"key":"value"}`},
		{"嵌套JSON", `{"a":{"b":1}}`, `{"a":{"b":1}}`},
		{"非JSON", `plain text`, `plain text`},
		{"空字符串", ``, ``},
		{"JSON数组", `[1,2,3]`, `[1,2,3]`},
		{"带转义引号", `{"key":"val\"ue"}`, `{"key":"val\"ue"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractJSON(tt.input)
			if got != tt.expected {
				t.Errorf("extractJSON(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestFlattenJSON 验证 JSON 扁平化。
func TestFlattenJSON(t *testing.T) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": 1,
			"c": map[string]interface{}{
				"d": 2,
			},
		},
		"e": "value",
	}

	result := make(map[string]interface{})
	flattenJSON(data, "", result)

	if result["a.b"] != 1 {
		t.Errorf("a.b = %v, want 1", result["a.b"])
	}
	if result["a.c.d"] != 2 {
		t.Errorf("a.c.d = %v, want 2", result["a.c.d"])
	}
	if result["e"] != "value" {
		t.Errorf("e = %v, want value", result["e"])
	}
}

// TestGetNestedValue 验证按点分路径取嵌套值。
func TestGetNestedValue(t *testing.T) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "deep",
			},
		},
	}

	tests := []struct {
		path     string
		expected interface{}
	}{
		{"a.b.c", "deep"},
		{"a.x", nil},
		{"x.y.z", nil},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := getNestedValue(data, tt.path)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("getNestedValue(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}

	// 单独验证返回 map 的情况
	t.Run("a.b returns map", func(t *testing.T) {
		got := getNestedValue(data, "a.b")
		expectedMap := map[string]interface{}{"c": "deep"}
		if !reflect.DeepEqual(got, expectedMap) {
			t.Errorf("getNestedValue(\"a.b\") = %v, want %v", got, expectedMap)
		}
	})
}

// TestParseTimestamp 验证多种时间戳格式解析。
func TestParseTimestamp(t *testing.T) {
	// float64 Unix 秒
	t1 := ParseTimestamp(float64(1700000000))
	if t1.Unix() != 1700000000 {
		t.Errorf("Unix秒解析错误: got %v", t1.Unix())
	}

	// float64 Unix 毫秒
	t2 := ParseTimestamp(float64(1700000000000))
	if t2.Unix() != 1700000000 {
		t.Errorf("Unix毫秒解析错误: got %v", t2.Unix())
	}

	// 字符串标准格式
	t3 := ParseTimestamp("2024-01-01 12:00:00")
	expected := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	if !t3.Equal(expected) {
		t.Errorf("标准格式解析错误: got %v, want %v", t3, expected)
	}

	// 字符串 Unix 秒
	t4 := ParseTimestamp("1700000000")
	if t4.Unix() != 1700000000 {
		t.Errorf("字符串Unix秒解析错误: got %v", t4.Unix())
	}

	// 无效格式返回当前时间（不 panic）
	t5 := ParseTimestamp("invalid")
	if t5.IsZero() {
		t.Error("无效格式应返回当前时间，不应为零值")
	}
}

// TestConvertSyslogTimestamp 验证 syslog 时间戳转换。
func TestConvertSyslogTimestamp(t *testing.T) {
	// 标准 syslog 格式 "Jan  2" (两个空格)
	result := convertSyslogTimestamp("Jan  2 15:04:05")
	expectedPrefix := time.Now().Format("2006") + "-01-02 15:04:05"
	if result != expectedPrefix {
		t.Errorf("convertSyslogTimestamp('Jan  2 15:04:05') = %q, want %q", result, expectedPrefix)
	}

	// 单数字日期 "Jan 02"
	result2 := convertSyslogTimestamp("Jan 02 15:04:05")
	if result2 != expectedPrefix {
		t.Errorf("convertSyslogTimestamp('Jan 02 15:04:05') = %q, want %q", result2, expectedPrefix)
	}

	// 无法解析的格式原样返回
	invalid := "invalid time"
	result3 := convertSyslogTimestamp(invalid)
	if result3 != invalid {
		t.Errorf("convertSyslogTimestamp(%q) = %q, want 原样返回", invalid, result3)
	}
}

// ==================== 配置解析函数测试 ====================

// TestParseDelimiterConfig 验证分隔符配置解析。
func TestParseDelimiterConfig(t *testing.T) {
	// 空配置使用默认值
	config, err := ParseDelimiterConfig("")
	if err != nil {
		t.Fatalf("ParseDelimiterConfig(\"\") error: %v", err)
	}
	if config.Delimiter != constants.DefaultDelimiter {
		t.Errorf("默认 delimiter = %q, want %q", config.Delimiter, constants.DefaultDelimiter)
	}

	// 正常配置
	config, err = ParseDelimiterConfig(`{"delimiter":",","fields":["a","b"]}`)
	if err != nil {
		t.Fatalf("ParseDelimiterConfig error: %v", err)
	}
	if config.Delimiter != "," {
		t.Errorf("delimiter = %q, want ,", config.Delimiter)
	}
	if len(config.Fields) != 2 {
		t.Errorf("fields length = %d, want 2", len(config.Fields))
	}

	// 无效 JSON
	_, err = ParseDelimiterConfig(`invalid`)
	if err == nil {
		t.Error("无效 JSON 应返回错误")
	}
}

// TestParseKVDelimiterConfig 验证键值分隔符配置解析。
func TestParseKVDelimiterConfig(t *testing.T) {
	// 空配置使用默认值
	config, err := ParseKVDelimiterConfig("")
	if err != nil {
		t.Fatalf("ParseKVDelimiterConfig(\"\") error: %v", err)
	}
	if config.Delimiter != constants.DefaultDelimiter {
		t.Errorf("默认 delimiter = %q, want %q", config.Delimiter, constants.DefaultDelimiter)
	}
	if config.KVSeparator != constants.DefaultKVSeparator {
		t.Errorf("默认 kv_separator = %q, want %q", config.KVSeparator, constants.DefaultKVSeparator)
	}

	// 自定义配置
	config, err = ParseKVDelimiterConfig(`{"delimiter":";","kv_separator":"="}`)
	if err != nil {
		t.Fatalf("ParseKVDelimiterConfig error: %v", err)
	}
	if config.Delimiter != ";" {
		t.Errorf("delimiter = %q, want ;", config.Delimiter)
	}
	if config.KVSeparator != "=" {
		t.Errorf("kv_separator = %q, want =", config.KVSeparator)
	}
}

// TestParseSmartDelimiterConfig 验证智能分隔符配置解析。
func TestParseSmartDelimiterConfig(t *testing.T) {
	// 空配置使用默认值
	config, err := ParseSmartDelimiterConfig("")
	if err != nil {
		t.Fatalf("ParseSmartDelimiterConfig(\"\") error: %v", err)
	}
	if config.Delimiter != constants.DefaultDelimiter {
		t.Errorf("默认 delimiter = %q, want %q", config.Delimiter, constants.DefaultDelimiter)
	}

	// 正常配置
	jsonStr := `{"delimiter":"|!","typeField":0,"skipHeader":true,"subTemplates":{"ioc":{"alertNameField":1}}}`
	config, err = ParseSmartDelimiterConfig(jsonStr)
	if err != nil {
		t.Fatalf("ParseSmartDelimiterConfig error: %v", err)
	}
	if !config.SkipHeader {
		t.Error("skipHeader 应为 true")
	}
	if _, ok := config.SubTemplates["ioc"]; !ok {
		t.Error("应包含 ioc 子模板")
	}
}

// ==================== 字段映射与值转换测试 ====================

// TestApplyFieldMapping_Simple 验证简单字段映射。
func TestApplyFieldMapping_Simple(t *testing.T) {
	data := map[string]interface{}{
		"oldField": "value1",
		"keep":     "value2",
	}
	mapping := `{"oldField":"newField"}`

	result := ApplyFieldMapping(data, mapping, nil)

	// 新字段应存在
	if result["newField"] != "value1" {
		t.Errorf("newField = %v, want value1", result["newField"])
	}
	// 旧字段应保留（兼容性设计）
	if result["oldField"] != "value1" {
		t.Errorf("oldField 应保留 = %v", result["oldField"])
	}
}

// TestApplyFieldMapping_Complex 验证复杂字段映射（按 source 取值）。
func TestApplyFieldMapping_Complex(t *testing.T) {
	data := map[string]interface{}{
		"src_ip": "1.1.1.1",
		"nested": map[string]interface{}{"deep": "deepValue"},
	}
	mapping := `{
		"attackIP": {"source":"src_ip"},
		"deepField": {"source":"json","path":"nested.deep"}
	}`

	result := ApplyFieldMapping(data, mapping, nil)

	if result["attackIP"] != "1.1.1.1" {
		t.Errorf("attackIP = %v, want 1.1.1.1", result["attackIP"])
	}
	if result["deepField"] != "deepValue" {
		t.Errorf("deepField = %v, want deepValue", result["deepField"])
	}
}

// TestApplyFieldMapping_Empty 验证空映射返回原数据。
func TestApplyFieldMapping_Empty(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	result := ApplyFieldMapping(data, "", nil)
	if result["key"] != "value" {
		t.Errorf("空映射应返回原数据")
	}
}

// TestApplyValueTransform_Basic 验证值转换。
func TestApplyValueTransform_Basic(t *testing.T) {
	data := map[string]interface{}{
		"level":    "1",
		"severity": "high",
	}
	transform := `{
		"level": {"1": "低危", "2": "中危", "3": "高危"},
		"severity": {"high": "严重"}
	}`

	result := ApplyValueTransform(data, transform)

	if result["level"] != "低危" {
		t.Errorf("level = %v, want 低危", result["level"])
	}
	if result["levelRaw"] != "1" {
		t.Errorf("levelRaw = %v, want 1", result["levelRaw"])
	}
	if result["severity"] != "严重" {
		t.Errorf("severity = %v, want 严重", result["severity"])
	}
}

// TestApplyValueTransform_AlertTimeTimestamp 验证 alertTime 时间戳转换。
func TestApplyValueTransform_AlertTimeTimestamp(t *testing.T) {
	// 秒级时间戳
	data := map[string]interface{}{
		"alertTime": "1700000000",
	}
	result := ApplyValueTransform(data, `{}`)
	expected := time.Unix(1700000000, 0).Format("2006-01-02 15:04:05")
	if result["alertTime"] != expected {
		t.Errorf("alertTime(秒) = %v, want %v", result["alertTime"], expected)
	}

	// 毫秒级时间戳
	data2 := map[string]interface{}{
		"alertTime": "1700000000000",
	}
	result2 := ApplyValueTransform(data2, `{}`)
	if result2["alertTime"] != expected {
		t.Errorf("alertTime(毫秒) = %v, want %v", result2["alertTime"], expected)
	}
}

// TestApplyValueTransform_Empty 验证空转换返回原数据。
func TestApplyValueTransform_Empty(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	result := ApplyValueTransform(data, "")
	if result["key"] != "value" {
		t.Errorf("空转换应返回原数据")
	}
}

// ==================== 提取函数测试 ====================

// TestExtractKeyFields 验证关键字段提取。
func TestExtractKeyFields(t *testing.T) {
	data := map[string]interface{}{
		"attackIp":    "1.1.1.1",
		"victimIp":    "2.2.2.2",
		"description": "测试告警",
		"extraField":  "应被忽略",
	}

	result := ExtractKeyFields(data)
	if result == "" {
		t.Fatal("ExtractKeyFields 返回空字符串")
	}

	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(result), &fields); err != nil {
		t.Fatalf("结果不是有效 JSON: %v", err)
	}

	if fields["attackIp"] != "1.1.1.1" {
		t.Errorf("attackIp = %v, want 1.1.1.1", fields["attackIp"])
	}
	if _, exists := fields["extraField"]; exists {
		t.Error("extraField 不应出现在关键字段中")
	}
}

// TestExtractKeyFields_Empty 验证无关键字段时返回空字符串。
func TestExtractKeyFields_Empty(t *testing.T) {
	data := map[string]interface{}{
		"unrelated": "value",
	}
	result := ExtractKeyFields(data)
	if result != "" {
		t.Errorf("无关键字段应返回空字符串, got %q", result)
	}
}

// TestFormatAlertTime_LocalTimestamp 验证从 localTimestamp 格式化告警时间。
func TestFormatAlertTime_LocalTimestamp(t *testing.T) {
	data := map[string]interface{}{
		"localTimestamp": float64(1700000000000), // 毫秒
	}
	result := FormatAlertTime(data)
	expected := time.UnixMilli(1700000000000).Format("2006-01-02 15:04:05")
	if result != expected {
		t.Errorf("FormatAlertTime(localTimestamp) = %q, want %q", result, expected)
	}
}

// TestFormatAlertTime_TimestampString 验证从 timestamp 字符串格式化。
func TestFormatAlertTime_TimestampString(t *testing.T) {
	data := map[string]interface{}{
		"timestamp": "2024-01-01 12:00:00",
	}
	result := FormatAlertTime(data)
	if result != "2024-01-01 12:00:00" {
		t.Errorf("FormatAlertTime(timestamp) = %q, want 2024-01-01 12:00:00", result)
	}
}

// TestFormatAlertTime_Fallback 验证无时间字段时返回当前时间。
func TestFormatAlertTime_Fallback(t *testing.T) {
	data := map[string]interface{}{}
	result := FormatAlertTime(data)
	if result == "" {
		t.Error("无时间字段应返回当前时间，不应为空")
	}
}

// ==================== fixMalformedJSON 测试 ====================

// TestFixMalformedJSON 验证修复包含非法嵌套 JSON 字符串的字段。
func TestFixMalformedJSON(t *testing.T) {
	// 包含 fullTree 字段的非法 JSON
	input := `{"name":"test","fullTree":"[invalid json]","other":"value"}`
	result := fixMalformedJSON(input)

	// 修复后应能解析
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Errorf("修复后的 JSON 应可解析: %v, result=%s", err, result)
	}
}

// TestFixMalformedJSON_NoFullTree 验证无 fullTree 字段时原样返回。
func TestFixMalformedJSON_NoFullTree(t *testing.T) {
	input := `{"name":"test","value":123}`
	result := fixMalformedJSON(input)
	if result != input {
		t.Errorf("无 fullTree 应原样返回, got %q", result)
	}
}
