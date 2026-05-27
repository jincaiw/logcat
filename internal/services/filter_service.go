package services

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Matched bool              `json:"matched"`
	Action  string            `json:"action"`
	Policy  *models.FilterPolicy `json:"policy,omitempty"`
	Message string            `json:"message,omitempty"`
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

	// Policies can target specific device, device group, or parse template
	query = query.Where("device_id = ? OR device_group_id = ? OR (device_id IS NULL AND device_group_id IS NULL)",
		deviceID, device.GroupID)

	query.Order("priority DESC").Find(&policies)

	for _, policy := range policies {
		result := s.applyFilter(&policy, parsedData)
		if result.Matched {
			return result, nil
		}
	}

	return &FilterResult{Matched: false, Action: "keep", Message: "no matching policy, default keep"}, nil
}

type filterCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// applyFilter applies a single policy's conditions against parsed data
func (s *FilterService) applyFilter(policy *models.FilterPolicy, parsedData map[string]interface{}) *FilterResult {
	if policy.Conditions == "" {
		return &FilterResult{Matched: false, Action: policy.Action}
	}

	var conditions []filterCondition
	if err := json.Unmarshal([]byte(policy.Conditions), &conditions); err != nil {
		return &FilterResult{
			Matched: false,
			Action:  policy.Action,
			Message: fmt.Sprintf("failed to parse conditions: %v", err),
		}
	}

	if len(conditions) == 0 {
		return &FilterResult{Matched: false, Action: policy.Action}
	}

	// Check whitelist first if enabled
	if policy.WhitelistEnabled && policy.WhitelistField != "" {
		if s.checkWhitelist(policy, parsedData) {
			return &FilterResult{
				Matched:  true,
				Action:   "keep",
				Policy:   policy,
				Message:  "whitelist matched",
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
			Matched: true,
			Action:  policy.Action,
			Policy:  policy,
			Message: fmt.Sprintf("policy '%s' matched", policy.Name),
		}
	}

	return &FilterResult{Matched: false, Action: policy.Action}
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
	if !exists {
		return false
	}

	fieldStr := strings.ToLower(fmt.Sprintf("%v", fieldValue))
	valueStr := strings.ToLower(fmt.Sprintf("%v", cond.Value))

	switch cond.Operator {
	case "equals", "==":
		return fieldStr == valueStr
	case "not_equals", "!=":
		return fieldStr != valueStr
	case "contains":
		return strings.Contains(fieldStr, valueStr)
	case "not_contains":
		return !strings.Contains(fieldStr, valueStr)
	case "starts_with":
		return strings.HasPrefix(fieldStr, valueStr)
	case "ends_with":
		return strings.HasSuffix(fieldStr, valueStr)
	case "in":
		return s.checkIn(fieldStr, cond.Value)
	case "not_in":
		return !s.checkIn(fieldStr, cond.Value)
	case "regex_match":
		return s.checkRegex(fieldStr, valueStr)
	case "greater_than", ">":
		return s.compareValues(fieldValue, cond.Value) > 0
	case "less_than", "<":
		return s.compareValues(fieldValue, cond.Value) < 0
	case "greater_equal", ">=":
		return s.compareValues(fieldValue, cond.Value) >= 0
	case "less_equal", "<=":
		return s.compareValues(fieldValue, cond.Value) <= 0
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
	// Simple contains match for now; full regex support in pipeline
	return strings.Contains(value, strings.Trim(pattern, "/"))
}

func (s *FilterService) compareValues(a, b interface{}) int {
	aFloat := toFloat64(a)
	bFloat := toFloat64(b)
	if aFloat < bFloat {
		return -1
	} else if aFloat > bFloat {
		return 1
	}
	return 0
}

func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case json.Number:
		f, _ := val.Float64()
		return f
	default:
		// Try to parse string
		var f float64
		fmt.Sscanf(fmt.Sprintf("%v", val), "%f", &f)
		return f
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