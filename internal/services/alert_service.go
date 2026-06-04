package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// AlertService handles alert orchestration
type AlertService struct{}

// NewAlertService creates a new AlertService
func NewAlertService() *AlertService {
	return &AlertService{}
}

// ListAlertRecords returns paginated alert records
func (s *AlertService) ListAlertRecords(page, pageSize int, status, channelType string) ([]models.AlertRecord, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not available")
	}

	var records []models.AlertRecord
	var total int64

	query := db.Model(&models.AlertRecord{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if channelType != "" {
		query = query.Where("channel_type = ?", channelType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").
		Preload("AlertRule").Preload("PushConfig").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetAlertRecord returns a single alert record
func (s *AlertService) GetAlertRecord(id uint) (*models.AlertRecord, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var record models.AlertRecord
	if err := db.Preload("AlertRule").Preload("PushConfig").First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// ProcessAlert orchestrates the alert pipeline: parse -> filter -> push
func (s *AlertService) ProcessAlert(logID string, sourceIP, rawMessage string, deviceID *uint) (*models.AlertRecord, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	// Find active alert rules
	var rules []models.AlertRule
	if err := db.Where("enabled = ?", true).Preload("FilterPolicy").Preload("PushConfig").Preload("OutputTemplate").Find(&rules).Error; err != nil {
		return nil, err
	}

	var lastRecord *models.AlertRecord

	if len(rules) == 0 {
		return nil, nil
	}

	// For each rule, check if applicable, then push
	for _, rule := range rules {
		if rule.PushConfigID == nil {
			continue
		}

		now := time.Now()
		record := models.AlertRecord{
			LogID:        logID,
			AlertRuleID:  &rule.ID,
			PushConfigID: rule.PushConfigID,
			ChannelType:  rule.ChannelType,
			Status:       "pending",
			SentAt:       &now,
		}

		if err := db.Create(&record).Error; err != nil {
			continue
		}
		lastRecord = &record

		// Execute push
		pushSvc := NewPushService()
		data := map[string]interface{}{
			"log_id":      logID,
			"source_ip":   sourceIP,
			"raw_message": rawMessage,
			"device_id":   deviceID,
			"alert_rule":  rule.Name,
			"timestamp":   now.Format(time.RFC3339),
		}

		var pushConfig models.PushConfig
		if err := db.First(&pushConfig, *rule.PushConfigID).Error; err != nil {
			record.Status = "failed"
			record.ErrorMessage = fmt.Sprintf("push config not found: %v", err)
			db.Save(&record)
			continue
		}

		switch pushConfig.Type {
		case "http":
			result, err := pushSvc.ExecutePush(*rule.PushConfigID, data)
			if err != nil || !result.Success {
				record.Status = "failed"
				if err != nil {
					record.ErrorMessage = err.Error()
				} else {
					record.ErrorMessage = result.ErrorMessage
				}
				if result != nil {
					record.ResponseStatusCode = result.StatusCode
					record.ResponseSummary = result.ResponseBody
				}
			} else {
				record.Status = "success"
				record.ResponseStatusCode = result.StatusCode
				record.ResponseSummary = result.ResponseBody
			}
		case "email":
			emailSvc := NewEmailService()
			if result, _ := emailSvc.SendEmail(*rule.PushConfigID, data); result != nil && result.Success {
				record.Status = "success"
			} else {
				record.Status = "failed"
				if result != nil {
					record.ErrorMessage = result.ErrorMessage
				}
			}
		case "syslog":
			syslogSvc := NewSyslogForwardService()
			if result, _ := syslogSvc.Forward(*rule.PushConfigID, data); result != nil && result.Success {
				record.Status = "success"
			} else {
				record.Status = "failed"
				if result != nil {
					record.ErrorMessage = result.ErrorMessage
				}
			}
		}

		db.Save(&record)
	}

	return lastRecord, nil
}

// CreateDisposition creates an alert disposition record
func (s *AlertService) CreateDisposition(alertRecordID *uint, aggregatedAlertID *uint, status, note string, operatorID *uint, operatorName string) (*models.AlertDisposition, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	now := time.Now()
	disposition := models.AlertDisposition{
		AlertRecordID:     alertRecordID,
		AggregatedAlertID: aggregatedAlertID,
		Status:            status,
		Note:              note,
		OperatorID:        operatorID,
		OperatorName:      operatorName,
		OperatedAt:        &now,
	}

	if err := db.Create(&disposition).Error; err != nil {
		return nil, err
	}

	if alertRecordID != nil {
		if err := db.Model(&models.AlertRecord{}).
			Where("id = ?", *alertRecordID).
			Updates(map[string]interface{}{"disposition_status": status}).Error; err != nil {
			return nil, err
		}
	}

	if aggregatedAlertID != nil {
		if err := db.Model(&models.AggregatedAlert{}).
			Where("id = ?", *aggregatedAlertID).
			Updates(map[string]interface{}{"status": status}).Error; err != nil {
			return nil, err
		}
	}

	return &disposition, nil
}

// ListDispositions returns dispositions for an alert
func (s *AlertService) ListDispositions(alertRecordID uint) ([]models.AlertDisposition, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var dispositions []models.AlertDisposition
	if err := db.Where("alert_record_id = ?", alertRecordID).
		Order("created_at DESC").Find(&dispositions).Error; err != nil {
		return nil, err
	}

	return dispositions, nil
}

// ListAllDispositions returns paginated disposition records.
func (s *AlertService) ListAllDispositions(page, pageSize int, status string) ([]models.AlertDisposition, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not available")
	}

	var dispositions []models.AlertDisposition
	var total int64

	query := db.Model(&models.AlertDisposition{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&dispositions).Error; err != nil {
		return nil, 0, err
	}

	return dispositions, total, nil
}
