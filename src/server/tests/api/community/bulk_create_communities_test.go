package community_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestBulkCreateCommunitiesSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid bulk community creation request
		WHEN:  POST /community/bulk-create/ is called with valid communities data
		THEN:  A HTTP_201_CREATED status should be returned with the created communities
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// Create request with multiple communities
	communities := []*schemas.CreateCommunityRequest{
		{
			Name:     "Community 1",
			Purpose:  "Purpose 1",
			ImageUrl: "https://example.com/community1.jpg",
		},
		{
			Name:     "Community 2",
			Purpose:  "Purpose 2",
			ImageUrl: "https://example.com/community2.jpg",
		},
		{
			Name:     "Community 3",
			Purpose:  "Purpose 3",
			ImageUrl: "https://example.com/community3.jpg",
		},
	}

	bulkCreateRequest := schemas.BatchCreateCommunityRequest{
		Communities: communities,
	}

	requestBody, _ := json.Marshal(bulkCreateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community/bulk-create/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response []*schemas.Community
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response
	assert.Equal(t, len(communities), len(response))

	// Check that each community was created correctly
	for i, community := range response {
		assert.NotEmpty(t, community.Id)
		assert.Equal(t, communities[i].Name, community.Name)
		assert.Equal(t, communities[i].Purpose, community.Purpose)
		assert.Equal(t, communities[i].ImageUrl, community.ImageUrl)
	}
}

func TestBulkCreateCommunitiesEmptyList(t *testing.T) {
	/*
		GIVEN: A bulk community creation request with an empty list
		WHEN:  POST /community/bulk-create/ is called with an empty communities list
		THEN:  A HTTP_400_BAD_REQUEST status should be returned (API considers empty lists invalid)
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	bulkCreateRequest := schemas.BatchCreateCommunityRequest{
		Communities: []*schemas.CreateCommunityRequest{},
	}

	requestBody, _ := json.Marshal(bulkCreateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community/bulk-create/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// No need to parse response for error cases
}

func TestBulkCreateCommunitiesInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /community/bulk-create/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community/bulk-create/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
