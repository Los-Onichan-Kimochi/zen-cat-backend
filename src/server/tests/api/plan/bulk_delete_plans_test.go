package plan_test

import (
	"bytes"
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

func TestBulkDeletePlansSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple plans exist in the database
		WHEN:  DELETE /plan/bulk-delete/ is called with valid plan IDs
		THEN:  A HTTP_204_NO_CONTENT status should be returned and the plans should be deleted
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test plans using factory
	numPlans := 3
	testPlans := factories.NewPlanModelBatch(db, numPlans)

	// Extract plan IDs
	planIds := make([]uuid.UUID, numPlans)
	for i, plan := range testPlans {
		planIds[i] = plan.Id
	}

	bulkDeleteRequest := schemas.BulkDeletePlanRequest{
		Plans: planIds,
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/plan/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the plans were deleted
	for _, planId := range planIds {
		var count int64
		db.Model(&testPlans[0]).Where("id = ?", planId).Count(&count)
		assert.Equal(t, int64(0), count)
	}
}

func TestBulkDeletePlansEmptyList(t *testing.T) {
	/*
		GIVEN: A bulk plan deletion request with an empty list
		WHEN:  DELETE /plan/bulk-delete/ is called with an empty plans list
		THEN:  A HTTP_204_NO_CONTENT status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	bulkDeleteRequest := schemas.BulkDeletePlanRequest{
		Plans: []uuid.UUID{},
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/plan/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBulkDeletePlansInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  DELETE /plan/bulk-delete/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/plan/bulk-delete/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
