package service_local_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestBulkCreateServiceLocalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid service-local data for bulk creation
		WHEN:  BulkCreatePostgresqlServiceLocals is called
		THEN:  All associations are created and returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})
	local1 := factories.NewLocalModel(db, factories.LocalModelF{})
	local2 := factories.NewLocalModel(db, factories.LocalModelF{})

	serviceLocalsData := []*schemas.CreateServiceLocalRequest{
		{
			ServiceId: service1.Id,
			LocalId:   local1.Id,
		},
		{
			ServiceId: service2.Id,
			LocalId:   local2.Id,
		},
	}
	updatedBy := "test-admin"

	// WHEN
	serviceLocals, err := adapter.BulkCreatePostgresqlServiceLocals(serviceLocalsData, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceLocals)
	assert.Equal(t, 2, len(serviceLocals))

	// Verify first association
	assert.NotEmpty(t, serviceLocals[0].Id)
	assert.Equal(t, service1.Id, serviceLocals[0].ServiceId)
	assert.Equal(t, local1.Id, serviceLocals[0].LocalId)

	// Verify second association
	assert.NotEmpty(t, serviceLocals[1].Id)
	assert.Equal(t, service2.Id, serviceLocals[1].ServiceId)
	assert.Equal(t, local2.Id, serviceLocals[1].LocalId)
}

func TestBulkCreateServiceLocalsWithEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: Valid service-local data but empty updatedBy
		WHEN:  BulkCreatePostgresqlServiceLocals is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	serviceLocalsData := []*schemas.CreateServiceLocalRequest{
		{
			ServiceId: service.Id,
			LocalId:   local.Id,
		},
	}
	updatedBy := ""

	// WHEN
	serviceLocals, err := adapter.BulkCreatePostgresqlServiceLocals(serviceLocalsData, updatedBy)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, serviceLocals)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue, *err)
}

func TestFetchServiceLocalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple service-local associations exist
		WHEN:  FetchPostgresqlServiceLocals is called
		THEN:  All matching associations are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})
	local1 := factories.NewLocalModel(db, factories.LocalModelF{})
	local2 := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create associations
	serviceLocal1, err1 := adapter.CreatePostgresqlServiceLocal(service1.Id, local1.Id, "test-admin")
	assert.Nil(t, err1)
	serviceLocal2, err2 := adapter.CreatePostgresqlServiceLocal(service2.Id, local2.Id, "test-admin")
	assert.Nil(t, err2)

	// WHEN - Fetch all associations
	serviceLocals, err := adapter.FetchPostgresqlServiceLocals(nil, nil)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceLocals)
	assert.GreaterOrEqual(t, len(serviceLocals), 2)

	// Find our created associations
	foundServiceLocal1 := false
	foundServiceLocal2 := false
	for _, serviceLocal := range serviceLocals {
		if serviceLocal.Id == serviceLocal1.Id {
			foundServiceLocal1 = true
		}
		if serviceLocal.Id == serviceLocal2.Id {
			foundServiceLocal2 = true
		}
	}
	assert.True(t, foundServiceLocal1)
	assert.True(t, foundServiceLocal2)
}

