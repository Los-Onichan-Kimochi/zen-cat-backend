package service_professional_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestCreateServiceProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid service and professional exist
		WHEN:  POST /service-professional/ is called with valid association data
		THEN:  A HTTP_201_CREATED status should be returned with the created association
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies
	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	request := schemas.CreateServiceProfessionalRequest{
		ServiceId:      service.Id,
		ProfessionalId: professional.Id,
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-professional/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.ServiceProfessional
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct data
	assert.NotEmpty(t, response.Id)
	assert.Equal(t, request.ServiceId, response.ServiceId)
	assert.Equal(t, request.ProfessionalId, response.ProfessionalId)
}

func TestCreateServiceProfessionalInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /service-professional/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-professional/", strings.NewReader(`{"invalid": json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateServiceProfessionalNilUUIDs(t *testing.T) {
	/*
		GIVEN: A request with nil UUIDs
		WHEN:  POST /service-professional/ is called with nil service or professional ID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	request := schemas.CreateServiceProfessionalRequest{
		ServiceId:      uuid.Nil, // Invalid UUID
		ProfessionalId: uuid.New(),
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-professional/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateServiceProfessionalNonExistentService(t *testing.T) {
	/*
		GIVEN: A request with non-existent service ID
		WHEN:  POST /service-professional/ is called with invalid service ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create only professional, not service
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	request := schemas.CreateServiceProfessionalRequest{
		ServiceId:      uuid.New(), // Non-existent service
		ProfessionalId: professional.Id,
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-professional/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateServiceProfessionalDuplicateAssociation(t *testing.T) {
	/*
		GIVEN: An association already exists between service and professional
		WHEN:  POST /service-professional/ is called with duplicate association
		THEN:  A HTTP_409_CONFLICT status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies
	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create first association
	firstRequest := schemas.CreateServiceProfessionalRequest{
		ServiceId:      service.Id,
		ProfessionalId: professional.Id,
	}

	firstBody, _ := json.Marshal(firstRequest)
	req1 := httptest.NewRequest(http.MethodPost, "/service-professional/", bytes.NewBuffer(firstBody))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Accept", "application/json")

	rec1 := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec1, req1)
	assert.Equal(t, http.StatusCreated, rec1.Code)

	// Try to create duplicate association
	duplicateRequest := schemas.CreateServiceProfessionalRequest{
		ServiceId:      service.Id,
		ProfessionalId: professional.Id,
	}

	duplicateBody, _ := json.Marshal(duplicateRequest)

	// WHEN
	req2 := httptest.NewRequest(http.MethodPost, "/service-professional/", bytes.NewBuffer(duplicateBody))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Accept", "application/json")

	rec2 := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec2, req2)

	// THEN
	assert.Equal(t, http.StatusConflict, rec2.Code)
}
