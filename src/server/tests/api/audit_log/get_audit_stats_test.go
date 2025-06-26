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

func TestGetAuditStatsSuccessfully(t *testing.T) {
	/*
		GIVEN: An audit log exists in the database
		WHEN:  GET /audit-log/stats is called
		THEN:  A HTTP_200_OK status should be returned with the audit log stats
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
	req := httptest.NewRequest(http.MethodGet, "/audit-log/stats/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.AuditStats
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, response.TotalEvents, int64(1))
}

func TestGetAuditStatsEmpty(t *testing.T) {
	/*
		GIVEN: No audit logs exist in the database
		WHEN:  GET /audit-log/stats is called
		THEN:  A HTTP_200_OK status should be returned with empty stats
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/audit-log/stats/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.AuditStats
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, int64(0), response.TotalEvents)
} 