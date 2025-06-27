package membership_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestCreateMembershipSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid request to create a membership
		WHEN:  POST /membership/ is called
		THEN:  A new membership is created and a HTTP_201_CREATED status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	user := factories.NewUserModel(db, factories.UserModelF{})
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// Create the CommunityPlan association (required for membership validation)
	communityPlan := factories.NewCommunityPlanModel(db, factories.CommunityPlanModelF{
		CommunityId: &community.Id,
		PlanId:      &plan.Id,
	})
	_ = communityPlan // Use the variable to avoid unused variable error

	startDate := time.Now()
	endDate := startDate.AddDate(1, 0, 0)

	request := schemas.CreateMembershipRequest{
		UserId:      user.Id,
		CommunityId: community.Id,
		PlanId:      plan.Id,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      schemas.MembershipStatusActive,
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/membership/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Membership
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Id)
	assert.Equal(t, request.UserId, response.UserId)
	assert.Equal(t, request.CommunityId, response.CommunityId)
}

func TestCreateMembershipInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body to create a membership
		WHEN:  POST /membership/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/membership/", strings.NewReader(`{"invalid": json`)) // Invalid JSON - missing closing quote and brace
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
