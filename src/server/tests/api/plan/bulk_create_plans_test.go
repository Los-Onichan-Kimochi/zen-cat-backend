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

func TestBulkCreatePlansSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid bulk plan creation request
		WHEN:  POST /plan/bulk-create/ is called with valid plans data
		THEN:  A HTTP_201_CREATED status should be returned with the created plans
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// Create request with multiple plans
	reservationLimit1 := 10
	reservationLimit2 := 20
	plans := []*schemas.CreatePlanRequest{
		{
			Fee:              50.0,
			Type:             model.PlanTypeMonthly,
			ReservationLimit: &reservationLimit1,
		},
		{
			Fee:              500.0,
			Type:             model.PlanTypeAnual,
			ReservationLimit: &reservationLimit2,
		},
	}

	bulkCreateRequest := schemas.BulkCreatePlanRequest{
		Plans: plans,
	}

	requestBody, _ := json.Marshal(bulkCreateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/plan/bulk-create/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response []*schemas.Plan
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response
	assert.Len(t, response, len(plans))

	for i, plan := range response {
		assert.NotEmpty(t, plan.Id)
		assert.Equal(t, plans[i].Fee, plan.Fee)
		assert.Equal(t, plans[i].Type, plan.Type)
		assert.Equal(t, plans[i].ReservationLimit, plan.ReservationLimit)
	}
}

func TestBulkCreatePlansEmpty(t *testing.T) {
	/*
		GIVEN: A bulk plan creation request with an empty list
		WHEN:  POST /plan/bulk-create/ is called
		THEN:  A HTTP_400_BAD_REQUEST status should be returned (API considers empty lists invalid)
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	bulkCreateRequest := schemas.BulkCreatePlanRequest{
		Plans: []*schemas.CreatePlanRequest{},
	}

	requestBody, _ := json.Marshal(bulkCreateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/plan/bulk-create/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// No need to parse response for error cases
}

func TestBulkCreatePlansInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /plan/bulk-create/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/plan/bulk-create/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
