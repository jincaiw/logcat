package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// ParseService handles log parsing logic
type ParseService struct{}

// NewParseService creates a new ParseService
func NewParseService() *ParseService {
	return &ParseService{}
}

// ParseResult holds the result of parsing
type ParseResult struct {
	Success bool                   `json:"success"`
	Fields  map[string]interface{} `json:"fields"`
	Error   string                 `json:"error,omitempty"`
}

type subTemplateRoute struct {
	TemplateID  uint   `json:"templateId"`
	MatchType   string `json:"matchType"`
	MatchField  string `json:"matchField"`
	MatchValue  string `json:"matchValue"`
	Description string `json:"description"`
}

type subTemplateConfig struct {
	Routes []subTemplateRoute `json:"routes"`
}

type valueTransformRule struct {
	Field   string                 `json:"field"`
	Type    string                 `json:"type"`
	Mapping map[string]interface{} `json:"mapping"`
	Default interface{}            `json:"default"`
	Layout  string                 `json:"layout"`
}

// TestParse tests a parse template against sample log data
func (s *ParseService) TestParse(templateID uint, rawLog string) (*ParseResult, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var template models.ParseTemplate
	if err := db.First(&template, templateID).Error; err != nil {
		return nil, err
	}

	return s.parse(&template, rawLog, 0), nil
}

// ParseByTemplateID parses a raw log using a specific template
func (s *ParseService) ParseByTemplateID(templateID uint, rawLog string) (*ParseResult, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var template models.ParseTemplate
	if err := db.First(&template, templateID).Error; err != nil {
		return nil, err
	}

	return s.parse(&template, rawLog, 0), nil
}

// parse performs the actual parsing based on the template
func (s *ParseService) parse(template *models.ParseTemplate, rawLog string, depth int) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog
	if depth > 8 {
		return &ParseResult{
			Success: false,
			Fields:  fields,
			Error:   "sub-template recursion limit exceeded",
		}
	}

	switch template.ParseType {
	case "json":
		return s.parseJSON(rawLog, template)
	case "delimiter":
		return s.parseDelimiter(rawLog, template)
	case "kv":
		return s.parseKV(rawLog, template)
	case "regex":
		return s.parseRegex(rawLog, template)
	case "syslog_json":
		return s.parseSyslogJSON(rawLog, template)
	case "sub_template":
		return s.parseSubTemplate(rawLog, template, depth)
	default:
		return &ParseResult{
			Success: false,
			Fields:  fields,
			Error:   fmt.Sprintf("unsupported parse type: %s", template.ParseType),
		}
	}
}

func (s *ParseService) parseJSON(rawLog string, template *models.ParseTemplate) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(rawLog), &data); err != nil {
		return &ParseResult{Success: false, Fields: fields, Error: fmt.Sprintf("JSON parse failed: %v", err)}
	}

	// Apply field mapping if configured
	if template.FieldMapping != "" {
		var mapping map[string]string
		if err := json.Unmarshal([]byte(template.FieldMapping), &mapping); err == nil {
			for srcField, dstField := range mapping {
				if val, ok := data[srcField]; ok {
					fields[dstField] = val
				}
			}
		}
	} else {
		for k, v := range data {
			fields[k] = v
		}
	}

	return s.finalizeFields(template, fields)
}

func (s *ParseService) parseDelimiter(rawLog string, template *models.ParseTemplate) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	if template.Delimiter == "" {
		return &ParseResult{Success: false, Fields: fields, Error: "delimiter is not set"}
	}

	parts := strings.Split(rawLog, template.Delimiter)
	if template.FieldMapping != "" {
		var mapping []string
		if err := json.Unmarshal([]byte(template.FieldMapping), &mapping); err == nil {
			for i, fieldName := range mapping {
				if i < len(parts) {
					fields[fieldName] = strings.TrimSpace(parts[i])
				}
			}
		}
	} else {
		for i, part := range parts {
			fields[fmt.Sprintf("field_%d", i+1)] = strings.TrimSpace(part)
		}
	}

	return s.finalizeFields(template, fields)
}

func (s *ParseService) parseKV(rawLog string, template *models.ParseTemplate) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	if template.Delimiter == "" {
		return &ParseResult{Success: false, Fields: fields, Error: "delimiter is not set"}
	}

	pairs := strings.Split(rawLog, template.Delimiter)
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		idx := strings.Index(pair, "=")
		if idx > 0 {
			key := strings.TrimSpace(pair[:idx])
			value := strings.TrimSpace(pair[idx+1:])
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			fields[key] = value
		}
	}

	return s.finalizeFields(template, fields)
}

