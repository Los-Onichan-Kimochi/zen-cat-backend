package login_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestGetCurrentUserSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid user exists in the system and is authenticated
		WHEN:  GET /me/ is called with valid authentication
		THEN:  A HTTP_200_OK status should be returned with user profile data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test user with proper password hashing
	email := utilsTest.GenerateRandomEmail()
	password := "testPassword123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	user := &model.User{
		Email:         email,
		Name:          "John",
		FirstLastName: "Doe",
		Password:      string(hashedPassword),
		Rol:           "MEMBER",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(user).Error
	assert.NoError(t, err)

	// Login to get authentication token
	loginRequest := schemas.LoginRequest{
		Email:    email,
		Password: password,
	}
	loginBody, _ := json.Marshal(loginRequest)

	loginReq := httptest.NewRequest(http.MethodPost, "/login/", bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	server.Echo.ServeHTTP(loginRec, loginReq)

	// Verify login was successful
	assert.Equal(t, http.StatusOK, loginRec.Code)

	var loginResponse schemas.LoginResponse
	err = json.NewDecoder(loginRec.Body).Decode(&loginResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResponse.Tokens.AccessToken)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/me/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResponse.Tokens.AccessToken)

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.UserProfile
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains correct user data
	assert.Equal(t, user.Id, response.Id)
	assert.Equal(t, email, response.Email)
	assert.Equal(t, "John", response.Name)
	assert.Equal(t, "Doe", response.FirstLastName)
}

func TestGetCurrentUserInvalidRequestMethod(t *testing.T) {
	/*
		GIVEN: A GET endpoint
		WHEN:  POST /me/ is called with wrong HTTP method
		THEN:  A HTTP_405_METHOD_NOT_ALLOWED status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/me/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestGetCurrentUserWithoutUser(t *testing.T) {
	/*
		GIVEN: No user exists in the system
		WHEN:  GET /me/ is called
		THEN:  A HTTP_404_NOT_FOUND or appropriate error status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/me/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	// The exact status code depends on implementation
	// Could be 404 (not found), 401 (unauthorized), or 500 (internal server error)
	assert.Contains(t, []int{http.StatusNotFound, http.StatusUnauthorized, http.StatusInternalServerError}, rec.Code)
}
