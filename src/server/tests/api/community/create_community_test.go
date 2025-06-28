package community_test

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

func TestCreateCommunitySuccessfully(t *testing.T) {
	/*
		GIVEN: A valid community creation request
		WHEN:  POST /community/ is called with valid community data
		THEN:  A HTTP_201_CREATED status should be returned with the created community
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	createCommunityRequest := schemas.CreateCommunityRequest{
		Name:     "New Community",
		Purpose:  "A community for testing purposes",
		ImageUrl: "https://example.com/community.jpg",
	}

	requestBody, _ := json.Marshal(createCommunityRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Community
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response
	assert.NotEmpty(t, response.Id)
	assert.Equal(t, createCommunityRequest.Name, response.Name)
	assert.Equal(t, createCommunityRequest.Purpose, response.Purpose)
	assert.Equal(t, createCommunityRequest.ImageUrl, response.ImageUrl)
	assert.Equal(t, 0, response.NumberSubscriptions)
}

func TestCreateCommunityDuplicateName(t *testing.T) {
	/*
		GIVEN: A community with a specific name already exists
		WHEN:  POST /community/ is called with the same name
		THEN:  A HTTP_400_BAD_REQUEST status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create an existing community
	name := "Existing Community"
	existingCommunity := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name: &name,
	})

	createCommunityRequest := schemas.CreateCommunityRequest{
		Name:     existingCommunity.Name, // Same name as existing community
		Purpose:  "Another community with same name",
		ImageUrl: "https://example.com/community2.jpg",
	}

	requestBody, _ := json.Marshal(createCommunityRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCommunityInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /community/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
