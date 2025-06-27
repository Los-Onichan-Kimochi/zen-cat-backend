package forgot_password_test

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

func TestForgotPasswordSuccessfully(t *testing.T) {
	/*
		GIVEN: A user exists with a valid email
		WHEN:  POST /forgot-password/ is called with that email
		THEN:  A HTTP_200_OK status should be returned with a success message and PIN
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test user using factory
	email := utilsTest.GenerateRandomEmail()
	factories.NewUserModel(db, factories.UserModelF{
		Email: &email,
	})

	forgotPasswordRequest := schemas.ForgotPasswordRequest{
		Email: email,
	}

	requestBody, _ := json.Marshal(forgotPasswordRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/forgot-password/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ForgotPasswordResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains a message and PIN
	assert.NotEmpty(t, response.Message)
	assert.NotEmpty(t, response.Pin)
}

func TestForgotPasswordUserNotFound(t *testing.T) {
	/*
		GIVEN: No user exists with the provided email
		WHEN:  POST /forgot-password/ is called with a non-existent email
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentEmail := utilsTest.GenerateRandomEmail()

	forgotPasswordRequest := schemas.ForgotPasswordRequest{
		Email: nonExistentEmail,
	}

	requestBody, _ := json.Marshal(forgotPasswordRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/forgot-password/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestForgotPasswordInvalidEmail(t *testing.T) {
	/*
		GIVEN: An invalid email format is provided
		WHEN:  POST /forgot-password/ is called with an invalid email
		THEN:  A HTTP_400_BAD_REQUEST status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidEmail := "not-an-email"

	forgotPasswordRequest := schemas.ForgotPasswordRequest{
		Email: invalidEmail,
	}

	requestBody, _ := json.Marshal(forgotPasswordRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/forgot-password/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestForgotPasswordInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /forgot-password/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/forgot-password/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
