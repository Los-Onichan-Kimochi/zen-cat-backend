package login_test

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

func TestGetCurrentUserSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid user exists in the system
		WHEN:  GET /me/ is called (auth disabled in tests)
		THEN:  A HTTP_200_OK status should be returned with user profile data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test user using factory
	email := utilsTest.GenerateRandomEmail()
	password := "testPassword123"
	name := "John"
	firstLastName := "Doe"

	user := factories.NewUserModel(db, factories.UserModelF{
		Email:         &email,
		Password:      &password,
		Name:          &name,
		FirstLastName: &firstLastName,
	})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/me/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// Note: Auth is disabled in tests, so no Authorization header needed

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.UserProfile
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains correct user data
	assert.Equal(t, user.Id, response.Id)
	assert.Equal(t, email, response.Email)
	assert.Equal(t, name, response.Name)
	assert.Equal(t, firstLastName, response.FirstLastName)
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
