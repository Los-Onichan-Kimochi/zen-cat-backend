package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestFetchServicesSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple services exist in the database
		WHEN:  GET /service/ is called
		THEN:  A HTTP_200_OK status should be returned with all services
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test services using factory
	numServices := 3
	factories.NewServiceModelBatch(db, numServices)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Services
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, len(response.Services), numServices)
}

func TestFetchServicesWithIdsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple services exist in the database
		WHEN:  GET /service/ is called with a list of service IDs
		THEN:  A HTTP_200_OK status should be returned with the filtered services
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test services using factory
	numServices := 5
	services := factories.NewServiceModelBatch(db, numServices)

	// Get a subset of service IDs
	idsToFetch := []string{
		services[0].Id.String(),
		services[2].Id.String(),
		services[4].Id.String(),
	}

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service/?ids="+strings.Join(idsToFetch, ","), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Services
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Len(t, response.Services, len(idsToFetch))

	// Verify that the correct services were returned
	returnedIds := make(map[string]bool)
	for _, service := range response.Services {
		returnedIds[service.Id.String()] = true
	}

	for _, id := range idsToFetch {
		assert.True(t, returnedIds[id])
	}
}

func TestFetchServicesEmpty(t *testing.T) {
	/*
		GIVEN: No services exist in the database
		WHEN:  GET /service/ is called
		THEN:  A HTTP_200_OK status should be returned with an empty array
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service/", nil)
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
