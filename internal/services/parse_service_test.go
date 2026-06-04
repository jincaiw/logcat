package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupParseTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.ParseTemplate{}); err != nil {
		t.Fatalf("migrate parse template: %v", err)
	}

	database.DB = db
	return db
}

func TestParseServiceRegexAndValueTransform(t *testing.T) {
	service := NewParseService()
	template := &models.ParseTemplate{
		ParseType:    "regex",
		HeaderRegex:  `src=(?P<src_ip>\S+)\s+severity=(?P<severity>\S+)\s+count=(?P<count>\d+)`,
		FieldMapping: `{"src_ip":"source_ip","severity":"severity","count":"count"}`,
		ValueTransform: `{
			"severity":{"type":"upper"},
			"count":{"type":"int"}
		}`,
	}

	result := service.parse(template, `src=10.0.0.8 severity=high count=42`, 0)
	if !result.Success {
		t.Fatalf("expected success, got error: %s", result.Error)
	}
	if result.Fields["source_ip"] != "10.0.0.8" {
		t.Fatalf("unexpected source_ip: %#v", result.Fields["source_ip"])
	}
	if result.Fields["severity"] != "HIGH" {
		t.Fatalf("unexpected severity: %#v", result.Fields["severity"])
	}
	if result.Fields["count"] != 42 {
		t.Fatalf("unexpected count: %#v", result.Fields["count"])
	}
}

func TestParseServiceRegexOrderedFieldMapping(t *testing.T) {
	service := NewParseService()
	template := &models.ParseTemplate{
		ParseType:    "regex",
		HeaderRegex:  `user=(\S+)\s+status=(\S+)`,
		FieldMapping: `["username","status"]`,
	}

	result := service.parse(template, `user=alice status=ok`, 0)
	if !result.Success {
		t.Fatalf("expected success, got error: %s", result.Error)
	}
	if result.Fields["username"] != "alice" {
		t.Fatalf("unexpected username: %#v", result.Fields["username"])
	}
	if result.Fields["status"] != "ok" {
		t.Fatalf("unexpected status: %#v", result.Fields["status"])
	}
}

func TestParseServiceValueTransformVariants(t *testing.T) {
	service := NewParseService()
	template := &models.ParseTemplate{
		ParseType: "json",
		ValueTransform: `[
			{"field":"name","type":"trim"},
			{"field":"severity","type":"lower"},
			{"field":"category","type":"upper"},
			{"field":"count","type":"int"},
			{"field":"ratio","type":"float"},
			{"field":"enabled","type":"bool"},
			{"field":"status","type":"map","mapping":{"1":"online"},"default":"unknown"},
			{"field":"occurredAt","type":"datetime","layout":"2006-01-02 15:04:05"}
		]`,
	}

	result := service.parse(template, `{
		"name":"  Alice  ",
		"severity":"HIGH",
		"category":"warning",
		"count":"42",
		"ratio":"3.14",
		"enabled":"true",
		"status":"1",
		"occurredAt":"2026-05-28 10:20:30"
	}`, 0)
	if !result.Success {
		t.Fatalf("expected success, got error: %s", result.Error)
	}
	if result.Fields["name"] != "Alice" {
		t.Fatalf("unexpected trim result: %#v", result.Fields["name"])
	}
	if result.Fields["severity"] != "high" {
		t.Fatalf("unexpected lower result: %#v", result.Fields["severity"])
	}
	if result.Fields["category"] != "WARNING" {
		t.Fatalf("unexpected upper result: %#v", result.Fields["category"])
	}
	if result.Fields["count"] != 42 {
		t.Fatalf("unexpected int result: %#v", result.Fields["count"])
	}
	if result.Fields["ratio"] != 3.14 {
		t.Fatalf("unexpected float result: %#v", result.Fields["ratio"])
	}
	if result.Fields["enabled"] != true {
		t.Fatalf("unexpected bool result: %#v", result.Fields["enabled"])
	}
	if result.Fields["status"] != "online" {
		t.Fatalf("unexpected map result: %#v", result.Fields["status"])
	}
	expectedTime := time.Date(2026, 5, 28, 10, 20, 30, 0, time.UTC).Format(time.RFC3339)
	if result.Fields["occurredAt"] != expectedTime {
		t.Fatalf("unexpected datetime result: %#v", result.Fields["occurredAt"])
	}
}

