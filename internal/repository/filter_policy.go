package repository

import (
	"syslog-alert/internal/models"
	"syslog-alert/internal/service/cache"
	"syslog-alert/pkg/constants"
)

// ---- 筛选策略 ----

func CreateFilterPolicy(policy *models.FilterPolicy) error {
	err := DB().Create(policy).Error
	if err == nil {
		cache.InvalidateFilterPolicies()
		cache.InvalidateStatsCaches()
	}
	return err
}

func GetFilterPolicies() []models.FilterPolicy {
	var policies []models.FilterPolicy
	DB().Order("priority DESC").Find(&policies)
	return policies
}

func GetActiveFilterPolicyCount() int64 {
	var count int64
	DB().Model(&models.FilterPolicy{}).Where("is_active = ?", true).Count(&count)
	return count
}

func GetFilterPoliciesByDeviceID(deviceID uint) []models.FilterPolicy {
	var policies []models.FilterPolicy
	DB().Where("device_id = ? OR device_id = 0", deviceID).Order("priority DESC").Find(&policies)
	return policies
}

func GetFilterPoliciesByDeviceGroupID(groupID uint) []models.FilterPolicy {
	var policies []models.FilterPolicy
	DB().Where("device_group_id = ? OR device_group_id = 0", groupID).Order("priority DESC").Find(&policies)
	return policies
}

func GetFilterPolicyByID(id uint) (*models.FilterPolicy, error) {
	var policy models.FilterPolicy
	err := DB().First(&policy, id).Error
	return &policy, err
}

func UpdateFilterPolicy(policy *models.FilterPolicy) error {
	err := DB().Save(policy).Error
	if err == nil {
		cache.InvalidateFilterPolicies()
		cache.InvalidateStatsCaches()
	}
	return err
}

func DeleteFilterPolicy(id uint) error {
	err := DB().Delete(&models.FilterPolicy{}, id).Error
	if err == nil {
		cache.InvalidateFilterPolicies()
		cache.InvalidateStatsCaches()
	}
	return err
}

// ---- 告警策略 ----

func CreateAlertPolicy(policy *models.AlertPolicy) error {
	err := DB().Create(policy).Error
	if err == nil {
		cache.InvalidateFilterPolicies()
		cache.InvalidateStatsCaches()
	}
	return err
}

func GetAlertPolicies() []models.AlertPolicy {
	var policies []models.AlertPolicy
	DB().Find(&policies)
	return policies
}

func GetAlertPolicyByID(id uint) (*models.AlertPolicy, error) {
	var policy models.AlertPolicy
	err := DB().First(&policy, id).Error
	return &policy, err
}

func UpdateAlertPolicy(policy *models.AlertPolicy) error {
	err := DB().Save(policy).Error
	if err == nil {
		cache.InvalidateFilterPolicies()
		cache.InvalidateStatsCaches()
	}
	return err
}

func DeleteAlertPolicy(id uint) error {
	err := DB().Delete(&models.AlertPolicy{}, id).Error
	if err == nil {
		cache.InvalidateFilterPolicies()
		cache.InvalidateStatsCaches()
	}
	return err
}

func GetActiveAlertPolicies() []models.AlertPolicy {
	var policies []models.AlertPolicy
	DB().Where("is_active = ?", true).Find(&policies)
	return policies
}

func GetAlertPoliciesByFilterPolicyID(filterPolicyID uint) []models.AlertPolicy {
	var policies []models.AlertPolicy
	DB().Where("filter_policy_id = ? AND is_active = ?", filterPolicyID, true).Find(&policies)
	return policies
}

// ---- 告警规则 ----

func CreateAlertRule(rule *models.AlertRule) error {
	err := DB().Create(rule).Error
	if err == nil {
		cache.InvalidateAlertRules()
		cache.InvalidateStatsCaches()
	}
	return err
}

func GetAlertRules() []models.AlertRule {
	var rules []models.AlertRule
	DB().Find(&rules)
	return rules
}

func GetAlertRuleByID(id uint) (*models.AlertRule, error) {
	var rule models.AlertRule
	err := DB().First(&rule, id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func UpdateAlertRule(rule *models.AlertRule) error {
	err := DB().Save(rule).Error
	if err == nil {
		cache.InvalidateAlertRules()
		cache.InvalidateStatsCaches()
	}
	return err
}

func DeleteAlertRule(id uint) error {
	err := DB().Delete(&models.AlertRule{}, id).Error
	if err == nil {
		cache.InvalidateAlertRules()
		cache.InvalidateStatsCaches()
	}
	return err
}

func DeleteAlertRulesByRobotID(robotID uint) error {
	err := DB().Where("robot_id = ?", robotID).Delete(&models.AlertRule{}).Error
	if err == nil {
		cache.InvalidateAlertRules()
		cache.InvalidateStatsCaches()
	}
	return err
}

func GetAlertRulesByRobotID(robotID uint) []models.AlertRule {
	var rules []models.AlertRule
	DB().Where("robot_id = ? AND is_active = ?", robotID, true).Order("created_at ASC").Find(&rules)
	return rules
}

func GetAlertRulesByFilterPolicyID(filterPolicyID uint) []models.AlertRule {
	var rules []models.AlertRule
	DB().Where("filter_policy_id = ? AND is_active = ?", filterPolicyID, true).Find(&rules)
	return rules
}

func GetRobotsByFilterPolicyID(filterPolicyID uint) []models.Robot {
	var rules []models.AlertRule
	DB().Where("filter_policy_id = ? AND is_active = ?", filterPolicyID, true).Find(&rules)

	var robotIDs []uint
	for _, rule := range rules {
		robotIDs = append(robotIDs, rule.RobotID)
	}

	if len(robotIDs) == 0 {
		return []models.Robot{}
	}

	var robots []models.Robot
	DB().Where("id IN ? AND is_active = ? AND (platform IN ? OR platform = '')", robotIDs, true, constants.SupportedNotificationPlatforms()).Find(&robots)
	return robots
}
