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

func TestBulkCreateCommunityPlansSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid community and plan data
		WHEN:  POST /community-plan/bulk-create/ is called
		THEN:  New community-plan associations are created and a HTTP_201_CREATED status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	community1 := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan1 := factories.NewPlanModel(db, factories.PlanModelF{})
	community2 := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan2 := factories.NewPlanModel(db, factories.PlanModelF{})

	request := schemas.BatchCreateCommunityPlanRequest{
		CommunityPlans: []*schemas.CreateCommunityPlanRequest{
			{CommunityId: community1.Id, PlanId: plan1.Id},
			{CommunityId: community2.Id, PlanId: plan2.Id},
		},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-plan/bulk-create/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.CommunityPlans
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.CommunityPlans, 2)
}

func TestBulkCreateCommunityPlansConflict(t *testing.T) {
	/*
		GIVEN: A community-plan association already exists
		WHEN:  POST /community-plan/bulk-create/ is called with the same data
		THEN:  A HTTP_409_CONFLICT status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityPlan := factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})

	request := schemas.BatchCreateCommunityPlanRequest{
		CommunityPlans: []*schemas.CreateCommunityPlanRequest{
			{CommunityId: communityPlan.CommunityId, PlanId: communityPlan.PlanId},
		},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-plan/bulk-create/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestBulkCreateCommunityPlansInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /community-plan/bulk-create/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": "json"}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/community-plan/bulk-create/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
