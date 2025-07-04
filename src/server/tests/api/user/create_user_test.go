package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestCreateUserSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid user creation request
		WHEN:  POST /user/ is called with valid user data
		THEN:  A HTTP_201_CREATED status should be returned with the created user
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	createUserRequest := schemas.CreateUserRequest{
		Name:           "John",
		FirstLastName:  "Doe",
		SecondLastName: "Smith",
		Password:       "securePassword123",
		Email:          utilsTest.GenerateRandomEmail(),
		Rol:            "ADMINISTRATOR",
		ImageUrl:       "https://example.com/avatar.jpg",
		Memberships:    []*schemas.Membership{},
		Onboarding:     &schemas.Onboarding{},
	}

	requestBody, _ := json.Marshal(createUserRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/user/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.User
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response
	assert.NotEmpty(t, response.Id)
	assert.Equal(t, createUserRequest.Name, response.Name)
	assert.Equal(t, createUserRequest.FirstLastName, response.FirstLastName)
	assert.Equal(t, &createUserRequest.SecondLastName, response.SecondLastName)
	assert.Equal(t, createUserRequest.Email, response.Email)
	assert.Equal(t, schemas.UserRol(createUserRequest.Rol), response.Rol)
	assert.True(t, strings.HasPrefix(response.ImageUrl, createUserRequest.ImageUrl))
	// Password should not be returned
	// assert.Empty(t, response.Password)

	// Verify the user was created in the database
	var dbUser model.User
	err = db.First(&dbUser, "id = ?", response.Id).Error
	assert.NoError(t, err)
	assert.Equal(t, response.Id, dbUser.Id)
	assert.Equal(t, createUserRequest.Email, dbUser.Email)
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	/*
		GIVEN: A user with a specific email already exists
		WHEN:  POST /user/ is called with the same email
		THEN:  A HTTP_400_BAD_REQUEST status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create an existing user
	existingEmail := utilsTest.GenerateRandomEmail()
	existingUser := &model.User{
		Email:         existingEmail,
		Name:          "Existing",
		FirstLastName: "User",
		Password:      "hashedpassword",
		Rol:           "MEMBER",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(existingUser).Error
	assert.NoError(t, err)

	createUserRequest := schemas.CreateUserRequest{
		Name:           "John",
		FirstLastName:  "Doe",
		SecondLastName: "Smith",
		Password:       "securePassword123",
		Email:          existingEmail,
		Rol:            "MEMBER",
		ImageUrl:       "https://example.com/avatar.jpg",
		Memberships:    []*schemas.Membership{},
		Onboarding:     &schemas.Onboarding{},
	}

	requestBody, _ := json.Marshal(createUserRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/user/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateUserInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /user/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/user/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
