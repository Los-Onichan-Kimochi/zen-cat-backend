package professional_test

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

func TestFetchProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple professionals exist in the database
		WHEN:  GET /professional/ is called
		THEN:  A HTTP_200_OK status should be returned with the professionals data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test professionals using factory
	numProfessionals := 3
	testProfessionals := factories.NewProfessionalModelBatch(db, numProfessionals)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/professional/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Professionals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct number of professionals
	assert.GreaterOrEqual(t, len(response.Professionals), numProfessionals)

	// Verify that our created professionals are in the response
	foundProfessionals := make(map[string]bool)
	for _, professional := range response.Professionals {
		foundProfessionals[professional.Id.String()] = true
	}

	for _, testProfessional := range testProfessionals {
		assert.True(t, foundProfessionals[testProfessional.Id.String()],
			"Created professional ID %s not found in response", testProfessional.Id.String())
	}
}

func TestFetchProfessionalsEmpty(t *testing.T) {
	/*
		GIVEN: No professionals exist in the database
		WHEN:  GET /professional/ is called
		THEN:  A HTTP_200_OK status should be returned with an empty professionals list
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/professional/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Professionals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Should return empty array, not null
	assert.NotNil(t, response.Professionals)
}
