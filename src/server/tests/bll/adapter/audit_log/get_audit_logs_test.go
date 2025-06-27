package audit_log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestGetAuditLogsSuccessfully(t *testing.T) {
	/*
		GIVEN: Audit events are logged
		WHEN:  GetAuditLogs is called
		THEN:  A list of audit logs is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewAuditLogAdapterTestWrapper(t)

	// Create test users for audit logs
	user1 := factories.NewUserModel(db, factories.UserModelF{})
	user2 := factories.NewUserModel(db, factories.UserModelF{})

	userAgent1 := "TestAgent/1.0"
	userAgent2 := "TestAgent/2.0"
	entityName1 := "Test User 1"
	entityName2 := "Test User 2"

	// Log audit events
	context1 := schemas.AuditContext{
		UserId:    user1.Id,
		UserEmail: user1.Email,
		UserRole:  schemas.UserRol(user1.Rol),
		IPAddress: "192.168.1.1",
		UserAgent: &userAgent1,
	}
	event1 := schemas.AuditEvent{
		Action:     schemas.AuditActionCreate,
		EntityType: schemas.AuditEntityUser,
		EntityId:   &user1.Id,
		EntityName: &entityName1,
		Success:    true,
	}
	err1 := adapter.LogAuditEvent(context1, event1)
	assert.Nil(t, err1)

	context2 := schemas.AuditContext{
		UserId:    user2.Id,
		UserEmail: user2.Email,
		UserRole:  schemas.UserRol(user2.Rol),
		IPAddress: "192.168.1.2",
		UserAgent: &userAgent2,
	}
	event2 := schemas.AuditEvent{
		Action:     schemas.AuditActionUpdate,
		EntityType: schemas.AuditEntityUser,
		EntityId:   &user2.Id,
		EntityName: &entityName2,
		Success:    true,
	}
	err2 := adapter.LogAuditEvent(context2, event2)
	assert.Nil(t, err2)

	// WHEN
	filters := schemas.AuditLogFilters{
		Page:     1,
		PageSize: 10,
	}
	auditLogs, err := adapter.GetAuditLogs(filters)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, auditLogs)
	assert.GreaterOrEqual(t, len(auditLogs.AuditLogs), 2)
	assert.Greater(t, auditLogs.TotalCount, int64(0))
	assert.Equal(t, 1, auditLogs.Page)
	assert.Equal(t, 10, auditLogs.PageSize)
}

func TestGetAuditLogsWithFilters(t *testing.T) {
	/*
		GIVEN: Audit events with different actions are logged
		WHEN:  GetAuditLogs is called with action filter
		THEN:  Only matching audit logs are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewAuditLogAdapterTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})

	userAgent := "TestAgent/1.0"
	entityName := "Test User"

	// Log CREATE event
	context := schemas.AuditContext{
		UserId:    user.Id,
		UserEmail: user.Email,
		UserRole:  schemas.UserRol(user.Rol),
		IPAddress: "192.168.1.1",
		UserAgent: &userAgent,
	}
	createEvent := schemas.AuditEvent{
		Action:     schemas.AuditActionCreate,
		EntityType: schemas.AuditEntityUser,
		EntityId:   &user.Id,
		EntityName: &entityName,
		Success:    true,
	}
	updateEvent := schemas.AuditEvent{
		Action:     schemas.AuditActionUpdate,
		EntityType: schemas.AuditEntityUser,
		EntityId:   &user.Id,
		EntityName: &entityName,
		Success:    true,
	}

	err1 := adapter.LogAuditEvent(context, createEvent)
	assert.Nil(t, err1)
	err2 := adapter.LogAuditEvent(context, updateEvent)
	assert.Nil(t, err2)

	// WHEN - Filter by CREATE action
	filters := schemas.AuditLogFilters{
		Actions:  []string{string(schemas.AuditActionCreate)},
		Page:     1,
		PageSize: 10,
	}
	auditLogs, err := adapter.GetAuditLogs(filters)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, auditLogs)
	assert.GreaterOrEqual(t, len(auditLogs.AuditLogs), 1)

	// Verify all returned logs have CREATE action
	for _, log := range auditLogs.AuditLogs {
		assert.Equal(t, schemas.AuditActionCreate, log.Action)
	}
}

func TestGetAuditLogsWithUserIdFilter(t *testing.T) {
	/*
		GIVEN: Audit events for different users are logged
		WHEN:  GetAuditLogs is called with userId filter
		THEN:  Only audit logs for that user are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewAuditLogAdapterTestWrapper(t)

	user1 := factories.NewUserModel(db, factories.UserModelF{})
	user2 := factories.NewUserModel(db, factories.UserModelF{})

	userAgent1 := "TestAgent/1.0"
	userAgent2 := "TestAgent/2.0"
	entityName1 := "User 1"
	entityName2 := "User 2"

	// Log events for different users
	context1 := schemas.AuditContext{
		UserId:    user1.Id,
		UserEmail: user1.Email,
		UserRole:  schemas.UserRol(user1.Rol),
		IPAddress: "192.168.1.1",
		UserAgent: &userAgent1,
	}
	context2 := schemas.AuditContext{
		UserId:    user2.Id,
		UserEmail: user2.Email,
		UserRole:  schemas.UserRol(user2.Rol),
		IPAddress: "192.168.1.2",
		UserAgent: &userAgent2,
	}

	event1 := schemas.AuditEvent{
		Action:     schemas.AuditActionCreate,
		EntityType: schemas.AuditEntityUser,
		EntityId:   &user1.Id,
		EntityName: &entityName1,
		Success:    true,
	}
	event2 := schemas.AuditEvent{
		Action:     schemas.AuditActionCreate,
		EntityType: schemas.AuditEntityUser,
		EntityId:   &user2.Id,
		EntityName: &entityName2,
		Success:    true,
	}

	err1 := adapter.LogAuditEvent(context1, event1)
	assert.Nil(t, err1)
	err2 := adapter.LogAuditEvent(context2, event2)
	assert.Nil(t, err2)

	// WHEN - Filter by user1 ID
	filters := schemas.AuditLogFilters{
		UserIds:  []string{user1.Id.String()},
		Page:     1,
		PageSize: 10,
	}
	auditLogs, err := adapter.GetAuditLogs(filters)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, auditLogs)
	assert.GreaterOrEqual(t, len(auditLogs.AuditLogs), 1)

	// Verify all returned logs are for user1
	for _, log := range auditLogs.AuditLogs {
		assert.Equal(t, user1.Id, log.UserId)
	}
}

func TestGetAuditLogsEmpty(t *testing.T) {
	/*
		GIVEN: No audit logs exist in the database
		WHEN:  GetAuditLogs is called
		THEN:  An empty list is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewAuditLogAdapterTestWrapper(t)

	// WHEN
	filters := schemas.AuditLogFilters{
		Page:     1,
		PageSize: 10,
	}
	auditLogs, err := adapter.GetAuditLogs(filters)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, auditLogs)
	assert.Equal(t, 0, len(auditLogs.AuditLogs))
	assert.Equal(t, int64(0), auditLogs.TotalCount)
}

func TestLogAuditEventSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid audit context and event
		WHEN:  LogAuditEvent is called
		THEN:  The event is logged successfully
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewAuditLogAdapterTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})

	userAgent := "TestAgent/1.0"
	entityName := "Test User"
	additionalInfo := `{"source": "test"}`

	context := schemas.AuditContext{
		UserId:    user.Id,
		UserEmail: user.Email,
		UserRole:  schemas.UserRol(user.Rol),
		IPAddress: "192.168.1.1",
		UserAgent: &userAgent,
	}
	event := schemas.AuditEvent{
		Action:         schemas.AuditActionCreate,
		EntityType:     schemas.AuditEntityUser,
		EntityId:       &user.Id,
		EntityName:     &entityName,
		OldValues:      nil,
		NewValues:      map[string]interface{}{"name": "Test User"},
		AdditionalInfo: &additionalInfo,
		Success:        true,
	}

	// WHEN
	err := adapter.LogAuditEvent(context, event)

	// THEN
	assert.Nil(t, err)

	// Verify the event was logged
	filters := schemas.AuditLogFilters{
		UserIds:  []string{user.Id.String()},
		Page:     1,
		PageSize: 10,
	}
	auditLogs, getErr := adapter.GetAuditLogs(filters)
	assert.Nil(t, getErr)
	assert.GreaterOrEqual(t, len(auditLogs.AuditLogs), 1)
}
