package membership_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestGetMembershipsByUserIdSuccessfully(t *testing.T) {
	/*
		GIVEN: A user has memberships
		WHEN:  GET /membership/user/{userId}/ is called
		THEN:  A list of memberships is returned with a HTTP_200_OK status
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	user := factories.NewUserModel(db, factories.UserModelF{})
	factories.NewMembershipModel(db, factories.MembershipModelF{UserId: &user.Id})
	factories.NewMembershipModel(db, factories.MembershipModelF{UserId: &user.Id})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/membership/user/"+user.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Memberships
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.Memberships), 2)
}

func TestGetMembershipsByUserIdNotFound(t *testing.T) {
	/*
		GIVEN: A request to get memberships by user ID for a user with no memberships
		WHEN:  GET /membership/user/{userId}/ is called
		THEN:  A HTTP_200_OK status is returned with an empty list
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	user := factories.NewUserModel(db, factories.UserModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/membership/user/"+user.Id.String()+"/", nil)
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Memberships
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response.Memberships)
}

func TestGetMembershipsByUserIdInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid user ID is provided
		WHEN:  GET /membership/user/{userId}/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/membership/user/invalid-uuid/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestGetMembershipsByNonExistentUserId(t *testing.T) {
	/*
		GIVEN: A request to get memberships by a non-existent user ID
		WHEN:  GET /membership/user/{userId}/ is called
		THEN:  A HTTP_200_OK status is returned with an empty list
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentUserId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/membership/user/"+nonExistentUserId.String()+"/", nil)
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Memberships
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response.Memberships)
}
