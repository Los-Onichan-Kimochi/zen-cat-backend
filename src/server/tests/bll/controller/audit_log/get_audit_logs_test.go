package audit_log_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestGetAuditLogsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple audit logs exist in the database
		WHEN:  GetAuditLogs is called without filters
		THEN:  All audit logs should be returned with pagination
	*/
	// GIVEN
	auditLogController, _, db := controllerTest.NewAuditLogControllerTestWrapper(t)

	// Create user using factory
	user := factories.NewUserModel(db)

	// Create audit logs manually (no factory available)
	createAction := model.AuditActionCreate
	updateAction := model.AuditActionUpdate
	entityType := model.AuditEntityUser
	success := true

	auditLogs := []*model.AuditLog{
		{
			Id:         uuid.New(),
			UserId:     user.Id,
			UserEmail:  user.Email,
			UserRole:   user.Rol,
			Action:     createAction,
			EntityType: entityType,
			EntityId:   &user.Id,
			Success:    success,
			IPAddress:  "127.0.0.1",
		},
		{
			Id:         uuid.New(),
			UserId:     user.Id,
			UserEmail:  user.Email,
			UserRole:   user.Rol,
			Action:     updateAction,
			EntityType: entityType,
			EntityId:   &user.Id,
			Success:    success,
			IPAddress:  "127.0.0.1",
		},
	}
	err := db.Create(auditLogs).Error
	assert.NoError(t, err)

	// WHEN
	result, errResult := auditLogController.GetAuditLogs(
		"",   // userIds
		"",   // actions
		"",   // entityTypes
		"",   // userRoles
		"",   // startDate
		"",   // endDate
		"",   // success
		"1",  // page
		"10", // pageSize
	)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result.AuditLogs), 2)
	assert.GreaterOrEqual(t, result.TotalCount, int64(2))
}

func TestGetAuditLogsWithUserFilter(t *testing.T) {
	/*
		GIVEN: Multiple audit logs for different users exist
		WHEN:  GetAuditLogs is called with user ID filter
		THEN:  Only audit logs for the specified user should be returned
	*/
	// GIVEN
	auditLogController, _, db := controllerTest.NewAuditLogControllerTestWrapper(t)

	// Create users using factories
	user1 := factories.NewUserModel(db)
	user2 := factories.NewUserModel(db)

	// Create audit logs for both users manually
	createAction := model.AuditActionCreate
	entityType := model.AuditEntityUser
	success := true

	auditLogs := []*model.AuditLog{
		{
			Id:         uuid.New(),
			UserId:     user1.Id,
			UserEmail:  user1.Email,
			UserRole:   user1.Rol,
			Action:     createAction,
			EntityType: entityType,
			EntityId:   &user1.Id,
			Success:    success,
			IPAddress:  "127.0.0.1",
		},
		{
			Id:         uuid.New(),
			UserId:     user2.Id,
			UserEmail:  user2.Email,
			UserRole:   user2.Rol,
			Action:     createAction,
			EntityType: entityType,
			EntityId:   &user2.Id,
			Success:    success,
			IPAddress:  "127.0.0.1",
		},
	}
	err := db.Create(auditLogs).Error
	assert.NoError(t, err)

	// WHEN
	result, errResult := auditLogController.GetAuditLogs(
		user1.Id.String(), // userIds filter
		"",                // actions
		"",                // entityTypes
		"",                // userRoles
		"",                // startDate
		"",                // endDate
		"",                // success
		"1",               // page
		"10",              // pageSize
	)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.AuditLogs))
	if len(result.AuditLogs) > 0 {
		assert.Equal(t, user1.Id, result.AuditLogs[0].UserId)
	}
}

func TestGetAuditLogsWithActionFilter(t *testing.T) {
	/*
		GIVEN: Multiple audit logs with different actions exist
		WHEN:  GetAuditLogs is called with action filter
		THEN:  Only audit logs matching the action should be returned
	*/
	// GIVEN
	auditLogController, _, db := controllerTest.NewAuditLogControllerTestWrapper(t)

	// Create user using factory
	user := factories.NewUserModel(db)

	// Create audit logs with different actions
	createAction := model.AuditActionCreate
	updateAction := model.AuditActionUpdate
	entityType := model.AuditEntityUser
	success := true

	auditLogs := []*model.AuditLog{
		{
			Id:         uuid.New(),
			UserId:     user.Id,
			UserEmail:  user.Email,
			UserRole:   user.Rol,
			Action:     createAction,
			EntityType: entityType,
			Success:    success,
			IPAddress:  "127.0.0.1",
		},
		{
			Id:         uuid.New(),
			UserId:     user.Id,
			UserEmail:  user.Email,
			UserRole:   user.Rol,
			Action:     updateAction,
			EntityType: entityType,
			Success:    success,
			IPAddress:  "127.0.0.1",
		},
	}
	err := db.Create(auditLogs).Error
	assert.NoError(t, err)

	// WHEN
	result, errResult := auditLogController.GetAuditLogs(
		"",       // userIds
		"CREATE", // actions filter
		"",       // entityTypes
		"",       // userRoles
		"",       // startDate
		"",       // endDate
		"",       // success
		"1",      // page
		"10",     // pageSize
	)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.AuditLogs))
	if len(result.AuditLogs) > 0 {
		assert.Equal(t, string(model.AuditActionCreate), string(result.AuditLogs[0].Action))
	}
}
