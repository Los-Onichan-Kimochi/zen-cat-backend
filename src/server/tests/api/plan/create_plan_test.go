package plan_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestCreatePlanSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid plan creation request
		WHEN:  POST /plan/ is called with valid plan data
		THEN:  A HTTP_201_CREATED status should be returned with the created plan
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	reservationLimit := 10
	createPlanRequest := schemas.CreatePlanRequest{
		Fee:              100.50,
		Type:             model.PlanTypeMonthly,
		ReservationLimit: &reservationLimit,
	}

	requestBody, _ := json.Marshal(createPlanRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/plan/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Plan
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response
	assert.NotEmpty(t, response.Id)
	assert.Equal(t, createPlanRequest.Fee, response.Fee)
	assert.Equal(t, createPlanRequest.Type, response.Type)
	assert.Equal(t, createPlanRequest.ReservationLimit, response.ReservationLimit)
}

func TestCreatePlanInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /plan/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/plan/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
