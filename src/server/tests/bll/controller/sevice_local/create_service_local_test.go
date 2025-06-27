package sevice_local_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestCreateServiceLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid service and local IDs exist
		WHEN:  CreateServiceLocal is called
		THEN:  A new service-local association should be created
	*/
	// GIVEN
	serviceLocalController, _, db := controllerTest.NewServiceLocalControllerTestWrapper(t)

	// Create dependencies
	testService := factories.NewServiceModel(db, factories.ServiceModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})
	updatedBy := "ADMIN"

	// WHEN
	createRequest := schemas.CreateServiceLocalRequest{
		ServiceId: testService.Id,
		LocalId:   testLocal.Id,
	}
	result, err := serviceLocalController.CreateServiceLocal(
		createRequest,
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testService.Id, result.ServiceId)
	assert.Equal(t, testLocal.Id, result.LocalId)

	// Verify the association was created in the database
	var serviceLocal model.ServiceLocal
	dbErr := db.Where("service_id = ? AND local_id = ?", testService.Id, testLocal.Id).First(&serviceLocal).Error
	assert.NoError(t, dbErr)
	assert.Equal(t, testService.Id, serviceLocal.ServiceId)
	assert.Equal(t, testLocal.Id, serviceLocal.LocalId)
}

func TestCreateServiceLocalWithNonExistentService(t *testing.T) {
	/*
		GIVEN: A non-existent service ID
		WHEN:  CreateServiceLocal is called
		THEN:  An error should be returned
	*/
	// GIVEN
	serviceLocalController, _, db := controllerTest.NewServiceLocalControllerTestWrapper(t)

	// Create local but not service
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})
	nonExistentServiceId := uuid.New()
	updatedBy := "ADMIN"

	// WHEN
	result, err := serviceLocalController.CreateServiceLocal(
		schemas.CreateServiceLocalRequest{
			ServiceId: nonExistentServiceId,
			LocalId:   testLocal.Id,
		},
		updatedBy,
	)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Message, "not found")
}

func TestCreateServiceLocalWithNonExistentLocal(t *testing.T) {
	/*
		GIVEN: A non-existent local ID
		WHEN:  CreateServiceLocal is called
		THEN:  An error should be returned
	*/
	// GIVEN
	serviceLocalController, _, db := controllerTest.NewServiceLocalControllerTestWrapper(t)

	// Create service but not local
	testService := factories.NewServiceModel(db, factories.ServiceModelF{})
	nonExistentLocalId := uuid.New()
	updatedBy := "ADMIN"

	// WHEN
	result, err := serviceLocalController.CreateServiceLocal(
		schemas.CreateServiceLocalRequest{
			ServiceId: testService.Id,
			LocalId:   nonExistentLocalId,
		},
		updatedBy,
	)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Message, "not found")
}

func TestCreateServiceLocalEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: Valid service and local IDs but empty updatedBy
		WHEN:  CreateServiceLocal is called
		THEN:  An error should be returned
	*/
	// GIVEN
	serviceLocalController, _, db := controllerTest.NewServiceLocalControllerTestWrapper(t)

	// Create dependencies
	testService := factories.NewServiceModel(db, factories.ServiceModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// WHEN
	result, err := serviceLocalController.CreateServiceLocal(
		schemas.CreateServiceLocalRequest{
			ServiceId: testService.Id,
			LocalId:   testLocal.Id,
		},
		"", // Empty updatedBy
	)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}
