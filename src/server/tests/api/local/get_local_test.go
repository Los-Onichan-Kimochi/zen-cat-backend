package local_test

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

func TestGetLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: A local exists in the database
		WHEN:  GET /local/{localId}/ is called with a valid ID
		THEN:  The local is returned with a HTTP_200_OK status
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/local/"+local.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Local
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, local.Id, response.Id)
	assert.Equal(t, local.LocalName, response.LocalName)
}

func TestGetLocalNotFound(t *testing.T) {
	/*
		GIVEN: A local with the given ID does not exist
		WHEN:  GET /local/{localId}/ is called
		THEN:  A HTTP_404_NOT_FOUND status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/local/"+nonExistentId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetLocalInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided
		WHEN:  GET /local/{localId}/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/local/invalid-uuid/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
