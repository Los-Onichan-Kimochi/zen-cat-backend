package adapter

import (
	"math"
	"strconv"

	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type AuditLog struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

func NewAuditLogAdapter(
	logger logging.Logger,
	daoAstroCatPsql *daoPostgresql.AstroCatPsqlCollection,
) *AuditLog {
	return &AuditLog{
		logger:        logger,
		DaoPostgresql: daoAstroCatPsql,
	}
}

// LogAuditEvent logs an audit event
func (a *AuditLog) LogAuditEvent(
	context schemas.AuditContext,
	event schemas.AuditEvent,
) *errors.Error {

	if err := a.DaoPostgresql.AuditLog.LogAuditEvent(
		context.UserId,
		context.UserEmail,
		model.UserRol(context.UserRole),
		model.AuditActionType(event.Action),
		model.AuditEntityType(event.EntityType),
		event.EntityId,
		event.EntityName,
		event.OldValues,
		event.NewValues,
		context.IPAddress,
		context.UserAgent,
		event.AdditionalInfo,
		event.Success,
		event.ErrorMessage,
	); err != nil {
		a.logger.Error("Failed to log audit event: ", err)
		return &errors.InternalServerError.Default
	}
	return nil
}

// GetAuditLogs retrieves audit logs with filtering and pagination
func (a *AuditLog) GetAuditLogs(filters schemas.AuditLogFilters) (*schemas.AuditLogs, *errors.Error) {
	// Default pagination values
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 50
	}
	if filters.PageSize > 200 {
		filters.PageSize = 200 // Max page size
	}

	auditLogModels, totalCount, err := a.DaoPostgresql.AuditLog.GetAuditLogs(
		filters.UserIds,
		filters.Actions,
		filters.EntityTypes,
		filters.UserRoles,
		filters.StartDate,
		filters.EndDate,
		filters.Success,
		filters.Page,
		filters.PageSize,
	)
	if err != nil {
		a.logger.Error("Failed to get audit logs: ", err)
		return nil, &errors.InternalServerError.Default
	}

	// Convert models to schemas
	auditLogs := make([]*schemas.AuditLog, len(auditLogModels))
	for i, logModel := range auditLogModels {
		auditLog := &schemas.AuditLog{
			Id:             logModel.Id,
			UserId:         logModel.UserId,
			UserEmail:      logModel.UserEmail,
			UserRole:       schemas.UserRol(logModel.UserRole),
			Action:         schemas.AuditActionType(logModel.Action),
			EntityType:     schemas.AuditEntityType(logModel.EntityType),
			EntityId:       logModel.EntityId,
			EntityName:     logModel.EntityName,
			OldValues:      logModel.OldValues,
			NewValues:      logModel.NewValues,
			IPAddress:      logModel.IPAddress,
			UserAgent:      logModel.UserAgent,
			AdditionalInfo: logModel.AdditionalInfo,
			Success:        logModel.Success,
			ErrorMessage:   logModel.ErrorMessage,
			CreatedAt:      logModel.CreatedAt,
		}

		// Include user profile if available
		if logModel.User != nil {
			auditLog.User = &schemas.UserProfile{
				Id:             logModel.User.Id,
				Name:           logModel.User.Name,
				FirstLastName:  logModel.User.FirstLastName,
				SecondLastName: logModel.User.SecondLastName,
				Email:          logModel.User.Email,
				Rol:            schemas.UserRol(logModel.User.Rol),
				ImageUrl:       logModel.User.ImageUrl,
			}
		}

		auditLogs[i] = auditLog
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(filters.PageSize)))

	return &schemas.AuditLogs{
		AuditLogs:  auditLogs,
		TotalCount: totalCount,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAuditLogById retrieves a specific audit log by ID
func (a *AuditLog) GetAuditLogById(id uuid.UUID) (*schemas.AuditLog, *errors.Error) {
	auditLogModel, err := a.DaoPostgresql.AuditLog.GetAuditLogById(id)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.AuditLogNotFound
	}

	auditLog := &schemas.AuditLog{
		Id:             auditLogModel.Id,
		UserId:         auditLogModel.UserId,
		UserEmail:      auditLogModel.UserEmail,
		UserRole:       schemas.UserRol(auditLogModel.UserRole),
		Action:         schemas.AuditActionType(auditLogModel.Action),
		EntityType:     schemas.AuditEntityType(auditLogModel.EntityType),
		EntityId:       auditLogModel.EntityId,
		EntityName:     auditLogModel.EntityName,
		OldValues:      auditLogModel.OldValues,
		NewValues:      auditLogModel.NewValues,
		IPAddress:      auditLogModel.IPAddress,
		UserAgent:      auditLogModel.UserAgent,
		AdditionalInfo: auditLogModel.AdditionalInfo,
		Success:        auditLogModel.Success,
		ErrorMessage:   auditLogModel.ErrorMessage,
		CreatedAt:      auditLogModel.CreatedAt,
	}

	// Include user profile if available
	if auditLogModel.User != nil {
		auditLog.User = &schemas.UserProfile{
			Id:             auditLogModel.User.Id,
			Name:           auditLogModel.User.Name,
			FirstLastName:  auditLogModel.User.FirstLastName,
			SecondLastName: auditLogModel.User.SecondLastName,
			Email:          auditLogModel.User.Email,
			Rol:            schemas.UserRol(auditLogModel.User.Rol),
			ImageUrl:       auditLogModel.User.ImageUrl,
		}
	}

	return auditLog, nil
}

