package community_plan_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestCreateCommunityPlanSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid community and plan exist
		WHEN:  CreateCommunityPlan is called with valid parameters
		THEN:  The community-plan association should be created successfully
	*/
	// GIVEN
	communityPlanController, _, db := controllerTest.NewCommunityPlanControllerTestWrapper(t)

	// Create test community and plan
	testCommunity := factories.NewCommunityModel(db)
	testPlan := factories.NewPlanModel(db)

	updatedBy := "TEST_USER"
	req := schemas.CreateCommunityPlanRequest{
		CommunityId: testCommunity.Id,
		PlanId:      testPlan.Id,
	}

	// WHEN
	result, err := communityPlanController.CreateCommunityPlan(req, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testCommunity.Id, result.CommunityId)
	assert.Equal(t, testPlan.Id, result.PlanId)

	// Verify in database
	var communityPlan model.CommunityPlan
	dbErr := db.Where("community_id = ? AND plan_id = ?", testCommunity.Id, testPlan.Id).First(&communityPlan).Error
	assert.NoError(t, dbErr)
}

func TestCreateCommunityPlanWithNonExistentCommunity(t *testing.T) {
	/*
		GIVEN: Community does not exist
		WHEN:  CreateCommunityPlan is called with non-existent community ID
		THEN:  It should return community not found error
	*/
	// GIVEN
	communityPlanController, _, db := controllerTest.NewCommunityPlanControllerTestWrapper(t)

	// Create test plan
	testPlan := factories.NewPlanModel(db)
	nonExistentCommunityId := uuid.New()

	updatedBy := "TEST_USER"
	req := schemas.CreateCommunityPlanRequest{
		CommunityId: nonExistentCommunityId,
		PlanId:      testPlan.Id,
	}

	// WHEN
	result, err := communityPlanController.CreateCommunityPlan(req, updatedBy)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.CommunityNotFound.Code, err.Code)
}

func TestCreateCommunityPlanWithNonExistentPlan(t *testing.T) {
	/*
		GIVEN: Plan does not exist
		WHEN:  CreateCommunityPlan is called with non-existent plan ID
		THEN:  It should return plan not found error
	*/
	// GIVEN
	communityPlanController, _, db := controllerTest.NewCommunityPlanControllerTestWrapper(t)

	// Create test community
	testCommunity := factories.NewCommunityModel(db)
	nonExistentPlanId := uuid.New()

	updatedBy := "TEST_USER"
	req := schemas.CreateCommunityPlanRequest{
		CommunityId: testCommunity.Id,
		PlanId:      nonExistentPlanId,
	}

	// WHEN
	result, err := communityPlanController.CreateCommunityPlan(req, updatedBy)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.PlanNotFound.Code, err.Code)
}

func TestCreateCommunityPlanAlreadyExists(t *testing.T) {
	/*
		GIVEN: Community-plan association already exists
		WHEN:  CreateCommunityPlan is called with existing association
		THEN:  It should return community-plan already exists error
	*/
	// GIVEN
	communityPlanController, _, db := controllerTest.NewCommunityPlanControllerTestWrapper(t)

	// Create test community and plan
	testCommunity := factories.NewCommunityModel(db)
	testPlan := factories.NewPlanModel(db)

	// Create existing association
	existingCommunityPlan := &model.CommunityPlan{
		CommunityId: testCommunity.Id,
		PlanId:      testPlan.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "EXISTING_USER",
		},
	}
	err := db.Create(existingCommunityPlan).Error
	assert.NoError(t, err)

	updatedBy := "TEST_USER"
	req := schemas.CreateCommunityPlanRequest{
		CommunityId: testCommunity.Id,
		PlanId:      testPlan.Id,
	}

	// WHEN
	result, createErr := communityPlanController.CreateCommunityPlan(req, updatedBy)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, createErr)
	assert.Equal(t, errors.ConflictError.CommunityPlanAlreadyExists.Code, createErr.Code)
}
