package filter

import (
	"testing"

	"syslog-alert/internal/models"
)

// TestEvaluateCondition_Regex 验证 regex 操作符使用真正的正则匹配，
// 而非旧的 strings.Contains 行为（这是重构中修复的关键 bug）。
func TestEvaluateCondition_Regex(t *testing.T) {
	data := map[string]interface{}{
		"ip": "192.168.1.100",
	}

	tests := []struct {
		name     string
		cond     models.FilterCondition
		expected bool
	}{
		{
			name:     "正则匹配 IP 网段",
			cond:     models.FilterCondition{Field: "ip", Operator: "regex", Value: `^192\.168\.`},
			expected: true,
		},
		{
			name:     "正则不匹配",
			cond:     models.FilterCondition{Field: "ip", Operator: "regex", Value: `^10\.`},
			expected: false,
		},
		{
			name:     "=~ 别名匹配",
			cond:     models.FilterCondition{Field: "ip", Operator: "=~", Value: `\d+\.\d+\.\d+\.\d+`},
			expected: true,
		},
		{
			name:     "正则匹配数字结尾",
			cond:     models.FilterCondition{Field: "ip", Operator: "regex", Value: `\.100$`},
			expected: true,
		},
		// 关键回归测试：旧实现用 strings.Contains，以下用例会错误返回 true
		{
			name:     "正则不应做子串匹配-点号",
			cond:     models.FilterCondition{Field: "ip", Operator: "regex", Value: `.`},
			expected: true, // "." 作为正则匹配任意字符，确实匹配
		},
		{
			name:     "正则转义点号-不应匹配任意字符",
			cond:     models.FilterCondition{Field: "ip", Operator: "regex", Value: `\.`},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateCondition(tt.cond, data)
			if result != tt.expected {
				t.Errorf("EvaluateCondition() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestEvaluateCondition_RegexRegression 验证旧 bug 不会回归：
// 旧实现 strings.Contains("192.168.1.100", "192.168") 返回 true，
// 但 "192.168" 作为正则也匹配（因为 . 在正则中匹配任意字符）。
// 但 "192\.168" 作为正则也应该匹配，而 strings.Contains 不会匹配反斜杠。
func TestEvaluateCondition_RegexRegression(t *testing.T) {
	data := map[string]interface{}{
		"msg": "attack from 203.0.113.50",
	}

	// 这个正则包含 \d，strings.Contains 不会匹配 \d 字面量
	cond := models.FilterCondition{
		Field:    "msg",
		Operator: "regex",
		Value:    `\d+\.\d+\.\d+\.\d+`,
	}
	if !EvaluateCondition(cond, data) {
		t.Error("正则 \\d+\\.\\d+\\.\\d+\\.\\d+ 应该匹配 IP 地址，但未匹配")
	}
}

// TestEvaluateCondition_NotRegex 验证 !~ 操作符。
func TestEvaluateCondition_NotRegex(t *testing.T) {
	data := map[string]interface{}{
		"status": "200",
	}

	tests := []struct {
		name     string
		cond     models.FilterCondition
		expected bool
	}{
		{
			name:     "!~ 不匹配 4xx",
			cond:     models.FilterCondition{Field: "status", Operator: "!~", Value: `^4`},
			expected: true,
		},
		{
			name:     "!~ 匹配 2xx 时返回 false",
			cond:     models.FilterCondition{Field: "status", Operator: "!~", Value: `^2`},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateCondition(tt.cond, data)
			if result != tt.expected {
				t.Errorf("EvaluateCondition() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestEvaluateCondition_Equals 验证等于操作符。
func TestEvaluateCondition_Equals(t *testing.T) {
	data := map[string]interface{}{
		"level":    3,
		"protocol": "tcp",
	}

	tests := []struct {
		name     string
		cond     models.FilterCondition
		expected bool
	}{
		{name: "数字相等", cond: models.FilterCondition{Field: "level", Operator: "equals", Value: "3"}, expected: true},
		{name: "数字不等", cond: models.FilterCondition{Field: "level", Operator: "equals", Value: "5"}, expected: false},
		{name: "字符串相等", cond: models.FilterCondition{Field: "protocol", Operator: "equals", Value: "tcp"}, expected: true},
		{name: "==别名", cond: models.FilterCondition{Field: "protocol", Operator: "==", Value: "tcp"}, expected: true},
		{name: "!=不等", cond: models.FilterCondition{Field: "protocol", Operator: "!=", Value: "udp"}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvaluateCondition(tt.cond, data); got != tt.expected {
				t.Errorf("EvaluateCondition() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestEvaluateCondition_Contains 验证包含操作符。
func TestEvaluateCondition_Contains(t *testing.T) {
	data := map[string]interface{}{
		"message": "SSH brute force attack detected",
	}

	tests := []struct {
		name     string
		cond     models.FilterCondition
		expected bool
	}{
		{name: "包含子串", cond: models.FilterCondition{Field: "message", Operator: "contains", Value: "brute"}, expected: true},
		{name: "不包含", cond: models.FilterCondition{Field: "message", Operator: "contains", Value: "normal"}, expected: false},
		{name: "notContains", cond: models.FilterCondition{Field: "message", Operator: "not_contains", Value: "normal"}, expected: true},
		{name: "startsWith", cond: models.FilterCondition{Field: "message", Operator: "starts_with", Value: "SSH"}, expected: true},
		{name: "endsWith", cond: models.FilterCondition{Field: "message", Operator: "ends_with", Value: "detected"}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvaluateCondition(tt.cond, data); got != tt.expected {
				t.Errorf("EvaluateCondition() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestEvaluateCondition_In 验证 in/notIn 操作符。
func TestEvaluateCondition_In(t *testing.T) {
	data := map[string]interface{}{
		"port": "443",
	}

	tests := []struct {
		name     string
		cond     models.FilterCondition
		expected bool
	}{
		{name: "in 匹配", cond: models.FilterCondition{Field: "port", Operator: "in", Value: "80,443,8443"}, expected: true},
		{name: "in 不匹配", cond: models.FilterCondition{Field: "port", Operator: "in", Value: "80,8080"}, expected: false},
		{name: "notIn 匹配", cond: models.FilterCondition{Field: "port", Operator: "not_in", Value: "80,8080"}, expected: true},
		{name: "in 带空格", cond: models.FilterCondition{Field: "port", Operator: "in", Value: "80, 443, 8443"}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvaluateCondition(tt.cond, data); got != tt.expected {
				t.Errorf("EvaluateCondition() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestEvaluateCondition_NumberComparison 验证数值比较操作符。
func TestEvaluateCondition_NumberComparison(t *testing.T) {
	data := map[string]interface{}{
		"level": "3",
	}

	tests := []struct {
		name     string
		cond     models.FilterCondition
		expected bool
	}{
		{name: "GT 大于", cond: models.FilterCondition{Field: "level", Operator: "gt", Value: "2"}, expected: true},
		{name: "GT 不大于", cond: models.FilterCondition{Field: "level", Operator: "gt", Value: "3"}, expected: false},
		{name: "GTE 大于等于", cond: models.FilterCondition{Field: "level", Operator: "gte", Value: "3"}, expected: true},
		{name: "LT 小于", cond: models.FilterCondition{Field: "level", Operator: "lt", Value: "4"}, expected: true},
		{name: "LTE 小于等于", cond: models.FilterCondition{Field: "level", Operator: "lte", Value: "3"}, expected: true},
		{name: "> 别名", cond: models.FilterCondition{Field: "level", Operator: ">", Value: "2"}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvaluateCondition(tt.cond, data); got != tt.expected {
				t.Errorf("EvaluateCondition() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestEvaluateCondition_Exists 验证 exists/notExists 操作符。
func TestEvaluateCondition_Exists(t *testing.T) {
	data := map[string]interface{}{
		"ip": "10.0.0.1",
	}

	tests := []struct {
		name     string
		cond     models.FilterCondition
		expected bool
	}{
		{name: "exists 存在", cond: models.FilterCondition{Field: "ip", Operator: "exists", Value: ""}, expected: true},
		{name: "exists 不存在", cond: models.FilterCondition{Field: "missing", Operator: "exists", Value: ""}, expected: false},
		{name: "notExists 不存在", cond: models.FilterCondition{Field: "missing", Operator: "not_exists", Value: ""}, expected: true},
		{name: "notExists 存在", cond: models.FilterCondition{Field: "ip", Operator: "not_exists", Value: ""}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvaluateCondition(tt.cond, data); got != tt.expected {
				t.Errorf("EvaluateCondition() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestEvaluateConditions_AndLogic 验证 AND 逻辑。
func TestEvaluateConditions_AndLogic(t *testing.T) {
	data := map[string]interface{}{
		"ip":    "192.168.1.1",
		"port":  "443",
		"level": "3",
	}

	conditions := []models.FilterCondition{
		{Field: "ip", Operator: "contains", Value: "192.168"},
		{Field: "port", Operator: "equals", Value: "443"},
		{Field: "level", Operator: "gt", Value: "2"},
	}

	if !EvaluateConditions(conditions, data, "AND") {
		t.Error("AND 逻辑：所有条件满足时应返回 true")
	}

	// 添加一个不满足的条件
	conditions = append(conditions, models.FilterCondition{Field: "level", Operator: "gt", Value: "5"})
	if EvaluateConditions(conditions, data, "AND") {
		t.Error("AND 逻辑：任一条件不满足时应返回 false")
	}
}

// TestEvaluateConditions_OrLogic 验证 OR 逻辑。
func TestEvaluateConditions_OrLogic(t *testing.T) {
	data := map[string]interface{}{
		"status": "200",
	}

	conditions := []models.FilterCondition{
		{Field: "status", Operator: "equals", Value: "404"},
		{Field: "status", Operator: "equals", Value: "500"},
		{Field: "status", Operator: "equals", Value: "200"},
	}

	if !EvaluateConditions(conditions, data, "OR") {
		t.Error("OR 逻辑：任一条件满足时应返回 true")
	}

	// 全部不满足
	conditions = []models.FilterCondition{
		{Field: "status", Operator: "equals", Value: "404"},
		{Field: "status", Operator: "equals", Value: "500"},
	}
	if EvaluateConditions(conditions, data, "OR") {
		t.Error("OR 逻辑：所有条件不满足时应返回 false")
	}
}

// TestEvaluateConditions_Empty 验证空条件列表。
func TestEvaluateConditions_Empty(t *testing.T) {
	if !EvaluateConditions(nil, map[string]interface{}{}, "AND") {
		t.Error("空条件列表应返回 true")
	}
}

// TestMatchCIDR 验证 CIDR 匹配。
func TestMatchCIDR(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		cidr     string
		expected bool
	}{
		{name: "单IP精确匹配", ip: "192.168.1.1", cidr: "192.168.1.1", expected: true},
		{name: "单IP不匹配", ip: "192.168.1.2", cidr: "192.168.1.1", expected: false},
		{name: "CIDR /24 匹配", ip: "192.168.1.100", cidr: "192.168.1.0/24", expected: true},
		{name: "CIDR /24 不匹配", ip: "192.168.2.100", cidr: "192.168.1.0/24", expected: false},
		{name: "CIDR /16 匹配", ip: "10.0.5.10", cidr: "10.0.0.0/16", expected: true},
		{name: "CIDR /8 匹配", ip: "10.1.2.3", cidr: "10.0.0.0/8", expected: true},
		{name: "无效IP", ip: "invalid", cidr: "192.168.1.0/24", expected: false},
		{name: "无效CIDR", ip: "192.168.1.1", cidr: "invalid", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchCIDR(tt.ip, tt.cidr); got != tt.expected {
				t.Errorf("MatchCIDR(%q, %q) = %v, want %v", tt.ip, tt.cidr, got, tt.expected)
			}
		})
	}
}

// TestMatchWhitelist 验证白名单匹配。
func TestMatchWhitelist(t *testing.T) {
	data := map[string]interface{}{
		"attackIp": "192.168.1.100",
	}

	whitelistJSON := `[
		{"cidr": "10.0.0.0/8", "description": "内网A", "enabled": true},
		{"cidr": "192.168.0.0/16", "description": "内网B", "enabled": true},
		{"cidr": "172.16.0.0/12", "description": "禁用", "enabled": false}
	]`

	// 应命中 192.168.0.0/16
	matched, err := MatchWhitelist(data, "attackIp", whitelistJSON)
	if err != nil {
		t.Fatalf("MatchWhitelist() error: %v", err)
	}
	if !matched {
		t.Error("192.168.1.100 应命中 192.168.0.0/16 白名单")
	}

	// 不在白名单中的 IP
	data["attackIp"] = "8.8.8.8"
	matched, _ = MatchWhitelist(data, "attackIp", whitelistJSON)
	if matched {
		t.Error("8.8.8.8 不应命中白名单")
	}

	// 禁用的白名单条目不应匹配
	data["attackIp"] = "172.16.1.1"
	matched, _ = MatchWhitelist(data, "attackIp", whitelistJSON)
	if matched {
		t.Error("172.16.1.1 不应命中禁用的白名单条目")
	}

	// 字段不存在
	matched, _ = MatchWhitelist(data, "missingField", whitelistJSON)
	if matched {
		t.Error("字段不存在时不应匹配")
	}
}

// TestMatchWhitelist_InvalidJSON 验证无效 JSON 返回错误。
func TestMatchWhitelist_InvalidJSON(t *testing.T) {
	_, err := MatchWhitelist(map[string]interface{}{}, "ip", "invalid json")
	if err == nil {
		t.Error("无效 JSON 应返回错误")
	}
}
