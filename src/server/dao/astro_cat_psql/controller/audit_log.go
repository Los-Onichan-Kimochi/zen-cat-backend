package controller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type AuditLog struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewAuditLogController(logger logging.Logger, postgresqlDB *gorm.DB) *AuditLog {
	return &AuditLog{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// LogAuditEvent logs an audit event to the database
func (a *AuditLog) LogAuditEvent(
	userId uuid.UUID,
	userEmail string,
	userRole model.UserRol,
	action model.AuditActionType,
	entityType model.AuditEntityType,
	entityId *uuid.UUID,
	entityName *string,
	oldValues interface{},
	newValues interface{},
	ipAddress string,
	userAgent *string,
	additionalInfo *string,
	success bool,
	errorMessage *string,
) error {
	auditLog := &model.AuditLog{
		Id:             uuid.New(),
		UserId:         userId,
		UserEmail:      userEmail,
		UserRole:       userRole,
		Action:         action,
		EntityType:     entityType,
		EntityId:       entityId,
		EntityName:     entityName,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		AdditionalInfo: additionalInfo,
		Success:        success,
		ErrorMessage:   errorMessage,
		CreatedAt:      time.Now(),
	}

	// Serialize old values to JSON if provided
	if oldValues != nil {
		if oldJSON, err := json.Marshal(oldValues); err == nil {
			oldStr := string(oldJSON)
			auditLog.OldValues = &oldStr
		}
	}

	// Serialize new values to JSON if provided
	if newValues != nil {
		if newJSON, err := json.Marshal(newValues); err == nil {
			newStr := string(newJSON)
			auditLog.NewValues = &newStr
		}
	}
	if err := a.PostgresqlDB.Create(auditLog).Error; err != nil {
		a.logger.Error("Failed to log audit event: ", err)
		return err
	}

	return nil
}

// GetAuditLogs retrieves audit logs with filtering and pagination
func (a *AuditLog) GetAuditLogs(
	userIds []string,
	actions []string,
	entityTypes []string,
	userRoles []string,
	startDate *time.Time,
	endDate *time.Time,
	success *bool,
	page int,
	pageSize int,
) ([]*model.AuditLog, int64, error) {
	var auditLogs []*model.AuditLog
	var totalCount int64

	// Build query
	query := a.PostgresqlDB.Model(&model.AuditLog{}).Preload("User")

	// Apply filters
	if len(userIds) > 0 {
		query = query.Where("user_id IN ?", userIds)
	}

	if len(actions) > 0 {
		query = query.Where("action IN ?", actions)
	}

	if len(entityTypes) > 0 {
		query = query.Where("entity_type IN ?", entityTypes)
	}

	if len(userRoles) > 0 {
		query = query.Where("user_role IN ?", userRoles)
	}

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	if success != nil {
		query = query.Where("success = ?", *success)
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and order
	offset := (page - 1) * pageSize
	if err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}

	return auditLogs, totalCount, nil
}

// GetAuditLogById retrieves a specific audit log by ID
func (a *AuditLog) GetAuditLogById(id uuid.UUID) (*model.AuditLog, error) {
	var auditLog model.AuditLog
	if err := a.PostgresqlDB.
		Preload("User").
		First(&auditLog, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &auditLog, nil
}

// DeleteOldAuditLogs deletes audit logs older than the specified number of days
func (a *AuditLog) DeleteOldAuditLogs(daysOld int) error {
	cutoffDate := time.Now().AddDate(0, 0, -daysOld)
	result := a.PostgresqlDB.Where("created_at < ?", cutoffDate).Delete(&model.AuditLog{})

	if result.Error != nil {
		return result.Error
	}

	a.logger.Info(fmt.Sprintf("Deleted %d old audit logs", result.RowsAffected))
	return nil
}

// GetAuditStats returns statistics about audit logs
func (a *AuditLog) GetAuditStats(days int, successFilter *bool) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Build base query with date range
	baseQuery := a.PostgresqlDB.Model(&model.AuditLog{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate)

	// Add success filter if provided
	if successFilter != nil {
		baseQuery = baseQuery.Where("success = ?", *successFilter)
	}

	// Total events in period
	var totalEvents int64
	if err := baseQuery.Count(&totalEvents).Error; err != nil {
		return nil, err
	}
	stats["total_events"] = totalEvents

	// Events by action type
	var actionStats []struct {
		Action string
		Count  int64
	}
	actionQuery := a.PostgresqlDB.Model(&model.AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	if successFilter != nil {
		actionQuery = actionQuery.Where("success = ?", *successFilter)
	}
	if err := actionQuery.Group("action").Order("count DESC").Scan(&actionStats).Error; err != nil {
		return nil, err
	}
	stats["actions"] = actionStats

	// Events by user role
	var roleStats []struct {
		UserRole string
		Count    int64
	}
	roleQuery := a.PostgresqlDB.Model(&model.AuditLog{}).
		Select("user_role, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	if successFilter != nil {
		roleQuery = roleQuery.Where("success = ?", *successFilter)
	}
	if err := roleQuery.Group("user_role").Order("count DESC").Scan(&roleStats).Error; err != nil {
		return nil, err
	}
	stats["user_roles"] = roleStats

	// Events by entity type
	var entityStats []struct {
		EntityType string
		Count      int64
	}
	entityQuery := a.PostgresqlDB.Model(&model.AuditLog{}).
		Select("entity_type, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	if successFilter != nil {
		entityQuery = entityQuery.Where("success = ?", *successFilter)
	}
	if err := entityQuery.Group("entity_type").Order("count DESC").Scan(&entityStats).Error; err != nil {
		return nil, err
	}
	stats["entity_types"] = entityStats

	// Success/failure ratio
	var successCount, failureCount int64
	if successFilter != nil {
		// If filtering by success, all events have the same success status
		if *successFilter {
			successCount = totalEvents
			failureCount = 0
		} else {
			successCount = 0
			failureCount = totalEvents
		}
	} else {
		// Count both success and failure when no filter is applied
		a.PostgresqlDB.Model(&model.AuditLog{}).
			Where("created_at >= ? AND created_at <= ? AND success = ?", startDate, endDate, true).
			Count(&successCount)
		a.PostgresqlDB.Model(&model.AuditLog{}).
			Where("created_at >= ? AND created_at <= ? AND success = ?", startDate, endDate, false).
			Count(&failureCount)
	}

	stats["success_count"] = successCount
	stats["failure_count"] = failureCount

	// Calcular usuarios activos (únicos que no sean "sin autenticar")
	var activeUsers int64
	activeUsersQuery := a.PostgresqlDB.Model(&model.AuditLog{}).
		Select("COUNT(DISTINCT user_email) as active_users").
		Where("created_at >= ? AND created_at <= ? AND user_email != ? AND user_email != ?",
			startDate, endDate, "sin autenticar", "")
	if err := activeUsersQuery.Scan(&activeUsers).Error; err != nil {
		a.logger.Warn("Failed to calculate active users: ", err)
		activeUsers = 0
	}
	stats["active_users"] = activeUsers

	return stats, nil
}
