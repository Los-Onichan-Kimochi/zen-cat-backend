package local_test

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

func TestFetchLocalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple locals exist in the database
		WHEN:  GET /local/ is called
		THEN:  A list of locals is returned with a HTTP_200_OK status
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	factories.NewLocalModel(db, factories.LocalModelF{})
	factories.NewLocalModel(db, factories.LocalModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/local/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Locals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.Locals), 2)
}

func TestFetchLocalsEmpty(t *testing.T) {
	/*
		GIVEN: No locals exist in the database
		WHEN:  GET /local/ is called
		THEN:  An empty list is returned with a HTTP_200_OK status
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/local/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Locals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response.Locals)
}
