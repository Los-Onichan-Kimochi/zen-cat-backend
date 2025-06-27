package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestBulkDeleteUsersSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple users exist in the database
		WHEN:  DELETE /user/bulk-delete/ is called with valid user IDs
		THEN:  A HTTP_204_NO_CONTENT status should be returned and the users should be deleted
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test users using factory
	numUsers := 3
	testUsers := factories.NewUserModelBatch(db, numUsers)

	// Extract user IDs
	userIds := make([]uuid.UUID, numUsers)
	for i, user := range testUsers {
		userIds[i] = user.Id
	}

	bulkDeleteRequest := schemas.BulkDeleteUserRequest{
		Users: userIds,
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/user/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the users were deleted
	for _, userId := range userIds {
		var count int64
		db.Model(&testUsers[0]).Where("id = ?", userId).Count(&count)
		assert.Equal(t, int64(0), count)
	}
}

func TestBulkDeleteUsersEmptyList(t *testing.T) {
	/*
		GIVEN: A bulk user deletion request with an empty list
		WHEN:  DELETE /user/bulk-delete/ is called with an empty users list
		THEN:  A HTTP_204_NO_CONTENT status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	bulkDeleteRequest := schemas.BulkDeleteUserRequest{
		Users: []uuid.UUID{},
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/user/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBulkDeleteUsersInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  DELETE /user/bulk-delete/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/user/bulk-delete/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
