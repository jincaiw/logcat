package repository

import "syslog-alert/internal/models"

// ---- 筛选策略 ----

func CreateFilterPolicy(policy *models.FilterPolicy) error {
	return DB().Create(policy).Error
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
	return DB().Save(policy).Error
}

func DeleteFilterPolicy(id uint) error {
	return DB().Delete(&models.FilterPolicy{}, id).Error
}

// ---- 告警策略 ----

func CreateAlertPolicy(policy *models.AlertPolicy) error {
	return DB().Create(policy).Error
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
	return DB().Save(policy).Error
}

func DeleteAlertPolicy(id uint) error {
	return DB().Delete(&models.AlertPolicy{}, id).Error
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
	return DB().Create(rule).Error
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
	return DB().Save(rule).Error
}

func DeleteAlertRule(id uint) error {
	return DB().Delete(&models.AlertRule{}, id).Error
}

func DeleteAlertRulesByRobotID(robotID uint) error {
	return DB().Where("robot_id = ?", robotID).Delete(&models.AlertRule{}).Error
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
	DB().Where("id IN ? AND is_active = ?", robotIDs, true).Find(&robots)
	return robots
}