func TestFetchServiceLocalsWithServiceFilter(t *testing.T) {
	/*
		GIVEN: Associations exist for different services
		WHEN:  FetchPostgresqlServiceLocals is called with service filter
		THEN:  Only associations for that service are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})
	local1 := factories.NewLocalModel(db, factories.LocalModelF{})
	local2 := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create associations
	serviceLocal1, err1 := adapter.CreatePostgresqlServiceLocal(service1.Id, local1.Id, "test-admin")
	assert.Nil(t, err1)
	_, err2 := adapter.CreatePostgresqlServiceLocal(service2.Id, local2.Id, "test-admin")
	assert.Nil(t, err2)

	// WHEN - Filter by service1
	serviceLocals, err := adapter.FetchPostgresqlServiceLocals(&service1.Id, nil)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceLocals)
	assert.GreaterOrEqual(t, len(serviceLocals), 1)

	// Verify all returned associations are for service1
	for _, serviceLocal := range serviceLocals {
		assert.Equal(t, service1.Id, serviceLocal.ServiceId)
	}

	// Verify our association is in the results
	foundServiceLocal1 := false
	for _, serviceLocal := range serviceLocals {
		if serviceLocal.Id == serviceLocal1.Id {
			foundServiceLocal1 = true
		}
	}
	assert.True(t, foundServiceLocal1)
}

func TestFetchServiceLocalsWithLocalFilter(t *testing.T) {
	/*
		GIVEN: Associations exist for different locals
		WHEN:  FetchPostgresqlServiceLocals is called with local filter
		THEN:  Only associations for that local are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})
	local1 := factories.NewLocalModel(db, factories.LocalModelF{})
	local2 := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create associations
	serviceLocal1, err1 := adapter.CreatePostgresqlServiceLocal(service1.Id, local1.Id, "test-admin")
	assert.Nil(t, err1)
	_, err2 := adapter.CreatePostgresqlServiceLocal(service2.Id, local2.Id, "test-admin")
	assert.Nil(t, err2)

	// WHEN - Filter by local1
	serviceLocals, err := adapter.FetchPostgresqlServiceLocals(nil, &local1.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceLocals)
	assert.GreaterOrEqual(t, len(serviceLocals), 1)

	// Verify all returned associations are for local1
	for _, serviceLocal := range serviceLocals {
		assert.Equal(t, local1.Id, serviceLocal.LocalId)
	}

	// Verify our association is in the results
	foundServiceLocal1 := false
	for _, serviceLocal := range serviceLocals {
		if serviceLocal.Id == serviceLocal1.Id {
			foundServiceLocal1 = true
		}
	}
	assert.True(t, foundServiceLocal1)
}

func TestBulkDeleteServiceLocalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple service-local associations exist
		WHEN:  BulkDeletePostgresqlServiceLocals is called
		THEN:  All specified associations are deleted
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service3 := factories.NewServiceModel(db, factories.ServiceModelF{})
	local1 := factories.NewLocalModel(db, factories.LocalModelF{})
	local2 := factories.NewLocalModel(db, factories.LocalModelF{})
	local3 := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create associations
	_, err1 := adapter.CreatePostgresqlServiceLocal(service1.Id, local1.Id, "test-admin")
	assert.Nil(t, err1)
	_, err2 := adapter.CreatePostgresqlServiceLocal(service2.Id, local2.Id, "test-admin")
	assert.Nil(t, err2)
	_, err3 := adapter.CreatePostgresqlServiceLocal(service3.Id, local3.Id, "test-admin")
	assert.Nil(t, err3)

	deleteRequests := []*schemas.DeleteServiceLocalRequest{
		{
			ServiceId: service1.Id,
			LocalId:   local1.Id,
		},
		{
			ServiceId: service2.Id,
			LocalId:   local2.Id,
		},
	}

	// WHEN
	err := adapter.BulkDeletePostgresqlServiceLocals(deleteRequests)

	// THEN
	assert.Nil(t, err)

	// Verify deleted associations cannot be found
	_, getErr1 := adapter.GetPostgresqlServiceLocal(service1.Id, local1.Id)
	assert.NotNil(t, getErr1)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceLocalNotFound, *getErr1)

	_, getErr2 := adapter.GetPostgresqlServiceLocal(service2.Id, local2.Id)
	assert.NotNil(t, getErr2)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceLocalNotFound, *getErr2)

	// Verify non-deleted association still exists
	serviceLocal3Result, getErr3 := adapter.GetPostgresqlServiceLocal(service3.Id, local3.Id)
	assert.Nil(t, getErr3)
	assert.NotNil(t, serviceLocal3Result)
	assert.Equal(t, service3.Id, serviceLocal3Result.ServiceId)
	assert.Equal(t, local3.Id, serviceLocal3Result.LocalId)
}

func TestBulkDeleteServiceLocalsWithEmptyList(t *testing.T) {
	/*
		GIVEN: An empty list of delete requests
		WHEN:  BulkDeletePostgresqlServiceLocals is called
		THEN:  No error occurs
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	emptyRequests := []*schemas.DeleteServiceLocalRequest{}

	// WHEN
	err := adapter.BulkDeletePostgresqlServiceLocals(emptyRequests)

	// THEN
	assert.Nil(t, err)
}

func TestBulkDeleteServiceLocalsWithInvalidIds(t *testing.T) {
	/*
		GIVEN: Invalid service or local IDs in delete requests
		WHEN:  BulkDeletePostgresqlServiceLocals is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	invalidRequests := []*schemas.DeleteServiceLocalRequest{
		{
			ServiceId: uuid.Nil, // Invalid ID
			LocalId:   uuid.New(),
		},
	}

	// WHEN
	err := adapter.BulkDeletePostgresqlServiceLocals(invalidRequests)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, errors.UnprocessableEntityError.InvalidServiceLocalId, *err)
}
