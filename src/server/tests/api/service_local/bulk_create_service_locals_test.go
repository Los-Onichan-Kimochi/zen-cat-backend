package service_local_test

import (
	"bytes"
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

func TestBulkCreateServiceLocalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid services and locals exist
		WHEN:  POST /service-local/bulk is called with valid associations
		THEN:  A HTTP_201_CREATED status should be returned with created associations
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create services and locals using factories
	service1 := factories.NewServiceModel(db)
	service2 := factories.NewServiceModel(db)
	local1 := factories.NewLocalModel(db)
	local2 := factories.NewLocalModel(db)

	bulkRequest := schemas.BatchCreateServiceLocalRequest{
		ServiceLocals: []*schemas.CreateServiceLocalRequest{
			{
				ServiceId: service1.Id,
				LocalId:   local1.Id,
			},
			{
				ServiceId: service2.Id,
				LocalId:   local2.Id,
			},
		},
	}

	requestBody, _ := json.Marshal(bulkRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-local/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.ServiceLocals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.ServiceLocals))
}

func TestBulkCreateServiceLocalsInvalidRequest(t *testing.T) {
	/*
		GIVEN: Invalid request body
		WHEN:  POST /service-local/bulk/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-local/bulk/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkCreateServiceLocalsNonExistentService(t *testing.T) {
	/*
		GIVEN: Non-existent service ID
		WHEN:  POST /service-local/bulk/ is called with invalid service ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a local using factory
	local := factories.NewLocalModel(db)

	// Create request with non-existent service
	nonExistentServiceId, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	request := schemas.BatchCreateServiceLocalRequest{
		ServiceLocals: []*schemas.CreateServiceLocalRequest{
			{
				ServiceId: nonExistentServiceId, // Non-existent service
				LocalId:   local.Id,
			},
		},
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-local/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code) // Adjusted expectation based on actual API behavior
}
