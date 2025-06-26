package community_service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestCreateCommunityServiceSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid community and service exist
		WHEN:  POST /community-service/ is called with valid data
		THEN:  A HTTP_201_CREATED status should be returned with the new community-service
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service := factories.NewServiceModel(db, factories.ServiceModelF{})

	requestBody := schemas.CreateCommunityServiceRequest{
		CommunityId: community.Id,
		ServiceId:   service.Id,
	}

	body, _ := json.Marshal(requestBody)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-service/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.CommunityService
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, community.Id, response.CommunityId)
	assert.Equal(t, service.Id, response.ServiceId)
}

func TestCreateCommunityServiceConflict(t *testing.T) {
	/*
		GIVEN: A community-service association already exists
		WHEN:  POST /community-service/ is called with the same data
		THEN:  A HTTP_409_CONFLICT status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityService := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	requestBody := schemas.CreateCommunityServiceRequest{
		CommunityId: communityService.CommunityId,
		ServiceId:   communityService.ServiceId,
	}

	body, _ := json.Marshal(requestBody)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-service/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCreateCommunityServiceInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /community-service/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": "json"}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-service/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
