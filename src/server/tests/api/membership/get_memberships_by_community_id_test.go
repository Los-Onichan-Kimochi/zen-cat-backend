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

func TestGetMembershipsByCommunityIdSuccessfully(t *testing.T) {
	/*
		GIVEN: A community has members
		WHEN:  GET /membership/community/{communityId}/ is called
		THEN:  A list of memberships is returned with a HTTP_200_OK status
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	factories.NewMembershipModel(db, factories.MembershipModelF{CommunityId: &community.Id})
	factories.NewMembershipModel(db, factories.MembershipModelF{CommunityId: &community.Id})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/membership/community/"+community.Id.String()+"/", nil)
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

func TestGetMembershipsByCommunityIdNotFound(t *testing.T) {
	/*
		GIVEN: A request to get memberships by community ID for a community with no memberships
		WHEN:  GET /membership/community/{communityId}/ is called
		THEN:  A HTTP_200_OK status is returned with an empty list
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/membership/community/"+community.Id.String()+"/", nil)
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

func TestGetMembershipsByCommunityIdInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid community ID is provided
		WHEN:  GET /membership/community/{communityId}/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/membership/community/invalid-uuid/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestGetMembershipsByNonExistentCommunityId(t *testing.T) {
	/*
		GIVEN: A request to get memberships by a non-existent community ID
		WHEN:  GET /membership/community/{communityId}/ is called
		THEN:  A HTTP_200_OK status is returned with an empty list
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentCommunityId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/membership/community/"+nonExistentCommunityId.String()+"/", nil)
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
