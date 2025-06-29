package community_service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestFetchCommunityServicesSuccessfully(t *testing.T) {
	/*
		GIVEN: Community-service associations exist
		WHEN:  GET /community-service/ is called
		THEN:  A HTTP_200_OK status is returned with a list of community services
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})
	factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/community-service/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.CommunityServices
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.CommunityServices), 2)
}

func TestFetchCommunityServicesWithFilter(t *testing.T) {
	/*
		GIVEN: Community-service associations exist
		WHEN:  GET /community-service/ is called with a communityId filter
		THEN:  A HTTP_200_OK status is returned with the filtered community services
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityService1 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})
	factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// WHEN
	url := "/community-service/?communityId=" + communityService1.CommunityId.String()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.CommunityServices
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.CommunityServices, 1)
	assert.Equal(t, communityService1.CommunityId, response.CommunityServices[0].CommunityId)
}

func TestFetchCommunityServicesEmpty(t *testing.T) {
	/*
		GIVEN: No community-service associations exist
		WHEN:  GET /community-service/ is called
		THEN:  A HTTP_200_OK status is returned with an empty list
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/community-service/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.CommunityServices
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response.CommunityServices)
}
