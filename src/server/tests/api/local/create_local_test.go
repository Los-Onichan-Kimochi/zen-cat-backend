package local_test

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

func TestCreateLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid request to create a local
		WHEN:  POST /local/ is called
		THEN:  A new local is created and a HTTP_201_CREATED status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	request := schemas.CreateLocalRequest{
		LocalName:      "Test Local",
		StreetName:     "Test Street",
		BuildingNumber: "123",
		District:       "Test District",
		Province:       "Test Province",
		Region:         "Test Region",
		Reference:      "Near the test park",
		Capacity:       100,
		ImageUrl:       "http://example.com/image.jpg",
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/local/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Local
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Id)
	assert.Equal(t, request.LocalName, response.LocalName)
	assert.Equal(t, request.StreetName, response.StreetName)
}

func TestCreateLocalInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /local/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/local/", strings.NewReader(`{"invalid": json`)) // Invalid JSON - missing closing quote and brace
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
