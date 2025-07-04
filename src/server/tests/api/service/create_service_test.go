package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestCreateServiceSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid service creation request
		WHEN:  POST /service/ is called with valid service data
		THEN:  A HTTP_201_CREATED status should be returned with the created service
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	createServiceRequest := schemas.CreateServiceRequest{
		Name:        "New Service",
		Description: "A new service for testing",
		ImageUrl:    "http://example.com/image.png",
		IsVirtual:   true,
	}

	requestBody, _ := json.Marshal(createServiceRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Service
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response
	assert.NotEmpty(t, response.Id)
	assert.Equal(t, createServiceRequest.Name, response.Name)
	assert.Equal(t, createServiceRequest.Description, response.Description)
	assert.True(t, strings.HasPrefix(response.ImageUrl, createServiceRequest.ImageUrl))
	assert.Equal(t, createServiceRequest.IsVirtual, response.IsVirtual)
}

func TestCreateServiceInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /service/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
