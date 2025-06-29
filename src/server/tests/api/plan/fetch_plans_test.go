package plan_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestFetchPlansSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple plans exist in the database
		WHEN:  GET /plan/ is called
		THEN:  A HTTP_200_OK status should be returned with all plans
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test plans using factory
	numPlans := 3
	factories.NewPlanModelBatch(db, numPlans)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/plan/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Plans
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, len(response.Plans), numPlans)
}

func TestFetchPlansWithIdsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple plans exist in the database
		WHEN:  GET /plan/ is called with a list of plan IDs
		THEN:  A HTTP_200_OK status should be returned with the filtered plans
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test plans using factory
	numPlans := 5
	plans := factories.NewPlanModelBatch(db, numPlans)

	// Get a subset of plan IDs
	idsToFetch := []string{
		plans[0].Id.String(),
		plans[2].Id.String(),
		plans[4].Id.String(),
	}

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/plan/?ids="+strings.Join(idsToFetch, ","), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Plans
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Len(t, response.Plans, len(idsToFetch))

	// Verify that the correct plans were returned
	returnedIds := make(map[string]bool)
	for _, plan := range response.Plans {
		returnedIds[plan.Id.String()] = true
	}

	for _, id := range idsToFetch {
		assert.True(t, returnedIds[id])
	}
}

func TestFetchPlansEmpty(t *testing.T) {
	/*
		GIVEN: No plans exist in the database
		WHEN:  GET /plan/ is called
		THEN:  A HTTP_200_OK status should be returned with an empty array
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/plan/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Plans
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Empty(t, response.Plans)
}
