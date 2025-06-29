package service_local_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestCreateServiceLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid service and local IDs
		WHEN:  CreatePostgresqlServiceLocal is called
		THEN:  A new service-local association is created and returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})
	updatedBy := "test-admin"

	// WHEN
	serviceLocal, err := adapter.CreatePostgresqlServiceLocal(
		service.Id,
		local.Id,
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceLocal)
	assert.NotEmpty(t, serviceLocal.Id)
	assert.Equal(t, service.Id, serviceLocal.ServiceId)
	assert.Equal(t, local.Id, serviceLocal.LocalId)
}

func TestCreateServiceLocalWithEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: Valid service and local IDs but empty updatedBy
		WHEN:  CreatePostgresqlServiceLocal is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})
	updatedBy := ""

	// WHEN
	serviceLocal, err := adapter.CreatePostgresqlServiceLocal(
		service.Id,
		local.Id,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, serviceLocal)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue, *err)
}

func TestGetServiceLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: A service-local association exists in the database
		WHEN:  GetPostgresqlServiceLocal is called
		THEN:  The association is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create the association first
	createdServiceLocal, createErr := adapter.CreatePostgresqlServiceLocal(
		service.Id,
		local.Id,
		"test-admin",
	)
	assert.Nil(t, createErr)
	assert.NotNil(t, createdServiceLocal)

	// WHEN
	serviceLocal, err := adapter.GetPostgresqlServiceLocal(service.Id, local.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceLocal)
	assert.Equal(t, service.Id, serviceLocal.ServiceId)
	assert.Equal(t, local.Id, serviceLocal.LocalId)
	assert.Equal(t, createdServiceLocal.Id, serviceLocal.Id)
}

func TestGetServiceLocalNotFound(t *testing.T) {
	/*
		GIVEN: No service-local association exists
		WHEN:  GetPostgresqlServiceLocal is called with non-existent IDs
		THEN:  A not found error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	nonExistentServiceId := uuid.New()
	nonExistentLocalId := uuid.New()

	// WHEN
	serviceLocal, err := adapter.GetPostgresqlServiceLocal(nonExistentServiceId, nonExistentLocalId)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, serviceLocal)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceLocalNotFound, *err)
}

func TestDeleteServiceLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: A service-local association exists in the database
		WHEN:  DeletePostgresqlServiceLocal is called
		THEN:  The association is deleted
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create the association first
	createdServiceLocal, createErr := adapter.CreatePostgresqlServiceLocal(
		service.Id,
		local.Id,
		"test-admin",
	)
	assert.Nil(t, createErr)
	assert.NotNil(t, createdServiceLocal)

	// WHEN
	err := adapter.DeletePostgresqlServiceLocal(service.Id, local.Id)

	// THEN
	assert.Nil(t, err)

	// Verify association is deleted by trying to get it
	deletedServiceLocal, getErr := adapter.GetPostgresqlServiceLocal(service.Id, local.Id)
	assert.NotNil(t, getErr)
	assert.Nil(t, deletedServiceLocal)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceLocalNotFound, *getErr)
}

func TestDeleteServiceLocalNotFound(t *testing.T) {
	/*
		GIVEN: No service-local association exists
		WHEN:  DeletePostgresqlServiceLocal is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewServiceLocalAdapterTestWrapper(t)

	nonExistentServiceId := uuid.New()
	nonExistentLocalId := uuid.New()

	// WHEN
	err := adapter.DeletePostgresqlServiceLocal(nonExistentServiceId, nonExistentLocalId)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, errors.BadRequestError.ServiceLocalNotDeleted, *err)
}
