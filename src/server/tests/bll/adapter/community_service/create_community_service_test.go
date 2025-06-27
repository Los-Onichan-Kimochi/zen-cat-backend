package community_service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

func TestCommunityServiceAdapter_CreatePostgresqlCommunityService_Success(t *testing.T) {
	// GIVEN: Valid community and service exist
	communityServiceAdapter, _, db := adapterTest.NewCommunityServiceAdapterTestWrapper(t)

	// Create dependencies
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service := factories.NewServiceModel(db, factories.ServiceModelF{})

	// WHEN: CreatePostgresqlCommunityService is called
	result, err := communityServiceAdapter.CreatePostgresqlCommunityService(
		community.Id,
		service.Id,
		"test_user",
	)

	// THEN: A new community-service association is created
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, community.Id, result.CommunityId)
	assert.Equal(t, service.Id, result.ServiceId)
	assert.NotEqual(t, "", result.Id)
}

func TestCommunityServiceAdapter_CreatePostgresqlCommunityService_EmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid community and service but empty updatedBy
	communityServiceAdapter, _, db := adapterTest.NewCommunityServiceAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service := factories.NewServiceModel(db, factories.ServiceModelF{})

	// WHEN: CreatePostgresqlCommunityService is called with empty updatedBy
	result, err := communityServiceAdapter.CreatePostgresqlCommunityService(
		community.Id,
		service.Id,
		"",
	)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestCommunityServiceAdapter_GetPostgresqlCommunityService_Success(t *testing.T) {
	// GIVEN: A community-service association exists
	communityServiceAdapter, _, db := adapterTest.NewCommunityServiceAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service := factories.NewServiceModel(db, factories.ServiceModelF{})

	// Create the association
	created, createErr := communityServiceAdapter.CreatePostgresqlCommunityService(
		community.Id,
		service.Id,
		"test_user",
	)
	assert.Nil(t, createErr)

	// WHEN: GetPostgresqlCommunityService is called
	result, err := communityServiceAdapter.GetPostgresqlCommunityService(
		community.Id,
		service.Id,
	)

	// THEN: The association is returned
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, created.Id, result.Id)
	assert.Equal(t, community.Id, result.CommunityId)
	assert.Equal(t, service.Id, result.ServiceId)
}

func TestCommunityServiceAdapter_GetPostgresqlCommunityService_NotFound(t *testing.T) {
	// GIVEN: No community-service association exists
	communityServiceAdapter, _, db := adapterTest.NewCommunityServiceAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service := factories.NewServiceModel(db, factories.ServiceModelF{})

	// WHEN: GetPostgresqlCommunityService is called for non-existent association
	result, err := communityServiceAdapter.GetPostgresqlCommunityService(
		community.Id,
		service.Id,
	)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Message, "not found")
}

func TestCommunityServiceAdapter_DeletePostgresqlCommunityService_Success(t *testing.T) {
	// GIVEN: A community-service association exists
	communityServiceAdapter, _, db := adapterTest.NewCommunityServiceAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service := factories.NewServiceModel(db, factories.ServiceModelF{})

	// Create the association
	_, createErr := communityServiceAdapter.CreatePostgresqlCommunityService(
		community.Id,
		service.Id,
		"test_user",
	)
	assert.Nil(t, createErr)

	// WHEN: DeletePostgresqlCommunityService is called
	err := communityServiceAdapter.DeletePostgresqlCommunityService(
		community.Id,
		service.Id,
	)

	// THEN: The association is deleted successfully
	assert.Nil(t, err)

	// Verify it was deleted
	_, getErr := communityServiceAdapter.GetPostgresqlCommunityService(
		community.Id,
		service.Id,
	)
	assert.NotNil(t, getErr)
}

