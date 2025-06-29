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

func TestGetAuditLogByIdSuccessfully(t *testing.T) {
	/*
		GIVEN: An audit log exists in the database
		WHEN:  GET /audit-log/{auditLogId}/ is called with a valid ID
		THEN:  A HTTP_200_OK status should be returned with the audit log data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})
	auditLog := &model.AuditLog{
		Id:        uuid.New(),
		UserId:    user.Id,
		UserEmail: user.Email,
		UserRole:  user.Rol,
		Action:    "LOGIN",
		Success:   true,
	}
	db.Create(auditLog)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/audit-log/"+auditLog.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.AuditLog
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, auditLog.Id, response.Id)
	assert.Equal(t, auditLog.UserId, response.UserId)
	assert.Equal(t, string(auditLog.Action), string(response.Action))
}

func TestGetAuditLogByIdNotFound(t *testing.T) {
	/*
		GIVEN: No audit log exists with the given ID
		WHEN:  GET /audit-log/{auditLogId}/ is called with a non-existent ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/audit-log/"+nonExistentId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAuditLogByIdInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided
		WHEN:  GET /audit-log/{auditLogId}/ is called with an invalid ID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/audit-log/invalid-uuid/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
