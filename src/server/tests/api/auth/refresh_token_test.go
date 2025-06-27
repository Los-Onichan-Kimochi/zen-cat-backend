package auth_test

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

func TestRefreshTokenSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid user with a valid access token
		WHEN:  POST /auth/refresh/ is called with valid token
		THEN:  A HTTP_200_OK status should be returned with new tokens
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a user first
	user := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Password:      "$2a$10$hash", // Hashed password
		Rol:           "MEMBER",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(user).Error
	assert.NoError(t, err)

	// Login to get a valid token
	loginRequest := schemas.LoginRequest{
		Email:    user.Email,
		Password: "password123",
	}
	loginBody, _ := json.Marshal(loginRequest)

	loginReq := httptest.NewRequest(http.MethodPost, "/login/", bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	server.Echo.ServeHTTP(loginRec, loginReq)

	var loginResponse schemas.TokenResponse
	json.NewDecoder(loginRec.Body).Decode(&loginResponse)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh/", nil)
	req.Header.Set("Authorization", "Bearer "+loginResponse.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.TokenResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
}

func TestRefreshTokenMissingToken(t *testing.T) {
	/*
		GIVEN: No authorization token provided
		WHEN:  POST /auth/refresh/ is called without token
		THEN:  A HTTP_401_UNAUTHORIZED status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRefreshTokenInvalidToken(t *testing.T) {
	/*
		GIVEN: An invalid authorization token
		WHEN:  POST /auth/refresh/ is called with invalid token
		THEN:  A HTTP_401_UNAUTHORIZED status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh/", nil)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRefreshTokenMalformedToken(t *testing.T) {
	/*
		GIVEN: A malformed authorization token
		WHEN:  POST /auth/refresh/ is called with malformed token
		THEN:  A HTTP_401_UNAUTHORIZED status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh/", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.jwt.format")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
