package services

import (
	"testing"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupFilterTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.Device{}, &models.DeviceGroup{}, &models.FilterPolicy{}); err != nil {
		t.Fatalf("migrate filter models: %v", err)
	}

	database.DB = db
	return db
}

func TestFilterServiceMatchOperators(t *testing.T) {
	service := NewFilterService()
	data := map[string]interface{}{
		"source_ip": "10.0.0.5",
		"message":   "Blocked By Policy",
		"severity":  "high",
		"count":     "42",
	}

	tests := []struct {
		name string
		cond filterCondition
		want bool
	}{
		{name: "exists", cond: filterCondition{Field: "source_ip", Operator: "exists"}, want: true},
		{name: "not exists", cond: filterCondition{Field: "destination_ip", Operator: "not_exists"}, want: true},
		{name: "equals alias", cond: filterCondition{Field: "severity", Operator: "eq", Value: "HIGH"}, want: true},
		{name: "not equals", cond: filterCondition{Field: "severity", Operator: "!=", Value: "low"}, want: true},
		{name: "contains", cond: filterCondition{Field: "message", Operator: "contains", Value: "blocked"}, want: true},
		{name: "not contains", cond: filterCondition{Field: "message", Operator: "not_contains", Value: "allowed"}, want: true},
		{name: "starts with", cond: filterCondition{Field: "message", Operator: "prefix", Value: "blocked"}, want: true},
		{name: "ends with", cond: filterCondition{Field: "message", Operator: "suffix", Value: "policy"}, want: true},
		{name: "in", cond: filterCondition{Field: "severity", Operator: "in", Value: []interface{}{"low", "high"}}, want: true},
		{name: "not in", cond: filterCondition{Field: "severity", Operator: "not_in", Value: "critical,medium"}, want: true},
		{name: "regex raw case", cond: filterCondition{Field: "message", Operator: "regex", Value: "Blocked\\s+By\\s+Policy"}, want: true},
		{name: "gt", cond: filterCondition{Field: "count", Operator: "gt", Value: 10}, want: true},
		{name: "gte", cond: filterCondition{Field: "count", Operator: ">=", Value: 42}, want: true},
		{name: "lt", cond: filterCondition{Field: "count", Operator: "<", Value: 100}, want: true},
		{name: "lte", cond: filterCondition{Field: "count", Operator: "lte", Value: 42}, want: true},
		{name: "invalid numeric compare", cond: filterCondition{Field: "message", Operator: ">", Value: 5}, want: false},
	}

	for _, tt := range tests {
		if got := service.matchCondition(tt.cond, data); got != tt.want {
			t.Fatalf("%s: expected %v, got %v", tt.name, tt.want, got)
		}
	}
}

func TestFilterServiceLogicAndActions(t *testing.T) {
	service := NewFilterService()
	policy := &models.FilterPolicy{
		Name:           "and-policy",
		Action:         "drop",
		ConditionLogic: "AND",
		Conditions:     `[{"field":"severity","operator":"eq","value":"high"},{"field":"count","operator":">","value":10}]`,
	}

	result := service.applyFilter(policy, map[string]interface{}{
		"severity": "high",
		"count":    42,
	})
	if !result.Matched {
		t.Fatal("expected AND policy to match")
	}
	if result.Action != "drop" {
		t.Fatalf("expected drop action, got %s", result.Action)
	}

	orPolicy := &models.FilterPolicy{
		Name:           "or-policy",
		Action:         "keep",
		ConditionLogic: "OR",
		Conditions:     `[{"field":"severity","operator":"eq","value":"medium"},{"field":"count","operator":"gte","value":42}]`,
	}
	orResult := service.applyFilter(orPolicy, map[string]interface{}{
		"severity": "high",
		"count":    42,
	})
	if !orResult.Matched {
		t.Fatal("expected OR policy to match")
	}
	if orResult.Action != "keep" {
		t.Fatalf("expected keep action, got %s", orResult.Action)
	}
}

