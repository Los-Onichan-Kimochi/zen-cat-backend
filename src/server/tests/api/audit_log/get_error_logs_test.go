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

func TestGetErrorLogsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple error logs exist in the database
		WHEN:  GET /error-log/ is called
		THEN:  A HTTP_200_OK status should be returned with the error logs
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})
	errorLog1 := &model.AuditLog{
		Id:           uuid.New(),
		UserId:       user.Id,
		UserEmail:    user.Email,
		UserRole:     user.Rol,
		Action:       "LOGIN",
		Success:      false,
		ErrorMessage: strPtr("Invalid credentials"),
	}
	errorLog2 := &model.AuditLog{
		Id:           uuid.New(),
		UserId:       user.Id,
		UserEmail:    user.Email,
		UserRole:     user.Rol,
		Action:       "CREATE",
		Success:      false,
		ErrorMessage: strPtr("Validation failed"),
	}
	db.Create(errorLog1)
	db.Create(errorLog2)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/error-log/", nil)
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
	for _, log := range response.AuditLogs {
		assert.False(t, log.Success)
	}
}

func strPtr(s string) *string {
	return &s
}
