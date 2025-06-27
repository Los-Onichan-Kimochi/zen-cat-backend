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
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestCheckUserExistsUserFound(t *testing.T) {
	/*
		GIVEN: A user exists in the database with a specific email
		WHEN:  GET /user/exists?email={email} is called with that email
		THEN:  A HTTP_200_OK status should be returned with exists=true
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test user with a specific email
	email := utilsTest.GenerateRandomEmail()
	factories.NewUserModel(db, factories.UserModelF{
		Email: &email,
	})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/user/exists?email="+email, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.CheckUserExistsResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response
	assert.Equal(t, email, response.Email)
	assert.True(t, response.Exists)
}

func TestCheckUserExistsUserNotFound(t *testing.T) {
	/*
		GIVEN: No user exists with the specified email
		WHEN:  GET /user/exists?email={email} is called with a non-existent email
		THEN:  A HTTP_200_OK status should be returned with exists=false
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentEmail := utilsTest.GenerateRandomEmail()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/user/exists?email="+nonExistentEmail, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.CheckUserExistsResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response
	assert.Equal(t, nonExistentEmail, response.Email)
	assert.False(t, response.Exists)
}

func TestCheckUserExistsMissingEmail(t *testing.T) {
	/*
		GIVEN: No email parameter is provided
		WHEN:  GET /user/exists is called without an email parameter
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/user/exists", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCheckUserExistsInvalidEmail(t *testing.T) {
	/*
		GIVEN: An invalid email format is provided
		WHEN:  GET /user/exists?email={invalid-email} is called with an invalid email
		THEN:  A HTTP_200_OK status should be returned (treated as non-existent user)
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidEmail := "not-an-email"

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/user/exists?email="+invalidEmail, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.CheckUserExistsResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Invalid email should be treated as non-existent user
	assert.Equal(t, invalidEmail, response.Email)
	assert.False(t, response.Exists)
}