func TestParseServiceSyslogJSON(t *testing.T) {
	service := NewParseService()
	template := &models.ParseTemplate{ParseType: "syslog_json"}

	success := service.parse(template, `<34>1 2026-05-28T10:00:00Z host app - - - {"event":"login","severity":"high"}`, 0)
	if !success.Success {
		t.Fatalf("expected success, got error: %s", success.Error)
	}
	if success.Fields["event"] != "login" {
		t.Fatalf("unexpected event: %#v", success.Fields["event"])
	}
	if success.Fields["syslog_header"] == "" {
		t.Fatal("expected syslog header to be present")
	}

	failure := service.parse(template, `<34>1 2026-05-28T10:00:00Z host app - - - not-json`, 0)
	if failure.Success {
		t.Fatal("expected missing body to fail")
	}
	if failure.Error == "" {
		t.Fatal("expected error message for invalid syslog json")
	}
}

func TestParseServiceSubTemplateRouting(t *testing.T) {
	db := setupParseTestDB(t)
	service := NewParseService()

	child := models.ParseTemplate{
		Name:          "child",
		ParseType:     "regex",
		HeaderRegex:   `event=(?P<event>\S+)`,
		FieldMapping:  `{"event":"event_type"}`,
		ValueTransform: `{"event_type":"upper"}`,
		Enabled:       true,
	}
	if err := db.Create(&child).Error; err != nil {
		t.Fatalf("create child template: %v", err)
	}

	routes, _ := json.Marshal([]map[string]interface{}{
		{
			"templateId": child.ID,
			"matchType":  "contains",
			"matchField": "raw_message",
			"matchValue": "event=",
		},
	})
	parent := models.ParseTemplate{
		Name:         "parent",
		ParseType:    "sub_template",
		SubTemplates: string(routes),
		Enabled:      true,
	}
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("create parent template: %v", err)
	}

	result, err := service.TestParse(parent.ID, "foo event=login")
	if err != nil {
		t.Fatalf("test parse: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success, got error: %s", result.Error)
	}
	if result.Fields["event_type"] != "LOGIN" {
		t.Fatalf("unexpected event_type: %#v", result.Fields["event_type"])
	}
	if result.Fields["_sub_template_id"] != child.ID {
		t.Fatalf("unexpected sub template id: %#v", result.Fields["_sub_template_id"])
	}
}

func TestParseServiceSubTemplateRoutingByExtractedField(t *testing.T) {
	db := setupParseTestDB(t)
	service := NewParseService()

	child := models.ParseTemplate{
		Name:         "json-child",
		ParseType:    "json",
		FieldMapping: `{"event":"event_type"}`,
		Enabled:      true,
	}
	if err := db.Create(&child).Error; err != nil {
		t.Fatalf("create child template: %v", err)
	}

	routes, _ := json.Marshal([]map[string]interface{}{
		{
			"templateId": child.ID,
			"matchType":  "equals",
			"matchField": "event",
			"matchValue": "login",
		},
	})
	parent := models.ParseTemplate{
		Name:         "parent-json",
		ParseType:    "sub_template",
		SubTemplates: string(routes),
		Enabled:      true,
	}
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("create parent template: %v", err)
	}

	result, err := service.TestParse(parent.ID, `{"event":"login","result":"success"}`)
	if err != nil {
		t.Fatalf("test parse: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success, got error: %s", result.Error)
	}
	if result.Fields["event_type"] != "login" {
		t.Fatalf("unexpected event_type: %#v", result.Fields["event_type"])
	}
}

func TestParseServiceTransformFailure(t *testing.T) {
	service := NewParseService()
	template := &models.ParseTemplate{
		ParseType:      "json",
		ValueTransform: `{"count":{"type":"int"}}`,
	}

	result := service.parse(template, `{"count":"not-a-number"}`, 0)
	if result.Success {
		t.Fatal("expected transform failure")
	}
	if result.Error == "" {
		t.Fatal("expected transform failure error")
	}
}
