package community_service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestGetServicesByCommunityIdSuccessfully(t *testing.T) {
	/*
		GIVEN: Community has associated services
		WHEN:  GetServicesByCommunityId is called with valid community ID
		THEN:  All services for the community should be returned
	*/
	// GIVEN
	communityServiceController, _, db := controllerTest.NewCommunityServiceControllerTestWrapper(t)

	// Create test community and services
	testCommunity := factories.NewCommunityModel(db)
	testService1 := factories.NewServiceModel(db)
	testService2 := factories.NewServiceModel(db)

	// Create community-service associations
	associations := []*model.CommunityService{
		{
			Id:          uuid.New(),
			CommunityId: testCommunity.Id,
			ServiceId:   testService1.Id,
			AuditFields: model.AuditFields{UpdatedBy: "TEST_USER"},
		},
		{
			Id:          uuid.New(),
			CommunityId: testCommunity.Id,
			ServiceId:   testService2.Id,
			AuditFields: model.AuditFields{UpdatedBy: "TEST_USER"},
		},
	}

	for _, association := range associations {
		err := db.Create(association).Error
		assert.NoError(t, err)
	}

	// WHEN
	result, err := communityServiceController.GetServicesByCommunityId(testCommunity.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 2)

	// Verify the services are the ones we created
	serviceIds := make(map[uuid.UUID]bool)
	for _, service := range result.Services {
		serviceIds[service.Id] = true
	}
	assert.True(t, serviceIds[testService1.Id])
	assert.True(t, serviceIds[testService2.Id])
}

func TestGetServicesByCommunityIdWithNoServices(t *testing.T) {
	/*
		GIVEN: Community has no associated services
		WHEN:  GetServicesByCommunityId is called with valid community ID
		THEN:  Empty services list should be returned
	*/
	// GIVEN
	communityServiceController, _, db := controllerTest.NewCommunityServiceControllerTestWrapper(t)

	// Create test community but no services
	testCommunity := factories.NewCommunityModel(db)

	// WHEN
	result, err := communityServiceController.GetServicesByCommunityId(testCommunity.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 0)
}

func TestGetServicesByCommunityIdWithNonExistentCommunity(t *testing.T) {
	/*
		GIVEN: Community does not exist
		WHEN:  GetServicesByCommunityId is called with non-existent community ID
		THEN:  It should return empty services list
	*/
	// GIVEN
	communityServiceController, _, _ := controllerTest.NewCommunityServiceControllerTestWrapper(t)

	nonExistentCommunityId := uuid.New()

	// WHEN
	result, err := communityServiceController.GetServicesByCommunityId(nonExistentCommunityId)

	// THEN
	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.Len(t, result.Services, 0)
}
