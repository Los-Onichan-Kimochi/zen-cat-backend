package plan_test

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

func TestGetPlanSuccessfully(t *testing.T) {
	/*
		GIVEN: A plan exists in the database
		WHEN:  GET /plan/{planId} is called with a valid plan ID
		THEN:  A HTTP_200_OK status should be returned with the plan data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test plan using factory
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/plan/"+plan.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Plan
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct plan data
	assert.Equal(t, plan.Id, response.Id)
	assert.Equal(t, plan.Fee, response.Fee)
	assert.Equal(t, plan.Type, response.Type)
	assert.Equal(t, plan.ReservationLimit, response.ReservationLimit)
}

func TestGetPlanNotFound(t *testing.T) {
	/*
		GIVEN: No plan exists with the provided ID
		WHEN:  GET /plan/{planId} is called with a non-existent plan ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentPlanId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/plan/"+nonExistentPlanId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetPlanInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  GET /plan/{planId} is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidPlanId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/plan/"+invalidPlanId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
