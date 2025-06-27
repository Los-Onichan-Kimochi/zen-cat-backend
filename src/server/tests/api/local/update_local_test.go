package local_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestUpdateLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: A local exists and a valid update request is made
		WHEN:  PATCH /local/{localId}/ is called
		THEN:  The local is updated and a HTTP_200_OK status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	local := factories.NewLocalModel(db, factories.LocalModelF{})
	newName := "Updated Local Name"
	request := schemas.UpdateLocalRequest{
		LocalName: &newName,
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/local/"+local.Id.String()+"/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Local
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, newName, response.LocalName)
}

func TestUpdateLocalNotFound(t *testing.T) {
	/*
		GIVEN: A local with the given ID does not exist
		WHEN:  PATCH /local/{localId}/ is called
		THEN:  A HTTP_404_NOT_FOUND status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentId := uuid.New()
	update := "update"
	request := schemas.UpdateLocalRequest{
		LocalName: &update,
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/local/"+nonExistentId.String()+"/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateLocalInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided
		WHEN:  PATCH /local/{localId}/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	update := "update"
	request := schemas.UpdateLocalRequest{
		LocalName: &update,
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/local/invalid-uuid/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateLocalInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body is provided
		WHEN:  PATCH /local/{localId}/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/local/"+local.Id.String()+"/", strings.NewReader(`{"invalid": json`)) // Invalid JSON - missing closing quote and brace
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
