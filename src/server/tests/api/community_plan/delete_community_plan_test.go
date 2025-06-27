package community_plan_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestDeleteCommunityPlanSuccessfully(t *testing.T) {
	/*
		GIVEN: A community-plan association exists
		WHEN:  DELETE /community-plan/{communityId}/{planId}/ is called
		THEN:  The association is deleted and a HTTP_204_NO_CONTENT status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityPlan := factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})

	// WHEN
	url := "/community-plan/" + communityPlan.CommunityId.String() + "/" + communityPlan.PlanId.String() + "/"
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	var count int64
	db.Model(&model.CommunityPlan{}).Where("community_id = ? AND plan_id = ?", communityPlan.CommunityId, communityPlan.PlanId).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteCommunityPlanNotFound(t *testing.T) {
	/*
		GIVEN: A community-plan association does not exist
		WHEN:  DELETE /community-plan/{communityId}/{planId}/ is called
		THEN:  A HTTP_404_NOT_FOUND status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	communityId := uuid.New()
	planId := uuid.New()

	// WHEN
	url := "/community-plan/" + communityId.String() + "/" + planId.String() + "/"
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteCommunityPlanInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided for communityId or planId
		WHEN:  DELETE /community-plan/{communityId}/{planId}/ is called
		THEN:  A HTTP_400_BAD_REQUEST status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityPlan := factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})

	// WHEN
	url := "/community-plan/invalid-uuid/" + communityPlan.PlanId.String() + "/"
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
