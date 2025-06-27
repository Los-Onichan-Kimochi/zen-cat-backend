package sevice_local_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestGetServiceLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: A service-local association exists in the database
		WHEN:  GetServiceLocal is called with valid ID
		THEN:  The service-local association should be returned
	*/
	// GIVEN
	serviceLocalController, _, db := controllerTest.NewServiceLocalControllerTestWrapper(t)

	// Create dependencies
	testService := factories.NewServiceModel(db, factories.ServiceModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create service-local association
	serviceLocal := &model.ServiceLocal{
		ServiceId: testService.Id,
		LocalId:   testLocal.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(serviceLocal).Error
	assert.NoError(t, err)

	// WHEN
	result, errResult := serviceLocalController.GetServiceLocal(testService.Id.String(), testLocal.Id.String())

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, serviceLocal.Id, result.Id)
	assert.Equal(t, testService.Id, result.ServiceId)
	assert.Equal(t, testLocal.Id, result.LocalId)
}

func TestGetServiceLocalNotFound(t *testing.T) {
	/*
		GIVEN: No service-local association exists with the given ID
		WHEN:  GetServiceLocal is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	serviceLocalController, _, _ := controllerTest.NewServiceLocalControllerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	result, err := serviceLocalController.GetServiceLocal(nonExistentId.String(), nonExistentId.String())

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Message, "not found")
}

func TestGetServiceLocalWithNilId(t *testing.T) {
	/*
		GIVEN: A nil UUID
		WHEN:  GetServiceLocal is called with nil UUID
		THEN:  An error should be returned
	*/
	// GIVEN
	serviceLocalController, _, _ := controllerTest.NewServiceLocalControllerTestWrapper(t)
	nilId := uuid.Nil

	// WHEN
	result, err := serviceLocalController.GetServiceLocal(nilId.String(), nilId.String())

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
