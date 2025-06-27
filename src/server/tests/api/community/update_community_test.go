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

func TestUpdateCommunitySuccessfully(t *testing.T) {
	/*
		GIVEN: A community exists in the database
		WHEN:  PATCH /community/{communityId}/ is called with valid update data
		THEN:  A HTTP_200_OK status should be returned with the updated community
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test community using factory
	name := "Original Community"
	purpose := "Original Purpose"
	community := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name:    &name,
		Purpose: &purpose,
	})

	// Prepare update request
	newName := "Updated Community"
	newPurpose := "Updated Purpose"
	updateCommunityRequest := schemas.UpdateCommunityRequest{
		Name:    &newName,
		Purpose: &newPurpose,
	}

	requestBody, _ := json.Marshal(updateCommunityRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/community/"+community.Id.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Community
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the updated community data
	assert.Equal(t, community.Id, response.Id)
	assert.Equal(t, *updateCommunityRequest.Name, response.Name)
	assert.Equal(t, *updateCommunityRequest.Purpose, response.Purpose)
}

func TestUpdateCommunityNotFound(t *testing.T) {
	/*
		GIVEN: No community exists with the provided ID
		WHEN:  PATCH /community/{communityId}/ is called with a non-existent community ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentCommunityId := uuid.New()

	// Prepare update request
	newName := "Updated Community"
	updateCommunityRequest := schemas.UpdateCommunityRequest{
		Name: &newName,
	}

	requestBody, _ := json.Marshal(updateCommunityRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/community/"+nonExistentCommunityId.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateCommunityInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  PATCH /community/{communityId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidCommunityId := "invalid-uuid"

	// Prepare update request
	newName := "Updated Community"
	updateCommunityRequest := schemas.UpdateCommunityRequest{
		Name: &newName,
	}

	requestBody, _ := json.Marshal(updateCommunityRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/community/"+invalidCommunityId+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateCommunityInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: A community exists but the request body is invalid
		WHEN:  PATCH /community/{communityId}/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test community
	community := factories.NewCommunityModel(db)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/community/"+community.Id.String()+"/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
