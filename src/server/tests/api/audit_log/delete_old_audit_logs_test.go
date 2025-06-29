package audit_log_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestDeleteOldAuditLogsSuccessfully(t *testing.T) {
	/*
		GIVEN: An old audit log exists in the database
		WHEN:  DELETE /audit-log/cleanup/ is called
		THEN:  The old audit log should be deleted and a HTTP_204_NO_CONTENT status returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})
	oldAuditLog := &model.AuditLog{
		Id:        uuid.New(),
		UserId:    user.Id,
		UserEmail: user.Email,
		UserRole:  user.Rol,
		Action:    "LOGIN",
		Success:   true,
	}
	// Manually set CreatedAt to be in the past
	oldAuditLog.CreatedAt = time.Now().AddDate(0, 0, -100)
	db.Create(oldAuditLog)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/audit-log/cleanup/?days=90", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the log was deleted
	var count int64
	db.Model(&model.AuditLog{}).Where("id = ?", oldAuditLog.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteOldAuditLogsNoDeletion(t *testing.T) {
	/*
		GIVEN: A recent audit log exists in the database
		WHEN:  DELETE /audit-log/cleanup/ is called
		THEN:  The log should not be deleted and a HTTP_204_NO_CONTENT status returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})
	recentAuditLog := &model.AuditLog{
		Id:        uuid.New(),
		UserId:    user.Id,
		UserEmail: user.Email,
		UserRole:  user.Rol,
		Action:    "LOGIN",
		Success:   true,
	}
	db.Create(recentAuditLog)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/audit-log/cleanup/?days=90", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the log was not deleted
	var count int64
	db.Model(&model.AuditLog{}).Where("id = ?", recentAuditLog.Id).Count(&count)
	assert.Equal(t, int64(1), count)
}
