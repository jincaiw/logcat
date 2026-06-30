package filter

import (
	"fmt"
	"regexp"
	"strings"

	"syslog-alert/internal/models"
	"syslog-alert/pkg/constants"
	applogger "syslog-alert/pkg/logger"
)

// EvaluateConditions 按 AND/OR 逻辑评估条件列表。
// 空条件列表返回 true（视为匹配）。
func EvaluateConditions(conditions []models.FilterCondition, data map[string]interface{}, logic string) bool {
	if len(conditions) == 0 {
		return true
	}

	results := make([]bool, len(conditions))
	for i, cond := range conditions {
		results[i] = EvaluateCondition(cond, data)
	}

	if logic == constants.LogicOR {
		for _, r := range results {
			if r {
				return true
			}
		}
		return false
	}

	// 默认 AND
	for _, r := range results {
		if !r {
			return false
		}
	}
	return true
}

// EvaluateCondition 评估单个筛选条件。
//
// 修复说明：原 syslog_service.go 中的 regexpMatch 函数错误地使用 strings.Contains
// 代替正则匹配，导致 "regex" 和 "=~" 操作符实际只做子串匹配。
// 此处改用 regexp.MatchString 实现真正的正则匹配。
func EvaluateCondition(cond models.FilterCondition, data map[string]interface{}) bool {
	value, exists := data[cond.Field]
	if !exists {
		return cond.Operator == constants.OpNotExists
	}

	strValue := fmt.Sprintf("%v", value)

	switch cond.Operator {
	case constants.OpEquals, "==":
		return strValue == cond.Value
	case constants.OpNotEquals, "!=":
		return strValue != cond.Value
	case constants.OpContains:
		return strings.Contains(strValue, cond.Value)
	case constants.OpNotContains:
		return !strings.Contains(strValue, cond.Value)
	case constants.OpIn:
		return matchIn(strValue, cond.Value)
	case constants.OpNotIn:
		return !matchIn(strValue, cond.Value)
	case constants.OpStartsWith:
		return strings.HasPrefix(strValue, cond.Value)
	case constants.OpEndsWith:
		return strings.HasSuffix(strValue, cond.Value)
	case constants.OpRegex, "=~":
		// 修复：使用真正的正则匹配，替代原 strings.Contains
		matched, err := regexp.MatchString(cond.Value, strValue)
		if err != nil {
			applogger.Warn("Invalid regex pattern %q: %v", cond.Value, err)
			return false
		}
		return matched
	case constants.OpNotRegex, "!~":
		// 修复：使用真正的正则匹配
		matched, err := regexp.MatchString(cond.Value, strValue)
		if err != nil {
			applogger.Warn("Invalid regex pattern %q: %v", cond.Value, err)
			return false
		}
		return !matched
	case constants.OpExists:
		return exists
	case constants.OpNotExists:
		return !exists
	case constants.OpGT, ">":
		return compareNumbers(strValue, cond.Value) > 0
	case constants.OpGTE, ">=":
		return compareNumbers(strValue, cond.Value) >= 0
	case constants.OpLT, "<":
		return compareNumbers(strValue, cond.Value) < 0
	case constants.OpLTE, "<=":
		return compareNumbers(strValue, cond.Value) <= 0
	default:
		return false
	}
}

// matchIn 检查 strValue 是否在逗号分隔的值列表中。
func matchIn(strValue, csvValues string) bool {
	for _, v := range strings.Split(csvValues, ",") {
		if strings.TrimSpace(v) == strValue {
			return true
		}
	}
	return false
}

// compareNumbers 比较两个数值字符串。非数字时回退到字符串比较。
func compareNumbers(a, b string) int {
	var aNum, bNum float64
	_, err1 := fmt.Sscanf(a, "%f", &aNum)
	_, err2 := fmt.Sscanf(b, "%f", &bNum)

	if err1 != nil || err2 != nil {
		return strings.Compare(a, b)
	}

	switch {
	case aNum > bNum:
		return 1
	case aNum < bNum:
		return -1
	default:
		return 0
	}
}
