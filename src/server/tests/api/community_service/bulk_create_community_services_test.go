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

func TestBulkCreateCommunityServicesSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid community and service data
		WHEN:  POST /community-service/bulk-create/ is called
		THEN:  New community-service associations are created and a HTTP_201_CREATED status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	community1 := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	community2 := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})

	request := schemas.BatchCreateCommunityServiceRequest{
		CommunityServices: []*schemas.CreateCommunityServiceRequest{
			{CommunityId: community1.Id, ServiceId: service1.Id},
			{CommunityId: community2.Id, ServiceId: service2.Id},
		},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-service/bulk-create/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	if rec.Code != http.StatusCreated {
		t.Logf("Response body: %s", rec.Body.String())
		t.Logf("Expected status: %d, got: %d", http.StatusCreated, rec.Code)
	}
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.CommunityServices
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.CommunityServices, 2)
}

func TestBulkCreateCommunityServicesConflict(t *testing.T) {
	/*
		GIVEN: A community-service association already exists
		WHEN:  POST /community-service/bulk-create/ is called with the same data
		THEN:  A HTTP_409_CONFLICT status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityService := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	request := schemas.BatchCreateCommunityServiceRequest{
		CommunityServices: []*schemas.CreateCommunityServiceRequest{
			{CommunityId: communityService.CommunityId, ServiceId: communityService.ServiceId},
		},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-service/bulk-create/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestBulkCreateCommunityServicesInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /community-service/bulk-create/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": "json"}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-service/bulk-create/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
