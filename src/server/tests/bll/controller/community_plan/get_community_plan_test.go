package community_plan_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestGetCommunityPlanSuccessfully(t *testing.T) {
	/*
		GIVEN: Community-plan association exists
		WHEN:  GetCommunityPlan is called with valid IDs
		THEN:  The community-plan association should be returned
	*/
	// GIVEN
	communityPlanController, _, db := controllerTest.NewCommunityPlanControllerTestWrapper(t)

	// Create test community and plan
	testCommunity := factories.NewCommunityModel(db)
	testPlan := factories.NewPlanModel(db)

	// Create community-plan association
	communityPlan := &model.CommunityPlan{
		CommunityId: testCommunity.Id,
		PlanId:      testPlan.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "TEST_USER",
		},
	}
	err := db.Create(communityPlan).Error
	assert.NoError(t, err)

	// WHEN
	result, getErr := communityPlanController.GetCommunityPlan(
		testCommunity.Id.String(),
		testPlan.Id.String(),
	)

	// THEN
	assert.Nil(t, getErr)
	assert.NotNil(t, result)
	assert.Equal(t, testCommunity.Id, result.CommunityId)
	assert.Equal(t, testPlan.Id, result.PlanId)
}

func TestGetCommunityPlanNotFound(t *testing.T) {
	/*
		GIVEN: Community-plan association does not exist
		WHEN:  GetCommunityPlan is called with non-existent association
		THEN:  It should return community-plan not found error
	*/
	// GIVEN
	communityPlanController, _, _ := controllerTest.NewCommunityPlanControllerTestWrapper(t)

	nonExistentCommunityId := uuid.New()
	nonExistentPlanId := uuid.New()

	// WHEN
	result, err := communityPlanController.GetCommunityPlan(
		nonExistentCommunityId.String(),
		nonExistentPlanId.String(),
	)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.CommunityPlanNotFound.Code, err.Code)
}

func TestGetCommunityPlanWithInvalidCommunityId(t *testing.T) {
	/*
		GIVEN: Invalid community ID format
		WHEN:  GetCommunityPlan is called with invalid community ID
		THEN:  It should return invalid updated by value error
	*/
	// GIVEN
	communityPlanController, _, _ := controllerTest.NewCommunityPlanControllerTestWrapper(t)

	invalidCommunityId := "invalid-uuid"
	validPlanId := uuid.New().String()

	// WHEN
	result, err := communityPlanController.GetCommunityPlan(
		invalidCommunityId,
		validPlanId,
	)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue.Code, err.Code)
}

func TestGetCommunityPlanWithInvalidPlanId(t *testing.T) {
	/*
		GIVEN: Invalid plan ID format
		WHEN:  GetCommunityPlan is called with invalid plan ID
		THEN:  It should return invalid updated by value error
	*/
	// GIVEN
	communityPlanController, _, _ := controllerTest.NewCommunityPlanControllerTestWrapper(t)

	validCommunityId := uuid.New().String()
	invalidPlanId := "invalid-uuid"

	// WHEN
	result, err := communityPlanController.GetCommunityPlan(
		validCommunityId,
		invalidPlanId,
	)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue.Code, err.Code)
}
