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

func TestGetErrorStatsSuccessfully(t *testing.T) {
	/*
		GIVEN: An error log exists in the database
		WHEN:  GET /error-log/stats/ is called
		THEN:  A HTTP_200_OK status should be returned with the error log stats
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})
	errorLog := &model.AuditLog{
		Id:           uuid.New(),
		UserId:       user.Id,
		UserEmail:    user.Email,
		UserRole:     user.Rol,
		Action:       "LOGIN",
		Success:      false,
		ErrorMessage: strPtr("Invalid credentials"),
	}
	db.Create(errorLog)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/error-log/stats/", nil)
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
	assert.GreaterOrEqual(t, response.FailureCount, int64(1))
}

func TestGetErrorStatsEmpty(t *testing.T) {
	/*
		GIVEN: No error logs exist in the database
		WHEN:  GET /error-log/stats/ is called
		THEN:  A HTTP_200_OK status should be returned with empty stats
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/error-log/stats/", nil)
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
	assert.Equal(t, int64(0), response.FailureCount)
}
