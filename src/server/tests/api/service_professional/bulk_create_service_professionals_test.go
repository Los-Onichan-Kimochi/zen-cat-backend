package service_professional_test

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

func TestBulkCreateServiceProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid services and professionals exist
		WHEN:  POST /service-professional/bulk/ is called with valid associations
		THEN:  A HTTP_201_CREATED status should be returned with created associations
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create services and professionals using factories
	service1 := factories.NewServiceModel(db)
	service2 := factories.NewServiceModel(db)
	professional1 := factories.NewProfessionalModel(db)
	professional2 := factories.NewProfessionalModel(db)

	bulkRequest := schemas.BatchCreateServiceProfessionalRequest{
		ServiceProfessionals: []*schemas.CreateServiceProfessionalRequest{
			{
				ServiceId:      service1.Id,
				ProfessionalId: professional1.Id,
			},
			{
				ServiceId:      service2.Id,
				ProfessionalId: professional2.Id,
			},
		},
	}

	requestBody, _ := json.Marshal(bulkRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-professional/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.ServiceProfessionals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.ServiceProfessionals))
}

func TestBulkCreateServiceProfessionalsInvalidRequest(t *testing.T) {
	/*
		GIVEN: Invalid request body
		WHEN:  POST /service-professional/bulk/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-professional/bulk/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkCreateServiceProfessionalsNonExistentService(t *testing.T) {
	/*
		GIVEN: Non-existent service ID
		WHEN:  POST /service-professional/bulk/ is called with invalid service ID
		THEN:  A HTTP_400_BAD_REQUEST status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a professional using factory
	professional := factories.NewProfessionalModel(db)

	// Create request with non-existent service
	nonExistentServiceId, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	request := schemas.BatchCreateServiceProfessionalRequest{
		ServiceProfessionals: []*schemas.CreateServiceProfessionalRequest{
			{
				ServiceId:      nonExistentServiceId, // Non-existent service
				ProfessionalId: professional.Id,
			},
		},
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-professional/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
