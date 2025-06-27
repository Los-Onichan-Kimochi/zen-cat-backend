package community_plan_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestBulkDeleteCommunityPlansSuccessfully(t *testing.T) {
	/*
		GIVEN: Community-plan associations exist
		WHEN:  DELETE /community-plan/bulk-delete/ is called with a list of associations
		THEN:  The associations are deleted and a HTTP_204_NO_CONTENT status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityPlan1 := factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})
	communityPlan2 := factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{})

	request := schemas.BulkDeleteCommunityPlanRequest{
		CommunityPlans: []*schemas.DeleteCommunityPlanRequest{
			{CommunityId: communityPlan1.CommunityId, PlanId: communityPlan1.PlanId},
			{CommunityId: communityPlan2.CommunityId, PlanId: communityPlan2.PlanId},
		},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/community-plan/bulk-delete/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	var count int64
	db.Model(&model.CommunityPlan{}).Where("community_id IN ? AND plan_id IN ?",
		[]string{communityPlan1.CommunityId.String(), communityPlan2.CommunityId.String()},
		[]string{communityPlan1.PlanId.String(), communityPlan2.PlanId.String()},
	).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestBulkDeleteCommunityPlansInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  DELETE /community-plan/bulk-delete/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": "json"}`

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/community-plan/bulk-delete/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
