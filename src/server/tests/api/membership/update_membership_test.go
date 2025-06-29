package membership_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestUpdateMembershipSuccessfully(t *testing.T) {
	/*
		GIVEN: A membership exists and a valid update request is made
		WHEN:  PATCH /membership/{membershipId}/ is called
		THEN:  The membership is updated and a HTTP_200_OK status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	membership := factories.NewMembershipModel(db, factories.MembershipModelF{})
	newDescription := "Updated description"
	newStatus := schemas.MembershipStatusCancelled
	request := schemas.UpdateMembershipRequest{
		Description: &newDescription,
		Status:      &newStatus,
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/membership/"+membership.Id.String()+"/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Membership
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, newDescription, response.Description)
	assert.Equal(t, newStatus, response.Status)
}

func TestUpdateMembershipNotFound(t *testing.T) {
	/*
		GIVEN: A membership with the given ID does not exist
		WHEN:  PATCH /membership/{membershipId}/ is called
		THEN:  A HTTP_404_NOT_FOUND status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentId := uuid.New()
	update := "update"
	request := schemas.UpdateMembershipRequest{
		Description: &update,
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/membership/"+nonExistentId.String()+"/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateMembershipInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided
		WHEN:  PATCH /membership/{membershipId}/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	update := "update"
	request := schemas.UpdateMembershipRequest{
		Description: &update,
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/membership/invalid-uuid/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateMembershipInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body to update a membership
		WHEN:  PATCH /membership/{membershipId}/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	user := factories.NewUserModel(db, factories.UserModelF{})
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})
	membership := factories.NewMembershipModel(db, factories.MembershipModelF{
		UserId:      &user.Id,
		CommunityId: &community.Id,
		PlanId:      &plan.Id,
	})

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/membership/"+membership.Id.String()+"/", strings.NewReader(`{"invalid": json`)) // Invalid JSON - missing closing quote and brace
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
