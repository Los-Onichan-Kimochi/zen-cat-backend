package community_test

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

func TestBulkDeleteCommunitiesSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple communities exist in the database
		WHEN:  DELETE /community/bulk-delete/ is called with valid community IDs
		THEN:  A HTTP_204_NO_CONTENT status should be returned and the communities should be deleted
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test communities using factory
	numCommunities := 3
	testCommunities := factories.NewCommunityModelBatch(db, numCommunities)

	// Extract community IDs
	communityIds := make([]uuid.UUID, numCommunities)
	for i, community := range testCommunities {
		communityIds[i] = community.Id
	}

	bulkDeleteRequest := schemas.BulkDeleteCommunityRequest{
		Communities: communityIds,
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/community/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the communities were deleted
	for _, communityId := range communityIds {
		var count int64
		db.Model(&testCommunities[0]).Where("id = ?", communityId).Count(&count)
		assert.Equal(t, int64(0), count)
	}
}

func TestBulkDeleteCommunitiesEmptyList(t *testing.T) {
	/*
		GIVEN: A bulk community deletion request with an empty list
		WHEN:  DELETE /community/bulk-delete/ is called with an empty communities list
		THEN:  A HTTP_204_NO_CONTENT status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	bulkDeleteRequest := schemas.BulkDeleteCommunityRequest{
		Communities: []uuid.UUID{},
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/community/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBulkDeleteCommunitiesInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  DELETE /community/bulk-delete/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/community/bulk-delete/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
