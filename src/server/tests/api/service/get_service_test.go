package service_test

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

func TestGetServiceSuccessfully(t *testing.T) {
	/*
		GIVEN: A service exists in the database
		WHEN:  GET /service/{serviceId} is called with a valid service ID
		THEN:  A HTTP_200_OK status should be returned with the service data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test service using factory
	service := factories.NewServiceModel(db, factories.ServiceModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service/"+service.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Service
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct service data
	assert.Equal(t, service.Id, response.Id)
	assert.Equal(t, service.Name, response.Name)
	assert.Equal(t, service.Description, response.Description)
	assert.Equal(t, service.ImageUrl, response.ImageUrl)
	assert.Equal(t, service.IsVirtual, response.IsVirtual)
}

func TestGetServiceNotFound(t *testing.T) {
	/*
		GIVEN: No service exists with the provided ID
		WHEN:  GET /service/{serviceId} is called with a non-existent service ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentServiceId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service/"+nonExistentServiceId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetServiceInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  GET /service/{serviceId} is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidServiceId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service/"+invalidServiceId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
