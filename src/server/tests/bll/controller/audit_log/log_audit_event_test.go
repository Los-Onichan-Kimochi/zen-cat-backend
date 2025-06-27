package audit_log_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestLogAuditEventSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid audit context and event data
		WHEN:  LogAuditEvent is called with valid parameters
		THEN:  The audit event should be logged successfully
	*/
	// GIVEN
	auditLogController, _, db := controllerTest.NewAuditLogControllerTestWrapper(t)

	// Create a user for the audit event
	user := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           model.UserRolAdmin,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(user).Error
	assert.NoError(t, err)

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
	err = db.Where("user_id = ? AND action = ?", user.Id, schemas.AuditActionCreate).First(&auditLog).Error
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

	// Create a user for the audit event
	user := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           model.UserRolClient,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(user).Error
	assert.NoError(t, err)

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
	err = db.Where("user_id = ? AND action = ?", user.Id, schemas.AuditActionCreate).First(&auditLog).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Id, auditLog.UserId)
	assert.False(t, auditLog.Success)
	assert.NotNil(t, auditLog.ErrorMessage)
	assert.Equal(t, errorMessage, *auditLog.ErrorMessage)
}

func TestLogAuditEventWithSystemAction(t *testing.T) {
	/*
		GIVEN: System-level audit event (no specific user)
		WHEN:  LogAuditEvent is called with system context
		THEN:  The audit event should be logged successfully
	*/
	// GIVEN
	auditLogController, _, db := controllerTest.NewAuditLogControllerTestWrapper(t)

	// Create system audit context (using a system user ID)
	systemUserId := uuid.New()
	auditContext := schemas.AuditContext{
		UserId:    systemUserId,
		UserEmail: "system@example.com",
		UserRole:  schemas.UserRolAdmin,
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
	err := db.Where("user_id = ? AND action = ?", systemUserId, schemas.AuditActionDelete).First(&auditLog).Error
	assert.NoError(t, err)
	assert.Equal(t, systemUserId, auditLog.UserId)
	assert.Equal(t, "system@example.com", auditLog.UserEmail)
	assert.True(t, auditLog.Success)
	assert.NotNil(t, auditLog.AdditionalInfo)
}