func TestFilterServiceWhitelistResult(t *testing.T) {
	service := NewFilterService()
	policy := &models.FilterPolicy{
		Name:             "whitelist-policy",
		Action:           "drop",
		Conditions:       `[{"field":"message","operator":"contains","value":"blocked"}]`,
		WhitelistEnabled: true,
		WhitelistField:   "source_ip",
		WhitelistValues:  "10.0.0.5,10.0.0.6",
	}

	result := service.applyFilter(policy, map[string]interface{}{
		"source_ip": "10.0.0.5",
		"message":   "blocked by policy",
	})
	if !result.Matched {
		t.Fatal("expected policy to match through whitelist")
	}
	if result.Action != "keep" {
		t.Fatalf("expected keep action, got %s", result.Action)
	}
	if result.WhitelistResult != "matched" {
		t.Fatalf("expected whitelist matched, got %s", result.WhitelistResult)
	}

	notMatched := service.applyFilter(policy, map[string]interface{}{
		"source_ip": "10.0.0.9",
		"message":   "blocked by policy",
	})
	if notMatched.WhitelistResult != "not_matched" {
		t.Fatalf("expected whitelist not_matched, got %s", notMatched.WhitelistResult)
	}
}

func TestFilterServicePriorityAndScope(t *testing.T) {
	db := setupFilterTestDB(t)
	service := NewFilterService()

	group := models.DeviceGroup{Name: "core"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("create group: %v", err)
	}

	parseTemplateID := uint(99)
	device := models.Device{
		Name:            "fw-01",
		IPAddress:       "10.0.0.5",
		GroupID:         &group.ID,
		ParseTemplateID: &parseTemplateID,
	}
	if err := db.Create(&device).Error; err != nil {
		t.Fatalf("create device: %v", err)
	}

	lowPriority := models.FilterPolicy{
		Name:            "global-low",
		Action:          "keep",
		Priority:        10,
		Conditions:      `[{"field":"severity","operator":"eq","value":"high"}]`,
		ParseTemplateID: &parseTemplateID,
		Enabled:         true,
	}
	highPriority := models.FilterPolicy{
		Name:       "group-high",
		Action:     "drop",
		Priority:   100,
		Conditions: `[{"field":"severity","operator":"eq","value":"high"}]`,
		Enabled:    true,
	}
	highPriority.DeviceGroupID = &group.ID

	if err := db.Create(&lowPriority).Error; err != nil {
		t.Fatalf("create low priority policy: %v", err)
	}
	if err := db.Create(&highPriority).Error; err != nil {
		t.Fatalf("create high priority policy: %v", err)
	}

	result, err := service.FilterByDevice(device.ID, map[string]interface{}{
		"severity": "high",
	})
	if err != nil {
		t.Fatalf("filter by device: %v", err)
	}
	if !result.Matched {
		t.Fatal("expected scoped policy to match")
	}
	if result.Policy == nil || result.Policy.Name != "group-high" {
		t.Fatalf("expected highest priority policy to win, got %#v", result.Policy)
	}
	if result.Action != "drop" {
		t.Fatalf("expected drop action, got %s", result.Action)
	}
}

func TestFilterServiceInvalidConditions(t *testing.T) {
	service := NewFilterService()
	policy := &models.FilterPolicy{
		Name:       "bad-policy",
		Action:     "keep",
		Conditions: `not-json`,
	}

	result := service.applyFilter(policy, map[string]interface{}{"message": "anything"})
	if result.Matched {
		t.Fatal("expected invalid condition payload not to match")
	}
	if result.Message == "" {
		t.Fatal("expected invalid condition error message")
	}
}

func TestFilterServiceExistsAndRegex(t *testing.T) {
	service := NewFilterService()
	data := map[string]interface{}{
		"source_ip": "10.0.0.5",
		"message":   "blocked by policy",
	}

	if !service.matchCondition(filterCondition{Field: "source_ip", Operator: "exists"}, data) {
		t.Fatal("expected exists to match")
	}
	if !service.matchCondition(filterCondition{Field: "destination_ip", Operator: "not_exists"}, data) {
		t.Fatal("expected not_exists to match")
	}
	if !service.matchCondition(filterCondition{Field: "message", Operator: "regex", Value: "blocked\\s+by\\s+policy"}, data) {
		t.Fatal("expected regex to match")
	}
}
