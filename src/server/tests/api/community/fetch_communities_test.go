package community_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestFetchCommunitiesSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple communities exist in the database
		WHEN:  GET /community/ is called
		THEN:  A HTTP_200_OK status should be returned with all communities
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test communities using factory
	numCommunities := 3
	testCommunities := factories.NewCommunityModelBatch(db, numCommunities)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/community/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Communities
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains all communities
	assert.GreaterOrEqual(t, len(response.Communities), numCommunities)

	// Create a map of community IDs for easier verification
	communityMap := make(map[string]bool)
	for _, community := range testCommunities {
		communityMap[community.Id.String()] = false
	}

	// Check that all created communities are in the response
	for _, community := range response.Communities {
		if _, exists := communityMap[community.Id.String()]; exists {
			communityMap[community.Id.String()] = true
		}
	}

	// Verify all test communities were found in the response
	for id, found := range communityMap {
		assert.True(t, found, "Community with ID %s was not found in the response", id)
	}
}

func TestFetchCommunitiesEmpty(t *testing.T) {
	/*
		GIVEN: No communities exist in the database
		WHEN:  GET /community/ is called
		THEN:  A HTTP_200_OK status should be returned with an empty array
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/community/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Communities
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains an empty array
	assert.Empty(t, response.Communities)
}
