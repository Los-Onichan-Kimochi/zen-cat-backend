package membership_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestCreateMembershipForUserSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid request to create a membership for a specific user
		WHEN:  POST /membership/user/{userId}/ is called
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

	request := schemas.CreateMembershipForUserRequest{
		CommunityId: community.Id,
		PlanId:      plan.Id,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      schemas.MembershipStatusActive,
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/membership/user/"+user.Id.String()+"/", bytes.NewBuffer(body))
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
	assert.Equal(t, user.Id, response.UserId)
	assert.Equal(t, request.CommunityId, response.CommunityId)
}

func TestCreateMembershipForUserInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid user ID is provided
		WHEN:  POST /membership/user/{userId}/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	request := schemas.CreateMembershipForUserRequest{}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/membership/user/invalid-uuid/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateMembershipForUserNotFound(t *testing.T) {
	/*
		GIVEN: A non-existent user ID is provided
		WHEN:  POST /membership/user/{userId}/ is called
		THEN:  A HTTP_404_NOT_FOUND status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	request := schemas.CreateMembershipForUserRequest{}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/membership/user/"+uuid.NewString()+"/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
