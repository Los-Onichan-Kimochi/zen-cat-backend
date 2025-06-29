package audit_log_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestGetAuditLogsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple audit logs exist in the database
		WHEN:  GET /audit-log/ is called
		THEN:  A HTTP_200_OK status should be returned with the audit logs
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})
	auditLog1 := &model.AuditLog{
		Id:        uuid.New(),
		UserId:    user.Id,
		UserEmail: user.Email,
		UserRole:  user.Rol,
		Action:    "LOGIN",
		Success:   true,
	}
	auditLog2 := &model.AuditLog{
		Id:        uuid.New(),
		UserId:    user.Id,
		UserEmail: user.Email,
		UserRole:  user.Rol,
		Action:    "CREATE",
		Success:   true,
	}
	db.Create(auditLog1)
	db.Create(auditLog2)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/audit-log/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.AuditLogs
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, len(response.AuditLogs), 2)
}

func TestGetAuditLogsWithFiltersSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple audit logs exist in the database
		WHEN:  GET /audit-log/ is called with filters
		THEN:  A HTTP_200_OK status should be returned with the filtered audit logs
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	user1 := factories.NewUserModel(db, factories.UserModelF{})
	user2 := factories.NewUserModel(db, factories.UserModelF{})

	auditLog1 := &model.AuditLog{
		Id:         uuid.New(),
		UserId:     user1.Id,
		UserEmail:  user1.Email,
		UserRole:   user1.Rol,
		Action:     "CREATE",
		EntityType: "USER",
		Success:    true,
	}
	auditLog2 := &model.AuditLog{
		Id:        uuid.New(),
		UserId:    user2.Id,
		UserEmail: user2.Email,
		UserRole:  user2.Rol,
		Action:    "LOGIN",
		Success:   true,
	}
	db.Create(auditLog1)
	db.Create(auditLog2)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/audit-log/?userIds="+user1.Id.String()+"&actions=CREATE&entityTypes=USER", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.AuditLogs
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Len(t, response.AuditLogs, 1)
	assert.Equal(t, user1.Id, response.AuditLogs[0].UserId)
	assert.Equal(t, "CREATE", string(response.AuditLogs[0].Action))
	assert.Equal(t, "USER", string(response.AuditLogs[0].EntityType))
}

func TestGetAuditLogsEmpty(t *testing.T) {
	/*
		GIVEN: No audit logs exist in the database
		WHEN:  GET /audit-log/ is called
		THEN:  A HTTP_200_OK status should be returned with an empty array
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/audit-log/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.AuditLogs
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Empty(t, response.AuditLogs)
}
