package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestChangeUserRoleSuccessfully(t *testing.T) {
	/*
		GIVEN: A user exists in the database
		WHEN:  PATCH /user/{userId}/role/ is called with a valid role
		THEN:  A HTTP_200_OK status should be returned with the updated user
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test user with CLIENT role
	rol := model.UserRolClient
	user := factories.NewUserModel(db, factories.UserModelF{
		Rol: &rol,
	})

	// Prepare role change request
	changeRoleRequest := schemas.ChangeUserRoleRequest{
		Rol: schemas.UserRolAdmin,
	}

	requestBody, _ := json.Marshal(changeRoleRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/user/"+user.Id.String()+"/role/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.User
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the updated role
	assert.Equal(t, user.Id, response.Id)
	assert.Equal(t, changeRoleRequest.Rol, response.Rol)
}

func TestChangeUserRoleUserNotFound(t *testing.T) {
	/*
		GIVEN: No user exists with the provided ID
		WHEN:  PATCH /user/{userId}/role/ is called with a non-existent user ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentUserId := uuid.New()

	// Prepare role change request
	changeRoleRequest := schemas.ChangeUserRoleRequest{
		Rol: schemas.UserRolAdmin,
	}

	requestBody, _ := json.Marshal(changeRoleRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/user/"+nonExistentUserId.String()+"/role/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestChangeUserRoleInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  PATCH /user/{userId}/role/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidUserId := "invalid-uuid"

	// Prepare role change request
	changeRoleRequest := schemas.ChangeUserRoleRequest{
		Rol: schemas.UserRolAdmin,
	}

	requestBody, _ := json.Marshal(changeRoleRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/user/"+invalidUserId+"/role/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestChangeUserRoleInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: A user exists but the request body is invalid
		WHEN:  PATCH /user/{userId}/role/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test user
	user := factories.NewUserModel(db)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/user/"+user.Id.String()+"/role/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
