package service_professional_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestGetServiceProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: Service-professional association exists
		WHEN:  GetServiceProfessional is called with valid IDs
		THEN:  The service-professional association should be returned
	*/
	// GIVEN
	serviceProfessionalController, _, db := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	// Create test service and professional
	testService := factories.NewServiceModel(db)
	testProfessional := factories.NewProfessionalModel(db)

	// Create service-professional association
	serviceProfessional := &model.ServiceProfessional{
		ServiceId:      testService.Id,
		ProfessionalId: testProfessional.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "TEST_USER",
		},
	}
	err := db.Create(serviceProfessional).Error
	assert.NoError(t, err)

	// WHEN
	result, getErr := serviceProfessionalController.GetServiceProfessional(
		testService.Id.String(),
		testProfessional.Id.String(),
	)

	// THEN
	assert.Nil(t, getErr)
	assert.NotNil(t, result)
	assert.Equal(t, testService.Id, result.ServiceId)
	assert.Equal(t, testProfessional.Id, result.ProfessionalId)
}

func TestGetServiceProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: Service-professional association does not exist
		WHEN:  GetServiceProfessional is called with non-existent association
		THEN:  It should return service-professional not found error
	*/
	// GIVEN
	serviceProfessionalController, _, _ := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	nonExistentServiceId := uuid.New()
	nonExistentProfessionalId := uuid.New()

	// WHEN
	result, err := serviceProfessionalController.GetServiceProfessional(
		nonExistentServiceId.String(),
		nonExistentProfessionalId.String(),
	)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceProfessionalNotFound.Code, err.Code)
}

func TestGetServiceProfessionalWithInvalidServiceId(t *testing.T) {
	/*
		GIVEN: Invalid service ID format
		WHEN:  GetServiceProfessional is called with invalid service ID
		THEN:  It should return invalid service ID error
	*/
	// GIVEN
	serviceProfessionalController, _, _ := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	invalidServiceId := "invalid-uuid"
	validProfessionalId := uuid.New().String()

	// WHEN
	result, err := serviceProfessionalController.GetServiceProfessional(
		invalidServiceId,
		validProfessionalId,
	)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.UnprocessableEntityError.InvalidServiceId.Code, err.Code)
}

func TestGetServiceProfessionalWithInvalidProfessionalId(t *testing.T) {
	/*
		GIVEN: Invalid professional ID format
		WHEN:  GetServiceProfessional is called with invalid professional ID
		THEN:  It should return invalid professional ID error
	*/
	// GIVEN
	serviceProfessionalController, _, _ := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	validServiceId := uuid.New().String()
	invalidProfessionalId := "invalid-uuid"

	// WHEN
	result, err := serviceProfessionalController.GetServiceProfessional(
		validServiceId,
		invalidProfessionalId,
	)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.UnprocessableEntityError.InvalidProfessionalId.Code, err.Code)
}
