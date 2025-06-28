package service_test

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

func TestUpdateServiceSuccessfully(t *testing.T) {
	/*
		GIVEN: A service exists in the database
		WHEN:  PATCH /service/{serviceId}/ is called with valid update data
		THEN:  A HTTP_200_OK status should be returned with the updated service
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test service using factory
	service := factories.NewServiceModel(db, factories.ServiceModelF{})

	// Prepare update request
	newName := "Updated Service"
	newDescription := "Updated description"
	newIsVirtual := false
	updateServiceRequest := schemas.UpdateServiceRequest{
		Name:        &newName,
		Description: &newDescription,
		IsVirtual:   &newIsVirtual,
	}

	requestBody, _ := json.Marshal(updateServiceRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/service/"+service.Id.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Service
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the updated service data
	assert.Equal(t, service.Id, response.Id)
	assert.Equal(t, *updateServiceRequest.Name, response.Name)
	assert.Equal(t, *updateServiceRequest.Description, response.Description)
	assert.Equal(t, *updateServiceRequest.IsVirtual, response.IsVirtual)
}

func TestUpdateServiceNotFound(t *testing.T) {
	/*
		GIVEN: No service exists with the provided ID
		WHEN:  PATCH /service/{serviceId}/ is called with a non-existent service ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentServiceId := uuid.New()

	// Prepare update request
	newName := "Updated Service"
	updateServiceRequest := schemas.UpdateServiceRequest{
		Name: &newName,
	}

	requestBody, _ := json.Marshal(updateServiceRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/service/"+nonExistentServiceId.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateServiceInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided
		WHEN:  PATCH /service/{serviceId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidServiceId := "invalid-uuid"

	// Prepare update request
	newName := "Updated Service"
	updateServiceRequest := schemas.UpdateServiceRequest{
		Name: &newName,
	}

	requestBody, _ := json.Marshal(updateServiceRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/service/"+invalidServiceId+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateServiceInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: A service exists but the request body is invalid
		WHEN:  PATCH /service/{serviceId}/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test service
	service := factories.NewServiceModel(db)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/service/"+service.Id.String()+"/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
