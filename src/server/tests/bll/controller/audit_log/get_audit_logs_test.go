package audit_log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestGetAuditLogsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple audit logs exist in the database
		WHEN:  GetAuditLogs is called with no filters
		THEN:  A list of audit logs should be returned with proper pagination
	*/
	// GIVEN
	auditLogController, _, db := controllerTest.NewAuditLogControllerTestWrapper(t)

	// Create a user for the audit logs
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

	// Create multiple audit logs
	auditLogs := []*model.AuditLog{
		{
			UserId:     user.Id,
			UserEmail:  user.Email,
			UserRole:   user.Rol,
			Action:     model.AuditActionCreate,
			EntityType: model.AuditEntityUser,
			EntityId:   &user.Id,
			Success:    true,
			IPAddress:  "127.0.0.1",
		},
		{
			UserId:     user.Id,
			UserEmail:  user.Email,
			UserRole:   user.Rol,
			Action:     model.AuditActionUpdate,
			EntityType: model.AuditEntityUser,
			EntityId:   &user.Id,
			Success:    true,
			IPAddress:  "127.0.0.1",
		},
	}
	err = db.Create(auditLogs).Error
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

	// Create two users
	user1 := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           model.UserRolAdmin,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	user2 := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "Jane",
		FirstLastName: "Smith",
		Rol:           model.UserRolClient,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create([]*model.User{user1, user2}).Error
	assert.NoError(t, err)

	// Create audit logs for both users
	auditLogs := []*model.AuditLog{
		{
			UserId:     user1.Id,
			UserEmail:  user1.Email,
			UserRole:   user1.Rol,
			Action:     model.AuditActionCreate,
			EntityType: model.AuditEntityUser,
			EntityId:   &user1.Id,
			Success:    true,
			IPAddress:  "127.0.0.1",
		},
		{
			UserId:     user2.Id,
			UserEmail:  user2.Email,
			UserRole:   user2.Rol,
			Action:     model.AuditActionCreate,
			EntityType: model.AuditEntityUser,
			EntityId:   &user2.Id,
			Success:    true,
			IPAddress:  "127.0.0.1",
		},
	}
	err = db.Create(auditLogs).Error
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
	assert.Equal(t, user1.Id, result.AuditLogs[0].UserId)
}

func TestGetAuditLogsWithActionFilter(t *testing.T) {
	/*
		GIVEN: Multiple audit logs with different actions exist
		WHEN:  GetAuditLogs is called with action filter
		THEN:  Only audit logs matching the action should be returned
	*/
	// GIVEN
	auditLogController, _, db := controllerTest.NewAuditLogControllerTestWrapper(t)

	// Create a user for the audit logs
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

	// Create audit logs with different actions
	auditLogs := []*model.AuditLog{
		{
			UserId:     user.Id,
			UserEmail:  user.Email,
			UserRole:   user.Rol,
			Action:     model.AuditActionCreate,
			EntityType: model.AuditEntityUser,
			Success:    true,
			IPAddress:  "127.0.0.1",
		},
		{
			UserId:     user.Id,
			UserEmail:  user.Email,
			UserRole:   user.Rol,
			Action:     model.AuditActionUpdate,
			EntityType: model.AuditEntityUser,
			Success:    true,
			IPAddress:  "127.0.0.1",
		},
	}
	err = db.Create(auditLogs).Error
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
	assert.Equal(t, model.AuditActionCreate, result.AuditLogs[0].Action)
}
