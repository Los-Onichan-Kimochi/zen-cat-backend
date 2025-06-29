package community_plan_test

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

func TestGetCommunityPlanSuccessfully(t *testing.T) {
	/*
		GIVEN: A community-plan association exists
		WHEN:  GET /community-plan/{communityId}/{planId}/ is called
		THEN:  A HTTP_200_OK status is returned with the community-plan data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityPlan := factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})

	// WHEN
	url := "/community-plan/" + communityPlan.CommunityId.String() + "/" + communityPlan.PlanId.String() + "/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.CommunityPlan
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, communityPlan.CommunityId, response.CommunityId)
	assert.Equal(t, communityPlan.PlanId, response.PlanId)
}

func TestGetCommunityPlanNotFound(t *testing.T) {
	/*
		GIVEN: A community-plan association does not exist
		WHEN:  GET /community-plan/{communityId}/{planId}/ is called
		THEN:  A HTTP_404_NOT_FOUND status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	communityId := uuid.New()
	planId := uuid.New()

	// WHEN
	url := "/community-plan/" + communityId.String() + "/" + planId.String() + "/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetCommunityPlanInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided for communityId or planId
		WHEN:  GET /community-plan/{communityId}/{planId}/ is called
		THEN:  A HTTP_400_BAD_REQUEST status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityPlan := factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})

	// WHEN
	url := "/community-plan/invalid-uuid/" + communityPlan.PlanId.String() + "/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
