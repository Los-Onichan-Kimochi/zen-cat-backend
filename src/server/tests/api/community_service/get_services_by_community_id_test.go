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

func TestGetServicesByCommunityIdSuccessfully(t *testing.T) {
	/*
		GIVEN: A community has associated services
		WHEN:  GET /community-service/{communityId}/ is called
		THEN:  A HTTP_200_OK status is returned with a list of services
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})
	factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{CommunityId: &community.Id, ServiceId: &service1.Id})
	factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{CommunityId: &community.Id, ServiceId: &service2.Id})

	// WHEN
	url := "/community-service/" + community.Id.String() + "/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Services
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Services, 2)
}

func TestGetServicesByCommunityIdNotFound(t *testing.T) {
	/*
		GIVEN: A community has no associated services
		WHEN:  GET /community-service/{communityId}/ is called
		THEN:  A HTTP_200_OK status is returned with an empty services list
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})

	// WHEN
	url := "/community-service/" + community.Id.String() + "/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Services
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response.Services)
}

func TestGetServicesByCommunityIdInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid communityId is provided
		WHEN:  GET /community-service/{communityId}/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	url := "/community-service/invalid-uuid/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestGetServicesByNonExistentCommunityId(t *testing.T) {
	/*
		GIVEN: A communityId that does not exist is provided
		WHEN:  GET /community-service/{communityId}/ is called
		THEN:  A HTTP_200_OK status should be returned with an empty services list
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	url := "/community-service/" + uuid.New().String() + "/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Services
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response.Services)
}
