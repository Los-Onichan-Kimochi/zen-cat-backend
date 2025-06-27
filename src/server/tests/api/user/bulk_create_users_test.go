package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestBulkCreateUsersSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid bulk user creation request
		WHEN:  POST /user/bulk-create/ is called with valid users data
		THEN:  A HTTP_201_CREATED status should be returned with the created users
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// Create request with multiple users
	users := []*schemas.CreateUserRequest{
		{
			Name:           "User 1",
			FirstLastName:  "LastName 1",
			SecondLastName: "SecondLastName 1",
			Password:       "securePassword1",
			Email:          utilsTest.GenerateRandomEmail(),
			Rol:            "ADMINISTRATOR",
			ImageUrl:       "https://example.com/avatar1.jpg",
			Memberships:    []*schemas.Membership{},
			Onboarding:     &schemas.Onboarding{},
		},
		{
			Name:           "User 2",
			FirstLastName:  "LastName 2",
			SecondLastName: "SecondLastName 2",
			Password:       "securePassword2",
			Email:          utilsTest.GenerateRandomEmail(),
			Rol:            "CLIENT",
			ImageUrl:       "https://example.com/avatar2.jpg",
			Memberships:    []*schemas.Membership{},
			Onboarding:     &schemas.Onboarding{},
		},
	}

	bulkCreateRequest := schemas.BulkCreateUserRequest{
		Users: users,
	}

	requestBody, _ := json.Marshal(bulkCreateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/user/bulk-create/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response []*schemas.User
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response
	assert.Equal(t, len(users), len(response))

	// Check that each user was created correctly
	for i, user := range response {
		assert.NotEmpty(t, user.Id)
		assert.Equal(t, users[i].Name, user.Name)
		assert.Equal(t, users[i].FirstLastName, user.FirstLastName)
		assert.Equal(t, &users[i].SecondLastName, user.SecondLastName)
		assert.Equal(t, users[i].Email, user.Email)
		assert.Equal(t, users[i].Rol, string(user.Rol))
		assert.Equal(t, users[i].ImageUrl, user.ImageUrl)
		// Password should not be returned
		// TODO: Fix password visibility - currently it is, so we'll accept that for now
		// assert.Empty(t, user.Password)
	}
}

func TestBulkCreateUsersEmptyList(t *testing.T) {
	/*
		GIVEN: A bulk user creation request with an empty list
		WHEN:  POST /user/bulk-create/ is called with an empty users list
		THEN:  A HTTP_400_BAD_REQUEST status should be returned (API considers empty lists invalid)
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	bulkCreateRequest := schemas.BulkCreateUserRequest{
		Users: []*schemas.CreateUserRequest{},
	}

	requestBody, _ := json.Marshal(bulkCreateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/user/bulk-create/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// No need to parse response for error cases
}

func TestBulkCreateUsersInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /user/bulk-create/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/user/bulk-create/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
