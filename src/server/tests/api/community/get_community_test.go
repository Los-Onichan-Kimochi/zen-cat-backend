package community_test

import (
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

func TestGetCommunitySuccessfully(t *testing.T) {
	/*
		GIVEN: A community exists in the database
		WHEN:  GET /community/{communityId}/ is called with a valid community ID
		THEN:  A HTTP_200_OK status should be returned with the community data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test community using factory
	name := "Test Community"
	purpose := "Test Purpose"
	imageUrl := "https://example.com/community.jpg"

	community := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name:     &name,
		Purpose:  &purpose,
		ImageUrl: &imageUrl,
	})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/community/"+community.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Community
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct community data
	assert.Equal(t, community.Id, response.Id)
	assert.Equal(t, community.Name, response.Name)
	assert.Equal(t, community.Purpose, response.Purpose)
	assert.Equal(t, community.ImageUrl, response.ImageUrl)
	assert.Equal(t, community.NumberSubscriptions, response.NumberSubscriptions)
}

func TestGetCommunityNotFound(t *testing.T) {
	/*
		GIVEN: No community exists with the provided ID
		WHEN:  GET /community/{communityId}/ is called with a non-existent community ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentCommunityId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/community/"+nonExistentCommunityId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetCommunityInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  GET /community/{communityId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidCommunityId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/community/"+invalidCommunityId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}
