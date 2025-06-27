package service_local_test

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

func TestCreateServiceLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid service and local exist
		WHEN:  POST /service-local/ is called with valid association data
		THEN:  A HTTP_201_CREATED status should be returned with the created association
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies
	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	request := schemas.CreateServiceLocalRequest{
		ServiceId: service.Id,
		LocalId:   local.Id,
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-local/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.ServiceLocal
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct data
	assert.NotEmpty(t, response.Id)
	assert.Equal(t, request.ServiceId, response.ServiceId)
	assert.Equal(t, request.LocalId, response.LocalId)
}

func TestCreateServiceLocalInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /service-local/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-local/", strings.NewReader(`{"invalid": json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateServiceLocalNonExistentService(t *testing.T) {
	/*
		GIVEN: A request with non-existent service ID
		WHEN:  POST /service-local/ is called with invalid service ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create only local, not service
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	request := schemas.CreateServiceLocalRequest{
		ServiceId: uuid.New(), // Non-existent service
		LocalId:   local.Id,
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-local/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateServiceLocalDuplicateAssociation(t *testing.T) {
	/*
		GIVEN: An association already exists between service and local
		WHEN:  POST /service-local/ is called with duplicate association
		THEN:  A HTTP_409_CONFLICT status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies
	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create first association
	firstRequest := schemas.CreateServiceLocalRequest{
		ServiceId: service.Id,
		LocalId:   local.Id,
	}

	firstBody, _ := json.Marshal(firstRequest)
	req1 := httptest.NewRequest(http.MethodPost, "/service-local/", bytes.NewBuffer(firstBody))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Accept", "application/json")

	rec1 := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec1, req1)
	assert.Equal(t, http.StatusCreated, rec1.Code)

	// Try to create duplicate association
	duplicateRequest := schemas.CreateServiceLocalRequest{
		ServiceId: service.Id,
		LocalId:   local.Id,
	}

	duplicateBody, _ := json.Marshal(duplicateRequest)

	// WHEN
	req2 := httptest.NewRequest(http.MethodPost, "/service-local/", bytes.NewBuffer(duplicateBody))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Accept", "application/json")

	rec2 := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec2, req2)

	// THEN
	assert.Equal(t, http.StatusConflict, rec2.Code)
}
