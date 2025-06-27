package audit_log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestLogAuditEventSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid audit context and event data
		WHEN:  LogAuditEvent is called with valid parameters
		THEN:  The audit event should be logged successfully
	*/
	// GIVEN
	auditLogController, _, db := controllerTest.NewAuditLogControllerTestWrapper(t)

	// Create a user using factory
	user := factories.NewUserModel(db)

	// Create audit context
	userAgent := "Test User Agent"
	auditContext := schemas.AuditContext{
		UserId:    user.Id,
		UserEmail: user.Email,
		UserRole:  schemas.UserRol(user.Rol),
		IPAddress: "127.0.0.1",
		UserAgent: &userAgent,
	}

	// Create audit event
	auditEvent := schemas.AuditEvent{
		Action:     schemas.AuditActionCreate,
		EntityType: schemas.AuditEntityUser,
		EntityId:   &user.Id,
		EntityName: &user.Name,
		Success:    true,
	}

	// WHEN
	result := auditLogController.LogAuditEvent(auditContext, auditEvent)

	// THEN
	assert.Nil(t, result)

	// Verify the audit log was created
	var auditLog model.AuditLog
	err := db.Where("user_id = ? AND action = ?", user.Id, schemas.AuditActionCreate).First(&auditLog).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Id, auditLog.UserId)
	assert.Equal(t, user.Email, auditLog.UserEmail)
	assert.Equal(t, model.AuditActionCreate, auditLog.Action)
	assert.Equal(t, model.AuditEntityUser, auditLog.EntityType)
	assert.True(t, auditLog.Success)
}

func TestLogAuditEventWithFailure(t *testing.T) {
	/*
		GIVEN: Audit event with failure status
		WHEN:  LogAuditEvent is called with failure event
		THEN:  The audit event should be logged with error details
	*/
	// GIVEN
	auditLogController, _, db := controllerTest.NewAuditLogControllerTestWrapper(t)

	// Create a user using factory
	user := factories.NewUserModel(db)

	// Create audit context
	auditContext := schemas.AuditContext{
		UserId:    user.Id,
		UserEmail: user.Email,
		UserRole:  schemas.UserRol(user.Rol),
		IPAddress: "127.0.0.1",
	}

	// Create audit event with failure
	errorMessage := "Failed to create user"
	auditEvent := schemas.AuditEvent{
		Action:       schemas.AuditActionCreate,
		EntityType:   schemas.AuditEntityUser,
		Success:      false,
		ErrorMessage: &errorMessage,
	}

	// WHEN
	result := auditLogController.LogAuditEvent(auditContext, auditEvent)

	// THEN
	assert.Nil(t, result)

	// Verify the audit log was created with failure
	var auditLog model.AuditLog
	err := db.Where("user_id = ? AND action = ?", user.Id, schemas.AuditActionCreate).First(&auditLog).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Id, auditLog.UserId)
	assert.False(t, auditLog.Success)
	assert.NotNil(t, auditLog.ErrorMessage)
	assert.Equal(t, errorMessage, *auditLog.ErrorMessage)
}

func TestLogAuditEventWithSystemAction(t *testing.T) {
	/*
		GIVEN: System-level audit event with admin user
		WHEN:  LogAuditEvent is called with system context
		THEN:  The audit event should be logged successfully
	*/
	// GIVEN
	auditLogController, _, db := controllerTest.NewAuditLogControllerTestWrapper(t)

	// Create an admin user to represent system actions
	adminRole := model.UserRolAdmin
	systemUser := factories.NewUserModel(db, factories.UserModelF{
		Rol: &adminRole,
	})

	// Create system audit context
	auditContext := schemas.AuditContext{
		UserId:    systemUser.Id,
		UserEmail: systemUser.Email,
		UserRole:  schemas.UserRol(systemUser.Rol),
		IPAddress: "127.0.0.1",
	}

	// Create system audit event
	additionalInfo := "System cleanup operation"
	auditEvent := schemas.AuditEvent{
		Action:         schemas.AuditActionDelete,
		EntityType:     schemas.AuditEntityUser,
		Success:        true,
		AdditionalInfo: &additionalInfo,
	}

	// WHEN
	result := auditLogController.LogAuditEvent(auditContext, auditEvent)

	// THEN
	assert.Nil(t, result)

	// Verify the audit log was created
	var auditLog model.AuditLog
	err := db.Where("user_id = ? AND action = ?", systemUser.Id, schemas.AuditActionDelete).First(&auditLog).Error
	assert.NoError(t, err)
	assert.Equal(t, systemUser.Id, auditLog.UserId)
	assert.Equal(t, systemUser.Email, auditLog.UserEmail)
	assert.True(t, auditLog.Success)
	assert.NotNil(t, auditLog.AdditionalInfo)
}
