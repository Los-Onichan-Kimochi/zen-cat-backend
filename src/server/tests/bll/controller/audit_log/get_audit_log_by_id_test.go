package audit_log_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestGetAuditLogByIdSuccessfully(t *testing.T) {
	/*
		GIVEN: An audit log exists in the database
		WHEN:  GetAuditLogById is called with valid ID
		THEN:  The audit log should be returned
	*/
	// GIVEN
	auditLogController, _, db := controllerTest.NewAuditLogControllerTestWrapper(t)

	// Create a user for the audit log
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

	// Create an audit log
	auditLog := &model.AuditLog{
		UserId:     user.Id,
		UserEmail:  user.Email,
		UserRole:   user.Rol,
		Action:     model.AuditActionCreate,
		EntityType: model.AuditEntityUser,
		EntityId:   &user.Id,
		Success:    true,
		IPAddress:  "127.0.0.1",
	}
	err = db.Create(auditLog).Error
	assert.NoError(t, err)

	// WHEN
	result, errResult := auditLogController.GetAuditLogById(auditLog.Id)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, auditLog.Id, result.Id)
	assert.Equal(t, auditLog.UserId, result.UserId)
	assert.Equal(t, auditLog.UserEmail, result.UserEmail)
	assert.Equal(t, string(auditLog.Action), string(result.Action))
	assert.Equal(t, string(auditLog.EntityType), string(result.EntityType))
	assert.True(t, result.Success)
}

func TestGetAuditLogByIdNotFound(t *testing.T) {
	/*
		GIVEN: No audit log exists with the given ID
		WHEN:  GetAuditLogById is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	auditLogController, _, _ := controllerTest.NewAuditLogControllerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	result, errResult := auditLogController.GetAuditLogById(nonExistentId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
	assert.Contains(t, errResult.Message, "not found")
}

func TestGetAuditLogByIdWithNilId(t *testing.T) {
	/*
		GIVEN: A nil UUID
		WHEN:  GetAuditLogById is called with nil UUID
		THEN:  An error should be returned
	*/
	// GIVEN
	auditLogController, _, _ := controllerTest.NewAuditLogControllerTestWrapper(t)
	nilId := uuid.Nil

	// WHEN
	result, errResult := auditLogController.GetAuditLogById(nilId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
}
