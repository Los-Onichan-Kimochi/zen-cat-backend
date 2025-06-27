package user_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestDeleteUserSuccessfully(t *testing.T) {
	/*
		GIVEN: A user exists in the database
		WHEN:  DELETE /user/{userId}/ is called with a valid user ID
		THEN:  A HTTP_204_NO_CONTENT status should be returned and the user should be deleted
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test user using factory
	user := factories.NewUserModel(db)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/user/"+user.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the user was deleted
	var count int64
	db.Model(&user).Where("id = ?", user.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteUserNotFound(t *testing.T) {
	/*
		GIVEN: No user exists with the provided ID
		WHEN:  DELETE /user/{userId}/ is called with a non-existent user ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentUserId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/user/"+nonExistentUserId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteUserInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  DELETE /user/{userId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidUserId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/user/"+invalidUserId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
