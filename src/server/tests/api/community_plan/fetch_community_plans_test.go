package community_plan_test

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

func TestFetchCommunityPlansSuccessfully(t *testing.T) {
	/*
		GIVEN: Community-plan associations exist
		WHEN:  GET /community-plan/ is called
		THEN:  A HTTP_200_OK status is returned with a list of community plans
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})
	factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/community-plan/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.CommunityPlans
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.CommunityPlans), 2)
}

func TestFetchCommunityPlansWithFilter(t *testing.T) {
	/*
		GIVEN: Community-plan associations exist
		WHEN:  GET /community-plan/ is called with a communityId filter
		THEN:  A HTTP_200_OK status is returned with the filtered community plans
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityPlan1 := factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})
	factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})

	// WHEN
	url := "/community-plan/?communityId=" + communityPlan1.CommunityId.String()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.CommunityPlans
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.CommunityPlans, 1)
	assert.Equal(t, communityPlan1.CommunityId, response.CommunityPlans[0].CommunityId)
}

func TestFetchCommunityPlansEmpty(t *testing.T) {
	/*
		GIVEN: No community-plan associations exist
		WHEN:  GET /community-plan/ is called
		THEN:  A HTTP_200_OK status is returned with an empty list
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/community-plan/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.CommunityPlans
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response.CommunityPlans)
}