func TestCommunityServiceAdapter_BulkCreatePostgresqlCommunityServices_Success(t *testing.T) {
	// GIVEN: Multiple communities and services
	communityServiceAdapter, _, db := adapterTest.NewCommunityServiceAdapterTestWrapper(t)

	community1 := factories.NewCommunityModel(db, factories.CommunityModelF{})
	community2 := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})

	requests := []*schemas.CreateCommunityServiceRequest{
		{
			CommunityId: community1.Id,
			ServiceId:   service1.Id,
		},
		{
			CommunityId: community2.Id,
			ServiceId:   service2.Id,
		},
	}

	// WHEN: BulkCreatePostgresqlCommunityServices is called
	results, err := communityServiceAdapter.BulkCreatePostgresqlCommunityServices(
		requests,
		"test_admin",
	)

	// THEN: Multiple associations are created
	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results, 2)
	
	assert.Equal(t, community1.Id, results[0].CommunityId)
	assert.Equal(t, service1.Id, results[0].ServiceId)
	assert.Equal(t, community2.Id, results[1].CommunityId)
	assert.Equal(t, service2.Id, results[1].ServiceId)
}

func TestCommunityServiceAdapter_BulkCreatePostgresqlCommunityServices_EmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid data but empty updatedBy
	communityServiceAdapter, _, db := adapterTest.NewCommunityServiceAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service := factories.NewServiceModel(db, factories.ServiceModelF{})

	requests := []*schemas.CreateCommunityServiceRequest{
		{
			CommunityId: community.Id,
			ServiceId:   service.Id,
		},
	}

	// WHEN: BulkCreatePostgresqlCommunityServices is called with empty updatedBy
	results, err := communityServiceAdapter.BulkCreatePostgresqlCommunityServices(
		requests,
		"",
	)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, results)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestCommunityServiceAdapter_GetPostgresqlServicesByCommunityId_Success(t *testing.T) {
	// GIVEN: A community with associated services
	communityServiceAdapter, _, db := adapterTest.NewCommunityServiceAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})

	// Create associations
	_, err1 := communityServiceAdapter.CreatePostgresqlCommunityService(
		community.Id,
		service1.Id,
		"test_user",
	)
	assert.Nil(t, err1)

	_, err2 := communityServiceAdapter.CreatePostgresqlCommunityService(
		community.Id,
		service2.Id,
		"test_user",
	)
	assert.Nil(t, err2)

	// WHEN: GetPostgresqlServicesByCommunityId is called
	services, err := communityServiceAdapter.GetPostgresqlServicesByCommunityId(community.Id)

	// THEN: The associated services are returned
	assert.Nil(t, err)
	assert.NotNil(t, services)
	assert.GreaterOrEqual(t, len(services), 2)

	// Verify both services are included
	foundService1 := false
	foundService2 := false
	for _, service := range services {
		if service.Id == service1.Id {
			foundService1 = true
		}
		if service.Id == service2.Id {
			foundService2 = true
		}
	}
	assert.True(t, foundService1)
	assert.True(t, foundService2)
}

func TestCommunityServiceAdapter_FetchPostgresqlCommunityServices_Success(t *testing.T) {
	// GIVEN: Multiple community-service associations exist
	communityServiceAdapter, _, db := adapterTest.NewCommunityServiceAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})

	// Create associations
	created1, err1 := communityServiceAdapter.CreatePostgresqlCommunityService(
		community.Id,
		service1.Id,
		"test_user",
	)
	assert.Nil(t, err1)

	created2, err2 := communityServiceAdapter.CreatePostgresqlCommunityService(
		community.Id,
		service2.Id,
		"test_user",
	)
	assert.Nil(t, err2)

	// WHEN: FetchPostgresqlCommunityServices is called with community filter
	associations, err := communityServiceAdapter.FetchPostgresqlCommunityServices(
		&community.Id,
		nil,
	)

	// THEN: The associations are returned
	assert.Nil(t, err)
	assert.NotNil(t, associations)
	assert.GreaterOrEqual(t, len(associations), 2)

	// Verify both associations are included
	foundAssoc1 := false
	foundAssoc2 := false
	for _, assoc := range associations {
		if assoc.Id == created1.Id {
			foundAssoc1 = true
		}
		if assoc.Id == created2.Id {
			foundAssoc2 = true
		}
	}
	assert.True(t, foundAssoc1)
	assert.True(t, foundAssoc2)
} 