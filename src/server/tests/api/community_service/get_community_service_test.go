package community_service_test

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

func TestGetCommunityServiceSuccessfully(t *testing.T) {
	/*
		GIVEN: A community-service association exists
		WHEN:  GET /community-service/{communityId}/{serviceId}/ is called
		THEN:  A HTTP_200_OK status is returned with the community-service data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityService := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// WHEN
	url := "/community-service/" + communityService.CommunityId.String() + "/" + communityService.ServiceId.String() + "/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.CommunityService
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, communityService.CommunityId, response.CommunityId)
	assert.Equal(t, communityService.ServiceId, response.ServiceId)
}

func TestGetCommunityServiceNotFound(t *testing.T) {
	/*
		GIVEN: A community-service association does not exist
		WHEN:  GET /community-service/{communityId}/{serviceId}/ is called
		THEN:  A HTTP_404_NOT_FOUND status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	communityId := uuid.New()
	serviceId := uuid.New()

	// WHEN
	url := "/community-service/" + communityId.String() + "/" + serviceId.String() + "/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetCommunityServiceInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided for communityId or serviceId
		WHEN:  GET /community-service/{communityId}/{serviceId}/ is called
		THEN:  A HTTP_400_BAD_REQUEST status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityService := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// WHEN
	url := "/community-service/invalid-uuid/" + communityService.ServiceId.String() + "/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
