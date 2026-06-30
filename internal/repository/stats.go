package repository

import (
	"encoding/json"
	"fmt"

	"syslog-alert/internal/models"
	"syslog-alert/internal/service/parser"

	"gorm.io/gorm"
)

// GetFieldStats 获取字段统计
func GetFieldStats(req models.FieldStatsRequest) models.FieldStatsResult {
	result := models.FieldStatsResult{
		Field: req.Field,
	}

	query := DB().Model(&models.SyslogLog{}).Session(&gorm.Session{})

	if req.DeviceID > 0 {
		query = query.Where("device_id = ?", req.DeviceID)
	}
	if req.FilterPolicyID > 0 {
		query = query.Where("matched_policy_id = ?", req.FilterPolicyID)
	}
	if req.StartTime != "" {
		query = query.Where("received_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("received_at <= ?", req.EndTime)
	}

	var totalLogs int64
	query.Count(&totalLogs)
	result.TotalLogs = totalLogs

	if req.TopN <= 0 {
		req.TopN = 10
	}

	type FieldCount struct {
		Value    string `json:"value"`
		Count    int64  `json:"count"`
		LastSeen string `json:"lastSeen"`
	}

	var fieldCounts []FieldCount

	fieldExpr := fmt.Sprintf("json_extract(parsed_data, '$.%s')", req.Field)

	baseQuery := DB().Model(&models.SyslogLog{}).
		Where("parsed_data IS NOT NULL AND parsed_data != ''").
		Where(fieldExpr + " IS NOT NULL").
		Where(fieldExpr + " != ''")

	if req.DeviceID > 0 {
		baseQuery = baseQuery.Where("device_id = ?", req.DeviceID)
	}
	if req.FilterPolicyID > 0 {
		baseQuery = baseQuery.Where("matched_policy_id = ?", req.FilterPolicyID)
	}
	if req.StartTime != "" {
		baseQuery = baseQuery.Where("received_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		baseQuery = baseQuery.Where("received_at <= ?", req.EndTime)
	}

	var fieldTotal int64
	baseQuery.Count(&fieldTotal)

	rows, err := baseQuery.
		Select(
			fieldExpr+" as value",
			"COUNT(*) as count",
			"MAX(received_at) as last_seen",
		).
		Group(fieldExpr).
		Order("count DESC").
		Limit(req.TopN).
		Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var fc FieldCount
			if err := rows.Scan(&fc.Value, &fc.Count, &fc.LastSeen); err == nil {
				fieldCounts = append(fieldCounts, fc)
			}
		}
	}

	result.Items = make([]models.StatsItem, 0, len(fieldCounts))

	var uniqueCount int64
	countQuery := DB().Model(&models.SyslogLog{}).
		Where("parsed_data IS NOT NULL AND parsed_data != ''").
		Where(fieldExpr + " IS NOT NULL").
		Where(fieldExpr + " != ''")
	if req.DeviceID > 0 {
		countQuery = countQuery.Where("device_id = ?", req.DeviceID)
	}
	if req.FilterPolicyID > 0 {
		countQuery = countQuery.Where("matched_policy_id = ?", req.FilterPolicyID)
	}
	if req.StartTime != "" {
		countQuery = countQuery.Where("received_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		countQuery = countQuery.Where("received_at <= ?", req.EndTime)
	}
	countQuery.Distinct(fieldExpr).Count(&uniqueCount)
	result.UniqueCount = uniqueCount

	for _, fc := range fieldCounts {
		percent := "0%"
		if fieldTotal > 0 {
			p := float64(fc.Count) / float64(fieldTotal) * 100
			percent = fmt.Sprintf("%.1f%%", p)
		}

		result.Items = append(result.Items, models.StatsItem{
			Value:    fc.Value,
			Location: "",
			Count:    fc.Count,
			Percent:  percent,
			LastSeen: fc.LastSeen,
		})
	}

	return result
}

// GetAvailableStatsFields 获取可统计字段列表
func GetAvailableStatsFields(policyID uint) []models.StatsField {
	if policyID == 0 {
		return []models.StatsField{}
	}

	var policy models.FilterPolicy
	if err := DB().First(&policy, policyID).Error; err != nil {
		return []models.StatsField{}
	}

	fieldMap := make(map[string]string)
	hasSubTemplates := false

	if policy.ParseTemplateID > 0 {
		var parseTemplate models.ParseTemplate
		if err := DB().First(&parseTemplate, policy.ParseTemplateID).Error; err == nil {
			if parseTemplate.ParseType == "smart_delimiter" && parseTemplate.FieldMapping != "" {
				var fieldMappingData map[string]interface{}
				if err := json.Unmarshal([]byte(parseTemplate.FieldMapping), &fieldMappingData); err == nil {
					if subTemplatesRaw, ok := fieldMappingData["subTemplates"]; ok {
						hasSubTemplates = true
						if subTemplatesMap, ok := subTemplatesRaw.(map[string]interface{}); ok {
							fieldKeyMap := map[string]string{
								"alertNameField":    "alertName",
								"attackIPField":     "attackIP",
								"victimIPField":     "victimIP",
								"alertTimeField":    "alertTime",
								"severityField":     "severity",
								"attackResultField": "attackResult",
							}
							displayNameMap := map[string]string{
								"alertName":    "告警名称",
								"attackIP":     "攻击IP",
								"victimIP":     "受害IP",
								"alertTime":    "告警时间",
								"severity":     "威胁等级",
								"attackResult": "攻击结果",
							}
							for _, subRaw := range subTemplatesMap {
								if sub, ok := subRaw.(map[string]interface{}); ok {
									for fieldKey, fieldName := range fieldKeyMap {
										if _, exists := sub[fieldKey]; exists {
											fieldMap[fieldName] = displayNameMap[fieldName]
										}
									}
								}
							}
						}
					}
				}
			} else if parseTemplate.FieldMapping != "" {
				var simpleMapping map[string]string
				if err := json.Unmarshal([]byte(parseTemplate.FieldMapping), &simpleMapping); err == nil {
					for fieldName, displayName := range simpleMapping {
						if fieldName != "" {
							fieldMap[fieldName] = displayName
						}
					}
				} else {
					var complexMapping map[string]map[string]interface{}
					if err := json.Unmarshal([]byte(parseTemplate.FieldMapping), &complexMapping); err == nil {
						for targetField := range complexMapping {
							if targetField != "" {
								fieldMap[targetField] = targetField
							}
						}
					}
				}
			}

			if parseTemplate.SubTemplates != "" {
				var subTemplates []parser.SubTemplateConfig
				if err := json.Unmarshal([]byte(parseTemplate.SubTemplates), &subTemplates); err == nil {
					hasSubTemplates = true
					displayNameMap := map[string]string{
						"alertName":    "告警名称",
						"attackIp":     "攻击IP",
						"victimIp":     "受害IP",
						"alertTime":    "告警时间",
						"severity":     "威胁等级",
						"attackResult": "攻击结果",
					}
					for _, st := range subTemplates {
						if st.AlertNameField > 0 {
							fieldMap["alertName"] = displayNameMap["alertName"]
						}
						if st.AttackIPField > 0 {
							fieldMap["attackIp"] = displayNameMap["attackIp"]
						}
						if st.VictimIPField > 0 {
							fieldMap["victimIp"] = displayNameMap["victimIp"]
						}
						if st.AlertTimeField > 0 {
							fieldMap["alertTime"] = displayNameMap["alertTime"]
						}
						if st.SeverityField > 0 {
							fieldMap["severity"] = displayNameMap["severity"]
						}
						if st.AttackResultField > 0 {
							fieldMap["attackResult"] = displayNameMap["attackResult"]
						}
						for _, cf := range st.CustomFields {
							if cf.Name != "" {
								fieldMap[cf.Name] = cf.Name
							}
						}
					}
				}
			}
		}
	}

	var fields []models.StatsField

	if hasSubTemplates {
		fixedFields := []string{"alertName", "attackIP", "victimIP", "alertTime", "severity", "attackResult"}
		for _, f := range fixedFields {
			if displayName, ok := fieldMap[f]; ok {
				fields = append(fields, models.StatsField{Name: f, DisplayName: displayName})
			} else {
				fields = append(fields, models.StatsField{Name: f, DisplayName: f})
			}
		}
	} else {
		for fieldName, displayName := range fieldMap {
			fields = append(fields, models.StatsField{Name: fieldName, DisplayName: displayName})
		}
	}

	return fields
}
