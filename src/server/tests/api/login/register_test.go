package login_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestRegisterSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid user registration data
		WHEN:  POST /register/ is called with valid user data
		THEN:  A HTTP_201_CREATED status should be returned with user data and token
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	registerRequest := schemas.RegisterRequest{
		Email:          utilsTest.GenerateRandomEmail(),
		Password:       "testPassword123",
		Name:           "John",
		FirstLastName:  "Doe",
		SecondLastName: stringPtr("Smith"),
	}

	requestBody, _ := json.Marshal(registerRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/register/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.LoginResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains token and user data
	assert.NotEmpty(t, response.Tokens.AccessToken)
	assert.NotNil(t, response.User)
	assert.Equal(t, registerRequest.Email, response.User.Email)
	assert.Equal(t, registerRequest.Name, response.User.Name)
	assert.Equal(t, registerRequest.FirstLastName, response.User.FirstLastName)
}

func TestRegisterDuplicateEmail(t *testing.T) {
	/*
		GIVEN: A user already exists with the same email
		WHEN:  POST /register/ is called with duplicate email
		THEN:  A HTTP_409_CONFLICT status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create an existing user
	existingEmail := utilsTest.GenerateRandomEmail()
	existingUser := &model.User{
		Email:         existingEmail,
		Name:          "Existing",
		FirstLastName: "User",
		Rol:           "MEMBER",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(existingUser).Error
	assert.NoError(t, err)

	registerRequest := schemas.RegisterRequest{
		Email:         existingEmail, // Same email as existing user
		Password:      "testPassword123",
		Name:          "John",
		FirstLastName: "Doe",
	}

	requestBody, _ := json.Marshal(registerRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/register/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestRegisterInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: Invalid request body
		WHEN:  POST /register/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/register/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestRegisterMissingRequiredFields(t *testing.T) {
	/*
		GIVEN: Registration request missing required fields
		WHEN:  POST /register/ is called with incomplete data
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	registerRequest := schemas.RegisterRequest{
		Email: utilsTest.GenerateRandomEmail(),
		// Missing Password, Name, FirstLastName
	}

	requestBody, _ := json.Marshal(registerRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/register/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

// Helper function for creating string pointers
func stringPtr(s string) *string {
	return &s
}
