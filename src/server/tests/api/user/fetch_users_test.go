package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestFetchUsersSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple users exist in the database
		WHEN:  GET /user/ is called
		THEN:  A HTTP_200_OK status should be returned with all users
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test users using factory
	numUsers := 3
	testUsers := factories.NewUserModelBatch(db, numUsers)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/user/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Users
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains all users
	assert.GreaterOrEqual(t, len(response.Users), numUsers)

	// Create a map of user IDs for easier verification
	userMap := make(map[string]bool)
	for _, user := range testUsers {
		userMap[user.Id.String()] = false
	}

	// Check that all created users are in the response
	for _, user := range response.Users {
		if _, exists := userMap[user.Id.String()]; exists {
			userMap[user.Id.String()] = true
		}
	}

	// Verify all test users were found in the response
	for id, found := range userMap {
		assert.True(t, found, "User with ID %s was not found in the response", id)
	}
}

func TestFetchUsersEmpty(t *testing.T) {
	/*
		GIVEN: No users exist in the database
		WHEN:  GET /user/ is called
		THEN:  A HTTP_200_OK status should be returned with an empty array
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/user/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Users
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains an empty array
	assert.Empty(t, response.Users)
}
