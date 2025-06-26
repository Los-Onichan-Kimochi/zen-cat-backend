package plan_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestUpdatePlanSuccessfully(t *testing.T) {
	/*
		GIVEN: A plan exists in the database
		WHEN:  PATCH /plan/{planId}/ is called with valid update data
		THEN:  A HTTP_200_OK status should be returned with the updated plan
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test plan using factory
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// Prepare update request
	newFee := 200.0
	newType := model.PlanTypeAnual
	newReservationLimit := 20
	updatePlanRequest := schemas.UpdatePlanRequest{
		Fee:              &newFee,
		Type:             &newType,
		ReservationLimit: &newReservationLimit,
	}

	requestBody, _ := json.Marshal(updatePlanRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/plan/"+plan.Id.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Plan
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the updated plan data
	assert.Equal(t, plan.Id, response.Id)
	assert.Equal(t, *updatePlanRequest.Fee, response.Fee)
	assert.Equal(t, *updatePlanRequest.Type, response.Type)
	assert.Equal(t, *updatePlanRequest.ReservationLimit, *response.ReservationLimit)
}

func TestUpdatePlanNotFound(t *testing.T) {
	/*
		GIVEN: No plan exists with the provided ID
		WHEN:  PATCH /plan/{planId}/ is called with a non-existent plan ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentPlanId := uuid.New()

	// Prepare update request
	newFee := 200.0
	updatePlanRequest := schemas.UpdatePlanRequest{
		Fee: &newFee,
	}

	requestBody, _ := json.Marshal(updatePlanRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/plan/"+nonExistentPlanId.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdatePlanInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided
		WHEN:  PATCH /plan/{planId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidPlanId := "invalid-uuid"

	// Prepare update request
	newFee := 200.0
	updatePlanRequest := schemas.UpdatePlanRequest{
		Fee: &newFee,
	}

	requestBody, _ := json.Marshal(updatePlanRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/plan/"+invalidPlanId+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdatePlanInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: A plan exists but the request body is invalid
		WHEN:  PATCH /plan/{planId}/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test plan
	plan := factories.NewPlanModel(db)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/plan/"+plan.Id.String()+"/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