func (s *ParseService) parseRegex(rawLog string, template *models.ParseTemplate) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	if template.HeaderRegex == "" {
		return &ParseResult{Success: false, Fields: fields, Error: "header regex is not set"}
	}

	pattern, err := regexp.Compile(template.HeaderRegex)
	if err != nil {
		return &ParseResult{Success: false, Fields: fields, Error: fmt.Sprintf("invalid regex: %v", err)}
	}

	if rawLog != "" && len(rawLog) > 65536 {
		return &ParseResult{Success: false, Fields: fields, Error: "input exceeds maximum length"}
	}

	matches := pattern.FindStringSubmatch(rawLog)
	if matches == nil {
		return &ParseResult{Success: false, Fields: fields, Error: "regex did not match"}
	}

	captured := make(map[string]interface{})
	names := pattern.SubexpNames()
	for idx := 1; idx < len(matches); idx++ {
		value := matches[idx]
		captured[strconv.Itoa(idx)] = value
		if idx < len(names) && names[idx] != "" {
			captured[names[idx]] = value
		}
	}

	if template.FieldMapping != "" {
		var mapping map[string]string
		if err := json.Unmarshal([]byte(template.FieldMapping), &mapping); err == nil {
			for srcField, dstField := range mapping {
				if val, ok := captured[srcField]; ok {
					fields[dstField] = val
				}
			}
			return s.finalizeFields(template, fields)
		}

		var orderedFields []string
		if err := json.Unmarshal([]byte(template.FieldMapping), &orderedFields); err == nil {
			for idx, fieldName := range orderedFields {
				if idx+1 < len(matches) {
					fields[fieldName] = matches[idx+1]
				}
			}
			return s.finalizeFields(template, fields)
		}
	}

	for idx := 1; idx < len(matches); idx++ {
		key := fmt.Sprintf("group_%d", idx)
		if idx < len(names) && names[idx] != "" {
			key = names[idx]
		}
		fields[key] = matches[idx]
	}

	return s.finalizeFields(template, fields)
}

func (s *ParseService) parseSyslogJSON(rawLog string, template *models.ParseTemplate) *ParseResult {
	// Parse as syslog header + JSON body
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	// Find JSON body - look for first '{'
	jsonStart := strings.Index(rawLog, "{")
	if jsonStart < 0 {
		return &ParseResult{Success: false, Fields: fields, Error: "json body not found"}
	}

	// Header part
	header := strings.TrimSpace(rawLog[:jsonStart])
	fields["syslog_header"] = header

	// JSON body
	body := strings.TrimSpace(rawLog[jsonStart:])
	if body == "" || body == "{}" {
		return &ParseResult{Success: false, Fields: fields, Error: "json body is empty"}
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return &ParseResult{Success: false, Fields: fields, Error: fmt.Sprintf("JSON body parse failed: %v", err)}
	}
	for k, v := range data {
		fields[k] = v
	}

	return s.finalizeFields(template, fields)
}

func (s *ParseService) parseSubTemplate(rawLog string, template *models.ParseTemplate, depth int) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	if template.SubTemplates == "" {
		return &ParseResult{Success: false, Fields: fields, Error: "sub-templates are not configured"}
	}

	routes, err := parseSubTemplateRoutes(template.SubTemplates)
	if err != nil {
		return &ParseResult{Success: false, Fields: fields, Error: err.Error()}
	}

	candidateData := buildSubTemplateCandidateData(rawLog)
	for _, route := range routes {
		if route.TemplateID == 0 {
			continue
		}
		if !matchSubTemplateRoute(route, candidateData) {
			continue
		}

		db := database.GetDB()
		if db == nil {
			return &ParseResult{Success: false, Fields: fields, Error: "database not available"}
		}
		var childTemplate models.ParseTemplate
		if err := db.First(&childTemplate, route.TemplateID).Error; err != nil {
			return &ParseResult{Success: false, Fields: fields, Error: err.Error()}
		}
		result := s.parse(&childTemplate, rawLog, depth+1)
		if result.Fields == nil {
			result.Fields = map[string]interface{}{}
		}
		result.Fields["_sub_template_id"] = route.TemplateID
		if route.Description != "" {
			result.Fields["_sub_template_description"] = route.Description
		}
		return result
	}

	return &ParseResult{Success: false, Fields: fields, Error: "no sub-template route matched"}
}

func buildSubTemplateCandidateData(rawLog string) map[string]interface{} {
	candidateData := map[string]interface{}{
		"raw_message": rawLog,
	}

	if raw := strings.TrimSpace(rawLog); raw != "" {
		if strings.HasPrefix(raw, "{") {
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(raw), &jsonData); err == nil {
				for k, v := range jsonData {
					candidateData[k] = v
				}
			}
		} else if jsonStart := strings.Index(raw, "{"); jsonStart > 0 {
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(raw[jsonStart:]), &jsonData); err == nil {
				for k, v := range jsonData {
					candidateData[k] = v
				}
			}
		}
	}

	for _, token := range strings.Fields(rawLog) {
		kv := strings.SplitN(token, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.Trim(strings.TrimSpace(kv[1]), "\"'")
		if key == "" {
			continue
		}
		candidateData[key] = value
	}

	return candidateData
}

