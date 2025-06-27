package user_test

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

func TestChangePasswordSuccessfully(t *testing.T) {
	/*
		GIVEN: A user exists in the database
		WHEN:  POST /user/change-password/ is called with valid email and new password
		THEN:  A HTTP_200_OK status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test user
	email := utilsTest.GenerateRandomEmail()
	factories.NewUserModel(db, factories.UserModelF{
		Email: &email,
	})

	// Prepare change password request
	changePasswordRequest := schemas.ChangePasswordInput{
		Email:       email,
		NewPassword: "newSecurePassword123",
	}

	requestBody, _ := json.Marshal(changePasswordRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/user/change-password/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestChangePasswordUserNotFound(t *testing.T) {
	/*
		GIVEN: No user exists with the provided email
		WHEN:  POST /user/change-password/ is called with a non-existent email
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentEmail := utilsTest.GenerateRandomEmail()

	// Prepare change password request
	changePasswordRequest := schemas.ChangePasswordInput{
		Email:       nonExistentEmail,
		NewPassword: "newSecurePassword123",
	}

	requestBody, _ := json.Marshal(changePasswordRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/user/change-password/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestChangePasswordInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /user/change-password/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/user/change-password/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestChangePasswordMissingRequiredFields(t *testing.T) {
	/*
		GIVEN: A request body missing required fields
		WHEN:  POST /user/change-password/ is called with missing required fields
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned (since it's a validation error)
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// Missing email
	missingEmailRequest := schemas.ChangePasswordInput{
		NewPassword: "newSecurePassword123",
	}

	requestBody, _ := json.Marshal(missingEmailRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/user/change-password/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	// Missing password
	missingPasswordRequest := schemas.ChangePasswordInput{
		Email: utilsTest.GenerateRandomEmail(),
	}

	requestBody, _ = json.Marshal(missingPasswordRequest)

	// WHEN
	req = httptest.NewRequest(http.MethodPost, "/user/change-password/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec = httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
