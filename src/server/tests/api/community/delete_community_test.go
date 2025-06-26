package community_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestDeleteCommunitySuccessfully(t *testing.T) {
	/*
		GIVEN: A community exists in the database
		WHEN:  DELETE /community/{communityId}/ is called with a valid community ID
		THEN:  A HTTP_204_NO_CONTENT status should be returned and the community should be deleted
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test community using factory
	community := factories.NewCommunityModel(db)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/community/"+community.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the community was deleted
	var count int64
	db.Model(&community).Where("id = ?", community.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteCommunityNotFound(t *testing.T) {
	/*
		GIVEN: No community exists with the provided ID
		WHEN:  DELETE /community/{communityId}/ is called with a non-existent community ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentCommunityId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/community/"+nonExistentCommunityId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteCommunityInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  DELETE /community/{communityId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidCommunityId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/community/"+invalidCommunityId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
