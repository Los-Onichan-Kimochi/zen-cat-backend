package login_test

import (
	"bytes"
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

func TestLoginSuccessfully(t *testing.T) {
	/*
		GIVEN: A user exists with valid credentials
		WHEN:  POST /login/ is called with correct email and password
		THEN:  A HTTP_200_OK status should be returned with authentication token
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test user using factory
	email := utilsTest.GenerateRandomEmail()
	password := "testPassword123"

	factories.NewUserModel(db, factories.UserModelF{
		Email:    &email,
		Password: &password, // Pass plain text password, factory will hash it
	})

	loginRequest := schemas.LoginRequest{
		Email:    email,
		Password: password,
	}

	requestBody, _ := json.Marshal(loginRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/login/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.LoginResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains token and user data
	assert.NotEmpty(t, response.Tokens.AccessToken)
	assert.NotNil(t, response.User)
	assert.Equal(t, email, response.User.Email)
	// Password is not included in UserProfile
}

func TestLoginInvalidCredentials(t *testing.T) {
	/*
		GIVEN: A user exists with valid credentials
		WHEN:  POST /login/ is called with incorrect password
		THEN:  A HTTP_401_UNAUTHORIZED status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test user using factory
	email := utilsTest.GenerateRandomEmail()
	password := "testPassword123"

	factories.NewUserModel(db, factories.UserModelF{
		Email:    &email,
		Password: &password, // Pass plain text password, factory will hash it
	})

	loginRequest := schemas.LoginRequest{
		Email:    email,
		Password: "wrongPassword", // Wrong password
	}

	requestBody, _ := json.Marshal(loginRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/login/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginUserNotFound(t *testing.T) {
	/*
		GIVEN: No user exists with the provided email
		WHEN:  POST /login/ is called with non-existent email
		THEN:  A HTTP_401_UNAUTHORIZED status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	loginRequest := schemas.LoginRequest{
		Email:    utilsTest.GenerateRandomEmail(), // Non-existent email
		Password: "anyPassword",
	}

	requestBody, _ := json.Marshal(loginRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/login/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /login/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/login/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
