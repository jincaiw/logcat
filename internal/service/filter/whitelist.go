package filter

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"syslog-alert/internal/models"
	applogger "syslog-alert/pkg/logger"
)

// MatchWhitelist 检查数据中指定字段的 IP 是否命中白名单。
//
// 参数：
//   - data: 解析后的日志数据
//   - field: 白名单检查的字段名（如 "attackIp"）
//   - whitelistJSON: 白名单 JSON 配置
//
// 返回 true 表示命中白名单（应跳过告警）。
func MatchWhitelist(data map[string]interface{}, field, whitelistJSON string) (bool, error) {
	var whitelist []models.WhitelistItem
	if err := json.Unmarshal([]byte(whitelistJSON), &whitelist); err != nil {
		return false, fmt.Errorf("invalid whitelist: %v", err)
	}

	value, exists := data[field]
	if !exists {
		applogger.Debug("Whitelist field %s not found in parsed data", field)
		return false, nil
	}

	ipStr := fmt.Sprintf("%v", value)
	applogger.Debug("Checking whitelist: field=%s, ip=%s, items=%d", field, ipStr, len(whitelist))

	for _, item := range whitelist {
		if !item.Enabled {
			continue
		}
		if MatchCIDR(ipStr, item.CIDR) {
			applogger.Debug("Whitelist matched: IP %s matches CIDR %s", ipStr, item.CIDR)
			return true, nil
		}
	}

	applogger.Debug("No whitelist match for IP %s", ipStr)
	return false, nil
}

// MatchCIDR 检查 IP 是否匹配 CIDR 或单个 IP。
//
// 合并了原 filter.go 和 syslog_service.go 中重复的 matchCIDR 实现。
// 支持：
//   - 单个 IP（无 "/"）：精确匹配
//   - CIDR 表示法（如 192.168.1.0/24）：网段匹配
func MatchCIDR(ipStr, cidr string) bool {
	// 单个 IP：精确匹配
	if !strings.Contains(cidr, "/") {
		return ipStr == cidr
	}

	// CIDR：网段匹配
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	return ipNet.Contains(ip)
}
