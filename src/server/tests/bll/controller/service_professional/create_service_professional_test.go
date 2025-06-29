package service_professional_test

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

func TestCreateServiceProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid service and professional exist
		WHEN:  CreateServiceProfessional is called with valid parameters
		THEN:  The service-professional association should be created successfully
	*/
	// GIVEN
	serviceProfessionalController, _, db := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	// Create test service
	testService := factories.NewServiceModel(db)
	// Create test professional
	testProfessional := factories.NewProfessionalModel(db)

	updatedBy := "TEST_USER"
	req := schemas.CreateServiceProfessionalRequest{
		ServiceId:      testService.Id,
		ProfessionalId: testProfessional.Id,
	}

	// WHEN
	result, err := serviceProfessionalController.CreateServiceProfessional(req, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testService.Id, result.ServiceId)
	assert.Equal(t, testProfessional.Id, result.ProfessionalId)

	// Verify in database
	var serviceProfessional model.ServiceProfessional
	dbErr := db.Where("service_id = ? AND professional_id = ?", testService.Id, testProfessional.Id).First(&serviceProfessional).Error
	assert.NoError(t, dbErr)
}

func TestCreateServiceProfessionalWithNonExistentService(t *testing.T) {
	/*
		GIVEN: Service does not exist
		WHEN:  CreateServiceProfessional is called with non-existent service ID
		THEN:  It should return service not found error
	*/
	// GIVEN
	serviceProfessionalController, _, db := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	// Create test professional
	testProfessional := factories.NewProfessionalModel(db)
	nonExistentServiceId := uuid.New()

	updatedBy := "TEST_USER"
	req := schemas.CreateServiceProfessionalRequest{
		ServiceId:      nonExistentServiceId,
		ProfessionalId: testProfessional.Id,
	}

	// WHEN
	result, err := serviceProfessionalController.CreateServiceProfessional(req, updatedBy)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceNotFound.Code, err.Code)
}

func TestCreateServiceProfessionalWithNonExistentProfessional(t *testing.T) {
	/*
		GIVEN: Professional does not exist
		WHEN:  CreateServiceProfessional is called with non-existent professional ID
		THEN:  It should return professional not found error
	*/
	// GIVEN
	serviceProfessionalController, _, db := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	// Create test service
	testService := factories.NewServiceModel(db)
	nonExistentProfessionalId := uuid.New()

	updatedBy := "TEST_USER"
	req := schemas.CreateServiceProfessionalRequest{
		ServiceId:      testService.Id,
		ProfessionalId: nonExistentProfessionalId,
	}

	// WHEN
	result, err := serviceProfessionalController.CreateServiceProfessional(req, updatedBy)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.ProfessionalNotFound.Code, err.Code)
}

func TestCreateServiceProfessionalAlreadyExists(t *testing.T) {
	/*
		GIVEN: Service-professional association already exists
		WHEN:  CreateServiceProfessional is called with existing association
		THEN:  It should return service-professional already exists error
	*/
	// GIVEN
	serviceProfessionalController, _, db := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	// Create test service and professional
	testService := factories.NewServiceModel(db)
	testProfessional := factories.NewProfessionalModel(db)

	// Create existing association
	existingServiceProfessional := &model.ServiceProfessional{
		ServiceId:      testService.Id,
		ProfessionalId: testProfessional.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "EXISTING_USER",
		},
	}
	err := db.Create(existingServiceProfessional).Error
	assert.NoError(t, err)

	updatedBy := "TEST_USER"
	req := schemas.CreateServiceProfessionalRequest{
		ServiceId:      testService.Id,
		ProfessionalId: testProfessional.Id,
	}

	// WHEN
	result, createErr := serviceProfessionalController.CreateServiceProfessional(req, updatedBy)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, createErr)
	assert.Equal(t, errors.ConflictError.ServiceProfessionalAlreadyExists.Code, createErr.Code)
}
