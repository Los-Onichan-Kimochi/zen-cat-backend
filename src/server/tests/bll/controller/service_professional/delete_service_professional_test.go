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

func TestDeleteServiceProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: Service-professional association exists
		WHEN:  DeleteServiceProfessional is called with valid IDs
		THEN:  The service-professional association should be deleted
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
	deleteErr := serviceProfessionalController.DeleteServiceProfessional(
		testService.Id.String(),
		testProfessional.Id.String(),
	)

	// THEN
	assert.Nil(t, deleteErr)

	// Verify deletion in database
	var deletedServiceProfessional model.ServiceProfessional
	dbErr := db.Where("service_id = ? AND professional_id = ?", testService.Id, testProfessional.Id).First(&deletedServiceProfessional).Error
	assert.Error(t, dbErr) // Should not be found
}

func TestDeleteServiceProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: Service-professional association does not exist
		WHEN:  DeleteServiceProfessional is called with non-existent association
		THEN:  It should return service-professional not found error
	*/
	// GIVEN
	serviceProfessionalController, _, _ := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	nonExistentServiceId := uuid.New()
	nonExistentProfessionalId := uuid.New()

	// WHEN
	err := serviceProfessionalController.DeleteServiceProfessional(
		nonExistentServiceId.String(),
		nonExistentProfessionalId.String(),
	)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceProfessionalNotFound.Code, err.Code)
}

func TestDeleteServiceProfessionalWithInvalidServiceId(t *testing.T) {
	/*
		GIVEN: Invalid service ID format
		WHEN:  DeleteServiceProfessional is called with invalid service ID
		THEN:  It should return invalid service ID error
	*/
	// GIVEN
	serviceProfessionalController, _, _ := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	invalidServiceId := "invalid-uuid"
	validProfessionalId := uuid.New().String()

	// WHEN
	err := serviceProfessionalController.DeleteServiceProfessional(
		invalidServiceId,
		validProfessionalId,
	)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, errors.UnprocessableEntityError.InvalidServiceId.Code, err.Code)
}
