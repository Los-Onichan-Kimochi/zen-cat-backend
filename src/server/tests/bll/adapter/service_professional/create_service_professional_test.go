package service_professional_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestCreateServiceProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid service and professional IDs
		WHEN:  CreatePostgresqlServiceProfessional is called
		THEN:  A new service-professional association is created and returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	updatedBy := "test-admin"

	// WHEN
	serviceProfessional, err := adapter.CreatePostgresqlServiceProfessional(
		service.Id,
		professional.Id,
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceProfessional)
	assert.NotEmpty(t, serviceProfessional.Id)
	assert.Equal(t, service.Id, serviceProfessional.ServiceId)
	assert.Equal(t, professional.Id, serviceProfessional.ProfessionalId)
}

func TestCreateServiceProfessionalWithEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: Valid service and professional IDs but empty updatedBy
		WHEN:  CreatePostgresqlServiceProfessional is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	updatedBy := ""

	// WHEN
	serviceProfessional, err := adapter.CreatePostgresqlServiceProfessional(
		service.Id,
		professional.Id,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, serviceProfessional)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue, *err)
}

func TestGetServiceProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: A service-professional association exists in the database
		WHEN:  GetPostgresqlServiceProfessional is called
		THEN:  The association is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create the association first
	createdServiceProfessional, createErr := adapter.CreatePostgresqlServiceProfessional(
		service.Id,
		professional.Id,
		"test-admin",
	)
	assert.Nil(t, createErr)
	assert.NotNil(t, createdServiceProfessional)

	// WHEN
	serviceProfessional, err := adapter.GetPostgresqlServiceProfessional(service.Id, professional.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceProfessional)
	assert.Equal(t, service.Id, serviceProfessional.ServiceId)
	assert.Equal(t, professional.Id, serviceProfessional.ProfessionalId)
	assert.Equal(t, createdServiceProfessional.Id, serviceProfessional.Id)
}

func TestGetServiceProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: No service-professional association exists
		WHEN:  GetPostgresqlServiceProfessional is called with non-existent IDs
		THEN:  A not found error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	nonExistentServiceId := uuid.New()
	nonExistentProfessionalId := uuid.New()

	// WHEN
	serviceProfessional, err := adapter.GetPostgresqlServiceProfessional(nonExistentServiceId, nonExistentProfessionalId)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, serviceProfessional)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceProfessionalNotFound, *err)
}

func TestDeleteServiceProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: A service-professional association exists in the database
		WHEN:  DeletePostgresqlServiceProfessional is called
		THEN:  The association is deleted
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create the association first
	createdServiceProfessional, createErr := adapter.CreatePostgresqlServiceProfessional(
		service.Id,
		professional.Id,
		"test-admin",
	)
	assert.Nil(t, createErr)
	assert.NotNil(t, createdServiceProfessional)

	// WHEN
	err := adapter.DeletePostgresqlServiceProfessional(service.Id, professional.Id)

	// THEN
	assert.Nil(t, err)

	// Verify association is deleted by trying to get it
	deletedServiceProfessional, getErr := adapter.GetPostgresqlServiceProfessional(service.Id, professional.Id)
	assert.NotNil(t, getErr)
	assert.Nil(t, deletedServiceProfessional)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceProfessionalNotFound, *getErr)
}

func TestDeleteServiceProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: No service-professional association exists
		WHEN:  DeletePostgresqlServiceProfessional is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	nonExistentServiceId := uuid.New()
	nonExistentProfessionalId := uuid.New()

	// WHEN
	err := adapter.DeletePostgresqlServiceProfessional(nonExistentServiceId, nonExistentProfessionalId)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, errors.BadRequestError.ServiceProfessionalNotDeleted, *err)
}
