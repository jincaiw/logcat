package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// FilterService handles log filtering logic
type FilterService struct{}

// NewFilterService creates a new FilterService
func NewFilterService() *FilterService {
	return &FilterService{}
}

// FilterResult holds the result of filtering
type FilterResult struct {
	Matched         bool                 `json:"matched"`
	Action          string               `json:"action"`
	Policy          *models.FilterPolicy `json:"policy,omitempty"`
	Message         string               `json:"message,omitempty"`
	WhitelistResult string               `json:"whitelistResult,omitempty"`
}

// TestFilter tests a filter policy against sample parsed data
func (s *FilterService) TestFilter(policyID uint, parsedData map[string]interface{}) (*FilterResult, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var policy models.FilterPolicy
	if err := db.First(&policy, policyID).Error; err != nil {
		return nil, err
	}

	return s.applyFilter(&policy, parsedData), nil
}

// FilterByDevice finds and applies the matching filter policy for a device
func (s *FilterService) FilterByDevice(deviceID uint, parsedData map[string]interface{}) (*FilterResult, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	// Get device to find its group
	var device models.Device
	if err := db.First(&device, deviceID).Error; err != nil {
		return nil, err
	}

	// Find applicable policies
	var policies []models.FilterPolicy
	query := db.Where("enabled = ?", true)

	// Policies can target specific device, device group, parse template, or be global.
	query = query.Where(
		"device_id = ? OR device_group_id = ? OR parse_template_id = ? OR (device_id IS NULL AND device_group_id IS NULL AND parse_template_id IS NULL)",
		deviceID, device.GroupID, device.ParseTemplateID,
	)

	query.Order("priority DESC").Find(&policies)

	for _, policy := range policies {
		result := s.applyFilter(&policy, parsedData)
		if result.Matched {
			return result, nil
		}
	}

	return &FilterResult{Matched: false, Action: "keep", Message: "no matching policy, default keep", WhitelistResult: "not_enabled"}, nil
}

type filterCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// applyFilter applies a single policy's conditions against parsed data
func (s *FilterService) applyFilter(policy *models.FilterPolicy, parsedData map[string]interface{}) *FilterResult {
	whitelistResult := "not_enabled"
	if policy.Conditions == "" {
		return &FilterResult{Matched: false, Action: policy.Action, WhitelistResult: whitelistResult}
	}

	var conditions []filterCondition
	if err := json.Unmarshal([]byte(policy.Conditions), &conditions); err != nil {
		return &FilterResult{
			Matched:         false,
			Action:          policy.Action,
			Message:         fmt.Sprintf("failed to parse conditions: %v", err),
			WhitelistResult: whitelistResult,
		}
	}

	if len(conditions) == 0 {
		return &FilterResult{Matched: false, Action: policy.Action, WhitelistResult: whitelistResult}
	}

	// Check whitelist first if enabled
	if policy.WhitelistEnabled && policy.WhitelistField != "" {
		whitelistResult = "not_matched"
		if s.checkWhitelist(policy, parsedData) {
			return &FilterResult{
				Matched:         true,
				Action:          "keep",
				Policy:          policy,
				Message:         "whitelist matched",
				WhitelistResult: "matched",
			}
		}
	}

	// Evaluate conditions
	var matched bool
	if policy.ConditionLogic == "OR" {
		matched = s.evalOR(conditions, parsedData)
	} else {
		matched = s.evalAND(conditions, parsedData)
	}

	if matched {
		return &FilterResult{
			Matched:         true,
			Action:          policy.Action,
			Policy:          policy,
			Message:         fmt.Sprintf("policy '%s' matched", policy.Name),
			WhitelistResult: whitelistResult,
		}
	}

	return &FilterResult{Matched: false, Action: policy.Action, WhitelistResult: whitelistResult}
}

func (s *FilterService) evalAND(conditions []filterCondition, data map[string]interface{}) bool {
	for _, cond := range conditions {
		if !s.matchCondition(cond, data) {
			return false
		}
	}
	return true
}

