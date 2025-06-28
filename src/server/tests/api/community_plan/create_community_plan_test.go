package community_plan_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestCreateCommunityPlanSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid community and plan exist
		WHEN:  POST /community-plan/ is called with valid data
		THEN:  A HTTP_201_CREATED status should be returned with the new community-plan
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	requestBody := schemas.CreateCommunityPlanRequest{
		CommunityId: community.Id,
		PlanId:      plan.Id,
	}

	body, _ := json.Marshal(requestBody)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-plan/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.CommunityPlan
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, community.Id, response.CommunityId)
	assert.Equal(t, plan.Id, response.PlanId)
}

func TestCreateCommunityPlanConflict(t *testing.T) {
	/*
		GIVEN: A community-plan association already exists
		WHEN:  POST /community-plan/ is called with the same data
		THEN:  A HTTP_409_CONFLICT status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityPlan := factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})

	requestBody := schemas.CreateCommunityPlanRequest{
		CommunityId: communityPlan.CommunityId,
		PlanId:      communityPlan.PlanId,
	}

	body, _ := json.Marshal(requestBody)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-plan/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCreateCommunityPlanInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /community-plan/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": "json"}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-plan/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
