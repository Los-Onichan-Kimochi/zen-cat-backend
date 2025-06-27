package community_service_test

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

func TestCreateCommunityServiceSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid community and service exist
		WHEN:  CreateCommunityService is called with valid parameters
		THEN:  The community-service association should be created successfully
	*/
	// GIVEN
	communityServiceController, _, db := controllerTest.NewCommunityServiceControllerTestWrapper(t)

	// Create test community and service
	testCommunity := factories.NewCommunityModel(db)
	testService := factories.NewServiceModel(db)

	updatedBy := "TEST_USER"
	req := schemas.CreateCommunityServiceRequest{
		CommunityId: testCommunity.Id,
		ServiceId:   testService.Id,
	}

	// WHEN
	result, err := communityServiceController.CreateCommunityService(req, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testCommunity.Id, result.CommunityId)
	assert.Equal(t, testService.Id, result.ServiceId)

	// Verify in database
	var communityService model.CommunityService
	dbErr := db.Where("community_id = ? AND service_id = ?", testCommunity.Id, testService.Id).First(&communityService).Error
	assert.NoError(t, dbErr)
}

func TestCreateCommunityServiceWithNonExistentCommunity(t *testing.T) {
	/*
		GIVEN: Community does not exist
		WHEN:  CreateCommunityService is called with non-existent community ID
		THEN:  It should return community not found error
	*/
	// GIVEN
	communityServiceController, _, db := controllerTest.NewCommunityServiceControllerTestWrapper(t)

	// Create test service
	testService := factories.NewServiceModel(db)
	nonExistentCommunityId := uuid.New()

	updatedBy := "TEST_USER"
	req := schemas.CreateCommunityServiceRequest{
		CommunityId: nonExistentCommunityId,
		ServiceId:   testService.Id,
	}

	// WHEN
	result, err := communityServiceController.CreateCommunityService(req, updatedBy)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.CommunityNotFound.Code, err.Code)
}

func TestCreateCommunityServiceWithNonExistentService(t *testing.T) {
	/*
		GIVEN: Service does not exist
		WHEN:  CreateCommunityService is called with non-existent service ID
		THEN:  It should return service not found error
	*/
	// GIVEN
	communityServiceController, _, db := controllerTest.NewCommunityServiceControllerTestWrapper(t)

	// Create test community
	testCommunity := factories.NewCommunityModel(db)
	nonExistentServiceId := uuid.New()

	updatedBy := "TEST_USER"
	req := schemas.CreateCommunityServiceRequest{
		CommunityId: testCommunity.Id,
		ServiceId:   nonExistentServiceId,
	}

	// WHEN
	result, err := communityServiceController.CreateCommunityService(req, updatedBy)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceNotFound.Code, err.Code)
}

func TestCreateCommunityServiceAlreadyExists(t *testing.T) {
	/*
		GIVEN: Community-service association already exists
		WHEN:  CreateCommunityService is called with existing association
		THEN:  It should return community-service already exists error
	*/
	// GIVEN
	communityServiceController, _, db := controllerTest.NewCommunityServiceControllerTestWrapper(t)

	// Create test community and service
	testCommunity := factories.NewCommunityModel(db)
	testService := factories.NewServiceModel(db)

	// Create existing association
	existingCommunityService := &model.CommunityService{
		CommunityId: testCommunity.Id,
		ServiceId:   testService.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "EXISTING_USER",
		},
	}
	err := db.Create(existingCommunityService).Error
	assert.NoError(t, err)

	updatedBy := "TEST_USER"
	req := schemas.CreateCommunityServiceRequest{
		CommunityId: testCommunity.Id,
		ServiceId:   testService.Id,
	}

	// WHEN
	result, createErr := communityServiceController.CreateCommunityService(req, updatedBy)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, createErr)
	assert.Equal(t, errors.ConflictError.CommunityServiceAlreadyExists.Code, createErr.Code)
}
