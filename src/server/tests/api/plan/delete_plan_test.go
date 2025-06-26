package plan_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestDeletePlanSuccessfully(t *testing.T) {
	/*
		GIVEN: A plan exists in the database
		WHEN:  DELETE /plan/{planId}/ is called with a valid plan ID
		THEN:  A HTTP_204_NO_CONTENT status should be returned and the plan should be deleted
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test plan using factory
	plan := factories.NewPlanModel(db)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/plan/"+plan.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the plan was deleted
	var count int64
	db.Model(&plan).Where("id = ?", plan.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeletePlanNotFound(t *testing.T) {
	/*
		GIVEN: No plan exists with the provided ID
		WHEN:  DELETE /plan/{planId}/ is called with a non-existent plan ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentPlanId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/plan/"+nonExistentPlanId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeletePlanInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided
		WHEN:  DELETE /plan/{planId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidPlanId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/plan/"+invalidPlanId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