func (s *ParseService) finalizeFields(template *models.ParseTemplate, fields map[string]interface{}) *ParseResult {
	if err := applyValueTransforms(template.ValueTransform, fields); err != nil {
		return &ParseResult{Success: false, Fields: fields, Error: err.Error()}
	}
	return &ParseResult{Success: true, Fields: fields}
}

func parseSubTemplateRoutes(raw string) ([]subTemplateRoute, error) {
	var routes []subTemplateRoute
	if err := json.Unmarshal([]byte(raw), &routes); err == nil {
		return routes, nil
	}

	var config subTemplateConfig
	if err := json.Unmarshal([]byte(raw), &config); err == nil && len(config.Routes) > 0 {
		return config.Routes, nil
	}

	return nil, errors.New("invalid sub-template configuration")
}

func matchSubTemplateRoute(route subTemplateRoute, data map[string]interface{}) bool {
	field := route.MatchField
	if field == "" {
		field = "raw_message"
	}
	matchType := strings.ToLower(strings.TrimSpace(route.MatchType))
	if matchType == "" {
		matchType = "contains"
	}

	value, exists := data[field]
	valueStr := fmt.Sprintf("%v", value)
	switch matchType {
	case "default":
		return true
	case "exists":
		return exists
	case "equals", "==":
		return exists && valueStr == route.MatchValue
	case "contains":
		return exists && strings.Contains(valueStr, route.MatchValue)
	case "regex":
		if len(route.MatchValue) > 1024 {
			return false
		}
		pattern, err := regexp.Compile(route.MatchValue)
		if err != nil {
			return false
		}
		valueStrLen := len(valueStr)
		if valueStrLen > 65536 {
			return false
		}
		return exists && pattern.MatchString(valueStr)
	default:
		return false
	}
}

func applyValueTransforms(raw string, fields map[string]interface{}) error {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	rules, err := parseValueTransformRules(raw)
	if err != nil {
		return fmt.Errorf("invalid value transform: %w", err)
	}

	for _, rule := range rules {
		if rule.Field == "" {
			continue
		}
		value, exists := fields[rule.Field]
		if !exists {
			if rule.Default != nil {
				fields[rule.Field] = rule.Default
			}
			continue
		}

		transformed, err := transformValue(value, rule)
		if err != nil {
			return fmt.Errorf("field %s transform failed: %w", rule.Field, err)
		}
		fields[rule.Field] = transformed
	}

	return nil
}

func parseValueTransformRules(raw string) ([]valueTransformRule, error) {
	var rules []valueTransformRule
	if err := json.Unmarshal([]byte(raw), &rules); err == nil {
		return rules, nil
	}

	var compact map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &compact); err != nil {
		return nil, err
	}

	rules = make([]valueTransformRule, 0, len(compact))
	for field, item := range compact {
		switch typed := item.(type) {
		case string:
			rules = append(rules, valueTransformRule{Field: field, Type: typed})
		default:
			payload, _ := json.Marshal(typed)
			var rule valueTransformRule
			if err := json.Unmarshal(payload, &rule); err != nil {
				return nil, err
			}
			rule.Field = field
			rules = append(rules, rule)
		}
	}

	return rules, nil
}

func transformValue(value interface{}, rule valueTransformRule) (interface{}, error) {
	switch strings.ToLower(strings.TrimSpace(rule.Type)) {
	case "", "noop":
		return value, nil
	case "trim":
		return strings.TrimSpace(fmt.Sprintf("%v", value)), nil
	case "lower":
		return strings.ToLower(fmt.Sprintf("%v", value)), nil
	case "upper":
		return strings.ToUpper(fmt.Sprintf("%v", value)), nil
	case "string":
		return fmt.Sprintf("%v", value), nil
	case "int":
		return strconv.Atoi(strings.TrimSpace(fmt.Sprintf("%v", value)))
	case "float":
		return strconv.ParseFloat(strings.TrimSpace(fmt.Sprintf("%v", value)), 64)
	case "bool":
		return strconv.ParseBool(strings.TrimSpace(strings.ToLower(fmt.Sprintf("%v", value))))
	case "map":
		mapped, ok := rule.Mapping[fmt.Sprintf("%v", value)]
		if ok {
			return mapped, nil
		}
		if rule.Default != nil {
			return rule.Default, nil
		}
		return value, nil
	case "datetime":
		layout := rule.Layout
		if layout == "" {
			layout = time.RFC3339
		}
		parsed, err := time.Parse(layout, strings.TrimSpace(fmt.Sprintf("%v", value)))
		if err != nil {
			return nil, err
		}
		return parsed.Format(time.RFC3339), nil
	default:
		return nil, fmt.Errorf("unsupported transform type: %s", rule.Type)
	}
}