func (s *FilterService) evalOR(conditions []filterCondition, data map[string]interface{}) bool {
	for _, cond := range conditions {
		if s.matchCondition(cond, data) {
			return true
		}
	}
	return false
}

func (s *FilterService) matchCondition(cond filterCondition, data map[string]interface{}) bool {
	fieldValue, exists := data[cond.Field]
	operator := strings.ToLower(strings.TrimSpace(cond.Operator))
	if operator == "exists" {
		return exists && strings.TrimSpace(fmt.Sprintf("%v", fieldValue)) != ""
	}
	if operator == "not_exists" {
		return !exists || strings.TrimSpace(fmt.Sprintf("%v", fieldValue)) == ""
	}
	if !exists {
		return false
	}

	fieldRaw := strings.TrimSpace(fmt.Sprintf("%v", fieldValue))
	fieldStr := strings.ToLower(fieldRaw)
	valueRaw := strings.TrimSpace(fmt.Sprintf("%v", cond.Value))
	valueStr := strings.ToLower(valueRaw)

	switch operator {
	case "equals", "equal", "eq", "==":
		return fieldStr == valueStr
	case "not_equals", "not_equal", "ne", "!=":
		return fieldStr != valueStr
	case "contains":
		return strings.Contains(fieldStr, valueStr)
	case "not_contains":
		return !strings.Contains(fieldStr, valueStr)
	case "starts_with", "prefix":
		return strings.HasPrefix(fieldStr, valueStr)
	case "ends_with", "suffix":
		return strings.HasSuffix(fieldStr, valueStr)
	case "in":
		return s.checkIn(fieldStr, cond.Value)
	case "not_in":
		return !s.checkIn(fieldStr, cond.Value)
	case "regex", "regex_match":
		return s.checkRegex(fieldRaw, valueRaw)
	case "greater_than", "gt", ">":
		return s.compareValues(fieldValue, cond.Value, ">")
	case "less_than", "lt", "<":
		return s.compareValues(fieldValue, cond.Value, "<")
	case "greater_equal", "greater_than_or_equal", "gte", ">=":
		return s.compareValues(fieldValue, cond.Value, ">=")
	case "less_equal", "less_than_or_equal", "lte", "<=":
		return s.compareValues(fieldValue, cond.Value, "<=")
	default:
		return false
	}
}

func (s *FilterService) checkIn(value string, list interface{}) bool {
	switch v := list.(type) {
	case []interface{}:
		for _, item := range v {
			if strings.ToLower(fmt.Sprintf("%v", item)) == value {
				return true
			}
		}
	case string:
		// Comma-separated list
		items := strings.Split(v, ",")
		for _, item := range items {
			if strings.TrimSpace(strings.ToLower(item)) == value {
				return true
			}
		}
	}
	return false
}

func (s *FilterService) checkRegex(value, pattern string) bool {
	if pattern == "" {
		return false
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(value)
}

func (s *FilterService) compareValues(a, b interface{}, operator string) bool {
	aFloat, okA := toFloat64(a)
	bFloat, okB := toFloat64(b)
	if !okA || !okB {
		return false
	}

	switch operator {
	case ">":
		return aFloat > bFloat
	case "<":
		return aFloat < bFloat
	case ">=":
		return aFloat >= bFloat
	case "<=":
		return aFloat <= bFloat
	default:
		return false
	}
}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case json.Number:
		f, _ := val.Float64()
		return f, true
	default:
		if f, err := strconv.ParseFloat(strings.TrimSpace(fmt.Sprintf("%v", val)), 64); err == nil {
			return f, true
		}
		return 0, false
	}
}

func (s *FilterService) checkWhitelist(policy *models.FilterPolicy, data map[string]interface{}) bool {
	fieldValue, exists := data[policy.WhitelistField]
	if !exists {
		return false
	}

	valueStr := fmt.Sprintf("%v", fieldValue)
	if policy.WhitelistValues != "" {
		return s.checkIn(strings.ToLower(valueStr), policy.WhitelistValues)
	}

	return false
}
