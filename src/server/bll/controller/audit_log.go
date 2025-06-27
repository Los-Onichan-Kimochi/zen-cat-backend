package controller

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type AuditLog struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create AuditLog controller
func NewAuditLogController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *AuditLog {
	return &AuditLog{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// LogAuditEvent logs an audit event
func (a *AuditLog) LogAuditEvent(
	context schemas.AuditContext,
	event schemas.AuditEvent,
) *errors.Error {
	return a.Adapter.AuditLog.LogAuditEvent(context, event)
}

// GetAuditLogs retrieves audit logs with filtering and pagination (only successful operations)
func (a *AuditLog) GetAuditLogs(
	userIdsStr string,
	actionsStr string,
	entityTypesStr string,
	userRolesStr string,
	startDateStr string,
	endDateStr string,
	successStr string,
	pageStr string,
	pageSizeStr string,
) (*schemas.AuditLogs, *errors.Error) {
	filters := schemas.AuditLogFilters{
		Page:     1,
		PageSize: 50,
	}

	// Parse user IDs
	if userIdsStr != "" {
		filters.UserIds = strings.Split(userIdsStr, ",")
	}

	// Parse actions
	if actionsStr != "" {
		filters.Actions = strings.Split(actionsStr, ",")
	}

	// Parse entity types
	if entityTypesStr != "" {
		filters.EntityTypes = strings.Split(entityTypesStr, ",")
	}

	// Parse user roles
	if userRolesStr != "" {
		filters.UserRoles = strings.Split(userRolesStr, ",")
	}

	// Parse start date
	if startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = &startDate
		}
	}

	// Parse end date
	if endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filters.EndDate = &endDate
		}
	}

	// Parse success filter
	if successStr != "" {
		if success, err := strconv.ParseBool(successStr); err == nil {
			filters.Success = &success
		}
	}

	// Parse page
	if pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		}
	}

	// Parse page size
	if pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			filters.PageSize = pageSize
		}
	}

	return a.Adapter.AuditLog.GetAuditLogs(filters)
}

// GetAuditLogById retrieves a specific audit log by ID
func (a *AuditLog) GetAuditLogById(id uuid.UUID) (*schemas.AuditLog, *errors.Error) {
	return a.Adapter.AuditLog.GetAuditLogById(id)
}

// GetAuditStats returns statistics about audit logs
func (a *AuditLog) GetAuditStats(daysStr string) (*schemas.AuditStats, *errors.Error) {
	return a.Adapter.AuditLog.GetAuditStats(daysStr, nil) // No success filter
}

// GetErrorStats returns statistics about error logs
func (a *AuditLog) GetErrorStats(daysStr string) (*schemas.AuditStats, *errors.Error) {
	successFilter := false
	return a.Adapter.AuditLog.GetAuditStats(daysStr, &successFilter) // Filter for failed operations only
}

// DeleteOldAuditLogs deletes audit logs older than the specified number of days
func (a *AuditLog) DeleteOldAuditLogs(daysStr string) *errors.Error {
	days := 90 // Default to 90 days
	if daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	return a.Adapter.AuditLog.DeleteOldAuditLogs(days)
}