// GetAuditStats returns statistics about audit logs
func (a *AuditLog) GetAuditStats(daysStr string, successFilter *bool) (*schemas.AuditStats, *errors.Error) {
	days := 30 // Default to 30 days
	if daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	statsData, err := a.DaoPostgresql.AuditLog.GetAuditStats(days, successFilter)
	if err != nil {
		a.logger.Error("Failed to get audit stats: ", err)
		return nil, &errors.InternalServerError.Default
	}

	// Convert interface{} maps to proper structs
	stats := &schemas.AuditStats{
		TotalEvents:  statsData["total_events"].(int64),
		SuccessCount: statsData["success_count"].(int64),
		FailureCount: statsData["failure_count"].(int64),
		ActiveUsers:  statsData["active_users"].(int64),
	}

	// Convert actions
	if actionData, ok := statsData["actions"].([]struct {
		Action string
		Count  int64
	}); ok {
		stats.Actions = make([]schemas.ActionStat, len(actionData))
		for i, action := range actionData {
			stats.Actions[i] = schemas.ActionStat{
				Action: action.Action,
				Count:  action.Count,
			}
		}
	}

	// Convert user roles
	if roleData, ok := statsData["user_roles"].([]struct {
		UserRole string
		Count    int64
	}); ok {
		stats.UserRoles = make([]schemas.UserRoleStat, len(roleData))
		for i, role := range roleData {
			stats.UserRoles[i] = schemas.UserRoleStat{
				UserRole: role.UserRole,
				Count:    role.Count,
			}
		}
	}

	// Convert entity types
	if entityData, ok := statsData["entity_types"].([]struct {
		EntityType string
		Count      int64
	}); ok {
		stats.EntityTypes = make([]schemas.EntityTypeStat, len(entityData))
		for i, entity := range entityData {
			stats.EntityTypes[i] = schemas.EntityTypeStat{
				EntityType: entity.EntityType,
				Count:      entity.Count,
			}
		}
	}

	return stats, nil
}

// DeleteOldAuditLogs deletes audit logs older than the specified number of days
func (a *AuditLog) DeleteOldAuditLogs(days int) *errors.Error {
	if err := a.DaoPostgresql.AuditLog.DeleteOldAuditLogs(days); err != nil {
		a.logger.Error("Failed to delete old audit logs: ", err)
		return &errors.InternalServerError.Default
	}
	return nil
}
